package fofa

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
)

func FetchHosts(ctx context.Context, domain string, auth *structs.FofaAuth) []string {
	var result []string
	address := fmt.Sprintf("%sapi/v1/search/all?email=%s&key=%s&qbase64=%s&is_fraud=false&is_honeypot=false&page=1&size=10000&fields=domain",
		auth.Address, auth.Email, auth.Key, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("domain=\"%s\"", domain))))
	_, b, err := clients.NewSimpleGetRequest(address, http.DefaultClient)
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[subdomain] fofa fetch host is error: %s", err))
		return result
	}
	var fr structs.FofaResult
	json.Unmarshal(b, &fr)
	if !fr.Error && fr.Size > 0 {
		for _, resultRow := range fr.Results {
			result = append(result, resultRow[3])
		}
	}
	return result
}
