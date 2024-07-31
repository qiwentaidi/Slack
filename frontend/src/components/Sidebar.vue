<script setup lang="ts">
import { Refresh, Setting } from "@element-plus/icons-vue";
import { computed, onMounted, ref } from "vue";
import updateUI from "./Update.vue";
import global from "../global";
import menus from "../router/menu";
import { check } from "../util";
onMounted(() => {
  check.client();
  check.poc();
});

const menuStyle = computed(() => {
    return global.Theme.value ? { 
      backgroundColor: '#333333'
    } : {
      backgroundColor: '#F2F3F5' 
    };
})

const updateDialog = ref(false)
</script>

<template>
  <el-menu :collapse="true" router :default-active="$route.path" :style="menuStyle">
    <template v-for="menu in menus">
      <el-menu-item v-if="!menu.children" :index="menu.path">
        <el-icon size="24">
          <svg viewBox="0 0 16 16" fill="" xmlns="http://www.w3.org/2000/svg">
            <path fill=""
              d="M2 13.5V7h1v6.5a.5.5 0 0 0 .5.5h9a.5.5 0 0 0 .5-.5V7h1v6.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5zm11-11V6l-2-2V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5z" />
            <path fill=""
              d="M7.293 1.5a1 1 0 0 1 1.414 0l6.647 6.646a.5.5 0 0 1-.708.708L8 2.207 1.354 8.854a.5.5 0 1 1-.708-.708L7.293 1.5z" />
          </svg>
        </el-icon>
        <template #title><span>{{ $t(menu.name) }}</span></template>
      </el-menu-item>
      <el-sub-menu :index="menu.path" v-else>
        <template #title>
          <el-icon size="24">
            <component :is="menu.icon" />
          </el-icon>
          <span>{{ $t(menu.name) }}</span>
        </template>
        <el-menu-item v-for="item in menu.children" :key="item.path" :index="item.path">{{ $t(item.name)
          }}</el-menu-item>
      </el-sub-menu>
    </template>

    <div style="flex-grow: 1;"></div>

    <el-menu-item class="custom-menu-item" @click="updateDialog = true">
      <el-icon size="24">
        <Refresh />
      </el-icon>
      <template #title><span>{{ $t("aside.update") }}</span></template>
      <el-badge value="new" v-if="global.UPDATE.ClientStatus ||
        global.UPDATE.PocStatus
      " />
    </el-menu-item>

    <el-menu-item class="custom-menu-item" @click="$router.push('/settings')">
      <el-icon size="24">
        <setting />
      </el-icon>
      <template #title><span>{{ $t("aside.setting") }}</span></template>
    </el-menu-item>
  </el-menu>

  <!-- update -->
  <el-dialog v-model="updateDialog" title="更新通知" width="40%">
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
  height: calc(100vh - 35px);
}

.el-menu-item {
  font-size: 16px;
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

.custom-menu-item:hover {
  background-color: var(--sidebar-bg-color);
  color: #3875f6;
}
</style>
