package info

import (
	"net/http"
	"regexp"
	"slack-wails/lib/clients"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// 返回域名组
func Beianx(company string) []string {
	var ddddomain []string
	h := http.Header{}
	h.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	_, body, err := clients.NewRequest("GET", "https://www.beianx.cn/search/"+company, h, nil, 10, http.DefaultClient)
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	reg := regexp.MustCompile(`info/\d+`)
	links := reg.FindAllString(string(body), -1)
	for _, link := range links {
		_, b, err := clients.NewRequest("GET", "https://www.beianx.cn/"+link, nil, nil, 10, clients.DefaultClient())
		if err != nil {
			logger.NewDefaultLogger().Debug(err.Error())
		}
		reg_domain := regexp.MustCompile(`whois/.+" `)
		for _, v := range reg_domain.FindAllString(string(b), -1) {
			ddddomain = append(ddddomain, v[:len(v)-2][6:])
		}
	}
	time.Sleep(time.Second * 2)
	return ddddomain
}
