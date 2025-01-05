package clients

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"time"

	"golang.org/x/net/proxy"
)

type Proxy struct {
	Enabled  bool
	Mode     string
	Address  string
	Port     int
	Username string
	Password string
}

// 选择代理模式，返回http.client
func SelectProxy(pr *Proxy, client *http.Client) (*http.Client, error) {
	if pr.Mode == "HTTP" {
		urlproxy, _ := url.Parse(fmt.Sprintf("%v://%v:%v", pr.Mode, pr.Address, pr.Port)) //"https://127.0.0.1:9743"
		client.Transport = &http.Transport{
			Proxy:                 http.ProxyURL(urlproxy),
			TLSClientConfig:       TlsConfig, // 防止HTTPS报错
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	} else {
		auth := &proxy.Auth{User: pr.Username, Password: pr.Password}
		dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", pr.Address, pr.Port), auth, proxy.Direct) //"127.0.0.1:9742"
		if err != nil {
			return nil, errors.New("socks address not available")
		}
		client.Transport = &http.Transport{
			TLSClientConfig:       TlsConfig, // 防止HTTPS报错
			Dial:                  dialer.Dial,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// 设置sock5
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := dialer.Dial(network, addr)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		}

	}
	return client, nil
}
