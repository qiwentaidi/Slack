package common

import (
	"strings"

	"slack/lib/poc"

	"slack/lib/util"
)

var (
	ReverseCeyeApiKey = "ba446c3277a60555ad9e74a6f0cb4290"
	ReverseCeyeDomain = "xrn0nb.ceye.io"

	ReverseCeyeLive bool
)

type Options struct {
	// Pocs Directory
	PocsDirectory []string

	Targets util.SafeSlice

	// target URLs/hosts to scan
	Target []string

	// PoC file or directory to scan
	PocFile string

	// search PoC by keyword , eg: -s tomcat
	Search string

	SearchKeywords []string

	// pocs to run based on severity. Possible values: info, low, medium, high, critical
	Severity string

	SeverityKeywords []string

	// Scan count num(targets * allpocs)
	Count int

	// Current Scan count num
	CurrentCount uint32

	// maximum number of requests to send per second (default 150)
	RateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	Concurrency int

	// maximum number of requests to send per second (default 150)
	ReverseRateLimit int

	// maximum number of afrog-pocs to be executed in parallel (default 25)
	ReverseConcurrency int

	// number of times to retry a failed request (default 1)
	Retries int

	MaxHostError int

	// time to wait in seconds before timeout (default 10)
	Timeout int

	// http/socks5 proxy to use
	Proxy string

	MaxRespBodySize int
}

func NewOptions(target []string, keyword, severity, proxy string) *Options {
	options := &Options{}
	options.Target = target
	options.Search = keyword
	options.Severity = severity
	options.RateLimit = Profile.WebScan.Thread
	options.Concurrency = 25
	options.ReverseRateLimit = 50
	options.ReverseConcurrency = 20
	options.MaxRespBodySize = 2
	options.Retries = 1
	options.Timeout = DefaultWebTimeout
	options.MaxHostError = 3
	options.Proxy = proxy
	poc.SelectFolderReadPocPath(options.PocFile)
	return options
}

// 判断关键词是单个还是多个
func (o *Options) SetSearchKeyword() bool {
	if len(o.Search) > 0 {
		arr := strings.Split(o.Search, ",")
		if len(arr) > 0 {
			for _, v := range arr {
				o.SearchKeywords = append(o.SearchKeywords, strings.TrimSpace(v))
			}
			return true
		}
	}
	return false
}

func (o *Options) CheckPocKeywords(id string) bool {
	if len(o.SearchKeywords) > 0 {
		for _, v := range o.SearchKeywords {
			v = strings.ToLower(v)
			if strings.Contains(strings.ToLower(id), v) {
				return true
			}
		}
	}
	return false
}

// 风险等级筛选
func (o *Options) SetSeverityKeyword() bool {
	if len(o.Severity) > 0 {
		arr := strings.Split(o.Severity, ",")
		if len(arr) > 0 {
			o.SeverityKeywords = append(o.SeverityKeywords, arr...)
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

// 对poc的关键字进行筛选搜索关键字筛选
func (o *Options) FilterPocSeveritySearch(pocId, severity string) bool {
	var isShow bool
	if len(o.Search) > 0 && o.SetSearchKeyword() && len(o.Severity) > 0 && o.SetSeverityKeyword() {
		if o.CheckPocKeywords(pocId) && o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Severity) > 0 && o.SetSeverityKeyword() {
		if o.CheckPocSeverityKeywords(severity) {
			isShow = true
		}
	} else if len(o.Search) > 0 && o.SetSearchKeyword() {
		if o.CheckPocKeywords(pocId) {
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

func (o *Options) CreatePocList(active, finger2poc bool) []poc.Poc {
	var pocSlice []poc.Poc
	// 如果LocalTestList数量大于0，说明是指定了文件夹那就只扫文件夹下面的POC
	if len(poc.LocalTestList) > 0 {
		for _, pocYaml := range poc.LocalTestList {
			if p, err := poc.LocalReadPocByPath(pocYaml); err == nil {
				pocSlice = append(pocSlice, p)
			}
		}
		return pocSlice
	}
	// 如果开启仅指纹扫描，或者指纹POC扫描，屏蔽其他POC，且优先级大于扫描指纹POC
	if !active && !finger2poc {
		for _, pocEmbedYaml := range poc.EmbedFileList {
			if p, err := poc.LocalReadPocByPath(pocEmbedYaml); err == nil {
				pocSlice = append(pocSlice, p)
			}
		}
	}

	newPocSlice := []poc.Poc{}
	for _, poc := range pocSlice {
		// 筛选
		if o.FilterPocSeveritySearch(poc.Id, poc.Info.Severity) {
			newPocSlice = append(newPocSlice, poc)
		}
	}
	return newPocSlice
}
