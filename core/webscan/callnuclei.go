package webscan

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slack-wails/lib/util"
	"strings"
	"time"
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
	ReportName  string
	Interactsh  string
	Severity    string
}

const reportTemp = "/slack/web_report/"

func NewNucleiCaller(path, reportName string, interactsh bool, severity string) *NucleiCaller {
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
		ReportName: reportName,
		Interactsh: ni,
		Severity:   severity,
	}
}

// 检查存储报告的文件夹是否存在
func (nc *NucleiCaller) ReportDirStat() error {
	path := util.HomeDir() + reportTemp
	if _, err := os.Stat(path); err != nil {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func (nc *NucleiCaller) Enabled() bool {
	cmd := exec.Command(nc.NucleiPath, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "Nuclei Engine Version")
}

func (nc *NucleiCaller) TargetBindFingerPocs(target string, fingerprints []string) FingerPoc {
	var fp FingerPoc
	fp.URL = target
	for fn1, we := range WorkFlowDB {
		for _, fn2 := range fingerprints {
			if fn1 == fn2 {
				fp.PocFiles = append(fp.PocFiles, FullPocName(we.PocsName)...)
			}
		}
	}
	return fp
}

func (nc *NucleiCaller) ReadNucleiJson(reportName string) []VulnerabilityInfo {
	var vis []VulnerabilityInfo
	b, err := os.ReadFile(util.HomeDir() + reportTemp + "temp.json")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	var nr NucleiResult
	json.Unmarshal(b, &nr)
	for _, result := range nr {
		vis = append(vis, VulnerabilityInfo{
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
	return vis
}

// Finger POC
func (nc *NucleiCaller) CallerFP(pe FingerPoc) []VulnerabilityInfo {
	result := util.HomeDir() + reportTemp + "temp.json"
	nc.CommandLine = []string{"-duc", "-u", pe.URL, "-t", strings.Join(pe.PocFiles, ","), "-je", result, nc.Interactsh}
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	if err := cmd.Run(); err != nil {
		fmt.Printf("err: %v\n", err)
		return []VulnerabilityInfo{}
	}
	return nc.ReadNucleiJson(nc.ReportName)
}

var (
	pocFile = util.HomeDir() + "/slack/config/pocs"
	result  = util.HomeDir() + reportTemp + "temp.json"
)

// ALL POC
func (nc *NucleiCaller) CallerAP(target string, keywords []string) []VulnerabilityInfo {
	var pocs []string
	nc.CommandLine = []string{"-duc", "-u", target, "-je", result, nc.Interactsh}
	// 风险等级、关键词筛选
	if nc.Severity == "" && len(keywords) == 0 {
		nc.CommandLine = append(nc.CommandLine, []string{"-t", pocFile}...)
	} else {
		if nc.Severity != "" {
			nc.CommandLine = append(nc.CommandLine, []string{"-s", nc.Severity}...)
		}
		if len(keywords) != 0 {
			pocs = nc.FilterPoc(ALLPoc(), keywords)
			nc.CommandLine = append(nc.CommandLine, []string{"-t", strings.Join(pocs, ",")}...)
		}
	}
	cmd := exec.Command(nc.NucleiPath, nc.CommandLine...)
	if err := cmd.Run(); err != nil {
		fmt.Printf("callerap err: %v\n", err)
		return []VulnerabilityInfo{}
	}
	return nc.ReadNucleiJson(nc.ReportName)
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
