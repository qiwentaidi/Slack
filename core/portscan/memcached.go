package portscan

import (
	"strings"
	"time"
)

func MemcachedScan(host string) *Burte {
	client, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(10 * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						return &Burte{
							Status:   true,
							Host:     host,
							Protocol: "memcached",
							Username: "unauthorized",
							Password: "",
						}
					}
				}
			}
		}
	}
	return &Burte{
		Status:   false,
		Host:     host,
		Protocol: "memcached",
		Username: "",
		Password: "",
	}
}
