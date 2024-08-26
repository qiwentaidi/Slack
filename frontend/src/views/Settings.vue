<template>
  <el-container style="height: 100vh;">
    <el-aside width="200px">
      <el-menu default-active="0">
        <el-menu-item v-for="item in settingOption" :index="item.show" @click="currentDisplay = item.show">
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <template #title><span>{{ $t(item.name) }}</span></template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main>
      <el-form :model="global.webscan" label-width="auto" v-show="currentDisplay == '0'">
        <h3>{{ $t(settingOption[0].name) }}</h3>
        <el-form-item :label="$t('setting.engine')">
          <el-input v-model="global.webscan.nucleiEngine" :placeholder="$t('setting.nuclei_placeholder')" clearable>
            <template #suffix>
              <el-button type="primary" link @click="TestNuclei()">{{ $t('setting.engine_enable')
                }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('setting.network_list')">
          <el-select v-model="global.temp.defaultNetwork">
            <el-option v-for="value in global.temp.NetworkCardList" :label="value" :value="value" />
          </el-select>
        </el-form-item>
        <el-button type="primary" @click="saveConfig" style="float: right;">{{ $t('setting.save') }}</el-button>
      </el-form>
      <el-form :inline="true" :model="global.proxy" label-width="auto" class="demo-form-inline"
        v-show="currentDisplay == '1'">
        <h3>{{ $t(settingOption[1].name) }}</h3>
        <el-form-item :label="$t('setting.enable')">
          <el-switch v-model="global.proxy.enabled" />
          <el-button type="primary" size="small" @click="TestProxy(0)" style="margin-left: 20px;"
            v-if="global.proxy.enabled">{{ $t('setting.test_agent') }}</el-button>
        </el-form-item>
        <div>
          <el-form-item :label="$t('setting.mode')">
            <el-select v-model="global.proxy.mode">
              <el-option label="HTTP" value="HTTP" />
              <el-option label="SOCK5" value="SOCK5" />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('setting.address')">
            <el-input v-model="global.proxy.address" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('setting.port')">
            <el-input-number v-model="global.proxy.port" :min="1" :max="65535" style="width: 220px;"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('setting.username')">
            <el-input v-model="global.proxy.username" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('setting.password')">
            <el-input v-model="global.proxy.password" type="password" show-password></el-input>
          </el-form-item>
        </div>
      </el-form>
      <el-form :model="global.space" label-width="auto" v-show="currentDisplay == '2'">
        <h3>{{ $t(settingOption[2].name) }}</h3>
        <el-form-item label="FOFA">
          <el-input v-model="global.space.fofaapi" placeholder="API address" clearable />
          <el-input v-model="global.space.fofaemail" placeholder="Email" clearable style="margin-top: 5px;" />
          <el-input v-model="global.space.fofakey" placeholder="API key" clearable style="margin-top: 5px;" />
        </el-form-item>
        <el-form-item :label="$t('aside.hunter')">
          <el-input v-model="global.space.hunterkey" placeholder="API key" clearable />
        </el-form-item>
        <el-form-item :label="$t('aside.360quake')">
          <el-input v-model="global.space.quakekey" placeholder="API key" clearable />
        </el-form-item>
        <el-alert type="info" show-icon :closable="false" style="margin-bottom: 5px;">
          <span>下方API可用于收集子域名信息，在使用子域名收集模块请配置（优先推荐配置Chao，免费且不限次数）</span>
        </el-alert>
        <el-form-item label="Chaos">
          <el-input v-model="global.space.chaos" placeholder="API key">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled" @click="BrowserOpenURL(chaosURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="Bevigil">
          <el-input v-model="global.space.bevigil" placeholder="API key">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled" @click="BrowserOpenURL(bevigilURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="Securitytrails">
          <el-input v-model="global.space.securitytrails" placeholder="API key">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled"
                @click="BrowserOpenURL(securitytrailsURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="Zoomeye">
          <el-input v-model="global.space.zoomeye" placeholder="API key" />
        </el-form-item>
        <el-form-item label="Github">
          <el-input v-model="global.space.github"
            placeholder="Settings -> Developer settings -> Presonal access tokens" />
        </el-form-item>
        <el-button type="primary" @click="saveConfig" style="float: right;">{{ $t('setting.save') }}</el-button>
      </el-form>
      <el-form label-width="auto" v-show="currentDisplay == '3'">
        <h3>{{ $t(settingOption[3].name) }}</h3>
        <el-form-item :label="$t('aside.language')">
          <el-select v-model="global.Language.value" @change="changeLanguage" style="width: 150px;">
            <el-option label="中文" value="zh"></el-option>
            <el-option label="English" value="en"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('aside.theme')">
          <el-switch v-model="global.Theme.value" :active-action-icon="Moon" :inactive-action-icon="Sunny"
            style="--el-switch-on-color: #2C2C2C; --el-switch-off-color: " @change="toggle" />
        </el-form-item>
      </el-form>
      <div v-show="currentDisplay == '4'">
        <h3>{{ $t(settingOption[4].name) }}</h3>
        <el-alert type="info" :closable="false" show-icon>内置字典密码所有协议通用</el-alert>
        <el-table :data="global.dict.usernames" stripe style="width: 100%">
          <el-table-column prop="name" label="服务名称" />
          <el-table-column label="操作" width="250" align="center">
            <template #default="scope">
              <el-button type="primary" link :icon="Edit"
                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/username/' + scope.row.name + '.txt'; ReadDict(ctrl.currentPath)">{{
                  $t('setting.username') }}</el-button>
              <el-button type="primary" link :icon="Edit"
                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/password/password.txt'; ReadDict(ctrl.currentPath)">{{
                  $t('setting.password') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-main>
  </el-container>
  <el-drawer v-model="ctrl.innerDrawer" title="字典管理" :append-to-body="true">
    <el-input type="textarea" :rows="20" v-model="ctrl.currentDic"></el-input>
    <el-button type="primary" style="margin-top: 10px; float: right;" @click="SaveFile(ctrl.currentPath)">保存</el-button>
  </el-drawer>
</template>

<script lang="ts" setup>
import global from "@/global"
import { ElMessage, ElNotification } from 'element-plus';
import { TestProxy, TestNuclei } from "@/util";
import { Edit, Sunny, Moon, UserFilled } from '@element-plus/icons-vue';
import { reactive, ref } from "vue";
import { ReadFile, SaveDataToFile, WriteFile } from "wailsjs/go/main/File";
import { File } from '@/interface';
import { useI18n } from "vue-i18n";
import { useDark, useToggle } from '@vueuse/core'
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import scanIcon from "@/assets/icon/scan.svg"
import themeIcon from "@/assets/icon/theme.svg"
import proxyIcon from "@/assets/icon/proxy.svg"
import dictmanagerIcon from "@/assets/icon/dict.svg"
import layersIcon from "@/assets/icon/layers.svg"

const bevigilURL = "https://bevigil.com/osint-api"
const chaosURL = "https://cloud.projectdiscovery.io/"
const securitytrailsURL = "https://securitytrails.com/"

const isDark = useDark({
  storageKey: 'theme',
  valueDark: 'dark',
  valueLight: 'light',
})

const toggle = useToggle(isDark)

const { locale } = useI18n();

global.Language.value = locale.value
const changeLanguage = (lang: string) => {
  localStorage.setItem("language", lang);
  locale.value = lang;
}

const saveConfig = () => {
  // 获取space的所有value值
  let list = Object.entries(global.space).map(([key, value]) => value);
  // 去除不可见字符
  list = list.map(item => item.replace(/[\r\n\s]/g, ''));
  var data = { proxy: global.proxy, space: global.space, jsfind: global.jsfind, webscan: global.webscan };
  SaveDataToFile(data).then(result => {
    if (result) {
      ElNotification.success({
        message: 'Save successful',
        position: 'bottom-right'
      })
    }
  })
};

const ctrl = reactive({
  innerDrawer: false,
  currentDic: '',
  currentPath: '',
})

async function ReadDict(path: string) {
  let file: File = await ReadFile(global.PATH.homedir + global.PATH.PortBurstPath + path)
  ctrl.currentDic = file.Content!
}

async function SaveFile(path: string) {
  WriteFile('txt', global.PATH.homedir + global.PATH.PortBurstPath + path, ctrl.currentDic).then(result => {
    result ? ElMessage.success('保存成功!') : ElMessage.error('保存失败!')
  })
}

const settingOption = [
  {
    name: 'setting.scan',
    icon: scanIcon,
    show: '0',
  },
  {
    name: 'setting.proxy',
    icon: proxyIcon,
    show: '1',
  },
  {
    name: 'setting.mapping',
    icon: layersIcon,
    show: '2',
  },
  {
    name: 'aside.display',
    icon: themeIcon,
    show: '3',
  },
  {
    name: 'aside.dict',
    icon: dictmanagerIcon,
    show: '4',
  },
]

const currentDisplay = ref('0')
</script>

<style scoped>
.demo-form-inline .el-input {
  --el-input-width: 220px;
}

.demo-form-inline .el-select {
  --el-select-width: 220px;
}

.el-main {
  padding-top: 0px;
  padding-right: 0px;
}

.el-menu {
  border-right: 0px;
}

.el-menu-item.is-active {
  background-color: var(--el-menu-active-color);
  color: var(--el-menu-hover-bg-color);
  border-radius: 5px;
}

.el-menu-item:hover {
  border-radius: 5px;
}
</style>