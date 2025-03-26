package jsfind

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	regJS  []*regexp.Regexp
	JsLink = []string{
		"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
	}
	Filter = []string{".vue", ".jpeg", ".png", ".jpg", ".gif", ".css", ".svg", ".scss", ".eot", ".ttf", ".woff", ".js", ".ts", ".tsx", ".ico", ".less"}
	// 因为在提取API时经常会有无用的类似数据，所以需要过滤
	apiFilter = []string{
		"text/xml",
		"text/html",
		"text/plain",
		"text/css",
		"text/javascript",
		"multipart/form-data",
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"audio/",
		"video/",
		"application/",
		"YYYY",
		"/a/i",
		"/a/b",
		"./",
		"0.",
		"1.",
		"2.",
		"3.",
		"4.",
		"5.",
		"6.",
		"7.",
		"8.",
		"9.",
		"404",
		"403",
		"500",
		"M/D/yy",
		"text/csv",
	}
	Sensitive = regexp.MustCompile(`\b(access.{0,3}key|access.{0,3}Key|access.{0,3}Id|access.{0,3}id|access.{0,3}Secret|access.{0,3}secret|bucket|Bucket|endpoint|Endpoint|.{0,5}密码|.{0,5}账号|默认.{0,5}密码|password|username)\s*[:=]\s*["']([^"']+)["']`)
	Phone     = regexp.MustCompile(`(^|[^0-9a-zA-Z])(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}([^0-9a-zA-Z]|$)`)
	IDCard    = regexp.MustCompile(`(^|[^0-9a-zA-Z])((\d{8}(0\d|10|11|12)([0-2]\d|30|31)\d{3}$)|(\d{6}(18|19|20)\d{2}(0[1-9]|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)))([^0-9a-zA-Z]|$)`)
	Email     = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	Link      = regexp.MustCompile(`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:\w)(?:[\?|#][^"|']{0,}|)))(?:"|')`)
	IP_PORT   = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:\d{1,5}`)
)

func init() {
	for _, reg := range JsLink {
		regJS = append(regJS, regexp.MustCompile(reg))
	}
}

func ExtractJS(ctx context.Context, url string) (allJS []string) {
	_, body, err := clients.NewSimpleGetRequest(url, clients.NewHttpClient(nil, true))
	if err != nil || body == nil {
		gologger.Debug(ctx, err)
		return
	}
	content := string(body)
	for _, reg := range regJS {
		for _, item := range reg.FindAllString(content, -1) {
			item = strings.TrimLeft(item, "=")
			item = strings.Trim(item, "\"")
			item = strings.Trim(item, "'")
			item = strings.TrimLeft(item, ".")
			allJS = append(allJS, item)
		}
	}
	return util.RemoveDuplicates(allJS)
}

// FindInfo 使用线程池处理单个 URL 的信息提取
func FindInfo(ctx context.Context, url string) *structs.FindSomething {
	var fs = &structs.FindSomething{}
	_, body, err := clients.NewSimpleGetRequest(url, clients.NewHttpClient(nil, true))
	if err != nil || body == nil {
		gologger.Debug(ctx, err)
		return fs
	}

	content := string(body)
	// 提取信息
	urls, apis := urlInfoSeparate(Link.FindAllString(content, -1))
	// fs.JS = *AppendSource(url, js)
	fs.APIRoute = *AppendSource(url, apis)
	fs.IP_URL = *AppendSource(url, append(IP_PORT.FindAllString(content, -1), urls...))
	fs.IDCard = *AppendSource(url, clean(IDCard.FindAllString(content, -1)))
	fs.Phone = *AppendSource(url, clean(Phone.FindAllString(content, -1)))
	fs.Sensitive = *AppendSource(url, Sensitive.FindAllString(content, -1))
	fs.Email = *AppendSource(url, Email.FindAllString(content, -1))
	return fs
}

// MultiThreadJSFind 使用 ants 线程池进行多线程处理
func MultiThreadJSFind(ctx context.Context, target, prefixJsURL string, jsLinks []string) structs.FindSomething {
	var fs = structs.FindSomething{}
	var mu sync.Mutex // 用于保护 fs 的并发写操作
	var wg sync.WaitGroup
	if prefixJsURL != "" {
		target = prefixJsURL
	}
	// 创建线程池
	pool, _ := ants.NewPoolWithFunc(10, func(data interface{}) {
		defer wg.Done()
		jslink := data.(string)
		var newURL string
		host, _ := util.GetBasePath(target)

		if strings.HasPrefix(jslink, "http") {
			newURL = jslink
		} else {
			newURL = strings.TrimRight(host, "/") + "/" + strings.TrimLeft(jslink, "/")
		}

		fs2 := FindInfo(ctx, newURL)

		// 加锁以安全地更新共享资源
		mu.Lock()
		fs.IP_URL = append(fs.IP_URL, fs2.IP_URL...)
		fs.IDCard = append(fs.IDCard, fs2.IDCard...)
		fs.Phone = append(fs.Phone, fs2.Phone...)
		fs.Sensitive = append(fs.Sensitive, fs2.Sensitive...)
		fs.APIRoute = append(fs.APIRoute, fs2.APIRoute...)
		mu.Unlock()
	})
	defer pool.Release()
	wg.Add(1) // 确保 target 任务计入 WaitGroup

	// 提交主目标任务
	pool.Invoke(target)

	// 提交 JS 链接任务
	for _, jslink := range jsLinks {
		wg.Add(1)
		pool.Invoke(jslink)
	}

	// 等待所有任务完成
	wg.Wait()

	// 去重处理
	fs.APIRoute = RemoveDuplicatesInfoSource(fs.APIRoute)
	fs.IDCard = RemoveDuplicatesInfoSource(fs.IDCard)
	fs.Phone = RemoveDuplicatesInfoSource(fs.Phone)
	fs.Sensitive = RemoveDuplicatesInfoSource(fs.Sensitive)
	fs.IP_URL = FilterExt(RemoveDuplicatesInfoSource(fs.IP_URL))
	return fs
}

// clean 去除手机号和身份证中的其他字符，保留数字、+ 和 X
func clean(filed []string) (news []string) {
	// 定义只保留数字、+ 和 X 的正则表达式
	cleanRegex := regexp.MustCompile(`[^\d+X]`)
	for _, p := range filed {
		news = append(news, cleanRegex.ReplaceAllString(p, ""))
	}
	return
}

func AppendSource(source string, filed []string) *[]structs.InfoSource {
	is := []structs.InfoSource{}
	for _, f := range filed {
		is = append(is, structs.InfoSource{Filed: f, Source: source})
	}
	return &is
}

func RemoveDuplicatesInfoSource(iss []structs.InfoSource) []structs.InfoSource {
	encountered := map[string]bool{}
	result := []structs.InfoSource{}
	for _, is := range iss {
		if !encountered[is.Filed] {
			encountered[is.Filed] = true
			result = append(result, is)
		}
	}
	return result
}

func urlInfoSeparate(links []string) (urls, apis []string) {
	for _, link := range links {
		link = strings.Trim(link, "\"")
		link = strings.Trim(link, "'")

		// 判断是否是 URL
		if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "ws") || strings.HasPrefix(link, "//") {
			urls = append(urls, link)
			continue // 直接进入下一次循环
		}

		matched := false
		for _, r := range Filter {
			if strings.Contains(link, r) {
				matched = true // 过滤匹配成功，跳过后续逻辑
				break
			}
		}
		if matched {
			continue // 直接跳过，不执行后续 API 判断
		}

		// 如果 link 包含 apiFilter 中的任何字符串，则跳过，不添加到 apis
		isFiltered := false
		filters := append(apiFilter, Filter...)
		for _, filter := range filters {
			if strings.Contains(link, filter) {
				isFiltered = true
				break
			}
		}
		if !isFiltered {
			apis = append(apis, link)
		}
	}

	return urls, apis
}
func FilterExt(iss []structs.InfoSource) (news []structs.InfoSource) {
	for _, link := range iss {
		matched := false
		for _, r := range Filter {
			if strings.Contains(link.Filed, r) {
				matched = true
				break
			}
		}
		if !matched {
			news = append(news, link)
		}

	}
	return news
}

// 处理 API 逻辑
func AnalyzeAPI(ctx context.Context, o structs.JSFindOptions) {
	resp, body, err := clients.NewSimpleGetRequest(o.HomeURL, clients.NewHttpClient(nil, true))
	if err != nil || resp == nil {
		gologger.Error(ctx, fmt.Sprintf("[AnalyzeAPI] %s, error: %v", o.HomeURL, err))
		return
	}
	homeBody := string(body)

	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(10, func(data interface{}) {
		defer wg.Done()
		api := data.(string)

		fullURL := strings.TrimRight(o.BaseURL, "/") + "/" + strings.TrimLeft(api, "/")
		method, err := detectMethod(fullURL, o.Headers)
		if err != nil {
			runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[!] %s : %v", fullURL, err))
			return
		}

		param := completeParameters(method, fullURL, url.Values{})
		if param != nil {
			runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[+] %s 已补全参数 %s", fullURL, param.Encode()))
		}

		apiReq := APIRequest{
			URL:     fullURL,
			Method:  method,
			Headers: o.Headers,
			Params:  param,
		}
		for _, router := range o.HighRiskRouter {
			if strings.Contains(apiReq.URL, router) {
				runtime.EventsEmit(ctx, "jsfindlog", "[!!] "+fullURL+" 高风险API已跳过测试, 相关敏感词: "+router)
				return
			}
		}
		// 测试未授权访问
		vulnerable, body, err := testUnauthorizedAccess(homeBody, apiReq, o.Authentication)
		if err != nil || !vulnerable {
			runtime.EventsEmit(ctx, "jsfindlog", "[-] "+fullURL+" 不存在未授权")
			return
		}
		runtime.EventsEmit(ctx, "jsfindlog", "[+] "+fullURL+" 存在未授权")
		runtime.EventsEmit(ctx, "jsfindvulcheck", structs.JSFindResult{
			Method:   method,
			Param:    param.Encode(),
			Source:   fullURL,
			Response: string(body),
			Length:   len(body),
		})
	})
	defer pool.Release()

	// 提交 API 任务到线程池
	for _, api := range o.ApiList {
		wg.Add(1)
		pool.Invoke(api)
	}

	wg.Wait()
}
