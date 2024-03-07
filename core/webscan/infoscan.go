package webscan

import (
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"

	"strings"
)

// 代理过滤器，防止在走代理扫描时将不存活的网站由于走了代理导致状态码变成200
var filter = []string{"Burp Suite Professional"}

// response
type CheckDatas struct {
	StatusCode  int    // 状态码
	Headers     string // 响应头中的全部信息组成的字符串 key:value形式类似 Server:nginx 中间不含空格
	Title       string // 标题
	Body        []byte // 主体内容
	FaviconHash string // hash适用于fofa
}

// 接收返回数据到checkdatas
func RecvResponse(url string, client *http.Client) *CheckDatas {
	var checkdatas CheckDatas
	h := http.Header{}
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	resp, body, err := clients.NewRequest("GET", url, h, nil, 10, client)
	if err == nil && resp != nil {
		// 把响应包的内容，标题,状态码赋值给结构体
		if match := util.RegTitle.FindSubmatch(body); len(match) > 1 {
			checkdatas.Title = util.Str2UTF8(string(match[1]))
		} else {
			checkdatas.Title = ""
		}
		checkdatas.Body = body
		checkdatas.StatusCode = resp.StatusCode
		for key, value := range resp.Header {
			checkdatas.Headers += fmt.Sprintf("%v:%v", key, strings.Join(value, ""))
		}
		checkdatas.FaviconHash = FaviconHash(url, client)
	}
	for _, v := range filter {
		if checkdatas.Title == v {
			checkdatas.StatusCode = 0
		}
	}

	return &checkdatas
}

// 获取favicon hash值
func FaviconHash(url string, client *http.Client) string {
	resp, body, err := clients.NewRequest("GET", url+"favicon.ico", nil, nil, 10, client)
	if err != nil {
		return ""
	}
	if resp.StatusCode == 200 {
		return util.Mmh3Hash32(util.Base64Encode(body))
	} else {
		return ""
	}
}

func MatchRule(rule, str string) bool {
	if strings.Contains(rule, "||") && !strings.Contains(rule, "&&") && !strings.Contains(rule, "(") { // 只存在 ||
		conditions := strings.Split(rule, " || ")
		for _, condition := range conditions {
			if strings.Contains(str, condition) {
				return true
			}
		}
		return false
	} else if strings.Contains(rule, "&&") && !strings.Contains(rule, "||") && !strings.Contains(rule, "(") { // 只存在 &&
		id := 0
		conditions := strings.Split(rule, " && ")
		for _, condition := range conditions {
			if strings.Contains(str, condition) {
				id++
			}
		}
		// && 条件需要全规则匹配
		if id == len(conditions) {
			return true
		} else {
			return false
		}
		// 两种运算都存在仅支持 (abc && 456) || 123 , 不支持(wtf || qnd) && wcao
	} else if strings.Contains(rule, "||") && strings.Contains(rule, "&&") && strings.Contains(rule, "(") {
		conditions := strings.Split(rule, " || ")
		for _, condition := range conditions {
			if strings.Contains(condition, "(") { // 包含()说明是 && 运算
				condition = condition[1 : len(condition)-1] // 去除左右()
				id := 0
				cond := strings.Split(condition, " && ")
				for _, c1 := range cond {
					if strings.Contains(str, c1) {
						id++
					}
				}
				if id == len(cond) {
					return true
				}
			} else { // 仅需要 || 符一项为真
				if strings.Contains(str, condition) {
					return true
				}
			}
		}
		return false
	}
	// 不存在运算符时
	if strings.Contains(str, rule) {
		return true
	}
	return false
}
