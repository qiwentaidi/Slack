<template>
  <el-container>
    <el-scrollbar max-height="85vh">
      <el-main>
        <el-collapse model-value="1">
          <el-collapse-item name="1"><template #title>
              <h2>DNS Server</h2>
            </template>
            <el-form :model="global.scan" label-width="110px" label-position="top">
              <el-form-item>
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-input v-model="global.scan.dns1" placeholder="DNS1"><template #append>:53</template></el-input>
                  </el-col>
                  <el-col :span="12">
                    <el-input v-model="global.scan.dns2" placeholder="DNS2"><template #append>:53</template></el-input>
                  </el-col>
                </el-row>
              </el-form-item>
            </el-form>
          </el-collapse-item>
          <el-collapse-item name="2"><template #title>
              <h2>{{ $t('setting.proxy') }}</h2>
            </template>
            <el-form :inline="true" :model="global.proxy" label-width="80px" class="demo-form-inline">
              <el-form-item :label="$t('setting.enable')">
                <el-switch v-model="global.proxy.enabled" />
                <el-button type="primary" style="margin-left: 20px;" v-if="global.proxy.enabled">测试代理</el-button>
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
                  <el-input v-model="global.proxy.port" />
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
              <h2>{{ $t('aside.space_engine') }}</h2>
            </template>
            <el-form :model="global.space" label-width="100px" label-position="top">
              <el-form-item label="FOFA" style="margin-top: 10px;">
                <el-input v-model="global.space.fofaapi" placeholder="api address" clearable></el-input>
                <el-input v-model="global.space.fofaemail" placeholder="email" clearable style="margin-top: 5px;"></el-input>
                <el-input v-model="global.space.fofakey" placeholder="key" clearable style="margin-top: 5px;"></el-input>
              </el-form-item>
              <el-form-item :label="$t('aside.hunter')">
                <el-input v-model="global.space.hunterkey" placeholder="key" clearable></el-input>
              </el-form-item>
              <!-- <el-form-item :label="$t('aside.360quake')">
                <el-input v-model="global.space.quakekey" placeholder="key" clearable></el-input>
              </el-form-item> -->
            </el-form>
          </el-collapse-item>
        </el-collapse>
      </el-main>
    </el-scrollbar>
    <el-footer>
      <el-button type="primary" @click="saveConfig">{{ $t('setting.save') }}</el-button>
    </el-footer>
  </el-container>
</template>

<script lang="ts" setup>
import global from "../global"
import { ElMessage } from 'element-plus';
import { GoFetch } from "../../wailsjs/go/main/App";

const saveConfig = () => {
  global.space.fofaapi = global.space.fofaapi.replace(/[\r\n\s]/g, '');
  global.space.fofaemail = global.space.fofaemail.replace(/[\r\n\s]/g, '');
  global.space.fofakey = global.space.fofakey.replace(/[\r\n\s]/g, '');
  global.space.hunterkey = global.space.hunterkey.replace(/[\r\n\s]/g, '');
  global.space.quakekey = global.space.quakekey.replace(/[\r\n\s]/g, '');
  localStorage.setItem('scan', JSON.stringify(global.scan));
  localStorage.setItem('proxy', JSON.stringify(global.proxy));
  localStorage.setItem('space', JSON.stringify(global.space));
  ElMessage({
    message: '保存成功',
    type: 'success',
  })
};

</script>

<style>
.el-footer {
  position: fixed;
  bottom: 0;
  width: 100%;
}
.mf {
  margin-right: 5px;
}

.demo-form-inline .el-input {
  --el-input-width: 220px;
}

.demo-form-inline .el-select {
  --el-select-width: 220px;
}
</style>