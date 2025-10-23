package portscan

import (
	"context"
	"fmt"
	"net"
	"slack-wails/core/webscan"
	"slack-wails/lib/structs"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/qiwentaidi/gonmap"

	"github.com/panjf2000/ants/v2"
	"github.com/qiwentaidi/clients"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func TcpScan(ctx, ctrlCtx context.Context, taskId string, addresses <-chan Address, workers, timeout int, proxyURL string) {
	var id int32
	single := make(chan struct{})
	retChan := make(chan *structs.InfoResult)
	var wg sync.WaitGroup
	// openPorts := make(map[string]bool) // 记录开放的端口
	go func() {
		for pr := range retChan {
			runtime.EventsEmit(ctx, "webFingerScan", pr)
		}
		close(single)
	}()
	var (
		portCount     = make(map[string]int)
		servicePorts  = make(map[string][]int)
		openPortMutex sync.Mutex
	)

	portScan := func(add Address) {
		defer wg.Done()
		defer func() {
			atomic.AddInt32(&id, 1)
			runtime.EventsEmit(ctx, "progressID", id)
		}()
		if ctrlCtx.Err() != nil {
			return
		}
		pr := Connect(ctx, taskId, add.IP, add.Port, timeout, proxyURL)
		if pr == nil {
			return
		}

		ip := pr.Host
		port := pr.Port
		scheme := pr.Scheme

		openPortMutex.Lock()
		portCount[ip]++
		if scheme != "unknown" {
			servicePorts[ip] = append(servicePorts[ip], port)
			openPortMutex.Unlock()
			retChan <- pr
		} else {
			totalOpen := portCount[ip]
			openPortMutex.Unlock()
			// 超过30个unknown服务直接忽略不
			if totalOpen > 30 {
				// gologger.Debug(ctx, fmt.Sprintf("[FILTER] %s:%d 忽略无服务端口（疑似全端口开放）", ip, port))
				return
			}
			retChan <- pr
		}
	}
	threadPool, _ := ants.NewPoolWithFunc(workers, func(ipaddr interface{}) {
		ipa := ipaddr.(Address)
		portScan(ipa)
	})
	defer threadPool.Release()
	for add := range addresses {
		if ctrlCtx.Err() != nil {
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

func Connect(ctx context.Context, taskId, ip string, port, timeout int, proxyURL string) *structs.InfoResult {
	scanner := gonmap.New()
	status, response := scanner.Scan(ip, port, time.Second*time.Duration(timeout), proxyURL)

	// 端口关闭或未知，直接返回 nil
	if status == gonmap.Closed || status == gonmap.Unknown {
		return nil
	}
	var tcpfinger []string
	var raw string // 添加一个默认值
	// 默认协议设为 unknown
	scheme := "unknown"
	if response != nil && response.FingerPrint.Service != "" {
		scheme = response.FingerPrint.Service
		raw = response.Raw
	}
	tcpinfo := &webscan.WebInfo{
		Protocol: scheme,
		Banner:   strings.ToLower(raw),
	}
	if scheme != "http" && scheme != "https" {
		tcpfinger = webscan.Scan(ctx, tcpinfo, webscan.FingerprintDB)
	}
	// if scheme == "unknown" {
	// 	scheme = gonmap.GuessProtocol(port)
	// }
	result := &structs.InfoResult{
		TaskId:       taskId,
		Host:         ip,
		Port:         port,
		Scheme:       scheme,
		URL:          fmt.Sprintf("%s://%s:%d", scheme, ip, port),
		Fingerprints: tcpfinger,
		Detect:       "Default",
	}

	// 若是 HTTP/HTTPS，尝试请求获取状态码
	if scheme == "http" || scheme == "https" {
		resp, err := clients.SimpleGet(result.URL, clients.NewRestyClient(nil, true))
		if err == nil {
			result.StatusCode = resp.StatusCode()
		}
	}

	return result
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
