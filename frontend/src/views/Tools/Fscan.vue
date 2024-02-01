<script lang="ts" setup>
import { reactive } from 'vue'
import { Fscan2Txt } from '../../../wailsjs/go/main/App'

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
</script>


<template>
    <el-form :model="fscan">
        <el-form-item>
            <div class="head">
                <el-input v-model="fscan.content" :rows="7" resize='none' type="textarea"
                    placeholder='请输入fscan扫描的结果内容' />
                <el-button type="primary" style=" margin-left: 10px; height: 40px;"
                    @click="extractresult">提取关键结果</el-button>
            </div>
        </el-form-item>
        <el-form-item>
            <el-input v-model="fscan.result" :rows="20" type="textarea" placeholder='可提取内容如下:
FTP等协议暴破成功字段
MS17-010
POC字段
DC主机
INFO信息
Vcenter主机
海康摄像头主机' />
        </el-form-item>
    </el-form>
</template>