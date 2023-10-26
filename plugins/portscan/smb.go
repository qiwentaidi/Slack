package portscan

import (
	"errors"
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/stacktitan/smb/smb"
)

func SmbScan(host, domain string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["smb"])
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(host, user, pass)
			if flag && err == nil {
				common.PortBurstResult = append(common.PortBurstResult, []string{"SMB", host, user, pass, ""})
				if domain != "" {
					custom.Console.Append(fmt.Sprintf("[+] smb:%v:%v\\%v:%v", host, domain, user, pass))
				} else {
					custom.Console.Append(fmt.Sprintf("[+] smb:%v %v:%v", host, user, pass))
				}
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] smb %v %v:%v %v", host, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
			}
		}
	}
}

func SmblConn(host, user, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host := strings.Split(host, ":")[0]
	Port, _ := strconv.Atoi(strings.Split(host, ":")[1])
	options := smb.Options{
		Host:        Host,
		Port:        Port,
		User:        user,
		Password:    pass,
		Domain:      "",
		Workstation: "",
	}
	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}

func doWithTimeOut(host, user, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = SmblConn(host, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(time.Duration(common.Profile.PortScan.Timeout) * time.Second):
		return false, errors.New("time out")
	}
}
