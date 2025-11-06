import { reactive, ref } from 'vue'
import { portGroupOptions } from '@/stores/options'
import { structs } from 'wailsjs/go/models'
import { ActivityItem } from '@/stores/interface'
import usePagination from '@/usePagination'

// 仪表盘状态
export const dashboard = reactive({
    // 漏洞数量
    riskLevel: {
        CRITICAL: 0,
        HIGH: 0,
        MEDIUM: 0,
        LOW: 0,
        INFO: 0,
    },
    // 指纹数量
    fingerLength: 0,
    // 漏洞数量
    pocLength: 0,
    // 主动指纹需要发包总数
    activeCount: 0,
    // 主动指纹当前发包进度
    activePercentage: 0,
    // 漏洞扫描进度
    nucleiPercentage: 0,
    // 漏洞扫描目标总数
    nucleiCount: 0,
    // IP*端口扫描总数
    portscanCount: 0,
    // 端口扫描当前进度
    portscanPercentage: 0,
})

// 界面参数，仅作用于前端显示或者输入匹配，与后端无交互
export const param = reactive({
    inputType: 0, // 输入类型
    portGroup: portGroupOptions[1].text,
    builtInUsername: true,
    builtInPassword: true,
    username: '',
    password: '',
    allTemplate: <{ label: string, value: string }[]>[],
    allFingerprint: <{ label: string, value: string }[]>[],
    hostFilter: '',
})

// webscan 表单状态
export const form = reactive({
    input: '',
    tags: [] as string[],
    newWebscanDrawer: false,
    newHostscanDrawer: false,
    runnningStatus: false,
    scanStopped: false,
    taskName: '',
    taskId: '',
    hideRequest: false,
    hideResponse: false,
    showYamlPoc: false,
    pocContent: '',
    portlist: '',
})

// 空间引擎配置
export const spaceEngineConfig = reactive({
    fofaQuery: '',
    fofaDialog: false,
    fofaPageSize: 1000,
    hunterQuery: '',
    hunterDialog: false,
    hunterPageSize: "100",
})

// 扫描配置
export const config = reactive({
    customTemplate: <string[]>[],
    customTags: <string[]>[],
    screenhost: false,
    rootPathScan: true,
    webscanOption: 0,
    skipNucleiWithoutTags: false,
    generateLog4j2: false,
    crack: false, // 是否开启暴破
    customHeaders: '',
    vulscan: false,
    excludePrintPorts: false, // 排除打印机端口
})

// 对话框状态
export const detailDialog = ref(false)
export const historyDialog = ref(false)
export const exportDialog = ref(false)
export const selectedRow = ref<any>()

// Shodan 相关状态
export const shodanVisible = ref(false)
export const shodanIp = ref('')
export const shodanPercentage = ref(0)
export const shodanThread = ref(2)
export const shodanRunningstatus = ref(false)

// 报告导出状态
export const reportOption = ref('HTML')
export const reportName = ref('')

// 分页实例
export const fp = usePagination<structs.InfoResult>(50)
export const vp = usePagination<structs.VulnerabilityInfo>(50)
export const rp = usePagination<structs.TaskResult>(20)

// 活动列表
export const activities = ref<ActivityItem[]>([])

export function updatePorts(index: number) {
    if (index >= 0 && index < portGroupOptions.length) {
        form.portlist = portGroupOptions[index].value;
    }
}

