package jsfind

import (
	"errors"
	"net/http"
	"slack-wails/lib/clients"
	"strings"
)

func detectMethod(fullURL string, headers map[string]string) (string, error) {
	resp, body, err := clients.NewRequest("GET", fullURL, headers, nil, 5, false, http.DefaultClient)
	if err != nil {
		if strings.Contains(err.Error(), "doesn't contain any IP SANs") {
			return "", errors.New("证书中不包含使用的域名/IP, 请求失败")
		}
		return "", err
	}
	if resp == nil {
		return "", errors.New("响应内容为空")
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
