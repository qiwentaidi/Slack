<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'

</script>

<template>
  <el-container>
    <el-aside>
      <Sidebar></Sidebar>
    </el-aside>
    <el-main>
      <!-- 一定要使用插槽否则keey-alive不会生效 -->
      <router-view v-slot="{ Component }"> 
        <keep-alive>
          <component :is="Component"></component>
        </keep-alive>
      </router-view>
    </el-main>
  </el-container>
</template>

<style>
.el-aside {
  width: 70px;
  height: 100vh;
}
</style>
