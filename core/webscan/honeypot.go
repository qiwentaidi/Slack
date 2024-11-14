package webscan

import "strings"

func checkHoneypotWithHeaders(rawHeaders string) bool {
	count := strings.Count(rawHeaders, "Set-Cookie")
	return count > 5
}

func checkHoneypotWithFingerprintLength(length int) bool {
	return length > 15
}
