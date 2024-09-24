<template>
    <el-form :model="parameter" label-position="top" :inline="true">
        <el-form-item>
            <template #label>
                <span style="font-weight: bold;">主要内容</span>
            </template>
            <el-input v-model="parameter.target" 
            type="textarea" :rows="5" 
            resize="none" placeholder="可以搜索例如域名或者公司名称等信息，多关键词使用换行分割" style="width: 40vh;"></el-input>
        </el-form-item>
        <el-form-item>
            <template #label>
                <div style="display: flex;">
                    <span style="font-weight: bold; width: 50%;">关键词列表</span>
                    <el-select v-model="dictionary" size="small" style="height: 17px;">
                        <el-option v-for="item in List" :key="item.label" :value="item.label">
                            <span style="float: left">{{ item.label }}</span>
                            <span style="float: right">
                                {{ item.value.split("\n").length }}
                            </span>
                        </el-option>
                    </el-select>
                </div>
            </template>
            <el-input v-model="parameter.keyword" type="textarea" :rows="5" resize="none" style="width: 40vh;"></el-input>
        </el-form-item>
        <el-form-item>
            <template #label>
                <span style="font-weight: bold;">使用说明</span>
            </template>
            <div>
                <el-text>项目参考<el-link @click="BrowserOpenURL('https://github.com/obheda12/GitDorker')">GitDorks</el-link>，
                    需要配置<el-link @click="$router.push('/Settings')">Github API</el-link>，内容查看最多100条</el-text>
                <el-button @click="Collect" v-if="!parameter.runningStatus">开始收集</el-button>
                <el-button type="danger" @click="stopscan" v-else>停止收集</el-button>
            </div>
        </el-form-item>
    </el-form>
    <el-table :data="pagination.table.pageContent" stripe height="66vh">
        <el-table-column fixed type="index" label="#" width="60px" />
        <el-table-column prop="Query" label="关键词" width="180" />
        <el-table-column prop="Total" label="总数" width="100" />
        <el-table-column prop="Link" label="查询链接" :show-overflow-tooltip="true">
            <template #default="scope">
                <el-button link :icon="ChromeFilled"></el-button>
                {{ scope.row.Link}}
            </template>
        </el-table-column>
        <el-table-column prop="Items" label="相关链接" width="100" align="center">
            <template #default="scope">
                <el-button @click="showDialog(scope.row.Items)">查看</el-button>
            </template>
        </el-table-column>
        <template #empty>
            <el-empty />
        </template>
    </el-table>
    <div class="my-header" style="margin-top: 5px;">
        <el-progress :text-inside="true" :stroke-width="18" :percentage="parameter.percentage" 
            style="width: 40%;" />
        <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
            @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
            :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100]"
            :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
            :total="pagination.table.result.length">
        </el-pagination>
    </div>
    <el-dialog v-model="dialogTableVisible" title="相关链接" width="700">
        <el-table :data="gridData" style="width: 100%; height: 500px;">
            <el-table-column type="index" width="50" />
            <el-table-column label="Link" :show-overflow-tooltip="true">
                <template #default="scope">
                    {{ scope.row }}
                </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="center">
                <template #default="scope">
                    <el-button link :icon="ChromeFilled" @click="BrowserOpenURL(scope.row)"></el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-dialog>
</template>

<script lang="ts" setup>
import global from '@/global';
import { ISICResult } from '@/interface';
import usePagination from '@/usePagination';
import { sleep, SplitTextArea } from '@/util';
import { onMounted, reactive, ref, watch } from 'vue';
import { Callgologger, GitDorks } from 'wailsjs/go/main/App';
import dorks from '@/stores/dorks'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
import { ElMessage, ElNotification } from 'element-plus';
import { ChromeFilled } from '@element-plus/icons-vue';

const parameter = reactive({
    target: '',
    keyword: '',
    id: 0,
    percentage: 0,
    count: 0,
    runningStatus: false,
});

const dictionary = ref('small_dorks')
const dialogTableVisible = ref(false);
const gridData = ref([] as string[])

const List = [
    {
        label: "all_dorks",
        value: dorks.alldorksv3
    },
    {
        label: "medium_dorks",
        value: dorks.medium_dorks
    },
    {
        label: "small_dorks",
        value: dorks.smalldorks
    }
]

const result = ref([] as ISICResult[])

const pagination = usePagination(result.value, 20)

onMounted(() => {
    updateKeyword() // 在组件挂载时调用 updateKeyword
});

// 监听 dictionary 变化并更新 parameter.keyword
watch(dictionary, () => {
    updateKeyword()
});

// 更新 parameter.keyword 的函数
function updateKeyword() {
    const selected = List.find(item => item.label === dictionary.value);
    if (selected) {
        parameter.keyword = selected.value;
    }
}

async function Collect() {
    let targets = SplitTextArea(parameter.target)
    let dorks = SplitTextArea(parameter.keyword)
    if (targets.length == 0 || dorks.length == 0) {
        ElMessage.warning("目标地址和关键字均不能为空")
        return
    }
    parameter.count = targets.length * dorks.length
    parameter.runningStatus = true
    parameter.id = 0
    pagination.table.result = []
    pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    for (const t of targets) {
        for (const d of dorks) {
            if (!parameter.runningStatus) {
                return
            }
            let result:any = await GitDorks(t, d, global.space.github)
            if (result.Status && Number(result.Total) > 0) {
                pagination.table.result.push({
                    Query: t + " " + d,
                    Status: result.Status,
                    Total: result.Total,
                    Items: result.Items,
                    Link: result.Link,
                })
                pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
            } else {
                Callgologger("info", `${t} ${d} 搜索结果为空，已跳过`)
            }
            parameter.id++
            parameter.percentage = Number(((parameter.id / parameter.count) * 100).toFixed(2));
            await sleep(500);
        }
    }
    parameter.runningStatus = false
}

function stopscan() {
    parameter.runningStatus = false
    ElNotification.error({
        message: "用户已终止扫描",
        position: "bottom-right",
    })
}
function showDialog(list: string[]) {
    dialogTableVisible.value = true
    gridData.value = list
}

</script>


<style scoped>

</style>