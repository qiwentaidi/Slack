<script lang="ts" setup>
import { onMounted, reactive, ref } from "vue";
import { ChromeFilled, ArrowUpBold, ArrowDownBold } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, TycCheckLogin, Callgologger } from "wailsjs/go/services/App";
import { ExportAssetToXlsx } from '@/export'
import usePagination from "@/usePagination";
import { Copy, getRootDomain, transformArrayFields } from "@/util";
import exportIcon from '@/assets/icon/doucment-export.svg'
import global from "@/stores";
import { LinkSubdomain } from "@/linkage";
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
    linkSubdomain: false,
    company: '',
    defaultHold: 100,
    subcompanyLevel: 1,
    token: '',
    activeName: 'subcompany',
    runningStatus: false,
    machineStr: ''
})

let pc = usePagination<structs.CompanyInfo>(20) // paginationCompany
let pw = usePagination<structs.WechatReulst>(20) // paginationWehcat

function NewRunner() {
   let collect = new Collect()
   collect.Runner()
}

class Collect {
    allCompany: string[] = []
    allSubdomain: string[] = []
    public async Runner() {

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
            ElMessage.warning('MachineStr为空，无法进行子域名查询，请先配置该内容')
            return
        } else {
            from.machineStr = from.machineStr.replace(/[\r\n\s]/g, '')
        }
        if (from.linkSubdomain && global.space.bevigil == "" && global.space.chaos == "" && global.space.zoomeye == "" && global.space.securitytrails == "" && global.space.github == "") {
            ElMessage.warning('未配置任何域名收集模块API，请在设置中至少配置一个')
            return
        }
        from.newTask = false
        from.runningStatus = true
        showForm.value = false
        const lines = from.company.split(/[(\r\n)\r\n]+/);
        this.allCompany = lines.map(line => line.trim().replace(/\s+/g, ''));
        pc.initTable()
        pw.initTable()
        // 1. 收集子公司信息
        for (const companyName of this.allCompany) {
            await this.collectSubCompanies(companyName);
        }

        Callgologger("info", "已完成子公司查询任务")

        // 2. 处理子域名链接
        if (this.allSubdomain.length != 0) {
            await LinkSubdomain(this.allSubdomain)
        }

        // 3. 收集微信公众号信息
        if (from.wechat) {
            for (const companyName of this.allCompany) {
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
        this.finishTask()
    }

    public async collectSubCompanies(companyName: string) {
        Callgologger("info", `正在收集${companyName}的子公司信息`);

        if (typeof companyName === 'string') {
            const result = await SubsidiariesAndDomains(companyName, from.subcompanyLevel, from.defaultHold, from.domain, from.machineStr);

            if (result.length > 0) {
                const hasError = result.some(item => item.Investment === "error");
                if (hasError) {
                    const confirm = await showConfirm(companyName);
                    if (!confirm) {
                        this.finishTask();
                        return;
                    }
                    // 如果确认继续，递归重试
                    await this.collectSubCompanies(companyName);
                    return;
                }
                pc.table.result.push(...result);
                throttleUpdate();
                for (const item of result) {
                    this.allCompany.push(item.CompanyName!);
                    if (from.linkSubdomain && item.Domains!.length > 0) {
                        this.allSubdomain.push(...item.Domains!);
                    }
                }
            }
        }
    }
    private finishTask() {
        from.runningStatus = false
        showForm.value = true
    }
}

async function showConfirm(companyName: string) {
    return ElMessageBox.confirm(
        `检测到 ${companyName} 查询数据异常, 疑似触发天眼查人机校验, 请打开网站手动校验! 确认校验通过后点击确定按钮继续任务. 点击取消按钮则会退出任务.`,
        '警告',
        {
            type: 'warning',
            closeOnClickModal: false,
            closeOnPressEscape: false,
        }
    )
        .then(() => {
            return true; // 返回 true，表示确认
        })
        .catch(() => {
            ElMessage({
                type: 'error',
                message: '用户已取消, 人机校验未通过, 任务已退出.',
            })
            return false; // 返回 false，表示取消
        })
}

function recheckLinkSubdomain() {
    if (!from.domain) { from.linkSubdomain = false }
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
        <el-button round :icon="showForm ? ArrowUpBold : ArrowDownBold" @click="toggleFormVisibility">
            {{ showForm ? '隐藏参数' : '展开参数' }}
        </el-button>
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
                    <el-checkbox v-model="from.domain" label="域名查询" @change="recheckLinkSubdomain" />
                    <!-- <el-checkbox v-model="from.linkSubdomain" label="联动子域名收集" :disabled="!from.domain"/> -->
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
                    <el-button type="primary" @click="NewRunner()" v-if="!from.runningStatus">开始任务</el-button>
                    <el-button loading v-else>正在查询</el-button>
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
</style>