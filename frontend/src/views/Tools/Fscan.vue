<script lang="ts" setup>
import { reactive } from 'vue'
import { Fscan2Txt, SelectFile } from '../../../wailsjs/go/main/App'
import { GetFileContent } from '../../../wailsjs/go/main/File'

const fscan = reactive({
    content: '',
    result: '',
})

function extractresult() {
    if (fscan.content !== "") {
        Fscan2Txt(fscan.content).then(
            result => {
                fscan.result = result
            }
        )
    }
}

async function ReadFile() {
    let filepath = await SelectFile()
    fscan.content = await GetFileContent(filepath)
}
</script>


<template>
    <el-form :model="fscan">
        <el-form-item>
            <div class="head">
                <el-input v-model="fscan.content" :rows="7" resize='none' type="textarea"
                    placeholder='请输入fscan扫描的结果内容' />
                <el-space direction="vertical" style="margin-left: 5px;"> 
                    <el-button @click="ReadFile">加载文件内容</el-button>
                    <el-button type="primary" @click="extractresult">提取关键结果</el-button>
                </el-space>
            </div>
        </el-form-item>
    </el-form>
    <el-input v-model="fscan.result" :rows="5" type="textarea" placeholder='可提取内容如下:
FTP等协议暴破成功字段
MS17-010
POC字段
DC主机
INFO信息
Vcenter主机
海康摄像头主机' resize="none" style="height: 75%;" />
</template>
<style>
.el-textarea__inner {
    height: 100%;
}
</style>