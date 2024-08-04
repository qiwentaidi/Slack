<script lang="ts" setup>
import { Close, Minus, Search, Plus, InfoFilled, Back, Right } from '@element-plus/icons-vue';
import { Quit, WindowMinimise, WindowToggleMaximise } from "wailsjs/runtime/runtime";
import { IsMacOS } from "wailsjs/go/main/File";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import global from "@/global";
import { useRoute } from "vue-router";
import { ref, computed } from "vue";
import about from "./About.vue";
import updateUI from "./Update.vue";
import aboutIcon from "@/assets/icon/about.svg"
import onlineIcon from "@/assets/icon/online.svg"
import runnerIcon from "@/assets/icon/apprunner.svg"
import maxmizeIcon from "@/assets/icon/maximize.svg"
import reductionIcon from "@/assets/icon/reduction.svg"

const isMax = ref(false)
const isMacOS = ref(false)
const onlineDrawer = ref(false)
const localDrawer = ref(false)
const showLogger = ref(false)
const route = useRoute();
const aboutDialog = ref(false)
const updateDialog = ref(false)

window.addEventListener('resize', () => {
    if (screen.availWidth <= window.innerWidth && screen.availHeight <= window.innerHeight) {
        isMax.value = true
    } else {
        isMax.value = false
    }
});

IsMacOS().then(res => {
    isMacOS.value = res
})
function setTitle(path: string) {
    switch (path) {
        case "/":
            return "Home";
        case "/settings":
            return "Settings";
        default:
            return path.split('/').slice(-1)[0];
    }
}

const rightStyle = computed(() => {
    return isMacOS.value ? { marginRight: '3.5px' } : {};
})

const leftStyle = computed(() => {
    return !isMacOS.value ? { marginLeft: '3.5px' } : {};
})

const macStyle = computed(() => {
    return isMacOS.value ? { marginLeft: '5.5%' } : {};
})

const titlebarStyle = computed(() => {
    return global.Theme.value ? {
        backgroundColor: '#333333',
        borderBottom: '1px solid #3B3B3B'
    } : {
        backgroundColor: '#F9F9F9',
        borderBottom: '1px solid #E6E6E6'
    };
})

const svgReverseFill = computed(() => {
    return global.Theme.value ? '#333333' : '#F9F9F9'
})

const searchFilter = ref("");
const filteredOptions = computed(() => {
    if (!searchFilter.value) {
        return global.onlineOptions;
    }
    return global.onlineOptions.map((group) => ({
        ...group,
        value: group.value.filter((item) =>
            item.name.toLowerCase().includes(searchFilter.value.toLowerCase())
        ),
    }));
});

const localNavigationRef = ref();
function addGroup() {
    if (localNavigationRef.value) {
        (localNavigationRef.value as any).addGroup();
    }
}
</script>

<template>
    <div class="titlebar" :style="titlebarStyle">
        <div :style="macStyle">
            <el-divider direction="vertical" v-if="isMacOS" />
            <el-button-group :style="leftStyle">
                <el-tooltip content="返回">
                    <el-button text class="custom-button" @click="$router.back()">
                        <template #icon>
                            <el-icon size="16">
                                <Back />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="前进">
                    <el-button text class="custom-button" @click="$router.forward()">
                        <template #icon>
                            <el-icon size="16">
                                <Right />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
            </el-button-group>
        </div>
        <div class="unoccupied" @dblclick="WindowToggleMaximise">
            <span class="title">{{ setTitle(route.path) }}</span>
        </div>
        <div style="display: flex">
            <el-button-group :style="rightStyle">
                <el-tooltip :content="$t('aside.about')">
                    <el-button class="custom-button" text @click="aboutDialog = true">
                        <template #icon>
                            <el-icon size="16">
                               <aboutIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="运行日志">
                    <el-button class="custom-button" text @click="showLogger = true">
                        <template #icon>
                            <el-icon size="16">
                                <svg t="1722163018715" class="icon" viewBox="0 0 1024 1024" version="1.1"
                                    xmlns="http://www.w3.org/2000/svg" p-id="12142" width="200" height="200">
                                    <path
                                        d="M956.8 128H67.2A67.2 67.2 0 0 0 0 195.2v633.6A67.2 67.2 0 0 0 67.2 896h889.6a67.2 67.2 0 0 0 67.2-67.2V195.2A67.2 67.2 0 0 0 956.8 128z m3.2 704H64V192h896v640z"
                                        fill="" p-id="12143"></path>
                                    <path d="M96 224v576h832V224H96z" :fill="svgReverseFill" p-id="12144"></path>
                                    <path
                                        d="M294.624 694.624l-45.248-45.248L386.752 512l-137.376-137.376 45.248-45.248L477.248 512z"
                                        fill="" p-id="12144"></path>
                                    <path d="M768 704h-256v-64h256v64z" fill="" p-id="12144"></path>
                                </svg>
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="应用启动器(本功能会在wails v3进行重做)">
                    <el-button text class="custom-button" @click="localDrawer = true" v-show="!isMacOS">
                        <template #icon>
                            <el-icon size="16">
                               <runnerIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="在线导航">
                    <el-button text class="custom-button" @click="onlineDrawer = true">
                        <template #icon>
                            <el-icon size="16">
                                <onlineIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
            </el-button-group>
            <div v-if="!isMacOS">
                <el-divider direction="vertical" />
                <el-button-group>
                    <el-button text @click="WindowMinimise">
                        <template #icon>
                            <el-icon size="16">
                                <Minus />
                            </el-icon>
                        </template>
                    </el-button>
                    <el-button text @click="WindowToggleMaximise">
                        <template #icon>
                            <el-icon size="16" v-show="!isMax">
                                <maxmizeIcon />
                            </el-icon>
                            <el-icon size="16" v-show="isMax">
                                <reductionIcon />
                            </el-icon>
                        </template>
                    </el-button>
                    <el-button text class="close" @click="Quit">
                        <template #icon>
                            <el-icon size="16">
                                <Close />
                            </el-icon>
                        </template>
                    </el-button>
                </el-button-group>
            </div>
        </div>
    </div>
    <!-- 在线导航 -->
    <el-drawer v-model="onlineDrawer" direction="rtl" size="35%">
        <template #header>
            <el-input :suffix-icon="Search" placeholder="根据名称过滤" v-model="searchFilter" @input="filteredOptions" />
        </template>
        <div v-for="groups in filteredOptions" :key="groups.label" style="margin-bottom: 5px;">
            <el-card v-if="groups.value.length > 0">
                <div style="margin-bottom: 10px">
                    <span style="font-weight: bold">{{ $t(groups.label) }}</span>
                </div>
                <div style="display: grid; gap: 10px;">
                    <div v-for="item in groups.value" :key="item.name">
                        <el-tooltip :content="item.url" placement="top" :show-after="1000">
                            <div class="nav-item" @click="BrowserOpenURL(item.url)">
                                <img :src="item.icon"><span class="nav-text">{{ item.name }}</span>
                            </div>
                        </el-tooltip>
                    </div>
                </div>
            </el-card>
        </div>
    </el-drawer>
    <!-- 本地导航 -->
    <el-drawer v-model="localDrawer" direction="rtl" size="70%">
        <template #header>
            <el-tooltip>
                <template #content>
                    1、添加组后会出现卡片组，名称不支持重复<br />
                    2、卡片新增后可以从桌面或者文件夹拖入元素或者右上角添加元素<br />
                    3、右键元素可以编辑已有信息、打开文件夹、删除<br />
                    4、对于CMD应用，Windows会进入当前应用路径的CMD窗口<br />
                </template>
                <el-button bg text :icon="InfoFilled">使用须知</el-button>
            </el-tooltip>
            <el-button-group>
                <el-button @click="addGroup">
                    <template #icon>
                        <img src="../assets/icon/group.svg">
                    </template>
                    添加组</el-button>
                <el-button :icon="Plus" @click="global.temp.localAddItem = true">添加元素</el-button>
            </el-button-group>
        </template>
        <LocalNavigation ref="localNavigationRef" />
    </el-drawer>
    <!-- running logs -->
    <el-drawer v-model="showLogger" title="运行日志" direction="rtl" size="50%">
        <div class="log-textarea" v-html="global.Logger.value"></div>
    </el-drawer>
    <!-- about -->
    <el-dialog v-model="aboutDialog" width="36%" center>
        <about></about>
    </el-dialog>
    <!-- update -->
    <el-dialog v-model="updateDialog" title="更新通知" width="40%">
        <updateUI></updateUI>
    </el-dialog>
</template>

<style scoped>
.titlebar {
    display: flex;
    width: 100%;
    height: 35px;

    .el-button {
        height: 35px;
        border-radius: 0;
    }

    .custom-button {
        margin-top: 3.5px;
        margin-bottom: 3.5px;
        height: 28px;
        width: 35px;
        border-radius: 10px;
    }

    .custom-button img {
        width: 16px;
    }
}

.unoccupied {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-grow: 1;
    --wails-draggable: drag;
}

.title {
    -webkit-user-select: none;
    /* Safari */
    -moz-user-select: none;
    /* Firefox */
    -ms-user-select: none;
    /* IE10+/Edge */
    user-select: none;
    /* Standard syntax */
    margin-right: 5%;
}

.title:hover {
    cursor: default;
}

img {
    width: 14px;
}

html.light .el-button.is-text:not(.is-disabled):hover {
    background-color: #EDEDED;
}

.el-button.is-text:not(.is-disabled).close:hover {
    background-color: red;
    color: #fff;
}
</style>