package webscan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slack-wails/lib/bridge"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type NucleiResult []struct {
	TemplateID   string `json:"template-id"`
	TemplatePath string `json:"template-path"`
	Info         struct {
		Name        string   `json:"name"`
		Author      []string `json:"author"`
		Tags        []string `json:"tags"`
		Description string   `json:"description"`
		Reference   []string `json:"reference"`
		Severity    string   `json:"severity"`
		Metadata    struct {
			MaxRequest int `json:"max-request"`
		} `json:"metadata"`
	} `json:"info"`
	Type             string    `json:"type"`
	Host             string    `json:"host"`
	Port             string    `json:"port"`
	Scheme           string    `json:"scheme"`
	URL              string    `json:"url"`
	MatchedAt        string    `json:"matched-at"`
	Request          string    `json:"request"`
	Response         string    `json:"response"`
	IP               string    `json:"ip"`
	Timestamp        time.Time `json:"timestamp"`
	CurlCommand      string    `json:"curl-command,omitempty"`
	MatcherStatus    bool      `json:"matcher-status"`
	ExtractedResults []string  `json:"extracted-results,omitempty"`
	Meta             struct {
		SapPath string `json:"sap_path"`
	} `json:"meta,omitempty"`
}

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
	Mode int
	// Target      string
	// Fingerprint []string
	Engine     string
	Interactsh bool
	CustomTags []string // 全漏洞扫描时，使用自定义标签
	Severity   string
	Proxy      clients.Proxy
}

type NucleiCaller struct {
	CommandLine []string
	NucleiPath  string
	Interactsh  string
	Severity    string
	Proxy       string
}

var (
	pocFile    = util.HomeDir() + "/slack/config/pocs"
	reportTemp = util.HomeDir() + "/slack/web_report/"
)

func NewNucleiCaller(path string, interactsh bool, severity string, proxy clients.Proxy) *NucleiCaller {
	var nucleiPath, ni, nucleiProxy string
	// 存在环境变量
	if path == "" {
		nucleiPath = "nuclei"
	} else {
		nucleiPath = path
	}
	// 判断反连
	if interactsh {
		ni = "-ni"
	}
	if proxy.Enabled {
		nucleiProxy = fmt.Sprintf("%s://%s:%d", proxy.Mode, proxy.Address, proxy.Port)
	}
	return &NucleiCaller{
		NucleiPath: nucleiPath,
		Interactsh: ni,
		Severity:   severity,
		Proxy:      nucleiProxy,
	}
}

// 检查存储报告的文件夹是否存在
func (nc *NucleiCaller) ReportDirStat() error {
	if _, err := os.Stat(reportTemp); err != nil {
		return os.MkdirAll(reportTemp, 0755)
	}
	return nil
}

func (nc *NucleiCaller) Enabled(ctx context.Context) bool {
	cmd := exec.Command(nc.NucleiPath, "--version")
	bridge.HideExecWindow(cmd)
	out, err := cmd.CombinedOutput()
	gologger.Info(ctx, string(out))
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Nuclei Engine Version")
}

func (nc *NucleiCaller) ReadNucleiJson(ctx context.Context, reportJson string) error {
	b, err := os.ReadFile(reportJson)
	if err != nil {
		return err
	}
	var nr NucleiResult
	err = json.Unmarshal(b, &nr)
	if err != nil {
		return err
	}
	for _, result := range nr {
		runtime.EventsEmit(ctx, "nucleiResult", VulnerabilityInfo{
			ID:          result.TemplateID,
			Name:        result.Info.Name,
			Description: result.Info.Description,
			Reference:   strings.Join(result.Info.Reference, ","),
			URL:         result.MatchedAt,
			Request:     result.Request,
			Response:    result.Response,
			Extract:     strings.Join(result.ExtractedResults, " | "),
			Type:        result.Type,
			Risk:        result.Info.Severity,
		})
	}
	return nil
}

// Finger POC
func (nc *NucleiCaller) CallerFP(ctx context.Context, pe FingerPoc) error {
	if len(pe.Tags) == 0 {
		gologger.Info(ctx, fmt.Sprintf("No tags found in %s", pe.URL))
		return nil
	}
	reportJson := fmt.Sprintf("%s%d.json", reportTemp, time.Now().UnixMilli())
	nc.CommandLine = []string{"-duc", "-u", pe.URL, "-t", pocFile, "-tags", strings.Join(pe.Tags, ","), "-je", reportJson, nc.Interactsh}
	if nc.Proxy != "" {
		nc.CommandLine = append(nc.CommandLine, "-proxy", nc.Proxy)
	}
	cmdLine := fmt.Sprintf("CallerFP: %s %s", nc.NucleiPath, strings.Join(nc.CommandLine, " "))
	gologger.Info(ctx, cmdLine)
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	bridge.HideExecWindow(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	runtime.EventsEmit(ctx, "nuclei-pid", cmd.Process.Pid)
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nc.ReadNucleiJson(ctx, reportJson)
}

// ALL POC
func (nc *NucleiCaller) CallerAP(ctx context.Context, target string, tags []string) error {
	reportJson := fmt.Sprintf("%s%d.json", reportTemp, time.Now().UnixMilli())
	nc.CommandLine = []string{"-duc", "-t", pocFile, "-je", reportJson, nc.Interactsh}
	// 风险等级、关键词筛选
	if nc.Severity != "" {
		nc.CommandLine = append(nc.CommandLine, []string{"-s", nc.Severity}...)
	}
	if len(tags) != 0 {
		nc.CommandLine = append(nc.CommandLine, []string{"-tags", strings.Join(tags, ",")}...)
	}
	if nc.Proxy != "" {
		nc.CommandLine = append(nc.CommandLine, "-proxy", nc.Proxy)
	}
	nc.CommandLine = append(nc.CommandLine, "-u", target)
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	bridge.HideExecWindow(cmd) // 让windows执行cmd时无窗口
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("callerap err: %v", err)
	}
	// runtime.EventsEmit(ctx, "nuclei-pid", cmd.Process.Pid)
	// cmd.Wait()
	return nc.ReadNucleiJson(ctx, reportJson)
}

func (nc *NucleiCaller) FilterPoc(pocs, keywords []string) []string {
	news := []string{}
	for _, poc := range pocs {
		for _, key := range keywords {
			if strings.Contains(poc, key) {
				news = append(news, poc)
			}
		}
	}
	return news
}
