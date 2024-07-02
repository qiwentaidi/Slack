<script lang="ts" setup>
import { reactive } from 'vue';
import { ElMessage } from 'element-plus';
const oa = reactive({
    encrypt: '',
    decrypt: '',
    mode: 'FanruanOA',
    options: ['FanruanOA', "SeeyonOA"],
})
const PASSWORD_MASK_ARRAY: number[] = [19, 78, 10, 15, 100, 213, 43, 23]; // 掩码
function OACrypt(input: string, mode: string) {
    oa.decrypt = ''
    if (mode === "FanruanOA") { // 帆软
        if (input.length > 3) {
            input = input.substring(3); // 截断三位后
            for (let i = 0; i < input.length / 4; i++) {
                let c1: number = parseInt(input.substring(i * 4, (i + 1) * 4), 16);
                let c2: number = c1 ^ PASSWORD_MASK_ARRAY[i % 8];
                oa.decrypt += String.fromCharCode(c2);
            }
        }
    } else {
        let pass = input.replace(/\//g, "");
        let p = pass.split(".0");
        if (p.length > 1) {
            let iv = parseInt(p[0]);
            let password = atob(p[1]);
            for (let i = 0; i < password.length; i++) {
                let char = password.charCodeAt(i);
                oa.decrypt += String.fromCharCode(char - iv);
            }
            return
        }
        ElMessage({
            showClose: true,
            message: '无法解密.',
            type: 'warning',
        })
    }
}
</script>


<template>
    <div style="margin-bottom: 10px;">
        <span>选择解密模式：</span>
        <el-select class="choose" v-model="oa.mode">
            <el-option v-for="item in oa.options" :value="item" :label="item" />
        </el-select>
    </div>
    <el-descriptions title="示例" :column="1" border style="margin-bottom: 10px;">
        <el-descriptions-item label="帆软数据报表 privilege.xml">___0072002a00670066000a</el-descriptions-item>
        <el-descriptions-item label="致远OA datasourceCtp.properties">/1.0/UWJ0dHgxc2U=</el-descriptions-item>
    </el-descriptions>
    <div class="flex-box">
        <el-input v-model="oa.encrypt" resize='none' placeholder="请输入需要解密的密码"></el-input>
        <el-col :span="2" style="text-align: center;">
            <span>=></span>
        </el-col>
        <el-input v-model="oa.decrypt" resize='none'></el-input>
        <el-button class="button" type="warning" plain @click="OACrypt(oa.encrypt, oa.mode)"
            style="margin-left: 10px;">解密</el-button>
    </div>
</template>