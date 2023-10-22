package custom

import (
	"slack/gui/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewFileEntry(placeholder string) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeholder
	e.ActionItem = widget.NewButtonWithIcon("", theme.FileTextIcon(), func() {
		if !e.Disabled() {
			d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
				if uc != nil {
					e.SetText(uc.URI().Path())
					e.Refresh()
				}
			}, global.Win)
			d.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
			d.Show()
		}
	})
	return e
}
