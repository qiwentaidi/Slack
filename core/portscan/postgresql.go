package portscan

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func PostgresScan(host string, usernames, passwords []string) *Burte {
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := PostgresConn(host, user, pass)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     host,
					Protocol: "postgres",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "postgres",
		Username: "",
		Password: "",
	}
}

func PostgresConn(host, user, pass string) (flag bool, err error) {
	flag = false
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v/postgres?sslmode=disable", user, pass, host)
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(10 * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			flag = true
		}
	}
	return flag, err
}
