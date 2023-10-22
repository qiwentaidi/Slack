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
	_ "github.com/lib/pq"
)

func PostgresScan(host string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["postgresql"])
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := PostgresConn(host, user, pass)
			if flag && err == nil {
				custom.Console.Append(fmt.Sprintf("[+] psql://%v %v:%v\n", host, user, pass))
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] psql %v %v %v %v", host, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
			}
		}
	}
}

func PostgresConn(host, user, pass string) (flag bool, err error) {
	flag = false
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v/postgres?sslmode=disable", user, pass, host)
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(common.Profile.PortScan.Timeout) * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			flag = true
			common.PortBurstResult = append(common.PortBurstResult, []string{"Postgresql", host, user, pass})
		}
	}
	return flag, err
}
