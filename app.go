package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"slack-wails/core"
	"slack-wails/core/dirsearch"
	alibaba "slack-wails/core/druid"
	"slack-wails/core/info"
	"slack-wails/core/jsfind"
	"slack-wails/core/portscan"
	"slack-wails/core/space"
	"slack-wails/core/webscan"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

// App struct
type App struct {
	ctx              context.Context
	workflowFile     string
	webfingerFile    string
	activefingerFile string
	cdnFile          string
	qqwryFile        string
	avFile           string
	// bruteFile        string
	defaultPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	home := util.HomeDir()
	return &App{
		workflowFile:     home + "/slack/config/workflow.yaml",
		webfingerFile:    home + "/slack/config/webfinger.yaml",
		activefingerFile: home + "/slack/config/dir.yaml",
		cdnFile:          home + "/slack/config/cdn.yaml",
		qqwryFile:        home + "/slack/config/qqwry.dat",
		avFile:           home + "/slack/config/antivirues.yaml",
		// bruteFile:        home + "/slack/portburte",
		defaultPath: home + "/slack/",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
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
	default:
		gologger.Debug(a.ctx, msg)
	}
}

func (a *App) IsRoot() bool {
	return os.Getuid() == 0
}

func (a *App) GOOS() string {
	return runtime.GOOS
}

func (a *App) GoFetch(method, target string, body interface{}, headers []map[string]string, timeout int, proxy clients.Proxy) *structs.Response {
	if _, err := url.Parse(target); err != nil {
		return &structs.Response{
			Error:  true,
			Proto:  "",
			Header: nil,
			Body:   "",
		}
	}
	hhhhheaders := http.Header{}
	for _, head := range headers {
		for k, v := range head {
			hhhhheaders.Set(k, v)
		}
	}
	var content []byte
	// 判断body的类型
	if data, ok := body.(map[string]interface{}); ok {
		content, _ = json.Marshal(data)
	} else {
		content = []byte(body.(string))
	}
	resp, b, err := clients.NewRequest(method, target, hhhhheaders, bytes.NewReader(content), 10, true, clients.JudgeClient(proxy))
	if err != nil {
		return &structs.Response{
			Error:  true,
			Proto:  "",
			Header: nil,
			Body:   "",
		}
	}
	var headerArray []map[string]string
	for key, values := range resp.Header {
		// 对于每个键，创建一个新的 map 并添加键值对
		headerMap := make(map[string]string)
		headerMap["key"] = key
		headerMap["value"] = strings.Join(values, " ")

		// 将 map 添加到切片中
		headerArray = append(headerArray, headerMap)
	}
	return &structs.Response{
		Error:  false,
		Proto:  resp.Proto,
		Header: headerArray,
		Body:   string(b),
	}
}

// fscan
func (a *App) Fscan2Txt(content string) string {
	return core.ExtractFscanResult(content)
}

// thinkdict
func (a *App) ThinkDict(userNameCN, userNameEN, companyName, companyDomain, birthday, jobNumber, connectWord string, weakList []string) []string {
	return core.GenerateDict(userNameCN, userNameEN, companyName, companyDomain, birthday, jobNumber, connectWord, weakList)
}

func (a *App) System(content string, mode int) [][]string {
	if mode == 0 { // tasklist
		datas, _ := core.AntivirusIdentify(content, a.avFile)
		return datas
	} else { // systeminfo
		return core.Patch(content)
	}
}

var onec sync.Once

func (a *App) CyberChefLocalServer() {
	onec.Do(func() {
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
func (a *App) CheckCdn(domain string) string {
	ips, cnames, err := core.Resolution(domain, 10)
	if err != nil {
		return fmt.Sprintf("域名: %v 解析失败,%v", domain, err)
	}
	if len(cnames) == 0 {
		return fmt.Sprintf("域名: %v 未解析到CNAME信息", domain)
	}
	cdnData := core.ReadCDNFile(a.ctx, a.cdnFile)
	ipList := strings.Join(ips, " | ")
	for _, cname := range cnames {
		for name, cdns := range cdnData {
			for _, cdn := range cdns {
				if strings.Contains(cname, cdn) {
					return fmt.Sprintf("识别到CDN域名，CNAME: %v CDN名称: %v 解析到IP为: %v", cname, name, ipList)
				}
			}
		}
		if strings.Contains(cname, "cdn") {
			return fmt.Sprintf("CNAME中含有关键字: cdn，该域名可能使用了CDN技术 CNAME: %v 解析到IP为: %v", cname, ipList)
		}
	}
	return fmt.Sprintf("未识别到CDN信息，解析到IP为: %v CNAME: %v", ipList, strings.Join(cnames, ","))
}

// 初始化IP记录器
func (a *App) InitIPResolved() {
	core.IPResolved = make(map[string]int)
}

// subodomain
func (a *App) LoadSubDict(configPath string) []string {
	return util.LoadSubdomainDict(util.HomeDir()+configPath, "/dicc.txt")
}

func (a *App) Subdomain(subdomain string, timeout int) []string {
	sr := core.BurstSubdomain(a.ctx, subdomain, timeout, a.qqwryFile, a.cdnFile)
	return []string{sr.Subdomain, strings.Join(sr.Cname, " | "), strings.Join(sr.Ips, " | "), sr.Notes}
}

func (a *App) InitTycHeader(token string) {
	info.InitHEAD(token)
}

// mode 0 = 查询子公司 . mode 1 = 查询公众号
func (a *App) SubsidiariesAndDomains(companyName string, ratio int) []info.CompanyInfo {
	companyId, fuzzName := info.GetCompanyID(a.ctx, companyName) // 获得到一个模糊匹配后，关联度最高的名称
	if companyName != fuzzName {                                 // 如果传进来的名称与模糊匹配的不相同
		var isFuzz = fmt.Sprintf("天眼查模糊匹配名称为%v ——> %v,已替换原有名称进行查.", companyName, fuzzName)
		gologger.Info(a.ctx, isFuzz)
	}
	return info.SearchSubsidiary(a.ctx, fuzzName, companyId, ratio)
}

func (a *App) WechatOfficial(companyName string) []info.WechatReulst {
	companyId, fuzzName := info.GetCompanyID(a.ctx, companyName) // 获得到一个模糊匹配后，关联度最高的名称
	if companyName != fuzzName {
		var isFuzz = fmt.Sprintf("天眼查模糊匹配名称为%v ——> %v,已替换原有名称进行查.", companyName, fuzzName)
		gologger.Info(a.ctx, isFuzz)
	}
	return info.WeChatOfficialAccounts(a.ctx, fuzzName, companyId)
}

// type HunterSearch struct {
// 	Total string
// 	Info  string
// }

// mode = 0 campanyName, mode = 1 domain or ip
// func (a *App) AssetHunter(mode int, target, api string) HunterSearch {
// 	if mode == 0 {
// 		str := fmt.Sprintf("icp.name=\"%v\"", target)
// 		t, i := space.SearchTotal(api, str)
// 		return HunterSearch{
// 			Total: fmt.Sprint(t),
// 			Info:  i,
// 		}
// 	} else {
// 		var str string
// 		// 处理网站域名是IP的情况
// 		if util.RegIP.MatchString(target) {
// 			str = fmt.Sprintf("ip=\"%v\"", target)
// 		} else {
// 			str = fmt.Sprintf("domain.suffix=\"%v\"", target)
// 		}
// 		t, i := space.SearchTotal(api, str)
// 		return HunterSearch{
// 			Total: fmt.Sprint(t),
// 			Info:  i,
// 		}
// 	}
// }

// dirsearch

func (a *App) LoadDirsearchDict(configPath string, newExts []string) []string {
	return util.LoadDirsearchDict(util.HomeDir()+configPath, "/dicc.txt", "%EXT%", newExts)
}

func (a *App) DirScan(options dirsearch.Options) {
	dirsearch.ExitFunc = false
	dirsearch.NewScanner(a.ctx, options)
}

func (a *App) StopDirScan() {
	dirsearch.ExitFunc = true
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

func (a *App) NewTcpScanner(specialTargets []string, ips []string, ports []int, thread, timeout int) {
	portscan.ExitFunc = false
	addrs := portscan.ParseTarget(ips, ports)
	for _, target := range specialTargets {
		temp := strings.Split(target, ":")
		port, err := strconv.Atoi(temp[1]) // 如果后缀端口有误则继续
		if err != nil {
			continue
		}
		addrs = append(addrs, portscan.Address{
			IP:   temp[0],
			Port: port,
		})
	}
	portscan.TcpScan(a.ctx, addrs, thread, timeout)
}

func (a *App) NewSynScanner(specialTargets []string, ips []string, ports []int) {
	portscan.ExitFunc = false
	addrs := portscan.ParseTarget2(ips, util.IntArrayToUint16Array(ports))
	for _, target := range specialTargets {
		temp := strings.Split(target, ":")
		port, err := strconv.ParseUint(temp[1], 10, 16) // 如果后缀端口有误则继续
		if err != nil {
			continue
		}
		addrs = append(addrs, portscan.Address2{
			IP:   net.IP(temp[0]),
			Port: uint16(port),
		})
	}
	portscan.SynScan(a.ctx, addrs)
}

func (a *App) StopPortScan() {
	portscan.ExitFunc = true
}

// 端口暴破
func (a *App) PortBrute(host string, usernames, passwords []string) {
	portscan.PortBrute(a.ctx, host, usernames, passwords)
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
	json.Unmarshal(b, &ts)
	return &ts
}

func (a *App) FofaSearch(query, pageSzie, pageNum, address, email, key string, fraud, cert bool) *space.FofaSearchResult {
	return space.FofaApiSearch(a.ctx, query, pageSzie, pageNum, address, email, key, fraud, cert)
}

func (a *App) Sock5Connect(ip string, port, timeout int, username, password string) bool {
	client, err := clients.SelectProxy(&clients.Proxy{
		Enabled:  true,
		Mode:     "SOCK5",
		Address:  ip,
		Port:     port,
		Username: username,
		Password: password,
	}, clients.DefaultClient())
	if err != nil {
		return false
	}
	_, _, err = clients.NewRequest("GET", "http://www.baidu.com/", nil, nil, timeout, true, client)
	return err == nil
}

func (a *App) IconHash(target string) string {
	_, b, err := clients.NewSimpleGetRequest(target, clients.DefaultClient())
	if err != nil {
		return ""
	}
	return util.Mmh3Hash32(util.Base64Encode(b))
}

// infoscan

func (a *App) CheckTarget(host string, proxy clients.Proxy) *structs.Status {
	protocolURL, err := clients.IsWeb(host, clients.JudgeClient(proxy))
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
func (a *App) InitRule() bool {
	return webscan.InitAll(a.webfingerFile, a.activefingerFile, a.workflowFile)
}

// webscan

func (a *App) FingerLength() int {
	return len(webscan.FingerprintDB)
}

func (a *App) FingerScan(target []string, proxy clients.Proxy) {
	webscan.NewFingerScan(a.ctx, target, proxy)
}

func (a *App) ActiveFingerScan(target []string, proxy clients.Proxy) {
	webscan.NewActiveFingerScan(a.ctx, target, proxy)
}

func (a *App) IsHighRisk(fingerprint string) bool {
	for name, wfe := range webscan.WorkFlowDB {
		if fingerprint == name {
			return len(wfe.PocsName) > 0
		}
	}
	return false
}

// 在漏扫开始时，需要轮询结果
// mode = 0  表示扫指纹漏洞， mode = 1 表示扫全漏洞
func (a *App) NucleiScanner(mode int, target string, fingerprints []string, nucleiPath string, interactsh bool, keywords, severity string) {
	nc := webscan.NewNucleiCaller(nucleiPath, interactsh, severity)
	nc.ReportDirStat()
	if mode == 0 {
		pe := nc.TargetBindFingerPocs(target, fingerprints)
		nc.CallerFP(a.ctx, pe)
	} else {
		keys := []string{}
		if keywords != "" {
			keys = strings.Split(keywords, ",")
		}
		nc.CallerAP(a.ctx, target, keys)
	}
}

func (a *App) NucleiEnabled(nucleiPath string) bool {
	nc := webscan.NewNucleiCaller(nucleiPath, false, "")
	return nc.Enabled(a.ctx)
}

func (a *App) WebPocLength() int {
	return len(webscan.ALLPoc())
}

// hunter

func (a *App) HunterTips(query string) *space.HunterTipsResult {
	return space.SearchHunterTips(query)
}

func (a *App) HunterSearch(api, query, pageSize, pageNum, times, asset string, deduplication bool) *space.HunterResult {
	hr := space.HunterApiSearch(api, query, pageSize, pageNum, times, asset, deduplication)
	time.Sleep(time.Second * 2)
	return hr
}

// quake

func (a *App) QuakeTips(query string) *space.QuakeTipsResult {
	return space.SearchQuakeTips(query)
}

func (a *App) QuakeSearch(ipList []string, query string, pageNum, pageSize int, latest, invalid, honeypot, cdn bool, token string) *space.QuakeResult {
	option := &space.QuakeRequestOptions{
		IpList:   ipList,
		Query:    query,
		PageNum:  pageNum,
		PageSize: pageSize,
		Latest:   latest,
		Invalid:  invalid,
		Honeypot: honeypot,
		CDN:      cdn,
		Token:    token,
	}
	qk := space.QuakeApiSearch(option)
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
			newURL := host + "/" + js
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
		if _, err := url.Parse(target); err != nil {
			return ""
		}
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

func (a *App) UseDruid(target string, attackType int, proxy clients.Proxy) (result []string) {
	druid := alibaba.NewDruid()
	if attackType == 1 {
		result, _ = druid.LoginBrute(target, clients.JudgeClient(proxy))
	} else {
		result, _ = druid.GetSession(target, clients.JudgeClient(proxy))
	}
	return
}
