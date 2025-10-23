package portscan

import (
	"context"
	"fmt"
	"net/url"
	"slack-wails/lib/gologger"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type crackFunc func(context.Context, context.Context, string, string, []string, []string)

var crackScanners = map[string]crackFunc{
	"ftp":        FtpScan,
	"ssh":        SshScan,
	"telnet":     TelnetScan,
	"smb":        SmbScan,
	"oracle":     OracleScan,
	"mssql":      MssqlScan,
	"mysql":      MysqlScan,
	"rdp":        RdpScan,
	"postgresql": PostgresScan,
	"mongodb":    MongodbScan,
	"ldap":       LdapScan,
	"mqtt":       MqttScan,
	"socks5":     Socks5Scan,
	"vnc":        VncScan,
	"redis":      RedisScan,
	"memcached":  MemcachedScan,
	"jdwp":       JdwpScan,
	"adb":        AdbScan,
	"java-rmi":   RmiScan,
	"activemq":   ActiveMQScan,
	"rsync":      RsyncScan,
	"kafka":      KafkaScan,
}

func Runner(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	u, err := url.Parse(host)
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[!] Parse url error: %s\n", err))
		return
	}
	if scanFunc, ok := crackScanners[u.Scheme]; ok {
		scanFunc(ctx, ctrlCtx, taskId, u.Host, usernames, passwords)
	} else {
		gologger.Error(ctx, fmt.Sprintf("[!] No brute module registered for: %s\n", u.Scheme))
	}
	// 额外漏洞扫描
	switch u.Scheme {
	case "smb":
		MS17010(ctx, taskId, u.Host)
	}
	runtime.EventsEmit(ctx, fmt.Sprintf("crackDone::%s", host))
}
