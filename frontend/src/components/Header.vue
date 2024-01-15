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
      <el-menu-item index="1-5" @click="currentComponent = 'Postman'">Postman(未完成)</el-menu-item>
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
      <el-menu-item index="3-2" @click="currentComponent = 'Hunter'">鹰图</el-menu-item>
      <el-menu-item index="3-3">360夸克(没做)</el-menu-item>
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
    <div class="setting" style="--wails-draggable: drag" @dblclick="WindowToggleMaximise"></div>

    <el-popover placement="bottom-end" title="更新通知" :width="350" trigger="hover">
      <template #reference>
        <el-menu-item index="5">
          <el-icon>
            <Bell />
          </el-icon>
          <el-badge is-dot v-if="version.ClientStatus == true || version.PocStatus == true" />
        </el-menu-item>
      </template>
      <el-card class="box-card" shadow="never" style="width: 325px" v-if="version.PocStatus">
        <template #header>
          <div class="card-header">
            <el-text>POC&指纹{{ version.RemotePoc }} <br /> 当前{{ version.LocalPoc }}</el-text>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.PocStatus"
              @click="update.poc">立即下载</el-button>
          </div>
        </template>
        <el-input type="textarea" rows="5" v-model="version.PocUpdateContent" resize="none" readonly></el-input>
      </el-card>
      <el-card class="box-card" shadow="never" style="width: 325px; margin-top: 15px" v-if="version.ClientStatus">
        <template #header>
          <div class="card-header">
            <el-text>客户端{{ version.RemoteClient }} <br /> 当前v{{ version.LocalClient }}</el-text>
            <el-button class="button" :icon="Download" type="primary" :disabled="!version.ClientStatus"
              @click="update.client">立即下载</el-button>
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
    <el-tooltip content="关于" placement="bottom">
      <el-menu-item index="7" @click="centerDialogVisible = true">
        <el-icon>
          <Help />
        </el-icon>
      </el-menu-item>
    </el-tooltip>
    <el-divider class="divder" direction="vertical" />
    <!-- MAX MIN CLOSE -->
    <el-menu-item @click="WindowMinimise">
      <el-icon>
        <Minus />
      </el-icon>
    </el-menu-item>
    <el-menu-item @click="WindowToggleMaximise" style="width: 45px;">
      <el-image v-if="!isMax" style="height: 18px; width: 18px;" src='/max.svg' />
      <el-image v-else style="height: 18px; width: 18px;" src='/reduction.svg' />
    </el-menu-item>
    <el-menu-item @click="Quit">
      <el-icon>
        <Close />
      </el-icon>
    </el-menu-item>
    <!-- 弹窗界面 -->
    <el-dialog v-model="centerDialogVisible" title="关于" width="36%" center>
      <h4>
        工具目前存在内存GC问题，如有改善意见或其他问题可以通过vx或者issue联系，联系方式可点击首页LOGO处前往项目地址获取
      </h4>
      <h4>前端: Vue + Typescript + Vite + Element-Plus</h4>
      <h4>后端: Wails + Go</h4>
    </el-dialog>
  </el-menu>
</template>

<script setup lang="ts">
import { ref, inject, reactive } from "vue";
import {
  OfficeBuilding,
  Tools,
  Minus,
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
  CheckFileStat,
  GetFileContent,
  GoFetch,
  UpdatePocFile,
  UpdateClinetFile,
  Restart
} from "../../wailsjs/go/main/App";
import { Quit, WindowMinimise, WindowToggleMaximise } from "../../wailsjs/runtime"
import { onMounted } from "vue";
import { ElNotification, ElMessageBox } from "element-plus";
import { compareVersion } from "../util"
const lv = "./config/afrog-pocs/version"
onMounted(async () => {
  // 处理窗口大小改变的逻辑
  window.addEventListener('resize', () => {
    if (screen.availWidth == window.innerWidth && screen.availHeight == window.innerHeight) {
      isMax.value = true
    } else {
      isMax.value = false
    }
  });
  check.client()
  let cfg = await CheckFileStat("./config")
  if (!cfg) {
    ElNotification({
      title: "Warning",
      message: "config配置文件目录加载失败，会影响程序功能使用",
      type: "warning",
    });
  } else {
    check.poc()
  }
});

const centerDialogVisible = ref(false);
const currentComponent = inject("currentComponent");
const defultreqHeader = [
  {
    key: "User-Agent",
    value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
  },
]

const check = ({
  // poc
  poc: async function () {
    let pcfg = await CheckFileStat(lv)
    if (!pcfg) {
      version.LocalPoc = "版本文件不存在"
      version.PocStatus = false
      return
    } else {
      version.LocalPoc = "v" + await GetFileContent(lv)
    }
    let rp1 = await GoFetch("GET", download.RemotePocV, "", defultreqHeader, 10, null)
    if (rp1.Error == true) {
      version.PocUpdateContent = "检测更新失败"
      version.PocStatus = false
    } else {
      version.RemotePoc = rp1.Body
      if (compareVersion(version.LocalPoc, version.RemotePoc) == -1) {
        version.PocUpdateContent = (await GoFetch("GET", download.RemotePocCentent, "", defultreqHeader, 10, null)).Text
        version.PocStatus = true
      } else {
        version.PocUpdateContent = "当前已是最新版本"
        version.PocStatus = false
      }
    }
  },
  // client
  client: async function () {
    let rp2 = await GoFetch("GET", download.RemoteClientV, "", defultreqHeader, 10, null)
    if (rp2.Error == true) {
      version.ClientUpdateContent = "检测更新失败"
      version.ClientStatus = false
    } else {
      version.RemoteClient = rp2.Body
      if (compareVersion(version.LocalClient, version.RemoteClient) == -1) {
        version.ClientUpdateContent = (await GoFetch("GET", download.RemoteClientCentent, "", defultreqHeader, 10, null)).Text
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
    } else {
      ElNotification({
        title: "Error",
        message: "POC更新失败！！",
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
});

const download = {
  RemotePocV:
    "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version",
  RemoteClientV:
    "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version",
  RemotePocCentent: 'https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update',
  RemoteClientCentent: 'https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update',
};

const isMax = ref(false)
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

.el-badge {
  height: 120px;
}
</style>
