<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue';
import { ElMessage } from 'element-plus'
import { ExportToXlsx, CopyURLs, SplitTextArea, ReadLine, currentTime } from '../../util'
import async from 'async';
import { QuestionFilled, DocumentChecked, DocumentCopy, ChromeFilled } from '@element-plus/icons-vue';
import { PortParse, IPParse, PortCheck, HostAlive, PortBrute } from '../../../wailsjs/go/main/App'
import { Mkdir, WriteFile, CheckFileStat, GetFileContent } from '../../../wailsjs/go/main/File'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
import global from '../../global'
onMounted(async () => {
    if (!(await CheckFileStat(PortBurstPath))) {
        InitBurteDict()
    }
    table.pageContent = []
    table.burteResult = []
    updatePorts(); // 更新初始化显示
});
const form = reactive({
    activeName: '1',
    target: '',
    ips: [{}],
    portlist: '',
    portsList: [{}],
    percentage: 0,
    id: 0,
})

const options = [{
    text: "数据库端口",
    value: "1433,1521,3306,5432,6379,9200,11211,27017",
},
{
    text: "企业端口",
    value: "21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880",
},
{
    text: "高危端口",
    value: "21,22,23,53,80,443,8080,8000,139,445,3389,1521,3306,6379,7001,2375,27017,11211",
},
{
    text: "全端口",
    value: "1-65535"
},
{
    text: "自定义",
    value: ""
},
]

// 初始化字典
async function InitBurteDict() {
    if (await Mkdir(PortBurstPath)) {
        Mkdir(PortBurstPath + "/username")
        Mkdir(PortBurstPath + "/password")
        for (const item of global.dict.usernames) {
            WriteFile("txt", `${PortBurstPath}/username/${item.name}.txt`, item.dic.join("\n"))
        }
        WriteFile("txt", `${PortBurstPath}/password/password.txt`, global.dict.passwords.join("\n"))
    }
}

async function ReadDict(path: string) {
    ctrl.currentDic = await GetFileContent(PortBurstPath + path)
}

async function SaveFile(path: string) {
    let r = await WriteFile('txt', PortBurstPath + path, ctrl.currentDic)
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
    alive: false,
    ping: false,
    burte: false,
})
const radio = ref('2')

async function NewScanner() {
    global.Logger.value +=  `${currentTime} Portscan task is loading\n`
    const ps = new Scanner()
    await ps.PortScanner()
    if (config.burte) {
        ps.Burte()
    }
}

class Scanner {
    public init() {
        table.result = []
        table.pageContent = []
        ctrl.runningStatus = true
        form.id = 0
    }
    public async PortScanner() {
        this.init()
        var count = 0
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
        form.portsList = await PortParse(form.portlist)
        form.ips = await IPParse(conventionTarget)
        if (form.ips == null) {
            count = specialTarget.length
        } else {
            count = form.ips.length * form.portsList.length + specialTarget.length
        }
        if (count == 0) {
            ElMessage({
                showClose: true,
                message: '可用目标或端口为空',
                type: 'warning',
            })
            return
        }
        if (config.alive === true) {
            HostAlive((form.ips as any), config.alive).then(
                result => {
                    form.ips = result
                }
            )
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
                    form.percentage = Number(((form.id / count) * 100).toFixed(2));
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
        async.eachSeries(form.ips, (ip: string, callback: () => void) => {
            global.Logger.value += "[INFO] Portscan " + ip + "，port count " + form.portsList.length + "\n"
            async.eachLimit(form.portsList, config.thread, (port: number, portCallback: () => void) => {
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
                    form.percentage = Number(((form.id / count) * 100).toFixed(2));
                    portCallback();
                });
            }, (err: any) => {
                callback();
            });
        }, (err: any) => {
            ctrl.runningStatus = false
            global.Logger.value += `${currentTime} Task is ending\n`
        });
    }
    public async Burte() {
        table.burteResult = []
        global.dict.passwords = (await ReadLine(PortBurstPath + "/password/password.txt"))!
        for (var item of global.dict.usernames) {
            item.dic = (await ReadLine(PortBurstPath + "/username/" + item.name + ".txt"))!
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
                                time: currentTime,
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
    if (index >= 0 && index < options.length) {
        form.portlist = options[index].value;
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
                <el-radio label=1>数据库端口</el-radio>
                <el-radio label=2>企业端口</el-radio>
                <el-radio label=3>高危端口</el-radio>
                <el-radio label=4>全端口</el-radio>
                <el-radio label=5>自定义</el-radio>
            </el-radio-group>
        </el-form-item>
        <el-form-item label="端口列表:">
            <el-input type="textarea" rows="3" v-model="form.portlist" resize="none" />
        </el-form-item>
        <el-form-item label="功能:">
            <el-space>
                <el-button type="primary" @click="NewScanner" v-if="!ctrl.runningStatus">开始扫描</el-button>
                <el-button type="danger" @click="ctrl.stop" v-else>停止扫描</el-button>
                <el-tag type="info">线程:{{ config.thread }}</el-tag>
                <el-link type="primary" @click="ctrl.drawer = true">更多参数</el-link>
            </el-space>
        </el-form-item>
    </el-form>
    <el-drawer v-model="ctrl.drawer" size="30%">
        <template #title>
            <h4>设置高级参数</h4>
        </template>
        <el-form label-width="auto">
            <el-form-item label="存活探测:">
                <el-checkbox v-model="config.alive" />
                <el-switch v-model="config.ping" active-text="ICMP" inactive-text="Ping" v-if="config.alive"
                    style="margin-left: 10px;" />
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
                <el-table :data="table.pageContent" border style="height: 45vh;">
                    <el-table-column type="selection" width="55px" />
                    <el-table-column prop="host" label="主机" />
                    <el-table-column prop="port" label="端口" width="100px" />
                    <el-table-column prop="fingerprint" label="指纹" />
                    <el-table-column prop="link" label="目标">
                        <template #default="scope">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.link)">
                            </el-button>
                            {{ scope.row.link }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="title" label="网站标题" />
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <el-progress :text-inside="true" :stroke-width="20" :percentage="form.percentage" color="#5DC4F7"
                        style="width: 70%;" />
                    <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
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
        <div class="custom_eltabs_titlebar" v-if="form.activeName == '1'">
            <el-button-group>
                <el-button text :icon="DocumentCopy" @click="CopyURLs(table.result)">复制全部URL</el-button>
                <el-button text :icon="DocumentChecked"
                    @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', table.result)">导出</el-button>
            </el-button-group>
        </div>
    </div>
</template>


<style>
.el-drawer__header {
    margin-bottom: 0px;
}
</style>
