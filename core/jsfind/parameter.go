package jsfind

import "regexp"

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
		return ""
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
