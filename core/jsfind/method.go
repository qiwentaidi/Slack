package jsfind

import (
	"errors"
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"strings"
)

func detectMethod(fullURL string, headers map[string]string) (string, error) {
	resp, body, err := clients.NewRequest("GET", fullURL, headers, nil, 5, false, http.DefaultClient)
	if err != nil || resp == nil {
		fmt.Printf("err: %v\n", err)
		return "", errors.New("请求失败")
	}
	// 模式错误情况 1
	if (strings.Contains(string(body), "not supported") && strings.Contains(string(body), "Request method")) || resp.StatusCode == 405 {
		return "POST", nil
	}
	if resp.StatusCode != 200 {
		return "", errors.New("非正确API地址, 已忽略")
	}

	return "GET", nil
}
