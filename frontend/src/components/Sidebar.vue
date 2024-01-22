<script setup lang="ts">
import {
  OfficeBuilding,
  Tools,
  Refresh,
  Monitor,
  Smoking,
  Help,
  HomeFilled,
  Setting,
  Download,
} from "@element-plus/icons-vue";
import { GoFetch } from "../../wailsjs/go/main/App";
import {
  CheckFileStat,
  GetFileContent,
  UpdatePocFile,
  UpdateClinetFile,
  Restart,
  UserHomeDir,
  InitConfig
} from "../../wailsjs/go/main/File";
import { onMounted, reactive } from "vue";
import { ElNotification, ElMessageBox } from "element-plus";
import { compareVersion } from "../util"
import Loading from "./Loading.vue";

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
      version.LocalPoc = "v" + await GetFileContent(window.LocalPocVersion)
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
  <el-menu :collapse="true" route style="height: 100%;">
    <el-menu-item index="/" @click="$router.push('/')">
      <el-icon>
        <HomeFilled />
      </el-icon>
      <template #title><span>主页</span></template>
    </el-menu-item>
    <el-sub-menu index="1">
      <template #title><span>渗透测试</span><el-icon>
          <Smoking />
        </el-icon></template>
      <el-menu-item index="/Permeation/Webscan" @click="$router.push('/Permeation/Webscan')">网站扫描</el-menu-item>
      <el-menu-item index="/Permeation/Portscan" @click="$router.push('/Permeation/Portscan')">主机扫描</el-menu-item>
      <!-- <el-menu-item index="1-3">漏洞利用</el-menu-item> -->
      <el-menu-item index="/Permeation/Dirsearch" @click="$router.push('/Permeation/Dirsearch')">目录扫描</el-menu-item>
      <el-menu-item index="/Permeation/Postman" @click="$router.push('/Permeation/Postman')">Postman(未完成)</el-menu-item>
      <el-menu-item index="/Permeation/Pocdetail" @click="$router.push('/Permeation/Pocdetail')">漏洞详情</el-menu-item>
    </el-sub-menu>
    <el-sub-menu index="2">
      <template #title><span>资产收集</span><el-icon>
          <OfficeBuilding />
        </el-icon></template>
      <el-menu-item index="/Asset/Asset" @click="$router.push('/Asset/Asset')">公司名称查资产</el-menu-item>
      <el-menu-item index="/Asset/Subdomain" @click="$router.push('/Asset/Subdomain')">子域名暴破</el-menu-item>
      <el-menu-item index="/Asset/Ipdomain" @click="$router.push('/Asset/Ipdomain')">域名信息查询</el-menu-item>
    </el-sub-menu>

    <el-sub-menu index="3">
      <template #title><span>空间引擎</span><el-icon>
          <Monitor />
        </el-icon></template>
      <el-menu-item index="/SpaceEngine/Fofa" @click="$router.push('/SpaceEngine/Fofa')">FOFA</el-menu-item>
      <el-menu-item index="/SpaceEngine/Hunter" @click="$router.push('/SpaceEngine/Hunter')">鹰图</el-menu-item>
      <el-menu-item index="3-3">360夸克(没做)</el-menu-item>
    </el-sub-menu>

    <el-sub-menu index="4">
      <template #title><span>小工具</span><el-icon>
          <Tools />
        </el-icon></template>
      <el-menu-item index="/Tools/Codec" @click="$router.push('/Tools/Codec')">加解密模块</el-menu-item>
      <el-menu-item index="/Tools/System" @click="$router.push('/Tools/System')">杀软识别/提权补丁</el-menu-item>
      <el-menu-item index="/Tools/Fscan" @click="$router.push('/Tools/Fscan')">Fscan内容提取</el-menu-item>
      <el-menu-item index="/Tools/Reverse" @click="$router.push('/Tools/Reverse')">反弹shell备忘录</el-menu-item>
      <el-menu-item index="/Tools/Thinkdict" @click="$router.push('/Tools/Thinkdict')">联想字典生成器</el-menu-item>
      <el-menu-item index="/Tools/Wxappid" @click="$router.push('/Tools/Wxappid')">微信AppId校验</el-menu-item>
    </el-sub-menu>
    <div class="bottom-align">
      <el-menu-item index="5" @click="version.updateDialogVisible = true">
        <el-icon>
          <Refresh />
        </el-icon>
        <template #title><span>更新</span></template>
        <el-badge is-dot v-if="version.ClientStatus == true || version.PocStatus == true" />
      </el-menu-item>
      <el-menu-item index="/Settings" @click="$router.push('/Settings')">
        <el-icon>
          <setting />
        </el-icon>
        <template #title><span>设置</span></template>
      </el-menu-item>

      <el-menu-item index="7" @click="version.helpDialogVisible = true">
        <el-icon>
          <Help />
        </el-icon>
        <template #title><span>关于</span></template>
      </el-menu-item>
    </div>
    <el-dialog v-model="version.updateDialogVisible" title="更新通知" width="50%">
      <el-card class="box-card" shadow="never" v-if="version.PocStatus" style="width: 100%;">
        <template #header>
          <div class="card-header">
            <el-text>POC&指纹{{ version.RemotePoc }} <br /> 当前{{ version.LocalPoc }}</el-text>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.PocStatus"
              @click="update.poc">立即下载</el-button>
          </div>
        </template>
        <el-input type="textarea" rows="5" v-model="version.PocUpdateContent" resize="none" readonly></el-input>
      </el-card>
      <h3 v-else>当前POC已是最新版本v{{ version.RemotePoc }} :)</h3>
      <el-card class="box-card" shadow="never" style="width: 100%" v-if="version.ClientStatus">
        <template #header>
          <div class="card-header">
            <el-text>客户端{{ version.RemoteClient }} <br /> 当前v{{ version.LocalClient }}</el-text>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.ClientStatus"
              @click="update.client">立即下载</el-button>
          </div>
        </template>
        <el-input type="textarea" rows="5" v-model="version.ClientUpdateContent" resize="none" readonly></el-input>
      </el-card>
      <h3 v-else>当前客户端已是最新版本v{{ version.RemoteClient }} :)</h3>
    </el-dialog>

    <!-- 弹窗界面 -->
    <el-dialog v-model="version.helpDialogVisible" title="关于" width="36%" center>
      <h4>
        工具目前存在内存GC问题，如有改善意见或其他问题可以通过vx或者issue联系，联系方式可点击首页LOGO处前往项目地址获取
      </h4>
      <h4>前端: Vue + Typescript + Vite + Element-Plus</h4>
      <h4>后端: Wails + Go</h4>
    </el-dialog>
  </el-menu>
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
