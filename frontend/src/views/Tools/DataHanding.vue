<script lang="ts" setup>
import { reactive } from 'vue'
import { ExtractIP, Fscan2Txt, IpLocation } from 'wailsjs/go/main/App'
import { ReadFile, FileDialog } from 'wailsjs/go/main/File'
import { Search, QuestionFilled, Location } from '@element-plus/icons-vue';
import { ElMessage, ElNotification } from 'element-plus';
import { File } from '@/interface';
import extractIcon from '@/assets/icon/extract.svg'
import { SplitTextArea } from '@/util';
import async from 'async';

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
    let filepath = await FileDialog("*.txt")
    if (!filepath) {
        return
    }
    let file: File = await ReadFile(filepath)
    if (file.Error) {
        ElMessage({
            type: "warning",
            message: file.Message
        })
        return
    }
    form.input = file.Content!
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
    const urls = form.input.match(urlPattern);
    return urls ? urls : [];
}

const menus = [
    {
        label: "上传文件",
        click: uploadFile,
    }
]

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
        type: "warning",
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
            form.result = Array.from(new Set(getURLs())).join('\n')
        },
    },
]
</script>


<template>
    <el-scrollbar height="92vh">
        <el-main>
            <el-form label-width="50px">
                <el-form-item label="内容">
                    <el-input v-model="form.input" type="textarea" :rows="7" placeholder='请输入内容或者右键上传文件'
                        v-menus:right="menus" />
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
                        <el-input v-model="form.dedupSplit" style="width: 300px;">
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
                    </el-space>
                </el-form-item>
            </el-form>
        </el-main>
    </el-scrollbar>
</template>
<style>
.el-textarea__inner {
    height: 100%;
}

.fieldset {
    border-radius: 5px;
}
</style>