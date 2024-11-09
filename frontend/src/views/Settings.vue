<template>
  <el-container style="height: 100%;">
    <el-aside width="200px">
      <el-menu default-active="0">
        <el-menu-item v-for="(item, index) in setupOptions" :index="index.toString()" @click="selectItem">
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <template #title><span>{{ $t(item.name) }}</span></template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main>
      <el-form :model="global.webscan" label-width="auto" v-show="currentDisplay == '0'">
        <h3>{{ $t(setupOptions[0].name) }}</h3>
        <el-form-item label="网站指纹线程">
          <el-input-number controls-position="right" v-model="global.webscan.web_thread" :min="1" :max="200" />
        </el-form-item>
        <el-form-item label="端口指纹超时">
          <el-input-number controls-position="right" v-model="global.webscan.port_timeout" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="存活验证模式">
          <el-select v-model="global.webscan.default_alive_module">
            <el-option v-for="item in aliveGroupOptions" :key="item.value" :value="item.value">
              <span style="float: left">{{ item.value }}</span>
              <span style="float: right">
                {{ item.description }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('setting.network_list')">
          <el-select v-model="global.webscan.default_network">
            <el-option v-for="value in global.temp.NetworkCardList" :label="value" :value="value" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-form :inline="true" :model="global.proxy" label-width="auto" class="demo-form-inline"
        v-show="currentDisplay == '1'">
        <h3>{{ $t(setupOptions[1].name) }}</h3>
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
            <el-input-number v-model.number="global.proxy.port" :min="1" :max="65535"
              style="width: 220px;"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('setting.username')">
            <el-input v-model="global.proxy.username" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('setting.password')">
            <el-input v-model="global.proxy.password" type="password" show-password></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="SaveConfig" style="float: right;">{{ $t('setting.save') }}</el-button>
          </el-form-item>
        </div>
      </el-form>
      <el-form :model="global.space" label-width="auto" v-show="currentDisplay == '2'">
        <h3>{{ $t(setupOptions[2].name) }}<el-divider direction="vertical" />Ⓓ标识符API主要用于收集子域名信息
          <el-divider direction="vertical" />🎉标识符建议优先注册</h3>
        <el-form-item label="FOFA">
          <el-input v-model="global.space.fofaapi" placeholder="API address" clearable />
          <el-input v-model="global.space.fofaemail" placeholder="Email" clearable style="margin-top: 5px;" />
          <el-input v-model="global.space.fofakey" placeholder="API key" clearable style="margin-top: 5px;" />
        </el-form-item>
        <el-form-item label="Hunter">
          <el-input v-model="global.space.hunterkey" placeholder="API key" clearable />
        </el-form-item>
        <el-form-item label="Quake">
          <el-input v-model="global.space.quakekey" placeholder="API key" clearable />
        </el-form-item>
        <el-form-item label="🎉ChaosⒹ">
          <el-input v-model="global.space.chaos">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled" @click="BrowserOpenURL(chaosURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="BevigilⒹ">
          <el-input v-model="global.space.bevigil" placeholder="API key">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled" @click="BrowserOpenURL(bevigilURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="SecuritytrailsⒹ">
          <el-input v-model="global.space.securitytrails" placeholder="API key">
            <template #suffix>
              <el-button link type="primary" :icon="UserFilled"
                @click="BrowserOpenURL(securitytrailsURL)">注册</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="ZoomeyeⒹ">
          <el-input v-model="global.space.zoomeye" placeholder="API key" />
        </el-form-item>
        <el-form-item label="GithubⒹ">
          <el-input v-model="global.space.github"
            placeholder="Settings -> Developer settings -> Presonal access tokens" />
        </el-form-item>
        <el-button type="primary" @click="SaveConfig" style="float: right;">{{ $t('setting.save') }}</el-button>
      </el-form>
      <el-form :model="global.Theme" label-width="auto" v-show="currentDisplay == '3'">
        <h3>{{ $t(setupOptions[3].name) }}</h3>
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
        <h3>{{ $t(setupOptions[4].name) }}<el-divider direction="vertical" />密码所有协议通用</h3>
        <el-table :data="crackDict.usernames" stripe style="width: 100%">
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
      <div style="display: flex; justify-content: center; align-items: center;" v-show="currentDisplay == '5'">
        <about></about>
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
import { ElMessage, MenuItemRegistered } from 'element-plus';
import { TestProxy } from "@/util";
import { Edit, Sunny, Moon, UserFilled } from '@element-plus/icons-vue';
import { reactive, ref } from "vue";
import { ReadFile, WriteFile } from "wailsjs/go/main/File";
import { File } from '@/stores/interface';
import { useI18n } from "vue-i18n";
import { useDark, useToggle } from '@vueuse/core'
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import { SaveConfig } from "@/config";
import { aliveGroupOptions, crackDict, setupOptions } from "@/stores/options";

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

const currentDisplay = ref('0')

function selectItem(item: MenuItemRegistered) {
  currentDisplay.value = item.index
}


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