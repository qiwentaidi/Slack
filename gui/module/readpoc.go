package module

import (
	"fmt"
	"os"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/lib/poc"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	rulePath     = "./config/afrog-pocs/README.md"
	pocDirectory = "./config/afrog-pocs"
)

func MakeReadPocUI() *fyne.Container {
	temp := poc.EmbedFileList
	search := widget.NewEntry()
	button := &widget.Button{Text: "搜索", Icon: theme.SearchIcon(), Importance: widget.HighImportance}
	docs := container.NewDocTabs(
		container.NewTabItemWithIcon("首页", theme.HomeIcon(), HomeRule()),
	)
	pocNums := custom.NewCenterLable("")
	pocNums.SetText(fmt.Sprintf("共查询到POC数量: %v", len(temp)))
	mode := widget.NewCheck("是否自动删除已查看标签", nil)
	mode.SetChecked(true)
	list := widget.NewList(func() int {
		return len(temp)
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(lii widget.ListItemID, o fyne.CanvasObject) {
		slice := strings.Split(temp[lii], "\\")
		name := strings.Split(slice[len(slice)-1:][0], ".")
		o.(*widget.Label).SetText(name[0])
	})
	list.OnSelected = func(id widget.ListItemID) {
		p, err2 := poc.LocalReadPocByPath(temp[id])
		if err2 != nil {
			logger.Info(err2)
		}
		slice := strings.Split(temp[id], "\\")
		name := strings.Split(slice[len(slice)-1:][0], ".")
		if mode.Checked && len(docs.Items) >= 2 {
			for i := 1; i <= len(docs.Items); { // 删除到只剩首页和另一个标签
				if len(docs.Items) <= 1 {
					break
				}
				docs.RemoveIndex(i)
			}
		}
		docs.Append(container.NewTabItem(name[0], DetailPage(temp[id], p)))
		docs.SelectIndex(len(docs.Items) - 1)
	}
	button.OnTapped = func() {

		//go func() {
		temp = nil
		if search.Text != "" {
			for _, v := range poc.EmbedFileList {
				slice := strings.Split(v, "\\")
				name := strings.Split(slice[len(slice)-1:][0], ".")
				if strings.Contains(strings.ToLower(name[0]), strings.ToLower(search.Text)) {
					temp = append(temp, v)
				}

			}
		} else {
			temp = poc.EmbedFileList
		}
		pocNums.SetText(fmt.Sprintf("共查询到POC数量: %v", len(temp)))
		list.Refresh()
		//}()
	}
	hbox := container.NewHSplit(container.NewBorder(mode, pocNums, nil, nil, list), docs)
	hbox.Offset = 0.2
	return container.NewBorder(container.NewBorder(nil, nil, nil, button, search), nil, nil, nil, hbox)
}

func HomeRule() *fyne.Container {
	b, err := os.ReadFile(rulePath)
	if err != nil {
		logger.Info(err)
	}
	homePage := widget.NewRichTextFromMarkdown(string(b))
	homePage.Wrapping = fyne.TextWrapBreak
	return container.NewBorder(nil, nil, nil, nil, container.NewVScroll(homePage))
}

func DetailPage(AbsolutePath string, p poc.Poc) *fyne.Container {
	return container.NewBorder(nil, nil, nil, nil, widget.NewForm(
		widget.NewFormItem("绝对路径:", &widget.Entry{Text: AbsolutePath}),
		widget.NewFormItem("漏洞ID:", &widget.Entry{Text: p.Id}),
		widget.NewFormItem("漏洞名称:", &widget.Entry{Text: p.Info.Name}),
		widget.NewFormItem("作者:", &widget.Entry{Text: p.Info.Author}),
		widget.NewFormItem("危险等级:", &widget.Entry{Text: p.Info.Severity}),
		widget.NewFormItem("漏洞描述:", custom.NewMultiLineEntryText(p.Info.Description)),
		widget.NewFormItem("参考文档:", custom.NewMultiLineEntryText(strings.Join(p.Info.Reference, "\n"))),
		widget.NewFormItem("影响版本:", &widget.Entry{Text: p.Info.Affected}),
		widget.NewFormItem("解决方案:", custom.NewMultiLineEntryText(p.Info.Solutions)),
		widget.NewFormItem("创建时间:", &widget.Entry{Text: p.Info.Created}),
		widget.NewFormItem("标签:", &widget.Entry{Text: p.Info.Created}),
		//widget.NewFormItem("漏洞编号:", &widget.Entry{Text: strings.Join(p.Info.Classification, "\n")}),
	))
}
