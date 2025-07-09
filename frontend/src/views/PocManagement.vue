<script lang="ts" setup>
import { watch, onMounted, onUnmounted, ref, nextTick, toRaw, computed } from 'vue';
import { Search, CirclePlusFilled, Delete, DocumentCopy, Reading, Minus, EditPen, CollectionTag } from "@element-plus/icons-vue";
import { CheckFileStat, ReadFile, RemoveFile, SaveFileDialog, WriteFile } from 'wailsjs/go/services/File';
import global from '@/stores';
import { FingerprintList, GetFingerPocMap } from 'wailsjs/go/services/App';
import { Copy } from '@/util';
import { Matcher, PocDetail } from '@/stores/interface';
import usePagination from '@/usePagination';
import { ElMessage, ElMessageBox } from 'element-plus';
import CustomTabs from '@/components/CustomTabs.vue';
import saveIcon from '@/assets/icon/save.svg'
import aiIcon from '@/assets/icon/ai.svg'
import { SaveConfig } from '@/config';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import "monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution";
import { dslOptions, metadataOptions, pocdetailFilterOptions, sortSeverityOptions } from '@/stores/options';
import { FormData } from '@/stores/interface';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';

onMounted(async () => {
    const pocMap = await GetFingerPocMap();
    for (const [poc, tags] of Object.entries(pocMap)) {
        pagination.table.result.push({
            Name: poc,
            AssociatedFingerprint: Array.from(new Set(tags))
        })
    }
    pagination.ctrl.watchResultChange(pagination.table)
    pagination.table.filterTemp = pagination.table.result
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

const filter = ref('')
function filterPocList() {
    if (filter.value == '') {
        pagination.table.result = pagination.table.filterTemp
        pagination.ctrl.watchResultChange(pagination.table)
        return
    }
    pagination.table.result = []
    if (defaultFilter.value != "Name") {
        for (const item of pagination.table.filterTemp) {
            for (const finger of item.AssociatedFingerprint) {
                if (finger.toLowerCase().includes(filter.value.toLowerCase())) {
                    pagination.table.result.push(item)
                    break
                }
            }
        }
    } else {
        for (const item of pagination.table.filterTemp) {
            if (item.Name.toLowerCase().includes(filter.value.toLowerCase())) {
                pagination.table.result.push(item)
            }
        }
    }

    pagination.ctrl.watchResultChange(pagination.table)
}

const selectHighlightFinger = ref<string[]>([])

const fingerOptions = ref<{ label: string, value: string }[]>([])
var highlightFingerOptions = ref<{ label: string, value: string }[]>([])

function deletePoc(pocName: string) {
    ElMessageBox.confirm(
        '确定删除该POC?',
        '警告',
        {
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
    SaveConfig()
    selectHighlightFinger.value = []
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

// step 1 add poc, hide poclist
const step = ref(0)

const metadataTemp = ref({}); // 临时存储 key

const formData = ref<FormData>({
    id: '',
    name: '',
    author: '',
    description: '',
    severity: 'medium',
    body: '',
    matchers: [] as Matcher[],
    matchersCondition: 'and'
});
// 控制缩进函数
const indent = (level: number, content: string) => ' '.repeat(level * 2) + content;

const generatedPoc = computed(() => {

    const matchers = formData.value.matchers.map(matcher => {
        let matcherConfig = [indent(3, `- type: ${matcher.type}`)];

        // 处理需要 part 属性的类型
        if (["word", "regex"].includes(matcher.type) && matcher.part !== "all") {
            matcherConfig.push(indent(4, `part: ${matcher.part || "body"}`));
        }

        // 处理不同类型的字段
        const fieldType = matcher.type === "regex" ? "regex" : matcher.type === "word" ? "words" : matcher.type;
        matcherConfig.push(indent(4, `${fieldType}:`));

        // 处理匹配值
        const wrapWithQuotes = !(["size", "status", "binary"].includes(matcher.type));
        matcherConfig.push(
            matcher.words
                .map(word => indent(5, `- ${wrapWithQuotes ? `'${word}'` : word}`))
                .join("\n")
        );

        // 处理 condition
        if (matcher.words.length > 1) {
            matcherConfig.push(indent(4, `condition: ${matcher.condition}`));
        }

        return matcherConfig.join("\n");
    }).join("\n");

    const pocConfig = [
        `id: ${formData.value.id}\n`,
        'info:',
        indent(1, `name: ${formData.value.name}`),
        indent(1, `author: ${formData.value.author}`),
        indent(1, `severity: ${formData.value.severity}`),
        indent(1, `description: |`),
        indent(2, formData.value.description),
    ];

    if (formData.value.reference && Object.keys(formData.value.reference).length > 0) {
        const referenceArray = formData.value.reference.split('\n').filter(ref => ref.trim() !== '');
        if (referenceArray.length > 0) {
            pocConfig.push(indent(1, 'reference:'));
            referenceArray.forEach(ref => {
                pocConfig.push(indent(2, `- ${ref.trim()}`));
            });
        }
    }

    if (formData.value.tags && Object.keys(formData.value.tags).length > 0) {
        pocConfig.push(indent(1, `tags: ${formData.value.tags.join(',')}`))
    }

    // 处理 metadata 部分
    if (formData.value.metadata && Object.keys(formData.value.metadata).length > 0) {
        pocConfig.push(indent(1, 'metadata:'));
        Object.entries(formData.value.metadata).forEach(([key, value]) => {
            if (value) {
                pocConfig.push(indent(2, `${key}: ${value}`));
            }
        });
    }

    pocConfig.push(
        '\nhttp:',
        indent(1, '- raw:'),
        indent(3, '- |'),
    );

    if (formData.value.body && Object.keys(formData.value.body).length > 0) {
        const bodyArray = formData.value.body.split('\n');
        if (bodyArray.length > 0) {
            bodyArray.forEach(ref => {
                pocConfig.push(indent(4, ref.trim()));
            });
        }
    }

    pocConfig.push('\n')

    if (formData.value.matchers.length > 1) {
        pocConfig.push(indent(2, `matchers-condition: ${formData.value.matchersCondition}`));
    }

    pocConfig.push(indent(2, 'matchers:'));
    pocConfig.push(matchers);

    return pocConfig.join('\n');
});

const editorDialog = ref(false);

function openEditor() {
    editorDialog.value = true
    content.value = generatedPoc.value
    isEditable.value = true
    nextTick(() => {
        if (editorDialog.value) {
            initEditor();
        }
    });
}

watch(editorDialog, (newValue) => {
    if (!newValue) {
        disposeEditor(); // 在关闭时销毁实例
    }
});

const nuclei = ({
    fetchMetatedataSuggestions: function (query: string, cb: Function) {
        cb(metadataOptions.filter(k => k.includes(query)).map(k => ({ value: k })))
    },
    updateMetadataKey: function (oldKey: string, newKey: string) {
        if (newKey && oldKey !== newKey && !formData.value.metadata[newKey]) {
            formData.value.metadata[newKey] = formData.value.metadata[oldKey];
            delete formData.value.metadata[oldKey];
            metadataTemp.value[newKey] = newKey;
            delete metadataTemp.value[oldKey];
        }
    },
    addMetadata: function () {
        if (!formData.value.metadata) {
            formData.value.metadata = {};
        }
        const newKey = ''; // 默认空 key
        formData.value.metadata[newKey] = '';
        metadataTemp.value[newKey] = newKey;
    },
    removeMetadata: function (key: string) {
        if (formData.value.metadata) {
            delete formData.value.metadata[key];
            delete metadataTemp.value[key]; // 同步删除临时 key
        }
    },
    fetchDslSuggestions: function (query: string, cb: Function) {
        cb(dslOptions.filter(k => k.text.includes(query)))
    },
    updateMatcherWords: function (index: number) {
        nextTick(() => {
            const matcher = formData.value.matchers[index];
            if (matcher.wordsText) {
                matcher.words = matcher.wordsText.split('\n').filter(word => word.trim());
            } else {
                matcher.words = [];
            }
        })
    },
    addMatcher: function () {
        formData.value.matchers.push({
            type: 'word',
            part: 'all',
            words: [],
            condition: 'and',
            wordsText: '',
        });
    },
    removeMatcher: function (index: number) {
        formData.value.matchers.splice(index, 1);
    },
    replaceHostAll: function () {
        formData.value.body = formData.value.body.replace(/(Host:\s*)([^\n]+)/g, '{{Hostname}}');
    },
})

async function savePoc() {
    if (formData.value.id === "") {
        ElMessage.warning("请输入POC ID!")
        return;
    }
    const path = await SaveFileDialog(formData.value.id + ".yaml");
    if (!path) {
        return;
    }
    const result = await WriteFile("yaml", path, content.value);
    result ? ElMessage.success("保存成功!") : ElMessage.error("保存失败!");
}

const variables = `variables:
  num: "999999999"
  filename: "{{rand_base(8)}}"
  s1: "{{rand_int(40000, 44800)}}"`

const func = `随机整数: {{rand_int(40000, 44800)}}
随机字符串: {{rand_base(8)}} || {{randstr}}
转小写: {{to_lower(rand_text_alpha(5))}}
转大写: {{to_upper(rand_base(12))}}`

const path = `http:
  - method: GET
    path:
      - "{{BaseURL}}/libs/granite/offloading/content/view.html"`

const uploads = `aspx:
  <%@Page Language="C#" %>
        <% Response.Write({{s1}}*{{s2}}); System.IO.File.Delete(Request.PhysicalPath); %>

jsp:
  <%out.print({{s1}}*{{s2}});new java.io.File(application.getRealPath(request.getServletPath())).delete();%>
  
php:
  <?php print('{{randstr}}');unlink(__FILE__);?>`
</script>

<template>
    <CustomTabs v-show="step == 0">
        <el-tabs v-model="activeTabs" type="card">
            <el-tab-pane name="poc" label="POC管理">
                <el-card>
                    <div class="flex-between mb-10px">
                        <el-input :suffix-icon="Search" v-model="filter" @input="filterPocList()"
                            placeholder="根据规则过滤POC" class="w-1/2">
                            <template #prepend>
                                <el-select v-model="defaultFilter" style="width: 150px;">
                                    <el-option v-for="item in pocdetailFilterOptions" :key="item.value"
                                        :label="item.label" :value="item.value">
                                    </el-option>
                                </el-select>
                            </template>
                        </el-input>
                        <el-button class="ml-5px" type="primary" plain :icon="CirclePlusFilled" @click="step = 1">添加POC</el-button>
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
                    <div class="flex-between mt-5px">
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
                    <div class="flex-between mb-10px">
                        <el-select-v2 v-model="selectHighlightFinger" placeholder="请选择需要高亮的指纹" filterable
                            :options="highlightFingerOptions" multiple clearable />
                        <el-button :icon="CirclePlusFilled" type="primary" plain
                            @click="addFingerprint(selectHighlightFinger)" class="ml-5px">添加</el-button>
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
    </CustomTabs>
    <div v-show="step == 1">
        <el-page-header @back="step = 0">
            <template #content>
                <span class="font-bold mr-5px">Nuclei PoC 生成器</span>
                <el-tag>填写完基础信息后, 通过编辑器模式进行调整/保存</el-tag>
            </template>
            <template #extra>
                <el-button circle @click="BrowserOpenURL('https://cloud.projectdiscovery.io/templates')">
                    <template #icon>
                        <el-icon size="24">
                            <aiIcon />
                        </el-icon>
                    </template>
                </el-button>
                <el-button round :icon="EditPen" @click="openEditor">编辑器模式</el-button>
            </template>
        </el-page-header>
        <el-divider />
        <div class="flex gap-2">
            <el-form :model="formData" label-width="auto" class="w-full">
                <el-form-item label="漏洞ID">
                    <el-input v-model="formData.id"></el-input>
                </el-form-item>
                <el-form-item label="漏洞名称">
                    <el-input v-model="formData.name"></el-input>
                </el-form-item>
                <el-form-item label="作者">
                    <el-input v-model="formData.author"></el-input>
                </el-form-item>
                <el-form-item label="风险等级">
                    <el-select v-model="formData.severity">
                        <el-option v-for="severity in sortSeverityOptions" :key="severity" :label="severity"
                            :value="severity.toLowerCase()" />
                    </el-select>
                </el-form-item>
                <el-form-item label="漏洞描述">
                    <el-input v-model="formData.description" type="textarea" :rows="3"></el-input>
                </el-form-item>
                <el-form-item label="漏洞来源">
                    <el-input v-model="formData.reference" type="textarea" :rows="3"></el-input>
                    <span class="form-item-tips">可换行分割多个来源</span>
                </el-form-item>
                <el-form-item label="Matedata">
                    <el-button type="primary" size="small" @click="nuclei.addMetadata">添加 Matedata</el-button>
                    <div v-if="formData.metadata">
                        <div v-for="(_, key) in formData.metadata" :key="key"
                            style="display: flex; align-items: center; margin-bottom: 8px; width: 100%;">
                            <!-- 这里用 el-autocomplete 提供建议 -->
                            <el-autocomplete v-model="metadataTemp[key]"
                                @select="nuclei.updateMetadataKey(key, $event.value)"
                                :fetch-suggestions="nuclei.fetchMetatedataSuggestions" placeholder="请输入 Key"
                                style="flex: 1; margin-right: 8px;"
                                @change="nuclei.updateMetadataKey(key, metadataTemp[key])" />
                            <el-input v-model="formData.metadata[key]" placeholder="请输入 Value"
                                style="flex: 1; margin-right: 8px;" />
                            <el-button :icon="Minus" size="small" circle @click="nuclei.removeMetadata(key)" />
                        </div>
                    </div>
                </el-form-item>
                <el-form-item label="关联指纹">
                    <el-select-v2 v-model="formData.tags" :options="fingerOptions" filterable multiple clearable />
                    <span class="form-item-tips">必须关联指纹库中的至少1条指纹</span>
                </el-form-item>
                <el-form-item label="Raw数据包">
                    <el-input v-model="formData.body" type="textarea" :rows="8" placeholder="请输入请求体内容" />
                </el-form-item>
            </el-form>
            <!-- 右侧预览区域 -->
            <el-form :model="formData" label-width="auto" class="w-full">
                <el-form-item label="匹配规则">
                    <el-select v-model="formData.matchersCondition" class="w-full">
                        <el-option label="AND" value="and" />
                        <el-option label="OR" value="or" />
                    </el-select>
                    <div class="w-full mt-5px">
                        <el-card v-for="(matcher, index) in formData.matchers" :key="index" class="mb-5px">
                            <template #header>
                                <div class="card-header">
                                    <span>规则 #{{ index + 1 }}</span>
                                    <el-button size="small" type="danger" :icon="Delete"
                                        @click="nuclei.removeMatcher(index)">
                                    </el-button>
                                </div>
                            </template>

                            <el-form label-position="top">
                                <el-form-item label="类型">
                                    <el-select v-model="matcher.type">
                                        <el-option label="word" value="word" />
                                        <el-option label="status" value="status" />
                                        <el-option label="regex" value="regex" />
                                        <el-option label="dsl" value="dsl" />
                                        <el-option label="size" value="size" />
                                        <el-option label="binary" value="binary" />
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="匹配部分" v-if="matcher.type === 'word' || matcher.type === 'regex'">
                                    <el-select v-model="matcher.part">
                                        <el-option label="响应体" value="body" />
                                        <el-option label="响应头" value="header" />
                                        <el-option label="全部" value="all" />
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="词条">
                                    <el-autocomplete v-model="matcher.wordsText" type="textarea" :rows="3"
                                        :fetch-suggestions="matcher.type === 'dsl' ? nuclei.fetchDslSuggestions : () => []"
                                        placeholder="请输入词条, 每行一个" @input="nuclei.updateMatcherWords(index)">
                                        <template #default="{ item }">
                                            <div>
                                                <span style="color: #559FF8;">{{ item.text }}</span>
                                                <el-divider direction="vertical" />
                                                <span>{{ item.value }}</span>
                                            </div>
                                        </template>
                                    </el-autocomplete>
                                </el-form-item>

                                <el-form-item label="词条匹配条件">
                                    <el-radio-group v-model="matcher.condition">
                                        <el-radio value="and">AND</el-radio>
                                        <el-radio value="or">OR</el-radio>
                                    </el-radio-group>
                                </el-form-item>
                            </el-form>
                        </el-card>
                        <el-button type="primary" size="small" @click="nuclei.addMatcher">
                            添加匹配规则
                        </el-button>
                    </div>
                </el-form-item>
            </el-form>
        </div>
    </div>
    <el-drawer v-model="detailDialog" size="70%" :before-close="handleBeforeClose">
        <template #header>
            <span class="drawer-title">漏洞详情</span>
        </template>
        <div class="editor-container">
            <!-- 操作区 -->
            <div class="editor-toolbar">
                <el-space>
                    <el-tooltip content="点击切换状态">
                        <el-button link color="#000" :icon="isEditable ? EditPen : Reading" @click="toggleEditable">
                            {{ isEditable ? '当前状态: 编辑' : '当前状态: 只读' }}
                        </el-button>
                    </el-tooltip>
                    <el-tag type="warning" v-show="hasUnsavedChanges">未保存</el-tag>
                </el-space>
                <el-space>
                    <el-button link color="#000" :icon="DocumentCopy" @click="Copy(content)">复制</el-button>
                    <el-button link color="#000" :icon="saveIcon" @click="saveContent">保存</el-button>
                </el-space>
            </div>
            <!-- Monaco 编辑器容器 -->
            <div ref="editorContainer" class="monaco-editor"></div>
        </div>
    </el-drawer>
    <el-drawer v-model="editorDialog" size="70%">
        <template #header>
            <span class="drawer-title">编辑漏洞</span>
        </template>
        <div class="editor-container">
            <!-- 操作区 -->
            <div class="editor-toolbar">
                <el-popover placement="bottom-start" :width="800" trigger="click">
                    <template #reference>
                        <el-button link color="#000" :icon="CollectionTag">常用语法</el-button>
                    </template>
                    <el-scrollbar height="600px">
                        <el-descriptions :column="1" border>
                            <el-descriptions-item label="声明变量">
                                <highlightjs language="yaml" :code='variables'></highlightjs>
                            </el-descriptions-item>
                            <el-descriptions-item label="常用函数">
                                <highlightjs language="yaml" :code='func'></highlightjs>
                            </el-descriptions-item>
                            <el-descriptions-item label="路径匹配">
                                <highlightjs language="yaml" :code='path'></highlightjs>
                            </el-descriptions-item>
                            <el-descriptions-item label="匹配成功立即停止">
                                <highlightjs language="yaml" code='stop-at-first-match'></highlightjs>
                            </el-descriptions-item>
                            <el-descriptions-item label="无害化上传">
                                <highlightjs language="yaml" :code='uploads'></highlightjs>
                            </el-descriptions-item>
                        </el-descriptions>
                    </el-scrollbar>
                </el-popover>
                <el-space>
                    <el-button link color="#000" :icon="DocumentCopy" @click="Copy(content)">复制</el-button>
                    <el-button link color="#000" :icon="saveIcon" @click="savePoc">保存</el-button>
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
    border-bottom: 1px solid #181818;
}

.monaco-editor {
    flex-grow: 1;
    height: calc(100% - 60px);
    /* Adjust based on toolbar height */
}
</style>
