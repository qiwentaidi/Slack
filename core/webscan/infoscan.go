package webscan

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	rt "runtime"
	"slack-wails/core/subdomain"
	"slack-wails/core/waf"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/gonmap"
	"slack-wails/lib/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	iconRels = []string{"icon", "shortcut icon", "apple-touch-icon", "mask-icon"}
	ExitFunc = false
)

type TargetINFO struct {
	Protocol      string
	Port          int
	Path          string
	Title         string
	StatusCode    int
	IconHash      string // mmh3
	BodyString    string
	HeadeString   string
	ContentType   string
	Server        string
	ContentLength int
	Banner        string // tcp指纹
	Cert          string // TLS证书
	Waf           waf.WAF
}

type InfoResult struct {
	URL          string
	StatusCode   int
	Length       int
	Title        string
	Fingerprints []string
	IsWAF        bool
	WAF          string
}

type FingerScanner struct {
	ctx   context.Context
	urls  []*url.URL
	proxy clients.Proxy
}

func NewFingerScanner(ctx context.Context, target []string, proxy clients.Proxy) *FingerScanner {
	urls := []*url.URL{}
	for _, t := range target {
		u, err := url.Parse(t)
		if err != nil {
			continue
		}
		urls = append(urls, u)
	}
	return &FingerScanner{
		ctx:   ctx,
		urls:  urls,
		proxy: proxy,
	}
}

func (fs *FingerScanner) NewFingerScan() {
	var wg sync.WaitGroup
	client := clients.JudgeClient(fs.proxy)
	single := make(chan struct{})
	retChan := make(chan InfoResult, len(fs.urls))
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(fs.ctx, "webFingerScan", pr)
		}
		close(single)
	}()
	// 指纹扫描
	fscan := func(u *url.URL) {
		if ExitFunc {
			return
		}
		resp, body, _ := clients.NewSimpleGetRequest(u.String(), client)
		if resp == nil {
			retChan <- InfoResult{
				URL:        u.String(),
				StatusCode: 0,
			}
			return
		}
		title, server, content_type := GetHeaderInfo(body, resp)
		headers, _, _ := DumpResponseHeadersAndRaw(resp)
		ti := &TargetINFO{
			HeadeString:   string(headers),
			ContentType:   content_type,
			Cert:          GetTLSString(u.Scheme, u.Host),
			BodyString:    string(body),
			Path:          u.Path,
			Title:         title,
			Server:        server,
			ContentLength: len(body),
			Port:          fs.GetPort(u),
			IconHash:      FaviconHash(u, client),
			StatusCode:    resp.StatusCode,
			Banner:        fs.GetBanner(u),
			Waf:           *waf.IsWAF(u.Hostname(), subdomain.DefaultDnsServers),
		}
		retChan <- InfoResult{
			URL:          u.String(),
			StatusCode:   ti.StatusCode,
			Length:       ti.ContentLength,
			Title:        ti.Title,
			Fingerprints: FingerScan(fs.ctx, ti, FingerprintDB),
			IsWAF:        ti.Waf.Exsits,
			WAF:          ti.Waf.Name,
		}
	}
	threadPool, _ := ants.NewPoolWithFunc(50, func(target interface{}) {
		t := target.(*url.URL)
		fscan(t)
		wg.Done()
	})
	defer threadPool.Release()
	for _, target := range fs.urls {
		if ExitFunc {
			return
		}
		wg.Add(1)
		threadPool.Invoke(target)
	}
	wg.Wait()
	close(retChan)
	gologger.Info(fs.ctx, "FingerScan Finished")
	<-single
}

type TFP struct {
	Target      string
	Fingerprint string
	Path        string
}

func (fs *FingerScanner) NewActiveFingerScan() {
	var wg sync.WaitGroup
	var matched bool
	client := clients.JudgeClient(fs.proxy)
	single := make(chan struct{})
	retChan := make(chan InfoResult, len(fs.urls))
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(fs.ctx, "webFingerScan", pr)
		}
		close(single)
		// runtime.EventsEmit(ctx, "webFingerScanComplete", "done")
	}()
	// 主动指纹扫描
	threadPool, _ := ants.NewPoolWithFunc(10, func(tfp interface{}) {
		defer wg.Done()
		fp := tfp.(TFP)
		u, _ := url.Parse(fp.Target)
		resp, body, _ := clients.NewSimpleGetRequest(fp.Target+fp.Path, client)
		if resp == nil {
			return
		}
		title, server, content_type := GetHeaderInfo(body, resp)
		headers, _, _ := DumpResponseHeadersAndRaw(resp)
		ti := &TargetINFO{
			HeadeString:   string(headers),
			ContentType:   content_type,
			Cert:          "",
			BodyString:    string(body),
			Path:          u.Path,
			Title:         title,
			Server:        server,
			ContentLength: len(body),
			Port:          fs.GetPort(u),
			IconHash:      "",
			StatusCode:    resp.StatusCode,
			Banner:        "",
		}
		result := FingerScan(fs.ctx, ti, ActiveFingerprintDB)
		// 多路径匹配时如果某一路径匹配到就立刻停止
		if len(result) > 0 && fp.Fingerprint == result[0] {
			retChan <- InfoResult{
				URL:          fp.Target + fp.Path,
				StatusCode:   ti.StatusCode,
				Length:       ti.ContentLength,
				Title:        ti.Title,
				Fingerprints: []string{fp.Fingerprint},
			}
			matched = true
		}
	})
	defer threadPool.Release()
	for _, target := range fs.urls {
		for fingername, paths := range Sensitive {
			matched = false
			for _, path := range paths {
				if ExitFunc { // 控制程序退出
					return
				}
				if matched { // 如果已经有匹配成功的指纹需要跳出当层目录循环
					break
				}
				wg.Add(1)
				threadPool.Invoke(TFP{
					Target:      target.String(),
					Fingerprint: fingername,
					Path:        path,
				})
			}
		}
	}
	wg.Wait()
	close(retChan)
	gologger.Info(fs.ctx, "ActiveFingerScan Finished")
	<-single
}

// DumpResponseHeadersAndRaw returns http headers and response as strings
func DumpResponseHeadersAndRaw(resp *http.Response) (headers, fullresp []byte, err error) {
	// httputil.DumpResponse does not work with websockets
	if resp.StatusCode >= http.StatusContinue && resp.StatusCode <= http.StatusEarlyHints {
		raw := resp.Status + "\n"
		for h, v := range resp.Header {
			raw += fmt.Sprintf("%s: %s\n", h, v)
		}
		return []byte(raw), []byte(raw), nil
	}
	headers, err = httputil.DumpResponse(resp, false)
	if err != nil {
		return
	}
	// logic same as httputil.DumpResponse(resp, true) but handles
	// the edge case when we get both error and data on reading resp.Body
	var buf1, buf2 bytes.Buffer
	b := resp.Body
	if _, err = buf1.ReadFrom(b); err != nil {
		if buf1.Len() <= 0 {
			return
		}
	}
	if err == nil {
		_ = b.Close()
	}

	// rewind the body to allow full dump
	resp.Body = io.NopCloser(bytes.NewReader(buf1.Bytes()))
	err = resp.Write(&buf2)
	fullresp = buf2.Bytes()

	// rewind once more to allow further reuses
	resp.Body = io.NopCloser(bytes.NewReader(buf1.Bytes()))
	return
}

func FingerScan(ctx context.Context, ti *TargetINFO, targetDB []FingerPEntity) []string {
	var fingerPrintResults []string

	isWeb := ti.Path != "no#web"
	workers := rt.NumCPU() * 2
	inputChan := make(chan FingerPEntity, len(targetDB))
	defer close(inputChan)
	results := make(chan string, len(targetDB))
	defer close(results)

	var wg sync.WaitGroup

	//接收结果
	go func() {
		for found := range results {
			if found != "" {
				fingerPrintResults = append(fingerPrintResults, found)
			}
			wg.Done()
		}
	}()

	//多指纹同时扫描
	for i := 0; i < workers; i++ {
		go func() {
			for finger := range inputChan {
				rules := finger.Rule
				product := finger.ProductName
				expr := finger.AllString
				for _, singleRule := range rules {
					singleRuleResult := false
					if singleRule.Key == "header" {
						if isWeb && dataCheckString(singleRule.Op, ti.HeadeString, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "body" {
						if isWeb && dataCheckString(singleRule.Op, ti.BodyString, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "server" {
						if isWeb && dataCheckString(singleRule.Op, ti.Server, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "title" {
						if isWeb && dataCheckString(singleRule.Op, ti.Title, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "cert" {
						if dataCheckString(singleRule.Op, ti.Cert, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "port" {
						value, err := strconv.Atoi(singleRule.Value)
						if err == nil && dataCheckInt(singleRule.Op, ti.Port, value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "protocol" {
						if singleRule.Op == 0 {
							if ti.Protocol == singleRule.Value {
								singleRuleResult = true
							}
						} else if singleRule.Op == 1 {
							if ti.Protocol != singleRule.Value {
								singleRuleResult = true
							}
						}
					} else if singleRule.Key == "path" {
						if isWeb && dataCheckString(singleRule.Op, ti.Path, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "icon_hash" {
						value, err := strconv.Atoi(singleRule.Value)
						hashIcon, errHash := strconv.Atoi(ti.IconHash)
						if isWeb && err == nil && errHash == nil && dataCheckInt(singleRule.Op, hashIcon, value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "status" {
						value, err := strconv.Atoi(singleRule.Value)
						if isWeb && err == nil && dataCheckInt(singleRule.Op, ti.StatusCode, value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "content_type" {
						if isWeb && dataCheckString(singleRule.Op, ti.ContentType, singleRule.Value) {
							singleRuleResult = true
						}
					} else if singleRule.Key == "banner" {
						if dataCheckString(singleRule.Op, ti.Banner, singleRule.Value) {
							singleRuleResult = true
						}
					}
					if singleRuleResult {
						expr = expr[:singleRule.Start] + "T" + expr[singleRule.End:]
					} else {
						expr = expr[:singleRule.Start] + "F" + expr[singleRule.End:]
					}
				}

				r := boolEval(ctx, expr)
				if r {
					results <- product
				} else {
					results <- ""
				}

			}

		}()
	}
	//添加扫描目标
	for _, input := range targetDB {
		wg.Add(1)
		inputChan <- input
	}
	wg.Wait()
	return util.RemoveDuplicates(fingerPrintResults)
}

// parseIcons 解析HTML文档head中的<link>标签中rel属性包含icon信息的href链接
func parseIcons(doc *goquery.Document) []string {
	var icons []string
	doc.Find("head link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// 匹配ICON链接
			if rel, exists := s.Attr("rel"); exists && util.ArrayContains(rel, iconRels) {
				icons = append(icons, href)
			}
		}
	})
	// 找不到自定义icon链接就使用默认的favicon地址
	if len(icons) == 0 {
		icons = append(icons, "favicon.ico")
	}
	return icons
}

// 获取favicon hash值
func FaviconHash(u *url.URL, client *http.Client) string {
	_, body, err := clients.NewSimpleGetRequest(u.String(), client)
	if err != nil {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return ""
	}
	iconLink := parseIcons(doc)[0]
	var finalLink string
	// 如果是完整的链接，则直接请求
	if strings.HasPrefix(iconLink, "http") {
		finalLink = iconLink
		// 如果为 // 开头采用与网站同协议
	} else if strings.HasPrefix(iconLink, "//") {
		finalLink = u.Scheme + ":" + iconLink
	} else {
		finalLink = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, iconLink)
	}
	resp, body, err := clients.NewSimpleGetRequest(finalLink, client)
	if err == nil && resp.StatusCode == 200 {
		return util.Mmh3Hash32(util.Base64Encode(body))
	}
	return ""
}
func (fs *FingerScanner) GetPort(u *url.URL) int {
	if u.Port() == "" {
		switch u.Scheme {
		case "http":
			return 80
		case "https":
			return 443
		}
	}
	port, _ := strconv.Atoi(u.Port())
	return port
}

func (fs *FingerScanner) GetBanner(u *url.URL) string {
	if strings.HasPrefix(u.Scheme, "http") {
		return ""
	}
	scanner := gonmap.New()
	_, response := scanner.Scan(u.Host, fs.GetPort(u), time.Second*time.Duration(10))
	if response != nil {
		return response.Raw
	}
	return ""
}

func GetHeaderInfo(body []byte, resp *http.Response) (title, server, content_type string) {
	if match := util.RegTitle.FindSubmatch(body); len(match) > 1 {
		title = util.Str2UTF8(string(match[1]))
	}
	for k, v := range resp.Header {
		if k == "Server" {
			server = strings.Join(v, ";")
		}
		if k == "Content-Type" {
			content_type = strings.Join(v, ";")
		}
	}
	return
}
