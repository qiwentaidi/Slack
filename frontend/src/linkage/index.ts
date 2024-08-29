import { deduplicateUrlFingerMap } from '@/util'
import { Callgologger, PortBrute, FingerScan, ActiveFingerScan, NucleiScanner, NucleiEnabled, LoadDirsearchDict, DirScan, HunterSearch } from 'wailsjs/go/main/App'
import global from '@/global'
import async from 'async';
import { DirScanOptions, URLFingerMap } from '@/interface';
import { ElMessage, ElNotification } from 'element-plus';


export async function LinkWebscan(ips: string[]) {
    let id = 0
    global.temp.urlFingerMap = []
    Callgologger("info", `正在将WEB目标联动网站扫描，共计加载目标: ${ips.length}`)
    await FingerScan(ips, global.proxy)
    await ActiveFingerScan(ips, global.proxy)
    if (await NucleiEnabled(global.webscan.nucleiEngine)) {
        const filteredUrlFingerprints = global.temp.urlFingerMap
            .filter(item => item.finger.length > 0 && item.url)
            .map(item => ({ url: item.url, finger: item.finger }));
        async.eachLimit(deduplicateUrlFingerMap(filteredUrlFingerprints), 10, async (ufm: URLFingerMap, callback: () => void) => {
            if (ufm.finger.length == 0) {
                return
            }
            await NucleiScanner(0, ufm.url, ufm.finger, global.webscan.nucleiEngine, false, [], "")
            id++
            if (id == filteredUrlFingerprints.length) {
                callback()
            }
        }, (err: any) => {
            Callgologger("info", "Webscan Finished")
            ElNotification.success({
                message: "Webscan Finished",
                position: "bottom-right"
            })
        })
    } else {
        Callgologger("error", `Nuclei引擎无效，无法进行漏洞扫描，已结束！`)
    }
}

export function LinkCrack(ips: string[]) {
    let id = 0
    async.eachLimit(ips, 20, async (target: string, callback: () => void) => {
        let protocol = target.split("://")[0]
        let userDict = global.dict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
        if (global.dict.options.includes(protocol)) {
            Callgologger("info", target + " is start brute")
        }
        await PortBrute(target, userDict, global.dict.passwords)
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
    let option: DirScanOptions = {
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
    ElMessage.info("正在导入鹰图数据，请稍后...")
    let urls = <string[]>[]
    let result:any = await HunterSearch(global.space.hunterkey, query, count, "1", "0", "3", false)
    if (result.code !== 200) {
        if (result.code == 40205) {
            ElMessage(result.message)
        } else {
            ElMessage({
                message: result.message,
                type: "error",
            });
            return
        }
    }
    result.data.arr.forEach((item: any) => {
        urls.push(item.url)
    });
    return urls
}