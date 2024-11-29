package portscan

import (
	"context"
	"net/url"
	"strings"
)

type Burte struct {
	Status   bool
	Protocol string
	Host     string
	Username string
	Password string
}

var ExitBruteFunc = false

var defaultPorts = map[string]string{
	"ftp":        "21",
	"ssh":        "22",
	"telnet":     "23",
	"ldap":       "389",
	"smb":        "445", // SMB 通常使用445端口
	"socks5":     "1080",
	"rmi":        "1099",
	"oracle":     "1521",
	"mssql":      "1433",
	"mqtt":       "1883",
	"mysql":      "3306",
	"rdp":        "3389",
	"postgresql": "5432",
	"adb":        "5555",
	"vnc":        "5900",
	"redis":      "6379",
	"jdwp":       "8000",
	"memcached":  "11211",
	"mongodb":    "27017",
}

// 检查并为给定的主机添加默认端口号
func addDefaultPort(scheme, host string) string {
	if strings.Contains(host, ":") {
		return host
	}
	defaultPort := defaultPorts[scheme]
	return host + ":" + defaultPort
}

func PortBrute(ctx context.Context, host string, usernames, passwords []string) {
	u, err := url.Parse(host)
	if err != nil {
		return
	}
	u.Host = addDefaultPort(u.Scheme, u.Host)
	switch u.Scheme {
	case "ftp":
		FtpScan(ctx, u.Host, usernames, passwords)
	case "ssh":
		SshScan(ctx, u.Host, usernames, passwords)
	case "telnet":
		TelenetScan(ctx, u.Host, usernames, passwords)
	case "smb":
		SmbScan(ctx, u.Host, usernames, passwords)
	case "oracle":
		OracleScan(ctx, u.Host, usernames, passwords)
	case "mssql":
		MssqlScan(ctx, u.Host, usernames, passwords)
	case "mysql":
		MysqlScan(ctx, u.Host, usernames, passwords)
	case "rdp":
		RdpScan(ctx, u.Host, usernames, passwords)
	case "postgresql":
		PostgresScan(ctx, u.Host, usernames, passwords)
	case "vnc":
		VncScan(ctx, u.Host, passwords)
	case "redis":
		RedisScan(ctx, u.Host, passwords)
	case "memcached":
		MemcachedScan(ctx, u.Host)
	case "mongodb":
		MongodbScan(ctx, u.Host, usernames, passwords)
	case "ldap":
		LdapScan(ctx, u.Host, usernames, passwords)
	case "mqtt":
		MqttScan(ctx, u.Host, usernames, passwords)
	case "socks5":
		Socks5Scan(ctx, u.Host, usernames, passwords)
	case "jdwp":
		JdwpScan(ctx, u.Host)
	case "adb":
		AdbScan(ctx, u.Host)
	case "rmi":
		RmiScan(ctx, u.Host)
	}
}
