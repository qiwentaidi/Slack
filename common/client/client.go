package client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"slack/common"
	"time"
)

// 返回body
func NewHttpWithDefaultHead(method, url string, client *http.Client) (*http.Response, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultWebTimeout*time.Second)
	defer cancel()
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	if resp != nil && resp.StatusCode != 302 {
		defer resp.Body.Close()
		if b, err := io.ReadAll(resp.Body); err == nil {
			return resp, b, nil
		} else {
			return nil, nil, err
		}
	}
	return nil, nil, err
}

// 跟随页面跳转最多10次
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 防止HTTPS报错
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
