<template>
    <CustomTabs>
        <el-tabs type="card">
            <el-tab-pane>
                <template #label>
                    <el-text class="position-center">
                        <el-icon class="mr-5px">
                            <DataAnalysis />
                        </el-icon>
                        数据展示
                    </el-text>
                </template>
                <el-table :data="pagination.table.pageContent" stripe style="height: calc(100vh - 160px);">
                    <!-- <el-table-column fixed type="index" label="#" width="60px" /> -->
                    <el-table-column prop="Query" label="关键词" width="180" />
                    <el-table-column prop="Total" label="总数" width="100px" />
                    <el-table-column prop="Link" label="查询链接" :show-overflow-tooltip="true">
                        <template #default="scope">
                            <el-button link :icon="ChromeFilled"></el-button>
                            {{ scope.row.Link }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="Source" label="来源" width="120px" />
                    <el-table-column label="相关链接" width="100" align="center">
                        <template #default="scope">
                            <el-button @click="showDialog(scope.row.Items)">查看</el-button>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <div class="flex-between mt-5px">
                    <el-progress :text-inside="true" :stroke-width="18" :percentage="parameter.percentage"
                        class="w-40%" />
                    <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
                        @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                        :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100]"
                        :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                        :total="pagination.table.result.length">
                    </el-pagination>
                </div>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <div style="margin-right: -5px;">
                <el-button type="primary" :icon="Plus" plain @click="drawer = true"
                    v-if="!parameter.runningStatus">新建任务</el-button>
                <el-button type="danger" plain @click="stopscan" v-else>停止收集</el-button>
            </div>
        </template>
    </CustomTabs>
    <el-drawer v-model="drawer" size="50%">
        <template #header>
            <span class="drawer-title">新建任务</span>
        </template>
        <el-form :model="parameter" :rules="rules" ref="formRef" label-width="auto">
            <el-form-item label="主要内容" prop="target">
                <el-input v-model="parameter.target" type="textarea" :rows="5" resize="none"
                    placeholder="可以搜索例如域名或者公司名称等信息，多关键词使用换行分割"></el-input>
            </el-form-item>
            <el-form-item label="数据来源">
                <el-checkbox-group v-model="parameter.dataSource">
                    <el-checkbox value="github" :disabled="global.space.github == ''">
                        <el-tooltip content="需要在设置中配置GitHub API">
                            <el-text>
                                <el-icon :size="16" class="mr-5px">
                                    <githubIcon />
                                </el-icon>
                                Github</el-text>
                        </el-tooltip>
                    </el-checkbox>
                    <el-checkbox value="bing">
                        <el-text><el-icon :size="16" class="mr-5px">
                                <bingIcon />
                            </el-icon>
                            Bing</el-text>
                    </el-checkbox>
                </el-checkbox-group>
            </el-form-item>
            <el-form-item label="关键词列表" prop="keyword" v-show="parameter.dataSource.includes('github')">
                <el-select v-model="dictionary" class="mb-5px">
                    <el-option v-for="item in dorksOptions" :key="item.label" :value="item.label">
                        <span class="float-left">{{ item.label }}</span>
                        <span class="float-right">
                            {{ item.value.split("\n").length }}
                        </span>
                    </el-option>
                </el-select>
                <el-input v-model="parameter.keyword" type="textarea" :rows="5" resize="none"></el-input>
            </el-form-item>
            <el-form-item label="Bing搜素语法" prop="bingQuery">
                <el-input v-model="parameter.bingQuery"></el-input>
                <span class="form-item-tips">%s表示查询主要内容的占位符</span>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="position-center">
                <el-button type="primary" @click="Collect">开始任务</el-button>
            </div>
        </template>
    </el-drawer>
    <el-dialog v-model="dialogTableVisible" title="相关链接(最多显示10条内容)" width="700">
        <el-table :data="gridData" class="w-full" style="height: 500px;">
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
import global from '@/stores';
import usePagination from '@/usePagination';
import { sleep, ProcessTextAreaInput } from '@/util';
import { onMounted, reactive, ref, watch } from 'vue';
import { Callgologger, GitDorks, GoogleHackerBingSearch } from 'wailsjs/go/services/App';
import dorks from '@/stores/dorks'
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
import { ElMessage, ElNotification } from 'element-plus';
import { ChromeFilled, Plus } from '@element-plus/icons-vue';
import CustomTabs from '@/components/CustomTabs.vue';
import githubIcon from '@/assets/icon/github.svg';
import bingIcon from '@/assets/icon/bing.svg';
import { structs } from 'wailsjs/go/models';

const drawer = ref(false);

const formRef = ref(null)

const rules = {
    target: [
        { required: true, message: '目标不能为空', trigger: 'blur' },
    ],
    keyword: [
        { required: true, message: '选中GitHub搜索时, 请填写GitHub搜索关键字', trigger: 'blur' },
        {
            validator: (_: any, value: string, callback: (error?: string | Error) => void) => {
                if (parameter.dataSource.includes('github')) {
                    callback(new Error('选中GitHub搜索时, 请填写GitHub搜索关键字'));
                } else {
                    callback();
                }
            },
            trigger: 'blur'
        }
    ],
    bingQuery: [
        { required: true, message: '选中Bing搜索时, 查询语句不能为空，且必须带有%s占位符', trigger: 'blur' },
        {
            validator: (_: any, value: string, callback: (error?: string | Error) => void) => {
                if (parameter.dataSource.includes('bing') && !value.includes('%s')) {
                    callback(new Error('选中Bing搜索时, 查询语句不能为空，且必须带有%s占位符'));
                } else {
                    callback();
                }
            },
            trigger: 'blur'
        }
    ],
}

const parameter = reactive({
    target: '',
    keyword: '',
    bingQuery: 'site: %s',
    dataSource: ['bing'],
    id: 0,
    percentage: 0,
    count: 0,
    runningStatus: false,
});

const dictionary = ref('small_dorks')
const dialogTableVisible = ref(false);
const gridData = ref([] as string[])

const dorksOptions = [
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

const pagination = usePagination<structs.ISICollectionResult>(20)

onMounted(() => {
    updateKeyword() // 在组件挂载时调用 updateKeyword
});

// 监听 dictionary 变化并更新 parameter.keyword
watch(dictionary, () => {
    updateKeyword()
});

// 更新 parameter.keyword 的函数
function updateKeyword() {
    const selected = dorksOptions.find(item => item.label === dictionary.value);
    if (selected) {
        parameter.keyword = selected.value;
    }
}

async function Collect() {
    let targets = ProcessTextAreaInput(parameter.target)
    let dorks = ProcessTextAreaInput(parameter.keyword)
    if (targets.length == 0 || dorks.length == 0) {
        ElMessage.warning("目标地址和关键字均不能为空")
        return
    }
    drawer.value = false
    parameter.count = targets.length * dorks.length
    parameter.runningStatus = true
    parameter.id = 0
    pagination.initTable()
    for (const t of targets) {
        if (parameter.dataSource.includes('bing')) {
            let query = parameter.bingQuery.replace('%s', t)
            let result = await GoogleHackerBingSearch(query)
            if (result) {
                pagination.table.result.push(result)
                pagination.ctrl.watchResultChange(pagination.table)
            }
        }
        if (parameter.dataSource.includes('github')) {
            for (const d of dorks) {
                if (!parameter.runningStatus) {
                    return
                }
                let result = await GitDorks(t, d, global.space.github)
                if (Number(result.Total) > 0) {
                    pagination.table.result.push(result)
                    pagination.ctrl.watchResultChange(pagination.table)
                } else {
                    Callgologger("info", `${t} ${d} 搜索结果为空，已跳过`)
                }
                parameter.id++
                parameter.percentage = Number(((parameter.id / parameter.count) * 100).toFixed(2));
                await sleep(500);
            }
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


<style scoped></style>