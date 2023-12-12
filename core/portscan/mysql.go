package portscan

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func MysqlScan(host string, usernames, passwords []string) *Burte {
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MysqlConn(host, user, pass)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     host,
					Protocol: "mysql",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "mysql",
		Username: "",
		Password: "",
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
