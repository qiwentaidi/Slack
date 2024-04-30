<script setup lang="ts">
import { ref } from "vue";
import global from "./global";
import Sidebar from "./components/Sidebar.vue";
import { ArrowRight } from '@element-plus/icons-vue';
import { useRoute } from "vue-router";
const route = useRoute();
const showLogger = ref(false);

function breadcrumbItems(label: string) {
  return label.slice(1).split('/');
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
          <el-breadcrumb :separator-icon="ArrowRight" style="margin-left: 2.5vh;" v-if="breadcrumbItems(route.path).length > 1">
              <el-breadcrumb-item v-for="label in breadcrumbItems(route.path)">
                {{ label }}
              </el-breadcrumb-item>
          </el-breadcrumb>
          <div v-else></div>
          <el-button link @click="showLogger = true" style="margin-right: 5px">
            <template #icon>
              <svg class="bi bi-terminal" style="height: 100%" viewBox="0 0 16 16" fill="currentColor"
                xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd"
                  d="M14 2H2a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V3a1 1 0 0 0-1-1zM2 1a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2H2z" />
                <path fill-rule="evenodd"
                  d="M6 9a.5.5 0 0 1 .5-.5h3a.5.5 0 0 1 0 1h-3A.5.5 0 0 1 6 9zM3.146 4.146a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1 0 .708l-2 2a.5.5 0 1 1-.708-.708L4.793 6.5 3.146 4.854a.5.5 0 0 1 0-.708z" />
              </svg>
            </template>
            Console
          </el-button>
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
