<template>
    <el-form v-model="quake" @submit.native.prevent="tableCtrl.addTab(quake.query, false)">
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
                                <el-tab-pane v-for="item in options.keywordSearch" :name="item.title"
                                    :label="item.title">
                                    <el-table :data="item.data" class="keyword-search" @row-click="syntax.rowClick">
                                        <el-table-column width="300" property="key">
                                            <template #default="scope">
                                                {{ scope.row.key }}<el-tag type="success" effect="dark" color="#4CA87D"
                                                    v-if="scope.row.isVip" style="margin-left: 5px;">VIP</el-tag>
                                            </template>
                                        </el-table-column>
                                        <el-table-column property="description" />
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
                                <el-input v-model="quake.batchIps" type="textarea" rows="5"
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
                            <el-table :data="quake.syntaxData" @row-click="syntax.rowClick">
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
        <el-form-item style="display: flex; align-items: center; width: 100%;">
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
                        <el-dropdown-item :icon="Grid" @click="exportData">导出当前查询页数据</el-dropdown-item>
                        <el-dropdown-item :icon="Grid" @click="exportData">导出全部数据</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </el-form-item>
    </el-form>
    <el-tabs v-model="table.acvtiveNames" v-loading="table.loading" type="card" closable
        @tab-remove="tableCtrl.removeTab" class="quake-tabs">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column prop="IP" fixed label="IP" width="150" />
                <el-table-column prop="Host" fixed label="域名" width="200">
                    <template #default="scope">
                        <span v-if="scope.row.Host != scope.row.IP">{{ scope.row.Host }}</span>
                        <span v-else>--</span>
                    </template>
                </el-table-column>
                <el-table-column prop="Title" label="标题" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口 | 协议" width="170" :show-overflow-tooltip="true">
                    <template #default="scope">
                        {{ scope.row.Port }}
                        <el-divider direction="vertical" />
                        {{ scope.row.Protocol }}
                    </template>
                </el-table-column>
                <el-table-column prop="Component" label="产品应用/版本" width="260">
                    <template #default="scope">
                        <el-button type="success" plain
                            v-if="Array.isArray(scope.row.Component) && scope.row.Component.length > 0">
                            <template #icon v-if="scope.row.Component.length > 1">
                                <el-popover placement="bottom" :width="350" trigger="hover">
                                    <template #reference>
                                        <el-icon>
                                            <Histogram />
                                        </el-icon>
                                    </template>
                                    <div style="display: flex; flex-direction: column;">
                                        <el-tag type="success" v-for="component in scope.row.Component">{{ component
                                            }}</el-tag>
                                    </div>
                                </el-popover>
                            </template>
                            {{ scope.row.Component[0] }}
                        </el-button>
                    </template>
                </el-table-column>
                <el-table-column prop="IcpName" label="ICP名称" width="160" :show-overflow-tooltip="true" />
                <el-table-column prop="IcpNumber" label="ICP编号" width="150" :show-overflow-tooltip="true" />
                <el-table-column prop="ISP" label="运营商" width="100" :show-overflow-tooltip="true" />
                <el-table-column prop="Position" label="地理位置" width="200" :show-overflow-tooltip="true" />
                <el-table-column fixed="right" label="操作" width="100">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="top">
                            <el-button link :icon="ChromeFilled" @click.prevent="options.openHttpLink(scope.row)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="C段查询" placement="top">
                            <el-button link
                                @click.prevent="tableCtrl.addTab('ip: ' + CsegmentIpv4(scope.row.IP), false)">
                                <template #icon>
                                    <svg t="1719219479838" class="icon" viewBox="0 0 1450 1024" version="1.1"
                                        xmlns="http://www.w3.org/2000/svg" p-id="5099" width="200" height="200">
                                        <path
                                            d="M1055.229398 0c204.270977 0 352.3973 139.935005 352.3973 339.939672 0 5.972836-4.52229 10.83643-10.153821 10.83643h-122.869761c-9.04458 0-16.382635-7.423381-17.321223-17.065245-8.617949-114.166486-86.606116-192.325306-201.369886-192.325306-141.556204 0-221.250896 111.777352-221.250895 312.464628V575.098742c0 197.615532 79.950671 308.027664 221.165569 308.027664 114.337139 0 192.154654-73.039247 201.455212-179.35572 0.853262-9.471211 8.276644-16.894592 17.40655-16.894592h123.040413c5.631531 0 10.153821 4.863595 10.15382 10.83643 0 190.533456-148.808933 326.116824-352.3973 326.116824-238.401467 0-375.691359-168.775269-375.691359-449.498542V453.935505c0-282.771102 137.375219-453.850179 375.435381-453.850179zM552.31664 850.019832c20.136989 0 36.434297 16.638613 36.434297 37.202233v55.80335a36.860928 36.860928 0 0 1-36.434297 37.287559H79.097409a36.860928 36.860928 0 0 1-36.434298-37.287559v-55.80335c0-20.56362 16.297309-37.202233 36.434298-37.202233h473.219231zM461.358887 477.656195c20.051662 0 36.348971 16.638613 36.348971 37.202233v55.888676a36.860928 36.860928 0 0 1-36.348971 37.202234H79.097409A36.860928 36.860928 0 0 1 42.663111 570.832431v-55.888676c0-20.478293 16.297309-37.202233 36.434298-37.202233h382.261478z m90.957753-372.363636c20.136989 0 36.434297 16.72394 36.434297 37.287559v55.80335a36.860928 36.860928 0 0 1-36.434297 37.287559H79.097409A36.860928 36.860928 0 0 1 42.663111 198.383468v-55.80335c0-20.56362 16.297309-37.287559 36.434298-37.287559h473.219231z"
                                            fill="#A3A9B3" p-id="5100"></path>
                                    </svg>
                                </template>
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <span style="color: #4CA87D;">{{ quake.message }}</span>
                <el-pagination class="quake-pagin" background v-model:page-size="item.pageSize"
                    :page-sizes="[10, 20, 50, 100, 200, 500]" layout="total, sizes, prev, pager, next"
                    @size-change="tableCtrl.handleSizeChange" @current-change="tableCtrl.handleCurrentChange"
                    :total="item.total" />
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
                <el-button color="#4CA87D" :dark="true" style="color: #fff;" @click="syntax.submitStar(ruleFormRef)">
                    确定
                </el-button>
                <el-button @click="syntax.starDialog.value = false">取消</el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<script lang="ts" setup>
import { Search, ArrowDown, CopyDocument, Document, Grid, PictureRounded, Histogram, UploadFilled, Delete, Star, Collection, CollectionTag, ChromeFilled, QuestionFilled } from '@element-plus/icons-vue';
import { reactive, ref } from 'vue';
import { Copy, ReadLine, generateRandomString, splitInt, transformArrayFields, CsegmentIpv4, validateURL } from '../../util';
import { ExportToXlsx } from '../../export';
import { QuakeData, QuakeResult, QuakeTableTabs, QuakeTipsData, RuleForm } from '../../interface';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { FaviconMd5, QuakeSearch, QuakeTips } from '../../../wailsjs/go/main/App';
import global from '../../global';
import { ElMessage, ElNotification, FormInstance, FormRules } from 'element-plus';
import { FileDialog } from '../../../wailsjs/go/main/File';
import { InsertFavGrammarFiled, RemoveFavGrammarFiled, SelectAllSyntax } from '../../../wailsjs/go/main/Database';

const options = ({
    keywordActive: "基本信息",
    keywordSearch: [
        {
            title: "基本信息",
            data: [
                { key: 'ip:1.1.1.1/22', description: '基本信息' },
                { key: 'is_ipv6:true', description: '查询IPv6数据' },
                { key: 'port:80', description: '查询开放80端口的资产' },
                { key: 'port:[50 TO 60]', description: '查询开放端口在50至60之间的资产' },
                { key: 'transport:udp', description: '查询udp数据' },
                { key: 'domain:google.com', description: '查询资产域名' },
                { key: 'asn:12345', description: '查询自治域号码为"12345"的资产' },
                { key: 'org:No.31,Jin-rong Street', description: '查询自治域归属组织为"No.31,Jin-rong Street"的资产' },
                { key: 'hostname:unifiedlayer.com', description: '查询主机名包含"unifiedlayer.com"的资产' }
            ]
        },
        {
            title: "服务数据",
            data: [
                { key: 'service:http', description: '查询服务名称为"http"的资产' },
                { key: 'app:Apache', description: '查询应用名称为"Apache"的资产' },
                { key: 'response:奇虎科技', description: '查询端口原生返回数据中包含"奇虎科技"的资产' },
                { key: 'cert:奇虎科技', description: '查询包含"奇虎科技"的证书资产' },
                { key: 'catalog:IoT物联网', description: '查询应用类别为"IoT物联网"的资产', isVip: true },
                { key: 'type:VPN OR 防火墙', description: '查询应用类型为"VPN 或 防火墙"的资产', isVip: true },
                { key: 'level:硬件设备层', description: '查询应用层级为"硬件设备层"', isVip: true },
                { key: 'vendor:Sangfor OR 微软', description: '查询应用生产商为"Sangfor 或 微软"的资产', isVip: true }
            ]
        },
        {
            title: "IP归属与定位",
            data: [
                { key: 'country:China', description: '搜索国家，可使用简写CN' },
                { key: 'country_cn:中国', description: '搜索中文国家名称' },
                { key: 'province:Sichuan', description: '搜索英文省份名称' },
                { key: 'province_cn:四川', description: '搜索中文省份名称' },
                { key: 'city:Chengdu', description: '搜索英文城市名称' },
                { key: 'city_cn:成都', description: '搜索中文城市名称' },
                { key: 'owner:tencent.com', description: '搜索IP归属单位' },
                { key: 'isp:amazon.com', description: '搜索IP归属运营商' },
            ]
        },
        {
            title: "网页深度识别",
            data: [
                { key: 'icp:京ICP备08010314号', description: '查询ICP备案号', isVip: true },
                { key: 'icp_nature:企业', description: '查询备案主体性质' },
                { key: 'icp_keywords:奇虎', description: '查询备案网站中的关键词或域名', isVip: true },
                { key: 'link_url:/login', description: '查询网页url中包含"/login"的资产', isVip: true },
                { key: 'script_variable:csrfToken', description: '查询script标签变量中包含"csrfToken"的资产', isVip: true },
                { key: 'script_function:indexOf', description: '查询script标签函数中包含"indexOf"的资产', isVip: true },
                { key: 'css_class:login', description: '查询css标签class选择器中包含"login"的资产', isVip: true },
                { key: 'css_id:login', description: '查询css标签id选择器中包含"login"的资产', isVip: true },
                { key: 'iframe_title:企业微信登录', description: '查询iframe链接标题中包含"Enterprise wechat login"的资产', isVip: true },
                { key: 'iframe_url:/wwopen/sso/v1/', description: '查询iframe链接中包含"/wwopen/sso/v1/"的资产', isVip: true },
                { key: 'url_load:/login', description: '查询网页加载流中包含"/login"的资产', isVip: true },
                { key: 'page_type:门户网站', description: '查询类型是"门户网站"的资产', isVip: true }
            ]
        },
        {
            title: "布尔逻辑",
            data: [
                { key: 'port:"80" AND app:" 群晖NAS"', description: '查询开放80端口的"Synology NAS"应用资产' },
                { key: 'title:" 管理系统" OR body:" 登录页"', description: '查询网页标题包含"management system"或网页正文包含"登录页"的网站' },
                { key: 'title:"登录" AND(NOT title:"管理系统")', description: '查询网页标题包含"login"但不包含"管理系统"的网站' },
            ]
        }
    ],
    switch: reactive({
        latest: true,
        invalid: false,
        honeypot: false,
        cdn: false,
    }),
    openHttpLink: function (row: any) {
        if (row.Protocol == "http") {
            BrowserOpenURL("http://" + row.Host + ":" + row.Port)
        } else if (row.Protocol == "http/ssl") {
            BrowserOpenURL("https://" + row.Host + ":" + row.Port)
        }
    }
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
        if (!quake.query) {
            quake.query = row.key
            return
        }
        quake.query += " AND " + row.key
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
        syntax.ruleForm.desc = quake.query
    },
    submitStar: async (formEl: FormInstance | undefined) => {
        if (!formEl) return
        let result = await formEl.validate()
        if (!result) return
        InsertFavGrammarFiled("quake", syntax.ruleForm.name!, syntax.ruleForm.desc!).then((r: Boolean) => {
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
    syntaxData: [] as RuleForm[],
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
            content: [] as QuakeData[],
            total: 0,
            pageSize: 10,
            currentPage: 1,
            isBatch: isBatch,
            ipList: ipList,
        });
        table.loading = true
        QuakeSearch(ipList, query, 1, 10, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon).then(
            (result: QuakeResult) => {
                if (result.Code != 0) {
                    quake.message = result.Message!
                    table.loading = false
                    return
                }
                quake.message = "查询成功,目前剩余积分:" + result.Credit
                const tab = table.editableTabs.find(tab => tab.name === newTabName)!;
                result.Data?.forEach((item: QuakeData) => {
                    tab.content.push({
                        Host: item.Host,
                        IP: item.IP,
                        Port: item.Port,
                        Protocol: item.Protocol,
                        IcpName: item.IcpName,
                        Component: item.Components,
                        Title: item.Title,
                        IcpNumber: item.IcpNumber,
                        ISP: item.Isp,
                        Position: item.Position,
                    })
                });
                tab.total = result.Total!
                table.loading = false
            }
        )
        table.acvtiveNames = newTabName
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
        tab.pageSize = val
        tab.currentPage = 1
        table.loading = true
        QuakeSearch(tab.ipList, tab.title, 1, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon).then(
            (result: QuakeResult) => {
                if (result.Code != 0) {
                    quake.message = result.Message!
                    table.loading = false
                    return
                }
                quake.message = "查询成功,目前剩余积分:" + result.Credit
                tab.content = []
                result.Data?.forEach((item: QuakeData) => {
                    tab.content.push({
                        Host: item.Host,
                        IP: item.IP,
                        Port: item.Port,
                        Protocol: item.Protocol,
                        IcpName: item.IcpName,
                        Component: item.Components,
                        Title: item.Title,
                        IcpNumber: item.IcpNumber,
                        ISP: item.Isp,
                        Position: item.Position,
                    })
                });
                tab.total = result.Total!
                table.loading = false
            }
        )
    },
    handleCurrentChange: async (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
        table.loading = true
        QuakeSearch(tab.ipList, tab.title, tab.currentPage, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon).then(
            (result: QuakeResult) => {
                if (result.Code != 0) {
                    quake.message = result.Message!
                    table.loading = false
                    return
                }
                quake.message = "查询成功,目前剩余积分:" + result.Credit
                tab.content = []
                result.Data?.forEach((item: QuakeData) => {
                    tab.content.push({
                        Host: item.Host,
                        IP: item.IP,
                        Port: item.Port,
                        Protocol: item.Protocol,
                        IcpName: item.IcpName,
                        Component: item.Components,
                        Title: item.Title,
                        IcpNumber: item.IcpNumber,
                        ISP: item.Isp,
                        Position: item.Position,
                    })
                });
                tab.total = result.Total!
                table.loading = false
            }
        )
    },
    handleOptionChange: function () {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.pageSize = 10
        tab.currentPage = 1
        table.loading = true
        QuakeSearch(tab.ipList, tab.title, 1, tab.pageSize, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon).then(
            (result: QuakeResult) => {
                if (result.Code != 0) {
                    quake.message = result.Message!
                    table.loading = false
                    return
                }
                quake.message = "查询成功,目前剩余积分:" + result.Credit
                tab.content = []
                result.Data?.forEach((item: QuakeData) => {
                    tab.content.push({
                        Host: item.Host,
                        IP: item.IP,
                        Port: item.Port,
                        Protocol: item.Protocol,
                        IcpName: item.IcpName,
                        Component: item.Components,
                        Title: item.Title,
                        IcpNumber: item.IcpNumber,
                        ISP: item.Isp,
                        Position: item.Position,
                    })
                });
                tab.total = result.Total!
                table.loading = false
            }
        )
    },
    searchFavicon: function () {
    if (!quake.iconURL && !quake.iconFile) {
        ElMessage.warning("请输入URL或者上传图标");
        return;
    }
    let target = quake.iconFile || validateURL(quake.iconURL) ? quake.iconURL : quake.batchFile
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
        let temp = [{}]
        temp.pop()
        if (tab.isBatch) {
            ipList = await tableCtrl.getIpList()
        }
        let index = 0
        for (const num of splitInt(tab.total, 500)) {
            index += 1
            ElMessage("正在导出第" + index.toString() + "页");
            let result: QuakeResult = await QuakeSearch(ipList, tab.title, index, num, options.switch.latest, options.switch.invalid, options.switch.honeypot, options.switch.cdn, global.space.quakekey, quake.certcommon)
            if (result.Code != 0) {
                quake.message = result.Message!
                table.loading = false
                return
            }
            result.Data?.forEach((item: QuakeData) => {
                temp.push({
                    Host: item.Host,
                    IP: item.IP,
                    Port: item.Port,
                    Protocol: item.Protocol,
                    IcpName: item.IcpName,
                    Component: item.Components,
                    Title: item.Title,
                    IcpNumber: item.IcpNumber,
                    ISP: item.Isp,
                    Position: item.Position,
                })
            });
        }
        ExportToXlsx(["IP", "域名", "标题", "端口", "协议", "应用/组件", "证书申请单位", "证书域名", "运营商", "地理位置"], "asset", "quake_asset", transformArrayFields(temp))
        temp = []
    }
}
</script>

<style scoped>
.keyword-search :deep(.el-table__header-wrapper) {
    height: 0;
}

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