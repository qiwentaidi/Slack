<script lang="ts" setup>
import { reactive } from 'vue'
import { ExtractIP, Fscan2Txt } from '../../../wailsjs/go/main/App'
import { ReadFile, FileDialog } from '../../../wailsjs/go/main/File'
import { Search } from '@element-plus/icons-vue';
import { ElMessage, ElNotification } from 'element-plus';
import { File } from '../../interface';

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

function extractUrls() {
    form.result = Array.from(new Set(getURLs())).join('\n')
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

</script>


<template>
    <div class="head">
        <ContextMenu>
            <template #menu>
                <div class="nav-item" @click="uploadFile">
                    <img src="../../assets/icon/upload.svg" style="margin-left: 10px;">
                    <span class="nav-text">上传文件</span>
                </div>
            </template>
            <el-input v-model="form.input" resize='none' type="textarea" placeholder='粘贴文件内容或者右键上传'
                style="height: 100%;" />
        </ContextMenu>
        <el-space direction="vertical" style="margin-left: 5px; width: 25%; align-items:start;">
            <el-button @click="FscanExtract" style="width: 300px;" type="primary">
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
            <el-button @click="extract" style="width: 300px" type="success">
                IP提取
            </el-button>
            <el-button @click="extractDomains" style="width: 300px" type="warning">
                提取URL中的域名
            </el-button>
            <el-button @click="extractUrls" style="width: 300px" type="info">
                URL提取
            </el-button>
            <el-input v-model="form.dedupSplit" style="width: 300px">
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
    </div>
    <el-input v-model="form.result" type="textarea" style="height: 100%; margin-top: 10px;" />
</template>
<style>
.el-textarea__inner {
    height: 100%;
}

.fieldset {
    border-radius: 5px;
}
</style>