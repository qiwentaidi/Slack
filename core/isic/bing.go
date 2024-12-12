package isic

import (
	"bytes"
	"net/http"
	"slack-wails/lib/clients"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 根据网络环境，选择www.bing.com 还是cn.bing.com

func chooseBingEnvironment() (string, error) {
	resp, _, err := clients.NewSimpleGetRequest("https://www.bing.com", clients.NotFollowClient())
	if err != nil || resp == nil {
		return "", err
	}
	if resp.StatusCode == http.StatusFound {
		return "https://cn.bing.com", nil
	}
	return "https://www.bing.com", nil
}

func BingSearch(query string) (string, error) {
	bingUrl, err := chooseBingEnvironment()
	if err != nil {
		return "", err
	}
	link := bingUrl + "/search?q=site%3a+" + strings.TrimSpace(query)
	_, body, err := clients.NewSimpleGetRequest(link, clients.DefaultClient())
	if err != nil {
		return "请求失败", err
	}
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
