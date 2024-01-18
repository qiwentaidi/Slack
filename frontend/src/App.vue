<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'

</script>

<template>
  <el-container>
    <el-aside>
      <Sidebar></Sidebar>
    </el-aside>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
</template>

<style>
.el-aside {
  width: 70px;
  height: 100vh;
}
</style>
