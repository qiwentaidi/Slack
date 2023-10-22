package main

import (
	"fmt"
	"os"
	"slack/common/proxy"
	"slack/lib/util"
	"slack/plugins/webscan"

	"gopkg.in/yaml.v2"
)

func main() {
	target := "https://localhost:11123"
	client := proxy.DefaultClient()
	yamlData, err := os.ReadFile("./config/webfinger.yaml")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	RuleData := make(map[string]map[string]string)
	err = yaml.Unmarshal(yamlData, &RuleData)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	var fingerprints []string
	var matched bool
	data := webscan.RecvResponse(target, client)
	if data.StatusCode != 0 { // 响应正常
		for name, ruleType := range RuleData {
			for types, rule := range ruleType {
				switch types {
				case "header":
					matched = webscan.MatchRule(rule, data.Headers)
				case "iconhash":
					matched = webscan.MatchRule(rule, data.FaviconHash)
				case "title":
					matched = webscan.MatchRule(rule, data.Title)
				default:
					matched = webscan.MatchRule(rule, string(data.Body))
				}
				// 每个应用有一组指纹,一组指纹中只需要匹配一个即可
				if matched {
					fingerprints = append(fingerprints, name)
					matched = false
					continue
				}
			}
		}
		fingerprints = util.RemoveDuplicates[string](fingerprints)
		fmt.Printf("fingerprints: %v\n", fingerprints) // 观察指纹是否可以被正确输出
	}
}
