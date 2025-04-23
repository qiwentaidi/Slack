package info

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"

	"github.com/PuerkitoBio/goquery"
)

var (
	r1 = regexp.MustCompile("license =  '(.*?)' ;")
	r2 = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	r3 = regexp.MustCompile(`t0-p0-c0-i0-d0-s-(.*?)" target="`)
)

func SeoChinaz(ctx context.Context, domain string) (string, string, string) {
	h := map[string]string{
		"Content-Type": "text/html; charset=utf-8",
	}
	resp, err := clients.DoRequest("GET", "https://seo.chinaz.com/"+domain, h, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Warning(ctx, err)
	}
	html := string(resp.Body())
	registration := r1.FindStringSubmatch(html)
	ip := r2.FindString(html)
	company := r3.FindStringSubmatch(html)
	if len(registration) > 0 && len(company) > 0 {
		return company[1], registration[1], ip
	}
	return "疑似违规域名不可查询", "疑似违规域名不可查询", ip
}

func Ip138IpHistory(domain string) string {
	resp, err := clients.SimpleGet(fmt.Sprintf("https://site.ip138.com/%s/", domain), clients.NewRestyClient(nil, true))
	if err != nil {
		return "ip138网站响应异常!"
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "解析网站失败"
	}
	var ipHistory string
	doc.Find("#J_ip_history p").Each(func(i int, s *goquery.Selection) {
		// date := s.Find("span.date").Text()
		ip := s.Find("a").Text()
		if ip != "" {
			// ipHistory += fmt.Sprintf("%s\t%s\n", ip, date)
			ipHistory += ip + "\n"
		}
	})

	return ipHistory
}

func Ip138Subdomain(domain string) string {
	resp, err := clients.SimpleGet(fmt.Sprintf("https://site.ip138.com/%s/domain.htm", domain), clients.NewRestyClient(nil, true))
	if err != nil {
		return "ip138网站响应异常!"
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return "解析网站失败"
	}
	var subdomain string
	doc.Find("#J_subdomain p").Each(func(i int, s *goquery.Selection) {
		domain := s.Find("a").Text()
		if domain != "" {
			subdomain += domain + "\n"
		}
	})
	return subdomain
}
