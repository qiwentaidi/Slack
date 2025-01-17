package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
	"testing"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

func TestNucleiCaller(t *testing.T) {
	// proxys := []string{"http://127.0.0.1:8080"}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(),
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{util.HomeDir() + "/slack/config/pocs"},
		}), // -t
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{Tags: []string{}}),    // 过滤 poc
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
		nuclei.DisableUpdateCheck(), // -duc
		// nuclei.WithProxy(proxys, false), // -proxy
	)
	if err != nil {
		panic(err)
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{}, false)
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

func TestThreadNucleiCaller(t *testing.T) {
	option1 := structs.NucleiOption{
		URL:          "http://xxxxx",
		TemplateFile: []string{util.HomeDir() + "/slack/config/pocs/iis-shortname.yaml"},
	}
	option2 := structs.NucleiOption{
		URL:          "https://xxxx/",
		TemplateFile: []string{util.HomeDir() + "/slack/config/pocs/iis-shortname.yaml"},
	}
	NewThreadSafeNucleiEngine(context.Background(), []structs.NucleiOption{option1, option2})
}
