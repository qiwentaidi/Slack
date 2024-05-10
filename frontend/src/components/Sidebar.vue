<script setup lang="ts">
import { Refresh, Setting } from "@element-plus/icons-vue";
import { reactive, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import about from "./About.vue";
import updateUI from "./Update.vue";
import global from "../global";
import { check } from "../util";
onMounted(() => {
  check.client();
  check.poc();
});


const { locale } = useI18n();

const changeLanguage = (area: string) => {
  localStorage.setItem("language", area);
  locale.value = area;
};

const dg = reactive({
  helpDialogVisible: false,
  updateDialogVisible: false,
});
</script>

<template>
  <div class="flex-box-v">
    <el-menu :collapse="true" :router="true" :default-active="$route.path" active-text-color="#fff"
      background-color="#F2F3F5" text-color="#000">
      <el-menu-item index="/">
        <el-icon>
          <svg class="bi bi-house" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd"
              d="M2 13.5V7h1v6.5a.5.5 0 0 0 .5.5h9a.5.5 0 0 0 .5-.5V7h1v6.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5zm11-11V6l-2-2V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5z" />
            <path fill-rule="evenodd"
              d="M7.293 1.5a1 1 0 0 1 1.414 0l6.647 6.646a.5.5 0 0 1-.708.708L8 2.207 1.354 8.854a.5.5 0 1 1-.708-.708L7.293 1.5z" />
          </svg>
        </el-icon>
        <template #title><span>{{ $t("aside.home") }}</span></template>
      </el-menu-item>
      <el-sub-menu index="1">
        <template #title>
          <el-icon>
            <Smoking />
          </el-icon>
          <span>{{ $t("aside.penetration") }}</span>
        </template>
        <el-menu-item index="/Permeation/Webscan">{{
      $t("aside.webscan")
    }}</el-menu-item>
        <el-menu-item index="/Permeation/Portscan">{{
        $t("aside.portscan")
      }}</el-menu-item>
        <el-menu-item index="/Permeation/Dirsearch">{{
        $t("aside.dirscan")
      }}</el-menu-item>
        <el-menu-item index="/Permeation/Jsfinder">JSFinder</el-menu-item>
        <!-- <el-menu-item index="/Permeation/Pocdetail">{{
        $t("aside.pocdetail")
      }}</el-menu-item> -->
      </el-sub-menu>
      <el-sub-menu index="2">
        <template #title>
          <el-icon>
            <OfficeBuilding />
          </el-icon>
          <span>{{ $t("aside.asset_collection") }}</span>
        </template>
        <el-menu-item index="/Asset/Company">{{
      $t("aside.asset_from_company")
    }}</el-menu-item>
        <el-menu-item index="/Asset/Subdomain">{{
        $t("aside.subdomain_brute_force")
      }}</el-menu-item>
        <el-menu-item index="/Asset/Ipdomain">{{
        $t("aside.search_domain_info")
      }}</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="3">
        <template #title>
          <el-icon>
            <Monitor />
          </el-icon>
          <span>{{ $t("aside.space_engine") }}</span>
        </template>
        <el-menu-item index="/SpaceEngine/FOFA">FOFA</el-menu-item>
        <el-menu-item index="/SpaceEngine/Hunter">{{
      $t("aside.hunter")
    }}</el-menu-item>
        <el-menu-item index="/SpaceEngine/AgentPool">{{
        $t("aside.agent_pool")
      }}</el-menu-item>
      </el-sub-menu>
      <el-sub-menu index="4" style="max-height: 56px">
        <template #title>
          <el-icon>
            <Tools />
          </el-icon>
          <span>{{ $t("aside.tools") }}</span>
        </template>
        <el-menu-item index="/Tools/Codec">{{
      $t("aside.en_and_de")
    }}</el-menu-item>
        <el-menu-item index="/Tools/System">{{
        $t("aside.systeminfo")
      }}</el-menu-item>
        <el-menu-item index="/Tools/DataHanding">{{
        $t("aside.data_handing")
      }}</el-menu-item>
        <el-menu-item index="/Tools/Memo">{{
        $t("aside.memorandum")
      }}</el-menu-item>
        <el-menu-item index="/Tools/Thinkdict">{{
        $t("aside.associate_dictionary_generator")
      }}</el-menu-item>
        <el-menu-item index="/Tools/AKSK">{{ $t("aside.aksk") }}</el-menu-item>
      </el-sub-menu>
    </el-menu>
    <div class="copy-menu"></div>
    <el-menu :collapse="true" active-text-color="#fff" background-color="#F2F3F5" text-color="#000">
      <el-menu-item class="custom-menu-item" index="/update" @click="dg.updateDialogVisible = true">
        <el-icon>
          <Refresh />
        </el-icon>
        <template #title><span>{{ $t("aside.update") }}</span></template>
        <el-badge is-dot v-if="global.UPDATE.ClientStatus == true ||
      global.UPDATE.PocStatus == true
      " />
      </el-menu-item>

      <el-menu-item class="custom-menu-item" index="/settings" @click="$router.push('/settings')">
        <el-icon>
          <setting />
        </el-icon>
        <template #title><span>{{ $t("aside.setting") }}</span></template>
      </el-menu-item>
      <el-sub-menu index="7">
        <template #title>
          <el-icon>
            <svg style="width: 22px" t="1713366512233" class="icon" viewBox="0 0 1366 1024" version="1.1"
              xmlns="http://www.w3.org/2000/svg" p-id="5910">
              <path
                d="M1309.745268 0H56.12297a55.802067 55.802067 0 0 0 0 112.245537h1253.622298a55.802067 55.802067 0 1 0 0-112.245537zM1309.745268 455.877231H56.12297a55.802067 55.802067 0 0 0 0 112.245538h1253.622298a55.802067 55.802067 0 1 0 0-112.245538zM1309.745268 911.754463H56.12297a55.802067 55.802067 0 1 0 0 112.245537h1253.622298a55.802067 55.802067 0 0 0 0-112.245537z"
                fill="#4D4D4D" p-id="5911"></path>
            </svg>
          </el-icon>
          <span>{{ $t("aside.more") }}</span>
        </template>
        <el-sub-menu index="language">
          <template #title><span>{{ $t("aside.language") }}</span></template>
          <el-menu-item index="cn" @click="changeLanguage('zh')">{{
      $t("aside.zh")
    }}</el-menu-item>
          <el-menu-item index="en" @click="changeLanguage('en')">{{
      $t("aside.en")
    }}</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="about" @click="dg.helpDialogVisible = true">{{
      $t("aside.about")
    }}</el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>

  <!-- update -->
  <el-dialog v-model="dg.updateDialogVisible" title="更新通知" width="40%">
    <updateUI></updateUI>
  </el-dialog>

  <!-- about -->
  <el-dialog v-model="dg.helpDialogVisible" width="36%" center>
    <about></about>
  </el-dialog>
</template>

<style>
.el-badge {
  margin-bottom: 70px;
}

.el-sub-menu.is-active .el-sub-menu__title i {
  color: #3875f6;
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

.el-drawer__header {
  margin-bottom: 0px;
}

.el-icon svg {
  height: 1.5em;
  width: 1.5em;
}

.flex-box-v {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.copy-menu {
  flex-grow: 1;
  background-color: #f2f3f5;
  border-right: solid 1px rgb(220, 223, 230);
}
</style>
