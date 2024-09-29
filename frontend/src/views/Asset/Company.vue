<script lang="ts" setup>
import { onMounted, reactive } from "vue";
import { QuestionFilled, Plus } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, InitTycHeader, Callgologger, TycCheckLogin } from "wailsjs/go/main/App";
import { ExportAssetToXlsx } from '@/export'
import { CompanyInfo, WechatInfo } from "@/interface";
import usePagination from "@/usePagination";
import { transformArrayFields } from "@/util";
import exportIcon from '@/assets/icon/doucment-export.svg'
import global from "@/global";
import { LinkSubdomain } from "@/linkage";
import CustomTabs from "@/components/CustomTabs.vue";
import wechatIcon from "@/assets/icon/wechat.svg"
import { debounce } from "lodash"

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
    companyData: [] as CompanyInfo[],
    wehcatData: [] as WechatInfo[],
    machineStr: ''
})

let pc = usePagination(from.companyData, 20) // paginationCompany
let pw = usePagination(from.wehcatData, 20) // paginationWehcat

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
        InitTycHeader(from.token)
        let isLogin = await TycCheckLogin()
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
    const lines = from.company.split(/[(\r\n)\r\n]+/);
    let companys = lines.map(line => line.trim().replace(/\s+/g, ''));    
    pc.ctrl.initTable()
    pw.ctrl.initTable()
    let allCompany = [] as string[]
    let allSubdomain = [] as string[]
    const promises = companys.map(async companyName => {
        Callgologger("info", `正在收集${companyName}的子公司信息`)
        if (typeof companyName === 'string') {
            const result: CompanyInfo[] = await SubsidiariesAndDomains(companyName, from.subcompanyLevel, from.defaultHold, from.domain, from.machineStr);
            if (result.length > 0) {
                pc.table.result.push(...result)
                pc.table.pageContent = pc.ctrl.watchResultChange(pc.table)
                for (const item of result) {
                    allCompany.push(item.CompanyName!)
                    if (from.linkSubdomain && item.Domains!.length > 0) {
                        allSubdomain.push(...item.Domains!)
                    }
                }
            }
        }
    });
    Promise.all(promises).then(() => {
        Callgologger("info", "已完成子公司查询任务")
        if (allSubdomain.length != 0) {
            LinkSubdomain(allSubdomain)
        }
        if (from.wechat) {
            const promises2 = allCompany.map(async companyName => {
                Callgologger("info", `正在收集${companyName}的微信公众号资产`)
                if (typeof companyName === 'string') {
                    const result: WechatInfo[] = await WechatOfficial(companyName);
                    if (Array.isArray(result) && result.length > 0) {
                        pw.table.result.push(...result);
                        pw.table.pageContent = pw.ctrl.watchResultChange(pw.table)
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

function recheckLinkSubdomain() {
    if (!from.domain) { from.linkSubdomain = false}
}

const debouncedInput = debounce(() => {
    // 2秒后将数据存储到localStorage
    localStorage.setItem('tyc-token', from.token);
}, 2000);
</script>


<template>
    <el-drawer v-model="from.newTask" direction="rtl" size="40%">
        <template #header>
            <span class="drawer-title">新建任务</span>
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
            <el-form-item label="子公司层级:">
                <el-input-number v-model="from.subcompanyLevel" :min="1" :max="3"></el-input-number>
            </el-form-item>
            <el-form-item label="其他查询内容:">
                <el-checkbox v-model="from.wechat" label="公众号" />
                <el-checkbox v-model="from.domain" label="查询子域名" @change="recheckLinkSubdomain" />
                <el-checkbox v-model="from.linkSubdomain" label="联动子域名收集" :disabled="!from.domain"/>
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
                <el-input v-model="from.token" type="textarea" :rows="4" @input="debouncedInput"></el-input>
            </el-form-item>
            <el-form-item>
                <template #label>
                    MachineStr:
                    <el-tooltip placement="right">
                        <template #content>由于https://www.beianx.cn/备案查域名新增校验机制<br />
                            需要在此处填入Cookie头中machine_str字段的值</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input v-model="from.machineStr"></el-input>
            </el-form-item>
            <el-form-item class="align-right">
                <el-button type="primary" @click="Collect">开始查询</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
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
                    <el-pagination size="small" background @size-change="pc.ctrl.handleSizeChange"
                        @current-change="pc.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pc.table.currentPage" :page-sizes="[20, 50, 100]" :page-size="pc.table.pageSize"
                        layout="total, sizes, prev, pager, next" :total="pc.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="pw.table.pageContent" style="height: calc(100vh - 175px);" :cell-style="{ height: '23px' }">
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
            <el-tooltip content="新建任务" v-if="!from.runningStatus">
                <el-button :icon="Plus" round @click="from.newTask = true">新建任务</el-button>
            </el-tooltip>
            <el-button round loading v-else>正在查询</el-button>
            <el-tooltip content="导出Excel">
                <el-button :icon="exportIcon" style="margin-left: 5px;"
                    @click="ExportAssetToXlsx(transformArrayFields(pc.table.result), pw.table.result)">
                </el-button>
            </el-tooltip>
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