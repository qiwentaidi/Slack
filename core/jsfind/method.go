package jsfind

import (
	"errors"
	"slack-wails/lib/clients"
	"strings"
)

func detectMethod(fullURL string, headers map[string]string) (string, error) {
	resp, err := clients.DoRequest("GET", fullURL, headers, nil, 5, clients.NewRestyClient(nil, true))
	body := string(resp.Body())
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
	if (strings.Contains(body, "not supported") && strings.Contains(body, "Request method")) || resp.StatusCode() == 405 {
		return "POST", nil
	}

	if resp.StatusCode() == 401 {
		return "GET", errors.New("不存在未授权访问")
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("非正确API地址, 已忽略")
	}

	return "GET", nil
}

func detectContentType(url string, headers map[string]string) string {
	// 先浅拷贝一下headers，避免污染原map
	hdr := make(map[string]string)
	for k, v := range headers {
		hdr[k] = v
	}

	// 第一次，不带 Content-Type 测试
	resp, err := clients.DoRequest("POST", url, hdr, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil || resp == nil {
		return ""
	}

	// 如果第一次就能正常POST，不需要加Content-Type，默认返回 application/x-www-form-urlencoded
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return "application/x-www-form-urlencoded"
	}

	// 如果出错，或者响应提示Content-Type不对，再带上 application/x-www-form-urlencoded 重试一次
	hdr["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err = clients.DoRequest("POST", url, hdr, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil || resp == nil {
		return ""
	}
	body := string(resp.Body())

	// 判断返回内容
	if strings.Contains(body, "application/x-www-form-urlencoded") && strings.Contains(body, "not supported") {
		return "application/json"
	}

	if strings.Contains(body, "not a multipart request") {
		return "multipart/form-data;boundary=8ce4b16b22b58894aa86c421e8759df3"
	}

	// 默认
	return "application/x-www-form-urlencoded"
}
