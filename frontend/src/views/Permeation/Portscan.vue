<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue';
import { ElMessage, ElNotification } from 'element-plus'
import { CopyURLs, SplitTextArea, ReadLine, currentTime } from '../../util'
import { ExportToXlsx } from '../../export'
import async from 'async';
import { QuestionFilled, CopyDocument, ChromeFilled } from '@element-plus/icons-vue';
import { PortParse, IPParse, PortCheck, HostAlive, PortBrute, IsRoot, SynPortCheck } from '../../../wailsjs/go/main/App'
import { Mkdir, WriteFile, CheckFileStat, GetFileContent, UserHomeDir } from '../../../wailsjs/go/main/File'
import { BrowserOpenURL, EventsOn, EventsOff } from '../../../wailsjs/runtime'
import global from '../../global'
interface PortScanData {
    Status: boolean
    IP: string
    Port: number
    Server: string
    Link: string
    HttpTitle: string
}
onMounted(() => {
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
    EventsOn("synProgressID", (id: number) => {
        form.percentage = Number(((id / form.count) * 100).toFixed(2));
    });
    EventsOn("synScanComplete", (log: any) => {
        form.percentage = 100
        ctrl.runningStatus = false
    });
    return () => {
        EventsOff("synPortScanLoading");
        EventsOff("synScanComplete");
    };
});
onMounted(async () => {
    if (!await IsRoot()) {
        ElNotification({
            title: "提示",
            message: "当前应用非ROOT模式启动，无法调用SYN扫描模式",
            type: "warning",
        });
    }
    global.PATH.PortBurstPath = await UserHomeDir() + global.PATH.PortBurstPath
    if (!(await CheckFileStat(global.PATH.PortBurstPath))) {
        InitBurteDict()
    }
    table.pageContent = []
    table.burteResult = []
    updatePorts(); // 更新初始化显示
});
const form = reactive({
    activeName: '1',
    target: '',
    portlist: '',
    percentage: 0,
    id: 0,
    count: 0,
    aliveOptions: ["None", "ICMP", "Ping"],
    currentAliveMoudle: "None",
    scanModuleOptions: ["CONNECT", "SYN"],
    currentScanMoudle: "CONNECT",
})



// 初始化字典
async function InitBurteDict() {
    if (await Mkdir(global.PATH.PortBurstPath)) {
        Mkdir(global.PATH.PortBurstPath + "/username")
        Mkdir(global.PATH.PortBurstPath + "/password")
        for (const item of global.dict.usernames) {
            WriteFile("txt", `${global.PATH.PortBurstPath}/username/${item.name}.txt`, item.dic.join("\n"))
        }
        WriteFile("txt", `${global.PATH.PortBurstPath}/password/password.txt`, global.dict.passwords.join("\n"))
    }
}

async function ReadDict(path: string) {
    ctrl.currentDic = await GetFileContent(global.PATH.PortBurstPath + path)
}

async function SaveFile(path: string) {
    let r = await WriteFile('txt', global.PATH.PortBurstPath + path, ctrl.currentDic)
    if (r) {
        ElMessage({
            showClose: true,
            message: 'success',
            type: 'success',
        })
    } else {
        ElMessage({
            showClose: true,
            message: 'failed',
            type: 'error',
        })
    }
}

const config = reactive({
    thread: 1000,
    timeout: 7,
    ping: false,
    burte: false,
})
const radio = ref('2')

async function NewScanner() {
    global.Logger.value += `[INF]${currentTime()} Portscan ${form.currentScanMoudle} moudle is loading\n`
    const ps = new Scanner()
    await ps.PortScanner()
}

class Scanner {
    public init() {
        table.result = []
        table.pageContent = []
        form.id = 0
        form.count = 0
    }
    public async PortScanner() {
        this.init()
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
        portsList = await PortParse(form.portlist)
        ips = await IPParse(conventionTarget)
        if (ips == null) {
            form.count = specialTarget.length
        } else {
            form.count = ips.length * portsList.length + specialTarget.length
        }
        if (form.count == 0) {
            ElMessage({
                showClose: true,
                message: '可用目标或端口为空',
                type: 'warning',
            })
            return
        }
        ctrl.runningStatus = true
        if (form.currentAliveMoudle != "None") {
            form.currentAliveMoudle === 'ICMP' ? config.ping = false : config.ping = true
            ips = await HostAlive((ips as any), config.ping)
        }
        if (specialTarget.length != 0) {
            global.Logger.value += "[INFO] 存在 192.168.1.1:6379 此类特殊目标优先执行扫描\n"
            async.eachSeries(specialTarget, (ipport: string, callback: () => void) => {
                let temp = ipport.split(":")
                PortCheck(temp[0], Number(temp[1]), config.timeout).then(result => {
                    if (result.Status) {
                        table.result.push({
                            host: temp[0],
                            port: temp[1],
                            fingerprint: result.Server,
                            link: result.Link,
                            title: result.HttpTitle,
                        })
                        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
                    }
                    form.id++;
                    form.percentage = Number(((form.id / form.count) * 100).toFixed(2));
                    callback();
                })
            }, (err: any) => {
                if (err) {
                    ElMessage.error(err)
                } else {
                    global.Logger.value += "[INFO] 特殊目标扫描已完成\n"
                }
            })
        }
        if (form.currentScanMoudle == "SYN") {
            await SynPortCheck(ips, portsList, 10)
            return
        }
        async.eachSeries(ips, (ip: string, callback: () => void) => {
            global.Logger.value += "[INFO] Portscan " + ip + "，port count " + portsList.length + "\n"
            async.eachLimit(portsList, config.thread, (port: number, portCallback: () => void) => {
                if (!ctrl.runningStatus) {
                    return
                }
                PortCheck(ip as string, port, config.timeout).then(result => {
                    if (result.Status) {
                        global.Logger.value += `[+] ${ip}:${port} is open!\n`
                        table.result.push({
                            host: ip,
                            port: port,
                            fingerprint: result.Server,
                            link: result.Link,
                            title: result.HttpTitle,
                        })
                        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
                    }
                    form.id++;
                    form.percentage = Number(((form.id / form.count) * 100).toFixed(2));
                    portCallback();
                });
            }, (err: any) => {
                callback();
            });
        }, (err: any) => {
            if (config.burte) {
                this.Burte()
            }
            ctrl.runningStatus = false
            ips = []
            portsList = []
            global.Logger.value += `${currentTime()} Task is ending\n`
        });
    }
    public async Burte() {
        table.burteResult = []
        global.dict.passwords = (await ReadLine(global.PATH.PortBurstPath + "/password/password.txt"))!
        for (var item of global.dict.usernames) {
            item.dic = (await ReadLine(global.PATH.PortBurstPath + "/username/" + item.name + ".txt"))!
        }
        async.eachLimit(getColumnData("link"), 20, (target: string, callback: () => void) => {
            if (!ctrl.runningStatus) {
                return
            }
            global.dict.usernames.forEach(ele => {
                let protocol = target.split("://")[0]
                if (ele.name.toLowerCase() === protocol) {
                    global.Logger.value += "[INF] Start burte " + target + "\n"
                    PortBrute(target, ele.dic, global.dict.passwords).then((result) => {
                        if (result !== undefined && result.Status !== false) {
                            global.Logger.value += `[++] ${target} ${result.Username}:${result.Password}\n`
                            table.burteResult.push({
                                host: result.Host,
                                protocol: result.Protocol,
                                username: result.Username,
                                password: result.Password,
                                time: currentTime(),
                            })
                        }
                        callback();
                    });
                }
            });
        }, (err: any) => {
            if (err) {
                ElMessage.error(err)
            } else {
                global.Logger.value += "[END] All target port brute have been completed\n"
            }
            ctrl.runningStatus = false
        });
    }
}

function getColumnData(prop: string): any[] {
    return table.result.map((item: any) => item[prop]);
}

const ctrl = reactive({
    drawer: false,
    innerDrawer: false,
    runningStatus: false,
    stop: function () {
        if (ctrl.runningStatus) {
            ctrl.runningStatus = false
            ElMessage({
                showClose: true,
                message: '已停止任务',
                type: 'warning',
            })
        }
    },
    currentDic: '',
    currentPath: ''
})

const table = reactive({
    currentPage: 1,
    pageSize: 10,
    result: [{}],
    pageContent: [{}],
    burteResult: [{}],
})

function updatePorts() {
    const index = Number(radio.value) - 1;
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
</script>


<template>
    <el-form :model="form" label-width="35%" style="width: 75%;">
        <el-form-item>
            <template #label>IP:
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
                        ...<br />
                        <br />
                        如果端口遗漏多请在配置中调高端口超时时间
                    </template>
                    <el-icon>
                        <QuestionFilled size="24" />
                    </el-icon>
                </el-tooltip>
            </template>
            <el-input type="textarea" rows="3" v-model="form.target" resize="none" />
        </el-form-item>
        <el-form-item label="预设端口:">
            <el-radio-group v-model="radio" @change="updatePorts">
                <el-radio value="1">数据库端口</el-radio>
                <el-radio value="2">企业端口</el-radio>
                <el-radio value="3">高危端口</el-radio>
                <el-radio value="4">全端口</el-radio>
                <el-radio value="5">三高一弱</el-radio>
            </el-radio-group>
        </el-form-item>
        <el-form-item label="端口列表:">
            <el-input type="textarea" rows="3" v-model="form.portlist" resize="none" />
        </el-form-item>
        <el-form-item>
            <el-space>
                <el-button type="primary" @click="NewScanner" v-if="!ctrl.runningStatus" size="small">开始扫描</el-button>
                <el-button type="danger" @click="ctrl.stop" v-else size="small">停止扫描</el-button>
                <el-tag type="info">线程:{{ config.thread }}</el-tag>
                <el-link type="primary" @click="ctrl.drawer = true">更多参数</el-link>
            </el-space>
        </el-form-item>
    </el-form>
    <el-drawer v-model="ctrl.drawer" size="40%">
        <template #header>
            <h4>设置高级参数</h4>
        </template>
        <el-form label-width="auto">
            <el-form-item label="扫描模式:">
                <el-segmented v-model="form.currentScanMoudle" :options="form.scanModuleOptions" block
                    style="width: 100%;" />
            </el-form-item>
            <el-form-item label="存活探测:">
                <el-segmented v-model="form.currentAliveMoudle" :options="form.aliveOptions" block
                    style="width: 100%;" />
            </el-form-item>
            <el-form-item label="线程数量:">
                <el-input-number controls-position="right" v-model="config.thread" :min="1" :max="30000" />
            </el-form-item>
            <el-form-item label="指纹超时时长(s):">
                <el-input-number controls-position="right" v-model="config.timeout" :min="1" />
            </el-form-item>
            <el-form-item label="口令猜测:">
                <el-switch v-model="config.burte" />
            </el-form-item>
        </el-form>
        <div v-if="config.burte === true">
            <el-descriptions title="字典管理" :column="2" border>
                <el-descriptions-item v-for="item in global.dict.usernames" :label="item.name" :span="2">
                    <el-button
                        @click="ctrl.innerDrawer = true; ctrl.currentPath = '/username/' + item.name + '.txt'; ReadDict('/username/' + item.name + '.txt')">
                        用户名
                    </el-button>
                </el-descriptions-item>
                <el-descriptions-item label="密码(所有协议通用)" :span="2">
                    <el-button
                        @click="ctrl.innerDrawer = true; ctrl.currentPath = '/password/password.txt'; ReadDict('/password/password.txt')">
                        密码
                    </el-button>
                </el-descriptions-item>
            </el-descriptions>
        </div>
        <el-drawer v-model="ctrl.innerDrawer" title="字典管理" :append-to-body="true">
            <el-input type="textarea" rows="20" v-model="ctrl.currentDic"></el-input>
            <el-button type="primary" style="margin-top: 10px; float: right;"
                @click="SaveFile(ctrl.currentPath)">保存</el-button>
        </el-drawer>
    </el-drawer>

    <div style="position: relative;">
        <el-tabs v-model="form.activeName" type="card">
            <el-tab-pane label="端口扫描控制台" name="1">
                <el-table :data="table.pageContent" border style="height: 42vh;">
                    <el-table-column type="selection" width="55px" />
                    <el-table-column prop="host" label="主机" />
                    <el-table-column prop="port" label="端口" width="100px" />
                    <el-table-column prop="fingerprint" label="指纹" />
                    <el-table-column prop="link" label="目标">
                        <template #default="scope">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.link)"
                                v-show="scope.row.link != ''">
                            </el-button>
                            {{ scope.row.link }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="title" label="网站标题" />
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <el-progress :text-inside="true" :stroke-width="20" :percentage="form.percentage" color="#5DC4F7"
                        style="width: 50%;" />
                    <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                        :current-page="table.currentPage" :page-sizes="[10, 20, 50]" :page-size="table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
            <el-tab-pane label="脆弱账号" name="2">
                <el-table :data="table.burteResult" border style="height: 45vh;">
                    <el-table-column type="index" width="60px" label="#" />
                    <el-table-column prop="host" label="主机" />
                    <el-table-column prop="protocol" label="协议" />
                    <el-table-column prop="username" label="用户名" />
                    <el-table-column prop="password" label="密码" />
                    <el-table-column prop="time" label="扫描时间" />
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-button-group>
                <el-button :icon="CopyDocument" @click="CopyURLs(table.result)">复制全部URL</el-button>
                <el-button @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', table.result)">
                    <template #icon>
                        <img src="/excle.svg" width="16">
                    </template>
                    导出Excle</el-button>
            </el-button-group>
        </div>
    </div>
</template>


<style>
.el-drawer__header {
    margin-bottom: 0px;
}
</style>
