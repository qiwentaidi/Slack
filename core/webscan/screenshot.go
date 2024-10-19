package webscan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slack-wails/lib/util"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

var dir = filepath.Join(util.HomeDir(), "slack", "screenshot")

// GetScreenshot 获取指定URL的屏幕截图，并保存到本地文件。
// 返回文件路径和错误，如果错误不为nil，则文件路径为空。
func GetScreenshot(url string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-background-timer-throttling", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// 创建一个浏览器实例
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	// 运行任务
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second), // 等待页面加载完成
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		return "", errors.New("无法获取屏幕截图")
	}

	// 定义保存路径
	fp := filepath.Join(dir, renameOutput(url))

	// 确保目标目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// 将截图保存到文件
	if err := os.WriteFile(fp, buf, 0644); err != nil {
		return "", errors.New("无法保存屏幕截图: " + err.Error())
	}

	return fp, nil
}

// renameOutput 用于清理URL中的非法字符，以生成合法的文件名。
func renameOutput(filename string) string {
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "___", "_")
	return filename
}
