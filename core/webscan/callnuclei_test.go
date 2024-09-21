package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/util"
	"strings"
	"testing"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

func TestNucleiCaller(t *testing.T) {
	// home := util.HomeDir()
	// var templateDir = home + "/slack/config/pocs"
	// var webfingerFile = home + "/slack/config/webfinger.yaml"
	// var activefingerFile = home + "/slack/config/dir.yaml"
	// NewConfig().InitAll(context.TODO(), webfingerFile, activefingerFile, templateDir)
	// nc := NewNucleiCaller("", false, "", clients.Proxy{})
	// v := nc.CallerFP(context.TODO(), FingerPoc{
	// 	URL:  "http://erpdev.wingtech.com:50000/",
	// 	Tags: []string{"SAP-Web-App-Server"},
	// })
	// fmt.Printf("v: %v\n", v)
	proxys := []string{"http://127.0.0.1:8080"}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(),
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{util.HomeDir() + "/slack/config/pocs"},
		}), // -t
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{Tags: []string{"Jeecg-Boot"}}), // 过滤 poc
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}),          // optionally enable metrics server for better observability
		nuclei.DisableUpdateCheck(),     // -duc
		nuclei.WithProxy(proxys, false), // -proxy
	)
	if err != nil {
		panic(err)
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{"http://xxxx"}, false)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		fmt.Printf("[%s] [%s] %s\n", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched)
		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
			fmt.Printf("Reference: %s\n", event.Info.Reference.ToSlice())
		}
		fmt.Printf("ExtractedResults: %s\n", strings.Join(event.ExtractedResults, ","))
	})
	if err != nil {
		panic(err)
	}
	defer ne.Close()
}
