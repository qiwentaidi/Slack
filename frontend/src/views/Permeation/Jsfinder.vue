<script lang="ts" setup>
import { reactive } from 'vue';
import { CopyALL, SplitTextArea } from '../../util';
import { JSFind } from '../../../wailsjs/go/main/App';
import { Link, QuestionFilled, Search } from '@element-plus/icons-vue';
import async from 'async';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import global from "../../global";
import { ElNotification } from 'element-plus';
import Loading from "../../components/Loading.vue";
const config = reactive({
  urls: "",
  loading: false,
  otherURL: false,
  thread: 10,
  drawer: false,
  perfixURL: '',
})

const dashboard = reactive({
  jslink: [] as LinkSource[],
  phone: [] as LinkSource[],
  idcard: [] as LinkSource[],
  sensitive: [] as LinkSource[],
  apipath: [] as LinkSource[],
  urls: [] as LinkSource[],
})

interface LinkSource {
  Filed: string;
  Source: string;
}

function JSFinder() {
  var urls = [] as string[]
  urls = SplitTextArea(config.urls)
  ElNotification({
    duration: 0,
    message: '正在进行JS信息爬取...',
    icon: Loading,
  });
  async.eachLimit(urls, 10, (url: string) => {
    JSFind(url, config.perfixURL).then(result => {
      dashboard.jslink = result.JS
      dashboard.apipath = result.APIRoute
      dashboard.sensitive = result.SensitiveField
      dashboard.phone = result.ChinesePhone
      dashboard.idcard = result.ChineseIDCard
      dashboard.urls = result.IP_URL
      if (config.otherURL) {
        async.eachLimit(result.IP_URL, 10, (url2: string) => {
          for (const whitedomain of global.scan.whiteList.split("\n")) {
            if (!url2.includes(whitedomain)) {
              JSFind(url, config.perfixURL).then(result => {
                dashboard.jslink.push(result.JS)
                dashboard.apipath.push(result.APIRoute)
                dashboard.sensitive.push(result.SensitiveField)
                dashboard.phone.push(result.ChinesePhone)
                dashboard.idcard.push(result.ChineseIDCard)
                dashboard.urls.push(result.IP_URL)
              })
            }
          }
        })
      }
      ElNotification.closeAll()
    });
  })
}

const saveConfig = () => {
  localStorage.setItem('scan', JSON.stringify(global.scan));
  ElNotification({
    message: 'Save successful',
    type: 'success',
    position: 'bottom-right'
  })
}
</script>
<template>
  <el-form label-width="auto">
    <el-form-item>
      <template #label>URL地址:
        <el-tooltip placement="right-end">
          <template #content>
            爬取流程:<br />
            1、对目标进行JS提取<br />
            2、再将JS路径进行拼接访问<br />
            3、提取JS中的敏感数据<br />
            4、若开启URL爬取则会对获取到的URL进行1、2、3步重复
          </template>
          <el-icon>
            <QuestionFilled size="24" />
          </el-icon>
        </el-tooltip>
      </template>
      <div class="head">
        <el-input v-model="config.urls"></el-input>
        <el-button type="primary" :icon="Search" style="margin-left: 5px" @click="JSFinder"
          v-if="!config.loading">开始检测</el-button>
        <el-button type="primary" style="margin-left: 5px" loading v-else>正在检测</el-button>
        <el-button color="rgb(194, 194, 196)" style="margin-left: 5px" @click="config.drawer = true">参数设置</el-button>
      </div>
    </el-form-item>
  </el-form>
  <el-drawer v-model="config.drawer" size="50%">
    <template #title>
      <h3>设置高级参数</h3>
    </template>
    <el-form label-width="auto">
      <el-form-item>
        <template #label>
          <span>JS路径拼接前缀:</span>
          <el-tooltip placement="left">
              <template #content>示例: 比如你在爬取二级目录时<br />页面中获取的JS路径需要拼接一级目录访问</template>
              <el-icon>
                  <QuestionFilled size="24" />
              </el-icon>
          </el-tooltip>
        </template>
        <el-input v-model="config.perfixURL" />
      </el-form-item>
      <el-form-item label="爬取线程:" class="el-margin">
        <el-input-number v-model="config.thread" :min="1" :max="500" controls-position="right"></el-input-number>
      </el-form-item>
      <el-form-item label="是否开启URL爬取:" class="el-margin">
        <el-switch v-model="config.otherURL" inline-prompt active-text="关闭" inactive-text="开启"></el-switch>
      </el-form-item>
      <el-form-item label="白名单域名:" class="el-margin" v-show="config.otherURL">
        <el-input v-model="global.scan.whiteList" type="textarea" rows="5"></el-input>
      </el-form-item>
      <el-button type="primary" style="float: right;" v-show="config.otherURL" @click="saveConfig">保存</el-button>
    </el-form>
  </el-drawer>
  <el-row>
    <el-col :sm="12" :lg="6">
      <el-result>
        <template #title>
          <el-text><span style="font-weight: bold; margin-right: 5px;">JS路径: {{ dashboard.jslink.length }}</span>
            <el-button round type="primary" size="small"
              @click="CopyALL(dashboard.jslink.map(item => item.Filed))">Copy</el-button></el-text>
        </template>
        <template #extra>
          <el-card class="js-card">
            <el-scrollbar height="53vh">
              <div v-for="ls in dashboard.jslink">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
            </el-scrollbar>
          </el-card>
        </template>
      </el-result>
    </el-col>
    <el-col :sm="12" :lg="6">
      <el-result icon="warning">
        <template #title>
          <el-text style="font-weight: bold;">敏感信息: {{ dashboard.sensitive.length }}</el-text>
        </template>
        <template #extra>
          <el-card class="js-card" style="margin-top: 2px;">
            <el-scrollbar height="53vh">
              <span style="font-weight: bold;">Phone:</span><br />
              <div v-for="ls in dashboard.phone">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
              <span style="font-weight: bold;">IDCard:</span><br />
              <div v-for="ls in dashboard.idcard">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
              <span style="font-weight: bold;">Sensitive:</span><br />
              <div v-for="ls in dashboard.sensitive">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
            </el-scrollbar>
          </el-card>
        </template>
      </el-result>
    </el-col>
    <el-col :sm="12" :lg="6">
      <el-result icon="success">
        <template #title>
          <el-text><span style="font-weight: bold; margin-right: 5px;">API接口: {{ dashboard.apipath.length }}</span>
            <el-button round type="primary" size="small"
              @click="CopyALL(dashboard.apipath.map(item => item.Filed))">Copy</el-button></el-text>
        </template>
        <template #extra>
          <el-card class="js-card">
            <el-scrollbar height="53vh">
              <div v-for="ls in dashboard.apipath">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
            </el-scrollbar>
          </el-card>
        </template>
      </el-result>
    </el-col>
    <el-col :sm="12" :lg="6">
      <el-result>
        <template #title>
          <el-text><span style="font-weight: bold; margin-right: 5px;">IP或URL: {{ dashboard.urls.length }}</span>
            <el-button round type="primary" size="small"
              @click="CopyALL(dashboard.urls.map(item => item.Filed))">Copy</el-button></el-text>
        </template>
        <template #extra>
          <el-card class="js-card">
            <el-scrollbar height="53vh">
              <div v-for="ls in dashboard.urls">
                <el-text line-clamp="1">
                  <el-button link :icon="Link" @click="BrowserOpenURL(ls.Source)"></el-button>
                  {{ ls.Filed }}
                </el-text>
              </div>
            </el-scrollbar>
          </el-card>
        </template>
      </el-result>
    </el-col>
  </el-row>
</template>

<style>
.el-drawer__header {
  margin-bottom: 0px;
}

.js-card {
  width: 37vh;
  height: 65vh;
  text-align: left;
}
</style>