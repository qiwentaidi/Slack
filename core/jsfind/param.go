package jsfind

import (
	"fmt"
	"net/url"
	"regexp"
	"slack-wails/lib/clients"
)

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var extractMissingRegex = regexp.MustCompile(`Required (String|Int|Long|Double|Boolean|Date).*?'([^']+)'`)

// 从错误信息中提取缺失参数的名称
func extractMissingParams(message string) *Parameter {
	// 提取匹配内容
	matches := extractMissingRegex.FindStringSubmatch(message)
	// 输出结果
	if len(matches) > 2 {
		return &Parameter{
			Name: matches[2],
			Type: matches[1],
		}
	}
	return nil
}

// 根据参数类型生成默认值
func generateDefaultValue(paramType string) interface{} {
	switch paramType {
	case "String":
		return "test"
	case "Int":
		return 0
	case "Long":
		return int64(0)
	case "Double":
		return 0.0
	case "Boolean":
		return false
	case "Date":
		return "1970-01-01"
	default:
		return "defaultValue"
	}
}

// 参数补全
func completeParameters(method, apiURL string, params url.Values) url.Values {
	// 构造完整 URL
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// 发送 GET 请求
	resp, err := clients.DoRequest(method, fullURL, nil, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		fmt.Println("请求失败:", err)
		return nil
	}

	// 提取缺失参数
	missingParam := extractMissingParams(string(resp.Body()))
	if missingParam != nil {
		// 生成默认值并补全参数
		defaultValue := generateDefaultValue(missingParam.Type)
		params.Set(missingParam.Name, fmt.Sprint(defaultValue))
		// 递归调用，直到所有参数补全
		return completeParameters(method, apiURL, params)
	}
	// 如果没有缺失参数提示，返回nil
	return nil
}
