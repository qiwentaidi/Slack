import * as XLSX from "xlsx";
import { ElNotification } from "element-plus";
import global from "./global";
import { CheckTarget, SaveFile, GoFetch, Sock5Connect } from "../wailsjs/go/main/App";
import { GetFileContent, WriteFile } from "../wailsjs/go/main/File";
import Loading from "./components/Loading.vue";
// 单sheet导出
export async function ExportToXlsx(
  headers: string[],
  sheetName: string,
  filename: string,
  result: {}[]
) {
  let wb = XLSX.utils.book_new();
  let ws = XLSX.utils.json_to_sheet(result);
  XLSX.utils.sheet_add_aoa(ws, [headers], { origin: "A1" });
  // 将工作表添加到工作簿中
  XLSX.utils.book_append_sheet(wb, ws, sheetName);
  // 将工作簿写入到一个新的Excel文件中
  const b64 = XLSX.write(wb, { bookType: "xlsx", type: "base64" });
  await ExportFile("base64", filename + ".xlsx", b64);
}

// 资产导出
export async function ExportAssetToXlsx(r1: {}[], r2: {}[], r3: {}[]) {
  // 创建一个新的工作簿
  let wb = XLSX.utils.book_new();
  // 自定义表头
  let suheaders = ["公司名称", "股权比例", "投资数额", "域名"];
  let weheaders = ["公众号名称", "微信号", "Logo", "二维码", "简介"];
  let huheaders = ["公司域名或ICP名称", "资产数量"];
  let ws1 = XLSX.utils.json_to_sheet(r1);
  let ws2 = XLSX.utils.json_to_sheet(r2);
  let ws3 = XLSX.utils.json_to_sheet(r3);
  XLSX.utils.sheet_add_aoa(ws1, [suheaders], { origin: "A1" });
  XLSX.utils.sheet_add_aoa(ws2, [weheaders], { origin: "A1" });
  XLSX.utils.sheet_add_aoa(ws3, [huheaders], { origin: "A1" });
  XLSX.utils.book_append_sheet(wb, ws1, "子公司");
  XLSX.utils.book_append_sheet(wb, ws2, "公众号");
  XLSX.utils.book_append_sheet(wb, ws3, "鹰图资产数量");
  const b64 = XLSX.write(wb, { bookType: "xlsx", type: "base64" });
  await ExportFile("base64", "asset.xlsx", b64);
}

export async function ExportTXT(filename: string, result: string[]) {
  //文件内容
  var text = "";
  for (const item of result) {
    text += item + "\n";
  }
  await ExportFile("txt", filename + ".txt", text);
}

export async function ExportFile(
  filetype: string,
  filename: string,
  content: string
) {
  const path = await SaveFile(filename);
  if (path == "") {
    return;
  }
  const result = await WriteFile(filetype, path, content);
  if (result) {
    ElNotification({
      message: "数据保存成功，路径为:" + path,
      type: "success",
    });
  } else {
    ElNotification({
      message: "数据导出失败!",
      type: "warning",
    });
  }
}

// 复制端口扫描中的所有HTTP链接
export function CopyURLs(result: {}[]) {
  // 避免控制台报错
  if (result.length <= 1) {
    ElNotification({
      message: "复制内容条数需大于1",
      type: "warning",
      position: 'bottom-right',
    });
    return;
  }
  const temp = [];
  for (const line of result) {
    if ((line as any)["link"].includes("http")) {
      temp.push((line as any)["link"]);
    }
  }
  Copy(temp.join("\n"))
}

export function Copy(content: string) {
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
  if (filed.length == 0) {
    ElNotification({
      message: "Copy data can't be empty",
      type: "warning",
      position: 'bottom-right',
    });
    return
  }
  Copy(filed.join("\n"))
}

export function SplitTextArea(textarea: string) {
  let lines = textarea.split(/[(\r\n)\r\n]+/); // 根据换行或者回车进行识别
  lines = lines.filter(item => item.trim() !== ''); // 删除空项并去除左右空格
  lines = Array.from(new Set(lines)); // 去重
  return lines;
}

export function LoadConfig() {
  const savedScan = localStorage.getItem("scan");
  const savedProxy = localStorage.getItem("proxy");
  const savedSpace = localStorage.getItem("space");
  if (savedScan) {
    Object.assign(global.scan, JSON.parse(savedScan));
  }

  if (savedProxy) {
    Object.assign(global.proxy, JSON.parse(savedProxy));
  }

  if (savedSpace) {
    Object.assign(global.space, JSON.parse(savedSpace));
  }
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
  let temp: Array<string> = [];
  let urls: Array<string> = [];
  for (var target of SplitTextArea(host)) {
    if (target.slice(0, 4) !== "http") {
      const result = await CheckTarget(host, global.proxy);
      if (result.Status === true) {
        target = result.ProtocolURL;
      }
    }
    temp.push(target);
  }
  for (var item of temp) {
    const urlObj = new URL(item)
    // 存在路径就不交验/
    if (urlObj.pathname.length > 1) {
      urls.push(item);
    } else {
      if (item.slice(-1) !== "/") {
        urls.push((item += "/"));
      } else {
        urls.push(item);
      }
    }
  }
  return urls;
}

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
  if (res !== "文件不存在") {
    const result = res.replace(/\r\n/g, "\n"); // 避免windows unix系统差异
    return Array.from(result.split("\n"));
  }
}

// mode 0 is button click
export async function TestProxy(mode: number) {
  const proxyURL = global.proxy.mode.toLowerCase() + "://" + global.proxy.address + ":" + global.proxy.port
  if (global.proxy.enabled) {
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

export function currentTime() {
  var date = new Date();
  return date.toLocaleString()
}