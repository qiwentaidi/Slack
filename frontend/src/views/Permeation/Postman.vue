<script setup lang="ts">
import { GoFetch } from "../../../wailsjs/go/main/App";
import { reactive, nextTick, ref } from 'vue'
import { ElTable } from 'element-plus'
import { Plus, Minus } from '@element-plus/icons-vue';
import Prism from "prismjs";
import global from "../../global"
import { Splitpanes, Pane } from 'splitpanes'

const form = reactive({
    url: '',
    requestOptions: ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"],
    reqOption: "GET",
    repsonse: '',
    respHeader: [{}],
    reqParams: [{}],
    reqHeader: [
        {
            key: "User-Agent",
            value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
        },
    ],
})


function NewFetch() {
    GoFetch(form.reqOption, form.url, "", form.reqHeader, 10, global.proxy).then(result => {
        if (!result.Error) {
            form.repsonse = result.Body
            nextTick().then(() => {
                Prism.highlightAll(); //修改内容后重新渲染
            });
            form.respHeader = result.Header
        }
    })
}

interface reqsHeader {
    key: string,
    value: string,
}
const multipleSelection = ref<reqsHeader[]>([])
const handleSelectionChange = (val: reqsHeader[]) => {
    multipleSelection.value = val
}

function GetFocus(row: any, column: any, cell: any, event: any) {
    if (column.property === 'key' || column.property === 'value') {
        row.isEditCell = true;
    }
}

function LostFocus(row: any, column: any) {
    row.isEditCell = false
}

function AddItem(content: {}[]) {
    content.push({
        key: "",
        value: ""
    })
}

const deleteRow = (content: {}[], index: number) => {
    content.splice(index, 1)
}

</script>

<template>
    <el-form>
        <el-form-item>
            <div class="head">
                <el-select v-model=form.reqOption size="large" value=options>
                    <el-option v-for="item in form.requestOptions" :value="item" :label="item" />
                </el-select>
                <el-input v-model="form.url" placeholder="请输入URL" style="margin-right: 10px;" />
                <el-button type="primary" size="large" @click="NewFetch">Send</el-button>
            </div>
        </el-form-item>
    </el-form>
    <splitpanes horizontal>
        <pane class="default-theme" min-size="0" max-size="100">
            <!-- table 1 -->
            <div style="height: 200px; width: 100%;">
                <el-tabs>
                    <el-tab-pane label="Params">
                        <el-table :data="form.reqParams" border>
                            <el-table-column type="selection" width="55" />
                            <el-table-column prop="key" label="key" />
                            <el-table-column prop="value" label="value" />
                        </el-table>
                    </el-tab-pane>
                    <el-tab-pane label="Header">
                        <el-table :data="form.reqHeader" border @cell-click="GetFocus"
                            @selection-change="handleSelectionChange">
                            <el-table-column type="selection" width="55" />
                            <el-table-column prop="key" label="key" width="200">
                                <template #default="scope">
                                    <span v-if="!scope.row.isEditCell">{{ scope.row.key }}</span>
                                    <el-input v-else v-model="scope.row.key"
                                        @blur="LostFocus(scope.row, scope.row.key)"></el-input>
                                </template>
                            </el-table-column>
                            <el-table-column prop="value" label="value" show-overflow-tooltip>
                                <template #default="scope">
                                    <span v-if="!scope.row.isEditCell">{{ scope.row.value }}</span>
                                    <el-input v-else v-model="scope.row.value"
                                        @blur="LostFocus(scope.row, scope.row.value)"></el-input>
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="100px" align="center">
                                <template #default="scope">
                                    <el-button link :icon="Plus" @click="AddItem(form.reqHeader)"></el-button>
                                    <el-button link :icon="Minus"
                                        @click.prevent="deleteRow(form.reqHeader, scope.$index)"></el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </el-tab-pane>
                    <el-tab-pane label="Body">
                        <el-input type="textarea" rows="6"></el-input>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </pane>
        <pane min-size="0" max-size="100">
            <!-- table 2 -->
            <div style="height: 100%; width: 100%;">
                <el-tabs v-if="form.repsonse.length > 1">
                    <el-tab-pane label="Body">
                        <el-tabs type="border-card">
                            <el-tab-pane label="Pretty">
                                <el-scrollbar class="fillcontent">
                                    <pre
                                        class="pre-wrap"><code class="language-html line-numbers">{{ form.repsonse }}</code></pre>
                                </el-scrollbar>
                            </el-tab-pane>
                            <el-tab-pane label="Raw">
                                <el-scrollbar class="fillcontent">
                                    <pre class="pre-wrap"><code>{{ form.repsonse }}</code></pre>
                                </el-scrollbar>
                            </el-tab-pane>
                        </el-tabs>
                    </el-tab-pane>
                    <el-tab-pane label="Headers">
                        <el-table :data="form.respHeader" stripe border style="width: 100%; height: 350px;">
                            <el-table-column prop="key" label="key" width="300px" />
                            <el-table-column prop="value" label="value" />
                        </el-table>
                    </el-tab-pane>
                </el-tabs>

                <el-empty description="请输入URL点击Send发送获取响应" v-else />
            </div>

        </pane>
    </splitpanes>
</template>

<style scoped>
.pre-wrap {
    white-space: pre-wrap;
    word-wrap: break-word;
    word-break: break-all;
    overflow-wrap: break-word;
}

.fillcontent {
    height: calc(80vh - 300px);
}</style>