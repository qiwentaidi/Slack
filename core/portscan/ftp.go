package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func FtpScan(ctx context.Context, address string, usernames, passwords []string) {
	flag, err := FtpConn(address, "anonymous", "")
	if flag && err == nil {
		runtime.EventsEmit(ctx, "bruteResult", Burte{
			Status:   true,
			Host:     address,
			Protocol: "ftp",
			Username: "anonymous",
			Password: "",
		})
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("ftp://%s is no unauthorized access", address))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := FtpConn(address, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     address,
					Protocol: "ftp",
					Username: user,
					Password: pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("ftp://%s %s:%s is login failed", address, user, pass))
			}
		}
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
