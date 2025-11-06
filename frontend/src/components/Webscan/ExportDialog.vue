<script lang="ts" setup>
import { getSelectedIcon } from '@/stores/style'
import { webReportOptions } from '@/stores/options'
import { exportDialog, reportOption, reportName, rp } from '@/composables/useWebscanState'
import { taskManager } from '@/composables/useTaskManager'
</script>

<template>
    <el-dialog title="导出报告" v-model="exportDialog">
        <el-alert :title="'已选择' + rp.table.selectRows.length + '个任务'" type="info" show-icon :closable="false"
            style="margin-bottom: 5px;" />
        <el-form label-width="auto">
            <el-form-item label="报告类型">
                <el-select v-model="reportOption" style="width: 240px;">
                    <template #label>
                        <el-space :size="6">
                            <el-icon :size="18">
                                <component :is="getSelectedIcon(reportOption)" />
                            </el-icon>
                            <span class="font-bold">{{ reportOption }}</span>
                        </el-space>
                    </template>
                    <el-option v-for="item in webReportOptions" :key="item.label" :label="item.label"
                        :value="item.label">
                        <el-space :size="6">
                            <el-icon :size="18">
                                <component :is="item.icon" />
                            </el-icon>
                            <span>{{ item.label }}</span>
                        </el-space>
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="报告名称">
                <el-input v-model="reportName"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="exportDialog = false">取消</el-button>
                <el-button type="primary" @click="taskManager.exportTask()">
                    导出
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

