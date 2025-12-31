<script lang="ts" setup>
import { onMounted, ref, nextTick, toRaw, computed, shallowRef } from 'vue';
import { Search, CirclePlusFilled, Delete, DocumentCopy, Reading, CollectionTag } from "@element-plus/icons-vue";
import { CheckFileStat, ReadFile, RemoveFile, SaveFileDialog, WriteFile, FilepathJoin, List } from 'wailsjs/go/services/File';
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
import { dslOptions, httpMethodOptions, metadataOptions, pocdetailFilterOptions, sortSeverityOptions, variablesOptions } from '@/stores/options';
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

    await buildPocPathMap();
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

async function buildPocPathMap() {
    const paths: string[] = []
    const defaultPath = await FilepathJoin([global.PATH.homedir, "slack", "config", "pocs"])
    paths.push(defaultPath)
    if (global.webscan.append_pocfile) {
        paths.push(global.webscan.append_pocfile)
    }
    const files = await List(paths)
    const map: Record<string, string> = {}
    files.forEach(file => {
        if (!file.Path.endsWith(".yaml")) return
        const base = file.BaseName
        const key = base.replace(/\.ya?ml$/i, "")
        map[key] = file.Path
        map[base] = file.Path
    })
    pocPathMap.value = map
}

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

const isModified = ref(false);
const originalContent = ref('')
const editorRef = shallowRef()
const isSettingContent = ref(false) // 添加标志位，区分程序设置内容和用户修改

const MONACO_EDITOR_OPTIONS = {
    automaticLayout: true,
    minimap: { enabled: false }, // 禁用迷你地图
    scrollBeyondLastLine: false, // 禁用滚动到最后一行后额外空白
    contextmenu: false, // 禁用右键菜单
}

function handleMount(editorInstance: any) {
    editorRef.value = editorInstance
    editorRef.value.addCommand(
        monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS,
        (e: DragEvent) => {
            // 阻止默认事件
            e.preventDefault?.();
            // 调用保存方法
            saveContent();
            isModified.value = false
        }
    );
    editorInstance.onDidChangeModelContent(() => {
        // 只有在非程序设置内容时才触发修改状态
        if (!isSettingContent.value) {
            const currentValue = editorInstance.getValue()
            isModified.value = currentValue !== originalContent.value
        }
    })
}

const detailDialog = ref(false)
const content = ref('')
const currentFilepath = ref(""); // 记录当前被编辑的文件路径，方便后续保存时调用
const pocPathMap = ref<Record<string, string>>({});

async function readPocFile(filename: string) {
    detailDialog.value = true
    isModified.value = false
    const lookup = pocPathMap.value[filename] || pocPathMap.value[`${filename}.yaml`]
    const defaultPath = await FilepathJoin([global.PATH.homedir, "slack", "config", "pocs", `${filename}.yaml`])
    let filepath = lookup || defaultPath
    let isStat = await CheckFileStat(filepath)
    if (!isStat && global.webscan.append_pocfile) {
        filepath = await FilepathJoin([global.webscan.append_pocfile, `${filename}.yaml`])
        isStat = await CheckFileStat(filepath)
    }
    if (!isStat && lookup) {
        filepath = lookup
    }
    currentFilepath.value = filepath
    let file = await ReadFile(filepath)
    // 使用 nextTick 确保 content 更新后再初始化编辑器
    nextTick(() => {
        if (editorRef.value) {
            content.value = file.Content
            originalContent.value = file.Content; // 保存原始内容
        }
        isSettingContent.value = false // 设置完成，取消标记
    });
}

function handleBeforeClose(done: () => void) {
    // 如果有未保存的更改，提示保存
    if (isModified.value) {
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
                isModified.value = false; // 放弃更改，标记无更改
                done(); // 放弃后关闭
            });
    } else {
        done(); // 如果没有未保存更改，直接关闭
    }
}

async function saveContent() {
    try {
        // 当前编辑器状态未改动时，直接返回不用处理保存功能
        if (!isModified.value) {
            return
        }
        // 已知BUG，需要使用 toRaw 获取原始对象，直接获取会导致CPU占用率到100%，
        const updatedContent = toRaw(editorRef.value).getValue() || content.value;
        let isSuccess = await WriteFile("yaml", currentFilepath.value, updatedContent);
        if (!isSuccess) {
            ElMessage.error("保存失败")
            return
        }
        
        ElMessage.success('保存成功!');
    } catch (error) {
        ElMessage.error('保存失败，请检查文件路径或权限!');
    }
};

// step 1 add poc, hide poclist
const step = ref(0)

const metadataTemp = ref({}); // 临时存储 key
const headerKeys = ref({}); // 临时存储请求头 key
const variablesTemp = ref({}); // 临时存储变量 key

const formRef = ref(null)

const nucleiTeamplate = ref<FormData>({
    id: 'nuclei-test-teamplate',
    name: '',
    author: 'slack',
    description: '',
    reference: '',
    severity: 'high',
    flow: '',
    rawRequest: false,
    variables: {},
    requests: [{
        method: 'GET',
        path: [],
        pathText: '',
        headers: {},
        body: '',
        rawBody: '',
        matchers: [] as Matcher[],
        matchersCondition: 'and',
        stopAtFirstMatcher: false,
        cookieReuse: false,
        extractors: []
    }],
    tags: [] as string[],
    metadata: {}
});


const rules = {
    id: [
        { required: true, message: '请输入模板id', trigger: 'blur' },
        {
            validator: (_: any, value: string, callback: (error?: string | Error) => void) => {
                const validIdPattern = /^([a-zA-Z0-9]+[-_])*[a-zA-Z0-9]+$/;
                if (!validIdPattern.test(value)) {
                    callback(new Error('模板ID只能包含字母、数字、中划线、下划线，且不能以中划线/下划线开头或结尾'));
                } else {
                    callback();
                }
            },
            trigger: 'blur'
        }
    ],
    name: [{ required: true, message: '请输入漏洞名称', trigger: 'blur' }],
    author: [{ required: true, message: '请输入作者名称', trigger: 'blur' }],
    description: [{ required: true, message: '请输入漏洞描述', trigger: 'blur' }],
    tags: [{ required: true, message: '请选择绑定的指纹, 至少一个', trigger: 'blur' }],
}

// 控制缩进函数
const indent = (level: number, content: string) => ' '.repeat(level * 2) + content;
const generatedPoc = computed(() => {

    const pocConfig = [
        `id: ${nucleiTeamplate.value.id}\n`,
        'info:',
        indent(1, `name: ${nucleiTeamplate.value.name}`),
        indent(1, `author: ${nucleiTeamplate.value.author}`),
        indent(1, `severity: ${nucleiTeamplate.value.severity}`),
        indent(1, `description: |`),
        indent(2, `${nucleiTeamplate.value.description}`),
    ];

    if (nucleiTeamplate.value.reference && Object.keys(nucleiTeamplate.value.reference).length > 0) {
        const referenceArray = nucleiTeamplate.value.reference.split('\n').filter(ref => ref.trim() !== '');
        if (referenceArray.length > 0) {
            pocConfig.push(indent(1, 'reference:'));
            referenceArray.forEach(ref => {
                pocConfig.push(indent(2, `- ${ref.trim()}`));
            });
        }
    }

    // 处理 metadata 部分
    if (nucleiTeamplate.value.metadata && Object.keys(nucleiTeamplate.value.metadata).length > 0) {
        pocConfig.push(indent(1, 'metadata:'));
        Object.entries(nucleiTeamplate.value.metadata).forEach(([key, value]) => {
            if (value) {
                pocConfig.push(indent(2, `${key}: ${value}`));
            }
        });
    }

    if (nucleiTeamplate.value.tags && Object.keys(nucleiTeamplate.value.tags).length > 0) {
        pocConfig.push(indent(1, `tags: ${nucleiTeamplate.value.tags.join(',')}`))
    }

    if (nucleiTeamplate.value.variables && Object.keys(nucleiTeamplate.value.variables).length > 0) {
        pocConfig.push(indent(0, 'variables:'));
        Object.entries(nucleiTeamplate.value.variables).forEach(([key, value]) => {
            if (value) {
                pocConfig.push(indent(1, `${key}: ${value}`));
            }
        });
    }

    if (nucleiTeamplate.value.flow && nucleiTeamplate.value.requests.length > 1) {
        pocConfig.push('')
        pocConfig.push(indent(0, 'flow: ' + nucleiTeamplate.value.flow));
    }

    pocConfig.push('\nhttp:');

    // 处理多个请求
    nucleiTeamplate.value.requests.forEach((request, requestIndex) => {
        if (requestIndex > 0) {
            pocConfig.push(''); // 请求之间的空行
        }

        if (nucleiTeamplate.value.rawRequest) {
            // Raw 请求模式
            pocConfig.push(indent(1, '- raw:'));
            pocConfig.push(indent(3, '- |'));

            if (request.rawBody && Object.keys(request.rawBody).length > 0) {
                const bodyArray = request.rawBody.split('\n');
                if (bodyArray.length > 0) {
                    bodyArray.forEach(line => {
                        pocConfig.push(indent(4, line.trim()));
                    });
                }
            }
        } else {
            // 结构化请求模式
            pocConfig.push(indent(1, '- method: ' + request.method));
            pocConfig.push(indent(2, 'path:'));
            request.path.forEach(path => {
                pocConfig.push(indent(3, '- ' + '"{{BaseURL}}' + path.trim()) + "\"");
            });

            // 处理请求头
            if (request.headers && Object.keys(request.headers).length > 0) {
                pocConfig.push(indent(2, 'headers:'));
                Object.entries(request.headers).forEach(([key, value]) => {
                    if (key && value) {
                        pocConfig.push(indent(3, `${key}: ${value}`));
                    }
                });
            }

            // 处理请求体
            if (request.body && request.body.trim()) {
                pocConfig.push(indent(2, 'body: |'));
                const bodyLines = request.body.split('\n');
                bodyLines.forEach(line => {
                    pocConfig.push(indent(3, line));
                });
            }
        }

        // 处理 stop-at-first-matcher
        if (request.stopAtFirstMatcher) {
            pocConfig.push('');
            pocConfig.push(indent(2, 'stop-at-first-matcher: true'));
        }

        if (request.cookieReuse) {
            pocConfig.push(indent(2, 'cookie-reuse: true'));
        }

        // 处理匹配规则
        if (request.matchers && request.matchers.length > 0) {

            if (request.matchers.length > 1) {
                pocConfig.push(indent(2, `matchers-condition: ${request.matchersCondition}`));
            }

            pocConfig.push(indent(2, 'matchers:'));

            const matchers = request.matchers.map(matcher => {
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

            pocConfig.push(matchers);
        }
        if (request.extractors && request.extractors.length > 0) {
            pocConfig.push(indent(2, 'extractors:'));
            request.extractors.forEach(extractor => {
                pocConfig.push(indent(3, '- type: ' + extractor.type));
                if (extractor.name) {
                    pocConfig.push(indent(4, 'name: ' + extractor.name));
                }
                if (extractor.part) {
                    pocConfig.push(indent(4, 'part: ' + extractor.part));
                }
                if (extractor.internal) {
                    pocConfig.push(indent(4, 'internal: ' + extractor.internal));
                }
                pocConfig.push(indent(4, extractor.type + ':'));
                extractor.typeValue.forEach(value => {
                    pocConfig.push(indent(5, '- ' + value));
                });
            });
        }
    });

    return pocConfig.join('\n');
});


const nuclei = ({
    fetchMetatedataSuggestions: function (query: string, cb: Function) {
        cb(metadataOptions.filter(k => k.includes(query)).map(k => ({ value: k })))
    },
    fetchDslSuggestions: function (query: string, cb: Function) {
        cb(dslOptions.filter(k => k.text.includes(query)))
    },
    fetchVariablesSuggestions: function (query: string, cb: Function) {
        cb(variablesOptions.filter(k => k.text.includes(query)).map(item => ({ value: item.text, label: item.text, realValue: item.value })));
    },
    updateMetadataKey: function (oldKey: string, newKey: string) {
        if (newKey && oldKey !== newKey && !nucleiTeamplate.value.metadata[newKey]) {
            nucleiTeamplate.value.metadata[newKey] = nucleiTeamplate.value.metadata[oldKey];
            delete nucleiTeamplate.value.metadata[oldKey];
            metadataTemp.value[newKey] = newKey;
            delete metadataTemp.value[oldKey];
        }
    },
    addMetadata: function () {
        if (!nucleiTeamplate.value.metadata) {
            nucleiTeamplate.value.metadata = {};
        }
        const newKey = ''; // 默认空 key
        nucleiTeamplate.value.metadata[newKey] = '';
        metadataTemp.value[newKey] = newKey;
    },
    removeMetadata: function (key: string) {
        if (nucleiTeamplate.value.metadata) {
            delete nucleiTeamplate.value.metadata[key];
            delete metadataTemp.value[key]; // 同步删除临时 key
        }
    },
    updateMatcherWords: function (requestIndex: number, matcherIndex: number) {
        nextTick(() => {
            const matcher = nucleiTeamplate.value.requests[requestIndex].matchers[matcherIndex];
            if (matcher.wordsText) {
                matcher.words = matcher.wordsText.split('\n').filter(word => word.trim());
            } else {
                matcher.words = [];
            }
        })
    },
    addMatcher: function (requestIndex: number) {
        nucleiTeamplate.value.requests[requestIndex].matchers.push({
            type: 'word',
            part: 'all',
            words: [],
            condition: 'and',
            wordsText: '',
        });
    },
    removeMatcher: function (requestIndex: number, matcherIndex: number) {
        nucleiTeamplate.value.requests[requestIndex].matchers.splice(matcherIndex, 1);
    },
    replaceHostAll: function () {
        nucleiTeamplate.value.requests[0].rawBody = nucleiTeamplate.value.requests[0].rawBody.replace(/(Host:\s*)([^\n]+)/g, '{{Hostname}}');
    },
    addHeader: function () {
        const newKey = ''; // 默认空 key
        nucleiTeamplate.value.requests[0].headers[newKey] = '';
        headerKeys.value[newKey] = newKey;
    },
    removeHeader: function (key: string) {
        if (nucleiTeamplate.value.requests[0].headers) {
            delete nucleiTeamplate.value.requests[0].headers[key];
            delete headerKeys.value[key]; // 同步删除临时 key
        }
    },
    updateHeaderKey: function (oldKey: string, newKey: string) {
        if (newKey && oldKey !== newKey && !nucleiTeamplate.value.requests[0].headers[newKey]) {
            nucleiTeamplate.value.requests[0].headers[newKey] = nucleiTeamplate.value.requests[0].headers[oldKey];
            delete nucleiTeamplate.value.requests[0].headers[oldKey];
            headerKeys.value[newKey] = newKey;
            delete headerKeys.value[oldKey];
        }
    },
    addVariable: function () {
        const newKey = '';
        nucleiTeamplate.value.variables[newKey] = '';
        variablesTemp.value[newKey] = newKey;
    },
    updateVariablesKey: function (oldKey: string, newKey: string) {
        if (newKey && oldKey !== newKey && !nucleiTeamplate.value.variables[newKey]) {
            nucleiTeamplate.value.variables[newKey] = nucleiTeamplate.value.variables[oldKey];
            delete nucleiTeamplate.value.variables[oldKey];
            variablesTemp.value[newKey] = newKey;
            delete variablesTemp.value[oldKey];
        }
    },
    removeVariable: function (key: string) {
        if (nucleiTeamplate.value.variables) {
            delete nucleiTeamplate.value.variables[key];
            delete variablesTemp.value[key];
        }
    },
    addRequest: function () {
        nucleiTeamplate.value.requests.push({
            method: 'GET',
            path: [],
            pathText: '',
            headers: {},
            body: '',
            rawBody: '',
            matchers: [],
            matchersCondition: 'and',
            stopAtFirstMatcher: false,
            cookieReuse: false,
            extractors: [],
        });
    },
    removeRequest: function (index: number) {
        if (nucleiTeamplate.value.requests.length > 1) {
            nucleiTeamplate.value.requests.splice(index, 1);
        }
    },
    addExtractor: function (requestIndex: number) {
        nucleiTeamplate.value.requests[requestIndex].extractors.push({
            type: 'regex',
            name: '',
            part: 'body',
            typeValue: [],
            internal: false,
        });
    },
    updateExtractor: function (requestIndex: number, matcherIndex: number) {
        nextTick(() => {
            const extractor = nucleiTeamplate.value.requests[requestIndex].extractors[matcherIndex];
            if (extractor.typeText) {
                extractor.typeValue = extractor.typeText.split('\n').filter(word => word.trim());
            } else {
                extractor.typeValue = [];
            }
        })
    },
    removeExtractor: function (requestIndex: number, extractorIndex: number) {
        nucleiTeamplate.value.requests[requestIndex].extractors.splice(extractorIndex, 1);
    },
    updatePaths: function (index: number) {
        const request = nucleiTeamplate.value.requests[index];
        request.path = request.pathText.split('\n').map(path => path.trim());
    },
})

async function savePoc() {
    if (nucleiTeamplate.value.id === "") {
        ElMessage.warning("请输入POC ID!")
        return;
    }
    const path = await SaveFileDialog(nucleiTeamplate.value.id + ".yaml");
    if (!path) {
        return;
    }
    const result = await WriteFile("yaml", path, generatedPoc.value);
    result ? ElMessage.success("保存成功!") : ElMessage.error("保存失败!");
}

const func = `随机整数: {{rand_int(40000, 44800)}}
随机字符串: {{rand_base(8)}} || {{randstr}}
转小写: {{to_lower(rand_text_alpha(5))}}
转大写: {{to_upper(rand_base(12))}}`

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
                        <el-button class="ml-5px" type="primary" plain :icon="CirclePlusFilled"
                            @click="step = 1">添加POC</el-button>
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
            </template>
            <template #extra>
                <el-button @click="BrowserOpenURL('https://cloud.projectdiscovery.io/templates')">
                    <template #icon>
                        <el-icon size="18">
                            <aiIcon />
                        </el-icon>
                    </template>
                    AI生成
                </el-button>
            </template>
        </el-page-header>
        <el-divider />
        <div class="flex gap-2">
            <el-form :model="nucleiTeamplate" :rules="rules" ref="formRef" label-width="auto" class="w-1/2">
                <el-form-item label="模板id" prop="id">
                    <el-input v-model="nucleiTeamplate.id"></el-input>
                </el-form-item>
                <el-form-item label="漏洞名称" prop="name">
                    <el-input v-model="nucleiTeamplate.name"></el-input>
                </el-form-item>
                <el-form-item label="作者" prop="author">
                    <el-input v-model="nucleiTeamplate.author"></el-input>
                </el-form-item>
                <el-form-item label="风险等级" prop="severity">
                    <el-select v-model="nucleiTeamplate.severity">
                        <el-option v-for="severity in sortSeverityOptions" :key="severity" :label="severity"
                            :value="severity.toLowerCase()" />
                    </el-select>
                </el-form-item>
                <el-form-item label="漏洞描述" prop="description">
                    <el-input v-model="nucleiTeamplate.description" type="textarea" :rows="3"></el-input>
                </el-form-item>
                <el-form-item label="漏洞来源" prop="reference">
                    <el-input v-model="nucleiTeamplate.reference" placeholder="请输入来源, 每行一个" type="textarea"
                        :rows="3"></el-input>
                </el-form-item>
                <el-form-item label="Matedata">
                    <div v-if="nucleiTeamplate.metadata" class="w-full">
                        <div v-for="(_, key) in nucleiTeamplate.metadata" :key="key"
                            style="display: flex; align-items: center; margin-bottom: 8px; width: 100%;">
                            <!-- 这里用 el-autocomplete 提供建议 -->
                            <el-autocomplete v-model="metadataTemp[key]"
                                @select="nuclei.updateMetadataKey(key, $event.value)"
                                :fetch-suggestions="nuclei.fetchMetatedataSuggestions" placeholder="请输入 Key"
                                style="flex: 1; margin-right: 8px;"
                                @change="nuclei.updateMetadataKey(key, metadataTemp[key])" />
                            <el-input v-model="nucleiTeamplate.metadata[key]" placeholder="请输入 Value"
                                style="flex: 1; margin-right: 8px;" />
                            <el-button :icon="Delete" type="danger" plain @click="nuclei.removeMetadata(key)" />
                        </div>
                    </div>
                    <el-button type="primary" size="small" @click="nuclei.addMetadata">添加 Matedata</el-button>
                </el-form-item>
                <el-form-item label="关联指纹" prop="tags">
                    <el-select-v2 v-model="nucleiTeamplate.tags" filterable :options="fingerOptions"
                     multiple clearable />
                </el-form-item>
                <el-form-item label="全局变量">
                    <div v-if="nucleiTeamplate.variables" class="w-full">
                        <div v-for="(_, key) in nucleiTeamplate.variables" :key="key"
                            style="display: flex; align-items: center; margin-bottom: 8px; width: 100%;">
                            <el-autocomplete v-model="variablesTemp[key]"
                                :fetch-suggestions="nuclei.fetchVariablesSuggestions" placeholder="请输入 Key"
                                style="flex: 1; margin-right: 8px;"
                                @select="(item) => { variablesTemp[key] = item.value; nuclei.updateVariablesKey(key, item.value); nucleiTeamplate.variables[item.value] = item.realValue; delete nucleiTeamplate.variables[key]; }"
                                @change="nuclei.updateVariablesKey(key, variablesTemp[key])" />
                            <el-input v-model="nucleiTeamplate.variables[key]" placeholder="请输入 Value"
                                style="flex: 1; margin-right: 8px;" />
                            <el-button :icon="Delete" type="danger" plain @click="nuclei.removeVariable(key)" />
                        </div>
                    </div>
                    <el-button type="primary" size="small" @click="nuclei.addVariable">添加 Variables</el-button>
                </el-form-item>
                <el-form-item label="Raw请求">
                    <el-switch v-model="nucleiTeamplate.rawRequest" />
                </el-form-item>
                <el-form-item label="工作流" v-show="nucleiTeamplate.requests.length > 1">
                    <el-input v-model="nucleiTeamplate.flow"></el-input>
                </el-form-item>
                <!-- 请求组 -->
                <div v-for="(request, requestIndex) in nucleiTeamplate.requests" :key="requestIndex"
                    class="request-group">
                    <el-card class="mb-10px">
                        <template #header>
                            <div class="card-header">
                                <span>请求 #{{ requestIndex + 1 }}</span>
                                <el-button v-if="nucleiTeamplate.requests.length > 1" size="small" type="danger" plain
                                    :icon="Delete" @click="nuclei.removeRequest(requestIndex)">
                                </el-button>
                            </div>
                        </template>

                        <el-form-item label="Raw数据包" v-show="nucleiTeamplate.rawRequest">
                            <el-input v-model="request.rawBody" type="textarea" :rows="8" placeholder="请输入请求体内容" />
                        </el-form-item>

                        <div v-show="!nucleiTeamplate.rawRequest">
                            <el-form-item label="请求方式">
                                <el-select-v2 v-model="request.method" :options="httpMethodOptions">
                                </el-select-v2>
                            </el-form-item>
                            <el-form-item label="请求路径">
                                <el-input v-model="request.pathText" type="textarea" :rows="3" placeholder="请输入路径，每行一个"
                                    @input="nuclei.updatePaths(requestIndex)" />
                            </el-form-item>
                            <el-form-item label="请求头">
                                <div v-for="(value, key) in request.headers" :key="key"
                                    style="display: flex; align-items: center; margin-bottom: 8px; width: 100%;">
                                    <el-input v-model="headerKeys[key]" placeholder="请输入 Key"
                                        style="flex: 1; margin-right: 8px;"
                                        @change="nuclei.updateHeaderKey(key, headerKeys[key])" />
                                    <el-input v-model="request.headers[key]" placeholder="请输入 Value"
                                        style="flex: 1; margin-right: 8px;" />
                                    <el-button :icon="Delete" type="danger" plain @click="nuclei.removeHeader(key)" />
                                </div>
                                <el-button type="primary" size="small" @click="nuclei.addHeader">添加请求头</el-button>
                            </el-form-item>
                            <el-form-item label="请求体">
                                <el-input v-model="request.body" type="textarea" :rows="5"></el-input>
                            </el-form-item>
                        </div>
                        <el-form-item label="首次匹配停止">
                            <el-switch v-model="request.stopAtFirstMatcher" />
                        </el-form-item>
                        <el-form-item label="Cookie复用">
                            <el-switch v-model="request.cookieReuse" />
                        </el-form-item>
                        <!-- 匹配规则 -->
                        <el-form-item label="匹配规则">
                            <el-select v-model="request.matchersCondition" class="w-full">
                                <el-option label="AND" value="and" />
                                <el-option label="OR" value="or" />
                            </el-select>
                            <div class="w-full mt-5px">
                                <el-card v-for="(matcher, matcherIndex) in request.matchers" :key="matcherIndex"
                                    class="mb-5px">
                                    <template #header>
                                        <div class="card-header">
                                            <span>规则 #{{ matcherIndex + 1 }}</span>
                                            <el-radio-group v-model="matcher.condition">
                                                <el-radio value="and">AND</el-radio>
                                                <el-radio value="or">OR</el-radio>
                                            </el-radio-group>
                                            <el-button size="small" type="danger" plain :icon="Delete"
                                                @click="nuclei.removeMatcher(requestIndex, matcherIndex)">
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
                                        <el-form-item label="匹配部分"
                                            v-if="matcher.type === 'word' || matcher.type === 'regex'">
                                            <el-select v-model="matcher.part">
                                                <el-option label="响应体" value="body" />
                                                <el-option label="响应头" value="header" />
                                                <el-option label="全部" value="all" />
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="词条">
                                            <el-autocomplete v-model="matcher.wordsText" type="textarea" :rows="3"
                                                :fetch-suggestions="matcher.type === 'dsl' ? nuclei.fetchDslSuggestions : () => []"
                                                placeholder="请输入词条, 每行一个"
                                                @input="nuclei.updateMatcherWords(requestIndex, matcherIndex)" />
                                        </el-form-item>
                                    </el-form>
                                </el-card>
                                <el-button type="primary" size="small" @click="nuclei.addMatcher(requestIndex)">
                                    添加匹配规则
                                </el-button>
                            </div>
                        </el-form-item>
                        <el-form-item label="提取器">
                            <div class="w-full">
                                <el-card v-for="(extractor, extractorIndex) in request.extractors" :key="extractorIndex"
                                    class="mb-5px">
                                    <template #header>
                                        <div class="card-header">
                                            <span>提取器 #{{ extractorIndex + 1 }}</span>
                                            <el-button size="small" type="danger" plain :icon="Delete"
                                                @click="nuclei.removeExtractor(requestIndex, extractorIndex)">
                                            </el-button>
                                        </div>
                                    </template>

                                    <el-form label-position="top">
                                        <el-form-item label="类型">
                                            <el-select v-model="extractor.type">
                                                <el-option label="regex" value="regex" />
                                                <el-option label="xpath" value="xpath" />
                                                <el-option label="json" value="json" />
                                                <el-option label="kval" value="kval" />
                                                <el-option label="dsl" value="dsl" />
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="名称">
                                            <el-input v-model="extractor.name" placeholder="提取器名称（可选）" />
                                        </el-form-item>
                                        <el-form-item label="提取部分">
                                            <el-select v-model="extractor.part">
                                                <el-option label="响应体" value="body" />
                                                <el-option label="响应头" value="header" />
                                                <el-option label="全部" value="all" />
                                            </el-select>
                                        </el-form-item>
                                        <el-form-item label="内部提取">
                                            <el-switch v-model="extractor.internal" />
                                        </el-form-item>
                                        <el-form-item label="表达式">
                                            <el-input v-model="extractor.typeText" type="textarea" :rows="3"
                                                placeholder="请输入表达式, 每行一个"
                                                @input="nuclei.updateExtractor(requestIndex, extractorIndex)" />
                                        </el-form-item>
                                    </el-form>
                                </el-card>
                                <el-button type="primary" size="small" @click="nuclei.addExtractor(requestIndex)">
                                    添加提取器
                                </el-button>
                            </div>
                        </el-form-item>
                    </el-card>
                </div>
                <el-button type="primary" size="small" @click="nuclei.addRequest" class="mb-10px">
                    添加请求
                </el-button>
            </el-form>
            <el-card class="w-1/2">
                <div class="card-header">
                    <span>POC预览</span>
                    <el-space>
                        <el-popover placement="left-start" :width="800" trigger="click">
                        <template #reference>
                            <el-button size="small" :icon="CollectionTag">常用语法</el-button>
                        </template>
                        <el-descriptions :column="1" border>
                            <el-descriptions-item label="常用函数">
                                <highlightjs language="yaml" :code='func'></highlightjs>
                            </el-descriptions-item>
                            <el-descriptions-item label="无害化上传">
                                <highlightjs language="yaml" :code='uploads'></highlightjs>
                            </el-descriptions-item>
                        </el-descriptions>
                    </el-popover>
                        <el-button size="small" :icon="DocumentCopy" @click="Copy(content)">复制</el-button>
                        <el-button size="small" :icon="saveIcon" @click="savePoc">保存</el-button>
                    </el-space>
                </div>
                <highlightjs language="yaml" :code='generatedPoc'></highlightjs>
            </el-card>
        </div>
    </div>
    <el-drawer v-model="detailDialog" size="70%" :before-close="handleBeforeClose">
        <template #header>
            <span class="drawer-title">漏洞详情</span>
        </template>
        <div class="editor-container">
            <!-- 操作区 -->
            <div class="editor-toolbar">
                <div>
                    <el-tag type="warning" v-show="isModified">未保存</el-tag>
                </div>
                <el-space>
                    <el-button plain :icon="DocumentCopy" @click="Copy(content)">复制</el-button>
                    <el-button plain :icon="saveIcon" @click="saveContent">保存</el-button>
                </el-space>
            </div>
            <div style="height: calc(100vh - 130px);">
                <vue-monaco-editor v-model:value="content" language="yaml" theme="vs-dark"
                    :options="MONACO_EDITOR_OPTIONS" @mount="handleMount" />
            </div>
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
