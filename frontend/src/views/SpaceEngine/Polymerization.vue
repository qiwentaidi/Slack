<template>
    <el-input v-model="uncover.query">
        <template #prepend>
            <el-select style="width: 150px;" v-model="uncover.currentGroup">
                <el-option v-for="filed in group" :value="filed" :label="filed" />
            </el-select>
        </template>
        <template #suffix>
            <el-button link type="primary" @click="">查询</el-button>
        </template>
    </el-input>
</template>

<script lang="ts" setup>
import { reactive } from 'vue';
import usePagination from '../../usePagination';

// usePagination()

const group = ["域名", "IP", "Body", "备案名称", "备案号", "标题"]

const uncover = reactive({
    query: '',
    currentGroup: "IP"
})
</script>
