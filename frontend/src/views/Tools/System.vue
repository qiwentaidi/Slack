<script lang="ts" setup>
import { ref } from 'vue'
import { System } from 'wailsjs/go/main/App'
import { onMounted } from 'vue';
import { ElMessage } from 'element-plus';
// 初始化时调用
onMounted(() => {
    avlist.value = [];
    patch.value = [];
});
function checksystem(mode: number) {
    System(textarea.value, mode).then(
        result => {
            if (result == null) {
                if (mode == 0) {
                    ElMessage('未识别到杀软')
                } else {
                    ElMessage('未查询到提权补丁')
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

const activeName = ref('1')
const textarea = ref('')
const avlist = ref([{}])
const patch = ref([{}])

</script>

<template>
    <div class="my-header">
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
    <el-tabs v-model="activeName" type="border-card" style="margin-top: 10px;">
        <el-tab-pane label="杀软识别" name="1">
            <el-table :data="avlist" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="p_name" label="进程名称" />
                <el-table-column prop="pid" label="PID" />
                <el-table-column prop="a_name" label="杀软名称" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
        <el-tab-pane label="补丁检测" name="2">
            <el-table :data="patch" style="height: calc(100vh - 280px);">
                <el-table-column type="index" width="60px" />
                <el-table-column prop="msid" label="微软编号" />
                <el-table-column prop="kbid" label="补丁编号" />
                <el-table-column prop="des" label="描述" />
                <el-table-column prop="windows" label="影响系统" />
                <el-table-column prop="link" label="链接" />
                <template #empty>
                    <el-empty />
                </template>
            </el-table>
        </el-tab-pane>
    </el-tabs>
</template>