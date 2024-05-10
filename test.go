package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"
// )

// type NucleiResult []struct {
// 	TemplateID   string `json:"template-id"`
// 	TemplatePath string `json:"template-path"`
// 	Info         struct {
// 		Name        string   `json:"name"`
// 		Author      []string `json:"author"`
// 		Tags        []string `json:"tags"`
// 		Description string   `json:"description"`
// 		Severity    string   `json:"severity"`
// 		Metadata    struct {
// 			MaxRequest int `json:"max-request"`
// 		} `json:"metadata"`
// 	} `json:"info"`
// 	Type             string    `json:"type"`
// 	Host             string    `json:"host"`
// 	Port             string    `json:"port"`
// 	Scheme           string    `json:"scheme"`
// 	URL              string    `json:"url"`
// 	MatchedAt        string    `json:"matched-at"`
// 	Request          string    `json:"request"`
// 	Response         string    `json:"response"`
// 	IP               string    `json:"ip"`
// 	Timestamp        time.Time `json:"timestamp"`
// 	CurlCommand      string    `json:"curl-command,omitempty"`
// 	MatcherStatus    bool      `json:"matcher-status"`
// 	ExtractedResults []string  `json:"extracted-results,omitempty"`
// 	Meta             struct {
// 		SapPath string `json:"sap_path"`
// 	} `json:"meta,omitempty"`
// }

// func main() {
// 	b, _ := os.ReadFile("./result.json")
// 	var nr NucleiResult
// 	json.Unmarshal(b, &nr)
// 	fmt.Printf("nr: %v\n", len(nr))
// 	for _, result := range nr {
// 		fmt.Printf("漏洞名称: %v\n", result.TemplateID)
// 		fmt.Printf("漏洞类型: %v\n", result.Type)
// 		fmt.Printf("漏洞等级: %v\n", result.Info.Severity)
// 		fmt.Printf("漏洞地址: %v\n", result.MatchedAt)
// 		fmt.Printf("拓展信息: %v\n", strings.Join(result.ExtractedResults, " | "))
// 	}
// }
