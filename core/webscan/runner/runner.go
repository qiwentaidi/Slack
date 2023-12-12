package runner

import (
	"bytes"
	"fmt"
	"slack-wails/core/webscan/proto"
	"slack-wails/lib/clients"
	"slack-wails/lib/report"
	"slack-wails/lib/util"

	"strings"
	"time"
)

type Runner struct {
	options       *Options
	Report        *report.Report
	OnResult      func(*report.Result)
	PocsEmbedYaml []string
	engine        *Engine
}

func NewRunner(options *Options, pocList []string) (*Runner, error) {
	runner := &Runner{options: options}
	runner.engine = &Engine{
		options: options,
	}
	clients.Init(&clients.Options{
		Proxy:   options.Proxy,
		Timeout: options.Timeout,
	})
	// report, err := report.NewReport(report.DefaultTemplate)
	// if err != nil {
	// 	return runner, err
	// }
	// runner.Report = report
	runner.options.Targets.Append(runner.options.Target)
	runner.PocsEmbedYaml = pocList

	checkReversePlatform()

	return runner, nil
}

var (
	ReverseCeyeApiKey = "ba446c3277a60555ad9e74a6f0cb4290"
	ReverseCeyeDomain = "xrn0nb.ceye.io"
	ReverseCeyeLive   bool
)

func reverseCheck(r *proto.Reverse, timeout int64) bool {
	if r == nil || (len(r.Domain) == 0 && len(r.Ip) == 0) {
		return false
	}
	time.Sleep(time.Second * time.Duration(timeout))
	urlStr := ""
	sub := strings.Split(r.Domain, ".")[0]
	if ReverseCeyeLive {
		urlStr = fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s", ReverseCeyeApiKey, sub)
		_, body, err := clients.NewRequest("GET", urlStr, nil, nil, 10, clients.DefaultClient())
		if err != nil {
			return false
		}
		if !bytes.Contains(body, []byte(`"data": []`)) && bytes.Contains(body, []byte(`{"code": 200`)) {
			return true
		}
		if bytes.Contains(body, []byte(`<title>503`)) {
			return false
		}
	}

	return false
}

// 检测反连是否开启
func checkReversePlatform() {
	if !CeyeTest() {
		ReverseCeyeLive = false
	} else {
		ReverseCeyeLive = true
	}
}

func CeyeTest() bool {
	url := fmt.Sprintf("http://%s.%s", util.RandLetters(5), ReverseCeyeDomain)
	_, body, err := clients.NewRequest("GET", url, nil, nil, 10, clients.DefaultClient())
	if err != nil {
		return false
	}
	if strings.Contains(string(body), "\"meta\":") || strings.Contains(string(body), "201") {
		return true
	}
	return false
}
