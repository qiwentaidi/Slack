package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slack-wails/lib/util"
	"strings"
	"time"
)

// type HTTP_OPTIONS struct {
// 	HTTPS    bool
// 	Redirect bool
// 	Proxys   Proxy
// }

func TestErrorClient() *http.Client {
	client, _ := SelectProxy(&Proxy{
		Enabled: true,
		Mode:    "HTTP",
		Address: "127.0.0.1",
		Port:    8080}, DefaultClient())
	return client
}

// 跟随页面跳转最多10次
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 防止HTTPS报错
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

// 不跟随页面跳转
func NotFollowClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func NewRequest(method, url string, headers http.Header, body io.Reader, timeout int, closeRespBody bool, client *http.Client) (*http.Response, []byte, error) {
	requestTimeout := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	if headers != nil {
		r.Header = headers
	} else {
		r.Header.Set("User-Agent", util.RandomUA())
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
	if b, err := io.ReadAll(resp.Body); err == nil {
		return resp, b, nil
	} else {
		if err.Error() == "unexpected EOF" {
			return nil, b, err
		}
		return nil, nil, err
	}
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
