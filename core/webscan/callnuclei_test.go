package webscan

import (
	"context"
	"fmt"
	"net/url"
	"slack-wails/core/subdomain"
	"slack-wails/core/waf"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/netutil"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	syncutil "github.com/projectdiscovery/utils/sync"
)

func TestNucleiCaller(t *testing.T) {
	// proxys := []string{"http://127.0.0.1:8080"}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(),
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{util.HomeDir() + "/slack/config/pocs"},
		}), // -t
		nuclei.DisableUpdateCheck(), // -duc
		// nuclei.WithProxy(proxys, false), // -proxy
	)
	if err != nil {
		panic(err)
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{}, false)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		fmt.Printf("[%s] [%s] %s\n", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched)
		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
			fmt.Printf("Reference: %s\n", event.Info.Reference.ToSlice())
		}
		fmt.Printf("ExtractedResults: %s\n", strings.Join(event.ExtractedResults, ","))
	})
	if err != nil {
		panic(err)
	}
	defer ne.Close()
}

func TestThreadSafeNucleiCaller(t *testing.T) {
	ctx := context.Background()
	// when running nuclei in parallel for first time it is a good practice to make sure
	// templates exists first

	// create nuclei engine with options
	ne, err := nuclei.NewThreadSafeNucleiEngineCtx(ctx)
	if err != nil {
		panic(err)
	}
	// setup sizedWaitgroup to handle concurrency
	sg, err := syncutil.New(syncutil.WithSize(10))
	if err != nil {
		panic(err)
	}

	// scan 1 = run dns templates on scanme.sh
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{"scanme.sh"},
			nuclei.WithTemplateFilters(nuclei.TemplateFilters{ProtocolTypes: "dns"}),
			nuclei.WithHeaders([]string{"X-Bug-Bounty: pdteam"}),
			nuclei.EnablePassiveMode(),
		)
		if err != nil {
			panic(err)
		}
	}()

	// scan 2 = run templates with oast tags on honey.scanme.sh
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{"https://202.88.229.90/"}, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{Templates: []string{"/Users/qwtd/slack/config/pocs/dss-download-fileread.yaml"}}))
		if err != nil {
			panic(err)
		}
	}()

	// wait for all scans to finish
	sg.Wait()
	defer ne.Close()

	// Output:
	// [dns-saas-service-detection] scanme.sh
	// [nameserver-fingerprint] scanme.sh
	// [dns-saas-service-detection] honey.scanme.sh
}

func TestFingerscan(t *testing.T) {
	s := NewFingerScanEngine(context.Background(), "", clients.Proxy{}, structs.WebscanOptions{
		Target:   []string{""},
		Thread:   10,
		RootPath: true,
	})
	if s == nil {
		t.Log("engine is nil")
		return
	}

	var wg sync.WaitGroup
	single := make(chan struct{})
	retChan := make(chan structs.InfoResult, len(s.urls))
	go func() {
		for pr := range retChan {
			fmt.Printf("pr: %v\n", pr)
		}
		close(single)
	}()
	// 指纹扫描
	fscan := func(u *url.URL) {
		var (
			rawHeaders  []byte
			server      string
			contentType string
			statusCode  int
		)
		fmt.Printf("u.String(): %v\n", u.String())
		// 先进行一次不会重定向的扫描，可以获得重定向前页面的响应头中获取指纹
		resp, err := clients.DoRequest("GET", u.String(), s.headers, nil, 10, s.notFollowClient)
		if err == nil && resp.StatusCode() == 302 {
			rawHeaders = DumpResponseHeadersOnly(resp.RawResponse)
		}

		// 过滤CDN
		if resp.StatusCode() == 422 {
			retChan <- structs.InfoResult{
				URL:        u.String(),
				StatusCode: 422,
				Scheme:     u.Scheme,
			}
			return
		}

		// 正常请求指纹
		resp, err = clients.DoRequest("GET", u.String(), s.headers, nil, 10, s.client)
		body := resp.Body()
		if err != nil {
			if len(rawHeaders) > 0 {
				t.Logf("%s has error to 302, response headers: %s", u.String(), string(rawHeaders))
				gologger.Debug(s.ctx, fmt.Sprintf("%s has error to 302, response headers: %s", u.String(), string(rawHeaders)))
				statusCode = 302
			} else {
				retChan <- structs.InfoResult{
					URL:        u.String(),
					StatusCode: 0,
					Scheme:     u.Scheme,
				}
				return
			}
		}
		// 合并请求头数据
		rawHeaders = append(rawHeaders, DumpResponseHeadersOnly(resp.RawResponse)...)

		// 请求Logo
		faviconHash, faviconMd5 := FaviconHash(u, s.headers, s.client)

		// 发送shiro探测
		rawHeaders = append(rawHeaders, fmt.Appendf(nil, "Set-Cookie: %s", s.ShiroScan(u))...)

		// 跟随JS重定向，并替换成重定向后的数据
		redirectBody := s.GetJSRedirectResponse(u, string(body))
		if redirectBody != nil {
			// JS重定向后，body数据不应该直接覆盖 fix in 2.0.8
			body = append(body, redirectBody...)
			// body = redirectBody
		}
		// 网站正常响应
		title := clients.GetTitle(body)
		server = resp.Header().Get("Server")
		contentType = resp.Header().Get("Content-Type")
		statusCode = resp.StatusCode()
		web := &WebInfo{
			HeadeString:   string(rawHeaders),
			ContentType:   contentType,
			Cert:          GetTLSString(u.Scheme, u.Host),
			BodyString:    string(body),
			Path:          u.Path,
			Title:         title,
			Server:        server,
			ContentLength: len(body),
			Port:          netutil.GetPort(u),
			IconHash:      faviconHash,
			IconMd5:       faviconMd5,
			StatusCode:    statusCode,
		}

		wafInfo := *waf.ResolveAndWafIdentify(u.Hostname(), subdomain.DefaultDnsServers)

		s.aliveURLs = append(s.aliveURLs, u)

		fingerprints := s.Scan(s.ctx, web, FingerprintDB)

		if s.generateLog4j2 {
			fingerprints = append(fingerprints, "Generate-Log4j2")
		}

		if s.FastjsonScan(u) {
			fingerprints = append(fingerprints, "Fastjson")
		}

		if checkHoneypotWithHeaders(web.HeadeString) || checkHoneypotWithFingerprintLength(len(fingerprints)) {
			fingerprints = []string{"疑似蜜罐"}
		}

		// 截屏
		var screenshotPath string
		// 截屏条件要满足协议, fix in v2.0.8
		if s.screenshot && (u.Scheme == "https" || u.Scheme == "http") {
			if screenshotPath, err = GetScreenshot(u.String()); err != nil {
				t.Log(err)
			}
		}

		s.mutex.Lock()
		s.basicURLWithFingerprint[u.String()] = append(s.basicURLWithFingerprint[u.String()], fingerprints...)
		s.mutex.Unlock()

		retChan <- structs.InfoResult{
			TaskId:       s.taskId,
			URL:          u.String(),
			Scheme:       u.Scheme,
			Host:         u.Host,
			Port:         web.Port,
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
	threadPool, _ := ants.NewPoolWithFunc(s.thread, func(target interface{}) {
		t := target.(*url.URL)
		fscan(t)
		wg.Done()
	})
	defer threadPool.Release()
	for _, target := range s.urls {
		wg.Add(1)
		threadPool.Invoke(target)
	}
	wg.Wait()
	close(retChan)
	t.Log("FingerScan Finished")
	<-single
}
