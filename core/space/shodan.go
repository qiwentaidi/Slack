package space

import (
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/gologger"

	"github.com/qiwentaidi/clients"
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
	resp, err := clients.SimpleGet("https://internetdb.shodan.io/"+ip, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("[shodan] %s, error: %v", ip, err))
		return []float64{}
	}
	var shodan ShodanIpReport
	json.Unmarshal(resp.Body(), &shodan)
	return shodan.Ports
}
