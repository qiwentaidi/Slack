<script lang="ts" setup>
import global from "@/global";
import { Copy, ReadLine, transformArrayFields } from '@/util'
import { ExportToXlsx } from '@/export'
import { reactive, ref, onMounted } from "vue";
import { ExitScanner, Subdomain } from "wailsjs/go/main/App";
import { CheckFileStat, FileDialog } from "wailsjs/go/main/File";
import { ElMessage, ElNotification } from 'element-plus'
import { Menu, CopyDocument } from '@element-plus/icons-vue';
import exportIcon from '@/assets/icon/doucment-export.svg'
import usePagination from "@/usePagination";
import { SubdomainInfo } from "@/stores/interface";
import { EventsOn, EventsOff } from "wailsjs/runtime/runtime";
import { debounce } from "lodash"
import { validateSingleDomain } from "@/stores/validate";
import { structs } from "wailsjs/go/models";
import { dnsServerOptions, subdomainRunnerOptions } from "@/stores/options";

const debounceUpdate = debounce(() => {
  pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table);
}, 500);

onMounted(() => {
    EventsOn("subdomainLoading", (result: SubdomainInfo) => {
        if (result.Source == "Enumeration") {
            let r = pagination.table.result.find(item => item.Subdomain == result.Subdomain)
            if (r == undefined) {
                return
            }
        }
        pagination.table.result.push(result)
        debounceUpdate()
        // pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    });
    EventsOn("subdomainProgressID", (id: number) => {
        config.id = id
        config.percentage = Number(((config.id / config.count) * 100).toFixed(2));
    });
    EventsOn("subdomainCounts", (count: number) => {
        config.count = count
    });
    EventsOn("subdomainComplete", (message: string) => {
        ElMessage.info(message)
        config.runningStatus = false
        config.percentage = 100
    });
    return () => {
        EventsOff("subdomainLoading");
        EventsOff("subdomainComplete");
    };
});

const currentRunner = ref(1);
const result = ref([] as SubdomainInfo[]);
const pagination = usePagination(result.value, 50)
const input = ref("");

const selectDnsServer = ref(["223.6.6.6:53", "8.8.8.8:53"])

const config = reactive({
    thread: 600,
    timeout: 5,
    resolveExcludeTimes: 5,
    subs: [] as string[],
    subFilepath: "",
    percentage: 0,
    id: 0,
    count: 0,
    runningStatus: false,
    drawer: false,
});

async function NewTask() {
    if (!config.runningStatus) {
        let task = new Runner()
        if (await task.checkInput()) {
            await task.NewRunner()
        }
    } else {
        ExitScanner("[subdomain]")
        config.runningStatus = false
        ElNotification.error({
            message: "用户已终止扫描任务",
            position: 'bottom-right',
        });
    }
}

class Runner {
    public async checkInput() {
        if (input.value == '') {
            ElMessage.warning("请输入域名或者域名文件")
            return false
        }
        if (validateSingleDomain(input.value)) {
            return true
        }
        let stat = await CheckFileStat(input.value)
        if (!stat) {
            ElMessage.warning('输入的文件路径不存在')
        }
        return stat
    }

    public async NewRunner() {
        let domains = [] as string[]
        if (await CheckFileStat(input.value)) {
            domains = (await ReadLine(input.value))!
        } else {
            domains = [input.value]
        }
        if (currentRunner.value == 1 && (!global.space.chaos) && (!global.space.bevigil) && (!global.space.securitytrails) && (!global.space.zoomeye)) {
            ElMessage.warning("请至少填写一个API进行查询")
            return
        }
        config.runningStatus = true
        config.id = 0;
        if (currentRunner.value != 1 && config.subs.length == 0) {
            config.subs = (await ReadLine(global.PATH.homedir + "/slack/config/subdomain/dicc.txt"))!
        }
        let option: structs.SubdomainOption = {
            Mode: currentRunner.value,
            Domains: domains,
            Subs: config.subs,
            Thread: config.thread,
            Timeout: config.timeout,
            DnsServers: selectDnsServer.value,
            ResolveExcludeTimes: 5,
            BevigilApi: global.space.bevigil,
            ChaosApi: global.space.chaos,
            SecuritytrailsApi: global.space.securitytrails,
            ZoomeyeApi: global.space.zoomeye,
            GithubApi: global.space.github,
        }
        await Subdomain(option)
    }
}
async function handleFileChange() {
    config.subFilepath = await FileDialog("*.txt")
    config.subs = (await ReadLine(config.subFilepath))!
}

async function handleTargetFileChange() {
    input.value = await FileDialog("*.txt")
}

const CopyDomains = () => {
    let subdomains = pagination.table.result.map(item => {
        return item.Subdomain
    })
    Copy(subdomains.join("\n"))
}
</script>

<template>
    <el-form :model="config">
        <el-form-item>
            <div class="head">
                <el-input v-model="input" placeholder="请输入域名或域名文件列表" style="margin-right: 10px;">
                    <template #prepend>
                        <el-select v-model="currentRunner" style="width: 150px">
                            <el-option v-for="item in subdomainRunnerOptions" :key="item.value" :label="item.label"
                                :value="item.value" style="width: 260px;">
                                <span style="float: left">{{ item.label }}</span>
                                <span class="select-tips">
                                    {{ item.tips }}
                                </span>
                            </el-option>
                        </el-select>
                    </template>
                    <template #suffix>
                        <el-button link @click="handleTargetFileChange()">
                            <el-icon>
                                <Document />
                            </el-icon>
                        </el-button>
                    </template>
                </el-input>
                <el-button type="primary" @click="NewTask" v-if="!config.runningStatus">开始任务</el-button>
                <el-button type="danger" @click="NewTask" v-else>停止任务</el-button>
                <el-button text bg @click="config.drawer = true">参数设置</el-button>
            </div>
        </el-form-item>
        <el-form-item>
            <el-space>
                <el-tag>子域字典大小:{{ config.subs.length }}</el-tag>
                <el-tag>线程:{{ config.thread }}</el-tag>
                <el-tag type="info">Tips: 建议均使用查询模式进行子域名收集，优点在于不用消耗本地的网络资源且搜集快</el-tag>
            </el-space>
            <el-dropdown style="margin-left: auto;">
                <el-button :icon="Menu" text bg />
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item
                            @click="ExportToXlsx(['主域名', '子域名', 'IPS', '是否为CDN', 'CDN名称', '来源'], '子域名暴破', 'subdomain', transformArrayFields(pagination.table.result))"
                            :icon="exportIcon">导出Excel</el-dropdown-item>
                        <el-dropdown-item
                            @click="CopyDomains()"
                            :icon="CopyDocument" divided>复制所有子域名</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </el-form-item>
    </el-form>
    <el-table :data="pagination.table.pageContent" border :cell-style="{ textAlign: 'center' }" 
    :header-cell-style="{ 'text-align': 'center' }" style="height: calc(100vh - 205px);">
        <el-table-column type="index" label="#" width="60px" />
        <el-table-column prop="Domain" label="主域名" />
        <el-table-column prop="Subdomain" label="子域名" />
        <el-table-column prop="Ips" label="IPs" />
        <el-table-column prop="CdnName" label="CDN/WAF" width="150px">
            <template #default="scope">
                <el-tag type="danger" v-if="scope.row.IsCdn">{{ scope.row.CdnName }}</el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="Source" label="来源" width="120px" />
        <template #empty>
            <el-empty></el-empty>
        </template>
    </el-table>
    <div class="my-header" style="margin-top: 5px;">
        <el-progress :text-inside="true" :stroke-width="18" :percentage="config.percentage" 
            style="width: 40%;" />
        <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
            @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
            :current-page="pagination.table.currentPage" :page-sizes="[50, 100, 200, 500]"
            :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
            :total="pagination.table.result.length">
        </el-pagination>
    </div>
    <el-drawer v-model="config.drawer" size="50%">
        <template #header>
            <h3>设置高级参数</h3>
        </template>
        <el-form :model="config" label-width="auto">
            <el-form-item label="线程数量">
                <el-input-number v-model="config.thread" :min="600" :max="10000" controls-position="right">
                </el-input-number>
            </el-form-item>
            <el-form-item label="解析超时(s)">
                <el-input-number v-model="config.timeout" :min="1" :max="20" controls-position="right">
                </el-input-number>
            </el-form-item>
            <el-form-item>
                <template #label>
                    <span>过滤IP次数</span>
                    <el-tooltip placement="left">
                        <template #content>设置次数可以有效过滤泛解析数据，确保在泛解析<br />
                            的域名中获取到不同的IP信息，次数为0表示不过滤
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="config.resolveExcludeTimes" :min="0" :max="1000" controls-position="right">
                </el-input-number>
            </el-form-item>
            <el-form-item label="DNS Servers">
                <el-select v-model="selectDnsServer" multiple clearable collapse-tags collapse-tags-tooltip
                    :max-collapse-tags="3">
                    <el-option v-for="item in dnsServerOptions" :key="item.value"
                        :label="item.value"
                        :value="item.value">
                        <span style="float: left">{{ item.value }}</span>
                        <span class="select-tips">
                            {{ item.label }}
                        </span>
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="子域字典选择">
                <el-input v-model="config.subFilepath">
                    <template #suffix>
                        <el-button link @click="handleFileChange">
                            <el-icon>
                                <Document />
                            </el-icon>
                        </el-button>
                    </template>
                </el-input>
            </el-form-item>
        </el-form>
    </el-drawer>
</template>
