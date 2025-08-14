package portscan

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"slack-wails/lib/gologger"
	"slack-wails/lib/utils/arrayutil"

	"strings"
	"sync"
	"time"

	"golang.org/x/net/icmp"
)

var (
	AliveHosts []string
	ExistHosts = make(map[string]struct{})
	livewg     sync.WaitGroup
)

// CheckLive 检测主机存活状态
func CheckLive(ctx context.Context, hostslist []string, Ping bool) []string {
	// 创建主机通道
	chanHosts := make(chan string, len(hostslist))

	// 处理存活主机
	go handleAliveHosts(ctx, chanHosts, hostslist, Ping)

	// 根据Ping参数选择检测方式
	if Ping {
		// 使用ping方式探测
		RunPing(hostslist, chanHosts)
	} else {
		probeWithICMP(ctx, hostslist, chanHosts)
	}

	// 等待所有检测完成
	livewg.Wait()
	close(chanHosts)

	// 输出存活统计信息
	return AliveHosts
}

func handleAliveHosts(ctx context.Context, chanHosts chan string, hostslist []string, isPing bool) {
	for ip := range chanHosts {
		if _, ok := ExistHosts[ip]; !ok && arrayutil.ArrayContains(ip, hostslist) {
			ExistHosts[ip] = struct{}{}
			AliveHosts = append(AliveHosts, ip)
			gologger.Info(ctx, fmt.Sprintf("%s is alive!", ip))
		}
		livewg.Done()
	}
}

// probeWithICMP 使用ICMP方式探测
func probeWithICMP(ctx context.Context, hostslist []string, chanHosts chan string) {
	// 尝试监听本地ICMP
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err == nil {
		RunIcmp1(hostslist, conn, chanHosts)
		return
	}

	gologger.Error(ctx, "icmp_listen_failed"+err.Error())
	gologger.Info(ctx, "trying_no_listen_icmp")

	// 尝试无监听ICMP探测
	conn2, err := net.DialTimeout("ip4:icmp", "127.0.0.1", 3*time.Second)
	if err == nil {
		defer conn2.Close()
		RunIcmp2(hostslist, chanHosts)
		return
	}

	// 降级使用ping探测
	RunPing(hostslist, chanHosts)
}

// RunIcmp1 使用ICMP批量探测主机存活(监听模式)
func RunIcmp1(hostslist []string, conn *icmp.PacketConn, chanHosts chan string) {
	endflag := false

	// 启动监听协程
	go func() {
		for {
			if endflag {
				return
			}
			// 接收ICMP响应
			msg := make([]byte, 100)
			_, sourceIP, _ := conn.ReadFrom(msg)
			if sourceIP != nil {
				livewg.Add(1)
				chanHosts <- sourceIP.String()
			}
		}
	}()

	// 发送ICMP请求
	for _, host := range hostslist {
		dst, _ := net.ResolveIPAddr("ip", host)
		IcmpByte := makemsg(host)
		conn.WriteTo(IcmpByte, dst)
	}

	// 等待响应
	start := time.Now()
	for {
		// 所有主机都已响应则退出
		if len(AliveHosts) == len(hostslist) {
			break
		}

		// 根据主机数量设置超时时间
		since := time.Since(start)
		wait := time.Second * 6
		if len(hostslist) <= 256 {
			wait = time.Second * 3
		}

		if since > wait {
			break
		}
	}

	endflag = true
	conn.Close()
}

// RunIcmp2 使用ICMP并发探测主机存活(无监听模式)
func RunIcmp2(hostslist []string, chanHosts chan string) {
	// 控制并发数
	num := 1000
	if len(hostslist) < num {
		num = len(hostslist)
	}

	var wg sync.WaitGroup
	limiter := make(chan struct{}, num)

	// 并发探测
	for _, host := range hostslist {
		wg.Add(1)
		limiter <- struct{}{}

		go func(host string) {
			defer func() {
				<-limiter
				wg.Done()
			}()

			if icmpalive(host) {
				livewg.Add(1)
				chanHosts <- host
			}
		}(host)
	}

	wg.Wait()
	close(limiter)
}

// icmpalive 检测主机ICMP是否存活
func icmpalive(host string) bool {
	startTime := time.Now()

	// 建立ICMP连接
	conn, err := net.DialTimeout("ip4:icmp", host, 6*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 设置超时时间
	if err := conn.SetDeadline(startTime.Add(6 * time.Second)); err != nil {
		return false
	}

	// 构造并发送ICMP请求
	msg := makemsg(host)
	if _, err := conn.Write(msg); err != nil {
		return false
	}

	// 接收ICMP响应
	receive := make([]byte, 60)
	if _, err := conn.Read(receive); err != nil {
		return false
	}

	return true
}

func RunPing(hostslist []string, chanHosts chan string) {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 50)
	for _, host := range hostslist {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if ExecCommandPing(host) {
				livewg.Add(1)
				chanHosts <- host
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}

// ExecCommandPing 执行系统Ping命令检测主机存活
func ExecCommandPing(ip string) bool {
	// 过滤黑名单字符
	forbiddenChars := []string{";", "&", "|", "`", "$", "\\", "'", "%", "\"", "\n"}
	for _, char := range forbiddenChars {
		if strings.Contains(ip, char) {
			return false
		}
	}

	var command *exec.Cmd
	// 根据操作系统选择不同的ping命令
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false")
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+ip+" && echo true || echo false")
	default: // linux
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+ip+" && echo true || echo false")
	}

	// 捕获命令输出
	var outinfo bytes.Buffer
	command.Stdout = &outinfo

	// 执行命令
	if err := command.Start(); err != nil {
		return false
	}

	if err := command.Wait(); err != nil {
		return false
	}

	// 分析输出结果
	output := outinfo.String()
	return strings.Contains(output, "true") && strings.Count(output, ip) > 2
}

func makemsg(host string) []byte {
	msg := make([]byte, 40)
	id0, id1 := genIdentifier(host)
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4], msg[5] = id0, id1
	msg[6], msg[7] = genSequence(1)
	check := checkSum(msg[0:40])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)
	return msg
}

func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = sum + (sum >> 16)
	answer := uint16(^sum)
	return answer
}

func genSequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genIdentifier(host string) (byte, byte) {
	return host[0], host[1]
}
