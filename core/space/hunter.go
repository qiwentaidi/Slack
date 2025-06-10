package space

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"

	"time"
)

const defaultHunterApi = "https://hunter.qianxin.com/"

func HunterApiSearch(ctx context.Context, api, key, query, pageSize, page, startTime, asset string, deduplication bool) *structs.HunterResult {
	var beforeTime string
	currentTime := time.Now().Format("2006-01-02")
	switch startTime {
	case "0":
		beforeTime = time.Now().AddDate(0, 0, -29).Format("2006-01-02")
	case "1":
		beforeTime = time.Now().AddDate(0, 0, -179).Format("2006-01-02")
	case "2":
		beforeTime = time.Now().AddDate(-1, 0, -0).Format("2006-01-02")
	}
	address := api + "openApi/search?api-key=" + key + "&search=" + HunterBaseEncode(query) + "&page=" +
		page + "&page_size=" + pageSize + "&is_web=" + asset + "&port_filter=" + fmt.Sprint(deduplication) + "&start_time=" + beforeTime + "&end_time=" + currentTime
	var hr structs.HunterResult
	resp, err := clients.SimpleGet(address, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.DualLog(ctx, gologger.Level_ERROR, fmt.Sprintf("[hunter] fetch host error: %v", err))
		return &hr
	}
	json.Unmarshal(resp.Body(), &hr)
	return &hr
}

// hunter base64加密接口
func HunterBaseEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func SearchHunterTips(query string) *structs.HunterTips {
	var ts structs.HunterTips
	resp, err := clients.SimpleGet("https://hunter.qianxin.com/api/recommend?keyword="+HunterBaseEncode(query), clients.NewRestyClient(nil, true))
	if err != nil {
		return &ts
	}
	json.Unmarshal(resp.Body(), &ts)
	return &ts
}
