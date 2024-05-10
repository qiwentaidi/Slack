<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { GoFetch, LoadDirsearchDict, PathRequest, SelectFile } from "../../../wailsjs/go/main/App";
import { ReadLine, SplitTextArea, Copy } from '../../util'
import { ElMessage, ElNotification } from 'element-plus'
import { BrowserOpenURL } from '../../../wailsjs/runtime'
import async from 'async';
import { Loading, QuestionFilled } from '@element-plus/icons-vue';
import { onMounted } from 'vue';
import global from '../../global';
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
    paths: [] as string[],
    percentage: 0,
    id: 0,
    currentRate: 0,
    errorCounts: 0,
    redirectClient: false,
    alive: false,
    respDialog: false,
    content: ""
})
const dir = ref([{}])

async function handleFileChange() {
    let path = await SelectFile()
    let result = (await ReadLine(path))!
    const extensions = from.exts.split(',');
    for (const line of result) {
        if (line.includes("%EXT%")) {
            for (const ext of extensions) {
                from.paths.push(line.replace('%EXT%', ext))
            }
        } else {
            from.paths.push(line)
        }
    }
    from.paths = Array.from(new Set(from.paths))
}

async function dirscan() {
    // true is in running 
    if (!control.runningStatus) {
        let ds = new Dirsearch()
        if (await ds.checkURL()) {
            ds.scanner()
        }
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
        if (config.customDict != "") {
            from.paths = SplitTextArea(config.customDict)
        } else {
            await LoadDirsearchDict(global.PATH.ConfigPath + "/dirsearch", from.exts.split(',')).then(result => {
                from.paths = result;
            });
        }
        dir.value = []
        from.id = 0
        from.errorCounts = 0
        let startTime = Date.now();
        let statusCounts: Record<number, number> = {};
        let statuscodeFilter = control.psc()
        let redirect = false
        if (from.redirectClient) {
            redirect = true
        }
        control.runningStatus = true
        async.eachLimit(from.paths, config.thread, (path: string, callback: (err?: Error) => void) => {
            from.id++;
            from.currentRate = Math.round(from.id / ((Date.now() - startTime) / 1000));
            from.percentage = Number(((from.id / from.paths.length) * 100).toFixed(2));
            if (control.runningStatus == false) {
                callback(new Error('用户已终止扫描任务'));
                return;
            }
            if (config.failedCounts != 0 && config.failedCounts <= from.errorCounts) {
                callback(new Error(`失败次数超过${config.failedCounts}次，扫描任务已停止`));
                return;
            }
            PathRequest(from.defaultOption, from.url + path, config.timeout, config.exclude, redirect, config.headers).then(result => {
                if (result.Status == 0) {
                    from.errorCounts++
                } else if (result.Status !== 1) {
                    if (config.timeout != 0) {
                        statusCounts[result.Length] = (statusCounts[result.Length] || 0) + 1;
                        if (statusCounts[result.Length] <= config.times) {
                            if (statuscodeFilter.length > 0) {
                                if (statuscodeFilter.find(number => number === result.Status) != undefined) {
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
                }
                callback()
            })
        }, (err: any) => {
            if (err) {
                ElNotification({
                    message: err,
                    type: 'error',
                    position: 'bottom-right',
                });
                control.runningStatus = false
            } else {
                ElNotification({
                    message: `目录扫描任务已结束`,
                    type: 'success',
                    position: 'bottom-right',
                });
                control.runningStatus = false
            }
        });
    }
}

const control = reactive({
    runningStatus: false,
    stop: function () {
        if (control.runningStatus == true) {
            control.runningStatus = false
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

async function GetResponse(url: string) {
    from.respDialog = true
    let result = await GoFetch("GET", url, "", [{}], 10, null);
    if (result.Error) {
        ElMessage({
            showClose: true,
            message: '目的地址响应超时',
            type: 'warning',
        })
        return
    }
    try {
        const jsonResult = JSON.parse(result.Body);
        from.content = jsonResult;
    } catch (error) {
        from.content = result.Body
    }
}

const config = reactive({
    drawer: false,
    thread: 50,
    timeout: 8,
    times: 5,
    failedCounts: 0,
    exclude: "",
    headers: "",
    customDict: "",
})
</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head">
                <el-select v-model=from.defaultOption value=options style="width: 20vh;">
                    <el-option v-for="item in from.options" :value="item" :label="item" />
                </el-select>
                <el-input v-model="from.url" placeholder="请输入URL地址" style="margin-right: 5px; width: 100%;" />
                <el-button type="primary" @click="dirscan" v-if="!control.runningStatus">开始扫描</el-button>
                <el-button type="danger" @click="control.stop" v-else>停止扫描</el-button>
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
                <el-link type="primary" @click="config.drawer = true">更多参数</el-link>
                <el-tag>字典大小:{{ from.paths.length }}</el-tag>
                <el-tag>线程:{{ config.thread }}</el-tag>
                <el-tooltip placement="bottom" content="请求失败数量">
                    <el-tag type="danger">ERROR:{{ from.errorCounts }}</el-tag>
                </el-tooltip>
            </el-space>
        </el-form-item>
    </el-form>
    <el-drawer v-model="config.drawer" size="50%">
        <template #header>
            <h3>设置高级参数</h3>
        </template>
        <el-form label-width="auto">
            <el-form-item label="线程(MAX 500):">
                <el-input-number v-model="config.thread" :min="1" :max="500" />
            </el-form-item>
            <el-form-item label="超时时长(s):" class="el-margin">
                <el-input-number v-model="config.timeout" :min="1" :max="20" />
            </el-form-item>
            <el-form-item class="el-margin">
                <template #label>
                    <span>过滤长度重复数据:</span>
                    <el-tooltip placement="left">
                        <template #content>值为0时不过滤数据</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="config.times" :min="1" :max="10000" />
            </el-form-item>
            <el-form-item class="el-margin">
                <template #label>
                    <span>失败阈值:</span>
                    <el-tooltip placement="left">
                        <template #content>当目标未响应次数超过失败阈值时，目录扫描任务自动停止<br />值为0时不限制次数</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="config.failedCounts" :min="0" />
            </el-form-item>
            <el-form-item class="el-margin">
                <template #label>
                    <span>扩展名:</span>
                    <el-tooltip placement="left">
                        <template #content>会将字典中%EXT%字段替换，不指定则去除有关%EXT%字段</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input v-model="from.exts"></el-input>
            </el-form-item>
            <el-form-item class="el-margin">
                <template #label>
                    <span>过滤响应体:</span>
                    <el-tooltip placement="left">
                        <template #content>过滤响应包中某些关键字段存在的数据</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input v-model="config.exclude"></el-input>
            </el-form-item>
            <el-form-item label="状态码过滤:" class="el-margin">
                <el-input v-model="from.statusFilter" placeholder="支持200,300 | 200-300,400-500"></el-input>
            </el-form-item>
            <el-form-item label="自定义请求头:" class="el-margin">
                <el-input v-model="config.headers" placeholder="以键:值形式输入，多行请用换行分割" type="textarea" rows="3"></el-input>
            </el-form-item>
            <el-form-item class="el-margin">
                <template #label>
                    <span>自定义字典:</span>
                    <el-tooltip placement="left">
                        <template #content>若文本框中存在内容，则加载其内容目录为字典，不使用内置字典</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input v-model="config.customDict" type="textarea" rows="5"></el-input>
                <el-button type="primary" :icon="Loading" @click="handleFileChange()"
                    style="margin-top: 10px;">选择字典(不选择加载默认字典)</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
    <el-table :data="dir" border style="height: 74vh;">
        <el-table-column type="index" label="#" width="60px" />
        <el-table-column prop="status" width="100px" label="状态码"
            :sort-method="(a: any, b: any) => { return a.status - b.status }" sortable />
        <el-table-column prop="length" width="100px" label="长度"
            :sort-method="(a: any, b: any) => { return a.length - b.length }" sortable />
        <el-table-column prop="path" label="目录路径" :show-overflow-tooltip="true">
            <template #default="scope">
                <el-tooltip placement="top">
                    <template #content>Redirect to {{ scope.row.location }}</template>
                    <el-button link @click.prevent="BrowserOpenURL(scope.row.location)" v-show="scope.row.location != ''">
                    <template #icon>
                        <svg class="bi bi-shuffle" width="1em" height="1em" viewBox="0 0 16 16" fill="currentColor"
                            xmlns="http://www.w3.org/2000/svg">
                            <path fill-rule="evenodd"
                                d="M12.646 1.146a.5.5 0 0 1 .708 0l2.5 2.5a.5.5 0 0 1 0 .708l-2.5 2.5a.5.5 0 0 1-.708-.708L14.793 4l-2.147-2.146a.5.5 0 0 1 0-.708zm0 8a.5.5 0 0 1 .708 0l2.5 2.5a.5.5 0 0 1 0 .708l-2.5 2.5a.5.5 0 0 1-.708-.708L14.793 12l-2.147-2.146a.5.5 0 0 1 0-.708z" />
                            <path fill-rule="evenodd"
                                d="M0 4a.5.5 0 0 1 .5-.5h2c3.053 0 4.564 2.258 5.856 4.226l.08.123c.636.97 1.224 1.865 1.932 2.539.718.682 1.538 1.112 2.632 1.112h2a.5.5 0 0 1 0 1h-2c-1.406 0-2.461-.57-3.321-1.388-.795-.755-1.441-1.742-2.055-2.679l-.105-.159C6.186 6.242 4.947 4.5 2.5 4.5h-2A.5.5 0 0 1 0 4z" />
                            <path fill-rule="evenodd"
                                d="M0 12a.5.5 0 0 0 .5.5h2c3.053 0 4.564-2.258 5.856-4.226l.08-.123c.636-.97 1.224-1.865 1.932-2.539C11.086 4.93 11.906 4.5 13 4.5h2a.5.5 0 0 0 0-1h-2c-1.406 0-2.461.57-3.321 1.388-.795.755-1.441 1.742-2.055 2.679l-.105.159C6.186 9.758 4.947 11.5 2.5 11.5h-2a.5.5 0 0 0-.5.5z" />
                        </svg>
                    </template>
                </el-button>
                </el-tooltip>
                {{ scope.row.path }}
            </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="180px" align="center">
            <template #default="scope">
                <el-button type="primary" link @click.prevent="Copy(scope.row.path)">复制</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="BrowserOpenURL(scope.row.path)">打开</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="GetResponse(scope.row.path)">查看</el-button>
            </template>
        </el-table-column>
        <template #empty>
            <el-empty description="点击开始扫描获取数据"></el-empty>
        </template>
    </el-table>
    <el-dialog v-model="from.respDialog" title="Response" width="800">
        <pre class="pretty-response"><code>{{ from.content }}</code></pre>
    </el-dialog>
    <el-progress :text-inside="true" :stroke-width="18" :percentage="from.percentage" :format="control.format"
        color="#5DC4F7" style="margin-top: 10px;" />
</template>

<style>
.el-drawer__header {
    margin-bottom: 0px;
}

.el-margin {
    margin-top: 18px;
}
</style>