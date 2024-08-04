<script lang="ts" setup>
import { reactive, onMounted, ref, computed } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { Copy, SplitTextArea, deduplicateUrlFingerMap } from '../../util'
import { ExportToXlsx } from '../../export'
import { QuestionFilled, ChromeFilled, Menu, Promotion, CopyDocument, Grid, Search, ArrowUpBold, ArrowDownBold } from '@element-plus/icons-vue';
import { PortParse, IPParse, NewTcpScanner, HostAlive, IsRoot, NewSynScanner, StopPortScan, Callgologger, PortBrute, FingerScan, ActiveFingerScan, NucleiScanner, NucleiEnabled } from '../../../wailsjs/go/main/App'
import { ReadFile, FileDialog } from '../../../wailsjs/go/main/File'
import { BrowserOpenURL, EventsOn, EventsOff } from '../../../wailsjs/runtime'
import global from '../../global'
import async from 'async';
import { URLFingerMap, PortScanData, File } from '../../interface';
import usePagination from '../../usePagination';

// syn 扫描模式
onMounted(() => {
    // 扫描状态，把结果从后端传输到前端
    EventsOn("portScanLoading", (p: PortScanData) => {
        pagination.table.result.push({
            IP: p.IP,
            Port: p.Port,
            Server: p.Server,
            Link: p.Link,
            HttpTitle: p.HttpTitle,
        })
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table.result, pagination.table.currentPage, pagination.table.pageSize)
    });
    // 进度条
    EventsOn("progressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    // 扫描结束时让进度条为100
    EventsOn("scanComplete", () => {
        form.percentage = 100
        ctrl.runningStatus = false
    });
    return () => {
        EventsOff("portScanLoading");
        EventsOff("scanComplete");
        EventsOff("progressID");
    };
});

onMounted(async () => {
    form.isRoot = await IsRoot()
    selectItem(1); // 更新初始化显示
});


const form = reactive({
    activeName: '1',
    target: '',
    portlist: '',
    percentage: 0,
    id: 0,
    count: 0,
    radio: "2",
    currentAliveMoudle: "None",
    isSYN: false,
    isRoot: false,
    filter: '',
    defaultFilterGroup: "Fingerprint",
    hideDashboard: false,
})

const table = reactive({
    result: [] as PortScanData[],
    temp: [] as PortScanData[], // 用于存储过滤之前的数据，后续需要还原给result
    filterId: 0
})

let pagination = usePagination(table.result, 20)

async function uploadFile() {
    let path = await FileDialog("*.txt")
    if (!path) {
        return
    }
    let file: File = await ReadFile(path)
    if (file.Error) {
        ElMessage.warning(file.Message)
        return
    }
    form.target = file.Content!
}

const options = ({
    filterGroup: ["Host", "Port", "Fingerprint", "Link", "WebTitle"],
    aliveGroup: ["None", "ICMP", "Ping"],
    filterField: function () {
        const filter = form.filter.trim();
        if (table.filterId == 0) {
            console.log("filter id" + table.filterId)
            table.temp = pagination.table.result
            table.filterId++
        }
        if (filter) {
            switch (form.defaultFilterGroup) {
                case "Host":
                    pagination.table.result = table.temp.filter((p: PortScanData) => p.IP.includes(filter));
                    break
                case "Port":
                    pagination.table.result = table.temp.filter((p: PortScanData) => p.Port.toString().includes(filter));
                    break
                case "Fingerprint":
                    pagination.table.result = table.temp.filter((p: PortScanData) => p.Server.includes(filter));
                    break
                case "Link":
                    pagination.table.result = table.temp.filter((p: PortScanData) => p.Link.includes(filter));
                    break
                case "WebTitle":
                    pagination.table.result = table.temp.filter((p: PortScanData) => p.HttpTitle.includes(filter));
            }
        } else {
            pagination.table.result = table.temp;
        }
        pagination.table.currentPage = 1; // 重置分页
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table.result, pagination.table.currentPage, pagination.table.pageSize)
    }
})

const config = reactive({
    thread: 1000,
    timeout: 7,
    ping: false,
    webscan: false,
    crack: false,
})

function validateInput() {
    const ipPatterns = [
        /^(\d{1,3}\.){3}\d{1,3}$/, // 192.168.1.1
        /^(\d{1,3}\.){3}\d{1,3}\/(\d{1,2})$/, // 192.168.1.1/8, 192.168.1.1/16, 192.168.1.1/24
        /^(\d{1,3}\.){3}\d{1,3}-(\d{1,3}\.){3}\d{1,3}|(\d{1,3}\.){2}\d{1,3}|\d{1,3}\)$/, // 192.168.1.1-192.168.255.255, 192.168.1.1-255
        /^(\d{1,3}\.){3}\d{1,3}:\d{1,5}$/, // 192.168.0.1:6379
        /^!((\d{1,3}\.){3}\d{1,3}(\/\d+)?|(\d{1,3}\.){2}\d{1,3}|\d{1,3})$/, // !192.168.1.6/28
        /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/, // domain
    ];
    const lines = form.target.split('\n');
    return lines.every(line =>
        ipPatterns.some(pattern => pattern.test(line.trim()))
    );
}

async function NewScanner() {
    let mode = form.isSYN ? "syn" : "fully connected"
    Callgologger("info", "Port scan task loaded, current mode: " + mode)
    if (!validateInput()) {
            ElMessage({
                type: "warning",
                message: "输入目标格式不正确",
            })
            return
        }
    const ps = new Scanner()
    await ps.PortScanner()
}

class Scanner {
    public async PortScanner() {
        // 初始化
        pagination.ctrl.initTable()
        form.id = 0
        form.count = 0
        form.percentage = 0
        table.temp = []
        var ips = [] as string[]
        var portsList = [] as number[]
        if (!validateInput()) {
            ElMessage({
                type: "warning",
                message: "输入目标格式不正确",
            })
            return
        }
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
        // Callgologger("info", "targets cout: " + ips.length)
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
        if (form.isSYN) {
            await NewSynScanner(specialTarget,ips, portsList)
            return
        }
        await NewTcpScanner(specialTarget, ips, portsList, config.thread, config.timeout)
        // 恢复配置
        ctrl.runningStatus = false
        ips = []
        portsList = []
        Callgologger("info", "Portscan task is ending")
        if (config.webscan) {
            ips = moreOperate.getHTTPLinks(pagination.table.result)
            EzWebscan(ips)
        }
        if (config.crack) {
            ips = pagination.ctrl.getColumnData("Link")
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

async function EzWebscan(ips: string[]) {
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
            await NucleiScanner(0, ufm.url, ufm.finger, global.webscan.nucleiEngine, false, "", "")
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
        ElNotification.success({
            message: "Crack Finished",
            position: "bottom-right"
        })
    });
}

const moreOperate = ({
    getHTTPLinks: function (result: PortScanData[]): string[] {
        return result
            .filter(line => line.Link.startsWith("http"))
            .map(line => line.Link);
    },

    getBruteLinks: function (result: PortScanData[]): string[] {
        return result
            .filter(line => needBrute.some(brute => line.Link.startsWith(brute)))
            .map(line => line.Link);
    },
    // 复制端口扫描中的所有HTTP链接
    CopyURLs: function (type: string, result: PortScanData[]) {
        if (result.length <= 1) {
            ElNotification.warning({
                message: "复制内容条数需大于1",
                position: 'bottom-right',
            });
            return;
        }
        if (type == "url") {
            Copy(this.getHTTPLinks(result).join("\n"))
        } else {
            Copy(this.getBruteLinks(result).join("\n"))
        }
    },
    CopySelectLinks: function () {
        const selectRows = JSON.parse(JSON.stringify(pagination.table.selectRows));
        let targets: string[] = selectRows.map((item: any) => item.Link)
        Copy(targets.join("\n"))
    },
    // 联动
    Linkage: function (mode: string) {
        // 处理对象，不然map拿不到值
        const selectRows = JSON.parse(JSON.stringify(pagination.table.selectRows));
        let targets = selectRows.map((item: any) => item.Link)
        if (targets.length == 0) {
            ElMessage("至少选择1个联动目标")
            return
        }
        if (mode == "webscan") {
            ElNotification.success(`已发送${targets.length}个目标到网站扫描`)
            EzWebscan(targets)
        } else {
            ElNotification.success(`已发送${targets.length}个目标到暴破与未授权检测`)
            EzCrack(targets)
        }
    },
})

function changeTableHeigth() {
    form.hideDashboard = !form.hideDashboard
    var portscanTable = document.getElementById('portscan-table')!
    if (!form.hideDashboard) {
        portscanTable.style.height = '51.5vh'
    } else {
        portscanTable.style.height = '80.5vh'
    }
}

const titleStyle = computed(() => {
    return global.Theme.value ? {
        backgroundColor: '#333333',
    } : {
        backgroundColor: '#eee',
    };
})
</script>

<template>
    <div v-show="!form.hideDashboard" style="margin-bottom: 5px;">
        <el-card shadow="never">
            <el-row :gutter="8">
                <el-col :span="6" style="display: flex; align-items: center;">
                    <el-checkbox v-model="form.isSYN" :disabled="!form.isRoot">
                        <el-tooltip placement="right">
                            <template #content>
                                需要ROOT权限
                            </template>
                            SYN
                        </el-tooltip>
                    </el-checkbox>
                    <el-checkbox v-model="config.crack" label="口令猜测"></el-checkbox>
                    <el-checkbox v-model="config.webscan" label="网站扫描"></el-checkbox>
                </el-col>
                <el-divider direction="vertical" style="height: 4vh;" />
                <el-col :span="7" style="display: flex; align-items: center;">
                    <el-button link>存活探测</el-button>
                    <el-radio-group v-model="form.currentAliveMoudle" style="margin-left: 10px;">
                        <el-radio v-for="item in options.aliveGroup" :value="item">{{ item }}</el-radio>
                    </el-radio-group>
                </el-col>
                <el-divider direction="vertical" style="height: 4vh;" />
                <el-col :span="4">
                    <el-input v-model="config.timeout">
                        <template #prepend>
                            指纹超时
                        </template>
                    </el-input>
                </el-col>
                <el-col :span="4">
                    <el-input v-model="config.thread">
                        <template #prepend>
                            线程数量
                        </template>
                    </el-input>
                </el-col>
                <el-col :span="2">
                    <el-button type="primary" @click="NewScanner" v-if="!ctrl.runningStatus">开始扫描</el-button>
                    <el-button type="danger" @click="ctrl.stop" v-else>停止扫描</el-button>
                </el-col>
            </el-row>
        </el-card>
        <el-row :gutter="8" style="margin-top: 10px;">
            <el-col :span="10">
                <div class="my-header" :style="titleStyle">
                    <span>IP:
                        <el-tooltip placement="right-end">
                            <template #content>
                                目标支持换行分割,IP支持如下格式:<br />
                                192.168.1.1<br />
                                192.168.1.1/8<br />
                                192.168.1.1/16<br />
                                192.168.1.1/24<br />
                                192.168.1.1-192.168.255.255<br />
                                192.168.1.1-255<br /><br />
                                如果IP输入模式为192.168.0.1:6379此类形式，则只扫描该端口<br />
                                <br />
                                排除IP可以在可支持输入的IP格式前加!:<br />
                                !192.168.1.6/28<br />
                                <br />
                                域名格式: www.expamle.com
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
                <span class="my-header" :style="titleStyle">
                    预设端口:
                </span>
                <el-scrollbar class="list-container" max-height="130px" style="width: 100%">
                    <div class="list-item" v-for="(item, index) in global.portGroup"
                        :class="{ 'selected': selectedIndex === index }" @click="selectItem(index)">{{ item.text }}
                    </div>
                </el-scrollbar>
            </el-col>
            <el-col :span="10">
                <div class="my-header" :style="titleStyle">
                    端口列表:
                    <el-button size="small" @click="form.portlist = ''">清空</el-button>
                </div>
                <el-input class="input" type="textarea" rows="3" v-model="form.portlist" resize="none" />
            </el-col>
        </el-row>
    </div>
    <div style="position: relative;">
        <el-tabs v-model="form.activeName">
            <el-tab-pane label="结果输出" name="1">
                <el-table :data="pagination.table.pageContent" border id="portscan-table"
                    @selection-change="pagination.ctrl.handleSelectChange" style="height: 51.5vh;">
                    <el-table-column type="selection" width="42px" />
                    <el-table-column prop="IP" label="Host" />
                    <el-table-column prop="Port" label="Port" width="100px" />
                    <el-table-column prop="Server" label="Fingerprint" />
                    <el-table-column prop="Link" label="Link">
                        <template #default="scope">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.Link)"
                                v-show="scope.row.Link != ''">
                            </el-button>
                            {{ scope.row.Link }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="HttpTitle" label="WebTitle" />
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <el-progress :text-inside="true" :stroke-width="20" :percentage="form.percentage" color="#5DC4F7"
                        style="width: 40%;" />
                    <el-pagination background @size-change="pagination.ctrl.handleSizeChange"
                        @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100, 200, 500]"
                        :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                        :total="pagination.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-space>
                <el-button link @click="changeTableHeigth">
                    <template #icon>
                        <ArrowUpBold v-if="!form.hideDashboard" />
                        <ArrowDownBold v-else />
                    </template>
                </el-button>
                <el-input v-model="form.filter" placeholder="Filter">
                    <template #prepend>
                        <el-select v-model="form.defaultFilterGroup" style="width: 120px;">
                            <el-option v-for="name in options.filterGroup" :value="name">{{ name }}</el-option>
                        </el-select>
                    </template>
                    <template #suffix>
                        <el-button :icon="Search" @click="options.filterField" link></el-button>
                    </template>
                </el-input>
                <el-dropdown>
                    <el-button :icon="Menu" color="#D2DEE3" />
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="moreOperate.CopyURLs('url', pagination.table.result)"
                                :icon="CopyDocument">复制全部URL</el-dropdown-item>
                            <el-dropdown-item @click="moreOperate.CopyURLs('brute', pagination.table.result)"
                                :icon="CopyDocument">复制全部可爆破协议</el-dropdown-item>
                            <el-dropdown-item @click="moreOperate.CopySelectLinks()"
                                :icon="CopyDocument">复制选中目标</el-dropdown-item>
                            <el-dropdown-item :icon="Grid"
                                @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', pagination.table.result)" divided>
                                导出Excel</el-dropdown-item>
                            <el-dropdown-item @click="moreOperate.Linkage('webscan')" :icon="Promotion"
                                divided>发送至网站扫描</el-dropdown-item>
                            <el-dropdown-item @click="moreOperate.Linkage('crack')"
                                :icon="Promotion">发送至暴破与未授权检测</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </el-space>
        </div>
    </div>
</template>


<style>
.el-textarea__inner {
    height: 100%;
}

.list-container {
    width: 20vh;   
    border-radius: 4px;
    text-align: center;
    max-height: 120px;
}

.list-item.selected {
    background-color: #a4c1e7;
}

.input {
    height: 120px;
}
</style>
