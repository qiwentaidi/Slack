package util

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func HomeDir() string {
	s, _ := os.UserHomeDir()
	return s
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func GetNowDateTime() string {
	now := time.Now()
	return now.Format("01-02 15:04:05")
}

func GetNowDateTimeReportName() string {
	now := time.Now()
	return now.Format("20060102-150405")
}

// 字符串转 utf 8
func Str2UTF8(str string) string {
	if len(str) == 0 {
		return ""
	}
	if !utf8.ValidString(str) {
		utf8Bytes, _ := io.ReadAll(transform.NewReader(
			strings.NewReader(str),
			simplifiedchinese.GBK.NewDecoder(),
		))
		return string(utf8Bytes)
	}
	return str
}

func GetItemInArray(a []string, s string) int {
	for index, v := range a {
		if v == s {
			return index
		}
	}
	return -1
}

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
	// if !strings.HasSuffix(dirPath, "/") {
	// 	dirPath += "/"
	// }

	// 组合完整 URL
	basePath := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, dirPath)
	return basePath, nil
}

// renameOutput 用于清理URL中的非法字符，以生成合法的文件名。
func RenameOutput(filename string) string {
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "___", "_")
	return filename
}
