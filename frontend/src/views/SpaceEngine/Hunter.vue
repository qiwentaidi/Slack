<script setup lang="ts">
import { ExportToXlsx } from '../../export'
import { reactive, ref } from 'vue';
import { ElNotification, ElMessage } from "element-plus";
import { Menu, Search, ChatLineRound, ArrowDown, ChromeFilled, CopyDocument, Share } from '@element-plus/icons-vue';
import { TableTabs, ApiSyntaxCheck, splitInt, SplitTextArea, validateIP, validateDomain, Copy } from '../../util'
import global from "../../global"
import {
    WebIconMd5,
    HunterSearch,
    HunterTips,
} from '../../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
const form = reactive({
    query: '',
    syntaxDialog: false,
    hunterImg: 'hunter_syntax.png',
    optionsTime: [
        {
            value: '0',
            label: '最近一个月',
        },
        {
            value: '1',
            label: '最近半年',
        },
        {
            value: '2',
            label: '最近一年',
        },
    ],
    defaultTime: '0',
    optionsServer: [
        {
            value: '3',
            label: '全部资产',
        },
        {
            value: '1',
            label: 'WEB服务资产',
        },
        {
            value: '2',
            label: '非WEB服务资产',
        },
    ],
    defaultSever: '3',
    deduplication: false,
    tips: '',
    loadAll: [] as EntryTips[],
    icondialog: false,
    hashURL: '',
    batchdialog: false,
    batchURL: '',
})

interface EntryTips {
    value: string
    assetNum: number
    tags: string[]
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
        HunterTips(queryString).then(result => {
            form.loadAll = []
            if (result.code == 200) {
                for (const item of result.data.app) {
                    form.loadAll.push({
                        value: item.name,
                        assetNum: item.asset_num,
                        tags: item.tags
                    })
                }
            }
        })
    },
    handleSelect: (item: Record<string, any>) => {
        form.query = `app.name="${item.value}"`
    }
})

const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as TableTabs[],
    addTab: (query: string) => {
        if (ApiSyntaxCheck(1, "", global.space.hunterkey, query) === false) {
            return
        }
        const newTabName = `${++table.tabIndex}`
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: [{}],
            total: 0,
            pageSize: 10,
            currentPage: 1,
        });
        loading.value = true
        HunterSearch(global.space.hunterkey, query, "10", "1", form.defaultTime, form.defaultSever, form.deduplication).then(result => {
            if (result.code !== 200) {
                switch (result.code) {
                    case 0:
                        ElNotification({
                            message: '请求超时',
                            type: "error",
                        });
                        loading.value = false
                        return
                    case 40205:
                        ElNotification({
                            message: result.message,
                            type: "info",
                        });
                    default:
                        ElNotification({
                            message: result.message,
                            type: "error",
                        });
                        loading.value = false
                        return
                }
            }
            form.tips = result.message + " 共查询到数据:" + result.data.total + "条," + result.data.rest_quota
            const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
            tab.content!.pop()
            if (result.data.arr == null) {
                ElNotification({
                    message: "暂未查询到相关数据",
                    type: "warning",
                });
                loading.value = false
                return
            }
            result.data.arr.forEach((item: any) => {
                tab.content?.push({
                    URL: item.url,
                    IP: item.ip,
                    Port: item.port,
                    Protocol: item.protocol,
                    Domain: item.domain,
                    Component: item.component,
                    Title: item.web_title,
                    Status: item.status_code,
                    ICP: item.company,
                    ISP: item.isp,
                    Position: item.country + "/" + item.province,
                    UpdateTime: item.updated_at,
                })
            });
            tab.total = result.data.total
            table.acvtiveNames = newTabName
            loading.value = false
        })
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
        HunterSearch(global.space.hunterkey, tab.title, val.toString(), "1", form.defaultTime, form.defaultSever, form.deduplication).then(result => {
            if (result.code !== 200) {
                if (result.code == 40205) {
                    ElNotification({
                        title: "提示",
                        message: result.message,
                        type: "info",
                    });
                } else {
                    ElNotification({
                        title: "提示",
                        message: result.message,
                        type: "error",
                    });
                    loading.value = false
                    return
                }
            }
            form.tips = result.message + " 共查询到数据:" + result.data.total + "条," + result.data.rest_quota
            tab.content = [{}]
            tab.content.pop()
            result.data.arr.forEach((item: any) => {
                tab.content?.push({
                    URL: item.url,
                    IP: item.ip,
                    Port: item.port,
                    Protocol: item.protocol,
                    Domain: item.domain,
                    Component: item.component,
                    Title: item.web_title,
                    Status: item.status_code,
                    ICP: item.company,
                    ISP: item.isp,
                    Position: item.country + "/" + item.province,
                    UpdateTime: item.updated_at,
                })
            });
            tab.total = result.data.total
            loading.value = false
        })
    },
    handleCurrentChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        loading.value = true
        HunterSearch(global.space.hunterkey, tab.title, tab.pageSize.toString(), val.toString(), form.defaultTime, form.defaultSever, form.deduplication).then(result => {
            if (result.code !== 200) {
                if (result.code == 40205) {
                    ElNotification({
                        title: "提示",
                        message: result.message,
                        type: "info",
                    });
                } else {
                    ElNotification({
                        title: "提示",
                        message: result.message,
                        type: "error",
                    });
                    loading.value = false
                    return
                }
            }
            form.tips = result.message + " 共查询到数据:" + result.data.total + "条," + result.data.rest_quota
            tab.content = [{}]
            tab.content.pop()
            result.data.arr.forEach((item: any) => {
                tab.content?.push({
                    URL: item.url,
                    IP: item.ip,
                    Port: item.port,
                    Protocol: item.protocol,
                    Domain: item.domain,
                    Component: item.component,
                    Title: item.web_title,
                    Status: item.status_code,
                    ICP: item.company,
                    ISP: item.isp,
                    Position: item.country + "/" + item.province,
                    UpdateTime: item.updated_at,
                })
            });
            tab.total = result.data.total
            loading.value = false
        })
    },
})
const loading = ref(false)

async function IconHashSearch() {
    let hash = await WebIconMd5(form.hashURL.trim())
    if (hash == "") {
        ElNotification({
            message: "目标不可达或者URL格式错误",
            type: "warning",
        });
        return
    }
    form.icondialog = false
    table.addTab(`web.icon=="${hash}"`)
}

function BatchSearch() {
    const lines = SplitTextArea(form.batchURL)
    var temp = ''
    for (const line of lines) {
        if (validateIP(line)) {
            temp += `ip="${line}"||`
        } else if (validateDomain(line)) {
            temp += `domain.suffix="${line}"||`
        }
    }
    if (temp == '') {
        ElNotification({
            message: "目标为空",
            type: "warning",
        });
        return
    }
    table.addTab(temp.slice(0, -2))
}


async function SaveData(mode: number) {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        if (mode == 0) {
            ExportToXlsx(["URL", "IP", "端口", "协议", "域名", "应用/组件", "标题", "状态码", "备案号", "运营商", "地理位置", "更新时间"], "asset", "hunter_asset", tab.content!)
        } else {
            ElNotification({
                title: "提示",
                message: "正在进行全数据导出，API每页最大查询限度100，请稍后。",
                type: "info",
            });
            let temp = [{}]
            temp.pop()
            let index = 0
            for (const num of splitInt(tab.total, 100)) {
                index += 1
                ElMessage("正在导出第" + index.toString() + "页");
                await HunterSearch(global.space.hunterkey, tab.title, "100", index.toString(), form.defaultTime, form.defaultSever, form.deduplication).then(result => {
                    result.data.arr.forEach((item: any) => {
                        temp.push({
                            URL: item.url,
                            IP: item.ip,
                            Port: item.port,
                            Protocol: item.protocol,
                            Domain: item.domain,
                            Component: item.component,
                            Title: item.web_title,
                            Status: item.status_code,
                            ICP: item.company,
                            ISP: item.isp,
                            Position: item.country + "/" + item.province,
                            UpdateTime: item.updated_at,
                        })
                    });
                })
            }
            ExportToXlsx(["URL", "IP", "端口", "协议", "域名", "应用/组件", "标题", "状态码", "备案号", "运营商", "地理位置", "更新时间"], "asset", "hunter_asset", temp)
            temp = []
        }
    }
}

// 0 当前页 1 100条
async function CopyURL(mode: number) {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        var temp = [];
        if (mode == 0) {
            for (const line of tab.content!) {
                temp.push((line as any)["URL"])
            }
        } else {
            await HunterSearch(global.space.hunterkey, tab.title, "100", "1", form.defaultTime, form.defaultSever, form.deduplication).then(result => {
                result.data.arr.forEach((item: any) => {
                    temp.push(item.url)
                })
            })
        }
        Copy(temp.join("\n"))
        temp = []
    }
}
</script>

<template>
    <el-form v-model="form" @submit.native.prevent="table.addTab(form.query)">
        <el-form-item label="查询条件">
            <div class="head">
                <el-autocomplete v-model="form.query" placeholder="Search..."
                    :fetch-suggestions="entry.querySearchAsync" @select="entry.handleSelect" :trigger-on-focus="false"
                    :debounce="1000" style="width: 100%;">
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
                            <span>({{ item.assetNum }} 条数据)</span>
                            <el-tag v-for="label in item.tags">{{ label }}</el-tag>
                        </el-space>
                    </template>
                </el-autocomplete>
                <el-dialog v-model="form.icondialog" title="输入目标favicon地址会自动计算并搜索相关资产" width="50%" center>
                    <el-input v-model="form.hashURL"></el-input>
                    <template #footer>
                        <span>
                            <el-button type="primary" @click="IconHashSearch">
                                搜索
                            </el-button>
                        </span>
                    </template>
                </el-dialog>
                <el-dialog v-model="form.batchdialog" title="批量查询: 请输入IP/网段/域名(MAX 5)" width="40%" center>
                    <el-input v-model="form.batchURL" type="textarea" rows="10"></el-input>
                    <template #footer>
                        <span>
                            <el-button type="primary" @click="BatchSearch">
                                搜索
                            </el-button>
                        </span>
                    </template>
                </el-dialog>
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
                            <el-dropdown-item :icon="CopyDocument" @click="CopyURL(0)"
                                divided>复制当前页URL</el-dropdown-item>
                            <el-dropdown-item :icon="CopyDocument" @click="CopyURL(1)">复制前100条URL</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
        </el-form-item>
    </el-form>
    <div class="my-header">
        <el-space>
            <el-select v-model="form.defaultTime" style="width: 120px;">
                <el-option v-for="item in form.optionsTime" :key="item.value" :label="item.label" :value="item.value"
                    style="text-align: center;" />
            </el-select>
            <el-select v-model="form.defaultSever" style="width: 160px;">
                <el-option v-for="item in form.optionsServer" :key="item.value" :label="item.label" :value="item.value"
                    style="text-align: center;" />
            </el-select>
            <el-checkbox v-model="form.deduplication" size="large">数据去重(需权益积分)</el-checkbox>
        </el-space>
        <el-link @click="form.syntaxDialog = true"><el-icon>
                <ChatLineRound />
            </el-icon>查询语法</el-link>
        <el-dialog v-model="form.syntaxDialog" title="查询语法参考" width="80%" center>
            <el-scrollbar height="400px">
                <el-image :src="form.hunterImg"></el-image>
            </el-scrollbar>
        </el-dialog>
    </div>
    <el-tabs v-model="table.acvtiveNames" v-loading="loading" type="card" style="margin-top: 10px;" closable
        @tab-remove="table.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column type="index" fixed label="#" width="60px" />
                <el-table-column prop="URL" fixed label="URL" width="200" :show-overflow-tooltip="true">
                    <template #default="scope">
                        <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)" v-show="scope.row.URL != ''" />
                        {{ scope.row.URL }}
                    </template>

                </el-table-column>
                <el-table-column prop="IP" fixed label="IP" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Port" fixed label="端口/服务" width="120">
                    <template #default="scope">
                        {{ scope.row.Port }}
                        <el-tag type="info">{{ scope.row.Protocol }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="Domain" label="域名" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Component" label="应用/组件" width="210xp">
                    <template #default="scope">
                        <el-space>
                            <el-tag v-if="Array.isArray(scope.row.Component) && scope.row.Component.length > 0">{{
        scope.row.Component[0].name + scope.row.Component[0].version }}</el-tag>
                            <el-popover placement="bottom" :width="350" trigger="hover">
                                <template #reference>
                                    <el-button round size="small"
                                        v-if="Array.isArray(scope.row.Component) && scope.row.Component.length > 0">共{{
        scope.row.Component.length }}条</el-button>
                                </template>
                                <template #default>
                                    <div style="display: flex; flex-direction: column;">
                                        <el-tag v-for="component in scope.row.Component">{{ component.name +
        component.version }}</el-tag>
                                    </div>
                                </template>
                            </el-popover>
                        </el-space>
                    </template>
                </el-table-column>
                <el-table-column prop="Title" label="标题" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Status" label="状态码" :show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="ISP" label="运营商" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Position" label="地理位置" width="120" :show-overflow-tooltip="true" />
                <el-table-column prop="UpdateTime" label="更新时间" width="150" :show-overflow-tooltip="true" />
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination background :page-size="10" :page-sizes="[10, 50, 100]" layout="sizes, prev, pager, next"
                    @size-change="table.handleSizeChange" @current-change="table.handleCurrentChange"
                    :total="item.total" />
            </div>
        </el-tab-pane>
        <el-image class="center" src="/loading.gif" alt="loading" v-else></el-image>
    </el-tabs>
</template>

<style>
.el-tabs__item {
    position: relative;
    display: inline-block;
    max-width: 300px;
    margin-bottom: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.el-tabs__item .el-icon {
        position: absolute !important;
        top: 13px !important;
        right: 3px !important;
    }

.el-tabs__nav {
    line-height: 255%;
}
</style>