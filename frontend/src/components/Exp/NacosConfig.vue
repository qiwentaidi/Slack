<script lang="ts" setup>
import { Document, DataAnalysis, Reading, InfoFilled } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { reactive } from 'vue';
import { structs } from 'wailsjs/go/models';
import { NacosCategoriesExtract } from 'wailsjs/go/services/App';
import { CheckFileStat, DirectoryDialog, ReadFile } from 'wailsjs/go/services/File';

const config = reactive({
    filePath: '',
    result: [] as structs.NacosConfig[],
    detailDialog: false,
    content: '',
    tipsDialog: false,
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
    config.detailDialog = true
    let isStat = await CheckFileStat(filepath)
    if (!isStat) {
        ElMessage.warning('文件不存在!')
        return
    }
    let file = await ReadFile(filepath)
    config.content = file.Content
}

const code = `var categories = map[string][]string{
    "Auth":     {"username", "password"},
    "OSS":      {"accesskey", "secret"},
    "Database": {"jdbc", "redis", "elasticsearch", "database", "mongo", "mssql", "mysql", "oracle", "postgres", "sqlserver"},
}`
</script>


<template>
    <div class="head" style="margin-bottom: 10px;">
        <el-button :icon="InfoFilled" @click="config.tipsDialog = true">模块介绍</el-button>
        <el-input v-model="config.filePath" placeholder="请输入解压后的文件夹路径" style="margin-inline: 5px;">
            <template #suffix>
                <el-button :icon="Document" link @click="selectDirectory"></el-button>
            </template>
        </el-input>
        <el-button type="primary" plain :icon="DataAnalysis" @click="startAnalysis">分析</el-button>
    </div>
    <el-table :data="config.result" :cell-style="{ textAlign: 'center' }"
        :header-cell-style="{ 'text-align': 'center' }" style="height: calc(100vh - 180px);">

        <!-- 序号列 -->
        <el-table-column type="index" label="#" width="60px" />

        <!-- 文件名称列 -->
        <el-table-column prop="Name" label="Filename">
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
        <el-table-column label="Operate" width="120px">
            <template #default="scope">
                <el-button size="small" :icon="Reading" @click="viewDetails(scope.row.Name)">
                    查看详情
                </el-button>
            </template>
        </el-table-column>
    </el-table>
    <el-drawer v-model="config.detailDialog" title="查看详情" size="70%">
        <highlightjs language="yaml" :code="config.content"></highlightjs>
    </el-drawer>
    <el-dialog v-model="config.tipsDialog" title="模块介绍" width="50%">
        <span>当Nacos处于<strong>未授权或者已鉴权登录状态时</strong>，可访问下面路径下载所有配置的zip文件<br/>
            /nacos/v1/cs/configs?export=true&group=&tenant=&appName=&ids=&dataId=<br/>
            应用会统计每个文件中关键次的出现次数
        </span>
        <highlightjs :code="code"></highlightjs>
    </el-dialog>
</template>