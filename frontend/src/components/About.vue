<script setup lang="ts">
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import global from "@/stores"
import githubIcon from "@/assets/icon/github.svg";

const wikiAddress = "https://github.com/qiwentaidi/Slack/wiki/%E6%9B%B4%E6%96%B0%E6%97%A5%E5%BF%97"
const projectAddress = "https://github.com/qiwentaidi/Slack/"
const issueAddress = "https://github.com/qiwentaidi/Slack/issues/new"

const relatedTools = [
    {
        Name: "LuckyCat",
        Description: "A chrome-extension for website fingerprint extraction",
        Link: "https://github.com/qiwentaidi/LuckyCat",
        Image: "/luckycat.png"
    },
    {
        Name: "initial",
        Description: "Slack webscan cli",
        Link: "https://github.com/qiwentaidi/initial",
        Image: "/initial.png"
    }
]
</script>

<template>
    <el-space direction="vertical" :size="8">
        <img src="/slack.svg" style="height: 8em; margin-bottom: 5px;">
        <span class="bold">Slack <el-tag round type="info" style="margin-bottom: 2px;">{{ "v" +
            global.LOCAL_VERSION }}</el-tag></span>
        <span class="slogan">{{ $t('setting.slogan') }}</span>

        <div>
            <el-link @click="BrowserOpenURL(projectAddress)">{{ $t('setting.source_code')
                }}</el-link>
            <el-divider direction="vertical" />
            <el-link @click="BrowserOpenURL(issueAddress)">{{ $t('setting.issue')
                }}</el-link>
            <el-divider direction="vertical" />
            <el-link @click="BrowserOpenURL(wikiAddress)">{{
                $t('setting.update_log') }}</el-link>
        </div>
        <div></div>
        <span>{{ $t('setting.related') }}</span>
        <div style="width: 500px;">
            <el-carousel height="350px" motion-blur>
                <el-carousel-item v-for="item in relatedTools" :key="item.Name">
                    <el-card shadow="hover" style="height: 90%; text-align: center;">
                        <el-image :src="item.Image" fit="contain" :preview-src-list="[item.Image]" preview-teleported
                            style="height: 180px;" />
                        <el-link @click="BrowserOpenURL(item.Link)">
                            <h4>{{ item.Name }}</h4>
                            <el-icon class="el-icon--right">
                                <githubIcon />
                            </el-icon>
                        </el-link>
                        <p>{{ item.Description }}</p>
                    </el-card>
                </el-carousel-item>
            </el-carousel>
        </div>
        <div></div>
        <span>{{ $t('setting.contact_us') }}</span>
        <img src="@/assets/icon/wechat.png" style="height: 208px;">
    </el-space>
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