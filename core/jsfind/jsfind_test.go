package jsfind

import (
	"net/url"
	"testing"
)

func TestFindInfo(t *testing.T) {
	var fs FindSomething
	target := "http://www.baidu.com"
	jsLinks := ExtractJS(target)
	u, _ := url.Parse(target)

	for _, js := range jsLinks {
		newURL := u.Scheme + "://" + u.Host + js
		fs2 := FindInfo(newURL)
		fs.APIRoute = append(fs.APIRoute, fs2.APIRoute...)
	}
	fs.APIRoute = RemoveDuplicatesInfoSource(fs.APIRoute)
}
