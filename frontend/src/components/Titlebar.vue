<script lang="ts" setup>
import { Close, Minus, Back, Right, RefreshRight } from '@element-plus/icons-vue';
import { Quit, WindowMinimise, WindowReload, WindowToggleMaximise } from "wailsjs/runtime/runtime";
import { IsMacOS } from "wailsjs/go/services/File";
import global from "@/global";
import { useRoute } from "vue-router";
import { ref, computed, onMounted } from "vue";
import updateUI from "./Update.vue";
import runnerIcon from "@/assets/icon/apprunner.svg"
import maxmizeIcon from "@/assets/icon/maximize.svg"
import reductionIcon from "@/assets/icon/reduction.svg"
import consoleIcon from "@/assets/icon/console.svg"
import { titlebarStyle, leftStyle, rightStyle, macStyle } from '@/stores/style';
import router from '@/router';

const showLogger = ref(false)
const route = useRoute();
const updateDialog = ref(false)

onMounted(() => {
    IsMacOS().then(res => {
        global.temp.isMacOS = res
    })
    isFullScreen()
})

window.addEventListener('resize', () => {
    isFullScreen()
});

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

const routerControl = [
    {
        label: "返回",
        icon: Back,
        action: () => {
            window.history.back();
        },
    },
    {
        label: "前进",
        icon: Right,
        action: () => {
            window.history.forward();
        }
    },
    {
        label: "刷新",
        icon: RefreshRight,
        action: () => {
            WindowReload()
        }
    },
]

const appControl = [
    {
        label: "运行日志",
        icon: consoleIcon,
        action: () => {
            showLogger.value = true
        },
    },
    {
        label: "应用启动器",
        icon: runnerIcon,
        action: () => {
            router.push('/AppStarter')
        },
    },
]

const windowsControl = computed(() => [
    {
        icon: Minus,
        action: () => {
            WindowMinimise();
        },
    },
    {
        icon: global.temp.isMax ? reductionIcon : maxmizeIcon,
        action: () => {
            WindowToggleMaximise();
        },
    },
    {
        icon: Close,
        action: () => {
            Quit();
        },
        class: 'close',
    },
]);


</script>

<template>
    <div class="titlebar" :style="titlebarStyle">
        <div :style="macStyle">
            <el-divider direction="vertical" v-if="global.temp.isMacOS && !global.temp.isMax" />
            <el-button-group :style="leftStyle">
                <el-tooltip v-for="item in routerControl" :content="item.label">
                    <el-button text class="custom-button" @click="item.action">
                        <el-icon :size="16">
                            <component :is="item.icon" />
                        </el-icon>
                    </el-button>
                </el-tooltip>
            </el-button-group>
        </div>
        <div class="unoccupied" @dblclick="WindowToggleMaximise">
            <span class="title">{{ setTitle(route.path) }}</span>
        </div>
        <div style="display: flex">
            <el-button-group :style="rightStyle">
                <el-tooltip v-for="item in appControl" :content="item.label">
                    <el-button class="custom-button" text @click="item.action">
                        <template #icon>
                            <el-icon :size="16">
                                <component :is="item.icon" />
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
    <el-drawer v-model="showLogger" title="运行日志" direction="rtl" size="50%">
        <div class="log-textarea" v-html="global.Logger.value"></div>
    </el-drawer>
    <!-- update -->
    <el-dialog v-model="updateDialog" title="更新通知" width="40%">
        <updateUI></updateUI>
    </el-dialog>
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

html.light .el-button.is-text:hover {
    background-color: #EDEDED;
}

.el-button.is-text.close:hover {
    background-color: red;
}
</style>