package portscan

import (
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/mitchellh/go-vnc"
)

func VncScan(host string, associate bool, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	for _, pass := range passwords {
		pass = strings.Replace(pass, "{user}", "vnc", -1)
		flag, err := VncConn(host, pass)
		if flag && err == nil {
			custom.Console.Append(fmt.Sprintf("[+] vnc %v %v\n", host, pass))
			return
		} else {
			atomic.AddInt64(&counter, 1)
			custom.LogProgress(counter, len(passwords), fmt.Sprintf("[-] vnc %v %v %v\n", host, pass, err))
			if time.Now().Unix()-custom.LogTime > int64(len(passwords)*common.Profile.PortScan.Timeout) {
				return
			}
		}
	}
}

func VncConn(host, pass string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", host, time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	if err != nil {
		return flag, err
	}
	// Create a VNC client connection
	cfg := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{&vnc.PasswordAuth{Password: pass}},
	}
	client, err := vnc.Client(conn, cfg)
	if err != nil {
		return flag, err
	}
	// Use the client connection to interact with the VNC server
	// ...

	// Close the client connection when done
	client.Close()
	flag = true
	common.PortBurstResult = append(common.PortBurstResult, []string{"VNC", host, "", pass})
	return flag, err
}
