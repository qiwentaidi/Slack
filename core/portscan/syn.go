package portscan

import (
	"context"
	"net"
	"sync"

	"github.com/XinRoom/go-portScan/core/port"
	"github.com/XinRoom/go-portScan/core/port/syn"
	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func SynScan(ctx context.Context, ips []string, ports []uint16, timeout int) {
	var id int
	single := make(chan struct{})
	retChan := make(chan port.OpenIpPort, len(ports))
	go func() {
		for ret := range retChan {
			pr := Connect(ret.Ip.To4().String(), int(ret.Port), timeout)
			pr.IP = ret.Ip.To4().String()
			pr.Port = int(ret.Port)
			runtime.EventsEmit(ctx, "synPortScanLoading", pr)
		}
		single <- struct{}{}
	}()
	startIp := net.ParseIP(ips[0])
	// scanner
	ss, err := syn.NewSynScanner(startIp, retChan, syn.DefaultSynOption)
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	// port scan func
	portScan := func(ip net.IP) {
		for _, port := range ports { // port
			id++
			runtime.EventsEmit(ctx, "synProgressID", id)
			ss.WaitLimiter()
			ss.Scan(ip, port) // syn 不能并发，默认以网卡和驱动最高性能发包
		}
	}
	// Pool - ping and port scan
	var wgPing sync.WaitGroup
	poolPing, _ := ants.NewPoolWithFunc(50, func(ip interface{}) {
		_ip := ip.(net.IP)
		portScan(_ip)
		wgPing.Done()
	})
	defer poolPing.Release()
	for _, ip := range ips {
		ipNet := net.ParseIP(ip)
		wgPing.Add(1)
		poolPing.Invoke(ipNet)
	}
	wgPing.Wait()
	ss.Wait()
	ss.Close()
	<-single
}
