package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func JdwpScan(ctx context.Context, address string) {
	client, err := WrapperTcpWithTimeout("tcp", address, time.Duration(6)*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("connect %s failed", address))
		return
	}
	err = client.SetDeadline(time.Now().Add(time.Duration(6) * time.Second))
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("connect %s failed", address))
		return
	}
	_, err = client.Write([]byte("JDWP-Handshake"))
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("write jdwp-handshake to %s failed", address))
		return
	}

	rev := make([]byte, 1024)
	n, errRead := client.Read(rev)
	if errRead != nil {
		gologger.Info(ctx, fmt.Sprintf("read %s err: %s", address, errRead))
		return
	}
	if !strings.Contains(string(rev[:n]), "JDWP-Handshake") {
		// 不是JDWP
		gologger.Info(ctx, fmt.Sprintf("%s is not jdwp", address))
		return
	}
	runtime.EventsEmit(ctx, "bruteResult", Burte{
		Status:   true,
		Host:     address,
		Protocol: "jdwp",
		Username: "unauthorized",
		Password: "",
	})
}
