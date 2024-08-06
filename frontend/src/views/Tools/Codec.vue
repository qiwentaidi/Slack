<script lang="ts" setup>
import { QuestionFilled } from '@element-plus/icons-vue'
import Timestamp from './transcode/Timestamp.vue';
import OA from './transcode/OA.vue';
import Transcode from './transcode/Transcode.vue';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
const cmd5 = "https://www.cmd5.com"
</script>

<template>
    <div style="position: relative;">
        <el-tabs style="height: 100%; width: 100%;">
            <el-tab-pane label="编码转换">
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
            <el-button text bg type="primary" @click="BrowserOpenURL('http://127.0.0.1:8731/')">
                <template #icon>
                    <el-tooltip placement="left">
                        <template #content>由于工具内容为iframe嵌入本地页面<br />
                            在切换其他路由时，iframe内容不会被keep-alive缓存<br />
                            所以每次路由换回来，页面都会重新加载导致内容都会丢失<br />
                            当你需要频繁切换页面时，推荐直接使用浏览器端的页面<br />
                        </template>
                        <QuestionFilled />
                    </el-tooltip>
                </template>
                打开本地环境
            </el-button>
            <el-button text bg type="primary" @click="BrowserOpenURL(cmd5)">CMD5</el-button>
        </div>
    </div>
</template>