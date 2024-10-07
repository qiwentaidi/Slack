<script lang="ts" setup>
import { reactive, onMounted, ref, h } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { Copy, SplitTextArea, UploadFileAndRead } from '@/util'
import { ExportToXlsx } from '@/export'
import { QuestionFilled, ChromeFilled, Promotion, CopyDocument, Search, Plus, Upload, CircleClose } from '@element-plus/icons-vue';
import { PortParse, IPParse, NewTcpScanner, HostAlive, IsRoot, NewSynScanner, StopPortScan, Callgologger, SpaceGetPort } from 'wailsjs/go/main/App'
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime'
import global from '@/global'
import { PortScanData } from '@/interface';
import usePagination from '@/usePagination';
import exportIcon from '@/assets/icon/doucment-export.svg'
import { defaultIconSize } from '@/stores/style';
import { LinkCrack, LinkWebscan } from '@/linkage';
import ContextMenu from '@imengyu/vue3-context-menu'
import CustomTabs from '@/components/CustomTabs.vue';
import { isPrivateIP, validateIp, validatePortscan } from '@/stores/validate';
import consoleIcon from '@/assets/icon/console.svg'
import { crackDict, portGroupOptions, portscanOptions } from '@/stores/options';

// syn 扫描模式
onMounted(async () => {
    form.isRoot = await IsRoot()
    updatePorts(1); // 更新初始化显示
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

const form = reactive({
    activeName: '1',
    target: '',
    portlist: '',
    percentage: 0,
    id: 0,
    count: 0,
    radio: "2",
    isRoot: false,
    filter: '',
    defaultFilterGroup: "Fingerprint",
    newscanner: false,
})

const table = reactive({
    result: [] as PortScanData[],
    temp: [] as PortScanData[], // 用于存储过滤之前的数据，后续需要还原给result
    filterId: 0
})

let pagination = usePagination(table.result, 20)

async function uploadFile() {
    form.target = await UploadFileAndRead()
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
    defaultPortscanOption: 1,
    ping: false,
    webscan: false,
    crack: false,
})


async function NewScanner() {
    form.newscanner = false
    Callgologger("info", "Port scanning start, current mode: " + config.defaultPortscanOption)
    if (!validatePortscan(form.target)) {
        ElMessage.warning("输入目标格式不正确")
        return
    }
    const ps = new Scanner()
    await ps.PortScanner()
}

class Scanner {
    public async PortScanner() {
        if (config.defaultPortscanOption == 0 && !form.isRoot) {
            ElMessage.warning("当前为非root用户，无法使用SYN扫描模式")
            return
        }
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
        if (config.defaultPortscanOption == 0) {
            await NewSynScanner(specialTarget, ips, portsList)
            Callgologger("info", "Portscan task is ending")
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


const currentPorts = ref(portGroupOptions[1].text)

function updatePorts(index: number) {
    if (index >= 0 && index < portGroupOptions.length) {
        form.portlist = portGroupOptions[index].value;
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
            .filter(line => crackDict.options.some(brute => line.Link.startsWith(brute)))
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
    <CustomTabs>
        <el-tabs v-model="form.activeName" type="border-card">
            <el-tab-pane label="结果输出" name="1">
                <el-table :data="pagination.table.pageContent" @selection-change="pagination.ctrl.handleSelectChange"
                    @row-contextmenu="handleContextMenu" style="height: calc(100vh - 175px)">
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
        <template #filter>
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
        </template>
        <template #ctrl>
            <el-space :size="3">
                <el-button type="primary" :icon="Plus" @click="form.newscanner = true"
                    v-if="!ctrl.runningStatus">新建任务</el-button>
                <el-button type="danger" :icon="CircleClose" @click="ctrl.stop" v-else>停止扫描</el-button>
                <el-tooltip content="导出Excel">
                    <el-button :icon="exportIcon"
                        @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', pagination.table.result)" />
                </el-tooltip>
            </el-space>
        </template>
    </CustomTabs>
    <el-drawer v-model="form.newscanner" size="40%">
        <template #header>
            <span class="drawer-title">新建扫描任务</span>
            <el-button text bg @click="shodanVisible = true">
                <template #icon>
                    <img src="/navigation/shodan.png" style="width: 14px; height: 14px;">
                </template>
                从Shodan导入
            </el-button>
        </template>
        <el-form :model="config" label-width="auto">
            <el-form-item label="扫描模式:">
                <el-segmented v-model="config.defaultPortscanOption" :options="portscanOptions" style="width: 100%;" />
                <el-tooltip>
                    <template #content>
                        Mac - sudo /Applications/Slack.app/Contents/MacOS/Slack<br />
                        Linux - sudo ./Slack<br />
                        Windows - 右键管理员运行
                    </template>
                    <el-button link size="small" style="margin-top: 5px;">SYN需要以管理员模式运行</el-button>
                </el-tooltip>
            </el-form-item>
            <el-form-item>
                <template #label><span>IP:
                        <el-tooltip>
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
                            <el-icon>
                                <QuestionFilled />
                            </el-icon>
                        </el-tooltip>
                    </span>
                </template>
                <el-input type="textarea" :rows="5" v-model="form.target" />
                <el-button link size="small" :icon="Upload" @click="uploadFile"
                    style="margin-top: 5px;">导入目标文件</el-button>
            </el-form-item>
            <el-form-item label="端口:">
                <el-select v-model="currentPorts" @change="updatePorts">
                    <el-option v-for="(item, index) in portGroupOptions" :label="item.text" :value="index" />
                </el-select>
                <el-input v-model="form.portlist" type="textarea" :rows="4" resize="none"
                    style="margin-top: 5px;"></el-input>
            </el-form-item>
            <el-form-item label="扫描线程:">
                <el-input-number controls-position="right" v-model="global.webscan.port_thread" :min="1" :max="10000" />
            </el-form-item>
            <el-form-item label="联动:">
                <el-checkbox v-model="config.crack">
                    口令猜测
                    <el-tooltip content="开启后会将ftp协议等应用支持暴破资产进行口令猜测">
                        <el-icon>
                            <QuestionFilled />
                        </el-icon>
                    </el-tooltip>
                </el-checkbox>
                <el-checkbox v-model="config.webscan">
                    网站扫描
                    <el-tooltip content="开启后会将http|https协议资产进行网站指纹漏洞扫描">
                        <el-icon>
                            <QuestionFilled />
                        </el-icon>
                    </el-tooltip>
                </el-checkbox>
            </el-form-item>
        </el-form>
        <div class="position-center">
            <el-button type="primary" @click="NewScanner" style="bottom: 10px; position: absolute;">开始任务</el-button>
        </div>
    </el-drawer>
    <el-dialog v-model="shodanVisible" width="500">
        <template #header>
            <span style="display: flex; align-items: center">
                <img src="/navigation/shodan.png" style="margin-right: 5px;">
                从Shodan拉取资产端口开放情况
            </span>
        </template>
        <el-text style="margin-bottom: 10px;"><strong>Tips: 不支持域名输入</strong>，可以通过日志处<el-icon>
                <consoleIcon />
            </el-icon>查看详细信息</el-text>
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
                    <el-button size="small" type="danger" @click="stopShodan"
                        v-if="shodanRunningstatus">终止探测</el-button>
                    <el-button size="small" type="primary" @click="LinkShodan" v-else>
                        开始收集
                    </el-button>
                </div>
            </div>
        </template>
    </el-dialog>
</template>


<style scoped>

</style>
