package jsfind

import (
	"context"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
	"sync"
)

var (
	regJS     []*regexp.Regexp
	regFilter []*regexp.Regexp
	JsLink    = []string{
		"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
	}
	Phone          = regexp.MustCompile(`[^\w]((?:(?:\+|00)86)?1(?:(?:3[\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\d])|(?:9[189]))\d{8})[^\w]`)
	IDCard         = regexp.MustCompile(`[^0-9]((\d{8}(0\d|10|11|12)([0-2]\d|30|31)\d{3}$)|(\d{6}(18|19|20)\d{2}(0[1-9]|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)))[^0-9]`)
	SensitiveField = []string{
		// sensitive-filed
		`((\[)?('|")?([\w]{0,10})((key)|(secret)|(token)|(config)|(auth)|(access)|(admin))([\w]{0,10})('|")?(\])?( |)(:|=)( |)('|")(.*?)('|")(|,))`,
		// username-filed
		`((|'|")(|[\w]{1,10})(([u](ser|name|sername))|(account)|((creat|updat)(|ed|or|er)(|by|on|at)))(|[\w]{1,10})(|'|")(:|=)( |)('|")(.*?)('|")(|,))`,
		// password-filed
		`((|'|")(|[\w]{1,10})([p](ass|wd|asswd|assword))(|[\w]{1,10})(|'|")(:|=)( |)('|")(.*?)('|")(|,))`,
	}
	Filter = []string{".vue", ".jpeg", ".png", ".jpg", ".ts", ".gif", ".css", ".svg", ".scss"}
)

type InfoSource struct {
	Filed  string
	Source string
}

type FindSomething struct {
	JS             []InfoSource
	APIRoute       []InfoSource
	IP_URL         []InfoSource
	ChineseIDCard  []InfoSource
	ChinesePhone   []InfoSource
	SensitiveField []InfoSource
}

func init() {
	for _, reg := range JsLink {
		regJS = append(regJS, regexp.MustCompile(reg))
	}
	for _, f := range Filter {
		regFilter = append(regFilter, regexp.MustCompile(f))
	}
}

func ExtractJS(ctx context.Context, url string) (allJS []string) {
	_, body, err := clients.NewSimpleGetRequest(url, clients.DefaultClient())
	if err != nil || body == nil {
		gologger.Debug(ctx, err)
		return
	}
	content := string(body)
	for _, reg := range regJS {
		for _, item := range reg.FindAllString(content, -1) {
			item = strings.TrimLeft(item, "=")
			item = strings.Trim(item, "\"")
			item = strings.Trim(item, "'")
			item = strings.TrimLeft(item, ".")
			allJS = append(allJS, item)
		}
	}
	regJSSecond := regexp.MustCompile(`(?i)<script[^>]+src=["']([^"']+\.js)["']?`)
	for _, item := range regJSSecond.FindAllStringSubmatch(content, -1) {
		allJS = append(allJS, item[1])
	}
	return util.RemoveDuplicates(allJS)
}

// setp 0 first need deep js
func FindInfo(ctx context.Context, url string, limiter chan bool, wg *sync.WaitGroup) *FindSomething {
	defer wg.Done()
	var fs FindSomething
	_, body, err := clients.NewSimpleGetRequest(url, clients.DefaultClient())
	if err != nil || body == nil {
		gologger.Debug(ctx, err)
		return &fs
	} else {
		content := string(body)
		// 先匹配其他信息
		urls, apis, js := urlInfoSeparate(util.RegLink.FindAllString(content, -1))
		fs.JS = *AppendSource(url, js)
		fs.APIRoute = *AppendSource(url, apis)
		fs.IP_URL = *AppendSource(url, append(util.RegIP_PORT.FindAllString(content, -1), urls...))
		fs.ChineseIDCard = *AppendSource(url, IDCard.FindAllString(content, -1))
		fs.ChinesePhone = *AppendSource(url, cleanPhoneNumber(Phone.FindAllString(content, -1)))
		for _, reg := range SensitiveField {
			regSen := regexp.MustCompile(reg)
			for _, item := range regSen.FindAllString(content, -1) {
				fs.SensitiveField = append(fs.SensitiveField, InfoSource{Filed: item, Source: url})
			}
		}
	}
	<-limiter
	return &fs
}

// 去除手机号中的其他字符
func cleanPhoneNumber(phone []string) (news []string) {
	// 定义只保留数字和+的正则表达式
	cleanRegex := regexp.MustCompile(`[^\d+]`)
	for _, p := range phone {
		news = append(news, cleanRegex.ReplaceAllString(p, ""))
	}
	return
}

func AppendSource(source string, filed []string) *[]InfoSource {
	is := []InfoSource{}
	for _, f := range filed {
		is = append(is, InfoSource{Filed: f, Source: source})
	}
	return &is
}

func RemoveDuplicatesInfoSource(iss []InfoSource) []InfoSource {
	encountered := map[string]bool{}
	result := []InfoSource{}
	for _, is := range iss {
		if !encountered[is.Filed] {
			encountered[is.Filed] = true
			result = append(result, is)
		}
	}
	return result
}

func urlInfoSeparate(links []string) (urls, apis, js []string) {
	for _, link := range links {
		link = strings.Trim(link, "\"")
		link = strings.Trim(link, "'")
		if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "ws") {
			urls = append(urls, link)
		} else {
			matched := false
			for _, r := range regFilter {
				if strings.Contains(link, ".js") {
					js = append(js, link)
					matched = true
					break
				}
				if r.MatchString(link) {
					matched = true // 匹配到过滤器后缀需要屏蔽
					break
				}
			}
			if !matched {
				apis = append(apis, link)
			}
		}
	}
	return urls, apis, js
}

func FilterExt(iss []InfoSource) (news []InfoSource) {
	for _, link := range iss {
		matched := false
		for _, r := range regFilter {
			if r.MatchString(link.Filed) {
				matched = true
				break
			}
		}
		if !matched {
			news = append(news, link)
		}

	}
	return news
}
