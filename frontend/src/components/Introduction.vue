<script setup lang="ts">
import { LoadConfig } from '../util'
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    LoadConfig() // 加载配置信息
});
</script>

<template>
    <el-container class="el-main">
        <el-space direction="vertical">
        <a href="https://github.com/qiwentaidi/Slack" target="_blank" title="前往Github仓库">
            <img src="/slack.svg" class="logo" alt="Slack-Wails logo" />
        </a>
        <h2>
            Slack-Wails 一体化渗透测试平台
        </h2>
        </el-space>
        <h4>前端: Vue + Typescript + Vite + Element-Plus</h4>
        <h4>后端: Wails + Go</h4>
    </el-container>
</template>

<style scoped>


.el-main {
  position:absolute;
  left:50%;
  top:50%;
  transform:translate(-50%,-50%);
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
