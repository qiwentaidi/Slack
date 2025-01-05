package space

import (
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
)

type ShodanIpReport struct {
	Ports     []float64 `json:"ports"`
	Tags      []string  `json:"tags"`
	Vulns     []string  `json:"vulns"`
	Cpes      []string  `json:"cpes"`
	Hostnames []string  `json:"hostnames"`
	Ip        string    `json:"ip"`
}

func GetShodanAllPort(ctx context.Context, ip string) []float64 {
	_, body, err := clients.NewSimpleGetRequest("https://internetdb.shodan.io/"+ip, clients.NewHttpClient(nil, true))
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("[shodan] %s, error: %v", ip, err))
		return []float64{}
	}
	var shodan ShodanIpReport
	json.Unmarshal(body, &shodan)
	return shodan.Ports
}
