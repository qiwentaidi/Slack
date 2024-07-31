import { reactive, ref } from 'vue'
import { URLFingerMap } from '../interface'

var space = reactive({
    fofaapi: 'https://fofa.info/',
    fofaemail: '',
    fofakey: '',
    hunterkey: '',
    quakekey: '',
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
    nucleiEngine: "",
})

// 临时全局变量但是不进行保存
var temp = reactive({
    urlFingerMap: [] as URLFingerMap[],
    dirsearchPathConut: 0,
    dirsearchStartTime: 0,
    thinkdict: [] as string[],
    localAddItem: false, 
})

const Logger = reactive({
    value: '',
    length: 100, // 日志显示条数
})

const LOCAL_VERSION = "1.5.8"

const Language = ref("zh")
const Theme = ref(false)

var PATH = {
    ConfigPath: "/slack/config",
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
})

const portGroup = [
    {
        text: "数据库端口",
        value: "1433,1521,3306,5432,6379,9200,11211,27017",
    },
    {
        text: "企业端口",
        value: "21,22,80,81,135,139,443,445,1433,1521,3306,5432,6379,7001,8000,8080,8089,9000,9200,11211,27017,80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880",
    },
    {
        text: "高危端口",
        value: "21,22,23,53,80,443,8080,8000,139,445,3389,1521,3306,6379,7001,2375,27017,11211",
    },
    {
        text: "全端口",
        value: "1-65535"
    },
    {
        text: "三高一弱",
        value: "21,22,23,25,53,69,110,111,135,139,143,161,389,445,873,1025,1433,1158,1521,3306,3389,3690,5432,5900,5901,6379,7001,7002,9000,9043,9200,9300,27017,27018,28017,50030,50060,50070,1099,2049,2181,2222,2375,2379,2888,3128,3888,4000,4040,4440,4848,4899,5000,5005,5601,5631,5632,5984,6123,7051,7077,7180,7182,7848,8019,8020,8042,8048,8051,8069,8080,8081,8083,8086,8088,8161,8443,8649,8848,8880,8888,9001,9042,9083,9092,9100,9990,10000,11000,11111,11211,18080,19888,20880,25000,25010,50000,50090,60000,60010,60030"
    },
    {
        text: "自定义",
        value: ""
    },
]

var dict = ({
    usernames: [
        {
            name: "FTP",
            dic: ["ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"],
        },
        {
            name: "SSH",
            dic: ["root", "admin"]
        },
        {
            name: "Telnet",
            dic: ["root", "admin"]
        },
        {
            name: "SMB",
            dic: ["administrator", "admin", "guest"]
        },
        {
            name: "Mssql",
            dic: ["sa", "sql"]
        },
        {
            name: "Oracle",
            dic: ["sys", "system", "admin", "test", "web", "orcl"]
        },
        {
            name: "Mysql",
            dic: ["root", "mysql"]
        },
        {
            name: "RDP",
            dic: ["administrator", "admin", "guest"]
        },
        {
            name: "Postgresql",
            dic: ["postgres", "admin"]
        },
        {
            name: "VNC",
            dic: ["admin", "administrator", "root"]
        },
    ],
    passwords: [] as string[],
})

var jsfind = reactive({
    whiteList: "github.com\ngoogle.com\namazon.com\ngitee.com\nw3.org\nqq.com",
    defaultType: ["info", "warning", "danger", "primary", "success", "info"]
})

var onlineOptions = [
    {
        label: 'aside.space_engine',
        value: [
            {
                name: "FOFA",
                url: "https://fofa.info/",
                icon: "/fofa.ico"
            },
            {
                name: "Hunter",
                url: "https://hunter.qianxin.com/",
                icon: "/hunter.ico"
            },
            {
                name: "Quake",
                url: "https://quake.360.net/",
                icon: "/navigation/360.ico"
            },
            {
                name: "Censys",
                url: "https://search.censys.io/",
                icon: "/navigation/censys.ico"
            },
            {
                name: "Shodan",
                url: "https://www.shodan.io/",
                icon: "/navigation/shodan.png"
            },
            {
                name: "Zoomeye",
                url: "https://www.zoomeye.org/",
                icon: "/navigation/zoomeye.ico"
            },
            {
                name: "零零信安",
                url: "https://0.zone/",
                icon: "/navigation/0zone.png"
            },
            {
                name: "DayDayMap",
                url: "https://www.daydaymap.com/",
                icon: "/navigation/daydaymap.svg"
            },
        ],
    },
    {
        label: 'navigator.enterprise',
        value: [
            {
                name: "工信部备案管理系统",
                url: "https://beian.miit.gov.cn/#/",
                icon: "/navigation/gxb.png",
            },
            {
                name: "天眼查",
                url: "https://www.tianyancha.com/",
                icon: "/navigation/tianyancha.png"
            },
            {
                name: "爱企查",
                url: "https://aiqicha.baidu.com/",
                icon: "/navigation/aiqicha.ico"
            },
            {
                name: "企查查",
                url: "https://www.qcc.com/",
                icon: "/navigation/qichacha.png"
            },
            {
                name: "小蓝本",
                url: "https://sou.xiaolanben.com/pc",
                icon: "/navigation/xiaolanben.png"
            },
        ],
    },
    {
        label: 'navigator.domaininfo',
        value: [
            {
                name: "Securitytrails*",
                url: "https://securitytrails.com/",
                icon: "/navigation/st.ico"
            },
            {
                name: "IP138",
                url: "https://site.ip138.com/",
                icon: "/navigation/ip138.ico"
            },
            {
                name: "爱站网",
                url: "https://dns.aizhan.com/",
                icon: "/navigation/azw.ico"
            },
            {
                name: "IPUU IP定位",
                url: "https://www.ipuu.net/Home",
                icon: "/navigation/awkj.ico"
            },
        ],
    },
    {
        label: 'navigator.threat_intelligence',
        value: [
            {
                name: "微步",
                url: "https://x.threatbook.com/",
                icon: "/navigation/x.ico"
            },
            {
                name: "VirusTotal",
                url: "https://www.virustotal.com/gui/home/upload",
                icon: "/navigation/vt.svg"
            },
            {
                name: "360安全大脑",
                url: "https://ti.360.cn/",
                icon: "/navigation/360aqdn.png"
            },
            {
                name: "安全星图平台",
                url: "https://ti.dbappsecurity.com.cn/",
                icon: "/navigation/dbapp.png"
            },
            {
                name: "深信服",
                url: "https://ti.sangfor.com.cn/",
                icon: "/navigation/sangfor.ico"
            },
            {
                name: "NTI",
                url: "https://ti.nsfocus.com/",
                icon: "/navigation/nsfocus.ico"
            },
            {
                name: "VenusEye",
                url: "https://www.venuseye.com.cn/",
                icon: "/navigation/venuseye.ico"
            },
        ],
    },
    {
        label: 'navigator.technical_forum',
        value: [
            {
                name: "先知社区",
                url: "https://xz.aliyun.com/",
                icon: "/navigation/xianzhi.ico"
            },
            {
                name: "奇安信攻防社区",
                url: "https://forum.butian.net/",
                icon: "/navigation/butian.png"
            },
            {
                name: "FREEBUF",
                url: "https://www.freebuf.com/",
                icon: "/navigation/freebuf.ico"
            },
            {
                name: "跳跳糖",
                url: "https://tttang.com/",
                icon: "/navigation/ttt.png"
            },
            {
                name: "Tools",
                url: "https://www.t00ls.com/",
                icon: "/navigation/tools.ico"
            },
            {
                name: "看雪",
                url: "https://bbs.kanxue.com/",
                icon: "/navigation/kanxue.ico"
            },
        ],
    },
    {
        label: 'navigator.dnslog',
        value: [
            {
                name: "Dnslog",
                url: "https://dnslog.cn/",
                icon: "/navigation/dnslog.ico"
            },
            {
                name: "Dig",
                url: "https://dig.pm/",
                icon: "/navigation/dig.png"
            },
            {
                name: "Eyes",
                url: "http://eyes.sh/",
                icon: "/navigation/eyes.png"
            },
        ],
    },
    {
        label: 'aside.en_and_de',
        value: [
            {
                name: "JS混淆解密",
                url: "https://dev-coco.github.io/Online-Tools/JavaScript-Deobfuscator.html#google_vignette",
                icon: "/navigation/jsob.png"
            },
            {
                name: "CyberChef*",
                url: "https://gchq.github.io/CyberChef/",
                icon: "/navigation/cyberchef.png"
            },
        ],
    },
]

export default {
    space,
    proxy,
    dict,
    Logger,
    LOCAL_VERSION,
    PATH,
    UPDATE,
    jsfind,
    webscan,
    onlineOptions,
    portGroup,
    temp,
    Language,
    Theme
};
