package hunter

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

func FetchHosts(ctx context.Context, domain, apikey string) []string {
	hunterStartTime := time.Now().AddDate(-1, 0, -0).Format("2006-01-02")
	address := "https://hunter.qianxin.com/openApi/search?api-key=" + apikey + "&search=" + base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf(`domain.suffix="%s"`, domain))) + "&page=1&page_size=100&is_web=3&port_filter=false&start_time=" + hunterStartTime + "&end_time=" + time.Now().Format("2006-01-02")
	var hr structs.HunterResult
	_, b, err := clients.NewSimpleGetRequest(address, clients.NewHttpClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[subdomain] hunter fetch host error: %v", err))
		return []string{}
	}
	json.Unmarshal(b, &hr)
	var result []string
	for _, item := range hr.Data.Arr {
		if item.Domain != "" {
			result = append(result, item.Domain)
		}
	}
	return result
}
