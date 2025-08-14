package jsfind

import (
	"context"
	"fmt"
	"regexp"
	"slack-wails/lib/gologger"
	"slack-wails/lib/utils/arrayutil"
	"strings"
	"time"

	"github.com/qiwentaidi/clients"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	regJS  []*regexp.Regexp
	JsLink = []string{
		"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
	}
)

func init() {
	for _, reg := range JsLink {
		regJS = append(regJS, regexp.MustCompile(reg))
	}
}

func ExtractAllJs(ctx context.Context, url string) (allJS []string) {
	staticJsLinks := extractStaticJs(ctx, url)
	dynamicJsLinks := extractDynamicsJs(ctx, url)
	// 如果动态加载到JS不为空，则获取到的完整URL一定是正确路径，静态提取的JS内容如果在动态JS中出现过则删除，可以避免后续二次访问
	allJS = append(allJS, dynamicJsLinks...)
	if len(dynamicJsLinks) > 0 {
		for _, static := range staticJsLinks {
			for _, dynamic := range dynamicJsLinks {
				if !strings.Contains(dynamic, static) {
					allJS = append(allJS, static)
				}
			}
		}
	} else {
		allJS = append(allJS, staticJsLinks...)
	}
	return arrayutil.RemoveDuplicates(allJS)
}

// 通过主页提取静态JS
func extractStaticJs(ctx context.Context, url string) []string {
	var staticJsLinks []string
	resp, err := clients.SimpleGet(url, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, err)
		return staticJsLinks
	}
	content := string(resp.Body())
	for _, reg := range regJS {
		for _, item := range reg.FindAllString(content, -1) {
			item = strings.TrimSpace(item)
			item = strings.TrimLeft(item, "=")
			item = strings.Trim(item, "\"")
			item = strings.Trim(item, "'")
			item = strings.TrimLeft(item, ".")
			staticJsLinks = append(staticJsLinks, item)
		}
	}
	return staticJsLinks
}

// 通过chromedp提取动态JS, 用于对静态的查漏补缺
func extractDynamicsJs(mainCtx context.Context, url string) []string {
	var dynamicJsLinks []string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	linkSet := make(map[string]bool) // 防止重复

	// 监听所有网络响应
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			url := ev.Response.URL
			if strings.Contains(url, ".js") {
				if !linkSet[url] {
					dynamicJsLinks = append(dynamicJsLinks, url)
					linkSet[url] = true
				}
			}
		}
	})

	// 启动网络监听
	err := chromedp.Run(ctx,
		network.Enable(),              // 启用网络监听
		chromedp.Navigate(url),        // 替换为目标 URL
		chromedp.Sleep(5*time.Second), // 等待页面加载 & JS 动态加载完成
	)
	if err != nil {
		runtime.EventsEmit(mainCtx, "jsfindlog", fmt.Sprintf("[-] chromedp 运行失败无法动态加载JS链接, 错误原因: (%v)", err))
		return dynamicJsLinks
	}

	return dynamicJsLinks
}
