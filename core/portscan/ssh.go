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
	"golang.org/x/crypto/ssh"
)

type task struct {
	User string
	Pass string
}

func SshScan(ctx, ctrlCtx context.Context, host string, usernames, passwords []string) {
	taskChan := make(chan task)
	doneChan := make(chan struct{})

	// 启动多线程暴破
	for range 5 {
		go func() {
			for t := range taskChan {
				// 检查是否需要停止
				if ctrlCtx.Err() != nil {
					continue
				}

				pass := strings.Replace(t.Pass, "{user}", t.User, -1)
				flag, err := SshConn(host, t.User, pass)
				if flag && err == nil {
					result, err := ExecuteSshCommand(host, t.User, pass, "whoami")
					if err != nil {
						result = err.Error()
					}
					runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
						ID:       "ssh weak password",
						Name:     "ssh weak password",
						URL:      host,
						Type:     "SSH",
						Severity: "HIGH",
						Extract:  t.User + "/" + pass,
						Request:  "[Command] whoami",
						Response: result,
					})
				} else {
					gologger.Info(ctx, fmt.Sprintf("ssh://%s %s:%s login failed", host, t.User, pass))
				}
			}
			doneChan <- struct{}{}
		}()
	}

	// 分发任务
	go func() {
		for _, user := range usernames {
			for _, pass := range passwords {
				// 如果中途终止，就不要继续塞任务了
				if ctrlCtx.Err() != nil {
					break
				}
				taskChan <- task{User: user, Pass: pass}
			}
		}
		close(taskChan)
	}()

	// 等待所有worker退出
	for range 5 {
		<-doneChan
	}
}

func SshConn(host, user, pass string) (bool, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return false, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	return true, nil
}

func ExecuteSshCommand(host, username, password, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 5 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	// Create a session to execute the command
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Run the command and capture the output
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}

	return string(output), nil
}
