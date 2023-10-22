package custom

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// 表格内容用superlabel
func NewTableWithUpdateHeader1(data *[][]string, width []float32) *widget.Table {
	table := widget.NewTable(
		func() (rows int, cols int) {
			return len((*data)[1:]), len((*data)[0])
		}, func() fyne.CanvasObject {
			return NewSuperLabel("")
		}, func(id widget.TableCellID, o fyne.CanvasObject) {
			if lb, ok := o.(*SuperLabel); ok {
				lb.SetText((*data)[1:][id.Row][id.Col])
			}
		})
	table.ShowHeaderRow = true
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if lb, ok := o.(*widget.Label); ok {
			lb.SetText((*data)[0][id.Col])
		}
	}
	for i, v := range width {
		table.SetColumnWidth(i, v)
	}
	return table
}

// 表格内容用普通标签
func NewTableWithUpdateHeader2(data *[][]string, width []float32) *widget.Table {
	table := widget.NewTable(
		func() (rows int, cols int) {
			return len((*data)[1:]), len((*data)[0])
		}, func() fyne.CanvasObject {
			return NewCenterLable("")
		}, func(id widget.TableCellID, o fyne.CanvasObject) {
			if lb, ok := o.(*widget.Label); ok {
				lb.SetText((*data)[1:][id.Row][id.Col])
			}
		})
	table.ShowHeaderRow = true
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		if lb, ok := o.(*widget.Label); ok {
			lb.SetText((*data)[0][id.Col])
		}
	}
	for i, v := range width {
		table.SetColumnWidth(i, v)
	}
	return table
}
