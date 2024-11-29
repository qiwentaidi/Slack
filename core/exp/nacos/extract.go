package nacos

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slack-wails/lib/structs"
	"strings"
)

// 定义关键词分类
var categories = map[string][]string{
	"Auth":     {"username", "password"},
	"OSS":      {"accesskey", "secret"},
	"Database": {"jdbc", "redis", "elasticsearch", "database", "mongo", "mssql", "mysql", "oracle", "postgres", "sqlserver"},
}

// 统计单个文件中每个类别关键词的出现次数
func countKeywordsInFile(filePath string) (structs.NacosNode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return structs.NacosNode{}, err
	}
	defer file.Close()

	// 初始化统计信息
	nodeInfo := structs.NacosNode{}

	// 逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text()) // 转换为小写
		for category, keywords := range categories {
			for _, keyword := range keywords {
				if strings.Contains(line, keyword) {
					switch category {
					case "Auth":
						nodeInfo.Auth++
					case "OSS":
						nodeInfo.OSS++
					case "Database":
						nodeInfo.Database++
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return structs.NacosNode{}, err
	}
	return nodeInfo, nil
}

// 遍历目录并统计每个文件的关键词出现次数，返回结果数组
func ProcessDirectory(dir string) []structs.NacosConfig {
	var results []structs.NacosConfig

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .yaml 或 .yml 文件
		if !info.IsDir() {
			nodeInfo, err := countKeywordsInFile(path)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", path, err)
				return nil
			}

			// 检查是否有关键词匹配
			if nodeInfo.Auth > 0 || nodeInfo.OSS > 0 || nodeInfo.Database > 0 {
				// 添加结果到数组中
				results = append(results, structs.NacosConfig{
					Name:     path,
					NodeInfo: nodeInfo,
				})
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	return results
}
