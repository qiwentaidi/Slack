<script lang="ts" setup>
import { reactive, ref } from 'vue';
import {
    Transcoding
} from '../../../../wailsjs/go/main/App'
import { ArrowDownBold, Close, Plus} from '@element-plus/icons-vue';
const autoEncrypt = ref(false);
const transcode = reactive({
    Encrypt: '',
    Decrypt: '',
    divs: [{ nkMode: 'Base64', nkOptions: ["Base64", "Base64 URL", "URLcode", "Unicode", "HEX", "HTML", "Ascii"] }],
})
function addDiv() {
    transcode.divs.push({ nkMode: 'Base64', nkOptions: ["Base64", "Base64 URL", "URLcode", "Unicode", "HEX", "HTML", "Ascii"] });
}
function removeDiv(index: number) {
    transcode.divs.splice(index, 1);
}

function handleButtonClick(mode: number) {
    // 获取所有el-select的当前选项值
    const crypt = transcode.divs.map(div => div.nkMode);
    Transcoding(crypt, transcode.Encrypt, mode).then(
        result => {
            transcode.Decrypt = result;
        }
    )
}

function handleInputChange() {
    if (transcode.Encrypt == '') {
        transcode.Decrypt = ''
        return
    }
    if (autoEncrypt.value) {
        const crypt = transcode.divs.map(div => div.nkMode);
        Transcoding(crypt, transcode.Encrypt, 0).then(
            result => {
                transcode.Decrypt = result;
            }
        )
    }
}
</script>


<template>
    <div class="flex-box">
        <div style="width: 235px;">
            <el-scrollbar height="460px" style="margin-top: 5px;">
                <div class="flex-box" v-for="(div, index) in transcode.divs"
                    :key="index">
                    <el-select v-model="div.nkMode">
                        <el-option v-for="item in div.nkOptions" :value="item" :label="item" />
                    </el-select>
                    <el-button :icon="Close" @click="removeDiv(index)"></el-button>
                </div>
            </el-scrollbar>
        </div>
        <div style="width: 100%; margin-left: 15px;">
            <el-input type="textarea" v-model="transcode.Encrypt" rows="10" resize='none'
                @input="handleInputChange" style="margin-bottom: 10px;"></el-input>
            <div class="my-header">
                <el-button color="#626aef" :icon="Plus" @click="addDiv">添加</el-button>
                <div class="flex-box">
                    <el-checkbox v-model="autoEncrypt" label="自动加密" border />
                    <el-button-group style="margin-left: 5px;">
                        <el-button type="primary" :icon="ArrowDownBold" @click="handleButtonClick(0)">加密</el-button>
                        <el-button type="primary" :icon="ArrowDownBold" @click="handleButtonClick(1)">解密</el-button>
                    </el-button-group>
                </div>
            </div>
            <el-input type="textarea" v-model="transcode.Decrypt" rows="10" resize='none'
                style="margin-top: 10px;"></el-input>
        </div>
    </div>     
</template>

<style>
.nkmode {
 display: flex;
 justify-content: space-between;
}
</style>