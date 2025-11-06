<script lang="ts" setup>
import { form, spaceEngineConfig, shodanVisible, shodanIp, shodanPercentage, shodanThread, shodanRunningstatus } from '@/composables/useWebscanState'
import { uncover } from '@/composables/useSpaceEngine'
import { Callgologger } from 'wailsjs/go/services/App'
import { SpaceGetPort } from 'wailsjs/go/services/App'
import { ProcessTextAreaInput } from '@/util'
import { validateIp, isPrivateIP } from '@/stores/validate'
import { IPParse } from 'wailsjs/go/core/Tools'
import async from 'async'

async function handleShodan() {
    if (!validateIp(shodanIp.value)) {
        return
    }
    const lines = ProcessTextAreaInput(shodanIp.value)
    for (const line of lines) {
        if (isPrivateIP(line)) {
            return
        }
    }

    let id = 0
    form.input = ""
    shodanRunningstatus.value = true
    let ips = await IPParse(lines)
    async.eachLimit(ips, shodanThread.value, async (ip: string, callback: () => void) => {
        if (!shodanRunningstatus.value) {
            return
        }
        let ports = await SpaceGetPort(ip)
        id++
        if (ports == null) {
            return
        }
        Callgologger("info", "[shodan] " + ip + " port: " + ports.join())
        shodanPercentage.value = Number(((id / ips.length) * 100).toFixed(2));
        for (const port of ports) {
            form.input += ip + ":" + port.toString() + "\n"
        }
    }, (err: any) => {
        shodanRunningstatus.value = false
        shodanVisible.value = false
    })
}
</script>

<template>
    <!-- FOFA Dialog -->
    <el-dialog v-model="spaceEngineConfig.fofaDialog">
        <template #header>
            <span class="drawer-title"><img src="/app/fofa.png">导入FOFA目标</span>
        </template>
        <el-form label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="spaceEngineConfig.fofaQuery"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-input-number v-model="spaceEngineConfig.fofaPageSize" :min="1" :max="10000" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.fofa">导入</el-button>
        </template>
    </el-dialog>

    <!-- Hunter Dialog -->
    <el-dialog v-model="spaceEngineConfig.hunterDialog">
        <template #header>
            <span class="drawer-title"><img src="/app/hunter.ico">导入鹰图目标，由于API查询限制，大数据推荐使用官网进行数据导出</span>
        </template>
        <el-form label-width="auto">
            <el-form-item label="查询条件">
                <el-input v-model="spaceEngineConfig.hunterQuery"></el-input>
            </el-form-item>
            <el-form-item label="导入数量">
                <el-select v-model="spaceEngineConfig.hunterPageSize" style="width: 150px;">
                    <el-option v-for='item in ["10", "20", "50", "100"]' :key="item" :label="item" :value="item" />
                </el-select>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="uncover.hunter">导入</el-button>
        </template>
    </el-dialog>

    <!-- Shodan Dialog -->
    <el-dialog v-model="shodanVisible" width="500">
        <template #header>
            <span class="drawer-title"><img src="/shodan.png">从Shodan拉取资产端口开放情况</span>
        </template>
        <el-form label-position="top">
            <el-form-item label="扫描线程:">
                <el-input-number v-model="shodanThread" :min="1" :max="3" />
            </el-form-item>
            <el-form-item>
                <el-input v-model="shodanIp" type="textarea" :rows="6" placeholder="输入规则与IP模式一致，但不支持域名和排除"></el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <div class="flex-between">
                <el-progress :text-inside="true" :stroke-width="18" :percentage="shodanPercentage"
                    style="width: 300px;" />
                <el-button size="small" type="danger" @click="shodanRunningstatus = false"
                    v-if="shodanRunningstatus">终止探测</el-button>
                <el-button size="small" type="primary" @click="handleShodan" v-else>
                    开始收集
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

