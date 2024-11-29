<script setup lang="ts">
import async from 'async';
import { Search, Filter, QuestionFilled } from '@element-plus/icons-vue';
import { splitInt } from '@/util'
import { ExportTXT } from '@/export'
import { reactive, onMounted } from 'vue';
import { Socks5Conn, FofaSearch } from 'wailsjs/go/services/App'
import global from "@/global"
import { SearchAgentPool, ExecSqlStatement } from 'wailsjs/go/services/Database';
import { ElMessage } from 'element-plus';
import CustomTabs from '@/components/CustomTabs.vue';

onMounted(async () => {
    let hosts = await SearchAgentPool()
    if (Array.isArray(hosts)) form.pool = hosts
});

const form = reactive({
    socksLogger: '',
    socksMax: 1,
    socksNum: 1,
    socksThreshold: 50,
    percentage: 0,
    activeNames: "1",
    pool: [] as string[],
    currentTableName: "0",
})

async function NewSock5Crawl(step: number) {
    let query = `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2021"  && country="CN" && region!="HK" && region!="TW"`
    let sc = new Socks5Crawling()
    if (step == 0) {
        form.socksLogger += `正在查询数据量中...\n`
        form.socksMax = await sc.SearchTotal(query)
        form.socksNum = form.socksMax
    } else if (step == 1) {
        // 如果数据库中有代理池则先验证
        let tempHosts = [] as string[]
        if (form.pool.length > 0) {
            form.socksLogger += "正在验证数据库代理存活数量...\n"
            await async.eachLimit(form.pool, 50, async (temp: string, callback: (err?: Error) => void) => {
                let t = temp.split(":")
                if (await Socks5Conn(t[0], Number(t[1]), 3, "", "")) {
                    form.socksLogger += `[+] ${t[0]}:${t[1]} is alive!\n`
                    tempHosts.push(t[0] + ":" + t[1])
                }
            })
            form.socksLogger += `共验证存活代理数量:${tempHosts.length.toString()}条\n`
        }
        // 如果已存活的数量大于阈值则退出
        if (tempHosts.length >= form.socksThreshold) {
            form.socksLogger += "已收集的代理存活数量大于等于阈值，程序退出\n"
            return
        }
        let results = await sc.TargetExtraction(query)
        let id = 0
        if (results != undefined) {
            let filteredResults = [] as Hosts[]
            // 排除已有的重复结果
            if (form.pool.length > 0) {
                let resultsStr = results.map(host => host.IP + ":" + host.Port);
                let filteredResultsStr = resultsStr.filter(hostStr => !tempHosts.includes(hostStr));
                filteredResults = filteredResultsStr.map(hostStr => {
                    let parts = hostStr.split(':');
                    return { IP: parts[0], Port: parts[1] };
                });
            } else {
                filteredResults = results
            }
            let count = filteredResults.length
            form.socksLogger += `筛选存活资产数据:${filteredResults.length.toString()}条，已过滤历史探活资产:${(results.length - filteredResults.length).toString()}条\n`
            // 开始执行
            await async.eachLimit(filteredResults, 50, async (temp: Hosts) => {
                id++
                form.percentage = Number(((id / count) * 100).toFixed(2));
                // 达到阈值 || 正常执行结束 停止
                if (tempHosts.length >= form.socksThreshold) {
                    return
                }
                if (await Socks5Conn(temp.IP, Number(temp.Port), 3, "", "")) {
                    let host = temp.IP + ":" + temp.Port
                    form.socksLogger += `[+] ${host} is unauthorized!\n`
                    tempHosts.push(host)
                }
            })
        }
        if (form.pool.length > 0) {
            await ExecSqlStatement("DELETE FROM agent_pool", ["agent_pool"])
            form.pool = []
        }
        let insertStmt = "INSERT INTO agent_pool(hosts) VALUES(?)"
        for (const host of tempHosts) {
            form.pool.push(host)
            await ExecSqlStatement(insertStmt, [host])
        }
    } else if (step == 2) {
        let hosts = [] as string[]
        for (const host of form.pool) {
            hosts.push(host)
        }
        ExportTXT("socks_unauth_asset", hosts)
    }
}

interface Hosts {
    IP: string
    Port: string
}

class Socks5Crawling {
    public async SearchTotal(query: string) {
        let result = await FofaSearch(query, "1", "1", global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, false)
        form.socksLogger += "共查询到数据:" + result.Size! + "条" + result.Message! + "\n"
        return Number(result.Size!)
    }

    public async TargetExtraction(query: string) {
        let index = 0
        let temps = [] as Hosts[]
        for (const num of splitInt(form.socksNum, 10000)) {
            index += 1
            let result = await FofaSearch(query, num.toString(), index.toString(), global.space.fofaapi, global.space.fofaemail, global.space.fofakey, true, false)
            if (result.Error) {
                form.socksLogger += "查询异常，已退出存活测试" + result.Message + "\n"
                return
            }
            for (const temp of result.Results!) {
                temps.push({
                    IP: temp.IP,
                    Port: temp.Port
                })
            }
        }
        return temps
    }
}

async function TestConnection(host: string) {
    let t = host.split(":")
    let result = await Socks5Conn(t[0], Number(t[1]), 5, "", "")
    if (result) {
        ElMessage({
            message: 'This proxy is reachable',
            type: 'success',
        })
    } else {
        ElMessage({
            message: 'Oops, this proxy is unreachable.',
            type: 'error',
        })
    }
}

async function Delelte(host: string) {
    let deleteStmt = `DELETE FROM agent_pool WHERE hosts = ?;`
    let result = await ExecSqlStatement(deleteStmt, [host])
    if (result) {
        form.pool = form.pool.filter(item => item !== host)
    } else {
        ElMessage({
            message: 'failed',
            type: 'error',
        })
    }
}
</script>

<template>
    <CustomTabs>
        <el-tabs v-model="form.currentTableName" type="card">
            <el-tab-pane name="0">
                <template #label>
                    <el-popover placement="right-end" title="此模块需要配置FOFA Email+key" :width="350" trigger="hover">
                        ①<b>检测数量</b>表示拉取和需要验证的目标数量<br /><br />
                        ②<b>存储阈值</b>表示数据库中存储的上限数量，如果超出立刻停止检测<br /><br />
                        ③点击
                        <Search style="height: 16px; width: 16px;" />确定FOFA中的测绘资产数量<br /><br />
                        ④点击
                        <Filter style="height: 16px; width: 16px;" />即开始检测
                        <br /><br />
                        操作顺序③①②④
                        <template #reference>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </template>
                    </el-popover>
                    代理池爬取
                </template>
                <el-input v-model="form.socksLogger" type="textarea" resize="none" style="height: 80vh;" />
                <el-progress :percentage="form.percentage" :text-inside="true" :stroke-width="18"
                    style="margin-top: 5px" />
            </el-tab-pane>
            <el-tab-pane label="历史记录" name="1">
                <el-table :data="form.pool" border style="height: 83vh;">
                    <el-table-column type="index" width="50" label="#" align="center" />
                    <el-table-column label="主机地址">
                        <template #default="scope">
                            {{ scope.row }}
                        </template>
                    </el-table-column>
                    <el-table-column fixed="right" label="操作" width="120" align="center">
                        <template #default="scope">
                            <el-button link type="primary" size="small"
                                @click.prevent="TestConnection(scope.row)">测试连接</el-button>
                            <el-button link type="primary" size="small"
                                @click.prevent="Delelte(scope.row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <template #ctrl>
            <el-space v-show="form.currentTableName == '0'">
                <span>检测数量:</span><el-input-number v-model="form.socksNum" :min="1" :max="form.socksMax"
                    controls-position="right" style="width: 100px;"></el-input-number>
                <span>存储阈值:</span><el-input-number v-model="form.socksThreshold" :min="1" controls-position="right"
                    style="width: 100px;"></el-input-number>
                <el-button-group>
                    <el-tooltip content="查询数据量" placement="left">
                        <el-button :icon="Search" type="primary" @click="NewSock5Crawl(0)"></el-button>
                    </el-tooltip>
                    <el-tooltip content="筛选存活" placement="left">
                        <el-button :icon="Filter" type="primary" @click="NewSock5Crawl(1)"></el-button>
                    </el-tooltip>
                </el-button-group>
            </el-space>
            <el-button type="primary" @click="NewSock5Crawl(2)" v-show="form.currentTableName == '1'">导出存活目标</el-button>
        </template>
    </CustomTabs>
</template>

<style scoped></style>