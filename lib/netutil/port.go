package netutil

import (
	"net/url"
	"strconv"
)

func GetPort(u *url.URL) int {
	if u.Port() == "" {
		switch u.Scheme {
		case "http":
			return 80
		case "https":
			return 443
		}
	}
	port, _ := strconv.Atoi(u.Port())
	return port
}
