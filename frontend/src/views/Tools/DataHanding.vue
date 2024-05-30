<script lang="ts" setup>
import { reactive } from 'vue'
import { ExtractIP, Fscan2Txt, SelectFile } from '../../../wailsjs/go/main/App'
import { GetFileContent } from '../../../wailsjs/go/main/File'
import { Delete, Files } from '@element-plus/icons-vue';
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

// async function HunterCSVremoveDuplicates() {
//     let filepath = await SelectFile()
//     let result = await HunterRemoveDuplicates(filepath)
//     if (result) {
//         ElNotification({
//             message: `已去重数据${result.length}条`,
//             type: 'success'
//         })
//     }
// }


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

            <el-button @click="extract" style="width: 250px">
                <template #icon>
                    <el-tooltip placement="right">
                        <template #content>IP提取:
                            <br />输入任意内容会自动匹配IPV4地址会进行提取并统计C段数量</template>
                        <QuestionFilled />
                    </el-tooltip>
                </template>
                IP提取
            </el-button>
            <!-- <el-button @click="HunterCSVremoveDuplicates" style="width: 250px">
                <template #icon>
                    <el-tooltip placement="right">
                        <template #content>去重Hunter从WEB端导出的CSV数据
                            <br />通过IP+端口+网站标题确定唯一性</template>
                        <QuestionFilled />
                    </el-tooltip>
                </template>
                Hunter资产去重
            </el-button> -->
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
    <el-input v-model="form.result" :rows="5" type="textarea" resize="none" style="height: 70%; margin-top: 10px;" />
</template>
<style>
.el-textarea__inner {
    height: 100%;
}

.fieldset {
    border-radius: 5px;
}
</style>