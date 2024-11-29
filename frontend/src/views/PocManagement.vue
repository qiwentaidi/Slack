<script lang="ts" setup>
import { watch, onMounted, onUnmounted, ref, nextTick, toRaw } from 'vue';
import { Search, CirclePlusFilled, Delete, DocumentCopy, Reading, Edit } from "@element-plus/icons-vue";
import { CheckFileStat, ReadFile, RemoveFile, WriteFile } from 'wailsjs/go/services/File';
import global from '@/global';
import { FingerprintList, GetFingerPocMap } from 'wailsjs/go/services/App';
import { Copy } from '@/util';
import { PocDetail } from '@/stores/interface';
import usePagination from '@/usePagination';
import { ElMessage, ElMessageBox } from 'element-plus';
import CustomTabs from '@/components/CustomTabs.vue';
import saveIcon from '@/assets/icon/save.svg'
import { SaveConfig } from '@/config';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import "monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution";

onMounted(async () => {
    const pocMap = await GetFingerPocMap();
    for (const [poc, tags] of Object.entries(pocMap)) {
        pagination.table.result.push({
            Name: poc,
            AssociatedFingerprint: Array.from(new Set(tags))
        })
    }
    pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    pagination.table.temp = pagination.table.result
    let fingers = await FingerprintList()
    const uniqueValues = new Set();
    fingers.forEach(item => {
        if (!uniqueValues.has(item)) {
            fingerOptions.value.push({
                label: item,
                value: item
            });
            uniqueValues.add(item);
        }
    });
    highlightFingerOptions.value = fingerOptions.value.filter(
        option => !global.webscan.highlight_fingerprints.includes(option.value)
    );

});

const pagination = usePagination<PocDetail>(20)

const defaultFilter = ref('Name')

const filterOptions = [
    {
        label: '名称',
        value: 'Name'
    },
    {
        label: '关联指纹',
        value: 'Fingerprint'
    },
]

const filter = ref('')
function filterPocList() {
    if (filter.value == '') {
        pagination.table.result = pagination.table.temp
        pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
        return
    }
    pagination.table.result = []
    if (defaultFilter.value != "Name") {
        for (const item of pagination.table.temp) {
            for (const finger of item.AssociatedFingerprint) {
                if (finger.toLowerCase().includes(filter.value.toLowerCase())) {
                    pagination.table.result.push(item)
                    break
                }
            }
        }
    } else {
        for (const item of pagination.table.temp) {
            if (item.Name.toLowerCase().includes(filter.value.toLowerCase())) {
                pagination.table.result.push(item)
            }
        }
    }

    pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
}

const selectHighlightFinger = ref<string[]>([])

const fingerOptions = ref<{ label: string, value: string }[]>([])
var highlightFingerOptions = ref<{ label: string, value: string }[]>([])

function deletePoc(pocName: string) {
    ElMessageBox.confirm(
        '确定删除该POC?',
        '警告',
        {
            confirmButtonText: '确认',
            cancelButtonText: '取消',
            type: 'warning',
        }
    )
        .then(async () => {
            let filepath = global.PATH.homedir + "/slack/config/pocs/" + pocName + ".yaml"
            if (await RemoveFile(filepath)) {
                ElMessage.success("Poc deleted successfully")
                pagination.table.result = pagination.table.result.filter(item => item.Name != pocName)
            } else {
                ElMessage.warning("Poc delete failed")
            }
        })
        .catch(() => {

        })

}

const activeTabs = ref("poc")
function addFingerprint(fingerprints: string[]) {
    if (fingerprints.length == 0) {
        ElMessage({
            message: "请至少选择一个指纹",
            showClose: true,
            duration: 2000,
            plain: true,
            grouping: true
        })
        return
    }
    ElMessage.success("添加成功")
    global.webscan.highlight_fingerprints.push(...fingerprints)
}
function deleteFingerprint(fingerprint: string) {
    ElMessage.success("删除成功")
    highlightFingerOptions.value.push({
        label: fingerprint,
        value: fingerprint
    })
    global.webscan.highlight_fingerprints.splice(global.webscan.highlight_fingerprints.indexOf(fingerprint), 1)
}

const editorContainer = ref<HTMLElement | null>(null);
const editor = ref<any | null>(null);
const isEditable = ref(false);
const hasUnsavedChanges = ref(false); // 标记是否有未保存的更改

async function initEditor() {
    if (editorContainer.value && !editor.value) {
        editor.value = monaco.editor.create(editorContainer.value, {
            value: content.value,
            language: 'yaml',
            theme: 'vs-dark',
            readOnly: !isEditable.value,
            automaticLayout: true,
            minimap: { enabled: false }, // 禁用迷你地图
            scrollBeyondLastLine: false, // 禁用滚动到最后一行后额外空白
            contextmenu: false, // 禁用右键菜单
        });

        // 监听编辑器内容变化
        editor.value.onDidChangeModelContent(() => {
            hasUnsavedChanges.value = true; // 标记为有未保存的更改
        });
    }
};

function disposeEditor() {
    if (editor.value) {
        // 解除所有事件监听器
        editor.value.getModel()?.dispose();
        editor.value.dispose(); // 销毁编辑器实例
        editor.value = null;
    }
};

function toggleEditable() {
    isEditable.value = !isEditable.value;
    if (editor.value) {
        editor.value.updateOptions({ readOnly: !isEditable.value });
    }
};


const detailDialog = ref(false)
const content = ref('')
const currentFilepath = ref(""); // 记录当前被编辑的文件路径，方便后续保存时调用

async function readPocFile(filename: string) {
    detailDialog.value = true
    hasUnsavedChanges.value = false
    isEditable.value = false
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + filename + ".yaml"
    let isStat = await CheckFileStat(filepath)
    if (!isStat) {
        filepath = global.webscan.append_pocfile + "/" + filename + ".yaml"
    }
    currentFilepath.value = filepath
    let file = await ReadFile(filepath)
    content.value = file.Content
    // 使用 nextTick 确保 content 更新后再初始化编辑器
    nextTick(() => {
        if (detailDialog.value) {
            initEditor();
        }
    });
}

function handleBeforeClose(done: () => void) {
    // 如果有未保存的更改，提示保存
    if (hasUnsavedChanges.value) {
        ElMessageBox.confirm(
            '您有未保存的更改，是否保存？',
            '未保存的更改',
            {
                confirmButtonText: '保存',
                cancelButtonText: '放弃',
                type: 'warning',
            }
        )
            .then(() => {
                saveContent()
                    .then(() => {
                        done(); // 保存成功后关闭
                    })
                    .catch(() => {
                        ElMessage.error('保存失败，请重试');
                    });
            })
            .catch(() => {
                hasUnsavedChanges.value = false; // 放弃更改，标记无更改
                done(); // 放弃后关闭
            });
    } else {
        done(); // 如果没有未保存更改，直接关闭
    }
}

async function saveContent() {
    try {
        // 当前编辑器状态未改动时，直接返回不用处理保存功能
        if (!hasUnsavedChanges.value) {
            return
        }
        // 已知BUG，需要使用 toRaw 获取原始对象，直接获取会导致CPU占用率到100%，
        const updatedContent = toRaw(editor.value).getValue() || content.value;
        let isSuccess = await WriteFile("yaml", currentFilepath.value, updatedContent);
        if (!isSuccess) {
            ElMessage.error("保存失败")
            return
        }
        hasUnsavedChanges.value = false; // 重置未保存状态
        ElMessage.success('保存成功!');
    } catch (error) {
        ElMessage.error('保存失败，请检查文件路径或权限!');
    }
};

watch(detailDialog, (newValue) => {
    if (!newValue) {
        disposeEditor(); // 在关闭时销毁实例
    }
});

onUnmounted(() => {
    disposeEditor();
});
</script>

<template>
    <CustomTabs>
        <el-tabs v-model="activeTabs" type="card" class="demo-tabs">
            <el-tab-pane name="poc" label="POC管理">
                <el-card>
                    <div style="margin-bottom: 10px;">
                        <el-input :suffix-icon="Search" v-model="filter" @input="filterPocList()"
                            placeholder="根据规则过滤POC" style="width: 50%;">
                            <template #prepend>
                                <el-select v-model="defaultFilter" style="width: 150px;">
                                    <el-option v-for="item in filterOptions" :key="item.value" :label="item.label"
                                        :value="item.value">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-input>
                    </div>
                    <el-table :data="pagination.table.pageContent" style="height: calc(100vh - 225px);">
                        <el-table-column prop="Name" label="名称" />
                        <el-table-column label="关联指纹" width="400px">
                            <template #default="scope">
                                <div class="finger-container">
                                    <el-tag v-for="item in scope.row.AssociatedFingerprint">{{ item }}</el-tag>
                                </div>
                            </template>
                        </el-table-column>
                        <el-table-column label="Operate" width="200px" align="center">
                            <template #default="scope">
                                <el-button size="small" :icon="Reading"
                                    @click="readPocFile(scope.row.Name)">详情</el-button>
                                <el-popconfirm title="Are you sure to delete poc?" @confirm="deletePoc(scope.row.Name)">
                                    <template #reference>
                                        <el-button type="danger" plain :icon="Delete" size="small">删除</el-button>
                                    </template>
                                </el-popconfirm>
                            </template>
                        </el-table-column>
                    </el-table>
                    <div class="my-header" style="margin-top: 5px;">
                        <div></div>
                        <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
                            @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                            :current-page="pagination.table.currentPage" :page-sizes="[20, 50, 100]"
                            :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next, jumper"
                            :total="pagination.table.result.length">
                        </el-pagination>
                    </div>
                </el-card>
            </el-tab-pane>
            <el-tab-pane name="finger" label="指纹管理">
                <el-card>
                    <div class="my-header" style="margin-bottom: 10px;">
                        <el-select-v2 v-model="selectHighlightFinger" placeholder="请选择需要高亮的指纹" filterable
                            :options="highlightFingerOptions" multiple clearable />
                        <el-space style="margin-left: 10px">
                            <el-button :icon="CirclePlusFilled" type="primary"
                                @click="addFingerprint(selectHighlightFinger)">添加</el-button>
                            <el-button :icon="saveIcon" type="primary" @click="SaveConfig">保存</el-button>
                        </el-space>
                    </div>
                    <el-table :data="global.webscan.highlight_fingerprints" style="height: calc(100vh - 195px);">
                        <el-table-column type="index" width="50" />
                        <el-table-column label="Fingerprint">
                            <template #default="scope">
                                {{ scope.row }}
                            </template>
                        </el-table-column>
                        <el-table-column label="Operate" width="150" align="center">
                            <template #default="scope">
                                <el-button type="danger" plain size="small" :icon="Delete"
                                    @click="deleteFingerprint(scope.row)">删除</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-tag :hit="true" v-show="activeTabs == 'finger'">指纹添加或删除成功后记得保存</el-tag>
        </template>
    </CustomTabs>
    <el-drawer v-model="detailDialog" title="漏洞详情" size="70%" :before-close="handleBeforeClose">
        <div class="editor-container">
            <!-- 操作区 -->
            <div class="editor-toolbar">
                <el-space>
                    <el-tooltip content="点击切换状态">
                        <el-button :icon="isEditable ? Edit : Reading" size="small" @click="toggleEditable">
                        {{ isEditable ? '当前状态: 编辑' : '当前状态: 只读' }}
                    </el-button>
                    </el-tooltip>
                    <el-tag type="warning" v-show="hasUnsavedChanges">未保存</el-tag>
                </el-space>
                <el-space>
                    <el-button :icon="DocumentCopy" size="small" @click="Copy(content)">复制</el-button>
                    <el-button :icon="saveIcon" size="small" @click="saveContent">保存</el-button>
                </el-space>
            </div>
            <!-- Monaco 编辑器容器 -->
            <div ref="editorContainer" class="monaco-editor"></div>
        </div>
    </el-drawer>
</template>

<style>
.editor-container {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.editor-toolbar {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 8px;
    background-color: #2D2D2D;
    border-bottom: 1px solid #0E1116;
}

.monaco-editor {
    flex-grow: 1;
    height: calc(100% - 60px);
    /* Adjust based on toolbar height */
}
</style>
