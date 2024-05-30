<script setup lang="ts">
import { ref } from "vue";
import global from "./global";
import Sidebar from "./components/Sidebar.vue";
import { ArrowRight } from '@element-plus/icons-vue';
import { useRoute } from "vue-router";
const route = useRoute();
const showLogger = ref(false);

function breadcrumbItems(label: string) {
  switch (label) {
    case "/":
      return ["Home"];
    case "/settings":
      return ["Settings"];
    default:
      return label.slice(1).split('/');
  }
}
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
      <div class="console-log">
        <div>
          <el-breadcrumb :separator-icon="ArrowRight" style="margin-left: 2.5vh;">
            <el-breadcrumb-item v-for="label in breadcrumbItems(route.path)">
              {{ label }}
            </el-breadcrumb-item>
          </el-breadcrumb>
          <div style="margin-right: 2.5vh">
            <el-button link @click="showLogger = true">
              <template #icon>
                <img src="/console.svg">
              </template>
              Console
            </el-button>
          </div>

        </div>
      </div>
    </el-container>
  </el-container>
  <!-- running logs -->
  <el-drawer v-model="showLogger" direction="ltr" size="30%">
    <template #header>
      <h4>运行日志</h4>
    </template>
    <el-input class="log-textarea" v-model="global.Logger.value" type="textarea" style="height: 100%"
      resize="none"></el-input>
  </el-drawer>
</template>

<style>
.el-aside {
  width: 64px;
}
</style>
