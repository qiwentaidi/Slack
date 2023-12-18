<template>
  <el-menu index="0" mode="horizontal" :ellipsis="false" route text-color="#000">
    <el-menu-item @click="currentComponent = 'Introduction'">
      <el-icon>
        <HomeFilled />
      </el-icon>
    </el-menu-item>
    <el-divider class="divder" direction="vertical" />
    <el-sub-menu index="1">
      <template #title><span>渗透测试</span><el-icon>
          <Smoking />
        </el-icon></template>
      <el-menu-item index="1-1" @click="currentComponent = 'Webscan'">网站扫描</el-menu-item>
      <el-menu-item index="1-2" @click="currentComponent = 'Portscan'">主机扫描</el-menu-item>
      <!-- <el-menu-item index="1-3">漏洞利用</el-menu-item> -->
      <el-menu-item index="1-4" @click="currentComponent = 'Dirsearch'">目录扫描</el-menu-item>
      <!-- <el-menu-item index="1-5" @click="currentComponent = 'Postman'">Postman</el-menu-item> -->
      <el-menu-item index="1-6" @click="currentComponent = 'Pocdetail'">漏洞详情</el-menu-item>
    </el-sub-menu>
    <el-sub-menu index="2">
      <template #title><span>资产收集</span><el-icon>
          <OfficeBuilding />
        </el-icon></template>
      <el-menu-item index="2-1" @click="currentComponent = 'Asset'">公司名称查资产</el-menu-item>
      <el-menu-item index="2-2" @click="currentComponent = 'Subdomain'">子域名暴破</el-menu-item>
      <el-menu-item index="2-3" @click="currentComponent = 'Ipdomain'">域名信息查询</el-menu-item>
    </el-sub-menu>

    <el-sub-menu index="3">
      <template #title><span>空间引擎</span><el-icon>
          <Monitor />
        </el-icon></template>
      <el-menu-item index="3-1" @click="currentComponent = 'Fofa'">FOFA</el-menu-item>
      <el-menu-item index="3-2">鹰图</el-menu-item>
      <el-menu-item index="3-3">360夸克</el-menu-item>
    </el-sub-menu>

    <el-sub-menu index="4">
      <template #title><span>小工具</span><el-icon>
          <Tools />
        </el-icon></template>
      <el-menu-item index="4-1" @click="currentComponent = 'Codec'">加解密模块</el-menu-item>
      <el-menu-item index="4-2" @click="currentComponent = 'System'">杀软识别/提权补丁</el-menu-item>
      <el-menu-item index="4-3" @click="currentComponent = 'Fscan'">Fscan内容提取</el-menu-item>
      <el-menu-item index="4-4" @click="currentComponent = 'Reverse'">反弹shell备忘录</el-menu-item>
      <el-menu-item index="4-5" @click="currentComponent = 'Thinkdict'">联想字典生成器</el-menu-item>
      <el-menu-item index="4-6" @click="currentComponent = 'Wxappid'">微信AppId校验</el-menu-item>
    </el-sub-menu>
    <!-- 右对齐 -->
    <div class="setting" style="--wails-draggable: drag" @dblclick="ToggleMaximise"></div>

    <el-popover placement="bottom-end" title="更新通知" :width="350" trigger="hover">
      <template #reference>
        <el-menu-item index="5">
          <el-icon>
            <Bell />
          </el-icon>
        </el-menu-item>
      </template>
      <el-card class="box-card" shadow="never" style="width: 325px">
        <template #header>
          <div class="card-header">
            <span>POC&指纹-v{{ version.LocalPoc }}</span>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.PocStatus"
              @click="">立即下载</el-button>
          </div>
        </template>
        <el-input type="textarea" rows="5" v-model="version.PocUpdateContent" resize="none" readonly></el-input>
      </el-card>
      <el-card class="box-card" shadow="never" style="width: 325px; margin-top: 15px">
        <template #header>
          <div class="card-header">
            <span>客户端-v{{ version.LocalClient }}</span>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.ClientStatus"
              @click="">立即下载</el-button>
          </div>
        </template>
        <el-input type="textarea" rows="5" v-model="version.ClientUpdateContent" resize="none" readonly></el-input>
      </el-card>
    </el-popover>

    <el-divider class="divder" direction="vertical" />
    <el-tooltip content="设置" placement="bottom">
      <el-menu-item index="6" @click="currentComponent = 'Settings'">
        <el-icon>
          <setting />
        </el-icon>
      </el-menu-item>
    </el-tooltip>
    <el-tooltip content="帮助" placement="bottom">
      <el-menu-item index="7" @click="centerDialogVisible = true">
        <el-icon>
          <Help />
        </el-icon>
      </el-menu-item>
    </el-tooltip>
    <el-divider class="divder" direction="vertical" />
    <!-- MAX MIN CLOSE -->
    <el-menu-item @click="Minimise">
      <el-icon>
        <Minus />
      </el-icon>
    </el-menu-item>
    <el-menu-item @click="ToggleMaximise">
      <el-icon>
        <FullScreen />
      </el-icon>
    </el-menu-item>
    <el-menu-item @click="Quit">
      <el-icon>
        <Close />
      </el-icon>
    </el-menu-item>
    <!-- 弹窗界面 -->
    <el-dialog v-model="centerDialogVisible" title="关于" width="36%" center>
      <h4>
        工具目前存在内存GC问题，如有改善意见或其他问题可以通过下方vx联系，项目地址可点击首页LOGO处前往
      </h4>
      <img src="/wechat.png" class="center" />
    </el-dialog>
  </el-menu>
</template>

<script setup lang="ts">
import { ref, inject, reactive } from "vue";
import {
  OfficeBuilding,
  Tools,
  Minus,
  FullScreen,
  Monitor,
  Smoking,
  Help,
  HomeFilled,
  Close,
  Setting,
  Bell,
  Download,
} from "@element-plus/icons-vue";
import {
  Quit,
  Minimise,
  ToggleMaximise,
  CheckFileStat,
  GetFileContent,
  GoSimpleFetch,
  UpdatePocFile
} from "../../wailsjs/go/main/App";
import { onMounted } from "vue";
import { ElNotification } from "element-plus";
import { compareVersion } from "../util"
const lv = "./config/afrog-pocs/version"
onMounted(async () => {
  if (!CheckFileStat("./config")) {
    ElNotification({
      title: "Warning",
      message: "config配置文件目录加载失败，会影响程序功能使用",
      type: "warning",
    });
    return
  }
  check.poc()
  check.client()
});

const centerDialogVisible = ref(false);
const currentComponent = inject("currentComponent");

const check = ({
  // poc
  poc: async function () {
    if (!CheckFileStat(lv)) {
      version.LocalPoc = "版本文件不存在"
      version.PocStatus = false
      return
    } else {
      version.LocalPoc = await GetFileContent(lv)
    }
    let rp1 = await GoSimpleFetch(download.RemotePocV)
    if (rp1.Status !== 200) {
      version.PocUpdateContent = "检测更新失败"
      version.PocStatus = false
    } else {
      if (compareVersion(version.LocalPoc, rp1.Text) == -1) {
        version.PocUpdateContent = (await GoSimpleFetch(download.RemotePocCentent)).Text
        version.PocStatus = true
      } else {
        version.PocUpdateContent = "当前已是最新版本"
        version.PocStatus = false
      }
    }
  },
  // client
  client: async function () {
    let rp2 = await GoSimpleFetch(download.RemoteClientV)
    if (rp2.Status !== 200) {
      version.ClientUpdateContent = "检测更新失败"
      version.ClientStatus = false
    } else {
      if (compareVersion(version.LocalClient, rp2.Text) == -1) {
        version.ClientUpdateContent = (await GoSimpleFetch(download.RemoteClientCentent)).Text
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
    let err = await UpdatePocFile("v" + version.RemotePoc)
    if (err !== "") {
      ElNotification({
        title: "Success",
        message: "POC更新成功！",
        type: "success",
      });
    }else {
      ElNotification({
        title: "Error",
        message: "POC更新失败！！",
        type: "error",
      });
    }
  },
  client: function () {

  } 
})

const version = reactive({
  LocalPoc: "",
  LocalClient: "1.4.5",
  RemotePoc: "",
  RemoteClient: "",
  PocUpdateContent: "",
  ClientUpdateContent: "",
  PocStatus: false,
  ClientStatus: false,
});

const download = {
  RemotePocV:
    "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version",
  RemoteClientV:
    "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version",
  RemotePocCentent: 'https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update',
  RemoteClientCentent: 'https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update',
};

</script>
<style>
.el-menu.el-menu--horizontal {
  height: 40px;
  margin-left: -18px;
  margin-right: -18px;
}

.el-menu-item {
  padding: 8px;
}

.setting {
  flex-grow: 1;
}

.el-menu-item.is-active.el-tooltip__trigger {
  border-bottom: #fff;
}
</style>
