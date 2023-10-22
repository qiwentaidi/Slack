package memu

import (
	"slack/common"
	"slack/gui/custom"
	"slack/gui/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BurstDict() {
	showDict := widget.NewMultiLineEntry()
	protocal := &widget.Select{Options: []string{"ftp", "ssh", "telnet", "smb", "mssql", "oracle", "mysql", "postgresql", "vnc", "redis"}, Selected: "ftp"}
	filed := &widget.Select{Options: []string{"username", "password"}}
	protocal.OnChanged = func(s string) {
		showDict.SetText("")
		if filed.Selected == "username" {
			for _, v := range common.Userdict[s] {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		} else {
			for _, v := range common.Passwords {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		}
	}
	filed.OnChanged = func(s string) {
		showDict.SetText("")
		if s == "username" {
			for _, v := range common.Userdict[protocal.Selected] {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		} else {
			for _, v := range common.Passwords {
				showDict.Text += v + "\n"
				showDict.Refresh()
			}
		}
	}
	filed.SetSelectedIndex(0)
	c := container.NewBorder(container.NewGridWithColumns(2, protocal, filed), nil, nil, nil, showDict)
	custom.ShowCustomDialog(mytheme.DictIcon(), "内置字典查看器", "暂不支持拓展内置字典", c, nil, fyne.NewSize(500, 600))
}
