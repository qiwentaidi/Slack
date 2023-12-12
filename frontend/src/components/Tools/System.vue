<script lang="ts" setup>
import { ref } from 'vue'
import {
    System
} from '../../../wailsjs/go/main/App'

function checksystem(mode: number) {
    System(textarea.value, mode).then(
        result => {
            if (result.length == 0) {
                if (mode == 0) {
                    alert('并未识别到杀软')
                } else {
                    alert('并未查询到提权补丁')
                }
            } else {
                if (mode == 0) {
                    const mappedResult = result.map(item => {
                        return {
                            p_name: item[0],
                            pid: item[1],
                            a_name: item[2]
                        }
                    })
                    avlist.value = mappedResult
                } else {
                    const mappedResult = result.map(item => {
                        return {
                            msid: item[0],
                            kbid: item[1],
                            des: item[2],
                            windows: item[3],
                            link: item[4]
                        }
                    })
                    patch.value = mappedResult
                }
            }
        }
    )
}

const activeName = ref('杀软识别')
const textarea = ref('')
const avlist = ref([{}])
const patch = ref([{}])

</script>

<template>
    <div style="display: flex; margin-bottom: 20px;">
        <el-input v-model="textarea" :rows="5" type="textarea" placeholder="tasklist /svc or systeminfo" />
        <div>
            <el-col :span="2" style="margin-left: 10px;">
                <el-button color="#626aef" @click="checksystem(0)">杀软识别</el-button>
            </el-col>
            <el-col :span="2" style="margin-left: 10px; margin-top: 10px;">
                <el-button color="#1e3864" @click="checksystem(1)">补丁检测</el-button>
            </el-col>
        </div>
    </div>
    
    <el-tabs v-model="activeName" class="demo-tabs" type="border-card">
        <el-tab-pane label="杀软识别" name="杀软识别">
            <el-table :data="avlist" border style="width: 100%" height="58vh" v-if="avlist.length > 1">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="p_name" label="进程名称" />
                <el-table-column prop="pid" label="PID" />
                <el-table-column prop="a_name" label="杀软名称" />
            </el-table>
            <el-empty description="暂无数据" v-else />
        </el-tab-pane>
        <el-tab-pane label="补丁检测" name="补丁检测">
            <el-table :data="patch" border style="width: 100%" height="58vh" v-if="patch.length > 1">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="msid" label="微软编号" />
                <el-table-column prop="kbid" label="补丁编号" />
                <el-table-column prop="des" label="描述" />
                <el-table-column prop="windows" label="影响系统" />
                <el-table-column prop="link" label="链接" />
            </el-table>
            <el-empty description="暂无数据" v-else />
        </el-tab-pane>
    </el-tabs>

</template>

<style>
.demo-tabs>.el-tabs__content {
    padding: 32px;
    color: #6b778c;
    font-size: 32px;
    font-weight: 600;
}
</style>