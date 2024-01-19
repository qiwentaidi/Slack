<script lang="ts" setup>
import { reactive, ref } from "vue";
import { QuestionFilled, Tickets } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus'
import { AssetWechat, AssetSubcompany, InitTycHeader, AssetHunter } from "../../../wailsjs/go/main/App";
import global from "../../global"
import { ExportAssetToXlsx } from '../../util'
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    su.value = [];
    hu.value = [];
    we.value = [];
});
const from = reactive({
    subcompany: false,
    wechat: false,

    company: '',
    defaultHold: 100,
    token: '',
    rows: 1,
    activeName: 'subcompany',

    companys: [{}],
    log: '',

    domains: [{}],
    correctName: [{}],
    getColumnData(prop: string): any[] {
        return su.value.map((item: any) => item[prop]);
    }
})
const su = ref([{}])
const hu = ref([{}])
const we = ref([{}])
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
    if (from.subcompany) {
        su.value = [];
        const promises = from.companys.map(async companyName => {
            from.log += `[INF] 正在收集${companyName}的子公司信息\n`
            if (typeof companyName === 'string') {
                const result = await AssetSubcompany(companyName, from.defaultHold);
                if (result.Shareholding.length > 0) {
                    const mappedResult = result.Shareholding.map((item: any) => {
                        return {
                            company: item[0],
                            ratio: item[1],
                            sunums: item[2],
                            domains: item[3],
                        };
                    });
                    su.value.push(...mappedResult);
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
        we.value = [];
        const promises = from.companys.map(async companyName => {
            from.log += `[INF] 正在收集${companyName}的微信公众号资产\n`
            if (typeof companyName === 'string') {
                const result = await AssetWechat(companyName);
                if (result.OfficialAccounts.length > 0) {
                    const mappedResult = result.OfficialAccounts.map((item: any) => {
                        return {
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
    hu.value = []
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
        await sleep(2000);
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
        await sleep(2000);
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

function sleep(time: number) {
    return new Promise((resolve) => setTimeout(resolve, time));
}


</script>

<template>
    <el-form :model="from" label-width="30%">
        <el-form-item label="公司名称:">
            <el-input v-model="from.company" type="textarea" :rows=from.rows @focus="from.rows = 4" @blur="from.rows = 1"
                resize="none" style="width: 50%;"></el-input>
            <el-button type="primary" @click="Collect" style="margin-left: 10px;">开始查询</el-button>
            <el-tooltip placement="left">
                <template
                    #content>需要控股企业查询的资产数量>=1<br />查询公司名或者域名在鹰图中的资产数量，一次查询消耗1积分<br />对应的鹰图查询语句为icp.name=""和domain.suffix=""</template>
                <el-button color="#fbcdef" @click="huntersearch" style="margin-left: 10px;">查询鹰图资产数量</el-button>
            </el-tooltip>
        </el-form-item>
        <el-form-item label="查询条件:">
            <div>
                <el-checkbox v-model="from.subcompany" label="查询子公司并反查域名" />
                <el-tooltip placement="left">
                    <template #content>控股率>=(%)</template>
                    <el-icon>
                        <QuestionFilled />
                    </el-icon>
                </el-tooltip>
                <el-input-number v-model="from.defaultHold" controls-position="right" :min="1" :max="100"
                    style="margin-left: 10px;" v-if="from.subcompany"></el-input-number>
            </div>
            <el-checkbox v-model="from.wechat" label="查询公众号" style="margin-right: 20px; margin-left: 20px;" />
        </el-form-item>
        <el-form-item>
            <template #label>
                Token:
                <el-tooltip placement="right">
                    <template #content>由于天眼查登录校验机制，为了确保数据准确，需要在此处填入网页登录后Cookie头中的 X-Auth-Token 字段或者 auth_token 字段</template>
                    <el-icon>
                        <QuestionFilled size="24" />
                    </el-icon>
                </el-tooltip>
            </template>
            <el-input v-model="from.token" style="width: 50%;"></el-input>
            <el-button color="#abcdef" @click="ExportAssetToXlsx(su, we, hu)" style="margin-left: 10px;">数据导出</el-button>
        </el-form-item>
    </el-form>
    <el-tabs v-model="from.activeName" type="card">
        <el-tab-pane label="控股企业" name="subcompany">
            <el-table :data="su" height="60vh" border style="width: 100%">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="company" label="公司名称" show-overflow-tooltip="true" />
                <el-table-column prop="ratio" label="股权比例" show-overflow-tooltip="true" />
                <el-table-column prop="sunums" label="投资数额" show-overflow-tooltip="true" />
                <el-table-column prop="domains" label="域名" show-overflow-tooltip="true" />
            </el-table>
        </el-tab-pane>

        <el-tab-pane label="公众号" name="wechat">
            <el-table :data="we" height="60vh" border style="width: 100%">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="weName" label="公众号名称" show-overflow-tooltip="true" />
                <el-table-column prop="weNums" label="微信号" show-overflow-tooltip="true" />
                <el-table-column prop="logo" label="Logo" show-overflow-tooltip="true" />
                <el-table-column prop="qrcode" label="二维码" show-overflow-tooltip="true" />
                <el-table-column prop="introduction" label="简介" show-overflow-tooltip="true" />
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="鹰图资产数量" name="hunter">
            <el-table :data="hu" height="60vh" border style="width: 100%">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="name" label="公司域名或ICP名称" show-overflow-tooltip="true" />
                <el-table-column prop="hunums" label="资产数量"
                    :sort-method="(a: any, b: any) => { return a.hunums - b.hunums }" sortable
                    show-overflow-tooltip="true" />
            </el-table>
        </el-tab-pane>
        <el-tab-pane>
            <template #label>
                <el-icon>
                    <Tickets />
                </el-icon>
                <span>运行日志</span>
            </template>
            <el-input class="log-textarea" v-model="from.log" type="textarea" rows="20" readonly></el-input>
        </el-tab-pane>
    </el-tabs>
</template>