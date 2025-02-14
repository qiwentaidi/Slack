package core

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	// 设置目标目录
	dir := "./configs" // 替换为实际的配置文件目录路径
	tools := &Tools{}
	// 获取所有文件的统计结果
	results := tools.NacosCategoriesExtract(dir)

	// 输出总结果
	fmt.Println("Final Results:")
	for _, result := range results {
		fmt.Printf("File: %s\n", result.Name)
		fmt.Printf("  Auth (账号密码): %d\n", result.NodeInfo.Auth)
		fmt.Printf("  OSS: %d\n", result.NodeInfo.OSS)
		fmt.Printf("  Database: %d\n", result.NodeInfo.Database)
	}
}
