<script lang="ts" setup>
import { reactive, onMounted, ref, nextTick } from 'vue'
import { VideoPause, QuestionFilled, Plus, DocumentCopy, ChromeFilled, Filter, View, Clock, Delete, Share, DArrowRight, DArrowLeft, Picture, Reading, FolderOpened, Tickets, CloseBold, UploadFilled, Edit, Refresh } from '@element-plus/icons-vue';
import { InitRule, FingerprintList, NewWebScanner, GetFingerPocMap, ExitScanner, Callgologger, SpaceGetPort, HostAlive, NewTcpScanner, PortBrute } from 'wailsjs/go/services/App'
import { ElMessage, ElMessageBox } from 'element-plus';
import { TestProxy, Copy, generateRandomString, ProcessTextAreaInput, getProxy, ReadLine } from '@/util'
import global from "@/stores"
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import usePagination from '@/usePagination';
import { LinkFOFA, LinkHunter } from '@/linkage';
import { getBadgeClass, getSelectedIcon, getTagTypeBySeverity, highlightFingerprints, typeToIcon } from '@/stores/style';
import CustomTabs from '@/components/CustomTabs.vue';
import { CheckFileStat, DirectoryDialog, FileDialog, List, ReadFile, SaveFileDialog } from 'wailsjs/go/services/File';
import { isPrivateIP, validateIp, validateIpAndDomain } from '@/stores/validate';
import { structs } from 'wailsjs/go/models';
import { portGroupOptions, webReportOptions, webscanOptions, crackDict, WebsiteInputTips, HostInputTips } from '@/stores/options'
import { nanoid as nano } from 'nanoid'
import {
    RemovePocscanResult, RemoveScanTask, ExportWebReportWithHtml, ExportWebReportWithJson, RetrieveAllScanTasks, AddFingerscanResult,
    AddPocscanResult, AddScanTask, ReadWebReportWithJson, RenameScanTask, RetrieveFingerscanResults, RetrievePocscanResults, UpdateScanTaskWithResults,
    RemoveFingerprintResult, ExportWebReportWithExcel
} from 'wailsjs/go/services/Database';
import saveIcon from '@/assets/icon/save.svg'
import githubIcon from '@/assets/icon/github.svg'
import hostIcon from "@/assets/icon/host.svg";
import websiteIcon from "@/assets/icon/website.svg";
import { SaveConfig } from '@/config';
import throttle from 'lodash/throttle';
import async from 'async'
import { handleWebscanContextMenu } from '@/linkage/contextMenu';
import CustomTextarea from '@/components/CustomTextarea.vue';
import { IPParse, PortParse } from 'wailsjs/go/core/Tools';
import { ActivityItem } from '@/stores/interface';
import Loading from '@/components/Loading.vue';

// 增加较快的数据需要节流刷新
const throttleFingerscanUpdate = throttle(() => {
    fp.ctrl.watchResultChange(fp.table)
}, 1000);

// 热加载节流触发，防止反序列化过于频繁导致闪退
const throttleInitialize = throttle(() => {
    initialize()
}, 2000);

// 仪表盘内容
const dashboard = reactive({
    // 漏洞数量
    riskLevel: {
        CRITICAL: 0,
        HIGH: 0,
        MEDIUM: 0,
        LOW: 0,
        INFO: 0,
    },
    // 指纹数量
    fingerLength: 0,
    // 漏洞数量
    pocLength: 0,
    // 主动指纹需要发包总数
    activeCount: 0,
    // 主动指纹当前发包进度
    activePercentage: 0,
    // 漏洞扫描进度
    nucleiPercentage: 0,
    // 漏洞扫描目标总数
    nucleiCount: 0,
    // IP*端口扫描总数
    portscanCount: 0,
    // 端口扫描当前进度
    portscanPercentage: 0,
})

// 界面参数，仅作用于前端显示或者输入匹配，与后端无交互
const param = reactive({
    inputType: 0, // 输入类型
    portGroup: portGroupOptions[1].text,
    builtInUsername: true,
    builtInPassword: true,
    username: '',
    password: '',
    allTemplate: <{ label: string, value: string }[]>[],
    allFingerprint: <{ label: string, value: string }[]>[],
    hostFilter: '',
})

function updatePorts(index: number) {
    if (index >= 0 && index < portGroupOptions.length) {
        form.portlist = portGroupOptions[index].value;
    }
}

// activities 列表
const activities = ref<ActivityItem[]>([])

// timeline 滚动容器
const timelineContainer = ref<HTMLElement | null>(null)

function addActivity(newActivity: ActivityItem) {
    activities.value.push({
        content: newActivity.content,
        timestamp: new Date().toLocaleTimeString(),
        type: newActivity.type,
        icon: typeToIcon[newActivity.type]
    })

    nextTick(() => {
        if (timelineContainer.value) {
            timelineContainer.value.scrollTop = timelineContainer.value.scrollHeight
        }
    })
}

// webscan

const form = reactive({
    input: '',
    tags: [] as string[],
    newWebscanDrawer: false,
    newHostscanDrawer: false,
    runnningStatus: false,
    scanStopped: false,
    taskName: '',
    taskId: '',
    hideRequest: false,
    hideResponse: false,
    showYamlPoc: false,
    pocContent: '',
    portlist: '',
})

const spaceEngineConfig = reactive({
    fofaQuery: '',
    fofaDialog: false,
    fofaPageSize: 1000,
    hunterQuery: '',
    hunterDialog: false,
    hunterPageSize: "100",
})


const config = reactive({
    customTemplate: <string[]>[],
    customTags: <string[]>[],
    screenhost: false,
    rootPathScan: true,
    webscanOption: 0,
    skipNucleiWithoutTags: false,
    generateLog4j2: false,
    crack: false, // 是否开启暴破
    customHeaders: '',
    vulscan: false,
})

const detailDialog = ref(false)
const historyDialog = ref(false)

const selectedRow = ref();

let fp = usePagination<structs.InfoResult>(50)
let vp = usePagination<structs.VulnerabilityInfo>(50)
let rp = usePagination<structs.TaskResult>(20)

// 初始化时调用
async function initialize() {
    let isSuccess = await InitRule(global.webscan.append_pocfile);
    if (!isSuccess) {
        ElMessage.error({
            showClose: true,
            message: "初始化指纹规则失败，请检查配置文件",
        });
        return;
    }
    let FingperintDBList = await FingerprintList();
    dashboard.fingerLength = FingperintDBList?.length || 0;
    param.allFingerprint = Array.from(new Set(FingperintDBList)).map(item => ({ label: item, value: item }))

    // 获取POC数量
    let pocMap = await GetFingerPocMap();
    dashboard.pocLength = pocMap ? Object.keys(pocMap).length : 0;

    // 遍历模板
    let files = await List([global.PATH.homedir + "/slack/config/pocs", global.webscan.append_pocfile]);
    param.allTemplate = files
        .filter(file => file.Path.endsWith(".yaml"))
        .map(file => ({ label: file.BaseName, value: file.Path }));

    let result = await RetrieveAllScanTasks();
    if (Array.isArray(result)) {
        rp.table.result = result;
        rp.ctrl.watchResultChange(rp.table);
    }
}

// 初始化时调用
onMounted(() => {
    updatePorts(1);
    initialize()
    // 获得结果回调
    EventsOn("nucleiResult", (result: structs.VulnerabilityInfo) => {
        // 更新漏洞数量
        const riskLevelKey = result.Severity as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        AddPocscanResult(result)
        // 前段漏洞表格，需要当任务ID与结果的任务ID一致时，才更新漏洞表格
        if (form.taskId == result.TaskId) {
            vp.table.result.push(result)
            vp.ctrl.watchResultChange(vp.table)
        }
        taskManager.updateTaskTable(result.TaskId)
    });
    EventsOn("webFingerScan", (result: structs.InfoResult) => {
        if ((result.Scheme == "http" || result.Scheme == "https") && result.StatusCode == 0) {
            addActivity({
                content: result.URL + " 访问失败",
                type: "warning",
            })
            // fix in 2.1。0 端口扫描/网站扫描访问失败的目标应该也要输出结果
            // return
        }
        if ((result.Scheme == "http" || result.Scheme == "https") && result.StatusCode == 422) {
            addActivity({
                content: result.URL + " 为云防护地址, 已过滤",
                type: "warning",
            })
            return
        }
        AddFingerscanResult(result)
        if (form.taskId == result.TaskId) {
            fp.table.result.push(result)
            throttleFingerscanUpdate()
        }
        taskManager.updateTaskTable(result.TaskId)
    });
    EventsOn("ActiveCounts", (count: number) => {
        dashboard.activeCount = count
    });
    EventsOn("ActiveProgressID", (id: number) => {
        dashboard.activePercentage = Number(((id / dashboard.activeCount) * 100).toFixed(2));
    });
    EventsOn("NucleiCounts", (count: number) => {
        dashboard.nucleiCount = count
    });
    EventsOn("NucleiProgressID", (id: number) => {
        dashboard.nucleiPercentage = Number(((id / dashboard.nucleiCount) * 100).toFixed(2));
    });
    // 进度条
    EventsOn("progressID", (id: number) => {
        dashboard.portscanPercentage = Number(((id / dashboard.portscanCount) * 100).toFixed(2));
    });
    return () => {
        EventsOff("nucleiResult");
        EventsOff("webFingerScan");
        EventsOff("ActiveCounts");
        EventsOff("ActiveProgressID");
        EventsOff("NucleiCounts");
        EventsOff("NucleiProgressID");
        EventsOff("portScanLoading");
        EventsOff("progressID");
    };
});

async function startScan() {
    let engine = new Engine
    if (!await engine.checkOptions()) return
    param.inputType == 0 ? form.newWebscanDrawer = false : form.newHostscanDrawer = false
    if (param.inputType == 1 && (!await engine.checkHostscanOptions())) {
        return
    }
    // 检查是否能写入任务
    if (!(await taskManager.writeTask())) {
        form.runnningStatus = false
        return
    }
    engine.clearDashboard()
    form.runnningStatus = true
    await engine.Runner()
}

function stopScan() {
    if (!form.runnningStatus) return
    ElMessage.error("正在停止任务, 请稍后!")
    ExitScanner("[webscan]")
    ExitScanner("[portscan]")
    // 新增一个标志变量来确保setTimeout只执行一次
    if (!form.scanStopped) {
        form.scanStopped = true; // 设置标志为true，表示扫描已停止

        // 增加10s的暂停时间，先立即退出扫描，等10s之后再将扫描状态停止，以应对一些数据依旧在增加等问题
        setTimeout(() => {
            form.runnningStatus = false
            form.scanStopped = false; // 重置标志，以便下次可以再次调用stopScan
            addActivity({
                content: "用户已退出扫描任务",
                type: "danger",
            })
        }, 10000);
    }
}

// 自动化扫描引擎
class Engine {
    inputLines = [] as string[] // 输入的目标行
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
            return
        }

        // 检测代理连通性
        if (!await TestProxy()) {
            return
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
            return
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
        this.ips = await IPParse(this.conventionTarget)
        if (this.ips == null) {
            dashboard.portscanCount = this.specialTarget.length
        } else {
            dashboard.portscanCount = this.ips.length * this.portsList.length + this.specialTarget.length
        }
        if (dashboard.portscanCount == 0) {
            ElMessage.warning('可用目标或端口为空')
            return
        }
        // 判断是否进行存活探测
        if (global.webscan.default_alive_module != "None") {
            addActivity({
                content: "正在加载主机探活引擎, 更多信息请查看日志",
                type: "primary",
            })
            form.runnningStatus = true
            global.webscan.ping_check_alive = global.webscan.default_alive_module != 'ICMP'
            this.ips = await HostAlive(this.ips, global.webscan.ping_check_alive)
            if (this.ips == null) {
                addActivity({
                    content: "未发现存活主机, 任务已结束",
                    type: "warning",
                })
                form.runnningStatus = false
                return
            }
            addActivity({
                content: "主机存活目标数量: " + this.ips.length.toString(),
                type: "info",
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
            dashboard.riskLevel[key] = 0;
        });
        dashboard.nucleiPercentage = 0
        dashboard.activePercentage = 0
    }

    public async Runner() {
        // 先扫端口
        if (param.inputType == 1) {
            addActivity({
                content: "正在加载端口扫描引擎, 目标数量: " + dashboard.portscanCount.toString(),
                type: "primary",
            })
            await NewTcpScanner(form.taskId, this.specialTarget, this.ips, this.portsList, global.webscan.port_thread, global.webscan.port_timeout, getProxy())
            if (!config.vulscan) {
                addActivity({
                    content: "只需要进行端口扫描, 任务已结束",
                    type: "success",
                })
                form.runnningStatus = false
                return
            }
            this.inputLines = fp.table.result.filter(line => line.Scheme === "http" || line.Scheme === "https").map(item => item.URL);
            if (!form.runnningStatus || form.scanStopped) {
                return
            }
            // 对web服务先去提取再去除，之后扫描
            if (this.inputLines.length != 0) {
                addActivity({
                    content: "提取并移除需要深度探测网站指纹目标, 目标数量: " + this.inputLines.length.toString() + ", 正在进行网站指纹检测",
                    type: "primary",
                })
                fp.table.result = fp.table.result.filter(
                    (line) => line.Scheme !== "http" && line.Scheme !== "https"
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
        addActivity({
            content: "任务已结束",
            type: "success",
        })
        form.runnningStatus = false
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
        addActivity({
            content: "正在加载网站扫描引擎, 当前模式: " + webscanOptions.find(item => item.value == config.webscanOption).label + " 已加载目标数: " + this.inputLines.length,
            type: "primary",
        })
        if (global.proxy.enabled) {
            addActivity({
                content: "代理已启用, Nuclei改为单线程执行",
                type: "primary",
            })
            await NewWebScanner(form.taskId, options, getProxy(), false)
        } else {
            await NewWebScanner(form.taskId, options, getProxy(), true)
        }
    }

    public async CrackRunner() {
        if (!form.runnningStatus || form.scanStopped) {
            return
        }
        let crackLinks = fp.table.result.filter(
            (line) => crackDict.options.includes(line.Scheme.toLowerCase())
        ).map(item => item.URL);
        if (crackLinks.length == 0) {
            addActivity({
                content: "未发现可被暴破的目标",
                type: "warning",
            })
            return
        }
        let passDict = [] as string[]
        let userDict = [] as string[]
        if (param.builtInUsername) {
            for (var item of crackDict.usernames) {
                item.dic = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/username/" + item.name + ".txt"))!
            }
        }
        if (param.builtInPassword) {
            passDict = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/password/password.txt"))!
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
        addActivity({
            content: "正在加载暴破引擎, 当前目标数量:" + crackLinks.length.toString(),
            type: "primary",
        })
        // **使用 Promise 包装 async.eachLimit，确保 CrackRunner() 被 await**
        await new Promise((resolve, reject) => {
            async.eachLimit(crackLinks, global.webscan.crack_thread, (target: string, callback: (err?: any) => void) => {
                if (!form.runnningStatus || form.scanStopped) {
                    return callback();  // 结束当前任务
                }

                let protocol = target.split("://")[0];
                userDict = crackDict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!;

                Callgologger("info", target + " is start weak password cracking");
                PortBrute(form.taskId, target, userDict, passDict);

                callback();  // 任务完成后调用 callback
            },
                (err: any) => {
                    if (err) {
                        reject(err);
                    } else {
                        Callgologger("info", "Crack Finished");
                        resolve(null); // 任务全部完成
                    }
                }
            );
        });
    }
}

// 联动空间引擎

const uncover = {
    fofa: async function () {
        spaceEngineConfig.fofaDialog = false
        let urls = await LinkFOFA(spaceEngineConfig.fofaQuery, spaceEngineConfig.fofaPageSize)
        if (urls) form.input = urls.join("\n")
    },
    hunter: async function () {
        spaceEngineConfig.hunterDialog = false
        let urls = await LinkHunter(spaceEngineConfig.hunterQuery, spaceEngineConfig.hunterPageSize)
        if (urls) form.input = urls.join("\n")
    },
    shodan: async function () {
        if (!validateIp(shodanIp.value)) {
            ElMessage.warning("目标输入格式不正确!")
            return
        }
        const lines = ProcessTextAreaInput(shodanIp.value)
        for (const line of lines) {
            if (isPrivateIP(line)) {
                ElMessage.warning(line + " 为内网地址不支持扫描!")
                return
            }
        }

        let id = 0
        form.input = ""
        shodanRunningstatus.value = true
        let ips = await IPParse(lines)
        async.eachLimit(ips, shodanThread.value, async (ip: string, callback: () => void) => {
            if (!shodanRunningstatus.value) {
                return
            }
            let ports = await SpaceGetPort(ip)
            id++
            if (ports == null) {
                return
            }
            Callgologger("info", "[shodan] " + ip + " port: " + ports.join())
            shodanPercentage.value = Number(((id / ips.length) * 100).toFixed(2));
            for (const port of ports) {
                form.input += ip + ":" + port.toString() + "\n"
            }
        }, (err: any) => {
            shodanRunningstatus.value = false
            shodanVisible.value = false
        })
    }
}

async function selectFolder() {
    global.webscan.append_pocfile = await DirectoryDialog()
}

function pictrueSRC(filepath: string): string {
    if (filepath == '') return ''
    const filename = filepath.split(/[/\\]/).pop(); // 适配 Windows 和 Linux 路径
    return `http://127.0.0.1:8732/screenhost/${filename}`;
}

const reportOption = ref('HTML')
const reportName = ref('')
const exportDialog = ref(false)

// 任务管理
const taskManager = {
    generateUniqueTaskName: (baseName: string): string => {
        let name = baseName.trim()
        let counter = 1
        while (rp.table.result.some(item => item.TaskName === name)) {
            name = `${baseName}-${counter}`
            counter++
        }
        return name
    },
    writeTask: async function () {
        const uniqueName = taskManager.generateUniqueTaskName(form.taskName)
        form.taskId = nano()
        let isSuccess = await AddScanTask(form.taskId, uniqueName, form.input, 0, 0)
        if (!isSuccess) {
            ElMessage.error("添加任务失败")
            return false
        }
        rp.table.result.push({
            TaskId: form.taskId,
            TaskName: uniqueName,
            Targets: form.input,
            Failed: 0,
            Vulnerability: 0,
        })
        rp.ctrl.watchResultChange(rp.table)
        return true
    },
    viewTask: async function (row: any) {
        historyDialog.value = false;
        form.input = row.Targets;
        form.taskName = row.TaskName;
        form.taskId = row.TaskId
        const [fingerResult, nucleiResult] = await Promise.all([
            RetrieveFingerscanResults(row.TaskId),
            RetrievePocscanResults(row.TaskId)
        ]);

        if (fingerResult) {
            fp.table.result = fingerResult;
            fp.ctrl.watchResultChange(fp.table);
        }
        // 初始化风险等级计数
        Object.keys(dashboard.riskLevel).forEach(key => {
            dashboard.riskLevel[key] = 0;
        });
        if (!nucleiResult) {
            vp.table.result = []
            vp.ctrl.watchResultChange(vp.table);
            return
        }
        vp.table.result = nucleiResult;
        vp.ctrl.watchResultChange(vp.table);
        // 遍历结果，统计每个风险等级的数量
        vp.table.result.forEach(item => {
            const riskLevelKey = item.Severity as keyof typeof dashboard.riskLevel;
            if (dashboard.riskLevel[riskLevelKey] !== undefined) {
                dashboard.riskLevel[riskLevelKey]++;
            }
        });
    },
    deleteTask: function (alltaskids: string[]) {
        ElMessageBox.confirm(
            "确定删除所选中的任务记录",
            '警告',
            {
                type: 'warning',
            }
        )
            .then(async () => {
                for (const taskid of alltaskids) {
                    let isSuccess = await RemoveScanTask(taskid)
                    if (!isSuccess) {
                        ElMessage.error(`任务ID: ${taskid}, 删除失败`)
                        return
                    }
                    rp.table.result = rp.table.result.filter(item => item.TaskId != taskid)
                }
                ElMessage.success("删除成功")
                rp.ctrl.watchResultChange(rp.table)
            })
            .catch(() => {
                Callgologger("error", "[webscan] delete task failed")
            })
    },
    renameTask: async function (taskId: string) {
        ElMessageBox.prompt('重命名任务', '编辑', {})
            .then(async ({ value }) => {
                let isSuccess = await RenameScanTask(taskId, value)
                if (!isSuccess) {
                    ElMessage.error("修改失败")
                }
                const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
                if (taskIndex !== -1) {
                    rp.table.result[taskIndex] = { ...rp.table.result[taskIndex], TaskName: value };
                    rp.ctrl.watchResultChange(rp.table);
                }
                ElMessage.success("修改成功")
            })
    },
    importTask: async function () {
        const filepath = await FileDialog("*.json")
        if (!filepath) {
            return
        }
        const result = await ReadWebReportWithJson(filepath)
        if (result) {
            const id = nano()
            let isSuccess = await AddScanTask(id, id, result.Targets, 0, 0)
            if (!isSuccess) {
                ElMessage.error("添加任务失败")
                return false
            }
            result.Fingerprints.forEach(item => {
                AddFingerscanResult(item)
            })
            if (result.POCs) {
                result.POCs.forEach(item => {
                    AddPocscanResult(item)
                })
                rp.table.result.push({
                    TaskId: id,
                    TaskName: id,
                    Targets: result.Targets,
                    Failed: 0,
                    Vulnerability: result.POCs.length,
                })
            } else {
                rp.table.result.push({
                    TaskId: id,
                    TaskName: id,
                    Targets: result.Targets,
                    Failed: 0,
                    Vulnerability: 0,
                })
            }
            rp.ctrl.watchResultChange(rp.table)
            ElMessage.success("添加成功")
        }
    },
    showExportDialog: function () {
        reportName.value = ""
        if (rp.table.selectRows.length >= 2) {
            reportName.value = "合并报告-"
        }
        reportName.value += rp.table.selectRows[0].TaskName
        exportDialog.value = true
    },
    exportTask: async function () {
        let isSuccess = false
        let taskids = rp.table.selectRows.map(item => item.TaskId)
        let filepath = await SaveFileDialog(reportName.value)
        if (!filepath) {
            return
        }
        switch (reportOption.value) {
            case "EXCEL":
                isSuccess = await ExportWebReportWithExcel(filepath + ".xlsx", rp.table.selectRows)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
                break
            case "JSON":
                isSuccess = await ExportWebReportWithJson(filepath + ".json", rp.table.selectRows)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
                break
            default:
                isSuccess = await ExportWebReportWithHtml(filepath + ".html", taskids)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
        }
        exportDialog.value = false
    },
    updateTaskTable: function (taskId: string) {
        let vulcount = dashboard.riskLevel.CRITICAL + dashboard.riskLevel.HIGH + dashboard.riskLevel.MEDIUM + dashboard.riskLevel.LOW + dashboard.riskLevel.INFO
        UpdateScanTaskWithResults(taskId, 0, vulcount)
        // 更新对应taskId的表格列中的漏洞数量
        const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
        if (taskIndex !== -1) {
            rp.table.result[taskIndex].Vulnerability = vulcount;
            rp.ctrl.watchResultChange(rp.table);
        }
    },
}

async function deleteVuln(row: any) {
    let isSuccess = await RemovePocscanResult(form.taskId, row.ID, row.URL)
    if (!isSuccess) {
        ElMessage.error("删除失败")
        return
    }
    ElMessage.success({
        message: "删除成功",
        grouping: true,
    })
    const riskLevelKey = row.Severity as keyof typeof dashboard.riskLevel;
    if (dashboard.riskLevel[riskLevelKey] !== undefined) {
        dashboard.riskLevel[riskLevelKey]--;
    }
    vp.table.result = vp.table.result.filter(item => !(item.URL == row.URL && item.Severity == row.Severity && item.Name == row.Name));
    vp.ctrl.watchResultChange(vp.table);
    taskManager.updateTaskTable(form.taskId)
}


async function showPocDetail(filename: string) {
    form.showYamlPoc = !form.showYamlPoc
    if (!form.showYamlPoc) {
        return
    }
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + filename + ".yaml"
    let isStat = await CheckFileStat(filepath)
    if (!isStat) {
        filepath = global.webscan.append_pocfile + "/" + filename + ".yaml"
    }
    let file = await ReadFile(filepath)
    form.pocContent = file.Content
}


const shodanVisible = ref(false)
const shodanIp = ref('')
const shodanPercentage = ref(0)
const shodanThread = ref(2)
const shodanRunningstatus = ref(false)

</script>

<template>
    <el-card class="mb-10px">
        <div class="flex-between">
            <span class="font-bold ml-5px" style="font-size: 20px;">
                Overview
            </span>
            <div class="risk-level-display">
                <el-tooltip v-for="(value, level) in dashboard.riskLevel" :key="level" :content="level">
                    <div class="risk-badge" :class="getBadgeClass(level)">
                        {{ value }}
                    </div>
                </el-tooltip>
            </div>
            <div>
                <el-popover placement="left-start" trigger="click" :width="244" v-if="!form.runnningStatus">
                    <template #reference>
                        <el-button type="primary" :icon="Plus">新建任务</el-button>
                    </template>
                    <el-space>
                        <el-button :icon="websiteIcon"
                            @click="form.newWebscanDrawer = true; param.inputType = 0">网站扫描</el-button>
                        <el-button :icon="hostIcon"
                            @click="form.newHostscanDrawer = true; param.inputType = 1">主机扫描</el-button>
                    </el-space>
                </el-popover>
                <el-button type="danger" :icon="VideoPause" @click="stopScan" v-else>停止任务</el-button>
            </div>
        </div>
        <el-row :gutter="10" class="mt-10px">
            <el-col :span="12">
                <div class="flex">
                    <Loading style="margin-left: 40px;" />
                    <div class="progress-details">
                        <div class="progress-item">
                            <el-icon class="icon icon-blue">
                                <Search />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>端口扫描</span>
                                    <span>{{ dashboard.portscanPercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.portscanPercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>

                        <div class="progress-item">
                            <el-icon class="icon icon-green">
                                <Monitor />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>主动探测</span>
                                    <span>{{ dashboard.activePercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.activePercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>

                        <div class="progress-item">
                            <el-icon class="icon icon-red">
                                <Warning />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>漏洞检测</span>
                                    <span>{{ dashboard.nucleiPercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.nucleiPercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>
                    </div>
                </div>
                <div class="summary-box">
                    <span class="summary-value">{{ dashboard.fingerLength }}</span>
                    <span class="summary-label">指纹总数</span><br />
                    <span class="summary-value">{{ dashboard.pocLength }}</span>
                    <span class="summary-label">漏洞总数</span>
                    <el-tooltip content="热加载指纹和POC">
                        <el-button link @click="throttleInitialize()">
                            <template #icon>
                                <el-icon :size="20">
                                    <Refresh />
                                </el-icon>
                            </template>
                        </el-button>
                    </el-tooltip>
                </div>
            </el-col>
            <el-col :span="12">
                <div ref="timelineContainer" class="timelineContainer">
                    <el-timeline v-if="activities.length >= 1"
                        style="text-align: left; padding-left: 5px;">
                        <el-timeline-item v-for="(activity, index) in activities" :key="index" :icon="activity.icon"
                            :type="activity.type" :timestamp="activity.timestamp">
                            {{ activity.content }}
                        </el-timeline-item>
                    </el-timeline>
                    <span class="position-center" v-else>No active tasks</span>
                </div>
            </el-col>
        </el-row>
    </el-card>
    <CustomTabs>
        <el-tabs type="border-card">
            <el-tab-pane label="信息">
                <el-table :data="fp.table.pageContent" stripe height="100vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }" @row-contextmenu="handleWebscanContextMenu"
                    @sort-change="fp.ctrl.sortChange">
                    <el-table-column fixed prop="URL" label="Link" width="350px" />
                    <el-table-column width="150px" label="Port & Protocol" :show-overflow-tooltip="true">
                        <template #default="scope">
                            <el-tag type="primary" round effect="plain">{{ scope.row.Port }}</el-tag>
                            <el-tag type="primary" round effect="plain" class="ml-5px">{{ scope.row.Scheme
                                }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="StatusCode" width="100px" label="Code" sortable="custom" />
                    <el-table-column prop="Length" width="100px" label="Length" sortable="custom" />
                    <el-table-column prop="Title" label="Title" width="220px" />
                    <el-table-column prop="Fingerprints" label="Technologies" :min-width="300">
                        <template #default="scope">
                            <div class="finger-container">
                                <el-tag v-for="finger in scope.row.Fingerprints" :key="finger"
                                    :effect="scope.row.Detect === 'Default' ? 'light' : 'dark'"
                                    :type="highlightFingerprints(finger)">{{
                                        finger }}</el-tag>
                                <el-tag type="warning" v-if="scope.row.IsWAF">{{ scope.row.WAF }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="Screen" width="150">
                        <template #default="scope">
                            <el-image :src="pictrueSRC(scope.row.Screenshot)"
                                :preview-src-list="[pictrueSRC(scope.row.Screenshot)]" :initial-index="0"
                                preview-teleported :max-scale="1" v-if="scope.row.Screenshot != ''">
                                <template #error>
                                    <div class="image-slot">
                                        <el-icon :size="16">
                                            <Picture />
                                        </el-icon>
                                    </div>
                                </template>
                            </el-image>
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="flex-between mt-5px">
                    <div></div>
                    <el-pagination size="small" background @size-change="fp.ctrl.handleSizeChange"
                        @current-change="fp.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="fp.table.currentPage" :page-sizes="[10, 20, 50, 100, 200]"
                        :page-size="fp.table.pageSize" layout="total, sizes, prev, pager, next"
                        :total="fp.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
            <el-tab-pane label="漏洞">
                <el-table :data="vp.table.pageContent" stripe height="100vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }"
                    @sort-change="((data: any) => vp.ctrl.sortChange(data, true))"
                    @filter-change="vp.ctrl.filterChange">
                    <el-table-column prop="ID" label="Template" width="250px" />
                    <el-table-column prop="Type" label="Type" width="150px" />
                    <el-table-column prop="Severity" width="150px" column-key="Severity" label="Severity"
                        :filters="!form.runnningStatus ? vp.ctrl.getColumnFilters('Severity') : []" sortable="custom">
                        <template #filter-icon>
                            <Filter />
                        </template>
                        <!-- 需要通过设置动态类图来设置自定义的tag类型 -->
                        <template #default="scope">
                            <el-tag :type="getTagTypeBySeverity(scope.row.Severity)"
                                :class="{ 'el-tag--critical': scope.row.Severity === 'CRITICAL' }">
                                {{ scope.row.Severity }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="URL" label="URL" />
                    <el-table-column label="Operate" width="120px">
                        <template #default="scope">
                            <el-tooltip content="查看详情">
                                <el-button :icon="Reading" type="primary" link
                                    @click="detailDialog = true; selectedRow = scope.row" />
                            </el-tooltip>
                            <el-tooltip content="删除">
                                <el-button :icon="Delete" link @click="deleteVuln(scope.row)" />
                            </el-tooltip>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="flex-between mt-5px">
                    <div></div>
                    <el-pagination size="small" background @size-change="vp.ctrl.handleSizeChange"
                        @current-change="vp.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="vp.table.currentPage" :page-sizes="[10, 20, 50, 100, 200]"
                        :page-size="vp.table.pageSize" layout="total, sizes, prev, pager, next"
                        :total="vp.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-button :icon="Clock" @click="historyDialog = true">任务管理</el-button>
        </template>
    </CustomTabs>
    <el-drawer v-model="form.newHostscanDrawer" size="50%">
        <template #header>
            <span class="drawer-title">新建主机扫描</span>
            <el-button link @click="shodanVisible = true">
                <template #icon>
                    <img src="/shodan.png" style="width: 14px; height: 14px;">
                </template>
                Shodan
            </el-button>
        </template>
        <template #footer>
            <div class="position-center">
                <el-button type="primary" @click="startScan" style="bottom: 10px;">开始任务</el-button>
            </div>
        </template>
        <el-form label-width="auto">
            <el-form-item label="任务名称:">
                <el-input v-model="form.taskName" />
            </el-form-item>
            <el-form-item label="目标地址:">
                <CustomTextarea v-model="form.input" :rows="param.inputType == 0 ? 6 : 11"
                    :placeholder="param.inputType == 0 ? WebsiteInputTips : HostInputTips"></CustomTextarea>
            </el-form-item>
            <el-form-item label="端口:">
                <el-select v-model="param.portGroup" @change="updatePorts">
                    <el-option v-for="(item, index) in portGroupOptions" :label="item.text" :value="index" />
                </el-select>
                <el-input v-model="form.portlist" type="textarea" :rows="4" resize="none"
                    class="mt-5px"></el-input>
            </el-form-item>
            <el-form-item label="漏洞扫描:">
                <el-switch v-model="config.vulscan" class="w-full" />
                <span class="form-item-tips">开启后端口扫描结束会进一步提取WEB应用指纹, 并调用指纹漏洞扫描模式扫描漏洞</span>
            </el-form-item>
            <el-form-item label="高级配置:" v-show="config.vulscan">
                <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                    <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                </el-tooltip>
                <el-checkbox label="无指纹目标跳过漏扫" v-model="config.skipNucleiWithoutTags" />
                <el-checkbox label="网站截图" v-model="config.screenhost" />
            </el-form-item>
            <el-form-item label="口令暴破:" v-show="config.vulscan">
                <el-switch v-model="config.crack" class="w-full" />
                <span class="form-item-tips" v-show="config.crack">默认字典可通过 设置->
                    字典管理处修改, 由于RDP暴破可能存在闪退, 暂时不支持暴破</span>
            </el-form-item>
            <el-form-item label="用户字典:" v-show="config.crack">
                <CustomTextarea v-model="param.username" :rows="5"
                    @input="param.builtInUsername = param.username.length === 0"></CustomTextarea>
                <el-checkbox v-model="param.builtInUsername"
                    :disabled="param.username.length == 0">使用默认用户字典</el-checkbox>
            </el-form-item>
            <el-form-item label="密码字典:" v-show="config.crack">
                <CustomTextarea v-model="param.password" :rows="5"
                    @input="param.builtInPassword = param.password.length === 0"></CustomTextarea>
                <el-checkbox v-model="param.builtInPassword"
                    :disabled="param.password.length == 0">使用默认密码字典</el-checkbox>
            </el-form-item>
        </el-form>
    </el-drawer>
    <el-drawer v-model="form.newWebscanDrawer" size="50%">
        <template #header>
            <span class="drawer-title">新建网站扫描</span>
            <el-button link @click="spaceEngineConfig.fofaDialog = true">
                <template #icon>
                    <img src="/app/fofa.png" style="width: 16px; height: 16px;">
                </template>
                FOFA
            </el-button>
            <el-divider direction="vertical" />
            <el-button link @click="spaceEngineConfig.hunterDialog = true">
                <template #icon>
                    <img src="/app/hunter.ico" style="width: 16px; height: 16px;">
                </template>
                Hunter
            </el-button>
        </template>
        <template #footer>
            <div class="position-center">
                <el-button type="primary" @click="startScan" style="bottom: 10px; position: absolute;">开始任务</el-button>
            </div>
        </template>
        <el-form label-width="auto">
            <el-form-item label="任务名称:">
                <el-input v-model="form.taskName" />
            </el-form-item>
            <el-form-item label="目标地址:">
                <CustomTextarea v-model="form.input" :rows="6"
                    :placeholder="param.inputType == 0 ? WebsiteInputTips : HostInputTips"></CustomTextarea>
            </el-form-item>
            <el-form-item label="追加POC:">
                <el-input v-model="global.webscan.append_pocfile">
                    <template #suffix>
                        <el-tooltip content="选择POC文件夹">
                            <el-button :icon="FolderOpened" link @click="selectFolder()" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="保存">
                            <el-button :icon="saveIcon" link @click="SaveConfig" />
                        </el-tooltip>
                    </template>
                </el-input>
                <span class="form-item-tips">额外追加一个文件夹下的所有YAML POC文件，便于管理自添加POC</span>
            </el-form-item>
            <el-form-item>
                <template #label>模式:
                    <el-tooltip>
                        <template #content>
                            1、指纹扫描: 只进行简单的指纹探测，不会探测敏感目录<br />
                            2、全指纹扫描: 会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                            3、指纹漏洞扫描: 指纹+主动敏感目录探测，扫描完成后扫描指纹对应POC，如果网站未识别到指纹会扫描全漏洞<br />
                            4、专项扫描: 只扫描简单指纹和已选择的漏洞，如果不指定漏洞会扫指纹漏洞(忽略主动探测功能)
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-segmented v-model="config.webscanOption" :options="webscanOptions" class="w-full">
                    <template #default="{ item }">
                        <el-space :size="3">
                            <el-icon :size="18">
                                <component :is="item.icon" :key="item.value" />
                            </el-icon>
                            <div>{{ item.label }}</div>
                        </el-space>
                    </template>
                </el-segmented>
            </el-form-item>
            <el-form-item label="请求头:">
                <el-input v-model="config.customHeaders" :rows="3" type="textarea"
                    :placeholder="$t('tips.customHeaders')"></el-input>
            </el-form-item>
            <div v-if="config.webscanOption == 3">
                <el-form-item label="指定指纹:">
                    <el-select-v2 v-model="config.customTags" :options="param.allFingerprint" filterable multiple
                        clearable />
                    <span class="form-item-tips">类似Jeecg-Boot指纹漏洞都在API路径下, 可通过填写API地址并指定指纹来进行扫描(POC 需要进行适配)</span>
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 3">
                <el-form-item label="指定漏洞:">
                    <el-select-v2 v-model="config.customTemplate" :options="param.allTemplate" filterable multiple
                        clearable />
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 2">
                <el-form-item label="Log4j2:">
                    <el-switch v-model="config.generateLog4j2" class="w-full" />
                    <span class="form-item-tips">开启后会将所有目标添加Generate-Log4j2指纹</span>
                </el-form-item>
            </div>
            <el-form-item label="高级配置:">
                <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                    <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                </el-tooltip>
                <el-checkbox label="无指纹目标跳过漏扫" v-model="config.skipNucleiWithoutTags" />
                <el-checkbox label="网站截图" v-model="config.screenhost" />
            </el-form-item>
        </el-form>
    </el-drawer>
    <el-drawer v-model="detailDialog" size="80%" @close="form.showYamlPoc = false">
        <template #header>
            <el-button :icon="Reading" link>漏洞详情</el-button>
            <el-button :icon="githubIcon" link
                @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/issues/new?title=' + selectedRow.ID)">误报提交</el-button>
            <el-divider direction="vertical" />
        </template>
        <div v-if="selectedRow">
            <el-descriptions :column="1" border class="w-full mt-10px">
                <el-descriptions-item label="Name:">{{ selectedRow.Name }}</el-descriptions-item>
                <el-descriptions-item label="Description:">
                    <span class="all-break">{{ selectedRow.Description }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="Reference:">
                    <div v-for="item in selectedRow.Reference.split(',')">
                        <el-link>{{ item }}</el-link>
                    </div>
                </el-descriptions-item>
                <el-descriptions-item label="Extracted:">{{ selectedRow.Extract }}</el-descriptions-item>
            </el-descriptions>
            <div class="flex">
                <el-card v-show="!form.hideRequest" class="flex-1 mr-5px">
                    <div class="flex-between">
                        <span class="font-bold">Request</span>
                        <el-button-group>
                            <el-button :icon="form.hideResponse ? DArrowLeft : DArrowRight" link
                                @click="form.hideResponse = !form.hideResponse" />
                            <el-button :icon="DocumentCopy" link @click="Copy(selectedRow.Request)" />
                            <el-button :icon="ChromeFilled" link @click="BrowserOpenURL(selectedRow.URL)" />
                        </el-button-group>
                    </div>
                    <highlightjs language="http" :code="selectedRow.Request" class="font-small"></highlightjs>
                </el-card>
                <el-card v-show="!form.hideResponse" class="flex-1">
                    <div class="flex-between">
                        <span class="font-bold">Response
                            <el-tag type="success" style="margin-left: 2px;">{{ selectedRow.ResponseTime ?
                                selectedRow.ResponseTime : '0' }}s</el-tag>
                        </span>
                        <el-button-group>
                            <el-tooltip content="查看/关闭POC内容">
                                <el-button :icon="Tickets" link @click="showPocDetail(selectedRow.ID)" />
                            </el-tooltip>
                            <el-button :icon="form.hideRequest ? DArrowRight : DArrowLeft" link
                                @click="form.hideRequest = !form.hideRequest" />
                            <el-button :icon="DocumentCopy" link @click="Copy(selectedRow.Response)" />
                        </el-button-group>
                    </div>
                    <highlightjs language="http" :code="selectedRow.Response" class="font-small"></highlightjs>
                </el-card>
            </div>
            <el-card v-if="form.showYamlPoc" class="mt-10px">
                <div class="flex-between">
                    <span class="font-bold">Yaml Poc Content</span>
                    <el-button :icon="CloseBold" link @click="form.showYamlPoc = false" />
                </div>
                <highlightjs language="yaml" :code="form.pocContent" class="font-small"></highlightjs>
            </el-card>
        </div>
    </el-drawer>

    <el-drawer v-model="historyDialog" size="70%">
        <template #header>
            <el-text class="font-bold" style="font-size: 16px;"><el-icon :size="18" class="mr-5px">
                    <Clock />
                </el-icon><span>任务管理</span></el-text>
        </template>
        <el-table :data="rp.table.pageContent" stripe @selection-change="rp.ctrl.handleSelectChange"
            :cell-style="{ textAlign: 'center' }" :header-cell-style="{ 'text-align': 'center' }"
            style="height: calc(100vh - 115px)">
            <el-table-column type="selection" width="50px" />
            <el-table-column prop="TaskName" label="任务名称" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="Targets" label="目标" :show-overflow-tooltip="true">
                <template #default="scope">
                    <el-link @click="taskManager.viewTask(scope.row)">{{ scope.row.Targets.includes('\n') ?
                        scope.row.Targets.split('\n')[0] : scope.row.Targets
                        }}</el-link>
                </template>
            </el-table-column>
            <el-table-column label="资产" width="100px">
                <template #default="scope">
                    <span>{{ scope.row.Targets.includes('\n') ? scope.row.Targets.split('\n').length : 1 }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="Vulnerability" label="漏洞" width="100px" />
            <el-table-column label="操作" width="180px" align="center">
                <template #default="scope">
                    <el-button-group>
                        <el-tooltip content="查看">
                            <el-button :icon="View" link @click="taskManager.viewTask(scope.row)" />
                        </el-tooltip>
                        <el-tooltip content="重命名">
                            <el-button :icon="Edit" link @click="taskManager.renameTask(scope.row.TaskId)" />
                        </el-tooltip>
                        <el-tooltip content="删除">
                            <el-button :icon="Delete" link @click="taskManager.deleteTask([scope.row.TaskId])" />
                        </el-tooltip>
                    </el-button-group>
                </template>
            </el-table-column>
            <template #empty>
                <el-empty></el-empty>
            </template>
        </el-table>
        <div class="flex-between mt-5px">
            <el-space>
                <el-button :icon="UploadFilled" size="small" @click="taskManager.importTask()">导入任务</el-button>
                <el-button :icon="Share" size="small" @click="taskManager.showExportDialog"
                    :disabled="rp.table.selectRows.length < 1">导出报告</el-button>
                <el-button :icon="Delete" size="small"
                    @click="taskManager.deleteTask(rp.table.selectRows.map(item => item.TaskId))"
                    :disabled="rp.table.selectRows.length < 1">批量删除</el-button>
            </el-space>
            <el-pagination size="small" background @size-change="rp.ctrl.handleSizeChange"
                @current-change="rp.ctrl.handleCurrentChange" :pager-count="5" :current-page="rp.table.currentPage"
                :page-sizes="[20, 50, 100]" :page-size="rp.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="rp.table.result.length">
            </el-pagination>
        </div>
    </el-drawer>
    <el-dialog title="导出报告" v-model="exportDialog">
        <el-alert :title="'已选择' + rp.table.selectRows.length + '个任务'" type="info" show-icon :closable="false"
            class="mr-5px" />
        <el-form :model="form" label-width="auto">
            <el-form-item label="报告类型">
                <el-select v-model="reportOption" style="width: 240px;">
                    <!-- 自定义选中项的显示 -->
                    <template #label>
                        <el-space :size="6">
                            <el-icon :size="18">
                                <component :is="getSelectedIcon(reportOption)" />
                            </el-icon>
                            <span class="font-bold">{{ reportOption }}</span>
                        </el-space>
                    </template>

                    <!-- 选项列表 -->
                    <el-option v-for="item in webReportOptions" :key="item.label" :label="item.label"
                        :value="item.label">
                        <el-space :size="6">
                            <el-icon :size="18">
                                <component :is="item.icon" />
                            </el-icon>
                            <span>{{ item.label }}</span>
                        </el-space>
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="报告名称">
                <el-input v-model="reportName"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="exportDialog = false">取消</el-button>
                <el-button type="primary" @click="taskManager.exportTask()">
                    导出
                </el-button>
            </div>
        </template>
    </el-dialog>
    <!-- 联动模块 dialog -->
    <!-- 联动模块 dialog -->
    <!-- 联动模块 dialog -->
    <!-- 联动模块 dialog -->
    <el-dialog v-model="spaceEngineConfig.fofaDialog">
        <template #header>
            <span class="drawer-title"><img src="/app/fofa.png">导入FOFA目标</span>
        </template>
        <el-form :model="form" label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="spaceEngineConfig.fofaQuery"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-input-number v-model="spaceEngineConfig.fofaPageSize" :min="1" :max="10000" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.fofa">导入</el-button>
        </template>
    </el-dialog>

    <el-dialog v-model="spaceEngineConfig.hunterDialog">
        <template #header>
            <span class="drawer-title"><img src="/app/hunter.ico">导入鹰图目标，由于API查询限制，大数据推荐使用官网进行数据导出</span>
        </template>
        <el-form :model="form" label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="spaceEngineConfig.hunterQuery"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-select v-model="spaceEngineConfig.hunterPageSize" style="width: 150px;">
                    <el-option v-for='item in ["10", "20", "50", "100"]' :label="item" :value="item" />
                </el-select>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.hunter">导入</el-button>
        </template>
    </el-dialog>

    <el-dialog v-model="shodanVisible" width="500">
        <template #header>
            <span class="drawer-title"><img src="/shodan.png">从Shodan拉取资产端口开放情况</span>
        </template>
        <el-form :model="form" label-position="top">
            <el-form-item label="扫描线程:">
                <el-input-number v-model="shodanThread" :min="1" :max="3" />
            </el-form-item>
            <el-form-item>
                <el-input v-model="shodanIp" type="textarea" :rows="6" placeholder="输入规则与IP模式一致，但不支持域名和排除"></el-input>
            </el-form-item>
        </el-form>

        <template #footer>
            <div class="flex-between">
                <el-progress :text-inside="true" :stroke-width="18" :percentage="shodanPercentage"
                    style="width: 300px;" />
                <el-button size="small" type="danger" @click="shodanRunningstatus = false"
                    v-if="shodanRunningstatus">终止探测</el-button>
                <el-button size="small" type="primary" @click="uncover.shodan" v-else>
                    开始收集
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<style>
.timelineContainer {
    background-color: var(--timeline-bg-color);
    padding: 1.5rem;
    border-radius: 0.5rem;
    height: 172px;
    overflow-y: auto;
}

.summary-box {
    justify-content: center;
    margin-top: 10px;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    /* gap-2 = 8px */
    background-color: var(--timeline-bg-color);
    /* bg-gray-50 */
    padding: 0.5rem 1rem;
    /* py-2 px-4 = 8px 16px */
    border-radius: 0.5rem;
    /* rounded-lg = 8px */
}

.summary-value {
    font-size: 1.125rem;
    /* text-lg = 18px */
    font-weight: 500;
    /* font-medium */
}

.summary-label {
    font-size: 0.875rem;
    /* text-sm = 14px */
    color: #6B7280;
    /* text-gray-500 */
}

.progress-details {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-left: 48px;
    flex: 1;
}

.progress-item {
    display: flex;
    align-items: center;
    gap: 16px;
}

.icon {
    font-size: 18px;
}

.icon-blue {
    color: #3B82F6;
    /* Tailwind text-blue-500 */
}

.icon-green {
    color: #22C55E;
    /* Tailwind text-green-500 */
}

.icon-red {
    color: #EF4444;
    /* Tailwind text-red-500 */
}

.progress-content {
    flex: 1;
}

.progress-labels {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
}
</style>