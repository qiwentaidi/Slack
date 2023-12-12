package portscan

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func MssqlScan(host string, usernames, passwords []string) *Burte {
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(host, user, pass)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     host,
					Protocol: "mssql",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "mssql",
		Username: "",
		Password: "",
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
