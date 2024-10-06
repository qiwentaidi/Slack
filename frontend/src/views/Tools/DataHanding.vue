<script lang="ts" setup>
import { reactive } from 'vue'
import { ExtractIP, Fscan2Txt, IpLocation } from 'wailsjs/go/main/App'
import { Search, QuestionFilled, Location, Cellphone, Postcard, Upload } from '@element-plus/icons-vue';
import { ElNotification } from 'element-plus';
import extractIcon from '@/assets/icon/extract.svg'
import { SplitTextArea, UploadFileAndRead } from '@/util';
import async from 'async';
import { regexpIdCard, regexpPhone } from '@/stores/validate';

const form = reactive({
    result: '',
    input: '',
    dedupSplit: '',
})

function FscanExtract() {
    if (form.input) {
        Fscan2Txt(form.input).then(
            result => {
                form.result = result
            }
        )
    }
}

async function uploadFile() {
    form.input = await UploadFileAndRead()
}

function Deduplication() {
    let lines = [] as string[]
    if (form.dedupSplit == "\\n") {
        lines = form.input.split(/[(\r\n)\r\n]+/) // 根据换行或者回车进行识别
    } else {
        lines = form.input.split(form.dedupSplit)
    }
    lines = lines.filter(item => item.trim() !== '') // 删除空项并去除左右空格
    let uniqueArray = Array.from(new Set(lines))
    if (lines.length === uniqueArray.length) {
        ElNotification('不存在重复数据')
        return
    }
    ElNotification.success(`已去重数据${lines.length - uniqueArray.length}条`)
    form.result = uniqueArray.join('\n')
}

function extractDomains() {
    const urls = getURLs();
    const domains = urls.map(url => {
        try {
            const parsedUrl = new URL(url);
            return parsedUrl.hostname;
        } catch (e) {
            console.error(`Invalid URL: ${url}`);
            return '';
        }
    }).filter(domain => domain); // 过滤掉空字符串
    form.result = Array.from(new Set(domains)).join('\n')
}

function getURLs(): string[] {
    const urlPattern = /https?:\/\/[^\s/$.?#].[^\s]*/g;
    const urls = form.input.match(urlPattern) || [];
    return Array.from(new Set(urls))
}

function getPhoneNumbers(): string[] {
    const phoneNumbers = form.input.match(regexpPhone) || [];
    return Array.from(new Set(phoneNumbers))
}

function getIdCards(): string[] {
    const idcards = form.input.match(regexpIdCard) || [];
    return Array.from(new Set(idcards))
}


const appControl = [
    {
        label: "IP提取",
        type: "success",
        icon: extractIcon,
        action: () => {
            if (form.input != "") {
                ExtractIP(form.input).then(
                    result => {
                        form.result = result
                    }
                )
            }
        },
    },
    {
        label: "IP定位查询",
        type: "success",
        icon: Location,
        action: () => {
            let lines = SplitTextArea(form.input)
            form.result = ""
            async.eachLimit(lines, 20, async (ip: string, callback: () => void) => {
                let result = await IpLocation(ip)
                form.result += `${ip}  |  ${result}\n`
            })
        },
    },
    {
        label: "提取URL中的域名",
        type: "info",
        icon: extractIcon,
        action: () => {
            extractDomains()
        },
    },
    {
        label: "URL提取",
        type: "info",
        icon: extractIcon,
        action: () => {
            form.result = getURLs().join('\n')
        },
    },
    {
        label: "手机号提取",
        type: "warning",
        icon: Cellphone,
        action: () => {
            form.result = getPhoneNumbers().join('\n')
        },
    },
    {
        label: "身份证提取",
        type: "warning",
        icon: Postcard,
        action: () => {
            form.result = getIdCards().join('\n')
        },
    },
]
</script>


<template>
    <el-scrollbar height="92vh">
        <el-main>
            <el-form label-width="50px">
                <el-form-item label="内容">
                    <el-input v-model="form.input" type="textarea" :rows="7" placeholder='请输入内容' />
                    <el-button link size="small" :icon="Upload" @click="uploadFile"
                        style="margin-top: 5px;">导入文件内容</el-button>
                </el-form-item>
                <el-form-item label="结果">
                    <el-input v-model="form.result" type="textarea" :rows="15" />
                </el-form-item>

                <el-form-item>
                    <el-space>
                        <el-button @click="FscanExtract" type="primary">
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
                        <el-button v-for="item in appControl" @click="item.action" :type="item.type">
                            <template #icon>
                                <el-icon :size="20">
                                    <component :is="item.icon" />
                                </el-icon>
                            </template>
                            {{ item.label }}
                        </el-button>
                    </el-space>
                    <el-input v-model="form.dedupSplit" style="width: 300px; margin-top: 10px;">
                        <template #prepend>
                            数据去重
                            <el-tooltip placement="left">
                                <template #content>输入分隔字符后转换成数组，然后去重，换行输入\n</template>
                                <el-icon>
                                    <QuestionFilled size="24" />
                                </el-icon>
                            </el-tooltip>
                        </template>
                        <template #suffix>
                            <el-button :icon="Search" link @click="Deduplication"></el-button>
                        </template>
                    </el-input>
                </el-form-item>
            </el-form>
        </el-main>
    </el-scrollbar>
</template>
