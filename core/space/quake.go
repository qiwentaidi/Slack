package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strconv"
	"strings"
)

const (
	quakeTipsApi   = "https://quake.360.net/api/visitor/search/app"
	quakeServerApi = "https://quake.360.net/api/v3/search/quake_service"
	quakeUserApi   = "https://quake.360.net/api/v3/user/info"
)

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

func QuakeApiSearch(o *QuakeRequestOptions) *QuakeResult {
	var startIndex int
	if o.PageNum == 1 {
		startIndex = 0
	} else {
		startIndex = (o.PageNum - 1) * o.PageSize
	}
	data := make(map[string]interface{})
	if len(o.IpList) > 0 {
		data["ip_list"] = o.IpList
	} else {
		data["query"] = o.Query
	}
	data["start"] = startIndex
	data["size"] = o.PageSize
	data["latest"] = o.Latest
	data["shortcuts"] = getShortcuts(o)
	bytesData, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type": "application/json",
		"X-QuakeToken": o.Token,
	}
	_, body, err := clients.NewRequest("POST", quakeServerApi, header, bytes.NewReader(bytesData), 10, false, clients.DefaultClient())
	if err != nil {
		return &QuakeResult{}
	}
	if string(body) == "/quake/login" {
		return &QuakeResult{
			Code:    302,
			Message: "Token is error",
		}
	}
	if string(body) == "暂不支持搜索该内容" {
		return &QuakeResult{
			Code:    302,
			Message: "暂不支持搜索该内容",
		}
	}
	var qrk QuakeRawResult
	json.Unmarshal(body, &qrk)
	qk := &QuakeResult{
		Message: qrk.Message,
		Total:   qrk.Meta.Pagination.Total,
		Credit:  QuakeUserSearch(o.Token),
	}
	if code, err := strconv.Atoi(fmt.Sprintf("%v", qrk.Code)); err == nil {
		qk.Code = code
	} else {
		qk.Code = 500
	}
	for _, item := range qrk.Data {
		var components []string
		var protocol string
		for _, v := range item.Components {
			if v.ProductNameEn == "" {
				components = append(components, util.MergeNonEmpty([]string{v.ProductNameCn, v.Version}, "/"))
			} else {
				components = append(components, util.MergeNonEmpty([]string{v.ProductNameEn, v.Version}, "/"))
			}
		}
		if item.Service.Name == "http/ssl" {
			protocol = "https"
		} else {
			protocol = item.Service.Name
		}
		qk.Data = append(qk.Data, QuakeData{
			URL:        mergeURL(protocol, item.Service.HTTP.Host, item.IP, item.Port),
			Components: components,
			Port:       item.Port,
			Protocol:   protocol,
			Host:       item.Service.HTTP.Host,
			Title:      item.Service.HTTP.Title,
			IcpName:    item.Service.HTTP.Icp.Main_licence.Unit,
			IcpNumber:  item.Service.HTTP.Icp.Main_licence.Licence,
			IP:         item.IP,
			Isp:        item.Location.Isp,
			Position: util.MergePosition(structs.Position{
				Province:  item.Location.ProvinceCn,
				City:      item.Location.CityCn,
				District:  item.Location.DistrictCn,
				Connector: "/",
			}),
		})
	}
	qrk = QuakeRawResult{} // 清空内存
	return qk
}

// 查询剩余积分
func QuakeUserSearch(token string) int {
	header := map[string]string{
		"Content-Type": "application/json",
		"X-QuakeToken": token,
	}
	_, body, err := clients.NewRequest("GET", quakeUserApi, header, nil, 10, true, clients.DefaultClient())
	if err != nil {
		return 0
	}
	var qui QuakeUserInfo
	if err := json.Unmarshal(body, &qui); err != nil {
		return 0
	}
	return qui.Data.Credit
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

func SearchQuakeTips(query string) *QuakeTipsResult {
	var qs QuakeTipsResult
	jsonData := fmt.Sprintf(`{"app_name":"%v","device":{"UUID":"aa963dba-1bfa-54cf-9fdd-7b9be5a30890"}}`, query)
	header := map[string]string{
		"Content-Type": "application/json",
	}
	_, b, err := clients.NewRequest("POST", quakeTipsApi, header, strings.NewReader(jsonData), 10, true, clients.DefaultClient())
	if err != nil {
		return &qs
	}
	json.Unmarshal(b, &qs)
	return &qs
}

// 首次请求不能是带Shortcuts，需要在请求一次quake之后，获取到正确的cert_common值
func getShortcuts(o *QuakeRequestOptions) []string {
	if o.CertCommon == "" {
		return []string{}
	}
	var (
		shortcutsMeta         ShortcutsMeta
		cdn, honepot, invalid string
		shortcuts             []string
	)
	header := map[string]string{
		"Cookie": "cert_common=" + o.CertCommon,
	}
	_, body, err := clients.NewRequest("GET", "https://quake.360.net/api/search/shortcuts/quake_service_unique", header, nil, 10, true, clients.DefaultClient())
	if err != nil {
		return shortcuts
	}
	json.Unmarshal(body, &shortcutsMeta)
	for _, v := range shortcutsMeta.Data {
		if v.Name == "排除CDN" {
			cdn = v.Id
		}
		if v.Name == "排除蜜罐" {
			honepot = v.Id
		}
		if v.Name == "过滤无效请求" {
			invalid = v.Id
		}
	}

	if o.CDN {
		shortcuts = append(shortcuts, cdn)
	}
	if o.Honeypot {
		shortcuts = append(shortcuts, honepot)
	}
	if o.Invalid {
		shortcuts = append(shortcuts, invalid)
	}
	return shortcuts
}

func mergeURL(protocol, domain, ip string, port int) string {
	host := domain
	if host == "" {
		host = ip
	}
	if port == 80 || port == 443 {
		return fmt.Sprintf("%s://%s", protocol, host)
	}
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}
