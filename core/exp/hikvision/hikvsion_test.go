package hikvision

import (
	"fmt"
	"testing"
)

func TestHikvision(t *testing.T) {
	URL := "http://47.150.37.246:81/"
	// resp := CVE_2017_7921_Snapshot(URL, clients.DefaultClient())
	// fmt.Printf("resp: %v\n", resp)

	// resp := CVE_2017_7921_Config(URL, clients.DefaultClient())
	// fmt.Printf("resp: %v\n", resp)
	resp := CheckLogin(URL, []string{"hik12345+", "Hik12345+", "admin12345", "12345"})
	fmt.Printf("resp: %v\n", resp)
}
