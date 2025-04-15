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

func AdbScan(ctx context.Context, address string, usernames, passwords []string) {
	result := "ADB> host::features=shell_v2,cmd,stat_v2,ls_v2,fixed_push_mkdir,apex,abb,fixed_push_symlink_timestamp,abb_exec,remount_shell,track_app,sendrecv_v2,sendrecv_v2_brotli,sendrecv_v2_lz4,sendrecv_v2_zstd,sendrecv_v2_dry_run_send,openscreen_mdns\n"
	conn, err := WrapperTcpWithTimeout("tcp", address, time.Duration(6)*time.Second)
	if err == nil {
		defer func() {
			if conn != nil {
				_ = conn.Close()
			}
		}()
	} else {
		gologger.Info(ctx, fmt.Sprintf("connect %s failed", address))
		return
	}

	_, err = conn.Write([]byte{0x43, 0x4e, 0x58, 0x4e, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x10, 0x00, 0xea, 0x00, 0x00, 0x00,
		0x44, 0x5b, 0x00, 0x00, 0xbc, 0xb1, 0xa7, 0xb1,
		0x68, 0x6f, 0x73, 0x74, 0x3a, 0x3a, 0x66, 0x65,
		0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x3d, 0x73,
		0x68, 0x65, 0x6c, 0x6c, 0x5f, 0x76, 0x32, 0x2c,
		0x63, 0x6d, 0x64, 0x2c, 0x73, 0x74, 0x61, 0x74,
		0x5f, 0x76, 0x32, 0x2c, 0x6c, 0x73, 0x5f, 0x76,
		0x32, 0x2c, 0x66, 0x69, 0x78, 0x65, 0x64, 0x5f,
		0x70, 0x75, 0x73, 0x68, 0x5f, 0x6d, 0x6b, 0x64,
		0x69, 0x72, 0x2c, 0x61, 0x70, 0x65, 0x78, 0x2c,
		0x61, 0x62, 0x62, 0x2c, 0x66, 0x69, 0x78, 0x65,
		0x64, 0x5f, 0x70, 0x75, 0x73, 0x68, 0x5f, 0x73,
		0x79, 0x6d, 0x6c, 0x69, 0x6e, 0x6b, 0x5f, 0x74,
		0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
		0x2c, 0x61, 0x62, 0x62, 0x5f, 0x65, 0x78, 0x65,
		0x63, 0x2c, 0x72, 0x65, 0x6d, 0x6f, 0x75, 0x6e,
		0x74, 0x5f, 0x73, 0x68, 0x65, 0x6c, 0x6c, 0x2c,
		0x74, 0x72, 0x61, 0x63, 0x6b, 0x5f, 0x61, 0x70,
		0x70, 0x2c, 0x73, 0x65, 0x6e, 0x64, 0x72, 0x65,
		0x63, 0x76, 0x5f, 0x76, 0x32, 0x2c, 0x73, 0x65,
		0x6e, 0x64, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x76,
		0x32, 0x5f, 0x62, 0x72, 0x6f, 0x74, 0x6c, 0x69,
		0x2c, 0x73, 0x65, 0x6e, 0x64, 0x72, 0x65, 0x63,
		0x76, 0x5f, 0x76, 0x32, 0x5f, 0x6c, 0x7a, 0x34,
		0x2c, 0x73, 0x65, 0x6e, 0x64, 0x72, 0x65, 0x63,
		0x76, 0x5f, 0x76, 0x32, 0x5f, 0x7a, 0x73, 0x74,
		0x64, 0x2c, 0x73, 0x65, 0x6e, 0x64, 0x72, 0x65,
		0x63, 0x76, 0x5f, 0x76, 0x32, 0x5f, 0x64, 0x72,
		0x79, 0x5f, 0x72, 0x75, 0x6e, 0x5f, 0x73, 0x65,
		0x6e, 0x64, 0x2c, 0x6f, 0x70, 0x65, 0x6e, 0x73,
		0x63, 0x72, 0x65, 0x65, 0x6e, 0x5f, 0x6d, 0x64,
		0x6e, 0x73})
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("write %s failed", address))
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(6) * time.Second))
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("set read deadline for %s failed", address))
		return
	}

	buf := make([]byte, 0x1000)
	n, err := conn.Read(buf)
	if err != nil {
		gologger.Info(ctx, fmt.Sprintf("read from %s failed", address))
		return
	}

	if n > 4 && string(buf[:4]) != "CNXN" {
		gologger.Info(ctx, "ADB需要授权/非ADB服务")
		return
	}

	if strings.Contains(string(buf[:n]), "ro.product.name") {
		result += string(buf[24:n]) + "\n"
	} else {
		buf = make([]byte, 0x1000)
		n, err = conn.Read(buf)
		if err != nil {
			gologger.Info(ctx, fmt.Sprintf("read from %s failed", address))
			return
		}

		result += string(buf[:n]) + "\n"
	}

	if result != "" {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			ID:       "adb unauthorized",
			Name:     "adb unauthorized",
			URL:      address,
			Type:     "ADB",
			Severity: "HIGH",
		})
		return
	}
}
