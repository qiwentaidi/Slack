<script setup lang="ts">
import { OfficeBuilding, Tools, Refresh, Monitor, Smoking, Setting } from "@element-plus/icons-vue";
import { reactive } from "vue";
import { useI18n } from "vue-i18n";
import about from "./About.vue"
import updateUI from "./Update.vue"
import global from "../global"

const { locale } = useI18n()

const changeLanguage = (area: string) => {
  localStorage.setItem('language', area)
  locale.value = area
};

const dg = reactive({
  helpDialogVisible: false,
  updateDialogVisible: false,
  showLogger: false
});

</script>

<template>
  <el-menu class="my-menu" :collapse="true" route active-text-color="#fff" background-color="#F2F3F5" text-color="#000">
    <el-menu-item index="/" @click="$router.push('/')">
      <el-icon>
      <svg class="bi bi-house" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" d="M2 13.5V7h1v6.5a.5.5 0 0 0 .5.5h9a.5.5 0 0 0 .5-.5V7h1v6.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5zm11-11V6l-2-2V2.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5z"/>
        <path fill-rule="evenodd" d="M7.293 1.5a1 1 0 0 1 1.414 0l6.647 6.646a.5.5 0 0 1-.708.708L8 2.207 1.354 8.854a.5.5 0 1 1-.708-.708L7.293 1.5z"/>
      </svg>
    </el-icon>
      <template #title><span>{{ $t('aside.home') }}</span></template>
    </el-menu-item>
    <el-sub-menu index="1">
      <template #title>
        <el-icon>
          <Smoking />
        </el-icon>
        <span>{{ $t('aside.penetration') }}</span>
      </template>
      <el-menu-item index="/Permeation/Webscan" @click="$router.push('/Permeation/Webscan')">{{ $t('aside.webscan')
      }}</el-menu-item>
      <el-menu-item index="/Permeation/Portscan" @click="$router.push('/Permeation/Portscan')">{{ $t('aside.portscan')
      }}</el-menu-item>
      <!-- <el-menu-item index="1-3">漏洞利用</el-menu-item> -->
      <el-menu-item index="/Permeation/Dirsearch" @click="$router.push('/Permeation/Dirsearch')">{{ $t('aside.dirscan')
      }}</el-menu-item>
      <el-menu-item index="/Permeation/Jsfinder" @click="$router.push('/Permeation/Jsfinder')">JSFinder</el-menu-item>
      <el-menu-item index="/Permeation/Pocdetail" @click="$router.push('/Permeation/Pocdetail')">{{
        $t('aside.pocdetail') }}</el-menu-item>
    </el-sub-menu>
    <el-sub-menu index="2">
      <template #title>
        <el-icon>
          <OfficeBuilding />
        </el-icon>
        <span>{{ $t('aside.asset_collection') }}</span>
      </template>
      <el-menu-item index="/Asset/Asset" @click="$router.push('/Asset/Asset')">{{ $t('aside.asset_from_company')
      }}</el-menu-item>
      <el-menu-item index="/Asset/Subdomain" @click="$router.push('/Asset/Subdomain')">{{
        $t('aside.subdomain_brute_force') }}</el-menu-item>
      <el-menu-item index="/Asset/Ipdomain" @click="$router.push('/Asset/Ipdomain')">{{ $t('aside.search_domain_info')
      }}</el-menu-item>
    </el-sub-menu>

    <el-sub-menu index="3">
      <template #title>
        <el-icon>
          <Monitor />
        </el-icon>
        <span>{{ $t('aside.space_engine') }}</span>
      </template>
      <el-menu-item index="/SpaceEngine/Fofa" @click="$router.push('/SpaceEngine/Fofa')">FOFA</el-menu-item>
      <el-menu-item index="/SpaceEngine/Hunter" @click="$router.push('/SpaceEngine/Hunter')">{{ $t('aside.hunter')
      }}</el-menu-item>
      <!-- <el-menu-item index="3-3">{{ $t('aside.360quake') }}</el-menu-item> -->
      <el-menu-item index="/SpaceEngine/AgentPool" @click="$router.push('/SpaceEngine/AgentPool')">{{
        $t('aside.agent_pool') }}</el-menu-item>
    </el-sub-menu>
    <el-sub-menu index="4" style="max-height: 56px;">
      <template #title>
        <el-icon>
          <Tools />
        </el-icon>
        <span>{{ $t('aside.tools') }}</span>
      </template>
      <el-menu-item index="/Tools/Codec" @click="$router.push('/Tools/Codec')">{{ $t('aside.en_and_de')
      }}</el-menu-item>
      <el-menu-item index="/Tools/System" @click="$router.push('/Tools/System')">{{ $t('aside.systeminfo')
      }}</el-menu-item>
      <el-menu-item index="/Tools/Fscan" @click="$router.push('/Tools/Fscan')">{{ $t('aside.fscan') }}</el-menu-item>
      <el-menu-item index="/Tools/Memo" @click="$router.push('/Tools/Memo')">{{ $t('aside.memorandum')
      }}</el-menu-item>
      <el-menu-item index="/Tools/Thinkdict" @click="$router.push('/Tools/Thinkdict')">{{
        $t('aside.associate_dictionary_generator') }}</el-menu-item>
      <el-menu-item index="/Tools/AKSK" @click="$router.push('/Tools/AKSK')">{{ $t('aside.aksk')
      }}</el-menu-item>
    </el-sub-menu>

    <el-menu-item class="custom-menu-item" index="/console" @click="dg.showLogger = true">
      <el-icon>
      <svg class="bi bi-file-earmark-ruled" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" d="M13 9H3V8h10v1zm0 3H3v-1h10v1z"/>
        <path fill-rule="evenodd" d="M5 14V9h1v5H5z"/>
        <path d="M4 1h5v1H4a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V6h1v7a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2z"/>
        <path d="M9 4.5V1l5 5h-3.5A1.5 1.5 0 0 1 9 4.5z"/>
      </svg>
    </el-icon>
      <template #title><span>{{ $t('aside.yx_log') }}</span></template>
    </el-menu-item>

    <el-menu-item class="custom-menu-item" index="/update" @click="dg.updateDialogVisible = true">
      <el-icon>
        <Refresh />
      </el-icon>
      <template #title><span>{{ $t('aside.update') }}</span></template>
      <el-badge is-dot v-if="global.UPDATE.ClientStatus == true || global.UPDATE.PocStatus == true" />
    </el-menu-item>

    <el-menu-item class="custom-menu-item" index="/settings" @click="$router.push('/Settings')">
      <el-icon>
        <setting />
      </el-icon>
      <template #title><span>{{ $t('aside.setting') }}</span></template>
    </el-menu-item>
    <el-sub-menu index="7">
      <template #title>
        <el-icon>
          <svg style="width: 22px;" t="1713366512233" class="icon" viewBox="0 0 1366 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="5910"><path d="M1309.745268 0H56.12297a55.802067 55.802067 0 0 0 0 112.245537h1253.622298a55.802067 55.802067 0 1 0 0-112.245537zM1309.745268 455.877231H56.12297a55.802067 55.802067 0 0 0 0 112.245538h1253.622298a55.802067 55.802067 0 1 0 0-112.245538zM1309.745268 911.754463H56.12297a55.802067 55.802067 0 1 0 0 112.245537h1253.622298a55.802067 55.802067 0 0 0 0-112.245537z" fill="#4D4D4D" p-id="5911"></path></svg>        
        </el-icon>
        <span>{{ $t('aside.more') }}</span>
      </template>
      <el-sub-menu index="language">
        <template #title><span>{{ $t('aside.language') }}</span></template>
        <el-menu-item index="cn" @click="changeLanguage('zh')">{{ $t('aside.zh') }}</el-menu-item>
        <el-menu-item index="en" @click="changeLanguage('en')">{{ $t('aside.en') }}</el-menu-item>
      </el-sub-menu>
      <el-menu-item index="about" @click="dg.helpDialogVisible = true">{{ $t('aside.about') }}</el-menu-item>
    </el-sub-menu>
  </el-menu>

  <!-- update -->
  <el-dialog v-model="dg.updateDialogVisible" title="更新通知" width="40%">
    <updateUI></updateUI>
  </el-dialog>

  <!-- about -->
  <el-dialog v-model="dg.helpDialogVisible" width="36%" center>
    <about></about>
  </el-dialog>

  <!-- running logs -->
  <el-drawer
    v-model="dg.showLogger"
    direction="ltr"
    size="30%"
  >
  <template #title>
    <h4>运行日志</h4>
  </template>
  <el-input class="log-textarea" v-model="global.Logger.value" type="textarea" style="height: 100%;" resize="none"></el-input>
  </el-drawer>
</template>

<style>
.el-badge {
  margin-bottom: 70px;
}

.my-menu {
  display: grid;
  grid-template-rows: auto auto auto auto 1fr;
  height: 100vh;
}

/* 暗色css */
.el-sub-menu.is-active .el-sub-menu__title i {
  color: #3875F6;
}

.el-menu-item.is-active {
  color: #000;
}

.el-menu-item.is-active::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 5px;
  /*  色块的宽度 */
  height: 100%;
  background-color: #3875F6;
  /*  色块的颜色 */
  border-radius: 0 3px 3px 0;
  /* 轨道的形状 */
}

.el-menu-item {
  font-size: 16px
}

.custom-menu-item:hover {
 /* 定义鼠标悬停时的样式 */
 background-color: #f5f5f5; /* 例如，使用一个浅灰色作为背景 */
 color: #3875F6;
}

.custom-menu-item.is-active::before {
 /* 定义激活状态时的样式 */
 background-color: inherit; /* 使用继承的背景颜色 */
 color: inherit; /* 使用继承的文本颜色 */
}

.el-drawer__header {
    margin-bottom: 0px;
}

.el-icon svg {
  height: 1.5em;
  width: 1.5em;
}
</style>
