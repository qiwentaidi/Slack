package info

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"slack-wails/lib/clients"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 如果请求到开头为
var acwscv2 = ""

// 返回域名组，延时防止请求过快
func Beianx(company string) ([]string, error) {
	h := map[string]string{
		"Cookie": "acw_sc__v2=" + acwscv2,
	}
	_, body, err := clients.NewRequest("GET", "https://www.beianx.cn/search/"+company, h, nil, 10, true, http.DefaultClient)
	if err != nil && len(body) == 17132 { // 符合长度表示存在acw_sc__v2校验，需要获取acw_sc__v2的值，再次执行函数即可
		arg1 := getArg1FromHTML(string(body))
		acwscv2 = getAcwScV2(arg1)
		time.Sleep(time.Second)
		return Beianx(company)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return []string{}, errors.New("解析网站内容失败")
	}
	var domains []string
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		domain := s.Find("td").Eq(5).Find("a").Text()
		if domain != "" {
			domains = append(domains, domain)
		}
	})
	time.Sleep(time.Second)
	return domains, nil
}

// 参考https://github.com/qiwentaidi/acw_sc_v2计算cookie

// int2base function in Golang
func int2base(x int, base int) string {
	digs := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if x == 0 {
		return string(digs[0])
	}
	sign := 1
	if x < 0 {
		sign = -1
		x = -x
	}
	var digits []byte
	for x > 0 {
		digits = append([]byte{digs[x%base]}, digits...)
		x = x / base
	}
	if sign < 0 {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

// hexXor function in Golang
func hexXor(str, key string) string {
	var result string
	for i := 0; i < len(str) && i < len(key); i += 2 {
		s1, _ := intFromHex(str[i : i+2])
		s2, _ := intFromHex(key[i : i+2])
		xor := int2base(int(s1)^int(s2), 16)
		if len(xor) == 1 {
			xor = "0" + xor
		}
		result += xor
	}
	return result
}

func intFromHex(hexStr string) (int64, error) {
	return strconv.ParseInt(hexStr, 16, 64)
}

// unsbox function in Golang
func unsbox(str string) string {
	box := []int{0xf, 0x23, 0x1d, 0x18, 0x21, 0x10, 0x1, 0x26, 0xa, 0x9, 0x13, 0x1f, 0x28, 0x1b, 0x16, 0x17, 0x19, 0xd, 0x6, 0xb, 0x27, 0x12, 0x14, 0x8, 0xe, 0x15, 0x20, 0x1a, 0x2, 0x1e, 0x7, 0x4, 0x11, 0x5, 0x3, 0x1c, 0x22, 0x25, 0xc, 0x24}
	result := make([]string, len(box))
	for i, ch := range str {
		for j, b := range box {
			if b == i+1 {
				result[j] = string(ch)
			}
		}
	}
	return strings.Join(result, "")
}

func getAcwScV2(arg1 string) string {
	key := "3000176000856006061501533003690027800375"
	unsboxStr := unsbox(arg1)
	return hexXor(unsboxStr, key)
}

func getArg1FromHTML(html string) string {
	re := regexp.MustCompile(`arg1='([^']+)'`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
