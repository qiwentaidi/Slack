package core

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
)

func Analysis(entry string) (result string) {
	var IP_analysis = make(map[string]int)
	result += "---提取IP资产---\n"
	for _, ip := range util.RemoveDuplicates[string](util.RegIP.FindAllString(entry, -1)) {
		result += ip + "\n"
		ip = ip[:len(ip)-len(path.Ext(ip))]
		IP_analysis[ip+".0"]++
	}
	result += "\n\n\n---提取C段资产---\n"
	for _, p := range util.SortMap(IP_analysis) {
		result += fmt.Sprintf("%v/24(%v)\n", p.Key, p.Value)
	}
	return result
}

func SeoChinaz(ctx context.Context, domain string) (string, string, string) {
	var (
		r1 = regexp.MustCompile("license =  '(.*?)' ;")
		r2 = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
		r3 = regexp.MustCompile(`t0-p0-c0-i0-d0-s-(.*?)" target="`)
	)
	h := http.Header{}
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.97")
	h.Add("Content-Type", "text/html; charset=utf-8")
	_, b, err := clients.NewRequest("GET", "https://seo.chinaz.com/"+domain, h, nil, 10, true, clients.DefaultClient())
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
