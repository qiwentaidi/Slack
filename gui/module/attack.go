package module

import (
	"slack/gui/custom"
	"slack/gui/global"
	"slack/plugins/vulattack"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	ExecModule = iota
	UploadModule
	reversePlace = "无回显模式，请填写反弹shell"
)

var (
	VulnerabilityMap = map[string][]string{
		"Apache": {"Apache-Hadoop-RCE(无回显)"},
		//"Struts2":       {"S2-061"},
		"Tomcat": {"CVE-2017-12615"},
		// "Elasticsearch": {"CVE-2014-3120", "CVE-2015-1427"},
		// "GIPI":          {"CVE-2022-35914"},
	}
	CurrentVulName = ""
	Target         *widget.Entry
)

func AttacktUI() *fyne.Container {
	keys := make([]string, 0, len(VulnerabilityMap))
	for k := range VulnerabilityMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := container.NewDocTabs()
	Target = widget.NewEntry()
	Target.ActionItem = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		Target.SetText("")
	})
	vulname := widget.NewSelect(nil, func(name string) {
		CurrentVulName = name
		switch name {
		case "Apache-Hadoop-RCE(无回显)":
			result.Items = NewMoudle(ExecModule)
		case "CVE-2017-12615":
			result.Items = NewMoudle(UploadModule)
		// case "S2-061":
		// 	result.Items = NewMoudle(ExecModule)
		default:
			result.Items = NewMoudle()
		}
		result.Refresh()
		result.SelectIndex(0)
	})
	assembly := widget.NewSelect(keys, func(s string) {
		for _, asse := range keys {
			if s == asse {
				vulname.Options = VulnerabilityMap[asse]
				vulname.SetSelectedIndex(0)
				vulname.Refresh()
			}
		}
	})
	assembly.SetSelectedIndex(0)
	vulname.SetSelectedIndex(0)
	return container.NewBorder(container.NewVBox(container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, widget.NewLabel("漏洞组件:"), nil, assembly),
		container.NewBorder(nil, nil, widget.NewLabel("漏洞名称:"), nil, vulname)),
		container.NewBorder(nil, nil, widget.NewLabel("漏洞地址:"), nil, Target),
	), nil, nil, nil, result)
}

func NewMoudle(args ...int) (items []*container.TabItem) {
	if len(args) == 0 {
		return nil
	}
	for _, i := range args {
		switch i {
		case ExecModule:
			items = append(items, container.NewTabItem("命令执行", ExecHub()))
		case UploadModule:
			items = append(items, container.NewTabItem("文件上传", UploadHub()))
		}
	}
	return items
}

func ExecHub() *fyne.Container {
	cmd := &widget.Entry{}
	result := widget.NewMultiLineEntry()
	return container.NewBorder(container.NewBorder(nil, nil, widget.NewLabel("命令:"), &widget.Button{Text: "执行命令", Importance: widget.HighImportance, OnTapped: func() {
		go func() {
			if err := TestTarget(Target.Text); err != nil {
				dialog.ShowError(err, global.Win)
				return
			}
			switch CurrentVulName {
			case "Apache-Hadoop-RCE(无回显)":
				vulattack.ApacheHadoopRCE(Target.Text, cmd.Text, result)
			}
		}()
	}}, cmd), nil, nil, nil, result)
}

func UploadHub() *fyne.Container {
	name := widget.NewEntry()
	result := widget.NewEntry()
	upload := custom.NewMultiLineEntryPlaceHolder("请输入需要上传的内容")
	do := widget.NewButtonWithIcon("上传", theme.UploadIcon(), func() {
		go func() {
			if err := TestTarget(Target.Text); err != nil {
				dialog.ShowError(err, global.Win)
				return
			}
			switch CurrentVulName {
			case "CVE-2017-12615":
				vulattack.CVE_2017_12615(Target.Text, upload.Text, result)
			}
		}()
	})
	return container.NewBorder(widget.NewForm(
		widget.NewFormItem("文件名称:", container.NewBorder(nil, nil, nil, do, name)),
		widget.NewFormItem("上传结果:", result),
	), nil, nil, nil, upload)
}
