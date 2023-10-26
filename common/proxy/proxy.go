package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"slack/common"
	"slack/common/logger"
	"time"

	"golang.org/x/net/proxy"
)

// 选择代理模式，返回http.client
func SelectProxy(profile *common.Profiles) (client *http.Client) {
	if profile.Proxy.Mode == "HTTP" {
		urli := url.URL{}
		urlproxy, _ := urli.Parse(fmt.Sprintf("%v://%v:%v", profile.Proxy.Mode, profile.Proxy.Address, profile.Proxy.Port)) //"https://127.0.0.1:9743"
		client = &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyURL(urlproxy),
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // 防止HTTPS报错
				TLSHandshakeTimeout: time.Second * 3,
			},
		}
	} else {
		auth := &proxy.Auth{User: profile.Proxy.Username, Password: profile.Proxy.Password}
		dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", profile.Proxy.Address, profile.Proxy.Port), auth, proxy.Direct) //"127.0.0.1:9742"
		if err != nil {
			logger.Debug(err)
			return nil
		}
		httpTransport := &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // 防止HTTPS报错
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: time.Second * 3,
		}
		client = &http.Client{Transport: httpTransport}
		// 设置sock5
		httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.Dial(network, addr)
			if err != nil {
				logger.Debug(err)
				return nil, err
			}
			return conn, nil
		}
	}
	return client
}
