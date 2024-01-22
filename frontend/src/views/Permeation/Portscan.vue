<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { ElMessage } from 'element-plus'
import { ExportToXlsx, CopyURLs, SplitTextArea } from '../../util'
import async from 'async';
import { QuestionFilled, ArrowDown, Promotion, ChromeFilled, Tickets } from '@element-plus/icons-vue';
import {
    PortParse,
    IPParse,
    PortCheck,
    HostAlive,
    PortBrute,
} from '../../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    table.pageContent = []
    table.burteResult = []
    updatePorts();
});
const form = reactive({
    activeName: '1',
    target: '',
    ips: [{}],
    portlist: '',
    portsList: [{}],
    percentage: 0,
    id: 0,
    log: '',
})

const config = reactive({
    options: [{
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
    ],
    thread: 1000,
    timeout: 12,
    changeIp: false,
    alive: false,
    ping: false,
    burte: false,
})
const radio = ref('2')

function NewScanner() {
    const ps = new Scanner()
    ps.PortScanner()
}

class Scanner {
    public init() {
        table.result = []
        table.pageContent = []
        ctrl.exit = false
        form.id = 0
    }
    public async PortScanner() {
        this.init()
        if (config.changeIp === false) {
            form.portsList = await PortParse(form.portlist)
            form.ips = await IPParse(form.target)
            if (form.portlist.length == 0 || form.ips.length == 0) {
                ElMessage({
                    showClose: true,
                    message: '目标或端口为空',
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
            const count = (form.ips.length * form.portsList.length)
            async.eachSeries(form.ips, (ip: string, callback: () => void) => {
                form.log += "[INFO] Portscan " + ip + "，port count " + form.portsList.length + "\n"
                async.eachLimit(form.portsList, config.thread, (port: number, callback: () => void) => {
                    if (ctrl.exit === true) {
                        return
                    }
                    PortCheck(ip as string, port, config.timeout).then(result => {
                        if (result.Status) {
                            form.log += `[+] Portscan ${ip}:${port} is open!\n`
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
                        callback();
                    });
                }, (err: any) => {
                    if (config.burte === true) {
                        this.burteport()
                    }
                    if (err) {
                        ElMessage.error(err)
                    } else {
                        ctrl.exit = false
                        form.log += "[END] 端口扫描任务已完成\n"
                    }
                    callback();
                });
            });
        } else {
            const lines = SplitTextArea(form.target)
            for (const line of lines) {
                if (ctrl.exit === true) {
                    return
                }
                form.id++;
                form.percentage = Number(((form.id / lines.length) * 100).toFixed(2));
                let temp = line.split(":")
                PortCheck(temp[0], Number(temp[1]), config.timeout).then((result) => {
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
                });
            }
        }
    }
    public burteport() {
        table.burteResult = []
        var date = new Date();
        async.eachLimit(getColumnData("link"), 20, (target: string, callback: () => void) => {
            if (ctrl.exit === true) {
                return
            }
            ctrl.dict.forEach(ele => {
                let pro = target.split("://")[0]
                if (ele.name.toLowerCase() === pro) {
                    form.log += "[INF] Start burte " + target + "\n"
                    PortBrute(target, ele.dic, ctrl.passwords).then((result) => {
                        if (result !== undefined && result.Status !== false) {
                            form.log += `[++] ${target} ${result.Username}:${result.Password}\n`
                            table.burteResult.push({
                                host: result.Host,
                                protocol: result.Protocol,
                                username: result.Username,
                                password: result.Password,
                                time: date.toLocaleString(),
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
                form.log += "[END] All target port brute have been completed\n"
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
    exit: false,
    stop: function () {
        if (ctrl.exit === false) {
            ctrl.exit = true
            ElMessage({
                showClose: true,
                message: '已停止任务',
                type: 'warning',
            })
            form.log += "[STOP] 已停止扫描任务\n"
        }
    },
    dict: [
        {
            name: "FTP",
            dic: ["ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"],
        },
        {
            name: "SSH",
            dic: ["root", "admin"]
        },
        {
            name: "Telnet",
            dic: ["root", "admin"]
        },
        {
            name: "SMB",
            dic: ["administrator", "admin", "guest"]
        },
        {
            name: "Mssql",
            dic: ["sa", "sql"]
        },
        {
            name: "Oracle",
            dic: ["sys", "system", "admin", "test", "web", "orcl"]
        },
        {
            name: "Mysql",
            dic: ["root", "mysql"]
        },
        {
            name: "RDP",
            dic: ["administrator", "admin", "guest"]
        },
        {
            name: "Postgresql",
            dic: ["postgres", "admin"]
        },
        {
            name: "VNC",
            dic: ["admin", "administrator", "root"]
        },
    ],
    passwords: ["123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "1234567", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789"],
    currentDic: '',
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
    if (index >= 0 && index < config.options.length) {
        form.portlist = config.options[index].value;
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
    <el-tabs v-model="form.activeName" type="card">
        <el-tab-pane label="端口扫描控制台" name="1">
            <el-form :model="form" label-width="20%">
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
                                192.168.1.1-255<br />
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
                    <el-input type="textarea" rows="3" v-model="form.target" style="width: 70%;" />
                    <div>
                        <el-button type="primary" class="left" @click="NewScanner">开始扫描</el-button>
                        <el-button type="primary" class="left" @click="ctrl.stop">停止</el-button>
                    </div>
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
                    <el-input type="textarea" rows="3" v-model="form.portlist" style="width: 70%;" />
                </el-form-item>
                <el-form-item>
                    <template #label>
                        功能:
                    </template>
                    <el-space>
                        <el-button color="#DEACBD" @click="CopyURLs">复制全部URL目标</el-button>
                        <el-button color="#7FEDF9"
                            @click="ExportToXlsx(['主机', '端口', '指纹', '目标', '网站标题'], '端口扫描', 'portscan', table.result)">导出Excle</el-button>

                        <el-dropdown>
                            <el-button :icon="Promotion" color="#A29EDE">
                                联动<el-icon class="el-icon--right"><arrow-down /></el-icon>
                            </el-button>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item>发送到网站扫描(没做)</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                        <el-tag type="success">端口超时: {{ config.timeout }}s</el-tag>
                        <el-link type="primary" @click="ctrl.drawer = true">更多参数</el-link>
                        <el-checkbox v-model="config.changeIp" label="更改IP输入模式" />
                        <el-tooltip placement="right">
                            <template #content>更改IP输入模式为192.168.0.1:6379此类形式<br />仅支持换行分割，开启后端口列表无效且扫描为单线程模式</template>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </el-tooltip>
                    </el-space>
                    <el-drawer v-model="ctrl.drawer" title="设置高级参数">
                        <el-form label-width="100px">
                            <el-form-item label="存活探测:" class="bottom">
                                <el-checkbox v-model="config.alive" label="开启" />
                                <el-switch v-model="config.ping" active-text="ICMP" inactive-text="Ping" v-if="config.alive"
                                    style="margin-left: 10px;" />
                            </el-form-item>
                            <el-form-item label="线程数量:" class="bottom">
                                <el-input-number controls-position="right" v-model="config.thread" :min="1" :max="30000" />
                            </el-form-item>
                            <el-form-item label="超时时长(s):" class="bottom">
                                <el-input-number controls-position="right" v-model="config.timeout" :min="1" />
                            </el-form-item>
                            <el-form-item label="口令猜测:" class="bottom">
                                <el-switch v-model="config.burte" />
                            </el-form-item>
                            <div v-if="config.burte === true">
                                <el-descriptions title="用户名字典管理" :column="2" border>
                                    <el-descriptions-item v-for="item in ctrl.dict" :label="item.name" :span="2">
                                        <el-button @click="ctrl.innerDrawer = true; ctrl.currentDic = item.dic.join('\n')">
                                            口令字典
                                        </el-button>
                                    </el-descriptions-item>
                                </el-descriptions>
                                <el-descriptions title="密码管理(全协议通用)" :column="2" border>
                                    <el-descriptions-item label="Password" :span="2">
                                        <el-button
                                            @click="ctrl.innerDrawer = true; ctrl.currentDic = ctrl.passwords.join('\n')">
                                            口令字典
                                        </el-button>
                                    </el-descriptions-item>
                                </el-descriptions>
                            </div>
                            <el-drawer v-model="ctrl.innerDrawer" title="账号管理" :append-to-body="true">
                                <el-input type="textarea" rows="20" v-model="ctrl.currentDic"></el-input>
                                <!-- <el-button type="primary" style="margin-top: 10px;">保存</el-button> -->
                            </el-drawer>
                        </el-form>
                    </el-drawer>
                </el-form-item>
            </el-form>
            <el-table :data="table.pageContent" border style="height: 45vh;">
                <el-table-column type="selection" width="55px" />
                <el-table-column prop="host" label="主机" />
                <el-table-column prop="port" label="端口" width="100px" />
                <el-table-column prop="fingerprint" label="指纹" />
                <el-table-column prop="link" label="目标" />
                <el-table-column prop="title" label="网站标题" />
                <el-table-column fixed="right" label="操作" width="55px">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="left">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.link)">
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
                :current-page="table.currentPage" :page-sizes="[10, 20, 50]" :page-size="table.pageSize"
                layout="total, sizes, prev, pager, next" :total="table.result.length"
                style="margin-top: 5px; float: right;">
            </el-pagination>
            <el-progress :text-inside="true" :stroke-width="18" :percentage="form.percentage"
                style="margin-top: 11px; width: 50%;" />
        </el-tab-pane>
        <el-tab-pane label="脆弱账号" name="2">
            <el-table :data="table.burteResult" border style="width: 100%; height: 80vh;">
                <el-table-column type="index" width="60px" label="#" />
                <el-table-column prop="host" label="主机" />
                <el-table-column prop="protocol" label="协议" />
                <el-table-column prop="username" label="用户名" />
                <el-table-column prop="password" label="密码" />
                <el-table-column prop="time" label="扫描时间" />
            </el-table>
        </el-tab-pane>
        <el-tab-pane name="3">
            <template #label>
                <el-icon>
                    <Tickets />
                </el-icon>
                <span>运行日志</span>
            </template>
            <el-input class="log-textarea" v-model="form.log" type="textarea" rows="20" readonly></el-input>
        </el-tab-pane>
        <el-tab-pane label="ceshisssssssssssssssssssssssssss" name="4">

        </el-tab-pane>
    </el-tabs>
</template>


<style>
.left {
    margin-left: 10px;
}

.bottom {
    margin-bottom: 20px;
}

.el-tabs__item {
    position: relative;
    display: inline-block;
    max-width: 130px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    .el-tooltip {
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .el-icon-close {
      position: absolute !important;
      top: 13px !important;
      right: 3px !important;
    }
  }
</style>
