package portscan

import (
	"fmt"
	"net"
	"slack-wails/lib/clients"
	"slack-wails/lib/gonmap"
	"slack-wails/lib/util"
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

type PortResult struct {
	Status    bool
	IP        string
	Port      int
	Server    string
	Link      string
	HttpTitle string
}

func Connect(ip string, port, timeout int) PortResult {
	var pr PortResult
	scanner := gonmap.New()
	status, response := scanner.Scan(ip, port, time.Second*time.Duration(timeout))
	switch status {
	case gonmap.Closed:
		pr.Status = false
	// filter 未知状态
	case gonmap.Unknown:
		pr.Status = true
		pr.Server = "filter"
	default:
		pr.Status = true
	}
	if response != nil {
		if response.FingerPrint.Service != "" {
			pr.Server = response.FingerPrint.Service
		} else {
			pr.Server = "unknow"
		}
		pr.Link = fmt.Sprintf("%v://%v:%v", pr.Server, ip, port)
		if pr.Server == "http" || pr.Server == "https" {
			if _, b, err := clients.NewSimpleGetRequest(pr.Link, clients.DefaultClient()); err == nil {
				if match := util.RegTitle.FindSubmatch(b); len(match) > 1 {
					pr.HttpTitle = util.Str2UTF8(string(match[1]))
				} else {
					pr.HttpTitle = "-"
				}
			}
		}
	}
	return pr
}
