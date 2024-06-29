package info

import (
	"fmt"
	"testing"
)

func TestBeianx(t *testing.T) {
	ips := Ip138IpHistory("bilibili.com")
	fmt.Printf("ips: %v\n", ips)
}
