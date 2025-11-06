<script lang="ts" setup>
import CustomTextarea from '@/components/CustomTextarea.vue'
import { portGroupOptions, HostInputTips } from '@/stores/options'
import { form, param, config, updatePorts } from '@/composables/useWebscanState'
import { shodanVisible } from '@/composables/useWebscanState'

defineEmits<{
    startScan: []
}>()
</script>

<template>
    <el-drawer v-model="form.newHostscanDrawer" size="50%">
        <template #header>
            <span class="drawer-title">新建主机扫描</span>
            <el-button link @click="shodanVisible = true">
                <template #icon>
                    <img src="/shodan.png" style="width: 14px; height: 14px;">
                </template>
                Shodan
            </el-button>
        </template>
        <template #footer>
            <div class="position-center">
                <el-button type="primary" @click="$emit('startScan')" style="bottom: 10px;">开始任务</el-button>
            </div>
        </template>
        <el-form label-width="auto">
            <el-form-item label="任务名称:">
                <el-input v-model="form.taskName" />
            </el-form-item>
            <el-form-item label="目标地址:">
                <CustomTextarea v-model="form.input" :rows="11" :placeholder="HostInputTips"></CustomTextarea>
            </el-form-item>
            <el-form-item label="端口:">
                <el-select v-model="param.portGroup" @change="(val: number) => updatePorts(val)">
                    <el-option v-for="(item, index) in portGroupOptions" :key="index" :label="item.text" :value="index" />
                </el-select>
                <el-input v-model="form.portlist" type="textarea" :rows="4" resize="none" class="mt-5px"></el-input>
            </el-form-item>
            <el-form-item label="其他:">
                <el-tooltip content="排除9100端口">
                    <el-checkbox v-model="config.excludePrintPorts">不扫描网络打印机</el-checkbox>
                </el-tooltip>
            </el-form-item>
            <el-form-item label="漏洞扫描:">
                <el-switch v-model="config.vulscan" class="w-full" />
                <span class="form-item-tips">开启后端口扫描结束会进一步提取WEB应用指纹, 并调用指纹漏洞扫描模式扫描漏洞</span>
            </el-form-item>
            <el-form-item label="高级配置:" v-show="config.vulscan">
                <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                    <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                </el-tooltip>
                <el-checkbox label="无指纹目标跳过漏扫" v-model="config.skipNucleiWithoutTags" />
                <el-checkbox label="网站截图" v-model="config.screenhost" />
            </el-form-item>
            <el-form-item label="口令暴破:" v-show="config.vulscan">
                <el-switch v-model="config.crack" class="w-full" />
                <span class="form-item-tips" v-show="config.crack">默认字典可通过 设置->
                    字典管理处修改, 由于RDP暴破可能存在闪退, 暂时不支持暴破</span>
            </el-form-item>
            <el-form-item label="用户字典:" v-show="config.crack">
                <CustomTextarea v-model="param.username" :rows="5"
                    @input="param.builtInUsername = param.username.length === 0"></CustomTextarea>
                <el-checkbox v-model="param.builtInUsername"
                    :disabled="param.username.length == 0">使用默认用户字典</el-checkbox>
            </el-form-item>
            <el-form-item label="密码字典:" v-show="config.crack">
                <CustomTextarea v-model="param.password" :rows="5"
                    @input="param.builtInPassword = param.password.length === 0"></CustomTextarea>
                <el-checkbox v-model="param.builtInPassword"
                    :disabled="param.password.length == 0">使用默认密码字典</el-checkbox>
            </el-form-item>
        </el-form>
    </el-drawer>
</template>

