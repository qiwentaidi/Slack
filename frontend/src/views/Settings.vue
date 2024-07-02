<template>
  <el-scrollbar max-height="87vh" style="height: 89vh;">
    <el-collapse model-value="1">
      <el-collapse-item name="1"><template #title>
          <h2>{{ $t('aside.webscan') }}</h2>
        </template>
        <el-form :inline="true" :model="global.proxy" label-width="auto" class="demo-form-inline">
          <el-form-item>
            <template #label>{{ $t('setting.engine') }}
              <el-tooltip placement="left">
                <template #content>
                  {{ $t('setting.nuclei_placeholder1') }}
                </template>
                <el-icon>
                  <QuestionFilled size="24" />
                </el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="global.webscan.nucleiEngine" :placeholder="$t('setting.nuclei_placeholder2')" clearable
              style="width: 115vh;"></el-input>
            <el-button type="primary" style="margin-left: 10px;" @click="TestNuclei()">{{ $t('setting.enable')
              }}</el-button>
          </el-form-item>
          <el-form-item :label="$t('setting.proxy')">
            <el-switch v-model="global.proxy.enabled" />
            <el-button type="primary" size="small" @click="TestProxy(0)" style="margin-left: 20px;"
              v-if="global.proxy.enabled">{{ $t('setting.test_agent') }}</el-button>
          </el-form-item>
          <div v-if="global.proxy.enabled">
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
      <el-collapse-item name="2"><template #title>
          <h2>{{ $t('aside.space_engine') }}</h2>
        </template>
        <el-form :model="global.space" label-width="auto">
          <el-form-item label="FOFA" style="margin-top: 10px;">
            <el-input v-model="global.space.fofaapi" placeholder="api address" clearable></el-input>
            <el-input v-model="global.space.fofaemail" placeholder="email" clearable
              style="margin-top: 5px;"></el-input>
            <el-input v-model="global.space.fofakey" placeholder="key" clearable style="margin-top: 5px;"></el-input>
          </el-form-item>
          <el-form-item :label="$t('aside.hunter')">
            <el-input v-model="global.space.hunterkey" placeholder="key" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('aside.360quake')">
            <el-input v-model="global.space.quakekey" placeholder="token" clearable></el-input>
          </el-form-item>
        </el-form>
      </el-collapse-item>
      <el-collapse-item name="3"><template #title>
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
import { Edit } from '@element-plus/icons-vue';
import { reactive } from "vue";
import { ReadFile, WriteFile } from "../../wailsjs/go/main/File";
import { File } from '../interface';

function saveConfig() {
  global.space.fofaapi = global.space.fofaapi.replace(/[\r\n\s]/g, '');
  global.space.fofaemail = global.space.fofaemail.replace(/[\r\n\s]/g, '');
  global.space.fofakey = global.space.fofakey.replace(/[\r\n\s]/g, '');
  global.space.hunterkey = global.space.hunterkey.replace(/[\r\n\s]/g, '');
  global.space.quakekey = global.space.quakekey.replace(/[\r\n\s]/g, '');
  localStorage.setItem('proxy', JSON.stringify(global.proxy));
  localStorage.setItem('space', JSON.stringify(global.space));
  localStorage.setItem('webscan', JSON.stringify(global.webscan));
  ElNotification.success({
    message: 'Save successful',
    position: 'bottom-right'
  })
};

const ctrl = reactive({
  drawer: false,
  innerDrawer: false,
  currentDic: '',
  currentPath: '',
})

async function ReadDict(path: string) {
  let file: File = await ReadFile(global.PATH.PortBurstPath + path)
  ctrl.currentDic = file.Content!
}

async function SaveFile(path: string) {
  let r = await WriteFile('txt', global.PATH.PortBurstPath + path, ctrl.currentDic)
  if (r) {
    ElMessage.success({
      showClose: true,
      message: 'success',
    })
  } else {
    ElMessage.error({
      showClose: true,
      message: 'failed',
    })
  }
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