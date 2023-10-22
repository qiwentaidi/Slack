package portscan

import (
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/lib/gotelnet"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
)

func TelenetScan(host string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["telnet"])
	h := strings.Split(host, ":")[0]
	p, _ := strconv.Atoi(strings.Split(host, ":")[1])
	serverType := getTelnetServerType(h, p)
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := TelnetConn(h, user, pass, p, serverType)
			if flag && err == nil {
				custom.Console.Append(fmt.Sprintf("[+] telnet %v %v:%v", host, user, pass))
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] telnet %v %v:%v %v", host, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
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
	common.PortBurstResult = append(common.PortBurstResult, []string{"Telnet", fmt.Sprintf("%v:%v", addr, port), user, pass})
	return flag, err
}
