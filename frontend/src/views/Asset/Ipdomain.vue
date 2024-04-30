<script lang="ts" setup>
import { reactive } from 'vue';
import { MapLocation } from '@element-plus/icons-vue';
import {
    DomainInfo,
    CheckCdn,
} from '../../../wailsjs/go/main/App'
import async from 'async';
import global from "../../global"
import { SplitTextArea } from '../../util';
const info = `配域名进行IP解析以及CDN判断(可批量)
===如果域名解析的非常慢，请考虑是否是本机网络不佳===
`;
const from = reactive({
    input: '',
    result: '',
});

const domainRegex = /^(?=.{1,253}$)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)+[a-z]{2,}$/i;

function searchdomain() {
    DomainInfo(from.input).then(
        result => {
            from.result = result
        }
    )
}

function cdncheck() {
    let lines = SplitTextArea(from.input)
    let domains = [] as string[]
    for (const domain of lines) {
        domainRegex.test(domain) ? domains.push(domain) : console.log(domain + ' is not a domain')
    }
    from.result += "---域名解析(CDN查询)---:\n"
    async.eachLimit(domains, 100, (domain: string) => {
        let result = CheckCdn(from.input, global.scan.dns1, global.scan.dns2)
        from.result += result + "\n"
    })
}

</script>

<template>
    <div class="head" style="margin-top: 10px; margin-bottom: 10px;">
        <el-input type="textarea" v-model="from.input" rows="10" resize="none" :placeholder="info"
            style="height: 40vh; margin-right: 5px;"></el-input>
        <el-space direction="vertical">
            <el-button :icon="MapLocation" @click="cdncheck" size="large">域名CDN查询</el-button>
            <el-button @click="searchdomain" size="large">
                <template #icon>
                    <img src="/chinaz.ico">
                </template>
                域名信息查询
            </el-button>
        </el-space>

    </div>
    <el-input type="textarea" v-model="from.result" rows="10" readonly resize="none" style="height: 50vh;"></el-input>
</template>

<style>
.el-textarea__inner {
    height: 100%;
}
</style>