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

func MemcachedScan(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	client, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(10 * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
							TaskId:   taskId,
							ID:       "memcached unauthorized",
							Name:     "memcached unauthorized",
							URL:      host,
							Type:     "Memcached",
							Severity: "HIGH",
							Request:  "stats\n",
							Response: string(rev[:n]),
						})
						return
					}
				} else {
					gologger.Info(ctx, fmt.Sprintf("memcached://%s is no unauthorized access", host))
				}
			}
		}
	}
}
