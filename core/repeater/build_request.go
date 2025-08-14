package repeater

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/qiwentaidi/clients"
)

func SendRequestWithRaw(raw string, forceHttps, redirect bool, proxyURL string) (*resty.Response, int64, error) {
	// 1. 统一换行符为 \n，方便处理
	raw = strings.ReplaceAll(raw, "\r\n", "\n")

	// 2. 去掉每行前面的空格和 tab
	lines := strings.Split(raw, "\n")
	for i := range lines {
		lines[i] = strings.TrimLeft(lines[i], " \t")
	}
	raw = strings.Join(lines, "\n")

	// 3. 统一换回 HTTP 标准的 \r\n
	raw = strings.ReplaceAll(raw, "\n", "\r\n")

	// 4. 确保以 \r\n\r\n 结尾（表示 header 结束）
	if !strings.HasSuffix(raw, "\r\n\r\n") {
		raw += "\r\n\r\n"
	}
	reader := bufio.NewReader(strings.NewReader(raw))

	req, err := http.ReadRequest(reader)
	if err != nil {
		return nil, 0, err
	}
	scheme := "http"
	if forceHttps {
		scheme = "https"
	}
	client := clients.NewRestyClientWithProxy(nil, redirect, proxyURL)
	// 不跟随重定向
	if !redirect {
		client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}))
	}

	// ===== 记录开始时间 =====
	start := time.Now()
	resp, err := ConvertHttpRequestToResty(client, req, scheme)
	if err != nil {
		return nil, 0, err
	}
	// ===== 计算耗时 =====
	duration := time.Since(start).Milliseconds()
	return resp, duration, nil
}
func ConvertHttpRequestToResty(client *resty.Client, req *http.Request, scheme string) (*resty.Response, error) {
	client.SetDoNotParseResponse(true) // 禁用自动解析响应
	// 构造 headers
	headers := make(map[string]string)
	for key, values := range req.Header {
		// resty.SetHeaders 只支持一个值，这里取第一个，多个值可自行拼接
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// body
	var body io.Reader
	if req.Body != nil {
		// 注意：不能直接传 req.Body，因为可能会被消费掉，后续还要用
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置 body
		body = bytes.NewBuffer(bodyBytes)
	}

	// query 参数拼接在 URL 上
	rawURL := buildFullURL(req, scheme)

	// 调用封装好的 DoRequest
	return clients.DoRequest(req.Method, rawURL, headers, body, 0, client)
}

// buildFullURL 根据 http.Request 和 scheme 构造完整 URL
func buildFullURL(req *http.Request, scheme string) string {
	var host string
	if scheme != "" {
		host = scheme + "://" + req.Host
	} else {
		referer := req.Referer()
		if strings.HasPrefix(referer, "http") {
			host = referer[:strings.Index(referer, "://")+3] + req.Host
		} else {
			host = "https://" + req.Host
		}
	}
	return host + req.URL.RequestURI()
}
