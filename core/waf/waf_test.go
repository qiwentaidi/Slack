package waf

import (
	"fmt"
	"testing"
)

func TestIsWAF(t *testing.T) {
	fmt.Println(IsWAF("www.syrmyy.com.cn:443"))
}
