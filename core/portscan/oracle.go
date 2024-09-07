package portscan

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func OracleScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     host,
					Protocol: "oracle",
					Username: user,
					Password: pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("oracle://%s %s:%s is login failed", host, user, pass))
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
		}
	}
	return flag, err
}
