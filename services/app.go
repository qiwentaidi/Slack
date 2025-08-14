package services

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slack-wails/core/dirsearch"
	"slack-wails/core/dumpall"
	"slack-wails/core/info/icp"
	"slack-wails/core/info/tianyancha"
	"slack-wails/core/isic"
	"slack-wails/core/jsfind"
	"slack-wails/core/portscan"
	"slack-wails/core/repeater"
	"slack-wails/core/space"
	"slack-wails/core/subdomain"
	"slack-wails/core/webscan"
	"slack-wails/lib/control"
	"slack-wails/lib/gologger"
	"slack-wails/lib/gomessage"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils"
	"slack-wails/lib/utils/arrayutil"
	"slack-wails/lib/utils/fileutil"
	"slack-wails/lib/utils/httputil"
	"slack-wails/lib/utils/netutil"
	"slack-wails/lib/utils/randutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qiwentaidi/clients"
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
	cyberCherDir     string
}

// NewApp creates a new App application struct
func NewApp() *App {
	home := utils.HomeDir()
	return &App{
		webfingerFile:    home + "/slack/config/webfinger.yaml",
		activefingerFile: home + "/slack/config/dir.yaml",
		cdnFile:          home + "/slack/config/cdn.yaml",
		qqwryFile:        home + "/slack/config/qqwry.dat",
		templateDir:      home + "/slack/config/pocs",
		defaultPath:      home + "/slack/",
		cyberCherDir:     filepath.Join(home, "slack", "CyberChef"),
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

func (a *App) GoFetch(method, target string, body interface{}, headers map[string]string, timeout int, proxyURL string) *structs.Response {
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
	resp, err := clients.DoRequest(method, target, headers, bytes.NewReader(content), 10, clients.NewRestyClientWithProxy(nil, true, proxyURL))
	if err != nil {
		return &structs.Response{
			Error: true,
		}
	}
	headerMap := make(map[string]string)
	for key, values := range resp.Header() {
		// 对于每个键，创建一个新的 map 并添加键值对
		headerMap["key"] = key
		headerMap["value"] = strings.Join(values, " ")
	}
	return &structs.Response{
		Error:     false,
		Proto:     resp.Proto(),
		StatsCode: resp.StatusCode(),
		Header:    headerMap,
		Body:      string(resp.Body()),
	}
}

var CyberChefLoader sync.Once

func (a *App) CyberChefLocalServer() {
	CyberChefLoader.Do(func() {
		go func() {
			// 定义要服务的目录
			// 创建文件服务器
			fs := http.FileServer(http.Dir(a.cyberCherDir))
			// 创建独立的 ServeMux
			mux := http.NewServeMux()
			mux.Handle("/", fs)

			err := http.ListenAndServe(fmt.Sprintf(":%d", 8731), mux)
			if err != nil {
				return
			}
		}()
	})
}

var qqwryLoader sync.Once

func (a *App) IpLocation(ip string) string {
	qqwryLoader.Do(func() {
		subdomain.InitQqwry(a.ctx, a.qqwryFile)
	})
	result, err := subdomain.Database.Find(ip)
	if err != nil {
		return ""
	}
	return result.String()
}

var cdndataLoader sync.Once

func (a *App) Subdomain(o structs.SubdomainOption) {
	ctrlCtx, _ := control.GetScanContext(control.Subdomain) // 标识任务
	qqwryLoader.Do(func() {
		subdomain.InitQqwry(a.ctx, a.qqwryFile)
	})
	cdndataLoader.Do(func() {
		subdomain.Cdndata = netutil.ReadCDNFile(a.ctx, a.cdnFile)
	})
	engine := subdomain.NewSubdomainEngine(a.ctx, o)
	switch o.Mode {
	case structs.EnumerationMode:
		for _, domain := range o.Domains {
			engine.Runner(ctrlCtx, domain, []string{}, "Enumeration")
		}
	case structs.ApiMode:
		engine.ApiPolymerization(ctrlCtx)
	case structs.MixedMode:
		engine.ApiPolymerization(ctrlCtx)
		for _, domain := range o.Domains {
			engine.Runner(ctrlCtx, domain, []string{}, "Enumeration")
		}
	default:
		engine.Runner(ctrlCtx, "", []string{}, "")
	}
}

func (a *App) ExitScanner(scanType string) {
	switch scanType {
	case "[subdomain]":
		control.CancelScanContext(control.Subdomain)
	case "[dirsearch]":
		control.CancelScanContext(control.Dirseach)
	case "[portscan]":
		control.CancelScanContext(control.Portscan)
		control.CancelScanContext(control.Crack)
	case "[webscan]":
		control.CancelScanContext(control.Webscan)
	}
}

func (a *App) FetchCompanyInfo(companyName string, ratio int, ds *structs.DataSource, maxDepth int) structs.CompanyInfo {
	var result structs.CompanyInfo
	if ds.Tianyancha.Enable {
		tyc := tianyancha.NewClient(a.ctx, ds.Tianyancha.Token, ds.Tianyancha.Token)
		if tyc.CheckLogin() {
			info, err := a.fetchCompanyRecursiveByTianyancha(tyc, companyName, ratio, 1, maxDepth)
			if err != nil {
				gologger.Error(a.ctx, fmt.Sprintf("[tianyancha] fetch company info error: %s", err.Error()))
			}
			a.WriteCompanyInfoToJson(info)
			result = info
		} else {
			gomessage.Warning(a.ctx, "tianyancha token is invalid")
			gologger.Warning(a.ctx, "tianyancha token is invalid")
		}
	} else {
		result.CompanyName = companyName
	}

	if ds.Miit.API != "" {
		ds.Miit.API = strings.TrimRight(ds.Miit.API, "/")
		// 给主公司填充ICP信息
		a.EnrichCompanyWithMiit(&result, ds.Miit.API)
		time.Sleep(randutil.SleepRandTime(2))
		// 给所有子公司递归填充
		var enrichSubsidiaries func(subs []structs.CompanyInfo)
		enrichSubsidiaries = func(subs []structs.CompanyInfo) {
			for i := range subs {
				a.EnrichCompanyWithMiit(&subs[i], ds.Miit.API)
				if len(subs[i].Subsidiaries) > 0 {
					enrichSubsidiaries(subs[i].Subsidiaries)
				}
			}
		}
		enrichSubsidiaries(result.Subsidiaries)
	}
	a.WriteCompanyInfoToJson(result)
	return result
}

func (a *App) EnrichCompanyWithMiit(company *structs.CompanyInfo, miitApi string) {
	// 吊销或注销则跳过
	if company.RegStatus == "吊销" || company.RegStatus == "注销" {
		return
	}
	gologger.Info(a.ctx, fmt.Sprintf("[icp] 正在查询%s域名信息", company.CompanyName))
	if webResp, err := icp.FetchWebInfo(a.ctx, miitApi, company.CompanyName); err == nil {
		var domains []string
		for _, data := range webResp.Params.List {
			domains = append(domains, data.Domain)
		}
		company.Domains = domains
	} else {
		gologger.Warning(a.ctx, fmt.Sprintf("%s fetch web info error: %s", company.CompanyName, err))
	}

	time.Sleep(2 * time.Second)
	gologger.Info(a.ctx, fmt.Sprintf("[icp] 正在查询%sApp信息", company.CompanyName))
	if appResp, err := icp.FetchAppInfo(a.ctx, miitApi, company.CompanyName); err == nil {
		company.Apps = appResp.Params.List
	} else {
		gologger.Warning(a.ctx, fmt.Sprintf("%s fetch app info error: %s", company.CompanyName, err))
	}

	time.Sleep(2 * time.Second)
	gologger.Info(a.ctx, fmt.Sprintf("[icp] 正在查询%s小程序信息", company.CompanyName))
	if appletResp, err := icp.FetchAppletInfo(a.ctx, miitApi, company.CompanyName); err == nil {
		company.Applets = appletResp.Params.List
	} else {
		gologger.Warning(a.ctx, fmt.Sprintf("%s fetch applet info error: %s", company.CompanyName, err))
	}
}

func (a *App) ResumeAfterHumanCheck() {
	go func() {
		tianyancha.HumanCheckChan <- struct{}{}
	}()
}

// func (a *App) fetchCompanyRecursiveByRiskbird(rb *riskbird.RiskbirdClient, company string, ratio int, currentDepth, maxDepth int) (structs.CompanyInfo, error) {
// 	var companyInfo structs.CompanyInfo

// 	info, orderNo, err := rb.FetchBasicCompanyInfo(company)
// 	if err != nil {
// 		return companyInfo, err
// 	}
// 	companyInfo = info

// 	// Step 2: 查询 App、小程序、公众号等
// 	if apps, err := rb.FetchApp(orderNo); err == nil {
// 		companyInfo.Apps = apps
// 	}
// 	time.Sleep(1 * time.Second)

// 	if applets, err := rb.FetchApplet(orderNo); err == nil {
// 		// applets 需要转换为 OfficialAccounts，如果你有对应函数可以调用，否则跳过
// 		for _, ap := range applets {
// 			companyInfo.OfficialAccounts = append(companyInfo.OfficialAccounts, structs.OfficialAccount{
// 				Name: ap.Name, Logo: ap.Logo, Qrcode: fmt.Sprintf("%v", ap.Qrcode),
// 			})
// 		}
// 	}
// 	time.Sleep(1 * time.Second)

// 	// Step 3: 查询子公司
// 	subs, err := rb.FetchSubsidiary(orderNo)
// 	if err == nil && currentDepth <= maxDepth {
// 		for _, sub := range subs {
// 			gq, _ := strconv.Atoi(strings.TrimSuffix(sub.FunderRatio, "%"))
// 			if gq < ratio {
// 				continue
// 			}
// 			child, err := a.fetchCompanyRecursiveByRiskbird(rb, sub.EntName, ratio, currentDepth+1, maxDepth)
// 			if err != nil {
// 				gologger.Error(a.ctx, fmt.Sprintf("[riskbird] %s fetch sub error: %s", sub.EntName, err.Error()))
// 				continue
// 			}
// 			child.Investment = sub.FunderRatio
// 			child.Amount = sub.RegCapFormat
// 			child.RegStatus = sub.EntStatus
// 			companyInfo.Subsidiaries = append(companyInfo.Subsidiaries, child)
// 		}
// 	}

//		return companyInfo, nil
//	}
func (a *App) fetchCompanyRecursiveByTianyancha(tyc *tianyancha.TycClient, company string, ratio int, currentDepth, maxDepth int) (structs.CompanyInfo, error) {
	var companyInfo structs.CompanyInfo

	// Step 1: 获取公司基本信息
	suggest, err := tyc.CheckKeyMap(company)
	if err != nil {
		return companyInfo, err
	}
	companyInfo.CompanyName = suggest.ComName
	companyInfo.Investment = "母公司"
	companyInfo.RegStatus = tyc.GetRegStatus(suggest.RegStatus)
	companyInfo.Trademark = suggest.Logo

	// Step 2: 非注销/吊销状态的公司需要获取公众号信息
	if officialAccounts, err := tyc.FetchWeChatOfficialAccounts(suggest.ComName, suggest.GraphID); err == nil {
		companyInfo.OfficialAccounts = officialAccounts
	}
	time.Sleep(randutil.SleepRandTime(2))

	// Step 3: 获取子公司信息
	subsidiaries, err := tyc.FetchSubsidiary(suggest.ComName, suggest.GraphID, ratio)
	if err != nil {
		return companyInfo, err
	}

	for _, subs := range subsidiaries {
		child := structs.CompanyInfo{
			CompanyName: subs.Name,
			Investment:  subs.Percent,
			Amount:      subs.Amount,
			RegStatus:   subs.RegStatus,
			Trademark:   fmt.Sprint(subs.Logo),
		}
		// 跳过已注销或吊销的子公司
		if suggest.RegStatus != 1 && suggest.RegStatus != 2 && currentDepth < maxDepth {
			for {
				time.Sleep(2 * time.Second)
				subInfo, err := a.fetchCompanyRecursiveByTianyancha(tyc, subs.Name, ratio, currentDepth+1, maxDepth)
				if err != nil && strings.Contains(err.Error(), "账号存在风险请人机验证") {
					// 通知前端进行人机验证
					runtime.EventsEmit(a.ctx, "tyc-human-check", "天眼查出现人机校验，请手动处理")
					gologger.DualLog(a.ctx, gologger.Level_DEBUG, "天眼查出现人机校验，请手动处理")
					<-tianyancha.HumanCheckChan
					gologger.DualLog(a.ctx, gologger.Level_DEBUG, "收到用户确认，继续查询")
					continue // 重试
				}
				if err == nil {
					child.OfficialAccounts = subInfo.OfficialAccounts
					child.Subsidiaries = subInfo.Subsidiaries
				}
				break
			}
		}

		// 当前深度已达最大，或者递归处理完都要追加子公司
		companyInfo.Subsidiaries = append(companyInfo.Subsidiaries, child)
	}

	return companyInfo, nil
}

var companyPath = filepath.Join(utils.HomeDir(), "slack", "company_info")

func (a *App) WriteCompanyInfoToJson(info structs.CompanyInfo) bool {
	os.Mkdir(companyPath, 0777)
	fp := filepath.Join(companyPath, fmt.Sprintf("%s-%s.json", info.CompanyName, info.RegStatus))
	return fileutil.SaveJsonWithFormat(a.ctx, fp, info)
}

// dirsearch
func (a *App) LoadDirsearchDict(dictPath, newExts []string) []string {
	var dicts []string
	for _, dict := range dictPath {
		dicts = append(dicts, LoadDirDict(dict, "%EXT%", newExts)...)
	}
	return arrayutil.RemoveDuplicates(dicts)
}

func LoadDirDict(filepath, old string, new []string) (dict []string) {
	file, _ := os.Open(filepath)
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" { // 去除空行
			if len(new) > 0 {
				if strings.Contains(s.Text(), old) { // 如何新数组不为空,将old字段替换成new数组
					for _, n := range new {
						dict = append(dict, strings.ReplaceAll(s.Text(), old, n))
					}
				} else {
					dict = append(dict, s.Text())
				}
			} else {
				if !strings.Contains(s.Text(), old) {
					dict = append(dict, s.Text())
				}
			}
		}
	}
	return dict
}

func (a *App) NewDirsearchScanner(options dirsearch.Options) {
	ctrlCtx, _ := control.GetScanContext(control.Dirseach) // 标识任务
	engine := dirsearch.NewDirsearchEngine(a.ctx, ctrlCtx, options)
	if options.Backupscan {
		engine.BackupRunner(ctrlCtx)
	} else {
		engine.Runner(ctrlCtx)
	}
}

// portscan

func (a *App) HostAlive(targets []string, Ping bool) []string {
	return portscan.CheckLive(a.ctx, targets, Ping)
}

func (a *App) SpaceGetPort(ip string) []float64 {
	return space.GetShodanAllPort(a.ctx, ip)
}

func (a *App) NewTcpScanner(taskId string, specialTargets []string, ips []string, ports []int, thread, timeout int, proxyURL string) {
	ctrlCtx, _ := control.GetScanContext(control.Portscan) // 标识任务
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
	portscan.TcpScan(a.ctx, ctrlCtx, taskId, addresses, thread, timeout, proxyURL)
}

// 端口暴破
func (a *App) NewCrackScanenr(taskId, host string, usernames, passwords []string) {
	ctrlCtx, _ := control.GetScanContext(control.Crack) // 标识任务
	portscan.Runner(a.ctx, ctrlCtx, taskId, host, usernames, passwords)
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

func (a *App) Socks5Conn(ip string, port, timeout int, username, password, aliveURL string) bool {
	return portscan.Socks5Conn(ip, port, timeout, username, password, aliveURL)
}

func (a *App) IconHash(target string) string {
	resp, err := clients.SimpleGet(target, clients.NewRestyClient(nil, true))
	if err != nil {
		return ""
	}
	return webscan.Mmh3Hash32(resp.Body())
}

// 仅在执行时调用一次
func (a *App) InitRule(appendTemplateFolder string) bool {
	templateFolders := []string{a.templateDir, appendTemplateFolder}
	config := &webscan.Config{
		TemplateFolders:     templateFolders,
		ActiveRuleFile:      a.activefingerFile,
		FingerprintRuleFile: a.webfingerFile,
	}
	return config.InitAll(a.ctx)
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
func (a *App) NewWebScanner(taskId string, options structs.WebscanOptions, proxyURL string, threadSafe bool) {
	ctrlCtx, cancel := control.GetScanContext(control.Webscan) // 标识任务
	defer cancel()
	webscan.IsRunning = true
	gologger.Info(a.ctx, fmt.Sprintf("Load web scanner, targets number: %d", len(options.Target)))
	gologger.Info(a.ctx, "Fingerscan is running ...")

	engine := webscan.NewWebscanEngine(a.ctx, taskId, proxyURL, options)
	if engine == nil {
		gologger.Error(a.ctx, "Init fingerscan engine failed")
		webscan.IsRunning = false
		return
	}

	// 指纹识别
	engine.FingerScan(ctrlCtx)
	if options.DeepScan && ctrlCtx.Err() == nil {
		engine.ActiveFingerScan(ctrlCtx)
	}

	if options.CallNuclei && ctrlCtx.Err() == nil {
		gologger.Info(a.ctx, "Init nuclei engine, vulnerability scan is running ...")

		// 准备模板目录
		var allTemplateFolders = []string{a.templateDir}
		if options.AppendTemplateFolder != "" {
			allTemplateFolders = append(allTemplateFolders, options.AppendTemplateFolder)
		}

		// 提取所有目标和标签
		fpm := engine.URLWithFingerprintMap()
		allOptions := []structs.NucleiOption{}
		for target, tags := range fpm {
			allOptions = append(allOptions, structs.NucleiOption{
				URL:                   target,
				Tags:                  arrayutil.RemoveDuplicates(tags),
				TemplateFile:          options.TemplateFiles,
				SkipNucleiWithoutTags: options.SkipNucleiWithoutTags,
				TemplateFolders:       allTemplateFolders,
				CustomTags:            options.Tags,
				CustomHeaders:         options.CustomHeaders,
				Proxy:                 proxyURL,
			})
		}
		counts := len(allOptions)
		if counts == 0 {
			gologger.Warning(a.ctx, "nuclei scan no targets")
			webscan.IsRunning = false
			return
		}
		runtime.EventsEmit(a.ctx, "NucleiCounts", counts)

		if threadSafe {
			webscan.NewThreadSafeNucleiEngine(a.ctx, ctrlCtx, taskId, allOptions)
		} else {
			webscan.NewNucleiEngine(a.ctx, ctrlCtx, taskId, allOptions)
		}

		gologger.Info(a.ctx, "Vulnerability scan has ended")
	}
	webscan.IsRunning = false
}

func (a *App) GetFingerPocMap() map[string][]string {
	return webscan.WorkFlowDB
}

// hunter

func (a *App) HunterTips(query string) *structs.HunterTips {
	return space.SearchHunterTips(query)
}

func (a *App) HunterSearch(api, key, query, pageSize, pageNum, times, asset string, deduplication bool) *structs.HunterResult {
	hr := space.HunterApiSearch(a.ctx, api, key, query, pageSize, pageNum, times, asset, deduplication)
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
	return jsfind.ExtractAllJs(a.ctx, url)
}

func (a *App) JSFind(target, prefixJsURL string, jsLinks, blackDomainList []string) structs.FindSomething {
	return jsfind.Scan(a.ctx, target, prefixJsURL, jsLinks, blackDomainList)
}

func (a *App) AnalyzeAPI(homeURL, baseURL string, apiList []string, headers, lowPrivilegeHeaders map[string]string, authentication []string, highRiskRouter []string) {
	options := structs.JSFindOptions{
		HomeURL:             homeURL,
		BaseURL:             baseURL,
		ApiList:             apiList,
		Headers:             headers,
		Authentication:      authentication,
		HighRiskRouter:      highRiskRouter,
		LowPrivilegeHeaders: lowPrivilegeHeaders,
	}
	jsfind.AnalyzeAPI(a.ctx, options)
}

// 允许目标传入文件或者目标favicon地址
func (a *App) FaviconMd5(target string) string {
	hasher := md5.New()
	if _, err := os.Stat(target); err != nil {
		resp, err := clients.SimpleGet(target, clients.NewRestyClient(nil, true))
		if err != nil {
			return ""
		}
		hasher.Write(resp.Body())
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

func (a *App) GitDorks(target, dork, apikey string) *structs.ISICollectionResult {
	return isic.GithubApiQuery(a.ctx, fmt.Sprintf("%s %s", target, dork), apikey)
}

func (a *App) GoogleHackerBingSearch(query string) *structs.ISICollectionResult {
	result, total, err := isic.GoogleHackerBingSearch(query)
	if err != nil {
		return nil
	}
	items := []string{}
	for _, item := range result {
		items = append(items, item.URL)
	}
	return &structs.ISICollectionResult{
		Items:  items,
		Link:   fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(query)),
		Source: "Bing",
		Total:  float64(total),
	}
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

func (a *App) SendRequest(raw string, forceHttps, redirect bool, proxyURL string) structs.RawResponse {
	resp, t, err := repeater.SendRequestWithRaw(raw, forceHttps, redirect, proxyURL)
	if err != nil {
		if errors.Is(err, http.ErrUseLastResponse) {
			return structs.RawResponse{
				StatusCode:   0,
				Error:        "",
				Response:     string(httputil.DumpResponseHeadersOnly(resp.RawResponse)),
				ResponseTime: 0,
			}
		}
		return structs.RawResponse{
			StatusCode:   0,
			Error:        err.Error(),
			Response:     "",
			ResponseTime: 0,
		}
	}
	rawReponse, err := httputil.DumpResponseHeadersAndDecodedBody(resp.RawResponse)
	if err != nil {
		return structs.RawResponse{
			StatusCode:   resp.StatusCode(),
			Error:        err.Error(),
			Response:     "",
			ResponseTime: t,
		}
	}
	return structs.RawResponse{
		Error:        "",
		Response:     rawReponse,
		ResponseTime: t,
		StatusCode:   resp.StatusCode(),
	}
}
