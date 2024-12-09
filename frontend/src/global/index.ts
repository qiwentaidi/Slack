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
    port_thread: 300,
    port_timeout: 7,
    ping_check_alive: false,
    default_alive_module: "None",
    default_network: "Auto",
    highlight_fingerprints: ["畅捷通-TPlus", "泛微-协同办公OA", "致远互联-OA", "易宝OA", "用友-NC-Cloud", "用友-GRP-U8", "用友-U8-CRM", "用友U8", "泛微-EOffice", "泛微-云桥e-Bridge", "泛微-EMobile", "泛微-E-message", "GitLab", "金蝶云星空", "K8s-Etcd", "Jeecg-Boot", "KindEditor"],
    append_pocfile: "",
})

var database = reactive({
    columnsNameKeywords: 'phone,telephone,idcard,id_number,id_card,password,username,mobile,sfz,secret,address,birth'
})

// 临时全局变量但是不进行保存
var temp = reactive({
    dirsearchPathConut: 0,
    dirsearchConut: 0,
    dirsearchStartTime: 0,
    thinkdict: [] as string[],
    NetworkCardList: ["Auto"],
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

const LOCAL_VERSION = "1.7.4"

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


var jsfind = reactive({
    whiteList: "github.com\ngoogle.com\namazon.com\ngitee.com\nw3.org\nqq.com",
    defaultType: ["info", "warning", "danger", "primary", "success", "info"]
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
    jsfind,
    webscan,
    temp,
    Language,
    Theme,
    database,
    syntaxRules,
};
