package isic

import (
	"bytes"
	"fmt"
	"slack-wails/lib/clients"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func BingSearch(query string) (string, error) {
	link := "https://www.bing.com/search?q=site%3a+" + strings.TrimSpace(query)
	_, body, err := clients.NewSimpleGetRequest(link, clients.DefaultClient())
	if err != nil {
		return "请求失败", err
	}
	fmt.Printf("string(body): %v\n", string(body))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "解析网站失败", err
	}
	var counts string
	doc.Find("span.sb_count").Each(func(i int, s *goquery.Selection) {
		counts = s.Text()
	})
	return counts, nil
}
