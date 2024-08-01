<template>
  <el-scrollbar style="height: 87vh;">
    <el-collapse model-value="1">
      <el-collapse-item name="1"><template #title>
          <h2>{{ $t('setting.scan') }}</h2>
        </template>
        <el-form label-width="auto">
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
        </el-form>
      </el-collapse-item>
      <el-collapse-item name="2"><template #title>
          <h2>{{ $t('setting.proxy') }}</h2>
        </template>
        <el-form :inline="true" :model="global.proxy" label-width="auto" class="demo-form-inline">
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
              <el-input-number v-model="global.proxy.port" :min="1" :max="65535"
                style="width: 220px;"></el-input-number>
            </el-form-item>
            <el-form-item :label="$t('setting.username')">
              <el-input v-model="global.proxy.username" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('setting.password')">
              <el-input v-model="global.proxy.password" clearable></el-input>
            </el-form-item>
          </div>
        </el-form>
      </el-collapse-item>
      <el-collapse-item name="3"><template #title>
          <h2>{{ $t('setting.mapping') }}</h2>
        </template>
        <el-form :model="global.space" label-width="auto">
          <el-form-item label="FOFA">
            <el-input v-model="global.space.fofaapi" placeholder="api address" clearable />
            <el-input v-model="global.space.fofaemail" placeholder="email" clearable style="margin-top: 5px;" />
            <el-input v-model="global.space.fofakey" placeholder="key" clearable style="margin-top: 5px;" />
          </el-form-item>
          <el-form-item :label="$t('aside.hunter')">
            <el-input v-model="global.space.hunterkey" placeholder="key" clearable />
          </el-form-item>
          <el-form-item :label="$t('aside.360quake')">
            <el-input v-model="global.space.quakekey" placeholder="token" clearable />
          </el-form-item>
        </el-form>
      </el-collapse-item>
      <el-collapse-item name="4"><template #title>
          <h2>{{ $t('aside.display') }}</h2>
        </template>
        <el-form label-width="auto">
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
      </el-collapse-item>
      <el-collapse-item name="5"><template #title>
          <h2>{{ $t('aside.dict') }}</h2>
        </template>
        <div>
          <el-alert type="info" :closable="false" show-icon>内置字典密码所有协议通用</el-alert>
          <el-table :data="global.dict.usernames" stripe style="width: 100%" @row-dblclick="">
            <el-table-column prop="name" label="服务名称" />
            <el-table-column label="操作" width="250" align="center">
              <template #default="scope">
                <el-button type="primary" link :icon="Edit"
                  @click="ctrl.innerDrawer = true; ctrl.currentPath = '/username/' + scope.row.name + '.txt'; ReadDict(ctrl.currentPath)">用户名</el-button>
                <el-button type="primary" link :icon="Edit"
                  @click="ctrl.innerDrawer = true; ctrl.currentPath = '/password/password.txt'; ReadDict(ctrl.currentPath)">密码</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-collapse-item>
    </el-collapse>
  </el-scrollbar>
  <el-button type="primary" @click="saveConfig" style="float: right;">{{ $t('setting.save') }}</el-button>
  <el-drawer v-model="ctrl.innerDrawer" title="字典管理" :append-to-body="true">
    <el-input type="textarea" rows="20" v-model="ctrl.currentDic"></el-input>
    <el-button type="primary" style="margin-top: 10px; float: right;" @click="SaveFile(ctrl.currentPath)">保存</el-button>
  </el-drawer>
</template>

<script lang="ts" setup>
import global from "../global"
import { ElMessage, ElNotification } from 'element-plus';
import { TestProxy, TestNuclei } from "../util";
import { Edit, Sunny, Moon } from '@element-plus/icons-vue';
import { reactive } from "vue";
import { ReadFile, SaveDataToFile, UserHomeDir, WriteFile } from "../../wailsjs/go/main/File";
import { File } from '../interface';
import { useI18n } from "vue-i18n";
import { useDark, useToggle } from '@vueuse/core'

const isDark = useDark({
  storageKey: 'theme',
  valueDark: 'dark',
  valueLight: 'light',
})

const toggle = useToggle(isDark)

const { locale } = useI18n();
const changeLanguage = (lang: string) => {
  localStorage.setItem("language", lang);
  locale.value = lang;
}
const removeInvisibleFiled = (filed: string) => {
  filed = filed.replace(/[\r\n\s]/g, '');
}
async function saveConfig() {
  let list = [global.space.fofaapi, global.space.fofaemail, global.space.fofakey, global.space.hunterkey, global.space.quakekey]
  list.forEach(item => {
    removeInvisibleFiled(item)
  })
  var data = { proxy: global.proxy, space: global.space, jsfind: global.jsfind, webscan: global.webscan };
  let result = await SaveDataToFile(data);
  if (result) {
    ElNotification.success({
      message: 'Save successful',
      position: 'bottom-right'
    })
  }
};

const ctrl = reactive({
  innerDrawer: false,
  currentDic: '',
  currentPath: '',
  theme: 0,
})

async function ReadDict(path: string) {
  let home = await UserHomeDir()
  let file: File = await ReadFile(home + global.PATH.PortBurstPath + path)
  ctrl.currentDic = file.Content!
}

async function SaveFile(path: string) {
  let home = await UserHomeDir()
  WriteFile('txt', home + global.PATH.PortBurstPath + path, ctrl.currentDic).then(result => {
    result ? ElMessage.success('保存成功!') : ElMessage.error('保存失败!')
  })
}
</script>

<style>
.demo-form-inline .el-input {
  --el-input-width: 220px;
}

.demo-form-inline .el-select {
  --el-select-width: 220px;
}
</style>