package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func RmiScan(ctx context.Context, host string) {
	// 使用10秒超时建立TCP连接
	client, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	if err != nil {
		err = client.SetDeadline(time.Now().Add(10 * time.Second))
		if err == nil {
			// RMI协议握手包
			handshake := []byte{
				0x4a, 0x52, 0x4d, 0x49, // "JRMI"
				0x00, 0x02, // 协议版本
				0x4b, // 握手标识
			}

			_, err = client.Write(handshake)
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					// 检查返回的数据是否包含RMI响应特征
					if n > 4 && string(rev[:4]) == "JRMI" {
						runtime.EventsEmit(ctx, "bruteResult", Burte{
							Status:   true,
							Host:     host,
							Protocol: "rmi",
							Username: "unauthorized",
							Password: "",
						})
						return
					}
				}
			}
		}
	}

	gologger.Info(ctx, fmt.Sprintf("rmi://%s is no unauthorized access", host))
}
