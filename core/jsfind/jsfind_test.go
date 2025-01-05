package jsfind

import (
	"fmt"
	"net/url"
	"slack-wails/lib/clients"
	"testing"
)

func TestFindInfo(t *testing.T) {

}

// 发送 GET 请求，直到参数补全
func sendRequest(apiURL string, params url.Values) url.Values {
	// 构造完整 URL
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// 发送 GET 请求
	_, body, err := clients.NewSimpleGetRequest(fullURL, clients.NewHttpClient(nil, true))
	if err != nil {
		fmt.Println("请求失败:", err)
		return params
	}

	// 提取缺失参数
	missingParam := extractMissingParams(string(body))
	if missingParam != nil {
		// 生成默认值并补全参数
		defaultValue := generateDefaultValue(missingParam.Type)
		params.Set(missingParam.Name, fmt.Sprint(defaultValue))
		// 递归调用，直到所有参数补全
		return sendRequest(apiURL, params)
	}
	// 如果没有缺失参数提示，返回当前参数集
	fmt.Println(fullURL)
	return params
}
