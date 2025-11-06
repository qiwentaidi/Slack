<script lang="ts" setup>
import { onMounted, onUnmounted } from 'vue'
import { VideoPause, Plus, Clock } from '@element-plus/icons-vue'
import { ExitScanner } from 'wailsjs/go/services/App'
import { ElMessage } from 'element-plus';
import CustomTabs from '@/components/CustomTabs.vue';
import hostIcon from "@/assets/icon/host.svg";
import websiteIcon from "@/assets/icon/website.svg";
import Dashboard from '@/components/Webscan/Dashboard.vue'
import FingerprintTable from '@/components/Webscan/FingerprintTable.vue'
import VulnerabilityTable from '@/components/Webscan/VulnerabilityTable.vue'
import NewHostScanDrawer from '@/components/Webscan/NewHostScanDrawer.vue'
import NewWebScanDrawer from '@/components/Webscan/NewWebScanDrawer.vue'
import VulnerabilityDetailDrawer from '@/components/Webscan/VulnerabilityDetailDrawer.vue'
import TaskManagementDrawer from '@/components/Webscan/TaskManagementDrawer.vue'
import ExportDialog from '@/components/Webscan/ExportDialog.vue'
import SpaceEngineDialogs from '@/components/Webscan/SpaceEngineDialogs.vue'
import { WebscanEngine } from '@/composables/useWebscanEngine'
import { taskManager } from '@/composables/useTaskManager'
import { 
    form, param, config, activities, updatePorts, historyDialog
} from '@/composables/useWebscanState'
import { 
    initialize, setupEventListeners, cleanupEventListeners, throttleInitialize, timelineContainer
} from '@/composables/useWebscan'

async function startScan() {
    let engine = new WebscanEngine
    if (!await engine.checkOptions()) return
    param.inputType == 0 ? form.newWebscanDrawer = false : form.newHostscanDrawer = false
    if (param.inputType == 1 && (!await engine.checkHostscanOptions())) {
        return
    }
    // 检查是否能写入任务
    if (!(await taskManager.writeTask())) {
        form.runnningStatus = false
        return
    }
    engine.clearDashboard()
    form.runnningStatus = true
    await engine.Runner()
}

function stopScan() {
    if (!form.runnningStatus) return
    ElMessage.error("正在停止任务, 请稍后!")
    ExitScanner("[webscan]")
    ExitScanner("[portscan]")
    // 新增一个标志变量来确保setTimeout只执行一次
    if (!form.scanStopped) {
        form.scanStopped = true;
        setTimeout(() => {
            form.runnningStatus = false
            form.scanStopped = false;
            activities.value.push({
                content: "用户已退出扫描任务",
                timestamp: new Date().toLocaleTimeString(),
                type: "danger",
                icon: null
            })
        }, 10000);
    }
}

onMounted(() => {
    updatePorts(1);
    initialize()
    setupEventListeners()
});

onUnmounted(() => {
    cleanupEventListeners()
});
</script>

<template>
    <Dashboard @refresh="throttleInitialize()">
        <template #actions>
            <div>
                <el-popover placement="left-start" trigger="click" :width="244" v-if="!form.runnningStatus">
                    <template #reference>
                        <el-button type="primary" :icon="Plus">新建任务</el-button>
                    </template>
                    <el-space>
                        <el-button :icon="websiteIcon"
                            @click="form.newWebscanDrawer = true; param.inputType = 0">网站扫描</el-button>
                        <el-button :icon="hostIcon"
                            @click="form.newHostscanDrawer = true; param.inputType = 1">主机扫描</el-button>
                    </el-space>
                </el-popover>
                <el-button type="danger" :icon="VideoPause" @click="stopScan" v-else>停止任务</el-button>
            </div>
        </template>
        <template #timeline>
            <div ref="timelineContainer" class="timelineContainer">
                <el-timeline v-if="activities.length >= 1" style="text-align: left; padding-left: 5px;">
                    <el-timeline-item v-for="(activity, index) in activities" :key="index" :icon="activity.icon"
                        :type="activity.type" :timestamp="activity.timestamp">
                        {{ activity.content }}
                    </el-timeline-item>
                </el-timeline>
                <span class="position-center" v-else>No active tasks</span>
            </div>
        </template>
    </Dashboard>
    <CustomTabs>
        <el-tabs type="border-card">
            <el-tab-pane label="信息">
                <FingerprintTable />
            </el-tab-pane>
            <el-tab-pane label="漏洞">
                <VulnerabilityTable />
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-button :icon="Clock" @click="historyDialog = true">任务管理</el-button>
        </template>
    </CustomTabs>
    <NewHostScanDrawer @start-scan="startScan" />
    <NewWebScanDrawer @start-scan="startScan" />
    <VulnerabilityDetailDrawer />
    <TaskManagementDrawer />
    <ExportDialog />
    <SpaceEngineDialogs />
</template>

<style scoped>
.timelineContainer {
    background-color: var(--timeline-bg-color);
    padding: 1.5rem;
    border-radius: 0.5rem;
    height: 172px;
    overflow-y: auto;
}
</style>
