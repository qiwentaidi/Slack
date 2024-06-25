package space

import (
	"fmt"
	"testing"
)

func TestQuake(t *testing.T) {
	qk := QuakeApiSearch(&QuakeRequestOptions{
		`app:"Zabbix"`, 1, 10, true, false, false, false, "b13a9af7-f6b5-43cc-bff1-c807ea92dd72",
	})
	fmt.Printf("qk: %v\n", qk)
	// qr := SearchQuakeTips("jeecg")
	// fmt.Printf("qr: %v\n", qr)
}
