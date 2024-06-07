package core

import (
	"fmt"
	"os"
	"slack-wails/lib/qqwry"
	"slack-wails/lib/util"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

type SubdomainResult struct {
	Subdomain string
	Cname     []string
	Ips       []string
	Notes     string
}

var (
	IPResolved map[string]int
	mutex      sync.Mutex
	database   *qqwry.QQwry
	cdndata    map[string][]string
	onec       sync.Once
)

// 初始化IP纯真库
func InitQqwry(qqwryFile string) {
	fs, err := os.OpenFile(qqwryFile, os.O_RDONLY, 0777)
	if err != nil {
		logger.NewDefaultLogger().Debug("qqwry open err:" + err.Error())
		return
	}
	if d, err := qqwry.NewQQwryFS(fs); err != nil {
		logger.NewDefaultLogger().Debug("qqwry init err:" + err.Error())
		return
	} else {
		database = d
	}
}

// 采用递归判断暴破层级
func BurstSubdomain(subdomains string, timeout int, qqwryFile, cdnFile string) *SubdomainResult {
	onec.Do(func() {
		cdndata = ReadCDNFile(cdnFile)
		InitQqwry(qqwryFile)
	})
	var sr SubdomainResult
	addrs, cnames, err := Resolution(subdomains, timeout)
	if err == nil {
		sr.Cname = cnames
	outloop:
		for _, cdns := range cdndata {
			for _, cdn := range cdns {
				for _, cname := range cnames {
					if strings.Contains(cname, cdn) { // 识别到cdn
						sr.Notes = fmt.Sprintf("在CNAME中识别到CDN字段%v", cdn)
						break outloop
					} else if strings.Contains(cname, "cdn") {
						sr.Notes = fmt.Sprintf("在CNAME %v中检测到cdn关键字", cname)
						break outloop
					}
				}
			}
		}
		for _, ip := range addrs {
			flag, result, _ := FindWithIP(ip)
			if flag {
				addrs = util.RemoveElement(addrs, ip)
				sr.Notes = fmt.Sprintf("%v在IP纯真库中识别到cdn字段%v", ip, result)
				continue
			}
			mutex.Lock()
			IPResolved[ip]++
			if IPResolved[ip] > 5 { // 解析到该IP5次以上加入黑名单
				addrs = util.RemoveElement(addrs, ip)
			}
			mutex.Unlock()
		}
		sr.Subdomain = subdomains
		sr.Ips = addrs
	}
	return &sr
}

func FindWithIP(query string) (bool, string, error) {
	result, err := Find(query)
	if strings.Contains(result, "CDN") {
		return true, result, err
	}
	return false, "", err
}

func Find(query string) (string, error) {
	result, err := database.Find(query)
	if err != nil {
		return "", err
	}
	if strings.Contains(result.String(), "对方和您在同一内部网") {
		return "", err
	}
	return result.String(), err
}
