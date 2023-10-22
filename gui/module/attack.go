package module

import (
	"fmt"
	"slack/lib/util"
	"slack/plugins/vulattack"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func AttacktUI() *fyne.Container {
	vul := []string{"Apache Hadoop Yarn RPC RCE"}
	vulList := widget.NewSelectEntry(vul)
	target := widget.NewEntry()
	cmd := widget.NewEntry()
	result := widget.NewMultiLineEntry()
	vulList.OnChanged = func(s string) {
		if vulList.Text == "Apache Hadoop Yarn RPC RCE" {
			cmd.PlaceHolder = fmt.Sprintf("已选择%v漏洞,请填写反弹shell命令", vulList.Text)
			cmd.Refresh()
		} else {
			cmd.PlaceHolder = ""
			cmd.Refresh()
		}
	}
	check := &widget.Button{Text: "开始检测", Importance: widget.HighImportance, OnTapped: func() {
		if util.ArrayContains(vulList.Text, vul) {
			go attack(vulList.Text, target.Text, cmd.Text, result)
		} else {
			result.Text = "请选择正确的漏洞名称进行利用"
			result.Refresh()
		}
	}}
	head := container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("选择漏洞:"), check, vulList),
		container.NewBorder(nil, nil, widget.NewLabel("漏洞地址:"), nil, target),
		container.NewBorder(nil, nil, widget.NewLabel("执行命令:"), nil, cmd),
	)

	return container.NewBorder(head, nil, nil, nil, result)
}

func attack(name, target, cmd string, result *widget.Entry) {
	result.SetText("")
	switch name {
	case "Apache Hadoop Yarn RPC RCE":
		vulattack.HadoopUnauthRCE(target, cmd, result)
	}
}
