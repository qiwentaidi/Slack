package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gonmap"
	"slack-wails/lib/util"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var ExitFunc = false

type PortResult struct {
	Status    bool
	IP        string
	Port      int
	Server    string
	Link      string
	HttpTitle string
}

func TcpScan(ctx context.Context, ips []string, ports []int, workers, timeout int) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan PortResult, len(ips)*len(ports))
	var wg sync.WaitGroup
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "tcpPortScanLoading", pr)
		}
		close(single)
		runtime.EventsEmit(ctx, "tcpScanComplete", "done")
	}()
	// port scan func
	portScan := func(port int) {
		for _, ip := range ips {
			if ExitFunc {
				return
			}
			pr := Connect(ip, port, timeout)
			runtime.EventsEmit(ctx, "tcpProgressID", id)
			atomic.AddInt32(&id, 1)
			if pr.Status {
				pr.IP = ip
				pr.Port = port
				retChan <- pr
			}
			// gologger.Info(ctx, pr)
		}
	}
	threadPool, _ := ants.NewPoolWithFunc(workers, func(ports interface{}) {
		port := ports.(int)
		portScan(port)
		wg.Done()
	})
	defer threadPool.Release()
	for _, port := range ports {
		if ExitFunc {
			return
		}
		wg.Add(1)
		threadPool.Invoke(port)
	}
	wg.Wait()
	close(retChan)
	<-single
}

type Address struct {
	IP   string
	Port int
}

// 处理 192.168.1.1:6379 这种单独IP端口组模式
func CorrespondsScan(ctx context.Context, address []Address, timeout int) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan PortResult, len(address))
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "csPortScanLoading", pr)
		}
		close(single)
	}()
	// port scan func
	portScan := func(addr Address) {
		pr := Connect(addr.IP, addr.Port, timeout)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(ctx, "csProgressID", id)
		if pr.Status {
			pr.IP = addr.IP
			pr.Port = addr.Port
			retChan <- pr
		}
	}
	var wg sync.WaitGroup
	threadPool, _ := ants.NewPoolWithFunc(20, func(task interface{}) {
		addr := task.(Address)
		portScan(addr)
		wg.Done()
	})
	defer threadPool.Release()
	for _, addr := range address {
		if ExitFunc {
			return
		}
		wg.Add(1)
		threadPool.Invoke(addr)
	}
	wg.Wait()
	close(retChan)
	<-single
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
