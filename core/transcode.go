package core

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"
)

func RecursiveEncrypt(encrypt []string, input string, output *string) {
	if len(encrypt) == 0 {
		*output = input
		return
	}
	var encrypted string
	switch encrypt[0] {
	case "Base64":
		encrypted = base64.StdEncoding.EncodeToString([]byte(input))
	case "Base64 URL":
		encrypted = base64.URLEncoding.EncodeToString([]byte(input))
	case "URLcode":
		encrypted = url.QueryEscape(input)
	case "Unicode":
		temp := strconv.QuoteToASCII(input)
		encrypted = temp[1 : len(temp)-1]
	case "HEX":
		encrypted = hex.EncodeToString([]byte(input))
	case "HTML":
		encrypted = html.EscapeString(input)
	case "Ascii":
		var temp string
		for _, v := range input {
			temp += fmt.Sprintf("%v ", v)
		}
		encrypted = temp
	default:
		encrypted = input
	}
	RecursiveEncrypt(encrypt[1:], encrypted, output)
}

func RecursiveDecrypt(decrypt []string, input string, output *string) {
	if len(decrypt) == 0 {
		*output = input
		return
	}
	var decrypted string
	switch decrypt[0] {
	case "Base64":
		data, _ := base64.StdEncoding.DecodeString(input)
		decrypted = string(data)
	case "Base64 URL":
		data, _ := base64.URLEncoding.DecodeString(input)
		decrypted = string(data)
	case "URLcode":
		data, _ := url.QueryUnescape(input)
		decrypted = data
	case "Unicode":
		data, _ := strconv.Unquote("\"" + input + "\"")
		decrypted = data
	case "HEX":
		data, _ := hex.DecodeString(input)
		decrypted = string(data)
	case "HTML":
		html.UnescapeString(input)
	case "Ascii":
		temp := ""
		if strings.Contains(input, " ") {
			for _, s := range strings.Split(input, " ") {
				if i, err := strconv.Atoi(s); err == nil {
					temp += fmt.Sprintf("%c", i)
				} else {
					temp += fmt.Sprintf("what?(%v)", s)
				}
			}
		} else {
			if i, err := strconv.Atoi(input); err == nil {
				temp = fmt.Sprintf("%c", i)
			}
		}
		decrypted = temp
	default:
		decrypted = input
	}
	RecursiveDecrypt(decrypt[1:], decrypted, output)
}
