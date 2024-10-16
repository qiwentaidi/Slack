export namespace clients {
	
	export class Proxy {
	    Enabled: boolean;
	    Mode: string;
	    Address: string;
	    Port: number;
	    Username: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Enabled = source["Enabled"];
	        this.Mode = source["Mode"];
	        this.Address = source["Address"];
	        this.Port = source["Port"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	    }
	}

}

export namespace dirsearch {
	
	export class Options {
	    Method: string;
	    URLs: string[];
	    Paths: string[];
	    Workers: number;
	    Timeout: number;
	    BodyExclude: string;
	    BodyLengthExcludeTimes: number;
	    StatusCodeExclude: number[];
	    Redirect: boolean;
	    Interval: number;
	    CustomHeader: string;
	    Recursion: number;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Method = source["Method"];
	        this.URLs = source["URLs"];
	        this.Paths = source["Paths"];
	        this.Workers = source["Workers"];
	        this.Timeout = source["Timeout"];
	        this.BodyExclude = source["BodyExclude"];
	        this.BodyLengthExcludeTimes = source["BodyLengthExcludeTimes"];
	        this.StatusCodeExclude = source["StatusCodeExclude"];
	        this.Redirect = source["Redirect"];
	        this.Interval = source["Interval"];
	        this.CustomHeader = source["CustomHeader"];
	        this.Recursion = source["Recursion"];
	    }
	}

}

export namespace info {
	
	export class CompanyInfo {
	    CompanyName: string;
	    Holding: string;
	    Investment: string;
	    RegStatus: string;
	    Domains: string[];
	    CompanyId: string;
	
	    static createFrom(source: any = {}) {
	        return new CompanyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CompanyName = source["CompanyName"];
	        this.Holding = source["Holding"];
	        this.Investment = source["Investment"];
	        this.RegStatus = source["RegStatus"];
	        this.Domains = source["Domains"];
	        this.CompanyId = source["CompanyId"];
	    }
	}
	export class WechatReulst {
	    CompanyName: string;
	    WechatName: string;
	    WechatNums: string;
	    Logo: string;
	    Qrcode: string;
	    Introduction: string;
	
	    static createFrom(source: any = {}) {
	        return new WechatReulst(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CompanyName = source["CompanyName"];
	        this.WechatName = source["WechatName"];
	        this.WechatNums = source["WechatNums"];
	        this.Logo = source["Logo"];
	        this.Qrcode = source["Qrcode"];
	        this.Introduction = source["Introduction"];
	    }
	}

}

export namespace isic {
	
	export class GithubResult {
	    Status: boolean;
	    Total: number;
	    Items: string[];
	    Link: string;
	
	    static createFrom(source: any = {}) {
	        return new GithubResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.Total = source["Total"];
	        this.Items = source["Items"];
	        this.Link = source["Link"];
	    }
	}

}

export namespace jsfind {
	
	export class InfoSource {
	    Filed: string;
	    Source: string;
	
	    static createFrom(source: any = {}) {
	        return new InfoSource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Filed = source["Filed"];
	        this.Source = source["Source"];
	    }
	}
	export class FindSomething {
	    JS: InfoSource[];
	    APIRoute: InfoSource[];
	    IP_URL: InfoSource[];
	    ChineseIDCard: InfoSource[];
	    ChinesePhone: InfoSource[];
	    SensitiveField: InfoSource[];
	
	    static createFrom(source: any = {}) {
	        return new FindSomething(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.JS = this.convertValues(source["JS"], InfoSource);
	        this.APIRoute = this.convertValues(source["APIRoute"], InfoSource);
	        this.IP_URL = this.convertValues(source["IP_URL"], InfoSource);
	        this.ChineseIDCard = this.convertValues(source["ChineseIDCard"], InfoSource);
	        this.ChinesePhone = this.convertValues(source["ChinesePhone"], InfoSource);
	        this.SensitiveField = this.convertValues(source["SensitiveField"], InfoSource);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class FileInfo {
	    Error: boolean;
	    Message: string;
	    Content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Error = source["Error"];
	        this.Message = source["Message"];
	        this.Content = source["Content"];
	    }
	}
	export class FileListInfo {
	    Path: string;
	    Name: string;
	    BaseName: string;
	    ModTime: string;
	    Size: number;
	
	    static createFrom(source: any = {}) {
	        return new FileListInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Name = source["Name"];
	        this.BaseName = source["BaseName"];
	        this.ModTime = source["ModTime"];
	        this.Size = source["Size"];
	    }
	}
	export class PathInfo {
	    Name: string;
	    Ext: string;
	    Dir: string;
	
	    static createFrom(source: any = {}) {
	        return new PathInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Ext = source["Ext"];
	        this.Dir = source["Dir"];
	    }
	}
	export class Syntax {
	    Name: string;
	    Content: string;
	
	    static createFrom(source: any = {}) {
	        return new Syntax(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Content = source["Content"];
	    }
	}
	export class pathTimes {
	    Path: string;
	    Times: number;
	
	    static createFrom(source: any = {}) {
	        return new pathTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Times = source["Times"];
	    }
	}

}

export namespace mongo {
	
	export class Client {
	
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace space {
	
	export class Data {
	    name: string;
	    company: string;
	    r_code: string;
	
	    static createFrom(source: any = {}) {
	        return new Data(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.company = source["company"];
	        this.r_code = source["r_code"];
	    }
	}
	export class Results {
	    URL: string;
	    Host: string;
	    Title: string;
	    IP: string;
	    Port: string;
	    Domain: string;
	    Protocol: string;
	    Region: string;
	    ICP: string;
	    Product: string;
	
	    static createFrom(source: any = {}) {
	        return new Results(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.Host = source["Host"];
	        this.Title = source["Title"];
	        this.IP = source["IP"];
	        this.Port = source["Port"];
	        this.Domain = source["Domain"];
	        this.Protocol = source["Protocol"];
	        this.Region = source["Region"];
	        this.ICP = source["ICP"];
	        this.Product = source["Product"];
	    }
	}
	export class FofaSearchResult {
	    Error: boolean;
	    Message: string;
	    Size: number;
	    Results: Results[];
	
	    static createFrom(source: any = {}) {
	        return new FofaSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Error = source["Error"];
	        this.Message = source["Message"];
	        this.Size = source["Size"];
	        this.Results = this.convertValues(source["Results"], Results);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuakeData {
	    Components: string[];
	    Port: number;
	    Protocol: string;
	    Host: string;
	    Title: string;
	    IcpName: string;
	    IcpNumber: string;
	    IP: string;
	    Isp: string;
	    Position: string;
	
	    static createFrom(source: any = {}) {
	        return new QuakeData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Components = source["Components"];
	        this.Port = source["Port"];
	        this.Protocol = source["Protocol"];
	        this.Host = source["Host"];
	        this.Title = source["Title"];
	        this.IcpName = source["IcpName"];
	        this.IcpNumber = source["IcpNumber"];
	        this.IP = source["IP"];
	        this.Isp = source["Isp"];
	        this.Position = source["Position"];
	    }
	}
	export class QuakeResult {
	    Code: number;
	    Message: string;
	    Data: QuakeData[];
	    Total: number;
	    Credit: number;
	
	    static createFrom(source: any = {}) {
	        return new QuakeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.Message = source["Message"];
	        this.Data = this.convertValues(source["Data"], QuakeData);
	        this.Total = source["Total"];
	        this.Credit = source["Credit"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class QuakeTipsData {
	    product_name: string;
	    vul_count: number;
	    vendor_name: string;
	    ip_count: number;
	
	    static createFrom(source: any = {}) {
	        return new QuakeTipsData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.product_name = source["product_name"];
	        this.vul_count = source["vul_count"];
	        this.vendor_name = source["vendor_name"];
	        this.ip_count = source["ip_count"];
	    }
	}
	export class QuakeTipsResult {
	    code: number;
	    message: string;
	    data: QuakeTipsData[];
	
	    static createFrom(source: any = {}) {
	        return new QuakeTipsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], QuakeTipsData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Result {
	    URL: string;
	    IP: string;
	    Domain: string;
	    Port: string;
	    Protocol: string;
	    Title: string;
	    Components: string;
	    Source: string;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.IP = source["IP"];
	        this.Domain = source["Domain"];
	        this.Port = source["Port"];
	        this.Protocol = source["Protocol"];
	        this.Title = source["Title"];
	        this.Components = source["Components"];
	        this.Source = source["Source"];
	    }
	}
	
	export class TipsResult {
	    code: number;
	    message: string;
	    data: Data[];
	
	    static createFrom(source: any = {}) {
	        return new TipsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], Data);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace structs {
	
	export class AntivirusResult {
	    Process: string;
	    Pid: string;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new AntivirusResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Process = source["Process"];
	        this.Pid = source["Pid"];
	        this.Name = source["Name"];
	    }
	}
	export class AuthPatch {
	    MS: string;
	    Patch: string;
	    Description: string;
	    System: string;
	    Reference: string;
	
	    static createFrom(source: any = {}) {
	        return new AuthPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MS = source["MS"];
	        this.Patch = source["Patch"];
	        this.Description = source["Description"];
	        this.System = source["System"];
	        this.Reference = source["Reference"];
	    }
	}
	export class Children {
	    Name: string;
	    Type: string;
	    Path: string;
	    Target: string;
	
	    static createFrom(source: any = {}) {
	        return new Children(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Path = source["Path"];
	        this.Target = source["Target"];
	    }
	}
	export class DatabaseConnection {
	    Nanoid: string;
	    Scheme: string;
	    Host: string;
	    Port: number;
	    Username: string;
	    Password: string;
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseConnection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Nanoid = source["Nanoid"];
	        this.Scheme = source["Scheme"];
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.Notes = source["Notes"];
	    }
	}
	export class HunterResult {
	    code: number;
	    // Go type: struct { AccountType string "json:\"account_type\""; Arr []struct { AsOrg string "json:\"as_org\""; Banner string "json:\"banner\""; BaseProtocol string "json:\"base_protocol\""; City string "json:\"city\""; Company string "json:\"company\""; Component []struct { Name string "json:\"name\""; Version string "json:\"version\"" } "json:\"component\""; Country string "json:\"country\""; Domain string "json:\"domain\""; IP string "json:\"ip\""; IsRisk string "json:\"is_risk\""; IsRiskProtocol string "json:\"is_risk_protocol\""; IsWeb string "json:\"is_web\""; Isp string "json:\"isp\""; Number string "json:\"number\""; Os string "json:\"os\""; Port int64 "json:\"port\""; Protocol string "json:\"protocol\""; Province string "json:\"province\""; StatusCode int64 "json:\"status_code\""; UpdatedAt string "json:\"updated_at\""; URL string "json:\"url\""; WebTitle string "json:\"web_title\"" } "json:\"arr\""; ConsumeQuota string "json:\"consume_quota\""; RestQuota string "json:\"rest_quota\""; SyntaxPrompt string "json:\"syntax_prompt\""; Time int64 "json:\"time\""; Total int64 "json:\"total\"" }
	    data: any;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], Object);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HunterTips {
	    code: number;
	    // Go type: struct { App []struct { Name string "json:\"name\""; AssetNum int "json:\"asset_num\""; Tags []string "json:\"tags\"" } "json:\"app\""; Collect []interface {} "json:\"collect\"" }
	    data: any;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterTips(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], Object);
	        this.message = source["message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Navigation {
	    Name: string;
	    Children: Children[];
	
	    static createFrom(source: any = {}) {
	        return new Navigation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Children = this.convertValues(source["Children"], Children);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response {
	    Error: boolean;
	    Proto: string;
	    Header: {[key: string]: string};
	    Body: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Error = source["Error"];
	        this.Proto = source["Proto"];
	        this.Header = source["Header"];
	        this.Body = source["Body"];
	    }
	}
	export class RowData {
	    Columns: string[];
	    Rows: any[][];
	
	    static createFrom(source: any = {}) {
	        return new RowData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Columns = source["Columns"];
	        this.Rows = source["Rows"];
	    }
	}
	export class SpaceOption {
	    FofaApi: string;
	    FofaEmail: string;
	    FofaKey: string;
	    HunterKey: string;
	    QuakeKey: string;
	
	    static createFrom(source: any = {}) {
	        return new SpaceOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FofaApi = source["FofaApi"];
	        this.FofaEmail = source["FofaEmail"];
	        this.FofaKey = source["FofaKey"];
	        this.HunterKey = source["HunterKey"];
	        this.QuakeKey = source["QuakeKey"];
	    }
	}
	export class Status {
	    Error: boolean;
	    Msg: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Error = source["Error"];
	        this.Msg = source["Msg"];
	    }
	}
	export class SubdomainOption {
	    Mode: number;
	    Domains: string[];
	    Subs: string[];
	    ChaosApi: string;
	    ZoomeyeApi: string;
	    SecuritytrailsApi: string;
	    BevigilApi: string;
	    GithubApi: string;
	    Thread: number;
	    Timeout: number;
	    ResolveExcludeTimes: number;
	    DnsServers: string[];
	
	    static createFrom(source: any = {}) {
	        return new SubdomainOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mode = source["Mode"];
	        this.Domains = source["Domains"];
	        this.Subs = source["Subs"];
	        this.ChaosApi = source["ChaosApi"];
	        this.ZoomeyeApi = source["ZoomeyeApi"];
	        this.SecuritytrailsApi = source["SecuritytrailsApi"];
	        this.BevigilApi = source["BevigilApi"];
	        this.GithubApi = source["GithubApi"];
	        this.Thread = source["Thread"];
	        this.Timeout = source["Timeout"];
	        this.ResolveExcludeTimes = source["ResolveExcludeTimes"];
	        this.DnsServers = source["DnsServers"];
	    }
	}
	export class WebscanOptions {
	    Target: string[];
	    Thread: number;
	    Screenshot: boolean;
	    DeepScan: boolean;
	    RootPath: boolean;
	    CallNuclei: boolean;
	    TemplateFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new WebscanOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Target = source["Target"];
	        this.Thread = source["Thread"];
	        this.Screenshot = source["Screenshot"];
	        this.DeepScan = source["DeepScan"];
	        this.RootPath = source["RootPath"];
	        this.CallNuclei = source["CallNuclei"];
	        this.TemplateFiles = source["TemplateFiles"];
	    }
	}

}

