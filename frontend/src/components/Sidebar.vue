<script setup lang="ts">
import { OfficeBuilding, Tools, Refresh, Monitor, Smoking, Menu, HomeFilled, Setting, Download } from "@element-plus/icons-vue";
import { GoFetch } from "../../wailsjs/go/main/App";
import { CheckFileStat, GetFileContent, UpdatePocFile, UpdateClinetFile, Restart, UserHomeDir, InitConfig } from "../../wailsjs/go/main/File";
import { onMounted, reactive } from "vue";
import { ElNotification, ElMessageBox } from "element-plus";
import { compareVersion } from "../util"
import Loading from "./Loading.vue";
import { useI18n } from "vue-i18n";

const { locale } = useI18n()

const changeEN = () => {
  localStorage.setItem('language', 'en')
  locale.value = 'en'

};
const changeZH = () => {
  localStorage.setItem('language', 'zh')
  locale.value = 'zh'
};

onMounted(async () => {
  // 初始赋值
  window.ActivePathPoc = "/slack/active-detect"
  window.AFGPathPoc = "/slack/afrog-pocs"
  window.PocVersion = "/slack/afrog-pocs/version"
  let home = await UserHomeDir()
  let cfg = await CheckFileStat(home + "/slack")
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
      window.LocalPocVersion = await UserHomeDir() + window.PocVersion
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
    window.LocalPocVersion = await UserHomeDir() + window.PocVersion
    check.client()
    check.poc()
  }
});

const check = ({
  // poc
  poc: async function () {
    let pcfg = await CheckFileStat(window.LocalPocVersion)
    if (!pcfg) {
      version.LocalPoc = "版本文件不存在"
      version.PocStatus = false
      return
    } else {
      version.LocalPoc = await GetFileContent(window.LocalPocVersion)
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
        version.PocUpdateContent = "当前已是最新版本"
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
        version.ClientUpdateContent = "当前已是最新版本"
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

const version = reactive({
  LocalPoc: "",
  LocalClient: "1.4.6",
  RemotePoc: "",
  RemoteClient: "",
  PocUpdateContent: "",
  ClientUpdateContent: "",
  PocStatus: false,
  ClientStatus: false,
  helpDialogVisible: false,
  updateDialogVisible: false,
});

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
  <div>
    <el-menu :collapse="true" route style="height: 100vh;">
      <el-menu-item index="/" @click="$router.push('/')">
        <el-icon>
          <HomeFilled />
        </el-icon>
        <template #title><span>{{ $t('aside.home') }}</span></template>
      </el-menu-item>
      <el-sub-menu index="1">
        <template #title><span>{{ $t('aside.penetration_test') }}</span><el-icon>
            <Smoking />
          </el-icon></template>
        <el-menu-item index="/Permeation/Webscan"
          @click="$router.push('/Permeation/Webscan')">{{ $t('aside.webscan') }}</el-menu-item>
        <el-menu-item index="/Permeation/Portscan"
          @click="$router.push('/Permeation/Portscan')">{{ $t('aside.portscan') }}</el-menu-item>
        <!-- <el-menu-item index="1-3">漏洞利用</el-menu-item> -->
        <el-menu-item index="/Permeation/Dirsearch"
          @click="$router.push('/Permeation/Dirsearch')">{{ $t('aside.dirscan') }}</el-menu-item>
        <el-menu-item index="/Permeation/Pocdetail"
          @click="$router.push('/Permeation/Pocdetail')">{{ $t('aside.pocdetail') }}</el-menu-item>
      </el-sub-menu>
      <el-sub-menu index="2">
        <template #title><span>{{ $t('aside.asset_collection') }}</span><el-icon>
            <OfficeBuilding />
          </el-icon></template>
        <el-menu-item index="/Asset/Asset"
          @click="$router.push('/Asset/Asset')">{{ $t('aside.asset_from_company') }}</el-menu-item>
        <el-menu-item index="/Asset/Subdomain"
          @click="$router.push('/Asset/Subdomain')">{{ $t('aside.subdomain_brute_force') }}</el-menu-item>
        <el-menu-item index="/Asset/Ipdomain"
          @click="$router.push('/Asset/Ipdomain')">{{ $t('aside.search_domain_info') }}</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="3">
        <template #title><span>{{ $t('aside.space_engine') }}</span><el-icon>
            <Monitor />
          </el-icon></template>
        <el-menu-item index="/SpaceEngine/Fofa" @click="$router.push('/SpaceEngine/Fofa')">FOFA</el-menu-item>
        <el-menu-item index="/SpaceEngine/Hunter"
          @click="$router.push('/SpaceEngine/Hunter')">{{ $t('aside.hunter') }}</el-menu-item>
        <!-- <el-menu-item index="3-3">{{ $t('aside.360quake') }}</el-menu-item> -->
        <el-menu-item index="/SpaceEngine/AgentPool"
          @click="$router.push('/SpaceEngine/AgentPool')">{{ $t('aside.agent_pool') }}</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="4">
        <template #title><span>{{ $t('aside.tools') }}</span><el-icon>
            <Tools />
          </el-icon></template>
        <el-menu-item index="/Tools/Codec" @click="$router.push('/Tools/Codec')">{{ $t('aside.en_and_de') }}</el-menu-item>
        <el-menu-item index="/Tools/System"
          @click="$router.push('/Tools/System')">{{ $t('aside.systeminfo') }}</el-menu-item>
        <el-menu-item index="/Tools/Fscan" @click="$router.push('/Tools/Fscan')">{{ $t('aside.fscan') }}</el-menu-item>
        <el-menu-item index="/Tools/Reverse"
          @click="$router.push('/Tools/Reverse')">{{ $t('aside.memorandum') }}</el-menu-item>
        <el-menu-item index="/Tools/Thinkdict"
          @click="$router.push('/Tools/Thinkdict')">{{ $t('aside.associate_dictionary_generator') }}</el-menu-item>
        <el-menu-item index="/Tools/Wxappid"
          @click="$router.push('/Tools/Wxappid')">{{ $t('aside.wechat_appid') }}</el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>

  <div class="bottom-align">
    <el-menu :collapse="true" route>
      <el-menu-item index="5" @click="version.updateDialogVisible = true">
        <el-icon>
          <Refresh />
        </el-icon>
        <template #title><span>{{ $t('aside.update') }}</span></template>
        <el-badge is-dot v-if="version.ClientStatus == true || version.PocStatus == true" />
      </el-menu-item>
      <el-menu-item index="/Settings" @click="$router.push('/Settings')">
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
          <el-menu-item index="cn" @click="changeZH">{{ $t('aside.zh') }}</el-menu-item>
          <el-menu-item index="en" @click="changeEN">{{ $t('aside.en') }}</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="about" @click="version.helpDialogVisible = true">{{ $t('aside.about') }}</el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>

  <!-- 更新界面 -->
  <el-dialog v-model="version.updateDialogVisible" title="更新通知" width="50%">
    <el-card class="box-card" shadow="never" v-if="version.PocStatus" style="width: 100%;">
      <template #header>
        <div class="card-header">
          POC&指纹{{ version.RemotePoc }} <br /> 当前{{ "v" + version.LocalPoc }}
          <el-button class="button" :icon="Download" type="primary" :disabled="!version.PocStatus"
            @click="update.poc">立即下载</el-button>
        </div>
      </template>
      <el-input type="textarea" rows="5" v-model="version.PocUpdateContent" resize="none" readonly></el-input>
    </el-card>
    <h3 v-else>当前POC已是最新版本{{ "v" + version.RemotePoc + " :)" }}</h3>
    <el-card class="box-card" shadow="never" style="width: 100%" v-if="version.ClientStatus">
      <template #header>
        <div class="card-header">
          客户端{{ version.RemoteClient }} <br /> 当前{{ "v" + version.LocalClient }}
          <el-button class="button" :icon="Download" type="primary" :disabled="!version.ClientStatus"
            @click="update.client">立即下载</el-button>
        </div>
      </template>
      <el-input type="textarea" rows="5" v-model="version.ClientUpdateContent" resize="none" readonly></el-input>
    </el-card>
    <h3 v-else>当前客户端已是最新版本{{ "v" + version.RemoteClient + ":)" }}</h3>
  </el-dialog>

  <!-- 弹窗界面 -->
  <el-dialog v-model="version.helpDialogVisible" :title="$t('aside.about')" width="36%" center>
    <h3>{{ $t('aside.suggestions') }}</h3>
    <br />
    <h4>{{ $t('aside.technology') }}: Vue + Typescript + Vite + Wails + Go</h4>
  </el-dialog>
</template>

<style>
.setting {
  flex-grow: 1;
}

.el-badge {
  margin-bottom: 70px;
}

.bottom-align {
  width: 100%;
  position: absolute;
  bottom: 0;
}
</style>
