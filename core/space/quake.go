package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"strings"
)

type QuakeRequestOptions struct {
	Query    string
	PageNum  int
	PageSize int
	Latest   bool
	CDN      bool
	Invalid  bool
	Honeypot bool
	Token    string
}

// 原始数据中有用的字段
type QuakeRawResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Components []struct {
			ProductNameEn string `json:"product_name_en"`
			Version       string `json:"version"`
		} `json:"components"`
		Port    int `json:"port"`
		Service struct {
			Name string `json:"name"`
			HTTP struct {
				Server string `json:"server"` // 中间件
				Host   string `json:"host"`
				Title  string `json:"title"`
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
	Components  []string
	Port        int
	Protocol    string // 协议类型
	Host        string
	Title       string
	CertCompany string // 证书申请单位
	CertDomain  string // 证书域名
	IP          string
	Isp         string
	Position    string
}

type QuakeUserInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

// cdn -> 635fcbaacc57190bd8826d0b
// honeypot -> 635fcb52cc57190bd8826d09
// invalid -> 63734bfa9c27d4249ca7261c
func QuakeApiSearch(o *QuakeRequestOptions) *QuakeResult {
	var startIndex int
	var shortcuts []string
	if o.PageNum == 1 {
		startIndex = 0
	} else {
		startIndex = (o.PageNum - 1) * o.PageSize
	}
	if o.CDN {
		shortcuts = append(shortcuts, "635fcbaacc57190bd8826d0b")
	}
	if o.Honeypot {
		shortcuts = append(shortcuts, "635fcb52cc57190bd8826d09")
	}
	if o.Invalid {
		shortcuts = append(shortcuts, "635fcb0acc57190bd8826d0c")
	}
	data := make(map[string]interface{})
	data["query"] = o.Query
	data["start"] = startIndex
	data["size"] = o.PageSize
	data["latest"] = o.Latest
	data["shortcuts"] = shortcuts
	bytesData, _ := json.Marshal(data)
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("X-QuakeToken", o.Token)
	_, body, err := clients.NewRequest("POST", "https://quake.360.net/api/v3/search/quake_service", header, bytes.NewReader(bytesData), 10, true, clients.DefaultClient())
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
			Message: "Token is error",
		}
	}
	var qrk QuakeRawResult
	json.Unmarshal(body, &qrk)
	qk := &QuakeResult{
		Code:    qrk.Code,
		Message: qrk.Message,
		Total:   qrk.Meta.Pagination.Total,
		Credit:  QuakeUserSearch(o.Token),
	}
	for _, item := range qrk.Data {
		var components []string
		for _, v := range item.Components {
			components = append(components, util.MergeNonEmpty([]string{v.ProductNameEn, v.Version}, "/"))
		}
		qk.Data = append(qk.Data, QuakeData{
			Components:  components,
			Port:        item.Port,
			Protocol:    item.Service.Name,
			Host:        item.Service.HTTP.Host,
			Title:       item.Service.HTTP.Title,
			CertCompany: strings.Join(item.Service.TLS.Handshake_log.Server_certificates.Certificate.Parsed.Subject.Organization, ","),
			CertDomain:  strings.Join(item.Service.TLS.Handshake_log.Server_certificates.Certificate.Parsed.Subject.Common_name, ","),
			IP:          item.IP,
			Isp:         item.Location.Isp,
			Position:    util.MergeNonEmpty([]string{item.Location.ProvinceCn, item.Location.CityCn, item.Location.DistrictCn}, "/"),
		})
	}
	qrk = QuakeRawResult{} // 清空内存
	return qk
}

// 查询剩余积分
func QuakeUserSearch(token string) int {
	header := http.Header{}
	header.Set("X-QuakeToken", token)
	_, body, err := clients.NewRequest("GET", "https://quake.360.net/api/v3/user/info", header, nil, 10, true, clients.DefaultClient())
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
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.112 Safari/537.36")
	_, b, err := clients.NewRequest("POST", "https://quake.360.net/api/visitor/search/app", header, strings.NewReader(jsonData), 10, true, clients.DefaultClient())
	if err != nil {
		return &qs
	}
	json.Unmarshal(b, &qs)
	return &qs
}
