package module

import (
	"fmt"
	"net"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/mytheme"
	"slack/lib/result"
	"slack/lib/util"
	"slack/plugins/fofa"
	"slack/plugins/hunter"
	"slack/plugins/info"
	"slack/plugins/webscan"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const webInfo = `1、关键字: 根据yamlpoc的id和info.name判断

2、指定路径: 扫描指定文件夹或单.yaml | .yml POC，文件夹可递归

3、风险等级: 根据security判断

4、仅指纹扫描: 只扫指纹，不会进行敏感目录探测

5、指纹POC扫描: 会进行敏感目录探测，且打指纹POC

6、主动指纹探测: 单选无效需要与仅指纹扫描一起勾选使用，会在指纹扫描基础上增加主动目录探测

7、FastJson: 由于默认是不会进行通用FastJson漏洞扫描如果需要进行FastJson漏洞检测请输入关键字fastjson-1`

// 全部的UI内容
func WebScanUI() *fyne.Container {
	global.WebScanTarget = custom.NewMultiLineEntryPlaceHolder("目标支持如下格式数据仅支持换行分割\ne.g\nhttp://www.baidu.com\nhttp://www.sogou.com\n127.0.0.1:8080")
	scan := widget.NewButtonWithIcon("开始", theme.MediaPlayIcon(), nil)
	stop := widget.NewButtonWithIcon("停止", theme.MediaStopIcon(), nil)
	fofa := widget.NewButton("FOFA", func() {
		NewImport(mytheme.FofaIcon(), "FOFA", 10000)
	})
	hunter := widget.NewButton("鹰图", func() {
		NewImport(mytheme.HunterIcon(), "鹰图", 500)
	})
	global.ProgerssWebscan = custom.NewCenterLable("0/0")
	distribute := widget.NewButtonWithIcon("任务下发", theme.GridIcon(), nil)
	table := custom.NewTableWithUpdateHeader1(&common.ScanResult, []float32{50, 300, 50, 80, 200, 430})
	vultable := custom.NewTableWithUpdateHeader1(&common.VulResult, []float32{50, 300, 80, 680})
	keyword := widget.NewEntry()
	pocnums := custom.NewCenterLable("查看当前条件下可用POC数量")
	info := widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		custom.ShowCustomDialog(theme.InfoIcon(), "提示", "", widget.NewLabel(webInfo), nil, fyne.NewSize(400, 300))
	})
	risk := custom.NewSelectList("点击选择", []string{"critical", "high", "medium", "low", "info"})
	checkpocnums := widget.NewButtonWithIcon("", theme.SearchIcon(), nil)
	global.VulnerabilityText = widget.NewEntry()
	rule := widget.NewForm(
		widget.NewFormItem("搜关键字:", container.NewBorder(nil, nil, nil, info, keyword)),
		widget.NewFormItem("指定路径:", container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
			d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
				if uc != nil {
					global.VulnerabilityText.SetText(uc.URI().Path())
				}
			}, global.Win)
			d.SetFilter(storage.NewExtensionFileFilter([]string{".yml", ".yaml"}))
			d.Show()
		}), global.VulnerabilityText)),
		widget.NewFormItem("风险等级:", risk),
		widget.NewFormItem("POC数量:", container.NewBorder(nil, nil, nil, checkpocnums, pocnums)),
	)
	alive := widget.NewCheck("仅指纹扫描", nil)
	finger2poc := widget.NewCheck("指纹POC扫描", nil)
	active := widget.NewCheck("主动指纹探测", nil)
	checkpocnums.OnTapped = func() {
		go func() {
			pocnums.SetText("正在查询...")
			options := common.NewOptions([]string{}, keyword.Text, strings.Join(risk.Checked, ","), "")
			pocSlice := options.CreatePocList(false, false)
			pocnums.SetText(fmt.Sprintf("%v", len(pocSlice)))
		}()
	}
	scan.OnTapped = func() {
		Scanner(scan, alive, active, finger2poc, keyword, risk)
	}
	distribute.OnTapped = func() {
		DistributeTask()
	}
	c := container.NewAppTabs(
		container.NewTabItem("指纹", table),
		container.NewTabItem("漏洞", vultable),
	)
	global.Refresh(c)
	c.SetTabLocation(1)
	s := container.NewHSplit(container.NewBorder(nil, container.NewGridWithColumns(3, alive, active, finger2poc), nil, nil, rule), container.NewBorder(nil, nil, nil,
		container.NewGridWithRows(4, scan, stop, container.NewGridWithColumns(2, fofa, hunter), distribute), global.WebScanTarget))
	s.Offset = 0.25
	return container.NewBorder(s, global.ProgerssWebscan, nil, nil, custom.Frame(c))
}

func testproxy() error {
	if common.Profile.Proxy.Enable {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%v", common.Profile.Proxy.Address, common.Profile.Proxy.Port), time.Second*time.Duration(10))
		if err != nil {
			dialog.ShowError(fmt.Errorf("proxy is not available\n%s:%v", common.Profile.Proxy.Address, common.Profile.Proxy.Port), global.Win)
			return err
		}
	}
	return nil
}

func NewImport(icon fyne.Resource, text string, num int) {
	var finalsize string
	var size fyne.CanvasObject
	e := widget.NewEntry()
	if text == "FOFA" {
		size = widget.NewEntry()
		size.(*widget.Entry).SetText(fmt.Sprintf("%v", num))
		size.(*widget.Entry).OnChanged = func(s string) {
			max, err := strconv.Atoi(s)
			if err != nil || max > num {
				size.(*widget.Entry).SetText("1000")
				finalsize = "1000"
			} else {
				finalsize = size.(*widget.Entry).Text
			}
		}
	} else {
		size = &widget.Select{Options: []string{"10", "20", "50", "100", "500"}, PlaceHolder: "50", Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
			finalsize = s
		}}
	}
	c := container.NewBorder(nil, nil, nil, nil, container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("请输入导入数量"), nil, size),
		container.NewBorder(nil, nil, widget.NewLabel("请输入查询语句"), nil, e),
	))
	custom.ShowCustomDialog(icon, fmt.Sprintf("联动%v导入目标，最多导入%v条结果", text, num), "导入", c, func() {
		if text == "FOFA" {
			urls, _ := fofa.Import(e.Text, finalsize, false)
			global.WebScanTarget.SetText(strings.Join(urls, "\n"))
		} else {
			urls := hunter.Import(e.Text, finalsize)
			global.WebScanTarget.SetText(strings.Join(urls, "\n"))
		}
	}, fyne.NewSize(400, 0))
}

func Scanner(scan *widget.Button, alive, active, finger2poc *widget.Check, keyword *widget.Entry, risk *custom.SelectList) {
	if err := testproxy(); err != nil {
		return
	}
	if !common.ScanDone {
		if global.WebScanTarget.Text != "" {
			targets := common.ParseTarget(global.WebScanTarget.Text, common.Mode_Url)
			common.VulResult = common.VulResult[:1]
			common.ScanDone = true
			scan.Icon = theme.MediaPlayIcon()
			go func() {
				global.ProgerssWebscan.SetText("正在初始化扫描任务")
				common.UrlFingerMap = make(map[string][]string) // 初始化对象
				if alive.Checked {
					webscan.FingerScan(targets, finger2poc.Checked, global.ProgerssWebscan)
				} else if finger2poc.Checked { // 如果开启了指纹POC扫描，任务需要执行两次，第一次扫主动指纹，第二次扫POC
					webscan.FingerScan(targets, finger2poc.Checked, global.ProgerssWebscan)
					global.ProgerssWebscan.SetText("正在进行主动探测")
					webscan.WebScan(targets, keyword.Text, strings.Join(risk.Checked, ","), true, false, global.VulnerabilityText)
					global.ProgerssWebscan.SetText("正在进行指纹POC扫描")
					webscan.WebScan(targets, keyword.Text, strings.Join(risk.Checked, ","), false, true, global.VulnerabilityText)
				} else {
					if active.Checked {
						webscan.FingerScan(targets, finger2poc.Checked, global.ProgerssWebscan)
					}
					global.ProgerssWebscan.SetText("正在进行漏洞扫描")
					webscan.WebScan(targets, keyword.Text, strings.Join(risk.Checked, ","), alive.Checked, finger2poc.Checked, global.VulnerabilityText)
				}
				fyne.CurrentApp().SendNotification(fyne.NewNotification("提示", "网站扫描任务已完成"))
				global.ProgerssWebscan.SetText("网站扫描任务已完成")
				common.ScanDone = false
				scan.Icon = theme.MediaPlayIcon()
				scan.SetText("开始")
			}()
		} else {
			dialog.ShowInformation("提示", "目标中存在错误的URL或者目标为空", global.Win)
		}
	} else {
		dialog.ShowInformation("提示", "任务正在执行中...", global.Win)
	}
}

func DistributeTask() {
	t := widget.NewEntry()
	custom.ShowCustomDialog(theme.GridIcon(), "请输入公司名称(执行进度可在日志中心查看)", "开始任务", t, func() {
		var (
			final_urls, final_domains, final_ips []string
			webscan_targets, portburst_targets   []string
			domainIP                             = make(map[string][]string)
		)
		dialog.ShowInformation("", "任务已开始执行，执行进度可以从控制台中观察", global.Win)
		common.WaitPortScan = []string{}
		go func() {
			custom.Console.Append("[*]任务下发完毕，公司名称: " + t.Text + "\n")
			fuzzname, _ := info.SearchSubsidiary(t.Text) // 此时已经获取到域名目标在WaitSearchDomains
			if fuzzname != t.Text {
				dialog.ShowConfirm("提示", "天眼查检测名称与输入名称不一致\n模糊查询结果为: "+fuzzname, nil, global.Win)
				return
			}
			custom.Console.Append("[*] 已查询完备案域名...\n")
			// for循环执行获得final_urls, final_domains, final_ips目标继续执行
			for _, domain := range info.WaitSearchDomain {
				total := fofa.SearchTotal(fmt.Sprintf("domain=\"%v\"", domain)) // 查询可以搜索到的域名数量
				if total == 0 {
					continue
				}
				urls, domainIP := fofa.Import(fmt.Sprintf("domain=\"%v\"", domain), fmt.Sprintf("%v", total), true)
				final_urls = append(final_urls, urls...)
				for i, d := range domainIP {
					final_domains = append(final_domains, d...)
					final_ips = append(final_ips, i)
				}
			}
			// 将域名、IP、URL结果进行保存
			result.DistributeTaskReport(fuzzname, "urls", util.RemoveDuplicates[string](final_urls))
			result.DistributeTaskReport(fuzzname, "domains", util.RemoveDuplicates[string](final_domains))
			result.DistributeTaskReport(fuzzname, "ips", util.RemoveDuplicates[string](final_ips))
			custom.Console.Append("[+] 任务目标已保存在./report/distribute/" + fuzzname + "目录下\n")
			// 将IP进行扫描端口，扫描结果在端口扫描的结果框中
			Run(common.WaitPortScan, common.ParsePort("1-65535"), global.PortscanProgress)
			// 筛选结果将http进行web扫描，可暴破协议进行暴破
			custom.Console.Append("[*] 正在进行目标分组...\n")
			lines := common.ParseTarget(PortScanResult.Text, common.Mode_Other)
			if len(lines) <= 0 {
				dialog.ShowError(fmt.Errorf("目标为空"), global.Win)
				return
			} else {
				global.WebScanTarget.SetText("")
				for _, line := range lines {
					if strings.Contains(line, "http") {
						webscan_targets = append(webscan_targets, line)
						for i, ds := range domainIP {
							for _, d := range ds {
								if strings.Contains(line, i) {
									// 把对应域名对应的IP扫到的结果，都将IP替换成域名放入待扫描目标，防止出现域名能访问IP不能访问到服务的情况
									webscan_targets = append(webscan_targets, strings.ReplaceAll(line, i, d))
								}
							}
						}
					}
					for protocol := range common.Userdict {
						if strings.Contains(line, protocol) {
							portburst_targets = append(portburst_targets, line)
						}
					}
				}
			}
			go func() {
				if len(webscan_targets) > 0 {
					webscan.FingerScan(webscan_targets, true, global.ProgerssWebscan)
					global.ProgerssWebscan.SetText("正在进行主动探测")
					webscan.WebScan(webscan_targets, "", "", true, false, global.VulnerabilityText)
					global.ProgerssWebscan.SetText("正在进行指纹POC扫描")
					webscan.WebScan(webscan_targets, "", "", false, true, global.VulnerabilityText)
					fyne.CurrentApp().SendNotification(fyne.NewNotification("提示", "网站扫描任务已完成"))
				}
			}()
			go func() {
				if len(portburst_targets) > 0 {
					for _, host := range portburst_targets {
						GoGo(host, false)
					}
				}
				fyne.CurrentApp().SendNotification(fyne.NewNotification("提示", "端口暴破任务已完成"))
			}()
		}()
	}, fyne.NewSize(400, 0))
}
