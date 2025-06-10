package subdomain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slack-wails/core/subdomain/bevigil"
	"slack-wails/core/subdomain/chaos"
	"slack-wails/core/subdomain/fofa"
	"slack-wails/core/subdomain/github"
	"slack-wails/core/subdomain/hunter"
	"slack-wails/core/subdomain/ip138"
	"slack-wails/core/subdomain/quake"
	"slack-wails/core/subdomain/securitytrails"
	"slack-wails/core/subdomain/zoomeye"
	"slack-wails/core/waf"
	"slack-wails/lib/gologger"
	"slack-wails/lib/qqwry"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils/arrayutil"
	"slack-wails/lib/utils/netutil"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Subdomain struct {
	ctx     context.Context
	options *structs.SubdomainOption
}

type SubdomainResult struct {
	Domain    string
	Subdomain string
	Ips       []string
	IsCdn     bool
	CdnName   string
	Source    string
}

var (
	DefaultDnsServers = []string{"223.6.6.6:53", "8.8.8.8:53"}
	Database          *qqwry.QQwry
	Cdndata           map[string][]string
)

func NewSubdomainEngine(ctx context.Context, option structs.SubdomainOption) *Subdomain {
	subdomain := &Subdomain{
		ctx:     ctx,
		options: &option,
	}
	if len(subdomain.options.DnsServers) == 0 {
		subdomain.options.DnsServers = DefaultDnsServers
	}
	return subdomain
}

// 初始化IP纯真库
func InitQqwry(ctx context.Context, qqwryFile string) {
	fs, err := os.OpenFile(qqwryFile, os.O_RDONLY, 0777)
	if err != nil {
		gologger.DualLog(ctx, gologger.Level_DEBUG, "qqwry open err: "+err.Error())
		return
	}
	if d, err := qqwry.NewQQwryFS(fs); err != nil {
		gologger.DualLog(ctx, gologger.Level_DEBUG, "qqwry init err: "+err.Error())
		return
	} else {
		Database = d
	}
}

// subdomains 为完整的域名列表
func (s *Subdomain) Runner(ctrlCtx context.Context, domain string, subdomains []string, source string) {
	ipResolved := make(map[string]int)
	single := make(chan struct{})
	retChan := make(chan SubdomainResult)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var id int32
	go func() {
		for sr := range retChan {
			runtime.EventsEmit(s.ctx, "subdomainLoading", sr)
		}
		close(single)
		gologger.Info(s.ctx, fmt.Sprintf("已完成 %s 的解析", domain))
		runtime.EventsEmit(s.ctx, "subdomainComplete", fmt.Sprintf("已完成 %s 的解析", domain))
	}()

	resolutionScan := func(subdomain string) {
		ips, cnames, err := netutil.Resolution(subdomain, s.options.DnsServers, s.options.Timeout)
		if err != nil {
			return
		}
		for _, ip := range ips {
			mutex.Lock()
			ipResolved[ip]++
			if ipResolved[ip] > s.options.ResolveExcludeTimes { // 解析到固定IP n次以上就不再显示
				ips = arrayutil.RemoveElement(ips, ip)
			}
			mutex.Unlock()
		}
		flag, name := s.DetectCdnORWAF(cnames)
		// 解析时移除内网地址
		for _, ip := range ips {
			_, err := Find(ip)
			if err != nil {
				ips = arrayutil.RemoveElement(ips, ip)
			}
		}
		retChan <- SubdomainResult{
			Domain:    domain,
			Subdomain: subdomain,
			Ips:       ips,
			IsCdn:     flag,
			CdnName:   name,
			Source:    source,
		}
	}
	threadPool, _ := ants.NewPoolWithFunc(s.options.Thread, func(p interface{}) {
		domain := p.(string)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(s.ctx, "subdomainProgressID", id)
		resolutionScan(domain)
		wg.Done()
	})
	defer threadPool.Release()
	// 只进行CDN/WAF识别
	if s.options.Mode == 3 {
		runtime.EventsEmit(s.ctx, "subdomainCounts", len(s.options.Domains))
		for _, domain := range s.options.Domains {
			if ctrlCtx.Err() != nil {
				return
			}
			wg.Add(1)
			threadPool.Invoke(domain)
		}
		// runtime.EventsEmit(ctx, "subdomainCounts", len(s.options.Subs))
	} else {
		// 枚举模式
		if len(s.options.Subs) > 0 {
			runtime.EventsEmit(s.ctx, "subdomainCounts", len(s.options.Subs))
			for _, sub := range s.options.Subs {
				if ctrlCtx.Err() != nil {
					return
				}
				wg.Add(1)
				threadPool.Invoke(sub + "." + domain)
			}
		} else { // API 模式
			runtime.EventsEmit(s.ctx, "subdomainCounts", len(subdomains))
			for _, subdomain := range subdomains {
				if ctrlCtx.Err() != nil {
					return
				}
				wg.Add(1)
				threadPool.Invoke(subdomain)
			}
		}
	}
	wg.Wait()
	close(retChan)
	<-single
}

func (s *Subdomain) DetectCdnORWAF(cnames []string) (bool, string) {
	wafs := waf.CheckWAF(cnames)
	if wafs.Exsits {
		return true, wafs.Name
	}
	for _, cname := range cnames {
		if strings.Contains(cname, "cdn") {
			return true, "CNames contains cdn filed"
		}
		for name, cdns := range Cdndata {
			for _, cdn := range cdns {
				if strings.Contains(cname, cdn) {
					return true, name
				}
			}
		}
	}
	return false, ""
}

func Find(query string) (string, error) {
	result, err := Database.Find(query)
	if err != nil {
		return "", err
	}
	if strings.Contains(result.String(), "局域网") {
		return "", errors.New("局域网IP")
	}
	return result.String(), err
}

func (s *Subdomain) ApiPolymerization(ctrlCtx context.Context) {
	for _, domain := range s.options.Domains {
		var subdomains []string
		if s.options.ChaosApi != "" {
			ch := chaos.FetchHosts(s.ctx, domain, s.options.ChaosApi)
			if ch != nil {
				for _, sub := range ch.Subdomains {
					if strings.Contains(sub, "*") {
						continue
					}
					subdomains = append(subdomains, sub+"."+domain)
				}
			}
		}
		if arrayutil.ArrayContains("FOFA", s.options.AppendEngines) && s.options.FofaApi != "" {
			fh := fofa.FetchHosts(s.ctx, domain, structs.FofaAuth{
				Address: s.options.FofaAddress,
				Email:   s.options.FofaEmail,
				Key:     s.options.FofaApi,
			})
			subdomains = append(subdomains, fh...)
		}

		if arrayutil.ArrayContains("Hunter", s.options.AppendEngines) && s.options.HunterApi != "" {
			hh := hunter.FetchHosts(s.ctx, domain, s.options.HunterApi)
			subdomains = append(subdomains, hh...)
		}

		if arrayutil.ArrayContains("Quake", s.options.AppendEngines) && s.options.QuakeApi != "" {
			qh := quake.FetchHosts(s.ctx, domain, s.options.QuakeApi)
			subdomains = append(subdomains, qh...)
		}
		if s.options.BevigilApi != "" {
			bh := bevigil.FetchHosts(s.ctx, domain, s.options.BevigilApi)
			if bh != nil {
				subdomains = append(subdomains, bh.Subdomains...)
			}
		}
		if s.options.ZoomeyeApi != "" {
			zh := zoomeye.FetchHosts(s.ctx, domain, s.options.ZoomeyeApi)
			if zh != nil {
				for _, item := range zh.List {
					subdomains = append(subdomains, item.Name)
				}
			}
		}
		if s.options.GithubApi != "" {
			result := github.FetchHosts(s.ctx, domain, s.options.GithubApi)
			// github 会返回一些不正确的域名信息，需要判断根域名是否包含
			for _, sub := range result {
				if strings.Contains(sub, "."+domain) {
					subdomains = append(subdomains, sub)
				}
			}
		}
		if s.options.SecuritytrailsApi != "" {
			sh := securitytrails.FetchHosts(s.ctx, domain, s.options.SecuritytrailsApi)
			if sh != nil {
				for _, sub := range sh.Subdomains {
					subdomains = append(subdomains, sub+"."+domain)
				}
			}
		}
		domains := ip138.FetchHosts(s.ctx, domain)
		if len(domains) > 0 {
			subdomains = append(subdomains, domains...)
		}
		subdomains = arrayutil.RemoveDuplicates(subdomains)
		gologger.Info(s.ctx, fmt.Sprintf("已从API获取到[%s]的子域名: %d个，正在验证存活", domain, len(subdomains)))
		s.Runner(ctrlCtx, domain, subdomains, "API")
		time.Sleep(time.Second)
	}
}
