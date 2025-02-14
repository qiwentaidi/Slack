<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { Copy, parseHeaders, ProcessTextAreaInput } from '@/util';
import { AnalyzeAPI, ExtractAllJSLink, JSFind } from 'wailsjs/go/services/App';
import { ArrowUpBold, ArrowDownBold, Delete, DocumentCopy } from '@element-plus/icons-vue';
import global from "@/stores";
import { ElNotification, ElMessage } from 'element-plus';
import CustomTextarea from '@/components/CustomTextarea.vue';
import saveIcon from '@/assets/icon/save.svg'
import usePagination from '@/usePagination';
import { JSFindOptions } from '@/stores/options';
import { JSFindData } from '@/stores/interface';
import { EventsOff, EventsOn } from 'wailsjs/runtime/runtime';

onMounted(() => {
  EventsOn("jsfindlog", (msg: any) => {
    config.consoleLog += msg + "\n";
  });
  EventsOn("jsfindvulcheck", (result: any) => {
    if (result.IsUnauth) {
      pagination.table.result.push({
        Target: result.Target,
        Method: result.Method,
        Source: "",
        VulType: "Unauth",
        Param: result.Param,
        Length: result.Length,
        Filed: "",
        Response: result.Response,
      })
    }
  });
  return () => {
    EventsOff("jsfindlog");
    EventsOff("webFingerScan");
  };
})

const value = ref(0)

const config = reactive({
  urls: "",
  loading: false,
  otherURL: false,
  prefixURL: '',
  headers: '',
  consoleLog: '',
})

const pagination = usePagination<JSFindData>(20)

async function JSFinder() {
  let blackList = global.jsfind.whiteList.split("\n")
  let urls = ProcessTextAreaInput(config.urls)
  if (urls.length == 0) {
    ElMessage.warning("可用目标为空");
    return
  }
  showForm.value = false
  config.loading = true
  pagination.initTable()
  for (const url of urls) {
    let apiRoute = [] as string[]
    config.consoleLog += `[*] 正在提取${url}的JS链接...\n`
    let jslinks = await ExtractAllJSLink(url)
    config.consoleLog += `[+] 共提取到JS链接: ${getLength(jslinks)}个\n`
    config.consoleLog += jslinks.join("\n")
    config.consoleLog += "\n\n\n"
    let somethings = await JSFind(url, jslinks)
    config.consoleLog += `[+] 共提取到API: ${getLength(somethings.APIRoute)}个\n`
    somethings.APIRoute.forEach(item => {
      apiRoute.push(item.Filed)
      config.consoleLog += `${item.Filed}\n`
    })
    config.consoleLog += "\n\n"
    config.consoleLog += `[+] 共提取到身份证: ${getLength(somethings.IDCard)}个\n`
    somethings.IDCard.forEach(item => {
      config.consoleLog += `${item.Filed} [ Source:: ${item.Source} ]\n`
      pagination.table.result.push({
        Source: item.Source,
        Method: "GET",
        Filed: item.Filed,
        VulType: "身份证号信息泄露",
        Length: 0,
        Param: "",
        Target: url,
        Response: "",
      })
    })
    config.consoleLog += "\n\n"
    config.consoleLog += `[+] 共提取到手机号: ${getLength(somethings.Phone)}个\n`
    somethings.Phone.forEach(item => {
      config.consoleLog += `${item.Filed} [ Source:: ${item.Source} ]\n`
      somethings.IDCard.forEach(item => {
        config.consoleLog += `${item.Filed} [ Source:: ${item.Source} ]\n`
        pagination.table.result.push({
          Source: item.Source,
          Method: "GET",
          Filed: item.Filed,
          VulType: "手机号信息泄露",
          Length: 0,
          Param: "",
          Target: url,
          Response: "",
        })
      })
    })
    config.consoleLog += "\n\n"
    config.consoleLog += `[+] 共提取到邮箱: ${getLength(somethings.Email)}个\n`
    if (somethings.Email) {
      const emails = somethings.Email.map(item => item.Filed)
      config.consoleLog += emails.join("\n")
    }
    config.consoleLog += "\n\n"
    if (somethings.Sensitive) {
      config.consoleLog += `[+] 共提取到敏感字段: ${getLength(somethings.Sensitive)}个\n`
      somethings.Sensitive.forEach(item => {
        pagination.table.result.push({
          Source: item.Source,
          Method: "GET",
          Filed: item.Filed,
          VulType: "敏感字段泄露",
          Length: 0,
          Param: "",
          Target: url,
          Response: "",
        })
        config.consoleLog += `${item.Filed} [ Source: ${item.Source} ]\n`
      })
    }
    config.consoleLog += "\n\n"
    if (somethings.IP_URL && somethings.IP_URL.length > 0) {
      const filteredIPs = somethings.IP_URL
        .filter(item => !blackList.some(black => item.Filed.includes(black))) // 过滤黑名单
        .map(item => item.Filed); // 提取字段

      if (filteredIPs.length > 0) {
        config.consoleLog += `[+] 共提取到IP/URL: ${filteredIPs.length}个\n` + filteredIPs.join("\n") + "\n";
      }
    } 
    config.consoleLog += "\n\n"
    config.consoleLog += "[*] 正在分析API漏洞中...\n"
    let baseURL = ""

    config.prefixURL != "" ? baseURL = config.prefixURL : baseURL = url

    await AnalyzeAPI(url, baseURL, apiRoute, parseHeaders(config.headers))
  }
  config.consoleLog += "[*] 任务运行结束\n"
  config.loading = false
}

function getLength(arr: any) {
  if (Array.isArray(arr)) {
    return arr.length;
  } else {
    return 0;
  }
}

const showForm = ref(true);

function toggleFormVisibility() {
  showForm.value = !showForm.value;
}

function saveConfig() {
  localStorage.setItem('jsfind', JSON.stringify(global.jsfind));
  ElNotification.success({
    message: 'Save successful',
    position: 'bottom-right'
  })
}
</script>

<template>
  <el-divider>
    <el-button round :icon="showForm ? ArrowUpBold : ArrowDownBold" @click="toggleFormVisibility">
      {{ showForm ? '隐藏参数' : '展开参数' }}
    </el-button>
  </el-divider>
  <el-collapse-transition>
    <div style="display: flex; gap: 10px;" v-show="showForm">
      <el-form :model="config" label-width="auto" style="width: 50%;">
        <el-form-item label="目标地址:">
          <CustomTextarea v-model="config.urls" :rows="5" />
        </el-form-item>
        <el-form-item label="路径前缀:">
          <el-input v-model="config.prefixURL" placeholder="部分API无法准确拼接时，自定义路径前缀" />
        </el-form-item>
        <el-form-item label="请求头:">
          <el-input v-model="config.headers" type="textarea" :rows="5"
            :placeholder="$t('tips.customHeaders')"></el-input>
        </el-form-item>
      </el-form>
      <el-form :model="config" label-width="auto" style="width: 50%;">
        <el-form-item label="黑名单域名:">
          <div class="textarea-container">
            <el-input v-model="global.jsfind.whiteList" type="textarea" :rows="5"></el-input>
            <el-tooltip content="保存">
              <el-button class="action-area" :icon="saveIcon" link @click="saveConfig()"></el-button>
            </el-tooltip>
          </div>
          <span class="form-item-tips">过滤IP/URL内容</span>
        </el-form-item>
        <el-form-item label=" ">
          <el-button type="primary" @click="JSFinder">开始任务</el-button>
        </el-form-item>
      </el-form>
    </div>
  </el-collapse-transition>
  <el-card shadow="never" style="margin-bottom: 10px;">
    <template #header>
      <div class="card-header">
        <el-segmented v-model="value" :options="JSFindOptions">
          <template #default="{ item }">
            <el-space :size="3">
              <el-icon>
                <component :is="item.icon" />
              </el-icon>
              <div>{{ item.label }}</div>
            </el-space>
          </template>
        </el-segmented>
        <div v-show="value == 0">
          <el-button :icon="DocumentCopy" link @click="Copy(config.consoleLog)" />
          <el-button :icon="Delete" link @click="config.consoleLog = ''" />
        </div>
      </div>
    </template>
    <pre class="pretty-response" style="margin-top: 0; margin-bottom: 0;" v-show="value == 0"><code>{{ config.consoleLog }}</code></pre>
    <el-table :data="pagination.table.result" style="height: calc(100vh - 200px)" v-show="value == 1">
      <el-table-column type="expand">
          <template #default="props">
              请求方法: {{ props.Method }}
              参数: {{ props.Param }}
              提取内容: {{ props.Filed }}
              响应内容: 
              {{ props.Response }}
          </template>
      </el-table-column>
      <el-table-column prop="Target" label="Target"></el-table-column>
      <el-table-column prop="Length" label="Length" width="120"></el-table-column>
      <el-table-column prop="Source" label="Source"></el-table-column>
      <el-table-column prop="VulType" label="Vulnerability"></el-table-column>
      <template #empty>
        <el-empty />
      </template>
    </el-table>
    <div class="my-header" style="margin-top: 5px;" v-show="value == 1">
      <div></div>
      <el-pagination size="small" background @size-change="pagination.ctrl.handleSizeChange"
        @current-change="pagination.ctrl.handleCurrentChange" :pager-count="5"
        :current-page="pagination.table.currentPage" :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.table.pageSize" layout="total, sizes, prev, pager, next"
        :total="pagination.table.result.length">
      </el-pagination>
    </div>
  </el-card>
</template>

<style scoped></style>