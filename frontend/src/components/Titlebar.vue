<script lang="ts" setup>
import { Close, Minus, Search, Plus, InfoFilled } from '@element-plus/icons-vue';
import { Quit, WindowMinimise, WindowToggleMaximise } from "../../wailsjs/runtime/runtime";
import { IsMacOS } from "../../wailsjs/go/main/File";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import global from "../global";
import { useRoute } from "vue-router";
import { ref, computed } from "vue";
import about from "./About.vue";

const isMax = ref(false)
const isMacOS = ref(false)
const onlineDrawer = ref(false)
const localDrawer = ref(false)
const showLogger = ref(false)
const route = useRoute();
const aboutDialog = ref(false)

window.addEventListener('resize', () => {
    if (screen.availWidth == window.innerWidth && screen.availHeight == window.innerHeight) {
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
    return isMacOS ? { marginRight: '3.5px' } : {};
})

const leftStyle = computed(() => {
    return !isMacOS ? { marginLeft: '3.5px' } : {};
})

const macStyle = computed(() => {
    return isMacOS ? { marginLeft: '5.5%' } : {};
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
    <div class="titlebar">
        <div :style="macStyle">
            <el-divider direction="vertical" v-if="isMacOS" />
            <el-button-group :style="leftStyle">
                <el-tooltip content="关于">
                    <el-button class="custom-button" text @click="aboutDialog = true">
                        <template #icon>
                            <img src="../assets/icon/about.svg">
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="运行日志">
                    <el-button class="custom-button" text @click="showLogger = true">
                        <template #icon>
                            <img src="../assets/icon/console.svg">
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
                <el-tooltip content="应用启动器(本功能会在wails v3进行重做)">
                    <el-button text class="custom-button" @click="localDrawer = true" v-if="!isMacOS">
                        <template #icon>
                            <img src="../assets/icon/apprunner.svg">
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="在线导航">
                    <el-button text class="custom-button" @click="onlineDrawer = true">
                        <template #icon>
                            <img src="../assets/icon/online.svg">
                        </template>
                    </el-button>
                </el-tooltip>
            </el-button-group>
            <div v-if="!isMacOS">
                <el-divider direction="vertical" />
                <el-button-group>
                    <el-button text :icon="Minus" @click="WindowMinimise"></el-button>
                    <el-button text @click="WindowToggleMaximise">
                        <template #icon v-if="!isMax">
                            <img src="../assets/icon/maximize.svg">
                        </template>
                        <template #icon v-else>
                            <img src="../assets/icon/reduction.svg">
                        </template>
                    </el-button>
                    <el-button text :icon="Close" class="close" @click="Quit"></el-button>
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
</template>

<style scoped>
.titlebar {
    display: flex;
    width: 100%;
    height: 35px;
    background-color: #fff;
    border-bottom: 1px solid #dcdcdc;

    .el-button {
        height: 35px;
        border-radius: 0;
    }

    .el-button:hover {
        background-color: #dcdcdc;
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

.el-button.is-text:not(.is-disabled).close:hover {
    background-color: red;
    color: #fff;
}
</style>