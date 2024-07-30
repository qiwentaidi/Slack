<script setup lang="ts">
import { Refresh, Setting } from "@element-plus/icons-vue";
import { onMounted, ref } from "vue";
import updateUI from "./Update.vue";
import global from "../global";
import menus from "../router/menu";
import { check } from "../util";
onMounted(() => {
  check.client();
  check.poc();
});



const updateDialog = ref(false)
</script>

<template>
  <div class="flex-box-v" style="height: calc(100vh - 35px);">
    <el-menu :collapse="true" :router="true" :default-active="$route.path" active-text-color="--sidebar-active-text-color"
      background-color="--sidebar-bg-color" text-color="--sidebar-text-color">
      <template v-for="menu in menus">
        <el-menu-item v-if="!menu.children" :index="menu.path">
          <el-icon>
            <img :src="menu.icon" class="custom-svg">
          </el-icon>
          <template #title><span>{{ $t(menu.name) }}</span></template>
        </el-menu-item>
        <el-sub-menu :index="menu.path" v-else>
          <template #title>
            <el-icon class="aside">
              <component :is="menu.icon" />
            </el-icon>
            <span>{{ $t(menu.name) }}</span>
          </template>
          <el-menu-item v-for="item in menu.children" :key="item.path" :index="item.path">{{ $t(item.name)
            }}</el-menu-item>
        </el-sub-menu>
      </template>
    </el-menu>

    <div class="copy-menu"></div>

    <el-menu :collapse="true" active-text-color="#fff" background-color="--sidebar-bg-color" text-color="#000">
      <el-menu-item class="custom-menu-item" index="/update" @click="updateDialog = true">
        <el-icon class="aside">
          <Refresh />
        </el-icon>
        <template #title><span>{{ $t("aside.update") }}</span></template>
        <el-badge value="new" v-if="global.UPDATE.ClientStatus ||
      global.UPDATE.PocStatus
      " />
      </el-menu-item>

      <el-menu-item class="custom-menu-item" @click="$router.push('/settings')">
        <el-icon class="aside">
          <setting />
        </el-icon>
        <template #title><span>{{ $t("aside.setting") }}</span></template>
      </el-menu-item>
    </el-menu>
  </div>

  <!-- update -->
  <el-dialog v-model="updateDialog" title="更新通知" width="40%">
    <updateUI></updateUI>
  </el-dialog>
</template>

<style>
.el-badge {
  margin-bottom: 80px;
  left: 15px;
  position: absolute;
}

.el-menu-item.is-active {
  color: #000;
}

.el-menu-item.is-active::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 5px;
  /*  色块的宽度 */
  height: 100%;
  background-color: #3875f6;
  /*  色块的颜色 */
  border-radius: 0 3px 3px 0;
  /* 轨道的形状 */
}

.el-menu-item {
  font-size: 16px;
}
.el-sub-menu__title {
  font-size: 16px;
}

.custom-menu-item:hover {
  /* 定义鼠标悬停时的样式 */
  background-color: #f5f5f5;
  /* 例如，使用一个浅灰色作为背景 */
  color: #3875f6;
}

.custom-menu-item.is-active::before {
  /* 定义激活状态时的样式 */
  background-color: inherit;
  /* 使用继承的背景颜色 */
  color: inherit;
  /* 使用继承的文本颜色 */
}

.aside svg {
  height: 1.5em;
  width: 1.5em;
}

.custom-svg {
  height: 1.5em;
  width: 24px;
  color: var(--sidebar-text-color);
}

.flex-box-v {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.copy-menu {
  flex-grow: 1;
  background-color: var(--sidebar-bg-color);
  border-right: solid 1px var(--menu-board-color);
}
</style>
