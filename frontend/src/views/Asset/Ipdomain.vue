<script lang="ts" setup>
import { reactive } from 'vue';
import { ICPInfo, CheckCdn, Ip138IpHistory, Ip138Subdomain } from 'wailsjs/go/services/App'
import { Document, UploadFilled, Promotion, Search } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { Copy, ProcessTextAreaInput } from '@/util';
import async from 'async'
import router from '@/router';

const domain = reactive({
    input: '',
    icp: '',
    cdn: '',
    history: '',
    subdomain: '',
    batch: '',
});

const domainRegex = /^(?=.{1,253}$)([a-z0-9]([-a-z0-9]*[a-z0-9])?\.)+[a-z]{2,}$/i;
function startSearch() {
    if (!domainRegex.test(domain.input)) {
        ElMessage.warning("请输入正确域名")
        return
    }
    ElMessage.info("正在查询...")
    CheckCdn(domain.input).then(result => {
        domain.cdn = result
    })
    ICPInfo(domain.input).then(
        result => {
            domain.icp = result
        }
    )
    Ip138IpHistory(domain.input).then(
        result => {
            if (!result) {
                domain.history = "暂无历史解析"
                return
            }
            domain.history = result
        }
    )
    Ip138Subdomain(domain.input).then(
        result => {
            if (!result) {
                domain.subdomain = "暂无子域名记录"
                return
            }
            domain.subdomain = result
        }
    )
}

function batchQuery() {
    let id = 0
    const lines = ProcessTextAreaInput(domain.batch)
    ElMessage.info("正在执行CDN批量检测，目标数: " + lines.length)
    async.eachLimit(lines, 20, (target: string, callback: () => void) => {
        CheckCdn(target).then(result => {
            domain.cdn += result + "\n\n"
            id++
        })
        if (id == lines.length) {
            callback()
        }
    }, (err: any) => {
        ElMessage.success("CheckCdn Finished")
    });
}

function CopyAndJump() {
    Copy(domain.cdn)
    router.push("/Tools/DataHanding")
}
</script>

<template>
    <div class="position-center">
        <el-input v-model="domain.input" style="width: 60%; height: 40px;">
            <template #prepend>
                域名
            </template>
            <template #suffix>
                <el-divider direction="vertical" />
                <el-button type="primary" link :icon="Search" @click="startSearch()">查询</el-button>
            </template>
        </el-input>
    </div>
    <div class="grid-container">
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/chinaz.ico"></el-icon><span class="title">备案信息:</span></el-text>
            <el-input v-model="domain.icp" type="textarea" :rows="13" resize="none"
                class="mt-10px"></el-input>
        </el-card>
        <el-card shadow="never" class="grid-item">
            <div class="flex-between">
                <span class="title">CDN信息:</span>
                <div class="flex">
                    <el-popover placement="left" :width="350" trigger="click">
                        <template #reference>
                            <div>
                                <el-tooltip content="批量查询">
                                    <el-button :icon="Document"></el-button>
                                </el-tooltip>
                            </div>
                        </template>
                        <el-button text bg :icon="UploadFilled" class="w-full mb-5px">选择域名文件</el-button>
                        <el-input v-model="domain.batch" type="textarea" :rows="5" placeholder="请输入域名，按换行分割"></el-input>
                        <div class="flex-between mt-5px">
                            <div></div>
                            <el-button type="primary" @click="batchQuery">开始批量查询</el-button>
                        </div>
                    </el-popover>
                    <el-tooltip content="复制内容并跳转到数据处理">
                        <el-button :icon="Promotion" @click="CopyAndJump" />
                    </el-tooltip>
                </div>
            </div>
            <el-input v-model="domain.cdn" type="textarea" :rows="13" resize="none"
                class="mt-10px"></el-input>

        </el-card>
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/ip138.ico"></el-icon><span class="title">子域名:</span></el-text>
            <el-input v-model="domain.subdomain" type="textarea" :rows="13" resize="none"
                class="mt-10px"></el-input>
        </el-card>
        <el-card shadow="never" class="grid-item">
            <el-text><el-icon><img src="/ip138.ico"></el-icon><span class="title">历史解析:</span></el-text>
            <el-input v-model="domain.history" type="textarea" :rows="13" resize="none"
                class="mt-10px"></el-input>
        </el-card>
    </div>
</template>

<style scoped>
.title {
    margin-left: 5px;
    font-weight: bold;
    font-size: 18px;
}

.grid-container {
    margin-top: 10px;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    gap: 10px;
    height: 100%;
    /* 高度可以根据需要调整 */
}

.grid-item {
    justify-content: center;
    align-items: center;
    border-radius: 5px;
}
</style>