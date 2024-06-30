
import { ElMessage, ElNotification } from "element-plus";
import global from "./global";
import { CheckTarget, GoFetch, NucleiEnabled, Sock5Connect } from "../wailsjs/go/main/App";
import { CheckFileStat, ReadFile, UserHomeDir } from "../wailsjs/go/main/File";
import Loading from "./components/Loading.vue";
import { URLFingerMap, ProxyOptions, File } from "./interface";

export var proxys: ProxyOptions // wails2.9之后替换原来的null

export function Copy(content: string) {
  if (content == "") {
    ElNotification.warning({
      message: "Copy data can't be empty",
      position: 'bottom-right',
    });
    return;
  }
  navigator.clipboard.writeText(content).then(
    function () {
      ElNotification.success({
        message: "Copy Finished",
        position: 'bottom-right',
      });
    },
    function (err) {
      ElNotification.error({
        message: "Copy Failed: " + err,
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
    ElMessage.warning("请输入正确的URL")
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
      const result: any = await CheckTarget(host, global.proxy)
      if (!result.Error) {
        target = result.Msg;
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

// hunter
export function ApiSyntaxCheck(
  key: string,
) {
  if (key == "") {
    ElNotification.warning("请在设置处填写Hunter Key");
    return false;
  }
  return true;
}

export async function ReadLine(filepath: string) {
  let file: File = await ReadFile(filepath);
  if (file.Error) {
    ElNotification.warning(file.Message);
    return;
  }
  const result = file.Content!.replace(/\r\n/g, "\n"); // 避免windows unix系统差异
  return Array.from(result.split("\n"));
}

// mode 0 is button click
export async function TestProxy(mode: number) {
  if (global.proxy.enabled) {
    const proxyURL = global.proxy.mode.toLowerCase() + "://" + global.proxy.address + ":" + global.proxy.port
    if (global.proxy.mode == "HTTP") {
      let resp: any = await GoFetch("GET", proxyURL, "", [{}], 10, proxys)
      if (resp.Error) {
        ElNotification.warning("The proxy is unreachable");
        return false
      }
      if (mode == 0) {
        ElNotification.success("The proxy is enabled");
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
        ElNotification.error("The sock5 proxy is unreachable");
        return false
      }
      ElNotification.closeAll()
      ElNotification.success("Connect to http://www.baidu.com is success");
    }
  }
  return true
}

export async function TestNuclei() {
  NucleiEnabled(global.webscan.nucleiEngine).then(result => {
    if (result) {
      ElNotification.success("Nuclei engine is enabled");
    } else {
      ElNotification.error("Nuclei engine is disable");
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
      let file: File = await ReadFile(await UserHomeDir() + global.PATH.LocalPocVersionFile)
      global.UPDATE.LocalPocVersion = file.Content!
    }
    let resp: any = await GoFetch("GET", download.RemotePocVersion, "", [{}], 10, proxys)
    if (resp.Error == true) {
      global.UPDATE.PocContent = "检测更新失败"
      global.UPDATE.PocStatus = false
    } else {
      global.UPDATE.RemotePocVersion = resp.Body!
      if (compareVersion(global.UPDATE.LocalPocVersion, global.UPDATE.RemotePocVersion) == -1) {
        let result: any = GoFetch("GET", download.PocUpdateCentent, "", [{}], 10, proxys)
        global.UPDATE.PocContent = result.Body
        global.UPDATE.PocStatus = true
      } else {
        global.UPDATE.PocContent = "已是最新版本"
        global.UPDATE.PocStatus = false
      }
    }
  },
  // client
  client: async function () {
    let resp: any = await GoFetch("GET", download.RemoteClientVersion, "", [{}], 10, proxys)
    if (resp.Error) {
      global.UPDATE.RemoteClientVersion = "检测更新失败"
      global.UPDATE.ClientStatus = false
    } else {
      global.UPDATE.RemoteClientVersion = resp.Body!
      if (compareVersion(global.LOCAL_VERSION, global.UPDATE.RemoteClientVersion) == -1) {
        let result: any = await GoFetch("GET", download.ClientUpdateCentent, "", [{}], 10, proxys)
        console.log(result)
        global.UPDATE.ClientContent = result.Body!
        global.UPDATE.ClientStatus = true
      } else {
        global.UPDATE.ClientContent = "已是最新版本"
        global.UPDATE.ClientStatus = false
      }
    }
  }
})

export async function sleep(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time));
}


export function deduplicateUrlFingerMap(urlFingerMap: URLFingerMap[]): URLFingerMap[] {
  const seenUrls = new Set<string>();
  return urlFingerMap.filter(item => {
    if (seenUrls.has(item.url)) {
      return false;
    } else {
      seenUrls.add(item.url);
      return true;
    }
  });
}

export function generateRandomString(length: number) :string {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}