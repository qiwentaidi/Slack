<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { CloseBold, CopyDocument, Delete, DocumentAdd, FolderOpened } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import CustomTextarea from "@/components/CustomTextarea.vue";
import { Copy } from "@/util";
import { CheckFileStat, FileDialog, ReadFile } from "wailsjs/go/services/File";
import { codecCategories, codecOperations, codecOperationMap, createAddedCodecOperation, runCodecPipeline, type AddedCodecOperation, type OperationMode } from "@/lib/codec";

type CodecWorkspaceTab = {
    id: string;
    title: string;
    search: string;
    input: string;
    output: string;
    tips: string;
    running: boolean;
    steps: AddedCodecOperation[];
};

let runToken = 0;
let tabCounter = 1;

const createWorkspaceTab = (): CodecWorkspaceTab => ({
    id: `codec-tab-${Date.now()}-${tabCounter}`,
    title: `标签 ${tabCounter++}`,
    search: "",
    input: "",
    output: "",
    tips: "选择左侧操作后会按顺序处理输入内容，可叠加多步流水线。",
    running: false,
    steps: [],
});

const tabs = reactive<CodecWorkspaceTab[]>([createWorkspaceTab()]);
const activeTabId = ref(tabs[0].id);
const currentTab = computed(() => tabs.find((tab) => tab.id === activeTabId.value) || tabs[0]);

const expandedCategories = ref(codecCategories.map((item) => item.id));
const draggingLibraryOp = ref<{ operationId: string; mode: OperationMode } | null>(null);
const draggingStepIndex = ref<number | null>(null);
const flowDropIndex = ref<number | null>(null);

const filteredOperations = computed(() => {
    const query = currentTab.value.search.trim().toLowerCase();
    if (!query) return codecOperations;
    return codecOperations.filter((item) =>
        item.name.toLowerCase().includes(query) || item.description.toLowerCase().includes(query)
    );
});

const groupedOperations = computed(() =>
    codecCategories.map((category) => ({
        ...category,
        operations: filteredOperations.value.filter((item) => item.category === category.id),
    }))
);

const resolveInputContent = async () => {
    const raw = currentTab.value.input.trim();
    if (!raw) {
        throw new Error("请输入待处理的内容或者文件");
    }
    if (raw.startsWith("{{file://") && raw.endsWith("}}")) {
        let filepath = raw.replace("{{file://", "").replace("}}", "");
        const exists = await CheckFileStat(filepath);
        if (!exists) {
            throw new Error("文件不存在");
        }
        const file = await ReadFile(filepath);
        return file.Content || "";
    }
    return currentTab.value.input;
};

const executePipeline = async () => {
    const targetTab = currentTab.value;
    const currentToken = ++runToken;
    targetTab.running = true;
    try {
        const input = await resolveInputContent();
        const output = await runCodecPipeline(input, targetTab.steps);
        if (currentToken !== runToken) return;
        targetTab.output = output;
        targetTab.tips = targetTab.steps.length ? `执行完成，共运行 ${targetTab.steps.length} 个步骤` : "没有添加步骤，结果保持原始输入";
    } catch (error) {
        if (currentToken !== runToken) return;
        targetTab.output = "";
        targetTab.tips = error instanceof Error ? error.message : "执行失败";
        ElMessage.warning(targetTab.tips);
    } finally {
        if (currentToken === runToken) {
            targetTab.running = false;
        }
    }
};

const hasRunnableInput = computed(() => {
    const raw = currentTab.value.input.trim();
    return Boolean(raw);
});

watch(
    () => [activeTabId.value, currentTab.value.input, JSON.stringify(currentTab.value.steps)],
    () => {
        if (!hasRunnableInput.value) {
            currentTab.value.output = "";
            currentTab.value.running = false;
            return;
        }
        executePipeline();
    },
    { deep: true }
);

const addOperation = (operationId: string, mode: OperationMode, index?: number) => {
    const step = createAddedCodecOperation(operationId, mode);
    if (!step) return;
    if (typeof index === "number") {
        currentTab.value.steps.splice(index, 0, step);
        return;
    }
    currentTab.value.steps.push(step);
};

const removeOperation = (index: number) => {
    currentTab.value.steps.splice(index, 1);
};

const updateOption = (step: AddedCodecOperation, key: string, value: string) => {
    step.options[key] = value;
};

const uploadFile = async () => {
    const filepath = await FileDialog("");
    if (!filepath) return;
    currentTab.value.input = `{{file://${filepath}}}`;
};

const resetAll = () => {
    currentTab.value.input = "";
    currentTab.value.output = "";
    currentTab.value.steps = [];
    currentTab.value.tips = "已清空输入、输出和流水线";
};

const clearPipeline = () => {
    currentTab.value.steps = [];
};

const startLibraryDrag = (operationId: string, mode: OperationMode) => {
    draggingLibraryOp.value = {
        operationId,
        mode,
    };
    draggingStepIndex.value = null;
};

const startStepDrag = (index: number) => {
    draggingStepIndex.value = index;
    draggingLibraryOp.value = null;
};

const clearDragState = () => {
    draggingLibraryOp.value = null;
    draggingStepIndex.value = null;
    flowDropIndex.value = null;
};

const setFlowDropIndex = (index: number) => {
    flowDropIndex.value = index;
};

const handleFlowDrop = () => {
    if (draggingLibraryOp.value) {
        addOperation(draggingLibraryOp.value.operationId, draggingLibraryOp.value.mode, flowDropIndex.value ?? currentTab.value.steps.length);
        clearDragState();
        return;
    }
    if (draggingStepIndex.value === null) return;
    const from = draggingStepIndex.value;
    const rawTarget = flowDropIndex.value ?? currentTab.value.steps.length;
    const target = rawTarget > from ? rawTarget - 1 : rawTarget;
    if (target === from || target < 0 || target > currentTab.value.steps.length - 1) {
        clearDragState();
        return;
    }
    const [item] = currentTab.value.steps.splice(from, 1);
    currentTab.value.steps.splice(target, 0, item);
    clearDragState();
};

const handleDropOutsideFlow = () => {
    if (draggingStepIndex.value === null) return;
    currentTab.value.steps.splice(draggingStepIndex.value, 1);
    clearDragState();
};

const useOutputAsInput = () => {
    if (!currentTab.value.output) return;
    currentTab.value.input = currentTab.value.output;
    currentTab.value.tips = "已将结果写回输入区";
};

const addTab = () => {
    const tab = createWorkspaceTab();
    tabs.push(tab);
    activeTabId.value = tab.id;
};

const removeTab = (targetId: string) => {
    if (tabs.length === 1) {
        resetAll();
        return;
    }
    const index = tabs.findIndex((tab) => tab.id === targetId);
    if (index < 0) return;
    tabs.splice(index, 1);
    if (activeTabId.value === targetId) {
        activeTabId.value = tabs[Math.max(0, index - 1)]?.id || tabs[0].id;
    }
};

const handleTabEdit = (targetName: string | number | undefined, action: "add" | "remove") => {
    if (action === "add") {
        addTab();
        return;
    }
    if (action === "remove" && typeof targetName === "string") {
        removeTab(targetName);
    }
};

const stepName = (step: AddedCodecOperation) => codecOperationMap.get(step.operationId)?.name || step.operationId;
const stepOptions = (step: AddedCodecOperation) => codecOperationMap.get(step.operationId)?.options || [];
</script>

<template>
    <div class="codec-page">
        <div class="workspace-tabs">
            <el-tabs v-model="activeTabId" type="card" editable @edit="handleTabEdit">
                <el-tab-pane v-for="tab in tabs" :key="tab.id" :name="tab.id" :label="tab.title" :closable="tabs.length > 1" />
            </el-tabs>
        </div>

        <el-splitter class="workspace-splitter">
            <el-splitter-panel
                size="20%"
                min="240px"
                class="pane panel-surface library-panel"
                :class="{ 'drop-delete-target': draggingStepIndex !== null }"
                @dragover.prevent="draggingStepIndex !== null"
                @drop.prevent="handleDropOutsideFlow"
            >
                <div class="library-body">
                    <div class="library-search">
                        <el-input v-model="currentTab.search" placeholder="搜索操作..." clearable />
                    </div>

                    <div class="library-list">
                        <el-collapse v-model="expandedCategories" class="op-groups">
                            <el-collapse-item
                                v-for="group in groupedOperations"
                                :key="group.id"
                                :name="group.id"
                                v-show="group.operations.length > 0"
                            >
                                <template #title>
                                    <div class="group-title">
                                        <span>{{ group.name }}</span>
                                        <el-tag size="small" effect="plain">{{ group.operations.length }}</el-tag>
                                    </div>
                                </template>

                                <div class="op-list compact">
                                    <article
                                        v-for="operation in group.operations"
                                        :key="operation.id"
                                        class="op-row clickable"
                                        :draggable="Boolean(operation.transform || (operation.encode && !operation.decode) || (!operation.encode && operation.decode))"
                                        @click="operation.transform ? addOperation(operation.id, 'transform') : operation.encode && !operation.decode ? addOperation(operation.id, 'encode') : !operation.encode && operation.decode ? addOperation(operation.id, 'decode') : undefined"
                                        @dragstart="
                                            operation.transform
                                                ? startLibraryDrag(operation.id, 'transform')
                                                : operation.encode && !operation.decode
                                                  ? startLibraryDrag(operation.id, 'encode')
                                                  : !operation.encode && operation.decode
                                                    ? startLibraryDrag(operation.id, 'decode')
                                                    : undefined
                                        "
                                        @dragend="clearDragState"
                                    >
                                        <el-tooltip :content="operation.name" placement="top" :show-after="300">
                                            <div class="op-copy compact">
                                                <strong>{{ operation.name }}</strong>
                                            </div>
                                        </el-tooltip>
                                        <div class="op-actions inline">
                                            <el-button
                                                v-if="operation.encode && operation.decode"
                                                size="small"
                                                text
                                                type="success"
                                                draggable="true"
                                                @click.stop="addOperation(operation.id, 'encode')"
                                                @dragstart="startLibraryDrag(operation.id, 'encode')"
                                                @dragend="clearDragState"
                                            >
                                                编码
                                            </el-button>
                                            <el-button
                                                v-if="operation.encode && operation.decode"
                                                size="small"
                                                text
                                                type="warning"
                                                draggable="true"
                                                @click.stop="addOperation(operation.id, 'decode')"
                                                @dragstart="startLibraryDrag(operation.id, 'decode')"
                                                @dragend="clearDragState"
                                            >
                                                解码
                                            </el-button>
                                        </div>
                                    </article>
                                </div>
                            </el-collapse-item>
                        </el-collapse>
                    </div>
                </div>
            </el-splitter-panel>

            <el-splitter-panel size="20%" min="220px" class="pane panel-surface">
                <div class="pane-head flow-head">
                    <h3>操作流程</h3>
                    <div class="toolbar single-line">
                        <el-button size="small" :icon="Delete" text @click="clearPipeline" />
                    </div>
                </div>

                <div v-if="currentTab.steps.length" class="flow-body">
                    <div
                        class="flow-list"
                        @dragover.prevent="setFlowDropIndex(currentTab.steps.length)"
                        @drop.prevent="handleFlowDrop"
                    >
                        <template v-for="(step, index) in currentTab.steps" :key="step.id">
                            <div
                                class="flow-dropzone"
                                :class="{ active: flowDropIndex === index }"
                                @dragover.prevent="setFlowDropIndex(index)"
                                @drop.prevent="handleFlowDrop"
                            />
                            <article
                                class="flow-card unified"
                                draggable="true"
                                @dragstart="startStepDrag(index)"
                                @dragend="clearDragState"
                                @dragover.stop.prevent="setFlowDropIndex(index)"
                                @drop.stop.prevent="handleFlowDrop"
                            >
                                <div class="flow-card-head">
                                    <div class="flow-card-main">
                                        <button class="drag-handle" type="button" aria-label="拖拽排序">
                                            <span /><span /><span /><span /><span /><span />
                                        </button>
                                        <div class="flow-index">{{ index + 1 }}</div>
                                        <div class="flow-copy">
                                            <strong>{{ stepName(step) }}</strong>
                                        </div>
                                    </div>
                                    <div class="flow-tools">
                                        <el-button size="small" text :icon="CloseBold" @click="removeOperation(index)" />
                                    </div>
                                </div>

                                <div v-if="stepOptions(step).length" class="step-options single-column integrated">
                                    <div v-for="option in stepOptions(step)" :key="option.key" class="option-item">
                                        <span>{{ option.label }}</span>
                                        <el-select
                                            v-if="option.type === 'select'"
                                            :model-value="step.options[option.key]"
                                            @update:model-value="(value) => updateOption(step, option.key, String(value))"
                                        >
                                            <el-option v-for="item in option.options || []" :key="item" :label="item" :value="item" />
                                        </el-select>
                                        <el-input
                                            v-else
                                            :model-value="step.options[option.key]"
                                            @update:model-value="(value) => updateOption(step, option.key, value)"
                                        />
                                    </div>
                                </div>
                            </article>
                        </template>
                        <div
                            class="flow-dropzone tail"
                            :class="{ active: flowDropIndex === currentTab.steps.length }"
                            @dragover.prevent="setFlowDropIndex(currentTab.steps.length)"
                            @drop.prevent="handleFlowDrop"
                        />
                    </div>

                </div>
                <div v-else class="flow-empty">
                    <strong>从左侧添加操作</strong>
                    <span>支持链式处理</span>
                    <span>支持拖拽删除</span>
                </div>
            </el-splitter-panel>

            <el-splitter-panel
                class="pane panel-surface"
                :class="{ 'drop-delete-target': draggingStepIndex !== null }"
                @dragover.prevent="draggingStepIndex !== null"
                @drop.prevent="handleDropOutsideFlow"
            >
                <el-splitter layout="vertical" class="io-splitter">
                    <el-splitter-panel size="52%" min="220px" class="io-pane">
                        <div class="pane-head io-head">
                            <h3>输入</h3>
                            <div class="toolbar single-line">
                                <el-tooltip content="选择文件" placement="top">
                                    <el-button size="small" :icon="FolderOpened" text @click="uploadFile" />
                                </el-tooltip>
                                <el-tooltip content="复制输入" placement="top">
                                    <el-button size="small" :icon="CopyDocument" text @click="Copy(currentTab.input)" />
                                </el-tooltip>
                                <el-tooltip content="清空输入" placement="top">
                                    <el-button size="small" :icon="Delete" text @click="currentTab.input = ''" />
                                </el-tooltip>
                            </div>
                        </div>
                        <el-input
                            v-model="currentTab.input"
                            type="textarea"
                            resize="none"
                            class="editor-area"
                            placeholder="输入待处理内容。"
                        />
                    </el-splitter-panel>

                    <el-splitter-panel class="io-pane">
                        <div class="pane-head io-head">
                            <h3>输出</h3>
                            <div class="toolbar single-line">
                                <el-tooltip content="复制输出" placement="top">
                                    <el-button size="small" :icon="CopyDocument" text @click="Copy(currentTab.output)" />
                                </el-tooltip>
                                <el-tooltip content="转为输入" placement="top">
                                    <el-button size="small" :icon="DocumentAdd" text @click="useOutputAsInput" />
                                </el-tooltip>
                                <el-tooltip content="全部清空" placement="top">
                                    <el-button size="small" :icon="Delete" text @click="resetAll" />
                                </el-tooltip>
                            </div>
                        </div>
                        <CustomTextarea
                            v-model="currentTab.output"
                            class="editor-area output-area"
                            :rows="10"
                            :readonly="true"
                            :hide-readonly-action="true"
                            resize="none"
                        />
                    </el-splitter-panel>
                </el-splitter>
            </el-splitter-panel>
        </el-splitter>
    </div>
</template>

<style scoped>
.codec-page {
    height: calc(100vh - 94px);
    min-height: 680px;
    display: flex;
    flex-direction: column;
}

.workspace-tabs {
    flex: 0 0 auto;
}

.workspace-tabs :deep(.el-tabs__nav-wrap::after) {
    display: none;
}

.workspace-splitter,
.io-splitter {
    height: 100%;
}

.workspace-splitter {
    flex: 1;
    min-height: 0;
}

.pane,
.io-pane {
    min-width: 0;
    min-height: 0;
}

.panel-surface {
    height: 100%;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 14px;
    overflow: hidden;
}

.pane-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 14px 16px;
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.pane-head h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 700;
    color: var(--el-text-color-primary);
}

.library-body,
.flow-list {
    height: calc(100% - 66px);
    padding: 12px;
}

.library-body {
    display: flex;
    flex-direction: column;
    gap: 0;
}

.library-panel .library-body {
    height: 100%;
}

.library-search {
    flex: 0 0 auto;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.library-list {
    flex: 1;
    min-height: 0;
    overflow: auto;
    padding-top: 12px;
    scrollbar-gutter: stable;
    padding-right: 18px;
}

.op-groups {
    border: none;
}

.group-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    font-weight: 700;
    padding-right: 28px;
}

.flow-body {
    height: calc(100% - 66px);
    display: flex;
    flex-direction: column;
    min-height: 0;
}

.flow-list {
    flex: 1;
    min-height: 0;
    overflow: auto;
}

.flow-empty {
    height: calc(100% - 66px);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    color: var(--el-text-color-secondary);
}

.flow-empty strong {
    color: var(--el-text-color-primary);
    font-size: 16px;
}

.flow-empty span {
    font-size: 13px;
}

.op-list.compact {
    display: flex;
    flex-direction: column;
}

.op-groups :deep(.el-collapse-item__header) {
    display: flex;
    align-items: center;
    gap: 8px;
}

.op-groups :deep(.el-collapse-item__arrow) {
    order: -1;
    margin: 0 2px 0 0;
}

.op-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding: 10px 12px;
    border-top: none;
    border-radius: 10px;
    cursor: pointer;
    transition: background-color 0.18s ease, color 0.18s ease;
}

.op-row:first-child {
    border-top: none;
}

.op-row.clickable:hover {
    background: var(--el-fill-color);
    color: var(--el-color-primary);
}

.op-copy.compact {
    min-width: 0;
    flex: 1;
}

.op-copy.compact strong {
    display: block;
    color: var(--el-text-color-primary);
    font-size: 14px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.op-actions.inline,
.toolbar.single-line,
.flow-tools,
.flow-tags {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: nowrap;
}

.op-actions.inline {
    gap: 2px;
}

.op-actions.inline :deep(.el-button) {
    padding-inline: 4px;
}

.flow-card {
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 12px;
    background: var(--el-fill-color-light);
    margin-bottom: 0;
}

.flow-card.unified {
    padding: 10px 12px;
}

.flow-card-head,
.flow-card-main {
    display: flex;
    align-items: center;
    gap: 10px;
}

.flow-card-head {
    justify-content: space-between;
}

.drag-handle {
    width: 16px;
    height: 28px;
    display: grid;
    grid-template-columns: repeat(2, 4px);
    grid-template-rows: repeat(3, 4px);
    gap: 2px;
    padding: 0;
    border: 0;
    background: transparent;
    cursor: grab;
    align-self: center;
    place-content: center;
    margin-top: 0;
}

.drag-handle span {
    display: block;
    width: 4px;
    height: 4px;
    border-radius: 999px;
    background: var(--el-text-color-secondary);
}

.flow-dropzone {
    height: 4px;
    border-radius: 999px;
    margin: 0 4px 4px;
    transition: background 0.2s ease;
}

.flow-dropzone.active {
    background: color-mix(in srgb, var(--el-color-primary) 35%, transparent);
}

.drop-delete-target {
    transition: border-color 0.18s ease, box-shadow 0.18s ease;
    border-color: color-mix(in srgb, var(--el-color-danger) 20%, var(--el-border-color-lighter));
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--el-color-danger) 10%, transparent);
}

.drop-delete-target .pane-head h3 {
    color: var(--el-color-danger);
}

.flow-index {
    width: 28px;
    height: 28px;
    border-radius: 8px;
    background: var(--el-fill-color);
    color: var(--el-text-color-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: 700;
}

.flow-copy {
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.flow-copy strong {
    font-size: 14px;
    color: var(--el-text-color-primary);
}

.step-options.single-column {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.step-options.single-column.integrated {
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid var(--el-border-color-lighter);
}

.option-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.option-item span {
    font-size: 12px;
    color: var(--el-text-color-secondary);
}

.editor-area {
    height: calc(100% - 67px);
}

.editor-area :deep(.el-textarea),
.editor-area :deep(.el-textarea__inner) {
    height: 100%;
}

.editor-area :deep(.el-textarea__inner) {
    resize: none;
    border: none;
    border-radius: 0;
    box-shadow: none;
    padding: 16px;
    background: transparent;
}

.output-area {
    height: calc(100% - 67px);
}

.io-pane {
    height: 100%;
    background: var(--el-bg-color);
}

@media (max-width: 1100px) {
    .codec-page {
        height: auto;
        min-height: auto;
    }
}
</style>
