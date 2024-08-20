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
	Code int `json:"code"`
	Data struct {
		App []struct {
			Name     string   `json:"name"`
			AssetNum int      `json:"asset_num"`
			Tags     []string `json:"tags"`
		} `json:"app"`
		Collect []interface{} `json:"collect"`
	} `json:"data"`
	Message string `json:"message"`
}

// Hunter数据的结构体
type HunterResult struct {
	Code int64 `json:"code"`
	Data struct {
		AccountType string `json:"account_type"`
		Arr         []struct {
			AsOrg        string `json:"as_org"`
			Banner       string `json:"banner"`
			BaseProtocol string `json:"base_protocol"`
			City         string `json:"city"`
			Company      string `json:"company"`
			Component    []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"component"`
			Country        string `json:"country"`
			Domain         string `json:"domain"`
			IP             string `json:"ip"`
			IsRisk         string `json:"is_risk"`
			IsRiskProtocol string `json:"is_risk_protocol"`
			IsWeb          string `json:"is_web"`
			Isp            string `json:"isp"`
			Number         string `json:"number"`
			Os             string `json:"os"`
			Port           int64  `json:"port"`
			Protocol       string `json:"protocol"`
			Province       string `json:"province"`
			StatusCode     int64  `json:"status_code"`
			UpdatedAt      string `json:"updated_at"`
			URL            string `json:"url"`
			WebTitle       string `json:"web_title"`
		} `json:"arr"`
		ConsumeQuota string `json:"consume_quota"`
		RestQuota    string `json:"rest_quota"`
		SyntaxPrompt string `json:"syntax_prompt"`
		Time         int64  `json:"time"`
		Total        int64  `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
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
	ChaosApi            string
	ZoomeyeApi          string
	SecuritytrailsApi   string
	BevigilApi          string
	GethubApi           string
	Thread              int // 解析线程
	Timeout             int // 仅枚举模式启用时生效
	ResolveExcludeTimes int // 解析过滤IP次数
	DnsServers          []string
}
