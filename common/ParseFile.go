package common

import (
	"bufio"
	"net/url"
	"os"
	"path"
	"slack/common/logger"

	"strings"

	"fyne.io/fyne/v2/widget"
)

const (
	Mode_Other = iota
	Mode_Url
)

// 将文件读取成每行的目标
func ParseFile(filepath string) (targets []string) {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Info("[ERR]" + filepath + "打开失败")
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" { // 去除空行
			targets = append(targets, s.Text())
		}
	}
	return targets
}

// 把TEXT文本转换为扫描目标对象,mode=1表示目标为url
func ParseTarget(text string, mode int) []string {
	var temp, targets []string
	temp = strings.Split(text, "\n") // 通过查找换行去分割每个目标
	if mode == 1 {
		for _, t := range temp {
			t = strings.ReplaceAll(strings.ReplaceAll(t, "\r", ""), " ", "") // 去除所有空格
			if t != "" {
				t = strings.TrimSuffix(t, "/") //如果末尾是/结尾则去除/
				targets = append(targets, t)
			}
		}
	} else {
		for _, t := range temp {
			t = strings.ReplaceAll(strings.ReplaceAll(t, "\r", ""), " ", "")
			if t != "" {
				targets = append(targets, t)
			}
		}
	}
	return targets
}

func ParseURLWithoutSlash(text string) (string, error) {
	if _, err := url.ParseRequestURI(text); err != nil {
		return "", err
	}
	text = strings.ReplaceAll(text, " ", "")
	if strings.HasSuffix(text, "/") {
		return text, nil
	} else {
		return text + "/", nil
	}
}

// 分辨参数是.txt字典或者是单个用户名
func ParseDict(e *widget.Entry, defaults []string) []string {
	var dicts []string
	if !e.Disabled() && e.Text != "" { // 使用自定义字典
		if path.Ext(e.Text) == ".txt" { // 文件
			dicts = append(dicts, ParseFile(e.Text)...)
		} else {
			dicts = append(dicts, e.Text)
		}
	} else {
		dicts = defaults // 默认使用内置字典
	}
	return dicts
}

func ParseDirectoryDict(filepath, old string, new []string) (dict []string) {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Info(err)
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
