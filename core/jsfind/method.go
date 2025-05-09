package jsfind

import (
	"errors"
	"maps"
	"slack-wails/lib/clients"
	"strings"
)

func detectMethod(fullURL string, headers map[string]string) (string, error) {
	resp, err := clients.DoRequest("GET", fullURL, headers, nil, 5, clients.NewRestyClient(nil, true))
	if err != nil {
		if strings.Contains(err.Error(), "doesn't contain any IP SANs") {
			return "", errors.New("证书中不包含使用的域名/IP, 请求失败")
		}
		return "", err
	}
	body := string(resp.Body())
	// 模式错误情况 1
	if (strings.Contains(body, "not supported") && strings.Contains(body, "Request method")) || strings.Contains(body, "请求方式不支持") || strings.Contains(body, "请求方式错误") || resp.StatusCode() == 405 || strings.Contains(body, "Method Not Allowed") {
		return "POST", nil
	}
	switch resp.StatusCode() {
	case 401:
		return "GET", errors.New("响应码401, 不存在未授权访问")
	case 404:
		return "", errors.New("非正确API地址, 已忽略")
	default:
		return "GET", nil
	}
}

func detectContentType(url string, headers map[string]string) string {
	// 先浅拷贝一下 headers，避免污染原 headers
	hdr := make(map[string]string)
	maps.Copy(hdr, headers)

	// 第一次，不带 Content-Type 直接测试
	resp, err := clients.DoRequest("POST", url, hdr, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return ""
	}

	body := string(resp.Body())

	if strings.Contains(body, "not a multipart request") {
		return "multipart/form-data;boundary=8ce4b16b22b58894aa86c421e8759df3"
	}

	// 参数体缺失，一般都为json请求
	if strings.Contains(body, "Required request body is missing") {
		return "application/json"
	}

	// 第一次响应：如果返回 text/plain不支持则需要 content-type 字段
	if !strings.Contains(body, "not supported") && !strings.Contains(body, "text/plain") {
		return "text/plain"
	}

	// 需要重试，带上 application/x-www-form-urlencoded 重新请求
	hdr["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err = clients.DoRequest("POST", url, hdr, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return ""
	}

	body = string(resp.Body())

	if strings.Contains(body, "application/x-www-form-urlencoded") && strings.Contains(body, "not supported") {
		return "application/json"
	}

	// 默认返回 application/x-www-form-urlencoded
	return "application/x-www-form-urlencoded"
}
