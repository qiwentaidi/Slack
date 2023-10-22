package info

import (
	"io"
	"net/http"
	"regexp"
	"slack/common/logger"
)

func Beianx(company string) []string {
	var ddddomain []string
	client := &http.Client{}
	r, err := http.NewRequest("GET", "https://www.beianx.cn/search/"+company, nil)
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	r.Header.Add("Cookie", "machine_str=undefined")
	if err != nil {
		logger.Debug(err)
	}
	resp, err2 := client.Do(r)
	if err2 != nil {
		logger.Debug(err)
	}
	b, _ := io.ReadAll(resp.Body)
	reg := regexp.MustCompile(`info/\d+`)
	links := reg.FindAllString(string(b), -1)
	for _, link := range links {
		r2, err := http.NewRequest("GET", "https://www.beianx.cn/"+link, nil)
		r2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
		if err != nil {
			logger.Debug(err)
		}
		resp2, err2 := client.Do(r2)
		if err2 != nil {
			logger.Debug(err)
		}
		b, _ = io.ReadAll(resp2.Body)
		reg_domain := regexp.MustCompile(`whois/.+" `)
		for _, v := range reg_domain.FindAllString(string(b), -1) {
			ddddomain = append(ddddomain, v[:len(v)-2][6:])
		}
	}
	return ddddomain
}
