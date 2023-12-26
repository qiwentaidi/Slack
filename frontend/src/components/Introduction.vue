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
    <el-space  direction="vertical">
      <a href="https://github.com/qiwentaidi/Slack" target="_blank" title="前往Github仓库">
        <img src="/slack.svg" class="logo" />
      </a>
      <h2>Slack 信安人员集成化工具平台</h2>
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
