<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { Search, CirclePlusFilled, Delete, Upload } from "@element-plus/icons-vue";
import { FileDialog, ReadFile, RemoveFile, WriteFile } from 'wailsjs/go/main/File';
import global from '@/global';
import { FingerprintList, GetFingerPocMap } from 'wailsjs/go/main/App';
import { Copy } from '@/util';
import { PocDetail } from '@/stores/interface';
import usePagination from '@/usePagination';
import { ElMessage } from 'element-plus';
import CustomTabs from '@/components/CustomTabs.vue';
import saveIcon from '@/assets/icon/save.svg'
import bugIcon from '@/assets/icon/bug.svg'
import fingerpirntIcon from '@/assets/icon/fingerprint.svg'
import { SaveConfig } from '@/config';

onMounted(async () => {
    const pocMap = await GetFingerPocMap();
    for (const [poc, tags] of Object.entries(pocMap)) {
        pagination.table.result.push({
            Name: poc,
            AssociatedFingerprint: Array.from(new Set(tags))
        })
    }
    pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
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

const pocs = ref<PocDetail[]>([])

const pagination = usePagination(pocs.value, 20)

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

const filterId = ref(0)
function filterPocList() {
    if (filterId.value == 0) {
        pagination.table.temp = pagination.table.result
        filterId.value++
    }
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

const detailDialog = ref(false)
const content = ref('')
async function readPoc(filename: string) {
    detailDialog.value = true
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + filename + ".yaml"
    let file: any = await ReadFile(filepath)
    content.value = file.Content!
}

// step 1 add poc, hide poclist
const step = ref(0)


const poc = reactive({
    Name: '',
    Content: '',
})

async function checkDuplicate() {
    const pocMap = await GetFingerPocMap();
    if (pocMap[poc.Name]) {
        ElMessage.warning("Poc name already exists")
        return true
    }
    return false
}

const selectFinger = ref<string[]>([])
const selectHighlightFinger = ref<string[]>([])

const fingerOptions = ref<{ label: string, value: string }[]>([])
var highlightFingerOptions = ref<{ label: string, value: string }[]>([])

async function importFile() {
    let path = await FileDialog("*.yaml")
    if (!path) {
        return
    }
    let file: any = await ReadFile(path)
    poc.Content = file.Content!
    poc.Name = poc.Content.split('\n')[0].replaceAll("id:", "").trim()
}

async function savePoc() {
    if (await checkDuplicate()) {
        return
    }
    if (!poc.Content.includes("tags:")) {
        ElMessage.warning("Nuclei poc must contain tags")
        return
    }
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + poc.Name + ".yaml"
    let result = await WriteFile("yaml", filepath, poc.Content)
    pagination.table.result.push({
        Name: poc.Name,
        AssociatedFingerprint: selectFinger.value
    })
    pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table)
    result ? ElMessage.success("Poc saved successfully") : ElMessage.error("Poc save failed")
}

async function deletePoc(pocName: string) {
    let filepath = global.PATH.homedir + "/slack/config/pocs/" + pocName + ".yaml"
    if (await RemoveFile(filepath)) {
        ElMessage.success("Poc deleted successfully")
    } else {
        ElMessage.warning("Poc delete failed")
    }
}

const activeTabs = ref("poc")

function addFingerprint(fingerprints: string[]) {
    ElMessage.success("添加成功")
    global.webscan.highlight_fingerprints.push(...fingerprints)
}
function deleteFingerprint(fingerprint: string) {
    ElMessage.success("删除成功")
    global.webscan.highlight_fingerprints.splice(global.webscan.highlight_fingerprints.indexOf(fingerprint), 1)
}
</script>

<template>
    <CustomTabs>
        <el-tabs v-model="activeTabs" type="card" class="demo-tabs">
            <el-tab-pane name="poc">
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon>
                            <bugIcon />
                        </el-icon>
                        <span>POC管理</span>
                    </span>
                </template>
                <el-card v-show="step == 0">
                    <template #header>
                        <div class="my-header">
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
                            <el-button :icon="CirclePlusFilled" @click="step = 1">POC</el-button>
                        </div>
                    </template>
                    <el-table :data="pagination.table.pageContent" style="height: calc(100vh - 230px);">
                        <el-table-column prop="Name" label="名称" />
                        <el-table-column label="关联指纹" width="400px">
                            <template #default="scope">
                                <div class="finger-container">
                                    <el-tag v-for="item in scope.row.AssociatedFingerprint">{{ item }}</el-tag>
                                </div>
                            </template>
                        </el-table-column>
                        <el-table-column label="操作" width="200px" align="center">
                            <template #default="scope">
                                <el-button type="primary" link @click="readPoc(scope.row.Name)">查看</el-button>
                                <el-popconfirm title="Are you sure to delete poc?" @confirm="deletePoc(scope.row.Name)">
                                    <template #reference>
                                        <el-button type="primary" link>移除</el-button>
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
                <div v-show="step == 1">
                    <el-page-header @back="step = 0">
                        <template #content>
                            <span>添加POC</span>
                        </template>
                        <template #extra>
                            <el-button-group>
                                <el-button type="primary" :icon="Upload" @click="importFile">导入</el-button>
                                <el-button type="primary" :icon="saveIcon" @click="savePoc">保存</el-button>
                            </el-button-group>
                        </template>
                    </el-page-header>
                    <el-divider />
                    <el-form :model="poc" label-width="auto">
                        <el-form-item label="POC名称">
                            <el-input v-model="poc.Name" placeholder="优先填写CVE、CNVD等编号，导入会自动读取id值" />
                        </el-form-item>
                        <el-form-item label="指纹查询">
                            <el-select-v2 v-model="selectFinger" placeholder="POC必须要填写tags信息，相关指纹可通过此处查找" filterable
                                :options="fingerOptions" multiple clearable />
                        </el-form-item>
                        <el-form-item label="POC内容">
                            <el-input v-model="poc.Content" type="textarea" resize="none" placeholder="Nulcei Yaml POC"
                                style="height: 65vh;"></el-input>
                        </el-form-item>
                    </el-form>
                </div>
            </el-tab-pane>
            <el-tab-pane name="finger">
                <template #label>
                    <span class="custom-tabs-label">
                        <el-icon>
                            <fingerpirntIcon />
                        </el-icon>
                        <span>指纹管理</span>
                    </span>
                </template>
                <el-card>
                    <template #header>
                        <div class="my-header">
                            <el-select-v2 v-model="selectHighlightFinger" placeholder="请选择需要高亮的指纹" filterable
                                :options="highlightFingerOptions" multiple clearable />
                            <el-space style="margin-left: 10px">
                                <el-button :icon="CirclePlusFilled" type="primary" @click="addFingerprint(selectHighlightFinger)">添加</el-button>
                                <el-button :icon="saveIcon" type="primary" @click="SaveConfig">保存</el-button>
                            </el-space>
                        </div>
                    </template>
                    <el-table :data="global.webscan.highlight_fingerprints" style="height: calc(100vh - 200px);">
                        <el-table-column type="index" width="50" />
                        <el-table-column label="Fingerprint">
                            <template #default="scope">
                                {{ scope.row }}
                            </template>
                        </el-table-column>
                        <el-table-column label="Operation" width="150" align="center">
                            <template #default="scope">
                                <el-button type="danger" :icon="Delete" @click="deleteFingerprint(scope.row)" />
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

    <el-drawer v-model="detailDialog" size="70%">
        <template #header>
            <el-button text bg>
                <template #icon>
                    <Notebook />
                </template>漏洞详情</el-button>
        </template>
        <div class="controls">
            <el-button type="primary" link @click="Copy(content)">复制</el-button>
        </div>
        <highlightjs language="yaml" :code="content"></highlightjs>
    </el-drawer>
</template>

<style scoped>
:deep(.el-card__header) {
    padding: 5px;
}

.controls {
    position: absolute;
    top: 90px;
    right: 30px;
    z-index: 1000;
    /* Ensure it's above the code */
}
</style>