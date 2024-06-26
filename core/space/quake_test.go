package space

import (
	"fmt"
	"testing"
)

func TestQuake(t *testing.T) {
	qk := QuakeApiSearch(&QuakeRequestOptions{
		`app:"Zabbix"`, []string{"183.131.92.179",
			"183.131.92.171",
			"183.131.92.172"}, 1, 100, true, false, false, false, "",
	})
	fmt.Printf("qk: %v\n", qk)
	// qr := SearchQuakeTips("jeecg")
	// fmt.Printf("qr: %v\n", qr)
}
