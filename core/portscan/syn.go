package portscan

import (
	"context"
	"net"
	"slack-wails/lib/gologger"
	"sync"
	"sync/atomic"

	"github.com/XinRoom/go-portScan/core/port"
	"github.com/XinRoom/go-portScan/core/port/syn"
	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const synTimeout = 12

type Address2 struct {
	IP   net.IP
	Port uint16
}

func SynScan(ctx context.Context, startIp net.IP, address <-chan Address2) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan port.OpenIpPort)
	go func() {
		for ret := range retChan {
			pr := Connect(ret.Ip.To4().String(), int(ret.Port), synTimeout)
			pr.IP = ret.Ip.To4().String()
			pr.Port = int(ret.Port)
			runtime.EventsEmit(ctx, "portScanLoading", pr)
			atomic.AddInt32(&id, 1)
			runtime.EventsEmit(ctx, "progressID", id)
		}
		single <- struct{}{}
		runtime.EventsEmit(ctx, "scanComplete", "done")
	}()
	// scanner
	ss, err := syn.NewSynScanner(startIp, retChan, syn.DefaultSynOption)
	if err != nil {
		gologger.Error(ctx, "Permission denied, please run with sudo")
		return
	}
	// port scan func
	portScan := func(addr Address2) {
		ss.WaitLimiter()
		ss.Scan(addr.IP, addr.Port) // syn 不能并发，默认以网卡和驱动最高性能发包
	}
	// Pool - ping and port scan
	var wgPing sync.WaitGroup
	poolPing, _ := ants.NewPoolWithFunc(50, func(addr interface{}) {
		add := addr.(Address2)
		portScan(add)
		wgPing.Done()
	})
	defer poolPing.Release()
	for addr := range address {
		if ExitFunc {
			return
		}
		wgPing.Add(1)
		poolPing.Invoke(addr)
	}
	wgPing.Wait()
	ss.Wait()
	ss.Close()
	<-single
}
