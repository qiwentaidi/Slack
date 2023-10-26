package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewDetailDialog(r1, r2 string) {
	req := NewMultiLineEntryText(r1)
	req.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	resp := NewMultiLineEntryText(r2)
	resp.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	hbox := container.NewHSplit(req, resp)
	ShowCustomDialog(nil, "漏洞细节", "", container.NewStack(hbox), nil, fyne.NewSize(800, 800))
}
