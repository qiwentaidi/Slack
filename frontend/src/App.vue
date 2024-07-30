<script setup lang="ts">
import { onMounted, computed } from "vue";
import global from "./global";
import Sidebar from "./components/Sidebar.vue";
import Titlebar from "./components/Titlebar.vue";
import { EventsOn } from "../wailsjs/runtime/runtime";
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import { LogInfo } from "./interface";

const locale = computed(() => (global.Language.value === 'zh' ? zhCn : en))

const logArray = [] as string[]

onMounted(() => {
  const levelClassMap: { [key: string]: string } = {
    "[INF]": "log-info",
    "[WRN]": "log-warning",
    "[ERR]": "log-error",
    "[DEB]": "log-debug",
    "[SUC]": "log-success"
  };
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
})
</script>

<template>
  <Titlebar />
  <el-container>
    <el-aside>
      <Sidebar />
    </el-aside>
    <el-main>
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