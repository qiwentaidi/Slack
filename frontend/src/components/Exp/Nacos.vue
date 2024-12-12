<script lang="ts" setup>
import { AlibabaNacos } from 'wailsjs/go/services/App';
import { reactive } from 'vue';
import { AddRightSubString, getProxy } from '@/util';
const form = reactive({
    url: "",
    username: "qioxzcio",
    password: "Aoioxzcoio123.",
    command: "whoami",
    header: "",
    content: "",
    selectVulnerability: 0
})


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
    // {
    //     name: "Derby SQLi 条件竞争 RCE",
    //     value: 3
    // },
]

async function useVulnerability() {
    form.url = AddRightSubString(form.url, "/")
    form.content = await AlibabaNacos(form.url,form.header ,form.selectVulnerability, form.username, form.password, form.command, "", getProxy())
}
</script>

<template>
    <el-form @submit.native.prevent="">
        <el-form-item>
            <div class="head">
                <el-input v-model="form.url" placeholder="请输入Nacos主页路径，例如 https://192.168.1.1/nacos/">
                    <template #prepend>
                        <el-select style="width: 300px;" placeholder="请选择漏洞" v-model="form.selectVulnerability">
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
    </el-form>
    <pre class="pretty-response" style="height: 63vh;"><code>{{ form.content }}</code></pre>
</template>

<style scoped>
.ml5 {
    margin-left: 5px;
}
</style>