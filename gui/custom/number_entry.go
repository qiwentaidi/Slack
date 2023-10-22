package custom

import (
	"strconv"

	"fyne.io/fyne/v2/widget"
)

type NumberEntry struct {
	widget.Entry
	Number int
}

// 限制用户输入的值必须为数字型
func NewNumEntry(defaultNum string) *NumberEntry {
	e := &NumberEntry{}
	e.Text = defaultNum
	e.Number, _ = strconv.Atoi(defaultNum)
	e.OnChanged = func(s string) {
		if _, err := strconv.Atoi(s); err != nil {
			e.SetText(defaultNum)
			e.Number, _ = strconv.Atoi(defaultNum)
		} else {
			e.Number, _ = strconv.Atoi(e.Text)
		}
	}
	e.ExtendBaseWidget(e)
	return e
}
