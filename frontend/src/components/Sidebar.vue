<script setup lang="ts">
import { Refresh, Setting } from "@element-plus/icons-vue";
import updateUI from "./Update.vue";
import global from "@/stores";
import menus from "@/router/menu";
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
    path: "/PocManagement",
  },
  {
    label: "aside.setting",
    icon: Setting,
    path: "/Settings"
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
        <el-menu-item class="react-menu-item" v-for="item in menu.children" :key="item.path" :index="menu.path + item.path"><span>{{
          $t(item.name)
            }}</span></el-menu-item>
      </el-sub-menu>
    </template>

    <div style="flex-grow: 1;"></div>

    <el-menu-item v-for="(item, index) in bottomControl" :index="item.path!" @click="item.action">
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

.react-menu-item {
  font-size: 16px;
  align-items: center;
  display: flex;
  margin: 3%;
}

.el-menu-item:hover {
  background-color: var(--sidebar-bg-color);
  border-radius: 10px;
}

.el-menu-item.is-active {
  background-color: var(--sidebar-bg-color);
  border-radius: 10px;
}
</style>
