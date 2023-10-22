package runner

import (
	"fmt"
	"slack/common"
	"slack/lib/poc"
	"slack/lib/protocols/retryhttpclient"
	"slack/lib/report"

	"slack/lib/result"

	"strings"
)

type Runner struct {
	options       *common.Options
	Report        *report.Report
	OnResult      func(*result.Result)
	PocsEmbedYaml []string
	engine        *Engine
}

func NewRunner(options *common.Options) (*Runner, error) {
	runner := &Runner{options: options}
	runner.engine = NewEngine(options)
	retryhttpclient.Init(&retryhttpclient.Options{
		Proxy:           options.Proxy,
		Timeout:         options.Timeout,
		Retries:         options.Retries,
		MaxRespBodySize: options.MaxRespBodySize,
	})

	report, err := report.NewReport(report.DefaultTemplate)
	if err != nil {
		return runner, err
	}
	runner.Report = report

	if len(runner.options.Target) > 0 {
		for _, t := range runner.options.Target {
			runner.options.Targets.Append(t)
		}
	}
	// init pocs
	if len(runner.options.PocFile) > 0 {
		runner.options.PocsDirectory = append(runner.options.PocsDirectory, runner.options.PocFile)
	}
	runner.PocsEmbedYaml = poc.EmbedFileList

	checkReversePlatform()

	return runner, nil
}

// 检测反连
func checkReversePlatform() {
	if !CeyeTest() {
		common.ReverseCeyeLive = false
	} else {
		common.ReverseCeyeLive = true
	}
}

func CeyeTest() bool {
	url := fmt.Sprintf("http://%s.%s", "test", common.ReverseCeyeDomain)
	resp, _, err := retryhttpclient.Get(url)
	if err != nil {
		return false
	}
	if strings.Contains(string(resp), "\"meta\":") || strings.Contains(string(resp), "201") {
		return true
	}
	return false
}
