<script lang="ts" setup>
import { reactive } from 'vue'
import { Monitor, ArrowDown, Search, QuestionFilled, Plus, ZoomIn } from '@element-plus/icons-vue';
import {
    FingerScan,
    InitRule,
    LocalWalkFiles,
    Webscan,
    PocNums,
    GetFingerPoc,
    FofaSearch,
    HunterSearch,
    GoFetch
} from '../../../wailsjs/go/main/App'
import { ExecutionPath } from '../../../wailsjs/go/main/File'
import { ElMessage } from 'element-plus';
import { formatURL, ApiSyntaxCheck, splitInt } from '../../util'
import async from 'async';
import global from "../../global"
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    form.fingerResult = []
});

interface Vulnerability {
    vulName: string
    severity: string
    vulURL: string
    request: string
    response: string
    extInfo: string
}

const form = reactive({
    url: '',
    keyword: '',
    path: '',
    risk: [],
    riskOptions: ["critical", "high", "medium", "low", "info"],
    pocnums: '',
    newscanner: false,
    currentModule: '仅指纹扫描',
    module: ["仅指纹扫描", "仅指纹扫描+主动目录探测", "指纹漏洞扫描", "全部漏洞扫描"],
    thread: 50,
    vulResult: [] as Vulnerability[],
    fingerResult: [{}],
    urlFingerMap: [] as uf[],
    numTips: '筛选当前条件下的POC数量',
    dialogVisible: false,
    currentLoadPath: [] as string[],
    fofaDialog: false,
    fofaNum: 1000,
    hunterDialog: false,
    hunterNum: ["10", "20", "50", "100", "1000"],
    defaultNum: "100",
    query: '',
})

interface uf {
    url: string,
    finger: string[]
}

const dashboard = reactive({
    reqErrorURLs: [] as string[],
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    info: 0,
    runningStatus: "" as string,
    count: 0,
    logger: '',
    request: '',
    response: '',
    extInfo: '',
})

const pathActive = "/config/active-detect"
const pathAFG = "/config/afrog-pocs"

const ctrl = reactive({
    exit: false,
    buttonDisabled: false,
})

async function startScan() {
    let ws = new Scanner
    form.newscanner = false
    ws.initData()
    await ws.infoScanner()
}

function stopScan() {
    if (ctrl.exit === false) {
        ctrl.exit = true
    }
}

interface WebResult {
    VulName: string
    Severity: string
    VulURL: string
    Request: string
    Response: string
    ExtInfo: string
}

class Scanner {
    urls = [] as string[]
    public initData() {
        form.urlFingerMap = [] as uf[]
        form.fingerResult = []
        form.vulResult = []
        dashboard.critical = 0
        dashboard.high = 0
        dashboard.medium = 0
        dashboard.low = 0
        dashboard.info = 0
        dashboard.reqErrorURLs = []
        dashboard.runningStatus = ""
        dashboard.logger = ''
        return
    }

    public async infoScanner() {
        // 检查先行条件
        this.urls = await formatURL(form.url)
        dashboard.count = this.urls.length
        if (this.urls.length == 0) {
            ElMessage({
                showClose: true,
                message: "可用目标为空",
                type: "warning",
            });
            return
        }
        const proxyURL = global.proxy.mode.toLowerCase() + "://" + global.proxy.address + ":" + global.proxy.port
        if (global.proxy.enabled && (await GoFetch("GET", proxyURL, "", [{}], 10, null))) {
            ElMessage({
                showClose: true,
                message: "代理地址不可达",
                type: "warning",
            });
            return
        }
        var date = new Date();
        // 
        // 
        // 指纹扫描 =====================================
        //         
        // 
        await InitRule()
        dashboard.logger += `${date.toLocaleString()} 已加载任务目标: ${this.urls.length}个\n`
        dashboard.logger += '[INFO] 正在进行指纹扫描\n'
        let count = 0
        async.eachLimit(this.urls, form.thread, async (target: string, callback: () => void) => {
            if (ctrl.exit === true) {
                return
            }
            dashboard.runningStatus = target
            let result = await FingerScan(target, global.proxy)
            if (result.StatusCode == 0) {
                dashboard.reqErrorURLs.push(target)
            } else {
                form.urlFingerMap.push({
                    url: result.URL,
                    finger: result.Fingerprints
                })
                form.fingerResult.push({
                    url: result.URL,
                    status: result.StatusCode,
                    length: result.Length,
                    title: result.Title,
                    fingerprint: result.Fingerprints,
                })
            }
            count++
            if (count == this.urls.length) { // 等任务全部执行完毕调用主动指纹探测
                dashboard.logger += `[END] 主动目录探测已结束\n`
                dashboard.logger += `[INFO] 正在初始化主动指纹探测任务，已加载主动指纹: ${form.currentLoadPath.length}个\n`
                count = 0
                form.currentLoadPath = await LocalWalkFiles(await ExecutionPath() + pathActive) // 初始化主动指纹目录
                callback();
            }
        }, (err: any) => {
            if (form.currentModule !== "仅指纹扫描") {
                // 主动指纹探测
                async.eachLimit(this.urls, form.thread, (target: string, callback: () => void) => {
                    if (ctrl.exit === true) {
                        return
                    }
                    dashboard.logger += `[INFO] ${target}，正在进行主动指纹探测\n`
                    dashboard.runningStatus = target
                    Webscan(target, "", "", form.currentLoadPath, global.proxy).then((result: WebResult[]) => {
                        if (Array.isArray(result)) {
                            for (const item of result) {
                                switch (item.Severity) {
                                    case "CRITICAL":
                                        dashboard.critical += 1
                                        break
                                    case "HIGH":
                                        dashboard.high += 1
                                        break
                                    case "MEDIUM":
                                        dashboard.medium += 1
                                        break
                                    case "LOW":
                                        dashboard.low += 1
                                        break
                                    case "INFO":
                                        dashboard.info += 1
                                }
                                form.vulResult.push({
                                    vulName: item.VulName,
                                    severity: item.Severity,
                                    vulURL: item.VulURL,
                                    request: item.Request,
                                    response: item.Response,
                                    extInfo: item.ExtInfo
                                })
                                if (item.Severity == "INFO") {
                                    for (let i = 0; i < form.urlFingerMap.length; i++) {
                                        if (form.urlFingerMap[i].url === target) {
                                            form.urlFingerMap[i].finger.push(item.VulName.split("-")[0]);
                                            break;
                                        }
                                    }
                                }
                            }
                        }
                    })
                    count++
                    if (count == this.urls.length) { // 等任务全部执行完毕调用主动指纹探测
                        dashboard.logger += `[END] 主动目录探测已结束\n`
                        callback();
                    }
                }, (err: any) => {
                    this.webScanner()
                })
            }
        });

    }

    public async webScanner() {
        if (form.currentModule == "指纹漏洞扫描") {
            dashboard.logger += `[INFO] 正在进行指纹漏洞扫描\n`
            let count = 0
            for (const uf of form.urlFingerMap) { // 统计能扫到指纹的目标数量
                if (uf.finger.length > 0) {
                    count++
                } 
            }
            async.eachLimit(form.urlFingerMap, 5, async (ufm: uf, callback: () => void) => {
                if (ctrl.exit === true) {
                    return
                }
                if (ufm.finger.length > 0) {
                    form.currentLoadPath = await GetFingerPoc(ufm.finger)
                    dashboard.logger += `[INFO] ${ufm.url}，已加载漏洞: ${form.currentLoadPath.length}个\n`
                    dashboard.runningStatus = ufm.url
                    Webscan(ufm.url, "", "", form.currentLoadPath, global.proxy).then((result: WebResult[]) => {
                        if (Array.isArray(result)) {
                            for (const item of result) {
                                switch (item.Severity) {
                                    case "CRITICAL":
                                        dashboard.critical += 1
                                        break
                                    case "HIGH":
                                        dashboard.high += 1
                                        break
                                    case "MEDIUM":
                                        dashboard.medium += 1
                                        break
                                    case "LOW":
                                        dashboard.low += 1
                                        break
                                    case "INFO":
                                        dashboard.info += 1
                                }
                                form.vulResult.push({
                                    vulName: item.VulName,
                                    severity: item.Severity,
                                    vulURL: item.VulURL,
                                    request: item.Request,
                                    response: item.Response,
                                    extInfo: item.ExtInfo
                                })
                            }
                        }
                    })
                    count--
                    if (count == 0) {
                        dashboard.logger += `[END] 指纹漏洞扫描已结束 :)\n`
                    }
                }
            })
        } else if (form.currentModule == "全部漏洞扫描") {
            form.currentLoadPath = await LocalWalkFiles(await ExecutionPath() + pathAFG)
            dashboard.logger += `[INFO] 正在初始化全漏洞扫描任务，已加载POC: ${form.currentLoadPath.length}个\n`
            let count = 0
            async.eachLimit(this.urls, form.thread, (target: string, callback: () => void) => {
                if (ctrl.exit === true) {
                    return
                }
                Webscan(target, form.risk.join(","), form.keyword, form.currentLoadPath, global.proxy).then((result: WebResult[]) => {
                    dashboard.runningStatus = target
                    for (const item of result) {
                        switch (item.Severity) {
                            case "CRITICAL":
                                dashboard.critical += 1
                                break
                            case "HIGH":
                                dashboard.high += 1
                                break
                            case "MEDIUM":
                                dashboard.medium += 1
                                break
                            case "LOW":
                                dashboard.low += 1
                                break
                            case "INFO":
                                dashboard.info += 1
                        }
                        form.vulResult.push({
                            vulName: item.VulName,
                            severity: item.Severity,
                            vulURL: item.VulURL,
                            request: item.Request,
                            response: item.Response,
                            extInfo: item.ExtInfo
                        })
                    }
                })
                count++
                if (count == this.urls.length) {
                    dashboard.logger += `[END] 全部漏洞扫描结束 :)\n`
                }
            })
        }
    }
}

// 显示信息
function showInfo(row: any) {
    form.dialogVisible = true
    dashboard.request = row.request
    dashboard.response = row.response
    dashboard.extInfo = row.extInfo
}

// 联动空间引擎

const coordination = reactive({
    fofa: function () {
        if (ApiSyntaxCheck(0, global.space.fofaemail, global.space.fofakey, form.query) === false) {
            return
        }
        FofaSearch(form.query, form.fofaNum.toString(), "1", global.space.fofaemail, global.space.fofakey, true, true).then(result => {
            if (result.Status == false) {
                return
            }
            form.url = ""
            for (const item of result.Results) {
                form.url += item.URL + "\n"
            }
            form.fofaDialog = false
        })
    },
    hunter: async function () {
        if (ApiSyntaxCheck(1, "", global.space.hunterkey, form.query) === false) {
            return
        }
        form.url = ""
        if (form.defaultNum == "1000") {
            let index = 0
            for (const num of splitInt(Number(form.defaultNum), 100)) {
                index += 1
                await HunterSearch(global.space.hunterkey, form.query, "100", "1", "0", "3", false).then(result => {
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
                        form.url += item.url + "\n"
                    });
                })
            }
        } else {
            await HunterSearch(global.space.hunterkey, form.query, form.defaultNum, "1", "0", "3", false).then(result => {
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
                    form.url += item.url + "\n"
                });
            })
        }
        form.hunterDialog = false
    }
})

async function countPocNums() {
    ctrl.buttonDisabled = true
    form.numTips = (await PocNums(form.risk.join(","), form.keyword)).toString()
    ctrl.buttonDisabled = false
}

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

function getClassBySeverity(row: any) {
    switch (row.severity) {
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
    <el-card style="margin-bottom: 10px;">
        <template #header>
            <div class="card-header">
                <span>仪表盘</span>
                <el-button type="primary" :icon="Plus" @click="form.newscanner = true">新建任务</el-button>
            </div>
        </template>
        <el-row>
            <el-col :span="2">
                <el-statistic title="紧急" :value="dashboard.critical" />
            </el-col>
            <el-col :span="2">
                <el-statistic title="高危" :value="dashboard.high" />
            </el-col>
            <el-col :span="2">
                <el-statistic title="中危" :value="dashboard.medium" />
            </el-col>
            <el-col :span="2">
                <el-statistic title="低危" :value="dashboard.low" />
            </el-col>
            <el-col :span="2">
                <el-statistic title="信息" :value="dashboard.info" />
            </el-col>
            <el-divider direction="vertical" style="height: 7vh;" />
            <el-col :span="3">
                <el-statistic title="目标总数" :value="dashboard.count" />
            </el-col>
            <el-col :span="3">
                <el-statistic :value="dashboard.reqErrorURLs.length">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            请求失败目标
                            <el-popover placement="left" title="可能是由于无法识别到协议而无法扫描的地址" :width="350" trigger="hover">
                                    <el-input v-model="dashboard.reqErrorURLs" type="textarea" rows="5"></el-input>
                                <template #reference>
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <ZoomIn />
                                    </el-icon>
                                </template>
                            </el-popover>
                        </div>
                    </template>
                </el-statistic>
            </el-col>
            <el-divider direction="vertical" style="height: 7vh;" />
            <el-col :span="7">
                <el-statistic title="当前目标" :value="dashboard.runningStatus" />
            </el-col>
        </el-row>
    </el-card>
    <el-tabs type="card" tab-position="left">
        <el-tab-pane label="指纹">
            <el-table :data="form.fingerResult" border height="65vh" :cell-style="{ textAlign: 'center' }"
                :header-cell-style="{ 'text-align': 'center' }">
                <el-table-column type="index" label="#" width="60px" />
                <el-table-column prop="url" label="网站地址" />
                <el-table-column prop="status" width="100px" label="状态码"
                    :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable
                    :show-overflow-tooltip="true" />
                <el-table-column prop="length" width="100px" label="长度"
                    :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable
                    :show-overflow-tooltip="true" />
                <el-table-column prop="title" label="标题" width="300px" />
                <el-table-column prop="fingerprint" label="指纹">
                    <template #default="scope">
                        <div style="flex-wrap: wrap; display: flex;">
                            <el-tag v-for="finger in scope.row.fingerprint" style="margin: 3px;">{{ finger }}</el-tag>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="漏洞">
            <el-table :data="form.vulResult" border height="65vh" :cell-style="{ textAlign: 'center' }"
                :header-cell-style="{ 'text-align': 'center' }">
                <el-table-column type="index" label="#" width="60px" />
                <el-table-column prop="vulName" label="漏洞名称" width="250px" :show-overflow-tooltip="true" />
                <el-table-column prop="severity" width="150px" label="风险等级" :sort-method="sortMethod" sortable
                    :show-overflow-tooltip="true">
                    <template #default="scope">
                        <span :class="getClassBySeverity(scope.row)">{{ scope.row.severity }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="vulURL" label="漏洞地址" :show-overflow-tooltip="true" />
                <el-table-column fixed="right" label="详情" width="55px">
                    <template #default="scope">
                        <el-button link :icon="ZoomIn" @click.prevent="showInfo(scope.row)">
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="日志">
            <el-input class="log-textarea" v-model="dashboard.logger" type="textarea" rows="20" readonly></el-input>
        </el-tab-pane>
    </el-tabs>

    <el-drawer v-model="form.newscanner">
        <template #header>
            <span>新建扫描任务</span>
            <el-dropdown>
                <el-button :icon="Monitor" size="small">联动<el-icon><arrow-down /></el-icon>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item @click="form.fofaDialog = true">FOFA</el-dropdown-item>
                        <el-dropdown-item @click="form.hunterDialog = true">鹰图</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </template>
        <el-dialog v-model="form.fofaDialog" title="导入FOFA目标(MAX 10000)" width="50%" center>
            <el-form :model="form" label-width="70px">
                <el-form-item label="查询条件">
                    <el-input v-model="form.query"></el-input>
                </el-form-item>
                <el-form-item label="导入数量">
                    <el-input-number v-model="form.fofaNum" :min="1" :max="10000" />
                </el-form-item>
            </el-form>
            <template #footer>
                <span>
                    <el-button type="primary" @click="coordination.fofa">
                        导入
                    </el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="form.hunterDialog" title="导入鹰图目标(API单次查询最大100，导入1000条数据需等待)" width="50%" center>
            <el-form :model="form" label-width="70px">
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
                <span>
                    <el-button type="primary" @click="coordination.hunter">
                        导入
                    </el-button>
                </span>
            </template>
        </el-dialog>
        <el-form label-width="80px" style="margin-top: auto;">
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
                            1、仅指纹扫描: 只扫指纹，不会进行敏感目录探测<br />
                            2、指纹POC扫描: 会进行敏感目录探测，且打指纹POC<br />
                            2、主动指纹探测: 勾选仅指纹扫描之后会出现，开启会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                            4、FastJson: 由于默认是不会进行通用FastJson漏洞扫描如果需要进行FastJson漏洞检测请输入关键字fastjson-1`<br />
                            5、全部漏洞扫描模式下具有更多的参数支持自定义<br />
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-select v-model="form.currentModule" style="width: 100%;">
                    <el-option v-for="item in form.module" :label="item" :value="item" />
                </el-select>
            </el-form-item>
            <div v-if="form.currentModule == '全部漏洞扫描'">
                <el-form-item label="关键字:" class="bottom">
                    <el-input v-model="form.keyword" placeholder="根据id和info.name判断','分割关键字" clearable></el-input>
                </el-form-item>
                <el-form-item label="风险等级:" class="bottom">
                    <el-select v-model="form.risk" placeholder="未选择即默认不进行筛选" multiple clearable style="width: 100%;">
                        <el-option v-for="item in form.riskOptions" :label="item" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="POC数量:">
                    <el-text style="width: 80%; text-align: center;">
                        {{ form.numTips }}
                    </el-text>
                    <el-button :icon="Search" circle style="margin-left: auto;" @click="countPocNums"
                        :disabled="ctrl.buttonDisabled" />
                </el-form-item>
            </div>
            <el-form-item label="指纹线程:">
                <el-input-number controls-position="right" v-model="form.thread" :min="1" :max="2000" />
            </el-form-item>
            <div style="margin-top: 10px;">
                <el-button type="primary" @click="startScan">开始任务</el-button>
                <el-button type="danger" @click="stopScan">停止任务</el-button>
            </div>
        </el-form>
    </el-drawer>
    <el-dialog v-model="form.dialogVisible" title="数据包详情" width="80%">
        <h4>拓展信息: {{ dashboard.extInfo }}</h4>
        <el-space>
            <el-input v-model="dashboard.request" :rows="20" type="textarea" resize="none" wrap="off"
                style="width: 63vh;" />
            <el-input v-model="dashboard.response" :rows="20" type="textarea" resize="none" wrap="off"
                style="width: 63vh;" />
        </el-space>
    </el-dialog>
</template>

<style scoped>
.el-col {
    text-align: center;
}

.severity-critical {
    color: purple;
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
</style>