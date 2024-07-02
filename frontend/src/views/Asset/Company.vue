<script lang="ts" setup>
import { reactive } from "vue";
import { QuestionFilled, Search } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, InitTycHeader, Callgologger } from "../../../wailsjs/go/main/App";
import { ExportAssetToXlsx } from '../../export'
import { onMounted } from 'vue';
import { CompanyInfo, WechatInfo } from "../../interface";

const from = reactive({
    subcompany: false,
    wechat: false,

    company: '',
    defaultHold: 100,
    token: '',
    activeName: 'subcompany',

    companys: [] as string[],
    runningStatus: false,
    domains: [] as string[],
})

const table = reactive({
    company: [] as CompanyInfo[],
    wehcat: [] as WechatInfo[]
})

function Collect() {
    if (from.company === "") {
        ElMessage({
            showClose: true,
            message: '查询目标不能为空',
            type: 'warning',
        })
        return
    }
    if (from.subcompany === false && from.wechat === false) {
        ElMessage({
            showClose: true,
            message: '请选择你需要查询的条件',
            type: 'warning',
        })
        return
    }
    const lines = from.company.split(/[(\r\n)\r\n]+/);
    from.companys = lines.map(line => line.trim().replace(/\s+/g, ''));
    if (from.token !== "") {
        InitTycHeader(from.token)
    } else {
        ElMessage({
            showClose: true,
            message: '天眼查Token为空，大概率会影响爬取结果，请先填写Token信息',
            type: 'error',
        })
        return
    }
    table.company = []
    table.wehcat = []
    if (from.subcompany) {
        const promises = from.companys.map(async companyName => {
            Callgologger("info", `正在收集${companyName}的子公司信息`)
            if (typeof companyName === 'string') {
                const result: CompanyInfo[] = await SubsidiariesAndDomains(companyName, from.defaultHold);
                if (result.length > 0) {
                    for (const item of result) {
                        table.company.push({
                            CompanyName: item.CompanyName,
                            Holding: item.Holding,
                            Investment: item.Investment,
                            Domains: item.Domains,
                        })
                    }
                }
            }
        });
        Promise.all(promises).then(() => {
            Callgologger("info", "已完成子公司查询任务")
        });
    }
    if (from.wechat) {
        const promises = from.companys.map(async companyName => {
            Callgologger("info", `正在收集${companyName}的微信公众号资产`)
            if (typeof companyName === 'string') {
                const result: WechatInfo[] = await WechatOfficial(companyName);
                if (result.length > 0) {
                    for (const item of result) {
                        table.wehcat.push({
                            CompanyName: item.CompanyName,
                            WechatName: item.WechatName,
                            WechatNums: item.WechatNums,
                            Qrcode: item.Qrcode,
                            Logo: item.Logo,
                            Introduction: item.Introduction,
                        });
                    }
                }
            }
        });
        Promise.all(promises).then(() => {
            Callgologger("info", "已完成微信公众查询任务")
        });
    }
}

const dataHandle = ({
    companyInfo: function (ci: CompanyInfo[]) {
        return ci.map(item => ({
            CompanyName: item.CompanyName,
            Holding: item.Holding,
            Investment: item.Investment,
            Domains: item.Domains?.map(it => JSON.stringify(it)).join('|')
        }))
    }
})
</script>

<template>
    <el-form :model="from" label-width="30%">
        <el-form-item label="公司名称:">
            <el-input v-model="from.company" type="textarea" :rows="1" style="width: 50%;"></el-input>
            <el-button type="primary" :icon="Search" @click="Collect" style="margin-left: 10px;"
                v-if="!from.runningStatus">开始查询</el-button>
            <el-button type="danger" loading style="margin-left: 10px;" v-else>正在查询</el-button>
        </el-form-item>
        <el-form-item label="查询条件:">
            <div class="flex-box">
                <el-checkbox v-model="from.subcompany" label="子公司(控股率>=)并反查域名" />
                <el-input-number v-model="from.defaultHold" size="small" :min="1" :max="100" style="margin-left: 10px;"
                    v-if="from.subcompany"></el-input-number>
            </div>
            <el-checkbox v-model="from.wechat" label="公众号" style="margin-right: 20px; margin-left: 20px;" />
        </el-form-item>
        <el-form-item>
            <template #label>
                Token:
                <el-tooltip placement="right">
                    <template #content>由于天眼查登录校验机制，为了确保爬取数据准确<br />
                        需要在此处填入网页登录后Cookie头中auth_token字段</template>
                    <el-icon>
                        <QuestionFilled size="24" />
                    </el-icon>
                </el-tooltip>
            </template>
            <el-input v-model="from.token" style="width: 50%;"></el-input>
        </el-form-item>
    </el-form>
    <div style="position: relative;">
        <el-tabs v-model="from.activeName" type="card">
            <el-tab-pane label="控股企业" name="subcompany">
                <el-table :data="table.company" height="65vh" border>
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="CompanyName" label="公司名称" :show-overflow-tooltip="true" />
                    <el-table-column prop="Holding" width="100px" label="股权比例" />
                    <el-table-column prop="Investment" width="160px" label="投资数额" />
                    <el-table-column prop="Domains" label="域名">
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
            </el-tab-pane>

            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="table.wehcat" height="65vh" border :cell-style="{ height: '23px' }">
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
                                    <svg t="1712903814884" viewBox="0 0 1024 1024" version="1.1"
                                        xmlns="http://www.w3.org/2000/svg" p-id="2765" width="1.5em" height="1.5em">
                                        <path
                                            d="M468 128H160c-17.7 0-32 14.3-32 32v308c0 4.4 3.6 8 8 8h332c4.4 0 8-3.6 8-8V136c0-4.4-3.6-8-8-8z m-56 284H192V192h220v220z m-138-74h56c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8h-56c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8z m194 210H136c-4.4 0-8 3.6-8 8v308c0 17.7 14.3 32 32 32h308c4.4 0 8-3.6 8-8V556c0-4.4-3.6-8-8-8z m-56 284H192V612h220v220z m-138-74h56c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8h-56c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8z m590-630H556c-4.4 0-8 3.6-8 8v332c0 4.4 3.6 8 8 8h332c4.4 0 8-3.6 8-8V160c0-17.7-14.3-32-32-32z m-32 284H612V192h220v220z m-138-74h56c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8h-56c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8z m194 210h-48c-4.4 0-8 3.6-8 8v134h-78V556c0-4.4-3.6-8-8-8H556c-4.4 0-8 3.6-8 8v332c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8V644h78v102c0 4.4 3.6 8 8 8h190c4.4 0 8-3.6 8-8V556c0-4.4-3.6-8-8-8zM746 832h-48c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8z m142 0h-48c-4.4 0-8 3.6-8 8v48c0 4.4 3.6 8 8 8h48c4.4 0 8-3.6 8-8v-48c0-4.4-3.6-8-8-8z"
                                            fill="#1296DB" p-id="2766"></path>
                                    </svg>
                                </template>
                                <template #default>
                                    <img :src="scope.row.Qrcode" style="width: 150px; height: 150px">
                                </template>
                            </el-popover>
                        </template>
                    </el-table-column>
                    <el-table-column prop="Introduction" label="简介" :show-overflow-tooltip="true" />
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-button @click="ExportAssetToXlsx(dataHandle.companyInfo(table.company), table.wehcat)">
                <template #icon>
                    <img src="/excel.svg" width="16">
                </template>导出Excel</el-button>
        </div>
    </div>

</template>

<style scoped>
.finger-container {
    flex-wrap: wrap;
    display: flex;
    gap: 7px;
}
</style>