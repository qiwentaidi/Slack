<script setup lang="ts">
import { Download } from "@element-plus/icons-vue";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import global from "../global"
import { UpdatePocFile } from "../../wailsjs/go/main/File";
import { ElNotification } from "element-plus";


const update = ({
    poc: async function () {
        let err = await UpdatePocFile()
        if (err == "") {
            ElNotification({
                message: "POC update success!",
                type: "success",
            });
        } else {
            ElNotification({
                message: "POC update failed! " + err,
                type: "error",
            });
        }
    }
})

</script>

<template>
    <el-card class="box-card">
        <template #header>
            <div class="card-header">
                <el-text>
                    <span style="font-weight: bold;">POC&指纹{{ global.UPDATE.RemotePocVersion }}</span>
                    <br />当前{{ "v" + global.UPDATE.LocalPocVersion }}
                </el-text>
                <el-button class="button" :icon="Download" text @click="update.poc"
                    v-if="global.UPDATE.PocStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.PocContent }}</span>
            </div>
        </template>
        <el-scrollbar style="height: 100px;" class="pretty-response" v-if="global.UPDATE.PocStatus">
            {{ global.UPDATE.PocContent }}
        </el-scrollbar>
    </el-card>

    <el-card class="box-card" style="margin-top: 10px;">
        <template #header>
            <div class="card-header">
                <el-text>
                    <span style="font-weight: bold;">客户端{{ global.UPDATE.RemoteClientVersion }}</span>
                    <br />当前{{ "v" + global.LOCAL_VERSION }}
                </el-text>
                <el-button class="button" :icon="Download" text
                    @click="BrowserOpenURL('https://github.com/qiwentaidi/Slack/releases')"
                    v-if="global.UPDATE.ClientStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.ClientContent }}</span>
            </div>
        </template>
        <el-scrollbar style="height: 100px;" class="pretty-response" v-if="global.UPDATE.ClientStatus">
            {{ global.UPDATE.ClientContent }}
        </el-scrollbar>
    </el-card>
</template>