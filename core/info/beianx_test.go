package info

import (
	"fmt"
	"testing"
)

func TestBeianx(t *testing.T) {
	// d := Beianx("深圳市鲸瀚数字文化有限公司")
	// fmt.Printf("d: %v\n", d)
	ips := Ip138IpHistory("bilibili.com")
	fmt.Printf("ips: %v\n", ips)
}
