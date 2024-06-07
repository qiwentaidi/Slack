<script lang="ts" setup>
import { reactive } from 'vue'
import { VideoPause, QuestionFilled, Plus, ZoomIn, CopyDocument, Link, Grid, RefreshRight, Menu } from '@element-plus/icons-vue';
import {
    InitRule,
    FofaSearch,
    HunterSearch,
    FingerLength,
    FingerScan,
    ActiveFingerScan,
    IsHighRisk,
    NucleiScanner,
    NucleiEnabled,
    WebPocLength,
    Callgologger
} from '../../../wailsjs/go/main/App'
import async from 'async';
import { ElMessage, ElNotification } from 'element-plus';
import { formatURL, ApiSyntaxCheck, TestProxy, Copy, CopyALL } from '../../util'
import { ExportWebScanToXlsx } from '../../export'
import global from "../../global"
import { onMounted } from 'vue';
import { BrowserOpenURL, EventsOn } from '../../../wailsjs/runtime/runtime';
// 初始化时调用
onMounted(() => {
    InitRule().then(err => {
        if (err) {
            FingerLength().then(leng => {
                dashboard.fingerLength = leng
            })
            WebPocLength().then((leng: number) => {
                dashboard.yamlPocsLength = leng
            })
        } else {
            ElMessage({
                showClose: true,
                message: "初始化指纹规则失败，请检查配置文件",
                type: "error"
            })
        }
    })

    NucleiEnabled(global.webscan.nucleiEngine).then(enable => {
        enable ? dashboard.nucleiEnabled = true : dashboard.nucleiEnabled = false
    })
});

onMounted(() => {
    EventsOn("nucleiResult", (result: any) => {
        const riskLevelKey = result.Risk as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        form.vulResult.push({
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
    });
    global.webscan.urlFingerMap = []
    EventsOn("webFingerScan", async (result: any) => {
        if (result.StatusCode == 0) {
            dashboard.reqErrorURLs.push(result.URL)
        } else {
            global.webscan.urlFingerMap.push({
                url: result.URL,
                finger: result.Fingerprints
            })
            let temp = await sortFinger(result.Fingerprints)
            form.fingerResult.push({
                url: result.URL,
                status: result.StatusCode,
                length: result.Length,
                title: result.Title,
                detect: "Default",
                existsWaf: result.IsWAF,
                waf: "WAF: " + result.WAF,
                fingerprint: temp,
            })
        }
    });
})

interface Vulnerability {
    vulID: string
    vulName: string
    protocoltype: string
    severity: string
    vulURL: string
    request: string
    response: string
    extInfo: string
    reference: string
    description: string
}

interface FingerprintTable {
    url: string
    status: string
    length: string
    title: string
    detect: string
    existsWaf: boolean
    waf: string
    fingerprint: FingerLevel[]
}

interface uf {
    url: string,
    finger: string[]
}

interface FingerLevel {
    name: string
    level: number // level 0 is default , level 1 is high risk
}

const form = reactive({
    url: '',
    keyword: '',
    risk: [] as string[],
    riskOptions: ["critical", "high", "medium", "low"],
    newscanner: false,
    currentModule: 0,
    module: [
        {
            label: "指纹扫描",
            value: 0
        },
        {
            label: "指纹+敏感目录扫描",
            value: 1
        },
        {
            label: "指纹漏洞扫描",
            value: 2
        },
        {
            label: "全部漏洞扫描",
            value: 3
        }
    ],
    thread: 50,
    vulResult: [] as Vulnerability[],
    fingerResult: [] as FingerprintTable[],
    currentLoadPath: [] as string[],
    fofaDialog: false,
    fofaNum: 1000,
    hunterDialog: false,
    hunterNum: ["10", "20", "50", "100"],
    defaultNum: "100",
    query: '',
    runnningStatus: false,
    noInteractsh: false,
})



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
    nucleiEnabled: false
})

// 导出文件
function convertFT(ft: FingerprintTable[]) {
    return ft.map(item => ({
        url: item.url,
        status: item.status,
        length: item.length,
        title: item.title,
        detect: item.detect,
        fingerprint: item.fingerprint.map(fingerprintItem => JSON.stringify(fingerprintItem)).join('|')
    }))
}

function convertVUL(vul: Vulnerability[]) {
    return vul.map(item => ({
        vulName: item.vulName,
        protocoltype: item.protocoltype,
        severity: item.severity,
        vulURL: item.vulURL
    }))
}

async function sortFinger(Fingerprints: any) {
    const temp = [] as FingerLevel[]
    for (const finger of Fingerprints) {
        if (await IsHighRisk(finger)) {
            temp.push({
                name: finger,
                level: 1
            })
        } else {
            temp.push({
                name: finger,
                level: 0
            })
        }
    }
    return temp
}

async function startScan() {
    let ws = new Scanner
    form.newscanner = false // 隐藏界面
    ws.init()
    await ws.infoScanner()
}

function stopScan() {
    if (form.runnningStatus) {
        form.runnningStatus = false
        ElMessage({
            showClose: true,
            message: "任务已停止",
        });
    }
}

class Scanner {
    urls = [] as string[]
    public init() {
        global.webscan.urlFingerMap = []
        form.fingerResult = []
        form.vulResult = []
        dashboard.riskLevel.critical = 0
        dashboard.riskLevel.high = 0
        dashboard.riskLevel.medium = 0
        dashboard.riskLevel.low = 0
        dashboard.riskLevel.info = 0
        dashboard.reqErrorURLs = []
        dashboard.currentModule = ""
        form.runnningStatus = true
    }
    public async infoScanner() {
        // 检查先行条件
        if (!await TestProxy(1)) {
            return
        }
        this.urls = await formatURL(form.url)
        this.urls = Array.from(new Set(this.urls))
        dashboard.count = this.urls.length
        if (this.urls.length == 0) {
            ElMessage({
                showClose: true,
                message: "可用目标为空",
                type: "warning",
            });
            return
        }
        dashboard.currentModule = form.module[form.currentModule].label
        // 指纹扫描      
        Callgologger("info", `Webscan task loaded, current: ${this.urls.length}，当前扫描模式: ${dashboard.currentModule}`)
        Callgologger("info", '正在进行指纹扫描 ...')
        await FingerScan(this.urls, global.proxy)
        if (dashboard.currentModule == "指纹+敏感目录扫描" || dashboard.currentModule == "指纹漏洞扫描") {
            const urlArray: string[] = global.webscan.urlFingerMap.map(item => item.url);
            Callgologger("info", '正在进行主动指纹探测 ...')
            await ActiveFingerScan(urlArray, global.proxy)
        }
        let count = 0
        let mode = 0
        if (dashboard.currentModule == "指纹漏洞扫描" || dashboard.currentModule == "全部漏洞扫描") {
            if (dashboard.currentModule == "全部漏洞扫描") {
                mode = 1
            }
            Callgologger("info", `正在进行${dashboard.currentModule} ...`)
            async.eachLimit(global.webscan.urlFingerMap, 10, async (ufm: uf, callback: () => void) => {
                if (!form.runnningStatus) {
                    return
                }
                await NucleiScanner(mode, ufm.url, ufm.finger, global.webscan.nucleiEngine, form.noInteractsh, form.keyword, form.risk.join(","))
                count++
                if (count == global.webscan.urlFingerMap.length) {
                    callback()
                }
            }, () => {
                Callgologger("info", "漏洞扫描已扫描结束")
                this.done()
                form.runnningStatus = false
            })
        } else {
            Callgologger("info", "漏洞扫描已扫描结束")
            this.done()
            form.runnningStatus = false
        }

    }
    public async done() {
        ElNotification({
            message: "Scanner Finished",
            position: "bottom-right",
        })
    }
}

// 联动空间引擎

const uncover = {
    fofa: function () {
        if (ApiSyntaxCheck(0, global.space.fofaemail, global.space.fofakey, form.query) === false) {
            return
        }
        form.fofaDialog = false
        ElMessage("正在导入FOFA数据，请稍后...")
        FofaSearch(form.query, form.fofaNum.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, true).then(result => {
            if (result.Status == false) {
                return
            }
            form.url = ""
            for (const item of result.Results) {
                form.url += item.URL + "\n"
            }
        })
    },
    hunter: function () {
        if (ApiSyntaxCheck(1, "", global.space.hunterkey, form.query) === false) {
            return
        }
        form.hunterDialog = false
        ElMessage("正在导入鹰图数据，请稍后...")
        HunterSearch(global.space.hunterkey, form.query, form.defaultNum, "1", "0", "3", false).then(result => {
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
            form.url = ""
            result.data.arr.forEach((item: any) => {
                form.url += item.url + "\n"
            });
        })
    }
}

// 排序
const levelMap: Record<string, number> = {
    'critical': 1,
    'high': 2,
    'medium': 3,
    'low': 4,
    'info': 5
};

const sortMethod = (a: any, b: any) => {
    return levelMap[a.level as keyof typeof levelMap] - levelMap[b.level as keyof typeof levelMap];
};

// 设置css样式
function getClassBySeverity(severity: string) {
    switch (severity) {
        case 'CRITICAL':
            return 'severity-critical';
        case 'HIGH':
            return 'severity-high';
        case 'MEDIUM':
            return 'severity-medium';
        case 'LOW':
            return 'severity-low';
        default:
            return 'severity-info';
    }
}
</script>

<template>
    <el-card>
        <template #header>
            <div class="card-header">
                <span>Dashboard</span>
                <el-alert title="Nuclei engine is enabled" type="success" :closable="false" center show-icon
                    style="width: 50%; height: 30px" v-if="dashboard.nucleiEnabled" />
                <el-alert title="Nuclei engine is underfind, please configure the nuclei path in the settings"
                    type="error" :closable="false" center show-icon style="width: 50%; height: 30px" v-else />
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
    <div style="position: relative; margin-top: 10px;">
        <el-tabs type="card">
            <el-tab-pane label="指纹">
                <el-table :data="form.fingerResult" border height="63vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column fixed type="index" label="#" width="60px" />
                    <el-table-column fixed prop="url" label="URL" width="300px" :show-overflow-tooltip="true">
                        <template #default="scope">
                            <el-popover placement="bottom" trigger="click">
                                <template #reference>
                                    <el-link :underline="false">{{ scope.row.url }}</el-link>
                                </template>
                                <template #default>
                                    <div class="context-menu">
                                        <div class="click-div" @click="Copy(scope.row.url)">
                                            <el-text class="align">
                                                <el-icon>
                                                    <CopyDocument />
                                                </el-icon>
                                                复制
                                            </el-text>
                                        </div>
                                        <div class="click-div" @click="BrowserOpenURL(scope.row.url)"
                                            style="width: 100%">
                                            <el-text class="align">
                                                <el-icon>
                                                    <Link />
                                                </el-icon>
                                                打开链接
                                            </el-text>
                                        </div>
                                    </div>
                                </template>
                            </el-popover>
                        </template>
                    </el-table-column>
                    <el-table-column prop="status" width="100px" label="Code"
                        :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="length" width="100px" label="Length"
                        :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable
                        :show-overflow-tooltip="true" />
                    <el-table-column prop="title" label="Title" />
                    <!-- <el-table-column prop="waf" label="WAF" width="120px" /> -->
                    <el-table-column fixed="right" prop="fingerprint" width="350px">
                        <template #header>
                            <el-tooltip placement="left">
                                <template #content>
                                    · 普通指纹会呈现浅蓝色，workflow中存在漏洞的指纹会呈现浅红色<br />
                                    · 若指纹标签会呈现填充红色，表示该指纹为敏感目录扫描得到<br />
                                </template>
                                <el-icon>
                                    <QuestionFilled size="24" />
                                </el-icon>
                            </el-tooltip>
                            Fingerprint
                        </template>
                        <template #default="scope">
                            <div class="finger-container">
                                <el-tag v-for="finger in scope.row.fingerprint" :key="finger.name"
                                    :type="finger.level === 1 ? 'danger' : 'primary'"
                                    :effect="scope.row.detect === 'Default' ? 'light' : 'dark'">{{ finger.name
                                    }}</el-tag>
                                <el-tag type="danger" v-if="scope.row.existsWaf">{{ scope.row.waf }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
            </el-tab-pane>
            <el-tab-pane label="漏洞">
                <el-table :data="form.vulResult" border height="63vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column type="expand">
                        <template #default="props">
                            <div>
                                <el-descriptions :column="2" border style="margin-bottom: 10px;">
                                    <el-descriptions-item label="Name:">{{ props.row.vulName }}</el-descriptions-item>
                                    <el-descriptions-item label="Extracted:">{{ props.row.extInfo
                                        }}</el-descriptions-item>
                                    <el-descriptions-item label="Description:" :span="2">{{ props.row.description
                                        }}</el-descriptions-item>
                                    <el-descriptions-item label="Reference:" :span="2" label-class-name="description">
                                        <div v-for="item in props.row.reference.split(',')">
                                            {{ item }}
                                        </div>
                                    </el-descriptions-item>
                                </el-descriptions>

                                <el-space>
                                    <div class="pretty-response" style="width: 80vh; height: 50vh;">{{
                        props.row.request }}</div>
                                    <div class="pretty-response" style="width: fit-content; height: 50vh;">{{
                        props.row.response }}</div>
                                </el-space>
                            </div>

                        </template>
                    </el-table-column>
                    <el-table-column prop="vulID" label="ID" width="250px" :show-overflow-tooltip="true" />
                    <el-table-column prop="protocoltype" label="Type" width="150px" />
                    <el-table-column prop="severity" width="150px" label="Risk" :sort-method="sortMethod" sortable>
                        <template #default="scope">
                            <span :class="getClassBySeverity(scope.row.severity)">{{ scope.row.severity }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="vulURL" label="URL" :show-overflow-tooltip="true" />
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <el-space class="custom_eltabs_titlebar" :size="5">

            <el-button :icon="RefreshRight" text bg type="danger" @click="NucleiEnabled(global.webscan.nucleiEngine)"
                v-show="!dashboard.nucleiEnabled">Reload
                Engine</el-button>

            <el-dropdown>
                <el-button :icon="Menu" color="#D2DEE3" />
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item @click="CopyALL(dashboard.reqErrorURLs)"
                            :icon="CopyDocument">复制全部失败目标</el-dropdown-item>
                        <el-dropdown-item :icon="Grid"
                            @click="ExportWebScanToXlsx(convertFT(form.fingerResult), convertVUL(form.vulResult))">
                            导出Excel</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </el-space>
    </div>


    <el-drawer v-model="form.newscanner" size="44%">
        <template #header>
            <h4>新建扫描任务</h4>
            <el-button link @click="form.fofaDialog = true">
                <template #icon>
                    <img src="/fofa.ico">
                </template>
                FOFA
            </el-button>
            <el-divider direction="vertical" />
            <el-button link @click="form.hunterDialog = true">
                <template #icon>
                    <img src="/hunter.ico" style="width: 16px; height: 16px;">
                </template>
                Hunter
            </el-button>
        </template>
        <el-form label-width="auto">
            <el-form-item>
                <template #label>目标:
                    <el-tooltip placement="left">
                        <template #content>
                            支持如下输入方式:<br />
                            192.168.0.1<br />
                            192.168.0.1:443<br />
                            https://192.168.0.1
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input v-model="form.url" type="textarea" rows="6" clearable></el-input>
            </el-form-item>
            <el-form-item>
                <template #label>模式:
                    <el-tooltip placement="left">
                        <template #content>
                            1、指纹扫描: 只进行指纹，不会探测敏感目录<br />
                            2、指纹漏洞扫描: 指纹+敏感目录探测，扫描完成后扫描指纹对应POC<br />
                            3、敏感目录扫描: 会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                            4、全部漏洞扫描模式下具有更多的参数支持自定义<br />
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-segmented v-model="form.currentModule" :options="form.module" style="width: 100%;" />
            </el-form-item>
            <div v-if="dashboard.currentModule == '全部漏洞扫描'">
                <el-form-item label="关键字:" class="bottom">
                    <el-input v-model="form.keyword" placeholder="根据id判断','分割关键字" clearable></el-input>
                </el-form-item>
                <el-form-item label="风险等级:" class="bottom">
                    <el-select v-model="form.risk" placeholder="未选择即默认不进行筛选" multiple clearable style="width: 100%;">
                        <el-option v-for="item in form.riskOptions" :label="item.toLocaleUpperCase()" :value="item" />
                    </el-select>
                </el-form-item>
            </div>
            <el-form-item label="指纹线程:">
                <el-input-number controls-position="right" v-model="form.thread" :min="1" :max="2000" />
            </el-form-item>
            <el-form-item label="禁用反连:">
                <el-checkbox v-model="form.noInteractsh" />
            </el-form-item>
        </el-form>
        <div>
            <el-button type="primary" @click="startScan"
                style="left: 45%;bottom: 10px; position: absolute;">开始任务</el-button>
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
            <el-button type="primary" @click="uncover.fofa">
                导入
            </el-button>
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
            <el-button type="primary" @click="uncover.hunter">
                导入
            </el-button>
        </template>
    </el-dialog>
</template>

<style>
.el-col {
    text-align: center;
}

.severity-critical {
    color: #AF63F6;
}

.severity-high {
    color: red;
}

.severity-medium {
    color: orange;
}

.severity-low {
    color: rgb(0, 145, 255);
}

.severity-info {
    color: green;
}

.scrollbar-demo-item {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 50px;
    margin: 10px;
    text-align: center;
    border-radius: 4px;
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
}

.align {
    display: flex;
    justify-content: left;
    align-items: center;
}

.align .el-icon {
    margin-right: 5px;
}

.finger-container {
    flex-wrap: wrap;
    display: flex;
    gap: 7px;
}

.description {
    width: 100px;
}

.el-text {
    display: flex;
    align-items: center;
    text-align: center
}
</style>