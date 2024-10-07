<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Search, ChromeFilled, CopyDocument, CollectionTag, Delete, Document, PictureRounded, Star, Collection } from '@element-plus/icons-vue';
import { SplitTextArea, splitInt, Copy, CsegmentIpv4 } from '@/util'
import { TableTabs, Results } from "@/interface"
import { ExportToXlsx } from '@/export'
import { FofaTips, FofaSearch, IconHash, Callgologger } from 'wailsjs/go/main/App'
import { BrowserOpenURL } from 'wailsjs/runtime'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus';
import global from "@/global"
import { InsertFavGrammarFiled, RemoveFavGrammarFiled, SelectAllSyntax } from 'wailsjs/go/main/Database';
import exportIcon from '@/assets/icon/doucment-export.svg'
import csegmentIcon from '@/assets/icon/csegment.svg'
import { validateSingleDomain, validateSingleIP } from '@/stores/validate';
import { main } from 'wailsjs/go/models';
import { fofaOptions } from '@/stores/options';

const form = reactive({
    query: '',
    fraud: false,
    cert: false,
    tips: '',
    syntaxData: [] as main.Syntax[],
})

// ref得单独校验
const ruleFormRef = ref<FormInstance>()

const syntax = ({
    keywordActive: "基础类",
    rowClick: function (row: any, column: any, event: Event) {
        if (!form.query) {
            form.query = row.key
            return
        }
        form.query += " && " + row.key
    },
    rowClick2: function (row: any, column: any, event: Event) {
        if (!form.query) {
            form.query = row.Content
            return
        }
        form.query += " && " + row.Content
    },
    starDialog: ref(false),
    ruleForm: reactive<main.Syntax>({
        Name: '',
        Content: '',
    }),
    createStarDialog: () => {
        syntax.starDialog.value = true
        syntax.ruleForm.Name = ""
        syntax.ruleForm.Content = form.query
    },
    submitStar: async (formEl: FormInstance | undefined) => {
        if (!formEl) return
        let result = await formEl.validate()
        if (!result) return
        InsertFavGrammarFiled("fofa", syntax.ruleForm.Name!, syntax.ruleForm.Content!).then((r: Boolean) => {
            if (r) {
                ElMessage.success('添加语法成功')
            } else {
                ElMessage.error('添加语法失败')
            }
            syntax.starDialog.value = false
        })
    },
    deleteStar: (name: string, content: string) => {
        RemoveFavGrammarFiled("fofa", name, content).then((r: Boolean) => {
            if (r) {
                ElMessage.success('删除语法成功,重新打开刷新')
            } else {
                ElMessage.error('删除语法失败')
            }
        })
    },
    searchStarSyntax: async () => {
        form.syntaxData = await SelectAllSyntax("fofa")
    },
})

// 输入框Tips提示

const entry = ({
    querySearchAsync: async (queryString: string, cb: any) => {
        if (queryString.includes("=") || !queryString) {
            cb([]);
            return
        }
        let tips = await entry.getTips(queryString)
        cb(tips);
    },
    getTips: async function (queryString: string) {
        let tips = <{text: string, value: string}[]>[]
        let result: any = await FofaTips(queryString)
        if (result.code == 0) {
            for (const item of result.data) {
                tips.push({
                    value: item.name,
                    text: item.company
                })
            }
        }
        return tips
    },
    handleSelect: (item: Record<string, any>) => {
        form.query = `app="${item.value}"`
    },
})


const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as TableTabs[],
    loading: false,
})

const tableCtrl = ({
    addTab: async (query: string) => {
        const newTabName = `${++table.tabIndex}`
        table.loading = true
        let result = await FofaSearch(query, "100", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
        if (result.Error) {
            ElMessage.warning(result.Message);
            table.loading = false
            return
        }
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: [],
            total: 0,
            pageSize: 100,
            currentPage: 1,
        });
        table.acvtiveNames = newTabName
        form.tips = result.Message!
        const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
        tab.content = result.Results!;
        tab.total = result.Size!
        table.loading = false
    },
    removeTab: (targetName: string) => {
        const tabs = table.editableTabs
        let activeName = table.acvtiveNames
        if (activeName === targetName) {
            tabs.forEach((tab, index) => {
                if (tab.name === targetName) {
                    tab.content = [] // 清理内存
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
    handleSizeChange: async (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        table.loading = true
        let result = await FofaSearch(tab.title, val.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
        form.tips = result.Message!
        table.loading = false
        if (result.Error) {
            return
        }
        tab.content = result.Results!;
        tab.total = result.Size!
    },
    handleCurrentChange: async (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        table.loading = true
        let result = await FofaSearch(tab.title, tab.pageSize.toString(), val.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
        form.tips = result.Message!
        table.loading = false
        if (result.Error) {
            return
        }
        tab.content = result.Results!;
        tab.total = result.Size!
    }
})


function IconHashSearch() {
    ElMessageBox.prompt('输入目标Favicon地址会自动计算并搜索相关资产', 'ICON搜索', {
        confirmButtonText: '检索',
        inputPattern: /^(https?:\/\/)?((([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,})|localhost|(\d{1,3}\.){3}\d{1,3})(:\d+)?(\/[^\s]*)?$/,
        inputErrorMessage: 'Invalid URL',
        showCancelButton: false,
    })
        .then(async ({ value }) => {
            let hash = await IconHash(value.trim())
            if (!hash) {
                ElMessage("目标错误或不可达");
                return
            }
            tableCtrl.addTab(`icon_hash="${hash}"`)
        })
}

function BatchSearch() {
    ElMessageBox.prompt('请输入IP/网段/域名(MAX 100)', '批量查询', {
        confirmButtonText: '检索',
        inputType: 'textarea',
        showCancelButton: false,
    })
        .then(async ({ value }) => {
            const lines = SplitTextArea(value);
            const temp = [];

            for (const line of lines) {
                if (validateSingleIP(line)) {
                    temp.push(`ip="${line}"`);
                } else if (validateSingleDomain(line)) {
                    temp.push(`domain="${line}"`);
                }
            }

            if (temp.length === 0) {
                ElMessage.warning("目标为空");
                return;
            }

            tableCtrl.addTab(temp.join('||'));
        })
}

// mode = 0 导出当前查询数据 
async function SaveData(mode: number) {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        if (mode == 0) {
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "国家", "省份", "城市", "备案号"], "asset", "fofa_asset", tab.content!)
        } else {
            let temp = [] as Results[]
            ElMessage("正在进行全数据导出，API每页最大查询限度10000，请稍后...");
            let index = 0
            for (const num of splitInt(tab.total, 10000)) {
                index += 1
                ElMessage("正在导出第" + index.toString() + "页");
                let result = await FofaSearch(tab.title, num.toString(), index.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
                if (result.Error) {
                    Callgologger("error",`[fofa] ${tab.title} page ${index} search error: ` + result.Error);
                    ElMessage.error("查询异常，已退出导出")
                    continue
                }
                temp.push(...result.Results!)
            }
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "组件", "国家", "省份", "城市", "备案号"], "asset", "fofa_asset", temp)
            temp = []
        }
    }
}

function getColumnData(prop: string): any[] {
    const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
    if (tab.content == null || tab.content == undefined) {
        return []
    }
    let protocols = new Set(tab.content.map((item: any) => item[prop]));
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

async function CopyDeduplicationURL() {
    let result = await FofaSearch(form.query, "10000", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
    if (result.Error) {
        ElMessage.warning("查询失败")
        return
    }
    var temp = [];
    let content = result.Results
    let seen = new Set(); // 用于存储已处理过的 IP:Port 组合
    if (Array.isArray(content)) {
        for (const line of content) {
            let ip = (line as any)["IP"];
            let port = (line as any)["Port"];
            let ipPort = `${ip}:${port}`;
            if (!seen.has(ipPort)) {
                seen.add(ipPort); // 将新的 IP:Port 组合加入到 set 中
                temp.push(line["URL"]); // 添加 URL 到 temp 数组中
            }
        }
    }
    Copy(temp.join("\n"))
    temp = []
}

function formatProduct(raw: string): string[] {
    return !raw ? [] : raw.split(",")
}

function searchCsegmentIpv4(ip: string) {
    let ipv4 = CsegmentIpv4(ip)
    tableCtrl.addTab(`ip="${ipv4}"`)
}
</script>

<template>
    <el-form :model="form" @keydown.enter.native.prevent="tableCtrl.addTab(form.query)">
        <el-form-item>
            <el-autocomplete v-model="form.query" placeholder="Search..." :fetch-suggestions="entry.querySearchAsync"
                @select="entry.handleSelect" :debounce="1000" style="width: 100%;">
                <template #prepend>
                    查询条件
                </template>
                <template #suffix>
                    <el-space :size="2">
                        <el-popover placement="bottom-end" :width="800" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="语法检索" placement="bottom">
                                        <el-button :icon="CollectionTag" link />
                                    </el-tooltip>
                                </div>
                            </template>
                            <el-tabs v-model="syntax.keywordActive">
                                <el-tab-pane v-for="item in fofaOptions.Syntax" :name="item.title"
                                    :label="item.title">
                                    <el-table :data="item.data" stripe class="keyword-search"
                                        @row-click="syntax.rowClick">
                                        <el-table-column label="例句" width="300" property="key" />
                                        <el-table-column label="用途说明" width="350" property="description">
                                            <template #default="scope">
                                                {{ scope.row.description }}<el-tag v-if="scope.row.level"
                                                    style="margin-left: 5px;">{{ scope.row.level }}</el-tag>
                                            </template>
                                        </el-table-column>
                                        <el-table-column label="=" property="filed1" />
                                        <el-table-column label="!=" property="filed2" />
                                        <el-table-column label="*=" property="filed3" />
                                    </el-table>
                                </el-tab-pane>
                                <el-tab-pane label="连接符">
                                    <el-table :data="fofaOptions.Advanced" stripe class="keyword-search">
                                        <el-table-column label="逻辑连接符" width="100px" property="key" />
                                        <el-table-column label="具体含义" property="description">
                                            <template #default="scope">
                                                {{ scope.row.description }}<el-tag v-if="scope.row.level"
                                                    style="margin-left: 5px;">{{ scope.row.level }}</el-tag>
                                            </template>
                                        </el-table-column>
                                    </el-table>
                                </el-tab-pane>
                            </el-tabs>
                        </el-popover>
                        <el-tooltip content="使用网页图标搜索" placement="bottom">
                            <el-button :icon="PictureRounded" link @click="IconHashSearch()" />
                        </el-tooltip>
                        <el-tooltip content="IP/域名批量搜索" placement="bottom">
                            <el-button :icon="Document" link @click="BatchSearch()" />
                        </el-tooltip>
                    </el-space>
                    <el-divider direction="vertical" />
                    <el-space :size="2">
                        <el-tooltip content="清空语法" placement="bottom">
                            <el-button :icon="Delete" link @click="form.query = ''" />
                        </el-tooltip>
                        <el-tooltip content="收藏语法" placement="bottom">
                            <el-button :icon="Star" link @click="syntax.createStarDialog" />
                        </el-tooltip>
                        <el-tooltip content="复制语法" placement="bottom">
                            <el-button :icon="CopyDocument" link @click="Copy(form.query)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                    </el-space>
                    <el-button link :icon="Search" @click="tableCtrl.addTab(form.query)"
                        style="height: 40px;">查询</el-button>
                </template>
                <template #append>
                    <el-space :size="25">
                        <el-popover placement="bottom-end" :width="550" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="我收藏的语法" placement="left">
                                        <el-button :icon="Collection" @click="syntax.searchStarSyntax" />
                                    </el-tooltip>
                                </div>
                            </template>
                            <el-table :data="form.syntaxData" @row-click="syntax.rowClick2"
                                class="keyword-search">
                                <el-table-column width="150" prop="Name" label="语法名称" />
                                <el-table-column prop="Content" label="语法内容" />
                                <el-table-column label="操作" width="100">
                                    <template #default="scope">
                                        <el-button link
                                            @click="syntax.deleteStar(scope.row.Name, scope.row.Content)">删除
                                        </el-button>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-popover>
                    </el-space>
                </template>
                <template #default="{ item }">
                    <div>
                        <span>{{ item.value }}</span>
                        <el-divider direction="vertical" />
                        <span style="color: #559FF8;">{{ item.text }}</span>
                    </div>
                </template>
            </el-autocomplete>
        </el-form-item>
        <el-form-item>
            <div>
                <el-checkbox v-model="form.fraud">排除干扰(专业版)</el-checkbox>
                <el-checkbox v-model="form.cert">证书(个人版)</el-checkbox>
            </div>
            <div style="flex: 1;"></div>
            <el-dropdown>
                <el-button text bg>
                    更多功能<el-icon class="el-icon--right">
                        <ArrowDown />
                    </el-icon>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item :icon="exportIcon" @click="SaveData(0)">导出当前查询页数据</el-dropdown-item>
                        <el-dropdown-item :icon="exportIcon" @click="SaveData(1)">导出全部数据</el-dropdown-item>
                        <el-dropdown-item :icon="CopyDocument" @click="CopyURL" divided>复制当前页URL</el-dropdown-item>
                        <el-dropdown-item :icon="CopyDocument"
                            @click="CopyDeduplicationURL">去重复制前1w条URL</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </el-form-item>
    </el-form>
    <el-tabs class="editor-tabs" v-model="table.acvtiveNames" v-loading="table.loading" type="card" closable
        @tab-remove="tableCtrl.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%; height: calc(100vh - 280px);">
                <el-table-column type="index" fixed label="#" width="60px" />
                <el-table-column prop="URL" fixed label="URL" width="300" :show-overflow-tooltip="true" />
                <el-table-column prop="Title" label="标题" :filters='getColumnData("Title")'
                    :filter-method="filterHandlerTitle" width="200" :show-overflow-tooltip="true">
                    <template #filter-icon>
                        <Filter />
                    </template>
                </el-table-column>
                <el-table-column prop="IP" label="IP" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口" width="100"
                    :sort-method="(a: any, b: any) => { return a.Port - b.Port }" sortable
                    :show-overflow-tooltip="true" />
                <el-table-column prop="Domain" label="域名" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Protocol" label="协议" :filters='getColumnData("Protocol")'
                    :filter-method="filterHandlerProtocol" width="100" :show-overflow-tooltip="true">
                    <template #filter-icon>
                        <Filter />
                    </template>
                </el-table-column>
                <el-table-column prop="Product" label="组件" width="200">
                    <template #default="scope">
                        <el-button type="primary" plain v-if="formatProduct(scope.row.Product).length > 0">
                            <template #icon v-if="formatProduct(scope.row.Product).length > 1">
                                <el-popover placement="bottom" :width="350" trigger="hover">
                                    <template #reference>
                                        <el-icon>
                                            <Histogram />
                                        </el-icon>
                                    </template>
                                    <el-space direction="vertical">
                                        <el-tag round type="primary" v-for="component in formatProduct(scope.row.Product)"
                                        style="width: 320px;">{{ component }}</el-tag>
                                    </el-space>
                                </el-popover>
                            </template>
                            {{ formatProduct(scope.row.Product)[0] }}
                        </el-button>
                    </template>
                </el-table-column>
                <el-table-column prop="Region" label="地区" :show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="200" :show-overflow-tooltip="true" />
                <el-table-column fixed="right" label="操作" width="100" align="center">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="top">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="C段查询" placement="top">
                            <el-button link :icon="csegmentIcon"
                                @click.prevent="searchCsegmentIpv4(scope.row.IP)">
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination size="small" background v-model:page-size="item.pageSize" :page-sizes="[100, 500, 1000]"
                    layout="total, sizes, prev, pager, next, jumper" @size-change="tableCtrl.handleSizeChange"
                    @current-change="tableCtrl.handleCurrentChange" :total="item.total" />
            </div>
        </el-tab-pane>
        <el-empty v-else></el-empty>
    </el-tabs>
    <el-dialog v-model="syntax.starDialog.value" title="收藏语法" width="40%" center>
        <!-- 一定要用:model v-model校验会失效 -->
        <el-form ref="ruleFormRef" :model="syntax.ruleForm" :rules="global.syntaxRules" status-icon>
            <el-form-item label="语法名称" prop="Name">
                <el-input v-model="syntax.ruleForm.Name" maxlength="30" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="语法内容" prop="Content">
                <el-input v-model="syntax.ruleForm.Content" type="textarea" :rows="10" maxlength="1024"
                    show-word-limit></el-input>
            </el-form-item>
            <el-form-item class="align-right">
                <el-button type="primary" @click="syntax.submitStar(ruleFormRef)">
                    确定
                </el-button>
                <el-button @click="syntax.starDialog.value = false">取消</el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<style scoped>

</style>
