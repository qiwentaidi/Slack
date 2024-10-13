package webscan

import (
	"fmt"
	"net/url"
	"slack-wails/lib/clients"
	"testing"
)

func TestInfoscan(t *testing.T) {
	u, _ := url.Parse("x")
	hash, md5 := FaviconHash(u, clients.DefaultClient())
	fmt.Printf("hash: %v\n", hash)
	fmt.Printf("md5: %v\n", md5)
	resp, body, err := clients.NewSimpleGetRequest("http://27.71.26.188:8072/", clients.DefaultClient())
	fmt.Printf("resp.StatusCode: %v\n", resp.StatusCode)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("string(body): %v\n", string(body))
}
