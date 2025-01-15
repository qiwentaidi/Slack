<script lang="ts" setup>
import { InfoFilled, Search } from "@element-plus/icons-vue";
import { ref } from "vue";
import githubIcon from "@/assets/icon/github.svg";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import { GoFetch, NewDSStoreEngine } from "wailsjs/go/services/App";
import { ElMessage } from "element-plus";
import async from "async";
import { getProxy } from "@/util";

const introduceDialog = ref(false);
const url = ref("")
const log = ref("")
const isRunning = ref(false)

async function Start() {
    if (!url.value.startsWith("http")) {
        ElMessage.warning("请输入正确的URL")
        return
    }
    
    if (url.value.endsWith(".DS_Store")) {
        isRunning.value = true
        log.value += "[*] 正在进行.DS_Store信息提取\n"
        let links = await NewDSStoreEngine(url.value)
        if (links.length == 0) {
            log.value += "[!] 未检测到无链接，请检查文件内容是否正确\n"
            isRunning.value = false
            return
        }
        log.value += "[*] 检测到 " + links.length + " 个链接，正在探测存活状态\n"
        async.eachLimit(links, 5, async (link: string, callback: () => void) => {
            let respone = await GoFetch("GET", link, "", null, 10, getProxy())
            if (!respone.Error) {
                log.value += `[+] ${link} [${respone.StatsCode}] [${respone.Body.length}]\n`
            }
        },
            (err: any) => {
                if (err) {
                    log.value += "[-] dumpall has error: " + err + "\n"
                } else {
                    log.value += "[+] 已处理完成数据\n"
                }
                isRunning.value = false
            }
        );
    } else {
        ElMessage("目前仅支持.DS_Store文件检测")
    }
}
</script>


<template>
    <div class="my-header" style="margin-bottom: 10px;">
        <el-button plain :icon="InfoFilled" @click="introduceDialog = true">模块介绍</el-button>
        <el-input v-model="url" placeholder="请输入URL地址" style="margin-inline: 5px"></el-input>
        <el-button type="primary" :icon="Search" @click="Start" v-if="!isRunning">开始任务</el-button>
        <el-button type="primary" loading v-else>正在检测</el-button>
    </div>
    <highlightjs language="customlog" :code="log"></highlightjs>
    <el-dialog v-model="introduceDialog" title="模块介绍" width="60%">
        <el-descriptions :column="1" border>
            <el-descriptions-item label="描述">
                <el-text>适用于 .git | .svn | .DS_Store 信息泄漏时使用，功能详情参考
                    <el-link @click="BrowserOpenURL('https://github.com/0xHJK/dumpall')">dumpall<el-icon
                            class="el-icon--right">
                            <githubIcon />
                        </el-icon></el-link></el-text>
            </el-descriptions-item>
            <el-descriptions-item label=".git源代码泄漏(未完成)">所有代码Git操作后的记录</el-descriptions-item>
            <el-descriptions-item label=".svn源代码泄漏(未完成)">Linux版本控制系统,用于管理和跟踪文件和目录的变化</el-descriptions-item>
            <el-descriptions-item label=".DS_Store信息泄漏">MacOS文件夹图标展示文件，可以还原出当前文件夹下的文件名称 </el-descriptions-item>
        </el-descriptions>
    </el-dialog>
</template>


<style scoped></style>