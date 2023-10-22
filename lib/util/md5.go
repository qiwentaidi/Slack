package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5Value(content []byte) string {
	hash := md5.New()
	io.WriteString(hash, string(content))
	return hex.EncodeToString(hash.Sum(nil))
}
