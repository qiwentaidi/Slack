<template>
    <el-form @submit.native.prevent="tableCtrl.addTab(uncover.query)">
        <el-form-item>
            <div class="head">
                <el-input v-model="uncover.query" placeholder="多目标用英文逗号,分割" style="margin-bottom: 10px;">
                    <template #prepend>
                        <el-select size="large" style="width: 150px;" v-model="uncover.currentGroup">
                            <el-option v-for="filed in group" :value="filed" :label="filed" />
                        </el-select>
                    </template>
                    <template #suffix>
                        <el-divider direction="vertical" />
                        <el-button link type="primary" :icon="Search"
                            @click="tableCtrl.addTab(uncover.query)">查询</el-button>
                    </template>
                    <template #prefix>
                        <el-popover placement="bottom-end" :width="500" trigger="click">
                            <template #reference>
                                <div>
                                    <el-tooltip content="使用须知" placement="bottom">
                                        <el-button link :icon="Management"></el-button>
                                    </el-tooltip>
                                </div>
                            </template>
                            <el-table :data="datas" stripe>
                                <el-table-column label="字段" prop="filed" />
                                <el-table-column label="FOFA" prop="fofa" />
                                <el-table-column label="Hunter" prop="hunter" />
                                <el-table-column label="Quake" prop="quake" />
                            </el-table>
                        </el-popover>
                    </template>
                </el-input>
                <el-input v-model="uncover.size" @change="handleInput" style="width: 240px; height: 40px; margin-left: 5px;">
                    <template #prepend>
                        查询数量
                    </template>
                </el-input>
            </div>
        </el-form-item>
    </el-form>
    <el-tabs v-model="table.acvtiveNames" v-loading="table.loading" type="card" closable
        @tab-remove="tableCtrl.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 72vh;">
                <el-table-column prop="URL" label="URL" width="200px" :show-overflow-tooltip="true" />
                <el-table-column prop="IP" label="IP" width="170px" />
                <el-table-column prop="Domain" label="域名" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="Port" label="端口" width="100" />
                <el-table-column prop="Protocol" label="协议" width="100" />
                <el-table-column prop="Title" label="标题" width="200" :show-overflow-tooltip="true" />
                <el-table-column prop="Components" label="产品应用/版本" width="200">
                    <template #default="scope">
                        <el-space>
                            <el-tag effect="light" type="info" round
                                v-if="handleComponents(scope.row.Components) > 0">{{
                                    scope.row.Components.split(',')[0] }}</el-tag>
                            <el-popover placement="bottom" :width="350" trigger="hover">
                                <template #reference>
                                    <el-button round size="small" v-if="handleComponents(scope.row.Components) > 0">共{{
                                        scope.row.Components.split(',').length }}条</el-button>
                                </template>
                                <template #default>
                                    <el-space direction="vertical">
                                        <el-tag round v-for="component in scope.row.Components.split(',')"
                                            style="width: 320px;">{{ component }}</el-tag>
                                    </el-space>
                                </template>
                            </el-popover>
                        </el-space>
                    </template>
                </el-table-column>
                <el-table-column prop="Source" width="100" label="来源" />
                <el-table-column fixed="right" width="100" label="操作" align="center">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="top">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="C段查询" placement="top">
                            <el-button link @click.prevent="tableCtrl.addTab(CsegmentIpv4(scope.row.IP))">
                                <template #icon>
                                    <svg t="1719219479838" class="icon" viewBox="0 0 1450 1024" version="1.1"
                                        xmlns="http://www.w3.org/2000/svg" p-id="5099" width="200" height="200">
                                        <path
                                            d="M1055.229398 0c204.270977 0 352.3973 139.935005 352.3973 339.939672 0 5.972836-4.52229 10.83643-10.153821 10.83643h-122.869761c-9.04458 0-16.382635-7.423381-17.321223-17.065245-8.617949-114.166486-86.606116-192.325306-201.369886-192.325306-141.556204 0-221.250896 111.777352-221.250895 312.464628V575.098742c0 197.615532 79.950671 308.027664 221.165569 308.027664 114.337139 0 192.154654-73.039247 201.455212-179.35572 0.853262-9.471211 8.276644-16.894592 17.40655-16.894592h123.040413c5.631531 0 10.153821 4.863595 10.15382 10.83643 0 190.533456-148.808933 326.116824-352.3973 326.116824-238.401467 0-375.691359-168.775269-375.691359-449.498542V453.935505c0-282.771102 137.375219-453.850179 375.435381-453.850179zM552.31664 850.019832c20.136989 0 36.434297 16.638613 36.434297 37.202233v55.80335a36.860928 36.860928 0 0 1-36.434297 37.287559H79.097409a36.860928 36.860928 0 0 1-36.434298-37.287559v-55.80335c0-20.56362 16.297309-37.202233 36.434298-37.202233h473.219231zM461.358887 477.656195c20.051662 0 36.348971 16.638613 36.348971 37.202233v55.888676a36.860928 36.860928 0 0 1-36.348971 37.202234H79.097409A36.860928 36.860928 0 0 1 42.663111 570.832431v-55.888676c0-20.478293 16.297309-37.202233 36.434298-37.202233h382.261478z m90.957753-372.363636c20.136989 0 36.434297 16.72394 36.434297 37.287559v55.80335a36.860928 36.860928 0 0 1-36.434297 37.287559H79.097409A36.860928 36.860928 0 0 1 42.663111 198.383468v-55.80335c0-20.56362 16.297309-37.287559 36.434298-37.287559h473.219231z"
                                            fill="#A3A9B3" p-id="5100"></path>
                                    </svg>
                                </template>
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="my-header" style="margin-top: 10px;">
                <el-button :icon="Grid"
                    @click='ExportToXlsx(["URL", "IP", "域名", "端口", "协议", "标题", "应用/组件", "来源"], "asset", "数据聚合", item.content)'>保存数据</el-button>
                <span>共查询到数据: {{ item.total }} 条</span>
            </div>
        </el-tab-pane>
        <el-empty v-else></el-empty>
    </el-tabs>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import { Search, Management, ChromeFilled, Grid } from '@element-plus/icons-vue';
import { Uncover } from '@/interface';
import { UncoverSearch } from 'wailsjs/go/main/App';
import global from '@/global';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
import { CsegmentIpv4 } from '@/util';
import { ExportToXlsx } from '@/export';

const group = ["IP", "域名", "标题", "Body", "备案名称", "备案号"]

const uncover = reactive({
    query: '',
    currentGroup: "IP",
    size: 100,
})

const datas = [
    { filed: "IP", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "域名", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "标题", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "Body", fofa: "✓", hunter: "✓", quake: "✓" },
    { filed: "备案名称", fofa: "✓", hunter: "✓", quake: "VIP" },
    { filed: "备案号", fofa: "✓", hunter: "✓", quake: "VIP" },
    { filed: "条件拼接数量", fofa: "100", hunter: "5", quake: "1000" },
    { filed: "单页API查询数", fofa: "10000", hunter: "100", quake: "500" }
]

const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as Uncover[],
    loading: false,
})

const tableCtrl = ({
    addTab: async (query: string) => {
        const newTabName = `${++table.tabIndex}`
        table.loading = true
        let result: any = await UncoverSearch(query, uncover.currentGroup, uncover.size, {
            fofaapi: global.space.fofaapi,
            fofaemail: global.space.fofaemail,
            fofakey: global.space.fofakey,
            hunterkey: global.space.hunterkey,
            quakekey: global.space.quakekey,
        })
        table.editableTabs.push({
            title: query,
            name: newTabName,
            content: result,
            total: result.length,
            pageSize: 100,
            currentPage: 1,
        });
        table.acvtiveNames = newTabName
        table.loading = false
    },
    removeTab: (targetName: string) => {
        const tabs = table.editableTabs
        let activeName = table.acvtiveNames
        if (activeName === targetName) {
            tabs.forEach((tab, index) => {
                if (tab.name === targetName) {
                    tab.content = [] // 清理内存
                    const nextTab = tabs[index + 1] || tabs[index - 1]
                    if (nextTab) {
                        activeName = nextTab.name
                    }
                }
            })
        }
        table.acvtiveNames = activeName
        table.editableTabs = tabs.filter((tab) => tab.name !== targetName)
    },
})

function handleInput(val: string) {
    // Remove all non-numeric characters
    console.log(val)
    const numericValue = val.replace(/\D/g, '');
    // Update the model value
    uncover.size = Number(numericValue)
}

function handleComponents(info: string) {
    if (info.includes(",")) {
        return info.split(",").length
    } else {
        return 0
    }
}

</script>