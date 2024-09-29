import { reactive, ref } from 'vue'

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
    port_thread: 1000,
    port_timeout: 7,
    ping_check_alive: false,
    default_alive_module: "None",
    default_network: "Auto",
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
    isGrid: true,
})

const Logger = reactive({
    value: '',
    length: 100, // 日志显示条数
})

const LOCAL_VERSION = "1.6.4"

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
        text: "两高一弱",
        value: "21,22,23,25,53,69,110,111,135,139,143,161,389,445,873,1025,1433,1158,1521,3306,3389,3690,5432,5900,5901,6379,7001,7002,9000,9043,9200,9300,27017,27018,28017,50030,50060,50070,1099,2049,2181,2222,2375,2379,2888,3128,3888,4000,4040,4440,4848,4899,5000,5005,5601,5631,5632,5984,6123,7051,7077,7180,7182,7848,8019,8020,8042,8048,8051,8069,8080,8081,8083,8086,8088,8161,8443,8649,8848,8880,8888,9001,9042,9083,9092,9100,9990,10000,11000,11111,11211,18080,19888,20880,25000,25010,50000,50090,60000,60010,60030"
    },
    {
        text: "自定义",
        value: ""
    },
]

var dict = ({
    options: ["ftp", "ssh", "telnet", "smb", "oracle", "mssql", "mysql", "rdp", "postgresql", "vnc", "redis", "memcached", "mongodb"],
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
    portGroup,
    temp,
    Language,
    Theme
};
