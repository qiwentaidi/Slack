<script lang="ts" setup>
import { reactive, ref } from "vue";
import { GoFetch } from "../../../wailsjs/go/main/App";
import { Delete, QuestionFilled } from "@element-plus/icons-vue";
import { ElNotification } from "element-plus";
const wechat = reactive({
  appid: "",
  secert: "",
  accessToken: "",
});
const result = ref("");

const wechatErrorDesc = [
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
];

const wechatAPI = {
  // 运维中心
  ywzx: [
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
};

async function checksecert() {
  let url =
    "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" +
    wechat.appid +
    "&secret=" +
    wechat.secert;
  let response = await GoFetch("GET", url, "", [{}], 10, null);
  if (response.Error) {
    result.value += "请求失败\n";
    return;
  }
  const jsonResult = JSON.parse(response.Body);
  if (response.Body.includes("access_token")) {
    wechat.accessToken = jsonResult.access_token;
    ElNotification({
      type: "success",
      message: "AccessToken set successfully",
    });
  }
  result.value = jsonResult;
}

async function doRequest(method: string, url: string, parameter: string) {
  if (wechat.accessToken == "") {
    ElNotification({
      message: "AccessToken is null",
      type: "warning",
    });
    return;
  }
  if (parameter != null) {
    let response = await GoFetch(
      method,
      url + wechat.accessToken + "&" + parameter,
      "",
      [{}],
      10,
      null
    );
    result.value = JSON.parse(response.Body);
    return;
  }
  let response = await GoFetch(
    method,
    url + wechat.accessToken,
    "",
    [{}],
    10,
    null
  );
  result.value = JSON.parse(response.Body);
}
</script>

<template>
  <el-tabs type="border-card" style="height: 40%">
    <el-tab-pane label="Wechat">
      <el-form :model="wechat" label-width="auto">
        <el-form-item label="Appid">
          <el-input v-model="wechat.appid" style="width: 40%" />
          <span style="margin-left: 10px; margin-right: 10px">Secert</span>
          <el-input v-model="wechat.secert" style="width: 40%" />
          <el-button
            type="primary"
            @click="checksecert"
            style="width: 13%; margin-left: 10px"
            >Check Token</el-button
          >
        </el-form-item>
        <el-form-item>
          <template #label>
            AccessToken
            <el-popover
              placement="right"
              :width="630"
              :height="300"
              trigger="hover"
            >
              <template #reference>
                <el-icon>
                  <QuestionFilled size="24" />
                </el-icon>
              </template>
              <el-table :data="wechatErrorDesc" style="height: 50vh">
                <el-table-column width="100" property="code" label="错误码" />
                <el-table-column
                  width="200"
                  property="describe"
                  label="错误描述"
                />
                <el-table-column
                  width="300"
                  property="solution"
                  label="解决方案"
                />
              </el-table>
            </el-popover>
          </template>
          <el-input v-model="wechat.accessToken" style="width: 100%" />
        </el-form-item>
        <el-form-item label="Operation">
          <el-button
            v-for="yw in wechatAPI.ywzx"
            @click="doRequest(yw.method, yw.url, '')"
            >{{ yw.name }}</el-button
          >
        </el-form-item>
      </el-form>
    </el-tab-pane>
    <el-tab-pane label="DingDing"> </el-tab-pane>
  </el-tabs>
  <div style="position: relative; margin-top: 5px;">
    <el-tabs>
      <el-tab-pane label="Console">
        <pre
          class="pretty-response"
          style="height: 40vh"
        ><code>{{ result }}</code></pre>
      </el-tab-pane>
    </el-tabs>
      <el-button class="custom_eltabs_titlebar" :icon="Delete" @click="result = ''" />
  </div>
</template>

<style>
.el-textarea__inner {
  height: 100%;
}
</style>
