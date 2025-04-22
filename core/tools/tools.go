package core

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"slack-wails/lib/util"
	"strconv"
)

type Tools struct{}

func (t *Tools) IsRoot() bool {
	if runtime.GOOS == "windows" {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	} else {
		return os.Getuid() == 0
	}
}

func (t *Tools) GOOS() string {
	return runtime.GOOS
}
func (t *Tools) IPParse(ipList []string) []string {
	return util.ParseIPs(ipList)
}

func (t *Tools) PortParse(text string) []int {
	return util.ParsePort(text)
}

func (t *Tools) ExtractIP(text string) string {
	var result string
	var IP_analysis = make(map[string]int)
	result += "---提取IP资产---\n"
	deupIPs := util.RemoveDuplicates(util.RegIP.FindAllString(text, -1))
	for _, ip := range deupIPs {
		result += ip + "\n"
		ip = ip[:len(ip)-len(path.Ext(ip))]
		IP_analysis[ip+".0"]++
	}
	result += "共计提取到IP资产" + strconv.Itoa(len(deupIPs)) + "个\n"
	result += "\n\n\n---提取C段资产---\n"
	for _, p := range util.SortMap(IP_analysis) {
		result += fmt.Sprintf("%v/24(%v)\n", p.Key, p.Value)
	}
	return result
}
