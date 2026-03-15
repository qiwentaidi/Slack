<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { CopyDocument } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { Copy } from "@/util";

type CompareTab = {
    id: string;
    name: string;
    text1: string;
    text2: string;
};

type DiffType = "equal" | "delete" | "insert" | "modify";

type LineDiff = {
    line1: string | null;
    line2: string | null;
    type: DiffType;
};

type CharDiff = {
    char1: string | null;
    char2: string | null;
    type: DiffType;
};

let tabCounter = 1;

const createTab = () => ({
    id: `compare-tab-${Date.now()}-${tabCounter}`,
    name: `标签-${tabCounter++}`,
    text1: "",
    text2: "",
});

const tabs = reactive<CompareTab[]>([createTab()]);
const activeTabId = ref(tabs[0].id);
const currentTab = computed(() => tabs.find((tab) => tab.id === activeTabId.value) || tabs[0]);

const computeDiff = (text1: string, text2: string): LineDiff[] => {
    const lines1 = text1.split("\n");
    const lines2 = text2.split("\n");
    const maxLines = Math.max(lines1.length, lines2.length);
    const diff: LineDiff[] = [];

    for (let index = 0; index < maxLines; index++) {
        const line1 = index < lines1.length ? lines1[index] : null;
        const line2 = index < lines2.length ? lines2[index] : null;

        if (line1 === null && line2 !== null) {
            diff.push({ line1: null, line2, type: "insert" });
            continue;
        }

        if (line1 !== null && line2 === null) {
            diff.push({ line1, line2: null, type: "delete" });
            continue;
        }

        if (line1 === line2) {
            diff.push({ line1, line2, type: "equal" });
            continue;
        }

        diff.push({ line1, line2, type: "modify" });
    }

    return diff;
};

const computeCharDiff = (line1: string, line2: string): CharDiff[] => {
    const chars1 = line1.split("");
    const chars2 = line2.split("");
    const maxLength = Math.max(chars1.length, chars2.length);
    const diff: CharDiff[] = [];

    for (let index = 0; index < maxLength; index++) {
        const char1 = index < chars1.length ? chars1[index] : null;
        const char2 = index < chars2.length ? chars2[index] : null;

        if (char1 === null && char2 !== null) {
            diff.push({ char1: null, char2, type: "insert" });
            continue;
        }

        if (char1 !== null && char2 === null) {
            diff.push({ char1, char2: null, type: "delete" });
            continue;
        }

        if (char1 === char2) {
            diff.push({ char1, char2, type: "equal" });
            continue;
        }

        diff.push({ char1, char2, type: "modify" });
    }

    return diff;
};

const lineDiff = computed(() => computeDiff(currentTab.value.text1, currentTab.value.text2));

const copyText = async (text: string, label: string) => {
    if (!text) {
        ElMessage.warning(`${label}暂无内容可复制`);
        return;
    }
    await Copy(text);
    ElMessage.success(`${label}已复制`);
};

const addTab = () => {
    const tab = createTab();
    tabs.push(tab);
    activeTabId.value = tab.id;
};

const removeTab = (targetId: string) => {
    if (tabs.length === 1) {
        ElMessage.warning("至少保留一个标签");
        return;
    }

    const index = tabs.findIndex((tab) => tab.id === targetId);
    if (index < 0) return;

    tabs.splice(index, 1);
    if (activeTabId.value === targetId) {
        activeTabId.value = tabs[Math.max(index - 1, 0)]?.id || tabs[0].id;
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

const leftCharDiff = (item: LineDiff) =>
    item.line1 !== null && item.line2 !== null && item.type === "modify"
        ? computeCharDiff(item.line1, item.line2)
        : null;
</script>

<template>
    <div class="compare-page">
        <div class="compare-tabs">
            <el-tabs v-model="activeTabId" type="card" editable @edit="handleTabEdit">
                <el-tab-pane
                    v-for="tab in tabs"
                    :key="tab.id"
                    :name="tab.id"
                    :label="tab.name"
                    :closable="tabs.length > 1"
                />
            </el-tabs>
        </div>

        <div class="compare-body">
            <section class="input-grid">
                <article class="surface-card editor-card">
                    <div class="card-head">
                        <h3>文本 1</h3>
                        <el-tooltip content="复制文本 1" placement="top">
                            <el-button size="small" text :icon="CopyDocument" @click="copyText(currentTab.text1, '文本 1')" />
                        </el-tooltip>
                    </div>
                    <el-input
                        v-model="currentTab.text1"
                        type="textarea"
                        resize="none"
                        class="editor-area"
                        placeholder="请输入第一段文本"
                    />
                </article>

                <article class="surface-card editor-card">
                    <div class="card-head">
                        <h3>文本 2</h3>
                        <el-tooltip content="复制文本 2" placement="top">
                            <el-button size="small" text :icon="CopyDocument" @click="copyText(currentTab.text2, '文本 2')" />
                        </el-tooltip>
                    </div>
                    <el-input
                        v-model="currentTab.text2"
                        type="textarea"
                        resize="none"
                        class="editor-area"
                        placeholder="请输入第二段文本"
                    />
                </article>
            </section>

            <section class="surface-card result-card">
                <div class="card-head">
                    <h3>对比结果</h3>
                </div>

                <div v-if="!currentTab.text1 && !currentTab.text2" class="result-empty">
                    输入两段文本后会在这里显示差异结果
                </div>

                <div v-else class="result-view">
                    <div class="result-column">
                        <div class="result-column-head">文本 1</div>
                        <div
                            v-for="(item, index) in lineDiff"
                            :key="`left-${index}`"
                            class="diff-line"
                            :class="{
                                delete: item.type === 'delete',
                                modify: item.type === 'modify',
                                spacer: item.type === 'insert',
                            }"
                        >
                            <template v-if="item.type === 'insert'">
                                <span>&nbsp;</span>
                            </template>
                            <template v-else-if="leftCharDiff(item)">
                                <span
                                    v-for="(charItem, charIndex) in leftCharDiff(item)"
                                    :key="`left-char-${index}-${charIndex}`"
                                    :class="{ 'char-change-delete': charItem.type === 'delete' || charItem.type === 'modify' }"
                                >
                                    {{ charItem.char1 || '' }}
                                </span>
                            </template>
                            <template v-else>
                                <span :class="{ 'line-delete-text': item.type === 'delete' }">{{ item.line1 || "" }}</span>
                            </template>
                        </div>
                    </div>

                    <div class="result-column">
                        <div class="result-column-head">文本 2</div>
                        <div
                            v-for="(item, index) in lineDiff"
                            :key="`right-${index}`"
                            class="diff-line"
                            :class="{
                                insert: item.type === 'insert',
                                modify: item.type === 'modify',
                                spacer: item.type === 'delete',
                            }"
                        >
                            <template v-if="item.type === 'delete'">
                                <span>&nbsp;</span>
                            </template>
                            <template v-else-if="leftCharDiff(item)">
                                <span
                                    v-for="(charItem, charIndex) in leftCharDiff(item)"
                                    :key="`right-char-${index}-${charIndex}`"
                                    :class="{ 'char-change-insert': charItem.type === 'insert' || charItem.type === 'modify' }"
                                >
                                    {{ charItem.char2 || '' }}
                                </span>
                            </template>
                            <template v-else>
                                <span>{{ item.line2 || "" }}</span>
                            </template>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<style scoped>
.compare-page {
    min-height: calc(100vh - 94px);
    display: flex;
    flex-direction: column;
}

.compare-tabs {
    flex: 0 0 auto;
}

.compare-tabs :deep(.el-tabs__nav-wrap::after) {
    display: none;
}

.compare-body {
    flex: 1;
    padding: 14px;
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.surface-card {
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 14px;
}

.input-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
}

.editor-card,
.result-card {
    overflow: hidden;
}

.card-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 14px 16px;
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.card-head h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 700;
    color: var(--el-text-color-primary);
}

.editor-area {
    height: 400px;
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
    font-family: ui-monospace, SFMono-Regular, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    font-size: 13px;
    line-height: 1.6;
}

.result-empty {
    height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--el-text-color-secondary);
}

.result-view {
    height: 400px;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
    padding: 14px;
    overflow: auto;
}

.result-column {
    min-width: 0;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 10px;
    overflow: hidden;
}

.result-column-head {
    position: sticky;
    top: 0;
    z-index: 1;
    padding: 10px 12px;
    font-size: 12px;
    font-weight: 700;
    color: var(--el-text-color-secondary);
    background: var(--el-fill-color-light);
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.diff-line {
    min-height: 22px;
    padding: 3px 12px;
    font-size: 13px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
    font-family: ui-monospace, SFMono-Regular, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.diff-line.delete {
    background: color-mix(in srgb, var(--el-color-danger) 16%, transparent);
}

.diff-line.insert {
    background: color-mix(in srgb, var(--el-color-success) 16%, transparent);
}

.diff-line.modify {
    background: color-mix(in srgb, var(--el-color-warning) 16%, transparent);
}

.diff-line.spacer {
    background: transparent;
}

.line-delete-text {
    text-decoration: line-through;
    opacity: 0.7;
}

.char-change-delete {
    background: color-mix(in srgb, var(--el-color-danger) 30%, transparent);
}

.char-change-insert {
    background: color-mix(in srgb, var(--el-color-success) 30%, transparent);
}

@media (max-width: 1100px) {
    .compare-page {
        height: auto;
        min-height: auto;
    }

    .input-grid,
    .result-view {
        grid-template-columns: 1fr;
    }
}
</style>
