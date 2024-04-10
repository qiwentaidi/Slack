<script lang="ts" setup>
import { reactive } from 'vue';
import {
    ExtractIP,
    DomainInfo,
    CheckCdn,
} from '../../../wailsjs/go/main/App'
import global from "../../global"
const info = `IP提取:
输入任意内容会自动匹配IPV4地址会进行提取并统计C段数量

域名解析(CDN查询):
输入任意内容会自动匹配域名进行IP解析以及CDN判断(可批量)
===如果域名解析的非常慢，请考虑是否是本机网络不佳===

域名备案&whois查询(数据来源: 站长之家)(不可批量)

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
});

function format() {
    if (from.current == 0) {
        from.result = convertCIDR(from.input)
    }else {
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

function searchdomain(){
    DomainInfo(from.input).then(
        result => {
            from.result = result
        }
    )
}

function cdncheck(){
    CheckCdn(from.input,global.scan.dns1,global.scan.dns2).then(
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
</script>

<template>
    <div>
        <el-input type="textarea" v-model="from.input" rows="10" :placeholder="info" style="height: 45vh;"></el-input>
        <div class="my-header" style="margin-top: 10px; margin-bottom: 10px;">
            <el-space>
                <el-button type="primary" @click="cdncheck">域名解析(CDN查询)</el-button>
                <el-button type="primary" @click="searchdomain">域名备案&whois查询</el-button>
            </el-space>
            <el-space>
                <el-button color="#5FF5B6" @click="extract">IP筛选统计</el-button>
                <el-button color="#5FF5B6" @click="format">IP转换</el-button>
                <el-select v-model="from.current" style="width: 350px;">
                    <el-option v-for="item in from.options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
            </el-space>
        </div>
        <el-input type="textarea" v-model="from.result" rows="10" readonly style="height: 45vh;"></el-input>
    </div>
</template>

<style>
.el-textarea__inner {
  height: 100%;
}
</style>