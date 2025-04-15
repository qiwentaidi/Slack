package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/mitchellh/go-vnc"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func VncScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, pass := range passwords {
		if ExitFunc {
			return
		}
		pass = strings.Replace(pass, "{user}", "vnc", -1)
		flag, err := VncConn(host, pass)
		if flag && err == nil {
			runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
				ID:       "vnc weak password",
				Name:     "vnc weak password",
				URL:      host,
				Type:     "VNC",
				Severity: "HIGH",
				Extract:  pass,
			})
			return
		} else {
			gologger.Info(ctx, fmt.Sprintf("vnc://%s %s is login failed", host, pass))
		}
	}
}

func VncConn(host, pass string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return flag, err
	}
	// Create a VNC client connection
	cfg := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{&vnc.PasswordAuth{Password: pass}},
	}
	client, err := vnc.Client(conn, cfg)
	if err != nil {
		return flag, err
	}
	// Use the client connection to interact with the VNC server
	// ...

	// Close the client connection when done
	client.Close()
	flag = true
	return flag, err
}
