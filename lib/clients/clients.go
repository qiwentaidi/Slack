package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"regexp"
	"slack-wails/lib/util"
	"strings"
	"time"
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
	},
}

func NewHttpClient(interfaceIp net.IP, followRedirect bool) *http.Client {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	if interfaceIp != nil {
		dialer.LocalAddr = &net.TCPAddr{
			IP: interfaceIp,
		}
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       TlsConfig,
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	return client
}

func NewHttpClientWithProxy(interfaceIp net.IP, followRedirect bool, proxy Proxy) *http.Client {
	client := NewHttpClient(interfaceIp, followRedirect)
	if proxy.Enabled {
		client, _ = SelectProxy(&proxy, client)
	}
	return client
}

func NewRequest(method, url string, headers map[string]string, body io.Reader, timeout int, closeRespBody bool, client *http.Client) (*http.Response, []byte, error) {
	requestTimeout := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("User-Agent", util.RandomUA())
	r.Header.Set("Connection", "close")
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, errors.New("response is nil, possible network error or timeout")
	}
	if closeRespBody {
		defer resp.Body.Close()
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle "unexpected EOF" as a specific error case
		if err.Error() == "unexpected EOF" {
			return resp, bodyBytes, err
		}
		return nil, nil, err
	}

	return resp, bodyBytes, nil
}

func NewSimpleGetRequest(url string, client *http.Client) (*http.Response, []byte, error) {
	return NewRequest("GET", url, nil, nil, 10, true, client)
}

var regTitle = regexp.MustCompile(`(?is)<title\b[^>]*>(.*?)</title>`)

func GetTitle(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	if match := regTitle.FindSubmatch(body); len(match) > 1 {
		return strings.TrimSpace(util.Str2UTF8(string(match[1])))
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
