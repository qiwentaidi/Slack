<script setup lang="ts">
import async from 'async';
import { Search, Filter, QuestionFilled } from '@element-plus/icons-vue';
import { splitInt, ExportTXT } from '../../util'
import { reactive, onMounted } from 'vue';
import { Sock5UnauthScan, FofaSearch } from '../../../wailsjs/go/main/App'
import { CreateTable } from '../../../wailsjs/go/main/File'
import global from "../../global"

// onMounted(async () => {
//     await CreateTable()
// });

const form = reactive({
    socksLogger: '',
    socksMax: 1,
    socksNum: 1,
    percentage: 0,
    activeNames: "1",
    pool: [] as string[],
    currentTableName: "0",
})


let hosts = [] as string[]

async function NewSock5Crawl(step: number) {
    let query = `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2021"  && country="CN" && region!="HK" && region!="TW"`
    let sc = new Socks5Crawling()
    let date = new Date()
    if (step == 0) {
        form.socksLogger += date.toLocaleString() + " 正在查询数据量中...\n"
        form.socksMax = await sc.SearchTotal(query)
        form.socksNum = form.socksMax
    } else if (step == 1) {
        let results = await sc.TargetExtraction(query)
        let id = 0
        if (results != undefined) {
            let count = results.length
            form.socksLogger += "筛选存活资产数据:" + count.toString() + "条\n"
            async.eachLimit(results, 50, async (temp: Hosts, callback: (err?: Error) => void) => {
                id++
                form.percentage = Number(((id / count) * 100).toFixed(2));
                if (await Sock5UnauthScan(temp.IP, Number(temp.Port), 3)) {
                    form.socksLogger += temp.IP + ":" + temp.Port + " is unauthorized!\n"
                    hosts.push(temp.IP + ":" + temp.Port)
                }
            })
        }
    } else if (step == 2) {
        ExportTXT("socks_unauth_asset", hosts)
        hosts = []
    }
}

interface Hosts {
    IP: string
    Port: string
}

class Socks5Crawling {
    public async SearchTotal(query: string) {
        let result = await FofaSearch(query, "1", "1", global.space.fofaemail, global.space.fofakey, true, false)
        form.socksLogger += "共查询到数据:" + result.Total + "条" + result.Message + "\n"
        return Number(result.Total)
    }

    public async TargetExtraction(query: string) {
        let index = 0
        let temps = [] as Hosts[]
        for (const num of splitInt(form.socksNum, 10000)) {
            index += 1
            let result = await FofaSearch(query, num.toString(), index.toString(), global.space.fofaemail, global.space.fofakey, true, false)
            if (result.Status == false) {
                form.socksLogger += "查询异常，已退出存活测试" + result.Message + "\n"
                return
            }
            for (const temp of result.Results) {
                temps.push({
                    IP: temp.IP,
                    Port: temp.Port
                })
            }
        }
        return temps
    }
}
</script>

<template>
    <div style="position: relative;">
        <el-tabs v-model="form.currentTableName" type="card">
            <el-tab-pane name="0">
                <template #label>
                    代理池爬取<el-popover placement="right-end" title="此模块需要配置FOFA Email+key" :width="350" trigger="hover">
                        ①点击<b>查询数据量</b>确定FOFA中的测绘资产数量<br /><br />
                        ②<b>检测数量</b>表示拉取和需要验证的目标数量<br /><br />
                        ③<b>存储阈值</b>表示数据库中存储的上限数量，如果超出立刻停止检测<br /><br />
                        ④<b>筛选存活</b>即开始检测
                        <template #reference>
                            <el-icon>
                                <QuestionFilled size="24" />
                            </el-icon>
                        </template>
                    </el-popover>
                </template>
                <el-input v-model="form.socksLogger" type="textarea" rows="20" resize="none" readonly
                    class="log-textarea"></el-input>
                <el-progress :percentage="form.percentage" :stroke-width="15" striped striped-flow :duration="20"
                    style="margin-top: 5px" />
            </el-tab-pane>
            <el-tab-pane label="数据库信息" name="1">
                <el-table :data="form.pool">
                    <el-table-column prop="host" label="主机地址" />
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <div class="custom_eltabs_titlebar" v-if="form.currentTableName == '0'">
            <el-space>
                <el-button :icon="Search" @click="NewSock5Crawl(0)">查询数据量</el-button>
                <span>检测数量:</span><el-input-number v-model="form.socksNum" :min="1" :max="form.socksMax"
                    controls-position="right" style="width: 100px;"></el-input-number>
                <span>存储阈值:</span><el-input-number v-model="form.socksNum" :min="1" :max="form.socksMax"
                    controls-position="right" style="width: 100px;"></el-input-number>
            </el-space>
            <el-button :icon="Filter" @click="NewSock5Crawl(1)">筛选存活</el-button>
        </div>
        <el-button type="primary" @click="NewSock5Crawl(2)" class="custom_eltabs_titlebar" v-else>导出存活目标</el-button>
    </div>
</template>

<style>
.custom_eltabs_titlebar {
    position: absolute;
    right: 0px;
    top: 4px;
}
</style>
