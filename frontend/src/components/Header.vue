<template>
  <el-menu index="0" mode="horizontal" :ellipsis="false" route text-color="#000">
    <el-menu-item @click="currentComponent = 'Introduction'">
      <el-icon>
        <HomeFilled />
      </el-icon>
    </el-menu-item>
    <el-divider class="divder" direction="vertical"/>
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
    <el-menu-item index="5" @click="">
        <el-icon>
          <Bell />
        </el-icon>
      </el-menu-item>
      <el-divider class="divder" direction="vertical"/>
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
    <el-divider class="divder" direction="vertical"/>
    <!-- MAX MIN CLOSE -->
    <el-menu-item @click="minimise">
      <el-icon>
        <Minus />
      </el-icon>
    </el-menu-item>
    <el-menu-item @click="ToggleMaximise">
      <el-icon>
        <FullScreen />
      </el-icon>
    </el-menu-item>
    <el-menu-item @click="quit">
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
import { ref, inject } from "vue";
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
} from "@element-plus/icons-vue";

import { Quit, Minimise, ToggleMaximise } from "../../wailsjs/go/main/App";

const centerDialogVisible = ref(false);
const currentComponent = inject("currentComponent");

function quit() {
  Quit();
}

function minimise() {
  Minimise();
}

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
</style>
