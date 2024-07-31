<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { InfoFilled, Search, Document, Link } from '@element-plus/icons-vue';
import { FileDialog, CheckFileStat, UserHomeDir } from '../../../wailsjs/go/main/File';
import { ElMessage } from 'element-plus';
import async from 'async';
import global from '../../global';
import { ReadLine } from '../../util';
import { PortBrute, Callgologger } from '../../../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';

interface BruteResult {
    Host: string
    Port: string
    Protocol: string
    Username: string
    Password: string
}

onMounted(() => {
    EventsOn("bruteResult", (result: BruteResult) => {
        let temp = result.Host.split(":")
        config.content.push({
            Host: temp[0],
            Port: temp[1],
            Protocol: result.Protocol,
            Username: result.Username,
            Password: result.Password,
        })
    });
    return () => {
        EventsOff("bruteResult");
    };
});

const options = ["ftp", "ssh", "telnet", "smb", "oracle", "mssql", "mysql", "rdp", "postgresql", "vnc", "redis", "memcached", "mongodb"]

const config = reactive({
    builtIn: true,
    association: false,
    username: '',
    password: '',
    target: '',
    content: [] as BruteResult[],
    defaultOption: "ftp",
    input: '',
    runningStatus: false,
})

const addTarget = () => {
    if (config.input == "") {
        ElMessage.warning("请输入目标")
        return
    }
    config.target += config.defaultOption + "://" + config.input
}

const placeholder = `
1、使用内置字典时，账号密码输入框无效，要使用自己的字典请取消勾选使用内置字典

2、账号密码框处支持单字段比如admin和txt字典路径，可以拖拽读取文件路径

3、联想模式意为根据关键字段如出生年月或公司名等固定规律组成字典，需要在联想
生成器模块生成结果才可以使用

4、支持协议名: ftp、ssh、telnet、smb、oracle、mssql、mysql、rdp、
postgresql、vnc、redis、memcached、mongodb
`

const placeholder2 = `e.g
redis://10.0.0.1:6379
mysql://10.0.0.1:3306
`

async function selectDict(input: string) {
    let filepath = await FileDialog("*.txt")
    if (input == "username") {
        config.username = filepath
        return
    }
    config.password = filepath
}

async function NewScanner() {
    let url = [] as string[]
    let passDict = [] as string[]
    let userDict = [] as string[]
    for (const item of config.target.split("\n")) {
        if (item.includes("://")) {
            url.push(item)
        }
    }
    if (url.length == 0) {
        ElMessage({
            message: "可用目标为空",
            type: "warning"
        })
        return
    }
    // 非内置模式处理字典
    if (!config.builtIn) {
        if (config.username != "") {
            let result = await CheckFileStat(config.username) // 文件不存在则为单用户名
            if (!result) {
                userDict.push(config.username)
            } else {
                userDict = (await ReadLine(config.username))!
            }
        }
        if (config.password != "") {
            let result = await CheckFileStat(config.password)
            if (!result) {
                passDict.push(config.password)
            } else {
                passDict = (await ReadLine(config.password))!
            }
        }
    }
    // 仅勾选联想模式
    if (config.association) {
        if (global.temp.thinkdict.length == 0) {
            ElMessage({
                message: "未生成联想字典",
                type: "warning"
            })
        } else {
            passDict.push(...global.temp.thinkdict)
        }
    }
    config.runningStatus = true
    config.content = []
    let home = await UserHomeDir()
    global.dict.passwords = (await ReadLine(home + global.PATH.PortBurstPath + "/password/password.txt"))!
    for (var item of global.dict.usernames) {
        item.dic = (await ReadLine(home + global.PATH.PortBurstPath + "/username/" + item.name + ".txt"))!
    }
    var id = 0

    async.eachLimit(url, 20, async (target: string, callback: () => void) => {
        let protocol = target.split("://")[0]
        // 如果原字典不为空则采用内置字典
        if (userDict.length == 0) {
            userDict = global.dict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
        }
        if (passDict.length == 0) {
            passDict = global.dict.passwords
        } else {
            passDict.push("") // 增加空密码
        }
        Callgologger("info", target + " is start brute")
        await PortBrute(target, userDict, passDict)
        id++
        if (id == url.length) {
            callback()
        }
    }, (err: any) => {
        Callgologger("info", "PortBrute Finished")
        config.runningStatus = false
    });
}

</script>

<template>
    <div class="my-header">
        <div style="display: flex;">
            <el-checkbox v-model="config.builtIn" label="内置字典" />
            <el-checkbox v-model="config.association" label="联想模式" />
            <el-button type="primary" link :icon="Link" style="margin-left: 10px;" v-if="config.association"
                @click="$router.push('/Tools/Thinkdict')">跳转到联想字典生成器</el-button>
        </div>
        <el-tooltip :content="placeholder" placement="bottom">
            <template #content>
                1、使用内置字典时，账号密码输入框无效，要使用自己的字典请取消勾选使用内置字典<br />
                2、账号密码框处支持单字段比如admin和txt字典路径，可以拖拽读取文件路径<br />
                3、联想模式意为根据关键字段如出生年月或公司名等固定规律组成字典<br />
                4、支持协议名: ftp、ssh、telnet、smb、oracle、mssql、mysql、rdp、vnc、redis<br />
                postgresql、memcached、mongodb<br />
            </template>
            <el-button :icon="InfoFilled" link>使用须知</el-button>
        </el-tooltip>
        <el-button type="primary" :icon="Search" @click="NewScanner">开始暴破</el-button>
    </div>
    <div style="margin-top: 10px; margin-bottom: 10px;">
        <el-input v-model="config.username" :disabled="config.builtIn" style="margin-bottom: 10px;">
            <template #prepend>
                账号
            </template>
            <template #suffix>
                <el-button link :icon="Document" @click="selectDict('username')" :disabled="config.builtIn"></el-button>
            </template>
        </el-input>
        <el-input v-model="config.password" :disabled="config.builtIn">
            <template #prepend>
                密码
            </template>
            <template #suffix>
                <el-button link :icon="Document" @click="selectDict" :disabled="config.builtIn"></el-button>
            </template>
        </el-input>
    </div>
    <div style="display: flex ;height: 75vh;">
        <div style="width: 30%; margin-right: 5px;">
            <el-input v-model="config.input" placeholder="主机地址">
                <template #prepend>
                    <el-select v-model="config.defaultOption" style="width: 15vh;">
                        <el-option v-for="value in options" :label="value" :value="value" />
                    </el-select>
                </template>
                <template #append>
                    <el-button @click="addTarget">添加</el-button>
                </template>
            </el-input>
            <el-input v-model="config.target" :placeholder="placeholder2" type="textarea" resize="none"
                style="height: 94.5%;"></el-input>
        </div>
        <div style="width: 70%;">
            <el-table border :data="config.content" :cell-style="{ textAlign: 'center' }"
                :header-cell-style="{ 'text-align': 'center' }" style="width: 100%; height: 100% ;">
                <el-table-column prop="Host" label="Host" />
                <el-table-column prop="Port" label="Port" />
                <el-table-column prop="Protocol" label="Protocol" />
                <el-table-column prop="Username" label="Username" />
                <el-table-column prop="Password" label="Password" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </div>
    </div>
</template>

<style>
.el-textarea__inner {
    height: 100%;
}
</style>