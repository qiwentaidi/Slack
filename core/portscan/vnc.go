package portscan

import (
	"strings"
	"time"

	"github.com/mitchellh/go-vnc"
)

func VncScan(host string, passwords []string) *Burte {
	for _, pass := range passwords {
		pass = strings.Replace(pass, "{user}", "vnc", -1)
		flag, err := VncConn(host, pass)
		if flag && err == nil {
			return &Burte{
				Status:   true,
				Host:     host,
				Protocol: "vnc",
				Username: "",
				Password: pass,
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "vnc",
		Username: "",
		Password: "",
	}
}

func VncConn(host, pass string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
	if err != nil {
		return flag, err
	}
	// Create a VNC client connection
	cfg := &vnc.ClientConfig{
		Auth: []vnc.ClientAuth{&vnc.PasswordAuth{Password: pass}},
	}
	client, err := vnc.Client(conn, cfg)
	if err != nil {
		return flag, err
	}
	// Use the client connection to interact with the VNC server
	// ...

	// Close the client connection when done
	client.Close()
	flag = true
	return flag, err
}
