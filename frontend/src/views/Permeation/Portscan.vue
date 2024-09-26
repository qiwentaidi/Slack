<script lang="ts" setup>
import { reactive, onMounted, ref, h } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { Copy, SplitTextArea } from '@/util'
import { ExportToXlsx } from '@/export'
import { QuestionFilled, ChromeFilled, Promotion, CopyDocument, Search } from '@element-plus/icons-vue';
import { PortParse, IPParse, NewTcpScanner, HostAlive, IsRoot, NewSynScanner, StopPortScan, Callgologger, SpaceGetPort } from 'wailsjs/go/main/App'
import { ReadFile, FileDialog } from 'wailsjs/go/main/File'
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime'
import global from '@/global'
import { PortScanData, File } from '@/interface';
import usePagination from '@/usePagination';
import exportIcon from '@/assets/icon/doucment-export.svg'
import { defaultIconSize, titleStyle } from '@/stores/style';
import { LinkCrack, LinkWebscan } from '@/linkage';
import ContextMenu from '@imengyu/vue3-context-menu'
import CustomTabs from '@/components/CustomTabs.vue';
import { isPrivateIP, validateIp, validatePortscan } from '@/stores/validate';
import consoleIcon from '@/assets/icon/console.svg'

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
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
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
    isSYN: false,
    isRoot: false,
    filter: '',
    defaultFilterGroup: "Fingerprint",
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
    filterField: function () {
        const filter = form.filter.trim();
        if (table.filterId == 0) {
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
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    }
})

const config = reactive({
    ping: false,
    webscan: false,
    crack: false,
})


async function NewScanner() {
    let mode = form.isSYN ? "syn" : "fully connected"
    Callgologger("info", "Port scan task loaded, current mode: " + mode)
    if (!validatePortscan(form.target)) {
        ElMessage.warning("输入目标格式不正确")
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
            ElMessage.warning('可用目标或端口为空')
            return
        }
        ctrl.runningStatus = true
        // 判断是否进行存活探测
        if (global.webscan.default_alive_module != "None") {
            global.webscan.default_alive_module === 'ICMP' ? config.ping = false : config.ping = true
            ips = await HostAlive((ips as any), config.ping)
        }
        if (form.isSYN) {
            await NewSynScanner(specialTarget, ips, portsList)
            return
        }
        await NewTcpScanner(specialTarget, ips, portsList, global.webscan.port_thread, global.webscan.port_timeout)
        // 恢复配置
        ctrl.runningStatus = false
        ips = []
        portsList = []
        Callgologger("info", "Portscan task is ending")
        if (config.webscan) {
            ips = moreOperate.getHTTPLinks(pagination.table.result)
            LinkWebscan(ips)
        }
        if (config.crack) {
            ips = pagination.ctrl.getColumnData("Link")
            LinkCrack(ips)
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


const moreOperate = ({
    getHTTPLinks: function (result: PortScanData[]): string[] {
        return result
            .filter(line => line.Link.startsWith("http"))
            .map(line => line.Link);
    },

    getBruteLinks: function (result: PortScanData[]): string[] {
        return result
            .filter(line => global.dict.options.some(brute => line.Link.startsWith(brute)))
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
            LinkWebscan(targets)
        } else {
            ElNotification.success(`已发送${targets.length}个目标到暴破与未授权检测`)
            LinkCrack(targets)
        }
    },
})

function handleContextMenu(row: any, column: any, e: MouseEvent) {
    //prevent the browser's default menu
    e.preventDefault();
    //show our menu
    ContextMenu.showContextMenu({
        x: e.x,
        y: e.y,
        items: [
            {
                label: "复制全部URL",
                icon: h(CopyDocument, defaultIconSize),
                onClick: () => {
                    moreOperate.CopyURLs('url', pagination.table.result)
                }
            },
            {
                label: "复制全部可爆破协议",
                icon: h(CopyDocument, defaultIconSize),
                onClick: () => {
                    moreOperate.CopyURLs('brute', pagination.table.result)
                }
            },
            {
                label: "复制选中目标",
                divided: true,
                icon: h(CopyDocument, defaultIconSize),
                onClick: () => {
                    moreOperate.CopySelectLinks()
                }
            },
            {
                label: "联动网站扫描",
                icon: h(Promotion, defaultIconSize),
                onClick: () => {
                    moreOperate.Linkage('webscan')
                }
            },
            {
                label: "联动暴破与未授权检测",
                icon: h(Promotion, defaultIconSize),
                onClick: () => {
                    moreOperate.Linkage('crack')
                }
            },
        ]
    });
}

const shodanVisible = ref(false)
const shodanIp = ref('')
const shodanPercentage = ref(0)
const shodanRunningstatus = ref(false)

async function LinkShodan() {
    if (!validateIp(shodanIp.value)) {
        ElMessage.warning("目标输入格式不正确!")
        return
    }
    const lines = SplitTextArea(shodanIp.value)
    for (const line of lines) {
        if (isPrivateIP(line)) {
            ElMessage.warning(line + " 为内网地址不支持扫描!")
            return
        }
    }
    
    let id = 0
    form.target = ""
    shodanRunningstatus.value = true
    let ips = await IPParse(lines)
    for (const ip of ips) {
        if (!shodanRunningstatus.value) {
            return
        }
        let ports = await SpaceGetPort(ip)
        id++
        if (ports == null) {
            continue
        }
        Callgologger("info", "[shodan] " + ip + " port: " + ports.join())
        shodanPercentage.value = Number(((id / ips.length) * 100).toFixed(2));
        if (ports.length > 0) {
            for (const port of ports) {
                form.target += ip + ":" + port.toString() + "\n"
            }
        }
    }
    shodanRunningstatus.value = false
    shodanVisible.value = false
}

function stopShodan() {
    shodanRunningstatus.value = false
    ElNotification.error({
        message: "用户已终止扫描任务",
        position: 'bottom-right',
    });
}
</script>

<template>
    <div style="margin-bottom: 5px;">
        <el-card shadow="never">
            <div class="my-header">
                <div>
                    <el-checkbox v-model="form.isSYN" :disabled="!form.isRoot">
                        <el-tooltip placement="right">
                            <template #content>
                                需要ROOT权限
                            </template>
                            SYN
                        </el-tooltip>
                    </el-checkbox>
                    <el-checkbox v-model="config.crack">口令猜测</el-checkbox>
                    <el-checkbox v-model="config.webscan">网站扫描</el-checkbox>
                </div>
                
                <el-button text bg @click="shodanVisible = true">
                    <template #icon>
                        <img src="/navigation/shodan.png" style="width: 14px; height: 14px;">
                    </template>
                    联动shodan
                </el-button>
                
                <el-button type="primary" @click="NewScanner" v-if="!ctrl.runningStatus">开始扫描</el-button>
                <el-button type="danger" @click="ctrl.stop" v-else>停止扫描</el-button>
            </div>
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
                                域名格式: www.example.com (SYN不支持)
                            </template>
                            <el-icon style="width: 13px;">
                                <QuestionFilled />
                            </el-icon>
                        </el-tooltip>
                    </span>
                    <el-button size="small" @click="uploadFile">IP导入</el-button>
                </div>
                <el-input class="input" type="textarea" :rows="3" v-model="form.target" resize="none" />
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
                <el-input class="input" type="textarea" :rows="3" v-model="form.portlist" resize="none" />
            </el-col>
        </el-row>
    </div>
    <CustomTabs>
        <el-tabs v-model="form.activeName" type="border-card">
            <el-tab-pane label="结果输出" name="1">
                <el-table :data="pagination.table.pageContent" @selection-change="pagination.ctrl.handleSelectChange"
                    @row-contextmenu="handleContextMenu" style="height: 100vh;">
                    <el-table-column type="selection" width="55px" align="center" />
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
                    <el-progress :text-inside="true" :stroke-width="18" :percentage="form.percentage"
                        style="width: 40%;" />
                    <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
                        @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100, 200, 500]"
                        :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                        :total="pagination.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-space :size="2">
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
                <el-tooltip content="导出Excel">
                    <el-button :icon="exportIcon"
                        @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', pagination.table.result)" />
                </el-tooltip>
            </el-space>
        </template>
    </CustomTabs>
    <el-dialog v-model="shodanVisible" width="500">
        <template #header>
            <span style="display: flex; align-items: center">
                <img src="/navigation/shodan.png" style="margin-right: 5px;">
                通过Shodan实现端口开放情况收集
            </span>
        </template>
        <el-text>Tips: 在此处输入要搜索的IP地址，<strong>不支持域名</strong>，可以通过日志处<el-icon><consoleIcon /></el-icon>查看详</el-text>
        <el-text style="margin-bottom: 5px;">细信息</el-text>
        <el-input v-model="shodanIp" type="textarea" :rows="6" placeholder="支持如下输入格式
1.1.1.1
1.1.1.1/24
1.1.1.1-1.1.255.255
1.1.1.1-255"></el-input>
        <template #footer>
            <div class="my-header">
                <el-progress :text-inside="true" :stroke-width="18" :percentage="shodanPercentage"
                    style="width: 300px;" />
                <div>
                    <el-button size="small" type="danger" @click="stopShodan" v-if="shodanRunningstatus">终止探测</el-button>
                    <el-button size="small" type="primary" @click="LinkShodan" v-else>
                        开始收集
                    </el-button>
                </div>
            </div>
        </template>
    </el-dialog>
</template>


<style>
.input {
    height: 120px;
}
</style>
