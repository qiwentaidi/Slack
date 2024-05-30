package space

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slack-wails/lib/clients"
)

type QuakeResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Components []struct {
			ProductNameCn string `json:"product_name_cn"`
			Version       string `json:"version"`
		} `json:"components"`
		Port    int `json:"port"`
		Service struct {
			Name string `json:"name"`
			HTTP struct {
				Server string `json:"server"`
				Host   string `json:"host"`
				Title  string `json:"title"`
			} `json:"http"`
		} `json:"service"`
		IP  string `json:"ip"`
		TLS struct {
			Certificate struct {
				Subject struct {
					Organization []string `json:"organization"`
					CommonName   []string `json:"common_name"`
				} `json:"subject"`
			} `json:"tls"`
		} `json:"data"`
		Location struct {
			Isp        string `json:"isp"`
			ProvinceCn string `json:"province_cn"`
			DistrictCn string `json:"district_cn"`
			CityCn     string `json:"city_cn"`
		} `json:"location"`
		Meta struct {
			Pagination struct {
				Count     int `json:"count"`
				PageIndex int `json:"page_index"`
				PageSize  int `json:"page_size"`
				Total     int `json:"total"`
			} `json:"pagination"`
		} `json:"meta"`
	}
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

func QuakeApiSearch(query string, startIndex, pageSize int, token string, latest bool) (*QuakeResult, int) {
	data := make(map[string]interface{})
	data["query"] = query
	data["start"] = startIndex
	data["size"] = pageSize
	// data["ignore_cache"] = false
	// data["include"] = []string{"ip", "port", "service.http.title", "service.name", "service.http.host", "service.http.server", "location.city_cn"}
	data["latest"] = latest
	bytesData, _ := json.Marshal(data)
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("X-QuakeToken", token)
	_, body, err := clients.NewRequest("POST", "https://quake.360.net/api/v3/search/quake_service", header, bytes.NewReader(bytesData), 10, true, clients.DefaultClient())
	if err != nil {
		return &QuakeResult{}, 0
	}
	var qk QuakeResult
	json.Unmarshal(body, &qk)
	return &qk, QuakeUserSearch(token)
}

func QuakeUserSearch(token string) (surplus int) {
	header := http.Header{}
	header.Set("X-QuakeToken", token)
	_, body, err := clients.NewRequest("GET", "https://quake.360.net/api/v3/user/info", header, nil, 10, true, clients.DefaultClient())
	if err != nil {
		return
	}
	var qui QuakeUserInfo
	if err := json.Unmarshal(body, &qui); err != nil {
		return
	}
	return qui.Data.Credit
}
