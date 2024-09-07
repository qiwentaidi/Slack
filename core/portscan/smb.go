package portscan

import (
	"context"
	"errors"
	"fmt"
	"slack-wails/lib/gologger"
	"strconv"
	"strings"
	"time"

	"github.com/stacktitan/smb/smb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SmbScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     host,
					Protocol: "smb",
					Username: user,
					Password: pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("smb://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func SmblConn(host, user, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host := strings.Split(host, ":")[0]
	Port, _ := strconv.Atoi(strings.Split(host, ":")[1])
	options := smb.Options{
		Host:        Host,
		Port:        Port,
		User:        user,
		Password:    pass,
		Domain:      "",
		Workstation: "",
	}
	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}

func doWithTimeOut(host, user, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = SmblConn(host, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(10 * time.Second):
		return false, errors.New("time out")
	}
}
