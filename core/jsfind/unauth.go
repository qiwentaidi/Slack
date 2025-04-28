package jsfind

import (
	"errors"
	"net/http"
	"net/url"
	"slack-wails/lib/clients"
	"strings"
	"unicode"
)

// APIRequest 结构体表示 API 请求信息
type APIRequest struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Params  url.Values        `json:"params"`
}

// 发送请求测试未授权访问
func testUnauthorizedAccess(homeBody string, apiReq APIRequest, authentication []string) (bool, string, error) {
	var requestURL string
	var requestBody *strings.Reader

	if apiReq.Method == http.MethodGet {
		requestURL = apiReq.URL + "?" + apiReq.Params.Encode()
		requestBody = strings.NewReader("")
	} else {
		requestURL = apiReq.URL
		requestBody = strings.NewReader(apiReq.Params.Encode())
		if strings.Contains(apiReq.Headers["Content-Type"], "application/json") {
			requestBody = strings.NewReader("{}")
		}
	}

	resp, err := clients.DoRequest(apiReq.Method, requestURL, apiReq.Headers, requestBody, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return false, "", err
	}
	body := string(resp.Body())

	// 1. HTTP状态码异常，直接返回
	if resp.StatusCode() == 401 || resp.StatusCode() == 403 {
		return false, "", nil
	}

	// 2. 页面内容完全相同，排除
	if homeBody == body {
		return false, "", nil
	}

	// 3. 页面相似度检查
	similarity := jaccardSimilarity(homeBody, body)
	if similarity >= 0.9 {
		return false, "", errors.New("页面内容相似度超过90%")
	}

	// 4. 关键词判断
	for _, auth := range authentication {
		if strings.Contains(body, auth) {
			return false, "", nil
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
