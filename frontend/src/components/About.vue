<script setup lang="ts">
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import global from "@/global"
const wikiAddress = "https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97"
const projectAddress = "https://github.com/qiwentaidi/Slack/"
const issueAddress = "https://github.com/qiwentaidi/Slack/issues/new"
</script>

<template>
    <div style="text-align: center;">
        <el-space direction="vertical" :size="8">
            <img src="/slack.svg" style="height: 8em; margin-bottom: 5px;">
            <span class="bold">Slack <el-tag round type="info" size="small" style="margin-bottom: 2px;">{{ "v" +
            global.LOCAL_VERSION }}</el-tag></span>
            <span class="slogan">{{ $t('aside.slogan') }}</span>

            <div>
                <el-link @click="BrowserOpenURL(projectAddress)">{{ $t('aside.source_code')
                    }}</el-link>
                <el-divider direction="vertical" />
                <el-link @click="BrowserOpenURL(issueAddress)">{{ $t('aside.issue')
                    }}</el-link>
                <el-divider direction="vertical" />
                <el-link @click="BrowserOpenURL(wikiAddress)">{{
            $t('aside.update_log') }}</el-link>
            </div>
        </el-space>
    </div>
</template>

<style>
.bold {
    font-weight: bold;
    font-size: 20px;
}
.slogan {
    font-weight: bold; 
    font-size: small;
}
</style>