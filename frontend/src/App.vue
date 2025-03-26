<script setup lang="ts">
import { onMounted, computed } from "vue";
import global from "./stores";
import Sidebar from "./components/Sidebar.vue";
import Titlebar from "./components/Titlebar.vue";
import { EventsOn } from "wailsjs/runtime/runtime";
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import { LogInfo } from "./stores/interface";
import { useDark } from '@vueuse/core'
import { NetworkCardInfo, UserHomeDir } from "wailsjs/go/services/File";
import { InitConfigFile } from "./config";
import { check } from "@/util";
import CyberChef from "./views/Tools/CyberChef.vue";
import { ElMessage } from "element-plus";
import { GOOS } from "wailsjs/go/core/Tools";

const levelClassMap: { [key: string]: string } = {
    "[INF]": "log-info",
    "[WRN]": "log-warning",
    "[ERR]": "log-error",
    "[DEB]": "log-debug",
    "[SUC]": "log-success"
};

const messageType = {
    "error": ElMessage.error,
    "warning": ElMessage.warning,
    "success": ElMessage.success,
    "info": ElMessage.info
};

// 初始化主题
useDark({
    storageKey: 'theme',
    valueDark: 'dark',
    valueLight: 'light',
})

// 初始化语言
const locale = computed(() => (global.Language.value === 'zh' ? zhCn : en))

var logArray = [] as string[]

onMounted(async () => {
    // 初始化目录
    global.PATH.homedir = await UserHomeDir();
    // 获取系统类型
    global.temp.goos = await GOOS();
    // 初始化配置文件
    await InitConfigFile(500);
    // 检测更新
    check.client();
    check.poc();
    // 初始化网卡
    let list = await NetworkCardInfo()
    global.temp.NetworkCardList.push(...list)
    // 初始化日志
    EventsOn("gologger", (log: LogInfo) => {
        const logClass = levelClassMap[log.Level];
        const logEntry = `<span class="${logClass}">${log.Level}</span> ${log.Msg}`;
        logArray.push(logEntry);
        // 最大存储数
        if (logArray.length > global.Logger.length) {
            logArray.shift();
        }
        global.Logger.value = logArray.join('\n');
    });
    EventsOn("gomessage", (log: LogInfo) => {
        messageType[log.Level]({
            message: log.Msg,
            plain: true,
            showClose: true,
        });
    })
});
</script>

<template>
    <Titlebar />
    <el-container>
        <el-aside style="width: 64px;">
            <Sidebar />
        </el-aside>
        <el-main :class="{ 'no-scroll': $route.path == '/Tools/CyberChef' }" class="content-wrapper">
            <el-config-provider :locale="locale">
                <!-- 一定要使用插槽否则keey-alive不会生效 -->
                <router-view v-slot="{ Component }">
                    <keep-alive>
                        <component :is="Component"></component>
                    </keep-alive>
                </router-view>
            </el-config-provider>
            <div v-show="$route.path == '/Tools/CyberChef'">
                <CyberChef />
            </div>
        </el-main>
    </el-container>
</template>

<style>
.content-wrapper {
    max-height: calc(100vh - 35px);
    overflow-y: auto;
    scrollbar-width: none;
}

.no-scroll {
    overflow-y: hidden;
    /* 禁用滚动 */
}
</style>