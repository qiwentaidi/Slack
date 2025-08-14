package httputil

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"regexp"
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

func DumpResponseHeadersAndDecodedBody(resp *http.Response) (string, error) {
	// 1. dump header（不包含 body）
	headerDump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return "", err
	}

	// 2. 读取并解压 body
	var reader io.Reader = resp.Body
	defer resp.Body.Close()

	switch strings.ToLower(resp.Header.Get("Content-Encoding")) {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("gzip decode failed: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		// deflate 可能是 zlib 包裹，也可能是原始 deflate 数据
		zlibReader, err := zlib.NewReader(resp.Body)
		if err == nil {
			defer zlibReader.Close()
			reader = zlibReader
		} else {
			flateReader := flate.NewReader(resp.Body)
			defer flateReader.Close()
			reader = flateReader
		}
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// 3. 转成字符串处理 header
	headerStr := string(headerDump)

	// 移除 Transfer-Encoding: chunked
	re := regexp.MustCompile(`(?i)Transfer-Encoding:\s*chunked\r\n`)
	headerStr = re.ReplaceAllString(headerStr, "")

	// 如果没有 Content-Length，就插入到 header 末尾（空行之前）
	if !regexp.MustCompile(`(?i)Content-Length:`).MatchString(headerStr) {
		if idx := strings.Index(headerStr, "\r\n\r\n"); idx != -1 {
			headerStr = headerStr[:idx] +
				fmt.Sprintf("\r\nContent-Length: %d", len(bodyBytes)) +
				headerStr[idx:]
		}
	}

	// 如果 Content-Encoding 存在且被解压了，需要移除
	if resp.Header.Get("Content-Encoding") != "" {
		reEnc := regexp.MustCompile(`(?i)Content-Encoding:\s*\S+\r\n`)
		headerStr = reEnc.ReplaceAllString(headerStr, "")
	}

	// 4. 拼接 header 和解压后的 body
	var buf bytes.Buffer
	buf.WriteString(headerStr)
	buf.WriteString("\r\n") // header-body 分隔
	buf.Write(bodyBytes)

	return buf.String(), nil
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
