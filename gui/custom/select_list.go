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
			sl.Icon = nil
			sl.SetText(strings.Join(s, " | "))
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
	h := fyne.CurrentApp().Driver().AbsolutePositionForObject(sl).Y + sl.Size().Height // 获取当前空间下方的Y轴定位
	w := fyne.CurrentApp().Driver().AbsolutePositionForObject(sl).X                    // X轴
	sl.PopUp.Resize(fyne.NewSize(sl.Size().Width, 0))                                  // 大小与控件大小相同
	sl.PopUp.ShowAtPosition(fyne.NewPos(w, h))
}
