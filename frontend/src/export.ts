import * as XLSX from "xlsx";
import { ElNotification } from "element-plus";
import { WriteFile, SaveFileDialog } from "wailsjs/go/services/File";
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
export async function ExportAssetToXlsx(r1: {}[], r2: {}[]) {
    // 创建一个新的工作簿
    let wb = XLSX.utils.book_new();
    // 自定义表头
    let suheaders = ["公司名称", "股权比例", "投资数额", "企业状态", "域名"];
    let weheaders = ["公司名称", "公众号名称", "微信号", "Logo", "二维码", "简介"];
    let ws1 = XLSX.utils.json_to_sheet(r1);
    let ws2 = XLSX.utils.json_to_sheet(r2);
    XLSX.utils.sheet_add_aoa(ws1, [suheaders], { origin: "A1" });
    XLSX.utils.sheet_add_aoa(ws2, [weheaders], { origin: "A1" });
    XLSX.utils.book_append_sheet(wb, ws1, "子公司");
    XLSX.utils.book_append_sheet(wb, ws2, "公众号");
    const b64 = XLSX.write(wb, { bookType: "xlsx", type: "base64" });
    await ExportFile("base64", "asset.xlsx", b64);
}

// export async function ExportWebScanToXlsx(r1: {}[], r2: {}[]) {
//   // 创建一个新的工作簿
//   let wb = XLSX.utils.book_new();
//   // 自定义表头
//   let fingerheaders = ["URL", "Code", "Length", "Title", "Detection", "isWAF", "WAF Name", "Fingerprint"];
//   let vulheaders = ["Template", "Name", "Type" ,"Severity", "URL", "ExtInfo"];
//   let ws1 = XLSX.utils.json_to_sheet(r1);
//   let ws2 = XLSX.utils.json_to_sheet(r2);
//   XLSX.utils.sheet_add_aoa(ws1, [fingerheaders], { origin: "A1" });
//   XLSX.utils.sheet_add_aoa(ws2, [vulheaders], { origin: "A1" });
//   XLSX.utils.book_append_sheet(wb, ws1, "指纹");
//   XLSX.utils.book_append_sheet(wb, ws2, "漏洞");
//   const b64 = XLSX.write(wb, { bookType: "xlsx", type: "base64" });
//   await ExportFile("base64", "webscan.xlsx", b64);
// }


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
    const path = await SaveFileDialog(filename);
    if (!path) {
        return;
    }
    const result = await WriteFile(filetype, path, content);
    if (result) {
        ElNotification.success("数据保存成功，路径为:" + path);
    } else {
        ElNotification.warning("数据导出失败!");
    }
}
