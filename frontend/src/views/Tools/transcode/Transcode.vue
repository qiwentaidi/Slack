<script lang="ts" setup>
import { reactive } from 'vue';
import { MagicStick, Key, QuestionFilled } from '@element-plus/icons-vue';
import OA from './OA.vue';
import Nokey from './Nokey.vue';
import CryptoJS from 'crypto-js';
import JSEncrypt from 'jsencrypt';
import { sm2, sm4 } from 'sm-crypto';
const yk = reactive({
    mode: 'AES',
    options: ["AES", "DES", "RSA", "SM2", "SM4"],
})
const aes = reactive({
    encryptedMessage: '',
    decryptedMessage: '',
    mode: ['ECB', 'CBC', 'CTR', 'OFB', 'CFB'],
    currentMode: 'ECB',
    padding: ['PKCS7', 'ZeroPadding', 'ISO10126', 'ISO97971', 'Ansix923', 'None'],
    currentPadding: 'PKCS7',
    key: '',
    iv: '',
    ende(mode: number) {
        aes.decryptedMessage == ""
        let options = {
            iv: CryptoJS.enc.Utf8.parse(aes.iv),
            mode: CryptoJS.mode.ECB,
            padding: CryptoJS.pad.Pkcs7
        };
        switch (aes.currentMode) {
            case "ECB":
                options.mode = CryptoJS.mode.ECB
            case "CBC":
                options.mode = CryptoJS.mode.CBC
            case "CTR":
                options.mode = CryptoJS.mode.CTR
            case "OFB":
                options.mode = CryptoJS.mode.OFB
            case "CFB":
                options.mode = CryptoJS.mode.CFB
        }
        switch (aes.currentPadding) {
            case "PKCS7":
                options.padding = CryptoJS.pad.Pkcs7
            case "Ansix923":
                options.padding = CryptoJS.pad.AnsiX923
            case "ISO10126":
                options.padding = CryptoJS.pad.Iso10126
            case "ZeroPadding":
                options.padding = CryptoJS.pad.ZeroPadding
            case "ISO97971":
                options.padding = CryptoJS.pad.Iso97971
            case "None":
                options.padding = CryptoJS.pad.NoPadding
        }
        let key = CryptoJS.enc.Utf8.parse(aes.key);
        if (yk.mode == "AES"){
            if (mode == 0) {
                let encryptedData = CryptoJS.AES.encrypt(aes.encryptedMessage, key, options);
                aes.decryptedMessage = encryptedData.toString();
            } else {
                try {
                    let encryptedData = CryptoJS.AES.decrypt(aes.encryptedMessage, key, options);
                    aes.decryptedMessage = encryptedData.toString(CryptoJS.enc.Utf8);
                } catch (error) {
                    aes.decryptedMessage = "无法解密";
                }
            }
        }else {
            if (mode == 0) {
                let encryptedData = CryptoJS.DES.encrypt(aes.encryptedMessage, key, options);
                aes.decryptedMessage = encryptedData.toString();
            } else {
                try {
                    let encryptedData = CryptoJS.DES.decrypt(aes.encryptedMessage, key, options);
                    aes.decryptedMessage = encryptedData.toString(CryptoJS.enc.Utf8);
                } catch (error) {
                    aes.decryptedMessage = "无法解密";
                }
            }
        }
    }
})
const rsa = reactive({
    publicKey: '',
    privateKey: '',
    encryptedMessage: '',
    decryptedMessage: '',
    ende(mode: number) {
        if (mode == 0){
            try {
                // 创建一个 JSEncrypt 实例
                const encryptor = new JSEncrypt();
                // 设置公钥
                encryptor.setPublicKey(this.publicKey);
                const encryptedData = encryptor.encrypt(this.encryptedMessage);
                this.decryptedMessage = encryptedData.toString()
            } catch (error: any) {
                rsa.decryptedMessage = "加密失败" + error.message;
            }
        }else {
            try {
                const encryptor = new JSEncrypt();
                encryptor.setPrivateKey(this.privateKey);
                const encryptedData = encryptor.decrypt(this.encryptedMessage);
                this.decryptedMessage = encryptedData.toString()
            } catch (error: any) {
                rsa.decryptedMessage = "解密失败";
            }
        }
    },
})
const sm = reactive({
    sm2encrypt: '',
    sm2decrypt: '',
    sm2publicKey: '',
    sm2privateKey: '',
    sm2CurrentMode: 'C1C2C3',
    sm2Mode: [ 'C1C2C3','C1C3C2',],
    sm2ende(mode: number){
        if (mode == 0) {
            try {
                if (this.sm2CurrentMode == 'C1C2C3'){
                    sm.sm2decrypt = sm2.doEncrypt(sm.sm2encrypt,sm.sm2publicKey,0);
                }else {
                    sm.sm2decrypt = sm2.doEncrypt(sm.sm2encrypt,sm.sm2publicKey,1);
                }
            } catch (error) {
                sm.sm2decrypt = "加密失败";
            }
        } else {
            try {
                if (this.sm2CurrentMode == 'C1C2C3'){
                    sm.sm2decrypt = sm2.doDecrypt(sm.sm2encrypt, sm.sm2privateKey,0);
                }else {
                    sm.sm2decrypt = sm2.doDecrypt(sm.sm2encrypt,sm.sm2publicKey,1);
                }
            } catch (error) {
                sm.sm2decrypt = "解密失败";
            }
        }
    },
    sm4encrypt: '',
    sm4decrypt: '',
    sm4key: '',
    sm4ende(mode: number) {
        if (mode == 0) {
            try {
                sm.sm4decrypt = sm4.encrypt(sm.sm4encrypt, sm.sm4key);
            } catch (error) {
                sm.sm4decrypt = "加密失败";
            }
        } else {
            try {
                sm.sm4decrypt = sm4.decrypt(sm.sm4encrypt, sm.sm4key);
            } catch (error) {
                sm.sm4decrypt = "解密失败";
            }
        }
    }
})

// mode 0 encode 1 decode
function handleButtonClick(mode: number) {
    switch (yk.mode) {
        case "AES":
            aes.ende(mode)
            break;
        case "DES":
            aes.ende(mode)
            break;
        case "RSA":
            rsa.ende(mode)
            break;
        case "SM4":
            sm.sm4ende(mode)
            break;
        case "SM2":
            sm.sm2ende(mode)
            break;
    }
}
</script>


<template>
    <el-card style="width: 100%;">
        <template #header>
            <div class="card-header">
                <span>编码转换</span>
            </div>
        </template>
        <el-tabs type="border-card" style="height: 570px;">
            <el-tab-pane label="无密钥编码">
                <Nokey></Nokey>
            </el-tab-pane>
            <el-tab-pane label="有密钥编码">
                <el-space style="margin-bottom: 10px;">
                    <el-select class="choose" v-model="yk.mode">
                        <el-option v-for="item in yk.options" :value="item" :label="item" />
                    </el-select>
                    <el-button color="#40CEED" :icon="MagicStick" @click="handleButtonClick(0)">加密</el-button>
                    <el-button color="#59EDAF" :icon="Key" @click="handleButtonClick(1)">解密</el-button>
                    <el-tooltip content="加密或者解密均是从内容到结果，RSA加解密失败会出现false" placement="right-start">
                        <el-icon>
                            <QuestionFilled />
                        </el-icon>
                    </el-tooltip>
                </el-space>
                <div style="width: 100%;" v-if="yk.mode == 'AES' || yk.mode == 'DES'">
                    <el-form :model="aes">
                        <el-form-item label="内容：">
                            <el-input type="textarea" v-model="aes.encryptedMessage" rows="9" resize='none'></el-input>
                        </el-form-item>
                        <el-form-item label="模式：">
                            <el-space>
                                <el-select v-model="aes.currentMode">
                                    <el-option v-for="item in aes.mode" :value="item" :label="item" />
                                </el-select>
                                <el-select v-model="aes.currentPadding">
                                    <el-option v-for="item in aes.padding" :value="item" :label="item" />
                                </el-select>
                                <el-input v-model="aes.key" placeholder="KEY密钥" style="width: 150px;"></el-input>
                                <el-input v-model="aes.iv" placeholder="IV偏移量" style="width: 140px;"
                                    v-if="aes.currentMode != 'ECB'"></el-input>
                            </el-space>
                        </el-form-item>
                        <el-form-item label="结果：">
                            <el-input type="textarea" v-model="aes.decryptedMessage" rows="9" resize='none'></el-input>
                        </el-form-item>
                    </el-form>
                </div>
                <div style="width: 100%;" v-else-if="yk.mode == 'RSA'">
                    <el-form v-model="rsa">
                        <el-form-item label="公钥">
                            <el-input  class="form-inline" v-model="rsa.publicKey" type="textarea" placeholder="-----BEGIN PUBLIC KEY-----
-----END PUBLIC KEY-----" resize='none' rows="2"></el-input>
                        </el-form-item>
                        <el-form-item label="私钥">
                            <el-input  class="form-inline" v-model="rsa.privateKey" type="textarea" placeholder="-----BEGIN RSA PRIVATE KEY-----
-----END RSA PRIVATE KEY-----" resize='none' rows="2"></el-input>
                        </el-form-item>
                        <el-form-item label="内容">
                            <el-input class="form-inline" v-model="rsa.encryptedMessage" type="textarea" resize='none' rows="7"></el-input>
                        </el-form-item>
                        <el-form-item label="结果">
                            <el-input  class="form-inline" v-model="rsa.decryptedMessage" type="textarea" resize='none' rows="7"></el-input>
                        </el-form-item>
                    </el-form>
                </div>
                <div style="width: 100%;" v-else-if="yk.mode == 'SM4'">
                    <el-form :model="sm">
                        <el-form-item label="内容：">
                            <el-input type="textarea" v-model="sm.sm4encrypt" rows="9" resize='none'></el-input>
                        </el-form-item>
                        <el-form-item label="密钥：">
                            <el-input v-model="sm.sm4key" placeholder="key长度需要满足128bits即32位，否则会加密失败"></el-input>
                        </el-form-item>
                        <el-form-item label="结果：">
                            <el-input type="textarea" v-model="sm.sm4decrypt" rows="9" resize='none'></el-input>
                        </el-form-item>
                    </el-form>
                </div>
                <div style="width: 100%;" v-else-if="yk.mode == 'SM2'">
                    <el-form :model="sm">
                        <el-form-item label="内容：">
                            <el-input type="textarea" v-model="sm.sm2encrypt" rows="6" resize='none'/>
                        </el-form-item>
                        <el-form-item label="公钥：">
                            <el-input v-model="sm.sm2publicKey"/>
                        </el-form-item>
                        <el-form-item label="私钥：">
                            <el-input v-model="sm.sm2privateKey"/>
                        </el-form-item>
                        <el-form-item label="模式：">
                            <el-select v-model="sm.sm2CurrentMode" placeholder="Select">
                                <el-option v-for="item in sm.sm2Mode" :value="item" :label="item" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="结果：">
                            <el-input type="textarea" v-model="sm.sm2decrypt" rows="7" resize='none'/>
                        </el-form-item>
                    </el-form>
                </div>
            </el-tab-pane>

            <el-tab-pane label="OA密码解密">
                <OA></OA>
            </el-tab-pane>
        </el-tabs>
    </el-card>
</template>

<style scoped>
.form-inline {
  width: 100%;
}
</style>