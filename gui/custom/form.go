package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 不能对其的FormItem，字数一致会更美观
func NewFormItem(text string, object fyne.CanvasObject) *fyne.Container {
	return container.NewBorder(nil, nil, widget.NewLabel(text), nil, object)
}
