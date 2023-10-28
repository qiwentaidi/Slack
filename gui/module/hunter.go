package module

import (
	"fmt"
	"io"
	"net/url"
	"slack/common"
	"slack/common/client"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/mytheme"
	"slack/lib/result"
	"slack/lib/util"
	"slack/plugins/hunter"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HunterSearch struct{}

func HunterUI() *fyne.Container {
	hs := &HunterSearch{}
	search := custom.NewHistoryEntry("hunter")
	searchTime := &widget.Select{Options: []string{"最近一个月", "最近半年", "最近一年"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		switch s {
		case "最近一个月":
			hunter.StartTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		case "最近半年":
			hunter.StartTime = time.Now().AddDate(0, 0, -179).Format("2006-01-02")
		case "最近一年":
			hunter.StartTime = time.Now().AddDate(-1, 0, -0).Format("2006-01-02")
		}
	}}
	searchTime.SetSelectedIndex(0)
	assets := &widget.Select{Options: []string{"全部资产", "web服务资产"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		switch s {
		case "全部资产":
			hunter.SelectAssets = "3"
		case "web服务资产":
			hunter.SelectAssets = "1"
		}
	}}
	assets.SetSelectedIndex(1)
	dataDeDuplication := widget.NewCheck("数据去重(需权益积分)", func(b bool) {
		if b {
			hunter.DeDuplication = "true"
		} else {
			hunter.DeDuplication = "false"
		}
	})
	// 查询功能实现
	searchButton := &widget.Button{Text: "查询", Importance: widget.HighImportance, Icon: theme.SearchIcon()}
	batchButton := widget.NewButtonWithIcon("批量查询", theme.SearchReplaceIcon(), nil)
	link, _ := url.Parse("https://hunter.qianxin.com/home/helpCenter?r=8-1")
	configItemHunter := container.NewHBox(searchTime, assets, dataDeDuplication, layout.NewSpacer(), widget.NewHyperlink("语法提示", link))
	doctabs := container.NewDocTabs(
		container.NewTabItem("首页", hs.HomePage()),
	)
	searchButton.OnTapped = func() {
		if search.Text != "" && len(util.RegCompliance.FindString(search.Text)) > 0 {
			search.WriteHistory("hunter.txt")
			doctabs.Append(container.NewTabItem(search.Text, hs.NewSeachPage(search.Text)))
			doctabs.SelectIndex(len(doctabs.Items) - 1)
		} else {
			dialog.ShowInformation("提示", "查询内容为空或语法错误", global.Win)
		}
	}
	batchButton.OnTapped = func() {
		var newsearch string
		e := widget.NewMultiLineEntry()
		e.PlaceHolder = "例如:\n192.168.10.1\n192.168.10.1-255\n192.168.0.192-192.168.0.255\n192.168.0.0/24\n192.168.0.0/255.255.255.0\nbaidu.com"
		custom.ShowCustomDialog(mytheme.HunterIcon(), "批量输入: 请输入IP/网段/域名", "查询", e, func() {
			if e.Text != "" {
				for _, ip := range common.ParseTarget(e.Text, common.Mode_Other) {
					if util.RegIP.MatchString(ip) {
						newsearch += fmt.Sprintf(`ip="%v"||`, ip)
					} else {
						newsearch += fmt.Sprintf(`domain.suffix="%v"||`, ip)
					}
				}
				doctabs.Append(container.NewTabItem("批量查询", hs.NewSeachPage(newsearch[:len(newsearch)-2])))
				doctabs.SelectIndex(len(doctabs.Items) - 1)
			} else {
				dialog.ShowInformation("提示", "请输入查询内容", global.Win)
			}
		}, fyne.NewSize(400, 400))
	}
	return container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, nil, container.NewHBox(searchButton, batchButton), search), configItemHunter), nil, nil, nil, doctabs)
}

// pageup - 0, pagedown - 1
func (hs *HunterSearch) PageTurning(query string, pageNum *custom.NumberEntry, usage *widget.Label, turnMode int, pageSize string, data *[][]string) {
	one, _ := strconv.ParseInt(pageSize, 10, 64)
	if hunter.HunterTotal > one {
		if turnMode == 0 {
			if pageNum.Number <= 1 {
				return
			}
			pageNum.SetText(fmt.Sprintf("%d", pageNum.Number-1))
			*data = (*data)[:1]
			go hs.Search(query, pageSize, pageNum.Text, data, usage)
		} else {
			if int64(pageNum.Number) >= hunter.HunterTotal/one+1 {
				return
			}
			pageNum.SetText(fmt.Sprintf("%d", pageNum.Number+1))
			*data = (*data)[:1]
			go hs.Search(query, pageSize, pageNum.Text, data, usage)
		}
	}
}

// 查询事件		 一页的大小10|50|100，页码
func (hs *HunterSearch) Search(query, pageSize, pageNum string, data *[][]string, usage *widget.Label) {
	if query != "" {
		usage.SetText("查询中...")
		hunter.HunterApiSearch(query, pageSize, pageNum, data)
		usage.SetText(fmt.Sprintf("共查询到%d条资产，用时%dms，当前剩余积分%v", hunter.HunterTotal, hunter.HunterTime, hunter.HunterSurplus))
	}
}

func (hs *HunterSearch) NewSeachPage(query string) *fyne.Container {
	data := [][]string{{"#", "URL", "IP", "端口/服务", "域名", "应用/组件", "站点标题", "状态码", "ICP备案企业", "地理位置", "更新时间"}}
	table := custom.NewTableWithUpdateHeader1(&data, []float32{50, 280, 120, 100, 200, 200, 200, 60, 200, 150, 150})
	pageNum := custom.NewNumEntry("1")
	usage := custom.NewCenterLable("")
	pageSize := &widget.Select{Options: []string{"10条/页", "50条/页", "100条/页"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		data = data[:1]
		switch s {
		case "10条/页":
			hunter.PageSize = 10
		case "50条/页":
			hunter.PageSize = 50
		case "100条/页":
			hunter.PageSize = 100
		}
		pageNum.SetText("1")
		go hs.Search(query, fmt.Sprintf("%v", hunter.PageSize), "1", &data, usage) // 初始化的时候会触发查询
	}}
	pageSize.SetSelectedIndex(0) // 必须再定义完后设置，才能触发查询
	adjustshowNum := container.NewBorder(nil, nil, widget.NewButtonWithIcon("数据导出", theme.DownloadIcon(), func() {
		hs.Export(query, fmt.Sprintf("%v", hunter.PageSize), data)
	}), container.NewHBox(
		pageSize,
		widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
			go hs.PageTurning(query, pageNum, usage, 0, fmt.Sprintf("%v", hunter.PageSize), &data)
		}),
		pageNum,
		widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
			go hs.PageTurning(query, pageNum, usage, 1, fmt.Sprintf("%v", hunter.PageSize), &data)
		}),
	), usage)
	return container.NewBorder(nil, adjustshowNum, nil, nil, custom.Frame(table))
}

func (hs *HunterSearch) HomePage() *fyne.Container {
	e := widget.NewEntry()
	l := widget.NewEntry()
	return container.NewVBox(container.NewBorder(nil, nil, widget.NewLabel("目标地址:"), &widget.Button{Text: "计算MD5", Importance: widget.WarningImportance, OnTapped: func() {
		c := client.DefaultClient()
		resp, err := c.Get(e.Text)
		if err != nil {
			logger.Info(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		l.SetText(util.Md5Value(b))
	}}, e), container.NewBorder(nil, nil, widget.NewLabel("计算结果:"), nil, l))
}

func (hs *HunterSearch) Export(query, pageSize string, data [][]string) {
	if len(data) > 1 {
		count := custom.NewCenterLable("")
		link, _ := url.Parse("https://hunter.qianxin.com/helpCenter?r=7-1")
		filename := widget.NewEntry()
		filename.SetText(fmt.Sprintf("hunter_asset_%v", time.Now().Format("20060102_150405")))
		method := &widget.Select{Options: []string{"导出全部", "导出当前数据"}, PlaceHolder: "导出全部", OnChanged: func(s string) {
			if s == "导出全部" {
				count.Text = fmt.Sprintf("%v", hunter.HunterTotal)
			} else {
				count.Text = pageSize
			}
			count.Refresh()
		}}
		export := widget.NewButton("导出", func() {
			if method.Selected == "导出全部" {
				go hunter.AssetExport(query, int(hunter.HunterTotal))
			} else {
				go result.ArrayOutput(data, "hunter_asset")
			}
		})
		cancel := widget.NewButton("取消", nil)
		p := widget.NewModalPopUp(container.NewBorder(widget.NewLabel("导出"), container.NewBorder(nil, nil, widget.NewHyperlink("导出扣分规则", link), container.NewHBox(export, cancel)), nil, widget.NewSeparator(), widget.NewForm(
			widget.NewFormItem("文件名称", filename),
			widget.NewFormItem("导出方式", method),
			widget.NewFormItem("导出条数", count),
		)), global.Win.Canvas())
		p.Resize(fyne.NewSize(300, 200))
		p.Show()
		cancel.OnTapped = func() {
			p.Hide()
		}
	}
}
