package module

import (
	"fmt"
	"net/url"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/plugins/quake"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type QuakeSearch struct{}

func QuakeUI() *fyne.Container {
	qk := &QuakeSearch{}
	search := custom.NewHistoryEntry("quake")
	latest := &widget.Select{Options: []string{"最新数据", "全量数据"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		if s == "最新数据" {
			quake.Latest = "true"
		} else {
			quake.Latest = "false"
		}
	}}
	latest.SetSelectedIndex(0)
	link, _ := url.Parse("https://quake.360.net/quake/#/help?id=5eb238f110d2e850d5c6aec8&title=%E6%A3%80%E7%B4%A2%E5%85%B3%E9%94%AE%E8%AF%8D")
	configItem := container.NewHBox(latest, layout.NewSpacer(), widget.NewHyperlink("语法提示", link))
	doctabs := container.NewDocTabs(
		container.NewTabItem("首页", qk.HomePage()),
	)
	searchButton := &widget.Button{Text: "查询", Icon: theme.SearchIcon(), Importance: widget.HighImportance, OnTapped: func() {
		if search.Text != "" && strings.Contains(search.Text, ":") {
			search.WriteHistory("quake.txt")
			doctabs.Append(container.NewTabItem(search.Text, qk.NewSeachPage(search.Text)))
			doctabs.SelectIndex(len(doctabs.Items) - 1)
		} else {
			dialog.ShowInformation("提示", "查询内容为空或语法错误", global.Win)
		}
	}}
	return container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, nil, searchButton, search), configItem), nil, nil, nil, doctabs)
}

func (qk *QuakeSearch) HomePage() *fyne.Container {
	e := widget.NewEntry()
	button := &widget.Button{Text: "favicon相似度查询(没做)", Importance: widget.WarningImportance, OnTapped: func() {

	}}
	return container.NewVBox(container.NewBorder(nil, nil, widget.NewLabel("目标favicon地址:"), button, e))
}

func (qk *QuakeSearch) NewSeachPage(query string) *fyne.Container {
	data := [][]string{{"#", "URL", "IP", "端口", "协议", "标题", "域名", "中间件", "城市"}}
	table := custom.NewTableWithUpdateHeader(&data, []float32{50, 280, 200, 60, 100, 200, 200, 150, 100}, custom.SuperClick)
	usage := custom.NewCenterLable("")
	pageNum := custom.NewNumEntry("1")
	pageSize := &widget.Select{Options: []string{"10条/页", "20条/页", "50条/页", "100条/页"}, Alignment: fyne.TextAlignCenter, OnChanged: func(s string) {
		data = data[:1]
		switch s {
		case "10条/页":
			quake.PageSize = 10
		case "20条/页":
			quake.PageSize = 20
		case "50条/页":
			quake.PageSize = 50
		case "100条/页":
			quake.PageSize = 100
		}
		pageNum.SetText("1")
		go qk.Search(query, quake.PageSize, 0, &data, usage)
	}}
	pageSize.SetSelectedIndex(0)
	adjustshowNum := container.NewBorder(nil, nil, nil, container.NewHBox(
		pageSize,
		widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
			go qk.PageTurning(query, pageNum, usage, 0, quake.PageSize, &data)
		}),
		pageNum,
		widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
			go qk.PageTurning(query, pageNum, usage, 1, quake.PageSize, &data)
		}),
	), usage)
	return container.NewBorder(nil, adjustshowNum, nil, nil, custom.Frame(table))
}

// pageup - 0, pagedown - 1
func (qk *QuakeSearch) PageTurning(query string, pageNum *custom.NumberEntry, usage *widget.Label, turnMode, pageSize int, data *[][]string) {
	// 翻页需要+-pageSize才是新的下标
	if quake.FinalTotal > pageSize { // 总数要大于页码才有翻页的必要
		if turnMode == 0 {
			if pageNum.Number <= 1 {
				return
			}
			pageNum.SetText(fmt.Sprintf("%d", pageNum.Number-1))
			*data = (*data)[:1]
			go qk.Search(query, pageSize, quake.StartIndex-pageSize, data, usage)
		} else {
			if pageNum.Number >= quake.FinalTotal/pageSize+1 {
				return
			}
			pageNum.SetText(fmt.Sprintf("%d", pageNum.Number+1))
			*data = (*data)[:1]
			go qk.Search(query, pageSize, quake.StartIndex+pageSize, data, usage)
		}
	}
}

func (qk *QuakeSearch) Search(query string, pageSize, startIndex int, data *[][]string, usage *widget.Label) {
	if query != "" {
		usage.SetText("查询中...")
		quake.QuakeApiSearch(query, startIndex, pageSize, data)
		usage.SetText(quake.QuakeUsage)
	}
}
