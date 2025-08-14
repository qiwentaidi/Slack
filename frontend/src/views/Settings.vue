<template>
    <el-container class="h-full">
        <el-aside width="200px">
            <el-menu default-active="0">
                <el-menu-item v-for="(item, index) in setupOptions" :index="index.toString()" @click="selectItem">
                    <el-icon>
                        <component :is="item.icon" />
                    </el-icon>
                    <template #title><span>{{ $t(item.name) }}</span></template>
                </el-menu-item>
            </el-menu>
        </el-aside>

        <el-main>
            <el-form :model="global.webscan" label-width="auto" v-show="currentDisplay == '0'">
                <h3>{{ $t(setupOptions[0].name) }}</h3>
                <el-form-item :label="$t('setting.webscan_thread')">
                    <el-input-number v-model="global.webscan.web_thread" :min="1" :max="200" />
                </el-form-item>
                <el-form-item :label="$t('setting.crack_thread')">
                    <el-input-number v-model="global.webscan.crack_thread" :min="1" :max="200" />
                </el-form-item>
                <el-form-item :label="$t('setting.portscan_thread')">
                    <el-input-number v-model="global.webscan.port_thread" :min="1" :max="10000" />
                </el-form-item>
                <el-form-item :label="$t('setting.portscan_timeout')">
                    <el-input-number v-model="global.webscan.port_timeout" :min="1" :max="20" />
                </el-form-item>
                <el-form-item :label="$t('setting.survival')">
                    <el-select v-model="global.webscan.default_alive_module">
                        <el-option v-for="item in aliveGroupOptions" :key="item.value" :value="item.value">
                            <span class="float-left">{{ item.value }}</span>
                            <span class="float-right">
                                {{ item.description }}
                            </span>
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('setting.network_list')">
                    <el-select v-model="global.webscan.default_network">
                        <el-option v-for="item in global.temp.NetworkCardList" :value="item.IP">
                            <span v-if="item.Name == ''">{{ item.IP }}</span>
                            <span v-else>{{ item.Name + '(' + item.IP + ')' }}</span>
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-button type="primary" @click="SaveConfig" class="float-right">{{ $t('setting.save') }}</el-button>
            </el-form>
            <el-form :model="global.proxy" label-width="auto" v-show="currentDisplay == '1'">
                <h3>{{ $t(setupOptions[1].name) }}</h3>
                <el-form-item :label="$t('setting.enable')">
                    <el-switch v-model="global.proxy.enabled" />
                    <el-button type="primary" size="small" @click="TestProxyWithNotify" class="ml-10px"
                        v-if="global.proxy.enabled">{{ $t('setting.test_agent') }}</el-button>
                </el-form-item>
                <el-form-item :label="$t('setting.mode')">
                    <el-select v-model="global.proxy.mode">
                        <el-option label="HTTP" value="HTTP" />
                        <el-option label="SOCK5" value="SOCK5" />
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('setting.address')">
                    <el-input v-model="global.proxy.address" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('setting.port')">
                    <el-input-number v-model.number="global.proxy.port" :min="1" :max="65535"></el-input-number>
                </el-form-item>
                <el-form-item :label="$t('setting.username')">
                    <el-input v-model="global.proxy.username" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('setting.password')">
                    <el-input v-model="global.proxy.password" type="password" show-password></el-input>
                </el-form-item>
                <el-button type="primary" @click="SaveConfig" class="float-right">{{ $t('setting.save') }}</el-button>
            </el-form>
            <el-form :model="global.space" label-width="auto" v-show="currentDisplay == '2'">
                <h3>{{ $t(setupOptions[2].name) }}<el-divider direction="vertical" />{{ $t('setting.identifier') }}</h3>
                <el-form-item label="FOFA">
                    <el-input v-model="global.space.fofaapi" placeholder="API address" clearable />
                    <el-input v-model="global.space.fofaemail" placeholder="Email" clearable class="mt-5px" />
                    <el-input v-model="global.space.fofakey" placeholder="API key" clearable class="mt-5px" />
                </el-form-item>
                <el-form-item label="Hunter">
                    <el-input v-model="global.space.hunterapi" placeholder="API address" clearable />
                    <el-input v-model="global.space.hunterkey" placeholder="API key" clearable class="mt-5px" />
                </el-form-item>
                <el-form-item label="Quake">
                    <el-input v-model="global.space.quakekey" placeholder="API key" clearable />
                </el-form-item>
                <el-form-item label="ChaosⒹ">
                    <el-input v-model="global.space.chaos">
                        <template #suffix>
                            <el-button link type="primary" :icon="User" @click="BrowserOpenURL(chaosURL)">{{
                                $t('setting.register') }}</el-button>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item label="BevigilⒹ">
                    <el-input v-model="global.space.bevigil" placeholder="API key">
                        <template #suffix>
                            <el-button link type="primary" :icon="User" @click="BrowserOpenURL(bevigilURL)">{{
                                $t('setting.register') }}</el-button>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item label="SecuritytrailsⒹ">
                    <el-input v-model="global.space.securitytrails" placeholder="API key">
                        <template #suffix>
                            <el-button link type="primary" :icon="User" @click="BrowserOpenURL(securitytrailsURL)">{{
                                $t('setting.register') }}</el-button>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item label="ZoomeyeⒹ">
                    <el-input v-model="global.space.zoomeye" placeholder="API key" />
                </el-form-item>
                <el-form-item label="GithubⒹ">
                    <el-input v-model="global.space.github"
                        placeholder="Settings -> Developer settings -> Presonal access tokens" />
                </el-form-item>
                <el-button type="primary" @click="SaveConfig" class="float-right">{{ $t('setting.save') }}</el-button>
            </el-form>
            <el-form :model="global.Theme" label-width="auto" v-show="currentDisplay == '3'">
                <h3>{{ $t(setupOptions[3].name) }}</h3>
                <el-form-item :label="$t('setting.poc_source')">
                    <el-input v-model="global.update.pocSource"></el-input>
                    <span class="form-item-tips">{{ $t('setting.poc_tips') }}</span>
                </el-form-item>
                <el-button type="primary" @click="SaveConfig" class="float-right">{{ $t('setting.save') }}</el-button>
            </el-form>
            <div v-show="currentDisplay == '4'">
                <h3>{{ $t(setupOptions[4].name) }}<el-divider direction="vertical" />{{ $t('setting.password_tips') }}
                </h3>
                <el-table :data="crackDict.usernames" stripe class="w-full">
                    <el-table-column prop="name" label="Protocol" />
                    <el-table-column label="Operate" width="250" align="center">
                        <template #default="scope">
                            <el-button type="primary" link :icon="Edit"
                                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/username/' + scope.row.name + '.txt'; ReadDict(ctrl.currentPath)">{{
                                    $t('setting.username') }}</el-button>
                            <el-button type="primary" link :icon="Edit"
                                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/password/password.txt'; ReadDict(ctrl.currentPath)">{{
                                    $t('setting.password') }}</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="position-center" v-show="currentDisplay == '5'">
                <about></about>
            </div>
        </el-main>
    </el-container>
    <el-drawer v-model="ctrl.innerDrawer" title="字典管理" :append-to-body="true">
        <el-input type="textarea" :rows="20" v-model="ctrl.currentDic"></el-input>
        <el-button type="primary" class="mt-10px float-right"
            @click="SaveFile(ctrl.currentPath)">保存</el-button>
    </el-drawer>
</template>

<script lang="ts" setup>
import global from "@/stores"
import { ElMessage, MenuItemRegistered } from 'element-plus';
import { TestProxyWithNotify } from "@/util";
import { Edit, User } from '@element-plus/icons-vue';
import { reactive, ref } from "vue";
import { ReadFile, WriteFile } from "wailsjs/go/services/File";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import { SaveConfig } from "@/config";
import { aliveGroupOptions, crackDict, setupOptions } from "@/stores/options";

const bevigilURL = "https://bevigil.com/osint-api"
const chaosURL = "https://cloud.projectdiscovery.io/"
const securitytrailsURL = "https://securitytrails.com/"

const ctrl = reactive({
    innerDrawer: false,
    currentDic: '',
    currentPath: '',
})

async function ReadDict(path: string) {
    let file = await ReadFile(global.PATH.homedir + global.PATH.PortBurstPath + path)
    ctrl.currentDic = file.Content
}

async function SaveFile(path: string) {
    let isSuccess = await WriteFile('txt', global.PATH.homedir + global.PATH.PortBurstPath + path, ctrl.currentDic)
    isSuccess ? ElMessage.success('保存成功!') : ElMessage.error('保存失败!')
}

const currentDisplay = ref('0')

function selectItem(item: MenuItemRegistered) {
    currentDisplay.value = item.index
}


</script>

<style scoped>
.el-main {
    padding-top: 0px;
    padding-right: 0px;
}

.el-menu {
    border-right: 0px;
}

.el-menu-item.is-active {
    background-color: var(--el-menu-active-color);
    color: var(--el-menu-hover-bg-color);
    border-radius: 5px;
}

.el-menu-item:hover {
    border-radius: 5px;
}
</style>