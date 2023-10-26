package module

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"regexp"
	"slack/common"
	"slack/common/client"
	"slack/gui/custom"
	"slack/lib/util"
	"slack/plugins/subdomain"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func IPAndCDN() *fyne.Container {
	entry := custom.NewMultiLineEntryPlaceHolder(`IP提取:
输入任意内容会自动匹配IPV4地址会进行提取并统计C段数量

域名解析(CDN查询):
输入任意内容会自动匹配域名进行IP解析以及CDN判断(可批量)
===如果域名解析的非常慢，请考虑是否是本机网络不佳===

域名备案&whois查询(数据来源: 站长之家)(不可批量)

IP转换
根据右边列表的形式进行IP转换
`)
	result := widget.NewMultiLineEntry()
	result.ActionItem = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		result.SetText("")
	})
	list := &widget.Select{Options: []string{"192.168.0.0/255.255.255.0 -> 192.168.0.0/24", "192.168.0.0/24 -> 192.168.0.0-192.168.0.255"},
		Alignment: fyne.TextAlignCenter, PlaceHolder: "192.168.0.0/255.255.255.0 -> 192.168.0.0/24"}
	ipconvert := widget.NewButton("IP转换", func() {
		go FormatIP(entry.Text, list.SelectedIndex(), result)
	})
	domainresolution := widget.NewButton("域名解析(CDN查询)", func() {
		go CdnCheck(entry, result)
	})
	ipfilter := widget.NewButton("IP筛选统计", func() {
		result.Text += Analysis(entry.Text)
		result.Refresh()
	})
	whois := widget.NewButton("域名备案&whois查询", func() {
		if util.RegDomain.FindString(entry.Text) != "" { //判断域名是否合规
			go func() {
				s, s2, s3 := SeoChinaz(entry.Text)
				s4, s5, s6, s7, s8, s9 := WhoisChinaz(entry.Text)
				result.SetText(fmt.Sprintf(`---备案查询---
公司名称: %v
	
备案号: %v
	
IP: %v
	
---whois查询---
更新时间: %v

创建时间: %v

过期时间: %v

注册商服务器: %v

DNS: %v

状态: %v`, s, s2, s3, s4, s5, s6, s7, s8, s9))
			}()
		}
	})
	center := container.NewBorder(nil, nil, container.NewGridWithColumns(4, domainresolution, whois, ipfilter, ipconvert), nil, list)

	vbox := container.NewVSplit(container.NewBorder(nil, center, nil, nil, entry), result)
	vbox.SetOffset(0.4)
	return container.NewBorder(nil, nil, nil, nil, vbox)
}

func FormatIP(text string, mode int, result *widget.Entry) {
	result.Text = ""
	if mode == 0 {
		for _, ipmask := range util.RegIPCompleteMask.FindAllString(text, -1) {
			ip := net.ParseIP(strings.Split(ipmask, "/")[0])
			subnetMask := net.IPMask(net.ParseIP(strings.Split(ipmask, "/")[1]).To4())
			// 合并IP地址和子网掩码
			ipNet := &net.IPNet{
				IP:   ip,
				Mask: subnetMask,
			}
			// 获取CIDR表示法
			result.Text += ipNet.String() + "\n"
			result.Refresh()
		}
	} else {
		for _, ipmask := range util.RegIPCIDR.FindAllString(text, -1) {
			_, ipNet, _ := net.ParseCIDR(ipmask)
			start := ipNet.IP
			end := net.IP(make([]byte, len(ipNet.IP)))
			copy(end, start)
			for i := len(end) - 1; i >= 0; i-- {
				if ipNet.Mask[i] != 0 {
					end[i] |= ^ipNet.Mask[i]
				}
			}
			result.Text += fmt.Sprintf("%s-%s\n", start.String(), end.String())
			result.Refresh()
		}
	}
}

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

func CdnCheck(entry, result *widget.Entry) {
	result.Text += "---域名解析(CDN查询)---:\n"
	result.Refresh()
	go func() {
	outerLoop:
		for _, domain := range util.RemoveDuplicates[string](util.RegDomain.FindAllString(entry.Text, -1)) {
			ips, cnames, err := subdomain.Resolution(domain, []string{common.Profile.Subdomain.DNS1 + ":53", common.Profile.Subdomain.DNS2 + ":53"})
			if err == nil {
				for name, cdns := range subdomain.ReadCDNFile() {
					for _, cdn := range cdns {
						for _, cname := range cnames {
							if strings.Contains(cname, cdn) { // 识别到cdn
								result.Text += fmt.Sprintf("域名: %v 识别到CDN域名，CNAME: %v CDN名称: %v 解析到IP为: %v\n", domain, cname, name, strings.Join(ips, ","))
								result.Refresh()
								break outerLoop
							} else if strings.Contains(cname, "cdn") {
								result.Text += fmt.Sprintf("域名: %v CNAME中含有关键字: cdn，该域名可能使用了CDN技术 CNAME: %v 解析到IP为: %v \n", domain, cname, strings.Join(ips, ","))
								result.Refresh()
								break outerLoop
							}
						}
					}
				}
				result.Text += fmt.Sprintf("域名: %v 解析到IP为: %v\n", domain, strings.Join(ips, ","))
				result.Refresh()
			} else {
				result.Text += fmt.Sprintf("域名: %v 解析失败,%v\n", domain, err)
				result.Refresh()
			}
		}
		result.Text += "\n\n"
		result.Refresh()
	}()
}

func SeoChinaz(domain string) (string, string, string) {
	var (
		r1 = regexp.MustCompile("license =  '(.*?)' ;")
		r2 = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
		r3 = regexp.MustCompile(`t0-p0-c0-i0-d0-s-(.*?)" target="`)
	)
	r, err := http.NewRequest("GET", "https://seo.chinaz.com/"+domain, nil)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.97")
	r.Header.Add("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		custom.Console.Append("[ERR] " + err.Error() + "\n")
	}
	c := client.DefaultClient()
	rx, err2 := c.Do(r)
	if err2 != nil {
		custom.Console.Append("[ERR] " + err2.Error() + "\n")
	}
	b, err3 := io.ReadAll(rx.Body)
	if err3 != nil {
		custom.Console.Append("[ERR] " + err3.Error() + "\n")
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

func WhoisChinaz(domain string) (string, string, string, string, string, string) {
	r, err := http.NewRequest("GET", "https://whois.chinaz.com/"+domain, nil)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.97")
	r.Header.Add("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		custom.Console.Append("[ERR] " + err.Error() + "\n")
	}
	c := client.DefaultClient()
	rx, err2 := c.Do(r)
	if err2 != nil {
		custom.Console.Append("[ERR] " + err2.Error() + "\n")
	}
	b, err3 := io.ReadAll(rx.Body)
	if err3 != nil {
		custom.Console.Append("[ERR] " + err3.Error() + "\n")
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
