<script lang="ts" setup>
import { Clock, View, Edit, Delete, UploadFilled, Share } from '@element-plus/icons-vue'
import { rp, historyDialog } from '@/composables/useWebscanState'
import { taskManager } from '@/composables/useTaskManager'
</script>

<template>
    <el-drawer v-model="historyDialog" size="70%">
        <template #header>
            <el-text class="font-bold" style="font-size: 16px;">
                <el-icon :size="18" class="mr-5px">
                    <Clock />
                </el-icon>
                <span>任务管理</span>
            </el-text>
        </template>
        <el-table :data="rp.table.pageContent" stripe @selection-change="rp.ctrl.handleSelectChange"
            :cell-style="{ textAlign: 'center' }" :header-cell-style="{ 'text-align': 'center' }"
            style="height: calc(100vh - 115px)">
            <el-table-column type="selection" width="50px" />
            <el-table-column prop="TaskName" label="任务名称" :show-overflow-tooltip="true"></el-table-column>
            <el-table-column prop="Targets" label="目标" :show-overflow-tooltip="true">
                <template #default="scope">
                    <el-link @click="taskManager.viewTask(scope.row)">
                        {{ scope.row.Targets.includes('\n') ?
                            scope.row.Targets.split('\n')[0] : scope.row.Targets }}
                    </el-link>
                </template>
            </el-table-column>
            <el-table-column label="资产" width="100px">
                <template #default="scope">
                    <span>{{ scope.row.Targets.includes('\n') ? scope.row.Targets.split('\n').length : 1 }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="Vulnerability" label="漏洞" width="100px" />
            <el-table-column label="操作" width="180px" align="center">
                <template #default="scope">
                    <el-button-group>
                        <el-tooltip content="查看">
                            <el-button :icon="View" link @click="taskManager.viewTask(scope.row)" />
                        </el-tooltip>
                        <el-tooltip content="重命名">
                            <el-button :icon="Edit" link @click="taskManager.renameTask(scope.row.TaskId)" />
                        </el-tooltip>
                        <el-tooltip content="删除">
                            <el-button :icon="Delete" link @click="taskManager.deleteTask([scope.row.TaskId])" />
                        </el-tooltip>
                    </el-button-group>
                </template>
            </el-table-column>
            <template #empty>
                <el-empty></el-empty>
            </template>
        </el-table>
        <div class="flex-between mt-5px">
            <el-space>
                <el-button :icon="UploadFilled" size="small" @click="taskManager.importTask()">导入任务</el-button>
                <el-button :icon="Share" size="small" @click="taskManager.showExportDialog"
                    :disabled="rp.table.selectRows.length < 1">导出报告</el-button>
                <el-button :icon="Delete" size="small"
                    @click="taskManager.deleteTask(rp.table.selectRows.map(item => item.TaskId))"
                    :disabled="rp.table.selectRows.length < 1">批量删除</el-button>
            </el-space>
            <el-pagination size="small" background @size-change="rp.ctrl.handleSizeChange"
                @current-change="rp.ctrl.handleCurrentChange" :pager-count="5" :current-page="rp.table.currentPage"
                :page-sizes="[20, 50, 100]" :page-size="rp.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="rp.table.result.length">
            </el-pagination>
        </div>
    </el-drawer>
</template>

