<script lang="ts" setup>
import { ElMessage } from 'element-plus';
import { CheckFileStat, UserHomeDir } from 'wailsjs/go/services/File';
import { DownloadCyberChef, CyberChefLocalServer } from 'wailsjs/go/services/App';
import { onMounted, ref, reactive } from 'vue';
import { EventsOn, EventsOff } from 'wailsjs/runtime/runtime';
import fishIcon from '@/assets/icon/fish.svg'
import bearIcon from '@/assets/icon/bear.svg'

const showIframe = ref(false);
const isRemote = ref(false);
const progress = ref(0);
const config = reactive({
    LocalENV: false,
    downloadRunningStatus: false
})

onMounted(() => {
    // 监听下载进度事件
    EventsOn("downloadProgress", (p: number) => {
        progress.value = p;
    });

    // 监听下载完成事件
    EventsOn("downloadComplete", (file: string) => {
        config.LocalENV = true
        config.downloadRunningStatus = false
        progress.value = 100;
    });

    // 清除事件监听器
    return () => {
        EventsOff("downloadProgress");
        EventsOff("downloadComplete");
    };
});

onMounted(async () => {
    if (await CheckFileStat(await UserHomeDir() + "/slack/CyberChef")) {
        config.LocalENV = true
    }
});

const url = "https://gitee.com/the-temperature-is-too-low/Slack/releases/download/v1/CyberChef.zip"

function startDownload() {
    if (!config.downloadRunningStatus) {
        config.downloadRunningStatus = true
    }else {
        ElMessage.warning({
            message: '正在下载中，请等待下载完成',
            duration: 2000,
        })
        return
    }
    // 调用后端下载函数
    DownloadCyberChef(url);
}

function loadRemote() {
    isRemote.value = true;
    showIframe.value = true;
}

async function loadLocal() {
    await CyberChefLocalServer()
    isRemote.value = false;
    showIframe.value = true;
}
</script>


<template>
    <div class="container" v-if="!showIframe">
        <el-result title="本地加载" sub-title="需要下载解密模块，后续会在本机8731端口部署CyberChef服务，适用部分内网环境，一次下载后续优先使用">
            <template #icon>
                <fishIcon />
            </template>
            <template #extra>
                <el-button type="primary" @click="startDownload()" v-if="!config.LocalENV">开始下载</el-button>
                <el-button type="primary" @click="loadLocal()" v-if="config.LocalENV">选择</el-button>
            </template>
        </el-result>
        <el-result title="远程加载" sub-title="远程内嵌官网站点，网络不好可能会加载失败，但不用下载环境">
            <template #icon>
                <bearIcon />
            </template>
            <template #extra>
                <el-button type="primary" @click="loadRemote()">选择</el-button>
            </template>
        </el-result>
    </div>
    <div class="external-page" v-if="showIframe">
        <iframe :src="isRemote ? 'https://gchq.github.io/CyberChef/' : 'http://127.0.0.1:8731/'" class="iframe"></iframe>
    </div>
    <el-progress class="download-progress" :percentage="progress" v-if="progress > 0 && progress < 100"></el-progress>
</template>

<style scoped>
.external-page {
    width: 100%;
    height: 91vh;
}

.container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    top: 50%;
    transform: translate(0, 50%);
}

.download-progress {
    position: absolute; 
    bottom: 15%; 
    width: 100%;
}
</style>