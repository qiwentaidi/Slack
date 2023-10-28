package module

import (
	"errors"
	"fmt"
	"os"
	"slack/common"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/memu"
	"slack/lib/qqwry"
	"slack/lib/result"
	"slack/lib/util"
	"slack/plugins/subdomain"
	"strings"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	IPResolved = make(map[string]int)
	mutex      sync.Mutex
	count      int
	id, num    int64
	database   *qqwry.QQwry
	onec       sync.Once
)

// 初始化IP纯真库
func InitQqwry() {
	fs, err := os.OpenFile(memu.Qqwrypath, os.O_RDONLY, 0400)
	if err != nil {
		logger.Debug("qqwry open err:" + err.Error())
		custom.Console.Append(fmt.Sprintf("[DEB] qqwry open err: %v\n", err))
		return
	}
	if d, err := qqwry.NewQQwryFS(fs); err != nil {
		custom.Console.Append(fmt.Sprintf("[DEB] qqwry init err: %v\n", err))
		return
	} else {
		database = d
	}
}

func SubdomainUI() *fyne.Container {
	global.SubdomainTarget = custom.NewMultiLineEntryPlaceHolder("请输入域名，目标仅支持换行分割\n\n如果域名解析的非常慢，请考虑是否是本机网络不佳")
	global.SubdomainText = custom.NewFileEntry("")
	table := custom.NewTableWithUpdateHeader1(&common.SubdomainResult, []float32{50, 200, 600, 0})
	progress := custom.NewCenterLable("0/0")
	level := custom.NewNumEntry("1")
	thread := custom.NewNumEntry("600")
	scan := &widget.Button{Text: "开始任务", OnTapped: func() {
		progress.SetText("任务正在初始化...")
		go func() {
			onec.Do(InitQqwry)
			domains := common.ParseTarget(global.SubdomainTarget.Text, common.Mode_Other)
			subs := common.ParseFile(global.SubdomainText.Text)
			servers := []string{common.Profile.Subdomain.DNS1 + ":53", common.Profile.Subdomain.DNS2 + ":53"}
			if len(domains) > 0 && len(subs) > 0 && len(servers) > 0 {
				var wg sync.WaitGroup
				limier := make(chan bool, thread.Number)
				data := subdomain.ReadCDNFile()
				count, id = len(subs)*len(domains), 0
				common.SubdomainResult = common.SubdomainResult[:1]
				burstSubdomain(domains, subs, servers, level.Number, thread.Number, limier, data, &wg, progress)
			}
			progress.SetText("任务已结束")
		}()
	}}
	export := widget.NewButtonWithIcon("数据导出", theme.MailForwardIcon(), func() {
		if len(common.SubdomainResult) > 1 {
			go result.ArrayOutput(common.SubdomainResult, "subdomain")
		} else {
			dialog.ShowError(errors.New("请先执行任务"), global.Win)
		}
	})
	hb := container.NewHSplit(container.NewBorder(container.NewBorder(scan, nil, nil, nil,
		container.NewGridWithRows(3,
			container.NewBorder(nil, nil, widget.NewLabel("字典路径:"), nil, global.SubdomainText),
			container.NewBorder(nil, nil, widget.NewLabel("线程数量:"), nil, thread),
			container.NewBorder(nil, nil, widget.NewLabel("暴破层级:"), nil, container.NewBorder(nil, nil, nil, export, level)),
		)), nil, nil, nil, global.SubdomainTarget), custom.Frame(table))
	hb.Offset = 0.25
	return container.NewBorder(nil, progress, nil, nil, hb)
}

// 采用递归判断暴破层级
func burstSubdomain(domains, subs, servers []string, level, tread int, limier chan bool, data map[string][]string, wg *sync.WaitGroup, progress *widget.Label) []string {
	if len(domains) == 0 || level == 0 {
		return nil
	}
	var temps []string
	level -= 1
	for _, domain := range domains {
		for _, sub := range subs {
			wg.Add(1)
			limier <- true
			go func(subdomains string) {
				defer wg.Done()
				addrs, cnames, err := subdomain.Resolution(subdomains, servers)
				atomic.AddInt64(&num, 1)
				if err == nil {
				outloop:
					for _, cdns := range data {
						for _, cdn := range cdns {
							for _, cname := range cnames {
								if strings.Contains(cname, cdn) { // 识别到cdn
									break outloop
								} else if strings.Contains(cname, "cdn") {
									break outloop
								}
							}
						}
					}
					for _, ip := range addrs {
						flag, result, _ := FindWithIP(ip)
						if flag {
							addrs = util.RemoveElement(addrs, ip)
							custom.Console.Append(fmt.Sprintf("[INF] %v识别到CDN字段%v", ip, result))
							continue
						}
						mutex.Lock()
						IPResolved[ip]++
						if IPResolved[ip] > 5 { // 解析到该IP5次以上加入黑名单
							addrs = util.RemoveElement(addrs, ip)
						}
						mutex.Unlock()
					}
					if len(addrs) > 0 {
						atomic.AddInt64(&id, 1)
						temps = append(temps, subdomains)
						common.SubdomainResult = append(common.SubdomainResult, []string{fmt.Sprintf("%v", id), subdomains, strings.Join(addrs, "	"), ""})
					}
				}
				progress.Text = fmt.Sprintf("%d/%d", num, count)
				progress.Refresh()
				<-limier
			}(sub + "." + domain)
		}
	}
	wg.Wait()
	return burstSubdomain(temps, subs, servers, level, tread, limier, data, wg, progress)
}

func FindWithIP(query string) (bool, string, error) {
	result, err := Find(query)
	if strings.Contains(result, "CDN") {
		return true, result, err
	}
	return false, "", err
}

func Find(query string) (string, error) {
	result, err := database.Find(query)
	if err != nil {
		return "", err
	}
	if strings.Contains(result.String(), "对方和您在同一内部网") {
		return "", err
	}
	return result.String(), err
}
