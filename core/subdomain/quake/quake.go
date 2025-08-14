package quake

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"

	"github.com/qiwentaidi/clients"
)

func FetchHosts(ctx context.Context, domain, apikey string) []string {
	var result []string
	data := make(map[string]interface{})
	data["query"] = fmt.Sprintf("domain: %s", domain)
	data["start"] = 0
	data["size"] = 500
	data["latest"] = true
	bytesData, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type": "application/json",
		"X-QuakeToken": apikey,
	}
	resp, err := clients.DoRequest("POST", "https://quake.360.net/api/v3/search/quake_service", header, bytes.NewReader(bytesData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[subdomain] quake fetch host error: %v", err))
		return result
	}
	body := resp.Body()
	if string(body) == "/quake/login" {
		gologger.Debug(ctx, "[subdomain] quake fetch host error: token is err")
		return result
	}
	if string(body) == "暂不支持搜索该内容" {
		gologger.Debug(ctx, "[subdomain] quake fetch host error: can't search this content")
		return result
	}
	var qrk structs.QuakeRawResult
	json.Unmarshal(body, &qrk)
	if qrk.Data == nil {
		return result
	}
	for _, item := range qrk.Data {
		result = append(result, item.Service.HTTP.Host)
	}
	return result
}
