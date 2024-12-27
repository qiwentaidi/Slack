import { Callgologger, PortBrute, NewWebScanner, LoadDirsearchDict, DirScan, HunterSearch, Subdomain, FofaSearch } from 'wailsjs/go/services/App'
import global from '@/global'
import async from 'async';
import { ElMessage, ElNotification } from 'element-plus';
import { dirsearch, structs } from 'wailsjs/go/models';
import { getProxy, getRootDomain } from '@/util';
import { crackDict } from '@/stores/options';


export async function LinkWebscan(urls: string[]) {
    let options: structs.WebscanOptions = {
        Target: urls,
        Thread: 50,
        Screenshot: false,
        Honeypot: true,
        DeepScan: true,
        RootPath: true,
        CallNuclei: true,
        TemplateFiles: [],
        SkipNucleiWithoutTags: false,
        GenerateLog4j2: true,
        AppendTemplateFolder: global.webscan.append_pocfile
    }
    await NewWebScanner(options, getProxy())
}

export function LinkCrack(ips: string[]) {
    let id = 0
    async.eachLimit(ips, 20, async (target: string, callback: () => void) => {
        let protocol = target.split("://")[0]
        let userDict = crackDict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
        if (crackDict.options.includes(protocol)) {
            Callgologger("info", target + " is start brute")
        }
        await PortBrute(target, userDict, crackDict.passwords)
        id++
        if (id == ips.length) {
            callback()
        }
    }, (err: any) => {
        Callgologger("info", "PortBrute Finished")
        ElNotification.success({
            message: "Crack Finished",
            position: "bottom-right"
        })
    });
}

export async function LinkDirsearch(url: string) {
    ElNotification.success({
        message: "已将目标联动至目录扫描",
        position: "bottom-right"
    })
    let dfp = global.PATH.homedir + "/slack/config/dirsearch/dicc.txt"
    let paths = await LoadDirsearchDict([dfp], "php,aspx,asp,jsp,html,js".split(','))
    global.temp.dirsearchPathConut = paths.length
    global.temp.dirsearchStartTime = Date.now()
    let option: dirsearch.Options = {
        Method: "GET",
        URLs: [url],
        Paths: paths,
        Workers: 25,
        Timeout: 8,
        BodyExclude: "",
        BodyLengthExcludeTimes: 5,
        StatusCodeExclude: [404],
        Redirect: false,
        Interval: 0,
        CustomHeader: "",
        Recursion: 0,
    }
    await DirScan(option)
}

export async function LinkHunter(query: string, count: string) {
    if (!global.space.hunterkey) {
        ElNotification.warning("请在设置处填写Hunter Key")
        return
    }
    ElMessage.info("正在查询鹰图数据，请稍后...")
    let result = await HunterSearch(global.space.hunterkey, query, count, "1", "0", "3", false)
    if (result.code !== 200) {
        if (result.code == 40205) {
            ElMessage(result.message)
        } else {
            ElMessage.error(result.message);
            return
        }
    }
    return result.data.arr.map(item => item.url);
}

export async function LinkFOFA(query: string, count: number) {
    if (global.space.fofakey == "" && global.space.fofaemail) {
        ElNotification.warning("请在设置处填写FOFA Key && FOFA Email")
        return
    }
    ElMessage.info("正在查询FOFA数据，请稍后...")
    let result = await FofaSearch(query, count.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, true)
    if (result.Error) {
        ElMessage.warning(result.Message)
        return
    }
    return result.Results.map(item => item.URL)
}

export async function LinkSubdomain(domains: string[]) {
    let rootDomains = <string[]>[]
    for (const domain of domains) {
        rootDomains.push(getRootDomain(domain))
    }
    Callgologger("info", `正在对${rootDomains.length}个域名进行子域名查询，请稍后...`)
    let option: structs.SubdomainOption = {
        Mode: 1,
        Domains: rootDomains,
        Subs: [],
        Thread: 10,
        Timeout: 5,
        DnsServers: ["223.6.6.6:53", "8.8.8.8:53"],
        ResolveExcludeTimes: 5,
        BevigilApi: global.space.bevigil,
        ChaosApi: global.space.chaos,
        SecuritytrailsApi: global.space.securitytrails,
        ZoomeyeApi: global.space.zoomeye,
        GithubApi: global.space.github,
        AppendEngines: ["FOFA", "Quake", "Hunter"],
        FofaAddress: global.space.fofaapi,
        FofaEmail: global.space.fofaemail,
        FofaApi: global.space.fofakey,
        HunterApi: global.space.hunterkey,
        QuakeApi: global.space.quakekey
    }
    await Subdomain(option)
}

