package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"

	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func NewNucleiEngine(ctx context.Context, proxy clients.Proxy, o structs.NucleiOption) {
	if o.SkipNucleiWithoutTags && len(o.Tags) == 0 {
		gologger.Info(ctx, fmt.Sprintf("[nuclei] %s does not have tags, scan skipped", o.URL))
		return
	}
	options := []nuclei.NucleiSDKOptions{
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
		nuclei.DisableUpdateCheck(), // -duc
		// nuclei.WithNetworkConfig(nuclei.NetworkConfig{SourceIP: ""}),
	}
	// 自定义请求头
	if o.CustomHeaders != "" {
		options = append(options, nuclei.WithHeaders(clients.Str2HeaderList(o.CustomHeaders)))
	}
	// 判断是使用指定poc文件还是根据标签
	if len(o.TemplateFile) == 0 {
		options = append(options, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: o.TemplateFolders,
		}))
		// 如果自定义标签不为空则使用
		options = append(options, nuclei.WithTemplateFilters(nuclei.TemplateFilters{
			Tags: finalTags(o.Tags, o.CustomTags),
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
		var reference string
		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
			reference = strings.Join(event.Info.Reference.ToSlice(), ",")
		}
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			ID:           event.TemplateID,
			Name:         event.Info.Name,
			Description:  event.Info.Description,
			Reference:    reference,
			URL:          showMatched(event),
			Request:      showRequest(event),
			Response:     showResponse(event),
			ResponseTime: limitDecimalPlaces(event.ResponseTime),
			Extract:      strings.Join(event.ExtractedResults, " | "),
			Type:         strings.ToUpper(event.Type),
			Severity:     strings.ToUpper(event.Info.SeverityHolder.Severity.String()),
		})
	})
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("nuclei execute callback err: %v", err))
		return
	}
	defer ne.Close()
}

func finalTags(detectTags, customTags []string) []string {
	if len(customTags) != 0 {
		return customTags
	}
	return detectTags
}

// func NewThreadSafeNucleiEngine(ctx context.Context, proxy clients.Proxy, o structs.NucleiOption) {
// 	if o.SkipNucleiWithoutTags && len(o.Tags) == 0 {
// 		gologger.Info(ctx, fmt.Sprintf("[nuclei] %s does not have tags, scan skipped", o.URL))
// 		return
// 	}
// 	options := []nuclei.NucleiSDKOptions{
// 		nuclei.DisableUpdateCheck(), // -duc
// 	}
// 	// 判断是使用指定poc文件还是根据标签
// 	if len(o.TemplateFile) == 0 {
// 		options = append(options, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
// 			Templates: o.TemplateFolders,
// 		}))
// 		options = append(options, nuclei.WithTemplateFilters(nuclei.TemplateFilters{
// 			Tags: o.Tags,
// 		}))
// 	} else {
// 		// 指定poc文件的时候就要删除tags标签
// 		options = append(options, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
// 			Templates: o.TemplateFile,
// 		}))
// 	}
// 	var proxys string
// 	if proxy.Enabled {
// 		proxys = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Mode, proxy.Username, proxy.Password, proxy.Address, proxy.Port)
// 		options = append(options, nuclei.WithProxy([]string{proxys}, false)) // -proxy
// 	}
// 	ne, err := nuclei.NewThreadSafeNucleiEngineCtx(context.Background(), options...)
// 	if err != nil {
// 		gologger.Error(ctx, fmt.Sprintf("nuclei init engine err: %v", err))
// 		return
// 	}
// 	// load targets and optionally probe non http/https targets
// 	gologger.Info(ctx, fmt.Sprintf("[nuclei] check vuln: %s", o.URL))
// 	ne.ExecuteNucleiWithOptsCtx(context.Background(), []string{o.URL}, options...)
// 	ne.GlobalResultCallback(func(event *output.ResultEvent) {
// 		gologger.Success(ctx, fmt.Sprintf("[%s] [%s] %s", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched))
// 		var reference string
// 		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
// 			reference = strings.Join(event.Info.Reference.ToSlice(), ",")
// 		}
// 		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
// 			ID:           event.TemplateID,
// 			Name:         event.Info.Name,
// 			Description:  event.Info.Description,
// 			Reference:    reference,
// 			URL:          showMatched(event),
// 			Request:      showRequest(event),
// 			Response:     showResponse(event),
// 			ResponseTime: limitDecimalPlaces(event.ResponseTime),
// 			Extract:      strings.Join(event.ExtractedResults, " | "),
// 			Type:         strings.ToUpper(event.Type),
// 			Severity:     strings.ToUpper(event.Info.SeverityHolder.Severity.String()),
// 		})
// 	})

// 	defer ne.Close()
// }

func showMatched(event *output.ResultEvent) string {
	if event.Matched != "" {
		return event.Matched
	}
	return event.URL
}

func showRequest(event *output.ResultEvent) string {
	if event.Request != "" {
		return event.Request
	}
	if event.Interaction != nil {
		return event.Interaction.RawRequest
	}
	return ""
}

func showResponse(event *output.ResultEvent) string {
	if event.Response != "" {
		byteResponse := []byte(event.Response)
		if len(byteResponse) > 1024*512 {
			return string(byteResponse[:1024*512]) + " ..."
		}
		return event.Response
	}
	if event.Interaction != nil {
		return event.Interaction.RawResponse
	}
	return ""
}

// 限制小数为2位，用于截取时间字符串
func limitDecimalPlaces(value string) string {
	parts := strings.Split(value, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		value = parts[0] + "." + parts[1][:2] // 截取小数点后两位
	}
	return value
}
