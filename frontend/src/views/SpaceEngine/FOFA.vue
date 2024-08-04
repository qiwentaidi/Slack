<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Search, ChromeFilled, CopyDocument, Grid, Operation, CollectionTag, Delete, Document, PictureRounded, Star, Collection } from '@element-plus/icons-vue';
import { SplitTextArea, validateIP, validateDomain, splitInt, Copy } from '@/util'
import { FofaResult, DefaultKeyValue, RuleForm, TableTabs } from "@/interface"
import { ExportToXlsx } from '@/export'
import { FofaTips, FofaSearch, IconHash } from 'wailsjs/go/main/App'
import { BrowserOpenURL } from 'wailsjs/runtime'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';
import global from "@/global"
import { InsertFavGrammarFiled, RemoveFavGrammarFiled, SelectAllSyntax } from 'wailsjs/go/main/Database';
const form = reactive({
    query: '',
    fraud: false,
    cert: false,
    tips: '',
    syntaxData: [] as RuleForm[],
})

// ref得单独校验
const ruleFormRef = ref<FormInstance>()

const syntax = ({
    keywordActive: "基础类",
    advSearch: [
        {
            key: '=',
            description: '匹配，=""时，可查询不存在字段或者值为空的情况。',
        },
        {
            key: '==',
            description: '完全匹配，==""时，可查询存在且值为空的情况。',
        },
        {
            key: '&&',
            description: '与',
        },
        {
            key: '||',
            description: '或',
        },
        {
            key: '!=',
            description: '不匹配，!=""时，可查询值不为空的情况。',
        },
        {
            key: '*=',
            description: '模糊匹配，使用*或者?进行搜索，比如banner*="mys??"。',
            level: "个人版"
        },
        {
            key: '()',
            description: '确认查询优先级，括号内容优先级最高。',
        }
    ],
    keywordSearch: [
        {
            title: "基础类",
            data: [
                { key: 'ip="1.1.1.1"', description: '通过单一IPv4地址进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'port="6379"', description: '通过端口号进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'domain="qq.com"', description: '通过根域名进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'host=".fofa.info"', description: '通过主机名进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'os="centos"', description: '通过操作系统进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'server="Microsoft-IIS/10"', description: '通过服务器进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'asn="19551"', description: '通过自治系统号进行搜索', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'org="LLC Baxet"', description: '通过所属组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'is_domain=(true/false)', description: '筛选(拥有/没有)域名的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'is_ipv6=(true/false)', description: '筛选是(ipv6/ipv4)的资产', filed1: "✓", filed2: "-", filed3: "-" },
            ]
        },
        {
            title: "标记类",
            data: [
                { key: 'app="Microsoft-Exchange"', description: '通过FOFA整理的规则进行查询', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'fid="sSXXGNUO2FefBTcCLIT/2Q=="', description: '通过FOFA聚合的站点指纹进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'product="NGINX"', description: '通过FOFA标记的产品名进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'category="服务"', description: '通过FOFA标记的分类进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
                { key: 'type="subdomain"', description: '筛选服务（网站类）资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'cloud_name="Aliyundun"', description: '通过云服务商进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'is_cloud=(true/false)', description: '筛选(是/不是)云服务的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'is_fraud=(true/false)', description: '筛选(是/不是)仿冒垃圾站群的资产', filed1: "✓", filed2: "-", filed3: "-", level: "专业版" },
                { key: 'is_honeypot=(true/false)', description: '筛选(是/不是)蜜罐的资产', filed1: "✓", filed2: "-", filed3: "-", level: "专业版" },
            ]
        },
        {
            title: "协议类",
            data: [
                { key: 'protocol="quic"', description: '通过协议名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'banner="users"', description: '通过协议返回信息进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'base_protocol="tcp"', description: '查询传输层为tcp协议的资产', filed1: "✓", filed2: "✓", filed3: "-" },
            ]
        },
        {
            title: "证书类",
            data: [
                { key: 'cert="baidu"', description: '通过证书进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject="Oracle Corporation"', description: '通过证书的持有者进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer="DigiCert"', description: '通过证书的颁发者进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject.org="Oracle Corporation"', description: '通过证书持有者的组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.subject.cn="baidu.com"', description: '通过证书持有者的通用名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer.org="cPanel, Inc."', description: '通过证书颁发者的组织进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.issuer.cn="Synology Inc. CA"', description: '通过证书颁发者的通用名称进行查询', filed1: "✓", filed2: "✓", filed3: "✓" },
                { key: 'cert.is_valid=(true/false)"', description: '筛选证书(是/不是)有效证书的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'cert.is_match=(true/false)', description: '筛选证书和域名(匹配/不匹配)的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'cert.is_expired=(true/false)', description: '筛选证书(已过期/未过期)的资产', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'tls.version="TLS 1.3"', description: '通过tls的协议版本进行查询', filed1: "✓", filed2: "✓", filed3: "-" },
            ]
        },
        {
            title: "时间类",
            data: [
                { key: 'after="2023-01-01"', description: '筛选某一时间之后有更新的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'before="2023-12-01"', description: '筛选某一时间之前有更新的资产', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'after="2023-01-01" && before="2023-12-01"', description: '筛选某一时间区间有更新的资产', filed1: "✓", filed2: "-", filed3: "-" }
            ]
        },
        {
            title: "独立IP语法（不可和上面其他语法共用）",
            data: [
                { key: 'port_size="6"', description: '筛选开放端口数量等于6个的独立IP', filed1: "✓", filed2: "✓", filed3: "-", level: "个人版" },
                { key: 'port_size_gt="6"', description: '筛选开放端口数量大于6个的独立IP', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'port_size_lt="12"', description: '筛选开放端口数量小于12个的独立IP', filed1: "✓", filed2: "-", filed3: "-", level: "个人版" },
                { key: 'ip_ports="80,161"', description: '筛选同时开放不同端口的独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_country="CN"', description: '通过国家的简称代码进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_region="Zhejiang"', description: '通过省份/地区英文名称进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_city="Hangzhou"', description: '通过城市英文名称进行查询独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_after="2021-03-18"', description: '筛选某一时间之后有更新的独立IP', filed1: "✓", filed2: "-", filed3: "-" },
                { key: 'ip_before="2019-09-09"', description: '筛选某一时间之前有更新的独立IP', filed1: "✓", filed2: "-", filed3: "-" }
            ]
        }
    ],
    rowClick: function (row: any, column: any, event: Event) {
        if (!form.query) {
            form.query = row.key
            return
        }
        form.query += " && " + row.key
    },
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
    querySearchAsync: async (queryString: string, cb: any) => {
        if (queryString.includes("=") || !queryString) {
            cb([]);
            return
        }
        let tips = await entry.getTips(queryString)
        cb(tips);
    },
    getTips: async function (queryString: string) {
        let tips = [] as DefaultKeyValue[]
        let result: any = await FofaTips(queryString)
        console.log(result)
        if (result.code == 0) {
            for (const item of result.data) {
                tips.push({
                    value: item.name,
                    label: item.company
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
        form.tips = result.Message!
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
            form.tips = result.Message!
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
            form.tips = result.Message!
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
            ExportToXlsx(["URL", "标签", "IP", "端口", "域名", "协议", "组件", "国家", "省份", "城市", "备案号"], "asset", "fofa_asset", temp)
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

function formatProduct(raw: string): string[] {
    return !raw ? [] : raw.split(",")
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
                        <el-popover placement="bottom-end" :width="800" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="语法检索" placement="bottom">
                                        <el-button :icon="CollectionTag" link />
                                    </el-tooltip>
                                </div>
                            </template>
                            <el-tabs v-model="syntax.keywordActive" class="quake">
                                <el-tab-pane v-for="item in syntax.keywordSearch" :name="item.title"
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
                                    <el-table :data="syntax.advSearch" stripe class="keyword-search">
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
                            <el-table :data="form.syntaxData" @row-click="syntax.rowClick" class="hunter-keyword-search">
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
                    </el-space>
                </template>
                <template #default="{ item }">
                    <div>
                        <span>{{ item.value }}</span>
                        <el-divider direction="vertical" />
                        <span style="color: #559FF8;">{{ item.label }}</span>
                    </div>
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
                        <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)"
                            v-show="scope.row.URL != ''">
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
                                    <div style="display: flex; flex-direction: column;">
                                        <el-tag type="primary" v-for="component in formatProduct(scope.row.Product)">{{
                                            component
                                            }}</el-tag>
                                    </div>
                                </el-popover>
                            </template>
                            {{ formatProduct(scope.row.Product)[0] }}
                        </el-button>
                    </template>
                </el-table-column>
                <el-table-column prop="Region" label="地区" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" :show-overflow-tooltip="true" />
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination background v-model:page-size="item.pageSize" :page-sizes="[100, 500, 1000]"
                    layout="total, sizes, prev, pager, next" @size-change="tableCtrl.handleSizeChange"
                    @current-change="tableCtrl.handleCurrentChange" :total="item.total" />
            </div>
        </el-tab-pane>
        <el-empty v-else></el-empty>
    </el-tabs>
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
.keyword-search :deep(.el-table__row:hover) {
    cursor: pointer;
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
