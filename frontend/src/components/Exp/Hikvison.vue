<template>
    <div class="head">
        <el-input v-model="form.url" type="textarea" rows="5" resize="none" placeholder="请输入URL根路径" style="width: 90%;">
        </el-input>
        <el-space direction="vertical" style="margin-left: 5px; align-items:start;">
            <el-select placeholder="请选择漏洞" v-model="form.selectVulnerability" style="width: 200px">
                <el-option v-for="vulnerability in vulnerabilityGroup" :value="vulnerability.value"
                    :label="vulnerability.name">
                    {{ vulnerability.name }} <el-tag v-show="vulnerability.isBatch">可批量</el-tag>
                </el-option>
            </el-select>
            <el-button type="primary" style="width: 200px;">执行任务</el-button>
            <el-button :icon="Picture" @click="getSnapshot" :disabled="form.selectVulnerability != 1"
                style="width: 200px;">查看快照</el-button>
        </el-space>
    </div>
    <pre class="pretty-response" style="height: 65vh;"><code>{{ form.content }}</code></pre>
    <el-dialog v-model="form.snapshotDialog" title="快照" width="530">
        <el-image style="width: 500px; height: 500px" :src="form.image" loading="eager" fit="fill" />
    </el-dialog>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { Picture } from '@element-plus/icons-vue';
import { HikvsionProduct } from '../../../wailsjs/go/main/App';
import global from '../../global';

const form = reactive({
    url: "",
    image: "",
    snapshotDialog: false,
    selectVulnerability: 1,
    content: "",
})

const vulnerabilityGroup = [
    {
        name: "CVE-2017-7921",
        isBatch: false,
        value: 1
    },
    {
        name: "CVE-2021-36260",
        isBatch: false,
        value: 2
    },
    {
        name: "弱口令检测",
        isBatch: true,
        value: 3
    },
]

async function useVulnerability() {
    // "http://47.150.37.246:81/"
    form.content = await HikvsionProduct(form.url, form.selectVulnerability, global.proxy)
}

async function getSnapshot() {
    form.snapshotDialog = true
    let body = await HikvsionProduct("http://47.150.37.246:81/", 0, global.proxy)
    const base64Image = `data:image/jpeg;base64,${body}`;
    form.image = base64Image
}
</script>