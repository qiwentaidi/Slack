package custom

import (
	"slack/lib/util"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CheckListBox struct {
	widget.List
	Options  []string // 所有的值
	Selected []string // 已选中的值
}

func NewCheckListBox(options []string) *CheckListBox {
	clb := &CheckListBox{Options: options}
	clb.Length = func() int {
		return len(clb.Options)
	}
	clb.CreateItem = func() fyne.CanvasObject {
		check := widget.NewCheck("", nil)
		check.OnChanged = func(b bool) {
			if b {
				options = append(options, check.Text)
			} else {
				options = util.RemoveElement(options, check.Text)
			}
		}
		return check
	}
	clb.UpdateItem = func(id widget.ListItemID, o fyne.CanvasObject) {
		spl := strings.Split(clb.Options[id], "\\")
		o.(*widget.Check).SetText(spl[len(spl)-1:][0])
	}
	clb.ExtendBaseWidget(clb)
	return clb
}
