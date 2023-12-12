package portscan

import (
	"fmt"
	"strings"
	"time"
)

func RedisScan(host string, passwords []string) *Burte {
	flag, err := RedisUnauth(host)
	if flag && err == nil {
		return &Burte{
			Status:   true,
			Host:     host,
			Protocol: "redis",
			Username: "unauthorized",
			Password: "",
		}
	}
	for _, pass := range passwords {
		pass = strings.Replace(pass, "{user}", "redis", -1)
		flag, err := RedisConn(host, pass)
		if flag && err == nil {
			return &Burte{
				Status:   true,
				Host:     host,
				Protocol: "redis",
				Username: "",
				Password: pass,
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "redis",
		Username: "",
		Password: "",
	}
}

func RedisConn(address, password string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, 10*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
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
	}
	return flag, err
}

func RedisUnauth(address string) (flag bool, err error) {
	flag = false
	conn, err := WrapperTcpWithTimeout("tcp", address, 10*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return flag, err
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
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
	}
	return flag, err
}
