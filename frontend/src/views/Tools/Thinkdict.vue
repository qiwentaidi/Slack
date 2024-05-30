<template>
    <el-container class="el-main">
        <el-row :gutter="20">
            <el-col :span="200">
                <el-form :model="dict" label-position="top">
                    <el-form-item  class="field" label="姓名">
                        <el-input v-model="dict.p_name" />
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
                    <el-form-item>
                        <el-button type="primary" @click="generate">生成</el-button>
                    </el-form-item>
                </el-form>
            </el-col>
            <el-col :span="200">
                <el-form :model="dict" label-position="top">
                    <el-form-item label="结果">
                        <el-input type="textarea" rows="19" v-model="r.result" />
                    </el-form-item>
                </el-form>
            </el-col>
        </el-row>
        <h3>共生成结果:{{r.count}}个</h3>
    </el-container>
</template>
  
<script lang="ts" setup>
import {reactive} from 'vue'
import {ThinkDict} from '../../../wailsjs/go/main/App'
const dict = reactive({
    p_name: '',
    c_name: '',
    c_domain: '',
    p_birthday: '',
    p_worknum: '',
})
const r = reactive({
    result: '',
    count: 0
})

function generate(){
    ThinkDict(dict.p_name,dict.c_name,dict.c_domain,dict.p_birthday,dict.p_worknum).then(
        result => {
            r.result = result.join("\n")
            r.count = result.length
        }
    )
}

</script>
  
<style scoped>
.field {
    width: 300px;
}

.el-main {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
}
</style>