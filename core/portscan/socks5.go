package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const defaultAliveURL = "http://www.baidu.com"

func Socks5Scan(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	hostwithoutport := strings.Split(host, ":")[0]
	port, err := strconv.Atoi(strings.Split(host, ":")[1])
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("socks5://%s is invalid port", host))
		return
	}
	flag := Socks5Conn(hostwithoutport, port, 3, "", "", defaultAliveURL)
	if flag {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			TaskId:   taskId,
			ID:       "socks5 unauthorized",
			Name:     "socks5 unauthorized",
			URL:      host,
			Type:     "socks5",
			Severity: "HIGH",
		})
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("socks5://%s is no unauthorized access", host))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[socks5] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag = Socks5Conn(hostwithoutport, port, 3, user, pass, defaultAliveURL)
			if flag {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "socks5 weak password",
					Name:     "socks5 weak password",
					URL:      host,
					Type:     "SOCKS5",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("socks5://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func Socks5Conn(ip string, port, timeout int, username, password, aliveURL string) bool {
	client := clients.NewRestyClientWithProxy(nil, true, clients.Proxy{
		Enabled:  true,
		Mode:     "SOCK5",
		Address:  ip,
		Port:     port,
		Username: username,
		Password: password,
	})
	_, err := clients.DoRequest("GET", aliveURL, nil, nil, timeout, client)
	return err == nil
}
