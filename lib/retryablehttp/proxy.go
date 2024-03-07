package retryablehttp

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	SOCKS5 = "socks5"
	HTTP   = "http"
	HTTPS  = "https"
)

var (
	// ProxyURL is the URL for the proxy server
	ProxyURL string
	// ProxySocksURL is the URL for the proxy socks server
	ProxySocksURL string
)

// loadProxyServers load list of proxy servers from file or comma seperated
func LoadProxyServers(proxy string) error {
	if len(proxy) == 0 {
		return nil
	}

	proxyURL, err := validateProxyURL(proxy)
	if err != nil {
		return fmt.Errorf("invalid proxy file or URL provided for %s", proxy)
	}

	if proxyURL.Scheme == HTTP || proxyURL.Scheme == HTTPS {
		ProxyURL = proxyURL.String()
		ProxySocksURL = ""
	} else if proxyURL.Scheme == SOCKS5 {
		ProxyURL = ""
		ProxySocksURL = proxyURL.String()
	}
	return nil
}

func validateProxyURL(proxy string) (url.URL, error) {
	if url, err := url.Parse(proxy); err == nil {
		return *url, nil
	}
	return url.URL{}, errors.New("invalid proxy format (It should be http[s]/socks5://[username:password@]host:port), ProxyURL: " + proxy)
}
