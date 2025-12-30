<script setup lang="ts">
import { onMounted, ref, onUnmounted, computed, watch, nextTick } from "vue";
import { useI18n } from "vue-i18n";
import * as echarts from "echarts";
import { GetVulnDashboard, GetTopRiskTargets } from "wailsjs/go/services/Database";
import type { structs } from "wailsjs/go/models";
import { EventsOn } from "wailsjs/runtime/runtime";
import global from "@/stores";
import { useRouter } from "vue-router";

const { t } = useI18n();
const dashboard = ref<structs.VulnDashboard | null>(null);
const riskTargets = ref<structs.TaskResult[]>([]);
const chartRef = ref<HTMLElement | null>(null);
const cpuUsage = ref<string>("-");
const memUsage = ref<string>("-");
const platform = ref<string>("未知");
const arch = ref<string>("未知");
const hostInfo = computed(() => ([
    { label: "CPU 使用率", value: cpuUsage.value },
    { label: "内存使用率", value: memUsage.value },
    { label: "平台", value: platform.value },
    { label: "架构", value: arch.value }
]));
let myChart: echarts.ECharts | null = null;
let statsOff: (() => void) | null = null;
const maxVuln = computed(() => {
    const vals = riskTargets.value.slice(0, 5).map((i) => i?.Vulnerability || 0);
    const max = Math.max(...vals, 0);
    return max > 0 ? max : 1;
});

const tt = (key: string, fallback: string) => {
    const res = t(key);
    return res === key ? fallback : res;
};

const readCssVar = (name: string, fallback: string) => {
    if (typeof window === "undefined") return fallback;
    const value = getComputedStyle(document.documentElement).getPropertyValue(name);
    return value?.trim() || fallback;
};

const initChart = () => {
    if (!chartRef.value) return;
    if (myChart) myChart.dispose();
    myChart = echarts.init(chartRef.value);
    const critical = readCssVar("--risk-critical", "#f74141");
    const high = readCssVar("--risk-high", "#e6a23c");
    const medium = readCssVar("--risk-medium", "#73a8ed");
    const low = readCssVar("--risk-low", "#909399");
    const info = readCssVar("--risk-info", "#8fa3c4");
    const panelBg = readCssVar("--el-bg-color", "#fff");
    const panelBorder = readCssVar("--el-border-color-lighter", "#ebeef5");
    const textPrimary = readCssVar("--el-text-color-primary", "#303133");
    const option = {
        tooltip: {
            trigger: "item",
            backgroundColor: panelBg,
            borderRadius: 8,
            textStyle: { color: textPrimary },
            borderColor: panelBorder
        },
        legend: {
            orient: "vertical",
            right: 6,
            top: "center",
            itemWidth: 10,
            itemHeight: 10,
            icon: "circle",
            textStyle: { color: textPrimary, fontSize: 13, fontWeight: 600 }
        },
        series: [{
            type: "pie",
            radius: ["64%", "84%"],
            center: ["50%", "50%"],
            avoidLabelOverlap: false,
            itemStyle: { borderRadius: 8, borderColor: panelBorder, borderWidth: 3, shadowBlur: 8, shadowColor: "rgba(0,0,0,0.03)" },
            label: { show: false },
            data: [
                { value: dashboard.value?.SeverityCritical || 0, name: "严重", itemStyle: { color: critical } },
                { value: dashboard.value?.SeverityHigh || 0, name: "高危", itemStyle: { color: high } },
                { value: dashboard.value?.SeverityMedium || 0, name: "中危", itemStyle: { color: medium } },
                { value: dashboard.value?.SeverityLow || 0, name: "低危", itemStyle: { color: low } },
                { value: dashboard.value?.SeverityInfo || 0, name: "信息", itemStyle: { color: info } }
            ]
        }]
    };
    myChart.setOption(option);
};

const loadData = async () => {
    const [dbRes, riskRes] = await Promise.all([GetVulnDashboard(), GetTopRiskTargets()]);
    dashboard.value = dbRes;
    riskTargets.value = Array.isArray(riskRes) ? riskRes : [];
    initChart();
};

const collectHostInfo = () => {
    const nav = typeof navigator !== "undefined" ? navigator : undefined;
    const ua = nav?.userAgent || "";
    const lowerUa = ua.toLowerCase();
    if (lowerUa.includes("mac")) {
        platform.value = "macOS";
    } else if (lowerUa.includes("win")) {
        platform.value = "Windows";
    } else if (lowerUa.includes("linux")) {
        platform.value = "Linux";
    } else {
        platform.value = "未知";
    }
    const uaData = (nav as any)?.userAgentData;
    const hintArch = uaData?.architecture || "";
    if (hintArch) {
        arch.value = hintArch;
    } else if (lowerUa.includes("arm") || lowerUa.includes("aarch64")) {
        arch.value = "ARM64";
    } else if (lowerUa.includes("x86_64") || lowerUa.includes("amd64") || lowerUa.includes("intel")) {
        arch.value = "x86_64";
    } else {
        arch.value = "未知";
    }
};

const handleResize = () => myChart?.resize();

const formatTargetDisplay = (raw: string | undefined) => {
    if (!raw) return { main: "", extra: 0 };
    const parts = raw.split(/[,;\s]+/).filter(Boolean);
    const main = parts[0] || raw;
    const extra = parts.length > 1 ? parts.length - 1 : 0;
    return { main, extra };
};

const router = useRouter();
const goTaskDetail = (task: structs.TaskResult) => {
    if (!task?.TaskId) return;
    router.push({ path: "/Permeation/Webscan", query: { taskId: task.TaskId } });
};

const parseUsage = (val: any) => {
    if (typeof val === "number") return `${val.toFixed(1)}%`;
    if (typeof val === "string") return val;
    return "-";
};

onMounted(() => {
    loadData();
    collectHostInfo();
    statsOff = EventsOn("system-stats", (payload: any) => {
        if (payload?.cpu !== undefined) cpuUsage.value = parseUsage(payload.cpu);
        if (payload?.mem !== undefined) memUsage.value = parseUsage(payload.mem);
        else if (payload?.memory !== undefined) memUsage.value = parseUsage(payload.memory);
    });
    window.addEventListener("resize", handleResize);
});

watch(() => global.Theme.value, () => {
    nextTick(() => initChart());
});

onUnmounted(() => {
    window.removeEventListener("resize", handleResize);
    statsOff?.();
    if (myChart) {
        myChart.dispose();
        myChart = null;
    }
});
</script>

<template>
    <div class="page-shell">

        <div class="content-grid">
            <section class="panel-card ranking-card compact">
                <div class="panel-head aligned">
                    <h3>风险任务排行</h3>
                    <el-tag round effect="plain" type="info">TOP5</el-tag>
                </div>
                <div class="ranking-list">
                    <div v-for="(item, idx) in riskTargets.slice(0, 5)" :key="idx" class="ranking-item">
                        <span class="r-idx" :class="'idx-' + (idx + 1)">0{{ idx + 1 }}</span>
                        <div class="r-info">
                            <div class="r-name-row">
                                <span class="r-name">{{ item.TaskName || item.Targets }}</span>
                                <span class="r-score">{{ item.Vulnerability }}</span>
                            </div>
                            <el-progress
                                :percentage="Math.min(((item.Vulnerability || 0) / maxVuln) * 100, 100)"
                                :stroke-width="6"
                                :show-text="false"
                                :color="idx === 0 ? 'var(--risk-critical)' : 'var(--risk-medium)'"
                            />
                        </div>
                    </div>
                </div>
            </section>

            <section class="panel-card chart-card small">
                <div class="panel-head aligned">
                    <h3>风险等级分布</h3>
                    <el-tag round effect="plain" type="info">实时更新</el-tag>
                </div>
                <div class="chart-wrapper">
                    <div ref="chartRef" class="chart-instance"></div>
                    <div class="chart-center-text">
                        <span class="text-total">{{ dashboard?.TotalVulnerabilities || 0 }}</span>
                        <span class="text-label">总漏洞</span>
                    </div>
                </div>
            </section>
        </div>

        <section class="panel-card activity-panel">
            <div class="panel-head">
                <h3>近期任务</h3>
                <el-tag round effect="plain" type="info">最新</el-tag>
            </div>
            <div v-if="dashboard?.RecentTasks?.length" class="activity-grid">
                <div v-for="task in dashboard?.RecentTasks?.slice(0, 3)" :key="task.TaskId" class="activity-card">
                    <div class="a-header">
                        <span class="a-title">{{ task.TaskName }}</span>
                        <span class="a-status">活跃</span>
                    </div>
                    <div class="a-body">
                        <code>
                            <span class="target-main">{{ formatTargetDisplay(task.Targets).main }}</span>
                            <span v-if="formatTargetDisplay(task.Targets).extra" class="target-extra">
                                +{{ formatTargetDisplay(task.Targets).extra }}
                            </span>
                        </code>
                    </div>
                    <div class="a-footer">
                        <span class="a-vuln-count">发现 <b>{{ task.Vulnerability }}</b> 个漏洞</span>
                        <el-button size="small" link type="primary" @click="goTaskDetail(task)">查看详情</el-button>
                    </div>
                </div>
            </div>
            <div v-else class="empty-state">暂无最新活动</div>
        </section>

        <section class="panel-card host-panel">
            <div class="panel-head">
                <h3>本机信息</h3>
                <el-tag round effect="plain" type="success">本地</el-tag>
            </div>
            <div class="host-grid">
                <div v-for="item in hostInfo" :key="item.label" class="host-item">
                    <span class="host-label">{{ item.label }}</span>
                    <span class="host-value">{{ item.value }}</span>
                </div>
            </div>
        </section>
    </div>
</template>

<style scoped>
:global(:root) {
    --risk-critical: #b30000;
    --risk-high: #ff4d4f;
    --risk-medium: #f5c542;
    --risk-low: #409eff;
    --risk-info: #8fa3c4;
}

/* trimmed hero removed for compact layout */
.page-shell {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 6px 0 12px;
}

.kpi-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 14px;
}

.kpi-card {
    padding: 16px;
    border-radius: 16px;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.05);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.kpi-card:hover { transform: translateY(-4px); box-shadow: 0 12px 28px rgba(0, 0, 0, 0.08); }
.kpi-label { font-size: 13px; color: var(--el-text-color-secondary); }
.kpi-num { font-size: 30px; font-weight: 800; margin: 6px 0; color: var(--el-text-color-primary); }
.kpi-card.small .kpi-num { font-size: 22px; }
.kpi-card small { color: var(--el-text-color-secondary); }

.content-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 12px;
}

.panel-card {
    background: var(--el-bg-color);
    border-radius: 20px;
    padding: 18px;
    border: 1px solid var(--el-border-color-lighter);
    box-shadow: 0 10px 28px rgba(0, 0, 0, 0.06);
}

.panel-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.panel-head h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 800;
    color: var(--el-text-color-primary);
}

.panel-head.aligned > * {
    margin: 0;
    line-height: 1.2;
    display: flex;
    align-items: center;
}

.chart-card { position: relative; overflow: hidden; }
.chart-card::after {
    content: "";
    position: absolute;
    bottom: -60px;
    right: -60px;
    width: 180px;
    height: 180px;
    background: radial-gradient(circle, color-mix(in srgb, var(--el-color-primary) 20%, transparent) 0%, transparent 70%);
}

.chart-wrapper { position: relative; height: 260px; }
.chart-instance { height: 100%; width: 100%; }
.chart-center-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-align: center;
    pointer-events: none;
}
.text-total { display: block; font-size: 30px; font-weight: 900; color: var(--el-text-color-primary); }
.text-label { font-size: 12px; color: var(--el-text-color-regular); font-weight: 600; }

.host-panel { padding-top: 16px; }
.host-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 10px;
}
.host-item {
    padding: 12px;
    border: 1px dashed var(--el-border-color-lighter);
    border-radius: 12px;
    background: var(--el-fill-color-light);
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.host-label { font-size: 12px; color: var(--el-text-color-secondary); }
.host-value { font-size: 14px; font-weight: 700; color: var(--el-text-color-primary); word-break: break-all; }

.ranking-list { display: flex; flex-direction: column; gap: 14px; }
.ranking-item {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 12px 14px;
    border-radius: 12px;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
}
.r-idx { font-weight: 800; font-size: 13px; color: var(--el-text-color-placeholder); width: 36px; text-align: center; }
.idx-1 { color: var(--risk-critical); }
.idx-2 { color: var(--risk-high); }
.idx-3 { color: var(--risk-medium); }
.r-info { flex: 1; }
.r-name-row { display: flex; justify-content: space-between; margin-bottom: 6px; align-items: center; }
.r-name { font-size: 13px; font-weight: 700; color: var(--el-text-color-primary); }
.r-score { font-size: 13px; font-weight: 800; color: var(--el-color-primary); }

.activity-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 14px;
}

.ranking-card.compact .panel-head { margin-bottom: 6px; }
.ranking-card.compact .ranking-item { padding: 10px 12px; }
.chart-card.small .panel-head { margin-bottom: 6px; }

.activity-card {
    padding: 14px;
    background: var(--el-bg-color);
    border-radius: 16px;
    border: 1px solid var(--el-border-color-lighter);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.activity-card:hover { transform: translateY(-3px); box-shadow: 0 12px 22px rgba(0, 0, 0, 0.08); }
.a-header { display: flex; justify-content: space-between; align-items: center; }
.a-title { font-size: 14px; font-weight: 700; color: var(--el-text-color-primary); }
.a-status { font-size: 12px; color: #52c41a; background: rgba(82, 196, 26, 0.12); padding: 4px 8px; border-radius: 999px; }
.a-body { margin: 10px 0; }
.a-body code {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    background: rgba(0, 0, 0, 0.04);
    padding: 6px 8px;
    border-radius: 8px;
    display: flex;
    gap: 6px;
    align-items: center;
    white-space: nowrap;
}
.target-main {
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
    display: inline-block;
    max-width: 100%;
}
.target-extra { color: var(--el-text-color-placeholder); }
.a-footer { font-size: 12px; color: var(--el-text-color-regular); display: flex; align-items: center; gap: 8px; }
.a-footer b { color: var(--el-color-danger); }

.empty-state {
    padding: 30px;
    text-align: center;
    color: var(--el-text-color-secondary);
    background: var(--el-fill-color-light);
    border-radius: 14px;
    border: 1px dashed var(--el-border-color-lighter);
}

@media (max-width: 1200px) {
    .page-hero { grid-template-columns: 1fr; }
    .content-grid { grid-template-columns: 1fr; }
}

@media (max-width: 820px) {
    .hero-actions { flex-direction: column; align-items: flex-start; }
    .hero-metrics { grid-template-columns: 1fr; }
    .kpi-grid { grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); }
}
</style>
