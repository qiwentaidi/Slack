<script setup lang="ts">
import {
    GoSimpleFetch,
} from "../../../wailsjs/go/main/App";
import { reactive, onMounted, onUpdated } from 'vue'
import Prism from "prismjs";
onUpdated(() => {
    Prism.highlightAll(); //修改内容后重新渲染
});
onMounted(() => {
    Prism.highlightAll(); //切换菜单重新渲染
})
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
    const parser = new DOMParser();
    const doc = parser.parseFromString(form.repsonse, 'text/html');
    return doc.body.innerHTML;
}
</script>

<template>
    <el-form>
        <el-form-item>
            <div class="head">
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
        <el-tab-pane label="Header">

        </el-tab-pane>
        <el-tab-pane label="Body">

        </el-tab-pane>
    </el-tabs>
    <el-divider content-position="left">Response</el-divider>
    <el-tabs type="border-card" v-if="form.repsonse.length > 1">
        <el-tab-pane label="Pretty">
            <el-scrollbar class="fillContent">
                <pre><code class="language-python line-numbers">{{ form.repsonse }}</code></pre>

                <!-- <pre v-highlight class="pre-wrap"><code class="html">{{ form.repsonse }}</code></pre> -->
            </el-scrollbar>
        </el-tab-pane>
        <el-tab-pane label="Raw">
            <el-scrollbar class="fillContent">
                <pre class="pre-wrap"><code>{{ form.repsonse }}</code></pre>
            </el-scrollbar>
        </el-tab-pane>
        <el-tab-pane label="Preview">
            <el-scrollbar class="fillContent">
                <div v-html="parsedHtml()"></div>
            </el-scrollbar>
        </el-tab-pane>
    </el-tabs>
    <el-empty description="请输入URL点击Send发送获取响应" v-else />
</template>

<style scoped>
.head {
    display: flex;
}

.pre-wrap {
    white-space: pre-wrap;
    word-wrap: break-word;
    word-break: break-all;
    overflow-wrap: break-word;
}

.fillContent {
    height: calc(100vh - 300px);
}
</style>