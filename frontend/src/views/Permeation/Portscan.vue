<script lang="ts" setup>
import { reactive, onMounted, ref } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { Copy, SplitTextArea } from '../../util'
import { ExportToXlsx } from '../../export'
import { QuestionFilled, ChromeFilled, Menu, Promotion, CopyDocument, Grid } from '@element-plus/icons-vue';
import { PortParse, IPParse, NewTcpScanner, NewCorrespondsScan, HostAlive, IsRoot, NewSynScanner, StopPortScan, Callgologger, PortBrute, FingerScan, ActiveFingerScan, NucleiScanner, NucleiEnabled } from '../../../wailsjs/go/main/App'
import { GetFileContent, FileDialog } from '../../../wailsjs/go/main/File'
import { BrowserOpenURL, EventsOn, EventsOff } from '../../../wailsjs/runtime'
import global from '../../global'
import async from 'async';
interface PortScanData {
    IP: string
    Port: number
    Server: string
    Link: string
    HttpTitle: string
}

// syn 扫描模式
onMounted(async () => {
    // 扫描状态，把结果从后端传输到前端
    EventsOn("synPortScanLoading", (p: PortScanData) => {
        table.result.push({
            host: p.IP,
            port: p.Port,
            fingerprint: p.Server,
            link: p.Link,
            title: p.HttpTitle,
        })
        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
    });
    // 进度条
    EventsOn("synProgressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    // 扫描结束时让进度条为100
    EventsOn("synScanComplete", () => {
        form.percentage = 100
        ctrl.runningStatus = false
    });
    return () => {
        EventsOff("synPortScanLoading");
        EventsOff("synScanComplete");
        EventsOff("synProgressID");
    };
});

// 全连接 扫描模式
onMounted(() => {
    EventsOn("tcpPortScanLoading", (p: PortScanData) => {
        table.result.push({
            host: p.IP,
            port: p.Port,
            fingerprint: p.Server,
            link: p.Link,
            title: p.HttpTitle,
        })
        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
    });
    EventsOn("tcpProgressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    EventsOn("tcpScanComplete", () => {
        form.percentage = 100
        ctrl.runningStatus = false
    });
    return () => {
        EventsOff("tcpPortScanLoading");
        EventsOff("tcpScanComplete");
        EventsOff("tcpProgressID");
    };
})

// 特殊目标扫描
onMounted(() => {
    EventsOn("csPortScanLoading", (p: PortScanData) => {
        table.result.push({
            host: p.IP,
            port: p.Port,
            fingerprint: p.Server,
            link: p.Link,
            title: p.HttpTitle,
        })
        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
    });
    EventsOn("csProgressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    return () => {
        EventsOff("csPortScanLoading");
        EventsOff("csProgressID");
    };
})


onMounted(async () => {
    form.isRoot = await IsRoot()
    table.result = []
    table.pageContent = []
    table.burteResult = []
    selectItem(1); // 更新初始化显示
});

const activeName = '1'

const form = reactive({
    target: '',
    portlist: '',
    percentage: 0,
    id: 0,
    count: 0,
    radio: "2",
    aliveOptions: ["None", "ICMP", "Ping"],
    currentAliveMoudle: "None",
    isSYN: false,
    isRoot: false,
})

interface uf {
    url: string,
    finger: string[]
}

async function uploadFile() {
    let path = await FileDialog()
    form.target = await GetFileContent(path)
}

const config = reactive({
    thread: 1000,
    timeout: 7,
    ping: false,
    webscan: false,
    crack: false,
})

async function NewScanner() {
    let mode = form.isSYN ? "syn扫描" : "全连接扫描"
    Callgologger("info", "Port scan task loaded, current mode: " + mode)
    const ps = new Scanner()
    await ps.PortScanner()
}

class Scanner {
    public async PortScanner() {
        // 初始化
        table.result = []
        table.pageContent = []
        form.id = 0
        form.count = 0
        var ips = [] as string[]
        var portsList = [] as number[]
        const lines = SplitTextArea(form.target)
        let specialTarget = [] as string[]
        let conventionTarget = [] as string[]
        // 处理目标 192.168.1.1:6379 与其他形式
        for (const line of lines) {
            if (line.includes(":")) {
                specialTarget.push(line)
            } else {
                conventionTarget.push(line)
            }
        }
        // 处理端口和IP组
        portsList = await PortParse(form.portlist)
        ips = await IPParse(conventionTarget)
        // 判断扫描总数
        ips == null ? form.count = specialTarget.length : form.count = ips.length * portsList.length + specialTarget.length
        if (form.count == 0) {
            ElMessage({
                showClose: true,
                message: '可用目标或端口为空',
                type: 'warning',
            })
            return
        }
        ctrl.runningStatus = true
        // 判断是否进行存活探测
        if (form.currentAliveMoudle != "None") {
            form.currentAliveMoudle === 'ICMP' ? config.ping = false : config.ping = true
            ips = await HostAlive((ips as any), config.ping)
        }
        if (specialTarget.length != 0) {
            await NewCorrespondsScan(specialTarget, config.timeout)
        }
        if (form.isSYN) {
            await NewSynScanner(ips, portsList)
            return
        }
        await NewTcpScanner(ips, portsList, config.thread, config.timeout)
        // 恢复配置
        ctrl.runningStatus = false
        ips = []
        portsList = []
        Callgologger("info", "Portscan task is ending")
        if (config.webscan) {
            ips = getHTTPLinks(table.result)
            EzWebscan(ips)
        }
        if (config.crack) {
            ips = getColumnData("link")
            EzCrack(ips)
        }
    }
}

const ctrl = reactive({
    runningStatus: false,
    stop: function () {
        if (ctrl.runningStatus) {
            ctrl.runningStatus = false
            StopPortScan()
            ElMessage({
                showClose: true,
                message: '已停止任务',
                type: 'warning',
            })
        }
    },
})

const table = reactive({
    currentPage: 1,
    pageSize: 10,
    result: [{}],
    pageContent: [{}],
    burteResult: [{}],
    selectRows: [] as PortScanData[],
})

const selectedIndex = ref<number | null>(null);
const selectItem = (index: number): void => {
    updatePorts(index)
    selectedIndex.value = index;
}

function updatePorts(radio: number) {
    const index = Number(radio);
    if (index >= 0 && index < global.portGroup.length) {
        form.portlist = global.portGroup[index].value;
    }
}

function handleSizeChange(val: any) {
    table.pageSize = val;
    table.currentPage = 1;
    table.pageContent = table.result.slice(0, val)
}

function handleCurrentChange(val: any) {
    table.currentPage = val;
    table.pageContent = table.result.slice((val - 1) * table.pageSize, (val - 1) * table.pageSize + table.pageSize)
}

function getColumnData(prop: string): any[] {
    return table.result.map((item: any) => item[prop]);
}

// 复制端口扫描中的所有HTTP链接
function CopyURLs(result: {}[]) {
    // 避免控制台报错
    if (result.length <= 1) {
        ElNotification({
            message: "复制内容条数需大于1",
            type: "warning",
            position: 'bottom-right',
        });
        return;
    }
    Copy(getHTTPLinks(result).join("\n"))
}

function getHTTPLinks(result: {}[]) {
    const temp = [];
    for (const line of result) {
        if ((line as any)["link"].includes("http")) {
            temp.push((line as any)["link"]);
        }
    }
    return temp
}

function handleChange(rows: PortScanData[]) {
    table.selectRows = rows
}

async function EzWebscan(ips: string[]) {
    let id = 0
    global.webscan.urlFingerMap = []
    Callgologger("info", `正在将WEB目标联动网站扫描，共计加载目标: ${ips.length}`)
    await FingerScan(ips, global.proxy)
    await ActiveFingerScan(ips, global.proxy)
    if (await NucleiEnabled(global.webscan.nucleiEngine)) {
        const filteredUrlFingerprints = global.webscan.urlFingerMap
            .filter(item => item.finger.length > 0 && item.url)
            .map(item => ({ url: item.url, finger: item.finger }));
        async.eachLimit(filteredUrlFingerprints, 10, async (ufm: uf, callback: () => void) => {
            await NucleiScanner(0, ufm.url, ufm.finger, global.webscan.nucleiEngine, false, "", "")
            id++
            if (id == filteredUrlFingerprints.length) {
                callback()
            }
        }, (err: any) => {
            Callgologger("info", "Webscan Finished")
            ElNotification({
                type: "success",
                message: "Webscan Finished",
                position: "bottom-right"
            })
        })
    } else {
        Callgologger("error", `Nuclei引擎无效，无法进行漏洞扫描，已结束！`)
    }
}

const needBrute = ["ftp", "ssh", "telnet", "smb", "oracle", "mssql", "mysql", "rdp", "postgresql", "vnc", "redis", "memcached", "mongodb"]

function EzCrack(ips: string[]) {
    let id = 0
    async.eachLimit(ips, 20, async (target: string, callback: () => void) => {
        let protocol = target.split("://")[0]
        let userDict = global.dict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
        if (needBrute.includes(protocol)) {
            Callgologger("info", target + " is start brute")
        }
        await PortBrute(target, userDict, global.dict.passwords)
        id++
        if (id == ips.length) {
            callback()
        }
    }, (err: any) => {
        Callgologger("info", "PortBrute Finished")
        ElNotification({
            type: "success",
            message: "Crack Finished",
            position: "bottom-right"
        })
    });
}

function linkage(mode: string) {
    // 处理对象，不然map拿不到值
    const selectRows = JSON.parse(JSON.stringify(table.selectRows));
    let targets = selectRows.map((item: any) => item.link)
    if (targets.length == 0) {
        ElMessage("至少选择1个联动目标")
        return
    }
    if (mode == "webscan") {
        ElNotification({
            type: "success",
            message: `已发送${targets.length}个目标到网站扫描`,
        })
        EzWebscan(targets)
    } else {
        ElNotification({
            type: "success",
            message: `已发送${targets.length}个目标到暴破与未授权检测`,
        })
        EzCrack(targets)
    }
}
</script>

<template>
    <el-card style="width: 100%;">
        <el-row :gutter="8">
            <el-col :span="6" style="display: flex; align-items: center;">
                <el-checkbox v-model="form.isSYN" :disabled="!form.isRoot">
                    <el-tooltip placement="right" v-if="!form.isRoot">
                        <template #content>
                            非ROOT模式启动时，无法调用SYN扫描模式
                        </template>
                        SYN
                    </el-tooltip>
                    <span v-else>SYN</span>
                </el-checkbox>
                <el-checkbox v-model="config.crack" label="口令猜测"></el-checkbox>
                <el-checkbox v-model="config.webscan" label="网站扫描"></el-checkbox>
            </el-col>
            <el-divider direction="vertical" style="height: 4vh;" />
            <el-col :span="7" style="display: flex; align-items: center;">
                <span style="width: 80px;">存活探测:</span>
                <el-radio-group v-model="form.currentAliveMoudle">
                    <el-radio v-for="item in form.aliveOptions" :value="item">{{ item }}</el-radio>
                </el-radio-group>
            </el-col>
            <el-divider direction="vertical" style="height: 4vh;" />
            <el-col :span="4" style="display: flex; align-items: center;">
                <span style="width: 120px;">指纹超时:</span>
                <el-input v-model="config.timeout" />
            </el-col>
            <el-col :span="4" style="display: flex; align-items: center;">
                <span style="width: 120px;">线程数量:</span>
                <el-input v-model="config.thread" />
            </el-col>
            <el-col :span="2">
                <el-button type="primary" @click="NewScanner" v-if="!ctrl.runningStatus">开始扫描</el-button>
                <el-button type="danger" @click="ctrl.stop" v-else>停止扫描</el-button>
            </el-col>
        </el-row>
    </el-card>
    <el-row :gutter="8" style="margin-top: 10px;">
        <el-col :span="10">
            <div class="my-header" style="background-color: #eee;">
                <span>IP:
                    <el-tooltip placement="right-end">
                        <template #content>
                            目标支持换行分割,IP支持如下格式:<br />
                            192.168.1.1<br />
                            192.168.1.1/8<br />
                            192.168.1.1/16<br />
                            192.168.1.1/24<br />
                            192.168.1.1,192.168.1.2<br />
                            192.168.1.1-192.168.255.255<br />
                            192.168.1.1-255<br /><br />
                            如果IP输入模式为192.168.0.1:6379此类形式，则只扫描该端口<br />
                            <br />
                            排除IP可以在可支持输入的IP格式前加!:<br />
                            !192.168.1.6/28<br />
                        </template>
                        <el-icon style="width: 13px;">
                            <QuestionFilled />
                        </el-icon>
                    </el-tooltip>
                </span>
                <el-button size="small" @click="uploadFile">IP导入</el-button>
            </div>
            <el-input class="input" type="textarea" rows="3" v-model="form.target" resize="none" />
        </el-col>
        <el-col :span="4">
            <span class="my-header" style="background-color: #eee;">
                预设端口:
            </span>
            <el-scrollbar class="list-container" max-height="130px" style="width: 100%">
                <div class="list-item" v-for="(item, index) in global.portGroup"
                    :class="{ 'selected': selectedIndex === index }" @click="selectItem(index)">{{ item.text }}
                </div>
            </el-scrollbar>
        </el-col>
        <el-col :span="10">
            <div class="my-header" style="background-color: #eee;">
                端口列表:
                <el-button size="small" @click="form.portlist = ''">清空</el-button>
            </div>
            <el-input class="input" type="textarea" rows="3" v-model="form.portlist" resize="none" />
        </el-col>
    </el-row>
    <div style="position: relative;">
        <el-tabs v-model="activeName">
            <el-tab-pane label="结果输出" name="1">
                <el-table :data="table.pageContent" border style="height: 51vh;" @selection-change="handleChange">
                    <el-table-column type="selection" width="55px" />
                    <el-table-column prop="host" label="Host" />
                    <el-table-column prop="port" label="Port" width="100px" />
                    <el-table-column prop="fingerprint" label="Fingerprint" />
                    <el-table-column prop="link" label="Link">
                        <template #default="scope">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.link)"
                                v-show="scope.row.link != ''">
                            </el-button>
                            {{ scope.row.link }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="title" label="WebTitle" />
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <el-progress :text-inside="true" :stroke-width="20" :percentage="form.percentage" color="#5DC4F7"
                        style="width: 40%;" />
                    <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                        :current-page="table.currentPage" :page-sizes="[10, 20, 50]" :page-size="table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-space>
                <el-input></el-input>
                <el-dropdown>
                    <el-button :icon="Menu" color="#D2DEE3" />
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="CopyURLs(table.result)"
                                :icon="CopyDocument">复制全部URL</el-dropdown-item>
                            <el-dropdown-item :icon="Grid"
                                @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', table.result)">
                                导出Excel</el-dropdown-item>
                            <el-dropdown-item @click="linkage('webscan')" :icon="Promotion"
                                divided>发送至网站扫描</el-dropdown-item>
                            <el-dropdown-item @click="linkage('crack')" :icon="Promotion">发送至暴破与未授权检测</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </el-space>
        </div>
    </div>
</template>


<style>
.list-container {
    width: 20vh;
    background-color: #ffffff;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    border-radius: 4px;
    text-align: center;
    max-height: 120px;
}

.list-item:hover {
    cursor: pointer;
    background-color: #EEF5FE;
}

.list-item.selected {
    background-color: #a4c1e7;
}

.input {
    height: 120px;
}

.el-textarea__inner {
    height: 100%;
}
</style>
