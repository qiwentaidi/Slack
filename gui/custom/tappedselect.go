package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TappedSelect struct {
	widget.Select
	Parent *fyne.Container
}

func NewTappedSelect(options []string, parent *fyne.Container) *TappedSelect {
	s := &TappedSelect{
		Parent: parent,
	}
	s.Options = options
	s.ExtendBaseWidget(s)
	return s
}

func (s *TappedSelect) TappedSecondary(ev *fyne.PointEvent) {
	s.Parent.Remove(s)
	s.Parent.Refresh()
}
