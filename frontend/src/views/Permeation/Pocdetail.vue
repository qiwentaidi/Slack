<script lang="ts" setup>
import { Search } from '@element-plus/icons-vue';
import { reactive } from 'vue';
import {
    LocalWalkFiles,
    ReadPocDetail
} from '../../../wailsjs/go/main/App'
import { UserHomeDir, PathBase } from '../../../wailsjs/go/main/File'
import { onMounted } from 'vue';
// 初始化时调用
onMounted(async () => {
    LoadPocList(window.AFGPathPoc);
});

var from = reactive({
    value: window.AFGPathPoc,
    pocOptions: [{
        label: "通用POC",
        value: window.AFGPathPoc
    }, {
        label: "主动目录探测",
        value: window.ActivePathPoc
    }],
    keyword: '',
    expands: [],
    filteredPocList: [{}],
    count: 0
})

// poc内容
const pd = reactive({
    names: '',
    risk: '',
    author: '',
    tags: '',
    createTime: '',
    description: '',
    reference: '',
    affected: '',
    solutions: '',
})

async function LoadPocList(filepath: string) {
    let poclist = LocalWalkFiles(await UserHomeDir() + filepath)
    let index = 0
    table.result = []
    for (const fullpath of await poclist) {
        const name = (await PathBase(fullpath)).split(".")[0];
        table.result.push(
            {
                id: index,
                names: name,
                fullPath: fullpath,
            }
        )
        index++
        table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
    }
    from.count = table.result.length
}

function readpoc(row: any) {
    ReadPocDetail(row.fullPath).then(
        result => {
            pd.names = result.Name;
            pd.risk = result.Risk;
            pd.author = result.Author;
            pd.tags = result.Tags;
            pd.description = result.Description;
            pd.reference = result.Reference;
            pd.affected = result.Affected;
            pd.solutions = result.Solutions;
        }
    )
}

// expandedRows 为展开的所有列
async function loadData(row: any, expandedRows: any) {
    from.expands = []
    if (expandedRows.length > 0) {
        from.expands.push(row.id as never)
    }
    readpoc(row);
}

function filterPocList() {
    let temp = []
    for (const item of table.result) {
        if (item.names.toLowerCase().includes(from.keyword.toLowerCase())){
            temp.push(item)
        }
    }
    from.filteredPocList = temp
    from.count = temp.length
    table.pageContent = from.filteredPocList.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
}

interface Poc {
 id: number;
 names: string;
 fullPath: string;
}

const table = reactive({
    currentPage: 1,
    pageSize: 10,
    pageContent: [{}],
    result: [] as Poc[]
})


function handleSizeChange(val: any) {
    table.pageSize = val;
    table.currentPage = 1;
    table.pageContent = table.result.slice(0, val)
}
function handleCurrentChange(val: any) {
    table.currentPage = val;
    table.pageContent = table.result.slice((val - 1) * table.pageSize, (val - 1) * table.pageSize + table.pageSize)
}
</script>

<template>
    <el-form :model="from" label-width="100px">
        <el-form-item label="关键字查询">
            <el-input v-model="from.keyword" placeholder="搜索" :prefix-icon="Search" clearable @input="filterPocList"/>
        </el-form-item>
        <el-form-item label="类型查询">
            <el-select v-model="from.value" collapse-tags collapse-tags-tooltip placeholder="选择需要查看的POC类型" @change="LoadPocList(from.value)">
                <el-option v-for="item in from.pocOptions" :label="item.label" :value="item.value" />
            </el-select>
        </el-form-item>
    </el-form>
    <div>
        <el-table :data="table.pageContent" border style="height: 70vh;" @expand-change="loadData" row-key="id"
            :expand-row-keys="from.expands">
            <el-table-column type="expand">
                <template #default="props">
                    <el-descriptions title="👾漏洞详情" :column="3" border>
                        <el-descriptions-item label="名称" :span="2">{{ pd.names }}</el-descriptions-item>
                        <el-descriptions-item label="风险等级" :span="1">{{ pd.risk }}</el-descriptions-item>
                        <el-descriptions-item label="作者" :span="2">{{ pd.author }}</el-descriptions-item>
                        <el-descriptions-item label="标签" :span="1">{{ pd.tags }}</el-descriptions-item>
                        <el-descriptions-item label="漏洞描述" :span="3">{{ pd.description }}</el-descriptions-item>
                        <el-descriptions-item label="参考文档" :span="3">{{ pd.reference }}</el-descriptions-item>
                        <el-descriptions-item label="影响版本" :span="3">{{ pd.affected }}</el-descriptions-item>
                        <el-descriptions-item label="解决方案" :span="3">{{ pd.solutions }}</el-descriptions-item>
                    </el-descriptions>
                </template>
            </el-table-column>
            <el-table-column prop="names" label="名称" />
        </el-table>
        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
            :current-page="table.currentPage" :page-sizes="[10, 20, 30]" :page-size="table.pageSize"
            layout="total, sizes, prev, pager, next" :total="from.count" style="margin-top: 5px; float: right;">
        </el-pagination>
    </div>
</template>
