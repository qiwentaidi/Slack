<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue';
import { Document, Edit, Setting, UploadFilled, Upload } from '@element-plus/icons-vue';
import resultIcon from '@/assets/icon/result.svg'
import { ElMessage } from 'element-plus';
import async from 'async';
import global from '@/global';
import { Copy, ReadLine, ProcessTextAreaInput, UploadFileAndRead } from '@/util';
import { PortBrute, Callgologger, ExitScanner } from 'wailsjs/go/services/App';
import { EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import { BruteResult } from '@/stores/interface';
import usePagination from '@/usePagination';
import { CheckFileStat, FileDialog } from 'wailsjs/go/services/File';
import { crackDict } from '@/stores/options';

onMounted(() => {
    EventsOn("bruteResult", (result: BruteResult) => {
        let temp = result.Host.split(":")
        pagination.table.result.push({
            Host: temp[0],
            Port: temp[1],
            Protocol: result.Protocol,
            Username: result.Username,
            Password: result.Password,
        })
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    });
    return () => {
        EventsOff("bruteResult");
    };
});

const config = reactive({
    builtInUsername: true,
    builtInPassword: true,
    username: '',
    password: '',
    target: '',
    defaultOption: "ftp",
    input: '',
    runningStatus: false,
    thread: 20,
})

async function addTarget() {
    if (!config.input) {
        ElMessage.warning("请输入目标")
        return
    }
    let stat = await CheckFileStat(config.input)
    if (stat) {
        let lines = (await ReadLine(config.input))!
        for (const line of lines) {
            config.target += config.defaultOption + "://" + line + "\n"
        }
        return
    }
    config.target += config.defaultOption + "://" + config.input
}
async function NewScanner() {
    let url = [] as string[]
    for (const item of config.target.split("\n")) {
        if (item.includes("://")) {
            url.push(item)
        }
    }
    if (url.length == 0) {
        ElMessage.warning("可用目标为空")
        return
    }
    let passDict = [] as string[]
    let userDict = [] as string[]
    pagination.initTable()
    if (config.builtInUsername) {
        for (var item of crackDict.usernames) {
            item.dic = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/username/" + item.name + ".txt"))!
        }
    }
    if (config.builtInPassword) {
        passDict = (await ReadLine(global.PATH.homedir + global.PATH.PortBurstPath + "/password/password.txt"))!
        passDict.push("")
    }
    if (config.username != "") {
        let users = ProcessTextAreaInput(config.username)
        for (var item of crackDict.usernames) {
            item.dic.push(...users)
        }
    }
    if (config.password != "") {
        let pass = ProcessTextAreaInput(config.password)
        passDict.push(...pass)
    }
    ElMessage.info("任务已开始")
    activeRef.value++
    config.runningStatus = true
    var id = 0
    async.eachLimit(url, config.thread, async (target: string, callback: () => void) => {
        if (!config.runningStatus) {
            callback()
        }
        let protocol = target.split("://")[0]
        userDict = crackDict.usernames.find(item => item.name.toLocaleLowerCase() === protocol)?.dic!
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

function stopScan() {
    if (config.runningStatus) {
        ExitScanner("[portbrute]")
        config.runningStatus = false
        ElMessage.error("用户已终止任务");
    }
}

const activeRef = ref(1)
function NextStep() {
    if (activeRef.value == 3) return
    activeRef.value++
}

function Back() {
    if (activeRef.value == 1) return
    activeRef.value--
}

async function uploadFile() {
    config.target = await UploadFileAndRead()
}

async function uploadUserFile() {
    config.username = await UploadFileAndRead()
}

async function uploadPassFile() {
    config.password = await UploadFileAndRead()
}

async function getFilepath() {
    let path = await FileDialog("*.txt")
    if (!path) return
    config.input = path
}

const pagination = usePagination<BruteResult>(20)

function checkinput() {
    if (config.username.length == 0) {
        config.builtInUsername = true
    }
    if (config.password.length == 0) {
        config.builtInPassword = true
    }
}
</script>

<template>
    <div style="height: 100%; position: relative;">
        <el-steps :active="activeRef" simple>
            <el-step title="输入目标" :icon="Edit" />
            <el-step title="设置参数" :icon="Setting" />
            <el-step title="结果输出" :icon="resultIcon" />
        </el-steps>
        <div style="margin-top: 20px;"></div>
        <el-form :model="config" label-width="auto" v-show="activeRef == 1" style="width: 60%;">
            <el-form-item label="目标:">
                <el-input v-model="config.target" placeholder="请输入目标，目标仅支持换行分割
扫描默认端口 ssh://10.0.0.1
指定端口 redis://10.0.0.1:6380

Memcachedb仅支持未授权检测
" type="textarea" resize="none" style="height: 50vh;" />
                <el-button link size="small" :icon="Upload" @click="uploadFile"
                    style="margin-top: 5px;">导入目标文件</el-button>
            </el-form-item>
            <el-form-item label="添加目标:">
                <el-input v-model="config.input" placeholder="请输入主机地址或者文件路径添加目标前缀">
                    <template #prepend>
                        <el-select v-model="config.defaultOption" style="width: 15vh;">
                            <el-option v-for="value in crackDict.options" :label="value" :value="value" />
                        </el-select>
                    </template>
                    <template #suffix>
                        <el-button link :icon="Document" @click="getFilepath"></el-button>
                    </template>
                    <template #append>
                        <el-button @click="addTarget">添加</el-button>
                    </template>
                </el-input>
            </el-form-item>
            <el-form-item label="复制前缀:">
                <el-row :gutter="10">
                    <el-col :span="4" v-for="item in crackDict.options" :key="item">
                        <el-button size="small" @click="Copy(item + '://')"
                            style="margin-bottom: 10px; width: 100%;">
                            {{ item }}
                        </el-button>
                    </el-col>
                </el-row>
            </el-form-item>
        </el-form>

        <el-form label-width="auto" v-show="activeRef == 2">
            <el-form-item label="目标并发:">
                <el-input-number v-model="config.thread"></el-input-number>
            </el-form-item>
            <el-form-item label="用户字典:">
                <el-input v-model="config.username" type="textarea" :rows="5" @input="checkinput"></el-input>
                <div style="display: flex;">
                    <el-checkbox v-model="config.builtInUsername"
                        :disabled="config.username.length == 0">使用默认用户字典</el-checkbox>
                    <el-button link type="primary" :icon="UploadFilled" @click="uploadUserFile"
                        style="margin-inline: 10px;">
                        上传文件
                    </el-button>
                    <el-tag type="info" style="margin-top: 5px;">Tips: 默认字典可通过 设置-> 字典管理处修改</el-tag>
                </div>
            </el-form-item>
            <el-form-item label="密码字典:">
                <el-input v-model="config.password" type="textarea" :rows="5" @input="checkinput"></el-input>
                <div style="display: flex;">
                    <el-checkbox v-model="config.builtInPassword"
                        :disabled="config.password.length == 0">使用默认密码字典</el-checkbox>
                    <el-button link type="primary" :icon="UploadFilled" @click="uploadPassFile"
                        style="margin-left: 10px;">
                        上传文件
                    </el-button>
                </div>
            </el-form-item>
        </el-form>
        <div v-show="activeRef == 3">
            <el-table border :data="pagination.table.pageContent" :cell-style="{ textAlign: 'center' }"
                :header-cell-style="{ 'text-align': 'center' }" style="height: calc(100vh - 175px)">
                <el-table-column prop="Host" label="Host" />
                <el-table-column prop="Port" label="Port" />
                <el-table-column prop="Protocol" label="Protocol" />
                <el-table-column prop="Username" label="Username" />
                <el-table-column prop="Password" label="Password" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
            <div class="my-header" style="margin-top: 5px;">
                <el-button type="primary" @click="Back" v-show="activeRef == 3 && !config.runningStatus">上一步</el-button>
                <el-button type="danger" @click="stopScan" v-show="config.runningStatus">结束任务</el-button>
                <el-pagination background @size-change="pagination.ctrl.handleSizeChange"
                    @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                    :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100]"
                    :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                    :total="pagination.table.result.length">
                </el-pagination>
            </div>
        </div>


        <div style="position: absolute; bottom: 0; right: 0;">
            <el-button type="primary" @click="Back" v-show="activeRef == 2">上一步</el-button>
            <el-button type="primary" @click="NextStep" v-show="activeRef == 1">下一步</el-button>
            <el-button type="primary" @click="activeRef = 3" v-show="activeRef == 1">查看联动结果</el-button>
            <el-button type="primary" @click="NewScanner" v-show="activeRef == 2">开始任务</el-button>
        </div>
    </div>
</template>
