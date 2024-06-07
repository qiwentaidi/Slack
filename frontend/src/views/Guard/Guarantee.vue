<script setup lang="ts">
import { reactive, ref } from 'vue'

const config = reactive({
    url: '',
    token: '',
    webhook: '',
    secret: '',
})

const checked1 = ref(false)
const checked2 = ref(false)
const checked3 = ref(false)
const checked4 = ref(false)
const checked5 = ref(false)


</script>

<template>
    <el-tabs type="card">
        <el-tab-pane>
            <template #label>
                <el-text>
                    <el-icon><img src="/navigation/ahapt.png" style="width: 17px; height: 17px;"></el-icon>
                    AH-APT(暂未上线)
                </el-text>
            </template>
        </el-tab-pane>
    </el-tabs>
    <div style="display: flex;">
        <el-card style="width: 49%;">
            <template #header>
                <div class="my-header">
                    <span>设置目标</span>
                    <el-button type="primary" class="multi-line-button">
                        开始运行
                    </el-button>
                </div>
            </template>
            <el-form label-width="auto">
                <el-form-item label="URL:" :rules="[
                    { required: true, message: 'url is required' },
                ]">
                    <el-input v-model="config.url"></el-input>
                </el-form-item>
                <el-form-item label="Token:" :rules="[
                    { required: true, message: 'token is required' },
                ]">
                    <el-input v-model="config.token"></el-input>
                </el-form-item>
                <el-form-item label="Webhook:" :rules="[
                    { required: true, message: 'webhook is required' },
                ]">
                    <el-input v-model="config.webhook"></el-input>
                </el-form-item>
                <el-form-item label="Secret:" :rules="[
                    { required: true, message: 'secret is required' },
                ]">
                    <el-input v-model="config.secret"></el-input>
                </el-form-item>
                <el-form-item label="攻击状态">
                    <el-space>
                        <el-check-tag :checked="checked1">全部</el-check-tag>
                        <el-check-tag :checked="checked2">成功</el-check-tag>
                        <el-check-tag :checked="checked3">尝试</el-check-tag>
                        <el-check-tag :checked="checked4">失败</el-check-tag>
                        <el-check-tag :checked="checked5">失陷</el-check-tag>
                    </el-space>
                </el-form-item>
            </el-form>
        </el-card>
        <div class="pretty-response" style="width: 49%; margin-left: 10px;"></div>
    </div>
</template>

<style>
.el-input {
    width: 100%;
}
</style>