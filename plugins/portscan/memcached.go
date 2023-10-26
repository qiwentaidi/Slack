package portscan

import (
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"strings"
	"time"
)

func MemcachedScan(host string) {
	client, err := WrapperTcpWithTimeout("tcp", host, time.Duration(common.Profile.PortScan.Timeout)*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(time.Duration(common.Profile.PortScan.Timeout) * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						custom.Console.Append(fmt.Sprintf("[+] Memcached %s unauthorized\n", host))
						common.PortBurstResult = append(common.PortBurstResult, []string{"Memcached", host, "", "unauthorized", ""})
					}
				} else {
					custom.Console.Append(fmt.Sprintf("[-] Memcached %v %v\n", host, err))
				}
			}
		}
	}
}
