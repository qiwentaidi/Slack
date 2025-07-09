import { space, structs } from "wailsjs/go/models";

export interface TableTabs {
    title: string;
    name: string;
    content: any[];
    total: number;
    pageSize: number;
    currentPage: number;
    message: string; // 查询完的提示信息
}

export interface QuakeTableTabs {
    title: string;
    name: string;
    content: structs.QuakeData[];
    total: number;
    pageSize: number;
    currentPage: number;
    isBatch: boolean;
    ipList: string[];
    message: string;
}

export interface DirseearchResult {
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


export interface HunterEntryTips {
    value: string;
    assetNum: number;
    tags: string[];
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
    content: space.Result[];
    total: number;
    pageSize: number;
    currentPage: number;
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

export interface PocDetail {
    Name: string
    AssociatedFingerprint: string[]
}


export interface TablePane {
    name: string
    content: string
    rowsCount: number
    matchedType: string
}

export interface DatabaseConnection {
    nanoid: string
    type: string
    host: string
    port: number
    username: string
    password: string
    servername: string
    notes: string
    connected: boolean
    loading: boolean
    databaseCount?: number
    tableCount?: number
    tablePanes: TablePane[]
    progress: number // 当前采集进度
}

export interface JSFindData {
    Target: string
    Method: string
    Source: string
    VulType: string
    Severity: string
    Request: string
    Length: number
    Filed: string
    Response: string
}

export interface TreeNode {
    id: string;
    label: string;
    isDir: boolean;
    children?: TreeNode[];
}

export interface DirectoryTab {
    name: string;
    title: string;
    path: string;
    status: 'select' | 'tree';
    treeData: TreeNode[];
    isCollapse: boolean;
}


export interface Matcher {
    type: string;
    part?: string;
    words: string[];
    condition?: string;
    wordsText?: string;
}

export interface Metadata {
    verified?: boolean;
    fofa?: string;
    google?: string;
    shodan?: string;
    hunter?: string;
    max_requests?: number;
}

export interface FormData {
    id: string;
    name: string;
    author: string;
    severity: string;
    description: string;
    reference?: string;
    tags?: string[];
    body: string;
    metadata?: Metadata,
    matchers: Matcher[];
    matchersCondition: string;
}

export interface ActivityItem {
    content: string
    type: 'primary' | 'success' | 'warning' | 'danger' | 'info'
    timestamp?: string
    icon?: any
}