package portscan

import (
	"context"
	"fmt"
	"net/url"
	"slack-wails/lib/gologger"
)

type crackFunc func(context.Context, string, []string, []string)

var crackScanners = map[string]crackFunc{
	"ftp":        FtpScan,
	"ssh":        SshScan,
	"telnet":     TelenetScan,
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
	"rmi":        RmiScan,
}

func Runner(ctx context.Context, host string, usernames, passwords []string) {
	u, err := url.Parse(host)
	if err != nil {
		gologger.Debug(ctx, fmt.Sprintf("[!] Parse url error: %s\n", err))
		return
	}
	if scanFunc, ok := crackScanners[u.Scheme]; ok {
		scanFunc(ctx, u.Host, usernames, passwords)
	} else {
		gologger.Error(ctx, fmt.Sprintf("[!] No brute module registered for: %s\n", u.Scheme))
	}
	// 额外漏洞扫描
	switch u.Scheme {
	case "smb":
		MS17010(ctx, u.Host)
	}
}
