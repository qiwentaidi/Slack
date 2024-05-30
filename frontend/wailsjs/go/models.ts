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

