package portscan

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MysqlScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MysqlConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "mysql weak password",
					Name:     "mysql weak password",
					URL:      host,
					Type:     "Mysql",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
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
	for _, database := range []string{"mysql", "information_schema"} {
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&timeout=%v", user, pass, host, database, 10*time.Second)
		db, err := sql.Open("mysql", dataSourceName)
		if err == nil {
			db.SetConnMaxLifetime(10 * time.Second)
			db.SetConnMaxIdleTime(10 * time.Second)
			db.SetMaxIdleConns(0)
			err = db.Ping()
			if err == nil {
				flag = true
				db.Close()
				break
			}
		}
		db.Close()
	}
	return flag, err
}
