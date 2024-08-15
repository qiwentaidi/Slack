<template>
    <el-card>
        <el-form :inline="true" :model="dict" label-width="auto">
            <el-form-item class="field" label="姓名">
                <el-input v-model="dict.p_cnName" />
            </el-form-item>
            <el-form-item class="field" label="英文名">
                <el-input v-model="dict.p_enName" />
            </el-form-item>
            <el-form-item class="field" label="公司名称">
                <el-input v-model="dict.c_name" />
            </el-form-item>
            <el-form-item class="field" label="公司域名">
                <el-input v-model="dict.c_domain" />
            </el-form-item>
            <el-form-item class="field" label="生日">
                <el-input v-model="dict.p_birthday" />
            </el-form-item>
            <el-form-item class="field" label="工号">
                <el-input v-model="dict.p_worknum" />
            </el-form-item>
            <el-form-item class="field" label="密码连接符">
                <el-input v-model="dict.connect" placeholder="@" />
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="generate" style="margin-left: 85px;">生成</el-button>
                <el-button type="primary" @click="Copy(r.result)" >复制结果</el-button>
            </el-form-item>
        </el-form>
    </el-card>
    <el-row style="margin-top: 10px;">
        <el-col :lg="5">
            拼接字典
            <el-input type="textarea" :rows="25" v-model="dict.easyContent" />
        </el-col>
        <el-col :lg="5" style="margin-left: 10px; margin-right: 10px"> 
            结果
            <el-input type="textarea" :rows="25" v-model="r.result" />
        </el-col>
            <h2>共生成结果:{{ r.count }}个</h2>
    </el-row>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
import { ThinkDict } from 'wailsjs/go/main/App'
import { onMounted } from 'vue';
import { Copy } from '@/util';
import global from '@/global';

onMounted(() => {
    dict.easyContent = dict.easyList.join("\n")
})

const dict = reactive({
    p_cnName: '',
    p_enName: '',
    c_name: '',
    c_domain: '',
    p_birthday: '',
    p_worknum: '',
    connect: '',
    easyList: ["123456", "123", "12345", "1234", "111", "123456789", "12345678"],
    easyContent: '',
})

const r = reactive({
    result: '',
    count: 0
})

function generate() {
    dict.easyList = dict.easyContent.split("\n")
    ThinkDict(dict.p_cnName, dict.p_enName, dict.c_name, dict.c_domain, dict.p_birthday, dict.p_worknum, dict.connect, dict.easyList).then(
        result => {
            global.temp.thinkdict = result
            r.result = result.join("\n")
            r.count = result.length
        }
    )
}

</script>

<style scoped>
.field {
    width: 22%;
}
</style>