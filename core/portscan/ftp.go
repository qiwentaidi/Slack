package portscan

import (
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

func FtpScan(address string, usernames, passwords []string) *Burte {
	flag, err := FtpConn(address, "anonymous", "")
	if flag && err == nil {
		return &Burte{
			Status:   true,
			Host:     address,
			Protocol: "ftp",
			Username: "anonymous",
			Password: "",
		}
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := FtpConn(address, user, pass)
			if flag && err == nil {
				return &Burte{
					Status:   true,
					Host:     address,
					Protocol: "ftp",
					Username: user,
					Password: pass,
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     address,
		Protocol: "ftp",
		Username: "",
		Password: "",
	}
}

func FtpConn(address, user, pass string) (flag bool, err error) {
	flag = false
	conn, err := ftp.Dial(address, ftp.DialWithTimeout(12*time.Second))
	if err == nil {
		err = conn.Login(user, pass)
		if err == nil {
			flag = true
		}
	}
	return flag, err
}
