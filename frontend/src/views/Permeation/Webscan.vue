<script lang="ts" setup>
import { reactive, onMounted, ref, h, nextTick } from 'vue'
import { VideoPause, QuestionFilled, Plus, ZoomIn, CopyDocument, ChromeFilled, Promotion, Filter, Upload, View, Clock } from '@element-plus/icons-vue';
import { InitRule, FingerprintList, NewWebScanner, GetFingerPocMap, ExitScanner, ViewPictrue } from 'wailsjs/go/main/App'
import { ElMessage, ElNotification } from 'element-plus';
import { TestProxy, Copy, CopyALL, transformArrayFields, FormatWebURL, getProxy, UploadFileAndRead } from '@/util'
import { ExportWebScanToXlsx } from '@/export'
import global from "@/global"
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import { Vulnerability, FingerprintTable } from '@/stores/interface';
import usePagination from '@/usePagination';
import exportIcon from '@/assets/icon/doucment-export.svg'
import { LinkDirsearch, LinkFOFA, LinkHunter } from '@/linkage';
import ContextMenu from '@imengyu/vue3-context-menu'
import { defaultIconSize, getClassBySeverity } from '@/stores/style';
import CustomTabs from '@/components/CustomTabs.vue';
import { List } from 'wailsjs/go/main/File';
import { validateWebscan } from '@/stores/validate';
import { structs } from 'wailsjs/go/models';
import { webscanOptions } from '@/stores/options'
import { nanoid as nano } from 'nanoid'
import { InsertFingerscanResult, InsertScanTask } from 'wailsjs/go/main/Database';

// webscan
const customTemplate = ref<string[]>([])
const allTemplate = ref<{ label: string, value: string }[]>([])
const templateLoading = ref(false)

const dashboard = reactive({
    reqErrorURLs: [] as string[],
    riskLevel: {
        critical: 0,
        high: 0,
        medium: 0,
        low: 0,
        info: 0,
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
    vulResult: [] as Vulnerability[],
    fingerResult: [] as FingerprintTable[],
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
})

const config = reactive({
    screenhost: false,
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

let fp = usePagination(form.fingerResult, 50)
let vp = usePagination(form.vulResult, 50)


// 初始化时调用
onMounted(async () => {
    let isSuccess = await InitRule()
    if (!isSuccess) {
        ElMessage.error({
            showClose: true,
            message: "初始化指纹规则失败，请检查配置文件",
        })
        return
    }
    FingerprintList().then(list => {
        dashboard.fingerLength = list.length
    })
    // 获取POC数量
    GetFingerPocMap().then(pocMap => {
        dashboard.yamlPocsLength = Object.keys(pocMap).length
    })
    // 遍历模板
    let files = await List(global.PATH.homedir + "/slack/config/pocs")
    files.forEach(file => {
        if (file.Path.endsWith(".yaml")) {
            allTemplate.value.push({
                label: file.BaseName,
                value: file.Path
            })
        }
    })
    templateLoading.value = false // 加载完毕
    // 获得结果回调
    EventsOn("nucleiResult", (result: any) => {
        const riskLevelKey = result.Risk as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        vp.table.result.push({
            vulID: result.ID,
            vulName: result.Name,
            protocoltype: result.Type.toLocaleUpperCase(),
            severity: result.Risk.toLocaleUpperCase(),
            vulURL: result.URL,
            request: result.Request,
            response: result.Response,
            extInfo: result.Extract,
            description: result.Description,
            reference: result.Reference,
        })
        vp.table.pageContent = vp.ctrl.watchResultChange(vp.table)
    });
    EventsOn("webFingerScan", (result: any) => {
        if (result.StatusCode == 0) {
            dashboard.reqErrorURLs.push(result.URL)
            return
        }
        fp.table.result.push({
            url: result.URL,
            status: result.StatusCode,
            length: result.Length,
            title: result.Title,
            detect: result.Detect,
            existsWaf: result.IsWAF,
            waf: "WAF: " + result.WAF,
            fingerprint: result.Fingerprints,
            screenshot: result.Screenshot
        })
        fp.table.pageContent = fp.ctrl.watchResultChange(fp.table)
        if (!config.writeDB) return
        InsertFingerscanResult(form.taskId, {
            URL: result.URL,
            StatusCode: result.StatusCode,
            Length: result.Length,
            Title: result.Title,
            Detect: result.Detect,
            IsWAF: result.IsWAF,
            WAF: result.WAF,
            Fingerprints: result.Fingerprints.join(","),
            Screenshot: result.Screenshot,
        })
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
        EventsOff("NucleiCounts");
        EventsOff("NucleiProgressID");
    };
});

async function startScan() {
    let ws = new Scanner
    if (config.writeDB && form.taskName == '') {
        ElMessage.warning({
            showClose: true,
            message: "任务名称不能为空",
        });
        return
    }
    if (!validateWebscan(form.url)) {
        return
    }
    form.taskId = nano()
    let isSuccess = await InsertScanTask(form.taskId, form.taskName, form.url, 0, 0)
    if (!isSuccess) {
        ElMessage.error("添加任务失败")
        return
    }
    form.newscanner = false // 隐藏界面
    ws.init()
    await ws.RunScanner()
}

function stopScan() {
    if (form.runnningStatus) {
        form.runnningStatus = false
        ExitScanner("[webscan]")
        ElMessage.warning("任务已停止");
    }
}

class Scanner {
    urls = [] as string[]
    public init() {
        fp.table.result = []
        vp.table.result = []
        vp.table.pageContent = vp.ctrl.watchResultChange(vp.table)
        dashboard.riskLevel.critical = 0
        dashboard.riskLevel.high = 0
        dashboard.riskLevel.medium = 0
        dashboard.riskLevel.low = 0
        dashboard.riskLevel.info = 0
        dashboard.reqErrorURLs = []
        dashboard.currentModule = ""
        form.percentage = 0
        form.nucleiPercentage = 0
        form.runnningStatus = true
    }
    public async RunScanner() {
        // 检查先行条件
        if (!await TestProxy(1)) {
            return
        }
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
            DeepScan: deepScan,
            RootPath: config.rootPathScan,
            CallNuclei: callNuclei,
            TemplateFiles: customTemplate.value,
            SkipNucleiWithoutTags: config.skipNucleiWithoutTags,
            GenerateLog4j2: config.generateLog4j2,
        }
        await NewWebScanner(options, getProxy())
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
    fofa: async function () {
        form.fofaDialog = false
        let urls = await LinkFOFA(form.query, form.fofaNum)
        if (urls != undefined) {
            form.url = urls.join("\n")
        }
    },
    hunter: async function () {
        form.hunterDialog = false
        let urls = await LinkHunter(form.query, form.defaultNum)
        if (urls != undefined) {
            form.url = urls.join("\n")
        }
    },

}

function filterHandlerSeverity(value: string, row: any): boolean {
    return row.severity === value;
}

function transformedResult() {
    return vp.table.result.map(({ vulID, vulName, protocoltype, severity, vulURL, extInfo }) => ({
        vulID,
        vulName,
        protocoltype,
        severity,
        vulURL,
        extInfo,
    }))
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
                icon: h(CopyDocument, defaultIconSize),
                onClick: () => {
                    Copy(row.url)
                }
            },
            {
                label: "复制选中链接",
                divided: true,
                icon: h(CopyDocument, defaultIconSize),
                onClick: () => {
                    CopyALL(fp.table.selectRows.map(item => item.url))
                }
            },
            {
                label: "打开链接",
                icon: h(ChromeFilled, defaultIconSize),
                onClick: () => {
                    BrowserOpenURL(row.url)
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
                    LinkDirsearch(row.url)
                }
            },
        ]
    });
}

function BatchBrowserOpenURL() {
    let targets = fp.table.selectRows.map(item => item.url)
    for (const item of targets) {
        BrowserOpenURL(item)
    }
}

// 选择目标文件读取
async function uploadFile() {
    form.url = await UploadFileAndRead()
}

async function ShowWebPictrue(filepath: string) {
    let bs64 = await ViewPictrue(filepath)
    if (bs64 == '') return
    screenDialog.value = true
    await nextTick()
    document.getElementById('webscan-img')!.setAttribute('src', bs64)
}


</script>

<template>
    <el-card style="margin-bottom: 10px;">
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
                                <el-button :icon="CopyDocument" @click="CopyALL(dashboard.reqErrorURLs)"
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
                    <el-table-column fixed prop="url" label="URL" width="300px" />
                    <el-table-column prop="status" width="100px" label="Code"
                        :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="length" width="100px" label="Length"
                        :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="title" label="Title" />
                    <el-table-column prop="fingerprint" label="Fingerprint" width="350px">
                        <template #default="scope">
                            <div class="finger-container">
                                <el-tag v-for="finger in scope.row.fingerprint" :key="finger"
                                    :effect="scope.row.detect === 'Default' ? 'light' : 'dark'"
                                    :type="global.webscan.highlight_fingerprints.includes(finger) ? 'danger' : 'primary'">{{
                                        finger }}</el-tag>
                                <el-tag type="danger" v-if="scope.row.existsWaf">{{ scope.row.waf }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column label="Screen" width="80">
                        <template #default="scope">
                            <el-button :icon="View" link @click="ShowWebPictrue(scope.row.screenshot)"
                                v-if="scope.row.screenshot != ''" />
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
                        :current-page="fp.table.currentPage" :page-sizes="[50, 100, 200]" :page-size="fp.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="fp.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
            <el-tab-pane label="漏洞">
                <el-table :data="vp.table.pageContent" stripe height="100vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column prop="vulID" label="Template" width="250px" />
                    <el-table-column prop="severity" width="150px" label="Severity"
                        :filter-method="filterHandlerSeverity" :filters="[
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
                            <span :class="getClassBySeverity(scope.row.severity)">{{ scope.row.severity }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="vulURL" label="URL" />
                    <el-table-column label="操作" width="120px">
                        <template #default="scope">
                            <el-button type="primary" link @click="detailDialog = true; selectedRow = scope.row">
                                查看详情
                            </el-button>
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
                        :current-page="vp.table.currentPage" :page-sizes="[50, 100, 200]" :page-size="vp.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="vp.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-space :size="2">
                <el-tooltip content="历史记录">
                    <el-button :icon="Clock" @click="historyDialog = true" />
                </el-tooltip>
                <el-tooltip content="导出Excel">
                    <el-button :icon="exportIcon"
                        @click="ExportWebScanToXlsx(transformArrayFields(fp.table.result), transformedResult())" />
                </el-tooltip>
            </el-space>
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
                    <el-button link size="small" :icon="Upload" @click="uploadFile"
                        style="margin-top: 5px;">导入目标文件</el-button>
                </el-form-item>
                <el-form-item>
                    <template #label>模式:
                        <el-tooltip>
                            <template #content>
                                1、指纹扫描: 只进行简单的指纹探测，不会探测敏感目录<br />
                                2、全指纹扫描: 会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                                3、指纹漏洞扫描: 指纹+主动敏感目录探测，扫描完成后扫描指纹对应POC<br />
                                4、专项扫描: 只扫描简单指纹和已选择的漏洞
                            </template>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <el-segmented v-model="config.webscanOption" :options="webscanOptions" block style="width: 100%;">
                        <template #default="{ item }">
                            <div style="display: flex; align-items: center">
                                <el-icon size="16" style="margin-right: 3px;">
                                    <component :is="item.icon" />
                                </el-icon>
                                <span>{{ item.label }}</span>
                            </div>
                        </template>
                    </el-segmented>
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
                    <el-input v-model="form.taskName" placeholder="必须输入不重复的任务名称" />
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
    <el-drawer v-model="detailDialog" size="70%">
        <template #header>
            <el-button text bg>
                <template #icon>
                    <Notebook />
                </template>漏洞详情</el-button>
        </template>
        <div v-if="selectedRow">
            <el-descriptions :column="1" border style="margin-bottom: 10px;">
                <el-descriptions-item label="Name:">{{ selectedRow.vulName }}</el-descriptions-item>
                <el-descriptions-item label="Extracted:">{{ selectedRow.extInfo }}</el-descriptions-item>
                <el-descriptions-item label="Description:">{{ selectedRow.description }}</el-descriptions-item>
                <el-descriptions-item label="Reference:" label-class-name="description">
                    <div v-for="item in selectedRow.reference.split(',')">
                        {{ item }}
                    </div>
                </el-descriptions-item>
            </el-descriptions>
            <div style="display: flex">
                <div class="pretty-response" style="font-size: small;">{{ selectedRow.request }}</div>
                <div class="pretty-response" style="font-size: small;">{{ selectedRow.response }}</div>
            </div>
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
    <el-drawer v-model="historyDialog" size="50%">
        <template #header>
            <el-text style="font-weight: bold; font-size: 16px;"><el-icon :size="18" style="margin-right: 5px;">
                    <Clock />
                </el-icon><span>扫描任务历史</span></el-text>
        </template>
        <el-table :cell-style="{ textAlign: 'center' }" :header-cell-style="{ 'text-align': 'center' }">
            <el-table-column prop="taskname" label="任务名称"></el-table-column>
            <el-table-column prop="targets" label="目标数量"></el-table-column>
            <el-table-column prop="failed" label="失败数量"></el-table-column>
            <el-table-column prop="vulnerability" label="漏洞数量"></el-table-column>
            <el-table-column label="操作" width="180px" align="center">
                <template #default="scope">

                </template>
            </el-table-column>
            <template #empty>
                <el-empty></el-empty>
            </template>
        </el-table>
    </el-drawer>
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
    padding-right: 20px;
    scrollbar-width: thin;
}
</style>