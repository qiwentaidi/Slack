package space

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/structs"

	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

func HunterApiSearch(api, query, pageSize, page, startTime, asset string, deduplication bool) *structs.HunterResult {
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
	var hr structs.HunterResult
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
	var hr structs.HunterResult
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

func SearchHunterTips(query string) *structs.HunterTips {
	var ts structs.HunterTips
	_, b, err := clients.NewSimpleGetRequest("https://hunter.qianxin.com/api/recommend?keyword="+HunterBaseEncode(query), clients.DefaultClient())
	if err != nil {
		return &ts
	}
	json.Unmarshal(b, &ts)
	return &ts
}
