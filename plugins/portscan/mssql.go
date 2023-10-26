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
	_ "github.com/denisenkom/go-mssqldb"
)

func MssqlScan(host string, associate bool, usertext, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["mssql"])
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(host, user, pass)
			if flag && err == nil {
				custom.Console.Append(fmt.Sprintf("[+] mssql://%v %v:%v\n", host, user, pass))
				return
			} else {
				atomic.AddInt64(&counter, 1)
				custom.LogProgress(counter, len(usernames)*len(passwords), fmt.Sprintf("[-] mssql %v %v:%v %v", host, user, pass, err))
				if time.Now().Unix()-custom.LogTime > int64(len(usernames)*len(passwords)*common.Profile.PortScan.Timeout) {
					return
				}
			}
		}
	}

}

func MssqlConn(host, user, pass string) (flag bool, err error) {
	flag = false
	Host, Port := strings.Split(host, ":")[0], strings.Split(host, ":")[1]
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", Host, user, pass, Port, time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(common.Profile.PortScan.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(common.Profile.PortScan.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			flag = true
			common.PortBurstResult = append(common.PortBurstResult, []string{"Mssql", host, user, pass, ""})
		}
	}
	return flag, err
}
