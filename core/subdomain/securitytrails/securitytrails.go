package securitytrails

import (
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
)

const securitytrailsURL = "https://api.securitytrails.com/"

type SecuritytrailsHost struct {
	Endpoint        string   `json:"endpoint"`
	Meta            Meta     `json:"meta"`
	Subdomain_count float64  `json:"subdomain_count"`
	Subdomains      []string `json:"subdomains"`
}

type Meta struct {
	Limit_reached bool `json:"limit_reached"`
}

// subdomains return is not complete subdomain, e.g www
func FetchHosts(ctx context.Context, domain, apikey string) *SecuritytrailsHost {
	header := map[string]string{
		"accept": "application/json",
		"APIKEY": apikey,
	}
	searchURL := fmt.Sprintf("%sv1/domain/%s/subdomains?children_only=false&include_inactive=true", securitytrailsURL, domain)
	resp, err := clients.DoRequest("GET", searchURL, header, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		fmt.Println(err)
	}
	var sh SecuritytrailsHost
	json.Unmarshal(resp.Body(), &sh)
	return &sh
}
