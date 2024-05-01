<script lang="ts" setup>
import { reactive } from 'vue'
import { VideoPause, Search, QuestionFilled, Plus, ZoomIn, CopyDocument, Link, Operation } from '@element-plus/icons-vue';
import {
    InitRule,
    FofaSearch,
    HunterSearch,
    FingerLength,
    FingerScan,
    ActiveFingerScan,
    IsHighRisk
} from '../../../wailsjs/go/main/App'
import { ElMessage } from 'element-plus';
import { formatURL, ApiSyntaxCheck, splitInt, TestProxy, currentTime, Copy } from '../../util'
import { ExportWebScanToXlsx } from '../../export'
import async from 'async';
import global from "../../global"
import { onMounted } from 'vue';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
// 初始化时调用
onMounted(async () => {
    let err = await InitRule()
    if (err == "") {
        dashboard.fingerLength = await FingerLength()
    } else {
        ElMessage({
            showClose: true,
            message: "初始化指纹规则失败，请检查配置文件",
            type: "error"
        })
    }
});

interface Vulnerability {
    vulName: string
    severity: string
    vulURL: string
    request: string
    response: string
    extInfo: string
}

interface FingerprintTable {
    url: string
    status: string
    length: string
    title: string
    detect: string
    fingerprint: FingerLevel[]
}

interface FingerLevel {
    name: string
    level: number // level 0 is default , level 1 is high risk
}

const form = reactive({
    url: '',
    keyword: '',
    path: '',
    risk: [],
    riskOptions: ["critical", "high", "medium", "low", "info"],
    pocnums: '',
    newscanner: false,
    currentModule: '指纹扫描',
    module: ["指纹扫描", "指纹+敏感目录扫描", "指纹漏洞扫描", "全部漏洞扫描"],
    thread: 50,
    vulResult: [] as Vulnerability[],
    fingerResult: [] as FingerprintTable[],
    urlFingerMap: [] as uf[],
    numTips: '筛选当前条件下的POC数量',
    currentLoadPath: [] as string[],
    fofaDialog: false,
    fofaNum: 1000,
    hunterDialog: false,
    hunterNum: ["10", "20", "50", "100", "1000"],
    defaultNum: "100",
    query: '',
    runnningStatus: false,
})

interface uf {
    url: string,
    finger: string[]
}

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
})

const ctrl = reactive({
    exit: false,
    buttonDisabled: false,
})

async function startScan() {
    let ws = new Scanner
    form.runnningStatus = true
    form.newscanner = false // 隐藏界面
    ws.init()
    await ws.infoScanner()
}

function stopScan() {
    if (ctrl.exit === false) {
        ctrl.exit = true
        ElMessage({
            showClose: true,
            message: "任务已停止",
        });
    }
}

class Scanner {
    urls = [] as string[]
    public init() {
        form.urlFingerMap = [] as uf[]
        form.fingerResult = []
        form.vulResult = []
        dashboard.riskLevel.critical = 0
        dashboard.riskLevel.high = 0
        dashboard.riskLevel.medium = 0
        dashboard.riskLevel.low = 0
        dashboard.riskLevel.info = 0
        dashboard.reqErrorURLs = []
        dashboard.currentModule = ""
    }

    public async sortFinger(Fingerprints: any) {
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

    public async infoScanner() {
        // 检查先行条件
        if (!await TestProxy(1)) {
            return
        }
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
        // 指纹扫描      
        global.Logger.value += `${currentTime()} 网站扫描任务已加载，目标数量: ${this.urls.length}\n`
        global.Logger.value += '正在进行指纹扫描 ...\n'
        let count = 0
        dashboard.currentModule = form.currentModule
        async.eachLimit(this.urls, form.thread, async (target: string, callback: () => void) => {
            if (ctrl.exit === true) {
                return
            }
            let result = await FingerScan(target, global.proxy)
            if (result.StatusCode == 0) {
                dashboard.reqErrorURLs.push(target)
            } else {
                form.urlFingerMap.push({
                    url: result.URL,
                    finger: result.Fingerprints
                })
                let temp = await this.sortFinger(result.Fingerprints)
                form.fingerResult.push({
                    url: result.URL,
                    status: result.StatusCode,
                    length: result.Length,
                    title: result.Title,
                    detect: "Default",
                    fingerprint: temp,
                })
            }
            count++
            if (count == this.urls.length) { // 等任务全部执行完毕调用主动指纹探测
                global.Logger.value += `${currentTime()} 指纹扫描已完成\n`
                count = 0
                callback();
            }
        }, async () => {
            if (form.currentModule == "指纹+敏感目录扫描") {
                async.eachLimit(form.urlFingerMap, form.thread, async (ufm: uf, callback2: () => void) => {
                    if (ctrl.exit === true) {
                        return
                    }
                    let activeResult = await ActiveFingerScan(ufm.url, global.proxy)
                    if (activeResult.length > 0) {
                        activeResult.forEach(async item => {
                            let temp = await this.sortFinger(item.Fingerprints)
                            form.fingerResult.push({
                                url: item.URL,
                                status: item.StatusCode,
                                length: item.Length,
                                title: item.Title,
                                detect: "Active",
                                fingerprint: temp,
                            })
                            for (const fm of form.urlFingerMap) {
                                if (fm.url == item.URL) {
                                    fm.finger = Array.from(new Set(item.Fingerprints.push(fm.finger)))
                                }
                            }
                        }
                        )
                    }

                })


            }
            form.runnningStatus = false
        })
    }
}

// 联动空间引擎

const coordination = reactive({
    fofa: function () {
        if (ApiSyntaxCheck(0, global.space.fofaemail, global.space.fofakey, form.query) === false) {
            return
        }
        FofaSearch(form.query, form.fofaNum.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, true).then(result => {
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
                <span>Dashboard</span>
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
                <el-statistic title="已加载指纹数量" :value="dashboard.fingerLength" />
            </el-col>
            <el-col :span="4">
                <el-statistic :value="dashboard.reqErrorURLs.length">
                    <template #title>
                        <div style="display: inline-flex; align-items: center">
                            失败数/目标数
                            <el-popover placement="left" :width="350" trigger="hover"
                                v-if="dashboard.reqErrorURLs.length >= 1">
                                <el-scrollbar height="150px">
                                    <p v-for="u in dashboard.reqErrorURLs" class="scrollbar-demo-item">{{ u }}</p>
                                </el-scrollbar>
                                <template #reference>
                                    <el-icon style="margin-left: 4px" :size="12">
                                        <ZoomIn />
                                    </el-icon>
                                </template>
                            </el-popover>
                        </div>
                    </template>
                    <template #suffix>/{{ dashboard.count }}</template>
                </el-statistic>
            </el-col>
            <el-divider direction="vertical" style="height: 7vh;" />
            <el-col :span="6">
                <el-statistic title="扫描模式" :value="dashboard.currentModule" />
            </el-col>
        </el-row>
    </el-card>
    <div style="position: relative;">
        <el-tabs type="card">
            <el-tab-pane label="指纹">
                <el-table :data="form.fingerResult" border height="63vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column fixed type="index" label="#" width="60px" />
                    <el-table-column fixed prop="url" label="URL" width="300px">
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
                    <el-table-column prop="title" label="Title" width="300px" />
                    <el-table-column prop="detect" label="Detection" sortable width="120px" />
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
                                    :type="finger.level === 1 ? 'danger' : 'default'"
                                    :effect="scope.row.detect === 'Default' ? 'light' : 'dark'">{{ finger.name
                                    }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>
            <el-tab-pane label="漏洞">
                <el-table :data="form.vulResult" border height="63vh" :cell-style="{ textAlign: 'center' }"
                    :header-cell-style="{ 'text-align': 'center' }">
                    <el-table-column type="expand">
                        <template #default="props">
                            <h4>拓展信息: {{ props.row.extInfo }}</h4>
                            <el-space>
                                <div class="pretty-response" wrap="off" style="width: 63vh; height: 50%;">{{
                            props.row.request }}</div>
                                <div class="pretty-response" wrap="off" style="width: 63vh; height: 50%;">{{
                            props.row.response }}</div>
                            </el-space>
                        </template>
                    </el-table-column>
                    <el-table-column prop="vulName" label="Name" width="250px" :show-overflow-tooltip="true" />
                    <el-table-column prop="severity" width="150px" label="Risk" :sort-method="sortMethod" sortable
                        :show-overflow-tooltip="true">
                        <template #default="scope">
                            <span :class="getClassBySeverity(scope.row)">{{ scope.row.severity }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="vulURL" label="URL" :show-overflow-tooltip="true" />
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-popover placement="left" :width="200" trigger="click">
                <template #reference>
                    <el-button text bg :icon="Operation" style="margin-right: -5px;"></el-button>
                </template>
                <template #default>
                    <div class="context-menu">
                        <div class="click-div" @click="">
                            <el-text class="align">
                                <el-icon>
                                    <CopyDocument />
                                </el-icon>
                                复制全部URL链接
                            </el-text>
                        </div>
                    </div>
                </template>
            </el-popover>

            <el-button @click="ExportWebScanToXlsx(form.fingerResult, form.vulResult)">
                <template #icon>
                    <img src="/excle.svg" width="16">
                </template>
                导出Excle</el-button>
        </div>
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
                            2、指纹漏洞扫描: 指纹+敏感目录探测，扫描完成后扫描指纹对于POC<br />
                            3、敏感目录扫描: 会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                            4、全部漏洞扫描模式下具有更多的参数支持自定义<br />
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-segmented v-model="form.currentModule" :options="form.module" size="default" />
            </el-form-item>
            <div v-if="form.currentModule == '全部漏洞扫描'">
                <el-form-item label="关键字:" class="bottom">
                    <el-input v-model="form.keyword" placeholder="根据id和info.name判断','分割关键字" clearable></el-input>
                </el-form-item>
                <el-form-item label="风险等级:" class="bottom">
                    <el-select v-model="form.risk" placeholder="未选择即默认不进行筛选" multiple clearable style="width: 100%;">
                        <el-option v-for="item in form.riskOptions" :label="item.toLocaleUpperCase()" :value="item" />
                    </el-select>
                </el-form-item>
                <el-form-item label="POC数量:">
                    <el-text style="width: 80%; text-align: center;">
                        {{ form.numTips }}
                    </el-text>
                    <el-button :icon="Search" circle style="margin-left: auto;" @click=""
                        :disabled="ctrl.buttonDisabled" />
                </el-form-item>
            </div>
            <el-form-item label="指纹线程:">
                <el-input-number controls-position="right" v-model="form.thread" :min="1" :max="2000" />
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
            <span>
                <el-button type="primary" @click="coordination.fofa">
                    导入
                </el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="form.hunterDialog" title="导入鹰图目标(API单次查询最大100，导入1000条数据需等待)" width="50%" center>
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
            <span>
                <el-button type="primary" @click="coordination.hunter">
                    导入
                </el-button>
            </span>
        </template>
    </el-dialog>
</template>

<style scoped>
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
    gap: 5px;
}
</style>