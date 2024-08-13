package subdomain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slack-wails/core/subdomain/bevigil"
	"slack-wails/core/subdomain/chaos"
	"slack-wails/core/subdomain/securitytrails"
	"slack-wails/core/subdomain/zoomeye"
	"slack-wails/core/waf"
	"slack-wails/lib/gologger"
	"slack-wails/lib/netutil"
	"slack-wails/lib/qqwry"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	DefaultDnsServers = []string{"223.5.5.5:53", "114.114.114.114:53", "8.8.8.8:53"}
	ExitFunc          = false
)

type SubdomainResult struct {
	Domain    string
	Subdomain string
	Ips       []string
	IsCdn     bool
	CdnName   string
	Source    string
}

var (
	Database *qqwry.QQwry
	Cdndata  map[string][]string
)

// 初始化IP纯真库
func InitQqwry(qqwryFile string) {
	fs, err := os.OpenFile(qqwryFile, os.O_RDONLY, 0777)
	if err != nil {
		logger.NewDefaultLogger().Debug("qqwry open err:" + err.Error())
		return
	}
	if d, err := qqwry.NewQQwryFS(fs); err != nil {
		logger.NewDefaultLogger().Debug("qqwry init err:" + err.Error())
		return
	} else {
		Database = d
	}
}

// subdomains 为完整的域名列表
func MultiThreadResolution(ctx context.Context, domain string, subdomains []string, source string, o structs.SubdomainOption) {
	ipResolved := make(map[string]int)
	single := make(chan struct{})
	retChan := make(chan SubdomainResult)
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var id int32
	go func() {
		for sr := range retChan {
			runtime.EventsEmit(ctx, "subdomainLoading", sr)
		}
		close(single)
		gologger.Info(ctx, fmt.Sprintf("已完成 %s 的解析", domain))
		runtime.EventsEmit(ctx, "subdomainComplete", fmt.Sprintf("已完成 %s 的解析", domain))
	}()

	resolutionScan := func(subdomain string) {
		ips, cnames, err := netutil.Resolution(subdomain, o.DnsServers, o.Timeout)
		if err != nil {
			return
		}
		for _, ip := range ips {
			mutex.Lock()
			ipResolved[ip]++
			if ipResolved[ip] > o.ResolveExcludeTimes { // 解析到该IP5次以上就不再显示
				ips = util.RemoveElement(ips, ip)
			}
			mutex.Unlock()
		}
		flag, name := CheckCdn(cnames)
		// 如果目标域名存在waf优先显示
		if len(cnames) > 0 {
			for n, domains := range waf.WAFs {
				for _, domain := range domains {
					if strings.Contains(cnames[0], domain) {
						flag = true
						name = n
						break
					}
				}
			}
		}
		// 解析时移除内网地址
		for _, ip := range ips {
			_, err := Find(ip)
			if err != nil {
				ips = util.RemoveElement(ips, ip)
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
	threadPool, _ := ants.NewPoolWithFunc(o.Thread, func(p interface{}) {
		domain := p.(string)
		atomic.AddInt32(&id, 1)
		runtime.EventsEmit(ctx, "subdomainProgressID", id)
		resolutionScan(domain)
		wg.Done()
	})
	defer threadPool.Release()
	// 枚举模式
	if len(o.Subs) > 0 {
		runtime.EventsEmit(ctx, "subdomainCounts", len(o.Subs))
		for _, sub := range o.Subs {
			if ExitFunc {
				return
			}
			wg.Add(1)
			threadPool.Invoke(sub + "." + domain)
		}
	} else { // API 模式
		runtime.EventsEmit(ctx, "subdomainCounts", len(subdomains))
		for _, subdomain := range subdomains {
			if ExitFunc {
				return
			}
			wg.Add(1)
			threadPool.Invoke(subdomain)
		}
	}
	wg.Wait()
	close(retChan)
	<-single
}

func CheckCdn(cnames []string) (bool, string) {
	for _, cname := range cnames {
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

func ApiPolymerization(ctx context.Context, o structs.SubdomainOption) {
	for _, domain := range o.Domains {
		var subdomains []string
		if o.ChaosApi != "" {
			ch := chaos.FetchHosts(ctx, domain, o.ChaosApi)
			if ch != nil {
				for _, sub := range ch.Subdomains {
					if sub == "*" {
						continue
					}
					subdomains = append(subdomains, sub+"."+domain)
				}
			}
		}
		if o.BevigilApi != "" {
			bh := bevigil.FetchHosts(ctx, domain, o.BevigilApi)
			if bh != nil {
				subdomains = append(subdomains, bh.Subdomains...)
			}
		}
		if o.ZoomeyeApi != "" {
			zh := zoomeye.FetchHosts(ctx, domain, o.ZoomeyeApi)
			if zh != nil {
				for _, item := range zh.List {
					subdomains = append(subdomains, item.Name)
				}
			}
		}
		if o.SecuritytrailsApi != "" {
			sh := securitytrails.FetchHosts(ctx, domain, o.SecuritytrailsApi)
			if sh != nil {
				for _, sub := range sh.Subdomains {
					subdomains = append(subdomains, sub+"."+domain)
				}
			}
		}
		subdomains = util.RemoveDuplicates(subdomains)
		gologger.Info(ctx, fmt.Sprintf("已从API获取到[%s]的子域名: %d个，正在验证存活", domain, len(subdomains)))
		MultiThreadResolution(ctx, domain, subdomains, "API", o)
		time.Sleep(time.Second)
	}
}
