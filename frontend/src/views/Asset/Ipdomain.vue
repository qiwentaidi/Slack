<script lang="ts" setup>
import { reactive } from 'vue';
import { ICPInfo, CheckCdn, Ip138IpHistory, Ip138Subdomain } from '../../../wailsjs/go/main/App'
import { ElMessage } from 'element-plus';

const domain = reactive({
    input: '',
    icp: '',
    cdn: '',
    history: '',
    subdomain: '',
});

const domainRegex = /^(?=.{1,253}$)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)+[a-z]{2,}$/i;
function startSearch() {
    if (!domainRegex.test(domain.input)) {
        ElMessage.warning("请输入正确域名")
        return
    }
    CheckCdn(domain.input).then(result => {
        domain.cdn = result
    })
    ICPInfo(domain.input).then(
        result => {
            domain.icp = result
        }
    )
    Ip138IpHistory(domain.input).then(
        result => {
            domain.history = result
        }
    )
    Ip138Subdomain(domain.input).then(
        result => {
            domain.subdomain = result
        }
    )
}
</script>

<template>
    <div style="margin-bottom: 10px; display: flex; justify-content: center;">
        <el-input v-model="domain.input" style="width: 60%; height: 40px;">
            <template #prepend>
                域名
            </template>
            <template #suffix>
                <el-button type="primary" link @click="startSearch()">开始查询</el-button>
            </template>
        </el-input>
    </div>
    <div class="grid-container">
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/chinaz.ico"></el-icon><span class="title">备案信息:</span></el-text>
            <el-input v-model="domain.icp" type="textarea" rows="12" resize="none" style="margin-top: 10px;"></el-input>
        </el-card>
        <el-card shadow="never" class="grid-item">
            <span class="title">CDN信息:</span>
            <el-input v-model="domain.cdn" type="textarea" rows="12" resize="none" style="margin-top: 10px;"></el-input>
        </el-card>
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/ip138.ico"></el-icon><span class="title">子域名:</span></el-text>
            <el-input v-model="domain.subdomain" type="textarea" rows="12" resize="none"
                style="margin-top: 10px;"></el-input>
        </el-card>
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/ip138.ico"></el-icon><span class="title">历史解析:</span></el-text>
            <el-input v-model="domain.history" type="textarea" rows="12" resize="none"
                style="margin-top: 10px;"></el-input>
        </el-card>
    </div>
</template>

<style scoped>
.title {
    margin-left: 5px;
    font-weight: bold;
    font-size: 18px;
}

.grid-container {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    gap: 10px;
    width: 100%;
    /* 宽度可以根据需要调整 */
    height: 98%;
    /* 高度可以根据需要调整 */
}

.grid-item {
    justify-content: center;
    align-items: center;
    border-radius: 5px;
    background-color: #f5f5f5;
}
</style>