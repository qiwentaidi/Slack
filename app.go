package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slack-wails/core"
	"slack-wails/core/info"
	"slack-wails/core/portscan"
	"slack-wails/core/space"
	"slack-wails/core/webscan"
	"slack-wails/core/webscan/poc"
	"slack-wails/core/webscan/runner"
	"slack-wails/lib/clients"
	"slack-wails/lib/gonmap"
	"slack-wails/lib/report"
	"slack-wails/lib/update"
	"slack-wails/lib/util"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func (a *App) Minimise() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) ToggleMaximise() {
	runtime.WindowToggleMaximise(a.ctx)
}

// 使用默认浏览器打开链接，避免使用webview2
func (a *App) DefaultOpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) CheckFileStat(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (a *App) GetFileContent(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "文件不存在"
	}
	return string(b)
}

type Response struct {
	OK     bool
	Status int
	Text   string
}

func (a *App) GoSimpleFetch(url string) *Response {
	var r Response
	resp, b, err := clients.NewRequest("GET", url, nil, nil, 10, clients.DefaultClient())
	if err == nil {
		r.OK = true
	}
	r.Status = resp.StatusCode
	r.Text = string(b)
	return &r
}

// wx
func (a *App) WechatAppid(appid, secert string) string {
	return core.CheckSecert(appid, secert)
}

// fscan
func (a *App) Fscan2Txt(content string) string {
	var result string
	lines := strings.Split(content, "\n")
	for name, reg := range core.FscanRegs {
		result += core.MatchLine(name, reg, lines)
	}
	return result
}

// thinkdict
func (a *App) ThinkDict(p_name, c_name, c_domain, p_birthday, p_worknum string) []string {
	return core.GenerateDict(p_name, c_name, c_domain, p_birthday, p_worknum)
}

func (a *App) System(content string, mode int) [][]string {
	if mode == 0 { // tasklist
		datas, _ := core.AntivirusIdentify(content)
		return datas
	} else { // systeminfo
		datas := core.Patch(content)
		return datas
	}
}

func (a *App) Transcoding(crypt []string, input string, mode int) string {
	var result string
	if mode == 0 { // 加密
		core.RecursiveEncrypt(crypt, input, &result)
	} else {
		core.RecursiveDecrypt(crypt, input, &result)
	}
	return result
}

func (a *App) ExtractIP(text string) string {
	return core.Analysis(text)
}

func (a *App) DomainInfo(text string) string {
	s, s2, s3 := core.SeoChinaz(text)
	s4, s5, s6, s7, s8, s9 := core.WhoisChinaz(text)
	return fmt.Sprintf(`---备案查询---
公司名称: %v
	
备案号: %v
	
IP: %v
	
---whois查询---
更新时间: %v

创建时间: %v

过期时间: %v

注册商服务器: %v

DNS: %v

状态: %v`, s, s2, s3, s4, s5, s6, s7, s8, s9)
}

func (a *App) CheckCdn(input, dns1, dns2 string) string {
	var result string
	result += "---域名解析(CDN查询)---:\n"
outerLoop:
	for _, domain := range util.RemoveDuplicates[string](util.RegDomain.FindAllString(input, -1)) {
		ips, cnames, err := core.Resolution(domain, []string{dns1 + ":53", dns2 + ":53"}, 5)
		if err == nil {
			for name, cdns := range core.ReadCDNFile() {
				for _, cdn := range cdns {
					for _, cname := range cnames {
						if strings.Contains(cname, cdn) { // 识别到cdn
							result += fmt.Sprintf("域名: %v 识别到CDN域名，CNAME: %v CDN名称: %v 解析到IP为: %v\n", domain, cname, name, strings.Join(ips, " | "))
							break outerLoop
						} else if strings.Contains(cname, "cdn") {
							result += fmt.Sprintf("域名: %v CNAME中含有关键字: cdn，该域名可能使用了CDN技术 CNAME: %v 解析到IP为: %v \n", domain, cname, strings.Join(ips, " | "))
							break outerLoop
						}
					}
				}
			}
			result += fmt.Sprintf("域名: %v 解析到IP为: %v\n", domain, strings.Join(ips, ","))
		} else {
			result = fmt.Sprintf("域名: %v 解析失败,%v\n", domain, err)
		}
	}
	return result
}

// 初始化IP记录器
func (a *App) InitIPResolved() {
	core.IPResolved = make(map[string]int)
}

var onec sync.Once

func (a *App) Subdomain(subdomain, dns1, dns2 string, timeout int) []string {
	var data map[string][]string
	onec.Do(func() {
		data = core.ReadCDNFile()
	})
	sr := core.BurstSubdomain(subdomain, []string{dns1 + ":53", dns2 + "53"}, data, timeout)
	return []string{sr.Subdomain, strings.Join(sr.Cname, " | "), strings.Join(sr.Ips, " | "), sr.Notes}
}

func (a *App) InitTycHeader(token string) {
	info.InitHEAD(token)
}

type SubcompanyInfo struct {
	Shareholding [][]string
	Prompt       string
}

func (*App) AssetSubcompany(companyName string, ratio int) SubcompanyInfo {
	asset, logs := info.SearchSubsidiary(companyName, ratio)
	return SubcompanyInfo{
		Shareholding: asset,
		Prompt:       logs,
	}
}

type WechatInfo struct {
	OfficialAccounts [][]string
	Prompt           string
}

func (*App) AssetWechat(companyName string) WechatInfo {
	asset, info := info.WeChatOfficialAccounts(companyName)
	return WechatInfo{
		OfficialAccounts: asset,
		Prompt:           info,
	}
}

type HunterSearch struct {
	Total string
	Info  string
}

// mode = 0 campanyName, mode = 1 domain or ip
func (a *App) AssetHunter(mode int, target, api string) HunterSearch {
	if mode == 0 {
		str := fmt.Sprintf("icp.name=\"%v\"", target)
		t, i := space.SearchTotal(api, str)
		return HunterSearch{
			Total: fmt.Sprint(t),
			Info:  i,
		}
	} else {
		var str string
		// 处理网站域名是IP的情况
		if util.RegIP.MatchString(target) {
			str = fmt.Sprintf("ip=\"%v\"", target)
		} else {
			str = fmt.Sprintf("domain.suffix=\"%v\"", target)
		}
		t, i := space.SearchTotal(api, str)
		return HunterSearch{
			Total: fmt.Sprint(t),
			Info:  i,
		}
	}
}

// dirsearch

// 测试URL是否可达
func (a *App) TestTarget(target string) bool {
	_, err := url.Parse(target)
	if err != nil {
		return false
	}
	if _, _, err2 := clients.NewRequest("GET", target, nil, nil, 10, clients.DefaultClient()); err2 != nil {
		return false
	}
	return true
}

func (a *App) InitDict(newExts []string) []string {
	return core.LoadDefaultDict("dicc.txt", "%EXT%", newExts)
}

type PathData struct {
	Status   int    // 状态码
	Location string // server信息
	Length   int    // 主体内容
}

func (a *App) PathRequest(method, url string, timeout int, bodyExclude string) PathData {
	var pd PathData // 将响应头和响应头的数据存储到结构体中
	resp, body, err := clients.NewRequest(method, url, nil, nil, timeout, clients.NotFollowClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return pd
	}
	if bodyExclude != "" && bytes.Contains(body, []byte(bodyExclude)) {
		pd.Status = 1 // status 1 表示被body排除在外，不计入ERROR请求中
	} else {
		pd.Status = resp.StatusCode
	}
	pd.Length = len(body)
	pd.Location = resp.Header.Get("Location")
	return pd
}

// portscan

func (a *App) IPParse(text string) []string {
	return util.ParseIPs(text)
}

func (a *App) PortParse(text string) []int {
	return util.ParsePort(text)
}

func (a *App) HostAlive(targets []string, Ping bool) []string {
	return portscan.CheckLive(targets, Ping)
}

type PortResult struct {
	Status    bool
	Server    string
	Link      string
	HttpTitle string
}

func (a *App) PortCheck(ip string, port, timeout int) PortResult {
	var pr PortResult
	scanner := gonmap.New()
	_, response := scanner.ScanTimeout(ip, port, time.Second*time.Duration(timeout))
	if response != nil {
		pr.Status = true
		if response.FingerPrint.Service != "" {
			pr.Server = response.FingerPrint.Service
		} else {
			pr.Server = "unknow"
		}
		pr.Link = fmt.Sprintf("%v://%v:%v", pr.Server, ip, port)
		if pr.Server == "http" || pr.Server == "https" {
			if _, b, err := clients.NewRequest("GET", pr.Link, nil, nil, 10, clients.DefaultClient()); err == nil {
				if match := util.RegTitle.FindSubmatch(b); len(match) > 1 {
					pr.HttpTitle = util.Str2UTF8(string(match[1]))
				} else {
					pr.HttpTitle = "-"
				}
			}
		}
	}
	return pr
}

// 漏洞详情

// 遍历文件夹下的yaml或者yml文件，获取所有绝对路径
func (a *App) LocalWalkFiles(folderPath string) []string {
	fileList := []string{}
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 检查是否为文件且以 .yaml 或 .yml 扩展名结尾
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			fileList = append(fileList, path)
		}

		return nil
	})
	return fileList
}

func (a *App) ReadPocDetail(absolutePath string) poc.VulnerabilityDetails {
	p, err := poc.LocalReadPocByPath(absolutePath)
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	return poc.VulnerabilityDetails{
		Name:        p.Info.Name,
		Risk:        p.Info.Severity,
		Author:      p.Info.Author,
		Tags:        p.Info.Tags,
		Description: p.Info.Description,
		Reference:   strings.Join(p.Info.Reference, "\n"),
		Affected:    p.Info.Affected,
		Solutions:   p.Info.Solutions,
	}
}

// 端口暴破
func (a *App) PortBrute(host string, usernames, passwords []string) *portscan.Burte {
	u, err := url.Parse(host)
	if err != nil {
		return nil
	}
	protocol := strings.Split(host, "://")[0]
	switch protocol {
	case "ftp":
		return portscan.FtpScan(u.Host, usernames, passwords)
	case "ssh":
		return portscan.SshScan(u.Host, usernames, passwords)
	case "telnet":
		return portscan.TelenetScan(u.Host, usernames, passwords)
	case "smb":
		return portscan.SmbScan(u.Host, usernames, passwords)
	case "oracle":
		return portscan.OracleScan(u.Host, usernames, passwords)
	case "mssql":
		return portscan.MssqlScan(u.Host, usernames, passwords)
	case "mysql":
		return portscan.MysqlScan(u.Host, usernames, passwords)
	case "rdp":
		return portscan.RdpScan(u.Host, usernames, passwords)
	case "postgresql":
		return portscan.PostgresScan(u.Host, usernames, passwords)
	case "vnc":
		return portscan.VncScan(u.Host, passwords)
	case "redis":
		return portscan.RedisScan(u.Host, passwords)
	case "memcached":
		return portscan.MemcachedScan(u.Host)
	case "mongodb":
		portscan.MongodbScan(u.Host)
	}
	return nil
}

// fofa

func (a *App) FofaTips(query string) *space.TipsResult {
	var ff space.FofaConfig
	var ts space.TipsResult
	ff.AppId = "9e9fb94330d97833acfbc041ee1a76793f1bc691"
	ff.PrivateKey = `MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC/TGN5+4FMXo7H3jRmostQUUEO1NwH10B8ONaDJnYDnkr5V0ZzUvkuola7JGSFgYVOUjgrmFGITG+Ne7AgR53Weiunlwp15MsnCa8/IWBoSHs7DX1O72xNHmEfFOGNPyJ4CsHaQ0B2nxeijs7wqKGYGa1snW6ZG/ZfEb6abYHI9kWVN1ZEVTfygI+QYqWuX9HM4kpFgy/XSzUxYE9jqhiRGI5f8SwBRVp7rMpGo1HZDgfMlXyA5gw++qRq7yHA3yLqvTPSOQMYJElJb12NaTcHKLdHahJ1nQihL73UwW0q9Zh2c0fZRuGWe7U/7Bt64gV2na7tlA62A9fSa1Dbrd7lAgMBAAECggEAPrsbB95MyTFc2vfn8RxDVcQ/dFCjEsMod1PgLEPJgWhAJ8HR7XFxGzTLAjVt7UXK5CMcHlelrO97yUadPAigHrwTYrKqEH0FjXikiiw0xB24o2XKCL+EoUlsCdg8GqhwcjL83Mke84c6Jel0vQBfdVQ+RZbetMCxqv1TpqpwW+iswlDY0+OKNxcDSnUyVkBko4M7bCqJ19DjzuHHLRmSuJhWLjX2PzdrVwIrRChxeJRR5AzrNE2BC/ssKasWjZfgkTOW6MS96q+wMLgwFGCQraU0f4AW5HA4Svg8iWT2uukcDg7VXXc/eEmkfmDGzmgsszUJZYb1hYsvjgbMP1ObwQKBgQDw1K0xfICYctiZ3aHS7mOk0Zt6B/3rP2z9GcJVs0eYiqH+lteLNy+Yx4tHtrQEuz16IKmM1/2Ghv8kIlOazpKaonk3JEwm1mCEXpgm4JI7UxPGQj/pFTCavKBBOIXxHJVSUSg0nKFkJVaoJiNy0CKwQNoFGdROk2fSYu8ReB/WlQKBgQDLWQR3RioaH/Phz8PT1ytAytH+W9M4P4tEx/2Uf5KRJxPQbN00hPnK6xxHAqycTpKkLkbJIkVWEKcIGxCqr6iGyte3xr30bt49MxIAYrdC0LtBLeWIOa88GTqYmIusqJEBmiy+A+DudM/xW4XRkgrOR1ZsagzI3FUVlei9DwFjEQKBgG8JH3EZfhDLoqIOVXXzA24SViTFWoUEETQAlGD+75udD2NaGLbPEtrV5ZmC2yzzRzzvojyVuQY1Z505VmKhq2YwUsLhsVqWrJlbI7uI/uLrQsq98Ml+Q5KUNS7c6KRqEU6KrIbVUHPj4zhTnTRqUhQBUoPXjNNNkyilBKSBReyhAoGAd3xGCIPdB17RIlW/3sFnM/o5bDmuojWMcw0ErvZLPCl3Fhhx3oNod9iw0/T5UhtFRV2/0D3n+gts6nFk2LbA0vtryBvq0C85PUK+CCX5QzR9Y25Bmksy8aBtcu7n27ttAUEDm1+SEuvmqA68Ugl7efwnBytFed0lzbo5eKXRjdECgYAk6pg3YIPi86zoId2dC/KfsgJzjWKVr8fj1+OyInvRFQPVoPydi6iw6ePBsbr55Z6TItnVFUTDd5EX5ow4QU1orrEqNcYyG5aPcD3FXD0Vq6/xrYoFTjZWZx23gdHJoE8JBCwigSt0KFmPyDsN3FaF66Iqg3iBt8rhbUA8Jy6FQA==`
	b, err := ff.GetTips(query)
	if err != nil {
		return &ts
	}
	json.Unmarshal([]byte(string(b)), &ts)
	return &ts
}

func (a *App) FofaSearch(query, pageSzie, pageNum, email, api string, fraud, cert bool) *space.FofaSearchResult {
	return space.FofaApiSearch(query, pageSzie, pageNum, email, api, fraud, cert)
}

func (a *App) IconHash(target string) string {
	_, b, err := clients.NewRequest("GET", target, nil, nil, 10, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	return util.Mmh3Hash32(util.Base64Encode(b))
}

// infoscan
type AliveTarget struct {
	Status      bool
	ProtocolURL string
}

func (a *App) CheckTarget(host string, proxy clients.Proxy) *AliveTarget {
	var client *http.Client
	if proxy.Enabled {
		client, _ = clients.SelectProxy(&proxy, clients.DefaultClient())
	}
	protocalURL, err := clients.CheckProtocol(host, client)
	if err != nil {
		return &AliveTarget{
			Status:      false,
			ProtocolURL: host,
		}
	}
	return &AliveTarget{
		Status:      true,
		ProtocolURL: protocalURL,
	}
}

type InfoResult struct {
	URL          string
	StatusCode   int
	Length       int
	Title        string
	Fingerprints []string
}

var RuleData map[string]map[string]string

// 仅在执行时调用一次
func (a *App) InitRule() []byte {
	yamlData, err := os.ReadFile("./config/webfinger.yaml")
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	RuleData = make(map[string]map[string]string)
	return yamlData
}

func (a *App) FingerScan(url string, yamlData []byte, proxy clients.Proxy) *InfoResult {
	if err := yaml.Unmarshal(yamlData, &RuleData); err != nil {
		logger.NewDefaultLogger().Debug("Failed to unmarshal YAML: " + err.Error())
	}
	var ir InfoResult
	var fingerprints []string
	var matched bool // 判断指纹匹配状态
	var client = clients.DefaultClient()
	if proxy.Enabled {
		client, _ = clients.SelectProxy(&proxy, clients.DefaultClient())
	}
	data := webscan.RecvResponse(url, client)
	ir.URL = url
	ir.StatusCode = data.StatusCode
	if data.StatusCode != 0 { // 响应正常
		for name, ruleType := range RuleData {
			for types, rule := range ruleType {
				switch types {
				case "header":
					matched = webscan.MatchRule(rule, data.Headers)
				case "iconhash":
					matched = webscan.MatchRule(rule, data.FaviconHash)
				case "title":
					matched = webscan.MatchRule(rule, data.Title)
				default:
					matched = webscan.MatchRule(rule, util.Str2UTF8(string(data.Body)))
				}
				// 每个应用有一组指纹,一组指纹中只需要匹配一个即代表匹配
				if matched {
					fingerprints = append(fingerprints, name)
					matched = false
					continue
				}
			}
		}
		ir.Fingerprints = util.RemoveDuplicates[string](fingerprints)
		ir.Length = len(data.Body)
		ir.Title = data.Title
	}
	return &ir
}

// webscan

type WebResult struct {
	VulName  string
	Severity string
	VulURL   string
	Request  string
	Response string
	ExtInfo  string
}

func (a *App) PocNums(severity, keyword string) int {
	o := runner.NewOptions("", keyword, severity, "")
	return len(o.CreatePocList(a.LocalWalkFiles("./config/afrog-pocs")))
}

func (a *App) GetFingerPoc(fingerprints []string) []string {
	s, err := poc.FingerPocFilepath(fingerprints)
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	return s
}

func (a *App) Webscan(url, severity, keyword string, pocpathList []string, pr clients.Proxy) *[]WebResult {
	var proxys string
	var wr []WebResult
	if pr.Enabled {
		proxys = fmt.Sprintf("%v://%v:%v", strings.ToLower(pr.Mode), pr.Address, pr.Port)
	}
	options := runner.NewOptions(url, keyword, severity, proxys)
	r, err := runner.NewRunner(options, pocpathList)
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	var (
		lock   = sync.Mutex{}
		number uint32
	)
	r.OnResult = func(result *report.Result) {
		if result.IsVul {
			lock.Lock()
			atomic.AddUint32(&number, 1)

			extinfo := "" // 输出拓展信息
			if len(result.Extractor) > 0 {
				for _, v := range result.Extractor {
					switch value := v.Value.(type) {
					case map[string]string:
					case string:
						extinfo += "," + v.Key.(string) + "=\"" + fmt.Sprintf("%v", value) + "\""
					}
				}
				extinfo = "[" + strings.TrimLeft(extinfo, ",") + "]"
			}
			wr = append(wr, WebResult{
				VulName:  result.PocInfo.Id,
				Severity: strings.ToUpper(result.PocInfo.Info.Severity),
				VulURL:   result.FullTarget,
				Request:  string(result.AllPocResult[0].ResultRequest.Raw),
				Response: result.AllPocResult[0].ReadFullResultResponseInfo(),
				ExtInfo:  extinfo,
			})

			// r.Report.SetResult(result)
			// r.Report.Append(fmt.Sprint(number))
			lock.Unlock()
		}
	}
	r.Execute(url, pocpathList)
	return &wr
}

func (a *App) UpdatePocFile(latestVersion string) string {
	if err := update.UpdatePoc(latestVersion); err != nil {
		return err.Error()
	}
	os.RemoveAll("./config/afrog-pocs.zip")
	return ""
}

func (a *App) UpdateClinetFile(latestVersion string) string {
	if err := update.UpdateClinet(latestVersion); err != nil {
		return err.Error()
	}
	return ""
}

func (a *App) Restart() {
	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		logger.NewDefaultLogger().Fatal(err.Error())
	}
	// 退出当前的进程
	os.Exit(0)
}

// hunter

func (a *App) HunterTips(query string) *space.HunterTipsResult {
	var ts space.HunterTipsResult
	bs64 := base64.StdEncoding.EncodeToString([]byte(query))
	_, b, err := clients.NewRequest("GET", "https://hunter.qianxin.com/api/recommend?keyword="+bs64, nil, nil, 10, clients.DefaultClient())
	if err != nil {
		return &ts
	}
	json.Unmarshal([]byte(string(b)), &ts)
	return &ts
}
