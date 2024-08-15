<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { LoadDirsearchDict, DirScan, StopDirScan } from "wailsjs/go/main/App";
import { SplitTextArea, Copy, ReadLine } from '@/util'
import { ElMessage, ElNotification } from 'element-plus'
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime'
import { QuestionFilled, RefreshRight, Document, FolderOpened } from '@element-plus/icons-vue';
import { onMounted } from 'vue';
import global from '@/global';
import { CheckFileStat, FileDialog, List, OpenFolder } from 'wailsjs/go/main/File';
import { Dir, DirScanOptions } from '@/interface';
import usePagination from '@/usePagination';
import redirectIcon from '@/assets/icon/redirect.svg'

onMounted(() => {
    // 获取当前全部字典
    getDictList()
    EventsOn("dirsearchLoading", (result: any) => {
        switch (result.Status) {
            case 0:
                from.errorCounts++
                break
            case 1:
                break
            default:
                pagination.table.result.push({
                    Status: result.Status,
                    Length: result.Length,
                    URL: result.URL,
                    Location: result.Location,
                    Body: result.Body,
                    Recursion: result.Recursion,
                })
                pagination.table.pageContent = pagination.ctrl.watchResultChange(pagination.table.result, pagination.table.currentPage, pagination.table.pageSize)
                break
        }
    });
    EventsOn("dirsearchProgressID", (id: number) => {
        from.id = id
        from.currentRate = Math.round(from.id / ((Date.now() - global.temp.dirsearchStartTime) / 1000));
        from.percentage = Number(((from.id / global.temp.dirsearchConut) * 100).toFixed(2));
    });
    EventsOn("dirsearchCounts", (count: number) => {
        global.temp.dirsearchConut = count
    });
    EventsOn("dirsearchComplete", () => {
        config.runningStatus = false
        from.percentage = 100
    });
    return () => {
        EventsOff("dirsearchLoading");
        EventsOff("dirsearchProgressID");
        EventsOff("dirsearchCounts");
        EventsOff("dirsearchComplete");
    };
});

const from = reactive({
    configPath: '',
    input: '',
    options: ['GET', 'POST', 'HEAD', 'OPTIONS'],
    defaultOption: 'GET',
    exts: 'php,aspx,asp,jsp,html,js',
    statusFilter: '404',
    paths: [] as string[],
    percentage: 0,
    id: 0,
    currentRate: 0,
    errorCounts: 0,
    respDialog: false,
    content: "",
    dictList: [] as string[],
    selectDict: [] as string[],
    checkAll: false,
    indeterminate: false,
    result: [] as Dir[],
})

let pagination = usePagination(from.result, 50)

async function getDictList() {
    from.selectDict = []
    from.configPath = global.PATH.homedir + "/slack/config/dirsearch"
    from.dictList = await List(from.configPath)
}

async function handleFileChange() {
    let path = await FileDialog("*.txt")
    if (!path) return
    from.selectDict.push(path)
}

async function GetFilepath() {
    let path = await FileDialog("*.txt")
    if (!path) return
    from.input = path
}
async function dirscan() {
    // true is in running 
    if (!config.runningStatus) {
        let ds = new Dirsearch()
        if (await ds.checkInput()) {
            ds.scanner()
        }
    } else {
        await StopDirScan()
        config.runningStatus = false
        ElNotification.error({
            message: "用户已终止扫描任务",
            position: 'bottom-right',
        });
    }
}


class Dirsearch {
    urls = [] as string[]
    public async checkInput() {
        if (!from.input) {
            ElMessage.warning('请输入URL或者文件路径')
            return false
        }

        try {
            new URL(from.input);
            return true;
        } catch (e) {

        }

        let stat = await CheckFileStat(from.input)
        if (!stat) {
            ElMessage.warning('输入的文件路径不存在')
        }
        return stat
    }
    public async init() {
        config.runningStatus = true
        if (config.customDict != "") {
            from.paths = SplitTextArea(config.customDict)
        } else {
            if (from.selectDict.length == 0) {
                from.selectDict = [from.configPath + "/dicc.txt"]
            }
            await LoadDirsearchDict(from.selectDict, from.exts.split(',')).then(result => {
                from.paths = result;
            });
        }
        global.temp.dirsearchPathConut = from.paths.length
        from.id = 0
        from.errorCounts = 0
        pagination.ctrl.initTable()
        global.temp.dirsearchStartTime = Date.now();
    }

    public async scanner() {
        await this.init()
        let statuscodeFilter = control.psc()
        if (await CheckFileStat(from.input)) {
            let lines = await ReadLine(from.input)
            if (!lines) {
                ElMessage.warning("文件不能为空")
                return
            }
            this.urls = lines
        } else {
            this.urls = [from.input]
        }
        
        for (let i = 0; i <= config.recursion; i++) {
            let option: DirScanOptions = {
                Method: from.defaultOption,
                URLs: this.urls,
                Paths: from.paths,
                Workers: config.thread,
                Timeout: config.timeout,
                BodyExclude: config.exclude,
                BodyLengthExcludeTimes: config.times,
                StatusCodeExclude: statuscodeFilter,
                Redirect: config.redirectClient,
                Interval: config.interval,
                CustomHeader: config.headers,
                Recursion: i,
            }
            if (i > 0) {
                option.URLs = pagination.table.result.filter(item => (item.Status == 200 && item.Recursion == i - 1))
                .map(line => { return line.URL })
            }
            if (option.URLs.length == 0) {
                return
            }
            await DirScan(option)
        }
        config.runningStatus = false
    }
}

const control = ({
    format: function () {
        return `${from.id}/${global.temp.dirsearchConut} (${from.currentRate}/s)`
    },
    // Processing status codes
    psc: function (): number[] {
        let temp: number[] = []
        if (from.statusFilter) {
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

function DispalyResponse(response: string) {
    from.respDialog = true
    try {
        const parsedContent = JSON.parse(response);
         // Converts the object back to a JSON string with pretty-printing
         // code must be string to highlightjs
        from.content = JSON.stringify(parsedContent, null, 2);
    } catch (error) {
        from.content = response
    }
}

const config = reactive({
    drawer: false,
    thread: 25,
    timeout: 8,
    times: 5,
    interval: 0,
    exclude: "",
    headers: "",
    customDict: "",
    redirectClient: false,
    runningStatus: false,
    recursion: 0,
})

const options = ["Pretty", "Raw"]
const optionIndex = ref("Pretty")
</script>

<template>
    <el-form :model="from">
        <el-form-item>
            <div class="head">
                <el-input v-model="from.input" placeholder="请输入URL地址或者选择目标文件">
                    <template #prepend>
                        <el-select v-model=from.defaultOption value=options style="width: 15vh;">
                            <el-option v-for="item in from.options" :value="item" :label="item" />
                        </el-select>
                    </template>
                    <template #suffix>
                        <el-button link :icon="Document" @click="GetFilepath" />
                    </template>
                </el-input>
                <el-space style="margin-left: 5px;">
                    <el-button :type="config.runningStatus ? 'danger' : 'primary'" @click="dirscan">
                        {{ config.runningStatus ? '停止扫描' : '开始扫描' }}
                    </el-button>
                    <el-button text bg @click="config.drawer = true">参数设置</el-button>
                </el-space>
            </div>
        </el-form-item>
        <el-form-item>
            <el-space>
                <div>
                    <span>重定向跟随：</span>
                    <el-switch v-model="config.redirectClient" />
                </div>
                <el-tag>递归层级:{{ config.recursion }}</el-tag>
                <el-tag>字典大小:{{ global.temp.dirsearchPathConut }}</el-tag>
                <el-tag>线程:{{ config.thread }}</el-tag>
                <el-tooltip placement="bottom" content="请求失败数量">
                    <el-tag type="danger">ERROR:{{ from.errorCounts }}</el-tag>
                </el-tooltip>
            </el-space>
        </el-form-item>
    </el-form>
    <el-table :data="pagination.table.pageContent" border style="height: 75vh;">
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
                            <el-icon>
                                <redirectIcon />
                            </el-icon>
                        </template>
                    </el-button>
                </el-tooltip>
                {{ scope.row.URL }}
            </template>
        </el-table-column>
        <el-table-column prop="Recursion" width="100px" label="递归层级" />
        <el-table-column label="操作" width="180px" align="center">
            <template #default="scope">
                <el-button type="primary" link @click.prevent="Copy(scope.row.URL)">复制</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="BrowserOpenURL(scope.row.URL)">打开</el-button>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click.prevent="DispalyResponse(scope.row.Body)">查看</el-button>
            </template>
        </el-table-column>
        <template #empty>
            <el-empty />
        </template>
    </el-table>
    <div class="my-header" style="margin-top: 5px;">
        <el-progress :text-inside="true" :stroke-width="20" :percentage="from.percentage" :format="control.format"
            color="#5DC4F7" style="width: 40%;" />
        <el-pagination background @size-change="pagination.ctrl.handleSizeChange"
            @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
            :current-page="pagination.table.currentPage" :page-sizes="[50, 100, 200, 500]"
            :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
            :total="pagination.table.result.length">
        </el-pagination>
    </div>
    <el-dialog v-model="from.respDialog" width="800">
        <template #header>
            <el-segmented v-model="optionIndex" :options="options"></el-segmented>
        </template>
        <pre class="pretty-response" v-if="optionIndex == 'Raw'"><code>{{ from.content }}</code></pre>
        <highlightjs autodetect :code="from.content" v-else></highlightjs>            
    </el-dialog>
    <el-drawer v-model="config.drawer" size="60%">
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
                    <span>递归层级</span>
                    <el-tooltip placement="left">
                        <template #content>对响应码为200的路径继续进行扫描</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="config.recursion" :min="0" :max="5" />
            </el-form-item>
            <el-form-item>
                <template #label>
                    <span>过滤长度次数:</span>
                    <el-tooltip placement="left">
                        <template #content>响应长度显示超过n次时不再显示，值为0时不过滤数据</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <el-input-number v-model="config.times" :min="1" :max="10000" />
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
                <el-input v-model="config.headers" placeholder="以键:值形式输入，多行请用换行分割" type="textarea" :rows="3"></el-input>
            </el-form-item>
            <el-form-item>
                <template #label>
                    <span>字典列表:</span>
                    <el-tooltip placement="left">
                        <template #content>若文本框中存在内容，则加载其内容目录为字典，不使用选中字典</template>
                        <el-icon>
                            <QuestionFilled size="24" />
                        </el-icon>
                    </el-tooltip>
                </template>
                <div class="head">
                    <el-select v-model="from.selectDict" multiple clearable collapse-tags collapse-tags-tooltip
                        placeholder="不选择默认加载dicc字典" :max-collapse-tags="1">
                        <el-option v-for="item in from.dictList" :label="item" :value="item" />
                    </el-select>
                    <el-space :size="2">
                        <el-tooltip content="加载自定义字典">
                            <el-button text bg :icon="Document" @click="handleFileChange()" />
                        </el-tooltip>
                        <el-tooltip content="打开文件夹">
                            <el-button text bg :icon="FolderOpened" @click="OpenFolder(from.configPath)" />
                        </el-tooltip>
                        <el-tooltip content="刷新字典列表">
                            <el-button text bg :icon="RefreshRight" @click="getDictList()" />
                        </el-tooltip>
                    </el-space>
                </div>
                <el-input v-model="config.customDict" type="textarea" :rows="8"></el-input>
            </el-form-item>
        </el-form>
    </el-drawer>
</template>