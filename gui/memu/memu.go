package memu

import (
	"net/url"
	"os"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/mytheme"
	"slack/lib/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const Qqwrypath = "./config/qqwry.dat"

type GlobalShortcut struct {
	dcs *desktop.CustomShortcut
	fun func(shortcut fyne.Shortcut)
}

func MyMenu() *fyne.MainMenu {
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault} // ctrl+,
	logShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierShortcutDefault}          // ctrl+l
	WindowsAddShortcut([]GlobalShortcut{
		{dcs: settingsShortcut, fun: func(shortcut fyne.Shortcut) {
			ChangeTheme()
		}},
		{dcs: logShortcut, fun: func(shortcut fyne.Shortcut) {
			custom.ConsoleWindow()
		}},
	})
	return fyne.NewMainMenu(
		fyne.NewMenu("文件", &fyne.MenuItem{Label: "打开目录", Icon: theme.FolderOpenIcon(), Action: func() {
			dir, _ := os.Getwd()
			util.OpenFolder(dir)
		}}),
		fyne.NewMenu("配置", &fyne.MenuItem{Label: "修改配置", Icon: theme.SettingsIcon(), Action: ConfigCenter}, &fyne.MenuItem{Label: "内置字典", Icon: mytheme.DictIcon(), Action: BurstDict}),
		fyne.NewMenu("更新",
			&fyne.MenuItem{Label: "IP纯真库更新", Icon: theme.DownloadIcon(), Action: DowdloadQqwry},
			&fyne.MenuItem{Label: "漏洞更新", Icon: theme.ViewRefreshIcon(), Action: ConfrimUpdatePoc},
			&fyne.MenuItem{Label: "客户端更新", Icon: mytheme.UpdateIcon(), Action: func() {
				ConfrimUpdateClient(1)
			}},
			&fyne.MenuItem{Label: "功能建议", Icon: mytheme.GithubIcon(), Action: func() {
				u, _ := url.ParseRequestURI(fyne.CurrentApp().Metadata().Custom["Issues"])
				fyne.CurrentApp().OpenURL(u)
			}},
		),
		fyne.NewMenu("主题", &fyne.MenuItem{Label: "修改主题", Icon: mytheme.ThemeIcon(), Shortcut: settingsShortcut, Action: ChangeTheme}),
		fyne.NewMenu("日志", &fyne.MenuItem{Label: "日志中心", Icon: theme.FileTextIcon(), Shortcut: logShortcut, Action: custom.ConsoleWindow}),
	)
}

// 窗口获得焦点时，可以使用的全局快捷键
func WindowsAddShortcut(gsc []GlobalShortcut) {
	for _, shortcut := range gsc {
		global.Win.Canvas().AddShortcut(shortcut.dcs, shortcut.fun)
	}
}
