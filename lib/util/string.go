package util

import "strings"

// 去除特殊字符
func RemoveIllegalChar(s string) string {
	for _, c := range []string{" ", "\r", "\n"} {
		if strings.Contains(s, c) {
			s = strings.ReplaceAll(s, c, "")
		}
	}
	return s
}
