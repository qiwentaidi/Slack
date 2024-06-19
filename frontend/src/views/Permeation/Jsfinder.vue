<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { CopyALL, formatURL } from '../../util';
import { JSFind } from '../../../wailsjs/go/main/App';
import { QuestionFilled, Search } from '@element-plus/icons-vue';
import async from 'async';
import global from "../../global";
import { ElNotification, ElMessage } from 'element-plus';
import ContextMenu from '../../components/ContextMenu.vue';
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
    children: [
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
    ]
  },
  {
    title: 'API接口',
    icon: 'SuccessFilled',
    tagType: ref(global.jsfind.defaultType[4]),
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
    ElMessage({
      showClose: true,
      message: "可用目标为空",
      type: "warning",
    });
    return
  }
  config.loading = true
  async.eachLimit(urls, 10, (url: string) => {
    JSFind(url, config.perfixURL).then((result: any) => {
      dashboardItems[0].data.value = result.JS
      dashboardItems[2].data.value = result.APIRoute
      dashboardItems[1].data.value = result.SensitiveField
      dashboardItems[1].children![0].data.value = result.ChinesePhone
      dashboardItems[1].children![1].data.value = result.ChineseIDCard
      dashboardItems[3].data.value = result.IP_URL
      if (config.otherURL) {
        async.eachLimit(result.IP_URL, 10, (url2: string) => {
          for (const whitedomain of global.jsfind.whiteList.split("\n")) {
            if (!url2.includes(whitedomain)) {
              JSFind(url, config.perfixURL).then((result: any) => {
                dashboardItems[0].data.value.push(result.JS)
                dashboardItems[2].data.value.push(result.APIRoute)
                dashboardItems[1].data.value.push(result.SensitiveField)
                dashboardItems[1].children![0].data.value.push(result.ChinesePhone)
                dashboardItems[1].children![1].data.value.push(result.ChineseIDCard)
                dashboardItems[3].data.value.push(result.IP_URL)
              })
            }
          }
        })
      }
      ElNotification({
        type: "success",
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
  ElNotification({
    message: 'Save successful',
    type: 'success',
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
        <el-input v-model="config.urls" style="margin-right: 5px;"></el-input>
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
            <span style="font-weight: bold; margin-right: 5px;">{{ item.title }}: {{ item.data.value.length }}</span>
          </div>
          <el-button round type="primary" size="small" @click="CopyALL(item.data.value.map(item => item.Filed))">Copy
            ALL</el-button>
        </div>
        <!-- CONTENT -->
        <div class="tag-container">
          <el-popover placement="right-end" trigger="contextmenu" v-for="(ls, index) in item.data.value" :key="index">
            <template #reference>
              <el-tag :type="item.tagType.value">{{ ls.Filed }}</el-tag>
            </template>
            <template #default>
              <ContextMenu :data="ls" />
            </template>
          </el-popover>
          <div class="tag-container" style="margin-top: 0;" v-if="item.children && item.children.length > 0"
            v-for="(child, index) in item.children" :key="index">
            <el-popover placement="right-end" trigger="contextmenu" v-for="ls in child.data.value">
              <template #reference>
                <el-tag :type="child.tagType.value">{{ ls.Filed }}</el-tag>
              </template>
              <template #default>
                <ContextMenu :data="ls" />
              </template>
            </el-popover>
          </div>
        </div>
      </div>
    </el-scrollbar>
  </el-card>
  <!-- adv -->
  <el-drawer v-model="config.drawer" size="56%">
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
        <el-input v-model="global.jsfind.whiteList" type="textarea" rows="5"></el-input>
      </el-form-item>
      <el-form-item label="标记颜色:">
        <div>
          <el-descriptions :column="1" border>
            <div v-for="(item, index) in dashboardItems" :key="index">
              <el-descriptions-item :label="item.title">
                <div class="flex-box" style="gap: 8px">
                  <el-segmented v-model="item.tagType.value" :options="typeOptions" @click="changeType(item)" />
                  <el-tag :type="item.tagType.value" size="large">example</el-tag>
                </div>
              </el-descriptions-item>
              <el-descriptions-item v-if="item.children && item.children.length > 0" v-for="item2 in item.children"
                :label="item2.title" @click="changeType(item2)">
                <div class="flex-box" style="gap: 8px">
                  <el-segmented v-model="item2.tagType.value" :options="typeOptions" />
                  <el-tag :type="item2.tagType.value" size="large">example</el-tag>
                </div>
              </el-descriptions-item>
            </div>
          </el-descriptions>
          <el-button type="primary" @click="saveConfig" style="float: right; margin-top: 10px;">保存</el-button>
        </div>
      </el-form-item>
    </el-form>
  </el-drawer>
</template>

<style>
.el-drawer__header {
  margin-bottom: 0px;
}

.tag-container {
  margin-top: 5px;
  margin-bottom: 10px;
  display: flex;
  flex-wrap: wrap;
  /* This allows the tags to wrap to a new line */
  gap: 8px;
  /* Adjust the gap between tags as needed */
}

.el-popper.is-light.el-popover {
  padding: 8px;
}
</style>