<script lang="ts" setup>
import { reactive, ref } from "vue";
import { QuestionFilled, Search } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { WechatOfficial, SubsidiariesAndDomains, InitTycHeader, AssetHunter } from "../../../wailsjs/go/main/App";
import global from "../../global"
import { ExportAssetToXlsx } from '../../export'
import { onMounted } from 'vue';
import { sleep } from "../../util";
// 初始化时调用
onMounted(() => {
    Init()
});
const from = reactive({
    subcompany: false,
    wechat: false,

    company: '',
    defaultHold: 100,
    token: '',
    activeName: 'subcompany',

    companys: [{}],
    log: '',
    runningStatus: false,
    domains: [{}],
    correctName: [{}],
    getColumnData(prop: string): any[] {
        return su.value.map((item: any) => item[prop]);
    }
})

interface CompanyInfo {
    CompanyName: string
    Holding: string
    Investment: string
    Domains: string[]
}

var su = ref([] as CompanyInfo[])
const hu = ref([{}])
const we = ref([{}])

function Init() {
    su.value = [];
    hu.value = [];
    we.value = [];
}

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
    Init()
    if (from.subcompany) {
        const promises = from.companys.map(async companyName => {
            from.log += `[INF] 正在收集${companyName}的子公司信息\n`
            if (typeof companyName === 'string') {
                const result = await SubsidiariesAndDomains(companyName, from.defaultHold);
                console.log(result)
                if (Array.isArray(result.Asset) && result.Asset.length > 0) {
                    for (const item of result.Asset) {
                        su.value.push({
                            CompanyName: item.CompanyName,
                            Holding: item.Holding,
                            Investment: item.Investment,
                            Domains: item.Domains,
                        })
                    }
                }
                if (result.Prompt.length > 0) { // 处理字符串
                    from.log += `[INF] ${result.Prompt}\n`;
                }
            }
        });
        Promise.all(promises).then(() => {
            from.log += `[SUCCESS] 已完成子公司查询任务\n`;
        });
    }
    if (from.wechat) {
        const promises = from.companys.map(async companyName => {
            from.log += `[INF] 正在收集${companyName}的微信公众号资产\n`
            if (typeof companyName === 'string') {
                const result = await WechatOfficial(companyName);
                if (result.Asset.length > 0) {
                    const mappedResult = result.Asset.map((item: any) => {
                        return {
                            companyName: companyName,
                            weName: item[0],
                            weNums: item[1],
                            logo: item[2],
                            qrcode: item[3],
                            introduction: item[4],
                        };
                    });
                    we.value.push(...mappedResult);
                }
                if (result.Prompt.length > 0) { // 处理字符串
                    from.log += `[INF] ${result.Prompt}\n`;
                }
            }
        });
        Promise.all(promises).then(() => {
            from.log += `[SUCCESS] 已完成微信公众查询任务\n`;
        });
    }
}
async function huntersearch() {
    if (global.space.hunterkey.length <= 1) {
        ElMessage({
            showClose: true,
            message: '请先在设置中配置鹰图key再使用该功能',
            type: 'warning',
        })
        return
    }
    if (su.value.length <= 1) {
        ElMessage({
            showClose: true,
            message: '请先查询控股企业信息再继续资产数量查询',
            type: 'warning',
        })
        return
    }
    // 初始化要查询的目标
    from.correctName = from.getColumnData("company")
    from.domains = []
    from.getColumnData("domains").forEach(dms => {
        if (dms != "") {
            if (dms.includes("|")) {
                let elements = dms.split(" | "); // 如果字符串包含 "|" 符号，则根据 "|" 符号进行分割
                // 将分割后的元素追加到 from.waitsearch 中
                elements.forEach((element: any) => {
                    from.domains.push(element);
                });
            } else {
                from.domains.push(dms); // 如果字符串不包含 "|" 符号，则直接追加到 from.waitsearch 中
            }
        }
    });
    // 查询icp
    for (let target of from.correctName) {
        await sleep(2500);
        AssetHunter(0, target as string, global.space.hunterkey).then(
            result => {
                hu.value.push(
                    {
                        name: target,
                        hunums: result.Total,
                    }
                )
                from.log += `[INF] ${result.Info}\n`
            }
        )
    }
    // 查询domain
    for (let target of from.domains) {
        await sleep(2500);
        AssetHunter(1, target as string, global.space.hunterkey).then(
            result => {
                hu.value.push(
                    {
                        name: target,
                        hunums: result.Total,
                    }
                )
                from.log += `[INF] ${result.Info}\n`
            }
        )
    }
}
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
                <el-checkbox v-model="from.subcompany" label="查询子公司并反查域名" />
                <el-tooltip placement="left">
                    <template #content>控股率>=(%)</template>
                    <el-icon>
                        <QuestionFilled />
                    </el-icon>
                </el-tooltip>
                <el-input-number v-model="from.defaultHold" size="small" :min="1" :max="100" style="margin-left: 10px;"
                    v-if="from.subcompany"></el-input-number>
            </div>
            <el-checkbox v-model="from.wechat" label="查询公众号" style="margin-right: 20px; margin-left: 20px;" />
        </el-form-item>
        <el-form-item>
            <template #label>
                Token:
                <el-tooltip placement="right">
                    <template #content>由于天眼查登录校验机制，为了确保公众号爬取数据准确<br />
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
                <el-table :data="su" height="60vh" border>
                    <el-table-column type="index" width="60px" />
                    <el-table-column prop="CompanyName" label="公司名称" :show-overflow-tooltip="true" />
                    <el-table-column prop="Holding" label="股权比例" :show-overflow-tooltip="true" />
                    <el-table-column prop="Investment" label="投资数额" :show-overflow-tooltip="true" />
                    <el-table-column prop="Domains" label="域名">
                        <template #default="scope">
                            <div class="finger-container" v-if="scope.row.Domains.length > 0">
                                <el-tag v-for="domain in scope.row.Domains" :key="domain">{{ domain
                                    }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>

            <el-tab-pane label="公众号" name="wechat">
                <el-table :data="we" height="60vh" border :cell-style="{ height: '23px' }">
                    <el-table-column type="index" width="60px" />
                    <el-table-column prop="companyName" label="公司名称" width="180px" />
                    <el-table-column prop="weName" label="公众号名称">
                        <template #default="scope">
                            <div class="flex-box">
                                <img :src="scope.row.logo" style="width: 25px; height: 25px; margin-right: 10px;">
                                <span>{{ scope.row.weName }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="weNums" label="微信号" :show-overflow-tooltip="true" />
                    <el-table-column prop="qrcode" width="80px" label="二维码" align="center">
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
                                    <img :src="scope.row.qrcode" style="width: 150px; height: 150px">
                                </template>
                            </el-popover>
                        </template>
                    </el-table-column>
                    <el-table-column prop="introduction" label="简介" :show-overflow-tooltip="true" />
                </el-table>
            </el-tab-pane>

            <el-tab-pane label="鹰图资产数量" name="hunter">
                <el-table :data="hu" height="60vh" border>
                    <el-table-column type="index" width="60px" />
                    <el-table-column prop="name" label="公司域名或ICP名称" :show-overflow-tooltip="true" />
                    <el-table-column prop="hunums" label="资产数量"
                        :sort-method="(a: any, b: any) => { return a.hunums - b.hunums }" sortable />
                </el-table>
            </el-tab-pane>
            <el-tab-pane label="运行日志">
                <el-input class="log-textarea" v-model="from.log" type="textarea" rows="20" readonly></el-input>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar">
            <el-button-group>
                <el-button @click="huntersearch">
                    <template #icon>
                        <img src="/hunter.ico" width="16">
                    </template>
                    测绘资产数量
                    <el-popover placement="left-start" :width="350" trigger="hover">
                        ①需要控股企业查询的<b>资产数量>=1</b><br /><br />
                        ②查询公司名或者域名在鹰图中的资产数量<br /><br />
                        ③一次查询消耗1积分对应的鹰图查询语句为<b>icp.name=""</b>和<b>domain.suffix=""</b>
                        <template #reference>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </template>
                    </el-popover>
                </el-button>
                <el-button @click="ExportAssetToXlsx(su, we, hu)">
                    <template #icon>
                        <img src="/excle.svg" width="16">
                    </template>导出Excle</el-button>
            </el-button-group>
        </div>
    </div>

</template>