<script setup lang="ts">
import { Download } from "@element-plus/icons-vue";
import global from "@/global"
import { UpdatePocFile, DownloadLastestClient, Restart, IsMacOS } from "wailsjs/go/services/File";
import { ElMessageBox, ElNotification } from "element-plus";
import { onMounted, h } from "vue";
import { EventsOn, EventsOff, WindowReload } from "wailsjs/runtime/runtime";
import { renderedMarkdown } from "@/util";

// 存储每个任务的通知实例和相关元素
const taskNotifications: Record<string, { notifyInstance: any, elProgressPlus: HTMLElement | null, elProgressText: HTMLElement | null }> = {}
const clientUpdate = "客户端更新"
const pocUpdate = "POC更新"

onMounted(() => {
    // 监听下载进度事件
    EventsOn("clientDownloadProgress", (p: number) => {
        LoadUpdate(clientUpdate, p, updateSuccess)
    });
    EventsOn("pocDownloadProgress", (p: number) => {
        LoadUpdate(pocUpdate, p, () => { })
    })
    // 清除事件监听器
    return () => {
        EventsOff("clientDownloadProgress");
        EventsOff("pocDownloadProgress");
    };
});


const update = ({
    poc: async function () {
        global.UPDATE.updateDialog = false
        LoadProgress(pocUpdate)
        let isSuccess = await UpdatePocFile(global.UPDATE.RemotePocVersion)
        if (isSuccess) {
            ElMessageBox.confirm(
            "POC更新成功，需要重载应用生效",
            {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'success',
                center: true,
            }
        )
            .then(() => {
                WindowReload()
            })
            .catch(() => {
                console.log('User cancelled or chose another option.')
            })
        } else {
            ElNotification.error("POC update failed!");
        }
    },
    client: async function () {
        global.UPDATE.updateDialog = false
        LoadProgress(clientUpdate)
        let result = await DownloadLastestClient()
        if (result.Error) {
            ElNotification.error("Download client failed! " + result.Msg);
        }
    },
})
function LoadProgress(title: string) {
    // 包含进度条和百分比的容器
    const progress = h(
        'div',
        {
            class: 'el-progress-container',
            style: {
                width: '230px',
                height: '6px',
                backgroundColor: '#f0f0f0',
                marginTop: '6px',
                borderRadius: '6px',
                position: 'relative',
            },
        },
        [
            h(
                'div',
                {
                    class: 'el-progress-plus',
                    style: {
                        width: '0%',
                        height: '100%',
                        backgroundColor: '#5e7ce0',
                        borderRadius: '6px',
                    },
                    percentage: '0',
                },
                ''
            ),
            h(
                'div',
                {
                    class: 'el-progress-text',
                    style: {
                        position: 'absolute',
                        top: '-20px',
                        right: '0',
                        fontSize: '12px',
                    },
                },
                '0%'
            ),
        ]
    )

    // 唯一的通知类名
    const className = `el-notification-plus-${title}`

    // 实例化通知
    const notifyInstance = ElNotification({
        title: title,
        customClass: className,
        message: h('div', {}, progress),
        duration: 0,
        position: 'bottom-right',
    })

    // 获取进度条和百分比文本的元素
    const elNotificationPlus = document.getElementsByClassName(className)[0] as HTMLElement
    const elProgressPlus = elNotificationPlus.getElementsByClassName('el-progress-plus')[0] as HTMLElement
    const elProgressText = elNotificationPlus.getElementsByClassName('el-progress-text')[0] as HTMLElement

    // 将实例保存到字典中
    taskNotifications[title] = { notifyInstance, elProgressPlus, elProgressText }
}

function LoadUpdate(taskId: string, p: number, callback: Function) {
    const task = taskNotifications[taskId]
    if (task && task.elProgressPlus && task.elProgressText) {
        const percentage = Math.min(100, p)
        task.elProgressPlus.setAttribute('percentage', percentage.toString())
        task.elProgressPlus.style.width = `${(230 * percentage) / 100}px`
        task.elProgressText.innerText = `${Math.round(percentage)}%`

        if (percentage === 100) {
            setTimeout(() => {
                task.notifyInstance.close()
                delete taskNotifications[taskId]  // 移除已完成的任务
            }, 500)
            callback()
        }
    }
}

function updateSuccess() {
    IsMacOS().then(ismac => {
        let message = ""
        if (ismac) {
            message = "更新成功，是否立即安装?"
        } else {
            message = "更新成功，是否重新启动?"
        }
        ElMessageBox.confirm(
            message,
            {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'success',
                center: true,
            }
        )
            .then(() => {
                Restart()
            })
            .catch(() => {
                console.log('User cancelled or chose another option.')
            })
    })
}
</script>

<template>
    <el-card class="box-card">
        <template #header>
            <div class="card-header">
                <span>POC&指纹: 最新{{ global.UPDATE.RemotePocVersion }}/当前{{
                    global.UPDATE.LocalPocVersion }}</span>
                <el-button type="primary" :icon="Download" text @click="update.poc"
                    v-if="global.UPDATE.PocStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.PocContent }}</span>
            </div>
        </template>
        <el-scrollbar v-if="global.UPDATE.PocStatus">
            <div v-html="renderedMarkdown(global.UPDATE.PocContent)" style="padding-inline: 10px;"></div>
        </el-scrollbar>
    </el-card>

    <el-card class="box-card" style="margin-top: 10px;">
        <template #header>
            <div class="card-header">
                <span>客户端: 最新{{ global.UPDATE.RemoteClientVersion }}/当前{{
                    global.LOCAL_VERSION }}</span>
                <el-button type="primary" :icon="Download" text @click="update.client"
                    v-if="global.UPDATE.ClientStatus">立即下载</el-button>
                <span v-else>{{ global.UPDATE.ClientContent }}</span>
            </div>
        </template>
        <el-scrollbar v-if="global.UPDATE.ClientStatus">
            <div v-html="renderedMarkdown(global.UPDATE.ClientContent)" style="padding-inline: 10px;"></div>
        </el-scrollbar>
    </el-card>
</template>

<style scoped>
.el-scrollbar {
    height: 150px;
}
</style>