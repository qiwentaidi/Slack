// https://github.com/shadow1ng/fscan/blob/main/Plugins/ActiveMQ.go

package portscan

import (
	"context"
	"fmt"
	"net"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ActiveMQScan(ctx, ctrlCtx context.Context, taskId, address string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[activemq] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := ActiveMQConn(address, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "activemq weak password",
					Name:     "activemq weak password",
					URL:      address,
					Type:     "ActiveMQ",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("activemq://%s %s:%s is login failed", address, user, pass))
			}
		}
	}
}

// ActiveMQConn 尝试ActiveMQ连接
func ActiveMQConn(address, user, pass string) (bool, error) {
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	dialer := &net.Dialer{Timeout: time.Duration(10) * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// 创建结果通道
	resultChan := make(chan struct {
		success bool
		err     error
	}, 1)

	// 在协程中处理认证
	go func() {
		// STOMP协议的CONNECT命令
		stompConnect := fmt.Sprintf("CONNECT\naccept-version:1.0,1.1,1.2\nhost:/\nlogin:%s\npasscode:%s\n\n\x00", user, pass)

		// 发送认证请求
		conn.SetWriteDeadline(time.Now().Add(time.Duration(10) * time.Second))
		if _, err := conn.Write([]byte(stompConnect)); err != nil {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success bool
				err     error
			}{false, err}:
			}
			return
		}

		// 读取响应
		conn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Second))
		respBuf := make([]byte, 1024)
		n, err := conn.Read(respBuf)
		if err != nil {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success bool
				err     error
			}{false, err}:
			}
			return
		}

		// 检查认证结果
		response := string(respBuf[:n])

		var success bool
		var resultErr error

		if strings.Contains(response, "CONNECTED") {
			success = true
			resultErr = nil
		} else if strings.Contains(response, "Authentication failed") || strings.Contains(response, "ERROR") {
			success = false
			resultErr = fmt.Errorf("认证失败")
		} else {
			success = false
			resultErr = fmt.Errorf("未知响应: %s", response)
		}

		select {
		case <-ctx.Done():
		case resultChan <- struct {
			success bool
			err     error
		}{success, resultErr}:
		}
	}()

	// 等待认证结果或上下文取消
	select {
	case result := <-resultChan:
		return result.success, result.err
	case <-ctx.Done():
		return false, ctx.Err()
	}
}
