package jsfind

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

// 加载网页全部内容，包括异步的js文件
func asyncLoader(url string) {
	// 创建 context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 用于存储页面内容
	var pageContent string

	// 运行任务
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#target-element`, chromedp.ByID), // 等待目标元素加载完成
		chromedp.OuterHTML("html", &pageContent),               // 获取页面内容
	)
	if err != nil {
		fmt.Println("chromedp 错误:", err)
		return
	}

	fmt.Println("页面内容:", pageContent)
}
