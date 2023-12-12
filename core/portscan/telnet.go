package portscan

import (
	"slack-wails/lib/gotelnet"
	"strconv"
	"strings"
)

func TelenetScan(host string, usernames, passwords []string) *Burte {
	h := strings.Split(host, ":")[0]
	p, _ := strconv.Atoi(strings.Split(host, ":")[1])
	serverType := getTelnetServerType(h, p)
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := TelnetConn(h, user, pass, p, serverType)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     host,
					Protocol: "telnet",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "telnet",
		Username: "",
		Password: "",
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
