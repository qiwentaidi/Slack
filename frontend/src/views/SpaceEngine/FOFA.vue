<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Search, ChromeFilled, CopyDocument, Grid, Operation, ChatLineSquare, Delete, Document, PictureRounded, Star, Collection } from '@element-plus/icons-vue';
import { SplitTextArea, validateIP, validateDomain, splitInt, Copy } from '../../util'
import { FofaResult, DefaultKeyValue, RuleForm, TableTabs } from "../../interface"
import { ExportToXlsx } from '../../export'
import { FofaTips, FofaSearch, IconHash } from '../../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';
import global from "../../global"
import { InsertFavGrammarFiled, RemoveFavGrammarFiled, SelectAllSyntax } from '../../../wailsjs/go/main/Database';
const form = reactive({
    query: '',
    fraud: false,
    cert: false,
    tips: '',
    loadAll: [] as DefaultKeyValue[],
    syntaxData: [] as RuleForm[],
})

// ref得单独校验
const ruleFormRef = ref<FormInstance>()

const syntax = ({
    fofaImg: [
        '/fofa_syntax/fofa1.png',
        '/fofa_syntax/fofa2.png',
        '/fofa_syntax/fofa3.png',
        '/fofa_syntax/fofa4.png',
        '/fofa_syntax/fofa5.png',
        '/fofa_syntax/fofa6.png',
    ],
    searchDialog: ref(false),
    starDialog: ref(false),
    rules: reactive<FormRules<RuleForm>>({
        name: [
            { required: true, message: '请输入语法名称', trigger: 'blur' },
        ],
        desc: [
            {
                required: true,
                message: '请输入语法内容',
                trigger: 'blur',
            },
        ],
    }),
    ruleForm: reactive<RuleForm>({
        name: '',
        desc: '',
    }),
    createStarDialog: () => {
        syntax.starDialog.value = true
        syntax.ruleForm.name = ""
        syntax.ruleForm.desc = form.query
    },
    submitStar: async (formEl: FormInstance | undefined) => {
        if (!formEl) return
        let result = await formEl.validate()
        if (!result) return
        InsertFavGrammarFiled("fofa", syntax.ruleForm.name!, syntax.ruleForm.desc!).then((r: Boolean) => {
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
    querySearchAsync: (queryString: string, cb: any) => {
        if (queryString.includes("=") || queryString == "") {
            cb(form.loadAll);
            return
        }
        entry.getTips(queryString)
        cb(form.loadAll);
    },
    getTips: function (queryString: string) {
        FofaTips(queryString).then(result => {
            form.loadAll = []
            if (result.code == 0) {
                for (const item of result.data) {
                    form.loadAll.push({
                        value: item.name,
                        label: item.company
                    })
                }
            }
        })
    },
    handleSelect: (item: Record<string, any>) => {
        form.query = `app="${item.value}"`
    },
    rowClick: function (row: any, column: any, event: Event) {
        if (form.query == "") {
            form.query = row.syntax
            return
        }
        form.query += " && " + row.syntax
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
        let result: FofaResult = await FofaSearch(query, "100", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert)
        if (result.Error) {
            ElMessage({
                showClose: true,
                message: result.Message,
                type: "warning",
            });
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
        form.tips = result.Message + " 共查询到数据:" + result.Size + "条"
        const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
        tab.content = result.Results!;
        tab.total = result.Size!
        table.acvtiveNames = newTabName
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
    handleSizeChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        table.loading = true
        FofaSearch(tab.title, val.toString(), "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then((result: FofaResult) => {
            form.tips = result.Message! + " 共查询到数据:" + result.Size! + "条"
            if (result.Error) {
                table.loading = false
                return
            }
            tab.content = result.Results!;
            tab.total = result.Size!
        })
        table.loading = false
    },
    handleCurrentChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        table.loading = true
        FofaSearch(tab.title, tab.pageSize.toString(), val.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then((result: FofaResult) => {
            form.tips = result.Message + " 共查询到数据:" + result.Size + "条"
            if (result.Error) {
                table.loading = false
                return
            }
            tab.content = result.Results!;
            tab.total = result.Size!
        })
        table.loading = false
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
            if (hash == "") {
                ElMessage({
                    showClose: true,
                    message: "目标错误或不可达",
                    type: "warning",
                });
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
                if (validateIP(line)) {
                    temp.push(`ip="${line}"`);
                } else if (validateDomain(line)) {
                    temp.push(`domain="${line}"`);
                }
            }

            if (temp.length === 0) {
                ElMessage({
                    showClose: true,
                    message: "目标为空",
                    type: "warning",
                });
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
            let temp = [{}]
            temp.pop()
            ElMessage("正在进行全数据导出，API每页最大查询限度10000，请稍后...");
            let index = 0
            for (const num of splitInt(tab.total, 10000)) {
                index += 1
                ElMessage("正在导出第" + index.toString() + "页");
                await FofaSearch(tab.title, num.toString(), index.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then((result: FofaResult) => {
                    if (result.Error) {
                        return
                    }
                    temp.push(...result.Results!)
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

async function CopyDeduplicationURL() {
    await FofaSearch(form.query, "10000", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, form.fraud, form.cert).then((result: FofaResult) => {
        if (result.Error) {
            ElMessage("查询失败")
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
    })
}
</script>

<template>
    <el-form :model="form" @submit.native.prevent="tableCtrl.addTab(form.query)">
        <el-form-item>
            <el-autocomplete v-model="form.query" placeholder="Search..." :fetch-suggestions="entry.querySearchAsync"
                @select="entry.handleSelect" :debounce="1000" style="width: 100%;">
                <template #prepend>
                    查询条件
                </template>
                <template #suffix>
                    <el-space :size="2">
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
                            <el-table :data="form.syntaxData" @row-click="entry.rowClick" class="hunter-keyword-search">
                                <el-table-column width="150" prop="Name" label="语法名称" />
                                <el-table-column prop="Content" label="语法内容" />
                                <el-table-column label="操作" width="100">
                                    <template #default="scope">
                                        <el-button type="text"
                                            @click="syntax.deleteStar(scope.row.Name, scope.row.Content)">删除
                                        </el-button>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-popover>
                        <el-tooltip content="查询语法" placement="left">
                            <el-button :icon="ChatLineSquare" @click="syntax.searchDialog.value = true" />
                        </el-tooltip>
                    </el-space>
                </template>
                <template #default="{ item }">
                    <el-space>
                        <span>{{ item.value }}</span>
                        <el-tag>{{ item.link }}</el-tag>
                    </el-space>
                </template>
            </el-autocomplete>
        </el-form-item>
    </el-form>
    <div class="my-header" style="margin-bottom: 10px;">
        <div>
            <el-checkbox v-model="form.fraud">排除干扰(专业版)</el-checkbox>
            <el-checkbox v-model="form.cert">证书(个人版)</el-checkbox>
        </div>
        <el-dropdown>
            <el-button :icon="Operation" color="#D2DEE3" />
            <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item :icon="Grid" @click="SaveData(0)">导出当前查询页数据</el-dropdown-item>
                    <el-dropdown-item :icon="Grid" @click="SaveData(1)">导出全部数据</el-dropdown-item>
                    <el-dropdown-item :icon="CopyDocument" @click="CopyURL" divided>复制当前页URL</el-dropdown-item>
                    <el-dropdown-item :icon="CopyDocument" @click="CopyDeduplicationURL">去重复制前1w条URL</el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>
    </div>
    <el-tabs v-model="table.acvtiveNames" v-loading="table.loading" type="card" closable
        @tab-remove="tableCtrl.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column type="index" label="#" width="60px" />
                <el-table-column prop="URL" label="URL" width="300" :show-overflow-tooltip="true">
                    <template #default="scope">
                        <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)">
                        </el-button>
                        {{ scope.row.URL }}
                    </template>
                </el-table-column>
                <el-table-column prop="Title" label="标题" :filters='getColumnData("Title")'
                    :filter-method="filterHandlerTitle" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="IP" label="IP" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口" width="100"
                    :sort-method="(a: any, b: any) => { return a.Port - b.Port }" sortable
                    :show-overflow-tooltip="true" />
                <el-table-column prop="Domain" label="域名" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="Protocol" label="协议" :filters='getColumnData("Protocol")'
                    :filter-method="filterHandlerProtocol" width="100" :show-overflow-tooltip="true" />
                <el-table-column prop="Region" label="地区" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" :show-overflow-tooltip="true" />
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination background v-model:page-size="item.pageSize" :page-sizes="[100, 500, 1000]"
                    layout="sizes, prev, pager, next" @size-change="tableCtrl.handleSizeChange"
                    @current-change="tableCtrl.handleCurrentChange" :total="item.total" />
            </div>
        </el-tab-pane>
        <el-empty v-else></el-empty>
    </el-tabs>
    <el-dialog v-model="syntax.searchDialog.value" title="查询语法参考" width="80%" center>
        <div class="demo-image__lazy">
            <el-image v-for="url in syntax.fofaImg" :key="url" :src="url" lazy></el-image>
        </div>
    </el-dialog>
    <el-dialog v-model="syntax.starDialog.value" title="收藏语法" width="40%" center>
        <!-- 一定要用:model v-model校验会失效 -->
        <el-form ref="ruleFormRef" :model="syntax.ruleForm" :rules="syntax.rules" status-icon>
            <el-form-item label="语法名称" prop="name">
                <el-input v-model="syntax.ruleForm.name" maxlength="30" show-word-limit></el-input>
            </el-form-item>
            <el-form-item label="语法内容" prop="desc">
                <el-input v-model="syntax.ruleForm.desc" type="textarea" rows="10" maxlength="1024"
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
