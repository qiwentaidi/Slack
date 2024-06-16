export interface TableTabs {
    title: string;
    name: string;
    content: null | [{}];
    total: number;
    pageSize: number;
    currentPage: number;
}

export interface URLFingerMap {
    url: string,
    finger: string[]
}

export interface PortScanData {
    IP: string
    Port: number
    Server: string
    Link: string
    HttpTitle: string
}

export interface Vulnerability {
    vulID: string
    vulName: string
    protocoltype: string
    severity: string
    vulURL: string
    request: string
    response: string
    extInfo: string
    reference: string
    description: string
}

export interface FingerprintTable {
    url: string
    status: string
    length: string
    title: string
    detect: string
    existsWaf: boolean
    waf: string
    fingerprint: FingerLevel[]
}


export interface FingerLevel {
    name: string
    level: number // level 0 is default , level 1 is high risk
}

export interface Dir {
    Status: number
    URL: string
    Location: string
    Length: number
}

export interface DirScanOptions {
    Method: string
    URL: string
    Paths: string[]
    Workers: number
    Timeout: number
    BodyExclude: string
    BodyLengthExcludeTimes: number
    StatusCodeExclude: number[]
    FailedCounts: number
    Redirect: boolean
    Interval: number
    CustomHeader: string
}