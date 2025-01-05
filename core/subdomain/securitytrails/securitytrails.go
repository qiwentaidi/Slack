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
	_, body, err := clients.NewRequest("GET", searchURL, header, nil, 10, true, clients.NewHttpClient(nil, true))
	if err != nil {
		fmt.Println(err)
	}
	var sh SecuritytrailsHost
	json.Unmarshal(body, &sh)
	return &sh
}
