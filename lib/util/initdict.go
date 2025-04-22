package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadDirsearchDict(filepath, old string, new []string) (dict []string) {
	file, _ := os.Open(filepath)
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

func ReadLine(filepath string) (dict []string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" { // 去除空行
			dict = append(dict, s.Text())
		}
	}
	return RemoveDuplicates(dict)
}
