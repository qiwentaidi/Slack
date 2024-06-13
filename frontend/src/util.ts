
import { ElMessage, ElNotification } from "element-plus";
import global from "./global";
import { CheckTarget, GoFetch, NucleiEnabled, Sock5Connect } from "../wailsjs/go/main/App";
import { CheckFileStat, GetFileContent, UserHomeDir } from "../wailsjs/go/main/File";
import Loading from "./components/Loading.vue";

export function Copy(content: string) {
  if (content == "") {
    ElNotification({
      message: "Copy data can't be empty",
      type: "warning",
      position: 'bottom-right',
    });
    return;
  }
  navigator.clipboard.writeText(content).then(
    function () {
      ElNotification({
        message: "Copy Finished",
        type: "success",
        position: 'bottom-right',
      });
    },
    function (err) {
      ElNotification({
        message: "Copy Failed: " + err,
        type: "error",
        position: 'bottom-right',
      });
    }
  );
}

export function CopyALL(filed: string[]) {
  Copy(filed.join("\n"))
}

export function SplitTextArea(textarea: string) {
  let lines = textarea.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
  lines = lines.filter(item => item.trim() !== ''); // 删除空项并去除左右空格
  lines = Array.from(new Set(lines)); // 去重
  return lines;
}

export function validateIP(ip: string): boolean {
  const regex =
    /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  return regex.test(ip);
}

export function validateDomain(domain: string): boolean {
  const regex = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
  return regex.test(domain);
}

export function validateURL(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch (e) {
    ElMessage("请输入正确的URL")
    return false
  }
}

export function splitInt(n: number, slice: number): number[] {
  let res: number[] = [];
  while (n > slice) {
    res.push(slice);
    n = n - slice;
  }
  res.push(n);
  return res;
}

export async function formatURL(host: string): Promise<string[]> {
  let urls: Array<string> = [];
  for (var target of SplitTextArea(host)) {
    if (!target.startsWith("http")) {
      const result = await CheckTarget(host, global.proxy);
      if (result.Status === true) {
        target = result.ProtocolURL;
      }
    }
    urls.push(AddRightSubString(target, "/"))
  }
  return urls;
}

function AddRightSubString(str: string, sub: string) {
  if (str.endsWith(sub)) {
    return str
  }
  return str + sub
}

// version1 > version2 return 1
export function compareVersion(version1: string, version2: string) {
  const v1 = version1.split(".").map(Number);
  const v2 = version2.split(".").map(Number);

  for (let i = 0; i < v1.length || i < v2.length; i++) {
    const num1 = v1[i] || 0;
    const num2 = v2[i] || 0;

    if (num1 > num2) {
      return 1;
    }
    if (num1 < num2) {
      return -1;
    }
  }

  return 0;
}

export interface TableTabs {
  title: string;
  name: string;
  content: null | [{}];
  total: number;
  pageSize: number;
  currentPage: number;
}

let RegCompliance = new RegExp('(\\w+)[!,=]{1,3}"([^"]+)"');

// mode 0 fofa , mode 1 hunter
export function ApiSyntaxCheck(
  mode: number,
  email: string,
  key: string,
  query: string
) {
  if (mode == 0) {
    if (email == "" || key == "") {
      ElNotification({
        message: "请在设置处填写FOFA Email和Key",
        type: "warning",
      });
      return false;
    }
  } else {
    if (key == "") {
      ElNotification({
        message: "请在设置处填写Hunter Key",
        type: "warning",
      });
      return false;
    }
  }
  if (RegCompliance.test(query) === false) {
    ElNotification({
      message: "Syntax error",
      type: "warning",
    });
    return false;
  }
  return true;
}

export async function ReadLine(filepath: string) {
  let res = await GetFileContent(filepath);
  if (res.length == 0) {
    ElNotification({
      message: "Reading file cannot be empty",
      type: "warning",
    });
    return;
  }
  if (res !== "") {
    const result = res.replace(/\r\n/g, "\n"); // 避免windows unix系统差异
    return Array.from(result.split("\n"));
  }
}

// mode 0 is button click
export async function TestProxy(mode: number) {
  if (global.proxy.enabled) {
    const proxyURL = global.proxy.mode.toLowerCase() + "://" + global.proxy.address + ":" + global.proxy.port
    if (global.proxy.mode == "HTTP") {
      let resp = await GoFetch("GET", proxyURL, "", [{}], 10, null)
      if (resp.Error == true) {
        ElNotification({
          message: "The proxy is unreachable",
          type: "warning",
        });
        return false
      }
      if (mode == 0) {
        ElNotification({
          message: "The proxy is enabled",
          type: "success",
        });
      }
    } else {
      ElNotification({
        duration: 0,
        message: 'Connecting to http://www.baidu.com',
        icon: Loading,
      });
      let resp = await Sock5Connect(global.proxy.address, global.proxy.port, 10, global.proxy.username, global.proxy.password)
      if (!resp) {
        ElNotification.closeAll()
        ElNotification({
          message: "Connect to http://www.baidu.com is timeout",
          type: "error",
        });
        return false
      }
      ElNotification.closeAll()
      ElNotification({
        title: 'Success',
        message: "Connect to http://www.baidu.com is success",
        type: "success",
      });
    }
  }
  return true
}

export async function TestNuclei() {
  NucleiEnabled(global.webscan.nucleiEngine).then(result => {
    if (result) {
      ElNotification({
        message: "Nuclei engine is enabled",
        type: "success",
        duration: 2000,
      });
    } else {
      ElNotification({
        type: 'error',
        message: "Nuclei engine is disable",
        duration: 2000,
      });
    }
  })
}

export function currentTimestamp() {
  const now = new Date();
  return Math.floor(now.getTime() / 1000).toString();
}

const download = {
  RemotePocVersion:
    "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version",
  RemoteClientVersion:
    "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version",
  PocUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update',
  ClientUpdateCentent: 'https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update',
};

export const check = ({
  // poc
  poc: async function () {
    let pcfg = await CheckFileStat(await UserHomeDir() + global.PATH.LocalPocVersionFile)
    if (!pcfg) {
      global.UPDATE.LocalPocVersion = "版本文件不存在"
      global.UPDATE.PocStatus = false
      return
    } else {
      global.UPDATE.LocalPocVersion = await GetFileContent(await UserHomeDir() + global.PATH.LocalPocVersionFile)
    }
    let resp = await GoFetch("GET", download.RemotePocVersion, "", [{}], 10, null)
    if (resp.Error == true) {
      global.UPDATE.PocContent = "检测更新失败"
      global.UPDATE.PocStatus = false
    } else {
      global.UPDATE.RemotePocVersion = resp.Body
      if (compareVersion(global.UPDATE.LocalPocVersion, global.UPDATE.RemotePocVersion) == -1) {
        global.UPDATE.PocContent = (await GoFetch("GET", download.PocUpdateCentent, "", [{}], 10, null)).Body
        global.UPDATE.PocStatus = true
      } else {
        global.UPDATE.PocContent = "已是最新版本"
        global.UPDATE.PocStatus = false
      }
    }
  },
  // client
  client: async function () {
    let resp = await GoFetch("GET", download.RemoteClientVersion, "", [{}], 10, null)
    if (resp.Error) {
      global.UPDATE.RemoteClientVersion = "检测更新失败"
      global.UPDATE.ClientStatus = false
    } else {
      global.UPDATE.RemoteClientVersion = resp.Body
      if (compareVersion(global.LOCAL_VERSION, global.UPDATE.RemoteClientVersion) == -1) {
        global.UPDATE.ClientContent = (await GoFetch("GET", download.ClientUpdateCentent, "", [{}], 10, null)).Body
        global.UPDATE.ClientStatus = true
      } else {
        global.UPDATE.ClientContent = "已是最新版本"
        global.UPDATE.ClientStatus = false
      }
    }
  }
})

export function sleep(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time));
}
