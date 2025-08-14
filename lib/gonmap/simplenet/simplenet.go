package simplenet

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// netloc 主机地址
func tcpSend(protocol string, netloc string, data string, duration time.Duration, size int) (string, error) {
	protocol = strings.ToLower(protocol)
	conn, err := net.DialTimeout(protocol, netloc, duration)
	if err != nil {
		//fmt.Println(conn)
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	defer conn.Close()
	_, err = conn.Write([]byte(data))
	if err != nil {
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}
	//读取数据
	var buf []byte              // big buffer
	var tmp = make([]byte, 256) // using small tmo buffer for demonstrating
	var length int
	for {
		//设置读取超时Deadline
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		length, err = conn.Read(tmp)
		buf = append(buf, tmp[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			break
		}
		if len(buf) > size {
			break
		}
	}
	if err != nil && err != io.EOF {
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	if len(buf) == 0 {
		return "", errors.New("STEP3:response is empty")
	}
	return string(buf), nil
}

func tcpSendWithProxy(protocol string, netloc string, data string, duration time.Duration, size int, socksAddress string, auth *proxy.Auth) (string, error) {
	if socksAddress == "" {
		return "", errors.New("socks address is empty")
	}
	dialer, err := proxy.SOCKS5("tcp", socksAddress, auth, &net.Dialer{
		Timeout: duration,
	})
	if err != nil {
		return "", err
	}
	protocol = strings.ToLower(protocol)
	conn, err := dialer.Dial(protocol, netloc)
	if err != nil {
		//fmt.Println(conn)
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	defer conn.Close()
	_, err = conn.Write([]byte(data))
	if err != nil {
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}
	//读取数据
	var buf []byte              // big buffer
	var tmp = make([]byte, 256) // using small tmo buffer for demonstrating
	var length int
	for {
		//设置读取超时Deadline
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		length, err = conn.Read(tmp)
		buf = append(buf, tmp[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			break
		}
		if len(buf) > size {
			break
		}
	}
	if err != nil && err != io.EOF {
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	if len(buf) == 0 {
		return "", errors.New("STEP3:response is empty")
	}
	return string(buf), nil
}

// tlsSend函数用于发送TLS请求
func tlsSend(protocol string, netloc string, data string, duration time.Duration, size int) (string, error) {
	// 将协议转换为小写
	protocol = strings.ToLower(protocol)
	// 创建TLS配置
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}
	// 创建网络拨号器
	dialer := &net.Dialer{
		Timeout:  duration,
		Deadline: time.Now().Add(duration * 2),
	}
	// 使用网络拨号器和TLS配置建立连接
	conn, err := tls.DialWithDialer(dialer, protocol, netloc, config)
	if err != nil {
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	defer conn.Close()
	// 向连接中写入数据
	_, err = io.WriteString(conn, data)
	if err != nil {
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}
	// 读取数据
	var buf []byte              // big buffer
	var tmp = make([]byte, 256) // using small tmo buffer for demonstrating
	var length int
	for {
		// 设置读取超时Deadline
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		length, err = conn.Read(tmp)
		buf = append(buf, tmp[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			break
		}
		if len(buf) > size {
			break
		}
	}
	if err != nil && err != io.EOF {
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	if len(buf) == 0 {
		return "", errors.New("STEP3:response is empty")
	}
	return string(buf), nil
}

func tlsSendWithProxy(protocol string, netloc string, data string, duration time.Duration, size int, socksAddress string, auth *proxy.Auth) (string, error) {
	if socksAddress == "" {
		return "", errors.New("socks address is empty")
	}

	// 创建 SOCKS5 代理 Dialer
	dialer, err := proxy.SOCKS5("tcp", socksAddress, auth, &net.Dialer{
		Timeout: duration,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	// 使用 Dialer 建立原始连接
	rawConn, err := dialer.Dial(protocol, netloc)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", netloc, err)
	}
	defer rawConn.Close()

	// 封装为 TLS 连接
	tlsConn := tls.Client(rawConn, &tls.Config{
		InsecureSkipVerify: true, // 跳过证书验证（生产环境需要谨慎）
		ServerName:         strings.Split(netloc, ":")[0],
	})

	// 执行 TLS 握手
	err = tlsConn.Handshake()
	if err != nil {
		return "", fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer tlsConn.Close()

	// 发送数据
	_, err = tlsConn.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("failed to send data: %w", err)
	}

	// 读取数据
	var buf []byte              // large buffer
	var tmp = make([]byte, 256) // small temporary buffer
	var length int
	for {
		// 设置超时
		_ = tlsConn.SetReadDeadline(time.Now().Add(3 * time.Second))
		length, err = tlsConn.Read(tmp)
		buf = append(buf, tmp[:length]...)
		if length < len(tmp) {
			break
		}
		if err != nil {
			break
		}
		if len(buf) > size {
			break
		}
	}
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	if len(buf) == 0 {
		return "", errors.New("response is empty")
	}

	return string(buf), nil
}

func Send(protocol string, tls bool, netloc string, data string, duration time.Duration, size int, proxyURL string) (string, error) {
	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		if u.User != nil {
			username := u.User.Username()
			password, _ := u.User.Password()
			if tls {
				return tlsSendWithProxy(protocol, netloc, data, duration, size, u.Host, &proxy.Auth{User: username, Password: password})
			} else {
				return tcpSendWithProxy(protocol, netloc, data, duration, size, u.Host, &proxy.Auth{User: username, Password: password})
			}
		}
	}
	if tls {
		return tlsSend(protocol, netloc, data, duration, size)
	} else {
		return tcpSend(protocol, netloc, data, duration, size)
	}
}
