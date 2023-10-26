package webscan

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slack/common"
	"slack/common/client"
	"slack/common/logger"
	"slack/common/proxy"
	"slack/lib/util"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v2"
)

// response
type CheckDatas struct {
	StatusCode  int    // 状态码
	Headers     string // 响应头中的全部信息组成的字符串 key:value形式类似 Server:nginx 中间不含空格
	Title       string // 标题
	Body        []byte // 主体内容
	FaviconHash string // hash适用于fofa
}

// 指纹扫描
func FingerScan(targets []string, finger2poc bool, progress *widget.Label) {
	var (
		wg             sync.WaitGroup
		count          = len(targets)
		id, progressid uint32
	)
	common.ScanResult = common.ScanResult[:1]
	limiter := make(chan bool, 50) // 限制协程数量
	client := client.DefaultClient()
	if common.Profile.Proxy.Enable {
		client = proxy.SelectProxy(&common.Profile)
	}
	yamlData, err := os.ReadFile("./config/webfinger.yaml")
	if err != nil {
		logger.Error(err)
	}
	RuleData := make(map[string]map[string]string)
	err = yaml.Unmarshal(yamlData, &RuleData)
	if err != nil {
		logger.Error(("Failed to unmarshal YAML: " + err.Error()))
	}
	for _, t := range targets {
		wg.Add(1)
		limiter <- true

		go func(url string) {
			var fingerprints []string
			var matched bool // 判断指纹匹配状态
			defer wg.Done()
			data := RecvResponse(url, client)
			if data.StatusCode != 0 { // 响应正常
				for name, ruleType := range RuleData {
					for types, rule := range ruleType {
						switch types {
						case "header":
							matched = MatchRule(rule, data.Headers)
						case "iconhash":
							matched = MatchRule(rule, data.FaviconHash)
						case "title":
							matched = MatchRule(rule, data.Title)
						default:
							matched = MatchRule(rule, string(data.Body))
						}
						// 每个应用有一组指纹,一组指纹中只需要匹配一个即可
						if matched {
							fingerprints = append(fingerprints, name)
							matched = false
							continue
						}
					}
				}
				atomic.AddUint32(&id, 1)
				fingerprints = util.RemoveDuplicates[string](fingerprints)
				if finger2poc { // 如果开启指纹poc扫描则加入待扫描的目标
					common.UrlFingerMap[url] = append(common.UrlFingerMap[url], fingerprints...)
				}
				common.ScanResult = append(common.ScanResult, []string{fmt.Sprintf("%d", id), url, fmt.Sprintf("%v", data.StatusCode), fmt.Sprintf("%v", len(data.Body)), data.Title, strings.Join(fingerprints, " | ")})
			}
			atomic.AddUint32(&progressid, 1)
			progress.SetText(fmt.Sprintf("%d/%d", progressid, count))
			<-limiter
		}(t)

	}
	wg.Wait()
}

// 接收返回数据到checkdatas
func RecvResponse(url string, client *http.Client) *CheckDatas {
	var checkdatas CheckDatas
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	if err != nil {
		logger.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), common.DefaultWebTimeout*time.Second)
	defer cancel()
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		logger.Error(err)
	}
	if resp != nil && resp.StatusCode != 302 { // 过滤重定向次数过多的
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
		}
		defer resp.Body.Close()
		// 把响应包的内容，标题,状态码赋值给结构体
		re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
		if match := re.FindSubmatch(body); len(match) > 1 {
			checkdatas.Title = string(match[1])
		} else {
			checkdatas.Title = ""
		}
		checkdatas.Body = body
		checkdatas.StatusCode = resp.StatusCode
		for key, value := range resp.Header {
			checkdatas.Headers += fmt.Sprintf("%v:%v", key, strings.Join(value, ""))
		}
		checkdatas.FaviconHash = FaviconHash(url, client, ctx)
	}
	return &checkdatas
}

// 获取favicon hash值
func FaviconHash(url string, client *http.Client, ctx context.Context) string {
	r, _ := http.NewRequest("GET", url+"/favicon.ico", nil)
	resp, err := client.Do(r.WithContext(ctx))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		b, err1 := io.ReadAll(resp.Body)
		if err1 != nil {
			return ""
		} else {
			return util.Mmh3Hash32(util.Base64Encode(b))
		}
	} else {
		return ""
	}
}

func MatchRule(rule, str string) bool {
	if strings.Contains(rule, "||") && !strings.Contains(rule, "&&") && !strings.Contains(rule, "(") { // 只存在 ||
		conditions := strings.Split(rule, " || ")
		for _, condition := range conditions {
			if strings.Contains(str, condition) {
				return true
			}
		}
		return false
	} else if !strings.Contains(rule, "||") && strings.Contains(rule, "&&") && !strings.Contains(rule, "(") { // 只存在 &&
		id := 0
		conditions := strings.Split(rule, " && ")
		for _, condition := range conditions {
			if strings.Contains(str, condition) {
				id++
			}
		}
		// && 条件需要全规则匹配
		if id == len(conditions) {
			return true
		} else {
			return false
		}
		// 两种运算都存在仅支持 (abc && 456) || 123 , 不支持(wtf || qnd) && wcao
	} else if strings.Contains(rule, "||") && strings.Contains(rule, "&&") && strings.Contains(rule, "(") {
		conditions := strings.Split(rule, " || ")
		for _, condition := range conditions {
			if strings.Contains(condition, "(") { // 包含()说明是 && 运算
				condition = condition[1 : len(condition)-1] // 去除左右()
				id := 0
				cond := strings.Split(condition, " && ")
				for _, c1 := range cond {
					if strings.Contains(str, c1) {
						id++
					}
				}
				if id == len(cond) {
					return true
				}
			} else { // 仅需要 || 符一项为真
				if strings.Contains(str, condition) {
					return true
				}
			}
		}
		return false
	}
	// 不存在运算符时
	if strings.Contains(str, rule) {
		return true
	}
	return false
}
