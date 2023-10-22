package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func Frame(content fyne.CanvasObject) fyne.CanvasObject {
	top := canvas.NewLine(theme.FocusColor())
	buttom := canvas.NewLine(theme.FocusColor())
	top.Resize(fyne.NewSize(1, 0))
	buttom.Resize(fyne.NewSize(1, 0))
	return container.NewBorder(top, buttom, nil, nil, content)
}
