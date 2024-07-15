<script lang="ts" setup>
import { reactive } from "vue";
import { QuestionFilled, Plus } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, InitTycHeader, Callgologger } from "../../../wailsjs/go/main/App";
import { ExportAssetToXlsx } from '../../export'
import { CompanyInfo, WechatInfo } from "../../interface";
import usePagination from "../../usePagination";
import { transformArrayFields } from "../../util";

const from = reactive({
    newTask: false,
    domain: true,
    wechat: true,
    company: '',
    defaultHold: 100,
    token: '',
    activeName: 'subcompany',
    runningStatus: false,
    companyData: [] as CompanyInfo[],
    wehcatData: [] as WechatInfo[]
})

let pc = usePagination(from.companyData, 20) // paginationCompany
let pw = usePagination(from.wehcatData, 20) // paginationWehcat

function Collect() {
    from.newTask = false
    from.runningStatus = true
    if (from.company === "") {
        ElMessage({
            showClose: true,
            message: '查询目标不能为空',
            type: 'warning',
        })
        return
    }
    if (from.token == "") {
        ElMessage({
            showClose: true,
            message: '天眼查Token为空，大概率会影响爬取结果，请先填写Token信息',
            type: 'error',
        })
        return
    } else {
        from.token = from.token.replace(/[\r\n\s]/g, '')
    }
    const lines = from.company.split(/[(\r\n)\r\n]+/);
    let companys = lines.map(line => line.trim().replace(/\s+/g, ''));
    InitTycHeader(from.token)
    pc.ctrl.initTable()
    pw.ctrl.initTable()
    let allCompany = [] as string[]
    const promises = companys.map(async companyName => {
        Callgologger("info", `正在收集${companyName}的子公司信息`)
        if (typeof companyName === 'string') {
            const result: CompanyInfo[] = await SubsidiariesAndDomains(companyName, from.defaultHold);
            if (result.length > 0) {
                for (const item of result) {
                    pc.table.result.push({
                        CompanyName: item.CompanyName,
                        Holding: item.Holding,
                        Investment: item.Investment,
                        RegStatus: item.RegStatus,
                        Domains: item.Domains,
                    })
                    allCompany.push(item.CompanyName!)
                    pc.table.pageContent = pc.ctrl.watchResultChange(pc.table.result, pc.table.currentPage, pc.table.pageSize)
                }
            }
        }
    });
    Promise.all(promises).then(() => {
        Callgologger("info", "已完成子公司查询任务")
        if (from.wechat) {
            const promises2 = allCompany.map(async companyName => {
                Callgologger("info", `正在收集${companyName}的微信公众号资产`)
                if (typeof companyName === 'string') {
                    const result: WechatInfo[] = await WechatOfficial(companyName);
                    if (Array.isArray(result) && result.length > 0) {
                        for (const item of result) {
                            pw.table.result.push({
                                CompanyName: item.CompanyName,
                                WechatName: item.WechatName,
                                WechatNums: item.WechatNums,
                                Qrcode: item.Qrcode,
                                Logo: item.Logo,
                                Introduction: item.Introduction,
                            });
                            pw.table.pageContent = pw.ctrl.watchResultChange(pw.table.result, pw.table.currentPage, pw.table.pageSize)
                        }
                    }
                }
            });
            Promise.all(promises2).then(() => {
                Callgologger("info", "已完成微信公众查询任务")
                from.runningStatus = false
            });
        } else {
            from.runningStatus = false
        }
    });

}

</script>


<template>
    <el-drawer v-model="from.newTask" direction="rtl" size="40%">
        <template #header>
            <h4>新建任务</h4>
        </template>
        <el-form :model="from" label-width="auto">
            <el-form-item label="公司名称:">
                <el-input v-model="from.company" type="textarea" :rows="5"></el-input>
            </el-form-item>
            <el-form-item>
                <template #label>
                    股权比例:
                    <el-tooltip placement="right" content="会自动反查域名">
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="from.defaultHold" :min="1" :max="100"></el-input-number>
            </el-form-item>
            <el-form-item label="其他查询内容:">
                <el-checkbox v-model="from.wechat" label="公众号" />
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
                <el-input v-model="from.token" type="textarea" rows="4"></el-input>
            </el-form-item>
            <el-form-item class="align-right">
                <el-button type="primary" @click="Collect">开始查询</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
    <div style="position: relative;">
        <el-tabs v-model="from.activeName">
            <el-tab-pane label="控股企业" name="subcompany">
                <el-table :data="pc.table.pageContent" height="80vh" border>
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="CompanyName" label="公司名称" :show-overflow-tooltip="true" />
                    <el-table-column prop="Holding" width="100px" label="股权比例" />
                    <el-table-column prop="Investment" width="160px" label="投资数额" />
                    <el-table-column prop="RegStatus" width="100px" label="企业状态" align="center">
                        <template #default="scope">
                            <el-tag v-if="scope.row.RegStatus === '存续'" type="success">{{ scope.row.RegStatus
                                }}</el-tag>
                            <el-tag v-else-if="scope.row.RegStatus === '注销'" type="danger">{{ scope.row.RegStatus
                                }}</el-tag>
                        </template>
                    </el-table-column>
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
                <div class="my-header" style="margin-top: 5px;">
                    <div></div>
                    <el-pagination background @size-change="pc.ctrl.handleSizeChange"
                        @current-change="pc.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pc.table.currentPage" :page-sizes="[20, 50, 100]" :page-size="pc.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="pc.table.result.length">
                    </el-pagination>
                </div>

            </el-tab-pane>
            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="pw.table.pageContent" height="80vh" border :cell-style="{ height: '23px' }">
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
                <div class="my-header" style="margin-top: 5px;">
                    <div></div>
                    <el-pagination background @size-change="pw.ctrl.handleSizeChange"
                        @current-change="pw.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pw.table.currentPage" :page-sizes="[20, 50, 100]" :page-size="pw.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="pw.table.result.length">
                    </el-pagination>
                </div>

            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-button type="primary" :icon="Plus" @click="from.newTask = true"
                v-if="!from.runningStatus">新建任务</el-button>
            <el-button type="primary" loading v-else>正在查询</el-button>
            <el-button style="margin-left: 5px;" color="#626aef"
                @click="ExportAssetToXlsx(transformArrayFields(pc.table.result), pw.table.result)">导出</el-button>
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