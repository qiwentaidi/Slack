package portscan

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

func OracleScan(host string, usernames, passwords []string) *Burte {
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(host, user, pass)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     host,
					Protocol: "oracle",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "oracle",
		Username: "",
		Password: "",
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
