package main

import (
	"fmt"
	"os"
	"slack/common"
	"slack/common/logger"
	"slack/gui/global"
	"slack/gui/memu"
	"slack/gui/module"
	"slack/gui/mytheme"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// 如果.old文件存在则删除
func init() {
	const oldPocZip = "./config/afrog-pocs.zip"
	currentMain := strings.Split(os.Args[0], "\\")
	dir, _ := os.Getwd()
	oldFile := fmt.Sprintf("%v\\.%v.old", dir, currentMain[len(currentMain)-1:][0])
	if _, err := os.Stat(oldFile); err == nil {
		if err2 := os.Remove(oldFile); err2 != nil {
			logger.Error(err)
		}
	}
	if _, err := os.Stat(oldPocZip); err == nil {
		if err2 := os.Remove(oldFile); err2 != nil {
			logger.Error(err)
		}
	}
}

func main() {
	a := app.NewWithID("golang.qwtd.slack")
	a.Settings().SetTheme(&mytheme.MyTheme{})
	global.Win = a.NewWindow("slack " + common.Version)
	memu.ConfrimUpdateClient(0)
	global.Win.SetMainMenu(memu.MyMenu())
	global.Win.SetContent(module.HomePage()) // 设置主题内容
	global.Win.Resize(fyne.NewSize(1200, 900))
	global.Win.SetCloseIntercept(a.Quit) // 退出时关闭全部窗口
	global.Win.SetIcon(mytheme.LogoIcon())
	global.Win.CenterOnScreen()
	global.Win.ShowAndRun()
}
