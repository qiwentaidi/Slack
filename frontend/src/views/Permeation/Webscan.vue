<script lang="ts" setup>
import { reactive, onMounted, ref, nextTick } from 'vue'
import {
    VideoPause, QuestionFilled, Plus, DocumentCopy, ChromeFilled, CircleCheck, Notification,
    Filter, View, Clock, Delete, Share, DArrowRight, DArrowLeft, Warning,
    Reading, FolderOpened, Tickets, CloseBold, UploadFilled, Edit, Refresh, MoreFilled,
} from '@element-plus/icons-vue';
import { InitRule, FingerprintList, NewWebScanner, GetFingerPocMap, ExitScanner, ViewPictrue, Callgologger, SpaceGetPort, HostAlive, NewTcpScanner, PortBrute, GetAllFinger } from 'wailsjs/go/services/App'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus';
import { TestProxy, Copy, generateRandomString, ProcessTextAreaInput, getProxy, ReadLine } from '@/util'
import global from "@/stores"
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import usePagination from '@/usePagination';
import { LinkFOFA, LinkHunter } from '@/linkage';
import { getBadgeClass, getTagTypeBySeverity } from '@/stores/style';
import CustomTabs from '@/components/CustomTabs.vue';
import { CheckFileStat, DirectoryDialog, FileDialog, List, ReadFile, SaveFileDialog } from 'wailsjs/go/services/File';
import { isPrivateIP, validateIp, validateIpAndDomain } from '@/stores/validate';
import { structs } from 'wailsjs/go/models';
import { portGroupOptions, webReportOptions, webscanOptions, crackDict } from '@/stores/options'
import { nanoid as nano } from 'nanoid'
import {
    RemovePocscanResult, RemoveScanTask, ExportWebReportWithHtml,
    ExportWebReportWithJson, RetrieveAllScanTasks, AddFingerscanResult,
    AddPocscanResult, AddScanTask, ReadWebReportWithJson, RenameScanTask,
    RetrieveFingerscanResults, RetrievePocscanResults, UpdateScanTaskWithResults,
    RemoveFingerprintResult,
    ExportWebReportWithExcel
} from 'wailsjs/go/services/Database';
import saveIcon from '@/assets/icon/save.svg'
import githubIcon from '@/assets/icon/github.svg'
import dashboardIcon from '@/assets/icon/dashboard.svg'
import activityIcon from '@/assets/icon/activity.svg'
import bugIcon from "@/assets/icon/bug.svg";
import fingerprintIcon from "@/assets/icon/fingerprint.svg";
import hostIcon from "@/assets/icon/host.svg";
import websiteIcon from "@/assets/icon/website.svg";
import { SaveConfig } from '@/config';
import throttle from 'lodash/throttle';
import async from 'async'
import { handleWebscanContextMenu } from '@/linkage/contextMenu';
import CustomTextarea from '@/components/CustomTextarea.vue';
import { IPParse, PortParse } from 'wailsjs/go/core/Tools';

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

// 增加较快的数据需要节流刷新
const throttleFingerscanUpdate = throttle(() => {
    fp.ctrl.watchResultChange(fp.table)
}, 1000);

// 热加载节流触发，防止反序列化过于频繁导致闪退
const throttleInitialize = throttle(() => {
    initialize()
}, 2000);


const activities = ref<{ content: string, type: string, timestamp: string, icon: any }[]>([])

const timelineContainer = ref<HTMLElement | null>(null); // 创建一个 ref

function updateActivities(newActivity: { content: string, type: string, icon: any }) {
    activities.value.push({
        content: newActivity.content,
        timestamp: (new Date()).toTimeString(),
        type: newActivity.type,
        icon: newActivity.icon
    });
    // if (timelineContainer.value) {
    //     timelineContainer.value.scrollTop = timelineContainer.value.scrollHeight; // 滚动到最底部
    // }
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
    honeypot: true,
    rootPathScan: true,
    webscanOption: 0,
    skipNucleiWithoutTags: false,
    generateLog4j2: false,
    crack: false, // 是否开启暴破
    customHeaders: '',
    vulscan: false,
})

const detailDialog = ref(false)
const screenDialog = ref(false)
const historyDialog = ref(false)

const selectedRow = ref();

let fp = usePagination<structs.InfoResult>(50)
let vp = usePagination<structs.VulnerabilityInfo>(50)
let rp = usePagination<structs.TaskResult>(10)

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
    let list = await FingerprintList();
    dashboard.fingerLength = list?.length || 0; // 使用可选链

    // 获取POC数量
    let pocMap = await GetFingerPocMap();
    dashboard.pocLength = pocMap ? Object.keys(pocMap).length : 0; // 使用可选链
    param.allFingerprint = (await GetAllFinger()).map(item => ({ label: item, value: item }))

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
    EventsOn("nucleiResult", (result: any) => {
        const riskLevelKey = result.Severity as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        vp.table.result.push(result)
        vp.ctrl.watchResultChange(vp.table)
        AddPocscanResult(form.taskId, result)
        taskManager.updateTaskTable(form.taskId)
    });
    EventsOn("webFingerScan", (result: any) => {
        if ((result.Scheme == "http" || result.Scheme == "https") && result.StatusCode == 0) {
            updateActivities({
                content: result.URL + " access failed",
                type: "warning",
                icon: Warning,
            })
            return
        }
        if ((result.Scheme == "http" || result.Scheme == "https") && result.StatusCode == 422) {
            updateActivities({
                content: result.URL + " is cloud protected and has been filtered",
                type: "warning",
                icon: Warning,
            })
            return
        }
        fp.table.result.push(result)
        throttleFingerscanUpdate()
        AddFingerscanResult(form.taskId, result)
        taskManager.updateTaskTable(form.taskId)
    });
    EventsOn("ActiveCounts", (count: number) => {
        dashboard.activeCount = count
    });
    EventsOn("ActiveProgressID", (id: number) => {
        dashboard.activePercentage = Number(((id / dashboard.activeCount) * 100).toFixed(2));
    });
    // 扫描结束时让进度条为100
    EventsOn("ActiveScanComplete", () => {
        dashboard.activePercentage = 100
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
    // 扫描结束时让进度条为100
    EventsOn("scanComplete", () => {
        dashboard.portscanPercentage = 100
    });
    return () => {
        EventsOff("nucleiResult");
        EventsOff("webFingerScan");
        EventsOff("ActiveCounts");
        EventsOff("ActiveProgressID");
        EventsOff("ActiveScanComplete");
        EventsOff("NucleiCounts");
        EventsOff("NucleiProgressID");
        EventsOff("portScanLoading");
        EventsOff("scanComplete");
        EventsOff("progressID");
    };
});

async function startScan() {
    let engine = new Engine
    if (!await engine.checkOptions()) return
    param.inputType == 0 ? form.newWebscanDrawer = false : form.newHostscanDrawer = false
    activities.value = [] // 清空前面的任务进度
    engine.clearDashboard()
    switch (param.inputType) {
        // 网站扫描
        case 0:
            if (!engine.checkWebscanOptions()) return
            break
        // 主机扫描
        case 1:
            if (!await engine.checkHostscanOptions()) return
            break
        default:
            ElMessage.warning("未知扫描类型")
            return
    }
    // 检查是否能写入任务
    if (!(await taskManager.writeTask())) {
        form.runnningStatus = false
        return
    }
    form.runnningStatus = true
    await engine.Runner()
}

function stopScan() {
    if (!form.runnningStatus) return
    ElMessage.error("正在停止任务, 请稍后!")
    ExitScanner("[webscan]")
    ExitScanner("[portscan]")
    ExitScanner("[portbrute]")
    // 新增一个标志变量来确保setTimeout只执行一次
    if (!form.scanStopped) {
        form.scanStopped = true; // 设置标志为true，表示扫描已停止
        // 增加10s的暂停时间，先立即退出扫描，等10s之后再将扫描状态停止，以应对一些数据依旧在增加等问题
        setTimeout(() => {
            form.runnningStatus = false
            form.scanStopped = false; // 重置标志，以便下次可以再次调用stopScan
            updateActivities({
                content: "The user has canceled the scanning task",
                type: "danger",
                icon: CloseBold
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

        // 检查是否需要写入任务
        if (form.taskName == '') {
            ElMessage.warning({
                message: "任务名称不能为空",
                grouping: true,
            });
            return
        }

        // 检测代理连通性
        if (!await TestProxy()) {
            return
        }

        this.inputLines = ProcessTextAreaInput(form.input)
        return true
    }

    public checkWebscanOptions() {
        this.inputLines = this.inputLines.filter(item => item.startsWith("http://") || item.startsWith("https://"))
        if (this.inputLines.length == 0) {
            ElMessage.warning('网站扫描URL输入目标有误, 请检查!')
            return
        }
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
            updateActivities({
                content: "Loading host alive scan engine, for more info check console",
                type: "primary",
                icon: MoreFilled,
            })
            form.runnningStatus = true
            global.webscan.ping_check_alive = global.webscan.default_alive_module != 'ICMP'
            this.ips = await HostAlive(this.ips, global.webscan.ping_check_alive)
            if (this.ips == null) {
                updateActivities({
                    content: "No alive host found, task end",
                    type: "warning",
                    icon: Warning
                })
                form.runnningStatus = false
                return
            }
            updateActivities({
                content: "Host alive count: " + this.ips.length.toString(),
                type: "info",
                icon: Notification
            })
        }
        return true
    }

    // 清空仪表盘以及表格数据
    public clearDashboard() {
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
            updateActivities({
                content: "Loading portscan engine, count: " + dashboard.portscanCount.toString(),
                type: "primary",
                icon: MoreFilled,
            })
            await NewTcpScanner(this.specialTarget, this.ips, this.portsList, global.webscan.port_thread, global.webscan.port_timeout, getProxy())
            if (!config.vulscan) {
                updateActivities({
                    content: "Only portscan, task is completed",
                    type: "success",
                    icon: CircleCheck,
                })
                form.runnningStatus = false
                return
            }
            this.inputLines = fp.table.result.filter(
                (line) => line.Scheme === "http" || line.Scheme === "https"
            ).map(item => item.URL);
            if (!form.runnningStatus || form.scanStopped) {
                return
            }
            // 对web服务先去提取再去除，之后扫描
            if (this.inputLines.length != 0) {
                updateActivities({
                    content: "Extracting and removing webpage activity detection information, totaling: " + this.inputLines.length.toString() + ", followed by fingerprint detection",
                    type: "primary",
                    icon: MoreFilled
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
        this.WebRunner()
        this.CrackRunner()
        updateActivities({
            content: "All task completed",
            type: "success",
            icon: CircleCheck,
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
            Honeypot: config.honeypot,
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
        updateActivities({
            content: "Loading webscan engine, current mode: " + config.webscanOption.toString(),
            type: "primary",
            icon: MoreFilled,
        })
        await NewWebScanner(options, getProxy())
    }

    public async CrackRunner() {
        if (!form.runnningStatus || form.scanStopped) {
            return
        }
        if (!config.crack) {
            updateActivities({
                content: "Not enable crack",
                type: "warning",
                icon: Warning,
            })
            form.runnningStatus = false
            return
        }
        let crackLinks = fp.table.result.filter(
            (line) => crackDict.options.includes(line.Scheme.toLowerCase())
        ).map(item => item.URL);
        if (crackLinks.length == 0) {
            updateActivities({
                content: "No protocol that can be crack detected",
                type: "warning",
                icon: Warning,
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
        updateActivities({
            content: "Loading crack engine, target count:" + crackLinks.length.toString(),
            type: "primary",
            icon: MoreFilled,
        })
        async.eachLimit(crackLinks, global.webscan.crack_thread, async (target: string, callback: () => void) => {
            if (!form.runnningStatus || form.scanStopped) {
                return
            }
            let protocol = target.split("://")[0]
            userDict = crackDict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
            Callgologger("info", target + " is start weak password cracking")
            await PortBrute(target, userDict, passDict)
        }, (err: any) => {
            Callgologger("info", "Crack Finished")
        });
    }
}

// 联动空间引擎

const uncover = {
    fofa: function () {
        spaceEngineConfig.fofaDialog = false
        LinkFOFA(spaceEngineConfig.fofaQuery, spaceEngineConfig.fofaPageSize).then(urls => {
            if (urls) {
                form.input = urls.join("\n")
            }
        })
    },
    hunter: function () {
        spaceEngineConfig.hunterDialog = false
        LinkHunter(spaceEngineConfig.hunterQuery, spaceEngineConfig.hunterPageSize).then(urls => {
            if (urls) {
                form.input = urls.join("\n")
            }
        })
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
            if (ports.length > 0) {
                for (const port of ports) {
                    form.input += ip + ":" + port.toString() + "\n"
                }
            }
        }, (err: any) => {
            shodanRunningstatus.value = false
            shodanVisible.value = false
        })
    }
}

function highlightFingerprints(fingerprint: string) {
    if (fingerprint == "疑似蜜罐") {
        return "warning"
    }
    if (global.webscan.highlight_fingerprints.includes(fingerprint)) {
        return "danger"
    }
    return "primary"
}

async function selectFolder() {
    global.webscan.append_pocfile = await DirectoryDialog()
}

async function ShowWebPictrue(filepath: string) {
    let bs64 = await ViewPictrue(filepath)
    if (bs64 == '') return
    screenDialog.value = true
    await nextTick()
    document.getElementById('webscan-img')!.setAttribute('src', bs64)
}

const reportOption = ref('HTML')
const reportName = ref('')
const exportDialog = ref(false)

// 任务管理
const taskManager = {
    writeTask: async function () {
        form.taskId = nano()
        let isSuccess = await AddScanTask(form.taskId, form.taskName, form.input, 0, 0)
        if (!isSuccess) {
            ElMessage.error("添加任务失败")
            return false
        }
        rp.table.result.push({
            TaskId: form.taskId,
            TaskName: form.taskName,
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
        if (nucleiResult) {
            vp.table.result = nucleiResult;
            vp.ctrl.watchResultChange(vp.table);
            // 遍历结果，统计每个风险等级的数量
            vp.table.result.forEach(item => {
                const riskLevelKey = item.Severity as keyof typeof dashboard.riskLevel;
                if (dashboard.riskLevel[riskLevelKey] !== undefined) {
                    dashboard.riskLevel[riskLevelKey]++;
                }
            });
        } else {
            vp.table.result = []
            vp.ctrl.watchResultChange(vp.table);
        }
    },
    deleteTask: function (taskId: string) {
        ElMessageBox.confirm(
            '确定删除该任务记录?',
            '警告',
            {
                type: 'warning',
            }
        )
            .then(async () => {
                let isSuccess = await RemoveScanTask(taskId)
                if (!isSuccess) {
                    ElMessage.error("删除失败")
                    return
                }
                ElMessage.success("删除成功")
                rp.table.result = rp.table.result.filter(item => item.TaskId != taskId)
                rp.ctrl.watchResultChange(rp.table)
                fp.initTable()
                vp.initTable()
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
        const result = await ReadWebReportWithJson(filepath)
        if (result) {
            const id = nano()
            let isSuccess = await AddScanTask(id, id, result.Targets, 0, 0)
            if (!isSuccess) {
                ElMessage.error("添加任务失败")
                return false
            }
            result.Fingerprints.forEach(item => {
                AddFingerscanResult(id, item)
            })
            if (result.POCs) {
                result.POCs.forEach(item => {
                    AddPocscanResult(id, item)
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
    appendTaskToDatabase: async function () {
        if (vp.table.result.length == 0 && fp.table.result.length == 0) {
            ElMessage.warning('请先添加扫描任务')
            return
        }
        form.taskId = nano()
        form.taskName = generateRandomString(10)
        let isSuccess = await AddScanTask(form.taskId, form.taskName, form.input, 0, 0)
        if (!isSuccess) {
            ElMessage.error('添加任务失败')
            return
        }
        rp.table.result.push({
            TaskId: form.taskId,
            TaskName: form.taskName,
            Targets: form.input,
            Failed: 0,
            Vulnerability: 0,
        })
        rp.ctrl.watchResultChange(rp.table)
        for (const item of fp.table.result) {
            let isSuccess = await AddFingerscanResult(form.taskId, item)
            if (!isSuccess) {
                ElMessage.error('添加指纹识别结果失败')
                return
            }
        }
        for (const item of vp.table.result) {
            let isSuccess = await AddPocscanResult(form.taskId, item)
            if (!isSuccess) {
                ElMessage.error('添加漏洞扫描结果失败')
                return
            }
        }
        taskManager.updateTaskTable(form.taskId)
        ElNotification.success("添加成功")
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


function stopShodan() {
    shodanRunningstatus.value = false
    ElNotification.error({
        message: "用户已终止扫描任务",
        position: 'bottom-right',
    });
}


const hostInputTips = `192.168.1.1
192.168.1.1/24
192.168.1.1-192.168.255.255
192.168.1.1-255
www.example.com

扫描特定端口:
192.168.0.1:6379

排除IP或网段可以在可支持输入的IP格式前加!:
!192.168.1.6/28`

const webInputTips = `https://www.example.com`

// 获取选中的 `icon`
const getSelectedIcon = (selectedLabel: string) => {
    const selectedItem = webReportOptions.find(item => item.label === selectedLabel);
    return selectedItem ? selectedItem.icon : null;
};
</script>

<template>
    <el-row :gutter="10" style="margin-bottom: 10px;">
        <el-col :span="14">
            <el-card>
                <div class="my-header">
                    <el-text class="position-center" style="font-size: 20px; font-weight: bold;">
                        <el-icon style="margin-right: 5px;">
                            <dashboardIcon />
                        </el-icon>
                        Overview
                    </el-text>
                    <el-text class="position-center" style="font-size: 20px; font-weight: bold;">
                        <el-tooltip content="指纹数量">
                            <el-icon style="margin-right: 5px;">
                                <fingerprintIcon />
                            </el-icon>
                        </el-tooltip>
                        {{ dashboard.fingerLength }}
                        <el-tooltip content="漏洞数量">
                            <el-icon style="margin-right: 5px; margin-left: 10px;">
                                <bugIcon />
                            </el-icon>
                        </el-tooltip>
                        {{ dashboard.pocLength }}
                        <el-tooltip content="热加载指纹和POC">
                            <el-button link @click="throttleInitialize()" style="margin-left: 5px;">
                                <template #icon>
                                    <el-icon :size="20">
                                        <Refresh />
                                    </el-icon>
                                </template>
                            </el-button>
                        </el-tooltip>
                    </el-text>
                    <div class="risk-level-display">
                        <el-tooltip v-for="(value, level) in dashboard.riskLevel" :key="level" :content="level">
                            <div class="risk-badge" :class="getBadgeClass(level)">
                                {{ value }}
                            </div>
                        </el-tooltip>
                    </div>
                </div>
                <el-row :gutter="10" style="margin-top: 10px;">
                    <el-col :span="8">
                        <el-progress type="dashboard" :percentage="dashboard.portscanPercentage">
                            <template #default="{ percentage }">
                                <span class="percentage-value">{{ percentage }}%</span>
                                <span class="percentage-label">Port scan</span>
                            </template>
                        </el-progress>
                    </el-col>
                    <el-col :span="8">
                        <el-progress type="dashboard" :percentage="dashboard.activePercentage">
                            <template #default="{ percentage }">
                                <span class="percentage-value">{{ percentage }}%</span>
                                <span class="percentage-label">Active detection</span>
                            </template>
                        </el-progress>
                    </el-col>
                    <el-col :span="8">
                        <el-progress type="dashboard" :percentage="dashboard.nucleiPercentage">
                            <template #default="{ percentage }">
                                <span class="percentage-value">{{ percentage }}%</span>
                                <span class="percentage-label">Vulnerability</span>
                            </template>
                        </el-progress>
                    </el-col>
                </el-row>
            </el-card>
        </el-col>
        <el-col :span="10">
            <el-card style="height: 192px">
                <div class="my-header">
                    <el-text class="position-center" style="font-size: 20px; font-weight: bold;">
                        <el-icon style="margin-right: 5px;">
                            <activityIcon />
                        </el-icon>
                        Process
                    </el-text>
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
                <div ref="timelineContainer" style="max-height: 150px; overflow-y: auto;">
                    <el-timeline v-if="activities.length >= 1"
                        style="margin-top: 10px; text-align: left; padding-left: 5px;">
                        <el-timeline-item v-for="(activity, index) in activities" :key="index" :icon="activity.icon"
                            :type="activity.type" :timestamp="activity.timestamp">
                            {{ activity.content }}
                        </el-timeline-item>
                    </el-timeline>
                    <span class="position-center" v-else>No active tasks</span>
                </div>
            </el-card>
        </el-col>
    </el-row>
    <CustomTabs>
        <el-tabs type="border-card">
            <el-tab-pane label="信息">
                <el-table :data="fp.table.pageContent" stripe height="100vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }" @row-contextmenu="handleWebscanContextMenu"
                    @sort-change="fp.ctrl.sortChange">
                    <el-table-column fixed prop="URL" label="Link" width="300px" />
                    <el-table-column width="150px" label="Port/Protocol" :show-overflow-tooltip="true">
                        <template #default="scope">
                            <el-tag type="primary" round effect="plain">{{ scope.row.Port }}</el-tag>
                            <el-tag type="primary" round effect="plain" style="margin-left: 5px;">{{ scope.row.Scheme
                                }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="StatusCode" width="100px" label="Code" sortable="custom"
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="Length" width="100px" label="Length" sortable="custom"
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="Title" label="Title" width="220px" />
                    <el-table-column prop="Fingerprints" label="Fingerprint" :min-width="350">
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

                    <el-table-column label="Screen" width="80">
                        <template #default="scope">
                            <el-button :icon="View" link @click="ShowWebPictrue(scope.row.Screenshot)"
                                v-if="scope.row.Screenshot != ''" />
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
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
                <div class="my-header" style="margin-top: 5px;">
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
                <el-button type="primary" @click="startScan" style="bottom: 10px; position: absolute;">开始任务</el-button>
            </div>
        </template>
        <el-form label-width="auto">
            <el-form-item label="任务名称:">
                <el-input v-model="form.taskName" />
            </el-form-item>
            <el-form-item label="目标地址:">
                <CustomTextarea v-model="form.input" :rows="param.inputType == 0 ? 6 : 11"
                    :placeholder="param.inputType == 0 ? webInputTips : hostInputTips"></CustomTextarea>
            </el-form-item>
            <el-form-item label="端口:">
                <el-select v-model="param.portGroup" @change="updatePorts">
                    <el-option v-for="(item, index) in portGroupOptions" :label="item.text" :value="index" />
                </el-select>
                <el-input v-model="form.portlist" type="textarea" :rows="4" resize="none"
                    style="margin-top: 5px;"></el-input>
            </el-form-item>
            <el-form-item label="漏洞扫描:">
                <el-switch v-model="config.vulscan" style="width: 100%;" />
                <span class="form-item-tips">开启后端口扫描结束会进一步提取WEB应用指纹, 并调用指纹漏洞扫描模式扫描漏洞</span>
            </el-form-item>
            <el-form-item label="高级配置:" v-show="config.vulscan">
                <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                    <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                </el-tooltip>
                <el-checkbox label="无指纹目标跳过漏扫" v-model="config.skipNucleiWithoutTags" />
                <el-checkbox label="网站截图" v-model="config.screenhost" />
            </el-form-item>
            <el-form-item label="口令暴破:">
                <el-switch v-model="config.crack" style="width: 100%;" />
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
                    <img src="/app/fofa.ico">
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
                    :placeholder="param.inputType == 0 ? webInputTips : hostInputTips"></CustomTextarea>
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
                <el-segmented v-model="config.webscanOption" :options="webscanOptions" style="width: 100%;">
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
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 3">
                <el-form-item label="指定漏洞:">
                    <el-select-v2 v-model="config.customTemplate" :options="param.allTemplate" placeholder="可选择1-10个漏洞"
                        filterable multiple clearable :multiple-limit="10" />
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 2">
                <el-form-item label="Log4j2:">
                    <el-switch v-model="config.generateLog4j2" style="width: 100%;" />
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
            <el-descriptions :column="1" border style="margin-bottom: 10px; width: 100%;">
                <el-descriptions-item label="Name:">{{ selectedRow.Name }}</el-descriptions-item>
                <el-descriptions-item label="Description:">
                    <span class="all-break">{{ selectedRow.Description }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="Reference:">
                    <span v-for="item in selectedRow.Reference.split(',')" class="text-overflow-break">
                        {{ item }}
                    </span>
                </el-descriptions-item>
                <el-descriptions-item label="Extracted:">{{ selectedRow.Extract }}</el-descriptions-item>
            </el-descriptions>
            <div style="display: flex">
                <el-card v-show="!form.hideRequest" style="flex: 1; margin-right: 5px;">
                    <div class="my-header">
                        <span style="font-weight: bold;">Request</span>
                        <el-button-group>
                            <el-button :icon="form.hideResponse ? DArrowLeft : DArrowRight" link
                                @click="form.hideResponse = !form.hideResponse" />
                            <el-button :icon="DocumentCopy" link @click="Copy(selectedRow.Request)" />
                            <el-button :icon="ChromeFilled" link @click="BrowserOpenURL(selectedRow.URL)" />
                        </el-button-group>
                    </div>
                    <highlightjs language="http" :code="selectedRow.Request" style="font-size: small;"></highlightjs>
                </el-card>
                <el-card v-show="!form.hideResponse" style="flex: 1;">
                    <div class="my-header">
                        <span style="font-weight: bold;">Response
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
                    <highlightjs language="http" :code="selectedRow.Response" style="font-size: small;"></highlightjs>
                </el-card>
            </div>
            <el-card v-if="form.showYamlPoc" style="margin-top: 10px;">
                <div class="my-header">
                    <span style="font-weight: bold;">Yaml Poc Content</span>
                    <el-button :icon="CloseBold" link @click="form.showYamlPoc = false" />
                </div>
                <highlightjs language="yaml" :code="form.pocContent" style="font-size: small;"></highlightjs>
            </el-card>
        </div>
    </el-drawer>

    <el-dialog v-model="screenDialog" width="50%" title="网站截图">
        <img id="webscan-img" src="" style="width: 100%; height: 400px;">
    </el-dialog>

    <el-drawer v-model="historyDialog" size="70%">
        <template #header>
            <el-text style="font-weight: bold; font-size: 16px;"><el-icon :size="18" style="margin-right: 5px;">
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
                        <el-tooltip content="编辑">
                            <el-button :icon="Edit" link @click="taskManager.renameTask(scope.row.TaskId)" />
                        </el-tooltip>
                        <el-tooltip content="删除">
                            <el-button :icon="Delete" link @click="taskManager.deleteTask(scope.row.TaskId)" />
                        </el-tooltip>
                    </el-button-group>
                </template>
            </el-table-column>
            <template #empty>
                <el-empty></el-empty>
            </template>
        </el-table>
        <div class="my-header" style="margin-top: 5px;">
            <el-space>
                <el-tooltip content="存储当前结果数据">
                    <el-button :icon="saveIcon" size="small" @click="taskManager.appendTaskToDatabase">临时入库</el-button>
                </el-tooltip>
                <el-button :icon="UploadFilled" size="small" @click="taskManager.importTask()">导入任务</el-button>
                <el-button :icon="Share" size="small" @click="taskManager.showExportDialog"
                    :disabled="rp.table.selectRows.length < 1">导出报告</el-button>
            </el-space>
            <el-pagination size="small" background @size-change="rp.ctrl.handleSizeChange"
                @current-change="rp.ctrl.handleCurrentChange" :pager-count="5" :current-page="rp.table.currentPage"
                :page-sizes="[10, 20]" :page-size="rp.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="rp.table.result.length">
            </el-pagination>
        </div>
    </el-drawer>
    <el-dialog title="导出报告" v-model="exportDialog">
        <el-alert :title="'已选择' + rp.table.selectRows.length + '个任务'" type="info" show-icon :closable="false"
            v-show="rp.table.selectRows.length >= 1" style="margin-bottom: 5px" />
        <el-form :model="form" label-width="auto">
            <el-form-item label="报告类型">
                <el-select v-model="reportOption" style="width: 240px;">
                    <!-- 自定义选中项的显示 -->
                    <template #label>
                        <el-space :size="6">
                            <el-icon :size="18">
                                <component :is="getSelectedIcon(reportOption)" />
                            </el-icon>
                            <span style="font-weight: bold">{{ reportOption }}</span>
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
            <span class="drawer-title"><img src="/app/fofa.ico">导入FOFA目标</span>
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
            <div class="my-header">
                <el-progress :text-inside="true" :stroke-width="18" :percentage="shodanPercentage"
                    style="width: 300px;" />
                <div>
                    <el-button size="small" type="danger" @click="stopShodan"
                        v-if="shodanRunningstatus">终止探测</el-button>
                    <el-button size="small" type="primary" @click="uncover.shodan" v-else>
                        开始收集
                    </el-button>
                </div>
            </div>
        </template>
    </el-dialog>
</template>

<style scoped></style>