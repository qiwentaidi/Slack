package portscan

import (
	"context"
	"fmt"
	"net"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/gonmap"
	"slack-wails/lib/structs"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var ExitFunc = false

func TcpScan(ctx context.Context, addresses <-chan Address, workers, timeout int, proxy clients.Proxy) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan *structs.InfoResult)
	var wg sync.WaitGroup
	openPorts := make(map[string]bool) // 记录开放的端口
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "webFingerScan", pr)
		}
		close(single)
		runtime.EventsEmit(ctx, "scanComplete", "done")
	}()
	// port scan func
	portScan := func(add Address) {
		if ExitFunc {
			return
		}
		pr := Connect(add.IP, add.Port, timeout, proxy)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(ctx, "progressID", id)
		if pr == nil {
			return
		}
		// 检查1-10端口的开放情况
		if pr.Port >= 1 && pr.Port <= 20 {
			openPorts[pr.Host] = true // 记录该IP有开放端口
			gologger.IntervalError(ctx, fmt.Sprintf("[portscan] %s 疑似cdn地址，会对未识别到服务的端口进行过滤", pr.Host))
		} else if openPorts[pr.Host] && pr.Scheme == "" {
			// 如果该IP在1-10端口有开放，后续端口必须识别到服务
			return // 如果没有识别到服务，则不返回
		}
		retChan <- pr
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

func Connect(ip string, port, timeout int, proxy clients.Proxy) *structs.InfoResult {
	scanner := gonmap.New()
	status, response := scanner.Scan(ip, port, time.Second*time.Duration(timeout), proxy)
	if status == gonmap.Closed || status == gonmap.Unknown {
		return nil
	}
	if response == nil || response.FingerPrint.Service == "" {
		return &structs.InfoResult{
			Host:   ip,
			Port:   port,
			Scheme: "unknow",
			URL:    fmt.Sprintf("%v://%v:%v", "unknow", ip, port),
		}
	}

	if response.FingerPrint.Service == "http" || response.FingerPrint.Service == "https" {
		result := &structs.InfoResult{
			Host:   ip,
			Port:   port,
			Scheme: response.FingerPrint.Service,
			URL:    fmt.Sprintf("%v://%v:%v", response.FingerPrint.Service, ip, port),
		}
		resp, _, err := clients.NewSimpleGetRequest(result.URL, clients.DefaultClient())
		if err == nil && resp != nil {
			result.StatusCode = resp.StatusCode
		}
		return result
	}

	return &structs.InfoResult{
		Host:   ip,
		Port:   port,
		Scheme: response.FingerPrint.Service,
		URL:    fmt.Sprintf("%v://%v:%v", response.FingerPrint.Service, ip, port),
	}
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
