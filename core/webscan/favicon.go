package webscan

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"net/url"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/twmb/murmur3"
)

var (
	iconDesktopRels = []string{"icon", "shortcut icon"}         // 桌面端 Logo 优先匹配
	iconMobileRels  = []string{"apple-touch-icon", "mask-icon"} // 移动｜其他端 Logo 其次
)

// 获取favicon Mmh3Hash32 和 MD5值
func FaviconHash(u *url.URL, headers map[string]string, client *http.Client) (string, string) {
	_, body, err := clients.NewRequest("GET", u.String(), headers, nil, 10, true, client)
	if err != nil {
		return "", ""
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", ""
	}
	iconLink := parseIcons(doc)[0]
	var finalLink string
	// 如果是完整的链接，则直接请求
	if strings.HasPrefix(iconLink, "http") {
		finalLink = iconLink
		// 如果为 // 开头采用与网站同协议
	} else if strings.HasPrefix(iconLink, "//") {
		finalLink = u.Scheme + ":" + iconLink
	} else {
		finalLink = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, iconLink)
	}
	resp, body, err := clients.NewRequest("GET", finalLink, headers, nil, 10, true, client)
	if err == nil && resp.StatusCode == 200 {
		hasher := md5.New()
		hasher.Write(body)
		sum := hasher.Sum(nil)
		return Mmh3Hash32(body), hex.EncodeToString(sum)
	}
	return "", ""
}

func GetFaviconFullLink(u *url.URL, client *http.Client) (string, error) {
	_, body, err := clients.NewSimpleGetRequest(u.String(), client)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", errors.New("goquery failed to parse " + u.String() + " content")
	}
	iconLink := parseIcons(doc)[0]
	var finalLink string
	// 如果是完整的链接，则直接请求
	if strings.HasPrefix(iconLink, "http") {
		finalLink = iconLink
		// 如果为 // 开头采用与网站同协议
	} else if strings.HasPrefix(iconLink, "//") {
		finalLink = u.Scheme + ":" + iconLink
	} else {
		finalLink = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, iconLink)
	}
	return finalLink, nil
}

// parseIcons 解析HTML文档head中的<link>标签中rel属性包含icon信息的href链接
func parseIcons(doc *goquery.Document) []string {
	var icons []string
	// 桌面端
	doc.Find("head link").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			// 匹配ICON链接
			if rel, exists := s.Attr("rel"); exists && util.ArrayContains(rel, iconDesktopRels) {
				icons = append(icons, href)
			}
		}
	})
	// 移动端
	if len(icons) == 0 {
		doc.Find("head link").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				// 匹配ICON链接
				if rel, exists := s.Attr("rel"); exists && util.ArrayContains(rel, iconMobileRels) {
					icons = append(icons, href)
				}
			}
		})
	}

	// 找不到自定义icon链接就使用默认的favicon地址
	if len(icons) == 0 {
		icons = append(icons, "favicon.ico")
	}

	return icons
}

// Reference: https://github.com/Becivells/iconhash

// Mmh3Hash32 计算 mmh3 hash
func Mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	_, err := h32.Write(base64Encode(raw))
	if err == nil {
		return fmt.Sprint(int32(h32.Sum32()))
	} else {
		return "0"
	}
}

// base64 encode
func base64Encode(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}
