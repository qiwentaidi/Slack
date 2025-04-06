import { FormRules } from 'element-plus'
import { reactive, ref } from 'vue'
import { structs } from 'wailsjs/go/models'

var space = reactive({
    fofaapi: 'https://fofa.info/',
    fofaemail: '',
    fofakey: '',
    hunterkey: '',
    quakekey: '',
    chaos: '',
    bevigil: '',
    zoomeye: '',
    securitytrails: '',
    github: '',
})

var proxy = reactive({
    enabled: false,
    mode: 'HTTP',
    address: '127.0.0.1',
    port: 8080,
    username: '',
    password: '',
})

var webscan = reactive({
    web_thread: 50,
    crack_thread: 20, // 暴破任务线程
    port_thread: 1000,
    port_timeout: 7,
    ping_check_alive: false,
    default_alive_module: "None",
    default_network: "Auto",
    highlight_fingerprints: [
        "泛微-协同办公OA",
        "致远互联-OA",
        "易宝OA",
        "用友-NC-Cloud",
        "用友-GRP-U8",
        "用友-U8-CRM",
        "用友U8",
        "泛微-EOffice",
        "泛微-云桥e-Bridge",
        "泛微-EMobile",
        "泛微-E-message",
        "GitLab",
        "金蝶云星空",
        "K8s-Etcd",
        "Jeecg-Boot",
        "KindEditor",
        "致远OA-M3-Server",
        "PageOffice",
        "Shiro",
        "畅捷通-TPlus",
        "Dahua-ICC",
        "若依-管理系统",
        "Fastadmin",
        "亿赛通-电子文档安全管理系统",
        "TDXK-通达OA",
        "XXL-JOB-执行器",
        "RocketMQ",
        "ZABBIX-监控系统",
        "WebLogic",
        "MINIO-Console",
        "飞企互联-FE业务协作平台",
        "HIKVISION-综合安防管理平台",
        "契约锁-电子签署平台",
        "金万维",
        "XXL-JOB"
    ],
    append_pocfile: "",
})

var database = reactive({
    columnsNameKeywords: 'phone,telephone,idcard,id_number,id_card,password,username,mobile,sfz,secret,address,birth'
})

// 临时全局变量但是不进行保存
var temp = reactive({
    NetworkCardList: <structs.NetwordCard[]>[{ Name: "", IP: "Auto" }],
    nucleiEnabled: false,
    isMacOS: false,
    isMax: false,
    isGrid: true,
    goos: '',
})

const Logger = reactive({
    value: '',
    length: 100, // 日志显示条数
})

const LOCAL_VERSION = "2.0.7"

const Language = ref("zh")
const Theme = ref(false)

var PATH = {
    homedir: "",
    LocalPocVersionFile: "/slack/config/version",
    PortBurstPath: "/slack/portburte"
}

var UPDATE = reactive({
    PocStatus: false,
    ClientStatus: false,
    LocalPocVersion: "",
    RemotePocVersion: "",
    RemoteClientVersion: "",
    PocContent: "",
    ClientContent: "",
    updateDialog: false,
})


var jsfinder = reactive({
    whiteList: ["github.com", "google.com", "amazon.com", "gitee.com", "w3.org", "qq.com"],
    authFiled: ["token不能为空", "令牌不能为空", "令牌已过期", "Unauthorized", "Access Denied", "认证失败", "\"code\":\"401", "\"code\":401", "\"code\":\"403", "\"code\":403", "\"code\":\"404", "\"code\":404"],
    highRiskRouter: ["logout", "loginout", "insert", "update", "remove", "add", "change", "save", "import", "create", "enable", "del", "disable"],
})

var fileRetrieval = reactive({
    keywords: 'password,username,jdbc,accesskey,secret,jwt',
    blackList: '.exe,.dll,.so',
})

var syntaxRules = reactive<FormRules<structs.SpaceEngineSyntax>>({
    Name: [
        { required: true, message: '请输入语法名称', trigger: 'blur' },
    ],
    Content: [
        {
            required: true,
            message: '请输入语法内容',
            trigger: 'blur',
        },
    ],
})

export default {
    space,
    proxy,
    Logger,
    LOCAL_VERSION,
    PATH,
    UPDATE,
    jsfinder,
    webscan,
    temp,
    Language,
    Theme,
    database,
    syntaxRules,
    fileRetrieval
};
