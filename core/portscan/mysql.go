package portscan

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MysqlScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MysqlConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     host,
					Protocol: "mysql",
					Username: user,
					Password: pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("mysql://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func MysqlConn(host, user, pass string) (flag bool, err error) {
	flag = false
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/mysql?charset=utf8&timeout=%v", user, pass, host, 10*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(10 * time.Second)
		db.SetConnMaxIdleTime(10 * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			flag = true
		}
	}
	return flag, err
}
