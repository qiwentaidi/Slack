<script lang="ts" setup>
import { onMounted, reactive, ref } from "vue";
import { ChromeFilled, ArrowUpBold, ArrowDownBold, DocumentCopy } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, TycCheckLogin, Callgologger } from "wailsjs/go/services/App";
import { ExportAssetToXlsx } from '@/export'
import usePagination from "@/usePagination";
import { Copy, getRootDomain, transformArrayFields } from "@/util";
import exportIcon from '@/assets/icon/doucment-export.svg'
import throttle from 'lodash/throttle';
import CustomTabs from "@/components/CustomTabs.vue";
import wechatIcon from "@/assets/icon/wechatOfficialAccount.svg"
import { debounce } from "lodash"
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import { structs } from "wailsjs/go/models";

const throttleUpdate = throttle(() => {
    pc.ctrl.watchResultChange(pc.table);
}, 1000);


onMounted(() => {
    const storedToken = localStorage.getItem('tyc-token');
    if (storedToken) {
        from.token = storedToken;
    }
})

const from = reactive({
    newTask: false,
    wechat: true,
    domain: false,
    company: '',
    defaultHold: 100,
    subcompanyLevel: 1,
    token: '',
    activeName: 'subcompany',
    runningStatus: false,
    machineStr: '',
    errorCompanies: [] as string[]
})

let pc = usePagination<structs.CompanyInfo>(20) // paginationCompany
let pw = usePagination<structs.WechatReulst>(20) // paginationWehcat

async function Collect() {
    if (from.company == "") {
        ElMessage.warning('查询目标不能为空')
        return
    }
    if (from.token == "") {
        ElMessage.warning('天眼查Token为空，大概率会影响爬取结果，请先填写Token信息')
        return
    } else {
        from.token = from.token.replace(/[\r\n\s]/g, '')
        let isLogin = await TycCheckLogin(from.token)
        if (!isLogin) {
            ElMessage.warning('天眼查Token已失效')
            return
        }
    }

    if (from.domain && from.machineStr == "") {
        ElMessage.warning('MachineStr为空, 无法进行子域名查询, 请先配置该内容')
        return
    } else {
        from.machineStr = from.machineStr.replace(/[\r\n\s]/g, '')
    }

    from.newTask = false
    from.runningStatus = true
    showForm.value = false
    const lines = from.company.split(/[(\r\n)\r\n]+/);
    let companys = lines.map(line => line.trim().replace(/\s+/g, ''));
    pc.initTable()
    pw.initTable()
    from.errorCompanies = []
    let allCompany = [] as string[]
    // 1. 收集子公司信息
    for (let i = 0; i < companys.length; i++) {
        const companyName = companys[i];
        Callgologger("info", `正在收集${companyName}的子公司信息`);

        if (typeof companyName === 'string') {
            const result = await SubsidiariesAndDomains(companyName, from.subcompanyLevel, from.defaultHold, from.domain, from.machineStr);

            // 查询失败了，可能是封号/人机校验
            if (!result) {
                from.errorCompanies.push(companyName);

                // 将剩余未查询的公司全部加入 errorCompanies
                from.errorCompanies.push(...companys.slice(i + 1));

                Callgologger("error", "触发人机校验或封号，终止任务");
                from.runningStatus = false;
                showForm.value = true;
                return;
            }

            if (result.length > 0) {
                pc.table.result.push(...result);
                throttleUpdate();
                for (const item of result) {
                    allCompany.push(item.CompanyName!);
                }
            }
        }
    }

    Callgologger("info", "已完成子公司查询任务")

    // 3. 收集微信公众号信息
    if (from.wechat) {
        for (const companyName of allCompany) {
            if (companyName == "") return
            Callgologger("info", `正在收集${companyName}的微信公众号资产`)
            if (typeof companyName === 'string') {
                const result = await WechatOfficial(companyName);
                if (Array.isArray(result) && result.length > 0) {
                    pw.table.result.push(...result);
                    pw.ctrl.watchResultChange(pw.table)
                }
            }
        }
        Callgologger("info", "已完成微信公众查询任务")
    }

    // 4. 完成所有任务
    from.runningStatus = false
    showForm.value = true
}

const debouncedInput = debounce(() => {
    // 2秒后将数据存储到localStorage
    localStorage.setItem('tyc-token', from.token);
}, 2000);

function copyAllDomains() {
    const allDomains = pc.table.result
        .flatMap(item => item.Domains)
        .map(getRootDomain);
    Copy(allDomains.join('\n'));
}

const showForm = ref(true);

function toggleFormVisibility() {
    showForm.value = !showForm.value;
}
</script>


<template>
    <el-divider>
        <el-button round :icon="showForm ? ArrowUpBold : ArrowDownBold" @click="toggleFormVisibility"
            v-if="!from.runningStatus">
            {{ showForm ? '隐藏参数' : '展开参数' }}
        </el-button>
        <el-button round loading v-else>正在运行</el-button>
    </el-divider>
    <el-collapse-transition>
        <div style="display: flex; gap: 10px;" v-show="showForm">
            <el-form :model="from" label-width="auto" style="width: 50%;">
                <el-form-item label="公司名称:">
                    <el-input v-model="from.company" type="textarea" :rows="5"></el-input>
                </el-form-item>
                <el-form-item label="股权比例:">
                    <el-input-number v-model="from.defaultHold" :min="1" :max="100"></el-input-number>
                </el-form-item>
                <el-form-item label="子公司层级:">
                    <el-input-number v-model="from.subcompanyLevel" :min="1" :max="3"></el-input-number>
                </el-form-item>
                <el-form-item label="其他查询内容:">
                    <el-checkbox v-model="from.wechat" label="公众号" />
                    <el-checkbox v-model="from.domain" label="域名查询" />
                </el-form-item>
                <el-form-item label="Token:">
                    <el-input v-model="from.token" type="textarea" :rows="4" @input="debouncedInput"></el-input>
                    <span class="form-item-tips">填写TYC网页登录后Cookie头中auth_token字段</span>
                </el-form-item>
                <el-form-item label="MachineStr:">
                    <el-input v-model="from.machineStr">
                        <template #suffix>
                            <el-button :icon="ChromeFilled" link @click="BrowserOpenURL('https://www.beianx.cn/')" />
                        </template>
                    </el-input>
                    <span class="form-item-tips">
                        填写www.beianx.cn Cookie头中machine_str字段的值, 如果没有的话先进行一次查询
                    </span>
                </el-form-item>
                <el-form-item class="align-right">
                    <el-button type="primary" @click="Collect">开始任务</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-collapse-transition>
    <CustomTabs>
        <el-tabs v-model="from.activeName" type="border-card">
            <el-tab-pane label="控股企业" name="subcompany">
                <el-table :data="pc.table.pageContent" style="height: calc(100vh - 175px);">
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="CompanyName" label="公司名称" :show-overflow-tooltip="true" />
                    <el-table-column prop="Holding" width="100px" label="股权比例" />
                    <el-table-column prop="Investment" width="160px" label="投资数额" />
                    <el-table-column prop="RegStatus" width="100px" label="企业状态" align="center">
                        <template #default="scope">
                            <el-tag v-if="scope.row.RegStatus === '存续' || scope.row.RegStatus === 'ok'"
                                type="success">{{ scope.row.RegStatus
                                }}</el-tag>
                            <el-tag v-else type="danger">{{ scope.row.RegStatus
                                }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Domains">
                        <template #header>
                            <span class="position-center">域名
                                <el-divider direction="vertical" />
                                <el-button size="small" text bg @click="copyAllDomains()">复制全部根域名</el-button>
                            </span>
                        </template>
                        <template #default="scope">
                            <div class="finger-container" v-if="scope.row.Domains.length > 0">
                                <el-tag v-for="domain in scope.row.Domains" :key="domain">{{ domain
                                    }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <div></div>
                    <el-pagination size="small" background @size-change="pc.ctrl.handleSizeChange"
                        @current-change="pc.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pc.table.currentPage" :page-sizes="[20, 50, 100]" :page-size="pc.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="pc.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="pw.table.pageContent" style="height: calc(100vh - 175px);"
                    :cell-style="{ height: '23px' }">
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="CompanyName" label="公司名称" width="180px" />
                    <el-table-column prop="WechatName" label="公众号名称">
                        <template #default="scope">
                            <div class="flex-box">
                                <img :src="scope.row.Logo" style="width: 25px; height: 25px; margin-right: 10px;">
                                <span>{{ scope.row.WechatName }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="WechatNums" label="微信号" :show-overflow-tooltip="true" />
                    <el-table-column prop="Qrcode" width="80px" label="二维码" align="center">
                        <template #default="scope">
                            <el-popover :width="180" placement="left">
                                <template #reference>
                                    <wechatIcon />
                                </template>
                                <template #default>
                                    <img :src="scope.row.Qrcode" style="width: 150px; height: 150px">
                                </template>
                            </el-popover>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Introduction" label="简介" />
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="my-header" style="margin-top: 5px;">
                    <div></div>
                    <el-pagination size="small" background @size-change="pw.ctrl.handleSizeChange"
                        @current-change="pw.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pw.table.currentPage" :page-sizes="[20, 50, 100]" :page-size="pw.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="pw.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-popover title="失败列表" :width="300" trigger="click">
                <template #reference>
                    <el-button>失败列表: {{ from.errorCompanies.length }}</el-button>
                </template>
                <el-scrollbar height="150px">
                    <p v-for="u in from.errorCompanies" class="scrollbar-demo-item">
                        {{ u }}</p>
                </el-scrollbar>
                <el-button :icon="DocumentCopy" @click="Copy(from.errorCompanies.join('\n'))"
                    style="width: 100%;">复制全部失败目标</el-button>
            </el-popover>
            <el-button :icon="exportIcon" style="margin-left: 5px;"
                @click="ExportAssetToXlsx(transformArrayFields(pc.table.result), pw.table.result)">
                结果导出
            </el-button>
        </template>
    </CustomTabs>

</template>

<style scoped>
.finger-container {
    flex-wrap: wrap;
    display: flex;
    gap: 7px;
}

.scrollbar-demo-item {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 50px;
    margin-block: 10px;
    text-align: center;
    border-radius: 4px;
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>