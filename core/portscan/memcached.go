package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MemcachedScan(ctx context.Context, host string) {
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
						runtime.EventsEmit(ctx, "bruteResult", Burte{
							Status:   true,
							Host:     host,
							Protocol: "memcached",
							Username: "unauthorized",
							Password: "",
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
