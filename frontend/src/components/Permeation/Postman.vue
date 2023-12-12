<script setup lang="ts">
import { ref } from 'vue'
const options = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"]
const value = ref('GET')
const input = ref('')

</script>

<template>
    <el-container>
        <el-header>
            <div class="head">
                <el-select v-model=value size="large" value=options>
                    <el-option v-for="item in options" :value="item" :label="item" />
                </el-select>
                    <el-input v-model="input" placeholder="Entry URL" style="margin-right: 10px;"/>
                    <el-button type="primary" size="large">Send</el-button>
            </div>
        </el-header>
        <el-main>
            <el-tabs>
                <el-tab-pane label="Header" name="first">Header</el-tab-pane>
                <el-tab-pane label="Body" name="second">Body</el-tab-pane>
            </el-tabs>
            <el-divider content-position="left">Response</el-divider>
            <el-tabs type="card">
                <el-tab-pane label="Pretty" name="first">Pretty</el-tab-pane>
                <el-tab-pane label="Raw" name="second">Raw</el-tab-pane>
                <el-tab-pane label="Preview" name="third">Preview</el-tab-pane>
            </el-tabs>
        </el-main>
        <el-footer>
        </el-footer>
    </el-container>
</template>

<style scoped>
.head {
    display: flex;
}
</style>
