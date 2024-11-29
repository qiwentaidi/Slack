<script lang="ts" setup>
import { reactive, onMounted, ref, h, nextTick } from 'vue'
import { VideoPause, QuestionFilled, Plus, ZoomIn, DocumentCopy, ChromeFilled, Promotion, Filter, Upload, View, Clock, Delete, Share, DArrowRight, DArrowLeft, Reading, FolderOpened, Tickets, CloseBold, UploadFilled, Edit } from '@element-plus/icons-vue';
import { InitRule, FingerprintList, NewWebScanner, GetFingerPocMap, ExitScanner, ViewPictrue, Callgologger } from 'wailsjs/go/services/App'
import { ElMessage, ElMessageBox, ElNotification, TableColumnCtx } from 'element-plus';
import { TestProxy, Copy, CopyALL, transformArrayFields, FormatWebURL, UploadFileAndRead, generateRandomString } from '@/util'
import { ExportWebScanToXlsx } from '@/export'
import global from "@/global"
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import usePagination from '@/usePagination';
import { LinkDirsearch, LinkFOFA, LinkHunter } from '@/linkage';
import ContextMenu from '@imengyu/vue3-context-menu'
import { defaultIconSize, getTagTypeBySeverity } from '@/stores/style';
import CustomTabs from '@/components/CustomTabs.vue';
import { CheckFileStat, DirectoryDialog, FileDialog, List, ReadFile, SaveFileDialog } from 'wailsjs/go/services/File';
import { validateWebscan } from '@/stores/validate';
import { structs } from 'wailsjs/go/models';
import { webReportOptions, webscanOptions } from '@/stores/options'
import { nanoid as nano } from 'nanoid'
import { DeletePocscanResult, DeleteScanTask, ExportWebReportWithHtml, ExportWebReportWithJson, GetAllScanTask, InsertFingerscanResult, InsertPocscanResult, InsertScanTask, ReadWebReportWithJson, RenameScanTask, SelectFingerscanResult, SelectPocscanResult, UpdateScanWithResult } from 'wailsjs/go/services/Database';
import saveIcon from '@/assets/icon/save.svg'
import githubIcon from '@/assets/icon/github.svg'
import { SaveConfig } from '@/config';
import throttle from 'lodash/throttle';


const throttleUpdate = throttle(() => {
    fp.table.pageContent = fp.ctrl.watchResultChange(fp.table)
}, 1000);

// webscan
const customTemplate = ref<string[]>([])
const allTemplate = ref<{ label: string, value: string }[]>([])
const templateLoading = ref(false)

const dashboard = reactive({
    reqErrorURLs: [] as string[],
    riskLevel: {
        CRITICAL: 0,
        HIGH: 0,
        MEDIUM: 0,
        LOW: 0,
        INFO: 0,
    },
    currentModule: "",
    count: 0,
    fingerLength: 0,
    yamlPocsLength: 0,
})

const form = reactive({
    url: '',
    tags: [] as string[],
    newscanner: false,
    fofaDialog: false,
    fofaNum: 1000,
    hunterDialog: false,
    hunterNum: ["10", "20", "50", "100"],
    defaultNum: "100",
    query: '',
    runnningStatus: false,
    percentage: 0,
    nucleiPercentage: 0,
    count: 0,
    nucleiCounts: 0,
    taskName: '',
    taskId: '',
    hideRequest: false,
    hideResponse: false,
    showYamlPoc: false,
    pocContent: '',
    exportTask: <structs.TaskResult>{},
})

const config = reactive({
    screenhost: false,
    honeypot: true,
    rootPathScan: true,
    webscanOption: 0,
    skipNucleiWithoutTags: false,
    generateLog4j2: false,
    writeDB: false,
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
    dashboard.yamlPocsLength = pocMap ? Object.keys(pocMap).length : 0; // 使用可选链

    // 遍历模板
    let files = await List([global.PATH.homedir + "/slack/config/pocs", global.webscan.append_pocfile]);
    if (Array.isArray(files)) {
        allTemplate.value = files
            .filter(file => file.Path.endsWith(".yaml"))
            .map(file => ({ label: file.BaseName, value: file.Path }));
    }
    templateLoading.value = false; // 加载完毕

    let result = await GetAllScanTask();
    if (Array.isArray(result)) {
        rp.table.result = result;
        rp.table.pageContent = rp.ctrl.watchResultChange(rp.table);
    }
}

// 初始化时调用
onMounted(() => {
    initialize()
    // 获得结果回调
    EventsOn("nucleiResult", (result: any) => {
        const riskLevelKey = result.Severity as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        vp.table.result.push(result)
        vp.table.pageContent = vp.ctrl.watchResultChange(vp.table)
        if (!config.writeDB) return
        InsertPocscanResult(form.taskId, result)
        taskManager.updateTaskTable(form.taskId)
    });
    EventsOn("webFingerScan", (result: any) => {
        if (result.StatusCode == 0 || result.StatusCode == 422) {
            dashboard.reqErrorURLs.push(result.URL)
            return
        }
        fp.table.result.push(result)
        throttleUpdate()
        if (!config.writeDB) return
        InsertFingerscanResult(form.taskId, result)
        taskManager.updateTaskTable(form.taskId)
    });
    EventsOn("ActiveCounts", (count: number) => {
        form.count = count
    });
    EventsOn("ActiveProgressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    // 扫描结束时让进度条为100
    EventsOn("ActiveScanComplete", () => {
        form.percentage = 100
    });
    EventsOn("NucleiCounts", (count: number) => {
        form.nucleiCounts = count
    });
    EventsOn("NucleiProgressID", (id: number) => {
        form.nucleiPercentage = Number(((id / form.nucleiCounts) * 100).toFixed(2));
    });
    return () => {
        EventsOff("nucleiResult");
        EventsOff("webFingerScan");
        EventsOff("ActiveCounts");
        EventsOff("ActiveProgressID");
        EventsOff("ActiveScanComplete");
        EventsOff("NucleiCounts");
        EventsOff("NucleiProgressID");
    };
});

async function startScan() {
    let engine = new WebscanEngine
    if (!await engine.checkOptions()) return
    form.newscanner = false // 隐藏界面
    await engine.Runner()
}

function stopScan() {
    if (!form.runnningStatus) return
    form.runnningStatus = false
    ExitScanner("[webscan]")
    ElMessage.warning("任务已停止")
}

class WebscanEngine {
    urls = [] as string[]
    public async checkOptions() {
        if (!validateWebscan(form.url)) {
            return false
        }
        // 检查先行条件
        if (!await TestProxy(1)) {
            return false
        }
        if (config.writeDB && form.taskName == '') {
            ElMessage.warning({
                showClose: true,
                message: "任务名称不能为空",
            });
            return false
        }
        if (config.writeDB) {
            form.taskId = nano()
            let isSuccess = await InsertScanTask(form.taskId, form.taskName, form.url, 0, 0)
            if (!isSuccess) {
                ElMessage.error("添加任务失败")
                return false
            }
            rp.table.result.push({
                TaskId: form.taskId,
                TaskName: form.taskName,
                Targets: form.url,
                Failed: 0,
                Vulnerability: 0,
            })
            rp.table.pageContent = rp.ctrl.watchResultChange(rp.table)
        }
        fp.table.result = []
        vp.table.result = []
        vp.table.pageContent = vp.ctrl.watchResultChange(vp.table)
        Object.keys(dashboard.riskLevel).forEach(key => {
            dashboard.riskLevel[key] = 0;
        });
        dashboard.reqErrorURLs = []
        dashboard.currentModule = ""
        form.percentage = 0
        form.nucleiPercentage = 0
        form.runnningStatus = true
        return true
    }
    public async Runner() {
        this.urls = await FormatWebURL(form.url) // 处理目标
        dashboard.count = this.urls.length
        if (this.urls.length == 0) {
            ElMessage.warning({
                showClose: true,
                message: "可用目标为空",
            });
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
                customTemplate.value = []
                break;
            case 3:
                callNuclei = true
                break;
        }
        if (config.webscanOption != 2) config.generateLog4j2 = false

        let options: structs.WebscanOptions = {
            Target: this.urls,
            Thread: global.webscan.web_thread,
            Screenshot: config.screenhost,
            Honeypot: config.honeypot,
            DeepScan: deepScan,
            RootPath: config.rootPathScan,
            CallNuclei: callNuclei,
            TemplateFiles: customTemplate.value,
            SkipNucleiWithoutTags: config.skipNucleiWithoutTags,
            GenerateLog4j2: config.generateLog4j2,
            AppendTemplateFolder: global.webscan.append_pocfile,
        }
        await NewWebScanner(options, {
            Enabled: global.proxy.enabled,
            Mode: global.proxy.mode,
            Address: global.proxy.address,
            Port: global.proxy.port,
            Username: global.proxy.username,
            Password: global.proxy.password,
        })
        this.done()
        form.runnningStatus = false
    }

    public async done() {
        ElNotification.success({
            message: "Scanner Finished",
            position: "bottom-right",
        })
    }
}

// 联动空间引擎

const uncover = {
    fofa: function () {
        form.fofaDialog = false
        LinkFOFA(form.query, form.fofaNum).then(urls => {
            if (urls) {
                form.url = urls.join("\n")
            }
        })
    },
    hunter: function () {
        form.hunterDialog = false
        LinkHunter(form.query, form.defaultNum).then(urls => {
            if (urls) {
                form.url = urls.join("\n")
            }
        })
    },
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

function filterSeverity(value: string, row: structs.VulnerabilityInfo, column: TableColumnCtx<structs.VulnerabilityInfo>) :boolean {
    return row.Severity === value
}

function handleContextMenu(row: any, column: any, e: MouseEvent) {
    //prevent the browser's default menu
    e.preventDefault();
    //show our menu
    ContextMenu.showContextMenu({
        x: e.x,
        y: e.y,
        items: [
            {
                label: "复制链接",
                icon: h(DocumentCopy, defaultIconSize),
                onClick: () => {
                    Copy(row.URL)
                }
            },
            {
                label: "复制选中链接",
                divided: true,
                icon: h(DocumentCopy, defaultIconSize),
                onClick: () => {
                    CopyALL(fp.table.selectRows.map(item => item.URL))
                }
            },
            {
                label: "打开链接",
                icon: h(ChromeFilled, defaultIconSize),
                onClick: () => {
                    BrowserOpenURL(row.URL)
                }
            },
            {
                label: "打开选中链接",
                divided: true,
                icon: h(ChromeFilled, defaultIconSize),
                onClick: () => {
                    BatchBrowserOpenURL()
                }
            },
            {
                label: "联动目录扫描",
                icon: h(Promotion, defaultIconSize),
                onClick: () => {
                    LinkDirsearch(row.URL)
                }
            },
        ]
    });
}

function BatchBrowserOpenURL() {
    let targets = fp.table.selectRows.map(item => item.URL)
    for (const item of targets) {
        BrowserOpenURL(item)
    }
}

async function selectFile() {
    form.url = await UploadFileAndRead()
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
const exportDialog = ref(false)

// 任务管理
const taskManager = {
    viewTask: async function (row: any) {
        historyDialog.value = false;
        form.url = row.Targets;
        form.taskName = row.TaskName;
        form.taskId = row.TaskId
        config.writeDB = true
        dashboard.reqErrorURLs = []
        if (row.Targets != undefined) {
            row.Targets!.includes('\n') ? dashboard.count = row.Targets.split('\n').length : dashboard.count = 1
        }
        const [fingerResult, nucleiResult] = await Promise.all([
            SelectFingerscanResult(row.TaskId),
            SelectPocscanResult(row.TaskId)
        ]);

        if (fingerResult) {
            fp.table.result = fingerResult;
            fp.table.pageContent = fp.ctrl.watchResultChange(fp.table);
        }
        // 初始化风险等级计数
        Object.keys(dashboard.riskLevel).forEach(key => {
            dashboard.riskLevel[key] = 0;
        });
        if (nucleiResult) {
            vp.table.result = nucleiResult;
            vp.table.pageContent = vp.ctrl.watchResultChange(vp.table);
            // 遍历结果，统计每个风险等级的数量
            vp.table.result.forEach(item => {
                const riskLevelKey = item.Severity as keyof typeof dashboard.riskLevel;
                if (dashboard.riskLevel[riskLevelKey] !== undefined) {
                    dashboard.riskLevel[riskLevelKey]++;
                }
            });
        } else {
            vp.table.result = []
            vp.table.pageContent = vp.ctrl.watchResultChange(vp.table);
        }
    },
    deleteTask: function (taskId: string) {
        ElMessageBox.confirm(
            '确定删除该任务记录?',
            '警告',
            {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'warning',
            }
        )
            .then(async () => {
                let isSuccess = await DeleteScanTask(taskId)
                if (isSuccess) {
                    ElMessage.success("删除成功")
                    rp.table.result = rp.table.result.filter(item => item.TaskId != taskId)
                    rp.table.pageContent = rp.ctrl.watchResultChange(rp.table)
                    return
                }
                ElMessage.error("删除失败")
            })
            .catch(() => {
                Callgologger("error", "[webscan] delete task failed")
            })
    },
    renameTask: async function (taskId: string) {
        ElMessageBox.prompt('重命名任务', '编辑', {
            confirmButtonText: '确定',
        })
            .then(async ({ value }) => {
                let isSuccess = await RenameScanTask(taskId, value)
                if (isSuccess) {
                    const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
                    if (taskIndex !== -1) {
                        rp.table.result[taskIndex] = { ...rp.table.result[taskIndex], TaskName: value };
                        rp.table.pageContent = rp.ctrl.watchResultChange(rp.table);
                    }
                    ElMessage.success("修改成功")
                } else {
                    ElMessage.error("修改失败")
                }
            })
    },
    importTask: async function () {
        const filepath = await FileDialog("*.json")
        const result = await ReadWebReportWithJson(filepath)
        if (result) {
            const id = nano()
            let isSuccess = await InsertScanTask(id, id, result.Targets, 0, 0)
            if (!isSuccess) {
                ElMessage.error("添加任务失败")
                return false
            }
            rp.table.result.push({
                TaskId: id,
                TaskName: id,
                Targets: result.Targets,
                Failed: 0,
                Vulnerability: 0,
            })
            rp.table.pageContent = rp.ctrl.watchResultChange(rp.table)
            result.Fingerprints.forEach(item => {
                InsertFingerscanResult(id, item)
            })
            if (result.POCs) {
                result.POCs.forEach(item => {
                    InsertPocscanResult(id, item)
                })
            }
        }
    },
    exportTask: async function () {
        let filepath = ""
        let isSuccess = false
        var [fingerResult, nucleiResult] = await Promise.all([
            SelectFingerscanResult(form.exportTask.TaskId),
            SelectPocscanResult(form.exportTask.TaskId)
        ]);
        if (!nucleiResult) {
            nucleiResult = []
        }
        switch (reportOption.value) {
            case "EXCEL":
                ExportWebScanToXlsx(transformArrayFields(fingerResult), nucleiResult.map(({ ID, Name, Type, Severity, URL, Extract }) => ({
                    ID,
                    Name,
                    Type,
                    Severity,
                    URL,
                    Extract,
                })))
                break
            case "JSON":
                filepath = await SaveFileDialog(form.exportTask.TaskId)
                isSuccess = await ExportWebReportWithJson(filepath + ".json", form.exportTask)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
                break
            default:
                filepath = await SaveFileDialog(form.exportTask.TaskId)
                isSuccess = await ExportWebReportWithHtml(filepath + ".html", form.exportTask.TaskId)
                isSuccess ? ElMessage.success("导出成功") : ElMessage.error("导出失败")
        }
        exportDialog.value = false
    },
    updateTaskTable: function (taskId: string) {
        let vulcount = dashboard.riskLevel.CRITICAL + dashboard.riskLevel.HIGH + dashboard.riskLevel.MEDIUM + dashboard.riskLevel.LOW + dashboard.riskLevel.INFO
        UpdateScanWithResult(taskId, dashboard.reqErrorURLs.length, vulcount)
        // 更新对应taskId的表格列中的漏洞数量
        const taskIndex = rp.table.result.findIndex(task => task.TaskId === taskId);
        if (taskIndex !== -1) {
            rp.table.result[taskIndex].Vulnerability = vulcount;
            rp.table.pageContent = rp.ctrl.watchResultChange(rp.table);
        }
    },
    appendTaskToDatabase: async function () {
        if (vp.table.result.length == 0 && fp.table.result.length == 0) {
            ElMessage.warning('请先添加扫描任务')
            return
        }
        form.taskId = nano()
        form.taskName = generateRandomString(10)
        let isSuccess = await InsertScanTask(form.taskId, form.taskName, form.url, 0, 0)
        if (!isSuccess) {
            ElMessage.error('添加任务失败')
            return
        }
        rp.table.result.push({
            TaskId: form.taskId,
            TaskName: form.taskName,
            Targets: form.url,
            Failed: 0,
            Vulnerability: 0,
        })
        rp.table.pageContent = rp.ctrl.watchResultChange(rp.table)
        for (const item of fp.table.result) {
            let isSuccess = await InsertFingerscanResult(form.taskId, item)
            if (!isSuccess) {
                ElMessage.error('添加指纹识别结果失败')
                return
            }
        }
        for (const item of vp.table.result) {
            let isSuccess = await InsertPocscanResult(form.taskId, item)
            if (!isSuccess) {
                ElMessage.error('添加漏洞扫描结果失败')
                return
            }
        }
        taskManager.updateTaskTable(form.taskId)
        ElNotification.success("添加成功")
    }
}

async function deleteVuln(row: any) {
    let isSuccess = await DeletePocscanResult(form.taskId, row.ID, row.URL)
    if (isSuccess) {
        ElMessage.success("删除成功")
        const riskLevelKey = row.Severity as keyof typeof dashboard.riskLevel;
        if (dashboard.riskLevel[riskLevelKey] !== undefined) {
            dashboard.riskLevel[riskLevelKey]--;
        }
        vp.table.result = vp.table.result.filter(item => !(item.URL == row.URL && item.Severity == row.Severity && item.Name == row.Name));
        vp.table.pageContent = vp.ctrl.watchResultChange(vp.table);
        taskManager.updateTaskTable(form.taskId)

        return
    }
    ElMessage.error("删除失败")
}

function toggleRequest() {
    form.hideResponse = !form.hideResponse
}

function toggleResponse() {
    form.hideRequest = !form.hideRequest
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
</script>

<template>
    <el-card style="margin-bottom: 10px; height: 140px">
        <template #header>
            <div class="card-header">
                <span class="title">仪表盘</span>
                <el-space>
                    <el-tag>主动指纹</el-tag>
                    <el-progress :text-inside="true" :stroke-width="18" :percentage="form.percentage"
                        style="width: 300px;" />
                    <el-divider direction="vertical" />
                    <el-tag>漏洞扫描</el-tag>
                    <el-progress :text-inside="true" :stroke-width="18" :percentage="form.nucleiPercentage"
                        style="width: 300px;" />
                </el-space>
                <div>
                    <el-button type="primary" :icon="Plus" @click="form.newscanner = true"
                        v-if="!form.runnningStatus">新建任务</el-button>
                    <el-button type="danger" :icon="VideoPause" @click="stopScan" v-else>停止任务</el-button>
                </div>
            </div>
        </template>
        <el-row>
            <el-col :span="2" v-for="(total, risk) in dashboard.riskLevel">
                <el-statistic :title="risk.toLocaleUpperCase()" :value="total" />
            </el-col>
            <el-divider direction="vertical" style="height: 7vh;" />
            <el-col :span="3">
                <el-statistic title="指纹数量" :value="dashboard.fingerLength" />
            </el-col>
            <el-col :span="3">
                <el-statistic title="漏洞数量" :value="dashboard.yamlPocsLength" />
            </el-col>
            <el-divider direction="vertical" style="height: 7vh;" />
            <el-col :span="7">
                <el-statistic :value="dashboard.reqErrorURLs.length">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            失败数/目标数
                            <el-popover placement="left" :width="350" trigger="hover"
                                v-if="dashboard.reqErrorURLs.length >= 1">
                                <el-scrollbar height="150px">
                                    <p v-for="u in dashboard.reqErrorURLs" class="scrollbar-demo-item">
                                        {{ u }}</p>
                                </el-scrollbar>
                                <el-button :icon="DocumentCopy" @click="CopyALL(dashboard.reqErrorURLs)"
                                    style="width: 100%;">复制全部失败目标</el-button>
                                <template #reference>
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <ZoomIn />
                                    </el-icon>
                                </template>
                            </el-popover>
                        </div>
                    </template>
                    <template #suffix>/{{ dashboard.count.toString() }}</template>
                </el-statistic>
            </el-col>
        </el-row>
    </el-card>
    <CustomTabs>
        <el-tabs type="border-card">
            <el-tab-pane label="网站">
                <el-table :data="fp.table.pageContent" stripe height="100vh"
                    @selection-change="fp.ctrl.handleSelectChange" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }" @row-contextmenu="handleContextMenu">
                    <el-table-column type="selection" width="55px" />
                    <el-table-column fixed prop="URL" label="URL" width="300px" />
                    <el-table-column prop="StatusCode" width="100px" label="Code"
                        :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="Length" width="100px" label="Length"
                        :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="Title" label="Title" />
                    <el-table-column prop="Fingerprints" label="Fingerprint" width="350px">
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
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column prop="ID" label="Template" width="250px" />
                    <el-table-column prop="Type" label="Type" width="150px" />
                    <el-table-column prop="Severity" width="150px" label="Severity" :filter-method="filterSeverity"
                        :filters="[
                            { text: 'INFO', value: 'INFO' },
                            { text: 'LOW', value: 'LOW' },
                            { text: 'MEDIUM', value: 'MEDIUM' },
                            { text: 'HIGH', value: 'HIGH' },
                            { text: 'CRITICAL', value: 'CRITICAL' },
                        ]">
                        <template #filter-icon>
                            <Filter />
                        </template>
                        <template #default="scope">
                            <el-tag :type="getTagTypeBySeverity(scope.row.Severity)">{{ scope.row.Severity }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="URL" label="URL" />
                    <el-table-column label="Operate" width="120px">
                        <template #default="scope">
                            <el-tooltip content="查看详情">
                                <el-button :icon="Reading" type="primary" link
                                    @click="detailDialog = true; selectedRow = scope.row" />
                            </el-tooltip>
                            <el-tooltip content="打开漏洞链接">
                                <el-button :icon="ChromeFilled" link @click="BrowserOpenURL(scope.row.URL)" />
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
            <el-tooltip content="扫描任务历史">
                <el-button :icon="Clock" @click="historyDialog = true" />
            </el-tooltip>
        </template>
    </CustomTabs>

    <el-drawer v-model="form.newscanner" size="50%">
        <template #header>
            <span class="drawer-title">新建扫描任务</span>
            <el-button link @click="form.fofaDialog = true">
                <template #icon>
                    <img src="/app/fofa.ico">
                </template>
                FOFA
            </el-button>
            <el-divider direction="vertical" />
            <el-button link @click="form.hunterDialog = true">
                <template #icon>
                    <img src="/app/hunter.ico" style="width: 16px; height: 16px;">
                </template>
                Hunter
            </el-button>
        </template>
        <div class="form-container">
            <el-form label-width="auto">
                <el-form-item>
                    <template #label>目标:
                        <el-tooltip placement="left">
                            <template #content>
                                支持如下输入方式:<br />
                                192.168.0.1<br />
                                192.168.0.1:443<br />
                                https://192.168.0.1<br />
                                www.baidu.com
                            </template>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <el-input v-model="form.url" type="textarea" :rows="6" clearable></el-input>
                    <el-button link size="small" :icon="Upload" @click="selectFile()"
                        style="margin-top: 5px;">导入目标文件</el-button>
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
                    <el-segmented v-model="config.webscanOption" :options="webscanOptions" block style="width: 100%;">
                        <template #default="{ item }">
                            <div style="display: flex; align-items: center">
                                <el-icon :size="16" style="margin-right: 3px;">
                                    <component :is="item.icon" />
                                </el-icon>
                                <span>{{ item.label }}</span>
                            </div>
                        </template>
                    </el-segmented>
                </el-form-item>
                <el-form-item label="蜜罐识别:">
                    <el-switch v-model="config.honeypot" />
                </el-form-item>
                <div v-if="config.webscanOption == 3">
                    <el-form-item label="指定漏洞:">
                        <el-select-v2 v-model="customTemplate" :options="allTemplate" placeholder="可选择1-10个漏洞"
                            filterable multiple clearable :multiple-limit="10" v-if="!templateLoading" />
                        <span v-else>Loading templates...</span>
                    </el-form-item>
                </div>
                <div v-if="config.webscanOption == 2">
                    <el-form-item label="Log4j2:">
                        <el-switch v-model="config.generateLog4j2" style="margin-right: 5px;" />
                        <el-tag>开启后会将所有目标添加Generate-Log4j2指纹</el-tag>
                    </el-form-item>
                </div>
                <el-form-item label="结果入库:">
                    <el-switch v-model="config.writeDB" />
                </el-form-item>
                <el-form-item label="任务名称:" v-if="config.writeDB">
                    <el-input v-model="form.taskName" />
                </el-form-item>
                <el-form-item label="其他配置:">
                    <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                        <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                    </el-tooltip>
                    <el-tooltip content="关闭状态无指纹目标会扫全漏洞">
                        <el-checkbox label="跳过无指纹目标漏扫" v-model="config.skipNucleiWithoutTags" />
                    </el-tooltip>
                    <el-checkbox label="网站截图" v-model="config.screenhost" />
                </el-form-item>
            </el-form>
            <div class="position-center">
                <el-button type="primary" @click="startScan" style="bottom: 10px; position: absolute;">开始任务</el-button>
            </div>
        </div>
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
                                @click="toggleRequest" />
                            <el-button :icon="DocumentCopy" link @click="Copy(selectedRow.Request)" />
                        </el-button-group>
                    </div>
                    <highlightjs language="http" :code="selectedRow.Request" style="font-size: small;"></highlightjs>
                </el-card>
                <el-card v-show="!form.hideResponse" style="flex: 1;">
                    <div class="my-header">
                        <span style="font-weight: bold;">Response
                            <el-tag type="success" style="margin-left: 2px;">{{ selectedRow.ResponseTime ? selectedRow.ResponseTime : '0' }}s</el-tag>
                        </span>
                        <el-button-group>
                            <el-tooltip content="查看/关闭POC内容">
                                <el-button :icon="Tickets" link @click="showPocDetail(selectedRow.ID)" />
                            </el-tooltip>
                            <el-button :icon="form.hideRequest ? DArrowRight : DArrowLeft" link
                                @click="toggleResponse" />
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
    <el-dialog v-model="form.fofaDialog" title="导入FOFA目标(MAX 10000)" width="50%" center>
        <el-form :model="form" label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="form.query"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-input-number v-model="form.fofaNum" :min="1" :max="10000" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.fofa">导入</el-button>
        </template>
    </el-dialog>
    <el-dialog v-model="form.hunterDialog" title="导入鹰图目标(MAX100)，由于API查询大小限制，大数据推荐使用官网进行数据导出" width="50%" center>
        <el-form :model="form" label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="form.query"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-select v-model="form.defaultNum" style="width: 150px;">
                    <el-option v-for="item in form.hunterNum" :label="item" :value="item" />
                </el-select>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.hunter">导入</el-button>
        </template>
    </el-dialog>

    <el-dialog v-model="screenDialog" width="50%" title="网站截图">
        <img id="webscan-img" src="" style="width: 100%; height: 400px;">
    </el-dialog>

    <el-drawer v-model="historyDialog" size="70%">
        <template #header>
            <el-text style="font-weight: bold; font-size: 16px;"><el-icon :size="18" style="margin-right: 5px;">
                    <Clock />
                </el-icon><span>扫描任务历史</span></el-text>
        </template>
        <el-table :data="rp.table.pageContent" stripe :cell-style="{ textAlign: 'center' }"
            :header-cell-style="{ 'text-align': 'center' }" style="height: calc(100vh - 120px)">
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
                        <el-tooltip content="导出">
                            <el-button :icon="Share" link @click="exportDialog = true; form.exportTask = scope.row" />
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
                    <el-button :icon="saveIcon" @click="taskManager.appendTaskToDatabase">临时入库</el-button>
                </el-tooltip>
                <el-button :icon="UploadFilled" @click="taskManager.importTask()">导入任务</el-button>
            </el-space>
            <el-pagination size="small" background @size-change="rp.ctrl.handleSizeChange"
                @current-change="rp.ctrl.handleCurrentChange" :pager-count="5" :current-page="rp.table.currentPage"
                :page-sizes="[10, 20]" :page-size="rp.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="rp.table.result.length">
            </el-pagination>
        </div>
    </el-drawer>
    <el-dialog title="选择导出格式" v-model="exportDialog">
        <el-radio-group v-model="reportOption">
            <el-radio v-for="item in webReportOptions" :key="item" :label="item">{{ item }}</el-radio>
        </el-radio-group>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="exportDialog = false">取消</el-button>
                <el-button type="primary" @click="taskManager.exportTask()">
                    导出
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<style scoped>
.scrollbar-demo-item {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 50px;
    margin-block: 10px;
    text-align: center;
    border-radius: 4px;
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    overflow: hidden;
    text-overflow: ellipsis;
}

.title {
    font-size: 16px;
    font-weight: bold;
}

.form-container {
    height: calc(100% - 40px);
    overflow-y: auto;
    scrollbar-width: thin;
}
</style>