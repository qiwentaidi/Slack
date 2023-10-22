package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewCenterLable(text string) *widget.Label {
	return &widget.Label{Text: text, Alignment: fyne.TextAlignCenter, TextStyle: fyne.TextStyle{}, Truncation: fyne.TextTruncateEllipsis}
}
