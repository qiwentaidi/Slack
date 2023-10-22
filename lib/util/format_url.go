package util

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	HTTP             = "http"
	HTTPS            = "https"
	SchemeSeparator  = "://"
	DefaultHTTPPort  = "80"
	DefaultHTTPSPort = "443"
)

func Hostname(s string) (string, error) {
	if !strings.HasPrefix(s, HTTP) && !strings.HasPrefix(s, HTTPS) {
		s = HTTP + SchemeSeparator + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func Host(s string) (string, error) {
	if !strings.HasPrefix(s, HTTP) && !strings.HasPrefix(s, HTTPS) {
		s = HTTP + SchemeSeparator + s
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

// URL Encode all characters
// URLEncodeAllChar("1' and sleep(5) and '1'='1")
// %31%27%20%61%6E%64%20%73%6C%65%65%70%28%35%29%20%61%6E%64%20%27%31%27%3D%27%31
func URLEncodeAllChar(s string) string {
	b := []byte(s)
	nb := ""
	for _, v := range b {
		nb += "%" + fmt.Sprintf("%X", v)
	}
	return nb
}
