package portscan

import (
	"context"
	"slack-wails/lib/structs"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 只要是nmap 扫描到jdwp协议，默认是 unauthorized (因为也是同样发JDWP-Handshake包检测)
func JdwpScan(ctx context.Context, address string) {
	// client, err := WrapperTcpWithTimeout("tcp", address, time.Duration(6)*time.Second)
	// defer func() {
	// 	if client != nil {
	// 		client.Close()
	// 	}
	// }()
	// if err != nil {
	// 	gologger.Info(ctx, fmt.Sprintf("connect %s failed", address))
	// 	return
	// }
	// err = client.SetDeadline(time.Now().Add(time.Duration(6) * time.Second))
	// if err != nil {
	// 	gologger.Info(ctx, fmt.Sprintf("connect %s failed", address))
	// 	return
	// }
	// _, err = client.Write([]byte("JDWP-Handshake"))
	// if err != nil {
	// 	gologger.Info(ctx, fmt.Sprintf("write jdwp-handshake to %s failed", address))
	// 	return
	// }

	// rev := make([]byte, 1024)
	// n, errRead := client.Read(rev)
	// if errRead != nil {
	// 	gologger.Info(ctx, fmt.Sprintf("read %s err: %s", address, errRead))
	// 	return
	// }
	// if !strings.Contains(string(rev[:n]), "JDWP-Handshake") {
	// 	// 不是JDWP
	// 	gologger.Info(ctx, fmt.Sprintf("%s is not jdwp", address))
	// 	return
	// }
	runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
		ID:          "jdwp unauthorized",
		Name:        "jdwp unauthorized",
		URL:         address,
		Type:        "JDWP",
		Severity:    "critical",
		Request:     "JDWP-Handshake",
		Response:    "JDWP-Handshake",
		Description: "检测到JDWP端口开放，自行尝试是否存在命令执行漏洞",
		Reference:   "https://forum.butian.net/share/1232,https://github.com/l3yx/jdwp-codeifier",
	})
}
