package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func RedisScan(ctx context.Context, host string, usernames, passwords []string) {
	flag, err := RedisUnauth(host)
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			ID:       "redis unauthorized",
			Name:     "redis unauthorized",
			URL:      host,
			Type:     "Redis",
			Severity: "HIGH",
		})
		gologger.Success(ctx, fmt.Sprintf("redis://%s is unauthorized access", host))
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("redis://%s is no unauthorized access", host))
	}
	for _, pass := range passwords {
		if ExitFunc {
			return
		}
		pass = strings.Replace(pass, "{user}", "redis", -1)
		flag, err := RedisConn(host, pass)
		if flag && err == nil {
			runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
				ID:       "redis weak password",
				Name:     "redis weak password",
				URL:      host,
				Type:     "Redis",
				Severity: "HIGH",
				Extract:  pass,
			})
			return
		} else {
			gologger.Info(ctx, fmt.Sprintf("redis://%s %s is login failed", host, pass))
		}
	}
}

func RedisConn(address, password string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, 10*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", password)))
	if err != nil {
		return flag, err
	}
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return flag, err
	}
	if strings.Contains(string(buffer[:n]), "+OK") {
		flag = true
	}
	return flag, err
}

func RedisUnauth(address string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, 10*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte("info\r\n"))
	if err != nil {
		return flag, err
	}
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return flag, err
	}
	if strings.Contains(string(buffer[:n]), "redis_version") {
		flag = true
	}
	return flag, err
}
