<template>
    <!-- 多分页区域 -->
    <el-tabs v-model="activeTab" type="card" editable @edit="handleTabsEdit">
        <el-tab-pane v-for="tab in tabs" :key="tab.name" :label="tab.label" :name="tab.name">
            <el-card class="h-full">
                <template #header>
                    <el-space>
                        <el-button type="primary" :icon="Promotion" @click="sendRequest(false)">发送请求</el-button>
                        <el-checkbox v-model="request.forceHttps">强制HTTPS</el-checkbox>
                    </el-space>
                </template>
                <el-splitter>
                    <el-splitter-panel size="50%">
                        <div class="h-full">
                            <div class="flex-between mr-5px ml-5px mb-5px">
                                <span class="font-bold">Request</span>
                                <div>
                                    <el-button size="small" v-show="tab.hasRedirect" @click="sendRequest(true)">Follow Redirect</el-button>
                                    <el-button-group>
                                        <!-- <el-button :icon="form.hideRequest ? DArrowRight : DArrowLeft" link
                                            @click="form.hideRequest = !form.hideRequest" /> -->
                                        <el-button :icon="DocumentCopy" link @click="Copy('')" />
                                    </el-button-group>
                                </div>
                            </div>
                            <vue-monaco-editor v-model:value="tab.requestCode" theme="httpTheme" language="http"
                                :options="Request_OPTIONS" @mount="handleRequestMount"
                                style="height: calc(100vh - 255px);" />
                        </div>
                    </el-splitter-panel>
                    <el-splitter-panel>
                        <div class="h-full">
                            <div class="flex-between mr-5px ml-5px mb-5px">
                                <span class="font-bold">
                                    Response
                                    <el-tag v-show="tab.responseTime > 0" size="small">响应时间：{{ tab.responseTime }} ms</el-tag>
                                </span>
                                <el-button-group>
                                    <!-- <el-button :icon="form.hideRequest ? DArrowRight : DArrowLeft" link
                                        @click="form.hideRequest = !form.hideRequest" /> -->
                                    <el-button :icon="DocumentCopy" link @click="Copy('')" />
                                </el-button-group>
                            </div>
                            <vue-monaco-editor 
                                v-loading="tab.loading"
                                element-loading-text="send request."
                                v-model:value="tab.responseCode" 
                                theme="httpTheme"
                                language="http"
                                :options="Response_OPTIONS"
                                @mount="handleResponseMount"
                                style="height: calc(100vh - 255px);"
                             />
                        </div>
                    </el-splitter-panel>
                </el-splitter>
            </el-card>
        </el-tab-pane>
        <!-- <template #extra>
      <el-button type="primary" icon="Plus" @click="addTab">新增分页</el-button>
    </template> -->
    </el-tabs>
</template>

<script lang="ts" setup>
import { reactive, ref, shallowRef } from 'vue'
import type { TabPaneName } from 'element-plus'
import { DocumentCopy, Promotion } from '@element-plus/icons-vue';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api.js';
import { Copy } from '@/util';
import { SendRequest } from 'wailsjs/go/services/App';
// 注册 http 语言和高亮规则
monaco.languages.register({ id: 'http' })

monaco.languages.setMonarchTokensProvider('http', {
    brackets: [],
    defaultToken: "",
    ignoreCase: true,
    includeLF: true,
    start: "",
    tokenPostfix: "",
    unicode: false,
    escapes: /\\(?:[abfnrtv\\"']|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,
    tokenizer: {
        root: [
            // HTTP请求方法
            // 基础 Fuzz 标签解析
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/(GET|POST|OPTIONS|DELETE|PUT|PATCH)/g, "http.method"],
            [/\s/, "delimiter", "@http_path"],
            // [/(html|div|src|\<\/?title\>|<alert>)/i, "keyword"],
            // [/(\<script\>|<alert>|<prompt>|<svg )/i, "keyword"],
            // [/(secret)|(access)|(password)|(verify)|(login)/i, "bold-keyword"],
        ],
        fuzz_tag: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag_second"],
            [/}}/, "fuzz.tag.inner", "@pop"],
            [/[\w:]+}}/, "fuzz.tag.inner", "@pop"],
            [/[\w:]+\(/, "fuzz.tag.inner", "@fuzz_tag_param"],
        ],
        fuzz_tag_second: [
            [/{{/, "fuzz.tag.second", "@fuzz_tag"],
            [/}}/, "fuzz.tag.second", "@pop"],
            [/[\w:]+}}/, "fuzz.tag.second", "@pop"],
            [/[\w:]+\(/, "fuzz.tag.second", "@fuzz_tag_param_second"],
        ],
        fuzz_tag_param: [
            [/\(/, "fuzz.tag.inner", "@fuzz_tag_param"],
            [/\\\)/, "bold-keyword"],
            [/\)/, "fuzz.tag.inner", "@pop"],
            [/{{/, "fuzz.tag.second", "@fuzz_tag_second"],
            [/./, "bold-keyword"]
        ],
        fuzz_tag_param_second: [
            [/\\\)/, "bold-keyword"],
            [/\)/, "fuzz.tag.second", "@pop"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/./, "bold-keyword"]
        ],
        http_path: [
            [/\s/, "delimiter", "@http_protocol"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            ["/(((http)|(https):)?\/\/[^\s]+?)/", "http.url"],
            [/#/, "http.anchor", "@http_anchor"],
            // [/\/[^\s^?^\/]+/, "http.path"],
            [/\?/, "http.query", "@query"],
        ],
        http_anchor: [
            [/\s/, "delimiter", "@http_protocol"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/./, "http.anchor"],
        ],
        http_protocol: [
            [/\n/, "delimiter", "@http_header"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/HTTP\/[0-9.]+/, "http.protocol"],
        ],
        http_header: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/^\n$/, 'body.delimiter', '@body'],
            [/(Cookie)(:)/g, ["http.header.info", { token: "delimiter", next: "@http_cookie" }]],
            [/(Content-Type)(:)/g, ["http.header.info", { token: "delimiter" }]],
            [/(Content-Length|Host|Origin|Referer)(:)/g, ["http.header.info", { token: "delimiter", next: "@http_header_value" }]],
            [/(Authorization|X-Forward|Real|User-Agent|Protection|CSP|X-Requested-With)(:)/g, ["http.header.info", { token: "delimiter", next: "@http_header_value" }]],
            [/Sec/, "http.header.info", "@sec_http_header"],
            [/:/, "delimiter", "@http_header_value"],
            [/\S/, "http.header.info"],
        ],
        sec_http_header: [
            [/\n/, 'body.delimiter', '@pop'],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/:/, "delimiter", "@http_header_value"],
            [/\S/, "http.header.info"],
        ],
        query: [
            [/\s/, "delimiter", "@http_protocol"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/#/, "http.anchor", "@http_anchor"],
            [/[^=&?\[\]\s]/, "http.query.params", "@http_query_params"],
            [/%[0-9ABCDEFabcdef]{2}/, "http.urlencoded"],
        ],
        http_query_params: [
            [/\s/, { token: "delimiter", next: "@pop", goBack: 1 }],
            [/#/, "http.anchor", "@http_anchor"],
            [/&/, 'delimiter', "@pop"],
            [/(\[)(\w+)(\])/, ["http.query.index", "http.query.index.values", "http.query.index"]],
            [/\=/, "http.query.equal", "http_query_params_values"],
            [/./, "http.query.params"],
        ],
        http_query_params_values: [
            [/\s/, { token: "delimiter", next: "@pop", goBack: 1 }],
            [/#/, "http.anchor", "@http_anchor"],
            [/&/, { token: 'delimiter', next: "@pop", goBack: 1 }],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/[^=&?\s]/, "http.query.values"],
            [/%[0-9ABCDEFabcdef]{2}/, "http.urlencoded"],
        ],
        http_cookie: [
            [/\n/, "delimiter", "@pop"],
            [/\s/, "delimiter"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/[^=;\s]+/, "http.query.params", "@http_cookie_params"],
            [/%[0-9ABCDEFabcdef]{2}/, "http.urlencoded"],
        ],
        http_cookie_params: [
            [/\n/, { token: "delimiter", next: "@pop", goBack: 1 }],
            [/[\s|;]/, "delimiter", "@pop"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/\=/, "http.query.equal"],
            [/[^=;?\s]/, "http.query.values"],
            [/%[0-9ABCDEFabcdef]{2}/, "http.urlencoded"],
        ],
        http_header_value: [
            [/\n/, { token: "delimiter", next: "@pop", goBack: 1 }],
            [/\s/, "delimiter"],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
        ],
        string_double: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/[^\\"]/, "json.key"],
            [/@escapes/, "string.escape"],
            [/\\./, "string.escape.invalid"],
            [/"/, "json.key", "@pop"],
        ],
        string_double_value: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/[^\\"]/, "json.value"],
            [/@escapes/, "string.escape"],
            [/\\./, "string.escape.invalid"],
            [/"/, "json.value", "@pop"],
        ],
        body: [
            [/(\{)("|\d|\n)/, ["json.start", { token: "json.key", next: "@body_json", goBack: 1 }]],
            // [/(\d+)(:)/, [{ token: "number", next: "@body_json" }, "delimiter"]],
            // [/(\d+\.\d*)(:)/, [{ token: "number", next: "@body_json" }, "delimiter"]],
            // [/"/, 'string', '@string_double'],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/-{2,}.*/, "body.boundary", "@body_form"],
            [/%[0-9ABCDEFabcdef]{2}/, "http.urlencoded"],
            [/./, "http.query.params", "@http_query_params"],
        ],
        body_json: [
            ["\}", { token: "json.end", next: "@pop" }],
            [/(\{)("|\d|\n)/, ["json.start.value", { token: "json.key", next: "@body_json", goBack: 1 }]],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/(:\s*)/, { token: "delimiter", next: "@body_json_value" }],
            [/(\d+\.\d*)/, "json.key"],
            [/(\d+)/, "json.key"],
            [/"/, 'json.key', '@string_double'],
            [/true|false|null/, "json.key"]
        ],
        body_json_value: [
            ["\}|\,", { token: "json.end.value", next: "@pop", goBack: 1 }],
            [/(\{)("|\d|\n)/, ["json.start.value", { token: "json.key", next: "@body_json", goBack: 1, }]],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/(\d+\.\d*)/, "json.value"],
            [/(\d+)/, "json.value"],
            [/"/, 'json.value', '@string_double_value'],
            [/\[/, "json.array.start", "@body_json_array_value"],
            [/true|false|null/, "json.value"],
        ],
        body_json_array_value: [
            [/\]/, "json.array.end", "@pop"],
            [/(\{)("|\d|\n)/, ["json.start.value", { token: "json.key", next: "@body_json", goBack: 1 }]],
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/(\d+\.\d*)/, "json.value"],
            [/(\d+)/, "json.value"],
            [/"/, 'json.value', '@string_double_value'],
        ],
        body_form: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/^\n$/, "body.delimiter", "@body_data"],
            [/([^:]*?)(:)/g, ["http.header.info", { token: "delimiter", next: "@http_header_value" }]],
        ],
        body_data: [
            [/{{/, "fuzz.tag.inner", "@fuzz_tag"],
            [/(-{2,}[a-zA-z0-9]+--)/, [{ token: "body.boundary.end", next: "@end" }]],
            [/(-{2,}[a-zA-z0-9]+)/, [{ token: "body.boundary", next: "@pop" }]],
        ],
        end: [],
    }
})


// 基于 Yakit 的设计风格定义增强主题
monaco.editor.defineTheme('httpTheme', {
    base: 'vs-dark',
    inherit: true,
    rules: [
        { token: 'http.url', foreground: '61AFEF', fontStyle: 'underline' },
        { token: 'http.anchor', foreground: 'E06C75' },
        { token: 'http.protocol', foreground: '56B6C2', fontStyle: 'italic' },

        { token: 'http.header.info', foreground: 'D5E7F7' },

        { token: 'http.header.mime.form', foreground: '98C379' },
        { token: 'http.header.mime.json', foreground: '61AFEF' },
        { token: 'http.header.mime.xml', foreground: 'C678DD' },
        { token: 'http.header.mime.urlencoded', foreground: 'E5C07B' },
        { token: 'http.header.mime.default', foreground: 'ABB2BF' },

        { token: 'http.query', foreground: 'C678DD' },
        { token: 'http.query.params', foreground: 'E5C07B' },
        { token: 'http.query.index', foreground: 'D19A66' },
        { token: 'http.query.index.values', foreground: '61AFEF' },
        { token: 'http.query.equal', foreground: '56B6C2' },
        { token: 'http.query.values', foreground: '98C379' },

        { token: 'json.key', foreground: 'E06C75' },
        { token: 'json.value', foreground: '98C379' },
        { token: 'json.start', foreground: 'D19A66' },
        { token: 'json.end', foreground: 'D19A66' },
        { token: 'json.start.value', foreground: 'C678DD' },
        { token: 'json.end.value', foreground: 'C678DD' },
        { token: 'json.array.start', foreground: 'D19A66' },
        { token: 'json.array.end', foreground: 'D19A66' },

        { token: 'body.delimiter', foreground: '5C6370' },
        { token: 'body.boundary', foreground: '61AFEF' },
        { token: 'body.boundary.end', foreground: 'E06C75' },


        { token: 'string.escape', foreground: '56B6C2' },
        { token: 'string.escape.invalid', foreground: 'BE5046' },

        { token: 'delimiter', foreground: 'ABB2BF' }
    ],
    colors: {

    }
})

// 分页逻辑
const tabs = ref([
  { name: 'tab-1', label: '1', requestCode: '', responseCode: '', responseTime: 0, hasRedirect: false, redirectReq: false, loading: false }
])

const activeTab = ref('tab-1')
let tabIndex = 2
function handleTabsEdit(targetName: TabPaneName | undefined, action: 'remove' | 'add') {
    if (action === 'add') {
        const newName = `tab-${tabIndex++}`
        tabs.value.push({ name: newName, label: `${tabIndex - 1}`, requestCode: '', responseCode: '', responseTime: 0, hasRedirect: false, redirectReq: false, loading: false })
        activeTab.value = newName
    } else if (action === 'remove') {
        const index = tabs.value.findIndex(t => t.name === targetName)
        if (index !== -1) {
            tabs.value.splice(index, 1)
            if (activeTab.value === targetName && tabs.value.length > 0) {
                activeTab.value = tabs.value[Math.max(0, index - 1)].name
            }
        }
    }
}

const Request_OPTIONS = {
    automaticLayout: true,
    formatOnType: true,
    formatOnPaste: true,
    readOnly: false,
    minimap: { enabled: false }, // 禁用迷你地图
    scrollBeyondLastLine: false, // 禁用滚动到最后一行后额外空白
    contextmenu: false, // 禁用右键菜单
    wordWrap: 'on', // ✅ 自动换行
}

const Response_OPTIONS = {
    automaticLayout: true,
    readOnly: true,
    minimap: { enabled: false },
    scrollBeyondLastLine: false,
    contextmenu: false,
    wordWrap: 'on'
}

const requestEditor = shallowRef()
const responseEditor = shallowRef()

const handleRequestMount = (instance: any) => (requestEditor.value = instance)
const handleResponseMount = (instance: any) => (responseEditor.value = instance)


async function sendRequest(needRedirect: boolean) {
    const currentTab = tabs.value.find(tab => tab.name === activeTab.value)
    if (!currentTab) return

    currentTab.loading = true
    currentTab.responseCode = ""
    if (needRedirect) {
        currentTab.redirectReq = true
        currentTab.hasRedirect = false
    }
    try {
        let rawResponse = await SendRequest(currentTab.requestCode, request.forceHttps, currentTab.redirectReq, request.proxyURL)
        if (rawResponse.Error != "") {
            currentTab.responseCode = rawResponse.Error
        } else {
            if (rawResponse.StatusCode >= 300 && rawResponse.StatusCode <= 399) {
                currentTab.hasRedirect = true
            }
            currentTab.responseCode = rawResponse.Response
            currentTab.responseTime = rawResponse.ResponseTime
        }
    } catch (err: any) {
        currentTab.responseCode = err.message
    } finally {
        currentTab.redirectReq = false
        currentTab.loading = false
    }
}

const request = reactive({
    forceHttps: false,
    proxyURL: ''
})
</script>

<style scoped></style>