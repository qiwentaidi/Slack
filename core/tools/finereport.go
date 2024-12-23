package core

import (
	"strings"
)

func cleanEncodedString(encoded string) string {
	encoded = strings.ReplaceAll(encoded, "\\u000d", "\r")
	encoded = strings.ReplaceAll(encoded, "\\u000a", "\n")
	encoded = strings.ReplaceAll(encoded, "\\r\\n", "\n")
	encoded = strings.ReplaceAll(encoded, "\\=", "=")
	return encoded
}
