<script setup lang="ts">
import { Download } from "@element-plus/icons-vue";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import global from "../global"
import { GoFetch } from "../../wailsjs/go/main/App";
import { UpdatePocFile, CheckFileStat, GetFileContent, UserHomeDir } from "../../wailsjs/go/main/File";
import { onMounted, reactive } from "vue";
import { ElNotification } from "element-plus";
import { compareVersion } from "../util"

onMounted(() => {
    check.client()
    check.poc()
});

const version = reactive({
    LocalPoc: "",
    RemotePoc: "",
    RemoteClient: "",
    PocUpdateContent: "",
    ClientUpdateContent: "",
});

const update = ({
    poc: async function () {
        let err = await UpdatePocFile()
        if (err == "") {
            ElNotification({
                title: "Success",
                message: "POC更新成功!",
                type: "success",
            });
        } else {
            ElNotification({
                title: "POC更新失败",
                message: err,
                type: "error",
            });
        }
    }
})


const check = ({
    // poc
    poc: async function () {
        let pcfg = await CheckFileStat(await UserHomeDir() + global.PATH.LocalPocVersionFile)
        if (!pcfg) {
            version.LocalPoc = "版本文件不存在"
            global.UPDATE.PocStatus = false
            return
        } else {
            version.LocalPoc = await GetFileContent(await UserHomeDir() + global.PATH.LocalPocVersionFile)
        }
        let resp = await GoFetch("GET", download.RemotePocVersion, "", [{}], 10, null)
        if (resp.Error == true) {
            version.PocUpdateContent = "检测更新失败"
            global.UPDATE.PocStatus = false
        } else {
            version.RemotePoc = resp.Body
            if (compareVersion(version.LocalPoc, version.RemotePoc) == -1) {
                version.PocUpdateContent = (await GoFetch("GET", download.PocUpdateCentent, "", [{}], 10, null)).Body
                global.UPDATE.PocStatus = true
            } else {
                version.PocUpdateContent = "已是最新版本"
                global.UPDATE.PocStatus = false
            }
        }
    },
    // client
    client: async function () {
        let resp = await GoFetch("GET", download.RemoteClientVersion, "", [{}], 10, null)
        if (resp.Error == true) {
            version.ClientUpdateContent = "检测更新失败"
            global.UPDATE.ClientStatus = false
        } else {
            version.RemoteClient = resp.Body
            if (compareVersion(global.LOCAL_VERSION, version.RemoteClient) == -1) {
                version.ClientUpdateContent = (await GoFetch("GET", download.ClientUpdateCentent, "", [{}], 10, null)).Body
                global.UPDATE.ClientStatus = true
            } else {
                version.ClientUpdateContent = "已是最新版本"
                global.UPDATE.ClientStatus = false
            }
        }
    }
})

const download = {
    RemotePocVersion:
        "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version",
    RemoteClientVersion:
        "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version",
    PocUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update',
    ClientUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update',
};

</script>

<template>
    <el-card class="box-card" style="width: 100%;">
        <template #header>
            <div class="card-header">
                <el-text>
                    <span style="font-weight: bold;">POC&指纹{{ version.RemotePoc }}</span>
                    <br />当前{{ "v" + version.LocalPoc }}
                </el-text>
                <el-button class="button" :icon="Download" text @click="update.poc"
                    v-if="global.UPDATE.PocStatus">立即下载</el-button>
                <span v-else>{{ version.PocUpdateContent }}</span>
            </div>
        </template>
        <el-scrollbar height="100px" style="white-space: pre-wrap;" v-if="global.UPDATE.PocStatus">
            {{ version.PocUpdateContent }}
        </el-scrollbar>
    </el-card>

    <el-card class="box-card" style="width: 100%; margin-top: 10px;">
        <template #header>
            <div class="card-header">
                <el-text>
                    <span style="font-weight: bold;">客户端{{ version.RemoteClient }}</span>
                    <br />当前{{ "v" + global.LOCAL_VERSION }}
                </el-text>
                <el-button class="button" :icon="Download" text
                    @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/releases')"
                    v-if="global.UPDATE.ClientStatus">立即下载</el-button>
                <span v-else>{{ version.ClientUpdateContent }}</span>
            </div>
        </template>
        <el-scrollbar height="100px" style="white-space: pre-wrap;" v-if="global.UPDATE.ClientStatus">
            {{ version.ClientUpdateContent }}
        </el-scrollbar>
    </el-card>
</template>