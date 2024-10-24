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
	SkipNucleiWithoutTags bool // 如果没有扫描到指纹，是否需要扫描全漏洞还是直接跳过
	URL                   string
	Tags                  []string // 全漏洞扫描时，使用自定义标签
	TemplateFile          []string
}

var pocFile = util.HomeDir() + "/slack/config/pocs"

func NewNucleiEngine(ctx context.Context, proxy clients.Proxy, o NucleiOption) {
	if o.SkipNucleiWithoutTags && len(o.Tags) == 0 {
		gologger.Info(ctx, fmt.Sprintf("[nuclei] %s does not have tags, scan skipped", o.URL))
		return
	}
	options := []nuclei.NucleiSDKOptions{
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
		nuclei.DisableUpdateCheck(), // -duc
	}
	// 判断是使用指定poc文件还是根据标签
	if len(o.TemplateFile) == 0 {
		options = append(options, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{pocFile},
		}))
		options = append(options, nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags: o.Tags,
		}))
	} else {
		// 指定poc文件的时候就要删除tags标签
		options = append(options, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: o.TemplateFile,
		}))
	}
	var proxys string
	if proxy.Enabled {
		proxys = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Mode, proxy.Username, proxy.Password, proxy.Address, proxy.Port)
		options = append(options, nuclei.WithProxy([]string{proxys}, false)) // -proxy
	}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(), options...)
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("nuclei init engine err: %v", err))
		return
	}
	// load targets and optionally probe non http/https targets
	gologger.Info(ctx, fmt.Sprintf("[nuclei] check vuln: %s", o.URL))
	ne.LoadTargets([]string{o.URL}, false)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		gologger.Success(ctx, fmt.Sprintf("[%s] [%s] %s", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched))
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
			Type:        strings.ToUpper(event.Type),
			Risk:        strings.ToUpper(event.Info.SeverityHolder.Severity.String()),
		})
	})
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("nuclei execute callback err: %v", err))
		return
	}
	defer ne.Close()
}
