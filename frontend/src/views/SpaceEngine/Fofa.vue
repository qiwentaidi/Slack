<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Menu, Search, ChatLineRound, ChromeFilled, ArrowDown, CopyDocument, Share } from '@element-plus/icons-vue';
import { SplitTextArea, validateIP, validateDomain, splitInt, TableTabs, ApiSyntaxCheck, Copy } from '../../util'
import { ExportToXlsx } from '../../export'
import {
    FofaTips,
    FofaSearch,
    IconHash,
} from '../../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
import { ElMessage, ElNotification } from 'element-plus';
import global from "../../global"
const form = reactive({
    query: '',
    fraud: false,
    cert: false,
    tips: '',
    fofaImg: [
        'fofa1.png',
        'fofa2.png',
        'fofa3.png',
        'fofa4.png',
        'fofa5.png',
        'fofa6.png',
    ],
    loadAll: [] as LinkItem[],
    syntaxdialog: false,
    icondialog: false,
    batchdialog: false,
    socksdialog: false,
    hashURL: '',
    batchURL: '',
    socksLogger: '',
    socksMax: 1,
    socksNum: 1,
    percentage: 0,
})
const loading = ref(false)

// 输入框Tips提示

interface LinkItem {
    value: string
    link: string
}
let timeout: ReturnType<typeof setTimeout>
const entry = reactive({
    querySearchAsync: (queryString: string, cb: (arg: any) => void) => {
        if (!queryString.includes("=")) {
            entry.getTips(queryString)
            clearTimeout(timeout)
            timeout = setTimeout(() => {
                cb(form.loadAll)
            }, 2000 * Math.random())
        }
    },
    getTips: function (queryString: string) {
        FofaTips(queryString).then(result => {
            form.loadAll = []
            if (result.code == 0) {
                for (const item of result.data) {
                    form.loadAll.push({
                        value: item.name,
                        link: item.company
                    })
                }
            }
        })
    },
    handleSelect: (item: Record<string, any>) => {
        form.query = `app="${item.value}"`
    }
})


const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as TableTabs[],
    addTab: async (query: string) => {
        if (ApiSyntaxCheck(0, global.space.fofaemail, global.space.fofakey, query) === false) {
            return
        }
        const newTabName = `${++table.tabIndex}`
        loading.value = true
        let result = await FofaSearch(query, "100", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
        if (!result.Status) {
            ElMessage({
                showClose: true,
                message: result.Message,
                type: "warning",
            });
            return
        }
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: [{}],
            total: 0,
            pageSize: 100,
            currentPage: 1,
        });
        form.tips = result.Message + " 共查询到数据:" + result.Total + "条"
        const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
        tab.content = result.Results;
        tab.total = result.Total
        table.acvtiveNames = newTabName
        loading.value = false
    },
    removeTab: (targetName: string) => {
        const tabs = table.editableTabs
        let activeName = table.acvtiveNames
        if (activeName === targetName) {
            tabs.forEach((tab, index) => {
                if (tab.name === targetName) {
                    tab.content = null // 清理内存
                    const nextTab = tabs[index + 1] || tabs[index - 1]
                    if (nextTab) {
                        activeName = nextTab.name
                    }
                }
            })
        }
        table.acvtiveNames = activeName
        table.editableTabs = tabs.filter((tab) => tab.name !== targetName)
    },
    handleSizeChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        loading.value = true
        FofaSearch(form.query, val.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then(result => {
            form.tips = result.Message + " 共查询到数据:" + result.Total + "条"
            if (result.Status == false) {
                return
            }
            tab.content = result.Results;
            tab.total = result.Total
        })
        loading.value = false
    },
    handleCurrentChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        loading.value = true
        FofaSearch(form.query, tab.pageSize.toString(), val.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then(result => {
            form.tips = result.Message + " 共查询到数据:" + result.Total + "条"
            if (result.Status == false) {
                return
            }
            tab.content = result.Results;
            tab.total = result.Total
        })
        loading.value = false
    }
})


async function IconHashSearch() {
    let hash = await IconHash(form.hashURL)
    if (hash == "") {
        ElMessage({
            showClose: true,
            message: "目标不可达",
            type: "warning",
        });
        return
    }
    table.addTab(`icon_hash="${hash}"`)
}

function BatchSearch() {
    const lines = SplitTextArea(form.batchURL)
    var temp = ''
    for (const line of lines) {
        if (validateIP(line) === true) {
            temp += `ip="${line}"||`
        } else if (validateDomain(line) === true) {
            temp += `domain="${line}"||`
        }
    }
    if (temp == '') {
        ElMessage({
            showClose: true,
            message: "目标为空",
            type: "warning",
        });
        return
    }
    table.addTab(temp.slice(0, -2))
}

// mode = 0 导出当前查询数据 
async function SaveData(mode: number) {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        if (mode == 0) {
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "国家", "省份", "城市", "备案号"], "asset", "fofa_asset", tab.content!)
        } else {
            let temp = [{}]
            temp.pop()
            ElNotification("正在进行全数据导出，API每页最大查询限度10000，请稍后。");
            let index = 0
            for (const num of splitInt(tab.total, 10000)) {
                index += 1
                ElMessage("正在导出第" + index.toString() + "页");
                await FofaSearch(tab.title, num.toString(), index.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then(result => {
                    if (result.Status == false) {
                        return
                    }
                    temp.push(...result.Results)
                })
            }
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "国家", "省份", "城市", "备案号"], "asset", "fofa_asset", temp)
            temp = []
        }
    }
}

function getColumnData(prop: string): any[] {
    const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
    let protocols = new Set(tab.content!.map((item: any) => item[prop]));
    let newArray = Array.from(protocols).map(protocol => ({ text: protocol, value: protocol }))
    return newArray
}

function filterHandlerProtocol(value: string, row: any): boolean {
    return row.Protocol === value;
}

function filterHandlerTitle(value: string, row: any): boolean {
    return row.Title === value;
}

async function CopyURL() {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        var temp = [];
        for (const line of tab.content!) {
            temp.push((line as any)["URL"])
        }
        Copy(temp.join("\n"))
        temp = []
    }
}
</script>

<template>
    <el-form :model="form" @submit.native.prevent="table.addTab(form.query)">
        <el-form-item label="查询条件">
            <div class="head">
                <el-autocomplete v-model="form.query" placeholder="Search..." :fetch-suggestions="entry.querySearchAsync"
                    @select="entry.handleSelect" :trigger-on-focus="false" debounce="1000" style="width: 100%;">
                    <template #append>
                        <el-dropdown>
                            <el-button :icon="Menu" />
                            <template #dropdown>
                                <el-dropdown-menu :hide-on-click="true">
                                    <el-dropdown-item @click="form.icondialog = true">icon搜索</el-dropdown-item>
                                    <el-dropdown-item @click="form.batchdialog = true">批量查询</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </template>
                    <template #default="{ item }">
                        <el-space>
                            <span>{{ item.value }}</span>
                            <el-tag>{{ item.link }}</el-tag>
                        </el-space>
                    </template>
                </el-autocomplete>
                <el-button type="primary" :icon="Search" @click="table.addTab(form.query)"
                    style="margin-left: 10px; margin-right: 10px;">查询</el-button>
                <el-dropdown>
                    <el-button color="#A29EDE">
                        数据导出/复制<el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item :icon="Share" @click="SaveData(0)">导出当前查询页数据</el-dropdown-item>
                            <el-dropdown-item :icon="Share" @click="SaveData(1)">导出全部数据</el-dropdown-item>
                            <el-dropdown-item :icon="CopyDocument" @click="CopyURL" divided>复制当前页URL</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <!-- icon搜索 -->
            <el-dialog v-model="form.icondialog" title="输入目标favicon地址会自动计算并搜索相关资产" width="50%" center>
                <el-input v-model="form.hashURL"></el-input>
                <template #footer>
                    <el-button type="primary" @click="IconHashSearch">搜索</el-button>
                </template>
            </el-dialog>
            <!-- 批量查询 -->
            <el-dialog v-model="form.batchdialog" title="批量查询: 请输入IP/网段/域名(MAX 100)" width="40%" center>
                <el-input v-model="form.batchURL" type="textarea" rows="10"></el-input>
                <template #footer>
                    <el-button type="primary" @click="BatchSearch">搜索</el-button>
                </template>
            </el-dialog>
        </el-form-item>
    </el-form>
    <div class="my-header">
        <div>
            <el-checkbox size="large" v-model="form.fraud">排除干扰(专业版)</el-checkbox>
            <el-checkbox size="large" v-model="form.cert">证书(个人版)</el-checkbox>
        </div>
        <el-link @click="form.syntaxdialog = true"><el-icon>
                <ChatLineRound />
            </el-icon>查询语法</el-link>
        <el-dialog v-model="form.syntaxdialog" title="查询语法参考" width="80%" center>
            <div class="demo-image__lazy">
                <el-image v-for="url in form.fofaImg" :key="url" :src="url" lazy></el-image>
            </div>
        </el-dialog>
    </div>
    <el-tabs v-model="table.acvtiveNames" type="card" style="margin-top: 10px;" closable @tab-remove="table.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column type="index" label="#" width="60px" />
                <el-table-column prop="URL" label="URL" width="200" show-overflow-tooltip="true">
                    <template #default="scope">
                        <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)">
                        </el-button>
                        {{ scope.row.URL }}
                    </template>
                </el-table-column>
                <el-table-column prop="Title" label="标题" :filters='getColumnData("Title")'
                    :filter-method="filterHandlerTitle" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="IP" label="IP" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口" width="100"
                    :sort-method="(a: any, b: any) => { return a.Port - b.Port }" sortable show-overflow-tooltip="true" />
                <el-table-column prop="Domain" label="域名" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Protocol" label="协议" :filters='getColumnData("Protocol")'
                    :filter-method="filterHandlerProtocol" width="100" show-overflow-tooltip="true" />
                <el-table-column prop="Country" label="国家" show-overflow-tooltip="true" />
                <el-table-column prop="Region" label="省份" show-overflow-tooltip="true" />
                <el-table-column prop="City" label="城市" show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" show-overflow-tooltip="true" />
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination :page-size="100" :page-sizes="[100, 500, 1000]" layout="sizes, prev, pager, next"
                    @size-change="table.handleSizeChange" @current-change="table.handleCurrentChange" :total="item.total" />
            </div>
        </el-tab-pane>
        <el-image class="center" src="/loading.gif" alt="loading" v-else></el-image>
    </el-tabs>
</template>

<style scoped>
.demo-image__lazy {
    height: 400px;
    overflow-y: auto;
}

.demo-image__lazy .el-image {
    display: block;
    min-height: 200px;
    margin-bottom: 10px;
}

.demo-image__lazy .el-image:last-child {
    margin-bottom: 0;
}

.el-tabs__item {
    position: relative;
    display: inline-block;
    max-width: 300px;
    margin-bottom: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    .el-tooltip {
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .el-icon-close {
        position: absolute !important;
        top: 13px !important;
        right: 3px !important;
    }
}

.el-tabs__nav {
    line-height: 230%;
}
</style>