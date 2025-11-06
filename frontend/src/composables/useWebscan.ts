import { onMounted, nextTick, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { InitRule, FingerprintList, GetFingerPocMap } from 'wailsjs/go/services/App'
import { AddFingerscanResult, AddPocscanResult } from 'wailsjs/go/services/Database'
import { EventsOn, EventsOff } from 'wailsjs/runtime/runtime'
import { FilepathJoin, List } from 'wailsjs/go/services/File'
import { RetrieveAllScanTasks } from 'wailsjs/go/services/Database'
import { structs } from 'wailsjs/go/models'
import { typeToIcon } from '@/stores/style'
import { ActivityItem } from '@/stores/interface'
import throttle from 'lodash/throttle'
import global from "@/stores"
import { dashboard, param, fp, vp, rp, form, activities } from './useWebscanState'
import { taskManager } from './useTaskManager'

// 增加较快的数据需要节流刷新
export const throttleFingerscanUpdate = throttle(() => {
    fp.ctrl.watchResultChange(fp.table)
}, 1000);

// 热加载节流触发，防止反序列化过于频繁导致闪退
export const throttleInitialize = throttle(() => {
    initialize()
}, 2000);

// timeline 滚动容器
export const timelineContainer = ref<HTMLElement | null>(null)

export function addActivity(newActivity: ActivityItem) {
    activities.value.push({
        content: newActivity.content,
        timestamp: new Date().toLocaleTimeString(),
        type: newActivity.type,
        icon: typeToIcon[newActivity.type]
    })

    nextTick(() => {
        if (timelineContainer.value) {
            timelineContainer.value.scrollTop = timelineContainer.value.scrollHeight
        }
    })
}

// 初始化时调用
export async function initialize() {
    let isSuccess = await InitRule(global.webscan.append_pocfile);
    if (!isSuccess) {
        ElMessage.error({
            showClose: true,
            message: "初始化指纹规则失败，请检查配置文件",
        });
        return;
    }
    let FingperintDBList = await FingerprintList();
    dashboard.fingerLength = FingperintDBList?.length || 0;
    param.allFingerprint = Array.from(new Set(FingperintDBList)).map(item => ({ label: item, value: item }))

    // 获取POC数量
    let pocMap = await GetFingerPocMap();
    dashboard.pocLength = pocMap ? Object.keys(pocMap).length : 0;

    // 遍历模板
    let defaultPath = await FilepathJoin([global.PATH.homedir, "/slack/config/pocs"]);
    let files = await List([defaultPath, global.webscan.append_pocfile]);
    param.allTemplate = files
        .filter(file => file.Path.endsWith(".yaml"))
        .map(file => ({ label: file.BaseName, value: file.Path }));

    let result = await RetrieveAllScanTasks();
    if (Array.isArray(result)) {
        rp.table.result = result;
        rp.ctrl.watchResultChange(rp.table);
    }
}

export function setupEventListeners() {
    // 获得结果回调
    EventsOn("nucleiResult", (result: structs.VulnerabilityInfo) => {
        // 更新漏洞数量
        const riskLevelKey = result.Severity as keyof typeof dashboard.riskLevel;
        dashboard.riskLevel[riskLevelKey]++;
        AddPocscanResult(result)
        // 前段漏洞表格，需要当任务ID与结果的任务ID一致时，才更新漏洞表格
        if (form.taskId == result.TaskId) {
            vp.table.result.push(result)
            vp.ctrl.watchResultChange(vp.table)
        }
        taskManager.updateTaskTable(result.TaskId)
    });
    EventsOn("webFingerScan", (result: structs.InfoResult) => {
        if ((result.Scheme == "http" || result.Scheme == "https") && result.StatusCode == 0) {
            addActivity({
                content: result.URL + " 访问失败",
                type: "warning",
            })
        }
        if (result.StatusCode == 422 || result.StatusCode == 402) {
            addActivity({
                content: result.URL + " 为云防护地址, 已过滤",
                type: "warning",
            })
            return
        }
        AddFingerscanResult(result)
        if (form.taskId == result.TaskId) {
            fp.table.result.push(result)
            throttleFingerscanUpdate()
        }
        taskManager.updateTaskTable(result.TaskId)
    });
    EventsOn("ActiveCounts", (count: number) => {
        dashboard.activeCount = count
    });
    EventsOn("ActiveProgressID", (id: number) => {
        dashboard.activePercentage = Number(((id / dashboard.activeCount) * 100).toFixed(2));
    });
    EventsOn("NucleiCounts", (count: number) => {
        dashboard.nucleiCount = count
    });
    EventsOn("NucleiProgressID", (id: number) => {
        dashboard.nucleiPercentage = Number(((id / dashboard.nucleiCount) * 100).toFixed(2));
    });
    // 进度条
    EventsOn("progressID", (id: number) => {
        dashboard.portscanPercentage = Number(((id / dashboard.portscanCount) * 100).toFixed(2));
    });
    // 端口扫描总数（由后端发送，确保前后端一致）
    EventsOn("portScanTotalCount", (count: number) => {
        dashboard.portscanCount = count;
    });
}

export function cleanupEventListeners() {
    EventsOff("nucleiResult");
    EventsOff("webFingerScan");
    EventsOff("ActiveCounts");
    EventsOff("ActiveProgressID");
    EventsOff("NucleiCounts");
    EventsOff("NucleiProgressID");
    EventsOff("portScanLoading");
    EventsOff("progressID");
    EventsOff("portScanTotalCount");
}

