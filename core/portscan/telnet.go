package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/gotelnet"
	"slack-wails/lib/structs"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func TelenetScan(ctx context.Context, host string, usernames, passwords []string) {
	h := strings.Split(host, ":")[0]
	p, _ := strconv.Atoi(strings.Split(host, ":")[1])
	serverType := getTelnetServerType(h, p)
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := TelnetConn(h, user, pass, p, serverType)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "telnet weak password",
					Name:     "telnet weak password",
					URL:      host,
					Type:     "Telnet",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("telnet://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func getTelnetServerType(ip string, port int) int {
	client := gotelnet.New(ip, port)
	err := client.Connect()
	if err != nil {
		return gotelnet.Closed
	}
	defer client.Close()
	return client.MakeServerType()
}

func TelnetConn(addr, user, pass string, port, serverType int) (flag bool, err error) {
	flag = false
	client := gotelnet.New(addr, port)
	err = client.Connect()
	if err != nil {
		return flag, err
	}
	defer client.Close()
	client.UserName = user
	client.Password = pass
	client.ServerType = serverType
	err = client.Login()
	if err != nil {
		return flag, err
	}
	flag = true
	return flag, err
}
