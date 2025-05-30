package portscan

import (
	"context"
	"errors"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strconv"
	"strings"
	"time"

	"github.com/stacktitan/smb/smb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SmbScan(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[smb] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "smb weak password",
					Name:     "smb weak password",
					URL:      host,
					Type:     "SMB",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
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
