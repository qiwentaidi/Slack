<script lang="ts" setup>
import { reactive } from "vue";
import { ExtractIP, IpLocation } from "wailsjs/go/services/App";
import { Search, QuestionFilled, Location, Cellphone, Postcard, Upload } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
import extractIcon from "@/assets/icon/extract.svg";
import { SplitTextArea, UploadFileAndRead } from "@/util";
import async from "async";
import { regexpIdCard, regexpPhone } from "@/stores/validate";

const form = reactive({
    activeTab: 'fscan',
    result: "",
    input: "",
    dedupSplit: "",
});

async function uploadFile() {
    form.input = await UploadFileAndRead();
}

function Deduplication() {
    let lines = [] as string[];
    if (form.dedupSplit == "\\n") {
        lines = form.input.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
    } else {
        lines = form.input.split(form.dedupSplit);
    }
    lines = lines.filter((item) => item.trim() !== ""); // 删除空项并去除左右空格
    let uniqueArray = Array.from(new Set(lines));
    if (lines.length === uniqueArray.length) {
        ElNotification("不存在重复数据");
        return;
    }
    ElNotification.success(`已去重数据${lines.length - uniqueArray.length}条`);
    form.result = uniqueArray.join("\n");
}

function extractDomains() {
    const urls = getURLs();
    const domains = urls
        .map((url) => {
            try {
                const parsedUrl = new URL(url);
                return parsedUrl.hostname;
            } catch (e) {
                console.error(`Invalid URL: ${url}`);
                return "";
            }
        })
        .filter((domain) => domain); // 过滤掉空字符串
    form.result = Array.from(new Set(domains)).join("\n");
}

function getURLs(): string[] {
    const urlPattern = /https?:\/\/[^\s/$.?#].[^\s]*/g;
    const urls = form.input.match(urlPattern) || [];
    return Array.from(new Set(urls));
}

function getPhoneNumbers(): string[] {
    const phoneNumbers = form.input.match(regexpPhone) || [];
    return Array.from(new Set(phoneNumbers));
}

function getIdCards(): string[] {
    const idcards = form.input.match(regexpIdCard) || [];
    return Array.from(new Set(idcards));
}

const appControl = [
    {
        label: "IP提取",
        type: "success",
        icon: extractIcon,
        action: () => {
            if (form.input != "") {
                ExtractIP(form.input).then((result) => {
                    form.result = result;
                });
            }
        },
    },
    {
        label: "IP定位查询",
        type: "success",
        icon: Location,
        action: () => {
            let lines = SplitTextArea(form.input);
            form.result = "";
            async.eachLimit(lines, 20, async (ip: string, callback: () => void) => {
                let result = await IpLocation(ip);
                form.result += `${ip}  |  ${result}\n`;
            });
        },
    },
    {
        label: "提取URL中的域名",
        type: "info",
        icon: extractIcon,
        action: () => {
            extractDomains();
        },
    },
    {
        label: "URL提取",
        type: "info",
        icon: extractIcon,
        action: () => {
            form.result = getURLs().join("\n");
        },
    },
    {
        label: "手机号提取",
        type: "warning",
        icon: Cellphone,
        action: () => {
            form.result = getPhoneNumbers().join("\n");
        },
    },
    {
        label: "身份证提取",
        type: "warning",
        icon: Postcard,
        action: () => {
            form.result = getIdCards().join("\n");
        },
    },
];

</script>

<template>
    <el-form :model="form" label-width="50px">
        <el-form-item label="内容">
            <div class="textarea-container">
                <el-input class="custom-textarea" v-model="form.input" type="textarea" :rows="7" placeholder="请输入内容" />
                <div class="action-area">
                    <el-button :icon="Upload" size="small" @click="uploadFile">Upload</el-button>
                </div>
            </div>
        </el-form-item>
        <el-form-item label="结果">
            <el-input v-model="form.result" type="textarea" :rows="15" />
        </el-form-item>
        <el-form-item>
            <el-space>
                <el-button v-for="item in appControl" @click="item.action" :type="item.type">
                    <template #icon>
                        <el-icon :size="20">
                            <component :is="item.icon" />
                        </el-icon>
                    </template>
                    {{ item.label }}
                </el-button>
            </el-space>
        </el-form-item>
        <el-form-item>
            <el-input v-model="form.dedupSplit" style="width: 400px;">
                <template #prepend>
                    数据去重
                    <el-tooltip placement="left">
                        <template #content>输入分隔字符后转换成数组，然后去重，换行输入\n</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <template #suffix>
                    <el-divider direction="vertical" />
                    <el-button type="primary" :icon="Search" link @click="Deduplication"></el-button>
                </template>
            </el-input>
        </el-form-item>
    </el-form>
</template>
