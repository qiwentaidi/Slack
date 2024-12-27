import { ElMessage, ElNotification } from "element-plus";
import global from "./global";
import { Callgologger, CheckTarget, GoFetch, NetDial, Socks5Conn } from "wailsjs/go/services/App";
import { CheckFileStat, FileDialog, ReadFile, RemoveOldClient } from "wailsjs/go/services/File";
import Loading from "./components/Loading.vue";
import { ProxyOptions } from "./stores/interface";
import { ClipboardSetText } from "wailsjs/runtime/runtime";
import { marked } from 'marked';
import async from "async";
import { clients } from "wailsjs/go/models";

export var proxys: ProxyOptions; // wails2.9之后替换原来的null

// 惰性函数，如果支持navigator.clipboard就用它，不支持就改成另一个
let copyText: (content: string) => void;

if (navigator.clipboard) {
  copyText = (content: string) => {
    navigator.clipboard.writeText(content);
    ElNotification.success({
      message: "Copy Finished",
      position: "bottom-right",
    });
  };
} else {
  copyText = (content: string) => {
    ClipboardSetText(content).then((result) => {
      if (result) {
        ElNotification.success({
          message: "Copy Finished",
          position: "bottom-right",
        });
      } else {
        ElNotification.error({
          message: "Copy Failed!",
          position: "bottom-right",
        });
      }
    });
  };
}

export function Copy(content: string) {
  if (isEmpty(content)) {
    return
  }
  copyText(content);
}

function isEmpty(obj: string) {
  if (!obj) {
    ElNotification.warning({
      message: "Copy data can't be empty",
      position: "bottom-right",
    });
    return true;
  }
  return false;
}

export function CopyALL(filed: string[]) {
  Copy(filed.join("\n"));
}

// Function to process the input target in the text area
// 处理文本域的目标输入
export function ProcessTextAreaInput(textarea: string) {
  let lines = textarea.split(/[(\r\n)\r\n]+/).map(line => line.trim()); // Recognize and remove leading and trailing spaces based on line breaks or carriage returns
  lines = lines.filter(item => item !== ""); // Remove empty items
  return Array.from(new Set(lines));;
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

export async function FormatWebURL(host: string): Promise<string[]> {
  let urls: Array<string> = [];
  const targets = ProcessTextAreaInput(host);
  Callgologger("info", "正在解析目标 ...")
  // Wrap the asynchronous processing of each target with async.eachLimit
  await new Promise((resolve, reject) => {
    async.eachLimit(targets, 50, async (target: string, callback: () => void) => {
      if (!target.startsWith("http")) {
        const result: any = await CheckTarget(target, getProxy());
        if (!result.Error) {
          target = result.Msg;
        }
      }
      urls.push(TrimRightSubString(target, "/"));
    },
      (err: any) => {
        if (err) {
          reject(err); // Reject the promise if any task encounters an error
        } else {
          resolve(urls); // Resolve the promise once all tasks are completed
        }
      }
    );
  });
  return urls;
}

export function TrimRightSubString(str: string, sub: string) {
  if (str.endsWith(sub)) {
    return str.substring(0, str.length - sub.length);
  }
  return str;
}

export function AddRightSubString(str: string, sub: string) {
  if (str.endsWith(sub)) {
    return str;
  }
  return str + sub;
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

export async function ReadLine(filepath: string) {
  let file = await ReadFile(filepath);
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
    const host = global.proxy.address + ":" + global.proxy.port;
    if (global.proxy.mode == "HTTP") {
      let error = await NetDial(host)
      if (!error) {
        ElNotification.warning("The proxy is unreachable");
        return false;
      }
      if (mode == 0) {
        ElNotification.success("The proxy is enabled");
      }
    } else {
      ElNotification({
        duration: 0,
        message: "Connecting to http://www.baidu.com",
        icon: Loading,
      });
      let error = await Socks5Conn(global.proxy.address, global.proxy.port, 10, global.proxy.username, global.proxy.password);
      if (!error) {
        ElNotification.closeAll();
        ElNotification.error("The socks5 proxy is unreachable");
        return false;
      }
      ElNotification.closeAll();
      ElNotification.success("The socks5 proxy is enabled");
    }
  }
  return true;
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
  PocUpdateCentent:
    "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update",
  ClientUpdateCentent:
    "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update",
};

export const check = {
  // poc
  poc: async function () {
    let pcfg = await CheckFileStat(global.PATH.homedir + global.PATH.LocalPocVersionFile);
    if (!pcfg) {
      global.UPDATE.LocalPocVersion = "版本文件不存在";
      global.UPDATE.PocStatus = false;
      return;
    } else {
      let file = await ReadFile(global.PATH.homedir + global.PATH.LocalPocVersionFile);
      global.UPDATE.LocalPocVersion = file.Content!;
    }
    let resp = await GoFetch("GET", download.RemotePocVersion, "", {}, 10, proxys);
    if (resp.Error == true) {
      global.UPDATE.PocContent = "检测更新失败";
      global.UPDATE.PocStatus = false;
    } else {
      global.UPDATE.RemotePocVersion = resp.Body!;
      if (compareVersion(global.UPDATE.LocalPocVersion, global.UPDATE.RemotePocVersion) == -1) {
        let result = await GoFetch("GET", download.PocUpdateCentent, "", {}, 10, proxys);
        global.UPDATE.PocContent = result.Body;
        global.UPDATE.PocStatus = true;
      } else {
        global.UPDATE.PocContent = "已是最新版本";
        global.UPDATE.PocStatus = false;
      }
    }
  },
  // client
  client: async function () {
    let resp = await GoFetch("GET", download.RemoteClientVersion, "", {}, 10, proxys);
    if (resp.Error) {
      global.UPDATE.RemoteClientVersion = "检测更新失败";
      global.UPDATE.ClientStatus = false;
    } else {
      global.UPDATE.RemoteClientVersion = resp.Body!;
      if (
        compareVersion(
          global.LOCAL_VERSION,
          global.UPDATE.RemoteClientVersion
        ) == -1
      ) {
        let result: any = await GoFetch("GET", download.ClientUpdateCentent, "", {}, 10, proxys);
        global.UPDATE.ClientContent = result.Body!;
        global.UPDATE.ClientStatus = true;
      } else {
        global.UPDATE.ClientContent = "已是最新版本";
        RemoveOldClient();
        global.UPDATE.ClientStatus = false;
      }
    }
  },
};

export async function sleep(time: number) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

export function generateRandomString(length: number): string {
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}


type AnyObject = { [key: string]: any };

// 将对象中的数组字段转换为自定义拼接符的字符串
export function transformArrayFields<T extends AnyObject>(data: T[], delimiter: string = '|'): T[] {
  return data.map(item => {
    const transformedItem: AnyObject = {};
    for (const key in item) {
      if (item.hasOwnProperty(key)) {
        if (Array.isArray(item[key])) {
          transformedItem[key] = item[key].map((subItem: any) => JSON.stringify(subItem)).join(delimiter);
        } else {
          transformedItem[key] = item[key];
        }
      }
    }
    return transformedItem as T;
  });
}

export function CsegmentIpv4(ip: string): string {
  return ip.split('.').slice(0, 3).join('.') + ".0/24";
}

export function renderedMarkdown(content: string) {
  return marked.parse(content);
}

export async function UploadFileAndRead() {
  let filepath = await FileDialog("*.txt")
  if (!filepath) {
    return ""
  }
  let file = await ReadFile(filepath)
  if (file.Error) {
    ElMessage({
      type: "warning",
      message: file.Message
    })
    return ""
  }
  return file.Content!
}

export function getBasicURL(rawURL: string) {
  try {
    const url = new URL(rawURL);
    return `${url.protocol}//${url.host}`;
  } catch (error) {
    console.error("Invalid URL:", error);
    return undefined;
  }
}

export function getProxy(): clients.Proxy {
  return {
    Enabled: global.proxy.enabled,
    Mode: global.proxy.mode,
    Address: global.proxy.address,
    Port: global.proxy.port,
    Username: global.proxy.username,
    Password: global.proxy.password,
  }
}

export function getRootDomain(hostname: string) {
  var parts = hostname.split('.'); // 拆分域名部分
  if (parts.length > 2) {
    // 如果域名有子域名（例如：sub.example.com）
    return parts.slice(-2).join('.'); // 返回最后两部分（例如：example.com）
  }
  return hostname; // 如果没有子域名，直接返回（例如：example.com）
}

// 处理 {{file:///xxxx}} 获取到文件路径
// export function getFilepath(raw: string) {
//   const startTag = "{{file://";
//   const endTag = "}}";

//   const startIndex = raw.indexOf(startTag);
//   const endIndex = raw.indexOf(endTag, startIndex);

//   if (startIndex !== -1 && endIndex !== -1) {
//     return raw.substring(startIndex + startTag.length, endIndex);
//   }
//   return ""
// }