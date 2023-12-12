package util

import (
	"bufio"
	"os"
)

// 将文件读取成每行的目标
func ParseFile(path string) ([]string, error) {
	var targets []string
	file, err := os.Open(path)
	if err != nil {
		return targets, err
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" { // 去除空行
			targets = append(targets, s.Text())
		}
	}
	return targets, nil
}
