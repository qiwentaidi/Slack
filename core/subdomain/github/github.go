package github

import (
	"context"
	"fmt"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"strings"
)

// subdomains return is complete subdomain
func FetchHosts(ctx context.Context, domain, apikey string) []string {
	headers := map[string]string{
		"Accept":        "application/vnd.github.v3.text-match+json",
		"Authorization": "token " + apikey,
	}
	searchURL := fmt.Sprintf("https://api.github.com/search/code?per_page=100&q=%s&sort=created&order=asc", domain)
	_, body, err := clients.NewRequest("GET", searchURL, headers, nil, 10, true, clients.NewHttpClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, err)
	}
	r := domainRegexp(domain)
	return r.FindAllString(string(body), -1)
}

func domainRegexp(domain string) *regexp.Regexp {
	rdomain := strings.ReplaceAll(domain, ".", "\\.")
	return regexp.MustCompile("(\\w[a-zA-Z0-9][a-zA-Z0-9-\\.]*)" + rdomain)
}
