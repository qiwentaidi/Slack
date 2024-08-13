package waf

import (
	"fmt"
	"testing"
)

func TestIsWAF(t *testing.T) {
	fmt.Println(IsWAF("www.syrmyy.com.cn:443", []string{"114.114.114.114:53", "223.5.5.5:53"}))
}
