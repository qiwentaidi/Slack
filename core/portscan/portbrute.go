package portscan

import (
	"context"
	"net"
	"net/url"
	"time"
)

type Burte struct {
	Status   bool
	Protocol string
	Host     string
	Username string
	Password string
}

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

func PortBrute(ctx context.Context, host string, usernames, passwords []string) {
	u, err := url.Parse(host)
	if err != nil {
		return
	}
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
