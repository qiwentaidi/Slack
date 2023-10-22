package memu

import (
	"slack/gui/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
)

func ChangeTheme() {
	s := settings.NewSettings()
	w := fyne.CurrentApp().NewWindow("主题设置")
	appearance := s.LoadAppearanceScreen(w)
	w.SetContent(appearance)
	w.SetIcon(mytheme.ThemeIcon())
	w.Resize(fyne.NewSize(480, 480))
	w.Show()
}
