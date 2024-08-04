<script lang="ts" setup>
import global from "@/global";
import async from 'async';
import { ReadLine } from '@/util'
import { ExportToXlsx } from '@/export'
import { reactive, ref } from "vue";
import { Subdomain, InitIPResolved, LoadSubDict } from "wailsjs/go/main/App";
import { FileDialog } from "wailsjs/go/main/File";
import { ElMessage } from 'element-plus'
import { onMounted } from 'vue';
import { Loading } from '@element-plus/icons-vue';

// 初始化时调用
onMounted(() => {
    sbr.value = [];
});
const from = reactive({
    domain: "",
    thread: 600,
    timeout: 5,
    tips: '选择子域字典(默认加载dicc.txt)',
    subs: [] as string[],
    percentage: 0,
    id: 0,
    runningStatus: false,
});
const sbr = ref([{}]);
async function BurstSubdomain() {
    from.runningStatus = true
    sbr.value = [];
    from.id = 0;
    InitIPResolved();
    if (from.subs.length === 0) {
        from.subs = await LoadSubDict(global.PATH.ConfigPath + "/subdomain")
        from.tips = `loaded ${from.subs.length} dicts`
    }
    async.eachLimit(from.subs, from.thread, (sub: string, callback: () => void) => {
        from.id++;
        from.percentage = Number(((from.id / from.subs.length) * 100).toFixed(2));
        if (!from.runningStatus) {
            return
        }
        Subdomain(sub + "." + from.domain, from.timeout).then((result) => {
            if (result[2].length > 0) {
                sbr.value.push({
                    subdomains: result[0],
                    cname: result[1],
                    ips: result[2],
                    notes: result[3],
                });
            }
            callback();
        });
    }, (err: any) => {
        if (err) {
            ElMessage.error(err)
        } else {
            ElMessage({
                showClose: true,
                message: '子域名暴破已完成',
                type: 'success',
            })
        }
        from.runningStatus = false
    });
}
function stop() {
    if (from.runningStatus == true) {
        from.runningStatus = false
        ElMessage({
            showClose: true,
            message: '已停止任务',
            type: 'warning',
        })
    }
}

async function handleFileChange() {
    let path = await FileDialog("*.txt")
    from.subs = (await ReadLine(path))!
    from.tips = `loaded ${from.subs.length} dicts`;
}
</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head">
                <el-input v-model="from.domain" placeholder="请输入域名,仅支持单域名" style="margin-right: 10px;" />
                <el-button type="primary" @click="BurstSubdomain" v-if="!from.runningStatus">开始任务</el-button>
                <el-button type="danger" @click="stop" v-else>停止任务</el-button>
            </div>
        </el-form-item>
        <el-form-item>
            <el-space>
                <div>
                    <span>线程数量：</span>
                    <el-input-number v-model="from.thread" :min="600" :max="10000" controls-position="right">
                    </el-input-number>
                </div>
                <div>
                    <span>解析超时(s)：</span>
                    <el-input-number v-model="from.timeout" :min="1" :max="20" controls-position="right">
                    </el-input-number>
                </div>
                <el-button type="primary" :icon="Loading" @click="handleFileChange()">{{ from.tips
                    }}</el-button>
            </el-space>
            <el-button style="margin-left: auto;"
                @click="ExportToXlsx(['子域名', 'CNAME', 'IPS', '备注'], '子域名暴破', 'subdomain', sbr)">
                <template #icon>
                    <img src="/excel.svg" width="16">
                </template>
                导出Excel</el-button>
        </el-form-item>
    </el-form>
    <el-table :data="sbr" border style="height: 75vh; margin-bottom: 10px;">
        <el-table-column type="index" label="#" width="60px" />
        <el-table-column prop="subdomains" label="子域名" :show-overflow-tooltip="true" />
        <el-table-column prop="cname" label="CNAME" :show-overflow-tooltip="true" />
        <el-table-column prop="ips" label="IPs" :show-overflow-tooltip="true" />
        <el-table-column prop="notes" label="备注" :show-overflow-tooltip="true" />
        <template #empty>
            <el-empty description="点击开始任务获取数据"></el-empty>
        </template>
    </el-table>
    <el-progress :text-inside="true" :stroke-width="18" :percentage="from.percentage" color="#5DC4F7" />
</template>
