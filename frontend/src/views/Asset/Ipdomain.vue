<script lang="ts" setup>
import { reactive } from 'vue';
import { MapLocation, Switch } from '@element-plus/icons-vue';
import {
    ExtractIP,
    DomainInfo,
    CheckCdn,
} from '../../../wailsjs/go/main/App'
import global from "../../global"
import { ElNotification } from 'element-plus';
const info = `IP提取:
输入任意内容会自动匹配IPV4地址会进行提取并统计C段数量

域名CDN查询:
输入任意内容会自动匹配域名进行IP解析以及CDN判断(可批量)
===如果域名解析的非常慢，请考虑是否是本机网络不佳===

域名信息查询(数据来源: 站长之家)(不可批量)

IP转换
根据右边列表的形式进行IP转换(不可批量)`;
const from = reactive({
    input: '',
    options: [{
        value: 0,
        label: "192.168.0.0/255.255.255.0 => 192.168.0.0/24",
    },
    {
        value: 1,
        label: "192.168.0.0/24 => 192.168.0.0-192.168.0.255"
    }],
    current: 0,
    result: '',
    dedupOptions: [
        {
            value: '\n',
            label: '换行'
        },
        {
            value: ',',
            label: '逗号'
        }
    ],
    dedupCurrent: '\n',
});

function format() {
    if (from.current == 0) {
        from.result = convertCIDR(from.input)
    } else {
        from.result = cidrToRange(from.input)
    }
}

function extract() {
    ExtractIP(from.input).then(
        result => {
            from.result = result
        }
    )
}

function searchdomain() {
    DomainInfo(from.input).then(
        result => {
            from.result = result
        }
    )
}

function cdncheck() {
    CheckCdn(from.input, global.scan.dns1, global.scan.dns2).then(
        result => {
            from.result = result
        }
    )
}


function convertCIDR(ip: string): string {
    const regex = /(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\/(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/;
    const matches = ip.match(regex);
    if (matches) {
        const ipAddress = matches[1];
        const subnetMask = matches[2];
        const cidr = maskToCidr(subnetMask);
        return `${ip} => ${ipAddress}/${cidr}`;
    }
    return ip;
}

function maskToCidr(mask: string): number {
    let maskArray = mask.split(".");
    // 将每个数字转换成二进制字符串，并拼接成一个32位的二进制字符串
    let maskBinary = maskArray.map((num) => parseInt(num).toString(2).padStart(8, "0")).join("");
    // 统计二进制字符串中1的个数，即为cidr的值
    let cidr = maskBinary.split("").filter((bit) => bit === "1").length;
    // 返回cidr的值
    return cidr;
}

function cidrToRange(cidr: string): string {
    let [network, prefix] = cidr.split("/");
    let networkArray = network.split(".");
    // 将每个数字转换成二进制字符串，并拼接成一个32位的二进制字符串
    let networkBinary = networkArray.map((num) => parseInt(num).toString(2).padStart(8, "0")).join("");
    // 根据网络前缀长度，将网络地址的二进制字符串分为网络部分和主机部分
    let networkPart = networkBinary.slice(0, parseInt(prefix));
    let hostPart = networkBinary.slice(parseInt(prefix));
    // 计算主机部分的取值范围，即从全0到全1
    let minHost = hostPart.replace(/./g, "0");
    let maxHost = hostPart.replace(/./g, "1");
    // 将网络部分和主机部分的最小值和最大值拼接起来，得到子网中的最小地址和最大地址的二进制字符串
    let minBinary = networkPart + minHost;
    let maxBinary = networkPart + maxHost;
    // 将最小地址和最大地址的二进制字符串按每8位分割，然后转换成十进制数字，并用点连接起来，得到子网中的最小地址和最大地址的点分十进制形式
    let minDecimal = minBinary.match(/.{8}/g)!.map((num) => parseInt(num, 2)).join(".");
    let maxDecimal = maxBinary.match(/.{8}/g)!.map((num) => parseInt(num, 2)).join(".");
    let range = minDecimal + "-" + maxDecimal;
    return `${cidr} => ${range}`;
}

function Deduplication() {
    let lines = [] as string[]
    if (from.dedupCurrent == "\n") {
        lines = from.input.split(/[(\r\n)\r\n]+/) // 根据换行或者回车进行识别
    }else {
        lines = from.input.split(',')
    }
    lines = lines.filter(item => item.trim() !== '') // 删除空项并去除左右空格
    let uniqueArray = Array.from(new Set(lines))
    if (lines.length === uniqueArray.length) {
        ElNotification('不存在重复数据')
        return
    }
    ElNotification({
        message: `已去重数据${lines.length-uniqueArray.length}条`,
        type: 'success'
    })
    from.result = uniqueArray.join('\n')
}
</script>

<template>
    <el-input type="textarea" v-model="from.input" rows="10" resize="none" :placeholder="info"
        style="height: 45vh;"></el-input>
    <div class="my-header" style="margin-top: 10px; margin-bottom: 10px;">
        <el-button-group>
            <el-button :icon="MapLocation" @click="cdncheck">域名CDN查询</el-button>
            <el-button @click="searchdomain">
                <template #icon>
                    <img src="/chinaz.ico" width="14" height="14">
                </template>
                域名信息查询
            </el-button>
        </el-button-group>
        <div>
            <el-button @click="Deduplication">
                <template #icon>
                    <svg class="bi bi-exclude" width="14px" height="14px" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                        <path fill-rule="evenodd" d="M1.5 0A1.5 1.5 0 0 0 0 1.5v9A1.5 1.5 0 0 0 1.5 12H4v2.5A1.5 1.5 0 0 0 5.5 16h9a1.5 1.5 0 0 0 1.5-1.5v-9A1.5 1.5 0 0 0 14.5 4H12V1.5A1.5 1.5 0 0 0 10.5 0h-9zM12 4H5.5A1.5 1.5 0 0 0 4 5.5V12h6.5a1.5 1.5 0 0 0 1.5-1.5V4z"/>
                    </svg>
                </template>
                数据去重
            </el-button>
            <el-select v-model="from.dedupCurrent" placeholder="选择分隔字符" style="width: 140px;">
                <el-option v-for="item in from.dedupOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
        </div>
        <div>
            <el-button-group>
                <el-button @click="extract">
                    <template #icon>
                        <svg t="1713165833152" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="9256" width="200" height="200"><path d="M192 800h736l-96 128H96z" fill="#F0F0F0" p-id="9257"></path><path d="M928 768v32l-96 128H96v-32h736z" fill="#175EA2" p-id="9258"></path><path d="M928 736v32l-96 128H96v-32h736z" fill="#3678C1" p-id="9259"></path><path d="M928 704v32l-96 128H96v-32h736z" fill="#4F81DB" p-id="9260"></path><path d="M928 672v32l-96 128H96v-32h736z" fill="#5D8AE7" p-id="9261"></path><path d="M928 640v32l-96 128H96v-32h736z" fill="#669AE9" p-id="9262"></path><path d="M928 608v32l-96 128H96v-32h736z" fill="#7DB4ED" p-id="9263"></path><path d="M192 608h736l-96 128H96z" fill="#95CBF0" p-id="9264"></path><path d="M96 224h64v64H96V224z m160-128h64v64H256V96z m352 128h64v64h-64V224z m224 160h96v96h-96v-96zM192 448h32v32H192v-32z m304-122.816V480a16 16 0 1 1-32 0v-152.032l-36.896 35.552a16 16 0 0 1-22.208-23.04l77.92-75.104 72.704 75.52a16 16 0 1 1-23.04 22.208l-36.48-37.92z" fill="#30AD98" p-id="9265"></path><path d="M416 64h128v128h-128V64zM288 288h64v64H288V288z m0 192h64v64H288v-64z m288 0h64v64h-64v-64z m160-32h32v32h-32v-32z" fill="#27A2DF" p-id="9266"></path></svg>
                    </template>
                    IP提取
                </el-button>
                <el-button :icon="Switch" @click="format">IP转换</el-button>
            </el-button-group>
            <el-select v-model="from.current" style="width: 350px;">
                <el-option v-for="item in from.options" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
        </div>
    </div>
    <el-input type="textarea" v-model="from.result" rows="10" readonly resize="none" style="height: 45vh;"></el-input>
</template>

<style>
.el-textarea__inner {
    height: 100%;
}
</style>