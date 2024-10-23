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
                            <el-table :data="uncoverSyntaxOptions" stripe>
                                <el-table-column label="字段" prop="filed" />
                                <el-table-column label="FOFA" prop="fofa" />
                                <el-table-column label="Hunter" prop="hunter" />
                                <el-table-column label="Quake" prop="quake" />
                            </el-table>
                        </el-popover>
                    </template>
                </el-input>
                <el-input v-model.number="uncover.size" @change="handleInput" style="width: 240px; height: 40px; margin-left: 5px;">
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
                        <el-tooltip content="打开链接">
                            <el-button link :icon="ChromeFilled" @click.prevent="BrowserOpenURL(scope.row.URL)" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="C段查询">
                            <el-button :icon="csegmentIcon" link @click.prevent="tableCtrl.addTab(CsegmentIpv4(scope.row.IP))">
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
import { Uncover } from '@/stores/interface';
import { UncoverSearch } from 'wailsjs/go/main/App';
import global from '@/global';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
import { CsegmentIpv4 } from '@/util';
import { ExportToXlsx } from '@/export';
import csegmentIcon from '@/assets/icon/csegment.svg'
import { structs } from 'wailsjs/go/models';
import { uncoverSyntaxOptions } from '@/stores/options';

const group = ["IP", "域名", "标题", "Body", "备案名称", "备案号"]

const uncover = reactive({
    query: '',
    currentGroup: "IP",
    size: 100,
})


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
        let options: structs.SpaceOption = {
            FofaApi: global.space.fofaapi,
            FofaEmail: global.space.fofaemail,
            FofaKey: global.space.fofakey,
            HunterKey: global.space.hunterkey,
            QuakeKey: global.space.quakekey,
        }
        let result: any = await UncoverSearch(query, uncover.currentGroup, uncover.size, options)
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