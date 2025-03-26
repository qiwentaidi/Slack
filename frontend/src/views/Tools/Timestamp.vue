<script lang="ts" setup>
import { Copy } from '@/util';
import { VideoPause, Switch, VideoPlay, DocumentCopy } from '@element-plus/icons-vue';
import { reactive, ref, onUnmounted } from 'vue'

const timestamp = reactive({
    options: ["秒/s", "毫秒/ms"],
    currentMode: '秒/s',
    now: '',
    time2bj: '',
    cptime: '',
    bj2time: '',
    cpbj: ''
})

let intervalId: ReturnType<typeof setInterval> | undefined;
let paused = ref(true);

const startPrintingTimestamp = () => {
    intervalId = setInterval(() => {
        timestamp.now = timestamp.currentMode === '秒/s' ? Math.floor(Date.now() / 1000).toString() : Date.now().toString();
    }, 1000);
}

const stopPrintingTimestamp = () => {
    if (intervalId !== undefined) {
        clearInterval(intervalId);
        intervalId = undefined;
    }
}

const togglePrintingTimestamp = () => {
    if (paused.value) {
        startPrintingTimestamp();
    } else {
        stopPrintingTimestamp();
    }
    paused.value = !paused.value;
}

const convertTimestampToTime = () => {
    const timestamps = parseInt(timestamp.cptime, 10);
    const date = timestamp.currentMode === '秒/s' ? new Date(timestamps * 1000) : new Date(timestamps);
    const time = date.toLocaleString('zh-CN', { timeZone: 'Asia/Shanghai' });
    timestamp.time2bj = time;
}

const convertTimeToTimestamp = () => {
    const date = new Date(timestamp.cpbj);
    const timestamps = timestamp.currentMode === '秒/s' ? Math.floor(date.getTime() / 1000) : date.getTime();
    timestamp.bj2time = timestamps.toString();
}

// 在组件被卸载时停止打印时间戳
onUnmounted(() => {
    stopPrintingTimestamp();
});
</script>

<template>
    <el-card>
        <template #header>
            <div class="card-header">
                <span>时间戳转换</span>
                <el-button class="button" :type="paused ? 'primary' : 'danger'" :icon="paused ? VideoPlay : VideoPause"
                    @click="togglePrintingTimestamp">{{ paused ? '开始' : '停止' }}</el-button>
            </div>
        </template>
        <el-form :model="timestamp" label-position="top">
            <el-form-item label="时间格式">
                <el-select v-model="timestamp.currentMode" placeholder="Select">
                    <el-option v-for="item in timestamp.options" :value="item" :label="item" />
                </el-select>
            </el-form-item>
            <el-form-item label="现在 (Unix 时间戳是从1970年1月1日开始所经过的秒数)">
                <el-input v-model="timestamp.now">
                    <template #suffix>
                        <el-button :icon="DocumentCopy" link @click="Copy(timestamp.now)"></el-button>
                    </template>
                </el-input>
            </el-form-item>

            <el-form-item label="时间戳 >> 北京时间">
                <el-col :span="11">
                    <el-input v-model="timestamp.cptime" />
                </el-col>
                <el-col :span="2">
                    <el-button :icon="Switch" @click="convertTimestampToTime">
                        转换
                    </el-button>
                </el-col>
                <el-col :span="11">
                    <el-input v-model="timestamp.time2bj" />
                </el-col>
            </el-form-item>
            <el-form-item label="北京时间 >> 时间戳">
                <el-col :span="11">
                    <el-input v-model="timestamp.cpbj" />
                </el-col>
                <el-col :span="2">
                    <el-button :icon="Switch" @click="convertTimeToTimestamp">
                        转换
                    </el-button>
                </el-col>
                <el-col :span="11">
                    <el-input v-model="timestamp.bj2time" />
                </el-col>
            </el-form-item>
        </el-form>
    </el-card>
</template>