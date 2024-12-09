<script lang="ts" setup>
import { ref } from 'vue'
import { AntivirusIdentify, PatchIdentify } from 'wailsjs/go/core/Tools'
import { ElMessage } from 'element-plus';
import { structs } from 'wailsjs/go/models';

const activeName = ref('1')
const textarea = ref('')
const avlist = ref(<structs.AntivirusResult[]>[])
const patch = ref(<structs.AuthPatch[]>[])

async function checkAVList() {
    let result = await AntivirusIdentify(textarea.value)
    if (result.length == 0) {
        ElMessage('未识别到杀软')
        avlist.value = []
        return
    }
    avlist.value = result
}

async function checkPatchList() { 
    let result = await PatchIdentify(textarea.value)
    if (result.length == 0) {
        ElMessage('未查询到提权补丁')
        patch.value = []
        return
    }
    patch.value = result
}
</script>

<template>
    <div class="my-header">
        <el-input v-model="textarea" :rows="5" type="textarea" placeholder="tasklist /svc or systeminfo" />
        <div>
            <el-col :span="2" style="margin-left: 10px;">
                <el-button color="#626aef" @click="checkAVList">杀软识别</el-button>
            </el-col>
            <el-col :span="2" style="margin-left: 10px; margin-top: 10px;">
                <el-button color="#1e3864" @click="checkPatchList">补丁检测</el-button>
            </el-col>
        </div>
    </div>
    <el-tabs v-model="activeName" type="border-card" style="margin-top: 10px;">
        <el-tab-pane label="杀软识别" name="1">
            <el-table :data="avlist" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="Process" label="进程名称" />
                <el-table-column prop="Pid" label="PID" />
                <el-table-column prop="Name" label="杀软名称" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="补丁检测" name="2">
            <el-table :data="patch" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="MS" label="微软编号" />
                <el-table-column prop="Patch" label="补丁编号" />
                <el-table-column prop="Description" label="描述" />
                <el-table-column prop="System" label="影响系统" />
                <el-table-column prop="Reference" label="链接" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
    </el-tabs>
</template>