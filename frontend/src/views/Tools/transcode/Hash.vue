<script lang="ts" setup>
import { reactive, ref } from 'vue';
import CryptoJS from 'crypto-js';
import { sm3 } from 'sm-crypto';
import { BrowserOpenURL } from '../../../../wailsjs/runtime'
const hashcode = reactive({
    currentMode: 'MD5',
    options: ["MD5", "SHA1", "SHA224", "SHA256", "SHA384", "SHA512", "SHA3", "SM3"],
    et: '',
    ct: ''
})

const autoEncrypt = ref(false);

function encrypt() {
    switch (hashcode.currentMode) {
        case 'MD5':
            hashcode.ct = CryptoJS.MD5(hashcode.et).toString();
            break;
        case 'SHA1':
            hashcode.ct = CryptoJS.SHA1(hashcode.et).toString();
            break;
        case 'SHA224':
            hashcode.ct = CryptoJS.SHA224(hashcode.et).toString();
            break;
        case 'SHA256':
            hashcode.ct = CryptoJS.SHA256(hashcode.et).toString();
            break;
        case 'SHA384':
            hashcode.ct = CryptoJS.SHA384(hashcode.et).toString();
            break;
        case 'SHA512':
            hashcode.ct = CryptoJS.SHA512(hashcode.et).toString();
            break;
        case 'SHA3':
            hashcode.ct = CryptoJS.SHA3(hashcode.et).toString();
            break;
        case 'SM3':
            hashcode.ct = sm3(hashcode.et);
            break;  
    }
}

function handleInputChange() {
    if (hashcode.et == '') {
        hashcode.ct = ''
        return
    }
    if (autoEncrypt.value) {
        encrypt();
    }
}

function handleButtonClick() {
    if (hashcode.et == '') {
        hashcode.ct = ''
        return
    }
    if (!autoEncrypt.value) {
        encrypt();
    }
}

const cmd5 = "https://www.cmd5.com"
</script>

<template>
    <el-card>
        <template #header>
            <div class="card-header">
                <span>Hash</span>
                <el-space>
                    <el-link @click=BrowserOpenURL(cmd5) type="primary" style="font-size: medium;">撞库</el-link>
                    <el-select class="choose" v-model="hashcode.currentMode" >
                        <el-option v-for="item in hashcode.options" :value="item" :label="item"/>
                    </el-select>
                    <el-checkbox v-model="autoEncrypt" label="自动" border />
                    <el-button class="button" type="primary" @click="handleButtonClick">加密</el-button>
                </el-space>
            </div>
        </template>
        <el-scrollbar height="200px">
            <el-input type="textarea" v-model="hashcode.et" rows="4" placeholder="请输入需要加密的内容" resize='none' @input="handleInputChange"></el-input>
            <el-input type="textarea" v-model="hashcode.ct" rows="4" resize='none' readonly></el-input>
        </el-scrollbar>
    </el-card>
</template>

<style>
.choose .el-input__inner {
    width: 80px;
}
</style>