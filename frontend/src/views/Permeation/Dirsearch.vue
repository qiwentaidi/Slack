<script lang="ts" setup>
import { reactive } from 'vue';
import { GoFetch, LoadDirsearchDict, DirScan, StopDirScan } from "../../../wailsjs/go/main/App";
import { ReadLine, SplitTextArea, Copy } from '../../util'
import { ElMessage, ElNotification } from 'element-plus'
import { BrowserOpenURL, EventsOn, EventsOff } from '../../../wailsjs/runtime'
import { QuestionFilled } from '@element-plus/icons-vue';
import { onMounted } from 'vue';
import global from '../../global';
import { FileDialog } from '../../../wailsjs/go/main/File';
import { Dir, DirScanOptions } from '../../interface';

onMounted(() => {
    EventsOn("dirsearchLoading", (result: any) => {
        switch (result.Status) {
            case 0:
                from.errorCounts++
                break
            case 1:
                break
            case 999:
                ElNotification({
                    message: result.Message,
                    type: 'error',
                    position: 'bottom-right',
                });
                break
            default:
                table.result.push({
                    Status: result.Status,
                    Length: result.Length,
                    URL: result.URL,
                    Location: result.Location,
                })
                table.pageContent = table.result.slice((table.currentPage - 1) * table.pageSize, (table.currentPage - 1) * table.pageSize + table.pageSize)
                break
        }
    });
    EventsOn("dirsearchProgressID", (id: number) => {
        from.id = id
        from.currentRate = Math.round(from.id / ((Date.now() - global.temp.dirsearchStartTime) / 1000));
        from.percentage = Number(((from.id / global.temp.dirsearchPathConut) * 100).toFixed(2));
    });
    EventsOn("dirsearchComplete ", () => {
        config.runningStatus = false
        from.percentage = 100
    });
    return () => {
        EventsOff("dirsearchLoading");
        EventsOff("dirsearchProgressID");
        EventsOff("dirsearchComplete");
    };
});

const from = reactive({
    url: '',
    options: ['GET', 'POST', 'HEAD', 'OPTIONS'],
    defaultOption: 'GET',
    exts: 'php,aspx,asp,jsp,html,js',
    statusFilter: '404',
    paths: [] as string[],
    percentage: 0,
    id: 0,
    currentRate: 0,
    errorCounts: 0,
    alive: false,
    respDialog: false,
    content: "",
})

const table = reactive({
    currentPage: 1,
    pageSize: 50,
    result: [] as Dir[],
    pageContent: [] as Dir[],
})

function handleSizeChange(val: any) {
    table.pageSize = val;
    table.currentPage = 1;
    table.pageContent = table.result.slice(0, val)
}

function handleCurrentChange(val: any) {
    table.currentPage = val;
    table.pageContent = table.result.slice((val - 1) * table.pageSize, (val - 1) * table.pageSize + table.pageSize)
}

async function handleFileChange() {
    let path = await FileDialog()
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
    if (!config.runningStatus) {
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
    public async init() {
        config.runningStatus = true
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
        global.temp.dirsearchPathConut = from.paths.length
        from.id = 0
        from.errorCounts = 0
        global.temp.dirsearchStartTime = Date.now();
    }

    public async scanner() {
        await this.init()
        let statuscodeFilter = control.psc()
        let option: DirScanOptions = {
            Method: from.defaultOption,
            URL: from.url,
            Paths: from.paths,
            Workers: config.thread,
            Timeout: config.timeout,
            BodyExclude: config.exclude,
            BodyLengthExcludeTimes: config.times,
            StatusCodeExclude: statuscodeFilter,
            FailedCounts: config.failedCounts,
            Redirect: config.redirectClient,
            Interval: config.interval,
            CustomHeader: config.headers,
        }
        await DirScan(option)
    }
}

const control = ({
    stop: async function () {
        if (config.runningStatus) {
            await StopDirScan()
            config.runningStatus = false
            ElNotification({
                message: "用户已终止扫描任务",
                type: 'error',
                position: 'bottom-right',
            });
        }
    },
    format: function () {
        return `${from.id}/${global.temp.dirsearchPathConut} (${from.currentRate}/s)`
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
        from.content = '目的地址响应超时'
    }
    try {
        from.content = JSON.parse(result.Body);
    } catch (error) {
        from.content = result.Body
    }
}

const config = reactive({
    drawer: false,
    thread: 25,
    timeout: 8,
    times: 5,
    interval: 0,
    failedCounts: 0,
    exclude: "",
    headers: "",
    customDict: "",
    redirectClient: false,
    runningStatus: false,
})

</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head">
                <el-select v-model=from.defaultOption value=options style="width: 20vh;">
                    <el-option v-for="item in from.options" :value="item" :label="item" />
                </el-select>
                <el-input v-model="from.url" placeholder="请输入URL地址" />
                <div class="two-end-space5">
                    <el-button type="primary" @click="dirscan" v-if="!config.runningStatus">开始扫描</el-button>
                    <el-button type="danger" @click="control.stop" v-else>停止扫描</el-button>
                </div>
                <el-button color="rgb(194, 194, 196)" @click="config.drawer = true">参数设置</el-button>
            </div>
        </el-form-item>
        <el-form-item>
            <el-space>
                <div>
                    <span>重定向跟随：</span>
                    <el-switch v-model="config.redirectClient" inline-prompt active-text="关闭" inactive-text="开启" />
                </div>
                <div>
                    <span>初始不判断存活：</span>
                    <el-switch v-model="from.alive" inline-prompt active-text="关闭" inactive-text="开启" />
                </div>
                <el-tag>字典大小:{{ global.temp.dirsearchPathConut }}</el-tag>
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
            <el-form-item label="超时时长(s):">
                <el-input-number v-model="config.timeout" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="请求间隔(s):">
                <el-input-number v-model="config.interval" :min="0" :max="60" />
            </el-form-item>
            <el-form-item>
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
            <el-form-item>
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
            <el-form-item>
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
            <el-form-item>
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
            <el-form-item label="排除状态码:">
                <el-input v-model="from.statusFilter" placeholder="支持200,300 | 200-300,400-500"></el-input>
            </el-form-item>
            <el-form-item label="自定义请求头:">
                <el-input v-model="config.headers" placeholder="以键:值形式输入，多行请用换行分割" type="textarea" rows="3"></el-input>
            </el-form-item>
            <el-form-item>
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
                <el-button link @click="handleFileChange()" style="margin-top: 10px;">选择字典(不选择加载默认字典)</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
    <el-table :data="table.pageContent" border style="height: 74vh;">
        <el-table-column type="index" label="#" width="60px" />
        <el-table-column prop="Status" width="100px" label="状态码"
            :sort-method="(a: any, b: any) => { return a.Status - b.Status }" sortable />
        <el-table-column prop="Length" width="100px" label="长度"
            :sort-method="(a: any, b: any) => { return a.Length - b.Length }" sortable />
        <el-table-column prop="URL" label="目录路径" :show-overflow-tooltip="true">
            <template #default="scope">
                <el-tooltip placement="top">
                    <template #content>Redirect to {{ scope.row.Location }}</template>
                    <el-button link @click.prevent="BrowserOpenURL(scope.row.URL)" v-show="scope.row.Location != ''">
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
                {{ scope.row.URL }}
            </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="180px" align="center">
            <template #default="scope">
                <el-button type="primary" link @click.prevent="Copy(scope.row.URL)">复制</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="BrowserOpenURL(scope.row.URL)">打开</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="GetResponse(scope.row.URL)">查看</el-button>
            </template>
        </el-table-column>
        <template #empty>
            <el-empty />
        </template>
    </el-table>
    <div class="my-header" style="margin-top: 10px;">
        <el-progress :text-inside="true" :stroke-width="20" :percentage="from.percentage" :format="control.format"
        color="#5DC4F7" style="width: 40%;" />
        <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
            :pager-count="5" :current-page="table.currentPage" :page-sizes="[50, 100, 200, 500]" :page-size="table.pageSize"
            layout="total, sizes, prev, pager, next" :total="table.result.length">
        </el-pagination>
    </div>
    <el-dialog v-model="from.respDialog" title="Response" width="800">
        <pre class="pretty-response"><code>{{ from.content }}</code></pre>
    </el-dialog>
</template>

<style>
.el-drawer__header {
    margin-bottom: 0px;
}
</style>