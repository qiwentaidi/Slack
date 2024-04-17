<script setup lang="ts">
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import global from "../global"
</script>

<template>
    <div style="text-align: center;">
        <el-space direction="vertical" :size="8">
            <img src="/slack.svg" style="height: 8em; margin-bottom: 5px;">
            <span style="font-size: 20px; font-weight: bold; color: black;">slack-wails</span>
            <span style="font-weight: bold;">{{ "v" + global.LOCAL_VERSION }}</span>
            <div>
                <el-link @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack')">{{ $t('aside.source_code')
                    }}</el-link>
                <el-divider direction="vertical" />
                <el-link @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/issues/new')">{{ $t('aside.issue')
                    }}</el-link>
                <el-divider direction="vertical" />
                <el-link
                    @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97')">{{
                    $t('aside.update_log') }}</el-link>
            </div>
        </el-space>
    </div>
</template>