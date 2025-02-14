package core

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
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

// 计算文件的 MD5 哈希值
func calculateFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil)), nil
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
func (t *Tools) NacosCategoriesExtract(dir string) []structs.NacosConfig {
	var results []structs.NacosConfig
	hashSet := make(map[string]bool) // 用于存储文件哈希值，防止重复

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .yaml 或 .yml 文件
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			// 计算文件哈希值
			hash, err := calculateFileHash(path)
			if err != nil {
				fmt.Printf("Error calculating hash for file %s: %v\n", path, err)
				return nil
			}

			// 检查是否已经处理过该文件
			if hashSet[hash] {
				return nil // 跳过重复文件
			}

			// 添加哈希值到集合中
			hashSet[hash] = true

			// 分析文件关键词
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
