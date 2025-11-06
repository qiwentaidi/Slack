<script lang="ts" setup>
import { WindowToggleMaximise, WindowGetSize } from "wailsjs/runtime/runtime";
import { IsMacOS } from "wailsjs/go/services/File";
import global from "@/stores";
import { ref, onMounted } from "vue";
import runnerIcon from "@/assets/icon/apprunner.svg"
import consoleIcon from "@/assets/icon/console.svg"
import languageIcon from "@/assets/icon/language.svg"
import { titlebarStyle, leftStyle, rightStyle, macStyle } from '@/stores/style';
import { routerControl, windowsControl } from '@/stores/options';
import throttle from 'lodash/throttle';
import { SaveWindowsScreenSize } from "wailsjs/go/services/Database";
import { Sunny, Moon } from '@element-plus/icons-vue';
import { useDark, useToggle } from '@vueuse/core'
import { useI18n } from "vue-i18n";

const showLogger = ref(false)

onMounted(() => {
    IsMacOS().then(res => {
        global.temp.isMacOS = res
    })
    isFullScreen()
})

window.addEventListener('resize', () => {
    isFullScreen()
    WindowGetSize().then(size => {
        throttleUpdate(size.w, size.h);
    })
});

const throttleUpdate = throttle((width: number, height: number) => {
    SaveWindowsScreenSize(width, height);
}, 1000);

function isFullScreen() {
    let height = window.innerHeight - screen.availHeight
    // 通用判断：当窗口高度比屏幕高度大于等于20时认为是全屏,
    global.temp.isMax = screen.availWidth == window.innerWidth && height >= 20;
}

function setTitle(path: string) {
    switch (path) {
        case "/":
            return "Home";
        default:
            return path.split('/').slice(-1)[0];
    }
}

const isDark = useDark({
    storageKey: 'theme',
    valueDark: 'dark',
    valueLight: 'light',
})

const toggle = useToggle(isDark)
function toggleTheme() {
    global.Theme.value = !global.Theme.value
    toggle(); // 保留你原来的切换逻辑（比如加/去掉 `.dark` 类）
}

const { locale } = useI18n();

global.Language.value = locale.value
function changeLanguage(lang: string) {
    localStorage.setItem("language", lang);
    locale.value = lang;
}
</script>

<template>
    <div class="titlebar" :style="titlebarStyle">
        <div :style="macStyle">
            <el-divider direction="vertical" v-if="global.temp.isMacOS && !global.temp.isMax" />
            <el-button-group :style="leftStyle">
                <el-tooltip v-for="item in routerControl" :content="$t(item.label)">
                    <el-button text class="custom-button" @click="item.action">
                        <el-icon :size="16">
                            <component :is="item.icon" />
                        </el-icon>
                    </el-button>
                </el-tooltip>
            </el-button-group>
            <el-button class="custom-button" text @click="toggleTheme" v-show="!global.temp.isMacOS">
                <el-icon :size="16">
                    <component :is="global.Theme.value ? Moon : Sunny" />
                </el-icon>
            </el-button>
        </div>
        <div class="unoccupied" @dblclick="WindowToggleMaximise">
            <div class="title-container">
                <span class="title">{{ setTitle($route.path) }}</span>
            </div>
        </div>
        <div class="flex">
            <el-dropdown @command="changeLanguage">
                <el-button class="custom-button" text>
                    <template #icon>
                        <el-icon :size="16">
                            <languageIcon />
                        </el-icon>
                    </template>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item command="zh">中文</el-dropdown-item>
                        <el-dropdown-item command="en">English</el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
            <el-button-group :style="rightStyle">
                <el-button class="custom-button" text @click="toggleTheme" v-show="global.temp.isMacOS">
                    <el-icon :size="16">
                        <component :is="global.Theme.value ? Moon : Sunny" />
                    </el-icon>
                </el-button>
                <el-tooltip :content="$t('titlebar.yx_log')">
                    <el-button class="custom-button" text @click="showLogger = true">
                        <template #icon>
                            <el-icon :size="16">
                                <consoleIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                <el-tooltip :content="$t('titlebar.app_launcher')">
                    <el-button class="custom-button" text @click="$router.push('/AppLauncher')">
                        <template #icon>
                            <el-icon :size="16">
                                <runnerIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
            </el-button-group>
            <div v-if="!global.temp.isMacOS">
                <el-divider direction="vertical" />
                <el-button-group>
                    <el-button v-for="item in windowsControl" :class="item.class!" text @click="item.action">
                        <template #icon>
                            <el-icon size="16">
                                <component :is="item.icon" />
                            </el-icon>
                        </template>
                    </el-button>
                </el-button-group>
            </div>
        </div>
    </div>
    <!-- running logs -->
    <el-drawer v-model="showLogger" :title="$t('titlebar.yx_log')" direction="rtl" size="50%">
        <div class="log-textarea" v-html="global.Logger.value"></div>
    </el-drawer>
</template>

<style scoped>
.titlebar {
    display: flex;
    width: 100%;
    height: var(--titlebar-height);

    .el-button {
        height: var(--titlebar-height);
        border-radius: 0;
    }

    .custom-button {
        margin-top: 3.5px;
        margin-bottom: 3.5px;
        height: 28px;
        width: 35px;
        border-radius: 10px;
    }
}

.unoccupied {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-grow: 1;
    --wails-draggable: drag;
}

.title-container {
    position: absolute;
    /* top: 0; */
    left: 50%;
    transform: translateX(-50%);
    /* height: 100%; */
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: none;
    /* 防止遮挡点击 */
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
    /* margin-right: 5%; */
    font-size: 13px;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 100%;
}

.title:hover {
    cursor: default;
}

html.light .el-button.is-text:hover {
    background-color: #EDEDED;
}

.el-button.is-text.close:hover {
    background-color: red;
}
</style>