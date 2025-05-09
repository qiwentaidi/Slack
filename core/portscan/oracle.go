package portscan

import (
	"context"
	"database/sql"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const defaultOracleServerName = "orcl"

func OracleScan(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[oracle] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(host, defaultOracleServerName, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "oracle weak password",
					Name:     "oracle weak password",
					URL:      host,
					Type:     "Oracle",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("oracle://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func OracleConn(host, servername, user, pass string) (flag bool, err error) {
	flag = false
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s/%s", user, pass, host, servername)
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
