<template>
    <el-form :model="form" label-width="100px" >
      <el-form-item label="wx_appid">
        <el-input v-model="form.appid" />
      </el-form-item>
      <el-form-item label="wx_secert">
        <el-input v-model="form.secert" />
      </el-form-item>
      <el-form-item label="result">
        <el-input v-model="form.result" type="textarea" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="checksecert">Check</el-button>
      </el-form-item>
    </el-form>
  </template>
  
<script lang="ts" setup>
import { reactive } from 'vue'
import {
  WechatAppid
} from '../../../wailsjs/go/main/App'
// do not use same name with ref
const form = reactive({
appid: '',
secert: '',
result: '',
})

function checksecert(){
  WechatAppid(form.appid,form.secert).then(
        result => {
            form.result = result
        }
    )
}

</script>