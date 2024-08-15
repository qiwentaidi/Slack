<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Copy, CopyALL, formatURL } from '@/util';
import { JSFind } from 'wailsjs/go/main/App';
import { QuestionFilled, Search } from '@element-plus/icons-vue';
import async from 'async';
import global from "@/global";
import { ElNotification, ElMessage } from 'element-plus';
import { BrowserOpenURL } from 'wailsjs/runtime/runtime';
const config = reactive({
  urls: "",
  loading: false,
  otherURL: false,
  drawer: false,
  perfixURL: '',
})

const typeOptions = [
  'primary',
  'success',
  'info',
  'warning',
  'danger',
]

var dashboardItems = [
  {
    title: 'JS路径',
    icon: 'InfoFilled',
    tagType: ref(global.jsfind.defaultType[0]),
    data: ref([] as LinkSource[]),
  },
  {
    title: '敏感字段',
    icon: 'WarningFilled',
    tagType: ref(global.jsfind.defaultType[3]),
    data: ref([] as LinkSource[]),
  },
  {
    title: 'API接口',
    icon: 'SuccessFilled',
    tagType: ref(global.jsfind.defaultType[4]),
    data: ref([] as LinkSource[]),
  },
  {
    title: '手机号',
    icon: 'WarningFilled',
    tagType: ref(global.jsfind.defaultType[1]),
    data: ref([] as LinkSource[]),
  },
  {
    title: '身份证',
    icon: 'WarningFilled',
    tagType: ref(global.jsfind.defaultType[2]),
    data: ref([] as LinkSource[]),
  },
  {
    title: 'IP或URL',
    icon: 'InfoFilled',
    tagType: ref(global.jsfind.defaultType[5]),
    data: ref([] as LinkSource[]),
  },
]

interface dashboardItem {
  title: string;
  icon: string;
  tagType: any;
  data: any;
  children?: any | undefined | null;
}

interface LinkSource {
  Filed: string;
  Source: string;
}

async function JSFinder() {
  var urls = [] as string[]
  urls = await formatURL(config.urls)
  if (urls.length == 0) {
    ElMessage.warning("可用目标为空");
    return
  }
  config.loading = true
  async.eachLimit(urls, 10, (url: string) => {
    JSFind(url, config.perfixURL).then((result: any) => {
      dashboardItems[0].data.value = result.JS
      dashboardItems[2].data.value = result.APIRoute
      dashboardItems[1].data.value = result.SensitiveField
      dashboardItems[3].data.value = result.ChinesePhone
      dashboardItems[4].data.value = result.ChineseIDCard
      dashboardItems[5].data.value = result.IP_URL
      if (config.otherURL) {
        async.eachLimit(result.IP_URL, 10, (url2: string) => {
          for (const whitedomain of global.jsfind.whiteList.split("\n")) {
            if (!url2.includes(whitedomain)) {
              JSFind(url, config.perfixURL).then((result: any) => {
                dashboardItems[0].data.value.push(result.JS)
                dashboardItems[2].data.value.push(result.APIRoute)
                dashboardItems[1].data.value.push(result.SensitiveField)
                dashboardItems[3].data.value.push(result.ChinesePhone)
                dashboardItems[4].data.value.push(result.ChineseIDCard)
                dashboardItems[5].data.value.push(result.IP_URL)
              })
            }
          }
        })
      }
      ElNotification.success({
        message: 'JSFind Finished!',
        position: "bottom-right"
      });
      config.loading = false
    });
  }
  )
}

const saveConfig = () => {
  localStorage.setItem('jsfind', JSON.stringify(global.jsfind));
  ElNotification.success({
    message: 'Save successful',
    position: 'bottom-right'
  })
}

function changeType(item: dashboardItem) {
  switch (item.title) {
    case "JS路径":
      global.jsfind.defaultType[0] = item.tagType.value
      break
    case "敏感字段":
      global.jsfind.defaultType[3] = item.tagType.value
      break
    case "手机号":
      global.jsfind.defaultType[1] = item.tagType.value
      break
    case "API接口":
      global.jsfind.defaultType[4] = item.tagType.value
      break
    case "IP或URL":
      global.jsfind.defaultType[5] = item.tagType.value
      break
    default:
      global.jsfind.defaultType[2] = item.tagType.value
  }
}

function getLength(arr: any) {
  if (Array.isArray(arr)) {
    return arr.length
  } else {
    return 0
  }
}

const menus = [
    {
        label: "复制",
        click: (menu: any, arg: any) => {
          Copy(arg.data.Filed)
        }
    },
    {
        label: "复制来源",
        click: (menu: any, arg: any) => {
          Copy(arg.data.Source)
        },
        divided: true,
    },
    {
        label: "打开来源链接",
        click: (menu: any, arg: any) => {
          BrowserOpenURL(arg.data.Source)
        }
    }
]
</script>
<template>
  <el-form label-width="auto">
    <el-form-item>
      <div class="head">
        <el-input v-model="config.urls" style="margin-right: 5px;">
          <template #prepend>
            URL地址:
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
        </el-input>
        <el-button type="primary" :icon="Search" @click="JSFinder" v-if="!config.loading">开始检测</el-button>
        <el-button type="primary" loading v-else>正在检测</el-button>
        <el-button color="rgb(194, 194, 196)" style="margin-left: 5px" @click="config.drawer = true">参数设置</el-button>
      </div>
    </el-form-item>
  </el-form>
  <!-- content -->
  <el-card style="height: 90%;">
    <el-scrollbar style="height: 80vh">
      <div v-for="item in dashboardItems">
        <!-- HEADER -->
        <div class="my-header">
          <div>
            <el-icon>
              <component :is="item.icon" />
            </el-icon>
            <span style="font-weight: bold; margin-right: 5px;">{{ item.title }}: {{ getLength(item.data.value)
              }}</span>
          </div>
          <el-button round type="primary" size="small" @click="CopyALL(item.data.value.map(item => item.Filed))">Copy
            ALL</el-button>
        </div>
        <!-- CONTENT -->
        <div class="tag-container">
            <el-tag v-for="(ls, index) in item.data.value" :data="ls"  :type="item.tagType.value" v-menus:right="menus">{{ ls.Filed }}</el-tag>
        </div>
      </div>
    </el-scrollbar>
  </el-card>
  <!-- adv -->
  <el-drawer v-model="config.drawer" size="53%">
    <template #header>
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
      <el-form-item label="是否开启URL爬取:" class="el-margin">
        <el-switch v-model="config.otherURL" inline-prompt active-text="关闭" inactive-text="开启"></el-switch>
      </el-form-item>
      <el-form-item label="白名单域名:" class="el-margin" v-show="config.otherURL">
        <el-input v-model="global.jsfind.whiteList" type="textarea" :rows="5"></el-input>
      </el-form-item>
      <el-form-item label="标记颜色:">
        <el-descriptions :column="1" border style="width: 100%;">
          <el-descriptions-item v-for="(item, index) in dashboardItems" :key="index" :label="item.title">
            <div class="flex-box" style="gap: 8px;">
              <el-segmented v-model="item.tagType.value" :options="typeOptions" @click="changeType(item)" />
              <el-tag :type="item.tagType.value" size="large">example</el-tag>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </el-form-item>
      <el-form-item class="align-right">
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<style>
.tag-container {
  margin-top: 5px;
  margin-bottom: 10px;
  display: flex;
  flex-wrap: wrap;
  /* This allows the tags to wrap to a new line */
  gap: 8px;
  /* Adjust the gap between tags as needed */
}
</style>