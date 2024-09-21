package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"

	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type VulnerabilityInfo struct {
	ID          string
	Name        string
	Description string
	Reference   string
	Type        string
	Risk        string
	URL         string
	Request     string
	Response    string
	Extract     string
}

type NucleiOption struct {
	URL      string
	Tags     []string // 全漏洞扫描时，使用自定义标签
	Severity string
}

var pocFile = util.HomeDir() + "/slack/config/pocs"

func NewNucleiEngine(ctx context.Context, proxy clients.Proxy, o NucleiOption) {
	options := []nuclei.NucleiSDKOptions{
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{pocFile},
		}), // -t
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags:     o.Tags,
			Severity: o.Severity,
		}), // 过滤 poc
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
		nuclei.DisableUpdateCheck(), // -duc
	}
	var proxys string
	if proxy.Enabled {
		proxys = fmt.Sprintf("%s://%s:%d", proxy.Mode, proxy.Address, proxy.Port)
		options = append(options, nuclei.WithProxy([]string{proxys}, false)) // -proxy
	}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(), options...)
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("nuclei init engine err: %v", err))
		return
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{o.URL}, false)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		// fmt.Printf("[%s] [%s] %s\n", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched)
		// fmt.Printf("Request: \n%s\n", event.Request)
		// fmt.Printf("Response: \n%s\n", event.Response)
		var reference string
		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
			reference = strings.Join(event.Info.Reference.ToSlice(), ",")
		}
		runtime.EventsEmit(ctx, "nucleiResult", VulnerabilityInfo{
			ID:          event.TemplateID,
			Name:        event.Info.Name,
			Description: event.Info.Description,
			Reference:   reference,
			URL:         event.Matched,
			Request:     event.Request,
			Response:    event.Response,
			Extract:     strings.Join(event.ExtractedResults, " | "),
			Type:        event.Type,
			Risk:        event.Info.SeverityHolder.Severity.String(),
		})
	})
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("nuclei execute callback err: %v", err))
		return
	}
	defer ne.Close()
}

// func Rename(filename string) string {
// 	filename = strings.ReplaceAll(filename, ":", "_")
// 	filename = strings.ReplaceAll(filename, "/", "_")
// 	filename = strings.ReplaceAll(filename, "___", "_")
// 	return filename
// }
