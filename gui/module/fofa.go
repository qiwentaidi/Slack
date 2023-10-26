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
	"slack/plugins/fofa"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FOFASearch struct{}

func FofaUI() *fyne.Container {
	fs := &FOFASearch{}
	search := custom.NewHistoryEntry("fofa")
	fraud := widget.NewCheck("排除干扰(专业版)", func(b bool) {
		if b {
			fofa.Fraud = "true"
		} else {
			fofa.Fraud = "false"
		}
	})
	cert := widget.NewCheck("证书(个人版)", func(b bool) {
		if b {
			fofa.Cert = "true"
		} else {
			fofa.Cert = "false"
		}
	})
	link, _ := url.Parse("https://fofa.info/")
	configItem := container.NewHBox(fraud, cert, layout.NewSpacer(), widget.NewHyperlink("语法提示", link))
	doctabs := container.NewDocTabs(
		container.NewTabItem("首页", fs.HomePage()),
	)
	searchButton := &widget.Button{Text: "查询", Icon: theme.SearchIcon(), Importance: widget.HighImportance, OnTapped: func() {
		if search.Text != "" && len(util.RegCompliance.FindString(search.Text)) > 0 {
			search.WriteHistory("fofa.txt")
			doctabs.Append(container.NewTabItem(search.Text, fs.NewSeachPage(search.Text)))
			doctabs.SelectIndex(len(doctabs.Items) - 1)
		} else {
			dialog.ShowInformation("提示", "查询内容为空或语法错误", global.Win)
		}
	}}
	batchButton := widget.NewButtonWithIcon("批量查询", theme.SearchReplaceIcon(), func() {
		var newsearch string
		e := widget.NewMultiLineEntry()
		e.PlaceHolder = "例如:\n192.168.10.1\n192.168.0.0/24\n192.168.0.0/255.255.255.0\nbaidu.com"
		custom.ShowCustomDialog(mytheme.FofaIcon(), "批量输入: 请输入IP/网段/域名", "查询", e, func() {
			if e.Text != "" {
				for _, ip := range common.ParseTarget(e.Text, common.Mode_Other) {
					if util.RegIP.MatchString(ip) {
						newsearch += fmt.Sprintf(`ip="%v"||`, ip)
					} else {
						newsearch += fmt.Sprintf(`domain="%v"||`, ip)
					}
				}
				doctabs.Append(container.NewTabItem("批量查询", fs.NewSeachPage(newsearch[:len(newsearch)-2])))
				doctabs.SelectIndex(len(doctabs.Items) - 1)
			} else {
				dialog.ShowInformation("提示", "请输入查询内容", global.Win)
			}
		}, fyne.NewSize(400, 400))
	})
	return container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, nil, container.NewHBox(searchButton, batchButton), search), configItem), nil, nil, nil, doctabs)
}

func (fs *FOFASearch) PageTurning(query string, pageNum *custom.NumberEntry, usage *widget.Label, turnMode int, pageSize string, data *[][]string) {
	one, _ := strconv.ParseInt(pageSize, 10, 64)
	if fofa.FofaTotal > one {
		if turnMode == 0 {
			if pageNum.Number <= 1 {
				return
			}
			pageNum.SetText(fmt.Sprint(pageNum.Number - 1))
			*data = (*data)[:1]
			go fs.Search(query, pageSize, pageNum.Text, data, usage)
		} else {
			if int64(pageNum.Number) >= fofa.FofaTotal/one+1 {
				return
			}
			pageNum.SetText(fmt.Sprint(pageNum.Number + 1))
			*data = (*data)[:1]
			go fs.Search(query, pageSize, pageNum.Text, data, usage)
		}
	}
}

func (fs *FOFASearch) Search(query, pageSize, pageNum string, data *[][]string, usage *widget.Label) {
	if query != "" {
		usage.SetText("查询中...")
		fofa.FofaApiSearch(query, pageSize, pageNum, data)
		usage.SetText(fmt.Sprintf("当前查询结果数量:%d条", fofa.FofaTotal))
	}
}

func (fs *FOFASearch) HomePage() *fyne.Container {
	e := widget.NewEntry()
	l := widget.NewEntry()
	return container.NewVBox(container.NewBorder(nil, nil, widget.NewLabel("目标地址:"), &widget.Button{Text: "计算HASH", Importance: widget.WarningImportance, OnTapped: func() {
		c := client.DefaultClient()
		resp, err := c.Get(e.Text)
		if err != nil {
			logger.Info(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		l.SetText(util.Mmh3Hash32(util.Base64Encode(b)))
	}}, e), container.NewBorder(nil, nil, widget.NewLabel("计算结果:"), nil, l))
}

func (fs *FOFASearch) NewSeachPage(query string) *fyne.Container {
	data := [][]string{{"#", "URL", "标题", "IP", "端口", "域名", "协议", "地理位置", "备案号"}}
	table := custom.NewTableWithUpdateHeader1(&data, []float32{50, 250, 200, 150, 60, 200, 70, 150, 200})
	usage := custom.NewCenterLable("")
	pageNum := custom.NewNumEntry("1")

	pageSize := &widget.Select{Options: []string{"10条/页", "100条/页", "500条/页", "1000条/页"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		data = data[:1]
		switch s {
		case "10条/页":
			fofa.PageSize = 10
		case "100条/页":
			fofa.PageSize = 100
		case "500条/页":
			fofa.PageSize = 500
		case "1000条/页":
			fofa.PageSize = 1000
		}
		pageNum.SetText("1")
		go fs.Search(query, fmt.Sprintf("%v", fofa.PageSize), "1", &data, usage)
	}}
	pageSize.SetSelectedIndex(1)
	adjustshowNum := container.NewBorder(nil, nil, widget.NewButtonWithIcon("数据导出", theme.DownloadIcon(), func() {
		fs.Export(query, fmt.Sprint(fofa.PageSize), data)
	}), container.NewHBox(
		pageSize,
		widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
			go fs.PageTurning(query, pageNum, usage, 0, fmt.Sprintf("%v", fofa.PageSize), &data)
		}),
		pageNum,
		widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
			go fs.PageTurning(query, pageNum, usage, 1, fmt.Sprintf("%v", fofa.PageSize), &data)
		}),
	), usage)
	return container.NewBorder(nil, adjustshowNum, nil, nil, custom.Frame(table))
}

func (fs *FOFASearch) Export(query, pageSize string, data [][]string) {
	if len(data) > 1 {
		count := custom.NewCenterLable("")
		filename := widget.NewEntry()
		filename.SetText(fmt.Sprintf("fofa_asset_%v", time.Now().Format("20060102_150405")))
		method := &widget.Select{Options: []string{"导出全部", "导出当前数据"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
			if s == "导出全部" {
				count.Text = fmt.Sprintf("%v", fofa.FofaTotal)
			} else {
				count.Text = pageSize
			}
			count.Refresh()
		}}
		method.SetSelectedIndex(0)
		export := widget.NewButton("导出", nil)
		cancel := widget.NewButton("取消", nil)
		p := widget.NewModalPopUp(container.NewBorder(widget.NewLabel("每次API查询会延时2s，10000条查一次API，导出结束会提示"), container.NewHBox(export, cancel), nil, nil, widget.NewForm(
			widget.NewFormItem("文件名称", filename),
			widget.NewFormItem("导出方式", method),
			widget.NewFormItem("导出条数", count),
		)), global.Win.Canvas())
		p.Resize(fyne.NewSize(300, 200))
		p.Show()
		export.OnTapped = func() {
			if method.Selected == "导出全部" {
				go fofa.AssetExport(query, int(fofa.FofaTotal))
			} else {
				go result.ArrayOutput(data, "fofa_asset")
			}
		}
		cancel.OnTapped = func() {
			p.Hide()
		}
	}
}
