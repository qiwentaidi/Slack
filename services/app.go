package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"slack-wails/core/dirsearch"
	"slack-wails/core/dumpall"
	"slack-wails/core/info"
	"slack-wails/core/isic"
	"slack-wails/core/jsfind"
	"slack-wails/core/portscan"
	"slack-wails/core/space"
	"slack-wails/core/subdomain"
	"slack-wails/core/webscan"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/netutil"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v2"
)

// App struct
type App struct {
	ctx              context.Context
	webfingerFile    string
	activefingerFile string
	cdnFile          string
	qqwryFile        string
	templateDir      string
	defaultPath      string
}

// NewApp creates a new App application struct
func NewApp() *App {
	home := util.HomeDir()
	return &App{
		webfingerFile:    home + "/slack/config/webfinger.yaml",
		activefingerFile: home + "/slack/config/dir.yaml",
		cdnFile:          home + "/slack/config/cdn.yaml",
		qqwryFile:        home + "/slack/config/qqwry.dat",
		templateDir:      home + "/slack/config/pocs",
		defaultPath:      home + "/slack/",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// 返回 true 将导致应用程序继续，false 将继续正常关闭
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	if !webscan.IsRunning {
		return false
	}
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Quit?",
		Message:       "Webscan is running are you sure you want to quit?",
		DefaultButton: "Confirm",
		CancelButton:  "Cancel",
		Buttons:       []string{"Confirm", "Cancel"},
	})
	if err != nil {
		return false
	}
	return dialog == "Cancel"
}

func (a *App) Callgologger(level, msg string) {
	switch level {
	case "info":
		gologger.Info(a.ctx, msg)
	case "warning":
		gologger.Warning(a.ctx, msg)
	case "error":
		gologger.Error(a.ctx, msg)
	case "success":
		gologger.Success(a.ctx, msg)
	default:
		gologger.Debug(a.ctx, msg)
	}
}

func (a *App) GoFetch(method, target string, body interface{}, headers map[string]string, timeout int, proxy clients.Proxy) *structs.Response {
	if _, err := url.Parse(target); err != nil {
		return &structs.Response{
			Error: true,
		}
	}
	var content []byte
	// 判断body的类型
	if data, ok := body.(map[string]interface{}); ok {
		content, _ = json.Marshal(data)
	} else {
		content = []byte(body.(string))
	}
	resp, b, err := clients.NewRequest(method, target, headers, bytes.NewReader(content), 10, true, clients.NewHttpClientWithProxy(nil, true, proxy))
	if err != nil || resp == nil {
		return &structs.Response{
			Error: true,
		}
	}
	headerMap := make(map[string]string)
	for key, values := range resp.Header {
		// 对于每个键，创建一个新的 map 并添加键值对
		headerMap["key"] = key
		headerMap["value"] = strings.Join(values, " ")
	}
	return &structs.Response{
		Error:     false,
		Proto:     resp.Proto,
		StatsCode: resp.StatusCode,
		Header:    headerMap,
		Body:      string(b),
	}
}

var CyberChefLoader sync.Once

func (a *App) CyberChefLocalServer() {
	CyberChefLoader.Do(func() {
		go func() {
			// 定义要服务的目录
			dir := util.HomeDir() + "/slack/CyberChef/"
			// 创建文件服务器
			fs := http.FileServer(http.Dir(dir))

			// 设置路由，将所有请求重定向到文件服务器
			http.Handle("/", fs)
			// 指定一个随机端口, 启动HTTP服务器
			err := http.ListenAndServe(fmt.Sprintf(":%d", 8731), nil)
			if err != nil {
				return
			}
		}()
	})
}

func (a *App) ICPInfo(domain string) string {
	name, icp, ip := info.SeoChinaz(a.ctx, domain)
	return fmt.Sprintf("公司名称: %v\n备案号: %v\nIP: %v", name, icp, ip)
}

func (App) Ip138IpHistory(domain string) string {
	return info.Ip138IpHistory(domain)
}

func (App) Ip138Subdomain(domain string) string {
	return info.Ip138Subdomain(domain)
}

var cdndataLoader sync.Once

func (a *App) CheckCdn(domain string) string {
	ips, cnames, err := netutil.Resolution(domain, subdomain.DefaultDnsServers, 10)
	if err != nil {
		return fmt.Sprintf("域名: %v 解析失败,%v", domain, err)
	}
	if len(ips) == 1 && len(cnames) == 0 {
		return fmt.Sprintf("域名: %v，解析到唯一IP %v，未解析到CNAME信息", domain, ips[0])
	}
	cdndataLoader.Do(func() {
		subdomain.Cdndata = netutil.ReadCDNFile(a.ctx, a.cdnFile)
	})
	ipList := strings.Join(ips, " | ")
	for _, cname := range cnames {
		for name, cdns := range subdomain.Cdndata {
			for _, cdn := range cdns {
				if strings.Contains(cname, cdn) {
					return fmt.Sprintf("域名: %v，识别到CDN域名，CNAME: %v CDN名称: %v 解析到IP为: %v", domain, cname, name, ipList)
				}
			}
		}
		if strings.Contains(cname, "cdn") {
			return fmt.Sprintf("域名: %v，CNAME中含有关键字: cdn，该域名可能使用了CDN技术 CNAME: %v 解析到IP为: %v", domain, cname, ipList)
		}
	}
	return fmt.Sprintf("域名: %v，未识别到CDN信息，解析到IP为: %v CNAME: %v", domain, ipList, strings.Join(cnames, ","))
}

var qqwryLoader sync.Once

func (a *App) IpLocation(ip string) string {
	qqwryLoader.Do(func() {
		subdomain.InitQqwry(a.qqwryFile)
	})
	result, err := subdomain.Database.Find(ip)
	if err != nil {
		return ""
	}
	return result.String()
}
func (a *App) Subdomain(o structs.SubdomainOption) {
	qqwryLoader.Do(func() {
		subdomain.InitQqwry(a.qqwryFile)
	})
	cdndataLoader.Do(func() {
		subdomain.Cdndata = netutil.ReadCDNFile(a.ctx, a.cdnFile)
	})
	subdomain.ExitFunc = false
	switch o.Mode {
	case structs.EnumerationMode:
		for _, domain := range o.Domains {
			subdomain.MultiThreadResolution(a.ctx, domain, []string{}, "Enumeration", o)
		}
	case structs.ApiMode:
		subdomain.ApiPolymerization(a.ctx, o)
	default:
		subdomain.ApiPolymerization(a.ctx, o)
		for _, domain := range o.Domains {
			subdomain.MultiThreadResolution(a.ctx, domain, []string{}, "Enumeration", o)
		}
	}
}

func (a *App) ExitScanner(scanType string) {
	switch scanType {
	case "[subdomain]":
		subdomain.ExitFunc = true
	case "[dirsearch]":
		dirsearch.ExitFunc = true
	case "[portscan]":
		portscan.ExitFunc = true
	case "[portbrute]":
		portscan.ExitBruteFunc = true
	case "[webscan]":
		webscan.ExitFunc = true
	}
}

func (a *App) SubsidiariesAndDomains(query string, subLevel, ratio int, searchDomain bool, machine string) []structs.CompanyInfo {
	tkm := info.CheckKeyMap(a.ctx, query)
	time.Sleep(util.SleepRandTime(1))
	result := info.SearchSubsidiary(a.ctx, tkm.CompanyName, tkm.CompanyId, ratio, false, searchDomain, machine)
	var secondCompanyNames []string
	if subLevel >= 2 {
		for _, r := range result {
			if r.CompanyName == tkm.CompanyName { // 跳过本级单位查询
				continue
			}
			secondCompanyNames = append(secondCompanyNames, r.CompanyName)
			secondResult := info.SearchSubsidiary(a.ctx, r.CompanyName, r.CompanyId, ratio, true, searchDomain, machine)
			result = append(result, secondResult...)
			time.Sleep(util.SleepRandTime(2))
		}
	}
	if subLevel == 3 {
		for _, r := range result {
			if util.ArrayContains(r.CompanyName, secondCompanyNames) { // 已经查询过的二级IP跳过
				continue
			}
			secondResult := info.SearchSubsidiary(a.ctx, r.CompanyName, r.CompanyId, ratio, true, searchDomain, machine)
			result = append(result, secondResult...)
			time.Sleep(util.SleepRandTime(1))
		}
	}
	return result
}

func (a *App) WechatOfficial(query string) []structs.WechatReulst {
	var companyId string
	for _, tkm := range info.TycKeyMap {
		if tkm.CompanyName == query {
			companyId = tkm.CompanyId
			break
		}
	}
	time.Sleep(time.Second)
	return info.WeChatOfficialAccounts(a.ctx, query, companyId)
}

func (a *App) TycCheckLogin(token string) bool {
	return info.CheckLogin(token)
}

// dirsearch
func (a *App) LoadDirsearchDict(dictPath, newExts []string) []string {
	var dicts []string
	for _, dict := range dictPath {
		dicts = append(dicts, util.LoadDirsearchDict(a.ctx, dict, "%EXT%", newExts)...)
	}
	return util.RemoveDuplicates(dicts)
}

func (a *App) DirScan(options dirsearch.Options) {
	dirsearch.ExitFunc = false
	dirsearch.NewScanner(a.ctx, options)
}

// portscan

func (a *App) HostAlive(targets []string, Ping bool) []string {
	return portscan.CheckLive(a.ctx, targets, Ping)
}

func (a *App) SpaceGetPort(ip string) []float64 {
	return space.GetShodanAllPort(a.ctx, ip)
}

func (a *App) NewTcpScanner(specialTargets []string, ips []string, ports []int, thread, timeout int, proxy clients.Proxy) {
	portscan.ExitFunc = false
	addresses := make(chan portscan.Address)

	go func() {
		defer close(addresses)
		// Generate addresses from ips and ports
		for _, ip := range ips {
			for _, port := range ports {
				addresses <- portscan.Address{IP: ip, Port: port}
			}
		}
		// Generate addresses from special targets
		for _, target := range specialTargets {
			temp := strings.Split(target, ":")
			port, err := strconv.Atoi(temp[1]) // Skip if port conversion fails
			if err != nil {
				continue
			}
			addresses <- portscan.Address{IP: temp[0], Port: port}
		}
	}()
	portscan.TcpScan(a.ctx, addresses, thread, timeout, proxy)
}

// 端口暴破
func (a *App) PortBrute(host string, usernames, passwords []string) {
	portscan.ExitBruteFunc = false
	portscan.PortBrute(a.ctx, host, usernames, passwords)
}

// fofa

func (a *App) FofaTips(query string) *structs.TipsResult {
	config := space.NewFofaConfig(nil)
	b, err := config.GetTips(query)
	if err != nil {
		gologger.Debug(a.ctx, err)
		return nil
	}
	var ts structs.TipsResult
	json.Unmarshal(b, &ts)
	return &ts
}

func (a *App) FofaSearch(query, pageSzie, pageNum, address, email, key string, fraud, cert bool) *structs.FofaSearchResult {
	config := space.NewFofaConfig(&structs.FofaAuth{
		Address: address,
		Email:   email,
		Key:     key,
	})
	return config.FofaApiSearch(a.ctx, query, pageSzie, pageNum, fraud, cert)
}

func (a *App) Socks5Conn(ip string, port, timeout int, username, password string) bool {
	return portscan.Socks5Conn(ip, port, timeout, username, password)
}

func (a *App) IconHash(target string) string {
	_, b, err := clients.NewSimpleGetRequest(target, clients.NewHttpClient(nil, true))
	if err != nil {
		return ""
	}
	return webscan.Mmh3Hash32(b)
}

// 仅在执行时调用一次
func (a *App) InitRule(appendTemplateFolder string) bool {
	templateFolders := []string{a.templateDir, appendTemplateFolder}
	config := webscan.NewConfig(templateFolders)
	return config.InitAll(a.ctx, a.webfingerFile, a.activefingerFile)
}

// webscan

func (a *App) FingerprintList() []string {
	var fingers []string
	for _, item := range webscan.FingerprintDB {
		fingers = append(fingers, item.ProductName)
	}
	return fingers
}

// 多线程 Nuclei 扫描，由于Nucli的设计问题，多线程无法调用代理，否则会导致扫描失败
func (a *App) NewThreadSafeWebScanner(options structs.WebscanOptions, proxy clients.Proxy) {
	webscan.ExitFunc = false
	webscan.IsRunning = true
	gologger.Info(a.ctx, fmt.Sprintf("Load web scanner, targets number: %d", len(options.Target)))
	gologger.Info(a.ctx, "Fingerscan is running ...")
	engine := webscan.NewFingerScanner(a.ctx, proxy, options)
	if engine == nil {
		gologger.Error(a.ctx, "Init fingerscan engine failed")
		webscan.IsRunning = false
		return
	}
	engine.NewFingerScan()
	if options.DeepScan && !webscan.ExitFunc {
		engine.NewActiveFingerScan()
	}
	if options.CallNuclei && !webscan.ExitFunc {
		gologger.Info(a.ctx, "Init nuclei engine, vulnerability scan is running ...")
		var allTemplateFolders = []string{a.templateDir}
		if options.AppendTemplateFolder != "" {
			allTemplateFolders = append(allTemplateFolders, options.AppendTemplateFolder)
		}
		fpm := engine.URLWithFingerprintMap()
		allOptions := []structs.NucleiOption{}
		for target, tags := range fpm {
			option := structs.NucleiOption{
				URL:                   target,
				Tags:                  util.RemoveDuplicates(tags),
				TemplateFile:          options.TemplateFiles,
				SkipNucleiWithoutTags: options.SkipNucleiWithoutTags,
				TemplateFolders:       allTemplateFolders,
				CustomTags:            options.Tags,
				CustomHeaders:         options.CustomHeaders,
			}
			allOptions = append(allOptions, option)
		}
		webscan.NewThreadSafeNucleiEngine(a.ctx, allOptions)
		gologger.Info(a.ctx, "Vulnerability scan has ended")
	}
	webscan.IsRunning = false
}

func (a *App) NewWebScanner(options structs.WebscanOptions, proxy clients.Proxy) {
	webscan.ExitFunc = false
	webscan.IsRunning = true
	gologger.Info(a.ctx, fmt.Sprintf("Load web scanner, targets number: %d", len(options.Target)))
	gologger.Info(a.ctx, "Fingerscan is running ...")
	engine := webscan.NewFingerScanner(a.ctx, proxy, options)
	if engine == nil {
		gologger.Error(a.ctx, "Init fingerscan engine failed")
		webscan.IsRunning = false
		return
	}
	engine.NewFingerScan()
	if options.DeepScan && !webscan.ExitFunc {
		engine.NewActiveFingerScan()
	}
	if options.CallNuclei && !webscan.ExitFunc {
		gologger.Info(a.ctx, "Init nuclei engine, vulnerability scan is running ...")
		var id = 0
		var allTemplateFolders = []string{a.templateDir}
		if options.AppendTemplateFolder != "" {
			allTemplateFolders = append(allTemplateFolders, options.AppendTemplateFolder)
		}
		fpm := engine.URLWithFingerprintMap()
		count := len(fpm)
		runtime.EventsEmit(a.ctx, "NucleiCounts", count)

		for target, tags := range fpm {
			if webscan.ExitFunc {
				gologger.Warning(a.ctx, "User exits vulnerability scanning")
				webscan.IsRunning = false
				return
			}
			id++
			gologger.Info(a.ctx, fmt.Sprintf("vulnerability scanning %d/%d", id, count))
			webscan.NewNucleiEngine(a.ctx, structs.NucleiOption{
				URL:                   target,
				Tags:                  util.RemoveDuplicates(tags),
				TemplateFile:          options.TemplateFiles,
				SkipNucleiWithoutTags: options.SkipNucleiWithoutTags,
				TemplateFolders:       allTemplateFolders,
				CustomTags:            options.Tags,
				CustomHeaders:         options.CustomHeaders,
				Proxy:                 clients.GetRawProxy(proxy),
			})
			runtime.EventsEmit(a.ctx, "NucleiProgressID", id)
		}
		gologger.Info(a.ctx, "Vulnerability scan has ended")
	}
	webscan.IsRunning = false
}

func (a *App) GetFingerPocMap() map[string][]string {
	return webscan.WorkFlowDB
}

func (a *App) GetAllFinger() []string {
	data, err := os.ReadFile(a.webfingerFile)
	if err != nil {
		return []string{}
	}

	fps := make(map[string]interface{})
	if err := yaml.Unmarshal(data, &fps); err != nil {
		return []string{}
	}
	var fingerprints []string
	for fingerprint := range fps {
		fingerprints = append(fingerprints, fingerprint)
	}
	return util.RemoveDuplicates(fingerprints)
}

// hunter

func (a *App) HunterTips(query string) *structs.HunterTips {
	return space.SearchHunterTips(query)
}

func (a *App) HunterSearch(api, query, pageSize, pageNum, times, asset string, deduplication bool) *structs.HunterResult {
	hr := space.HunterApiSearch(api, query, pageSize, pageNum, times, asset, deduplication)
	time.Sleep(time.Second * 2)
	return hr
}

// quake

func (a *App) QuakeTips(query string) *structs.QuakeTipsResult {
	return space.SearchQuakeTips(query)
}

func (a *App) QuakeSearch(ipList []string, query string, pageNum, pageSize int, latest, invalid, honeypot, cdn bool, token, certcommon string) *structs.QuakeResult {
	option := structs.QuakeRequestOptions{
		IpList:     ipList,
		Query:      query,
		PageNum:    pageNum,
		PageSize:   pageSize,
		Latest:     latest,
		Invalid:    invalid,
		Honeypot:   honeypot,
		CDN:        cdn,
		Token:      token,
		CertCommon: certcommon,
	}
	qk := space.QuakeApiSearch(&option)
	time.Sleep(time.Second * 1)
	return qk
}

func (a *App) ExtractAllJSLink(url string) []string {
	return jsfind.ExtractJS(a.ctx, url)
}

func (a *App) JSFind(target string, jsLinks []string) (fs *structs.FindSomething) {
	return jsfind.MultiThreadJSFind(a.ctx, target, jsLinks)
}

func (a *App) AnalyzeAPI(homeURL, baseURL string, apiList []string, headers map[string]string) {
	jsfind.AnalyzeAPI(a.ctx, homeURL, baseURL, apiList, headers)
}

// 允许目标传入文件或者目标favicon地址
func (a *App) FaviconMd5(target string) string {
	hasher := md5.New()
	if _, err := os.Stat(target); err != nil {
		_, body, err := clients.NewSimpleGetRequest(target, clients.NewHttpClient(nil, true))
		if err != nil {
			return ""
		}
		hasher.Write(body)
	} else {
		content, err := os.ReadFile(target)
		if err != nil {
			return ""
		}
		hasher.Write(content)
	}
	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum)
}

func (a *App) UncoverSearch(query, types string, option structs.SpaceOption) []space.Result {
	return space.Uncover(a.ctx, query, types, option)
}

func (a *App) GitDorks(target, dork, apikey string) *isic.GithubResult {
	return isic.GithubApiQuery(a.ctx, fmt.Sprintf("%s %s", target, dork), apikey)
}

func (a *App) ViewPictrue(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		gologger.Debug(a.ctx, err)
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
}

func (a *App) NetDial(host string) bool {
	_, err := net.Dial("tcp", host)
	return err == nil
}

func (a *App) NewDSStoreEngine(url string) []string {
	url = strings.TrimSpace(url)
	links, err := dumpall.ExtractDSStore(url)
	if err != nil {
		gologger.Debug(a.ctx, err)
		return nil
	}
	return links
}
