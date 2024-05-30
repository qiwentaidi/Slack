package core

import (
	"fmt"
	"net/http"
	"path"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"

	"github.com/wailsapp/wails/v2/pkg/logger"
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

func WhoisChinaz(domain string) (string, string, string, string, string, string) {
	h := http.Header{}
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.97")
	h.Add("Content-Type", "text/html; charset=utf-8")
	_, b, err := clients.NewRequest("GET", "https://whois.chinaz.com/"+domain, h, nil, 10, true, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	html := string(b)
	updateTimeRegex := regexp.MustCompile(`<div class="fl WhLeList-left">更新时间</div>\s*<div class="fr WhLeList-right">\s*<span>(.*?)</span>`)
	creationTimeRegex := regexp.MustCompile(`<div class="fl WhLeList-left">创建时间</div>\s*<div class="fr WhLeList-right">\s*<span>(.*?)</span>`)
	expirationTimeRegex := regexp.MustCompile(`<div class="fl WhLeList-left">过期时间</div>\s*<div class="fr WhLeList-right">\s*<span>(.*?)</span>`)
	registrarServerRegex := regexp.MustCompile(`<div class="fl WhLeList-left">注册商服务器</div>\s*<div class="fr WhLeList-right">\s*<span>(.*?)</span>`)
	dnsRegex := regexp.MustCompile(`<div class="fl WhLeList-left">DNS</div>\s*<div class="fr WhLeList-right">\s*<span>(.*?)</span>`)
	statusRegex := regexp.MustCompile(`<div class="fl WhLeList-left">状态</div>\s*<div class="fr WhLeList-right clearfix">\s*<span>(.*?)</span>`)
	ut := updateTimeRegex.FindStringSubmatch(html)
	ct := creationTimeRegex.FindStringSubmatch(html)
	et := expirationTimeRegex.FindStringSubmatch(html)
	rs := registrarServerRegex.FindStringSubmatch(html)
	dns := dnsRegex.FindStringSubmatch(html)
	server := statusRegex.FindStringSubmatch(html)
	if len(ut) > 0 && len(ct) > 0 && len(et) > 0 && len(rs) > 0 && len(dns) > 0 && len(server) > 0 {
		return ut[1], ct[1], et[1], rs[1], dns[1], server[1]
	}
	infoRegex := regexp.MustCompile(`<p id="info" class="col-red">(.*?)</p>`)
	info := infoRegex.FindStringSubmatch(html)
	if len(info) > 0 {
		return info[1], "", "", "", "", ""
	}
	return "", "", "", "", "", ""
}

func SeoChinaz(domain string) (string, string, string) {
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
		logger.NewDefaultLogger().Debug(err.Error())
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
