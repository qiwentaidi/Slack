import { space } from "wailsjs/go/models";

export interface TableTabs {
    title: string;
    name: string;
    content: any[];
    total: number;
    pageSize: number;
    currentPage: number;
}

export interface QuakeTableTabs {
    title: string;
    name: string;
    content: space.QuakeData[];
    total: number;
    pageSize: number;
    currentPage: number;
    isBatch: boolean;
    ipList: string[];
}

export interface PortScanData {
    IP: string;
    Port: number;
    Server: string;
    Link: string;
    HttpTitle: string;
}

export interface Vulnerability {
    vulID: string;
    vulName: string;
    protocoltype: string;
    severity: string;
    vulURL: string;
    request: string;
    response: string;
    extInfo: string;
    reference: string;
    description: string;
}

export interface FingerprintTable {
    url: string;
    status: string;
    length: string;
    title: string;
    detect: string;
    existsWaf: boolean;
    waf: string;
    fingerprint: string[];
    screenshot: string;
}

export interface Dir {
    Status: number;
    URL: string;
    Location: string;
    Length: number;
    Body: string;
    Recursion: number;
}

export interface ProxyOptions {
    Enabled: boolean;
    Mode: string;
    Address: string;
    Port: number;
    Username: string;
    Password: string;
}

export interface Results {
    URL: string;
    Host: string;
    Title: string;
    IP: string;
    Port: string;
    Domain: string;
    Protocol: string;
    Region: string;
    ICP: string;
}

export interface File {
    Error?: boolean;
    Message?: string;
    Content?: string;
}

export interface HunterEntryTips {
    value: string;
    assetNum: number;
    tags: string[];
}

export interface CompanyInfo {
    CompanyName?: string;
    Holding?: string;
    Investment?: string;
    RegStatus?: string;
    Domains?: string[];
}

export interface WechatInfo {
    CompanyName?: string;
    WechatName?: string;
    WechatNums?: string;
    Logo?: string;
    Qrcode?: string;
    Introduction?: string;
}

export interface QuakeTipsResult {
    Code?: number;
    Message?: string;
    Data?: QuakeTipsData[];
}

export interface QuakeTipsData {
    Product_name: string;
    Vul_count: number;
    Vendor_name: string;
    Ip_count: number;
}

export interface Uncover {
    title: string;
    name: string;
    content: UncoverData[];
    total: number;
    pageSize: number;
    currentPage: number;
}

interface UncoverData {
    URL: string
    IP: string
    Domain: string
    Port: string
    Protocol: string
    Component: string
    Source: string
}

export interface LogInfo {
    Level: string
    Msg: string
}

export interface BruteResult {
    Host: string
    Port: string
    Protocol: string
    Username: string
    Password: string
}

export interface SubdomainInfo {
    Domain: string
    Subdomain: string
    Ips: string[]
    IsCdn: boolean
    CdnName: string
    Source: string
}

export interface ISICResult {
    Query: string
    Status: boolean
    Total: number
    Items: string[]
    Link: string
}

export interface PocDetail {
    Name: string
    AssociatedFingerprint: string[]
}
