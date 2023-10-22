package custom

import (
	"slack/gui/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewDetailDialog(r1, r2 string) {
	req := NewMultiLineEntryText(r1)
	resp := NewMultiLineEntryText(r2)
	button := &widget.Button{Text: "退出", Icon: theme.LogoutIcon(), Importance: widget.HighImportance}
	dlog := dialog.NewCustomWithoutButtons("漏洞详情", container.NewBorder(nil, button, nil, nil, container.NewHSplit(container.NewBorder(
		NewCenterLable("Request"), nil, nil, nil, req), container.NewBorder(NewCenterLable("Response"), nil, nil, nil, resp))), global.Win)
	dlog.Resize(fyne.NewSize(600, 800))
	dlog.Show()
}
