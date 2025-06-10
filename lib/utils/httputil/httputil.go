package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func PrettyURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url
}

func PrettyPath(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimLeft(path, "/")
	}
	return path
}

// LimitResponse 通用响应截断函数
func LimitResponse(response string, size int, customMsg string) string {
	if size == 0 {
		return ""
	}
	runes := []rune(response)
	if len(runes) > size {
		if customMsg != "" {
			return customMsg
		}
		return string(runes[:size]) + " ..."
	}
	return response
}

func LimitResponseBytes(response []byte, size int) []byte {
	if len(response) > size {
		return response[:size]
	}
	return response
}

// DumpResponseHeadersOnly 只返回响应头
func DumpResponseHeadersOnly(resp *http.Response) []byte {
	headers, _ := httputil.DumpResponse(resp, false)
	return headers
}

// DumpResponseHeadersAndRaw returns http headers and response as strings
func DumpResponseHeadersAndRaw(resp *http.Response) (headers, fullresp []byte, err error) {
	// httputil.DumpResponse does not work with websockets
	if resp.StatusCode >= http.StatusContinue && resp.StatusCode <= http.StatusEarlyHints {
		raw := resp.Status + "\n"
		for h, v := range resp.Header {
			raw += fmt.Sprintf("%s: %s\n", h, v)
		}
		return []byte(raw), []byte(raw), nil
	}
	headers, err = httputil.DumpResponse(resp, false)
	if err != nil {
		return
	}
	// logic same as httputil.DumpResponse(resp, true) but handles
	// the edge case when we get both error and data on reading resp.Body
	var buf1, buf2 bytes.Buffer
	b := resp.Body
	if _, err = buf1.ReadFrom(b); err != nil {
		if buf1.Len() <= 0 {
			return
		}
	}
	if err == nil {
		b.Close()
	}

	// rewind the body to allow full dump
	resp.Body = io.NopCloser(bytes.NewReader(buf1.Bytes()))
	err = resp.Write(&buf2)
	fullresp = buf2.Bytes()

	// rewind once more to allow further reuses
	resp.Body = io.NopCloser(bytes.NewReader(buf1.Bytes()))
	return
}

// 获取基本URL
func GetBasicURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Scheme + "://" + u.Host
}

// 获取基本路径
func GetBasePath(inputURL string) (string, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	// 获取路径部分
	dirPath := path.Dir(parsedURL.Path)
	// 确保路径以 `/` 结尾
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	// 组合完整 URL
	basePath := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, dirPath)
	return basePath, nil
}

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
