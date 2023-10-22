package custom

import (
	"slack/gui/global"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SelectList struct {
	widget.Button
	Options []string
	Checked []string
	PopUp   *widget.PopUp
}

func NewSelectList(text string, options []string) *SelectList {
	sl := &SelectList{Options: options}
	sl.Text = text
	sl.Icon = theme.MenuDropDownIcon()
	sl.IconPlacement = widget.ButtonIconTrailingText // 使图标跟随在文本之后
	sl.Alignment = widget.ButtonAlignCenter
	cg := widget.NewCheckGroup(sl.Options, func(s []string) {
		if len(s) > 0 {
			sl.Text = ""
			sl.Text += strings.Join(s, " | ")
			sl.Icon = nil
			sl.Refresh()
			sl.Checked = s
		} else {
			sl.Text = text
			sl.Checked = options
			sl.Icon = theme.MenuDropDownIcon()
			sl.Refresh()
		}
	})
	cg.Horizontal = false
	sl.PopUp = widget.NewPopUp(cg, global.Win.Canvas())
	sl.ExtendBaseWidget(sl)
	return sl
}

func (sl *SelectList) Tapped(ev *fyne.PointEvent) {
	h := fyne.CurrentApp().Driver().AbsolutePositionForObject(sl).Y + sl.Size().Height
	sl.PopUp.Resize(fyne.NewSize(sl.Size().Width-2, 0))
	sl.PopUp.ShowAtPosition(fyne.NewPos(sl.Position().X+5, h))
}
