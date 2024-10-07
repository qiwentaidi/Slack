package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func Socks5Scan(ctx context.Context, host string, usernames, passwords []string) {
	hostwithoutport := strings.Split(host, ":")[0]
	port, err := strconv.Atoi(strings.Split(host, ":")[1])
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("socks5://%s is invalid port", host))
		return
	}
	flag := Socks5Conn(hostwithoutport, port, 3, "", "")
	if flag {
		runtime.EventsEmit(ctx, "bruteResult", Burte{
			Status:   false,
			Host:     host,
			Protocol: "socks5",
			Username: "",
			Password: "",
		})
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("socks5://%s is no unauthorized access", host))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag = Socks5Conn(hostwithoutport, port, 3, user, pass)
			if flag {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     host,
					Protocol: "socks5",
					Username: user,
					Password: pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("socks5://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func Socks5Conn(ip string, port, timeout int, username, password string) bool {
	client, err := clients.SelectProxy(&clients.Proxy{
		Enabled:  true,
		Mode:     "SOCK5",
		Address:  ip,
		Port:     port,
		Username: username,
		Password: password,
	}, clients.DefaultClient())
	if err != nil {
		return false
	}
	_, _, err = clients.NewRequest("GET", "http://www.baidu.com/", nil, nil, timeout, true, client)
	return err == nil
}
