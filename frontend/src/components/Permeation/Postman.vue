<script setup lang="ts">
import {
    GoSimpleFetch,
} from "../../../wailsjs/go/main/App";
import { reactive } from 'vue'
const form = reactive({
    input: '',
    requestOptions: ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"],
    requestDefault: "GET",
    repsonse: '',

})

function NewFetch() {
    GoSimpleFetch(form.input).then(result => {
        if (result.OK) {
            form.repsonse = result.Text
        }
    })
}

function parsedHtml() {
    // const parser = new DOMParser();
    // const doc = parser.parseFromString(form.repsonse, 'text/html');
    // return doc.body.innerHTML;
}
</script>

<template>
    <el-form>
        <el-form-item>
            <div class="head" style="margin-left: 0px;">
                <el-select v-model=form.requestDefault size="large" value=options>
                    <el-option v-for="item in form.requestOptions" :value="item" :label="item" />
                </el-select>
                <el-input v-model="form.input" placeholder="请输入URL" style="margin-right: 10px;" />
                <el-button type="primary" size="large" @click="NewFetch">Send</el-button>
            </div>
        </el-form-item>
    </el-form>

    <el-tabs>
        <el-tab-pane label="Params">

        </el-tab-pane>
        <el-tab-pane label="Header"></el-tab-pane>
        <el-tab-pane label="Body"></el-tab-pane>
    </el-tabs>
    <el-divider content-position="left">Response</el-divider>
    <el-tabs type="card" v-if="form.repsonse.length > 1">
        <el-tab-pane label="Pretty"></el-tab-pane>
        <el-scrollbar height="500px">
            <!-- <div v-highlight>
                <pre><code>{{ form.repsonse }}</code></pre>
            </div> -->
        </el-scrollbar>

        <el-tab-pane label="Raw">
            <!-- <div id="myHTML">
                {{form.repsonse}}
            </div> -->
            <!-- <el-input type="textarea" rows="20" v-model="form.repsonse"></el-input> -->
        </el-tab-pane>
        <el-tab-pane label="Preview">
            <el-scrollbar height="500px">
                <!-- <div v-html="parsedHtml()"></div> -->
            </el-scrollbar>

        </el-tab-pane>
    </el-tabs>
    <el-empty description="请输入URL点击Send发送获取响应" v-else />
</template>

<style scoped>
.head {
    display: flex;
}
</style>
