package portscan

import (
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strings"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
)

func RedisScan(address string, associate bool, passtext *widget.Entry) {
	var counter int64
	var passwords []string
	custom.LogTime = time.Now().Unix()
	flag, err := RedisUnauth(address)
	if flag && err == nil {
		custom.Console.Append(fmt.Sprintf("[+] redis %v %v\n", address, "unauthorized"))
		return
	}
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	for _, pass := range passwords {
		pass = strings.Replace(pass, "{user}", "redis", -1)
		flag, err := RedisConn(address, pass)
		if flag && err == nil {
			custom.Console.Append(fmt.Sprintf("[+] redis %v %v\n", address, pass))
			return
		} else {
			atomic.AddInt64(&counter, 1)
			custom.LogProgress(counter, len(passwords), fmt.Sprintf("[-] redis %v %v %v\n", address, pass, err))
			if time.Now().Unix()-custom.LogTime > int64(len(passwords)*common.Profile.PortScan.Timeout) {
				return
			}
		}
	}
}

func RedisConn(address, password string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(common.Profile.PortScan.Timeout) * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", password)))
	if err != nil {
		return flag, err
	}
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return flag, err
	}
	if strings.Contains(string(buffer[:n]), "+OK") {
		flag = true
		common.PortBurstResult = append(common.PortBurstResult, []string{"Redis", address, "", password})
	}
	return flag, err
}

func RedisUnauth(address string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(time.Duration(common.Profile.PortScan.Timeout) * time.Second))
	if err != nil {
		return flag, err
	}
	_, err = conn.Write([]byte("info\r\n"))
	if err != nil {
		return flag, err
	}
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return flag, err
	}
	if strings.Contains(string(buffer[:n]), "redis_version") {
		flag = true
		common.PortBurstResult = append(common.PortBurstResult, []string{"redis", address, "", "unauthorized"})
	}
	return flag, err
}
