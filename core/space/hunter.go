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
	var beforeTime string
	currentTime := time.Now().Format("2006-01-02")
	switch startTime {
	case "0":
		beforeTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	case "1":
		beforeTime = time.Now().AddDate(0, 0, -179).Format("2006-01-02")
	case "2":
		beforeTime = time.Now().AddDate(-1, 0, -0).Format("2006-01-02")
	}
	address := "https://hunter.qianxin.com/openApi/search?api-key=" + api + "&search=" + HunterBaseEncode(query) + "&page=" +
		page + "&page_size=" + pageSize + "&is_web=" + asset + "&port_filter=" + fmt.Sprint(deduplication) + "&start_time=" + beforeTime + "&end_time=" + currentTime
	var hr structs.HunterResult
	_, b, err := clients.NewSimpleGetRequest(address, clients.NewHttpClient(nil, true))
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return &hr
	}
	json.Unmarshal(b, &hr)
	return &hr
}

// hunter base64加密接口
func HunterBaseEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func SearchHunterTips(query string) *structs.HunterTips {
	var ts structs.HunterTips
	_, b, err := clients.NewSimpleGetRequest("https://hunter.qianxin.com/api/recommend?keyword="+HunterBaseEncode(query), clients.NewHttpClient(nil, true))
	if err != nil {
		return &ts
	}
	json.Unmarshal(b, &ts)
	return &ts
}
