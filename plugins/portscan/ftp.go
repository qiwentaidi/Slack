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
	"github.com/jlaffaye/ftp"
)

func FtpScan(address string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	flag, err := FtpConn(address, "anonymous", "")
	if flag && err == nil {
		custom.Console.Append(fmt.Sprintf("[+] ftp://%v %v\n", address, "anonymous"))
		return
	}
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["ftp"])
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := FtpConn(address, user, pass)
			if flag && err == nil {
				custom.Console.Append(fmt.Sprintf("[+] ftp://%v %v:%v\n", address, user, pass))
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] ftp://%v %v:%v %v", address, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
			}
		}
	}
}

func FtpConn(address, user, pass string) (flag bool, err error) {
	flag = false
	conn, err := ftp.Dial(address, ftp.DialWithTimeout(time.Duration(common.Profile.PortScan.Timeout)*time.Second))
	if err == nil {
		err = conn.Login(user, pass)
		if err == nil {
			flag = true
			common.PortBurstResult = append(common.PortBurstResult, []string{"FTP", address, user, pass})
		}
	}
	return flag, err
}
