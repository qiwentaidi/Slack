package portscan

import (
	"context"
	"net"
	"net/url"
	"strings"
	"time"
)

type Burte struct {
	Status   bool
	Protocol string
	Host     string
	Username string
	Password string
}

var ExitBruteFunc = false

func WrapperTcpWithTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := &net.Dialer{Timeout: timeout}
	return WrapperTCP(network, address, d)
}

func WrapperTCP(network, address string, forward *net.Dialer) (net.Conn, error) {
	//get conn
	conn, err := forward.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var DefaultPorts = map[string]string{
	"ftp":        "21",
	"ssh":        "22",
	"telnet":     "23",
	"smb":        "445", // SMB 通常使用445端口
	"oracle":     "1521",
	"mssql":      "1433",
	"mysql":      "3306",
	"rdp":        "3389",
	"postgresql": "5432",
	"vnc":        "5900",
	"redis":      "6379",
	"memcached":  "11211",
	"mongodb":    "27017",
}

// AddDefaultPort 检查并为给定的主机添加默认端口号
func AddDefaultPort(scheme, host string) string {
	if strings.Contains(host, ":") {
		return host
	}
	defaultPort := DefaultPorts[scheme]
	return host + ":" + defaultPort
}

func PortBrute(ctx context.Context, host string, usernames, passwords []string) {
	u, err := url.Parse(host)
	if err != nil {
		return
	}
	u.Host = AddDefaultPort(u.Scheme, u.Host)
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
		MongodbScan(ctx, u.Host)
	}
}
