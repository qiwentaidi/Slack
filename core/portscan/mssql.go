package portscan

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MssqlScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "mssql weak password",
					Name:     "mssql weak password",
					URL:      host,
					Type:     "Mssql",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("mssql://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func MssqlConn(host, user, pass string) (flag bool, err error) {
	flag = false
	Host, Port := strings.Split(host, ":")[0], strings.Split(host, ":")[1]
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", Host, user, pass, Port, 10*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
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
