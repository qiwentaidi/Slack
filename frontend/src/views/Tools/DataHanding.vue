<script lang="ts" setup>
import { h, reactive, ref } from "vue";
import { ExtractIP, IpLocation } from "wailsjs/go/services/App";
import { UserFilled, Filter, Location, Cellphone, Postcard, Upload, DeleteFilled } from "@element-plus/icons-vue";
import alibabaIcon from "@/assets/icon/alibaba.svg";
import onlineIcon from "@/assets/icon/online.svg";
import { ElMessage } from "element-plus";
import { ProcessTextAreaInput, UploadFileAndRead } from "@/util";
import async from "async";
import { regexpIdCard, regexpPhone, regexAlibabaDruidWebURI, regexURL, regexAlibabaDruidWebSession } from "@/stores/validate";
import ContextMenu from '@imengyu/vue3-context-menu'
import { defaultIconSize } from "@/stores/style";

const input = ref("");

const form = reactive({
    result: "",
    dedupSplit: "",
    tips: "",
});

async function uploadFile() {
    input.value = await UploadFileAndRead();
}

function isEmpty() {
    if (input.value.length == 0) {
        ElMessage.warning({
            message: "请输入待处理的内容或者文件",
            grouping: true,
        })
        return false
    }
    return true
}

function deduplication() {
    if (!isEmpty()) {
        return
    }
    let lines = [] as string[];
    if (form.dedupSplit == "\\n") {
        lines = input.value.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
    } else {
        lines = input.value.split(form.dedupSplit);
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
}

function getDomains(): string[] {
    const urls = getURLs();
    const domains = urls
        .map((url) => {
            try {
                const parsedUrl = new URL(url);
                return parsedUrl.hostname;
            } catch (e) {
                
            }
        })
        .filter((domain) => domain); // 过滤掉空字符串
    return Array.from(new Set(domains));
}

function getURLs(): string[] {
    if (!isEmpty()) {
        return []
    }
    const urls = input.value.match(regexURL) || [];
    return Array.from(new Set(urls));
}

function getPhoneNumbers(): string[] {
    if (!isEmpty()) {
        return []
    }
    const phoneNumbers = input.value.match(regexpPhone) || [];
    return Array.from(new Set(phoneNumbers));
}

function getIdCards(): string[] {
    if (!isEmpty()) {
        return []
    }
    const idcards = input.value.match(regexpIdCard) || [];
    return Array.from(new Set(idcards));
}

function getAlibabaDruidWebURI() {
    if (!isEmpty()) {
        return []
    }
    const uris = input.value.match(regexAlibabaDruidWebURI) || [];
    form.tips = `共提取 ${uris.length} 个结果`
    return Array.from(new Set(uris));
}

function getAlibabaDruidWebSession() {
    if (!isEmpty()) {
        return []
    }
    const sessions = input.value.match(regexAlibabaDruidWebSession) || [];
    form.tips = `共提取 ${sessions.length} 个结果`
    return Array.from(new Set(sessions));
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
                            form.result = getAlibabaDruidWebURI().join("\n")
                        }
                    },
                    { 
                        label: "提取WebSession", 
                        onClick: () => {
                            form.result = getAlibabaDruidWebSession().join("\n")
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
                            if (!isEmpty()) {
                                return
                            }
                            ExtractIP(input.value).then((result) => {
                                form.result = result;
                            });
                        }
                    },
                    { 
                        label: "提取URL", 
                        onClick: () => {
                            const urls = getURLs();
                            form.result = urls.join("\n");
                            form.tips = `共提取 ${urls.length} 个结果`
                        }
                    },
                    { 
                        label: "提取URL中的域名", 
                        divided: true,
                        onClick: () => {
                            const domains = getDomains();
                            form.result = domains.join("\n");
                            form.tips = `共提取 ${domains.length} 个结果`
                        }
                    },
                    { 
                        label: "IP定位查询", 
                        icon: h(Location, defaultIconSize),
                        onClick: () => {
                            if (!isEmpty()) {
                                return
                            }
                            let lines = ProcessTextAreaInput(input.value);
                            form.result = "";
                            async.eachLimit(lines, 20, async (ip: string, callback: () => void) => {
                                let result = await IpLocation(ip);
                                form.result += `${ip}  |  ${result}\n`;
                            });
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
                            const phones = getPhoneNumbers()
                            form.result = phones.join("\n");
                            form.tips = `共提取 ${phones.length} 个结果`
                        }
                    },
                    { 
                        label: "提取身份证", 
                        icon: h(Postcard, defaultIconSize),
                        onClick: () => {
                            const cards = getIdCards()
                            form.result = cards.join("\n");
                            form.tips = `共提取 ${cards.length} 个结果`
                        }
                    },
                ]
            },
            { 
                label: "清空文本框", 
                icon: h(DeleteFilled, defaultIconSize),
                onClick: () => {
                    input.value = ""
                }
            }
        ]
    });
}

const code = `请输入内容，输出处理等功能通过右键菜单进行调用

druid数据提取需要输出响应包内容

由于文本框性能问题，请减少大文本内容的改动，否则可能造成卡顿或卡死。
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
        <el-input v-model="input" type="textarea" :rows="10"
        :placeholder="code" @contextmenu.stop.prevent="handleContextMenu"></el-input>
        <div class="action-area">
            <el-button :icon="Upload" size="small" @click="uploadFile">Upload</el-button>
        </div>
    </div>
    <el-input v-model="form.result" type="textarea" :rows="20" readonly />
    <span class="form-item-tips">{{ form.tips }}</span>
</template>
