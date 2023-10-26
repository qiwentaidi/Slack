package portscan

import (
	"database/sql"
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
	_ "github.com/sijms/go-ora/v2"
)

func OracleScan(host string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["oracle"])
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(host, user, pass)
			if flag && err == nil {
				custom.Console.Append(fmt.Sprintf("[+] oracle://%v %v:%v\n", host, user, pass))
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] oracle %v %v:%v %v", host, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
			}
		}
	}
}

func OracleConn(host, user, pass string) (flag bool, err error) {
	flag = false
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s/orcl", user, pass, host)
	db, err := sql.Open("oracle", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(1) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(1) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			flag = true
			common.PortBurstResult = append(common.PortBurstResult, []string{"Oracle", host, user, pass, ""})
		}
	}
	return flag, err
}
