package raw

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const (
	HTTP_PROXY_ENV = "HTTP_PROXY"
	SOCKS5         = "socks5"
	HTTP           = "http"
	HTTPS          = "https"
	DefaultTimeout = 10
)

var (
	// ProxyURL is the URL for the proxy server
	ProxyURL string
	// ProxySocksURL is the URL for the proxy socks server
	ProxySocksURL string
)

var proxyURLList []url.URL

// loadProxyServers load list of proxy servers from file or comma seperated
func LoadProxyServers(proxy string) error {
	if len(proxy) == 0 {
		return nil
	}
	if len(strings.Split(proxy, ",")) > 1 {
		for _, proxy := range strings.Split(proxy, ",") {
			if strings.TrimSpace(proxy) == "" {
				continue
			}
			if proxyURL, err := validateProxyURL(proxy); err != nil {
				return err
			} else {
				proxyURLList = append(proxyURLList, proxyURL)
			}
		}
	} else if proxyURL, err := validateProxyURL(proxy); err == nil {
		proxyURLList = append(proxyURLList, proxyURL)
	} else {
		return fmt.Errorf("invalid proxy file or URL provided for %s", proxy)
	}
	return processProxyList()
}

func processProxyList() error {
	if len(proxyURLList) == 0 {
		return fmt.Errorf("could not find any valid proxy")
	} else {
		done := make(chan bool)
		exitCounter := make(chan bool)
		counter := 0

		if len(proxyURLList) > 0 {
			for {
				select {
				case <-done:
					{
						close(done)
						return nil
					}
				case <-exitCounter:
					{
						if counter += 1; counter == len(proxyURLList) {
							return errors.New("no reachable proxy found")
						}
						close(exitCounter)
					}
				}
			}
		}
	}
	return nil
}

func validateProxyURL(proxy string) (url.URL, error) {
	if url, err := url.Parse(proxy); err == nil && isSupportedProtocol(url.Scheme) {
		return *url, nil
	}
	return url.URL{}, errors.New("invalid proxy format (It should be http[s]/socks5://[username:password@]host:port), ProxyURL: " + proxy)
}

// isSupportedProtocol checks given protocols are supported
func isSupportedProtocol(value string) bool {
	return value == HTTP || value == HTTPS || value == SOCKS5
}

func RandomIntWithMin(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return int(rand.Intn(max-min) + min)
}
