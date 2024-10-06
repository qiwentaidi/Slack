package webscan

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
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
	"slack-wails/lib/netutil"
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

type WebInfo struct {
	Protocol      string
	Port          int
	Path          string
	Title         string
	StatusCode    int
	IconHash      string // mmh3
	IconMd5       string // md5
	BodyString    string
	HeadeString   string
	ContentType   string
	Server        string
	ContentLength int
	Banner        string // tcp指纹
	Cert          string // TLS证书
}

type InfoResult struct {
	URL          string
	StatusCode   int
	Length       int
	Title        string
	Fingerprints []string
	IsWAF        bool
	WAF          string
	Detect       string
	Screenshot   string // 截图图片路径
}

type FingerScanner struct {
	ctx                     context.Context
	urls                    []*url.URL
	aliveURLs               []*url.URL // 默认指纹扫描结束后，存活的URL，以便后续主动指纹过滤目标
	client                  *http.Client
	screenshot              bool
	thread                  int  // 指纹线程
	deepScan                bool // 代表主动指纹探测
	rootPath                bool // 主动指纹是否采取根路径扫描
	basicURLWithFingerprint map[string][]string
	mutex                   sync.RWMutex
}

func NewFingerScanner(ctx context.Context, target []string, proxy clients.Proxy, thread int, deepScan, rootPath, screenshot bool) *FingerScanner {
	urls := make([]*url.URL, 0, len(target)) // 提前分配容量
	for _, t := range target {
		u, err := url.Parse(t)
		if err != nil {
			continue
		}
		urls = append(urls, u)
	}
	if len(urls) == 0 {
		gologger.Error(ctx, "未发现可用目标，请检查输入")
		return nil
	}
	return &FingerScanner{
		ctx:                     ctx,
		urls:                    urls,
		client:                  clients.JudgeClient(proxy),
		screenshot:              screenshot,
		thread:                  thread,
		deepScan:                deepScan,
		rootPath:                rootPath,
		basicURLWithFingerprint: make(map[string][]string),
	}
}

func (s *FingerScanner) NewFingerScan() {
	var wg sync.WaitGroup
	single := make(chan struct{})
	retChan := make(chan InfoResult, len(s.urls))
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(s.ctx, "webFingerScan", pr)
		}
		close(single)
	}()
	// 指纹扫描
	fscan := func(u *url.URL) {
		if ExitFunc {
			return
		}
		resp, body, err := clients.NewSimpleGetRequest(u.String(), s.client)
		if err != nil || resp == nil {
			retChan <- InfoResult{
				URL:        u.String(),
				StatusCode: 0,
			}
			return
		}
		title, server, content_type := s.GetHeaderInfo(body, resp)
		headers, _, _ := DumpResponseHeadersAndRaw(resp)
		hashValue, md5Value := FaviconHash(u, s.client)
		web := &WebInfo{
			HeadeString:   string(headers),
			ContentType:   content_type,
			Cert:          GetTLSString(u.Scheme, u.Host),
			BodyString:    string(body),
			Path:          u.Path,
			Title:         title,
			Server:        server,
			ContentLength: len(body),
			Port:          netutil.GetPort(u),
			IconHash:      hashValue,
			IconMd5:       md5Value,
			StatusCode:    resp.StatusCode,
			Banner:        s.GetBanner(u),
		}

		wafInfo := *waf.IsWAF(u.Hostname(), subdomain.DefaultDnsServers)

		s.aliveURLs = append(s.aliveURLs, u)

		fingerprints := s.FingerScan(s.ctx, web, FingerprintDB)

		s.mutex.Lock()
		s.basicURLWithFingerprint[u.String()] = append(s.basicURLWithFingerprint[u.String()], fingerprints...)
		s.mutex.Unlock()

		screenshotPath := ""

		if s.screenshot {
			if screenshotPath, err = GetScreenshot(u.String()); err != nil {
				gologger.Debug(s.ctx, err)
			}
		}
		retChan <- InfoResult{
			URL:          u.String(),
			StatusCode:   web.StatusCode,
			Length:       web.ContentLength,
			Title:        web.Title,
			Fingerprints: fingerprints,
			IsWAF:        wafInfo.Exsits,
			WAF:          wafInfo.Name,
			Detect:       "Default",
			Screenshot:   screenshotPath,
		}
	}
	threadPool, _ := ants.NewPoolWithFunc(50, func(target interface{}) {
		t := target.(*url.URL)
		fscan(t)
		wg.Done()
	})
	defer threadPool.Release()
	for _, target := range s.urls {
		if ExitFunc {
			return
		}
		wg.Add(1)
		threadPool.Invoke(target)
	}
	wg.Wait()
	close(retChan)
	gologger.Info(s.ctx, "FingerScan Finished")
	<-single
}

type TFP struct {
	URL  *url.URL
	Fpe  []FingerPEntity
	Path string
}

func (s *FingerScanner) NewActiveFingerScan(rootPath bool) {
	if len(s.aliveURLs) == 0 {
		gologger.Warning(s.ctx, "未发现存活目标，已跳过主动指纹扫描")
		return
	}
	gologger.Info(s.ctx, "正在进行主动指纹探测 ...")
	var id = 0
	var wg sync.WaitGroup
	visited := make(map[string]bool) // 用于记录已访问的URL和路径组合

	single := make(chan struct{})
	retChan := make(chan InfoResult, len(s.urls))
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(s.ctx, "webFingerScan", pr)
		}
		close(single)
	}()
	// 主动指纹扫描
	threadPool, _ := ants.NewPoolWithFunc(5, func(tfp interface{}) {
		defer wg.Done()
		fp := tfp.(TFP)
		fullURL := fp.URL.String() + fp.Path

		// 确保并发唯一性：检查和标记已经访问过的URL和路径组合
		s.mutex.Lock()
		if visited[fullURL] {
			s.mutex.Unlock()
			return // 已经处理过，跳过
		}
		visited[fullURL] = true
		s.mutex.Unlock()

		resp, body, _ := clients.NewSimpleGetRequest(fullURL, s.client)
		if resp == nil {
			return
		}
		title, server, content_type := s.GetHeaderInfo(body, resp)
		headers, _, _ := DumpResponseHeadersAndRaw(resp)
		ti := &WebInfo{
			HeadeString:   string(headers),
			ContentType:   content_type,
			BodyString:    string(body),
			Path:          fp.Path,
			Title:         title,
			Server:        server,
			ContentLength: len(body),
			Port:          netutil.GetPort(fp.URL),
			StatusCode:    resp.StatusCode,
		}
		result := s.FingerScan(s.ctx, ti, fp.Fpe)
		if len(result) > 0 {
			s.mutex.Lock()
			s.basicURLWithFingerprint[fp.URL.String()] = append(s.basicURLWithFingerprint[fp.URL.String()], result...)
			s.mutex.Unlock()

			retChan <- InfoResult{
				URL:          fullURL,
				StatusCode:   ti.StatusCode,
				Length:       ti.ContentLength,
				Title:        ti.Title,
				Fingerprints: []string{fp.Fpe[0].ProductName},
				Detect:       "Active",
			}
		}
	})
	defer threadPool.Release()
	s.ActiveCounts()
	// 载入存活目标
	for _, target := range s.aliveURLs {
		for _, afdb := range ActiveFingerprintDB {
			for _, path := range afdb.Path {
				if ExitFunc { // 控制程序退出
					return
				}
				wg.Add(1)
				id += 1
				runtime.EventsEmit(s.ctx, "ActiveProgressID", id)
				if rootPath {
					target, _ = url.Parse(util.GetBasicURL(target.String()))
				}
				threadPool.Invoke(TFP{
					URL:  target,
					Fpe:  afdb.Fpe,
					Path: path,
				})
			}
		}
	}
	wg.Wait()
	close(retChan)
	gologger.Info(s.ctx, "ActiveFingerScan Finished")
	runtime.EventsEmit(s.ctx, "ActiveScanComplete", 100)
	<-single
}

// 统计主动指纹总共要扫描的目标
func (s *FingerScanner) ActiveCounts() {
	var id = 0
	for _, afdb := range ActiveFingerprintDB {
		id += len(afdb.Path)
	}
	count := len(s.aliveURLs) * id
	runtime.EventsEmit(s.ctx, "ActiveCounts", count)
}

func (s *FingerScanner) URLWithFingerprintMap() map[string][]string {
	return s.basicURLWithFingerprint
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

func (s *FingerScanner) FingerScan(ctx context.Context, web *WebInfo, targetDB []FingerPEntity) []string {
	var fingerPrintResults []string

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
				expr := finger.AllString
				for _, rule := range finger.Rule {
					var result bool
					switch rule.Key {
					case "header":
						result = dataCheckString(rule.Op, web.HeadeString, rule.Value)
					case "body":
						result = dataCheckString(rule.Op, web.BodyString, rule.Value)
					case "server":
						result = dataCheckString(rule.Op, web.Server, rule.Value)
					case "title":
						result = dataCheckString(rule.Op, web.Title, rule.Value)
					case "cert":
						result = dataCheckString(rule.Op, web.Cert, rule.Value)
					case "port":
						value, err := strconv.Atoi(rule.Value)
						if err == nil {
							result = dataCheckInt(rule.Op, web.Port, value)
						}
					case "protocol":
						result = (rule.Op == 0 && web.Protocol == rule.Value) || (rule.Op == 1 && web.Protocol != rule.Value)
					case "path":
						result = dataCheckString(rule.Op, web.Path, rule.Value)
					case "icon_hash":
						value, err := strconv.Atoi(rule.Value)
						hashIcon, errHash := strconv.Atoi(web.IconHash)
						if err == nil && errHash == nil {
							result = dataCheckInt(rule.Op, hashIcon, value)
						}
					case "icon_mdhash":
						result = dataCheckString(rule.Op, web.IconMd5, rule.Value)
					case "status":
						value, err := strconv.Atoi(rule.Value)
						if err == nil {
							result = dataCheckInt(rule.Op, web.StatusCode, value)
						}
					case "content_type":
						result = dataCheckString(rule.Op, web.ContentType, rule.Value)
					case "banner":
						result = dataCheckString(rule.Op, web.Banner, rule.Value)
					}

					if result {
						expr = expr[:rule.Start] + "T" + expr[rule.End:]
					} else {
						expr = expr[:rule.Start] + "F" + expr[rule.End:]
					}
				}
				r := boolEval(ctx, expr)
				if r {
					results <- finger.ProductName
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

// 获取favicon Mmh3Hash32 和 MD5值
func FaviconHash(u *url.URL, client *http.Client) (string, string) {
	_, body, err := clients.NewSimpleGetRequest(u.String(), client)
	if err != nil {
		return "", ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", ""
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
		hasher := md5.New()
		hasher.Write(body)
		sum := hasher.Sum(nil)
		return util.Mmh3Hash32(util.Base64Encode(body)), hex.EncodeToString(sum)
	}
	return "", ""
}

func (s *FingerScanner) GetBanner(u *url.URL) string {
	if strings.HasPrefix(u.Scheme, "http") {
		return ""
	}
	scanner := gonmap.New()
	_, response := scanner.Scan(u.Host, netutil.GetPort(u), time.Second*time.Duration(10))
	if response != nil {
		return response.Raw
	}
	return ""
}

func (s *FingerScanner) GetHeaderInfo(body []byte, resp *http.Response) (title, server, content_type string) {
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
