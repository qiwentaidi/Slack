<script setup lang="ts">
import { ArrowUpBold, ArrowDownBold, Connection, Delete, DocumentCopy } from '@element-plus/icons-vue';
import { ExportTXT } from '@/export'
import { reactive, onMounted, ref } from 'vue';
import { Socks5Conn } from 'wailsjs/go/services/App'
import { DeleteAgentPool, InsertAgentPool, SelectAllAgentPool, DeleteAllAgentPool } from 'wailsjs/go/services/Database';
import { ElMessage } from 'element-plus';
import exportIcon from '@/assets/icon/doucment-export.svg';
import CustomTextarea from '@/components/CustomTextarea.vue';
import Note from '@/components/Note.vue';
import { AgentPoolOptions } from '@/stores/options';
import { Copy, ProcessTextAreaInput } from '@/util';

onMounted(async () => {
    let hosts = await SelectAllAgentPool()
    if (Array.isArray(hosts)) form.pool = hosts
});

// 定义代理验证结果接口
interface ProxyResult {
    host: string;
    successCount: number;
    totalAttempts: number;
    isAlive: boolean;
}

// 定义验证状态
const validationStatus = reactive({
    isRunning: false,
    currentHost: '',
    totalHosts: 0,
    completedHosts: 0
});

const form = reactive({
    target: '',
    aliveURL: 'http://www.baidu.com',
    aliveTimes: 3,
    socksLogger: '',
    percentage: 0,
    activeNames: 0,
    pool: [] as string[],
})


// 代理验证函数
async function validateProxy(host: string, attempts: number): Promise<ProxyResult> {
    const [ip, port] = host.split(':');
    const promises = Array(attempts).fill(null).map(() => 
        Socks5Conn(ip, Number(port), 5, "", "", form.aliveURL)
    );
    
    const results = await Promise.all(promises);
    const successCount = results.filter(Boolean).length;
    
    return {
        host,
        successCount,
        totalAttempts: attempts,
        isAlive: successCount > 0
    };
}
        
// 批量验证函数
async function CheckSock5Unauth() {
    if (validationStatus.isRunning) {
        ElMessage.warning('验证任务正在运行中...');
        return;
    }
    showForm.value = false
    validationStatus.isRunning = true;
    form.socksLogger = '[*] 验证任务已开始运行，请稍等...\n';
    const lines = ProcessTextAreaInput(form.target);
    validationStatus.totalHosts = lines.length;
    validationStatus.completedHosts = 0;
    try {
        // 验证新输入的代理
        for (const host of lines) {
            if (!host.includes(':')) continue;

            validationStatus.currentHost = host;
            const result = await validateProxy(host, Number(form.aliveTimes));

            validationStatus.completedHosts++;
            form.percentage = Number(((validationStatus.completedHosts / validationStatus.totalHosts) * 100).toFixed(2))

            if (result.isAlive) {
                form.socksLogger += `[+] ${host} is alive! (${result.successCount}/${result.totalAttempts} successful attempts) | (${form.percentage})\n`;
                await insertPool(host)
            } else {
                form.socksLogger += `[-] ${host} is dead (${result.successCount}/${result.totalAttempts} successful attempts) | (${form.percentage})\n`;
            }
        }

        form.socksLogger += `[*] 验证完成`
    } catch (error) {
        form.socksLogger += `[!] Error during validation: ${error}\n`;
    } finally {
        validationStatus.isRunning = false;
    }
}

function exportHosts() {
    if (!form.pool) return
    ExportTXT("socks_unauth_asset", form.pool)
}

async function TestConnection(host: string) {
    const [ip, port] = host.split(':');
    let result = await Socks5Conn(ip, Number(port), 5, "", "", form.aliveURL)
    if (result) {
        ElMessage.success('This proxy is reachable')
    } else {
        ElMessage.error('Oops, this proxy is unreachable.')
    }
}

async function deleltePool(host: string) {
    let result = await DeleteAgentPool(host)
    if (result) {
        form.pool = form.pool.filter(item => item !== host)
        ElMessage.success('删除成功')
    } else {
        ElMessage.error('删除失败')
    }
}

async function deleteAllPool() {
    let result = await DeleteAllAgentPool()
    if (result) {
        form.pool = []
        ElMessage.success('删除成功')
    } else {
        ElMessage.error('删除失败')
    }
}

async function insertPool(host: string) {
    // 如果数据库中已存在该代理则删除
    if (form.pool.includes(host)) return

    let result = await InsertAgentPool(host);
    if (result) {
        form.pool.push(host)
        return
    }
    form.socksLogger += `[!] insert ${host} in the database failed`
}

const showForm = ref(true);

function toggleFormVisibility() {
    showForm.value = !showForm.value;
}

const quakeSearch = 'service:socks5 AND country: "CN" AND response:"No authentication"'
const hunterSearch = 'protocol=="socks5"&&protocol.banner="No authentication"&&ip.country="CN"'
const fofaSearch = 'protocol=="socks5" && country="CN" && banner="Method:No Authentication"'
</script>

<template>
    <el-divider>
        <el-button round :icon="showForm ? ArrowUpBold : ArrowDownBold" @click="toggleFormVisibility">
            {{ showForm ? '隐藏参数' : '展开参数' }}
        </el-button>
    </el-divider>
    <el-collapse-transition>
        <div style="display: flex; gap: 10px; margin-bottom: 10px;" v-show="showForm">
            <el-form :model="form" label-width="auto" style="width: 50%;">
                <el-form-item label="目标:">
                    <CustomTextarea v-model="form.target" :rows="5" placeholder="192.168.1.1:1111, 目标以换行分割" />
                </el-form-item>
                <el-form-item label="代理验证:">
                    <el-input v-model="form.aliveURL"></el-input>
                </el-form-item>
                <el-form-item label="校验次数:">
                    <el-input v-model="form.aliveTimes"></el-input>
                    <span class="form-item-tips">为了保证代理的可用性, 需要通过代理对存活目标访问成功次数</span>
                </el-form-item>
                <el-form-item class="align-right">
                    <el-button v-if="!validationStatus.isRunning" type="primary" @click="CheckSock5Unauth">开始任务</el-button>
                    <el-button v-else loading>程序正在运行中</el-button>
                </el-form-item>
            </el-form>
            <Note style="width: 50%;">
                本模块仅提供SOCKS5代理批量验证保存功能, 如需测绘可以通过空间引擎自行搜索<br /><br />
                FOFA:<br />
                <el-tag type="info" @click="Copy(fofaSearch)">{{ fofaSearch }}</el-tag><br />
                Quake: <br />
                <el-tag type="info" @click="Copy(quakeSearch)">{{ quakeSearch }}</el-tag><br />
                Hunter:<br />
                <el-tag type="info" @click="Copy(hunterSearch)">{{ hunterSearch }}</el-tag>
            </Note>
        </div>
    </el-collapse-transition>
    <el-card shadow="never" style="width: 100%;">
        <template #header>
            <div class="card-header">
                <el-segmented v-model="form.activeNames" :options="AgentPoolOptions">
                    <template #default="{ item }">
                        <el-space :size="3">
                            <el-icon>
                                <component :is="item.icon" />
                            </el-icon>
                            <div>{{ item.label }}</div>
                        </el-space>
                    </template>
                </el-segmented>
                <el-space v-show="form.activeNames == 1">
                    <el-button :icon="DocumentCopy" plain type="primary" @click="Copy(form.pool.join('\n'))">复制全部</el-button>
                    <el-button :icon="exportIcon" plain type="primary" @click="exportHosts">导出目标</el-button>
                    <el-button :icon="Delete" plain type="danger" @click="deleteAllPool">删除全部</el-button>
                </el-space>
            </div>
        </template>
        <pre class="pretty-response" style="margin-top: 0; margin-bottom: 0;" v-show="form.activeNames == 0"><code>{{
            form.socksLogger
        }}</code></pre>
        <el-table :data="form.pool" border style="height: 100vh;" v-show="form.activeNames == 1">
            <el-table-column type="index" width="50" label="#" align="center" />
            <el-table-column label="主机地址">
                <template #default="scope">
                    {{ scope.row }}
                </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" width="190" align="center">
                <template #default="scope">
                    <el-space>
                        <el-button :icon="Connection" size="small"
                        @click.prevent="TestConnection(scope.row)">测试连接</el-button>
                        <el-button :icon="Delete" type="danger" plain size="small" @click.prevent="deleltePool(scope.row)">删除</el-button>
                    </el-space>
                </template>
            </el-table-column>
        </el-table>
    </el-card>
</template>
