package isic

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// BingResult 表示单条Bing搜索结果
type BingResult struct {
	Title   string
	URL     string
	Snippet string
}

// BingSearch 使用 chromedp 模拟浏览器搜索 site:domain
func GoogleHackerBingSearch(query string) ([]BingResult, int, error) {
	// 设置 Chrome 执行选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-background-timer-throttling", false),
		chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Flag("ignore-certificate-errors", true),
	)

	// 创建执行上下文
	allocatorCtx, chromedpCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer chromedpCancel()

	// 创建主上下文
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()
	searchURL := "https://www.bing.com/search?q=" + url.QueryEscape(query)

	var htmlContent string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(`#b_content`, chromedp.ByID),
		chromedp.OuterHTML("body", &htmlContent),
	); err != nil {
		return nil, 0, err
	}
	results, _, totalNum := parseBingHTML(htmlContent)

	// 解析 HTML，提取搜索结果
	return results, totalNum, nil
}
func parseBingHTML(html string) ([]BingResult, string, int) {
	var results []BingResult
	var totalText string
	var totalNum int

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return results, "", 0
	}

	// ✅ 提取结果总数文本（支持中英文）
	if sel := doc.Find("#b_tween_searchResults > span"); sel.Length() > 0 {
		totalText = strings.TrimSpace(sel.Text())
	} else if sel := doc.Find("#b_tween .sb_count"); sel.Length() > 0 {
		totalText = strings.TrimSpace(sel.Text())
	}

	// ✅ 提取数字
	if totalText != "" {
		re := regexp.MustCompile(`[\d,]+`)
		if match := re.FindString(totalText); match != "" {
			match = strings.ReplaceAll(match, ",", "")
			fmt.Sscanf(match, "%d", &totalNum)
		}
	}

	// ✅ 提取搜索结果列表
	doc.Find("li.b_algo").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := strings.TrimSpace(s.Find("h2").Text())
		url, _ := s.Find("h2 a").Attr("href")
		snippet := strings.TrimSpace(s.Find(".b_caption p").Text())
		if title != "" && url != "" {
			results = append(results, BingResult{
				Title:   title,
				URL:     url,
				Snippet: snippet,
			})
		}
		return len(results) < 10
	})

	return results, totalText, totalNum
}
