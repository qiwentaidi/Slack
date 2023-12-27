<script setup lang="ts">
import { reactive } from 'vue';
import { Menu, Search, ChatLineRound, ArrowDown, ChromeFilled } from '@element-plus/icons-vue';
import { TableTabs, ApiSyntaxCheck } from '../../util'
import global from "../Global.vue"
import {
    DefaultOpenURL,
    HunterTips
} from '../../../wailsjs/go/main/App'
const form = reactive({
    query: '',
    syntaxDialog: false,
    hunterImg: 'hunter.png',
    optionsTime: [
        {
            value: '0',
            label: '最近一个月',
        },
        {
            value: '1',
            label: '最近半年',
        },
        {
            value: '2',
            label: '最近一年',
        },
    ],
    defaultTime: '0',
    optionsServer: [
        {
            value: '3',
            label: '全部资产',
        },
        {
            value: '1',
            label: 'WEB服务资产',
        },
        {
            value: '2',
            label: '非WEB服务资产',
        },
    ],
    defaultSever: '3',
    tips: '',
    loadAll: [] as EntryTips[],
})

interface EntryTips {
    value: string
    assetNum: number
    tags: string[]
}


let timeout: ReturnType<typeof setTimeout>
const entry = reactive({
    querySearchAsync: (queryString: string, cb: (arg: any) => void) => {
        if (queryString.length >= 1) {
            entry.getTips()
            clearTimeout(timeout)
            timeout = setTimeout(() => {
                cb(form.loadAll)
            }, 3000 * Math.random())
        } else {
            cb([])
        }
    },
    getTips: function () {
        HunterTips(form.query).then(result => {
            form.loadAll = []
            if (result.code == 200) {
                for (const item of result.data.app) {
                    form.loadAll.push({
                        value: item.name,
                        assetNum: item.asset_num,
                        tags: item.tags
                    })
                }
            }
        })
    },
    handleSelect: (item: Record<string, any>) => {
        form.query = `app.name="${item.value}"`
    }
})

const table = reactive({
    acvtiveNames: "1",
    tabIndex: 1,
    editableTabs: [] as TableTabs[],
    addTab: (query: string) => {
        if (ApiSyntaxCheck(1, "", global.space.hunterkey, query) === false) {
            return
        }
    },
    removeTab: (targetName: string) => {
        const tabs = table.editableTabs
        let activeName = table.acvtiveNames
        if (activeName === targetName) {
            tabs.forEach((tab, index) => {
                if (tab.name === targetName) {
                    tab.content = null // 清理内存
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
    handleSizeChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;

    },
    handleCurrentChange: (val: any) => {
        const tab = table.editableTabs.find(tab => tab.name === table.acvtiveNames)!;
        tab.currentPage = val
    },
})


</script>

<template>
    <el-form v-model="form">
        <el-form-item label="查询条件">
            <div class="head">
                <el-autocomplete v-model="form.query" placeholder="Search..." :fetch-suggestions="entry.querySearchAsync"
                    @select="entry.handleSelect"  :trigger-on-focus="false" style="width: 100%;">
                    <template #append>
                        <el-dropdown>
                            <el-button :icon="Menu" />
                            <template #dropdown>
                                <el-dropdown-menu :hide-on-click="true">
                                    <el-dropdown-item @click="">icon搜索</el-dropdown-item>
                                    <el-dropdown-item @click="">批量查询</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </template>
                    <template #default="{ item }">
                        
                        <el-space>
                            <span>{{ item.value }}</span>
                            <span>({{ item.assetNum }} 条数据)</span>
                            <el-tag v-for="label in item.tags">{{ label }}</el-tag>
                        </el-space>
                    </template>
                </el-autocomplete>
                <el-button type="primary" :icon="Search" @click="table.addTab(form.query)"
                    style="margin-left: 10px; margin-right: 10px;">查询</el-button>
                <el-dropdown>
                    <el-button color="#A29EDE">
                        数据导出<el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item>导出当前查询页数据</el-dropdown-item>
                            <el-dropdown-item>导出全部数据</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
        </el-form-item>
    </el-form>
    <div class="nkmode">
        <el-space>
            <el-select v-model="form.defaultTime" style="width: 120px;">
                <el-option v-for="item in form.optionsTime" :key="item.value" :label="item.label" :value="item.value"
                    style="text-align: center;" />
            </el-select>
            <el-select v-model="form.defaultSever" style="width: 160px;">
                <el-option v-for="item in form.optionsServer" :key="item.value" :label="item.label" :value="item.value"
                    style="text-align: center;" />
            </el-select>
            <el-checkbox size="large">数据去重(需权益积分)</el-checkbox>
        </el-space>
        <el-link @click="form.syntaxDialog = true"><el-icon>
                <ChatLineRound />
            </el-icon>查询语法</el-link>
        <el-dialog v-model="form.syntaxDialog" title="查询语法参考" width="80%" center>
            <el-scrollbar height="400px">
                <el-image :src="form.hunterImg"></el-image>
            </el-scrollbar>
        </el-dialog>
    </div>
    <el-tabs v-model="table.acvtiveNames" type="card" style="margin-top: 10px;" closable @tab-remove="table.removeTab">
        <el-tab-pane v-for="item in table.editableTabs" :key="item.name" :label="item.title" :name="item.name"
            v-if="table.editableTabs.length != 0">
            <el-table :data="item.content" border style="width: 100%;height: 65vh;">
                <el-table-column type="index" label="#" width="60px" />
                <el-table-column prop="URL" label="URL" width="200" show-overflow-tooltip="true" />
                <el-table-column prop="IP" label="IP" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="PortAndServer" label="端口/服务" width="120" show-overflow-tooltip="true" />
                <el-table-column prop="Domain" label="域名" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Assembly" label="应用/组件" width="200" show-overflow-tooltip="true" />
                <el-table-column prop="Title" label="标题" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Status" label="状态码" show-overflow-tooltip="true" />
                <el-table-column prop="ICP" label="备案号" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="Position" label="地理位置" width="150" show-overflow-tooltip="true" />
                <el-table-column prop="UpdateTime" label="更新时间" width="150" show-overflow-tooltip="true" />

                <el-table-column fixed="right" label="操作" width="55px">
                    <template #default="scope">
                        <el-tooltip content="打开链接" placement="left">
                            <el-button link :icon="ChromeFilled" @click.prevent="DefaultOpenURL(scope.row.link)">
                            </el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>
            <div class="nkmode" style="margin-top: 10px;">
                <span style="color: cornflowerblue;">{{ form.tips }}</span>
                <el-pagination :page-size="100" :page-sizes="[100, 500, 1000]" layout="sizes, prev, pager, next"
                    @size-change="table.handleSizeChange" @current-change="table.handleCurrentChange" :total="item.total" />
            </div>
        </el-tab-pane>
        <el-image class="center" src="/loading.gif" alt="loading" v-else></el-image>
    </el-tabs>
</template>

<style></style>