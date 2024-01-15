package main

import (
	"fmt"

	"github.com/zan8in/retryablehttp"
)

func main() {
	resp, err := retryablehttp.DefaultClientPool.Get("http://47.104.36.43:8425")
	fmt.Println(resp, err)
}
