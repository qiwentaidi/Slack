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
	_, b, err := clients.NewRequest("GET", "https://seo.chinaz.com/"+domain, h, nil, 10, true, clients.NewHttpClient(nil, true))
	if err != nil {
		gologger.Warning(ctx, err)
	}
	html := string(b)
	registration := r1.FindStringSubmatch(html)
	ip := r2.FindString(html)
	company := r3.FindStringSubmatch(html)
	if len(registration) > 0 && len(company) > 0 {
		return company[1], registration[1], ip
	}
	return "疑似违规域名不可查询", "疑似违规域名不可查询", ip
}

func Ip138IpHistory(domain string) string {
	_, body, err := clients.NewSimpleGetRequest(fmt.Sprintf("https://site.ip138.com/%s/", domain), clients.NewHttpClient(nil, true))
	if err != nil {
		return "ip138网站响应异常!"
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
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
	_, body, err := clients.NewSimpleGetRequest(fmt.Sprintf("https://site.ip138.com/%s/domain.htm", domain), clients.NewHttpClient(nil, true))
	if err != nil {
		return "ip138网站响应异常!"
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
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
