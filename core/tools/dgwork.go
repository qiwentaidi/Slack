// https://github.com/xjl662750/dgwork-go/blob/main/client.go
package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"slack-wails/lib/clients"
	"strconv"
	"strings"
	"time"
)

type DGworkClient struct {
	BaseURL   string
	AppKey    string
	AppSecret string
}

// NewClient NewClient
func NewDGworkClient(BaseURL, AppKey, AppSecret string) *DGworkClient {
	return &DGworkClient{
		BaseURL:   BaseURL,
		AppKey:    AppKey,
		AppSecret: AppSecret,
	}
}

// ComputeHmac256 ComputeHmac256
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetFirstIP 获取本机的ip地址
func GetFirstIP() (ip string) {
	ip = "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}

// GetFirstMacAddress 获取本机的MAC地址
func GetFirstMacAddress() (macAddress string) {
	macAddress = "00:00:00:00:00:00"
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range interfaces {
		macAddress = inter.HardwareAddr.String()
		return
	}
	return
}

// 获得完整的路径
func (c *DGworkClient) GETURL(uri string) string {
	c.BaseURL = strings.TrimRight(c.BaseURL, "/")
	return c.BaseURL + uri
}

func (c *DGworkClient) GetHeaders(getURL string, params url.Values) map[string]string {
	mac := GetFirstMacAddress()
	ip := GetFirstIP()
	timestamp := time.Now().Format("2006-01-02T15:04:05.999999+08:00")
	nonce := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	reqURL, _ := url.Parse(getURL)
	message := "GET\n" + timestamp + "\n" + nonce + "\n" + reqURL.Path
	message = message + "\n" + params.Encode()
	signature := ComputeHmac256(message, c.AppSecret)
	return map[string]string{
		"X-Hmac-Auth-Timestamp": timestamp,
		"X-Hmac-Auth-Version":   "1.0",
		"X-Hmac-Auth-Nonce":     nonce,
		"apiKey":                c.AppKey,
		"X-Hmac-Auth-Signature": signature,
		"X-Hmac-Auth-IP":        ip,
		"X-Hmac-Auth-MAC":       mac,
	}
}

func (c *DGworkClient) GetToken(getURL string) string {
	params := url.Values{}
	params.Add("appsecret", c.AppSecret)
	params.Add("appkey", c.AppKey)
	getURL = getURL + "?" + params.Encode()
	_, body, err := clients.NewRequest("GET", getURL, c.GetHeaders(getURL, params), nil, 10, false, http.DefaultClient)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func (t *Tools) GetToken(baseURL, appkey, appsecret string) string {
	c := NewDGworkClient(baseURL, appkey, appsecret)
	return c.GetToken(c.GETURL("/gettoken.json"))
}
