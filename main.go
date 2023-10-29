package main

import (
	"slack/gui/global"
	"slack/gui/memu"
	"slack/gui/module"
	"slack/gui/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("golang.qwtd.slack")
	a.Settings().SetTheme(&mytheme.MyTheme{})
	global.Win = a.NewWindow(a.Metadata().Name + " " + a.Metadata().Version)
	memu.ConfrimUpdateClient(0)
	global.Win.SetMainMenu(memu.MyMenu())
	global.Win.SetContent(module.HomePage()) // 设置主题内容
	global.Win.Resize(fyne.NewSize(1200, 900))
	global.Win.SetCloseIntercept(a.Quit) // 退出时关闭全部窗口
	global.Win.SetIcon(mytheme.LogoIcon())
	global.Win.CenterOnScreen()
	global.Win.ShowAndRun()
}
