<script lang="ts" setup>
import { Search, Monitor, Warning, Refresh } from '@element-plus/icons-vue'
import Loading from '@/components/Loading.vue'
import { getBadgeClass } from '@/stores/style'
import { dashboard } from '@/composables/useWebscanState'

defineEmits<{
    refresh: []
}>()
</script>

<template>
    <el-card class="mb-10px">
        <div class="flex-between">
            <span class="font-bold ml-5px" style="font-size: 20px;">
                Overview
            </span>
            <div class="risk-level-display">
                <el-tooltip v-for="(value, level) in dashboard.riskLevel" :key="level" :content="level">
                    <div class="risk-badge" :class="getBadgeClass(level)">
                        {{ value }}
                    </div>
                </el-tooltip>
            </div>
            <slot name="actions"></slot>
        </div>
        <el-row :gutter="10" class="mt-10px">
            <el-col :span="12">
                <div class="flex">
                    <Loading style="margin-left: 40px;" />
                    <div class="progress-details">
                        <div class="progress-item">
                            <el-icon class="icon icon-blue">
                                <Search />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>端口扫描</span>
                                    <span>{{ dashboard.portscanPercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.portscanPercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>

                        <div class="progress-item">
                            <el-icon class="icon icon-green">
                                <Monitor />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>主动探测</span>
                                    <span>{{ dashboard.activePercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.activePercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>

                        <div class="progress-item">
                            <el-icon class="icon icon-red">
                                <Warning />
                            </el-icon>
                            <div class="progress-content">
                                <div class="progress-labels">
                                    <span>漏洞检测</span>
                                    <span>{{ dashboard.nucleiPercentage }}%</span>
                                </div>
                                <el-progress :percentage="dashboard.nucleiPercentage" :show-text="false"
                                    :stroke-width="8" />
                            </div>
                        </div>
                    </div>
                </div>
                <div class="summary-box">
                    <span class="summary-value">{{ dashboard.fingerLength }}</span>
                    <span class="summary-label">指纹总数</span><br />
                    <span class="summary-value">{{ dashboard.pocLength }}</span>
                    <span class="summary-label">漏洞总数</span>
                    <el-tooltip content="热加载指纹和POC">
                        <el-button link @click="$emit('refresh')">
                            <template #icon>
                                <el-icon :size="20">
                                    <Refresh />
                                </el-icon>
                            </template>
                        </el-button>
                    </el-tooltip>
                </div>
            </el-col>
            <el-col :span="12">
                <slot name="timeline"></slot>
            </el-col>
        </el-row>
    </el-card>
</template>

<style scoped>
.summary-box {
    justify-content: center;
    margin-top: 10px;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: var(--timeline-bg-color);
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
}

.summary-value {
    font-size: 1.125rem;
    font-weight: 500;
}

.summary-label {
    font-size: 0.875rem;
    color: #6B7280;
}

.progress-details {
    display: flex;
    flex-direction: column;
    gap: 24px;
    margin-left: 48px;
    flex: 1;
}

.progress-item {
    display: flex;
    align-items: center;
    gap: 16px;
}

.icon {
    font-size: 18px;
}

.icon-blue {
    color: #3B82F6;
}

.icon-green {
    color: #22C55E;
}

.icon-red {
    color: #EF4444;
}

.progress-content {
    flex: 1;
}

.progress-labels {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
}
</style>

