<script lang="ts" setup>
import { reactive, ref } from "vue";
import { GoFetch } from "wailsjs/go/main/App";
import { QuestionFilled, InfoFilled } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { generateRandomString, proxys } from "@/util";
import { BrowserOpenURL } from "wailsjs/runtime/runtime";

const warning = "First, need to obtain the accesstoken"

const wechat = reactive({
  appid: "",
  secert: "",
  accessToken: "",
});

const dingding = reactive({
  appid: "",
  secert: "",
  accessToken: "",
  name: "",
  phone: "",
  userid: "",
})

const wechatOption = ({
  responseDescription: [
    {
      code: "-1",
      describe: "system error",
      solution: "系统繁忙，此时请开发者稍候再试",
    },
    {
      code: "40001",
      describe: "invalid credential  access_token isinvalid or not latest",
      solution:
        "获取 access_token 时 AppSecret 错误，或者 access_token 无效。请开发者认真比对 AppSecret 的正确性，或查看是否正在为恰当的公众号调用接口",
    },
    {
      code: "40013",
      describe: "invalid appid",
      solution:
        "不合法的 AppID ，请开发者检查 AppID 的正确性，避免异常字符，注意大小写",
    },
    {
      code: "40002",
      describe: "invalid grant_type",
      solution: "不合法的凭证类型",
    },
    {
      code: "40125",
      describe: "不合法的 secret",
      solution: "请检查 secret 的正确性，避免异常字符，注意大小写",
    },
    {
      code: "40164",
      describe: "调用接口的IP地址不在白名单中",
      solution: "请在接口IP白名单中进行设置",
    },
    {
      code: "41004",
      describe: "appsecret missing",
      solution: "缺少 secret 参数",
    },
    {
      code: "50004",
      describe: "禁止使用 token 接口",
      solution: "",
    },
    {
      code: "50007",
      describe: "账号已冻结",
      solution: "",
    },
    {
      code: "61024",
      describe: "第三方平台 API 需要使用第三方平台专用 token",
      solution: "",
    },
    {
      code: "40243",
      describe: "AppSecret已被冻结，请登录小程序平台解冻后再次调用。",
      solution: "",
    },
  ],
  api: [
    {
      name: "查询域名配置",
      method: "POST",
      url: "https://api.weixin.qq.com/wxa/getwxadevinfo?access_token=",
    },
    // {
    //   name: '查询实时日志',
    //   method: 'GET',
    //   url: 'https://api.weixin.qq.com/wxaapi/userlog/userlog_search?access_token=',
    //   parameter: 'date=&begintime=&endtime=',
    // },
    {
      name: "获取长期订阅用户",
      method: "POST",
      url: "https://api.weixin.qq.com/wxa/business/get_wxa_followers?access_token=",
    },
    {
      name: "获取用户反馈列表",
      method: "GET",
      url: "https://api.weixin.qq.com/wxaapi/feedback/list?access_token=",
    },
  ],
  checksecret: async function () {
    let url = `https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=${wechat.appid}&secret=${wechat.secert}`
    let response: any = await GoFetch("GET", url, "", [{}], 10, proxys);
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
    let response: any = await GoFetch(method, url + wechat.accessToken, "", [{}], 10, proxys);
    result.value = JSON.parse(response.Body);
  }
})

const dingdingOption = ({
  api: [
    {
      name: "获取员工人数",
      method: "POST",
      url: "https://oapi.dingtalk.com/topapi/user/count?access_token=",
      body: `{
 "only_active":"true"
}`,
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
      body: `{
	"dept_id":1,
	"cursor":0,
	"size":100
}`
    }
  ],
  checksecret: async function () {
    let response: any = await GoFetch("GET", `https://oapi.dingtalk.com/gettoken?appkey=${dingding.appid}&appsecret=${dingding.secert}`, "", [{}], 10, proxys)
    if (response.Error) {
      result.value += "请求失败\n";
      return;
    }
    const jsonResult = JSON.parse(response.Body);
    result.value = jsonResult;
    if (response.Body.includes("7200")) {
      dingding.accessToken = jsonResult.access_token;
      ElMessage.success("Successfully");
    }
  },
  doRequest: async function (method: string, url: string, parameter: string) {
    if (!dingding.accessToken) {
      ElMessage.warning(warning);
      return;
    }
    let body = ""
    if (parameter != "") {
      body = JSON.parse(parameter);
    } else {
      body = parameter
    }
    let response: any = await GoFetch(method, url + dingding.accessToken, body, [{}], 10, proxys);
    result.value = JSON.parse(response.Body);
  },
  addUser: async function () {
    let rdm = generateRandomString(12)
    let body = `{
	"name":"${dingding.name}",
	"mobile":"${dingding.phone}",
	"dept_title_list":[
		{
			"dept_id":1,
			"title":"普通员工"
		}
	],
	"dept_order_list":[
		{
			"dept_id":1,
			"order":1
		}
	],
	"dept_id_list":"1",
	"userid":"${rdm}"
}`
    let response: any = await GoFetch("POST", `https://oapi.dingtalk.com/topapi/v2/user/create?access_token=${dingding.accessToken}`, body, [{}], 10, proxys)
    result.value = `当前添加用户userid为: ${rdm}\n\n${JSON.parse(response.Body)}`
  },
  delUser: async function () {
    let body = `{ "userid":"${dingding.userid}" }`
    let response: any = await GoFetch("POST", `https://oapi.dingtalk.com/topapi/v2/user/delete?access_token=${dingding.accessToken}`, body, [{}], 10, proxys)
    result.value = JSON.parse(response.Body);
  }
})
const result = ref("");
</script>

<template>
  <div style="position: relative; margin-top: 10px;">
    <el-tabs type="border-card" style="height: 40%">
      <el-tab-pane label="Wechat">
        <el-form label-width="auto">
          <el-form-item label="Appid">
            <el-input v-model="wechat.appid" />
          </el-form-item>
          <el-form-item label="Secert">
            <el-input v-model="wechat.secert" />
          </el-form-item>
          <el-form-item label="AccessToken">
            <el-input v-model="wechat.accessToken" style="width: 100%">
              <template #suffix>
                <el-button type="primary" link @click="wechatOption.checksecret">获取Token</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Operation">
            <el-button v-for="yw in wechatOption.api" @click="wechatOption.doRequest(yw.method, yw.url)">{{ yw.name
              }}</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="DingDing">
        <el-form label-width="auto">
          <el-form-item label="Appid">
            <el-input v-model="dingding.appid" />
          </el-form-item>
          <el-form-item label="Secert">
            <el-input v-model="dingding.secert" />
          </el-form-item>
          <el-form-item label="AccessToken">
            <el-input v-model="dingding.accessToken">
              <template #suffix>
                <el-button type="primary" link @click="dingdingOption.checksecret">获取Token</el-button>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="Parameter">
            <div style="display: flex; width: 100%;">
              <el-input v-model="dingding.phone" placeholder="添加用户时所需的手机号" />
              <el-input v-model="dingding.name" placeholder="添加用户时所需的姓名" />
              <el-input v-model="dingding.userid" placeholder="userid添加/删除的时候需要匹配" />
            </div>
          </el-form-item>
          <el-form-item label="Operation">
            <el-button v-for="yw in dingdingOption.api"
              @click="dingdingOption.doRequest(yw.method, yw.url, yw.body)">{{
                yw.name
              }}</el-button>
            <el-button type="danger" @click="dingdingOption.addUser">添加用户</el-button>
            <el-button type="danger" @click="dingdingOption.delUser">删除用户</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
    <el-space class="custom_eltabs_titlebar" :size="5">
      <el-popover placement="left" :width="630" :height="300" trigger="hover">
        <template #reference>
          <el-button :icon="QuestionFilled" text>微信错误码详情</el-button>
        </template>
        <el-table :data="wechatOption.responseDescription" style="height: 50vh">
          <el-table-column width="100" property="code" label="错误码" />
          <el-table-column width="200" property="describe" label="错误描述" />
          <el-table-column width="300" property="solution" label="解决方案" />
        </el-table>
      </el-popover>
      <el-button :icon="InfoFilled" text type="warning"
        @click="BrowserOpenURL('https://open.dingtalk.com/document/orgapp/delete-a-user')">钉钉API使用详情</el-button>
    </el-space>
  </div>

  <pre class="pretty-response" style="height: 48vh"><code>{{ result }}</code></pre>
</template>

<style>
.el-textarea__inner {
  height: 100%;
}
</style>
