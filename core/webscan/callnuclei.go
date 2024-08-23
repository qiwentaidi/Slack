package webscan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slack-wails/lib/bridge"
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

type NucleiCaller struct {
	CommandLine []string
	NucleiPath  string
	Interactsh  string
	Severity    string
}

var (
	pocFile    = util.HomeDir() + "/slack/config/pocs"
	reportTemp = util.HomeDir() + "/slack/web_report/"
	result     = reportTemp + "temp.json"
)

func NewNucleiCaller(path string, interactsh bool, severity string) *NucleiCaller {
	var nucleiPath, ni string
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
	return &NucleiCaller{
		NucleiPath: nucleiPath,
		Interactsh: ni,
		Severity:   severity,
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
	gologger.Debug(ctx, string(out))
	if err != nil {
		gologger.Debug(ctx, err)
		return false
	}
	return strings.Contains(string(out), "Nuclei Engine Version")
}

func (nc *NucleiCaller) ReadNucleiJson(ctx context.Context) error {
	b, err := os.ReadFile(result)
	if err != nil {
		return err
	}
	var nr NucleiResult
	if err = json.Unmarshal(b, &nr); err != nil {
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
	nc.CommandLine = []string{"-duc", "-u", pe.URL, "-t", pocFile, "-tags", strings.Join(util.RemoveDuplicates(pe.Tags), ","), "-je", result, nc.Interactsh}
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	bridge.HideExecWindow(cmd)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nc.ReadNucleiJson(ctx)
}

// ALL POC
func (nc *NucleiCaller) CallerAP(ctx context.Context, target string, tags []string) error {
	nc.CommandLine = []string{"-duc", "-u", target, "-t", pocFile, "-je", result, nc.Interactsh}
	// 风险等级、关键词筛选
	if nc.Severity != "" {
		nc.CommandLine = append(nc.CommandLine, []string{"-s", nc.Severity}...)
	}
	if len(tags) != 0 {
		nc.CommandLine = append(nc.CommandLine, []string{"-tags", strings.Join(tags, ",")}...)
	}
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	bridge.HideExecWindow(cmd) // 让windows执行cmd时无窗口
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("callerap err: %v", err)
	}
	return nc.ReadNucleiJson(ctx)
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
