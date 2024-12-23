<script lang="ts" setup>
import { h, reactive } from "vue";
import { ExtractIP, IpLocation } from "wailsjs/go/services/App";
import { UserFilled, Filter, Location, Cellphone, Postcard, Upload, DeleteFilled, Unlock, DocumentCopy } from "@element-plus/icons-vue";
import alibabaIcon from "@/assets/icon/alibaba.svg";
import onlineIcon from "@/assets/icon/online.svg";
import { ElMessage, ElMessageBox } from "element-plus";
import { Copy, getFilepath, ProcessTextAreaInput } from "@/util";
import async from "async";
import { regexpIdCard, regexpPhone, regexAlibabaDruidWebURI, regexURL, regexAlibabaDruidWebSession, regexFileProtocol } from "@/stores/validate";
import ContextMenu from '@imengyu/vue3-context-menu'
import { defaultIconSize } from "@/stores/style";
import { CheckFileStat, FileDialog, ReadFile, SaveToTempFile } from "wailsjs/go/services/File";

let content = "" // 实际要处理的内容

const form = reactive({
    input: "",
    result: "",
    dedupSplit: "",
    tips: "",
});

async function uploadFile() {
    let filepath = await FileDialog("");
    if (!filepath) return
    form.input = `{{file://${filepath}}}`
}

async function processInput() {
    if (form.input.length == 0) {
        ElMessage.warning({
            message: "请输入待处理的内容或者文件",
            grouping: true,
        })
        return false
    }
    let path = form.input.match(regexFileProtocol) || []
    if (path && path.length == 0) {
        content = form.input
        return true
    }
    let filepath = getFilepath(path[0])
    let isStat = await CheckFileStat(filepath)
    if (!isStat) {
        ElMessage.warning('文件不存在!')
        return false
    }
    let file = await ReadFile(filepath)
    content = file.Content
    return true
}

function handlePaste(event: any) {
    const clipboardData = event.clipboardData
    const pastedData = clipboardData.getData('Text');
    if (pastedData.length > 100 * 1024) { // 检查输入内容是否大于100KB
        event.preventDefault(); // 阻止输入
        ElMessageBox.confirm('粘贴的内容过长，是已转换为临时文件存储？', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
        }).then(async () => {
            let tempfile = await SaveToTempFile(pastedData)
            form.input = `{{file://${tempfile}}}`; // 更新文本框内容为文件路径
        });
    }
}

function deduplication() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        let lines = [] as string[];
        if (form.dedupSplit == "\\n") {
            lines = content.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
        } else {
            lines = content.split(form.dedupSplit);
        }
        lines = lines.filter((item) => item.trim() !== ""); // 删除空项并去除左右空格
        let uniqueArray = Array.from(new Set(lines));
        if (lines.length === uniqueArray.length) {
            ElMessage.info({
                message: "不存在重复数据",
                grouping: true,
            });
            return;
        }
        ElMessage.success({
            message: `已去重数据${lines.length - uniqueArray.length}条`,
            grouping: true,
        });
        form.result = uniqueArray.join("\n");
    })
}

function getDomains() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const urls = content.match(regexURL) || [];
        const uniqueArray = Array.from(new Set(urls))
        const domains = urls
        .map((url: string) => {
            try {
                const parsedUrl = new URL(url);
                return parsedUrl.hostname;
            } catch (e) {

            }
        })
        .filter((domain) => domain); // 过滤掉空字符串
        form.result = domains.join("\n");
        form.tips = `共提取 ${domains.length} 个结果`
    })    
}

function getURLs() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const urls = content.match(regexURL) || [];
        const uniqueArray = Array.from(new Set(urls))
        form.result = uniqueArray.join("\n");
        form.tips = `共提取 ${uniqueArray.length} 个结果`
    })    
}

function getPhoneNumbers() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const phoneNumbers = content.match(regexpPhone) || [];
        const uniqueArray = Array.from(new Set(phoneNumbers))
        form.result = uniqueArray.join("\n");
        form.tips = `共提取 ${uniqueArray.length} 个结果`
    })  
}

function getIdCards() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const idcards = content.match(regexpIdCard) || [];
        const uniqueArray = Array.from(new Set(idcards))
        form.result = uniqueArray.join("\n");
        form.tips = `共提取 ${uniqueArray.length} 个结果`
    }) 
}

function getAlibabaDruidWebURI() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const uris = content.match(regexAlibabaDruidWebURI) || [];
        form.result = uris.join("\n");
        form.tips = `共提取 ${uris.length} 个结果`
    }) 
}

function getAlibabaDruidWebSession() {
    processInput().then(isSuccess => {
        if (!isSuccess) {
            return
        }
        const sessions = content.match(regexAlibabaDruidWebSession) || [];
        form.result = sessions.join("\n");
        form.tips = `共提取 ${sessions.length} 个结果`
    })
}

// 帆软v8解密
const PASSWORD_MASK_ARRAY: number[] = [19, 78, 10, 15, 100, 213, 43, 23]; // 掩码

function decryptWinSCPIniContent(iniContent: string) {
    const WSCP_SIMPLE_MAGIC = 0xA3;
    const WSCP_SIMPLE_STRING = '0123456789ABCDEF';
    const WSCP_SIMPLE_FLAG = 0xFF;
    const WSCP_SIMPLE_INTERNAL = 0x00;
    let WSCP_CHARS = [];

    function _simple_decrypt_next_char() {
        if (WSCP_CHARS.length == 0) {
            return WSCP_SIMPLE_INTERNAL;
        }

        const a = WSCP_SIMPLE_STRING.indexOf(WSCP_CHARS.shift());
        const b = WSCP_SIMPLE_STRING.indexOf(WSCP_CHARS.shift());

        return WSCP_SIMPLE_FLAG & ~(((a << 4) + b << 0) ^ WSCP_SIMPLE_MAGIC);
    }

    function decrypt(username: string, hostname: string, encodedPassword: string) {
        if (!encodedPassword.match(/[A-F0-9]+/)) {
            return '';
        }

        const result = [];
        const key = [username, hostname].join('');

        WSCP_CHARS = encodedPassword.split('');

        const flag = _simple_decrypt_next_char();
        let length;

        if (flag == WSCP_SIMPLE_FLAG) {
            _simple_decrypt_next_char();
            length = _simple_decrypt_next_char();
        } else {
            length = flag;
        }

        WSCP_CHARS = WSCP_CHARS.slice(_simple_decrypt_next_char() * 2);

        for (let i = 0; i < length; i++) {
            result.push(String.fromCharCode(_simple_decrypt_next_char()));
        }

        if (flag == WSCP_SIMPLE_FLAG) {
            const valid = result.slice(0, key.length).join('');

            if (valid != key) {
                return '';
            } else {
                return result.slice(key.length).join('');
            }
        }

        return result.join('');
    }

    const sessions = iniContent.split('\n[').filter(session => session.startsWith('Sessions\\'));
    const decodedSessions = sessions.map(session => {
        const sessionName = decodeURIComponent(session.split(']')[0].replace('Sessions\\', ''));
        const usernameLine = session.split('\n').find(line => line.startsWith('UserName='));
        const hostnameLine = session.split('\n').find(line => line.startsWith('HostName='));
        const passwordLine = session.split('\n').find(line => line.startsWith('Password='));

        if (usernameLine && hostnameLine && passwordLine) {
            const username = usernameLine.split('=')[1].trim();
            const hostname = hostnameLine.split('=')[1].trim();
            const encodedPassword = passwordLine.split('=')[1].trim();
            const decodedPassword = decrypt(username, hostname, encodedPassword);
            return {
                session: sessionName,
                username,
                hostname,
                decodedPassword: decodedPassword || '解密失败'
            };
        }
        return {
            session: sessionName,
            decodedPassword: '密码未找到'
        };
    });

    return decodedSessions;
}

function handleContextMenu(e: MouseEvent) {
    //prevent the browser's default menu
    e.preventDefault();
    //show our menu
    ContextMenu.showContextMenu({
        x: e.x,
        y: e.y,
        items: [
            {
                label: "Alibaba Druid",
                icon: h(alibabaIcon, defaultIconSize),
                children: [
                    {
                        label: "提取WebURI",
                        onClick: () => {
                            getAlibabaDruidWebURI()
                        }
                    },
                    {
                        label: "提取WebSession",
                        onClick: () => {
                            getAlibabaDruidWebSession()
                        }
                    },
                ]
            },
            {
                label: "网络信息",
                divided: true,
                icon: h(onlineIcon, defaultIconSize),
                children: [
                    {
                        label: "提取IP",
                        onClick: () => {
                            processInput().then(isSuccess => {
                                if (!isSuccess) {
                                    return
                                }
                                ExtractIP(content).then((result) => {
                                    form.result = result;
                                });
                            })
                            
                        }
                    },
                    {
                        label: "提取URL",
                        onClick: () => {
                            getURLs();
                        }
                    },
                    {
                        label: "提取URL中的域名",
                        divided: true,
                        onClick: () => {
                            getDomains();
                        }
                    },
                    {
                        label: "IP定位查询",
                        icon: h(Location, defaultIconSize),
                        onClick: () => {
                            // if (!processInput()) {
                            //     return
                            // }
                            // let lines = ProcessTextAreaInput(content);
                            // form.result = "";
                            // async.eachLimit(lines, 20, async (ip: string, callback: () => void) => {
                            //     let result = await IpLocation(ip);
                            //     form.result += `${ip}  |  ${result}\n`;
                            // });
                        }
                    },
                ]
            },
            {
                label: "个人敏感信息",
                icon: h(UserFilled, defaultIconSize),
                divided: true,
                children: [
                    {
                        label: "提取手机号",
                        icon: h(Cellphone, defaultIconSize),
                        onClick: () => {
                            getPhoneNumbers()
                        }
                    },
                    {
                        label: "提取身份证",
                        icon: h(Postcard, defaultIconSize),
                        onClick: () => {
                            getIdCards()
                        }
                    },
                ]
            },
            {
                label: "密码解密",
                icon: h(Unlock, defaultIconSize),
                divided: true,
                children: [
                    {
                        label: "WinSCP - WinSCP.ini",
                        onClick: async () => {
                            if (! await processInput()) {
                                return
                            }
                            const decodedSessions = decryptWinSCPIniContent(content);
                            form.result = decodedSessions.map(session =>
                                `Session: ${session.session}\nUsername: ${session.username}\nHostname: ${session.hostname}\nPassword: ${session.decodedPassword}\n`
                            ).join('\n');
                            if (decodedSessions.length === 0) {
                                form.result = "WinSCP - WinSCP.ini 解密失败，请输入完整的WinSCP.ini内容"
                                return
                            }
                            form.tips = `解密完成，共处理 ${decodedSessions.length} 个会话`;
                        }
                    },
                    {
                        label: "Finereport v8 - privilege.xml",
                        onClick: async() => {
                            if (! await processInput()) {
                                return
                            }
                            if (content.length <= 3) {
                                form.result = "Finereport v8 - privilege.xml 解密失败，密码示例: ___0072002a00670066000a"
                                return
                            }
                            let result = ""
                            let temp = content.substring(3); // 截断三位后
                            for (let i = 0; i < temp.length / 4; i++) {
                                let c1: number = parseInt(temp.substring(i * 4, (i + 1) * 4), 16);
                                let c2: number = c1 ^ PASSWORD_MASK_ARRAY[i % 8];
                                result += String.fromCharCode(c2);
                            }
                            form.result = result;
                            form.tips = `解密成功`
                        }
                    },
                    {
                        label: "Seeyon OA - datasourceCtp.properties",
                        onClick: async() => {
                            if (! await processInput()) {
                                return
                            }
                            let pass = content.replace(/\//g, "");
                            let p = pass.split(".0");
                            if (p.length <= 1) {
                                form.result = "Seeyon OA - datasourceCtp.properties 解密失败，密码示例: /1.0/UWJ0dHgxc2U="
                                return
                            }
                            let result = ""
                            let iv = parseInt(p[0]);
                            let password = atob(p[1]);
                            for (let i = 0; i < password.length; i++) {
                                let char = password.charCodeAt(i);
                                result += String.fromCharCode(char - iv);
                            }
                            form.result = result;
                            form.tips = `解密成功`
                        }
                    },
                ]
            },
            {
                label: "清空文本框",
                icon: h(DeleteFilled, defaultIconSize),
                onClick: () => {
                    form.input = ""
                }
            }
        ]
    });
}

const code = `请输入内容，大文本内容会转换成特定的文件形式处理，输出处理等功能通过右键菜单进行调用

druid数据提取需要输出响应包内容
`
</script>

<template>
    <el-input v-model.lazy="form.dedupSplit" placeholder="在此处输入分隔字符后会将数据转换成数组然后去重，换行输入\n">
        <template #prepend>
            数据去重
        </template>
        <template #suffix>
            <el-divider direction="vertical" />
            <el-button type="primary" :icon="Filter" link @click="deduplication">去重</el-button>
        </template>
    </el-input>
    <div class="textarea-container" style="margin-block: 10px;">
        <el-input v-model="form.input" type="textarea" :rows="10" :placeholder="code"
            @contextmenu.stop.prevent="handleContextMenu" @paste="handlePaste"></el-input>
        <div class="action-area">
            <el-button :icon="Upload" size="small" @click="uploadFile">Upload</el-button>
        </div>
    </div>
    <div class="textarea-container">
        <el-input v-model="form.result" type="textarea" :rows="19" readonly />
        <div class="action-area">
            <el-button :icon="DocumentCopy" size="small" @click="Copy(form.result)">Copy</el-button>
        </div>
    </div>
    <span class="form-item-tips">{{ form.tips }}</span>
</template>
