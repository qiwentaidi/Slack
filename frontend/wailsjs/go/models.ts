export namespace clients {
	
	export class Proxy {
	
	
	    static createFrom(source: any = {}) {
	        return new Proxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace dirsearch {
	
	export class Options {
	
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace info {
	
	export class CompanyInfo {
	
	
	    static createFrom(source: any = {}) {
	        return new CompanyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class WechatReulst {
	
	
	    static createFrom(source: any = {}) {
	        return new WechatReulst(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace jsfind {
	
	export class FindSomething {
	
	
	    static createFrom(source: any = {}) {
	        return new FindSomething(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace main {
	
	export class Children {
	
	
	    static createFrom(source: any = {}) {
	        return new Children(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class FileInfo {
	
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class FormatTarget {
	
	
	    static createFrom(source: any = {}) {
	        return new FormatTarget(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Navigation {
	
	
	    static createFrom(source: any = {}) {
	        return new Navigation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class PathInfo {
	
	
	    static createFrom(source: any = {}) {
	        return new PathInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Response {
	
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Syntax {
	
	
	    static createFrom(source: any = {}) {
	        return new Syntax(source);
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
	export class FofaSearchResult {
	
	
	    static createFrom(source: any = {}) {
	        return new FofaSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
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
	export class HunterTipsResult {
	    code: number;
	    // Go type: struct { App []struct { Name string "json:\"name\""; AssetNum int "json:\"asset_num\""; Tags []string "json:\"tags\"" } "json:\"app\""; Collect []interface {} "json:\"collect\"" }
	    data: any;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HunterTipsResult(source);
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
	export class QuakeResult {
	
	
	    static createFrom(source: any = {}) {
	        return new QuakeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
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

