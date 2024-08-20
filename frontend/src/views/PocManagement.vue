<template>
    <el-container>
        <el-aside>
            <el-form-item label="显示模式:">
                <el-radio-group v-model="displayMode" @change="">
                    <el-radio :value="1">POC列表</el-radio>
                    <el-radio :value="0">指纹漏洞树</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-input style="margin-bottom: 10px;">
                <template #suffix>
                    <el-tooltip content="搜索指纹用指纹漏洞树，POC用POC列表">
                        <el-icon>
                            <InfoFilled />
                        </el-icon>
                    </el-tooltip>
                </template>
            </el-input>
            <el-card>
                <el-tree :data="treeData" :highlight-current="true" @node-click="handleNodeClick"
                    v-show="displayMode == 0" />
                <el-table-v2 :columns="columns" :data="data" :width="700" :height="tableHeight" 
                fixed :row-event-handlers="{ onClick: tableSelect }" 
                v-show="displayMode == 1" />
            </el-card>
        </el-aside>
        <div style="margin-left: 10px; width: 100%; height: 100%;">
            <el-form-item>
                <div></div>
                <div style="flex: 1;"></div>
                <el-button text bg @click="Copy(content)">Copy</el-button>
            </el-form-item>
            <highlightjs language="yaml" :code="content"></highlightjs>
        </div>
    </el-container>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { GetFingerPocMap } from 'wailsjs/go/main/Database';
import { InfoFilled } from "@element-plus/icons-vue";
import { ReadFile } from 'wailsjs/go/main/File';
import global from '@/global';
import { WebPocFiles } from 'wailsjs/go/main/App';
import { Copy } from '@/util';
import { RowEventHandlerParams } from 'element-plus';

const displayMode = ref(1)
const content = ref('')

// 定义树的数据结构
interface TreeNode {
    label: string;
    children?: TreeNode[];
}
const tableHeight = ref(window.innerHeight * 0.78);

// 创建一个 ref 变量来存储树的数据
const treeData = ref<TreeNode[]>([]);
const listData = ref<string[]>([])

onMounted(async () => {
    const result = await GetFingerPocMap();
    // 将 result 转换为 el-tree 需要的数据格式
    treeData.value = Object.keys(result).map(key => {
        return {
            label: key,
            children: result[key].map(item => ({ label: item })),
        };
    });
    listData.value = await WebPocFiles();
    data.value = generateData(columns, listData.value);

    window.addEventListener('resize', () => {
        tableHeight.value = window.innerHeight * 0.78;
    });
});

async function tableSelect(row: RowEventHandlerParams) {
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + row.rowData.fileName + ".yaml"
    let file: any = await ReadFile(filepath)
    content.value = file.Content!
}

const handleNodeClick = async (data: TreeNode) => {
    if (!data.children) {
        let filepath = global.PATH.homedir + "/slack/config/pocs/" + data.label + ".yaml"
        let file: any = await ReadFile(filepath)
        content.value = file.Content!
    }
}

// Generate a single column
const generateColumns = (length = 1, props?: any) =>
    Array.from({ length }).map((_, columnIndex) => ({
        ...props,
        key: `${columnIndex}`,
        dataKey: 'fileName',  // Adjust to match the data key
        title: 'File Name',
        width: 300,
    }));

// Generate the table data from listData
const generateData = (columns: ReturnType<typeof generateColumns>, files: string[]) =>
    files.map((file, rowIndex) => ({
        id: `${rowIndex}`,
        fileName: file,  // Map the file name to the correct data key
    }));

const columns = generateColumns(1);
const data = ref<{ fileName: string }[]>([]);
</script>

<style scoped>
.el-card {
    --el-card-padding: 5px;
}

.el-tree {
    height: 78vh;
    overflow-y: auto;
    overflow-x: hidden;
}

:deep(.el-table__header-wrapper) {
    height: 0;
}
</style>