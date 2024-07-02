package jsfind

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
)

func TestFindInfo(t *testing.T) {
	var fs *FindSomething
	target := "https://www.baidu.com/"
	jsLinks := ExtractJS(context.TODO(), target)
	fmt.Printf("jsLinks: %v\n", jsLinks)
	var wg sync.WaitGroup
	limiter := make(chan bool, 100)
	wg.Add(1)
	limiter <- true
	go func() {
		fs = FindInfo(context.TODO(), target, limiter, &wg)
	}()
	wg.Wait()
	u, _ := url.Parse(target)
	fs.JS = *AppendSource(target, jsLinks)
	host := u.Scheme + "://" + u.Host
	for _, jslink := range jsLinks {
		wg.Add(1)
		limiter <- true
		go func(js string) {
			var newURL string
			if strings.HasPrefix(js, "http") {
				newURL = js
			} else {
				newURL = host + "/" + js
			}
			fs2 := FindInfo(context.TODO(), newURL, limiter, &wg)
			fs.IP_URL = append(fs.IP_URL, fs2.IP_URL...)
			fs.ChineseIDCard = append(fs.ChineseIDCard, fs2.ChineseIDCard...)
			fs.ChinesePhone = append(fs.ChinesePhone, fs2.ChinesePhone...)
			fs.SensitiveField = append(fs.SensitiveField, fs2.SensitiveField...)
			fs.APIRoute = append(fs.APIRoute, fs2.APIRoute...)
		}(jslink)
	}
	wg.Wait()
	fs.APIRoute = RemoveDuplicatesInfoSource(fs.APIRoute)
	fs.ChineseIDCard = RemoveDuplicatesInfoSource(fs.ChineseIDCard)
	fs.ChinesePhone = RemoveDuplicatesInfoSource(fs.ChinesePhone)
	fs.SensitiveField = RemoveDuplicatesInfoSource(fs.SensitiveField)
	fs.IP_URL = RemoveDuplicatesInfoSource(fs.IP_URL)
}
