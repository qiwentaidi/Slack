<template>
    <div v-show="fscan.showType == 'default'">
        <div class="head">
            <el-button :icon="InfoFilled" @click="fscan.tipsDialog = true">模块介绍</el-button>
            <el-input v-model="fscan.input" placeholder="移入或者选择文件路径" style="margin-inline: 5px;">
                <template #suffix>
                    <el-button link :icon="Document" @click="selectFileAndAssign(fscan, 'input', '*.txt')"></el-button>
                </template>
            </el-input>
            <el-button type="primary" @click="fscanParse">开始解析</el-button>
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
                <el-input v-model="fscan.result" type="textarea" resize="none" style="height: calc(100vh - 200px);" />
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
                    style="height: calc(100vh - 200px);">
                    <el-table-column type="expand">
                        <template #default="props">
                            <highlightjs language="bash" :code="props.row.extend" style="padding-inline: 10px;">
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
                        { text: 'smb', value: 'smb' },
                    ]" :filter-method="filter.weakpass">
                        <template #filter-icon>
                            <Filter />
                        </template>
                    </el-table-column>
                    <el-table-column prop="ip" label="IP" />
                    <el-table-column prop="port" label="Port" width="80" />
                    <el-table-column prop="username" label="Username" />
                    <el-table-column prop="password" label="Password" />
                    <el-table-column label="Operate" align="center">
                        <template #default="props">
                            <el-button :icon="consoleIcon" @click="connectAndExecute(props.row)">Command</el-button>
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
                <el-table :data="fscan.virus" style="height: calc(100vh - 200px);">
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
                    <el-table-column label="Operate" width="100" align="center">
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
                <el-table :data="fscan.fingerprint" style="height: calc(100vh - 200px);">
                    <el-table-column type="index" width="60px" />
                    <el-table-column prop="url" label="URL" :show-overflow-tooltip="true" />
                    <el-table-column prop="fingerprint" label="Fingerprints" />
                    <el-table-column label="Operate" width="100" align="center">
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
                <el-table :data="wip.table.pageContent" style="height: calc(100vh - 235px);">
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
    </div>
    <div id="main" v-show="fscan.showType == 'graph'" style="width: 100%; height: 100vh;"></div>
    <el-dialog v-model="fscan.tipsDialog" title="模块介绍" width="50%">
        <Note type="info">
            Raw: 按行匹配筛选规则<br /><br />
            其他模块: 只有在使用官网项目 <strong>https://github.com/shadow1ng/fscan</strong> 时才能正常使用<br /><br />
            <el-space><el-tag type="warning">其他魔改过输出内容的版本，其他模块无法适配</el-tag><el-tag type="warning">且版本需要小于2.0.0</el-tag></el-space>
        </Note>
    </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import passwordIcon from "@/assets/icon/password.svg";
import bugIcon from "@/assets/icon/bug.svg";
import fingerprintIcon from "@/assets/icon/fingerprint.svg";
import consoleIcon from "@/assets/icon/console.svg";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import { Document, ChromeFilled, DocumentCopy, Filter, Camera, Box, RefreshLeft, InfoFilled } from "@element-plus/icons-vue";
import { Copy, selectFileAndAssign } from "@/util";
import usePagination from "@/usePagination";
import { ReadFile } from "wailsjs/go/services/File";
import { ConnectAndExecute, FormatOutput } from "wailsjs/go/core/Tools";
import { ElMessage } from "element-plus";
import Note from "@/components/Note.vue";

const fscan = reactive({
    currentMode: 0,
    showType: "default",
    input: "",
    result: "",
    weakpass: [] as { type: string, ip: string, port: string, username: string, password: string, extend: string }[],
    virus: [] as { url: string, pocinfo: string, extend: string }[],
    fingerprint: [] as { url: string, fingerprint: string }[],
    tipsDialog: false,
});

const wip = usePagination<{ url: string, code: string, title: string, length: string, redirect: string }>(20)

const weekProtocol = ["ftp", "ssh", "telnet", "mysql", "oracle", "mssql", "postgres", "rdp", "mongodb", "redis", "smb"]
async function fscanParse() {
    if (!fscan.input) {
        ElMessage.warning("请输入文件路径")
        return
    }
    let file = await ReadFile(fscan.input);
    let result = await FormatOutput(file.Content)
    if (!result) return
    fscan.weakpass = []
    fscan.virus = []
    fscan.result = ""
    for (const [key, values] of Object.entries(result)) {
        if (values) {
            fscan.result += `[${key}]\n${values.join("\n")}\n\n`
        }
        // if (key == "NetInfo") {

        // }
        if (weekProtocol.includes(key.toLocaleLowerCase())) {
            const uniqueWeakpass = new Set();
            values.forEach(line => {
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
                const weakpassKey = `${protocol}-${ip}-${port}-${username}-${password}-${ext}`;
                uniqueWeakpass.add(weakpassKey);
            });
            Array.from(uniqueWeakpass).forEach((weakpassKey: string) => {
                const [protocol, ip, port, username, password, ext] = weakpassKey.split('-');
                fscan.weakpass.push({
                    type: protocol,
                    ip: ip,
                    port: port,
                    username: username,
                    password: password,
                    extend: ext
                });
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
            wip.ctrl.watchResultChange(wip.table)
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


async function connectAndExecute(row: any) {
    // Find the matching entry in fscan.weakpass by IP and port
    const matchedEntry = fscan.weakpass.find(entry => entry.ip === row.ip && entry.port === row.port);
    expandedRowKeys.value = [generateRowKey(row)]

    matchedEntry.extend = "正在连接请稍等..."
    const result = await ConnectAndExecute(row.type, row.ip, row.port, row.username, row.password);

    // If a matching entry is found, update its 'extend' field with the result
    if (matchedEntry) {
        matchedEntry.extend = result;
    }
}


const filter = ({
    hikvsion: ["2512", "600", "481", "480"],
    vm: ["ID_VC_Welcome", "ID_EESX_Welcome"],
    reset: () => {
        if (wip.table.filterTemp.length != 0) wip.table.result = wip.table.filterTemp
        wip.ctrl.watchResultChange(wip.table);
    },
    hikvisionFilter: () => {
        if (wip.table.filterTemp.length == 0) wip.table.filterTemp = wip.table.result
        wip.table.result = wip.table.filterTemp.filter(item => filter.hikvsion.includes(item.length))
        wip.ctrl.watchResultChange(wip.table)
    },
    vmFilter: () => {
        if (wip.table.filterTemp.length == 0) wip.table.filterTemp = wip.table.result
        wip.table.result = wip.table.filterTemp.filter(item => {
            for (const name of filter.vm) {
                if (item.title.includes(name)) return item
            }
        });
        wip.ctrl.watchResultChange(wip.table)
    },
    pocinfo(value: string, row: any): boolean {
        return row.pocinfo === value;
    },
    weakpass(value: string, row: any): boolean {
        return row.type.toLowerCase() == value;
    },
})

</script>

<style scoped></style>