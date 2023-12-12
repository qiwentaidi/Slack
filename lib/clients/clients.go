package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

func NewRequest(method, url string, headers http.Header, body io.Reader, timeout int, client *http.Client) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	if headers != nil {
		r.Header = headers
	} else {
		r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	}
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		if b, err := io.ReadAll(resp.Body); err == nil {
			return resp, b, nil
		} else {
			return nil, nil, err
		}
	}
	return nil, nil, errors.New("resp is nil")
}

var (
	HTTP_PREFIX  = "http://"
	HTTPS_PREFIX = "https://"
)

// return error if host is not living
// or if host is live return http(s) url
func CheckProtocol(host string, client *http.Client) (string, error) {
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
		_, _, err := NewRequest("GET", HTTP_PREFIX+host, nil, nil, 10, client)
		if err != nil {
			return result, err
		}
		return HTTP_PREFIX + host, nil
	case parsePort == "443":
		_, _, err := NewRequest("GET", HTTPS_PREFIX+host, nil, nil, 10, client)
		if err != nil {
			return result, err
		}

		return HTTPS_PREFIX + host, nil

	default:
		_, _, err := NewRequest("GET", HTTPS_PREFIX+host, nil, nil, 10, client)
		if err == nil {
			return HTTPS_PREFIX + host, err
		}

		_, body, err := NewRequest("GET", HTTP_PREFIX+host, nil, nil, 10, client)
		if err == nil {
			if strings.Contains(string(body), "<title>400 The plain HTTP request was sent to HTTPS port</title>") {
				return HTTPS_PREFIX + host, nil
			}
			return HTTP_PREFIX + host, nil
		}

	}
	return "", fmt.Errorf("host %q is empty", host)
}
