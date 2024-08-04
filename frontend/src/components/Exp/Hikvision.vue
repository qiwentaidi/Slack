<template>
    <div style="margin-bottom: 10px; display: flex;">
        <el-button text bg>漏洞名称:</el-button>
        <el-select v-model="form.selectVulnerability" style="width: 40%;">
            <el-option v-for="vulnerability in vulnerabilityGroup" :value="vulnerability.value"
                :label="vulnerability.name">
                {{ vulnerability.name }} <el-tag v-show="vulnerability.isBatch">可批量</el-tag>
            </el-option>
        </el-select>
        <el-button text bg>CMD:</el-button>
        <el-input v-model="form.cmd"></el-input>
        <el-button type="primary" @click="useVulnerability(form.selectVulnerability)"
            style="margin-left: 5px;">执行</el-button>
        <el-button :icon="Picture" @click="useVulnerability(0)" :disabled="form.selectVulnerability != 1"
            style="margin-left: 5px;">查看快照</el-button>
    </div>
    <el-row :gutter="8" style="margin-top: 10px;">
        <el-col :span="14">
            <div class="my-header" style="background-color: #eee;">
                <span>目标地址:

                </span>
                <el-button size="small" @click="uploadFile">目标导入</el-button>
            </div>
            <el-input v-model="form.url" type="textarea" rows="5" resize="none" placeholder="请输入URL根路径"></el-input>
        </el-col>
        <el-col :span="10">
            <div class="my-header" style="background-color: #eee;">
                密码口令:
                <el-button size="small" @click="form.passwordList = ''">清空</el-button>
            </div>
            <el-input v-model="form.passwordList" type="textarea" rows="5" resize="none"></el-input>
        </el-col>
    </el-row>


    <pre class="pretty-response" style="height: 58vh;"><code>{{ form.content }}</code></pre>
    <el-dialog v-model="form.snapshotDialog" title="快照" width="530">
        <el-image style="width: 500px; height: 500px" :src="form.image" loading="eager" fit="fill" />
    </el-dialog>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { Picture } from '@element-plus/icons-vue';
import { HikvsionCamera } from 'wailsjs/go/main/App';
import global from '@/global';
import { formatURL } from '@/util';
import { FileDialog, ReadFile } from 'wailsjs/go/main/File';
import { File } from '@/interface';
import { ElMessage } from 'element-plus';
import async from 'async';

const form = reactive({
    url: "",
    cmd: "",
    passwordList: "hik12345+\nHik12345+\n12345\nadmin12345",
    image: "",
    snapshotDialog: false,
    selectVulnerability: 1,
    content: "弱口令检测通过chromedp包实现，使用过程中可能会有些网站访问失败，所以在请求超过1次失败后就会直接停止目标爆破，如果网站可以打开请自行测试。",
})

const vulnerabilityGroup = [
    {
        name: "CVE-2017-7921",
        isBatch: false,
        value: 1
    },
    {
        name: "CVE-2021-36260",
        isBatch: false,
        value: 2
    },
    {
        name: "弱口令检测",
        isBatch: true,
        value: 3
    },
]

async function useVulnerability(mode: number) {
    let urls = await formatURL(form.url)
    if (urls.length == 0) {
        ElMessage.warning("请输入URL地址")
        return
    }
    if (mode == 0) {
        form.snapshotDialog = true
        let body = await HikvsionCamera(urls[0], 0, form.passwordList.split("\n"), form.cmd, global.proxy)
        const base64Image = `data:image/jpeg;base64,${body}`;
        form.image = base64Image
        return
    }
    if (mode == 3) {
        let id = 0
        async.eachLimit(urls, 10, async (url: string, callback: () => void) => {
            let result = await HikvsionCamera(url, 3, form.passwordList.split("\n"), form.cmd, global.proxy)
            form.content += result + "\n"
            id++
            if (id == urls.length) {
                callback()
            }
        }, () => {
            form.content += "弱口令检测结束！"
        })
    }
    // "http://47.150.37.246:81/"
    form.content = await HikvsionCamera(form.url, mode, form.passwordList.split("\n"), form.cmd,global.proxy)
}

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
    form.url = file.Content!
}
</script>