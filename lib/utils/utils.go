package utils

import (
	"io"
	"os"

	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func HomeDir() string {
	s, _ := os.UserHomeDir()
	return s
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

// renameOutput 用于清理URL中的非法字符，以生成合法的文件名。
func RenameOutput(filename string) string {
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "___", "_")
	return filename
}
