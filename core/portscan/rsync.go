// https://github.com/shadow1ng/fscan/blob/main/Plugins/Rsync.go
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

func RsyncScan(ctx, ctrlCtx context.Context, taskId, address string, usernames, passwords []string) {
	flag, moduleName, err := RsyncConn(address, "", "")
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			TaskId:   taskId,
			ID:       "rsync unauthorized",
			Name:     "rsync unauthorized",
			URL:      address,
			Type:     "Rsync",
			Severity: "HIGH",
			Extract:  "moduleName: " + moduleName,
		})
		gologger.Success(ctx, fmt.Sprintf("rsync://%s is unauthorized access", address))
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("rsync://%s is no unauthorized access", address))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[rsync] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, moduleName, err = RsyncConn(address, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "rsync weak password",
					Name:     "rsync weak password",
					URL:      address,
					Type:     "Rsync",
					Severity: "HIGH",
					Extract:  "moduleName: " + moduleName,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("rsync://%s %s:%s is login failed", address, user, pass))
			}
		}
	}
}

// RsyncConn 尝试Rsync连接
func RsyncConn(address, user, pass string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()
	timeout := time.Duration(10) * time.Second

	// 设置带有上下文的拨号器
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// 建立连接
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return false, "", err
	}
	defer conn.Close()

	// 创建结果通道用于超时控制
	resultChan := make(chan struct {
		success    bool
		moduleName string
		err        error
	}, 1)

	// 在协程中处理连接，以支持上下文取消
	go func() {
		buffer := make([]byte, 1024)

		// 1. 读取服务器初始greeting
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := conn.Read(buffer)
		if err != nil {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success    bool
				moduleName string
				err        error
			}{false, "", err}:
			}
			return
		}

		greeting := string(buffer[:n])
		if !strings.HasPrefix(greeting, "@RSYNCD:") {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success    bool
				moduleName string
				err        error
			}{false, "", fmt.Errorf("不是Rsync服务")}:
			}
			return
		}

		// 获取服务器版本号
		version := strings.TrimSpace(strings.TrimPrefix(greeting, "@RSYNCD:"))

		// 2. 回应相同的版本号
		conn.SetWriteDeadline(time.Now().Add(timeout))
		_, err = conn.Write([]byte(fmt.Sprintf("@RSYNCD: %s\n", version)))
		if err != nil {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success    bool
				moduleName string
				err        error
			}{false, "", err}:
			}
			return
		}

		// 3. 选择模块 - 先列出可用模块
		conn.SetWriteDeadline(time.Now().Add(timeout))
		_, err = conn.Write([]byte("#list\n"))
		if err != nil {
			select {
			case <-ctx.Done():
			case resultChan <- struct {
				success    bool
				moduleName string
				err        error
			}{false, "", err}:
			}
			return
		}

		// 4. 读取模块列表
		var moduleList strings.Builder
		for {
			// 检查上下文是否取消
			select {
			case <-ctx.Done():
				return
			default:
			}

			conn.SetReadDeadline(time.Now().Add(timeout))
			n, err = conn.Read(buffer)
			if err != nil {
				break
			}
			chunk := string(buffer[:n])
			moduleList.WriteString(chunk)
			if strings.Contains(chunk, "@RSYNCD: EXIT") {
				break
			}
		}

		modules := strings.Split(moduleList.String(), "\n")
		for _, module := range modules {
			if strings.HasPrefix(module, "@RSYNCD") || module == "" {
				continue
			}

			// 获取模块名
			moduleName := strings.Fields(module)[0]

			// 检查上下文是否取消
			select {
			case <-ctx.Done():
				return
			default:
			}

			// 5. 为每个模块创建新连接尝试认证
			authConn, err := dialer.DialContext(ctx, "tcp", address)
			if err != nil {
				continue
			}
			defer authConn.Close()

			// 重复初始握手
			authConn.SetReadDeadline(time.Now().Add(timeout))
			_, err = authConn.Read(buffer)
			if err != nil {
				authConn.Close()
				continue
			}

			authConn.SetWriteDeadline(time.Now().Add(timeout))
			_, err = authConn.Write([]byte(fmt.Sprintf("@RSYNCD: %s\n", version)))
			if err != nil {
				authConn.Close()
				continue
			}

			// 6. 选择模块
			authConn.SetWriteDeadline(time.Now().Add(timeout))
			_, err = authConn.Write([]byte(moduleName + "\n"))
			if err != nil {
				authConn.Close()
				continue
			}

			// 7. 等待认证挑战
			authConn.SetReadDeadline(time.Now().Add(timeout))
			n, err = authConn.Read(buffer)
			if err != nil {
				authConn.Close()
				continue
			}

			authResponse := string(buffer[:n])
			if strings.Contains(authResponse, "@RSYNCD: OK") {
				// 模块不需要认证
				if user == "" && pass == "" {
					authConn.Close()
					select {
					case <-ctx.Done():
					case resultChan <- struct {
						success    bool
						moduleName string
						err        error
					}{true, moduleName, nil}:
					}
					return
				}
			} else if strings.Contains(authResponse, "@RSYNCD: AUTHREQD") {
				if user != "" && pass != "" {
					// 8. 发送认证信息
					authString := fmt.Sprintf("%s %s\n", user, pass)
					authConn.SetWriteDeadline(time.Now().Add(timeout))
					_, err = authConn.Write([]byte(authString))
					if err != nil {
						authConn.Close()
						continue
					}

					// 9. 读取认证结果
					authConn.SetReadDeadline(time.Now().Add(timeout))
					n, err = authConn.Read(buffer)
					if err != nil {
						authConn.Close()
						continue
					}

					if !strings.Contains(string(buffer[:n]), "@ERROR") {
						authConn.Close()
						select {
						case <-ctx.Done():
						case resultChan <- struct {
							success    bool
							moduleName string
							err        error
						}{true, moduleName, nil}:
						}
						return
					}
				}
			}
			authConn.Close()
		}

		// 如果执行到这里，没有找到成功的认证
		select {
		case <-ctx.Done():
		case resultChan <- struct {
			success    bool
			moduleName string
			err        error
		}{false, "", fmt.Errorf("认证失败或无可用模块")}:
		}
	}()

	// 等待结果或上下文取消
	select {
	case result := <-resultChan:
		return result.success, result.moduleName, result.err
	case <-ctx.Done():
		return false, "", ctx.Err()
	}
}
