<script setup lang="ts">
import global from "../global";
import { onMounted } from 'vue';
import { BrowserOpenURL } from '../../wailsjs/runtime'
import { CheckFileStat, InitConfig, UserHomeDir } from '../../wailsjs/go/main/File';
import { ElNotification } from 'element-plus';
import Loading from '../components/Loading.vue';
// 初始化时调用
onMounted(async () => {
  LoadConfig() // 加载配置信息
  let cfg = await CheckFileStat(await UserHomeDir() + "/slack")
  if (!cfg) {
    ElNotification({
      duration: 0,
      message: '未检测到配置文件，正在初始化...',
      icon: Loading,
    });
    if (await InitConfig()) {
      ElNotification.closeAll()
      ElNotification({
        message: "配置文件初始化成功!",
        type: "success",
      });
    } else {
      ElNotification.closeAll()
      ElNotification({
        title: "无法下载配置文件",
        message: "请自行到https://gitee.com/the-temperature-is-too-low/slack-poc/releases/下载config.zip并解压到用户根目录下的/slack/文件夹下!",
        type: "warning",
      });
    }
  }
});

function LoadConfig() {
  const allLocaolStorage = [
    {
      key: "scan",
      value: global.scan,
    },
    {
      key: "proxy",
      value: global.proxy,
    },
    {
      key: "space",
      value: global.space,
    },
    {
      key: "jsfind",
      value: global.jsfind,
    },
    {
      key: "webscan",
      value: global.webscan,
    }
  ];
  allLocaolStorage.forEach(item => {
    const v = localStorage.getItem(item.key)
    if (v) {
      Object.assign(item.value, JSON.parse(v));
    }
  });
}
</script>

<template>
  <el-container class="el-main">
    <el-space direction="vertical">
      <a @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack')" title="前往Github仓库">
        <img src="/slack.svg" class="logo" />
      </a>
      <h2>{{ $t('aside.slogan') }}</h2>
    </el-space>
  </el-container>
</template>

<style scoped>
.el-main {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
</style>
