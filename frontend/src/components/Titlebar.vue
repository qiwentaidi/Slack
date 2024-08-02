<script lang="ts" setup>
import { Close, Minus, Search, Plus, InfoFilled, Back, Right } from '@element-plus/icons-vue';
import { Quit, WindowMinimise, WindowToggleMaximise } from "../../wailsjs/runtime/runtime";
import { IsMacOS } from "../../wailsjs/go/main/File";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import global from "../global";
import { useRoute } from "vue-router";
import { ref, computed } from "vue";
import about from "./About.vue";
import updateUI from "./Update.vue";

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
                                <svg t="1722350808221" class="icon" viewBox="0 0 1024 1024" version="1.1"
                                    xmlns="http://www.w3.org/2000/svg" p-id="5300" width="200" height="200">
                                    <path
                                        d="M512.032 0C229.216 0 0 229.216 0 512c0 282.752 229.216 512 512.032 512C794.816 1024 1024 794.752 1024 512c0-282.784-229.184-512-511.968-512z m-0.064 960C264.576 960 63.936 759.392 63.936 512c0-247.424 200.64-448 448.032-448C759.488 64 960 264.576 960 512c0 247.392-200.512 448-448.032 448z"
                                        fill="" p-id="5301"></path>
                                    <path
                                        d="M512 208a176.192 176.192 0 0 0-176 176 48 48 0 1 0 96 0c0-44.096 35.872-80 80-80s80 35.904 80 80c0 28.8-15.36 54.688-41.056 69.248a176.768 176.768 0 0 0-88.128 152.576 48 48 0 1 0 96 0c0-28.384 15.328-54.944 39.936-69.248l2.112-1.312A176.16 176.16 0 0 0 688 384c0-97.056-78.944-176-176-176z"
                                        fill="" p-id="5302"></path>
                                    <path d="M512 768m-64 0a64 64 0 1 0 128 0 64 64 0 1 0-128 0Z" fill="" p-id="5303">
                                    </path>
                                </svg>
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
                                <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                                    viewBox="0 0 512 512">
                                    <path d="M448 64L64 240.14h200a8 8 0 0 1 8 8V448z" fill="none" stroke="currentColor"
                                        stroke-linecap="round" stroke-linejoin="round" stroke-width="32"></path>
                                </svg>
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip content="在线导航">
                    <el-button text class="custom-button" @click="onlineDrawer = true">
                        <template #icon>
                            <el-icon size="16">
                                <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                                    viewBox="0 0 512 512">
                                    <path
                                        d="M464 256c0-114.87-93.13-208-208-208S48 141.13 48 256s93.13 208 208 208s208-93.13 208-208z"
                                        fill="none" stroke="currentColor" stroke-miterlimit="10" stroke-width="32">
                                    </path>
                                    <path
                                        d="M445.57 172.14c-16.06.1-14.48 29.73-34.49 15.75c-7.43-5.18-12-12.71-21.33-15c-8.15-2-16.5.08-24.55 1.47c-9.15 1.58-20 2.29-26.94 9.22c-6.71 6.67-10.26 15.62-17.4 22.33c-13.81 13-19.64 27.19-10.7 45.57c8.6 17.67 26.59 27.26 46 26c19.07-1.27 38.88-12.33 38.33 15.38c-.2 9.8 1.85 16.6 4.86 25.71c2.79 8.4 2.6 16.54 3.24 25.21c1.18 16.2 4.16 34.36 12.2 48.67l15-21.16c1.85-2.62 5.72-6.29 6.64-9.38c1.63-5.47-1.58-14.87-1.95-21s-.19-12.34-1.13-18.47c-1.32-8.59-6.4-16.64-7.1-25.13c-1.29-15.81 1.6-28.43-10.58-41.65c-11.76-12.75-29-15.81-45.47-13.22c-8.3 1.3-41.71 6.64-28.3-12.33c2.65-3.73 7.28-6.79 10.26-10.34c2.59-3.09 4.84-8.77 7.88-11.18s17-5.18 21-3.95s8.17 7 11.64 9.56a49.89 49.89 0 0 0 21.81 9.36c13.66 2 42.22-5.94 42-23.46c-.04-8.4-7.84-20.1-10.92-27.96z"
                                        fill="currentColor"></path>
                                    <path
                                        d="M287.45 316.3c-5.33-22.44-35.82-29.94-52.26-42.11c-9.45-7-17.86-17.81-30.27-18.69c-5.72-.41-10.51.83-16.18-.64c-5.2-1.34-9.28-4.14-14.82-3.41c-10.35 1.36-16.88 12.42-28 10.92c-10.55-1.42-21.42-13.76-23.82-23.81c-3.08-12.92 7.14-17.11 18.09-18.26c4.57-.48 9.7-1 14.09.67c5.78 2.15 8.51 7.81 13.7 10.67c9.73 5.33 11.7-3.19 10.21-11.83c-2.23-12.94-4.83-18.22 6.71-27.12c8-6.14 14.84-10.58 13.56-21.61c-.76-6.48-4.31-9.41-1-15.86c2.51-4.91 9.4-9.34 13.89-12.27c11.59-7.56 49.65-7 34.1-28.16c-4.57-6.21-13-17.31-21-18.83c-10-1.89-14.44 9.27-21.41 14.19c-7.2 5.09-21.22 10.87-28.43 3c-9.7-10.59 6.43-14.07 10-21.46s-8.27-21.36-14.61-24.9l-29.81 33.43a41.52 41.52 0 0 0 8.34 31.86c5.93 7.63 15.37 10.08 15.8 20.5c.42 10-1.14 15.12-7.68 22.15c-2.83 3-4.83 7.26-7.71 10.07c-3.53 3.43-2.22 2.38-7.73 3.32c-10.36 1.75-19.18 4.45-29.19 7.21C95.34 199.94 93.8 172.69 86.2 162l-25 20.19c-.27 3.31 4.1 9.4 5.29 13c6.83 20.57 20.61 36.48 29.51 56.16c9.37 20.84 34.53 15.06 45.64 33.32c9.86 16.2-.67 36.71 6.71 53.67c5.36 12.31 18 15 26.72 24c8.91 9.09 8.72 21.53 10.08 33.36a305.22 305.22 0 0 0 7.45 41.28c1.21 4.69 2.32 10.89 5.53 14.76c2.2 2.66 9.75 4.95 6.7 5.83c4.26.7 11.85 4.68 15.4 1.76c4.68-3.84 3.43-15.66 4.24-21c2.43-15.9 10.39-31.45 21.13-43.35c10.61-11.74 25.15-19.69 34.11-33c8.73-12.98 11.36-30.49 7.74-45.68zm-33.39 26.32c-6 10.71-19.36 17.88-27.95 26.39c-2.33 2.31-7.29 10.31-10.21 8.58c-2.09-1.24-2.8-11.62-3.57-14a61.17 61.17 0 0 0-21.71-29.95c-3.13-2.37-10.89-5.45-12.68-8.7c-2-3.53-.2-11.86-.13-15.7c.11-5.6-2.44-14.91-1.06-20c1.6-5.87-1.48-2.33 3.77-3.49c2.77-.62 14.21 1.39 17.66 2.11c5.48 1.14 8.5 4.55 12.82 8c11.36 9.11 23.87 16.16 36.6 23.14c9.86 5.46 12.76 12.37 6.46 23.62z"
                                        fill="currentColor"></path>
                                    <path
                                        d="M184.46 67.09c4.74 4.63 9.2 10.11 16.27 10.57c6.69.45 13-3.17 18.84 1.38c6.48 5 11.15 11.33 19.75 12.89c8.32 1.51 17.13-3.35 19.19-11.86c2-8.11-2.31-16.93-2.57-25.07c0-1.13.61-6.15-.17-7c-.58-.64-5.42.08-6.16.1q-8.13.24-16.22 1.12a207.1 207.1 0 0 0-57.18 14.65c2.43 1.68 5.48 2.35 8.25 3.22z"
                                        fill="currentColor"></path>
                                    <path
                                        d="M356.4 123.27c8.49 0 17.11-3.8 14.37-13.62c-2.3-8.23-6.22-17.16-15.76-12.72c-6.07 2.82-14.67 10-15.38 17.12c-.81 8.08 11.11 9.22 16.77 9.22z"
                                        fill="currentColor"></path>
                                    <path
                                        d="M349.62 166.24c8.67 5.19 21.53 2.75 28.07-4.66c5.11-5.8 8.12-15.87 17.31-15.86a15.4 15.4 0 0 1 10.82 4.41c3.8 3.93 3.05 7.62 3.86 12.54c1.81 11.05 13.66.63 16.75-3.65c2-2.79 4.71-6.93 3.8-10.56c-.84-3.39-4.8-7-6.56-10.11c-5.14-9-9.37-19.47-17.07-26.74c-7.41-7-16.52-6.19-23.55 1.08c-5.76 6-12.45 10.75-16.39 18.05c-2.78 5.13-5.91 7.58-11.54 8.91c-3.1.73-6.64 1-9.24 3.08c-7.24 5.7-3.12 19.39 3.74 23.51z"
                                        fill="currentColor"></path>
                                </svg>
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
                                <svg t="1721800148219" class="icon" viewBox="0 0 1024 1024" version="1.1"
                                    xmlns="http://www.w3.org/2000/svg" p-id="7164"
                                    xmlns:xlink="http://www.w3.org/1999/xlink" width="200" height="200">
                                    <path
                                        d="M812.3 959.4H213.7c-81.6 0-148-66.4-148-148V212.9c0-81.6 66.4-148 148-148h598.5c81.6 0 148 66.4 148 148v598.5C960.3 893 893.9 959.4 812.3 959.4zM213.7 120.9c-50.7 0-92 41.3-92 92v598.5c0 50.7 41.3 92 92 92h598.5c50.7 0 92-41.3 92-92V212.9c0-50.7-41.3-92-92-92H213.7z"
                                        fill="" p-id="7165"></path>
                                </svg>
                            </el-icon>
                            <el-icon size="16" v-show="isMax">
                                <svg t="1721800117386" class="icon" viewBox="0 0 1024 1024" version="1.1"
                                    xmlns="http://www.w3.org/2000/svg" p-id="6870"
                                    xmlns:xlink="http://www.w3.org/1999/xlink" width="200" height="200">
                                    <path
                                        d="M959.72 0H294.216a63.96 63.96 0 0 0-63.96 63.96v127.92H64.28A63.96 63.96 0 0 0 0.32 255.84V959.4a63.96 63.96 0 0 0 63.96 63.96h703.56a63.96 63.96 0 0 0 63.96-63.96V792.465h127.92a63.96 63.96 0 0 0 63.96-63.96V63.96A63.96 63.96 0 0 0 959.72 0zM767.84 728.505V959.4H64.28V255.84h703.56z m189.322 0H831.8V255.84a63.96 63.96 0 0 0-63.96-63.96H294.216V63.96H959.72z"
                                        p-id="6871"></path>
                                </svg>
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