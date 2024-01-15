package retryablehttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"
)

var (
	forceMaxRedirects int
	poolMutex         *sync.RWMutex
	normalClient      *Client
	clientPool        map[string]*Client
)

type PoolOptions struct {
	// Retries is the number of times to retry the request
	Retries int
	// MaxRedirects is the maximum numbers of redirects to be followed.
	MaxRedirects int
	// FollowRedirects enables following redirects for http request module
	FollowRedirects bool
	// FollowRedirects enables following redirects for http request module only on the same host
	FollowHostRedirects bool
	// DisableRedirects disables following redirects for http request module
	DisableRedirects bool
	// SNI custom hostname
	SNI string
	// ClientCertFile client certificate file (PEM-encoded) used for authenticating against scanned hosts
	ClientCertFile string
	// ClientKeyFile client key file (PEM-encoded) used for authenticating against scanned hosts
	ClientKeyFile string
	// ClientCAFile client certificate authority file (PEM-encoded) used for authenticating against scanned hosts
	ClientCAFile string
	// Timeout is the seconds to wait for a response from the server.
	Timeout int
	// description: |
	//   Threads specifies number of threads to use sending requests. This enables Connection Pooling.
	//
	//   Connection: Close attribute must not be used in request while using threads flag, otherwise
	//   pooling will fail and engine will continue to close connections after requests.
	// examples:
	//   - name: Send requests using 10 concurrent threads
	//     value: 10
	Threads int
	// NoTimeout disables http request timeout for context based usage
	NoTimeout bool
	// CookieReuse enables cookie reuse for the http client (cookiejar impl)
	CookieReuse bool
	// FollowRedirects specifies the redirects flow
	RedirectFlow RedirectFlow
	// Connection defines custom connection configuration
	Connection *ConnectionConfiguration
	// Proxy string defines
	Proxy string
}

var DefaultPoolOptions = PoolOptions{
	Retries:             1,
	MaxRedirects:        10,
	FollowRedirects:     false,
	FollowHostRedirects: false,
	DisableRedirects:    true,
	SNI:                 "",
	ClientCertFile:      "",
	ClientKeyFile:       "",
	ClientCAFile:        "",
	Timeout:             10,
	Threads:             0,
	NoTimeout:           false,
	CookieReuse:         false,
	RedirectFlow:        DontFollowRedirect,
	Connection:          nil,
	Proxy:               "",
}

func InitClientPool(options *PoolOptions) error {
	if normalClient != nil {
		return nil
	}

	// if len(options.Proxy) > 0 {
	// 	if err := LoadProxyServers(options.Proxy); err != nil {
	// 		return err
	// 	}
	// }

	forceMaxRedirects = options.MaxRedirects

	poolMutex = &sync.RWMutex{}
	clientPool = make(map[string]*Client)

	return nil
}

// ConnectionConfiguration contains the custom configuration options for a connection
type ConnectionConfiguration struct {
	// DisableKeepAlive of the connection
	DisableKeepAlive bool
	cookiejar        *cookiejar.Jar
	mu               sync.RWMutex
}

func (cc *ConnectionConfiguration) SetCookieJar(cookiejar *cookiejar.Jar) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.cookiejar = cookiejar
}

func (cc *ConnectionConfiguration) GetCookieJar() *cookiejar.Jar {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.cookiejar
}

func (cc *ConnectionConfiguration) HasCookieJar() bool {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return cc.cookiejar != nil
}

// Configuration contains the custom configuration options for a client
type Configuration struct {
	// Threads contains the threads for the client
	Threads int
	// MaxRedirects is the maximum number of redirects to follow
	MaxRedirects int
	// NoTimeout disables http request timeout for context based usage
	NoTimeout bool
	// CookieReuse enables cookie reuse for the http client (cookiejar impl)
	CookieReuse bool
	// FollowRedirects specifies the redirects flow
	RedirectFlow RedirectFlow
	// Connection defines custom connection configuration
	Connection *ConnectionConfiguration
}

type RedirectFlow uint8

const (
	DontFollowRedirect RedirectFlow = iota
	FollowSameHostRedirect
	FollowAllRedirect
)

// Hash returns the hash of the configuration to allow client pooling
func (c *PoolOptions) Hash() string {
	builder := &strings.Builder{}
	builder.Grow(16)
	builder.WriteString("t")
	builder.WriteString(strconv.Itoa(c.Threads))
	builder.WriteString("m")
	builder.WriteString(strconv.Itoa(c.MaxRedirects))
	builder.WriteString("n")
	builder.WriteString(strconv.FormatBool(c.NoTimeout))
	builder.WriteString("f")
	builder.WriteString(strconv.Itoa(int(c.RedirectFlow)))
	builder.WriteString("r")
	builder.WriteString(strconv.FormatBool(c.CookieReuse))
	builder.WriteString("c")
	builder.WriteString(strconv.FormatBool(c.Connection != nil))
	hash := builder.String()
	return hash
}

func (options *PoolOptions) EnableRedirect(flow RedirectFlow) {
	if flow == DontFollowRedirect {
		options.DisableRedirects = true
		options.FollowRedirects = false
		options.FollowHostRedirects = false
	}
	if flow == FollowSameHostRedirect {
		options.DisableRedirects = false
		options.FollowRedirects = false
		options.FollowHostRedirects = true
	}
	if flow == FollowAllRedirect {
		options.DisableRedirects = false
		options.FollowRedirects = true
		options.FollowHostRedirects = false
	}
	options.RedirectFlow = flow
}

// HasStandardOptions checks whether the configuration requires custom settings
func (c *PoolOptions) HasStandardOptions() bool {
	return c.Threads == 0 && c.MaxRedirects == 0 && c.RedirectFlow == DontFollowRedirect && !c.CookieReuse && c.Connection == nil && !c.NoTimeout
}

// Get creates or gets a client for the protocol based on custom configuration
func GetPool(options *PoolOptions) (*Client, error) {
	if options.HasStandardOptions() {
		return normalClient, nil
	}

	if len(options.Proxy) > 0 && len(ProxyURL) == 0 && len(ProxySocksURL) == 0 {
		if err := LoadProxyServers(options.Proxy); err != nil {
			return normalClient, err
		}
	}

	return wrappedGet(options)
}

// wrappedGet wraps a get operation without normal client check
func wrappedGet(options *PoolOptions) (*Client, error) {
	var err error

	hash := options.Hash()
	poolMutex.RLock()
	if client, ok := clientPool[hash]; ok {
		poolMutex.RUnlock()
		return client, nil
	}
	poolMutex.RUnlock()

	// Multiple Host
	retryableHttpOptions := DefaultOptionsSpraying
	disableKeepAlives := true
	maxIdleConns := 0
	maxConnsPerHost := 0
	maxIdleConnsPerHost := -1

	if options.Threads > 0 {
		// Single host
		retryableHttpOptions = DefaultOptionsSingle
		disableKeepAlives = false
		maxIdleConnsPerHost = 500
		maxConnsPerHost = 500
	}

	retryableHttpOptions.RetryWaitMax = 10 * time.Second
	retryableHttpOptions.RetryMax = options.Retries
	redirectFlow := options.RedirectFlow
	maxRedirects := options.MaxRedirects

	if forceMaxRedirects > 0 {
		// by default we enable general redirects following
		switch {
		case options.FollowHostRedirects:
			redirectFlow = FollowSameHostRedirect
		default:
			redirectFlow = FollowAllRedirect
		}
		maxRedirects = forceMaxRedirects
	}
	if options.DisableRedirects {
		options.FollowRedirects = false
		options.FollowHostRedirects = false
		redirectFlow = DontFollowRedirect
		maxRedirects = 0
	}

	// override connection's settings if required
	if options.Connection != nil {
		disableKeepAlives = options.Connection.DisableKeepAlive
	}

	// Set the base TLS configuration definition
	tlsConfig := &tls.Config{
		Renegotiation:      tls.RenegotiateOnceAsClient,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}

	if options.SNI != "" {
		tlsConfig.ServerName = options.SNI
	}
	// Add the client certificate authentication to the request if it's configured
	tlsConfig, err = AddConfiguredClientCertToRequest(tlsConfig, options)
	if err != nil {
		return nil, fmt.Errorf("%s could not create client certificate", err.Error())
	}

	dialer := &net.Dialer{}

	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		MaxConnsPerHost:     maxConnsPerHost,
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   disableKeepAlives,
	}

	if ProxyURL != "" {
		if proxyURL, err := url.Parse(ProxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	} else if ProxySocksURL != "" {
		socksURL, proxyErr := url.Parse(ProxySocksURL)
		if proxyErr != nil {
			return nil, proxyErr
		}
		dialer, err := proxy.FromURL(socksURL, proxy.Direct)
		if err != nil {
			return nil, err
		}

		dc := dialer.(interface {
			DialContext(ctx context.Context, network, addr string) (net.Conn, error)
		})
		if proxyErr == nil {
			transport.DialContext = dc.DialContext
			transport.DialTLSContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				// upgrade proxy connection to tls
				conn, err := dc.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				return tls.Client(conn, tlsConfig), nil
			}
		}
	}

	var jar *cookiejar.Jar
	if options.Connection != nil && options.Connection.HasCookieJar() {
		jar = options.Connection.GetCookieJar()
	} else if options.CookieReuse {
		if jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
			return nil, fmt.Errorf("%s could not create cookiejar", err.Error())
		}
	}

	httpclient := &http.Client{
		Transport:     transport,
		CheckRedirect: makeCheckRedirectFunc(redirectFlow, maxRedirects),
	}
	if !options.NoTimeout {
		httpclient.Timeout = time.Duration(options.Timeout) * time.Second
	}
	client := NewWithHTTPClient(httpclient, retryableHttpOptions)
	if jar != nil {
		client.HTTPClient.Jar = jar
	}
	client.CheckRetry = HostSprayRetryPolicy()

	// Only add to client pool if we don't have a cookie jar in place.
	if jar == nil {
		poolMutex.Lock()
		clientPool[hash] = client
		poolMutex.Unlock()
	}
	return client, nil
}

// AddConfiguredClientCertToRequest adds the client certificate authentication to the tls.Config object and returns it
func AddConfiguredClientCertToRequest(tlsConfig *tls.Config, options *PoolOptions) (*tls.Config, error) {
	// Build the TLS config with the client certificate if it has been configured with the appropriate options.
	// Only one of the options needs to be checked since the validation checks in main.go ensure that all three
	// files are set if any of the client certification configuration options are.
	if len(options.ClientCertFile) > 0 {
		// Load the client certificate using the PEM encoded client certificate and the private key file
		cert, err := tls.LoadX509KeyPair(options.ClientCertFile, options.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}

		// Load the certificate authority PEM certificate into the TLS configuration
		caCert, err := os.ReadFile(options.ClientCAFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}
	return tlsConfig, nil
}

const defaultMaxRedirects = 10

type checkRedirectFunc func(req *http.Request, via []*http.Request) error

func makeCheckRedirectFunc(redirectType RedirectFlow, maxRedirects int) checkRedirectFunc {
	return func(req *http.Request, via []*http.Request) error {
		switch redirectType {
		case DontFollowRedirect:
			return http.ErrUseLastResponse
		case FollowSameHostRedirect:
			var newHost = req.URL.Host
			var oldHost = via[0].Host
			if oldHost == "" {
				oldHost = via[0].URL.Host
			}
			if newHost != oldHost {
				// Tell the http client to not follow redirect
				return http.ErrUseLastResponse
			}
			return checkMaxRedirects(req, via, maxRedirects)
		case FollowAllRedirect:
			return checkMaxRedirects(req, via, maxRedirects)
		}
		return nil
	}
}

func checkMaxRedirects(req *http.Request, via []*http.Request, maxRedirects int) error {
	if maxRedirects == 0 {
		if len(via) > defaultMaxRedirects {
			return http.ErrUseLastResponse
		}
		return nil
	}

	if len(via) > maxRedirects {
		return http.ErrUseLastResponse
	}
	return nil
}
