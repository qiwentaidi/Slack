<script lang="ts" setup>
import { Close, Minus, Search, Back, Right, RefreshRight } from '@element-plus/icons-vue';
import { Quit, WindowMinimise, WindowReload, WindowToggleMaximise } from "wailsjs/runtime/runtime";
import { IsMacOS } from "wailsjs/go/main/File";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import global from "@/global";
import { onlineOptions }  from '@/stores/options';
import { useRoute } from "vue-router";
import { ref, computed } from "vue";
import updateUI from "./Update.vue";
import onlineIcon from "@/assets/icon/online.svg"
import runnerIcon from "@/assets/icon/apprunner.svg"
import maxmizeIcon from "@/assets/icon/maximize.svg"
import reductionIcon from "@/assets/icon/reduction.svg"
import consoleIcon from "@/assets/icon/console.svg"
import { titlebarStyle, leftStyle, rightStyle, macStyle } from '@/stores/style';
import router from '@/router';

const isMax = ref(false)
const onlineDrawer = ref(false)
const showLogger = ref(false)
const route = useRoute();
const updateDialog = ref(false)

window.addEventListener('resize', () => {
    if (screen.availWidth <= window.innerWidth && screen.availHeight <= window.innerHeight) {
        isMax.value = true
    } else {
        isMax.value = false
    }
});

IsMacOS().then(res => {
    global.temp.isMacOS = res
})
function setTitle(path: string) {
    switch (path) {
        case "/":
            return "Home";
        default:
            return path.split('/').slice(-1)[0];
    }
}

const searchFilter = ref("");
const filteredOptions = computed(() => {
    if (!searchFilter.value) return onlineOptions;
    return onlineOptions.map((group) => ({
        ...group,
        value: group.value.filter((item) =>
            item.name.toLowerCase().includes(searchFilter.value.toLowerCase())
        ),
    }));
});

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
    {
        label: "在线导航",
        icon: onlineIcon,
        action: () => {
            onlineDrawer.value = true
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
        icon: isMax.value ? reductionIcon : maxmizeIcon,
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
            <el-divider direction="vertical" v-if="global.temp.isMacOS" />
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