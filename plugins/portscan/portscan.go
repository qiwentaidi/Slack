package portscan

import (
	"fmt"
	"net"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"slack/lib/gonmap"
)

type Addr struct {
	ip   string
	port int
}

var (
	id    uint32
	count int
)

// 进行TCP扫描
func PortScanTCP(IPs []string, Ports []int, fe *custom.ForwordEntry) {
	var (
		workers int
		wg      sync.WaitGroup
	)
	count = len(IPs) * len(Ports)
	global.PortscanProgress.SetText(fmt.Sprintf("0/%v", count))
	Addrs := make(chan Addr, count)
	results := make(chan string, count)
	if common.Profile.PortScan.Thread > int(count) {
		workers = int(count)
	} else {
		workers = common.Profile.PortScan.Thread
	}
	//接收结果
	go func() {
		for range results {
			wg.Done()
		}
	}()
	id = 0
	//多线程扫描
	for i := 0; i < workers; i++ {
		go func() {
			for addr := range Addrs {
				PortConnect(addr, results, &wg, fe)
				wg.Done()
			}
		}()
	}

	//添加扫描目标
	for _, port := range Ports {
		for _, host := range IPs {
			wg.Add(1)
			Addrs <- Addr{host, port}
		}
	}
	wg.Wait()
	global.PortscanProgress.SetText("端口扫描任务结束")
	close(Addrs)
	close(results)
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

func PortConnect(addr Addr, respondingHosts chan<- string, wg *sync.WaitGroup, result *custom.ForwordEntry) {
	conn, err := WrapperTcpWithTimeout("tcp", fmt.Sprintf("%s:%v", addr.ip, addr.port), time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	atomic.AddUint32(&id, 1)
	global.PortscanProgress.SetText(fmt.Sprintf("%v/%v", id, count))
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil {
		PortCheck(addr.ip, addr.port, result)
		wg.Add(1)
		respondingHosts <- addr.ip + ":" + strconv.Itoa(addr.port)
	}
}

func PortCheck(ip string, port int, result *custom.ForwordEntry) {
	scanner := gonmap.New()
	_, response := scanner.ScanTimeout(ip, port, time.Second*time.Duration(common.Profile.PortScan.Timeout))
	if response != nil {
		if response.FingerPrint.Service != "" {
			result.Text += fmt.Sprintf("%v://%v:%v\n", response.FingerPrint.Service, ip, port)
		} else {
			result.Text += fmt.Sprintf("unkonw://%v:%v\n", ip, port)
		}
		result.Refresh()
	}
}
