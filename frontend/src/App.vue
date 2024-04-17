<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'

</script>

<template>
  <el-container>
    <el-aside>
      <Sidebar></Sidebar>
    </el-aside>
    <el-container>
      <el-main>
        <!-- 一定要使用插槽否则keey-alive不会生效 -->
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component"></component>
          </keep-alive>
        </router-view>
      </el-main>
      <!-- <el-footer style="height: 25px; background-color: #F2F3F5;">
        <el-breadcrumb separator="/" style="display: flex; align-items: center; height: 100%;">
          <el-breadcrumb-item>homepage</el-breadcrumb-item>
        </el-breadcrumb>
      </el-footer> -->
    </el-container>
    
  </el-container>
</template>

<style>
.el-aside {
  width: 64px;
}
</style>
