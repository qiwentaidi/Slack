package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewMultiLineEntryText(text string) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.SetText(text)
	e.Wrapping = fyne.TextWrapBreak
	return e
}

func NewMultiLineEntryPlaceHolder(placeHolder string) *widget.Entry {
	e := widget.NewMultiLineEntry()
	e.PlaceHolder = placeHolder
	e.Wrapping = fyne.TextWrapBreak
	return e
}
