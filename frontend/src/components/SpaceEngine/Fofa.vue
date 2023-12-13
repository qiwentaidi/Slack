<script lang="ts" setup>
import { reactive } from 'vue';
import { Menu, Search, ChatLineRound, ChromeFilled, ArrowDown } from '@element-plus/icons-vue';
import { SplitTextArea, validateIP, validateDomain, ExportToXlsx, splitInt } from '../../util'
import {
    FofaTips,
    FofaSearch,
    IconHash,
    TestTarget
} from '../../../wailsjs/go/main/App'
import { ElMessage } from 'element-plus';
import global from "../Global.vue"
const from = reactive({
    query: '',
    fraud: false,
    cert: false,
    tips: '',
    dialog: false,
    fofaImg: [
        'fofa1.png',
        'fofa2.png',
        'fofa3.png',
        'fofa4.png',
        'fofa5.png',
        'fofa6.png',
    ],
    loadAll: [] as LinkItem[],
    icondialog: false,
    batchdialog: false,
    hashURL: '',
    batchURL: '',
})

// 输入框Tips提示

interface LinkItem {
    value: string
    link: string
}
let timeout: ReturnType<typeof setTimeout>
const entry = reactive({
    querySearchAsync: (queryString: string, cb: (arg: any) => void) => {
        if (queryString.length > 1) {
            entry.getTips()
            clearTimeout(timeout)
            timeout = setTimeout(() => {
                cb(from.loadAll)
            }, 3000 * Math.random())
        } else {
            cb([])
        }
    },
    getTips: function () {
        FofaTips(from.query).then(result => {
            from.loadAll = []
            if (result.code == 0) {
                for (const item of result.data) {
                    from.loadAll.push({
                        value: item.name,
                        link: item.company
                    })
                }
            }
        })
    },
    handleInput: (value: string) => {
        entry.querySearchAsync(value, suggestions => {
            from.loadAll = suggestions;
        });
    },
    handleSelect: (item: Record<string, any>) => {
        from.query = `app="${item.value}"`
    }
})

interface tableTabs {
    title: string,
    name: string,
    content: null | [{}],
    total: number,
    pageSize: number,
    currentPage: number,
}

const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as tableTabs[],
    addTab: (query: string) => {
        let ff = new FOFA
        if (ff.CheckApi() === false) {
            return
        }
        if (ff.CheckSyntax(query) === false) {
            return
        }
        const newTabName = `${++table.tabIndex}`
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: [{}],
            total: 0,
            pageSize: 100,
            currentPage: 1,
        });
        FofaSearch(query, "100", "1", global.space.fofaemail, global.space.fofakey, from.fraud, from.cert).then(result => {
            from.tips = result.Message + " 共查询到数据:" + result.Total + "条"
            if (result.Status == false) {
                return
            }
            const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
            tab.content = result.Results;
            tab.total = result.Total
            table.acvtiveNames = newTabName
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
    openLink: (row: any) => {
        window.open(row.URL, '_blank');
    },
    handleSizeChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        FofaSearch(from.query, val.toString(), "1", global.space.fofaemail, global.space.fofakey, from.fraud, from.cert).then(result => {
            from.tips = result.Message + " 共查询到数据:" + result.Total + "条"
            if (result.Status == false) {
                return
            }
            tab.content = result.Results;
            tab.total = result.Total
        })
    },
    handleCurrentChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        FofaSearch(from.query, tab.pageSize.toString(), val.toString(), global.space.fofaemail, global.space.fofakey, from.fraud, from.cert).then(result => {
            from.tips = result.Message + " 共查询到数据:" + result.Total + "条"
            if (result.Status == false) {
                return
            }
            tab.content = result.Results;
            tab.total = result.Total
        })
    }
})

let RegCompliance = new RegExp("(\\w+)[!,=]{1,3}\"([^\"]+)\"");

class FOFA {
    public CheckApi() {
        if (global.space.fofaemail == "" || global.space.fofakey == "") {
            ElMessage({
                showClose: true,
                message: "请在设置处填写FOFA Email和Key",
                type: "error",
            });
            return false
        }
    }
    public CheckSyntax(query: string) {
        if (RegCompliance.test(query) === false) {
            ElMessage({
                showClose: true,
                message: "请输入正确的查询语法",
                type: "warning",
            });
            return false
        }
    }
}
async function HashSearch() {
    if (await TestTarget(from.hashURL) == false) {
        ElMessage({
            showClose: true,
            message: "目标不可达",
            type: "warning",
        });
        return
    }
    let hash = IconHash(from.hashURL)
    let query = `icon_hash="${(await hash).toString()}"`
    table.addTab(query)
}

function BatchSearch() {
    const lines = SplitTextArea(from.batchURL)
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
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "国家", "省份", "城市", "备案号"], table.acvtiveNames, "fofa_asset", tab.content!)
        }else {
            let temp = [{}]
            splitInt(tab.total,10000).forEach(async (val, idx, array) => {
                await FofaSearch(table.acvtiveNames, val.toString(), (idx+1).toString(), global.space.fofaemail, global.space.fofakey, from.fraud, from.cert).then(result => {
                    if (result.Status == false) {
                        return
                    }
                    temp.push(result.Results)
                })
            })
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "国家", "省份", "城市", "备案号"], table.acvtiveNames, "fofa_asset", temp)
            temp = []
        }
    }
}

</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <template #label>
                <el-tooltip placement="right">
                    <template #content>输入字符>=2会出现提示</template>
                    <span>查询条件</span>
                </el-tooltip>
            </template>
            <div class="head">
                <el-autocomplete v-model="from.query" placeholder="Search..." :fetch-suggestions="entry.querySearchAsync"
                    @select="entry.handleSelect" @input="entry.handleInput" :trigger-on-focus="false" style="width: 100%;">
                    <template #append>
                        <el-dropdown>
                            <el-button :icon="Menu" />
                            <template #dropdown>
                                <el-dropdown-menu :hide-on-click="true">
                                    <el-dropdown-item @click="from.icondialog = true">icon搜索</el-dropdown-item>
                                    <el-dropdown-item @click="from.batchdialog = true">批量查询</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </template>
                </el-autocomplete>
                <el-button type="primary" :icon="Search" @click="table.addTab(from.query)"
                    style="margin-left: 10px; margin-right: 10px;">查询</el-button>
                <el-dropdown>
                    <el-button color="#A29EDE">
                        数据导出<el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="SaveData(0)">导出当前查询页数据</el-dropdown-item>
                            <el-dropdown-item @click="SaveData(1)">导出全部数据</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <el-dialog v-model="from.icondialog" title="输入目标favicon地址会自动计算并搜索相关资产" width="50%" center>
                <el-input v-model="from.hashURL"></el-input>
                <template #footer>
                    <span>
                        <el-button type="primary" @click="HashSearch">
                            搜索
                        </el-button>
                    </span>
                </template>
            </el-dialog>
            <el-dialog v-model="from.batchdialog" title="批量查询: 请输入IP/网段/域名(MAX 100)" width="40%" center>
                <el-input v-model="from.batchURL" type="textarea" rows="10"></el-input>
                <template #footer>
                    <span>
                        <el-button type="primary" @click="BatchSearch">
                            搜索
                        </el-button>
                    </span>
                </template>
            </el-dialog>
        </el-form-item>
    </el-form>
    <div class="nkmode">
        <div>
            <el-checkbox size="large" v-model="from.fraud">排除干扰(专业版)</el-checkbox>
            <el-checkbox size="large" v-model="from.cert">证书(个人版)</el-checkbox>
        </div>
        <el-link @click="from.dialog = true"><el-icon>
                <ChatLineRound />
            </el-icon>查询语法</el-link>
        <el-dialog v-model="from.dialog" title="查询语法参考" width="80%" center>
            <div class="demo-image__lazy">
                <el-image v-for="url in from.fofaImg" :key="url" :src="url" lazy></el-image>
            </div>
        </el-dialog>
    </div>
    <el-tabs v-model="table.acvtiveNames" type="card" style="margin-top: 10px;" closable @tab-remove="table.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column prop="URL" label="URL" width="200" show-overflow-tooltip="true" />
                <el-table-column prop="Title" label="标题" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="IP" label="IP" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口" width="70" show-overflow-tooltip="true" />
                <el-table-column prop="Domain" label="域名" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Protocol" label="协议" width="100" show-overflow-tooltip="true" />
                <el-table-column prop="Country" label="国家" show-overflow-tooltip="true" />
                <el-table-column prop="Region" label="省份" show-overflow-tooltip="true" />
                <el-table-column prop="City" label="城市" show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" show-overflow-tooltip="true" />
                <el-table-column fixed="right" label="操作" width="55px">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="left">
                            <el-button link :icon="ChromeFilled" @click.prevent="table.openLink(scope.row)">
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="nkmode" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ from.tips }}</span>
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
</style>

