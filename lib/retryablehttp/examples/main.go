package main

import (
	"fmt"
	"io/ioutil"

	"github.com/zan8in/retryablehttp"
)

func main() {
	opts := retryablehttp.DefaultOptionsSpraying

	client := retryablehttp.NewClient(opts)
	resp, err := client.Get("http://example.com")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Data: %v\n", string(data))

}
