package portscan

import (
	"context"
	"fmt"
	"net"
	"slack-wails/lib/gologger"
	"sync"
	"sync/atomic"
	"time"

	"github.com/XinRoom/go-portScan/core/host"
	"github.com/XinRoom/go-portScan/core/port"
	"github.com/XinRoom/go-portScan/core/port/syn"
	"github.com/XinRoom/iprange"
	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const synTimeout = 5

func SynScan(ctx context.Context, ip string, ports []uint16, id *int32) {
	single := make(chan struct{})
	retChan := make(chan port.OpenIpPort, 65535)
	go func() {
		for ret := range retChan {
			result := PortResult{
				Status: true,
				IP:     ret.Ip.String(),
				Port:   int(ret.Port),
			}
			gologger.Success(ctx, fmt.Sprintf("[syn] %v is open ! ", ret))
			pr := Connect(ret.Ip.String(), int(ret.Port), synTimeout, nil)
			if pr.Status {
				result.HttpTitle = pr.HttpTitle
				result.Server = pr.Server
				result.Link = pr.Link
			}
			runtime.EventsEmit(ctx, "portScanLoading", result)
		}
		single <- struct{}{}
	}()

	// parse ip
	it, startIp, _ := iprange.NewIter(ip)

	// scanner
	ss, err := syn.NewSynScanner(startIp, retChan, syn.DefaultSynOption)
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("[syn] %s", err))
		// gologger.Error(ctx, "[syn] Permission denied, please run with sudo")
	}

	// port scan func
	portScan := func(ip net.IP) {
		for _, _port := range ports { // port
			ss.WaitLimiter()
			ss.Scan(ip, _port) // syn 不能并发，默认以网卡和驱动最高性能发包
			atomic.AddInt32(id, 1)
			runtime.EventsEmit(ctx, "progressID", id)
		}
	}

	// Pool - ping and port scan
	var wgPing sync.WaitGroup
	poolPing, _ := ants.NewPoolWithFunc(50, func(ip interface{}) {
		_ip := ip.(net.IP)
		if host.IsLive(_ip.String(), true, 800*time.Millisecond) {
			portScan(_ip)
		}
		wgPing.Done()
	})
	defer poolPing.Release()
	for i := uint64(0); i < it.TotalNum(); i++ { // ip索引
		ip := make(net.IP, len(it.GetIpByIndex(0)))
		copy(ip, it.GetIpByIndex(i)) // Note: dup copy []byte when concurrent (GetIpByIndex not to do dup copy)
		wgPing.Add(1)
		poolPing.Invoke(ip)
	}
	wgPing.Wait()
	ss.Wait()
	ss.Close()
	<-single
}
