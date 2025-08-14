package ip138

import (
	"bytes"
	"context"
	"fmt"
	"slack-wails/lib/gologger"

	"github.com/qiwentaidi/clients"

	"github.com/PuerkitoBio/goquery"
)

func FetchHosts(ctx context.Context, domain string) []string {
	result := []string{}
	resp, err := clients.SimpleGet(fmt.Sprintf("https://site.ip138.com/%s/domain.htm", domain), clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.DualLog(ctx, gologger.Level_ERROR, "[subdomain] ip138网站响应异常!")
		return result
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		gologger.DualLog(ctx, gologger.Level_ERROR, "[subdomain] ip138解析网站失败")
		return result
	}
	doc.Find("#J_subdomain p").Each(func(i int, s *goquery.Selection) {
		domain := s.Find("a").Text()
		if domain != "" {
			result = append(result, domain)
		}
	})
	return result
}
