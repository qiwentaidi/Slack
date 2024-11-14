<script lang="ts" setup>
import { reactive, ref } from "vue";
import { ExtractIP, IpLocation } from "wailsjs/go/main/App";
import { Search, QuestionFilled, Location, Cellphone, Postcard, Upload, Document, ChromeFilled, DocumentCopy, Filter, Camera, Box, RefreshLeft } from "@element-plus/icons-vue";
import { ElMessage, ElNotification } from "element-plus";
import extractIcon from "@/assets/icon/extract.svg";
import { Copy, SplitTextArea, UploadFileAndRead } from "@/util";
import async from "async";
import { regexpIdCard, regexpPhone } from "@/stores/validate";
import { FileDialog, ReadFile } from "wailsjs/go/main/File";
import passwordIcon from "@/assets/icon/password.svg";
import bugIcon from "@/assets/icon/bug.svg";
import fingerprintIcon from "@/assets/icon/fingerprint.svg";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import usePagination from "@/usePagination";
import { FormatOutput, ConnectAndExecute } from "wailsjs/go/core/Tools";

const form = reactive({
    activeTab: 'fscan',
    result: "",
    input: "",
    dedupSplit: "",
});

const fscan = reactive({
    input: "",
    result: "",
    weakpass: [] as { type: string, ip: string, port: string, username: string, password: string, extend: string }[],
    virus: [] as { url: string, pocinfo: string, extend: string }[],
    fingerprint: [] as { url: string, fingerprint: string }[],
});

const wip = usePagination<{ url: string, code: string, title: string, length: string, redirect: string }>(20)

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

async function selectFilePath() {
    fscan.input = await FileDialog("*.txt");
}

const weekProtocol = ["ftp", "ssh", "telnet", "mysql", "oracle", "mssql", "postgres", "rdp", "mongodb", "redis"]
async function fscanParse() {
    let file = await ReadFile(fscan.input);
    let result = await FormatOutput(file.Content)
    if (!result) return
    fscan.weakpass = []
    fscan.virus = []

    for (const [key, values] of Object.entries(result)) {
        fscan.result += `[${key}]\n${values.join("\n")}\n\n`
        if (weekProtocol.includes(key.toLocaleLowerCase())) {
            values.map(line => {
                line = line.slice(4)
                const parts = line.split(' ');
                const protocol = parts[0];
                let [ip, port, username] = parts[1].split(':');
                let password = parts[2] || '';
                if (password == "unauthorized") {
                    username = password;
                    password = ""
                }
                const ext = parts[3] || '';
                fscan.weakpass.push({
                    type: protocol,
                    ip: ip,
                    port: port,
                    username: username,
                    password: password,
                    extend: ext
                })
            });
        }
        if (key == "Web INFO") {
            wip.table.result = values.map(line => {
                const urlMatch = line.match(/WebTitle (http[^\s]+)/);
                const codeMatch = line.match(/code:(\d+)/);
                const lenMatch = line.match(/len:(\d+)/);
                const titleMatch = line.match(/title:(.*?)(?:\s+跳转url|$)/);  // 捕获完整的标题直到跳转url或行尾
                const redirectMatch = line.match(/跳转url: ([^\s]+)/);

                const url = urlMatch ? urlMatch[1].trim() : '';
                const code = codeMatch ? codeMatch[1].trim() : '';
                const length = lenMatch ? lenMatch[1].trim() : '';
                const title = titleMatch ? titleMatch[1].trim() : '';

                const redirect = redirectMatch ? redirectMatch[1].trim() : '';

                return {
                    url: url,
                    code: code,
                    title: title,
                    length: length,
                    redirect: redirect,
                };
            });
            wip.table.pageContent = wip.ctrl.watchResultChange(wip.table)
        }

        if (key == "POC") {
            values.map(line => {
                const urlMatch = line.match(/PocScan (http[^\s]+)/);
                const pocMatch = line.match(/poc([\w-]+)/);
                fscan.virus.push({
                    url: urlMatch ? urlMatch[1].trim() : '',
                    pocinfo: pocMatch ? "poc" + pocMatch[1].trim() : '',
                    extend: ""
                })
            });
        }

        if (key == "MS17-010") {
            values.map(line => {
                const hostMatch = line.match(/MS17-010 ([^\s]+)/);
                const versionMatch = line.match(/\((.+)\)/);
                fscan.virus.push({
                    url: hostMatch ? hostMatch[1] : "",
                    pocinfo: "MS17-010",
                    extend: versionMatch ? versionMatch[1] : "",
                })
            });
        }

        if (key == "INFO") {
            fscan.fingerprint = values.map(line => {
                line = line.slice(4); // 移除[+]
                const urlMatch = line.match(/InfoScan (http[^\s]+)/);
                const fingerprintsMatch = line.match(/\[(.+)\]/);
                return {
                    url: urlMatch ? urlMatch[1].trim() : '',
                    fingerprint: fingerprintsMatch ? fingerprintsMatch[1] : '',
                }
            });
        }
    }
}

const expandedRowKeys = ref([] as string[])
function generateRowKey(row: any) {
    return `${row.ip}-${row.port}`;
}


async function specialHW(row: any) {
    ElMessage.info("正在连接请稍等...")
    const result = await ConnectAndExecute(row.type, row.ip, row.port, row.username, row.password);

    // Find the matching entry in fscan.weakpass by IP and port
    const matchedEntry = fscan.weakpass.find(entry => entry.ip === row.ip && entry.port === row.port);

    // If a matching entry is found, update its 'extend' field with the result
    if (matchedEntry) {
        matchedEntry.extend = result;
    }
    expandedRowKeys.value = [generateRowKey(row)]
}


const filter = ({
    hikvsion: ["2512", "600", "481", "480"],
    vm: ["ID_VC_Welcome", "ID_EESX_Welcome"],
    reset: () => {
        if (wip.table.temp.length != 0) wip.table.result = wip.table.temp
        wip.table.pageContent = wip.ctrl.watchResultChange(wip.table);
    },
    hikvisionFilter: () => {
        if (wip.table.temp.length == 0) wip.table.temp = wip.table.result
        wip.table.result = wip.table.temp.filter(item => filter.hikvsion.includes(item.length))
        wip.table.pageContent = wip.ctrl.watchResultChange(wip.table)
    },
    vmFilter: () => {
        if (wip.table.temp.length == 0) wip.table.temp = wip.table.result
        wip.table.result = wip.table.temp.filter(item => {
            for (const name of filter.vm) {
                if (item.title.includes(name)) return item
            }
        });
        wip.table.pageContent = wip.ctrl.watchResultChange(wip.table)
    },
    pocinfo(value: string, row: any): boolean {
        return row.pocinfo === value;
    },
    weakpass(value: string, row: any): boolean {
        return row.type.toLowerCase() == value; 
    },
})
</script>

<template>

    <el-tabs v-model="form.activeTab" :stretch="true">
        <el-tab-pane label="Fscan结果提取" name="fscan">
            <div class="head">
                <el-input v-model="fscan.input" placeholder="移入或者选择文件路径">
                    <template #suffix>
                        <el-button link :icon="Document" @click="selectFilePath"></el-button>
                    </template>
                </el-input>
                <el-button type="primary" @click="fscanParse" style="margin-left: 10px">开始解析</el-button>
            </div>
            <el-tabs type="border-card" style="margin-top: 15px;" class="demo-tabs">
                <el-tab-pane>
                    <template #label>
                        <span class="custom-tabs-label">
                            <el-icon>
                                <Reading />
                            </el-icon>
                            <span>Raw</span>
                        </span>
                    </template>
                    <el-input v-model="fscan.result" type="textarea" resize="none"
                        style="height: calc(100vh - 250px);" />
                </el-tab-pane>

                <el-tab-pane>
                    <template #label>
                        <span class="custom-tabs-label">
                            <el-icon>
                                <passwordIcon />
                            </el-icon>
                            <span>WeakPass({{ fscan.weakpass.length }})</span>
                        </span>
                    </template>
                    <el-table :data="fscan.weakpass" :row-key="generateRowKey" :expand-row-keys="expandedRowKeys"
                        style="height: calc(100vh - 250px);">
                        <el-table-column type="expand">
                            <template #default="props">
                                <highlightjs language="bash" :code="props.row.extend" style="padding-inline: 20px; border: 1px solid #ccc;">
                                </highlightjs>
                            </template>
                        </el-table-column>
                        <el-table-column type="index" width="60px" />
                        <el-table-column prop="type" label="Type" width="120px" :filters="[
                            { text: 'ftp', value: 'ftp' },
                            { text: 'ssh', value: 'ssh' },
                            { text: 'mysql', value: 'mysql' },
                            { text: 'mssql', value: 'mssql' },
                            { text: 'oracle', value: 'oracle' },
                            { text: 'redis', value: 'redis' },
                            { text: 'mongodb', value: 'mongodb' },
                            { text: 'memcached', value: 'memcached' },
                        ]" :filter-method="filter.weakpass">
                        <template #filter-icon>
                            <Filter />
                        </template>
                        </el-table-column>
                        <el-table-column prop="ip" label="IP" />
                        <el-table-column prop="port" label="Port" width="80" />
                        <el-table-column prop="username" label="Username" />
                        <el-table-column prop="password" label="Password" />
                        <el-table-column label="Operations" align="center">
                            <template #default="props">
                                <el-button type="primary" link @click="specialHW(props.row)">Connect &
                                    Command</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>

                <el-tab-pane>
                    <template #label>
                        <span class="custom-tabs-label">
                            <el-icon>
                                <bugIcon />
                            </el-icon>
                            <span>Vulnerability({{ fscan.virus.length }})</span>
                        </span>
                    </template>
                    <el-table :data="fscan.virus" style="height: calc(100vh - 250px);">
                        <el-table-column type="index" width="60px" />
                        <el-table-column prop="url" label="URL" :show-overflow-tooltip="true" />
                        <el-table-column prop="pocinfo" label="PocInfo" :filters="[
                            { text: 'MS17-010', value: 'MS17-010' },
                            { text: 'poc-yaml-alibaba-nacos', value: 'poc-yaml-alibaba-nacos' },
                        ]" :filter-method="filter.pocinfo">
                        <template #filter-icon>
                            <Filter />
                        </template>
                        </el-table-column>
                        <el-table-column prop="extend" label="Extend" width="200px" />
                        <el-table-column label="Operations" width="100" align="center">
                            <template #default="scope">
                                <el-button link :icon="ChromeFilled" @click="BrowserOpenURL(scope.row.url)"></el-button>
                                <el-button link :icon="DocumentCopy" @click="Copy(scope.row.url)"></el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>

                <el-tab-pane>
                    <template #label>
                        <span class="custom-tabs-label">
                            <el-icon>
                                <fingerprintIcon />
                            </el-icon>
                            <span>Fingerprints({{ fscan.fingerprint.length }})</span>
                        </span>
                    </template>
                    <el-table :data="fscan.fingerprint" style="height: calc(100vh - 250px);">
                        <el-table-column type="index" width="60px" />
                        <el-table-column prop="url" label="URL" :show-overflow-tooltip="true" />
                        <el-table-column prop="fingerprint" label="Fingerprints" />
                        <el-table-column label="Operations" width="100" align="center">
                            <template #default="scope">
                                <el-button link :icon="ChromeFilled" @click="BrowserOpenURL(scope.row.url)"></el-button>
                                <el-button link :icon="DocumentCopy" @click="Copy(scope.row.url)"></el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>

                <el-tab-pane>
                    <template #label>
                        <span class="custom-tabs-label">
                            <el-icon>
                                <ChromeFilled />
                            </el-icon>
                            <span>Web({{ wip.table.result.length }})</span>
                        </span>
                    </template>
                    <el-table :data="wip.table.pageContent" style="height: calc(100vh - 285px);">
                        <el-table-column prop="url" width="400">
                            <template #header>
                                <el-text><span>URL</span>
                                    <el-divider direction="vertical" />
                                    <el-dropdown>
                                        <el-button :icon="Filter" size="small" bg text>筛选</el-button>
                                        <template #dropdown>
                                            <el-dropdown-menu>
                                                <el-dropdown-item :icon="Camera"
                                                    @click="filter.hikvisionFilter">海康摄像头</el-dropdown-item>
                                                <el-dropdown-item :icon="Box" @click="filter.vmFilter">Vcenter &
                                                    Exsi</el-dropdown-item>
                                                <el-dropdown-item :icon="RefreshLeft" divided
                                                    @click="filter.reset">重置</el-dropdown-item>
                                            </el-dropdown-menu>
                                        </template>
                                    </el-dropdown>
                                </el-text>
                            </template>
                        </el-table-column>
                        <el-table-column prop="title" label="Title" width="280" />
                        <el-table-column prop="code" label="Code" width="100" />
                        <el-table-column prop="length" label="Length" width="120" />
                        <el-table-column prop="redirect" label="Redirect" width="400" />
                        <el-table-column label="Operate" fixed="right" width="100" align="center">
                            <template #default="scope">
                                <el-button link :icon="ChromeFilled" @click="BrowserOpenURL(scope.row.url)"></el-button>
                                <el-button link :icon="DocumentCopy" @click="Copy(scope.row.url)"></el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                    <div class="my-header" style="margin-top: 10px;">
                        <div></div>
                        <el-pagination size="small" background @size-change="wip.ctrl.handleSizeChange"
                            @current-change="wip.ctrl.handleCurrentChange" :current-page="wip.table.currentPage"
                            :page-sizes="[50, 100, 200]" :page-size="wip.table.pageSize"
                            layout="total, sizes, prev, pager, next" :total="wip.table.result.length">
                        </el-pagination>
                    </div>
                </el-tab-pane>
            </el-tabs>

        </el-tab-pane>
        <el-tab-pane label="其他信息提取" name="other">
            <el-form :model="form" label-width="50px">
                <el-form-item label="内容">
                    <el-input v-model="form.input" type="textarea" :rows="7" placeholder="请输入内容" />
                    <el-button link size="small" :icon="Upload" @click="uploadFile"
                        style="margin-top: 5px">导入文件内容</el-button>
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
        </el-tab-pane>
    </el-tabs>
</template>
