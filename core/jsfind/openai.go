package jsfind

// type UnauthorizedChecker struct {
// 	client   *openai.Client
// 	messages []openai.ChatCompletionMessageParamUnion
// 	model    string
// 	ctx      context.Context
// }

// // 定义最大并发数
// const maxConcurrentRequests = 2

// // 用于控制并发的信号通道
// var sem = make(chan struct{}, maxConcurrentRequests)

// // NewChecker 初始化
// func NewChecker(apiKey, baseURL, model string) *UnauthorizedChecker {
// 	client := openai.NewClient(
// 		option.WithAPIKey(apiKey),
// 		option.WithBaseURL(baseURL),
// 		option.WithHTTPClient(&http.Client{
// 			Timeout: 5 * time.Second,
// 		}),
// 	)
// 	return &UnauthorizedChecker{
// 		client: &client,
// 		model:  model,
// 		ctx:    context.Background(),
// 		messages: []openai.ChatCompletionMessageParamUnion{
// 			openai.UserMessage("你是一个未授权访问检测助手，用户会提供一些接口的返回数据，你需要判断接口的返回数据来判断接口是否需要用户登录。回答用true/false，除了true或false不要输出任何其他内容。"),
// 		},
// 	}
// }

// // Check 输入响应文本，返回是否未授权访问
// func (c *UnauthorizedChecker) Check(responseBody string) (bool, error) {
// 	// 获取并发许可
// 	sem <- struct{}{}
// 	defer func() {
// 		// 释放并发许可
// 		<-sem
// 	}()
// 	// 加入新的用户输入
// 	c.messages = append(c.messages, openai.UserMessage(responseBody))
// 	chatCompletion, err := c.client.Chat.Completions.New(c.ctx, openai.ChatCompletionNewParams{
// 		Messages: c.messages,
// 		Model:    c.model,
// 	})
// 	if err != nil {
// 		return false, err
// 	}

// 	// 安全处理，防止下标越界
// 	if len(chatCompletion.Choices) == 0 {
// 		return false, fmt.Errorf("no choices returned from model")
// 	}

// 	reply := strings.TrimSpace(chatCompletion.Choices[0].Message.Content)

// 	// 只判断是否返回 true
// 	return strings.ToLower(reply) == "true", nil
// }
