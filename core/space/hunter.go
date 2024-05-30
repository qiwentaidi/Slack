package space

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"

	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// 存储Hunter数据的结构体
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

func HunterApiSearch(api, query, pageSize, page, startTime, asset string, deduplication bool) *HunterResult {
	var hunterStartTime string
	switch startTime {
	case "0":
		hunterStartTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	case "1":
		hunterStartTime = time.Now().AddDate(0, 0, -179).Format("2006-01-02")
	case "2":
		hunterStartTime = time.Now().AddDate(-1, 0, -0).Format("2006-01-02")
	}
	address := "https://hunter.qianxin.com/openApi/search?api-key=" + api + "&search=" + HunterBaseEncode(query) + "&page=" +
		page + "&page_size=" + pageSize + "&is_web=" + asset + "&port_filter=" + fmt.Sprint(deduplication) + "&start_time=" + hunterStartTime + "&end_time=" + time.Now().Format("2006-01-02")
	var hr HunterResult
	_, b, err := clients.NewSimpleGetRequest(address, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return &hr
	}
	json.Unmarshal(b, &hr)
	return &hr
}

func SearchTotal(api, search string) (total int64, message string) {
	current_time := time.Now()
	before_time := current_time.AddDate(0, -1, 0)
	addr := "https://hunter.qianxin.com/openApi/search?api-key=" + api + "&search=" + HunterBaseEncode(search) +
		"&page=1&page_size=1&is_web=3&port_filter=false&start_time=" + before_time.Format("2006-01-02") + "&end_time=" + current_time.Format("2006-01-02")
	_, b, err := clients.NewSimpleGetRequest(addr, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return
	}
	var hr HunterResult
	json.Unmarshal(b, &hr)
	if hr.Code == 200 {
		total = hr.Data.Total
		message = hr.Data.RestQuota
	}
	return total, message
}

// hunter base64加密接口
func HunterBaseEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

type HunterTipsResult struct {
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

func SearchHunterTips(query string) *HunterTipsResult {
	var ts HunterTipsResult
	_, b, err := clients.NewSimpleGetRequest("https://hunter.qianxin.com/api/recommend?keyword="+HunterBaseEncode(query), clients.DefaultClient())
	if err != nil {
		return &ts
	}
	json.Unmarshal(b, &ts)
	return &ts
}
