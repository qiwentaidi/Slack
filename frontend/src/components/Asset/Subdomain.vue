<script lang="ts" setup>
import global from "../Global.vue";
import async from 'async';
import { ExportToXlsx } from '../../util'
import { reactive, ref, } from "vue";
import { Subdomain, InitIPResolved, SelectFile, GetFileContent } from "../../../wailsjs/go/main/App";
import { ElMessage } from 'element-plus'
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    sbr.value = [];
});
const from = reactive({
    domain: "",
    thread: 600,
    timeout: 5,
    tips: '',
    subs: [{}],
    percentage: 0,
    id: 0,
});
const sbr = ref([{}]);
const exit = ref(false)
function BurstSubdomain() {
    exit.value = false
    sbr.value = [];
    from.id = 0;
    InitIPResolved();
    async.eachLimit(from.subs, from.thread, (sub: string, callback: () => void) => {
        from.id++;
        from.percentage = Number(((from.id / from.subs.length) * 100).toFixed(2));
        if (exit.value) {
            return
        }
        Subdomain(sub + "." + from.domain, global.scan.dns1, global.scan.dns2, from.timeout).then((result) => {
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
    });
}
function stop() {
    if (exit.value === false) {
        exit.value = true
        ElMessage({
            showClose: true,
            message: '已停止任务',
            type: 'warning',
        })
    }
}

async function handleFileChange() {
    let path = await SelectFile()
    let res = await GetFileContent(path)
    if (res.length == 0) {
        ElMessage({
            showClose: true,
            message: '不能上传空文件',
            type: 'warning',
        })
        return
    }
    if ( res !== "文件不存在") {
        const result = res.replace(/\r\n/g, '\n'); // 避免windows unix系统差异
        from.subs = Array.from(result.split('\n'))
        from.tips = `loaded ${from.subs.length} dicts`;
    }
}
</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head">
                <el-input v-model="from.domain" placeholder="请输入域名,仅支持单域名" />
                <el-button type="primary" style="margin-left: 10px" @click="BurstSubdomain">开始暴破</el-button>
                <el-button type="primary" style="margin-left: 10px" @click="stop">停止</el-button>
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
                <div style="display: flex;">
                    <el-button type="primary" @click="handleFileChange()">选择子域字典</el-button>
                    <span style="margin-left: 10px;">{{ from.tips }}</span>
                </div>
            </el-space>
            <el-button type="primary" style="margin-left: auto;" text
                @click="ExportToXlsx(['子域名', 'CNAME', 'IPS', '备注'], '子域名暴破', 'subdomain', sbr)"
                :disabled="sbr.length < 2">数据导出</el-button>
        </el-form-item>
        <el-form-item>
                <el-table :data="sbr" border max-height="calc(100vh - 220px)">
                    <el-table-column type="index" label="#" width="60px" />
                    <el-table-column prop="subdomains" label="子域名" show-overflow-tooltip="true" />
                    <el-table-column prop="cname" label="CNAME" show-overflow-tooltip="true" />
                    <el-table-column prop="ips" label="IPs" show-overflow-tooltip="true" />
                    <el-table-column prop="notes" label="备注" show-overflow-tooltip="true" />
                </el-table>
        </el-form-item>
    </el-form>
    <el-progress :text-inside="true" :stroke-width="18" :percentage="from.percentage"  />
</template>
