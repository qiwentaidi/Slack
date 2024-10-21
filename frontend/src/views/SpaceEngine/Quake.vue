<template>
    <el-form v-model="quake" @keydown.enter.native.prevent="tableCtrl.addTab(quake.query, false)">
        <el-form-item>
            <el-autocomplete v-model="quake.query" placeholder="Search..." :fetch-suggestions="syntax.querySearchAsync"
                @select="syntax.handleSelect" :trigger-on-focus="false" :debounce="1000" style="width: 100%;">
                <template #prepend>
                    查询条件
                </template>
                <template #suffix>
                    <el-space :size="2">
                        <el-popover placement="bottom-end" :width="700" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="语法检索" placement="bottom">
                                        <el-button :icon="CollectionTag" link />
                                    </el-tooltip>
                                </div>
                            </template>
                            <el-tabs v-model="options.keywordActive" class="quake">
                                <el-tab-pane v-for="item in quakeSyntaxOptions" :name="item.title" :label="item.title">
                                    <el-table :data="item.data" class="keyword-search" @row-click="syntax.rowClick">
                                        <el-table-column width="300" property="key" label="例句">
                                            <template #default="scope">
                                                {{ scope.row.key }}<el-tag type="success" effect="dark" color="#4CA87D"
                                                    v-if="scope.row.isVip" style="margin-left: 5px;">VIP</el-tag>
                                            </template>
                                        </el-table-column>
                                        <el-table-column property="description" label="用途说明" />
                                    </el-table>
                                </el-tab-pane>
                            </el-tabs>
                        </el-popover>
                        <el-popover placement="bottom-end" :width="400" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="使用网页图标搜索" placement="bottom">
                                        <el-button :icon="PictureRounded" link />
                                    </el-tooltip>
                                </div>
                            </template>
                            <div class="batch-search">
                                <el-alert type="info" :closable="false" title="如果上传文件与URL均不为空时优先检测文件" show-icon />
                                <el-input v-model="quake.iconURL" placeholder="Favicon URL地址"></el-input>
                                <el-button class="upload" :icon="UploadFilled">上传图片文件</el-button>
                                <div class="my-header">
                                    <div></div>
                                    <el-button color="#4CA87D" :dark="true" class="search"
                                        @click="tableCtrl.searchFavicon">检索</el-button>
                                </div>
                            </div>
                        </el-popover>
                        <el-popover placement="bottom-end" :width="400" title="IP/域名批量搜索" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="IP/域名批量搜索" placement="bottom">
                                        <el-button :icon="Document" link />
                                    </el-tooltip>
                                </div>
                            </template>
                            <div class="batch-search">
                                <el-alert type="info" :closable="false" title="上传包含IP/域名的.txt文件，数量不超过1000个" show-icon />
                                <el-button class="upload" :icon="UploadFilled">上传文件</el-button>
                                <el-input v-model="quake.batchIps" type="textarea" :rows="5"
                                    placeholder="请输入IP/域名，每行一个，多个请换行输入"></el-input>
                                <div class="my-header">
                                    <div></div>
                                    <el-button color="#4CA87D" :dark="true" class="search"
                                        @click="tableCtrl.addTab(generateRandomString(12), true)">检索</el-button>
                                </div>
                            </div>
                        </el-popover>
                    </el-space>
                    <el-divider direction="vertical" />
                    <el-space :size="2">
                        <el-tooltip content="清空语法" placement="bottom">
                            <el-button :icon="Delete" link @click="quake.query = ''" />
                        </el-tooltip>
                        <el-tooltip content="收藏语法" placement="bottom">
                            <el-button :icon="Star" link @click="syntax.starDialog.value = true" />
                        </el-tooltip>
                        <el-tooltip content="复制语法" placement="bottom">
                            <el-button :icon="CopyDocument" link @click="Copy(quake.query)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                    </el-space>
                    <el-button link :icon="Search" @click="tableCtrl.addTab(quake.query, false)"
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
                            <el-table :data="quake.syntaxData" @row-click="syntax.rowClick2" class="keyword-search">
                                <el-table-column width="150" prop="Name" label="语法名称" />
                                <el-table-column prop="Content" label="语法内容" />
                                <el-table-column label="操作" width="100">
                                    <template #default="scope">
                                        <el-button link style="color: #4CA87D"
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
                        <span>{{ item.Product_name }}</span>
                        <el-divider direction="vertical" />
                        <span v-if="item.Vendor_name != ''">{{ item.Vendor_name }}</span>
                        <span v-else>-</span>
                        <el-divider direction="vertical" />
                        <el-button link style="color: #4CA87D;">
                            测绘资产数量: {{ item.Ip_count }}
                        </el-button>
                        <el-divider direction="vertical" />
                        <el-button link style="color: #F56C6C;">
                            关联漏洞: {{ item.Vul_count }}
                        </el-button>
                    </div>
                </template>
            </el-autocomplete>
        </el-form-item>
        <el-form-item>
            <div>
                <span class="mr">最新数据</span><el-switch v-model="options.switch.latest"
                    @change="tableCtrl.handleOptionChange" style="--el-switch-on-color: #4CA87D;" />
            </div>
            <el-divider direction="vertical" />
            <div>
                <el-tooltip content="开启后，将过滤掉400、401、502等状态码和无法解析的协议/端口数据" placement="bottom">
                    <span class="mr">过滤无效请求</span>
                </el-tooltip>
                <el-switch v-model="options.switch.invalid" @change="tableCtrl.handleOptionChange"
                    :disabled="quake.certcommon.length == 0" style="--el-switch-on-color: #4CA87D;" />
            </div>
            <el-divider direction="vertical" />
            <div>
                <span class="mr">排除蜜罐</span><el-switch v-model="options.switch.honeypot"
                    :disabled="quake.certcommon.length == 0" @change="tableCtrl.handleOptionChange"
                    style="--el-switch-on-color: #4CA87D;" />
            </div>
            <el-divider direction="vertical" />
            <div>
                <span class="mr">排除CDN</span><el-switch v-model="options.switch.cdn"
                    :disabled="quake.certcommon.length == 0" @change="tableCtrl.handleOptionChange"
                    style="--el-switch-on-color: #4CA87D;" />
            </div>
            <el-divider direction="vertical" />
            <div>
                <span class="mr">CertCommon</span>
                <el-input v-model="quake.certcommon" size="small" style="width: 250px;">
                    <template #prefix>
                        <el-tooltip content="由于排除字段的值为按时间自动生成，请填写网页版登录后的Cookie中的cert_common字段，排除才可正常使用">
                            <el-icon>
                                <QuestionFilled />
                            </el-icon>
                        </el-tooltip>
                    </template>
                </el-input>
            </div>
            <div style="flex-grow: 1;"></div>
            <el-dropdown>
                <el-button :dark="true" text bg style="color: #4CA87D;">
                    更多功能<el-icon class="el-icon--right">
                        <ArrowDown />
                    </el-icon>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item :icon="exportIcon" @click="exportData">导出当前查询页数据</el-dropdown-item>
                        <el-dropdown-item :icon="exportIcon" @click="exportData">导出全部数据</el-dropdown-item>
                        <el-dropdown-item :icon="CopyDocument" @click="CopyURL" divided>复制当前页URL</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </el-form-item>
    </el-form>
    <el-tabs v-model="table.acvtiveNames" v-loading="table.loading" type="card" closable
        @tab-remove="tableCtrl.removeTab" class="quake-tabs">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: calc(100vh - 280px);">
                <el-table-column prop="URL" fixed label="URL" width="240" :show-overflow-tooltip="true" />
                <el-table-column prop="IP" fixed label="IP" width="130" />
                <el-table-column prop="Port" label="端口/协议" width="130">
                    <template #default="scope">
                        {{ scope.row.Port }}<el-tag type=info>{{ scope.row.Protocol }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="Host" label="域名" width="150" :show-overflow-tooltip="true">
                    <template #default="scope">
                        <span v-if="scope.row.Host != scope.row.IP">{{ scope.row.Host }}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column prop="Title" label="标题" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="Component" label="产品应用/版本" width="260">
                    <template #default="scope">
                        <el-button type="success" plain
                            v-if="Array.isArray(scope.row.Components) && scope.row.Components.length > 0">
                            <template #icon v-if="scope.row.Components.length > 1">
                                <el-popover placement="bottom" :width="350" trigger="hover">
                                    <template #reference>
                                        <el-icon>
                                            <Histogram />
                                        </el-icon>
                                    </template>
                                    <el-space direction="vertical">
                                        <el-tag round type="success" v-for="component in scope.row.Components"
                                            style="width: 320px;">
                                            {{ component }}</el-tag>
                                    </el-space>
                                </el-popover>
                            </template>
                            {{ scope.row.Components[0] }}
                        </el-button>
                    </template>
                </el-table-column>
                <el-table-column prop="IcpName" label="备案名称" width="160" :show-overflow-tooltip="true" />
                <el-table-column prop="IcpNumber" label="备案号" :show-overflow-tooltip="true" />
                <el-table-column prop="Isp" label="运营商" width="100" :show-overflow-tooltip="true" />
                <el-table-column prop="Position" label="地理位置" width="200" :show-overflow-tooltip="true" />
                <el-table-column fixed="right" label="操作" width="100" align="center">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="top">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="C段查询" placement="top">
                            <el-button link :icon="csegmentIcon"
                                @click.prevent="tableCtrl.addTab('ip: ' + CsegmentIpv4(scope.row.IP), false)">
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: #4CA87D; font-size: 14px;">{{ quake.message }}</span>
                <el-pagination size="small" class="quake-pagin" background v-model:page-size="item.pageSize"
                    :page-sizes="[10, 20, 50, 100, 200, 500]" layout="total, sizes, prev, pager, next, jumper"
                    @size-change="tableCtrl.handleSizeChange" @current-change="tableCtrl.handleCurrentChange"
                    :total="item.total" />
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
                <el-button color="#4CA87D" :dark="true" style="color: #fff;" @click="syntax.submitStar(ruleFormRef)">
                    确定
                </el-button>
                <el-button @click="syntax.starDialog.value = false">取消</el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<script lang="ts" setup>
import { Search, ArrowDown, CopyDocument, Document, PictureRounded, Histogram, UploadFilled, Delete, Star, Collection, CollectionTag, ChromeFilled, QuestionFilled } from '@element-plus/icons-vue';
import { reactive, ref } from 'vue';
import { Copy, ReadLine, generateRandomString, splitInt, transformArrayFields, CsegmentIpv4 } from '@/util';
import { ExportToXlsx } from '@/export';
import { QuakeTableTabs, QuakeTipsData } from '@/interface';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
import { Callgologger, FaviconMd5, QuakeSearch, QuakeTips } from 'wailsjs/go/main/App';
import global from '@/global';
import { ElMessage, ElNotification, FormInstance } from 'element-plus';
import { FileDialog } from 'wailsjs/go/main/File';
import { InsertFavGrammarFiled, RemoveFavGrammarFiled, SelectAllSyntax } from 'wailsjs/go/main/Database';
import exportIcon from '@/assets/icon/doucment-export.svg'
import csegmentIcon from '@/assets/icon/csegment.svg'
import { validateSingleURL } from '@/stores/validate';
import { main, space } from 'wailsjs/go/models';
import { quakeSyntaxOptions } from '@/stores/options';

const options = ({
    keywordActive: "基本信息",
    switch: reactive({
        latest: true,
        invalid: false,
        honeypot: false,
        cdn: false,
    }),
})

const ruleFormRef = ref<FormInstance>()

const syntax = ({
    querySearchAsync: (queryString: string, cb: Function) => {
        if (queryString.includes(":") || !queryString) {
            cb(quake.tips)
            return
        }
        syntax.getTips(queryString)
        cb(quake.tips);
    },
    getTips: async function (queryString: string) {
        quake.tips = []
        let result = await QuakeTips(queryString)
        if (result.code != 0) {
            return
        }
        for (const item of result.data!) {
            quake.tips.push({
                Product_name: item.product_name,
                Vul_count: item.vul_count,
                Vendor_name: item.vendor_name,
                Ip_count: item.ip_count,
            })
        }
    },
    handleSelect: (item: Record<string, any>) => {
        quake.query = `app:"${item.Product_name}"`
    },
    rowClick: function (row: any, column: any, event: Event) {
        if (quake.query == "") {
            quake.query = row.key
            return
        }
        quake.query += " AND " + row.key
    },
    rowClick2: function (row: any, column: any, event: Event) {
        if (quake.query == "") {
            quake.query = row.Content
            return
        }
        quake.query += " AND " + row.Content
    },
    starDialog: ref(false),
    ruleForm: reactive<main.Syntax>({
        Name: '',
        Content: '',
    }),
    createStarDialog: () => {
        syntax.starDialog.value = true
        syntax.ruleForm.Name = ""
        syntax.ruleForm.Content = quake.query
    },
    submitStar: async (formEl: FormInstance | undefined) => {
        if (!formEl) return
        let result = await formEl.validate()
        if (!result) return
        InsertFavGrammarFiled("quake", syntax.ruleForm.Name!, syntax.ruleForm.Content!).then((r: Boolean) => {
            if (r) {
                ElMessage.success('添加语法成功')
            } else {
                ElMessage.error('添加语法失败')
            }
            syntax.starDialog.value = false
        })
    },
    deleteStar: (name: string, content: string) => {
        RemoveFavGrammarFiled("quake", name, content).then((r: Boolean) => {
            if (r) {
                ElMessage.success('删除语法成功,重新打开刷新')
            } else {
                ElMessage.error('删除语法失败')
            }
        })
    },
    searchStarSyntax: async () => {
        quake.syntaxData = await SelectAllSyntax("quake")
    },
})

const quake = reactive({
    query: '',
    message: "",
    tips: [] as QuakeTipsData[],
    iconURL: "",
    iconFile: "",
    batchIps: "",
    batchFile: "",
    syntaxData: [] as main.Syntax[],
    certcommon: "",
})

const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as QuakeTableTabs[],
    loading: false,
})

const tableCtrl = ({
    addTab: async (query: string, isBatch: boolean) => {
        if (!query) {
            ElMessage.warning("请输入查询语句")
            return
        }
        const newTabName = `${++table.tabIndex}`
        let ipList = [] as string[]
        if (isBatch) {
            ipList = await tableCtrl.getIpList()
        }
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: [] as space.QuakeData[],
            total: 0,
            pageSize: 10,
            currentPage: 1,
            isBatch: isBatch,
            ipList: ipList,
        });
        table.loading = true
        table.acvtiveNames = newTabName
        let result = await QuakeSearch(ipList, query, 1, 10, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
        if (result.Code != 0) {
            quake.message = result.Message!
            table.loading = false
            return
        }
        quake.message = "查询成功,目前剩余积分:" + result.Credit
        const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
        tab.content = result.Data
        tab.total = result.Total!
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
        tab.pageSize = val
        tab.currentPage = 1
        table.loading = true
        let result = await QuakeSearch(tab.ipList, tab.title, 1, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
        if (result.Code != 0) {
            quake.message = result.Message!
            table.loading = false
            return
        }
        quake.message = "查询成功,目前剩余积分:" + result.Credit
        tab.content = result.Data
        tab.total = result.Total!
        table.loading = false
    },
    handleCurrentChange: async (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        table.loading = true
        let result = await QuakeSearch(tab.ipList, tab.title, tab.currentPage, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
        if (result.Code != 0) {
            quake.message = result.Message!
            table.loading = false
            return
        }
        quake.message = "查询成功,目前剩余积分:" + result.Credit
        tab.content = result.Data
        tab.total = result.Total!
        table.loading = false
    },
    handleOptionChange: async () => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.pageSize = 10
        tab.currentPage = 1
        table.loading = true
        let result = await QuakeSearch(tab.ipList, tab.title, 1, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
        if (result.Code != 0) {
            quake.message = result.Message!
            table.loading = false
            return
        }
        quake.message = "查询成功,目前剩余积分:" + result.Credit
        tab.content = result.Data
        tab.total = result.Total!
        table.loading = false
    },
    searchFavicon: function () {
        if (!quake.iconURL && !quake.iconFile) {
            ElMessage.warning("请输入URL或者上传图标");
            return;
        }
        let target = quake.iconFile || validateSingleURL(quake.iconURL) ? quake.iconURL : quake.batchFile
        if (!target) {
            return;
        }
        FaviconMd5(target).then((result: string) => {
            tableCtrl.addTab(`favicon:${result}`, false);
        });
    },
    // type 0 choose txt , type 1 choose img
    handleFileupload: async function (type: number) {
        if (type == 0) {
            quake.iconFile = await FileDialog("")
        } else {
            quake.batchFile = await FileDialog("*.txt")
            quake.batchIps = (await ReadLine(quake.batchFile))!.join("\n")
        }
    },
    getIpList: async function () {
        let ips = quake.batchIps.split("\n")
        if (ips.length > 1000) {
            ElMessage.warning("最多支持1000个IP")
            return []
        }
        return ips
    },
})

async function CopyURL() {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        const temp = tab.content.map(item => {
            if (item.URL != "") return item.URL
        })
        Copy(temp.join("\n"))
    }
}

async function exportData() {
    if (table.editableTabs.length != 0) {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        if (tab.total > 500) {
            ElNotification.info({
                title: "提示",
                message: "正在进行全数据导出，API每页最大查询限度500，请稍后。",
            });
        }
        let ipList = [] as string[]
        let temp = [] as space.QuakeData[]
        if (tab.isBatch) {
            ipList = await tableCtrl.getIpList()
        }
        let index = 0
        let numbs = splitInt(tab.total, 500)
        for (const num of numbs) {
            index += 1
            if (numbs.length != 1) {
                ElMessage("正在导出第" + index.toString() + "页");
            }
            let result = await QuakeSearch(ipList, tab.title, index, num, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
            if (result.Code != 0) {
                quake.message = result.Message!
                ElMessage.error(quake.message + " 已退出导出!")
                Callgologger("error", `[quake] ${tab.title} export data error: ${quake.message}`)
                table.loading = false
                return
            }
            temp.push(...result.Data)
        }
        ExportToXlsx(["URL", "应用/组件", "端口", "协议", "域名", "标题", "单位名称", "备案号", "IP", "运营商", "地理位置"], "asset", "quake_asset", transformArrayFields(temp))
        temp = []
    }
}
</script>

<style scoped>
.keyword-search :deep(.el-table__row:hover) {
    color: #4CA87D;
    cursor: pointer;
}

.quake :deep(.el-tabs__nav-scroll) {
    display: flex;
    justify-content: center;
    align-items: center;
}

.quake :deep(.el-tabs__item:hover) {
    color: #4CA87D;
}

.quake :deep(.el-tabs__item.is-active) {
    color: #4CA87D;
    /* 文本颜色 */
}

.quake :deep(.el-tabs__active-bar) {
    background-color: #4CA87D;
}

.el-alert .--el-alert-icon-large-size {
    width: 16px;
}

.batch-search {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .upload {
        height: 50px;
        width: 100%;
        border-style: dashed;
        color: #4CA87D;
    }

    .upload:hover {
        background-color: #F9FAFD;
    }

    .search {
        color: #fff;
        width: 20%;
    }

    .search:hover {
        color: #fff;
    }
}

.mr {
    margin-right: 10px;
}

.quake-tabs :deep(.el-tabs__item) {
    position: relative;
    display: inline-block;
    max-width: 300px;
    margin-bottom: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding-right: 26px !important;
    /* 给关闭按钮预留空间 */
}

.quake-tabs :deep(.el-tabs__item:hover) {
    color: #4CA87D;
    cursor: pointer;
}

.quake-tabs :deep(.el-tabs__item .el-icon) {
    position: absolute !important;
    top: 13px !important;
    right: 7px !important;
}

.quake-tabs :deep(.el-tabs__nav) {
    line-height: 255%;
}

.quake-tabs :deep(.el-tabs__item.is-active) {
    color: #4CA87D;
}

.quake-pagin :deep(.el-pager li.is-active) {
    background-color: #4CA87D;
}

.quake-pagin :deep(.el-pager li.is-active:hover) {
    color: #fff;
}

.quake-pagin :deep(.el-pager li:hover) {
    color: #4CA87D;
}

.quake-pagin :deep(.el-select-dropdown__item.is-selected) {
    color: #4CA87D;
    font-weight: bold;
}
</style>