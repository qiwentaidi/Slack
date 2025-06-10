package structs

type WindowsSize struct {
	Width  int
	Height int
}

// 返回后端执行状态
// Error: true/false
// Msg:   错误信息
type Status struct {
	Error bool
	Msg   string
}

type Response struct {
	Error     bool
	Proto     string
	StatsCode int
	Header    map[string]string
	Body      string
}

type Navigation struct {
	Name     string
	Children []Children
}

type Children struct {
	Name    string
	Type    string
	Path    string
	Target  string
	Favicon string // 纯图标路径
	Args    string // 占位符解析参数
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
	CdnMode         = 3
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
	Nanoid     string
	Scheme     string
	Host       string
	Port       int
	Username   string
	Password   string
	ServerName string // Oracle 的服务名称
	Notes      string
}

type RowData struct {
	Columns   []string
	Rows      [][]interface{}
	RowsCount int
}

type WebscanOptions struct {
	Target                []string
	TcpTarget             map[string][]string // tcp层的目标，兼容nuclei可以扫描
	Thread                int
	Screenshot            bool
	DeepScan              bool
	RootPath              bool
	CallNuclei            bool
	Tags                  []string
	TemplateFiles         []string
	SkipNucleiWithoutTags bool
	GenerateLog4j2        bool   // 开启后会将所有目标添加 Generate-Log4j2 的指纹
	AppendTemplateFolder  string // 追加模板文件夹
	NetworkCard           string // 指定扫描网卡
	CustomHeaders         string // 自定义请求头
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
		Components []Component `json:"components"`
		Hostname   string      `json:"hostname"`
		Org        string      `json:"org"`
		Port       int         `json:"port"`
		Service    Service     `json:"service"`
		Version    string      `json:"version"`
		IP         string      `json:"ip"`
		Location   Location    `json:"location"`
		IsIPv6     bool        `json:"is_ipv6"`
		Transport  string      `json:"transport"`
		Time       string      `json:"time"`
		ASN        int         `json:"asn"`
		ID         string      `json:"id"`
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

type Component struct {
	ProductLevel   string   `json:"product_level"`
	ProductType    []string `json:"product_type"`
	ProductVendor  string   `json:"product_vendor"`
	ProductNameCN  string   `json:"product_name_cn"`
	ProductNameEN  string   `json:"product_name_en"`
	ID             string   `json:"id"`
	ProductCatalog []string `json:"product_catalog"`
	Version        string   `json:"version"`
}

type Service struct {
	Response string `json:"response"`
	Name     string `json:"name"`
	HTTP     HTTP   `json:"http"`
	Cert     string `json:"cert"`
}

type HTTP struct {
	XPoweredBy      string  `json:"x_powered_by"`
	Server          string  `json:"server"`
	Path            string  `json:"path"`
	HTMLHash        string  `json:"html_hash"`
	ResponseHeaders string  `json:"response_headers"`
	StatusCode      int     `json:"status_code"`
	Favicon         Favicon `json:"favicon"`
	Host            string  `json:"host"`
	Body            string  `json:"body"`
	MetaKeywords    string  `json:"meta_keywords"`
	Title           string  `json:"title"`
	ICP             struct {
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
}

type Favicon struct {
	Data     string `json:"data"`
	Location string `json:"location"`
	Hash     string `json:"hash"`
	S3URL    string `json:"s3_url"`
}

type Location struct {
	Owner       string    `json:"owner"`
	ProvinceCN  string    `json:"province_cn"`
	ISP         string    `json:"isp"`
	ProvinceEN  string    `json:"province_en"`
	CountryEN   string    `json:"country_en"`
	DistrictCN  string    `json:"district_cn"`
	GPS         []float64 `json:"gps"`
	StreetCN    string    `json:"street_cn"`
	CityEN      string    `json:"city_en"`
	DistrictEN  string    `json:"district_en"`
	CountryCN   string    `json:"country_cn"`
	StreetEN    string    `json:"street_en"`
	CityCN      string    `json:"city_cn"`
	CountryCode string    `json:"country_code"`
	ASName      string    `json:"asname"`
	SceneCN     string    `json:"scene_cn"`
	SceneEN     string    `json:"scene_en"`
	Radius      float64   `json:"radius"`
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
	IcpName    string // 备案单位名称
	IcpNumber  string // 备案单位编号
	CertName   string // 证书申请单位
	FaviconURL string // 图标地址
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

type FofaSingleFiledResult struct {
	Error   bool     `json:"error"`
	Errmsg  string   `json:"errmsg"`
	Mode    string   `json:"mode"`
	Page    int64    `json:"page"`
	Query   string   `json:"query"`
	Results []string `json:"results"`
	Size    int64    `json:"size"`
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

type FofaUserInfo struct {
	Error           bool   `json:"error"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Category        string `json:"category"`
	Fcoin           int    `json:"fcoin"`
	FofaPoint       int    `json:"fofa_point"`
	RemainFreePoint int    `json:"remain_free_point"`
	RemainAPIQuery  int    `json:"remain_api_query"`
	RemainAPIData   int    `json:"remain_api_data"`
	Isvip           bool   `json:"isvip"`
	VipLevel        int    `json:"vip_level"`
	IsVerified      bool   `json:"is_verified"`
	Avatar          string `json:"avatar"`
	Message         string `json:"message"`
	FofacliVer      string `json:"fofacli_ver"`
	FofaServer      bool   `json:"fofa_server"`
	Expiration      string `json:"expiration"`
}

type VulnerabilityInfo struct {
	TaskId       string // 任务ID
	ID           string
	Name         string
	Description  string
	Reference    string
	Type         string
	Severity     string
	URL          string
	Request      string
	Response     string
	ResponseTime string
	Extract      string
}

type NucleiOption struct {
	SkipNucleiWithoutTags bool // 如果没有扫描到指纹，是否需要扫描全漏洞还是直接跳过
	URL                   string
	Tags                  []string // 指纹识别到的标签
	CustomTags            []string // 自定义的指纹标签
	TemplateFile          []string
	TemplateFolders       []string
	CustomHeaders         string
	Proxy                 string
}

type InfoResult struct {
	TaskId       string // 任务ID
	URL          string // 网站链接
	Scheme       string // 协议
	Host         string // 域名 或者 IP
	Port         int
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
	CompanyName      string
	Trademark        string // 商标
	Investment       string // 投资比例
	Amount           string // 注册资金
	RegStatus        string // 注册状态
	Domains          []string
	Subsidiaries     []CompanyInfo
	Apps             []App             // App
	Applets          []Applet          // 小程序
	OfficialAccounts []OfficialAccount // 公众号
}

type App struct {
	CityID           int    `json:"cityId"`
	CountyID         int    `json:"countyId"`
	DataID           int    `json:"dataId"`
	LeaderName       string `json:"leaderName"`
	MainID           string `json:"mainId"`
	MainLicence      string `json:"mainLicence"`
	MainUnitAddress  string `json:"mainUnitAddress"`
	MainUnitCertNo   string `json:"mainUnitCertNo"`
	MainUnitCertType int    `json:"mainUnitCertType"`
	NatureID         int    `json:"natureId"`
	NatureName       string `json:"natureName"`
	ProvinceID       int    `json:"provinceId"`
	ServiceID        int64  `json:"serviceId"`
	ServiceLicence   string `json:"serviceLicence"`
	ServiceName      string `json:"serviceName"`
	ServiceType      int    `json:"serviceType"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
	Version          string `json:"version"`
}

type Applet struct {
	CityID           int    `json:"cityId"`
	CountyID         int    `json:"countyId"`
	DataID           int    `json:"dataId"`
	LeaderName       string `json:"leaderName"`
	MainID           string `json:"mainId"`
	MainLicence      string `json:"mainLicence"`
	MainUnitAddress  string `json:"mainUnitAddress"`
	MainUnitCertNo   string `json:"mainUnitCertNo"`
	MainUnitCertType int    `json:"mainUnitCertType"`
	NatureID         int    `json:"natureId"`
	NatureName       string `json:"natureName"`
	ProvinceID       int    `json:"provinceId"`
	ServiceID        int64  `json:"serviceId"`
	ServiceLicence   string `json:"serviceLicence"`
	ServiceName      string `json:"serviceName"`
	ServiceType      int    `json:"serviceType"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
	Version          string `json:"version"`
}

type OfficialAccount struct {
	Name         string
	Numbers      string
	Logo         string
	Qrcode       string
	Introduction string
}

type SpaceEngineSyntax struct {
	Name    string
	Content string
}

type PathTimes struct {
	Path  string
	Times int
}

type NacosConfig struct {
	Name     string
	NodeInfo NacosNode
}

type NacosNode struct {
	Auth     int
	OSS      int
	Database int
}

type DatabaseInfo struct {
	Name  string
	Table []TableInfo
}

type TableInfo struct {
	Name      string
	RowsCount int
}

type NetwordCard struct {
	Name string
	IP   string
}

type InfoSource struct {
	Filed  string
	Source string
}

type FindSomething struct {
	JS        []InfoSource
	APIRoute  []InfoSource
	IP_URL    []InfoSource
	IDCard    []InfoSource
	Phone     []InfoSource
	Email     []InfoSource
	Sensitive []InfoSource
}

type JSFindResult struct {
	Target   string
	VulType  string
	Severity string
	Source   string
	Method   string
	Request  string
	Response string
	Length   int
	Filed    string
}

type JSFindOptions struct {
	HomeURL        string
	BaseURL        string
	ApiList        []string
	Authentication []string
	HighRiskRouter []string
	Headers        map[string]string
	// 低权限用户请求头
	LowPrivilegeHeaders map[string]string
}

type DataSource struct {
	Tianyancha Tianyancha
	Miit       Miit
}

type Tianyancha struct {
	Enable bool
	Token  string
	Id     string
}

type Miit struct {
	API string
}
