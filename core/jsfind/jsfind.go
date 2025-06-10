package jsfind

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils/httputil"
	"strings"
	"sync"

	"maps"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const maxResponseSize = 500 * 1024 // 500KB

var (
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
	Sensitive = regexp.MustCompile(`\b(access.{0,3}key|access.{0,3}Key|access.{0,3}Id|access.{0,3}id|access.{0,3}Secret|access.{0,3}secret|bucket|Bucket|endpoint|Endpoint|.{0,5}密码|.{0,5}账号|默认.{0,5}密码|password|username)\s*[:=]\s*["']?([^"'\s]+)["']?|\b((ey[A-Za-z0-9_-]{10,}\.[A-Za-z0-9._-]{10,}|ey[A-Za-z0-9_\/+-]{10,}\.[A-Za-z0-9._\/+-]{10,}))\b`)
	Phone     = regexp.MustCompile(`(^|[^0-9a-zA-Z])(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}([^0-9a-zA-Z]|$)`)
	IDCard    = regexp.MustCompile(`(^|[^0-9a-zA-Z])((\d{8}(0\d|10|11|12)([0-2]\d|30|31)\d{3}$)|(\d{6}(18|19|20)\d{2}(0[1-9]|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)))([^0-9a-zA-Z]|$)`)
	Email     = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	Link      = regexp.MustCompile(`(?:"|')(((?:[a-zA-Z]{1,10}://|//)[^"'/]{1,}\.[a-zA-Z]{2,}[^"']{0,})|((?:/|\.\./|\./)[^"'><,;|*()(%%$^/\\\[\]][^"'><,;|()]{1,})|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{1,}\.(?:[a-zA-Z]{1,4}|action)(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-/]{1,}/[a-zA-Z0-9_\-/]{3,}(?:[\?|#][^"|']{0,}|))|([a-zA-Z0-9_\-]{1,}\.(?:\w)(?:[\?|#][^"|']{0,}|)))(?:"|')`)
	IP_PORT   = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:\d{1,5}`)
)

// FindInfo 使用线程池处理单个 URL 的信息提取
func FindInfo(ctx context.Context, url string) *structs.FindSomething {
	var fs = &structs.FindSomething{}
	resp, err := clients.SimpleGet(url, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Debug(ctx, err)
		return fs
	}
	content := string(resp.Body())
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

// 对还原出来的 SourceMap 目录提取信息
func ExtractFromSourceMapDir(ctx context.Context, dirPath string) *structs.FindSomething {
	var fsResult = &structs.FindSomething{}

	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // 忽略错误继续
		}
		if d.IsDir() {
			return nil
		}

		// 只处理 .js .vue .ts .tsx 文件
		if !(strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".vue") || strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx")) {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(contentBytes)

		// 提取
		urls, apis := urlInfoSeparate(Link.FindAllString(content, -1))

		fsResult.APIRoute = append(fsResult.APIRoute, *AppendSource(path, apis)...)
		fsResult.IP_URL = append(fsResult.IP_URL, *AppendSource(path, append(IP_PORT.FindAllString(content, -1), urls...))...)
		fsResult.IDCard = append(fsResult.IDCard, *AppendSource(path, clean(IDCard.FindAllString(content, -1)))...)
		fsResult.Phone = append(fsResult.Phone, *AppendSource(path, clean(Phone.FindAllString(content, -1)))...)
		fsResult.Sensitive = append(fsResult.Sensitive, *AppendSource(path, Sensitive.FindAllString(content, -1))...)
		fsResult.Email = append(fsResult.Email, *AppendSource(path, Email.FindAllString(content, -1))...)

		return nil
	})

	// 去重
	fsResult.APIRoute = RemoveDuplicatesInfoSource(fsResult.APIRoute)
	fsResult.IP_URL = FilterExt(RemoveDuplicatesInfoSource(fsResult.IP_URL))
	fsResult.IDCard = RemoveDuplicatesInfoSource(fsResult.IDCard)
	fsResult.Phone = RemoveDuplicatesInfoSource(fsResult.Phone)
	fsResult.Sensitive = RemoveDuplicatesInfoSource(fsResult.Sensitive)

	return fsResult
}

func Scan(ctx context.Context, target, prefixJsURL string, jsLinks, blackDomainList []string) structs.FindSomething {
	var fs = structs.FindSomething{}
	var mu sync.Mutex // 用于保护 fs 的并发写操作
	var wg sync.WaitGroup
	var visited sync.Map // 并发安全的去重 map
	if prefixJsURL != "" {
		target = prefixJsURL
	}

	// 创建线程池
	pool, _ := ants.NewPoolWithFunc(10, func(data interface{}) {
		defer wg.Done()
		jslink := data.(string)
		newURL := formatURL(target, jslink)
		for _, blackDomain := range blackDomainList {
			if strings.Contains(newURL, blackDomain) {
				// 黑名单域名，跳过
				return
			}
		}
		// 判断是否已访问
		if _, loaded := visited.LoadOrStore(newURL, struct{}{}); loaded {
			// 已处理过，跳过
			return
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

	for _, jslink := range jsLinks {
		mapURL := formatURL(target, jslink) + ".map"
		resp, err := clients.SimpleGet(mapURL, clients.NewRestyClient(nil, true))
		if err == nil && resp.StatusCode() == 200 {
			// 检测到 .map 泄漏，尝试还原
			fp, err := RestoreWebpack(ctx, mapURL)
			if err == nil {
				runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[+] 发现JS SourceMap泄漏: %s, 恢复webpack成功: %s", mapURL, fp))
			}
			sourceMapInfo := ExtractFromSourceMapDir(ctx, fp)
			mu.Lock()
			fs.IP_URL = append(fs.IP_URL, sourceMapInfo.IP_URL...)
			fs.IDCard = append(fs.IDCard, sourceMapInfo.IDCard...)
			fs.Phone = append(fs.Phone, sourceMapInfo.Phone...)
			fs.Sensitive = append(fs.Sensitive, sourceMapInfo.Sensitive...)
			fs.APIRoute = append(fs.APIRoute, sourceMapInfo.APIRoute...)
			mu.Unlock()
		}
	}

	// 去重处理
	fs.APIRoute = RemoveDuplicatesInfoSource(fs.APIRoute)
	fs.IDCard = RemoveDuplicatesInfoSource(fs.IDCard)
	fs.Phone = RemoveDuplicatesInfoSource(fs.Phone)
	fs.Sensitive = RemoveDuplicatesInfoSource(fs.Sensitive)
	fs.IP_URL = FilterExt(RemoveDuplicatesInfoSource(fs.IP_URL))
	return fs
}

func formatURL(url string, jslink string) string {
	var newURL string
	host, _ := httputil.GetBasePath(url)
	if strings.HasPrefix(jslink, "http") {
		newURL = jslink
	} else {
		newURL = host + strings.TrimLeft(jslink, "/")
	}
	return newURL
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
	resp, err := clients.SimpleGet(o.HomeURL, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("[AnalyzeAPI] 请求首页失败 %s, 错误: %v", o.HomeURL, err))
		return
	}
	homeBody := string(resp.Body())

	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(10, func(data interface{}) {
		defer wg.Done()
		api := data.(string)

		// 为每个 API 独立生成 fullURL 和 headers 副本
		fullURL := strings.TrimRight(o.BaseURL, "/") + "/" + strings.TrimLeft(api, "/")

		// 拷贝 headers，避免竞争
		apiHeaders := make(map[string]string)
		maps.Copy(apiHeaders, o.Headers)

		// 检测请求方法
		method, err := detectMethod(fullURL, apiHeaders)
		if err != nil {
			runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[!] %s: %v", fullURL, err))
			return
		}
		body := ""
		// 如果是 POST 方法，动态探测 Content-Type
		if method == http.MethodPost {
			if contentType := detectContentType(fullURL, apiHeaders); contentType != "" {
				apiHeaders["Content-Type"] = contentType
				if contentType == "application/json" {
					body = "{}"
				}
			}
		}

		// 补全参数
		param := completeParameters(ctx, method, fullURL, url.Values{})
		if param.Encode() != "" {
			runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[+] %s 已补全参数: %s", fullURL, param.Encode()))
		}

		// 构建请求对象
		apiReq := APIRequest{
			URL:     fullURL,
			Method:  method,
			Headers: apiHeaders,
			Params:  param,
			Body:    body,
		}

		// 检查高风险路由，直接跳过测试
		for _, router := range o.HighRiskRouter {
			if strings.Contains(strings.ToLower(apiReq.URL), router) {
				runtime.EventsEmit(ctx, "jsfindlog", "[!!] "+fullURL+" 高风险API跳过测试, 触发敏感词: "+router)
				return
			}
		}

		// 测试未授权访问
		vulnerable, body, err := testUnauthorizedAccess(homeBody, apiReq, o.Authentication)
		if err != nil {
			runtime.EventsEmit(ctx, "jsfindlog", fmt.Sprintf("[-] %s 测试未授权错误: %v", fullURL, err))
			return
		}

		if !vulnerable {
			runtime.EventsEmit(ctx, "jsfindlog", "[-] "+fullURL+" 不存在未授权访问")
			return
		}

		// 存在未授权，记录漏洞信息
		runtime.EventsEmit(ctx, "jsfindlog", "[+] "+fullURL+" 存在未授权访问！")
		runtime.EventsEmit(ctx, "jsfindvulcheck", structs.JSFindResult{
			VulType:  "未授权访问",
			Method:   method,
			Request:  buildRawRequest(apiReq),
			Source:   fullURL,
			Response: httputil.LimitResponse(body, maxResponseSize, "Response too large, please manually open the link to view."),
			Length:   len(body),
		})
		// 检测越权
		if len(o.LowPrivilegeHeaders) != 0 && len(o.Headers) != 0 {
			lowPrivilegeReq := apiReq
			lowPrivilegeReq.Headers = o.LowPrivilegeHeaders
			isvulnerable, lowPrivBody, err := testPrivilegeEscalation(body, lowPrivilegeReq)
			if err != nil {
				runtime.EventsEmit(ctx, "jsfindlog", "[!] "+fullURL+" 检测越权访问失败："+err.Error())
				return
			}
			if !isvulnerable {
				runtime.EventsEmit(ctx, "jsfindlog", "[-] "+fullURL+" 不存在未授权访问")
				return
			}
			runtime.EventsEmit(ctx, "jsfindlog", "[+] "+fullURL+" 检测到越权访问")
			runtime.EventsEmit(ctx, "jsfindvulcheck", structs.JSFindResult{
				VulType:  "越权访问",
				Method:   method,
				Request:  buildRawRequest(lowPrivilegeReq),
				Source:   fullURL,
				Response: httputil.LimitResponse(lowPrivBody, maxResponseSize, "Response too large, please manually open the link to view."),
				Length:   len(lowPrivBody),
			})
		}
	})
	defer pool.Release()

	// 提交任务
	for _, api := range o.ApiList {
		wg.Add(1)
		pool.Invoke(api)
	}

	wg.Wait()
}

func buildRawRequest(req APIRequest) string {
	var sb strings.Builder

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return "" // URL解析失败直接返回空
	}

	// 如果是 GET 请求，需要将参数附加到 URL 后面
	if req.Method == http.MethodGet && req.Params != nil {
		// 编码查询参数并将其附加到 URL
		queryString := req.Params.Encode()
		if parsedURL.RawQuery == "" {
			parsedURL.RawQuery = queryString
		} else {
			parsedURL.RawQuery += "&" + queryString
		}
	}

	// 写入请求行，只保留Path+Query，不要域名
	pathWithQuery := parsedURL.Path
	if parsedURL.RawQuery != "" {
		pathWithQuery += "?" + parsedURL.RawQuery
	}
	sb.WriteString(fmt.Sprintf("%s %s HTTP/1.1\n", req.Method, pathWithQuery))

	// Host 头
	sb.WriteString(fmt.Sprintf("Host: %s\n", parsedURL.Host))

	// 其他 Headers
	for k, v := range req.Headers {
		// 避免重复写 Host
		if strings.ToLower(k) == "host" {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	sb.WriteString("\n") // 头结束空一行

	// 写入请求体
	if req.Method != http.MethodGet && req.Params != nil {
		// 如果是非 GET 请求，参数放到请求体中
		sb.WriteString(req.Params.Encode())
	}
	if req.Body != "" {
		// 如果请求体不为空，则直接写入
		sb.WriteString(req.Body)
	}
	return sb.String()
}
