import CryptoJS from "crypto-js";
import JSEncrypt from "jsencrypt";
import sm from "sm-crypto";
import { IpLocation } from "wailsjs/go/services/App";

export type OperationMode = "encode" | "decode" | "transform";

type MaybePromise<T> = T | Promise<T>;

export interface CodecOperationOption {
    key: string;
    label: string;
    default: string;
    type?: "input" | "select";
    options?: string[];
}

export interface CodecOperation {
    id: string;
    name: string;
    category: string;
    description: string;
    encode?: (input: string, options?: Record<string, string>) => MaybePromise<string>;
    decode?: (input: string, options?: Record<string, string>) => MaybePromise<string>;
    transform?: (input: string, options?: Record<string, string>) => MaybePromise<string>;
    options?: CodecOperationOption[];
}

export interface CodecCategory {
    id: string;
    name: string;
}

export interface AddedCodecOperation {
    id: string;
    operationId: string;
    mode: OperationMode;
    options: Record<string, string>;
}

const regexpPhone = /(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}/g;
const regexpIdCard = /[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]/g;
const regexURL = /https?:\/\/[^\s"'<>]+/g;
const regexIP = /\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b/g;
const regexDruidWebSession = /"SESSIONID":"([^"]+)"/g;
const regexDruidWebURI = /"URI":"([^"]+)"/g;
const PASSWORD_MASK_ARRAY = [19, 78, 10, 15, 100, 213, 43, 23];

const BASE32_ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";

const splitBySeparator = (input: string, separator?: string) => {
    let actualSeparator: string | RegExp;
    let outputSeparator: string;
    if (!separator || separator === "\\n") {
        actualSeparator = /\r?\n/;
        outputSeparator = "\n";
    } else if (separator === "\\t") {
        actualSeparator = "\t";
        outputSeparator = "\t";
    } else if (separator === "\\r\\n") {
        actualSeparator = /\r\n/;
        outputSeparator = "\n";
    } else {
        actualSeparator = separator;
        outputSeparator = separator;
    }
    return {
        parts: input.split(actualSeparator),
        outputSeparator,
    };
};

const validateSingleIP = (ip: string): boolean => {
    const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
    if (!ipRegex.test(ip)) return false;
    const parts = ip.split(".").map(Number);
    return parts.every((part) => part >= 0 && part <= 255);
};

const parseKeyOrIv = (value: string, format: string, targetLength: number) => {
    switch (format.toLowerCase()) {
        case "hex":
            return CryptoJS.enc.Hex.parse(value);
        case "base64":
            return CryptoJS.enc.Base64.parse(value);
        case "utf8":
        default:
            return CryptoJS.enc.Utf8.parse(value.padEnd(targetLength, "0").slice(0, targetLength));
    }
};

const getCryptoMode = (mode: string) => {
    const modes: Record<string, typeof CryptoJS.mode.CBC> = {
        CBC: CryptoJS.mode.CBC,
        ECB: CryptoJS.mode.ECB,
        CFB: CryptoJS.mode.CFB,
        OFB: CryptoJS.mode.OFB,
        CTR: CryptoJS.mode.CTR,
    };
    return modes[mode.toUpperCase()] || CryptoJS.mode.CBC;
};

const getCryptoPadding = (padding: string) => {
    const paddings: Record<string, typeof CryptoJS.pad.Pkcs7> = {
        Pkcs7: CryptoJS.pad.Pkcs7,
        Pkcs5: CryptoJS.pad.Pkcs7,
        ZeroPadding: CryptoJS.pad.ZeroPadding,
        NoPadding: CryptoJS.pad.NoPadding,
        Iso10126: CryptoJS.pad.Iso10126,
        Iso97971: CryptoJS.pad.Iso97971,
        AnsiX923: CryptoJS.pad.AnsiX923,
    };
    return paddings[padding] || CryptoJS.pad.Pkcs7;
};

const textToHex = (value: string, bytes: number) =>
    Array.from(value.padEnd(bytes, "\0").slice(0, bytes))
        .map((char) => char.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("");

const hexToBase64 = (hex: string) => {
    const cleaned = hex.replace(/\s+/g, "");
    let binary = "";
    for (let i = 0; i < cleaned.length; i += 2) {
        binary += String.fromCharCode(parseInt(cleaned.slice(i, i + 2), 16));
    }
    return btoa(binary);
};

const base64ToHex = (value: string) =>
    Array.from(atob(value.trim()))
        .map((char) => char.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("");

const base64Encode = (input: string, options?: Record<string, string>): string => {
    try {
        const mode = (options?.mode || "Standard").toUpperCase();
        const bytes = new TextEncoder().encode(input);
        const base64 = btoa(String.fromCharCode(...bytes));
        if (mode === "URL") {
            return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
        }
        if (mode === "MIME") {
            let result = "";
            for (let i = 0; i < base64.length; i += 76) {
                result += base64.slice(i, i + 76) + "\n";
            }
            return result.trimEnd();
        }
        return base64;
    } catch {
        return "[编码错误: 无效输入]";
    }
};

const base64Decode = (input: string, options?: Record<string, string>): string => {
    try {
        const mode = (options?.mode || "Standard").toUpperCase();
        let base64 = input.trim();
        if (mode === "URL") {
            base64 = base64.replace(/-/g, "+").replace(/_/g, "/");
            const padding = base64.length % 4;
            if (padding) base64 += "=".repeat(4 - padding);
        } else if (mode === "MIME") {
            base64 = base64.replace(/\n/g, "");
        }
        const binary = atob(base64);
        const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
        return new TextDecoder().decode(bytes);
    } catch {
        return "[解码错误: 无效的 Base64]";
    }
};

const base32Encode = (input: string, options?: Record<string, string>): string => {
    try {
        const padding = (options?.padding || "True").toUpperCase() !== "FALSE";
        const bytes = new TextEncoder().encode(input);
        let bits = "";
        let result = "";
        for (const byte of bytes) {
            bits += byte.toString(2).padStart(8, "0");
        }
        for (let i = 0; i < bits.length; i += 5) {
            result += BASE32_ALPHABET[parseInt(bits.slice(i, i + 5).padEnd(5, "0"), 2)];
        }
        if (padding) {
            const padLength = (8 - (result.length % 8)) % 8;
            result += "=".repeat(padLength);
        }
        return result;
    } catch {
        return "[编码错误: 无效输入]";
    }
};

const base32Decode = (input: string): string => {
    try {
        const base32 = input.trim().toUpperCase().replace(/=+$/, "");
        let bits = "";
        for (const char of base32) {
            const index = BASE32_ALPHABET.indexOf(char);
            if (index === -1) return "[解码错误: 无效的 Base32 字符]";
            bits += index.toString(2).padStart(5, "0");
        }
        const bytes: number[] = [];
        for (let i = 0; i + 8 <= bits.length; i += 8) {
            bytes.push(parseInt(bits.slice(i, i + 8), 2));
        }
        return new TextDecoder().decode(Uint8Array.from(bytes));
    } catch {
        return "[解码错误: 无效的 Base32]";
    }
};

const urlEncode = (input: string) => {
    try {
        return encodeURIComponent(input);
    } catch {
        return "[编码错误]";
    }
};

const urlDecode = (input: string) => {
    try {
        return decodeURIComponent(input);
    } catch {
        return "[解码错误: 无效的 URL 编码]";
    }
};

const htmlEncode = (input: string) => {
    const map: Record<string, string> = {
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#39;",
    };
    return input.replace(/[&<>"']/g, (matched) => map[matched]);
};

const htmlDecode = (input: string) => {
    const map: Record<string, string> = {
        "&amp;": "&",
        "&lt;": "<",
        "&gt;": ">",
        "&quot;": '"',
        "&#39;": "'",
        "&nbsp;": " ",
    };
    return input.replace(/&(amp|lt|gt|quot|#39|nbsp);/g, (matched) => map[matched] || matched);
};

const hexEncode = (input: string) =>
    Array.from(input)
        .map((char) => char.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("");

const hexDecode = (input: string) => {
    try {
        const hex = input.replace(/\s/g, "");
        if (hex.length % 2 !== 0) return "[解码错误: 无效的 Hex 长度]";
        let result = "";
        for (let i = 0; i < hex.length; i += 2) {
            result += String.fromCharCode(parseInt(hex.slice(i, i + 2), 16));
        }
        return result;
    } catch {
        return "[解码错误: 无效的 Hex]";
    }
};

const unicodeEncode = (input: string) =>
    Array.from(input)
        .map((char) => "\\u" + char.charCodeAt(0).toString(16).padStart(4, "0"))
        .join("");

const unicodeDecode = (input: string) => {
    try {
        return input.replace(/\\u([0-9a-fA-F]{4})/g, (_matched, code) =>
            String.fromCharCode(parseInt(code, 16))
        );
    } catch {
        return "[解码错误: 无效的 Unicode]";
    }
};

const jwtDecode = (input: string) => {
    try {
        const parts = input.trim().split(".");
        if (parts.length !== 3) return "[错误: 无效的 JWT 格式]";
        const header = JSON.parse(atob(parts[0]));
        const payload = JSON.parse(atob(parts[1]));
        return JSON.stringify({ header, payload }, null, 2);
    } catch {
        return "[解码错误: 无效的 JWT]";
    }
};

const sha256Hash = async (input: string) => {
    const buffer = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(input));
    return Array.from(new Uint8Array(buffer))
        .map((byte) => byte.toString(16).padStart(2, "0"))
        .join("");
};

const sha1Hash = async (input: string) => {
    const buffer = await crypto.subtle.digest("SHA-1", new TextEncoder().encode(input));
    return Array.from(new Uint8Array(buffer))
        .map((byte) => byte.toString(16).padStart(2, "0"))
        .join("");
};

const sha512Hash = async (input: string) => {
    const buffer = await crypto.subtle.digest("SHA-512", new TextEncoder().encode(input));
    return Array.from(new Uint8Array(buffer))
        .map((byte) => byte.toString(16).padStart(2, "0"))
        .join("");
};

const md5Hash = (input: string) => {
    try {
        return CryptoJS.MD5(input).toString();
    } catch (error) {
        return `[MD5哈希错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm3Hash = (input: string) => {
    try {
        return sm.sm3(input);
    } catch (error) {
        return `[SM3哈希错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const aesEncrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "0123456789abcdef";
        const iv = options?.iv || "0123456789abcdef";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const outputFormat = (options?.outputFormat || "Base64").toLowerCase();
        const keySize = parseInt(options?.keySize || "128", 10);
        const encrypted = CryptoJS.AES.encrypt(input, parseKeyOrIv(key, keyFormat, keySize / 8), {
            iv: parseKeyOrIv(iv, ivFormat, 16),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return outputFormat === "hex" ? encrypted.ciphertext.toString(CryptoJS.enc.Hex) : encrypted.toString();
    } catch (error) {
        return `[AES加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const aesDecrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "0123456789abcdef";
        const iv = options?.iv || "0123456789abcdef";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const inputFormat = (options?.inputFormat || "Base64").toLowerCase();
        const keySize = parseInt(options?.keySize || "128", 10);
        const cipherParams = inputFormat === "hex"
            ? CryptoJS.lib.CipherParams.create({ ciphertext: CryptoJS.enc.Hex.parse(input.trim()) })
            : input.trim();
        const decrypted = CryptoJS.AES.decrypt(cipherParams, parseKeyOrIv(key, keyFormat, keySize / 8), {
            iv: parseKeyOrIv(iv, ivFormat, 16),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return decrypted.toString(CryptoJS.enc.Utf8) || "[解密结果为空，请检查密钥/IV/模式是否正确]";
    } catch (error) {
        return `[AES解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const desEncrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "12345678";
        const iv = options?.iv || "12345678";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const outputFormat = (options?.outputFormat || "Base64").toLowerCase();
        const encrypted = CryptoJS.DES.encrypt(input, parseKeyOrIv(key, keyFormat, 8), {
            iv: parseKeyOrIv(iv, ivFormat, 8),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return outputFormat === "hex" ? encrypted.ciphertext.toString(CryptoJS.enc.Hex) : encrypted.toString();
    } catch (error) {
        return `[DES加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const desDecrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "12345678";
        const iv = options?.iv || "12345678";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const inputFormat = (options?.inputFormat || "Base64").toLowerCase();
        const cipherParams = inputFormat === "hex"
            ? CryptoJS.lib.CipherParams.create({ ciphertext: CryptoJS.enc.Hex.parse(input.trim()) })
            : input.trim();
        const decrypted = CryptoJS.DES.decrypt(cipherParams, parseKeyOrIv(key, keyFormat, 8), {
            iv: parseKeyOrIv(iv, ivFormat, 8),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return decrypted.toString(CryptoJS.enc.Utf8) || "[解密结果为空]";
    } catch (error) {
        return `[DES解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const tripleDesEncrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "123456789012345678901234";
        const iv = options?.iv || "12345678";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const outputFormat = (options?.outputFormat || "Base64").toLowerCase();
        const encrypted = CryptoJS.TripleDES.encrypt(input, parseKeyOrIv(key, keyFormat, 24), {
            iv: parseKeyOrIv(iv, ivFormat, 8),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return outputFormat === "hex" ? encrypted.ciphertext.toString(CryptoJS.enc.Hex) : encrypted.toString();
    } catch (error) {
        return `[3DES加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const tripleDesDecrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "123456789012345678901234";
        const iv = options?.iv || "12345678";
        const mode = options?.mode || "CBC";
        const padding = options?.padding || "Pkcs7";
        const keyFormat = options?.keyFormat || "UTF8";
        const ivFormat = options?.ivFormat || "UTF8";
        const inputFormat = (options?.inputFormat || "Base64").toLowerCase();
        const cipherParams = inputFormat === "hex"
            ? CryptoJS.lib.CipherParams.create({ ciphertext: CryptoJS.enc.Hex.parse(input.trim()) })
            : input.trim();
        const decrypted = CryptoJS.TripleDES.decrypt(cipherParams, parseKeyOrIv(key, keyFormat, 24), {
            iv: parseKeyOrIv(iv, ivFormat, 8),
            mode: getCryptoMode(mode),
            padding: getCryptoPadding(padding),
        });
        return decrypted.toString(CryptoJS.enc.Utf8) || "[解密结果为空]";
    } catch (error) {
        return `[3DES解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rc4Encrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "secretkey";
        const keyBytes = (options?.keyFormat || "UTF8").toLowerCase() === "hex"
            ? CryptoJS.enc.Hex.parse(key)
            : CryptoJS.enc.Utf8.parse(key);
        const encrypted = CryptoJS.RC4.encrypt(input, keyBytes);
        return (options?.outputFormat || "Base64").toLowerCase() === "hex"
            ? encrypted.ciphertext.toString(CryptoJS.enc.Hex)
            : encrypted.toString();
    } catch (error) {
        return `[RC4加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rc4Decrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "secretkey";
        const keyBytes = (options?.keyFormat || "UTF8").toLowerCase() === "hex"
            ? CryptoJS.enc.Hex.parse(key)
            : CryptoJS.enc.Utf8.parse(key);
        const cipherParams = (options?.inputFormat || "Base64").toLowerCase() === "hex"
            ? CryptoJS.lib.CipherParams.create({ ciphertext: CryptoJS.enc.Hex.parse(input.trim()) })
            : input.trim();
        const decrypted = CryptoJS.RC4.decrypt(cipherParams, keyBytes);
        return decrypted.toString(CryptoJS.enc.Utf8) || "[解密结果为空]";
    } catch (error) {
        return `[RC4解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rabbitEncrypt = (input: string, options?: Record<string, string>) => {
    try {
        const keyBytes = CryptoJS.enc.Utf8.parse(options?.key || "secretkey");
        const iv = options?.iv;
        const config: { iv?: CryptoJS.lib.WordArray } = {};
        if (iv) config.iv = CryptoJS.enc.Utf8.parse(iv);
        const encrypted = CryptoJS.Rabbit.encrypt(input, keyBytes, config);
        return (options?.outputFormat || "Base64").toLowerCase() === "hex"
            ? encrypted.ciphertext.toString(CryptoJS.enc.Hex)
            : encrypted.toString();
    } catch (error) {
        return `[Rabbit加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rabbitDecrypt = (input: string, options?: Record<string, string>) => {
    try {
        const keyBytes = CryptoJS.enc.Utf8.parse(options?.key || "secretkey");
        const iv = options?.iv;
        const config: { iv?: CryptoJS.lib.WordArray } = {};
        if (iv) config.iv = CryptoJS.enc.Utf8.parse(iv);
        const cipherParams = (options?.inputFormat || "Base64").toLowerCase() === "hex"
            ? CryptoJS.lib.CipherParams.create({ ciphertext: CryptoJS.enc.Hex.parse(input.trim()) })
            : input.trim();
        const decrypted = CryptoJS.Rabbit.decrypt(cipherParams, keyBytes, config);
        return decrypted.toString(CryptoJS.enc.Utf8) || "[解密结果为空]";
    } catch (error) {
        return `[Rabbit解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm4Encrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "0123456789abcdef";
        const iv = options?.iv || "0123456789abcdef";
        const keyHex = (options?.keyFormat || "UTF8").toLowerCase() === "hex" ? key.padEnd(32, "0").slice(0, 32) : textToHex(key, 16);
        const ivHex = (options?.ivFormat || "UTF8").toLowerCase() === "hex" ? iv.padEnd(32, "0").slice(0, 32) : textToHex(iv, 16);
        const encrypted = sm.sm4.encrypt(input, keyHex, {
            mode: (options?.mode || "CBC").toUpperCase() === "ECB" ? "ecb" : "cbc",
            iv: ivHex,
            output: "array",
        }) as number[];
        const hexResult = encrypted.map((byte) => byte.toString(16).padStart(2, "0")).join("");
        return (options?.outputFormat || "Hex").toLowerCase() === "base64" ? hexToBase64(hexResult) : hexResult;
    } catch (error) {
        return `[SM4加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm4Decrypt = (input: string, options?: Record<string, string>) => {
    try {
        const key = options?.key || "0123456789abcdef";
        const iv = options?.iv || "0123456789abcdef";
        const keyHex = (options?.keyFormat || "UTF8").toLowerCase() === "hex" ? key.padEnd(32, "0").slice(0, 32) : textToHex(key, 16);
        const ivHex = (options?.ivFormat || "UTF8").toLowerCase() === "hex" ? iv.padEnd(32, "0").slice(0, 32) : textToHex(iv, 16);
        const inputHex = (options?.inputFormat || "Hex").toLowerCase() === "base64" ? base64ToHex(input) : input.trim();
        return sm.sm4.decrypt(inputHex, keyHex, {
            mode: (options?.mode || "CBC").toUpperCase() === "ECB" ? "ecb" : "cbc",
            iv: ivHex,
            output: "string",
        }) as string || "[解密结果为空]";
    } catch (error) {
        return `[SM4解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm2Encrypt = (input: string, options?: Record<string, string>) => {
    try {
        const publicKey = options?.publicKey || "";
        if (!publicKey) return "[错误: 请提供SM2公钥]";
        return sm.sm2.doEncrypt(input, publicKey, 1);
    } catch (error) {
        return `[SM2加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm2Decrypt = (input: string, options?: Record<string, string>) => {
    try {
        const privateKey = options?.privateKey || "";
        if (!privateKey) return "[错误: 请提供SM2私钥]";
        return sm.sm2.doDecrypt(input.trim(), privateKey, 1) || "[解密结果为空]";
    } catch (error) {
        return `[SM2解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const sm2GenerateKeyPair = () => {
    try {
        const keypair = sm.sm2.generateKeyPairHex();
        return `=== SM2 密钥对 ===\n\n【公钥 - 用于加密】\n${keypair.publicKey}\n\n【私钥 - 用于解密】\n${keypair.privateKey}`;
    } catch (error) {
        return `[生成密钥对错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rsaGenerateKeyPair = (_input: string, options?: Record<string, string>) => {
    try {
        const encrypt = new JSEncrypt({ default_key_size: (options?.keySize || "2048") as never });
        encrypt.getKey();
        return `=== RSA ${(options?.keySize || "2048")} 密钥对 ===\n\n【公钥 - 用于加密】\n${encrypt.getPublicKey()}\n\n【私钥 - 用于解密】\n${encrypt.getPrivateKey()}`;
    } catch (error) {
        return `[生成密钥对错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rsaEncrypt = (input: string, options?: Record<string, string>) => {
    try {
        const publicKey = options?.publicKey || "";
        if (!publicKey) return "[错误: 请提供RSA公钥]";
        const encrypt = new JSEncrypt();
        encrypt.setPublicKey(publicKey);
        return encrypt.encrypt(input) || "[加密失败: 请检查公钥格式或输入数据长度]";
    } catch (error) {
        return `[RSA加密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const rsaDecrypt = (input: string, options?: Record<string, string>) => {
    try {
        const privateKey = options?.privateKey || "";
        if (!privateKey) return "[错误: 请提供RSA私钥]";
        const decrypt = new JSEncrypt();
        decrypt.setPrivateKey(privateKey);
        return decrypt.decrypt(input.trim()) || "[解密失败: 请检查私钥格式或密文是否正确]";
    } catch (error) {
        return `[RSA解密错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const extractIP = (input: string) => {
    const unique = Array.from(new Set((input.match(regexIP) || []).filter(validateSingleIP)));
    return unique.length ? unique.join("\n") : "[未找到IP地址]";
};

const extractIPPort = (input: string) => {
    const regexIPPort =
        /\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):[0-9]{1,5}\b/g;
    const unique = Array.from(
        new Set(
            (input.match(regexIPPort) || []).filter((item) => {
                const [ip, port] = item.split(":");
                const portNum = parseInt(port, 10);
                return validateSingleIP(ip) && portNum >= 1 && portNum <= 65535;
            })
        )
    );
    return unique.length ? unique.join("\n") : "[未找到IP:端口]";
};

const extractURL = (input: string) => {
    const unique = Array.from(new Set(input.match(regexURL) || []));
    return unique.length ? unique.join("\n") : "[未找到URL]";
};

const extractDomain = (input: string) => {
    const domains = Array.from(new Set(input.match(regexURL) || []))
        .map((url) => {
            try {
                return new URL(url).hostname;
            } catch {
                return null;
            }
        })
        .filter((item): item is string => !!item);
    const unique = Array.from(new Set(domains));
    return unique.length ? unique.join("\n") : "[未找到域名]";
};

const extractRootDomain = (input: string) => {
    const rootDomains = Array.from(new Set(input.match(regexURL) || []))
        .map((url) => {
            try {
                const parts = new URL(url).hostname.replace(/^www\./, "").split(".");
                return parts.length >= 2 ? parts.slice(-2).join(".") : parts[0];
            } catch {
                return null;
            }
        })
        .filter((item): item is string => !!item);
    const unique = Array.from(new Set(rootDomains));
    return unique.length ? unique.join("\n") : "[未找到根域名]";
};

const extractPhone = (input: string) => {
    const unique = Array.from(new Set(input.match(regexpPhone) || []));
    return unique.length ? unique.join("\n") : "[未找到手机号]";
};

const extractIdCard = (input: string) => {
    const unique = Array.from(new Set(input.match(regexpIdCard) || []));
    return unique.length ? unique.join("\n") : "[未找到身份证号]";
};

const extractDruidWebURI = (input: string) => {
    const matches: string[] = [];
    regexDruidWebURI.lastIndex = 0;
    let match;
    while ((match = regexDruidWebURI.exec(input)) !== null) {
        matches.push(match[1]);
    }
    const unique = Array.from(new Set(matches));
    return unique.length ? unique.join("\n") : "[未找到WebURI]";
};

const extractDruidWebSession = (input: string) => {
    const matches: string[] = [];
    regexDruidWebSession.lastIndex = 0;
    let match;
    while ((match = regexDruidWebSession.exec(input)) !== null) {
        matches.push(match[1]);
    }
    const unique = Array.from(new Set(matches));
    return unique.length ? unique.join("\n") : "[未找到WebSession]";
};

const ipLocationLookup = async (input: string) => {
    const ips = Array.from(new Set(input.split(/\r?\n/).map((line) => line.trim()).filter(validateSingleIP)));
    if (!ips.length) return "[未找到可定位的IP]";
    const results: string[] = [];
    for (const ip of ips) {
        const location = await IpLocation(ip);
        results.push(`${ip} | ${location}`);
    }
    return results.join("\n");
};

const decryptFinereport = (input: string) => {
    if (input.length <= 3) return "[解密失败: 密码格式错误，示例: ___0072002a00670066000a]";
    let result = "";
    const temp = input.substring(3);
    for (let i = 0; i < temp.length / 4; i++) {
        const c1 = parseInt(temp.substring(i * 4, (i + 1) * 4), 16);
        result += String.fromCharCode(c1 ^ PASSWORD_MASK_ARRAY[i % 8]);
    }
    return result || "[解密失败]";
};

const decryptSeeyon = (input: string) => {
    const pass = input.replace(/\//g, "");
    const parts = pass.split(".0");
    if (parts.length <= 1) return "[解密失败: 密码格式错误，示例: /1.0/UWJ0dHgxc2U=]";
    try {
        let result = "";
        const iv = parseInt(parts[0], 10);
        const password = atob(parts[1]);
        for (let i = 0; i < password.length; i++) {
            result += String.fromCharCode(password.charCodeAt(i) - iv);
        }
        return result || "[解密失败]";
    } catch {
        return "[解密失败: Base64解码错误]";
    }
};

const decryptAesDesJson = (input: string) => {
    try {
        const data = JSON.parse(input);
        if (!Array.isArray(data.words) || typeof data.sigBytes !== "number") {
            return '[JSON格式错误，示例: {"words": [1161312566,808858179], "sigBytes": 8}]';
        }
        const buffer = new ArrayBuffer(data.words.length * 4);
        const view = new DataView(buffer);
        data.words.forEach((num: number, index: number) => view.setUint32(index * 4, num, false));
        return new TextDecoder("utf-8").decode(new Uint8Array(buffer, 0, data.sigBytes));
    } catch {
        return "[JSON解析错误]";
    }
};

const decryptWinSCP = (input: string) => {
    const WSCP_SIMPLE_MAGIC = 0xa3;
    const WSCP_SIMPLE_STRING = "0123456789ABCDEF";
    const WSCP_SIMPLE_FLAG = 0xff;
    const WSCP_SIMPLE_INTERNAL = 0x00;

    function decrypt(username: string, hostname: string, encodedPassword: string): string {
        if (!encodedPassword.match(/[A-F0-9]+/)) return "";
        let chars = encodedPassword.split("");
        const nextChar = () => {
            if (chars.length === 0) return WSCP_SIMPLE_INTERNAL;
            const a = WSCP_SIMPLE_STRING.indexOf(chars.shift()!);
            const b = WSCP_SIMPLE_STRING.indexOf(chars.shift()!);
            return WSCP_SIMPLE_FLAG & ~(((a << 4) + b) ^ WSCP_SIMPLE_MAGIC);
        };

        const result: string[] = [];
        const key = username + hostname;
        const flag = nextChar();
        let length = flag;
        if (flag === WSCP_SIMPLE_FLAG) {
            nextChar();
            length = nextChar();
        }
        chars = chars.slice(nextChar() * 2);
        for (let i = 0; i < length; i++) {
            result.push(String.fromCharCode(nextChar()));
        }
        if (flag === WSCP_SIMPLE_FLAG) {
            const valid = result.slice(0, key.length).join("");
            if (valid !== key) return "";
            return result.slice(key.length).join("");
        }
        return result.join("");
    }

    const sessions = input.split("\n[").filter((session) => session.includes("Sessions\\"));
    if (!sessions.length) return "[未找到WinSCP会话信息]";
    return sessions
        .map((session) => {
            const sessionName = decodeURIComponent(session.split("]")[0].replace("Sessions\\", "").replace("[", ""));
            const lines = session.split("\n");
            const username = lines.find((line) => line.startsWith("UserName="))?.split("=")[1]?.trim() || "";
            const hostname = lines.find((line) => line.startsWith("HostName="))?.split("=")[1]?.trim() || "";
            const encodedPassword = lines.find((line) => line.startsWith("Password="))?.split("=")[1]?.trim() || "";
            if (!username || !hostname || !encodedPassword) return `[${sessionName}] 信息不完整`;
            return `[${sessionName}]\n用户: ${username}\n主机: ${hostname}\n密码: ${decrypt(username, hostname, encodedPassword) || "解密失败"}`;
        })
        .join("\n\n");
};

const rot13 = (input: string) =>
    input.replace(/[a-zA-Z]/g, (char) => {
        const base = char <= "Z" ? 65 : 97;
        return String.fromCharCode(((char.charCodeAt(0) - base + 13) % 26) + base);
    });

const reverseString = (input: string) => Array.from(input).reverse().join("");
const toUpperCase = (input: string) => input.toUpperCase();
const toLowerCase = (input: string) => input.toLowerCase();
const toTitleCase = (input: string) => input.replace(/\w\S*/g, (txt) => txt.charAt(0).toUpperCase() + txt.slice(1).toLowerCase());
const trimWhitespace = (input: string) => input.trim();
const removeAllSpaces = (input: string) => input.replace(/\s/g, "");
const removeNewlines = (input: string) => input.replace(/[\r\n]/g, "");

const deduplicateLines = (input: string, options?: Record<string, string>) => {
    const { parts, outputSeparator } = splitBySeparator(input, options?.separator);
    const unique = Array.from(new Set(parts.filter((part) => part.trim())));
    return unique.length ? unique.join(outputSeparator) : "[无内容]";
};

const sortLines = (input: string, options?: Record<string, string>) => {
    const { parts, outputSeparator } = splitBySeparator(input, options?.separator);
    const lines = parts.filter((line) => line.trim());
    if (!lines.length) return "[无内容]";
    if ((options?.mode || "alphabetical") === "natural") {
        lines.sort((a, b) => a.localeCompare(b, undefined, { numeric: true, sensitivity: "base" }));
    } else {
        lines.sort((a, b) => a.localeCompare(b));
    }
    if ((options?.order || "asc") === "desc") lines.reverse();
    return lines.join(outputSeparator);
};

const regexFindReplace = (input: string, options?: Record<string, string>) => {
    const find = options?.find ?? "";
    if (!find) return "[请输入查找内容]";
    try {
        let flags = "";
        if (options?.globalMatch === "true") flags += "g";
        if (options?.caseInsensitive === "true") flags += "i";
        if (options?.multiline === "true") flags += "m";
        if (options?.dotAll === "true") flags += "s";
        return input.replace(new RegExp(find, flags), options?.replace ?? "");
    } catch (error) {
        return `[正则表达式错误: ${error instanceof Error ? error.message : "未知错误"}]`;
    }
};

const jsonFormat = (input: string) => {
    try {
        return JSON.stringify(JSON.parse(input), null, 2);
    } catch {
        return "[格式化错误: 无效的 JSON]";
    }
};

const jsonMinify = (input: string) => {
    try {
        return JSON.stringify(JSON.parse(input));
    } catch {
        return "[压缩错误: 无效的 JSON]";
    }
};

export const codecCategories: CodecCategory[] = [
    { id: "encoding", name: "编码/解码" },
    { id: "hash", name: "哈希计算" },
    { id: "crypto", name: "加密/解密" },
    { id: "extract", name: "数据提取" },
    { id: "decrypt", name: "密码破解" },
    { id: "text", name: "文本处理" },
    { id: "format", name: "格式转换" },
];

export const codecOperations: CodecOperation[] = [
    {
        id: "base64",
        name: "Base64",
        category: "encoding",
        description: "Base64 编码与解码",
        encode: base64Encode,
        decode: base64Decode,
        options: [{ key: "mode", label: "模式", default: "Standard", type: "select", options: ["Standard", "URL", "MIME"] }],
    },
    {
        id: "base32",
        name: "Base32",
        category: "encoding",
        description: "Base32 编码与解码",
        encode: base32Encode,
        decode: base32Decode,
        options: [{ key: "padding", label: "填充", default: "True", type: "select", options: ["True", "False"] }],
    },
    { id: "url", name: "URL", category: "encoding", description: "URL 编码与解码", encode: urlEncode, decode: urlDecode },
    { id: "html", name: "HTML 实体", category: "encoding", description: "HTML 实体编码与解码", encode: htmlEncode, decode: htmlDecode },
    { id: "hex", name: "Hex", category: "encoding", description: "十六进制编码与解码", encode: hexEncode, decode: hexDecode },
    { id: "unicode", name: "Unicode", category: "encoding", description: "Unicode 转义编码与解码", encode: unicodeEncode, decode: unicodeDecode },
    { id: "jwt", name: "JWT 解码", category: "encoding", description: "解码 JSON Web Token", decode: jwtDecode },
    { id: "sha256", name: "SHA-256", category: "hash", description: "SHA-256 哈希计算", transform: sha256Hash },
    { id: "sha1", name: "SHA-1", category: "hash", description: "SHA-1 哈希计算", transform: sha1Hash },
    { id: "sha512", name: "SHA-512", category: "hash", description: "SHA-512 哈希计算", transform: sha512Hash },
    { id: "md5", name: "MD5", category: "hash", description: "MD5 哈希计算", transform: md5Hash },
    { id: "sm3", name: "SM3", category: "hash", description: "SM3 国密哈希算法", transform: sm3Hash },
    {
        id: "aes",
        name: "AES",
        category: "crypto",
        description: "AES 对称加密",
        encode: aesEncrypt,
        decode: aesDecrypt,
        options: [
            { key: "key", label: "密钥", default: "0123456789abcdef" },
            { key: "keyFormat", label: "Key格式", default: "UTF8", type: "select", options: ["UTF8", "Hex", "Base64"] },
            { key: "keySize", label: "密钥位数", default: "128", type: "select", options: ["128", "192", "256"] },
            { key: "iv", label: "IV向量", default: "0123456789abcdef" },
            { key: "ivFormat", label: "IV格式", default: "UTF8", type: "select", options: ["UTF8", "Hex", "Base64"] },
            { key: "mode", label: "模式", default: "CBC", type: "select", options: ["CBC", "ECB", "CFB", "OFB", "CTR"] },
            { key: "padding", label: "填充", default: "Pkcs7", type: "select", options: ["Pkcs7", "ZeroPadding", "NoPadding", "Iso10126", "AnsiX923"] },
            { key: "outputFormat", label: "输出格式", default: "Base64", type: "select", options: ["Base64", "Hex"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Base64", type: "select", options: ["Base64", "Hex"] },
        ],
    },
    {
        id: "des",
        name: "DES",
        category: "crypto",
        description: "DES 对称加密",
        encode: desEncrypt,
        decode: desDecrypt,
        options: [
            { key: "key", label: "密钥(8字节)", default: "12345678" },
            { key: "keyFormat", label: "Key格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "iv", label: "IV向量(8字节)", default: "12345678" },
            { key: "ivFormat", label: "IV格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "mode", label: "模式", default: "CBC", type: "select", options: ["CBC", "ECB", "CFB", "OFB"] },
            { key: "padding", label: "填充", default: "Pkcs7", type: "select", options: ["Pkcs7", "ZeroPadding", "NoPadding"] },
            { key: "outputFormat", label: "输出格式", default: "Base64", type: "select", options: ["Base64", "Hex"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Base64", type: "select", options: ["Base64", "Hex"] },
        ],
    },
    {
        id: "3des",
        name: "3DES",
        category: "crypto",
        description: "Triple DES 加密",
        encode: tripleDesEncrypt,
        decode: tripleDesDecrypt,
        options: [
            { key: "key", label: "密钥(24字节)", default: "123456789012345678901234" },
            { key: "keyFormat", label: "Key格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "iv", label: "IV向量(8字节)", default: "12345678" },
            { key: "ivFormat", label: "IV格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "mode", label: "模式", default: "CBC", type: "select", options: ["CBC", "ECB", "CFB", "OFB"] },
            { key: "padding", label: "填充", default: "Pkcs7", type: "select", options: ["Pkcs7", "ZeroPadding", "NoPadding"] },
            { key: "outputFormat", label: "输出格式", default: "Base64", type: "select", options: ["Base64", "Hex"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Base64", type: "select", options: ["Base64", "Hex"] },
        ],
    },
    {
        id: "rc4",
        name: "RC4",
        category: "crypto",
        description: "RC4 流加密算法",
        encode: rc4Encrypt,
        decode: rc4Decrypt,
        options: [
            { key: "key", label: "密钥", default: "secretkey" },
            { key: "keyFormat", label: "Key格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "outputFormat", label: "输出格式", default: "Base64", type: "select", options: ["Base64", "Hex"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Base64", type: "select", options: ["Base64", "Hex"] },
        ],
    },
    {
        id: "rabbit",
        name: "Rabbit",
        category: "crypto",
        description: "Rabbit 流加密算法",
        encode: rabbitEncrypt,
        decode: rabbitDecrypt,
        options: [
            { key: "key", label: "密钥", default: "secretkey" },
            { key: "iv", label: "IV向量(可选)", default: "" },
            { key: "outputFormat", label: "输出格式", default: "Base64", type: "select", options: ["Base64", "Hex"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Base64", type: "select", options: ["Base64", "Hex"] },
        ],
    },
    {
        id: "sm4",
        name: "SM4",
        category: "crypto",
        description: "SM4 国密对称加密",
        encode: sm4Encrypt,
        decode: sm4Decrypt,
        options: [
            { key: "key", label: "密钥(16字节)", default: "0123456789abcdef" },
            { key: "keyFormat", label: "Key格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "iv", label: "IV向量(16字节)", default: "0123456789abcdef" },
            { key: "ivFormat", label: "IV格式", default: "UTF8", type: "select", options: ["UTF8", "Hex"] },
            { key: "mode", label: "模式", default: "CBC", type: "select", options: ["CBC", "ECB"] },
            { key: "outputFormat", label: "输出格式", default: "Hex", type: "select", options: ["Hex", "Base64"] },
            { key: "inputFormat", label: "输入格式(解密)", default: "Hex", type: "select", options: ["Hex", "Base64"] },
        ],
    },
    {
        id: "sm2",
        name: "SM2",
        category: "crypto",
        description: "SM2 国密非对称加密",
        encode: sm2Encrypt,
        decode: sm2Decrypt,
        options: [
            { key: "publicKey", label: "公钥(Hex)", default: "" },
            { key: "privateKey", label: "私钥(Hex)", default: "" },
        ],
    },
    { id: "sm2-keygen", name: "SM2 密钥生成", category: "crypto", description: "生成 SM2 密钥对", transform: sm2GenerateKeyPair },
    {
        id: "rsa",
        name: "RSA",
        category: "crypto",
        description: "RSA 非对称加密算法",
        encode: rsaEncrypt,
        decode: rsaDecrypt,
        options: [
            { key: "publicKey", label: "公钥(加密)", default: "" },
            { key: "privateKey", label: "私钥(解密)", default: "" },
        ],
    },
    {
        id: "rsa-keygen",
        name: "RSA 密钥生成",
        category: "crypto",
        description: "生成 RSA 密钥对",
        transform: rsaGenerateKeyPair,
        options: [{ key: "keySize", label: "密钥长度", default: "2048", type: "select", options: ["1024", "2048", "4096"] }],
    },
    { id: "extract-ip", name: "IP", category: "extract", description: "从文本中提取所有IP地址", transform: extractIP },
    { id: "extract-ipport", name: "IP:端口", category: "extract", description: "从文本中提取 IP:端口", transform: extractIPPort },
    { id: "extract-url", name: "URL", category: "extract", description: "从文本中提取 URL 链接", transform: extractURL },
    { id: "extract-domain", name: "域名", category: "extract", description: "从 URL 中提取域名", transform: extractDomain },
    { id: "extract-rootdomain", name: "根域名", category: "extract", description: "从 URL 中提取根域名", transform: extractRootDomain },
    { id: "extract-phone", name: "手机号", category: "extract", description: "提取中国大陆手机号", transform: extractPhone },
    { id: "extract-idcard", name: "身份证", category: "extract", description: "提取中国大陆身份证号码", transform: extractIdCard },
    { id: "extract-druid-uri", name: "Druid WebURI", category: "extract", description: "提取 Druid WebURI", transform: extractDruidWebURI },
    { id: "extract-druid-session", name: "Druid WebSession", category: "extract", description: "提取 Druid WebSession", transform: extractDruidWebSession },
    { id: "ip-location", name: "IP 定位", category: "extract", description: "批量定位输入中的 IP 地址", transform: ipLocationLookup },
    { id: "decrypt-winscp", name: "WinSCP - WinSCP.ini", category: "decrypt", description: "提取并解密 WinSCP.ini 中保存的密码", transform: decryptWinSCP },
    { id: "decrypt-finereport", name: "FineReport v8 - privilege.xml", category: "decrypt", description: "解密 FineReport v8 密码", transform: decryptFinereport },
    { id: "decrypt-seeyon", name: "Seeyon OA - datasourceCtp.properties", category: "decrypt", description: "解密致远 OA 数据库密码", transform: decryptSeeyon },
    { id: "decrypt-aes-des-json", name: "AES/DES JS数组Key转字符串", category: "decrypt", description: "将前端 JS WordArray JSON 转为字符串", transform: decryptAesDesJson },
    { id: "rot13", name: "ROT13", category: "text", description: "ROT13 字母替换", transform: rot13 },
    { id: "reverse", name: "反转字符串", category: "text", description: "反转字符串顺序", transform: reverseString },
    { id: "uppercase", name: "转大写", category: "text", description: "将文本转换为大写", transform: toUpperCase },
    { id: "lowercase", name: "转小写", category: "text", description: "将文本转换为小写", transform: toLowerCase },
    { id: "titlecase", name: "首字母大写", category: "text", description: "每个单词首字母大写", transform: toTitleCase },
    { id: "trim", name: "去除首尾空白", category: "text", description: "去除首尾空白字符", transform: trimWhitespace },
    { id: "removespaces", name: "移除所有空格", category: "text", description: "移除所有空白字符", transform: removeAllSpaces },
    { id: "removenewlines", name: "移除换行", category: "text", description: "移除所有换行符", transform: removeNewlines },
    {
        id: "deduplicate",
        name: "去重",
        category: "text",
        description: "根据指定分隔符去重，留空默认按行",
        transform: deduplicateLines,
        options: [{ key: "separator", label: "分隔符", default: "" }],
    },
    {
        id: "sort",
        name: "排序",
        category: "text",
        description: "对文本内容按指定规则排序",
        transform: sortLines,
        options: [
            { key: "separator", label: "分隔符", default: "" },
            { key: "order", label: "顺序", default: "asc", type: "select", options: ["asc", "desc"] },
            { key: "mode", label: "模式", default: "alphabetical", type: "select", options: ["alphabetical", "natural"] },
        ],
    },
    {
        id: "regex-find-replace",
        name: "正则查找替换",
        category: "text",
        description: "使用正则表达式进行查找和替换",
        transform: regexFindReplace,
        options: [
            { key: "find", label: "查找(正则)", default: "" },
            { key: "replace", label: "替换为", default: "" },
            { key: "globalMatch", label: "全局匹配", default: "true", type: "select", options: ["true", "false"] },
            { key: "caseInsensitive", label: "忽略大小写", default: "false", type: "select", options: ["true", "false"] },
            { key: "multiline", label: "多行模式", default: "false", type: "select", options: ["true", "false"] },
            { key: "dotAll", label: ".匹配所有", default: "false", type: "select", options: ["true", "false"] },
        ],
    },
    { id: "jsonformat", name: "JSON 格式化", category: "format", description: "格式化 JSON 字符串", transform: jsonFormat },
    { id: "jsonminify", name: "JSON 压缩", category: "format", description: "压缩 JSON 字符串", transform: jsonMinify },
];

export const codecOperationMap = new Map(codecOperations.map((operation) => [operation.id, operation]));

export const createAddedCodecOperation = (operationId: string, mode: OperationMode): AddedCodecOperation | null => {
    const operation = codecOperationMap.get(operationId);
    if (!operation) return null;
    const options = Object.fromEntries((operation.options || []).map((option) => [option.key, option.default]));
    return {
        id: `${operationId}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        operationId,
        mode,
        options,
    };
};

export const runCodecPipeline = async (input: string, steps: AddedCodecOperation[]) => {
    let result = input;
    for (const step of steps) {
        const operation = codecOperationMap.get(step.operationId);
        if (!operation) continue;
        try {
            if (step.mode === "encode" && operation.encode) {
                result = await operation.encode(result, step.options);
            } else if (step.mode === "decode" && operation.decode) {
                result = await operation.decode(result, step.options);
            } else if (step.mode === "transform" && operation.transform) {
                result = await operation.transform(result, step.options);
            }
        } catch (error) {
            result = `[执行错误: ${error instanceof Error ? error.message : "未知错误"}]`;
            break;
        }
    }
    return result;
};
