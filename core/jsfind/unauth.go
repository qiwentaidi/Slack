package jsfind

import (
	"errors"
	"net/url"
	"slack-wails/lib/clients"
	"strings"
	"unicode"
)

// 一些提示需要鉴权时的返回信息
var authentication = []string{"token不能为空", "令牌不能为空", "令牌已过期", "Unauthorized", "Access Denied"}

// APIRequest 结构体表示 API 请求信息
type APIRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Params  url.Values        `json:"params"`
}

// 发送请求测试未授权访问
func testUnauthorizedAccess(homeBody string, apiReq APIRequest) (bool, []byte, error) {
	var requestURL string
	var requestBody *strings.Reader

	if apiReq.Method == "GET" {
		// 将参数附加到 URL 查询字符串
		requestURL = apiReq.URL + "?" + apiReq.Params.Encode()
		requestBody = strings.NewReader("")
	} else {
		// 其他请求方式（POST/PUT 等），参数作为请求体
		requestURL = apiReq.URL
		requestBody = strings.NewReader(apiReq.Params.Encode())
	}

	// 发送请求
	resp, body, err := clients.NewRequest(apiReq.Method, requestURL, apiReq.Headers, requestBody, 10, false, clients.NewHttpClient(nil, true))
	if err != nil || resp == nil {
		return false, nil, err
	}

	// 判断是否为未授权访问
	if resp.StatusCode == 401 {
		return false, nil, nil
	}
	if homeBody == string(body) {
		return false, nil, nil
	}
	// // 计算文本相似度
	similarity := jaccardSimilarity(homeBody, string(body))
	if similarity >= 0.9 {
		return false, nil, errors.New("页面内容相似度超过90%")
	}

	// 解析返回内容，检查是否包含未授权提示信息
	for _, auth := range authentication {
		if strings.Contains(string(body), auth) {
			return false, nil, nil
		}
	}

	return true, body, nil
}

// 将文本分割成 shingle（n-gram 片段），用于计算相似度
func tokenize(text string, n int) map[string]struct{} {
	// 预处理文本，去除空格和标点
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return ' '
	}, text)
	words := strings.Fields(cleaned)

	// 生成 n-gram 片段
	tokens := make(map[string]struct{})
	for i := 0; i < len(words)-n+1; i++ {
		token := strings.Join(words[i:i+n], " ")
		tokens[token] = struct{}{}
	}
	return tokens
}

// 计算 Jaccard 相似度
func jaccardSimilarity(text1, text2 string) float64 {
	set1 := tokenize(text1, 3) // 3-gram
	set2 := tokenize(text2, 3)

	// 计算交集大小
	intersection := 0
	for token := range set1 {
		if _, exists := set2[token]; exists {
			intersection++
		}
	}

	// 计算并集大小
	union := len(set1) + len(set2) - intersection

	if union == 0 {
		return 0.0
	}
	return float64(intersection) / float64(union)
}
