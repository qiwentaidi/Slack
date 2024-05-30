package waf

import (
	"slack-wails/core"
	"slack-wails/lib/util"
	"strings"
)

var (
	WAFs = map[string][]string{
		"sanfor-shield":  {".sangfordns.com"},
		"360panyun":      {".360panyun.com"},
		"baiduyun":       {".yunjiasu-cdn.net"},
		"chuangyudun":    {".365cyd.cn", ".cyudun.net", ".365cyd.net"},
		"knownsec":       {".jiashule.com", ".jiasule.org"},
		"huaweicloud":    {".huaweicloudwaf.com"},
		"xinliuyun":      {".ngaagslb.cn"},
		"chinacache":     {".chinacache.net", ".ccgslb.net", ".chinacache.com"},
		"nscloudwaf":     {".nscloudwaf.com"},
		"wangsu":         {".wsssec.com", ".lxdns.com", ".wscdns.com", ".cdn20.com", ".cdn30.com", ".ourplat.net", ".wsdvs.com", ".wsglb0.com", ".wswebcdn.com", ".wswebpic.com", ".wsssec.com", ".wscloudcdn.com", ".mwcloudcdn.com", ".chinanetcenter.com"},
		"qianxin":        {".360safedns.com", ".360cloudwaf.com", ".360wzb.com"},
		"baiduyunjiasu":  {".yunjiasu-cdn.net"},
		"anquanbao":      {".anquanbao.net", ".anquanbao.com"},
		"aliyun":         {"kunlun", "aliyunddos", "aliyunwaf", "aligaofang", "aliyundunwaf"},
		"xuanwudun":      {".saaswaf.com", ".dbappwaf.cn"},
		"yundun":         {".hwwsdns.cn", ".yunduncname.com"},
		"knownsec-ns":    {".jiasule.net"},
		"baiduyunjiasue": {".ns.yunjiasu.com"},
		"cloudflare":     {".ns.cloudflare.com"},
		"edns":           {".iidns.com"},
		"ksyun":          {".ksyunwaf.com"},
	}
)

type WAF struct {
	Exsits bool
	Name   string
}

func IsWAF(host string) *WAF {
	// 如果是IP则直接返回
	if util.RegIP.MatchString(host) {
		return &WAF{}
	}
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	cnames, err := core.LookupCNAME(host, core.DnsServers, 10)
	if err != nil || len(cnames) == 0 {
		return &WAF{}
	}
	for name, domains := range WAFs {
		for _, domain := range domains {
			if strings.Contains(cnames[0], domain) {
				return &WAF{Exsits: true, Name: name}
			}
		}
	}
	return &WAF{}
}
