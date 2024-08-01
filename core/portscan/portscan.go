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

func TcpScan(ctx context.Context, addresses <-chan Address, workers, timeout int) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan PortResult)
	var wg sync.WaitGroup
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "portScanLoading", pr)
		}
		close(single)
		runtime.EventsEmit(ctx, "scanComplete", "done")
	}()
	// port scan func
	portScan := func(add Address) {
		if ExitFunc {
			return
		}
		pr := Connect(add.IP, add.Port, timeout)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(ctx, "progressID", id)
		if pr.Status {
			pr.IP = add.IP
			pr.Port = add.Port
			retChan <- pr
		}
		// gologger.Info(ctx, pr)
	}
	threadPool, _ := ants.NewPoolWithFunc(workers, func(ipaddr interface{}) {
		ipa := ipaddr.(Address)
		portScan(ipa)
		wg.Done()
	})
	defer threadPool.Release()
	for add := range addresses {
		if ExitFunc {
			return
		}
		wg.Add(1)
		threadPool.Invoke(add)
	}
	wg.Wait()
	close(retChan)
	<-single
}

type Address struct {
	IP   string
	Port int
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
			if resp, b, err := clients.NewSimpleGetRequest(pr.Link, clients.DefaultClient()); err == nil {
				// 过滤云防护
				if resp.StatusCode == 422 {
					pr.Status = false
				}
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
