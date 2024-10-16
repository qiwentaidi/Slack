<script lang="ts" setup>
import { ref } from 'vue'
import { AntivirusIdentify, PatchIdentify } from 'wailsjs/go/core/Tools'
import { onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { structs } from 'wailsjs/go/models';
// 初始化时调用
onMounted(() => {
    avlist.value = [];
    patch.value = [];
});
async function checksystem(mode: number) {
    if (mode == 0) {
        let result = await AntivirusIdentify(textarea.value)
        if (result.length == 0) {
            ElMessage('未识别到杀软')
            return
        }
        avlist.value = result
        return
    }
    
    let result = await PatchIdentify(textarea.value)
    if (result.length == 0) {
        ElMessage('未查询到提权补丁')
        return
    }
    patch.value = result
}

const activeName = ref('1')
const textarea = ref('')
const avlist = ref(<structs.AntivirusResult[]>[])
const patch = ref(<structs.AuthPatch[]>[])

</script>

<template>
    <div class="my-header">
        <el-input v-model="textarea" :rows="5" type="textarea" placeholder="tasklist /svc or systeminfo" />
        <div>
            <el-col :span="2" style="margin-left: 10px;">
                <el-button color="#626aef" @click="checksystem(0)">杀软识别</el-button>
            </el-col>
            <el-col :span="2" style="margin-left: 10px; margin-top: 10px;">
                <el-button color="#1e3864" @click="checksystem(1)">补丁检测</el-button>
            </el-col>
        </div>
    </div>
    <el-tabs v-model="activeName" type="border-card" style="margin-top: 10px;">
        <el-tab-pane label="杀软识别" name="1">
            <el-table :data="avlist" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="p_name" label="进程名称" />
                <el-table-column prop="pid" label="PID" />
                <el-table-column prop="a_name" label="杀软名称" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="补丁检测" name="2">
            <el-table :data="patch" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="msid" label="微软编号" />
                <el-table-column prop="kbid" label="补丁编号" />
                <el-table-column prop="des" label="描述" />
                <el-table-column prop="windows" label="影响系统" />
                <el-table-column prop="link" label="链接" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
    </el-tabs>
</template>