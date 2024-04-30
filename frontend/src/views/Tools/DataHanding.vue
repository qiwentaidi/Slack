<script lang="ts" setup>
import { reactive } from 'vue'
import { ExtractIP, Fscan2Txt, SelectFile } from '../../../wailsjs/go/main/App'
import { GetFileContent } from '../../../wailsjs/go/main/File'
import { Delete, Files, Switch, Setting } from '@element-plus/icons-vue';
import { ElNotification } from 'element-plus';

const form = reactive({
    result: '',
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
    dedupOptions: [
        {
            value: '\n',
            label: '按换行'
        },
        {
            value: ',',
            label: '按逗号'
        }
    ],
    dedupCurrent: '\n',
})

function FscanExtract() {
    if (form.input !== "") {
        Fscan2Txt(form.input).then(
            result => {
                form.result = result
            }
        )
    }
}

async function ReadFile() {
    let filepath = await SelectFile()
    if (filepath != "") {
        form.input = await GetFileContent(filepath)
    }
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

function format() {
    if (form.input != "") {
        if (form.current == 0) {
            form.result = convertCIDR(form.input)
        } else {
            form.result = cidrToRange(form.input)
        }
    }
}

function extract() {
    if (form.input != "") {
        ExtractIP(form.input).then(
            result => {
                form.result = result
            }
        )
    }
}

function Deduplication() {
    let lines = [] as string[]
    if (form.dedupCurrent == "\n") {
        lines = form.input.split(/[(\r\n)\r\n]+/) // 根据换行或者回车进行识别
    } else {
        lines = form.input.split(',')
    }
    lines = lines.filter(item => item.trim() !== '') // 删除空项并去除左右空格
    let uniqueArray = Array.from(new Set(lines))
    if (lines.length === uniqueArray.length) {
        ElNotification('不存在重复数据')
        return
    }
    ElNotification({
        message: `已去重数据${lines.length - uniqueArray.length}条`,
        type: 'success'
    })
    form.result = uniqueArray.join('\n')
}
</script>


<template>
    <div class="head">
        <el-input v-model="form.input" :rows="7" resize='none' type="textarea"
            placeholder='Please input the content or load the file' />
        <el-space direction="vertical" style="margin-left: 5px; width: 25%; align-items:start;">
            <el-button @click="FscanExtract" style="width: 250px;">
                <template #icon>
                    <el-tooltip placement="left">
                        <template #content>可提取内容如下:<br />
                            NetInfo信息<br />
                            FTP等协议暴破成功字段<br />
                            MS17-010<br />
                            POC字段<br />
                            DC主机<br />
                            INFO信息<br />
                            Vcenter主机<br />
                            海康摄像头主机</template>
                        <QuestionFilled />
                    </el-tooltip>
                </template>
                Fscan结果提取
            </el-button>

            <el-button-group>
                <el-button @click="extract" style="width: 125px">
                    <template #icon>
                        <el-tooltip placement="right">
                            <template #content>IP提取:
                                <br />输入任意内容会自动匹配IPV4地址会进行提取并统计C段数量</template>
                            <QuestionFilled />
                        </el-tooltip>
                    </template>
                    IP提取
                </el-button>
                <el-popover placement="left" :width="375" trigger="click">
                    <template #reference>
                        <el-button :icon="Switch" style="width: 125px;">IP转换</el-button>
                    </template>
                    <el-select v-model="form.current">
                        <el-option v-for="item in form.options" :key="item.value" :label="item.label"
                            :value="item.value" @click="format" />
                    </el-select>
                </el-popover>
            </el-button-group>
            <div style="display: flex;">
                <el-button @click="Deduplication" style="width: 125px;">
                    <template #icon>
                        <svg class="bi bi-exclude" width="14px" height="14px" viewBox="0 0 16 16" fill="currentColor"
                            xmlns="http://www.w3.org/2000/svg">
                            <path fill-rule="evenodd"
                                d="M1.5 0A1.5 1.5 0 0 0 0 1.5v9A1.5 1.5 0 0 0 1.5 12H4v2.5A1.5 1.5 0 0 0 5.5 16h9a1.5 1.5 0 0 0 1.5-1.5v-9A1.5 1.5 0 0 0 14.5 4H12V1.5A1.5 1.5 0 0 0 10.5 0h-9zM12 4H5.5A1.5 1.5 0 0 0 4 5.5V12h6.5a1.5 1.5 0 0 0 1.5-1.5V4z" />
                        </svg>
                    </template>
                    数据去重
                </el-button>
                <el-select v-model="form.dedupCurrent" placeholder="选择分隔字符" style="width: 125px;">
                    <el-option v-for="item in form.dedupOptions" :key="item.value" :label="item.label"
                        :value="item.value" />
                </el-select>
            </div>

            <el-space>
                <el-tooltip content="Load File" placement="left">
                    <el-button type="primary" :icon="Files" circle size="large" @click="ReadFile"></el-button>
                </el-tooltip>
                <el-tooltip content="Clear input" placement="left">
                    <el-button type="primary" :icon="Delete" circle size="large" @click=""></el-button>
                </el-tooltip>
            </el-space>
        </el-space>
    </div>
    <el-input v-model="form.result" :rows="5" type="textarea" resize="none" style="height: 75%; margin-top: 10px;" />
</template>
<style>
.el-textarea__inner {
    height: 100%;
}

.fieldset {
    border-radius: 5px;
}
</style>