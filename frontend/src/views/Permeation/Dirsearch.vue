<script lang="ts" setup>
import { reactive } from 'vue';
import { LoadDirsearchDict, NewDirsearchScanner, ExitScanner } from "wailsjs/go/services/App";
import { ProcessTextAreaInput, Copy, ReadLine, selectFileAndAssign } from '@/util'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { BrowserOpenURL, EventsOn, EventsOff } from 'wailsjs/runtime'
import { RefreshRight, Document, FolderOpened, DocumentCopy, ChromeFilled, Reading, Setting, Delete, WarnTriangleFilled } from '@element-plus/icons-vue';
import { onMounted } from 'vue';
import global from '@/stores';
import { CheckFileStat, FilepathJoin, List, OpenFolder } from 'wailsjs/go/services/File';
import { DirseearchResult } from '@/stores/interface';
import usePagination from '@/usePagination';
import redirectIcon from '@/assets/icon/redirect.svg'
import { DeleteRecordByPath, DeleteRecordsWithTimesEqualOne, GetAllPathsAndTimes, UpdateOrInsertPath } from 'wailsjs/go/services/Database';
import { dirsearch } from 'wailsjs/go/models';
import throttle from 'lodash/throttle';

// 因为目录扫描的进度条更新比较快，使用节流函数每隔1s更新一次, 避免频繁更新导致卡顿
const updatePercentageThrottled = throttle((id: number) => {
    from.percentage = Number(((id / from.total) * 100).toFixed(2));
}, 1000);

onMounted(() => {
    // 获取当前全部字典
    getDictList()
    GetAllPathsAndTimes().then(result => {
        if (Array.isArray(result)) {
            for (const item of result) {
                pathTimes.table.result.push({
                    path: item.Path,
                    times: item.Times,
                })
            }
        }
        pathTimes.ctrl.watchResultChange(pathTimes.table)
    })
    EventsOn("dirsearchLoading", (result: any) => {
        switch (result.Status) {
            case 0:
                from.errorCounts++
                break
            case 1:
                break
            default:
                if (result.Status == 200 || result.Status == 500) {
                    try {
                        let url = new URL(result.URL);
                        UpdateOrInsertPath(url.pathname.substring(1))
                    } catch (error) {
                        console.log(error)
                    }
                }
                pagination.table.result.push(result)
                pagination.ctrl.watchResultChange(pagination.table)
                break
        }
    });
    EventsOn("dirsearchProgressID", (id: number) => {
        updatePercentageThrottled(id);
    });
    EventsOn("dirsearchCounts", (count: number) => {
        from.total = count
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
    errorCounts: 0,
    respDialog: false,
    content: "",
    dictList: [] as string[],
    selectDict: [] as string[],
    checkAll: false,
    indeterminate: false,
    // 待扫描的路径总数: 目标*字段数
    total: 0,
})

let pagination = usePagination<DirseearchResult>(50)
let pathTimes = usePagination<{ path: string, times: number }>(50)
async function getDictList() {
    from.selectDict = []
    from.configPath = await FilepathJoin([global.PATH.homedir, "/slack/config/dirsearch"])
    let files = await List([from.configPath])
    from.dictList = files.filter(item => item.Path.endsWith(".txt")).map(item => item.Path);
}

async function dirscan() {
    // true is in running 
    if (!config.runningStatus) {
        let ds = new Dirsearch()
        if (await ds.checkInput()) {
            ds.scanner()
        }
    } else {
        ExitScanner("[dirsearch]")
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
            from.paths = ProcessTextAreaInput(config.customDict)
        } else {
            if (from.selectDict.length == 0) {
                from.selectDict = [from.configPath + "/dicc.txt"]
            }
            from.paths = await LoadDirsearchDict(from.selectDict, from.exts.split(','))
        }
        from.errorCounts = 0
        pagination.initTable()
    }

    public async scanner() {
        await this.init()
        let statuscodeFilter = processStatusCode()
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
            let option: dirsearch.Options = {
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
                Backupscan: config.backupfileScan
            }
            if (i > 0) {
                option.URLs = pagination.table.result.filter(item => (item.Status == 200 && item.Recursion == i - 1))
                    .map(line => { return line.URL })
            }
            if (option.URLs.length == 0 || !config.runningStatus) {
                return
            }
            await NewDirsearchScanner(option)
        }
        config.runningStatus = false
        ElNotification.success({
            message: "目录扫描已完成",
            position: 'bottom-right',
        });
    }
}


function processStatusCode(): number[] {
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
        if (temp.length == 0) {
            ElMessage.warning("状态码过滤格式错误, 已忽略")
            return []
        }
    }
    return temp
}

function dispalyResponse(response: string) {
    from.respDialog = true
    try {
        const parsedContent = JSON.parse(response);
        // Converts the object back to a JSON string with pretty-printing
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
    backupfileScan: false,
})

function copyHistory(length: number) {
    // 对 items 进行排序，按 times 值降序排列
    const sortedItems = pathTimes.table.result.sort((a, b) => b.times - a.times);
    // 截取前 xxx 项
    const topItems = sortedItems.slice(0, length);
    // 提取 path 值
    const topPaths = topItems.map(item => item.path);
    config.customDict = topPaths.join("\n")
}

async function deleteRecordFormPath(path: string) {
    let isSuccess = await DeleteRecordByPath(path)
    if (isSuccess) {
        pathTimes.table.result = pathTimes.table.result.filter(item => item.path != path)
        pathTimes.ctrl.watchResultChange(pathTimes.table)
        ElMessage.success({
            message: "删除成功",
            grouping: true
        })
        return
    }
    ElMessage.error({
        message: "删除失败",
        grouping: true
    })
}

function deleteRecordsWithTimesEqualOne() {
    ElMessageBox.confirm('确定要删除所有次数为1的记录吗？', '警告', {
        type: 'warning'
    }).then(async () => {
        let isSuccess = await DeleteRecordsWithTimesEqualOne()
        if (isSuccess) {
            pathTimes.table.result = pathTimes.table.result.filter(item => item.times != 1)
            pathTimes.ctrl.watchResultChange(pathTimes.table)
            ElMessage.success({
                message: "删除成功",
                grouping: true
            })
            return
        }
        ElMessage.error({
            message: "删除失败",
            grouping: true
        })
    }).catch(() => {

    })
}
</script>

<template>
    <div class="flex mb-10px">
        <el-input v-model="from.input" placeholder="请输入URL地址或者选择目标文件">
            <template #prepend>
                <el-select v-model=from.defaultOption value=options style="width: 15vh;">
                    <el-option v-for="item in from.options" :value="item" :label="item" />
                </el-select>
            </template>
            <template #suffix>
                <el-button link :icon="Document" @click="selectFileAndAssign(from, 'input', '*.txt')" />
            </template>
        </el-input>
        <el-button :type="config.runningStatus ? 'danger' : 'primary'" @click="dirscan" class="ml-5px">
            {{ config.runningStatus ? '停止扫描' : '开始扫描' }}
        </el-button>
    </div>
    <el-card shadow="never">
        <div class="flex-between mb-5px">
            <el-space>
                <el-switch v-model="config.backupfileScan" inline-prompt active-text="开启备份文件扫描模式" inactive-text="关闭备份文件扫描模式" />
                <el-switch v-model="config.redirectClient" inline-prompt active-text="开启重定向跟随" inactive-text="关闭重定向跟随" />
                <el-tag>递归层级:{{ config.recursion }}</el-tag>
                <el-tag>
                    <span v-if="!config.backupfileScan">字典大小: {{ from.paths.length }}</span>
                    <span v-else>内置字典规则</span>
                </el-tag>
                <el-tooltip placement="bottom" content="请求失败数量">
                    <el-tag type="danger">ERROR:{{ from.errorCounts }}</el-tag>
                </el-tooltip>
            </el-space>
            <el-button :icon="Setting" @click="config.drawer = true">参数设置</el-button>
        </div>
        <el-table :data="pagination.table.pageContent" 
            @sort-change="pagination.ctrl.sortChange"
            :highlight-current-row="true"
            style="height: calc(100vh - 205px);">
            <el-table-column type="index" label="#" width="60px" />
            <el-table-column prop="Status" width="100px" label="Code" sortable="custom" />
            <el-table-column prop="Length" width="100px" label="Length" sortable="custom" />
            <el-table-column prop="URL" label="URL" :show-overflow-tooltip="true">
                <template #default="scope">
                    <el-tooltip placement="top">
                        <template #content>Redirect to {{ scope.row.Location }}</template>
                        <el-button link @click.prevent="BrowserOpenURL(scope.row.URL)"
                            v-show="scope.row.Location != ''">
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
            <el-table-column prop="Title" width="300px" label="Title" />
            <el-table-column prop="Recursion" width="100px" label="Recursion" />
            <el-table-column label="Operate" width="120px" align="center">
                <template #default="scope">
                    <el-tooltip content="查看响应包">
                        <el-button :icon="Reading" link @click.prevent="dispalyResponse(scope.row.Body)"></el-button>
                    </el-tooltip>
                    <el-tooltip content="打开链接">
                        <el-button :icon="ChromeFilled" link @click.prevent="BrowserOpenURL(scope.row.URL)"></el-button>
                    </el-tooltip>
                    <el-tooltip content="复制链接">
                        <el-button :icon="DocumentCopy" link @click.prevent="Copy(scope.row.URL)"></el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
            <template #empty>
                <el-empty />
            </template>
        </el-table>
        <div class="flex-between mt-5px">
            <el-progress :text-inside="true" :stroke-width="18" :percentage="from.percentage" class="w-40%" />
            <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
                @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
                :current-page="pagination.table.currentPage" :page-sizes="[50, 100, 200, 500]"
                :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
                :total="pagination.table.result.length">
            </el-pagination>
        </div>
    </el-card>
    <el-dialog v-model="from.respDialog" title="Response" width="800">
        <pre class="pretty-response"><code>{{ from.content }}</code></pre>
    </el-dialog>
    <el-drawer v-model="config.drawer" size="60%">
        <template #header>
            <span class="drawer-title">设置高级参数</span>
        </template>
        <el-form label-width="auto">
            <el-form-item label="线程:">
                <el-input-number v-model="config.thread" :min="1" :max="100" />
            </el-form-item>
            <el-form-item label="超时时长(s):">
                <el-input-number v-model="config.timeout" :min="1" :max="20" />
            </el-form-item>
            <el-form-item label="请求间隔(s):">
                <el-input-number v-model="config.interval" :min="0" :max="60" />
            </el-form-item>
            <el-form-item label="递归层级:">
                <div class="form-item-not-fill">
                    <el-input-number v-model="config.recursion" :min="0" :max="5" />
                    <span class="form-item-tips">对响应码为200的路径继续进行扫描</span>
                </div>
            </el-form-item>
            <el-form-item label="过滤长度次数:">
                <div class="form-item-not-fill">
                    <el-input-number v-model="config.times" :min="1" :max="10000" />
                    <span class="form-item-tips">响应长度显示超过n次时不再显示, 值为0时不过滤数据</span>
                </div>
            </el-form-item>
            <el-form-item label="扩展名:">
                <el-input v-model="from.exts"></el-input>
                <span class="form-item-tips">会将字典中%EXT%字段替换，不指定则去除有关%EXT%字段</span>
            </el-form-item>
            <el-form-item label="过滤响应体内容:">
                <el-input v-model="config.exclude"></el-input>
            </el-form-item>
            <el-form-item label="排除状态码:">
                <el-input v-model="from.statusFilter" placeholder="支持200,300 | 200-300,400-500"></el-input>
            </el-form-item>
            <el-form-item label="请求头:">
                <el-input v-model="config.headers" :placeholder="$t('tips.customHeaders')" type="textarea"
                    :rows="3"></el-input>
            </el-form-item>
            <el-form-item label="字典列表:">
                <el-select v-model="from.selectDict" multiple clearable collapse-tags collapse-tags-tooltip
                    placeholder="不选择默认加载dicc字典" :max-collapse-tags="1" class="mb-5px">
                    <template #prefix>
                        <el-button-group>
                            <el-tooltip content="加载自定义字典">
                                <el-button link :icon="Document"
                                    @click="selectFileAndAssign(from, 'selectDict', '*.txt')" />
                            </el-tooltip>
                            <el-tooltip content="打开文件夹">
                                <el-button link :icon="FolderOpened" @click="OpenFolder(from.configPath)" />
                            </el-tooltip>
                            <el-tooltip content="刷新字典列表">
                                <el-button link :icon="RefreshRight" @click="getDictList()" />
                            </el-tooltip>
                        </el-button-group>
                    </template>
                    <el-option v-for="item in from.dictList" :label="item" :value="item" />
                </el-select>
                <el-input v-model="config.customDict" type="textarea" :rows="8"
                    placeholder="若该文本框中存在内容，则加载其内容目录为字典，不使用选中字典"></el-input>
            </el-form-item>
            <el-form-item label="扫描记录:" class="align-right">
                <el-alert title="记录响应码为200、500路径出现次数, 数据非实时, 每次启动/刷新应用时会更新" type="info" show-icon
                    :closable="false"></el-alert>
                <el-button type="danger" plain :icon="WarnTriangleFilled" @click="deleteRecordsWithTimesEqualOne"
                    class="w-full my-5px">删除次数为1的记录</el-button>
                <el-table :data="pathTimes.table.pageContent" border @sort-change="pathTimes.ctrl.sortChange"
                    class="w-full" style="height: 50vh;">
                    <el-table-column prop="path">
                        <template #header>
                            <el-text><span>Path</span>
                                <el-divider direction="vertical" />
                                <el-button size="small" text bg @click="copyHistory(100)">复制Top100</el-button>
                                <el-divider direction="vertical" />
                                <el-button size="small" text bg @click="copyHistory(1000)">复制Top1000</el-button>
                            </el-text>
                        </template>
                    </el-table-column>
                    <el-table-column prop="times" label="Occurrences" width="150" sortable="custom" />
                    <el-table-column label="Operate" width="100" align="center">
                        <template #default="scope">
                            <el-button :icon="Delete" size="small" plain
                                @click="deleteRecordFormPath(scope.row.path)">删除</el-button>
                        </template>
                    </el-table-column>
                    <template #empty>
                        <el-empty />
                    </template>
                </el-table>
                <el-pagination size="small" background @size-change="pathTimes.ctrl.handleSizeChange"
                    @current-change="pathTimes.ctrl.handleCurrentChange" :pager-count="5"
                    :current-page="pathTimes.table.currentPage" :page-sizes="[50, 100, 200]"
                    :page-size="pathTimes.table.pageSize" layout="total, sizes, prev, pager, next"
                    :total="pathTimes.table.result.length" class="mt-5px">
                </el-pagination>
            </el-form-item>
        </el-form>
    </el-drawer>
</template>