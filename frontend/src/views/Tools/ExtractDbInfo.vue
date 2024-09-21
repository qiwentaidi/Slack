<script lang="ts" setup>
import { ref } from 'vue';
import disconnectIcon from '@/assets/icon/disconnect.svg'
import { Connection } from '@element-plus/icons-vue';

const connected = ref(false)
const databaseOption = ["mysql", "mssql", "oracle", "postgresql", "mongodb"]
</script>

<template>
    <div class="head">
        <el-input placeholder="127.0.0.1:3306@root:123321">
            <template #prepend>
                <el-select style="width: 200px;">
                    <el-option v-for="db in databaseOption" :value="db"
                        :label="db">
                    </el-option>
                </el-select>
            </template>
        </el-input>
        <el-button :icon="Connection" style="margin-left: 5px;" v-if="connected" />
        <el-button :icon="disconnectIcon" style="margin-left: 5px;" v-else />
    </div>
</template>