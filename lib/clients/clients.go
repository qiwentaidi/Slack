package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slack-wails/lib/util"
	"strings"
	"time"
)

var TlsConfig = &tls.Config{
	InsecureSkipVerify: true,             // 防止HTTPS报错
	MinVersion:         tls.VersionTLS10, // 最低支持TLS 1.0
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

var (
	HTTP_PREFIX  = "http://"
	HTTPS_PREFIX = "https://"
)

// return error if host is not living
// or if host is live return http(s) url
func IsWeb(host string, client *http.Client) (string, error) {
	var result string
	if len(strings.TrimSpace(host)) == 0 {
		return result, fmt.Errorf("host %q is empty", host)
	}
	u, err := url.Parse(HTTP_PREFIX + host)
	if err != nil {
		return result, err
	}
	parsePort := u.Port()
	switch {
	case parsePort == "80":
		_, _, err := NewSimpleGetRequest(HTTP_PREFIX+host, client)
		if err != nil {
			return result, err
		}
		return HTTP_PREFIX + host, nil
	case parsePort == "443":
		_, _, err := NewSimpleGetRequest(HTTPS_PREFIX+host, client)
		if err != nil {
			return result, err
		}

		return HTTPS_PREFIX + host, nil

	default:
		_, _, err := NewSimpleGetRequest(HTTPS_PREFIX+host, client)
		if err == nil {
			return HTTPS_PREFIX + host, err
		}

		_, body, err := NewSimpleGetRequest(HTTP_PREFIX+host, client)
		if err == nil {
			if strings.Contains(string(body), "<title>400 The plain HTTP request was sent to HTTPS port</title>") {
				return HTTPS_PREFIX + host, nil
			}
			return HTTP_PREFIX + host, nil
		}

	}
	return "", fmt.Errorf("host %q is empty", host)
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
