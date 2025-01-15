package util

import (
	"io"
	"net/url"
	"os"

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

// IsURL tests a string to determine if it is a well-structured url or not.
func IsURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
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
