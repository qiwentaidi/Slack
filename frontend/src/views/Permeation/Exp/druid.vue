<script lang="ts" setup>
import { ElMessage } from 'element-plus';
import { UseDruid } from '../../../../wailsjs/go/main/App';
import { reactive } from 'vue';
import global from "../../../global/index"
import { CopyDocument } from '@element-plus/icons-vue';
import { Copy, validateURL } from '../../../util';
const form = reactive({
    url: "",
    username: "",
    password: "",
    content: ""
})

const Login = 1
const Session = 2

function TestDruid(mode: number) {
    if (form.url == "") {
        ElMessage("url can't be empty")
        return
    }
    if (!validateURL(form.url)) {
        return
    }
    try {
        let url = form.url.split("/druid/")[0]
        if (mode == 1) {
            LoginDruid(url)
        }else {
            GetSeesion(url)
        }
    } catch {
        form.content = "please enter the correct druid path"
    }
}

async function GetSeesion(url: string) {
    const seesions = await UseDruid(url, Session, global.proxy)
    if (seesions.length == 0) {
        form.content = "[+] seesion count is 0"
    } else {
        form.content = seesions.join("\n")
        form.content = "[+] success get sessions: " + seesions.length + "\n"
    } 
}

async function LoginDruid(url: string) {
    const result = await UseDruid(url, Login, global.proxy)
    form.content = result.join("\n")
}
</script>

<template>
    <el-form label-width="auto" @submit.native.prevent="GetSeesion">
        <el-form-item label="URL:">
            <div class="head">
                <el-input v-model="form.url" placeholder="/druid/index.html or /druid/login.html"></el-input>
                <el-button @click="TestDruid(1)" class="ml5">Login Brute</el-button>
                <el-button @click="TestDruid(2)" class="ml5">Extract Session</el-button>
                <el-button @click="Copy(form.content)" class="ml5" :icon="CopyDocument"></el-button>
            </div>
        </el-form-item>
        <!-- <el-form-item label="Auth:">
            <el-space>
                <el-input style="width: 50%;">
                    <template #prepend>username</template>
                </el-input>
                <el-input style="width: 50%;">
                    <template #prepend>password</template>
                </el-input>
            </el-space>
        </el-form-item> -->
    </el-form>
    <pre class="pretty-response" style="height: 70vh;"><code>{{ form.content }}</code></pre>
</template>

<style scoped>
.ml5 {
  margin-left: 5px;
}
</style>