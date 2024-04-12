package webscan

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func TestInfoscan(t *testing.T) {
	// 使用http.Get获取网页内容
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// 将body转换为字符串
	bodyStr := string(body)

	// 创建CEL环境
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("body", decls.String),
		),
	)
	if err != nil {
		fmt.Printf("Error creating CEL environment: %v\n", err)
		return
	}

	// 编译CEL表达式
	ast, iss := env.Compile(`body.contains('baidu')`)
	if iss.Err() != nil {
		fmt.Printf("Error compiling CEL expression: %v\n", iss.Err())
		return
	}

	// 创建CEL程序
	program, err := env.Program(ast)
	if err != nil {
		fmt.Printf("Error creating CEL program: %v\n", err)
		return
	}

	// 执行CEL程序
	out, _, err := program.Eval(map[string]interface{}{
		"body": bodyStr,
	})
	if err != nil {
		fmt.Printf("Error evaluating CEL program: %v\n", err)
		return
	}

	fmt.Println("Check if body contains 'baidu':", out)
}
