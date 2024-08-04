<script lang="ts" setup>
import Timestamp from './transcode/Timestamp.vue';
import OA from './transcode/OA.vue';
import Transcode from './transcode/Transcode.vue';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
const cmd5 = "https://www.cmd5.com"
</script>

<template>
    <div style="position: relative;">
        <el-tabs style="height: 100%; width: 100%;">
            <el-tab-pane label="编码转换" style="height: 100vh;">
                <Transcode></Transcode>
            </el-tab-pane>
            <el-tab-pane label="时间戳">
                <Timestamp></Timestamp>
            </el-tab-pane>
            <el-tab-pane label="OA密码解密">
                <OA></OA>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-button text bg type="primary" @click="BrowserOpenURL(cmd5)">CMD5</el-button>
        </div>
    </div>
</template>