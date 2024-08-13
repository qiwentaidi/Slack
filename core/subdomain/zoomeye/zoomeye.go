package zoomeye

import (
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
)

const zoomeyeURL = "https://api.zoomeye.org/"

type List struct {
	Name      string   `json:"name"`
	Timestamp string   `json:"timestamp"`
	Ip        []string `json:"ip"`
}

type ZoomeyeHost struct {
	Status float64 `json:"status"`
	Total  float64 `json:"total"`
	List   []List  `json:"list"`
	Msg    string  `json:"msg"`
	Type   float64 `json:"type"`
}

// subdomains return is complete subdomain
func FetchHosts(ctx context.Context, domain, apikey string) *ZoomeyeHost {
	searchURL := fmt.Sprintf("%sdomain/search?q=%s&type=1&s=1000&page=1", zoomeyeURL, domain)
	header := map[string]string{
		"API-KEY": apikey,
	}
	_, body, err := clients.NewRequest("GET", searchURL, header, nil, 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Debug(ctx, err)
		return nil
	}
	var zh ZoomeyeHost
	json.Unmarshal(body, &zh)
	return &zh
}
