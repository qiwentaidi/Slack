<template>
    <div style="height: 100vh;">
        暂未完成
        <!-- <el-steps :active="currentActive" simple finish-status="success">
            <el-step title="输入目标" />
            <el-step title="选定参数" />
            <el-step title="结果输出" />
        </el-steps>
        <div style="display: flex; margin-top: 20px; grid-template-columns: repeat(2, 1fr); gap: 10px;">
            <el-form :mode="target" label-width="auto" style="width: 100%;">
                <el-form-item label="公司名称">
                    <el-input v-model="target" type="textarea" :rows="10"></el-input>
                </el-form-item>
                <el-form-item label="天眼查Token">
                    <el-input v-model="target" type="textarea" :rows="5"></el-input>
                </el-form-item>
            </el-form>
            <div>
                <el-image src="/workflow.jpg" style="width: 100%;"></el-image>
            </div>
        </div>
        <el-button style="float: right; bottom: 0;" @click="next">Next step</el-button> -->
    </div>
</template>


<script lang="ts" setup>
import { ref } from 'vue'
const currentActive = ref(1)
const target = ref('')

const next = () => {
    if (currentActive.value++ > 2) currentActive.value = 1
}


</script>


<style scoped></style>