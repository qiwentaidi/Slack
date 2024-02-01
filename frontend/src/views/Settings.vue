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
              <h2>代理配置(仅适用网站扫描)</h2>
            </template>
            <el-form :model="global.proxy" label-width="80px">
              <el-form-item label="开启代理:">
                <el-switch v-model="global.proxy.enabled" />
              </el-form-item>
              <div v-if="global.proxy.enabled">
                <el-form-item label="代理模式:">
                  <el-select v-model="global.proxy.mode">
                    <el-option label="HTTP" value="HTTP" />
                    <el-option label="SOCK5" value="SOCK5" />
                  </el-select>
                </el-form-item>
                <el-form-item label="代理地址:">
                  <el-input v-model="global.proxy.address" clearable></el-input>
                </el-form-item>
                <el-form-item label="代理端口:">
                  <el-input-number controls-position="right" v-model="global.proxy.port" :min="1" :max="65535" />
                </el-form-item>
                <el-form-item label="用户名:">
                  <el-input v-model="global.proxy.username" clearable></el-input>
                </el-form-item>
                <el-form-item label="密码:">
                  <el-input v-model="global.proxy.password" clearable></el-input>
                </el-form-item>
              </div>
            </el-form>
          </el-collapse-item>
          <el-collapse-item name="3"><template #title>
              <h2>资产测绘</h2>
            </template>
            <el-form :model="global.space" label-width="80px">
              <el-form-item label="FOFA:" style="margin-top: 10px;">
                <el-input v-model="global.space.fofaemail" placeholder="邮箱" clearable></el-input>
                <el-input v-model="global.space.fofakey" placeholder="key" clearable style="margin-top: 5px;"></el-input>
              </el-form-item>
              <el-form-item label="鹰图: ">
                <el-input v-model="global.space.hunterkey" placeholder="key" clearable></el-input>
              </el-form-item>
              <el-form-item label="夸克: ">
                <el-input v-model="global.space.quakekey" placeholder="key" clearable></el-input>
              </el-form-item>
            </el-form>
          </el-collapse-item>
        </el-collapse>
      </el-main>
    </el-scrollbar>
    <el-footer>
      <el-button type="primary" @click="saveConfig">保存</el-button>
    </el-footer>
  </el-container>
</template>

<script lang="ts" setup>
import global from "../global"
import { ElMessage } from 'element-plus';

const saveConfig = () => {
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
</style>