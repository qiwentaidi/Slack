<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { Copy, parseHeaders, ProcessTextAreaInput } from '@/util';
import { AnalyzeAPI, ExtractAllJSLink, JSFind } from 'wailsjs/go/services/App';
import { ArrowUpBold, ArrowDownBold, Delete, DocumentCopy } from '@element-plus/icons-vue';
import global from "@/stores";
import { ElNotification, ElMessage } from 'element-plus';
import CustomTextarea from '@/components/CustomTextarea.vue';
import saveIcon from '@/assets/icon/save.svg'
import usePagination from '@/usePagination';
import { JSFindOptions } from '@/stores/options';
import { JSFindData } from '@/stores/interface';
import { BrowserOpenURL, EventsOff, EventsOn } from 'wailsjs/runtime/runtime';
import { getTagTypeBySeverity } from '@/stores/style';
import { SaveConfig } from '@/config';

onMounted(() => {
    EventsOn("jsfindlog", (msg: any) => {
        config.consoleLog += msg + "\n";
    });
    EventsOn("jsfindvulcheck", (result: any) => {
        pagination.table.result.push({
            Target: "",
            Method: result.Method,
            Source: result.Source,
            VulType: "未授权访问",
            Severity: "HIGH",
            Param: result.Param,
            Length: result.Length,
            Filed: "",
            Response: result.Response,
        })
        pagination.ctrl.watchResultChange(pagination.table)
    });
    return () => {
        EventsOff("jsfindlog");
        EventsOff("webFingerScan");
    };
})

const value = ref(0)

const config = reactive({
    urls: "",
    loading: false,
    otherURL: false,
    prefixApiURL: '',
    prefixJsURL: '',
    headers: '',
    consoleLog: '',
})

const pagination = usePagination<JSFindData>(20)

async function JSFinder() {
    let blackList = global.jsfind.whiteList.split("\n")
    let urls = ProcessTextAreaInput(config.urls)
    if (urls.length == 0) {
        ElMessage.warning("可用目标为空");
        return
    }
    showForm.value = false
    config.loading = true
    pagination.initTable()
    config.consoleLog = ""
    for (const url of urls) {
        let apiRoute = [] as string[]
        config.consoleLog += `[*] 正在提取${url}的JS链接...\n`
        let jslinks = await ExtractAllJSLink(url)
        config.consoleLog += `[+] 共提取到JS链接: ${getLength(jslinks)}个\n`
        config.consoleLog += jslinks.join("\n")
        config.consoleLog += "\n\n[*] 正在提取JS信息中...\n"
        let somethings = await JSFind(url, config.prefixJsURL, jslinks)
        config.consoleLog += `[+] 共提取到API: ${getLength(somethings.APIRoute)}个\n`
        somethings.APIRoute.forEach(item => {
            apiRoute.push(item.Filed)
            config.consoleLog += `${item.Filed}\n`
        })
        config.consoleLog += "\n\n"
        config.consoleLog += `[+] 共提取到身份证: ${getLength(somethings.IDCard)}个\n`
        somethings.IDCard.forEach(item => {
            config.consoleLog += `${item.Filed} [ Source: ${item.Source} ]\n`
            pagination.table.result.push({
                Source: item.Source,
                Method: "GET",
                Filed: item.Filed,
                VulType: "身份证号信息泄露",
                Severity: "MEDIUM",
                Length: 0,
                Param: "",
                Target: url,
                Response: "",
            })
            pagination.ctrl.watchResultChange(pagination.table)
        })
        config.consoleLog += "\n\n"
        config.consoleLog += `[+] 共提取到手机号: ${getLength(somethings.Phone)}个\n`
        somethings.Phone.forEach(item => {
            config.consoleLog += `${item.Filed} [ Source: ${item.Source} ]\n`
            pagination.table.result.push({
                Source: item.Source,
                Method: "GET",
                Filed: item.Filed,
                VulType: "手机号信息泄露",
                Severity: "LOW",
                Length: 0,
                Param: "",
                Target: url,
                Response: "",
            })
            pagination.ctrl.watchResultChange(pagination.table)
        })
        config.consoleLog += "\n\n"
        config.consoleLog += `[+] 共提取到邮箱: ${getLength(somethings.Email)}个\n`
        if (somethings.Email) {
            const emails = somethings.Email.map(item => item.Filed)
            config.consoleLog += emails.join("\n")
        }
        config.consoleLog += "\n\n"
        if (somethings.Sensitive) {
            config.consoleLog += `[+] 共提取到敏感字段: ${getLength(somethings.Sensitive)}个\n`
            somethings.Sensitive.forEach(item => {
                config.consoleLog += `${item.Filed} [ Source: ${item.Source} ]\n`
                pagination.table.result.push({
                    Source: item.Source,
                    Method: "GET",
                    Filed: item.Filed,
                    VulType: "敏感字段泄露",
                    Severity: "INFO",
                    Length: 0,
                    Param: "",
                    Target: url,
                    Response: "",
                })
                pagination.ctrl.watchResultChange(pagination.table)
            })
        }
        config.consoleLog += "\n\n"
        if (somethings.IP_URL && somethings.IP_URL.length > 0) {
            const filteredIPs = somethings.IP_URL
                .filter(item => !blackList.some(black => item.Filed.includes(black))) // 过滤黑名单
                .map(item => item.Filed); // 提取字段

            if (filteredIPs.length > 0) {
                config.consoleLog += `[+] 共提取到IP/URL: ${filteredIPs.length}个\n` + filteredIPs.join("\n") + "\n";
            }
        }
        config.consoleLog += "\n\n"
        config.consoleLog += "[*] 正在分析API漏洞中...\n"
        let baseURL = ""

        config.prefixApiURL != "" ? baseURL = config.prefixApiURL : baseURL = url

        await AnalyzeAPI(url, baseURL, apiRoute, parseHeaders(config.headers), global.jsfind.authFiled.split(","), global.jsfind.highRiskRouter.split(","))
    }
    config.consoleLog += "[*] 任务运行结束\n"
    config.loading = false
    ElNotification.success({
        message: "JSFinder 任务已完成",
        position: "bottom-right",
    });
}

function getLength(arr: any) {
    if (Array.isArray(arr)) {
        return arr.length;
    } else {
        return 0;
    }
}

const showForm = ref(true);

function toggleFormVisibility() {
    showForm.value = !showForm.value;
}

const dialog = ref(false);
const detail = reactive({
    Target: '',
    Method: '',
    Source: '',
    VulType: '',
    Severity: '',
    Param: '',
    Filed: '',
    Response: '',
});
function openDialog(row: any) {
    dialog.value = true;
    detail.Target = row.Target;
    detail.Method = row.Method;
    detail.Source = row.Source;
    detail.VulType = row.VulType;
    detail.Severity = row.Severity;
    detail.Param = row.Param;
    detail.Filed = row.Filed;
    // 超过50kb不显示
    if (row.Response && row.Response.length >= 500 * 1024) {
        detail.Response = 'Response too large, please manually open the link to view.';
    } else {
        detail.Response = row.Response;
    }
}
</script>

<template>
    <el-divider>
        <el-button round :icon="showForm ? ArrowUpBold : ArrowDownBold" @click="toggleFormVisibility"
            v-if="!config.loading">
            {{ showForm ? '隐藏参数' : '展开参数' }}
        </el-button>
        <el-button round loading v-else>正在运行</el-button>
    </el-divider>
    <el-collapse-transition>
        <div style="display: flex; gap: 10px;" v-show="showForm">
            <el-form :model="config" label-width="auto" style="width: 50%;">
                <el-form-item label="目标地址:">
                    <CustomTextarea v-model="config.urls" :rows="5" />
                    <span class="form-item-tips">目前功能还在测试阶段, 尽量只使用单目标查询, 且需要输入目录根路径</span>
                </el-form-item>
                <el-form-item label="JS前缀:">
                    <el-input v-model="config.prefixJsURL" />
                    <span class="form-item-tips">部分二级路径采集的JS无法准确拼接时, 自定义路径前缀</span>
                </el-form-item>
                <el-form-item label="路径前缀:">
                    <el-input v-model="config.prefixApiURL" />
                    <span class="form-item-tips">默认将获取到的API拼接到URL路径中, 大部分API需要都拼接在接口路径, 需要自行获取</span>
                </el-form-item>
                <el-form-item label="请求头:">
                    <el-input v-model="config.headers" type="textarea" :rows="5" />
                    <span class="form-item-tips">{{ $t('tips.customHeaders') }}</span>
                </el-form-item>
            </el-form>
            <el-form :model="config" label-width="auto" style="width: 50%;">
                <el-form-item label="鉴权字段:">
                    <el-input v-model="global.jsfind.authFiled" type="textarea" :rows="5"></el-input>
                    <span class="form-item-tips">判断内容响应体是否需要鉴权, 以逗号分隔</span>
                </el-form-item>
                <el-form-item label="高危路由:">
                    <el-input v-model="global.jsfind.highRiskRouter" type="textarea" :rows="5"></el-input>
                    <span class="form-item-tips">以逗号分隔</span>
                </el-form-item>
                <el-form-item label="黑名单域名:">
                    <el-input v-model="global.jsfind.whiteList" type="textarea" :rows="5"></el-input>
                    <span class="form-item-tips">过滤IP/URL内容, 以换行分隔</span>
                </el-form-item>
                <el-form-item label=" ">
                    <el-button type="primary" @click="JSFinder">开始任务</el-button>
                    <el-button :icon="saveIcon" @click="SaveConfig()">保存配置</el-button>
                </el-form-item>
            </el-form>
        </div>
    </el-collapse-transition>
    <el-card shadow="never" style="margin-bottom: 10px;">
        <template #header>
            <div class="card-header">
                <el-segmented v-model="value" :options="JSFindOptions">
                    <template #default="{ item }">
                        <el-space :size="3">
                            <el-icon>
                                <component :is="item.icon" />
                            </el-icon>
                            <div>{{ item.label }}</div>
                        </el-space>
                    </template>
                </el-segmented>
                <div v-show="value == 0">
                    <el-button :icon="DocumentCopy" link @click="Copy(config.consoleLog)" />
                    <el-button :icon="Delete" link @click="config.consoleLog = ''" />
                </div>
            </div>
        </template>
        <pre class="pretty-response" style="margin-top: 0; margin-bottom: 0;" v-show="value == 0"><code>{{ config.consoleLog
        }}</code></pre>
        <el-table :data="pagination.table.pageContent" style="height: calc(100vh - 200px)"
            @sort-change="pagination.ctrl.sortChange" v-show="value == 1">
            <el-table-column label="INFO">
                <template #default="scope">
                    <el-space>
                        <el-tag type="info">{{ scope.row.Method }}</el-tag>
                        <el-tag type="info">{{ scope.row.Source }}</el-tag>
                        <el-tag type="warning" v-if="scope.row.Filed != ''">{{ scope.row.Filed }}</el-tag>
                    </el-space>
                </template>
            </el-table-column>
            <el-table-column prop="Length" label="Length" width="150" sortable="custom" />
            <el-table-column prop="VulType" label="Vulnerability" column-key="Severity" width="300">
                <template #filter-icon>
                    <Filter />
                </template>
                <template #default="scope">
                    <el-space>
                        <el-tag :type="getTagTypeBySeverity(scope.row.Severity)">
                            {{ scope.row.Severity }}
                        </el-tag>
                        <span>{{ scope.row.VulType }}</span>
                    </el-space>
                </template>
            </el-table-column>
            <el-table-column width="100" align="center">
                <template #default="scope">
                    <el-button type="primary" @click="openDialog(scope.row)">More</el-button>
                </template>
            </el-table-column>
            <template #empty>
                <el-empty />
            </template>
        </el-table>
        <div class="my-header" style="margin-top: 5px;" v-show="value == 1">
            <div></div>
            <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
                @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                :current-page="pagination.table.currentPage" :page-sizes="[10, 20, 50, 100]"
                :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="pagination.table.result.length">
            </el-pagination>
        </div>
    </el-card>
    <el-dialog v-model="dialog" title="漏洞详情">
        请求方式: {{ detail.Method }} <br /><br />
        <span v-if="detail.VulType != '未授权访问'">源目标: <el-link type="primary" @click="BrowserOpenURL(detail.Target)">{{
            detail.Target }}</el-link> <br /><br /></span>
        来源链接: <el-link type="primary" @click="BrowserOpenURL(detail.Source)">{{ detail.Source }}</el-link><br /><br />
        漏洞类型: {{ detail.VulType }} <br /><br />
        漏洞等级: <el-tag :type="getTagTypeBySeverity(detail.Severity)">
            {{ detail.Severity }}
        </el-tag> <br /><br />
        <span v-if="detail.Param != ''">漏洞参数: {{ detail.Param }}</span><br /><br />
        <span>匹配内容: <strong>{{ detail.Filed }}</strong></span><br /><br />
        <span><br /><br />响应内容: {{ detail.Response }}</span><br />
    </el-dialog>
</template>

<style scoped></style>