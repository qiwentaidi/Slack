<script setup lang="ts">
import { Refresh, Setting } from "@element-plus/icons-vue";
import updateUI from "./Update.vue";
import global from "@/global/index";
import menus from "@/router/menu";
import router from "@/router";
import pocIcon from '@/assets/icon/pocmanagement.svg'

const bottomControl = [
  {
    label: "aside.update",
    icon: Refresh,
    action: () => {
      global.UPDATE.updateDialog = true
    }
  },
  {
    label: "aside.poc_manage",
    icon: pocIcon,
    action: () => {
      router.push('/PocManagement')
    }
  },
  {
    label: "aside.setting",
    icon: Setting,
    action: () => {
      router.push('/Settings')
    }
  },
]
</script>

<template>
  <el-menu :collapse="true" router :default-active="$route.path">
    <template v-for="menu in menus">
      <el-menu-item v-if="!menu.children" :index="menu.path">
        <el-icon size="24">
          <component :is="menu.icon" />
        </el-icon>
        <template #title><span>{{ $t(menu.name) }}</span></template>
      </el-menu-item>
      <el-sub-menu :index="menu.path" v-else>
        <template #title>
          <el-icon :size="24">
            <component :is="menu.icon" />
          </el-icon>
          <span>{{ $t(menu.name) }}</span>
        </template>
        <el-menu-item v-for="item in menu.children" :key="item.path" :index="menu.path + item.path"><span>{{ $t(item.name)
          }}</span></el-menu-item>
      </el-sub-menu>
    </template>

    <div style="flex-grow: 1;"></div>

    <el-menu-item v-for="(item, index) in bottomControl" @click="item.action">
      <el-icon size="24">
        <component :is="item.icon" />
      </el-icon>
      <template #title><span>{{ $t(item.label) }}</span></template>
      <el-badge value="new" v-if="index == 0 && (global.UPDATE.ClientStatus || global.UPDATE.PocStatus)" />
    </el-menu-item>
  </el-menu>

  <!-- update -->
  <el-dialog v-model="global.UPDATE.updateDialog" title="更新通知" width="40%">
    <updateUI></updateUI>
  </el-dialog>
</template>

<style scoped>
.el-badge {
  margin-bottom: 80px;
  left: 15px;
  position: absolute;
}

.el-menu {
  display: flex;
  flex-direction: column;
  height: calc(100vh - var(--titlebar-height));
  background-color: var(--sidebar-bg-color);
}

.el-menu-item {
  font-size: 16px;
}

.el-menu-item:hover {
  background-color: var(--sidebar-bg-color);
  color: #3875f6;
}

.el-menu-item.is-active {
  color: var(--sidebar-text-color);
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

.el-sub-menu__title {
  font-size: 16px;
} 
</style>
