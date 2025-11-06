import { ElMessage } from 'element-plus';
import { NewWebScanner, NewTcpScanner, NewCrackScanenr, Callgologger, HostAlive } from 'wailsjs/go/services/App'
import { TestProxy, ProcessTextAreaInput, getProxy, generateRandomString, ReadLineWithoutNotify, ReadLine } from '@/util'
import { validateIpAndDomain } from '@/stores/validate';
import { IPParse, PortParse } from 'wailsjs/go/core/Tools';
import { FilepathJoin } from 'wailsjs/go/services/File';
import { RemoveFingerprintResult } from 'wailsjs/go/services/Database';
import { EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import global from "@/stores"
import { structs } from 'wailsjs/go/models';
import { crackDict, webscanOptions } from '@/stores/options'
import async from 'async'
import { dashboard, form, config, param, fp, vp, activities } from './useWebscanState'

export class WebscanEngine {
    inputLines = [] as string[] // 输入的目标行
    tcpLines = {} as {
        [key: string]: string[];
    }
    ips = [] as string[] // IP列表
    portsList = [] as number[] // 端口列表
    specialTarget = [] as string[] // IP:PORT 特殊目标
    conventionTarget = [] as string[] // 其他IP规则的目标

    // 检查基本条件
    public async checkOptions() {
        if (form.input == "") {
            ElMessage.warning({
                message: "目标不能为空",
                grouping: true
            })
            return false
        }

        // 检测代理连通性
        if (!await TestProxy()) {
            return false
        }

        if (form.taskName == '') {
            form.taskName = generateRandomString(8)
        }

        this.inputLines = ProcessTextAreaInput(form.input)
        return true
    }

    public async checkHostscanOptions() {
        const invalidLines = this.inputLines.filter(item => !validateIpAndDomain.some(validate => validate.test(item)));
        if (invalidLines.length > 0) {
            ElMessage.warning("主机扫描输入目标有误, 可通过控制台检查!")
            Callgologger("error", "以下行目标错误:\n" + invalidLines.join("\n"));
            return false
        }
        for (const line of this.inputLines) {
            if (line.includes(":")) {
                this.specialTarget.push(line)
            } else {
                this.conventionTarget.push(line)
            }
        }
        // 处理端口和IP组
        this.portsList = await PortParse(form.portlist)
        // 排除9100端口
        if (config.excludePrintPorts) {
            this.portsList = this.portsList.filter(port => port != 9100);
        }
        if (this.portsList.length == 0) {
            ElMessage.warning('端口列表为空')
            return false
        }
        this.ips = await IPParse(this.conventionTarget)
        // 检查是否有可用目标
        const estimatedCount = this.ips == null 
            ? this.specialTarget.length 
            : this.ips.length * this.portsList.length + this.specialTarget.length
        if (estimatedCount == 0) {
            ElMessage.warning('可用目标或端口为空')
            return false
        }
        // 判断是否进行存活探测
        if (global.webscan.default_alive_module != "None") {
            activities.value.push({
                content: "正在加载主机探活引擎, 更多信息请查看日志",
                timestamp: new Date().toLocaleTimeString(),
                type: "primary",
                icon: null
            })
            form.runnningStatus = true
            global.webscan.ping_check_alive = global.webscan.default_alive_module != 'ICMP'
            this.ips = await HostAlive(this.ips, global.webscan.ping_check_alive)
            if (this.ips == null) {
                activities.value.push({
                    content: "未发现存活主机, 任务已结束",
                    timestamp: new Date().toLocaleTimeString(),
                    type: "warning",
                    icon: null
                })
                form.runnningStatus = false
                return false
            }
            activities.value.push({
                content: "主机存活目标数量: " + this.ips.length.toString(),
                timestamp: new Date().toLocaleTimeString(),
                type: "info",
                icon: null
            })
        }
        return true
    }

    // 清空仪表盘以及表格数据
    public clearDashboard() {
        activities.value = [] // 清空前面的任务进度
        fp.initTable()
        vp.initTable()
        Object.keys(dashboard.riskLevel).forEach(key => {
            dashboard.riskLevel[key as keyof typeof dashboard.riskLevel] = 0;
        });
        dashboard.nucleiPercentage = 0
        dashboard.activePercentage = 0
    }

    public async Runner() {
        // 先扫端口
        if (param.inputType == 1) {
            // 计算估算值用于显示（实际总数由后端发送）
            const estimatedCount = this.ips == null 
                ? this.specialTarget.length 
                : this.ips.length * this.portsList.length + this.specialTarget.length
            activities.value.push({
                content: "正在加载端口扫描引擎, 目标数量: " + estimatedCount.toString(),
                timestamp: new Date().toLocaleTimeString(),
                type: "primary",
                icon: null
            })
            await NewTcpScanner(form.taskId, this.specialTarget, this.ips, this.portsList, global.webscan.port_thread, global.webscan.port_timeout, getProxy())
            if (!config.vulscan) {
                activities.value.push({
                    content: "只需要进行端口扫描, 任务已结束",
                    timestamp: new Date().toLocaleTimeString(),
                    type: "success",
                    icon: null
                })
                form.runnningStatus = false
                return
            }
            fp.table.result.map(line => {
                if (line.Scheme === "http" || line.Scheme === "https") {
                    this.inputLines.push(line.URL)
                } else {
                    this.tcpLines[line.URL] = line.Fingerprints
                }
            })
            if (!form.runnningStatus || form.scanStopped) {
                return
            }
            // 对web服务先去提取再去除，之后扫描
            if (this.inputLines.length != 0) {
                activities.value.push({
                    content: "提取并移除需要深度探测网站指纹目标, 目标数量: " + this.inputLines.length.toString() + ", 正在进行网站指纹检测",
                    timestamp: new Date().toLocaleTimeString(),
                    type: "primary",
                    icon: null
                })
                fp.table.result = fp.table.result.filter(
                    line => line.Scheme !== "http" && line.Scheme !== "https" && line.Scheme !== "apachemq"
                );
                fp.ctrl.watchResultChange(fp.table)
                RemoveFingerprintResult(form.taskId, this.inputLines)
            }
            // 设置后续的扫描类型
            config.webscanOption = 2
        }
        await this.WebRunner()
        if (config.crack) {
            await this.CrackRunner()
        }
        // fix 2.1.2 用户点击停止按钮过早结束，只有在任务未被手动停止时才提示结束
        if (!form.scanStopped) {
            activities.value.push({
                content: "任务已结束",
                timestamp: new Date().toLocaleTimeString(),
                type: "success",
                icon: null
            })
            form.runnningStatus = false
        }
    }

    public async WebRunner() {
        if (!form.runnningStatus || form.scanStopped) {
            return
        }
        // 指纹扫描      
        let deepScan = false
        let callNuclei = false
        switch (config.webscanOption) {
            case 1:
                deepScan = true
                break;
            case 2:
                callNuclei = true
                deepScan = true
                config.customTemplate = []
                break;
            case 3:
                callNuclei = true
                break;
        }
        if (config.webscanOption != 2) config.generateLog4j2 = false

        let options: structs.WebscanOptions = {
            Target: this.inputLines,
            TcpTarget: this.tcpLines,
            Thread: global.webscan.web_thread,
            Screenshot: config.screenhost,
            DeepScan: deepScan,
            RootPath: config.rootPathScan,
            CallNuclei: callNuclei,
            TemplateFiles: config.customTemplate,
            SkipNucleiWithoutTags: config.skipNucleiWithoutTags,
            GenerateLog4j2: config.generateLog4j2,
            AppendTemplateFolder: global.webscan.append_pocfile,
            NetworkCard: global.webscan.default_network,
            Tags: config.customTags,
            CustomHeaders: config.customHeaders,
        }
        activities.value.push({
            content: "正在加载网站扫描引擎, 当前模式: " + webscanOptions.find(item => item.value == config.webscanOption)!.label + " 已加载目标数: " + this.inputLines.length,
            timestamp: new Date().toLocaleTimeString(),
            type: "primary",
            icon: null
        })
        if (global.proxy.enabled) {
            activities.value.push({
                content: "代理已启用, Nuclei改为单线程执行",
                timestamp: new Date().toLocaleTimeString(),
                type: "primary",
                icon: null
            })
            await NewWebScanner(form.taskId, options, getProxy(), false)
        } else {
            await NewWebScanner(form.taskId, options, getProxy(), true)
        }
    }

    public async NewCrackScanenrAwait(taskId: string, target: string, userDict: string[], passDict: string[]) {
        return new Promise<void>((resolve) => {
            NewCrackScanenr(taskId, target, userDict, passDict)
            EventsOn(`crackDone::${target}`, () => {
                resolve()
            })
        });
    }

    public async CrackRunner() {
        if (!form.runnningStatus || form.scanStopped) {
            return
        }
        let crackLinks = fp.table.result.filter(line => crackDict.options.includes(line.Scheme.toLowerCase()))
            .map(item => item.URL);
        if (crackLinks.length == 0) {
            activities.value.push({
                content: "未发现可被暴破的目标",
                timestamp: new Date().toLocaleTimeString(),
                type: "warning",
                icon: null
            })
            return
        }
        let passDict = [] as string[]
        let userDict = [] as string[]
        if (param.builtInUsername) {
            for (var item of crackDict.usernames) {
                let filepath = await FilepathJoin([global.PATH.homedir, global.PATH.PortBurstPath, item.dicPath])
                item.dic = (await ReadLineWithoutNotify(filepath))!
            }
        }
        if (param.builtInPassword) {
            let filepath = await FilepathJoin([global.PATH.homedir, global.PATH.PortBurstPath, "/password/password.txt"])
            passDict = (await ReadLine(filepath))!
            passDict.push("")
        }
        if (param.username != "") {
            let users = ProcessTextAreaInput(param.username)
            for (var item of crackDict.usernames) {
                item.dic.push(...users)
            }
        }
        if (param.password != "") {
            let pass = ProcessTextAreaInput(param.password)
            passDict.push(...pass)
        }
        activities.value.push({
            content: "正在加载暴破引擎, 当前目标数量:" + crackLinks.length.toString(),
            timestamp: new Date().toLocaleTimeString(),
            type: "primary",
            icon: null
        })
        // **使用 Promise 包装 async.eachLimit，确保 CrackRunner() 被 await**
        await new Promise((resolve, reject) => {
            async.eachLimit(
                crackLinks,
                global.webscan.crack_thread,
                (target: string, callback: (err?: any) => void) => {
                    (async () => {
                        try {
                            if (!form.runnningStatus || form.scanStopped) return callback();

                            let protocol = target.split("://")[0];
                            userDict = crackDict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!;

                            Callgologger("info", target + " is start weak password cracking");
                            await this.NewCrackScanenrAwait(form.taskId, target, userDict, passDict);
                        } catch (err) {
                            Callgologger("error", "crack error: " + err);
                        } finally {
                            EventsOff(`crackDone::${target}`);
                            callback();
                        }
                    })();
                },
                (err: any) => {
                    if (err) {
                        reject(err);
                    } else {
                        resolve(null);
                    }
                }
            );
        });
    }
}

