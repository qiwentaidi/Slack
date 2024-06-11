<script setup lang="ts">
import { reactive, onMounted } from 'vue';
import { QuestionFilled, Search, Document } from '@element-plus/icons-vue';
import { FileDialog, CheckFileStat } from '../../../wailsjs/go/main/File';
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
        // Vue.set(config.content, config.content.length, newItem);
        // config.renderKey++
        console.log(config.content)
    });
    return () => {
        EventsOff("bruteResult");
    };
});

const config = reactive({
    builtIn: true,
    association: false,
    username: '',
    password: '',
    target: '',
    content: [] as BruteResult[],
    renderKey: 0
})

const placeholder = `
1、使用内置字典时，账号密码输入框无效，要使用自己的字典请取消勾选使用内置字典

2、账号密码框处支持单字段比如admin和txt字典路径，可以拖拽读取文件路径

3、联想模式意为根据关键字段如出生年月或公司名等固定规律组成字典

4、支持协议名: ftp、ssh、telnet、smb、oracle、mssql、mysql、rdp、
postgresql、vnc、redis、memcached、mongodb
`

const placeholder2 = `e.g
redis://10.0.0.1:6379
mysql://10.0.0.1:3306
`

async function selectDict(input: string) {
    let filepath = await FileDialog()
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
        ElMessage("可用目标为空")
        return
    }
    // 内置模式处理字典
    if (!config.builtIn) {
        if (config.username != "") {
            let result = await CheckFileStat(config.username) // 文件不存在则为单用户名
            if (!result) {
                userDict.push(config.username)
            }else {
                userDict = (await ReadLine(config.username))!
            }
        }
        if (config.password != "") {
            let result = await CheckFileStat(config.password)
            if (!result) {
                 passDict.push(config.password)
            }else {
                passDict = (await ReadLine(config.password))!
            }
        }
    }
    config.content = []
    global.dict.passwords = (await ReadLine(global.PATH.PortBurstPath + "/password/password.txt"))!
    for (var item of global.dict.usernames) {
        item.dic = (await ReadLine(global.PATH.PortBurstPath + "/username/" + item.name + ".txt"))!
    }
    var id = 0
    
    async.eachLimit(url, 20, async (target: string, callback: () => void) => {
        let protocol = target.split("://")[0]
        // 如果当个字典不为空
        if (userDict.length == 0) {
            userDict = global.dict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
        }
        if (passDict.length == 0) {
            passDict = global.dict.passwords
        }
        Callgologger("info", target + " is start brute")
        await PortBrute(target, userDict, passDict)
        id++
        if (id == url.length) {
            callback()
        }
    }, (err: any) => {
        Callgologger("info", "PortBrute Finished")
        // ctrl.runningStatus = false
    });
}``

</script>

<template>
    <div class="my-header">
        <div style="display: flex;">
            <el-checkbox v-model="config.builtIn" label="内置字典" />
            <!-- <el-checkbox v-model="config.association" label="联想模式" /> -->
            <el-tooltip :content="placeholder" placement="right-end">
                <template #content>
                    1、使用内置字典时，账号密码输入框无效，要使用自己的字典请取消勾选使用内置字典<br />
                    2、账号密码框处支持单字段比如admin和txt字典路径，可以拖拽读取文件路径<br />
                    3、联想模式意为根据关键字段如出生年月或公司名等固定规律组成字典<br />
                    4、支持协议名: ftp、ssh、telnet、smb、oracle、mssql、mysql、rdp、vnc、redis<br />
                    postgresql、memcached、mongodb<br />
                </template>
                <el-button :icon="QuestionFilled" link style="margin-bottom: 4px; margin-left: 5px;"></el-button>
            </el-tooltip>

        </div>
        <el-button type="primary" :icon="Search" @click="NewScanner">开始暴破</el-button>
    </div>
    <el-form style="margin-top: 10px;">
        <el-form-item label="账号">
            <el-input v-model="config.username" :disabled="config.builtIn">
                <template #suffix>
                    <el-button link :icon="Document" @click="selectDict('username')"></el-button>
                </template>
            </el-input>
        </el-form-item>
        <el-form-item label="密码">
            <el-input v-model="config.password" :disabled="config.builtIn">
                <template #suffix>
                    <el-button link :icon="Document" @click="selectDict"></el-button>
                </template>
            </el-input>
        </el-form-item>
    </el-form>
    <splitpanes class="default-theme" style="height: 73vh;">
        <pane size="30">
            <el-input v-model="config.target" :placeholder="placeholder2" type="textarea" resize="none"
                style="height: 100%;"></el-input>
        </pane>
        <pane size="70">
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
        </pane>
    </splitpanes>
</template>

<style>
.el-textarea__inner {
    height: 100%;
}
</style>