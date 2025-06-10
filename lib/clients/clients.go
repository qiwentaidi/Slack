package clients

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slack-wails/lib/utils"
	"slack-wails/lib/utils/randutil"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var TlsConfig = &tls.Config{
	InsecureSkipVerify: true,             // 防止HTTPS报错
	MinVersion:         tls.VersionTLS10, // 最低支持TLS 1.0
	CipherSuites: []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	},
}

func NewRestyClient(interfaceIp net.IP, followRedirect bool) *resty.Client {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	if interfaceIp != nil {
		dialer.LocalAddr = &net.TCPAddr{IP: interfaceIp}
	}

	transport := &http.Transport{
		TLSClientConfig:       TlsConfig,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := resty.New().
		SetTransport(transport).
		SetTimeout(10*time.Second).
		SetHeader("User-Agent", randutil.RandomUA()).
		SetHeader("Connection", "close")
	// 设置重定向规则
	if followRedirect {
		client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
	} else {
		client.SetRedirectPolicy(resty.NoRedirectPolicy())
	}
	return client
}

func NewRestyClientWithProxy(interfaceIp net.IP, followRedirect bool, proxy Proxy) *resty.Client {
	client := NewRestyClient(interfaceIp, followRedirect)
	if proxy.Enabled {
		proxyURL := GetRawProxy(proxy)
		client.SetProxy(proxyURL)
	}
	return client
}

func DoRequest(method, url string, headers map[string]string, body io.Reader, timeout int, client *resty.Client) (*resty.Response, error) {
	if timeout > 0 {
		client.SetTimeout(time.Duration(timeout) * time.Second)
	}

	req := client.R()

	if headers != nil {
		req.SetHeaders(headers)
	}

	if body != nil {
		// 将 io.Reader 转换为 []byte 读入
		data, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("read body failed: %w", err)
		}
		req.SetBody(data)
	}

	var resp *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = req.Get(url)
	case http.MethodPost:
		resp, err = req.Post(url)
	case http.MethodPut:
		resp, err = req.Put(url)
	case http.MethodDelete:
		resp, err = req.Delete(url)
	case http.MethodPatch:
		resp, err = req.Patch(url)
	case http.MethodOptions:
		resp, err = req.Options(url)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return resp, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func SimpleGet(url string, client *resty.Client) (*resty.Response, error) {
	return DoRequest("GET", url, nil, nil, 10, client)
}

var regTitle = regexp.MustCompile(`(?is)<title\b[^>]*>(.*?)</title>`)

func GetTitle(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	if match := regTitle.FindSubmatch(body); len(match) > 1 {
		return strings.TrimSpace(utils.Str2UTF8(string(match[1])))
	}
	return ""
}

func Str2HeadersMap(str string) map[string]string {
	headers := make(map[string]string)
	if str == "" {
		return headers
	}
	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if i := strings.IndexByte(line, ':'); i > 0 {
			headers[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
		}
	}
	return headers
}

func Str2HeaderList(str string) []string {
	headers := make([]string, 0)
	if str == "" {
		return headers
	}
	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		headers = append(headers, line)
	}
	return headers
}

const (
	HTTP_PREFIX  = "http://"
	HTTPS_PREFIX = "https://"
)

// 检测协议，返回完整URL（带 http(s)）
func CheckProtocol(host string, client *resty.Client) (string, error) {
	var result string

	if len(strings.TrimSpace(host)) == 0 {
		return result, fmt.Errorf("host %q is empty", host)
	}

	if strings.HasPrefix(host, HTTPS_PREFIX) || strings.HasPrefix(host, HTTP_PREFIX) {
		_, err := SimpleGet(host, client)
		if err != nil {
			return result, err
		}
		return host, nil
	}

	u, err := url.Parse(HTTP_PREFIX + host)
	if err != nil {
		return result, err
	}
	parsePort := u.Port()

	switch {
	case parsePort == "80":
		_, err := SimpleGet(HTTP_PREFIX+host, client)
		if err != nil {
			return result, err
		}
		return HTTP_PREFIX + host, nil

	case parsePort == "443":
		_, err := SimpleGet(HTTPS_PREFIX+host, client)
		if err != nil {
			return result, err
		}
		return HTTPS_PREFIX + host, nil

	default:
		// 先试 https
		_, err := SimpleGet(HTTPS_PREFIX+host, client)
		if err == nil {
			return HTTPS_PREFIX + host, nil
		}
		// 再试 http
		resp, err := SimpleGet(HTTP_PREFIX+host, client)
		if err == nil {
			if strings.Contains(string(resp.Body()), "<title>400 The plain HTTP request was sent to HTTPS port</title>") {
				return HTTPS_PREFIX + host, nil
			}
			return HTTP_PREFIX + host, nil
		}
	}

	return "", fmt.Errorf("both http and https check failed for host: %s", host)
}
