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
	"path"
	rt "runtime"
	"slack-wails/core/dirsearch"
	"slack-wails/core/exp/hikvision"
	"slack-wails/core/exp/nacos"
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

func (a *App) IsRoot() bool {
	if rt.GOOS == "windows" {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	} else {
		return os.Getuid() == 0
	}
}

func (a *App) GOOS() string {
	return rt.GOOS
}

func (a *App) GoFetch(method, target string, body interface{}, headers map[string]string, timeout int, proxy clients.Proxy) *structs.Response {
	if _, err := url.Parse(target); err != nil {
		return &structs.Response{
			Error:  true,
			Proto:  "",
			Header: nil,
			Body:   "",
		}
	}
	var content []byte
	// 判断body的类型
	if data, ok := body.(map[string]interface{}); ok {
		content, _ = json.Marshal(data)
	} else {
		content = []byte(body.(string))
	}
	resp, b, err := clients.NewRequest(method, target, headers, bytes.NewReader(content), 10, true, clients.DefaultWithProxyClient(proxy))
	if err != nil {
		return &structs.Response{
			Error:  true,
			Proto:  "",
			Header: nil,
			Body:   "",
		}
	}
	headerMap := make(map[string]string)
	for key, values := range resp.Header {
		// 对于每个键，创建一个新的 map 并添加键值对
		headerMap["key"] = key
		headerMap["value"] = strings.Join(values, " ")
	}
	return &structs.Response{
		Error:  false,
		Proto:  resp.Proto,
		Header: headerMap,
		Body:   string(b),
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

func (a *App) ExtractIP(text string) string {
	var result string
	var IP_analysis = make(map[string]int)
	result += "---提取IP资产---\n"
	for _, ip := range util.RemoveDuplicates(util.RegIP.FindAllString(text, -1)) {
		result += ip + "\n"
		ip = ip[:len(ip)-len(path.Ext(ip))]
		IP_analysis[ip+".0"]++
	}
	result += "\n\n\n---提取C段资产---\n"
	for _, p := range util.SortMap(IP_analysis) {
		result += fmt.Sprintf("%v/24(%v)\n", p.Key, p.Value)
	}
	return result
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
func (a *App) InitTycHeader(token string) {
	info.InitHEAD(token)
}

func (a *App) SubsidiariesAndDomains(query string, subLevel, ratio int, searchDomain bool, machine string) []structs.CompanyInfo {
	tkm := info.CheckKeyMap(a.ctx, query)
	time.Sleep(time.Second)
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
			time.Sleep(time.Second)
		}
	}
	if subLevel == 3 {
		for _, r := range result {
			if util.ArrayContains(r.CompanyName, secondCompanyNames) { // 已经查询过的二级IP跳过
				continue
			}
			secondResult := info.SearchSubsidiary(a.ctx, r.CompanyName, r.CompanyId, ratio, true, searchDomain, machine)
			result = append(result, secondResult...)
			time.Sleep(time.Second)
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

func (a *App) TycCheckLogin() bool {
	return info.CheckLogin()
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

func (a *App) IPParse(ipList []string) []string {
	return util.ParseIPs(ipList)
}

func (a *App) PortParse(text string) []int {
	return util.ParsePort(text)
}

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
	portscan.TcpScan(a.ctx, addresses, thread, timeout, &proxy)
}
func (a *App) NewSynScanner(specialTargets []string, ips []string, ports []uint16) {
	portscan.ExitFunc = false
	// Generate addresses from special targets
	var id int32 = 0
	for _, target := range specialTargets {
		if portscan.ExitFunc {
			return
		}
		temp := strings.Split(target, ":")
		port, err := strconv.Atoi(temp[1]) // Skip if port conversion fails
		if err != nil {
			continue
		}
		portscan.SynScan(a.ctx, temp[0], []uint16{uint16(port)}, &id)
	}
	// Generate addresses from ips and ports
	for _, ip := range ips {
		if portscan.ExitFunc {
			return
		}
		portscan.SynScan(a.ctx, ip, ports, &id)
	}
	runtime.EventsEmit(a.ctx, "scanComplete", id)

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
	_, b, err := clients.NewSimpleGetRequest(target, clients.DefaultClient())
	if err != nil {
		return ""
	}
	return webscan.Mmh3Hash32(b)
}

// infoscan

func (a *App) CheckTarget(host string, proxy clients.Proxy) *structs.Status {
	protocolURL, err := clients.IsWeb(host, clients.DefaultWithProxyClient(proxy))
	if err != nil {
		return &structs.Status{
			Error: true,
			Msg:   host,
		}
	}
	return &structs.Status{
		Error: false,
		Msg:   protocolURL,
	}
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

func (a *App) NewWebScanner(options structs.WebscanOptions, proxy clients.Proxy) {
	webscan.ExitFunc = false
	gologger.Info(a.ctx, fmt.Sprintf("Load web scanner, targets number: %d", len(options.Target)))
	gologger.Info(a.ctx, "Fingerscan is running ...")
	engine := webscan.NewFingerScanner(a.ctx, proxy, options)
	if engine == nil {
		return
	}
	engine.NewFingerScan()
	if options.DeepScan {
		engine.NewActiveFingerScan(options.RootPath)
	}
	if options.CallNuclei {
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
				return
			}
			id++
			gologger.Info(a.ctx, fmt.Sprintf("vulnerability scanning %d/%d", id, count))
			webscan.NewNucleiEngine(a.ctx, proxy, structs.NucleiOption{
				URL:                   target,
				Tags:                  util.RemoveDuplicates(tags),
				TemplateFile:          options.TemplateFiles,
				SkipNucleiWithoutTags: options.SkipNucleiWithoutTags,
				TemplateFolders:       allTemplateFolders,
			})
			runtime.EventsEmit(a.ctx, "NucleiProgressID", id)
		}
		gologger.Info(a.ctx, "Vulnerability scan has ended")
	}
}

func (a *App) GetFingerPocMap() map[string][]string {
	return webscan.WorkFlowDB
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

func (a *App) JSFind(target, customPrefix string) (fs *jsfind.FindSomething) {
	jsLinks := jsfind.ExtractJS(a.ctx, target)
	var wg sync.WaitGroup
	limiter := make(chan bool, 100)
	wg.Add(1)
	limiter <- true
	go func() {
		fs = jsfind.FindInfo(a.ctx, target, limiter, &wg)
	}()
	wg.Wait()
	u, _ := url.Parse(target)
	fs.JS = *jsfind.AppendSource(target, jsLinks)
	host := ""
	if customPrefix != "" {
		host = customPrefix
	} else {
		host = u.Scheme + "://" + u.Host
	}
	for _, jslink := range jsLinks {
		wg.Add(1)
		limiter <- true
		go func(js string) {
			var newURL string
			if strings.HasPrefix(js, "http") {
				newURL = js
			} else {
				newURL = host + "/" + js
			}
			fs2 := jsfind.FindInfo(a.ctx, newURL, limiter, &wg)
			fs.IP_URL = append(fs.IP_URL, fs2.IP_URL...)
			fs.ChineseIDCard = append(fs.ChineseIDCard, fs2.ChineseIDCard...)
			fs.ChinesePhone = append(fs.ChinesePhone, fs2.ChinesePhone...)
			fs.SensitiveField = append(fs.SensitiveField, fs2.SensitiveField...)
			fs.APIRoute = append(fs.APIRoute, fs2.APIRoute...)
		}(jslink)
	}
	wg.Wait()
	fs.APIRoute = jsfind.RemoveDuplicatesInfoSource(fs.APIRoute)
	fs.ChineseIDCard = jsfind.RemoveDuplicatesInfoSource(fs.ChineseIDCard)
	fs.ChinesePhone = jsfind.RemoveDuplicatesInfoSource(fs.ChinesePhone)
	fs.SensitiveField = jsfind.RemoveDuplicatesInfoSource(fs.SensitiveField)
	fs.IP_URL = jsfind.FilterExt(jsfind.RemoveDuplicatesInfoSource(fs.IP_URL))
	return fs
}

// 允许目标传入文件或者目标favicon地址
func (a *App) FaviconMd5(target string) string {
	hasher := md5.New()
	if _, err := os.Stat(target); err != nil {
		_, body, err := clients.NewSimpleGetRequest(target, clients.DefaultClient())
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

func (a *App) AlibabaNacos(target, headers string, attackType int, username, password, command, service string, proxy clients.Proxy) string {
	switch attackType {
	case 0:
		if nacos.CVE_2021_29441_Step1(target, username, password, clients.DefaultWithProxyClient(proxy)) {
			return "添加用户成功: \n username: " + username + "， password: " + password
		}
	case 1:
		if nacos.CVE_2021_29441_Step2(target, username, clients.DefaultWithProxyClient(proxy)) {
			return "删除用户成功!"
		}
	case 2:
		return nacos.CVE_2021_29442(target, clients.DefaultWithProxyClient(proxy))
	case 3:
		return nacos.DerbySqljinstalljarRCE(a.ctx, headers, target, command, service, clients.DefaultWithProxyClient(proxy))
	}
	return target + "不存在该漏洞"
}

func (a *App) NacosCategoriesExtract(filePath string) []structs.NacosConfig {
	return nacos.ProcessDirectory(filePath)
}

func (a *App) HikvsionCamera(target string, attackType int, passwordList []string, cmd string, proxy clients.Proxy) string {
	switch attackType {
	case 0:
		body := hikvision.CVE_2017_7921_Snapshot(target, clients.DefaultWithProxyClient(proxy))
		return base64.RawStdEncoding.EncodeToString(body)
	case 1:
		return hikvision.CVE_2017_7921_Config(target, clients.DefaultWithProxyClient(proxy))
	case 2:
		return hikvision.CVE_2021_36260(target, cmd, clients.DefaultWithProxyClient(proxy))
	case 3:
		return hikvision.CameraHandlessLogin(a.ctx, target, passwordList)
	}
	return ""
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
