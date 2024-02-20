<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { GoFetch, LoadDirsearchDict, PathRequest, SelectFile } from "../../../wailsjs/go/main/App";
import { GetFileContent } from "../../../wailsjs/go/main/File";
import { ElMessage } from 'element-plus'
import async from 'async';
import { QuestionFilled, Loading } from '@element-plus/icons-vue';
import { onMounted } from 'vue';
// 初始化时调用
onMounted(() => {
    dir.value = []
});
const from = reactive({
    url: '',
    options: ['GET', 'POST', 'HEAD', 'OPTIONS'],
    defaultOption: 'GET',
    exts: 'php,aspx,asp,jsp,html,js',
    statusFilter: '',
    header: '',
    paths: [] as string[],
    percentage: 0,
    id: 0,
    tips: '选择字典(默认加载dicc.txt)',
    currentRate: 0,
    errorCounts: 0,
    redirectClient: false,
    alive: false,
})
const dir = ref([{}])

async function handleFileChange() {
    let path = await SelectFile()
    let res = await GetFileContent(path)
    if (res.length == 0) {
        ElMessage({
            showClose: true,
            message: '不能上传空文件',
            type: 'warning',
        })
        return
    }
    if (res !== "文件不存在") {
        const result = res.replace(/\r\n/g, '\n'); // 避免windows unix系统差异
        const extensions = from.exts.split(',');
        for (const line of result.split('\n')) {
            if (line.includes("%EXT%")) {
                for (const ext of extensions) {
                    from.paths.push(line.replace('%EXT%', ext))
                }
            } else {
                from.paths.push(line)
            }
        }
        from.paths = Array.from(new Set(from.paths))
        from.tips = `loaded ${from.paths.length} dicts`;
    }
}

async function dirscan() {
    let ds = new Dirsearch()
    if (await ds.checkURL()) {
        ds.scanner()
    }
}


class Dirsearch {
    public async checkURL() {
        if (from.url == "") {
            ElMessage({
                showClose: true,
                message: '请输入URL',
                type: 'warning',
            })
            return false
        }
        if (!from.alive) {
            let result = await GoFetch("GET", from.url, "", [{}], 10, null);
            if (result.Error) {
                ElMessage({
                    showClose: true,
                    message: 'URL格式错误或目标不可达',
                    type: 'warning',
                })
                return false
            }
        }
        return true
    }
    public async scanner() {
        if (from.url[from.url.length - 1] !== "/") {
            from.url += "/"
        }
        if (from.paths.length === 0) {
            await LoadDirsearchDict(window.ConfigPath + "/dirsearch", from.exts.split(',')).then(result => {
                from.paths = result;
                from.tips = `loaded default (${from.paths.length} dicts)`;
            });
        }
        dir.value = []
        from.id = 0
        from.errorCounts = 0
        control.exit = false
        let startTime = Date.now();
        let statusCounts: Record<string, number> = {};
        let filter: number[] = []
        filter = control.psc()
        let redirect = false
        if (from.redirectClient) {
            redirect = true
        }

        async.eachLimit(from.paths, config.thread, (path: string, callback: (err?: Error) => void) => {
            from.id++;
            from.currentRate = Math.round(from.id / ((Date.now() - startTime) / 1000));
            from.percentage = Number(((from.id / from.paths.length) * 100).toFixed(2));
            if (control.exit === true) {
                callback();
                return;
            }
            PathRequest(from.defaultOption, from.url + path, config.timeout, config.exclude, redirect, config.headers).then(result => {
                if (result.Status == 0) {
                    from.errorCounts++
                } else if (result.Status !== 1) {
                    statusCounts[result.Status] = (statusCounts[result.Status] || 0) + 1;
                    if (statusCounts[result.Status] <= config.times) {
                        if (filter.length > 0) {
                            if (filter.find(number => number === result.Status) != undefined) {
                                dir.value.push({
                                    status: result.Status,
                                    length: result.Length,
                                    path: from.url + path,
                                    location: result.Location,
                                });
                            }
                        } else {
                            dir.value.push({
                                status: result.Status,
                                length: result.Length,
                                path: from.url + path,
                                location: result.Location,
                            });
                        }
                    }
                }
                callback()
            })
        }, (err: any) => {
            if (err) {
                ElMessage.error(err);
            } else {
                ElMessage({
                    showClose: true,
                    message: `${from.url}目录扫描结束`,
                    type: 'success',
                });
            }
        });
    }
}

const control = reactive({
    exit: false,
    stop: function () {
        if (control.exit === false) {
            control.exit = true
        }
    },
    format: function (percentage: any) {
        return `${from.id}/${from.paths.length} (${from.currentRate}/s)`
    },
    // Processing status codes
    psc: function (): number[] {
        let temp: number[] = []
        if (from.statusFilter !== "") {
            for (const block of from.statusFilter.split(",")) {
                if (block.indexOf("-") !== -1) {
                    let c = block.split("-")
                    for (var i = Number(c[0]); i <= Number(c[1]); i++) {
                        temp.push(Number(i))
                    }
                } else {
                    temp.push(Number(block))
                }
            }
        }
        return temp
    }
})

const config = reactive({
    drawer: false,
    thread: 50,
    timeout: 8,
    times: 5,
    exclude: "",
    headers: "",
})
</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head" style="width: 100%;">
                <el-select v-model=from.defaultOption value=options style="width: 20vh;">
                    <el-option v-for="item in from.options" :value="item" :label="item" />
                </el-select>
                <el-input v-model="from.url" placeholder="请输入URL地址" style="margin-right: 5px; width: 100%;" />
                <el-button type="primary" @click="dirscan">开始扫描</el-button>
                <el-button type="danger" @click="control.stop">停止</el-button>
            </div>
        </el-form-item>
        <el-form-item>
            <el-space>
                <div>
                    <span>重定向跟随：</span>
                    <el-switch v-model="from.redirectClient" inline-prompt active-text="关闭" inactive-text="开启" />
                </div>
                <div>
                    <span>初始不判断存活：</span>
                    <el-switch v-model="from.alive" inline-prompt active-text="关闭" inactive-text="开启" />
                </div>
                <el-tag type="success">线程:{{ config.thread }}</el-tag>
                <el-tag type="success">超时:{{ config.timeout }}s</el-tag>
                <el-tooltip placement="bottom" content="请求失败数量">
                    <el-tag type="danger">ERROR:{{ from.errorCounts }}</el-tag>
                </el-tooltip>
                <el-link type="primary" @click="config.drawer = true">更多参数</el-link>
            </el-space>
            <el-drawer v-model="config.drawer" title="设置高级参数">
                <el-form label-width="100px" label-position="top">
                    <el-form-item label="线程(MAX 500):" style="margin-bottom: 10px;">
                        <el-input-number v-model="config.thread" :min="1" :max="500" />
                    </el-form-item>
                    <el-form-item label="超时时长(s):" style="margin-bottom: 10px;">
                        <el-input-number v-model="config.timeout" :min="1" :max="20" />
                    </el-form-item>
                    <el-form-item label="过滤长度重复数据:" style="margin-bottom: 10px;">
                        <el-input-number v-model="config.times" :min="1" :max="10000" />
                    </el-form-item>
                    <el-form-item style="margin-bottom: 10px;">
                        <el-tooltip placement="left">
                            <template #content>会将字典中%EXT%字段替换，不指定则去除有关%EXT%字段</template>
                            <span>扩展名:<el-icon>
                                    <QuestionFilled size="24" />
                                </el-icon></span>
                        </el-tooltip>
                        <el-input v-model="from.exts"></el-input>
                    </el-form-item>
                    <el-form-item style="margin-bottom: 10px;">
                        <el-tooltip placement="left">
                            <template #content>过滤某些关键字段存在的数据</template>
                            <span>过滤body内容:<el-icon>
                                    <QuestionFilled size="24" />
                                </el-icon></span>
                        </el-tooltip>
                        <el-input v-model="config.exclude"></el-input>
                    </el-form-item>
                    <el-form-item label="状态码过滤:" style="margin-bottom: 10px;">
                        <el-input v-model="from.statusFilter" placeholder="支持200,300 | 200-300,400-500"></el-input>
                    </el-form-item>
                    <el-form-item label="自定义请求头:" style="margin-bottom: 10px;">
                        <el-input v-model="from.header" placeholder="仅支持单行请求头以键:值形式输入" type="textarea" rows="3"></el-input>
                    </el-form-item>
                    <el-form-item label="自定义字典:" style="margin-bottom: 10px;">
                        <el-button type="primary" :icon="Loading" @click="handleFileChange()">{{ from.tips
                        }}</el-button>
                    </el-form-item>
                </el-form>
            </el-drawer>
        </el-form-item>
    </el-form>
    <el-table :data="dir" border style="width: 100%; height: 79vh;">
        <el-table-column type="index" label="#" width="60px" />
        <el-table-column prop="status" width="100px" label="状态码"
            :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable show-overflow-tooltip="true" />
        <el-table-column prop="length" width="100px" label="长度"
            :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable show-overflow-tooltip="true" />
        <el-table-column prop="path" label="完整目录路径" show-overflow-tooltip="true" />
        <el-table-column prop="location" label="跳转路径" show-overflow-tooltip="true" />
    </el-table>
    <el-progress :text-inside="true" :stroke-width="18" :percentage="from.percentage" :format="control.format"
        color="#5DC4F7" style="margin-top: 10px;" />
</template>
