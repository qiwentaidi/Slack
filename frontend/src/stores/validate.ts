import { ProcessTextAreaInput } from "@/util";
import { ElMessage } from "element-plus";

export function validatePortscan(input: string) {
    const ipPatterns = [
        /^(\d{1,3}\.){3}\d{1,3}$/, // 192.168.1.1
        /^(\d{1,3}\.){3}\d{1,3}\/(\d{1,2})$/, // 192.168.1.1/8, 192.168.1.1/16, 192.168.1.1/24
        /^(\d{1,3}\.){3}\d{1,3}-(\d{1,3}\.){3}\d{1,3}|(\d{1,3}\.){2}\d{1,3}|\d{1,3}\)$/, // 192.168.1.1-192.168.255.255, 192.168.1.1-255
        /^(\d{1,3}\.){3}\d{1,3}:\d{1,5}$/, // 192.168.0.1:6379
        /^!((\d{1,3}\.){3}\d{1,3}(\/\d+)?|(\d{1,3}\.){2}\d{1,3}|\d{1,3})$/, // !192.168.1.6/28
        /^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?$/, // domain
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
        /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/;
    return regex.test(ip);
}

export function validateSingleDomain(domain: string) {
    const regex = /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/;
    return regex.test(domain);
}

export function isPrivateIP(ip: string) {
    const regex = /^(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(172\.(1[6-9]|2[0-9]|3[0-1])\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})|(127\.\d{1,3}\.\d{1,3}\.\d{1,3})/;
    return regex.test(ip);
}

export function validateWebscan(input: string) {
    const ipPatterns = /^[a-zA-Z0-9\-]+.[a-zA-Z0-9\-]+/ // 符合域名规范即可
    var lines = ProcessTextAreaInput(input)
    lines = lines.filter(line => line.trim() !== '');
    for (const line of lines) {
        if (!ipPatterns.test(line)) {
            ElMessage.warning(line + " 输入格式错误")
            return false
        }
    }
    return true
}

export const regexpPhone = /(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}/g;
export const regexpIdCard = /[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[012])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]/g;
export const regexpAKSK = /access/g;
