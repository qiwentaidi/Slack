import { SplitTextArea } from "@/util";
import { ElMessage } from "element-plus";

export function validatePortscan(input: string) {
    const ipPatterns = [
        /^(\d{1,3}\.){3}\d{1,3}$/, // 192.168.1.1
        /^(\d{1,3}\.){3}\d{1,3}\/(\d{1,2})$/, // 192.168.1.1/8, 192.168.1.1/16, 192.168.1.1/24
        /^(\d{1,3}\.){3}\d{1,3}-(\d{1,3}\.){3}\d{1,3}|(\d{1,3}\.){2}\d{1,3}|\d{1,3}\)$/, // 192.168.1.1-192.168.255.255, 192.168.1.1-255
        /^(\d{1,3}\.){3}\d{1,3}:\d{1,5}$/, // 192.168.0.1:6379
        /^!((\d{1,3}\.){3}\d{1,3}(\/\d+)?|(\d{1,3}\.){2}\d{1,3}|\d{1,3})$/, // !192.168.1.6/28
        /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/, // domain
        /^\s*$/ // 空行
    ];
    const lines = input.split('\n');
    return lines.every(line =>
        ipPatterns.some(pattern => pattern.test(line.trim()))
    );
}

export function validateIp(input: string) {
    const ipPatterns = [
        /^(\d{1,3}\.){3}\d{1,3}$/, // 192.168.1.1
        /^(\d{1,3}\.){3}\d{1,3}\/(\d{1,2})$/, // 192.168.1.1/8, 192.168.1.1/16, 192.168.1.1/24
        /^(\d{1,3}\.){3}\d{1,3}-(\d{1,3}\.){3}\d{1,3}|(\d{1,3}\.){2}\d{1,3}|\d{1,3}\)$/, // 192.168.1.1-192.168.255.255, 192.168.1.1-255
        /^!((\d{1,3}\.){3}\d{1,3}(\/\d+)?|(\d{1,3}\.){2}\d{1,3}|\d{1,3})$/, // !192.168.1.6/28
    ];
    const lines = input.split('\n');
    return lines.every(line =>
        ipPatterns.some(pattern => pattern.test(line.trim()))
    );
}

export function validateSingleIP(ip: string) {
    const regex =
        /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return regex.test(ip);
}

export function validateSingleDomain(domain: string) {
    const regex = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
    return regex.test(domain);
}

export function validateSingleURL(url: string): boolean {
    if (!url) {
        ElMessage.warning("URL不能为空");
        return false;
    }
    try {
        new URL(url);
        return true;
    } catch (e) {
        ElMessage.warning("请输入正确的URL");
        return false;
    }
}

export function isPrivateIP(ip: string) {
    const regex = /^(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(172\.(1[6-9]|2[0-9]|3[0-1])\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})|(127\.\d{1,3}\.\d{1,3}\.\d{1,3})/;
    return regex.test(ip);
}

export function validateWebscan(input: string) {
    const ipPatterns = /^[a-zA-Z0-9\-]+.[a-zA-Z0-9\-]+/ // 符合域名规范即可
    var lines = SplitTextArea(input)
    lines = lines.filter(line => line.trim() !== '');
    for (const line of lines) {
        if (!ipPatterns.test(line)) {
            ElMessage.warning(line + " 输入格式错误")
            return false
        }
    }
    return true
}

export const regexpPhone = /(?<![A-Za-z0-9])((?:(?:\+|00)86)?1(?:(?:3\d)|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8\d)|(?:9[189]))\d{8})(?![A-Za-z0-9])/g
export const regexpIdCard = /[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[012])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]/g;
export const regexpAKSK = /access/g;