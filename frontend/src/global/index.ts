import { reactive, ref } from 'vue'
var space = reactive({
    fofaapi: 'https://fofa.info/',
    fofaemail: '',
    fofakey: '',
    hunterkey: '',
    quakekey: ''
})
var scan = reactive({
    dns1: "114.114.114.114",
    dns2: "223.5.5.5",
    whiteList: "github.com\ngoogle.com\namazon.com\ngitee.com\nw3.org\nqq.com",
})
var proxy = reactive({
    enabled: false,
    mode: 'HTTP',
    address: '127.0.0.1',
    port: 8080,
    username: '',
    password: '',
})

const Logger = ref("")

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
    passwords: ["123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "1234567", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789"],
})

export default {
    space,
    scan,
    proxy,
    dict,
    Logger,
};
