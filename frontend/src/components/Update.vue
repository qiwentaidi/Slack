<script setup lang="ts">
import { Download } from "@element-plus/icons-vue";
import global from "../global"
import { UpdatePocFile, DownloadLastestClient, Restart } from "../../wailsjs/go/main/File";
import { ElMessageBox, ElNotification } from "element-plus";
import { onMounted, ref } from "vue";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";


onMounted(() => {
    // 监听下载进度事件
    EventsOn("clientDownloadProgress", (p: number) => {
        update.progress.value = p;
    });

    // 监听下载完成事件
    EventsOn("clientDownloadComplete", (msg: string) => {
        let message = ""
        if (msg == "mac-success") {
            message = "更新成功，是否立即安装?"
        } else {
            message = "更新成功，是否重新启动?"
        }
        update.downloadRunningStatus.value = false
        update.progress.value = 100;
            ElMessageBox.confirm(
                message,
                {
                    confirmButtonText: '确认',
                    cancelButtonText: '取消',
                    type: 'success',
                    center: true,
                }
            )
                .then(() => {
                    Restart()
                })
                .catch(() => {
                    console.log('User cancelled or chose another option.')
                })
    });

    // 清除事件监听器
    return () => {
        EventsOff("clientDownloadProgress");
        EventsOff("clientDownloadComplete");
    };
});

const update = ({
    poc: async function () {
        let err = await UpdatePocFile()
        if (err == "") {
            ElNotification.success("POC update success!");
        } else {
            ElNotification.error("POC update failed! " + err);
        }
    },
    client: async function () {
        update.downloadRunningStatus.value = true
        let result: any = await DownloadLastestClient()
        if (result.Error) {
            ElNotification.error("Download client failed! " + result.Msg);
        }
    },
    progress: ref(0),
    downloadRunningStatus: ref(false),
})

const customColorMethod = (percentage: number) => {
    if (percentage < 30) {
        return '#909399'
    }
    if (percentage < 70) {
        return '#e6a23c'
    }
    return '#67c23a'
}
</script>

<template>
    <el-card class="box-card">
        <template #header>
            <div class="card-header">
                <span style="font-weight: bold;">POC&指纹: 最新{{ global.UPDATE.RemotePocVersion }}/当前{{
                    global.UPDATE.LocalPocVersion }}</span>
                <el-button type="primary" :icon="Download" text @click="update.poc"
                    v-if="global.UPDATE.PocStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.PocContent }}</span>
            </div>
        </template>
        <el-scrollbar style="height: 100px;" class="pretty-response" v-if="global.UPDATE.PocStatus">
            {{ global.UPDATE.PocContent }}
        </el-scrollbar>
    </el-card>

    <el-card class="box-card" style="margin-top: 10px;">
        <template #header>
            <div class="card-header">
                <span style="font-weight: bold;">客户端: 最新{{ global.UPDATE.RemoteClientVersion }}/当前{{
                    global.LOCAL_VERSION }}</span>
                <el-button type="primary" :icon="Download" text @click="update.client"
                    v-if="global.UPDATE.ClientStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.ClientContent }}</span>
            </div>
        </template>
        <el-scrollbar style="height: 100px;" class="pretty-response" v-if="global.UPDATE.ClientStatus">
            {{ global.UPDATE.ClientContent }}
        </el-scrollbar>
    </el-card>
    <div style="margin-top: 5px;" v-if="update.downloadRunningStatus.value">正在下载中，请等待下载完成...
        <el-progress :percentage="update.progress.value" :color="customColorMethod"></el-progress>
    </div>
</template>