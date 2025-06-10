<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from "vue";
import { Plus, Share, ChromeFilled } from '@element-plus/icons-vue';
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { FetchCompanyInfo, GoFetch, ResumeAfterHumanCheck } from "wailsjs/go/services/App";
import CustomTabs from "@/components/CustomTabs.vue";
import wechatIcon from "@/assets/icon/wechatOfficialAccount.svg"
import { debounce } from "lodash"
import { BrowserOpenURL, EventsOff, EventsOn } from "wailsjs/runtime/runtime";
import { structs } from "wailsjs/go/models";
import { getProxy, ProcessTextAreaInput } from "@/util";
import { FilepathJoin, OpenFolder } from "wailsjs/go/services/File";
import global from '@/stores';

const companiesInfo = ref<structs.CompanyInfo[]>([])

onMounted(async () => {
    const storedToken = localStorage.getItem('tyc-token');
    if (storedToken) {
        ruleForm.tycToken = storedToken;
    }
    const storedId = localStorage.getItem('tyc-id');
    if (storedId) {
        ruleForm.tycId = storedId;
    }
    const storedMiit = localStorage.getItem('miit-api');
    if (storedMiit) {
        ruleForm.miitApi = storedMiit;
    }
    from.configPath = await FilepathJoin([global.PATH.homedir, "/slack/company_info"])
    EventsOn("tyc-human-check", (msg: string) => {
        from.humanCheckDialogVisible = true
        from.humanCheckMessage = msg
    })
    EventsOff("tyc-human-check")
})

const handleHumanCheckConfirmed = () => {
    from.humanCheckDialogVisible = false
    ResumeAfterHumanCheck() // 调用后端接口，释放阻塞
}

const storageConfig = debounce(() => {
    // 2秒后将数据存储到localStorage
    localStorage.setItem('tyc-token', ruleForm.tycToken);
    localStorage.setItem('tyc-id', ruleForm.tycId);
    localStorage.setItem('miit-api', ruleForm.miitApi);
}, 2000);


const from = reactive({
    newTask: false,
    active: 1,
    activeName: 'subcompany',
    runningStatus: false,
    humanCheckDialogVisible: false,
    humanCheckMessage: "",
    configPath: "",
})

interface RuleForm {
    company: string;
    invest: number;
    dataSource: string[];
    tycToken: string;
    tycId: string;
    subcompanyLevel: number;
    domain: boolean;
    miitApi: string;
}

const ruleFormRef = ref<FormInstance>()

const ruleForm = reactive<RuleForm>({
    company: '',
    invest: 100,
    dataSource: ['tyc', 'miit'],
    tycToken: '',
    tycId: '',
    subcompanyLevel: 1,
    domain: false,
    miitApi: 'http://127.0.0.1:16181/',
})

const rules = reactive<FormRules<RuleForm>>({
    company: [
        { required: true, message: '请输入公司名称, 多个公司换行处理', trigger: 'blur' },
    ],
    tycToken: [
        {
            required: true,
            validator: (_, value, callback) => {
                if (ruleForm.dataSource.includes("tyc") && !value.trim()) {
                    callback(new Error('请填写登录天眼查后请求头中的X-Auth-Token字段'));
                } else {
                    callback();
                }
            },
            trigger: 'blur',
        },
    ],
    tycId: [
        {
            required: true,
            validator: (_, value, callback) => {
                if (ruleForm.dataSource.includes("tyc") && !value.trim()) {
                    callback(new Error('请填写登录天眼查后请求头中的X-TYCID字段'));
                } else {
                    callback();
                }
            },
            trigger: 'blur',
        },
    ],
    miitApi: [
        { required: true, message: '接口地址不能为空', trigger: 'blur' }
    ]
})
async function Collect() {
    let response = await GoFetch("GET", ruleForm.miitApi, "", {}, 10, getProxy())
    if (response.Error) {
        ElMessage.warning("请求MIIT接口失败! 请根据ICP_Query项目搭建自身接口")
        return
    }
    try {
        let jsonResult = JSON.parse(response.Body)
        if (jsonResult.msg != "查询请访问http://0.0.0.0:16181/query/{name}") {
            ElMessage.warning("请求接口结果异常!")
            return
        }
    } catch (error) {
        ElMessage.warning(error)
        return
    }
    ruleForm.tycToken = ruleForm.tycToken.replace(/[\r\n\s]/g, '')
    ruleForm.tycId = ruleForm.tycId.replace(/[\r\n\s]/g, '')
    from.newTask = false

    let dataSource: structs.DataSource = {
        Tianyancha: {} as structs.Tianyancha,
        Miit: {} as structs.Miit,
        convertValues: () => { }
    };

    dataSource.Tianyancha = {
        Enable: ruleForm.dataSource.includes('tyc'),
        Token: ruleForm.tycToken,
        Id: ruleForm.tycId
    }

    dataSource.Miit.API = ruleForm.miitApi
    let compnaies = ProcessTextAreaInput(ruleForm.company)
    from.runningStatus = true
    companiesInfo.value = []
    for (const compamy of compnaies) {
        let result = await FetchCompanyInfo(compamy, ruleForm.invest, dataSource, ruleForm.subcompanyLevel)
        companiesInfo.value.push(result)
    }
    from.runningStatus = false
}

const allOfficialAccounts = computed(() =>
    extractAllChildren(companiesInfo.value, 'OfficialAccounts')
);

const allApps = computed(() =>
    extractAllChildren(companiesInfo.value, 'Apps')
);

const allApplets = computed(() =>
    extractAllChildren(companiesInfo.value, 'Applets')
);
// 给每个子项加上所属公司名
function extractAllChildren<T>(
    companies: any[],
    key: 'OfficialAccounts' | 'Apps' | 'Applets'
): (T & { BelongCompany: string })[] {
    const result: any[] = [];

    const dfs = (list: any[]) => {
        for (const company of list) {
            const children = company[key];
            if (Array.isArray(children)) {
                for (const item of children) {
                    result.push({
                        ...item,
                        BelongCompany: company.CompanyName,
                    });
                }
            }

            if (Array.isArray(company.Subsidiaries)) {
                dfs(company.Subsidiaries);
            }
        }
    };

    dfs(companies);
    return result;
}
</script>


<template>
    <CustomTabs>
        <el-tabs v-model="from.activeName" type="card">
            <el-tab-pane label="对外投资" name="subcompany">
                <el-table :data="companiesInfo" row-key="CompanyName" :highlight-current-row="true" border
                    :tree-props="{ children: 'Subsidiaries' }">
                    <el-table-column prop="CompanyName" label="公司名称">
                        <template #default="scope">
                            <div class="company-cell">
                                <div v-if="scope.row.Trademark == '<nil>'"></div>
                                <el-popover placement="right" :width="180" v-else>
                                    <template #reference>
                                        <el-image :src="scope.row.Trademark" class="avatar mr-5px">
                                            <template #error>
                                                <div></div>
                                            </template>
                                        </el-image>
                                    </template>
                                    <template #default><el-image :src="scope.row.Trademark" class="qr" /></template>
                                </el-popover>
                                <span>{{ scope.row.CompanyName }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Investment" label="投资比例" width="100px" />
                    <el-table-column prop="Amount" label="投资金额" width="160px" />
                    <el-table-column prop="RegStatus" label="状态" width="100px">
                        <template #default="scope">
                            <el-tag
                                :type="scope.row.RegStatus === '存续' || scope.row.RegStatus === 'ok' ? 'success' : 'danger'">
                                {{ scope.row.RegStatus }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="域名" align="center">
                        <template #default="scope">
                            <div class="finger-container"
                                v-if="scope.row.Domains != null && scope.row.Domains.length > 0">
                                <el-tag v-for="domain in scope.row.Domains" :key="domain">{{ domain
                                }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>
            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="allOfficialAccounts" border>
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="BelongCompany" label="所属单位" />
                    <el-table-column prop="Name" label="名称">
                        <template #default="scope">
                            <div class="flex">
                                <img :src="scope.row.Logo" class="avatar mr-5px" />
                                <span>{{ scope.row.Name }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Numbers" label="微信号" />
                    <el-table-column label="二维码" align="center" width="100px">
                        <template #default="scope">
                            <el-popover placement="left" :width="180">
                                <template #reference>
                                    <wechatIcon />
                                </template>
                                <template #default><img :src="scope.row.Qrcode" class="qr" /></template>
                            </el-popover>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Introduction" label="简介" show-overflow-tooltip />
                </el-table>
            </el-tab-pane>

            <el-tab-pane label="小程序" name="applet">
                <el-table :data="allApplets" border>
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="BelongCompany" label="所属单位" />
                    <el-table-column prop="serviceName" label="小程序名称" />
                    <el-table-column prop="serviceLicence" label="备案号" />
                    <el-table-column prop="updateRecordTime" label="更新日期" />
                </el-table>
            </el-tab-pane>

            <el-tab-pane label="APP" name="app">
                <el-table :data="allApps" border>
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="BelongCompany" label="所属单位" />
                    <el-table-column prop="serviceName" label="APP名称" />
                    <el-table-column prop="serviceLicence" label="备案号" />
                    <!-- <el-table-column prop="iconUrl" label="图标">
                        <template #default="scope"><img :src="scope.row.iconUrl" class="avatar" /></template>
                    </el-table-column> -->
                </el-table>
            </el-tab-pane>

        </el-tabs>
        <template #ctrl>
            <el-space style="margin-right: -5px;">
                <el-button :icon="Share" @click="OpenFolder(from.configPath)">结果保存</el-button>
                <el-button :icon="Plus" @click="from.newTask = true" v-if="!from.runningStatus">新建任务</el-button>
                <el-button loading v-else>正在查询</el-button>
            </el-space>
        </template>
    </CustomTabs>
    <!-- 
<el-button :icon="exportIcon" class="mr-5px"
    @click="ExportAssetToXlsx(transformArrayFields(pc.table.result), pw.table.result)">
    结果导出
</el-button> -->

    <el-drawer v-model="from.newTask" title size="50%">
        <template #header>
            <span class="drawer-title">新建任务</span>
        </template>
        <el-form ref="ruleFormRef" :model="ruleForm" :rules="rules" label-width="120px">
            <el-form-item label="公司名称" prop="company">
                <el-input v-model="ruleForm.company" type="textarea" :rows="5"></el-input>
                <span class="form-item-tips">调用天眼查时尽量减少目标数与层级</span>
            </el-form-item>
            <el-form-item label="对外投资比例">
                <el-input-number v-model="ruleForm.invest" :min="1" :max="100"
                    controls-position="right"></el-input-number>
            </el-form-item>
            <el-form-item label="子公司层级:">
                <el-input-number v-model="ruleForm.subcompanyLevel" :min="1" :max="2"
                    controls-position="right"></el-input-number>
            </el-form-item>
            <el-form-item label="数据来源:">
                <el-checkbox-group v-model="ruleForm.dataSource">
                    <!-- <el-checkbox value="riskbird">
                        <el-text><el-image src="/riskbird.ico" class="mr-5px" style="width: 16px;" />
                            风鸟</el-text>
                    </el-checkbox> -->
                    <el-checkbox value="tyc">
                        <el-tooltip content="对外投资/公众号">
                            <el-text><el-image src="/tyc.png" class="mr-5px" style="width: 16px;" />
                                天眼查</el-text>
                        </el-tooltip>
                    </el-checkbox>
                    <el-checkbox value="miit" disabled>
                        <el-tooltip content="小程序/App/域名, 单独使用时需要输入完整公司名称">
                            <el-text><el-image src="/icp.ico" class="mr-5px" style="width: 16px;" />
                                工信部</el-text>
                        </el-tooltip>
                    </el-checkbox>
                </el-checkbox-group>
            </el-form-item>
            <el-form-item label="天眼查Token:" prop="tycToken" v-show="ruleForm.dataSource.includes('tyc')">
                <el-input v-model="ruleForm.tycToken" type="textarea" :rows="4" @input="storageConfig"></el-input>
            </el-form-item>
            <el-form-item label="天眼查Id:" prop="tycId" v-show="ruleForm.dataSource.includes('tyc')">
                <el-input v-model="ruleForm.tycId" @input="storageConfig"></el-input>
            </el-form-item>
            <el-form-item label="工信部API:" prop="miitApi">
                <el-input v-model="ruleForm.miitApi" @input="storageConfig">
                    <template #suffix>
                        <el-button :icon="ChromeFilled" link
                            @click="BrowserOpenURL('https://github.com/HG-ha/ICP_Query')" />
                    </template>
                </el-input>
                <span class="form-item-tips">根据ICP_Query项目搭建自身接口</span>
            </el-form-item>
            <el-form-item class="align-right">
                <el-button type="primary" @click="Collect">开始任务</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
    <el-dialog v-model="from.humanCheckDialogVisible" title="人机校验提醒" width="400px">
        <p>{{ from.humanCheckMessage }}</p>
        <p>人机校验次数请勿在短时间内重复多次, 大于3次时需要关闭程序防止账号封禁</p>
        <template #footer>
            <el-button type="primary" @click="handleHumanCheckConfirmed">我已完成验证</el-button>
        </template>
    </el-dialog>
</template>

<style scoped>
.company-cell {
    display: inline-flex;
    /* 只占内容宽度 */
    align-items: center;
    /* 垂直居中 */
    vertical-align: middle;
}

.avatar {
    width: 20px;
    height: 20px;
    /* border-radius: 50%; */
}

.qr {
    width: 150px;
    height: 150px;
}
</style>