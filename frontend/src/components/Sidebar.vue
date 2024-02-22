<script setup lang="ts">
import { OfficeBuilding, Tools, Refresh, Monitor, Smoking, Menu, HomeFilled, Setting, Download } from "@element-plus/icons-vue";
import { GoFetch } from "../../wailsjs/go/main/App";
import { CheckFileStat, GetFileContent, UpdatePocFile, UpdateClinetFile, Restart, UserHomeDir, InitConfig } from "../../wailsjs/go/main/File";
import { onMounted, reactive } from "vue";
import { ElNotification, ElMessageBox } from "element-plus";
import { compareVersion } from "../util"
import Loading from "./Loading.vue";
import { useI18n } from "vue-i18n";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";

const { locale } = useI18n()

const changeLanguage = (area: string) => {
  localStorage.setItem('language', area)
  locale.value = area
};

const version = reactive({
  LocalPoc: "",
  LocalClient: '1.4.7',
  RemotePoc: "",
  RemoteClient: "",
  PocUpdateContent: "",
  ClientUpdateContent: "",
  PocStatus: false,
  ClientStatus: false,
  helpDialogVisible: false,
  updateDialogVisible: false,
});

onMounted(async () => {
  window.HomePath = await UserHomeDir()
  window.ConfigPath = "/slack/config"
  window.ActivePathPoc = ConfigPath + "/active-detect"
  window.AFGPathPoc = ConfigPath + "/afrog-pocs"
  window.PocVersion = ConfigPath + "/afrog-pocs/version"
  window.LocalPocVersionFile = HomePath + PocVersion
  window.PortBurstPath = HomePath + "/slack/portburte"
  let cfg = await CheckFileStat(HomePath + "/slack")
  if (!cfg) {
    ElNotification({
      duration: 0,
      message: '未检测到配置文件，正在初始化...',
      icon: Loading,
    });
    if (await InitConfig()) {
      ElNotification.closeAll()
      ElNotification({
        message: "配置文件初始化成功!",
        type: "success",
      });
      check.client()
      check.poc()
    } else {
      ElNotification.closeAll()
      ElNotification({
        title: "无法下载配置文件",
        message: "请自行到https://gitee.com/the-temperature-is-too-low/slack-poc/releases/下载slack.zip并解压到用户根目录!",
        type: "warning",
      });
    }
  } else {
    check.client()
    check.poc()
  }
});

const check = ({
  // poc
  poc: async function () {
    let pcfg = await CheckFileStat(LocalPocVersionFile)
    if (!pcfg) {
      version.LocalPoc = "版本文件不存在"
      version.PocStatus = false
      return
    } else {
      version.LocalPoc = await GetFileContent(LocalPocVersionFile)
    }
    let resp = await GoFetch("GET", download.RemotePocVersion, "", [{}], 10, null)
    if (resp.Error == true) {
      version.PocUpdateContent = "检测更新失败"
      version.PocStatus = false
    } else {
      version.RemotePoc = resp.Body
      if (compareVersion(version.LocalPoc, version.RemotePoc) == -1) {
        version.PocUpdateContent = (await GoFetch("GET", download.PocUpdateCentent, "", [{}], 10, null)).Body
        version.PocStatus = true
      } else {
        version.PocUpdateContent = "已是最新版本"
        version.PocStatus = false
      }
    }
  },
  // client
  client: async function () {
    let resp = await GoFetch("GET", download.RemoteClientVersion, "", [{}], 10, null)
    if (resp.Error == true) {
      version.ClientUpdateContent = "检测更新失败"
      version.ClientStatus = false
    } else {
      version.RemoteClient = resp.Body
      if (compareVersion(version.LocalClient, version.RemoteClient) == -1) {
        version.ClientUpdateContent = (await GoFetch("GET", download.ClientUpdateCentent, "", [{}], 10, null)).Body
        version.ClientStatus = true
      } else {
        version.ClientUpdateContent = "已是最新版本"
        version.ClientStatus = false
      }
    }
  }
})

const update = ({
  poc: async function () {
    let err = await UpdatePocFile()
    if (err == "") {
      ElNotification({
        title: "Success",
        message: "POC更新成功！",
        type: "success",
      });
    } else {
      ElNotification({
        title: "POC更新失败",
        message: err,
        type: "error",
      });
    }
  },
  client: async function () {
    ElNotification({
      title: "提示",
      message: "客户端后台自动下载中~",
      type: "info",
    });
    let err = await UpdateClinetFile("v" + version.RemoteClient)
    if (err == "") {
      ElMessageBox.confirm(
        '更新成功，是否重新启动?',
        {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'success',
          center: true,
        }
      )
        .then(() => {
          Restart()
        })
        .catch((action) => {
          // 在这里处理用户取消或者其他非 "OK" 的选项被点击时的操作
          console.log('User cancelled or chose another option.')
        })
    } else {
      ElNotification({
        title: "Error",
        message: "客户端更新失败！！",
        type: "error",
      });
    }
  }
})

const download = {
  RemotePocVersion:
    "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version",
  RemoteClientVersion:
    "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version",
  PocUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update',
  ClientUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update',
};

</script>

<template>
  <el-menu class="my-menu" :collapse="true" route active-text-color="#fff" background-color="#F2F3F5" text-color="#000">
    <el-menu-item index="/" @click="$router.push('/')">
      <el-icon>
        <HomeFilled />
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
      <el-menu-item index="/Tools/Reverse" @click="$router.push('/Tools/Reverse')">{{ $t('aside.memorandum')
      }}</el-menu-item>
      <el-menu-item index="/Tools/Thinkdict" @click="$router.push('/Tools/Thinkdict')">{{
        $t('aside.associate_dictionary_generator') }}</el-menu-item>
      <el-menu-item index="/Tools/Wxappid" @click="$router.push('/Tools/Wxappid')">{{ $t('aside.wechat_appid')
      }}</el-menu-item>
    </el-sub-menu>

    <el-menu-item class="noactive" index="/update" @click="version.updateDialogVisible = true">
      <el-icon>
        <Refresh />
      </el-icon>
      <template #title><span>{{ $t('aside.update') }}</span></template>
      <el-badge is-dot v-if="version.ClientStatus == true || version.PocStatus == true" />
    </el-menu-item>

    <el-menu-item index="/settings" @click="$router.push('/Settings')">
      <el-icon>
        <setting />
      </el-icon>
      <template #title><span>{{ $t('aside.setting') }}</span></template>
    </el-menu-item>
    <el-sub-menu index="7">
      <template #title>
        <el-icon>
          <Menu />
        </el-icon>
        <span>{{ $t('aside.more') }}</span>
      </template>
      <el-sub-menu index="language">
        <template #title><span>{{ $t('aside.language') }}</span></template>
        <el-menu-item index="cn" @click="changeLanguage('zh')">{{ $t('aside.zh') }}</el-menu-item>
        <el-menu-item index="en" @click="changeLanguage('en')">{{ $t('aside.en') }}</el-menu-item>
      </el-sub-menu>
      <el-menu-item index="about" @click="version.helpDialogVisible = true">{{ $t('aside.about') }}</el-menu-item>
    </el-sub-menu>
  </el-menu>

  <!-- update -->
  <el-dialog v-model="version.updateDialogVisible" title="更新通知" width="40%">
    <el-card class="box-card" style="width: 100%;">
      <template #header>
        <div class="card-header">
          <el-text>
            <span style="font-weight: bold;">POC&指纹{{ version.RemotePoc }}</span>
            <br />当前{{ "v" + version.LocalPoc }}
          </el-text>
          <el-button class="button" :icon="Download" text @click="update.poc"
            v-if="version.PocStatus">立即下载</el-button>
          <span v-else>{{ version.PocUpdateContent }}</span>
        </div>
      </template>
      <el-scrollbar height="100px" style="white-space: pre-wrap;" v-if="version.PocStatus">
        {{ version.PocUpdateContent }}
      </el-scrollbar>
    </el-card>

    <el-card class="box-card" style="width: 100%; margin-top: 10px;">
      <template #header>
        <div class="card-header">
          <el-text>
            <span style="font-weight: bold;">客户端{{ version.RemoteClient }}</span>
            <br />当前{{ "v" + version.LocalClient }}
          </el-text>
          <el-button class="button" :icon="Download" text @click="update.client"
            v-if="version.ClientStatus">立即下载</el-button>
          <span v-else>{{ version.ClientUpdateContent }}</span>
        </div>
      </template>
      <el-scrollbar height="100px" style="white-space: pre-wrap;" v-if="version.ClientStatus">
        {{ version.ClientUpdateContent }}
      </el-scrollbar>
    </el-card>
  </el-dialog>

  <!-- about -->
  <el-dialog v-model="version.helpDialogVisible" width="36%" center>
    <div style="text-align: center;">
      <el-space direction="vertical" :size="8">
        <img src="/slack.svg" style="height: 8em; margin-bottom: 5px;">
        <span style="font-size: 20px; font-weight: bold; color: black;">slack-wails</span>
        <span style="font-weight: bold;">{{ "v" + version.LocalClient }}</span>
        <div>
          <el-link @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack')">{{ $t('aside.source_code') }}</el-link>
          <el-divider direction="vertical" />
          <el-link @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/issues/new')">{{ $t('aside.issue')
          }}</el-link>
          <el-divider direction="vertical" />
          <el-link
            @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97')">{{
              $t('aside.update_log') }}</el-link>
        </div>
      </el-space>
    </div>
  </el-dialog>
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
}</style>
