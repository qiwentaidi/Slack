package jsfind

import (
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

var (
	regJS  []*regexp.Regexp
	JsLink = []string{
		"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
	}
	// UrlFind = []string{
	// 	"[\"'‘“`]\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
	// 	"=\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
	// 	"[\"'‘“`]\\s{0,6}([#,.]{0,2}/[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
	// 	"\"([-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?[/]{1}[-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?)\"",
	// 	"href\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})|action\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
	// }
	SensitiveField = []string{
		// sensitive-filed
		`((\[)?('|")?([\w]{0,10})((key)|(secret)|(token)|(config)|(auth)|(access)|(admin))([\w]{0,10})('|")?(\])?( |)(:|=)( |)('|")(.*?)('|")(|,))`,
		// username-filed
		`((|'|")(|[\w]{1,10})(([u](ser|name|sername))|(account)|((creat|updat)(|ed|or|er)(|by|on|at)))(|[\w]{1,10})(|'|")(:|=)( |)('|")(.*?)('|")(|,))`,
		// password-filed
		`((|'|")(|[\w]{1,10})([p](ass|wd|asswd|assword))(|[\w]{1,10})(|'|")(:|=)( |)('|")(.*?)('|")(|,))`,
	}
	Filter = []string{".vue", ".jpeg", ".png", ".jpg", ".ts", "gif", ".css", ".svg", ".scss"}
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
}

func ExtractJS(url string) (allJS []string) {
	_, body, err := clients.NewRequest("GET", url, nil, nil, 10, clients.DefaultClient())
	if err != nil || body == nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return
	}
	content := string(body)
	for _, reg := range regJS {
		for _, item := range reg.FindAllString(content, -1) {
			item = strings.TrimLeft(item, "=")
			item = strings.Trim(item, "\"")
			item = strings.TrimLeft(item, ".")
			if item[0:4] != "http" {
				allJS = append(allJS, item)
			}
		}
	}
	return util.RemoveDuplicates(allJS)
}

// setp 0 first need deep js
func FindInfo(url string) *FindSomething {
	var fs FindSomething
	_, body, err := clients.NewRequest("GET", url, nil, nil, 10, clients.DefaultClient())
	if err != nil || body == nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return &fs
	} else {
		content := string(body)
		// 先匹配其他信息
		urls, apis := urlInfoSeparate(util.RegLink.FindAllString(content, -1))
		fs.APIRoute = *AppendSource(url, apis)
		fs.IP_URL = *AppendSource(url, append(util.RegIP_PORT.FindAllString(content, -1), urls...))
		fs.ChineseIDCard = *AppendSource(url, util.RegIDCard.FindAllString(content, -1))
		fs.ChinesePhone = *AppendSource(url, util.RegPhone.FindAllString(content, -1))
		for _, reg := range SensitiveField {
			regSen := regexp.MustCompile(reg)
			for _, item := range regSen.FindAllString(content, -1) {
				fs.SensitiveField = append(fs.SensitiveField, InfoSource{Filed: item, Source: url})
			}
		}
	}
	return &fs
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

func urlInfoSeparate(links []string) (urls, apis []string) {
	for _, link := range links {
		link = strings.Trim(link, "\"")
		if strings.HasPrefix(link, "http") || strings.HasPrefix(link, "ws") {
			urls = append(urls, link)
		} else {
			matched := false
			for _, f := range Filter {
				if strings.HasSuffix(link, f) {
					matched = true
					break
				}
			}
			if !matched {
				apis = append(apis, link)
			}
		}
	}
	return urls, apis
}
