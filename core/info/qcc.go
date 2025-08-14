package info

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/qiwentaidi/clients"
)

var (
	regCompanyInfo = regexp.MustCompile(`KeyNo":"([a-zA-Z0-9]+)","Name":\s*"([^"]+)"`)
	regTotalCount  = regexp.MustCompile(`"TotalRecords":(\d+)`)
)

// 只能查60家 ...
func FetchCompanyNamesByArea(areaName string) {
	areaNameEscape := url.QueryEscape(areaName)
	searchIndexEscape := url.QueryEscape(fmt.Sprintf(`{"scope":"%s"}`, areaName))
	baseUrl := "https://www.qcc.com/web/search?key=%s&p=%d&searchIndex=%s"

	var headers = map[string]string{
		"Host":            "www.qcc.com",
		"Connection":      "keep-alive",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		"Referer":         fmt.Sprintf("https://www.qcc.com/web/search?key=%s", areaNameEscape),
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Cookie":          "QCCSESSID=xxx;",
	}

	// 第一次请求，获取总页数
	firstUrl := fmt.Sprintf(baseUrl, areaNameEscape, 1, searchIndexEscape)
	resp, err := clients.DoRequest("GET", firstUrl, headers, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		panic(err)
	}
	body := string(resp.Body())

	// 提取总数
	totalCountMatch := regTotalCount.FindStringSubmatch(body)
	if len(totalCountMatch) < 2 {
		fmt.Println("未提取到总记录数")
		return
	}
	totalCount, _ := strconv.Atoi(totalCountMatch[1])
	totalPages := (totalCount + 19) / 20 // 每页20条，向上取整

	fmt.Printf("共 %d 条记录，约 %d 页\n", totalCount, totalPages)

	// 遍历每一页
	for page := 1; page <= totalPages; page++ {
		pageUrl := fmt.Sprintf(baseUrl, areaNameEscape, page, searchIndexEscape)
		resp, err := clients.DoRequest("GET", pageUrl, headers, nil, 10, clients.NewRestyClient(nil, true))
		if err != nil {
			fmt.Printf("请求第 %d 页失败: %v\n", page, err)
			continue
		}

		pageBody := string(resp.Body())
		matches := regCompanyInfo.FindAllStringSubmatch(pageBody, -1)

		if len(matches) == 0 {
			fmt.Printf("第 %d 页无公司信息\n", page)
			continue
		}

		fmt.Printf("第 %d 页提取到的企业数据：\n", page)
		for _, match := range matches {
			keyNo := match[1]
			name := match[2]
			fmt.Printf("KeyNo: %s, Name: %s\n", keyNo, name)
		}

		// 可选：加延时防止被封
		time.Sleep(1 * time.Second)
	}
}
