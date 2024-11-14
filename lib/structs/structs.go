package structs

// 返回后端执行状态
// Error: true/false
// Msg:   错误信息
type Status struct {
	Error bool
	Msg   string
}

type Response struct {
	Error  bool
	Proto  string
	Header map[string]string
	Body   string
}

type Navigation struct {
	Name     string
	Children []Children
}

type Children struct {
	Name   string
	Type   string
	Path   string
	Target string
}

type TycCompanyInfo struct {
	CompanyId   string
	CompanyName string
}

type Position struct {
	Country   string
	Province  string
	City      string
	District  string
	Connector string // 连接符
}

type SpaceOption struct {
	FofaApi   string
	FofaEmail string
	FofaKey   string
	HunterKey string
	QuakeKey  string
}

type HunterTips struct {
	Code    int            `json:"code"`
	Data    HunterTipsData `json:"data"`
	Message string         `json:"message"`
}

type HunterTipsData struct {
	App     []HunterTipsApp `json:"app"`
	Collect []interface{}   `json:"collect"`
}

type HunterTipsApp struct {
	Name     string   `json:"name"`
	AssetNum int      `json:"asset_num"`
	Tags     []string `json:"tags"`
}

// Hunter数据的结构体
type HunterResult struct {
	Code    int64      `json:"code"`
	Data    HunterData `json:"data"`
	Message string     `json:"message"`
}

type HunterData struct {
	AccountType  string          `json:"account_type"`
	Arr          []HunterDataArr `json:"arr"`
	ConsumeQuota string          `json:"consume_quota"`
	RestQuota    string          `json:"rest_quota"`
	SyntaxPrompt string          `json:"syntax_prompt"`
	Time         int64           `json:"time"`
	Total        int64           `json:"total"`
}

type HunterDataArr struct {
	AsOrg          string            `json:"as_org"`
	Banner         string            `json:"banner"`
	BaseProtocol   string            `json:"base_protocol"`
	City           string            `json:"city"`
	Company        string            `json:"company"`
	Component      []HunterComponent `json:"component"`
	Country        string            `json:"country"`
	Domain         string            `json:"domain"`
	IP             string            `json:"ip"`
	IsRisk         string            `json:"is_risk"`
	IsRiskProtocol string            `json:"is_risk_protocol"`
	IsWeb          string            `json:"is_web"`
	Isp            string            `json:"isp"`
	Number         string            `json:"number"`
	Os             string            `json:"os"`
	Port           int64             `json:"port"`
	Protocol       string            `json:"protocol"`
	Province       string            `json:"province"`
	StatusCode     int64             `json:"status_code"`
	UpdatedAt      string            `json:"updated_at"`
	URL            string            `json:"url"`
	WebTitle       string            `json:"web_title"`
}

type HunterComponent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

const (
	EnumerationMode = 0
	ApiMode         = 1
	MixedMode       = 2
)

type SubdomainOption struct {
	Mode                int
	Domains             []string
	Subs                []string
	AppendEngines       []string
	FofaAddress         string
	FofaEmail           string
	FofaApi             string
	HunterApi           string
	QuakeApi            string
	ChaosApi            string
	ZoomeyeApi          string
	SecuritytrailsApi   string
	BevigilApi          string
	GithubApi           string
	Thread              int // 解析线程
	Timeout             int // 仅枚举模式启用时生效
	ResolveExcludeTimes int // 解析过滤IP次数
	DnsServers          []string
}

type DatabaseConnection struct {
	Nanoid   string
	Scheme   string
	Host     string
	Port     int
	Username string
	Password string
	Notes    string
}

type RowData struct {
	Columns []string
	Rows    [][]interface{}
}

type WebscanOptions struct {
	Target                []string
	Thread                int
	Screenshot            bool
	DeepScan              bool
	RootPath              bool
	CallNuclei            bool
	TemplateFiles         []string
	SkipNucleiWithoutTags bool
	GenerateLog4j2        bool // 开启后会将所有目标添加 Generate-Log4j2 的指纹
	AppendTemplateFolder  string
}

type AntivirusResult struct {
	Process string
	Pid     string
	Name    string
}

type AuthPatch struct {
	MS          string
	Patch       string
	Description string
	System      string
	Reference   string
}

type TaskResult struct {
	TaskId        string
	TaskName      string
	Targets       string
	Failed        int
	Vulnerability int
}

type QuakeRequestOptions struct {
	Query      string
	IpList     []string // 判断 IpList 是否为空决定是否为批量查询
	PageNum    int
	PageSize   int
	Latest     bool
	CDN        bool
	Invalid    bool
	Honeypot   bool
	Token      string
	CertCommon string // 让其他排除筛选的功能可以正常使用
}

// 原始数据中有用的字段
type QuakeRawResult struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Data    []struct {
		Components []struct {
			ProductNameEn string `json:"product_name_en"`
			ProductNameCn string `json:"product_name_cn"`
			Version       string `json:"version"`
		} `json:"components"`
		Port    int `json:"port"`
		Service struct {
			Name string `json:"name"`
			HTTP struct {
				Server string `json:"server"` // 中间件
				Host   string `json:"host"`
				Title  string `json:"title"`
				Icp    struct {
					Leader_name  string `json:"leader_name"`
					Domain       string `json:"domain"`
					Main_licence struct {
						Unit    string `json:"unit"`
						Nature  string `json:"nature"`
						Licence string `json:"licence"`
					} `json:"main_licence"`
					Content_type_name string `json:"content_type_name"`
					Limit_access      bool   `json:"limit_access"`
					Licence           string `json:"licence"`
				} `json:"icp"`
			} `json:"http"`
			TLS struct {
				Handshake_log struct {
					Server_certificates struct {
						Certificate struct {
							Parsed struct {
								Subject struct {
									Country      []string `json:"country"`
									Province     []string `json:"province"`
									Organization []string `json:"organization"`
									Common_name  []string `json:"common_name"`
								} `json:"subject"`
							} `json:"parsed"`
						} `json:"certificate"`
					} `json:"server_certificates"`
				} `json:"handshake_log"`
			} `json:"tls"`
		} `json:"service"`
		IP       string `json:"ip"`
		Location struct {
			Isp        string `json:"isp"`
			ProvinceCn string `json:"province_cn"`
			DistrictCn string `json:"district_cn"`
			CityCn     string `json:"city_cn"`
		} `json:"location"`
	}
	Meta struct {
		Pagination struct {
			Count     int `json:"count"`
			PageIndex int `json:"page_index"`
			PageSize  int `json:"page_size"`
			Total     int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

type QuakeResult struct {
	Code    int    // 响应状态信息，正常是0
	Message string // 提示信息
	Data    []QuakeData
	Total   int
	Credit  int // 剩余积分
}

type QuakeData struct {
	URL        string
	Components []string
	Port       int
	Protocol   string // 协议类型
	Host       string
	Title      string
	IcpName    string // 证书申请单位
	IcpNumber  string // 证书域名
	IP         string
	Isp        string
	Position   string
}

type QuakeUserInfo struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Data    struct {
		ID   string `json:"id"`
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Fullname string `json:"fullname"`
			Email    string `json:"email"`
		} `json:"user"`
		Baned            bool   `json:"baned"`
		BanStatus        string `json:"ban_status"`
		Credit           int    `json:"credit"`
		PersistentCredit int    `json:"persistent_credit"`
		Token            string `json:"token"`
		MobilePhone      string `json:"mobile_phone"`
		Source           string `json:"source"`
		PrivacyLog       struct {
			Status bool        `json:"status"`
			Time   interface{} `json:"time"`
		} `json:"privacy_log"`
		EnterpriseInformation struct {
			Name   interface{} `json:"name"`
			Email  interface{} `json:"email"`
			Status string      `json:"status"`
		} `json:"enterprise_information"`
		PersonalInformationStatus bool `json:"personal_information_status"`
		Role                      []struct {
			Fullname string `json:"fullname"`
			Priority int    `json:"priority"`
			Credit   int    `json:"credit"`
		} `json:"role"`
	} `json:"data"`
	Meta struct {
	} `json:"meta"`
}

type ShortcutsMeta struct {
	Meta    interface{} `json:"meta"`
	Code    float64     `json:"code"`
	Message string      `json:"message"`
	Data    []struct {
		Published      bool    `json:"published"`
		Order          float64 `json:"order"`
		Is_new         bool    `json:"is_new"`
		Put_more_tools bool    `json:"put_more_tools"`
		Id             string  `json:"id"`
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		Index          string  `json:"index"`
	} `json:"data"`
}

type QuakeTipsResult struct {
	Code    float64         `json:"code"`
	Message string          `json:"message"`
	Data    []QuakeTipsData `json:"data"`
}

type QuakeTipsData struct {
	Product_name string  `json:"product_name"`
	Vul_count    float64 `json:"vul_count"`
	Vendor_name  string  `json:"vendor_name"`
	Ip_count     float64 `json:"ip_count"`
}

type FofaAuth struct {
	Address string
	Email   string
	Key     string
}

type TipsResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	RCode   string `json:"r_code"`
}

type FofaResult struct {
	Error   bool       `json:"error"`
	Errmsg  string     `json:"errmsg"`
	Mode    string     `json:"mode"`
	Page    int64      `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int64      `json:"size"`
}

type FofaSearchResult struct {
	Error   bool
	Message string
	Size    int64
	Results []Results
}

type Results struct {
	URL      string
	Host     string
	Title    string
	IP       string
	Port     string
	Domain   string
	Protocol string
	Region   string
	ICP      string
	Product  string
}

type VulnerabilityInfo struct {
	ID          string
	Name        string
	Description string
	Reference   string
	Type        string
	Risk        string
	URL         string
	Request     string
	Response    string
	Extract     string
}

type NucleiOption struct {
	SkipNucleiWithoutTags bool // 如果没有扫描到指纹，是否需要扫描全漏洞还是直接跳过
	URL                   string
	Tags                  []string // 全漏洞扫描时，使用自定义标签
	TemplateFile          []string
	TemplateFolders       []string
}

type InfoResult struct {
	URL          string
	StatusCode   int
	Length       int
	Title        string
	Fingerprints []string
	IsWAF        bool
	WAF          string
	Detect       string
	Screenshot   string // 截图图片路径
}

type WebReport struct {
	Targets      string
	Fingerprints []InfoResult
	POCs         []VulnerabilityInfo
}

type CompanyInfo struct {
	CompanyName string
	Holding     string
	Investment  string // 投资比例
	RegStatus   string
	Domains     []string
	CompanyId   string
}

type WechatReulst struct {
	CompanyName  string
	WechatName   string
	WechatNums   string
	Logo         string
	Qrcode       string
	Introduction string
}
