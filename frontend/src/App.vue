<script setup lang="ts">
import { onMounted, computed } from "vue";
import global from "./global";
import Sidebar from "./components/Sidebar.vue";
import Titlebar from "./components/Titlebar.vue";
import { EventsOn } from "wailsjs/runtime/runtime";
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import { LogInfo } from "./interface";
import { useDark } from '@vueuse/core'
import { NetworkCardInfo, UserHomeDir } from "wailsjs/go/main/File";
import { InitConfigFile } from "./config";
import { check } from "@/util";

const levelClassMap: { [key: string]: string } = {
  "[INF]": "log-info",
  "[WRN]": "log-warning",
  "[ERR]": "log-error",
  "[DEB]": "log-debug",
  "[SUC]": "log-success"
};

// 初始化主题
useDark({
  storageKey: 'theme',
  valueDark: 'dark',
  valueLight: 'light',
})

// 初始化语言
const locale = computed(() => (global.Language.value === 'zh' ? zhCn : en))

const logArray = [] as string[]

onMounted(async () => {
  // 初始化目录
  global.PATH.homedir = await UserHomeDir();
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
});
</script>

<template>
  <Titlebar />
  <el-container>
    <el-aside style="width: 64px;">
      <Sidebar />
    </el-aside>
    <el-main :class="{ 'no-scroll': $route.path == '/Tools/Codec' || $route.path == '/Tools/Reverse' }" class="content-wrapper">
      <el-config-provider :locale="locale">
        <!-- 一定要使用插槽否则keey-alive不会生效 -->
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component"></component>
          </keep-alive>
        </router-view>
      </el-config-provider>
    </el-main>
  </el-container>
</template>

<style>
.content-wrapper {
  max-height: calc(100vh - 35px);
  overflow-y: auto;
}

.no-scroll {
    overflow-y: hidden; /* 禁用滚动 */
}
</style>