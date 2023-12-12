package runner

import (
	"slack-wails/core/webscan/poc"
	"slack-wails/lib/util"
	"strings"
)

type Options struct {
	Targets util.SafeSlice

	// target URLs/hosts to scan
	Target string

	// PoC file or directory to scan
	PocFile string

	// search PoC by keyword , eg: -s tomcat
	Search string

	SearchKeywords []string

	// pocs to run based on severity. Possible values: info, low, medium, high, critical
	Severity string

	SeverityKeywords []string

	// maximum number of requests to send per second (default 150)
	RateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	Concurrency int

	// maximum number of requests to send per second (default 150)
	ReverseRateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	ReverseConcurrency int

	MaxHostError int

	// time to wait in seconds before timeout (default 10)
	Timeout int

	// http/socks5 proxy to use
	Proxy string
}

func NewOptions(target, keyword, severity, proxy string) *Options {
	options := &Options{}
	options.Target = target
	options.Search = keyword
	options.Severity = severity
	options.RateLimit = 150
	options.Concurrency = 25
	options.ReverseRateLimit = 50
	options.ReverseConcurrency = 20
	options.Timeout = 10
	options.MaxHostError = 3
	options.Proxy = proxy
	return options
}

// 处理关键字
func (o *Options) SetSearchKeyword() {
	if len(o.Search) > 0 {
		for _, v := range strings.Split(o.Search, ",") {
			o.SearchKeywords = append(o.SearchKeywords, strings.TrimSpace(v))
		}
	}
}

// 处理风险等级
func (o *Options) SetSeverityKeyword() {
	o.SeverityKeywords = append(o.SeverityKeywords, strings.Split(o.Severity, ",")...)
}

// 判断poc id 或者 name值是否匹配关键字
func (o *Options) CheckPocKeywords(id, name string) bool {
	for _, v := range o.SearchKeywords {
		v = strings.ToLower(v)
		if strings.Contains(strings.ToLower(id), v) || strings.Contains(strings.ToLower(name), v) {
			return true
		}
	}
	return false
}

// 关键字筛选
func (o *Options) CheckPocSeverityKeywords(severity string) bool {
	if len(o.SeverityKeywords) > 0 {
		for _, v := range o.SeverityKeywords {
			if strings.EqualFold(severity, v) {
				return true
			}
		}
	}
	return false
}

// *适用全部漏洞扫描*
// 通过关键字和风险等级筛选POC
func (o *Options) FilterPocSeveritySearch(pocId, pocName, severity string) bool {
	var isShow bool
	o.SetSearchKeyword()
	o.SetSeverityKeyword()
	if len(o.Search) > 0 && len(o.Severity) > 0 {
		if o.CheckPocKeywords(pocId, pocName) && o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Severity) > 0 {
		if o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Search) > 0 {
		if o.CheckPocKeywords(pocId, pocName) {
			isShow = true
		}
	} else {
		isShow = true
	}
	return isShow
}

// 区分dnslog poc
func (o *Options) ReversePoCs(allpocs []poc.Poc) ([]poc.Poc, []poc.Poc) {
	reverse := []poc.Poc{}
	other := []poc.Poc{}
	for _, poc := range allpocs {
		flag := false
		for _, item := range poc.Set {
			key := item.Key.(string)
			if key == "reverse" {
				flag = true
				break
			}
		}
		if flag {
			reverse = append(reverse, poc)
		} else {
			other = append(other, poc)
		}
	}
	return reverse, other
}

// 返回目标所要扫描的POC
func (o *Options) CreatePocList(pocpathList []string) []poc.Poc {
	var pocSlice []poc.Poc
	for _, pocEmbedYaml := range pocpathList {
		if p, err := poc.LocalReadPocByPath(pocEmbedYaml); err == nil {
			pocSlice = append(pocSlice, p)
		}
	}
	// 筛选
	if len(o.Search) > 0 || len(o.Severity) > 0 {
		newPocSlice := []poc.Poc{}
		for _, poc := range pocSlice {
			if o.FilterPocSeveritySearch(poc.Id, poc.Info.Name, poc.Info.Severity) {
				newPocSlice = append(newPocSlice, poc)
			}
		}
		return newPocSlice
	}
	return pocSlice
}
