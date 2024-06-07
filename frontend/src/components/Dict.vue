<template>
    <el-descriptions title="字典管理" :column="2" border>
        <el-descriptions-item v-for="item in global.dict.usernames" :label="item.name" :span="2">
            <el-button
                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/username/' + item.name + '.txt'; ReadDict('/username/' + item.name + '.txt')">
                用户名
            </el-button>
        </el-descriptions-item>
        <el-descriptions-item label="密码(所有协议通用)" :span="2">
            <el-button
                @click="ctrl.innerDrawer = true; ctrl.currentPath = '/password/password.txt'; ReadDict('/password/password.txt')">
                密码
            </el-button>
        </el-descriptions-item>
    </el-descriptions>
    <el-drawer v-model="ctrl.innerDrawer" title="字典管理" :append-to-body="true">
        <el-input type="textarea" rows="20" v-model="ctrl.currentDic"></el-input>
        <el-button type="primary" style="margin-top: 10px; float: right;"
            @click="SaveFile(ctrl.currentPath)">保存</el-button>
    </el-drawer>
</template>


<script setup lang="ts">
import { reactive } from 'vue';
import global from '../global';
import { GetFileContent, WriteFile } from '../../wailsjs/go/main/File';
import { ElMessage } from 'element-plus';

const ctrl = reactive({
    drawer: false,
    innerDrawer: false,
    currentDic: '',
    currentPath: '',
})

async function ReadDict(path: string) {
    ctrl.currentDic = await GetFileContent(global.PATH.PortBurstPath + path)
}

async function SaveFile(path: string) {
    let r = await WriteFile('txt', global.PATH.PortBurstPath + path, ctrl.currentDic)
    if (r) {
        ElMessage({
            showClose: true,
            message: 'success',
            type: 'success',
        })
    } else {
        ElMessage({
            showClose: true,
            message: 'failed',
            type: 'error',
        })
    }
}
</script>