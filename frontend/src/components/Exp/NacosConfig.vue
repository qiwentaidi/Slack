<script lang="ts" setup>
import { Document, DataAnalysis, Reading } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { reactive } from 'vue';
import { structs } from 'wailsjs/go/models';
import { NacosCategoriesExtract } from 'wailsjs/go/services/App';
import { CheckFileStat, DirectoryDialog, ReadFile } from 'wailsjs/go/services/File';

const config = reactive({
    filePath: '',
    result: [] as structs.NacosConfig[],
    dialog: false,
    content: '',
})

async function selectDirectory() {
    let dir = await DirectoryDialog()
    if (!dir) {
        return
    }
    config.filePath = dir
}

async function startAnalysis() {
    if (!config.filePath) {
        ElMessage.warning('请选择文件夹路径')
        return
    }
    config.result = await NacosCategoriesExtract(config.filePath)
}

async function viewDetails(filepath: string) {
    config.dialog = true
    let isStat = await CheckFileStat(filepath)
    if (!isStat) {
        ElMessage.warning('文件不存在!')
        return
    }
    let file = await ReadFile(filepath)
    config.content = file.Content
}
</script>


<template>
    <div class="head" style="margin-bottom: 10px;">
        <el-input v-model="config.filePath" placeholder="请输入解压后的文件夹路径">
            <template #prepend>
                <span>文件夹路径</span>
            </template>
            <template #suffix>
                <el-button :icon="Document" link @click="selectDirectory"></el-button>
            </template>
        </el-input>
        <el-button type="primary" plain :icon="DataAnalysis" @click="startAnalysis"
            style="margin-left: 5px;">分析</el-button>
    </div>
    <el-table :data="config.result" :cell-style="{ textAlign: 'center' }"
        :header-cell-style="{ 'text-align': 'center' }" style="height: calc(100vh - 180px);">

        <!-- 序号列 -->
        <el-table-column type="index" label="#" width="60px" />

        <!-- 文件名称列 -->
        <el-table-column prop="Name" label="名称">
        </el-table-column>

        <!-- Auth 列 -->
        <el-table-column label="Auth (账号密码)" width="200px">
            <template #default="scope">
                {{ scope.row.NodeInfo.Auth }}
            </template>
        </el-table-column>

        <!-- OSS 列 -->
        <el-table-column label="OSS" width="200px">
            <template #default="scope">
                {{ scope.row.NodeInfo.OSS }}
            </template>
        </el-table-column>

        <!-- Database 列 -->
        <el-table-column label="Database" width="200px">
            <template #default="scope">
                {{ scope.row.NodeInfo.Database }}
            </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" width="120px">
            <template #default="scope">
                <el-button size="small" icon="Reading" @click="viewDetails(scope.row.Name)">
                    查看详情
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-drawer v-model="config.dialog" title="查看详情" size="70%">
        <highlightjs language="yaml" :code="config.content"></highlightjs>
    </el-drawer>
</template>