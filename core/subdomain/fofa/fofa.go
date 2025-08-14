package fofa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slack-wails/core/space"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"

	"github.com/qiwentaidi/clients"
)

func FetchHosts(ctx context.Context, domain string, auth structs.FofaAuth) []string {
	var result []string
	query := space.FOFABaseEncode(fmt.Sprintf("domain=\"%s\"", domain))
	address := fmt.Sprintf("%sapi/v1/search/all?email=%s&key=%s&qbase64=%s&is_fraud=false&is_honeypot=false&page=1&size=10&fields=link", auth.Address, auth.Email, auth.Key, query)
	resp, err := clients.SimpleGet(address, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[subdomain] fofa fetch host is error: %s", err))
		return result
	}
	var fr structs.FofaSingleFiledResult
	json.Unmarshal(resp.Body(), &fr)
	if !fr.Error && fr.Size > 0 {
		for _, link := range fr.Results {
			url, err := url.Parse(link)
			if err == nil {
				result = append(result, url.Hostname())
			}
		}
	}
	return result
}
