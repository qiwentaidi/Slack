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
			Proxy:               http.ProxyURL(urlproxy),
			TLSClientConfig:     TlsConfig, // 防止HTTPS报错
			TLSHandshakeTimeout: time.Second * 10,
		}
	} else {
		auth := &proxy.Auth{User: pr.Username, Password: pr.Password}
		dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", pr.Address, pr.Port), auth, proxy.Direct) //"127.0.0.1:9742"
		if err != nil {
			return nil, errors.New("socks address not available")
		}
		client.Transport = &http.Transport{
			TLSClientConfig:     TlsConfig, // 防止HTTPS报错
			Dial:                dialer.Dial,
			TLSHandshakeTimeout: time.Second * 3,
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

func JudgeClient(proxy Proxy) *http.Client {
	client := DefaultClient()
	if proxy.Enabled {
		client, _ = SelectProxy(&proxy, client)
	}
	return client
}
