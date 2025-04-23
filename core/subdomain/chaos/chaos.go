package chaos

import (
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
)

const chaosURL = "https://dns.projectdiscovery.io/"

type ChaosHosts struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Count      float64  `json:"count"`
}

// subdomains return is not complete subdomain, e.g www
func FetchHosts(ctx context.Context, domain, apikey string) *ChaosHosts {
	header := map[string]string{
		"Authorization": apikey,
	}
	searchURL := fmt.Sprintf("%sdns/%s/subdomains", chaosURL, domain)
	resp, err := clients.DoRequest("GET", searchURL, header, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, err)
		return nil
	}
	var ch ChaosHosts
	json.Unmarshal(resp.Body(), &ch)
	return &ch
}
