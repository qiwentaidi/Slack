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

export namespace services {
	
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
	export class Tree {
	    id: string;
	    label: string;
	    isDir: boolean;
	    hits?: {[key: string]: number};
	    children?: Tree[];
	
	    static createFrom(source: any = {}) {
	        return new Tree(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.isDir = source["isDir"];
	        this.hits = source["hits"];
	        this.children = this.convertValues(source["children"], Tree);
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

export namespace space {
	
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
	    Favicon: string;
	    Args: string;
	
	    static createFrom(source: any = {}) {
	        return new Children(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Type = source["Type"];
	        this.Path = source["Path"];
	        this.Target = source["Target"];
	        this.Favicon = source["Favicon"];
	        this.Args = source["Args"];
	    }
	}
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
	export class DatabaseConnection {
	    Nanoid: string;
	    Scheme: string;
	    Host: string;
	    Port: number;
	    Username: string;
	    Password: string;
	    ServerName: string;
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
	        this.ServerName = source["ServerName"];
	        this.Notes = source["Notes"];
	    }
	}
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
	    IDCard: InfoSource[];
	    Phone: InfoSource[];
	    Email: InfoSource[];
	    Sensitive: InfoSource[];
	
	    static createFrom(source: any = {}) {
	        return new FindSomething(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.JS = this.convertValues(source["JS"], InfoSource);
	        this.APIRoute = this.convertValues(source["APIRoute"], InfoSource);
	        this.IP_URL = this.convertValues(source["IP_URL"], InfoSource);
	        this.IDCard = this.convertValues(source["IDCard"], InfoSource);
	        this.Phone = this.convertValues(source["Phone"], InfoSource);
	        this.Email = this.convertValues(source["Email"], InfoSource);
	        this.Sensitive = this.convertValues(source["Sensitive"], InfoSource);
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
	export class HunterComponent {
	    name: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterComponent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	    }
	}
	export class HunterDataArr {
	    as_org: string;
	    banner: string;
	    base_protocol: string;
	    city: string;
	    company: string;
	    component: HunterComponent[];
	    country: string;
	    domain: string;
	    ip: string;
	    is_risk: string;
	    is_risk_protocol: string;
	    is_web: string;
	    isp: string;
	    number: string;
	    os: string;
	    port: number;
	    protocol: string;
	    province: string;
	    status_code: number;
	    updated_at: string;
	    url: string;
	    web_title: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterDataArr(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.as_org = source["as_org"];
	        this.banner = source["banner"];
	        this.base_protocol = source["base_protocol"];
	        this.city = source["city"];
	        this.company = source["company"];
	        this.component = this.convertValues(source["component"], HunterComponent);
	        this.country = source["country"];
	        this.domain = source["domain"];
	        this.ip = source["ip"];
	        this.is_risk = source["is_risk"];
	        this.is_risk_protocol = source["is_risk_protocol"];
	        this.is_web = source["is_web"];
	        this.isp = source["isp"];
	        this.number = source["number"];
	        this.os = source["os"];
	        this.port = source["port"];
	        this.protocol = source["protocol"];
	        this.province = source["province"];
	        this.status_code = source["status_code"];
	        this.updated_at = source["updated_at"];
	        this.url = source["url"];
	        this.web_title = source["web_title"];
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
	export class HunterData {
	    account_type: string;
	    arr: HunterDataArr[];
	    consume_quota: string;
	    rest_quota: string;
	    syntax_prompt: string;
	    time: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new HunterData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.account_type = source["account_type"];
	        this.arr = this.convertValues(source["arr"], HunterDataArr);
	        this.consume_quota = source["consume_quota"];
	        this.rest_quota = source["rest_quota"];
	        this.syntax_prompt = source["syntax_prompt"];
	        this.time = source["time"];
	        this.total = source["total"];
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
	
	export class HunterResult {
	    code: number;
	    data: HunterData;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], HunterData);
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
	export class HunterTipsApp {
	    name: string;
	    asset_num: number;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new HunterTipsApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.asset_num = source["asset_num"];
	        this.tags = source["tags"];
	    }
	}
	export class HunterTipsData {
	    app: HunterTipsApp[];
	    collect: any[];
	
	    static createFrom(source: any = {}) {
	        return new HunterTipsData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app = this.convertValues(source["app"], HunterTipsApp);
	        this.collect = source["collect"];
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
	    data: HunterTipsData;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterTips(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.data = this.convertValues(source["data"], HunterTipsData);
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
	
	
	export class InfoResult {
	    TaskId: string;
	    URL: string;
	    Scheme: string;
	    Host: string;
	    Port: number;
	    StatusCode: number;
	    Length: number;
	    Title: string;
	    Fingerprints: string[];
	    IsWAF: boolean;
	    WAF: string;
	    Detect: string;
	    Screenshot: string;
	
	    static createFrom(source: any = {}) {
	        return new InfoResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskId = source["TaskId"];
	        this.URL = source["URL"];
	        this.Scheme = source["Scheme"];
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	        this.StatusCode = source["StatusCode"];
	        this.Length = source["Length"];
	        this.Title = source["Title"];
	        this.Fingerprints = source["Fingerprints"];
	        this.IsWAF = source["IsWAF"];
	        this.WAF = source["WAF"];
	        this.Detect = source["Detect"];
	        this.Screenshot = source["Screenshot"];
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
	export class NetwordCard {
	    Name: string;
	    IP: string;
	
	    static createFrom(source: any = {}) {
	        return new NetwordCard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.IP = source["IP"];
	    }
	}
	export class PathTimes {
	    Path: string;
	    Times: number;
	
	    static createFrom(source: any = {}) {
	        return new PathTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Times = source["Times"];
	    }
	}
	export class QuakeData {
	    URL: string;
	    Components: string[];
	    Port: number;
	    Protocol: string;
	    Host: string;
	    Title: string;
	    IcpName: string;
	    IcpNumber: string;
	    CertName: string;
	    FaviconURL: string;
	    IP: string;
	    Isp: string;
	    Position: string;
	
	    static createFrom(source: any = {}) {
	        return new QuakeData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.URL = source["URL"];
	        this.Components = source["Components"];
	        this.Port = source["Port"];
	        this.Protocol = source["Protocol"];
	        this.Host = source["Host"];
	        this.Title = source["Title"];
	        this.IcpName = source["IcpName"];
	        this.IcpNumber = source["IcpNumber"];
	        this.CertName = source["CertName"];
	        this.FaviconURL = source["FaviconURL"];
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
	export class Response {
	    Error: boolean;
	    Proto: string;
	    StatsCode: number;
	    Header: {[key: string]: string};
	    Body: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Error = source["Error"];
	        this.Proto = source["Proto"];
	        this.StatsCode = source["StatsCode"];
	        this.Header = source["Header"];
	        this.Body = source["Body"];
	    }
	}
	
	export class RowData {
	    Columns: string[];
	    Rows: any[][];
	    RowsCount: number;
	
	    static createFrom(source: any = {}) {
	        return new RowData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Columns = source["Columns"];
	        this.Rows = source["Rows"];
	        this.RowsCount = source["RowsCount"];
	    }
	}
	export class SpaceEngineSyntax {
	    Name: string;
	    Content: string;
	
	    static createFrom(source: any = {}) {
	        return new SpaceEngineSyntax(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Content = source["Content"];
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
	    AppendEngines: string[];
	    FofaAddress: string;
	    FofaEmail: string;
	    FofaApi: string;
	    HunterApi: string;
	    QuakeApi: string;
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
	        this.AppendEngines = source["AppendEngines"];
	        this.FofaAddress = source["FofaAddress"];
	        this.FofaEmail = source["FofaEmail"];
	        this.FofaApi = source["FofaApi"];
	        this.HunterApi = source["HunterApi"];
	        this.QuakeApi = source["QuakeApi"];
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
	export class TaskResult {
	    TaskId: string;
	    TaskName: string;
	    Targets: string;
	    Failed: number;
	    Vulnerability: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskId = source["TaskId"];
	        this.TaskName = source["TaskName"];
	        this.Targets = source["Targets"];
	        this.Failed = source["Failed"];
	        this.Vulnerability = source["Vulnerability"];
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
	export class VulnerabilityInfo {
	    TaskId: string;
	    ID: string;
	    Name: string;
	    Description: string;
	    Reference: string;
	    Type: string;
	    Severity: string;
	    URL: string;
	    Request: string;
	    Response: string;
	    ResponseTime: string;
	    Extract: string;
	
	    static createFrom(source: any = {}) {
	        return new VulnerabilityInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TaskId = source["TaskId"];
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Reference = source["Reference"];
	        this.Type = source["Type"];
	        this.Severity = source["Severity"];
	        this.URL = source["URL"];
	        this.Request = source["Request"];
	        this.Response = source["Response"];
	        this.ResponseTime = source["ResponseTime"];
	        this.Extract = source["Extract"];
	    }
	}
	export class WebReport {
	    Targets: string;
	    Fingerprints: InfoResult[];
	    POCs: VulnerabilityInfo[];
	
	    static createFrom(source: any = {}) {
	        return new WebReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Targets = source["Targets"];
	        this.Fingerprints = this.convertValues(source["Fingerprints"], InfoResult);
	        this.POCs = this.convertValues(source["POCs"], VulnerabilityInfo);
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
	export class WebscanOptions {
	    Target: string[];
	    TcpTarget: {[key: string]: string[]};
	    Thread: number;
	    Screenshot: boolean;
	    DeepScan: boolean;
	    RootPath: boolean;
	    CallNuclei: boolean;
	    Tags: string[];
	    TemplateFiles: string[];
	    SkipNucleiWithoutTags: boolean;
	    GenerateLog4j2: boolean;
	    AppendTemplateFolder: string;
	    NetworkCard: string;
	    CustomHeaders: string;
	
	    static createFrom(source: any = {}) {
	        return new WebscanOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Target = source["Target"];
	        this.TcpTarget = source["TcpTarget"];
	        this.Thread = source["Thread"];
	        this.Screenshot = source["Screenshot"];
	        this.DeepScan = source["DeepScan"];
	        this.RootPath = source["RootPath"];
	        this.CallNuclei = source["CallNuclei"];
	        this.Tags = source["Tags"];
	        this.TemplateFiles = source["TemplateFiles"];
	        this.SkipNucleiWithoutTags = source["SkipNucleiWithoutTags"];
	        this.GenerateLog4j2 = source["GenerateLog4j2"];
	        this.AppendTemplateFolder = source["AppendTemplateFolder"];
	        this.NetworkCard = source["NetworkCard"];
	        this.CustomHeaders = source["CustomHeaders"];
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
	export class WindowsSize {
	    Width: number;
	    Height: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowsSize(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	    }
	}

}

