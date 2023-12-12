package core

import (
	"bufio"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

const (
	defaultDict = "./config/dirsearch/"
)

func LoadDefaultDict(filename, old string, new []string) (dict []string) {
	file, err := os.Open(defaultDict + filename)
	if err != nil {
		logger.NewFileLogger(filename).Debug(err.Error())
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" { // 去除空行
			if len(new) > 0 {
				if strings.Contains(s.Text(), old) { // 如何新数组不为空,将old字段替换成new数组
					for _, n := range new {
						dict = append(dict, strings.ReplaceAll(s.Text(), old, n))
					}
				} else {
					dict = append(dict, s.Text())
				}
			} else {
				if !strings.Contains(s.Text(), old) {
					dict = append(dict, s.Text())
				}
			}
		}
	}
	return dict
}
