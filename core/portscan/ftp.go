package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func FtpScan(ctx, ctrlCtx context.Context, taskId, address string, usernames, passwords []string) {
	flag, directories, err := FtpConn(address, "anonymous", "")
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			TaskId:   taskId,
			ID:       "ftp unauthorized",
			Name:     "ftp unauthorized",
			URL:      address,
			Type:     "FTP",
			Severity: "HIGH",
			Response: strings.Join(directories, "\n"),
		})
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("ftp://%s is no unauthorized access", address))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[ftp] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, directories, err := FtpConn(address, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "ftp weak password",
					Name:     "ftp weak password",
					URL:      address,
					Type:     "ftp",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
					Response: strings.Join(directories, "\n"),
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("ftp://%s %s:%s is login failed", address, user, pass))
			}
		}
	}
}

func FtpConn(address, user, pass string) (flag bool, directories []string, err error) {
	conn, err := ftp.Dial(address, ftp.DialWithTimeout(12*time.Second))
	if err != nil {
		return false, nil, err
	}
	defer func() {
		if conn != nil {
			conn.Quit()
		}
	}()

	err = conn.Login(user, pass)
	if err != nil {
		return false, nil, err
	}

	// 获取目录信息
	dirs, err := conn.List("")
	if err == nil && len(dirs) > 0 {
		directories = make([]string, 0, min(6, len(dirs)))
		for i := 0; i < len(dirs) && i < 6; i++ {
			name := dirs[i].Name
			if len(name) > 50 {
				name = name[:50]
			}
			directories = append(directories, name)
		}
	}

	return true, directories, nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
