<script lang="ts" setup>
import { QuestionFilled, FolderOpened } from '@element-plus/icons-vue'
import CustomTextarea from '@/components/CustomTextarea.vue'
import { webscanOptions, WebsiteInputTips } from '@/stores/options'
import { form, param, config, spaceEngineConfig } from '@/composables/useWebscanState'
import { DirectoryDialog } from 'wailsjs/go/services/File'
import { SaveConfig } from '@/config'
import global from "@/stores"
import saveIcon from '@/assets/icon/save.svg'

async function selectFolder() {
    global.webscan.append_pocfile = await DirectoryDialog()
}

defineEmits<{
    startScan: []
}>()
</script>

<template>
    <el-drawer v-model="form.newWebscanDrawer" size="50%">
        <template #header>
            <span class="drawer-title">新建网站扫描</span>
            <el-button link @click="spaceEngineConfig.fofaDialog = true">
                <template #icon>
                    <img src="/app/fofa.png" style="width: 16px; height: 16px;">
                </template>
                FOFA
            </el-button>
            <el-divider direction="vertical" />
            <el-button link @click="spaceEngineConfig.hunterDialog = true">
                <template #icon>
                    <img src="/app/hunter.ico" style="width: 16px; height: 16px;">
                </template>
                Hunter
            </el-button>
        </template>
        <template #footer>
            <div class="position-center">
                <el-button type="primary" @click="$emit('startScan')" style="bottom: 10px; position: absolute;">开始任务</el-button>
            </div>
        </template>
        <el-form label-width="auto">
            <el-form-item label="任务名称:">
                <el-input v-model="form.taskName" />
            </el-form-item>
            <el-form-item label="目标地址:">
                <CustomTextarea v-model="form.input" :rows="6" :placeholder="WebsiteInputTips"></CustomTextarea>
            </el-form-item>
            <el-form-item label="追加POC:">
                <el-input v-model="global.webscan.append_pocfile">
                    <template #suffix>
                        <el-tooltip content="选择POC文件夹">
                            <el-button :icon="FolderOpened" link @click="selectFolder()" />
                        </el-tooltip>
                        <el-divider direction="vertical" />
                        <el-tooltip content="保存">
                            <el-button :icon="saveIcon" link @click="SaveConfig" />
                        </el-tooltip>
                    </template>
                </el-input>
                <span class="form-item-tips">额外追加一个文件夹下的所有YAML POC文件，便于管理自添加POC</span>
            </el-form-item>
            <el-form-item>
                <template #label>模式:
                    <el-tooltip>
                        <template #content>
                            1、指纹扫描: 只进行简单的指纹探测，不会探测敏感目录<br />
                            2、全指纹扫描: 会在指纹扫描基础上增加主动敏感目录探测，例如Nacos、报错页面信息判断指纹等<br />
                            3、指纹漏洞扫描: 指纹+主动敏感目录探测，扫描完成后扫描指纹对应POC，如果网站未识别到指纹会扫描全漏洞<br />
                            4、专项扫描: 只扫描简单指纹和已选择的漏洞，如果不指定漏洞会扫指纹漏洞(忽略主动探测功能)
                        </template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-segmented v-model="config.webscanOption" :options="webscanOptions" class="w-full">
                    <template #default="{ item }">
                        <el-space :size="3">
                            <el-icon :size="18">
                                <component :is="item.icon" :key="item.value" />
                            </el-icon>
                            <div>{{ item.label }}</div>
                        </el-space>
                    </template>
                </el-segmented>
            </el-form-item>
            <el-form-item label="请求头:">
                <el-input v-model="config.customHeaders" :rows="3" type="textarea"
                    :placeholder="$t('tips.customHeaders')"></el-input>
            </el-form-item>
            <div v-if="config.webscanOption == 3">
                <el-form-item label="指定指纹:">
                    <el-select-v2 v-model="config.customTags" :options="param.allFingerprint" filterable multiple
                        clearable />
                    <span class="form-item-tips">类似Jeecg-Boot指纹漏洞都在API路径下, 可通过填写API地址并指定指纹来进行扫描(POC 需要进行适配)</span>
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 3">
                <el-form-item label="指定漏洞:">
                    <el-select-v2 v-model="config.customTemplate" :options="param.allTemplate" filterable multiple
                        clearable />
                </el-form-item>
            </div>
            <div v-if="config.webscanOption == 2">
                <el-form-item label="Log4j2:">
                    <el-switch v-model="config.generateLog4j2" class="w-full" />
                    <span class="form-item-tips">开启后会将所有目标添加Generate-Log4j2指纹</span>
                </el-form-item>
            </div>
            <el-form-item label="高级配置:">
                <el-tooltip content="启用后主动指纹只拼接根路径，否则会拼接输入的完整URL">
                    <el-checkbox label="根路径扫描" v-model="config.rootPathScan" />
                </el-tooltip>
                <el-checkbox label="无指纹目标跳过漏扫" v-model="config.skipNucleiWithoutTags" />
                <el-checkbox label="网站截图" v-model="config.screenhost" />
            </el-form-item>
        </el-form>
    </el-drawer>
</template>

