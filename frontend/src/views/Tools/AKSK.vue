<script lang="ts" setup>
import { reactive, ref } from "vue";
import { GoFetch } from "wailsjs/go/services/App";
import { QuestionFilled, ChromeFilled, DocumentCopy } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox, FormInstance, FormRules } from "element-plus";
import { Copy, generateRandomString, proxys } from "@/util";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";
import CustomTabs from "@/components/CustomTabs.vue";
import { wechatResponseDescription } from "@/stores/options";
import dingtalkIcon from "@/assets/icon/dingtalk.svg"
import wechatIcon from "@/assets/icon/wechat.svg"

const activeName = ref("wechat")

const warning = "First, need to obtain the accesstoken"

const wechat = reactive({
  appid: "",
  secert: "",
  accessToken: "",
});

const dingtalk = reactive({
  appid: "",
  secert: "",
  accessToken: "",
  name: "",
  phone: "",
  addUserDialog: false,
  delUserDialiog: false,
})

const result = ref("");

const wechatOption = ({
  api: [
    {
      name: "查询域名配置",
      method: "POST",
      url: "https://api.weixin.qq.com/wxa/getwxadevinfo?access_token=",
    },
    {
      name: "获取长期订阅用户",
      method: "POST",
      url: "https://api.weixin.qq.com/wxa/business/get_wxa_followers?access_token=",
    },
    {
      name: "获取用户列表(1w)",
      method: "POST",
      url: "https://api.weixin.qq.com/cgi-bin/user/get?count=10000&access_token=",
    },
    {
      name: "获取用户反馈列表",
      method: "GET",
      url: "https://api.weixin.qq.com/wxaapi/feedback/list?access_token=",
    },
  ],
  checksecret: async function () {
    let url = `https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=${wechat.appid}&secret=${wechat.secert}`
    let response: any = await GoFetch("GET", url, "", {}, 10, proxys);
    if (response.Error) {
      result.value += "请求失败\n";
      return;
    }
    const jsonResult = JSON.parse(response.Body);
    if (response.Body.includes("access_token")) {
      wechat.accessToken = jsonResult.access_token;
      ElMessage.success("Successfully");
    }
    result.value = jsonResult;
  },
  doRequest: async function (method: string, url: string) {
    if (!wechat.accessToken) {
      ElMessage.warning(warning);
      return;
    }
    let response = await GoFetch(method, url + wechat.accessToken, "", {}, 10, proxys);
    result.value = JSON.parse(response.Body);
  }
})

const ruleFormRef = ref<FormInstance>()

const dingtalkRules = reactive<FormRules>({
  phone: [
    { required: true, message: "Phone can't be empty", trigger: 'blur' },
  ],
  name: [
    { required: true, message: "Name can't be empty", trigger: 'blur' },
  ],
})

const dingdingOption = ({
  api: [
    {
      name: "获取员工人数",
      method: "POST",
      url: "https://oapi.dingtalk.com/topapi/user/count?access_token=",
      body: `{"only_active":"true"}`,
    },
    {
      name: "获取管理员列表",
      method: "GET",
      url: "https://oapi.dingtalk.com/topapi/user/listadmin?access_token=",
      body: ``
    },
    {
      name: "获取部门用户完整信息(100)",
      method: "POST",
      url: "https://oapi.dingtalk.com/topapi/v2/user/list?access_token=",
      body: `{"dept_id":1,"cursor":0,"size":100}`
    }
  ],
  checksecret: async function () {
    let response = await GoFetch("GET", `https://oapi.dingtalk.com/gettoken?appkey=${dingtalk.appid}&appsecret=${dingtalk.secert}`, "", {}, 10, proxys)
    if (response.Error) {
      result.value += "请求失败\n";
      return;
    }
    const jsonResult = JSON.parse(response.Body);
    result.value = jsonResult;
    if (response.Body.includes("7200")) {
      dingtalk.accessToken = jsonResult.access_token;
      ElMessage.success("Successfully");
    }
  },
  doRequest: async function (method: string, url: string, parameter: string) {
    if (!dingtalk.accessToken) {
      ElMessage.warning(warning);
      return;
    }
    let body = ""
    if (parameter != "") {
      body = JSON.parse(parameter);
    } else {
      body = parameter
    }
    let response = await GoFetch(method, url + dingtalk.accessToken, body, {}, 10, proxys);
    result.value = JSON.parse(response.Body);
  },
  addUser: async function (formEl: FormInstance | undefined) {
    if (!formEl) return
    let isValidate = await formEl.validate()
    if (!isValidate) return
    let rdm = generateRandomString(12)
    let body = `{"name":"${dingtalk.name}","mobile":"${dingtalk.phone}","dept_title_list":[{"dept_id":1,"title":"普通员工"}],"dept_order_list":[{"dept_id":1,"order":1}],"dept_id_list":"1","userid":"${rdm}"}`
    let response = await GoFetch("POST", `https://oapi.dingtalk.com/topapi/v2/user/create?access_token=${dingtalk.accessToken}`, body, {}, 10, proxys)
    result.value = `当前添加用户UserId为: ${rdm} \n这个一定要记住后续删除用户的时候要用到\n\n${JSON.parse(response.Body)}`
  },
  delUser: async function () {
    ElMessageBox.prompt('请输入添加用户时返回的UserId', "删除用户", {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: "UserId can't be empty",
    })
      .then(async ({ value }) => {
        let body = `{ "userid":"${value}" }`
        let response = await GoFetch("POST", `https://oapi.dingtalk.com/topapi/v2/user/delete?access_token=${dingtalk.accessToken}`, body, {}, 10, proxys)
        result.value = JSON.parse(response.Body);
      })
  },
  getAllUsers: function () {
    result.value = `POST方法请求下面接口, size最大值为100, 根据员工总数更改cursor值进行翻页

https://oapi.dingtalk.com/topapi/v2/user/list?access_token=${dingtalk.accessToken}
POST 参数: {"dept_id":1,"cursor":0,"size":100}`
  },
})

function CopyResult() {
  Copy(result.value)
}
</script>

<template>
  <CustomTabs>
    <el-tabs v-model="activeName" type="border-card">
      <el-tab-pane name="wechat">
        <template #label>
          <el-text class="position-center">
            <wechatIcon style="margin-right: 2px;" />Wechat
          </el-text>
        </template>
        <el-form label-width="auto">
          <el-form-item label="Appid">
            <el-input v-model="wechat.appid" />
          </el-form-item>
          <el-form-item label="Secert">
            <el-input v-model="wechat.secert" />
          </el-form-item>
          <el-form-item label="Token">
            <el-input v-model="wechat.accessToken" style="width: 100%">
              <template #suffix>
                <el-button type="primary" link @click="wechatOption.checksecret">获取Token</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Operate">
            <el-button v-for="yw in wechatOption.api" @click="wechatOption.doRequest(yw.method, yw.url)">{{ yw.name
              }}</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane name="dingtalk">
        <template #label>
          <el-text class="position-center">
            <dingtalkIcon style="margin-right: 2px;" />Dingtalk
          </el-text>
        </template>
        <el-form label-width="auto">
          <el-form-item label="Appid">
            <el-input v-model="dingtalk.appid" />
          </el-form-item>
          <el-form-item label="Secert">
            <el-input v-model="dingtalk.secert" />
          </el-form-item>
          <el-form-item label="Token">
            <el-input v-model="dingtalk.accessToken">
              <template #suffix>
                <el-button type="primary" link @click="dingdingOption.checksecret">获取Token</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Operate">
            <el-button v-for="yw in dingdingOption.api" @click="dingdingOption.doRequest(yw.method, yw.url, yw.body)">{{
              yw.name
              }}</el-button>
            <el-button @click="dingdingOption.getAllUsers">获取部门用户完整信息(所有)</el-button>
            <el-button type="danger" plain @click="dingtalk.addUserDialog = true">添加用户</el-button>
            <el-button type="danger" plain @click="dingdingOption.delUser">删除用户</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
    <template #ctrl>
      <el-space v-show="activeName == 'wechat'">
        <el-popover placement="left" :width="630" :height="300" trigger="hover">
          <template #reference>
            <el-button :icon="QuestionFilled" plain type="warning">错误码详情</el-button>
          </template>
          <el-table :data="wechatResponseDescription" style="height: 50vh">
            <el-table-column width="100" property="code" label="错误码" />
            <el-table-column width="200" property="describe" label="错误描述" />
            <el-table-column width="300" property="solution" label="解决方案" />
          </el-table>
        </el-popover>
        <el-button :icon="ChromeFilled" plain type="info" @click="BrowserOpenURL('https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html')">
          API文档
        </el-button>
        <el-button :icon="ChromeFilled" plain type="info" @click="BrowserOpenURL('https://mp.weixin.qq.com/debug/')">
          官方调试工具
        </el-button>
      </el-space>
      <el-button :icon="ChromeFilled" plain type="info" v-show="activeName == 'dingtalk'"
        @click="BrowserOpenURL('https://open.dingtalk.com/document/orgapp/delete-a-user')">API文档</el-button>
    </template>
  </CustomTabs>
  <div class="textarea-container">
    <pre class="pretty-response" style="height: 50vh"><code>{{ result }}</code></pre>
    <div class="action-area">
      <el-button :icon="DocumentCopy" size="small" @click="CopyResult">Copy</el-button>
    </div>
  </div>
  <el-dialog v-model="dingtalk.addUserDialog" title="Dingtalk添加用户" width="40%">
    <el-form ref="ruleFormRef" :model="dingtalk" :rules="dingtalkRules" status-icon label-width="auto">
      <el-form-item label="手机号" prop="phone">
        <el-input v-model="dingtalk.phone" />
      </el-form-item>
      <el-form-item label="姓名" prop="name">
        <el-input v-model="dingtalk.name" />
      </el-form-item>
      <el-form-item style="float: right;">
        <el-button>取消</el-button>
        <el-button type="primary" @click="dingdingOption.addUser" style="float: right;">添加</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
