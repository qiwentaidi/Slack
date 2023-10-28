package custom

import (
	"image/color"
	"slack/common"
	"slack/gui/mytheme"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	sortOff int = iota
	sortAsc
	sortDesc
)

// 表格内容用superlabel，并且支持排序
func NewTableWithUpdateHeader1(data *[][]string, width []float32) *widget.Table {
	var sorts = make([]int, len((*data)[0]))
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
	table.CreateHeader = func() fyne.CanvasObject { // 一定得先CreateHeader才能使得表格表头为其他类型控件
		return widget.NewButton("000", func() {})
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			b.SetText((*data)[0][id.Col])
			switch sorts[id.Col] {
			case sortAsc:
				b.Icon = theme.MoveUpIcon()
			case sortDesc:
				b.Icon = theme.MoveDownIcon()
			default:
				b.Icon = nil
			}
			b.Importance = widget.MediumImportance
			b.OnTapped = func() {
				applySort(sorts, data, id.Col, table)
			}
			b.Enable()
			b.Refresh()
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

func NewVulnerabilityTable(data *[]common.VulnerabilityInfo, width []float32) *widget.Table {
	var sorts = make([]int, 5)
	table := widget.NewTable(
		func() (rows int, cols int) {
			return len((*data)), 5
		}, func() fyne.CanvasObject {
			return container.NewStack(NewSuperLabel(""), canvas.NewText("", color.White), &widget.Button{Icon: theme.ZoomInIcon(), Importance: widget.LowImportance})
		}, func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*fyne.Container).Objects[0].(*SuperLabel)
			r := o.(*fyne.Container).Objects[1].(*canvas.Text) // canvas.Text可以设置任意字体颜色
			b := o.(*fyne.Container).Objects[2].(*widget.Button)
			l.Show()
			b.Hide()
			r.Hide()
			if id.Col == 4 {
				l.Hide()
				b.Show()
				b.OnTapped = func() {
					req := NewMultiLineEntryText((*data)[id.Row].Request)
					req.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
					resp := NewMultiLineEntryText((*data)[id.Row].Response)
					resp.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
					hbox := container.NewHSplit(req, resp)
					ShowCustomDialog(mytheme.DetailIcon(), "数据包详情", "", container.NewStack(hbox), nil, fyne.NewSize(800, 800))
				}
			} else if id.Col == 0 {
				l.SetText((*data)[id.Row].Id)
			} else if id.Col == 1 {
				l.SetText((*data)[id.Row].Name)
			} else if id.Col == 2 {
				l.Hide()
				r.Show()
				switch (*data)[id.Row].RiskLevel {
				case "CRITICAL":
					r.Color = &color.RGBA{75, 0, 130, 255}
				case "HIGH":
					r.Color = &color.RGBA{200, 0, 0, 200}
				case "MEDIUM":
					r.Color = &color.RGBA{255, 140, 0, 255}
				case "LOW":
					r.Color = &color.RGBA{0, 64, 128, 255}
				case "INFO":
					r.Color = &color.RGBA{0, 200, 0, 255}
				}
				r.Alignment = fyne.TextAlignCenter
				r.Text = (*data)[id.Row].RiskLevel
				r.Refresh()
			} else {
				l.SetText((*data)[id.Row].Url)
			}
		})
	table.ShowHeaderRow = true
	table.CreateHeader = func() fyne.CanvasObject { // 一定得先CreateHeader才能使得表格表头为其他类型控件
		return widget.NewButton("000", func() {})
	}
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		b := o.(*widget.Button)
		if id.Col == -1 {
			b.SetText(strconv.Itoa(id.Row))
			b.Importance = widget.LowImportance
			b.Disable()
		} else {
			b.SetText(common.VulHeader[id.Col])
			switch sorts[id.Col] {
			case sortAsc:
				b.Icon = theme.MoveUpIcon()
			case sortDesc:
				b.Icon = theme.MoveDownIcon()
			default:
				b.Icon = nil
			}
			b.Importance = widget.MediumImportance
			b.OnTapped = func() {
				applySort2(sorts, data, id.Col, table)
			}
			b.Enable()
			b.Refresh()
		}
	}
	for i, v := range width {
		table.SetColumnWidth(i, v)
	}
	return table
}

func applySort(sorts []int, data *[][]string, col int, t *widget.Table) {
	order := sorts[col]
	order++
	if order > sortDesc {
		order = sortOff
	}
	// reset all and assign tapped sort
	for i := 0; i < len((*data)[0]); i++ {
		sorts[i] = sortOff
	}
	sorts[col] = order
	sort.Slice((*data)[1:], func(i, j int) bool {
		a := (*data)[1:][i]
		b := (*data)[1:][j]
		// re-sort with no sort selected
		if order == sortOff {
			return a[col] < b[col]
		}
		if order == sortAsc {
			return strings.Compare(a[col], b[col]) < 0
		}
		return strings.Compare(a[col], b[col]) > 0

	})
	t.Refresh()
}

func applySort2(sorts []int, data *[]common.VulnerabilityInfo, col int, t *widget.Table) {
	order := sorts[col]
	order++
	if order > sortDesc {
		order = sortOff
	}
	// reset all and assign tapped sort
	for i := 0; i < 5; i++ {
		sorts[i] = sortOff
	}
	sorts[col] = order
	sort.Slice(*data, func(i, j int) bool {
		a := (*data)[i]
		b := (*data)[j]
		// re-sort with no sort selected
		if order == sortOff {
			return a.Id < b.Id
		}
		var comparisonResult int
		switch col {
		case 0: // Assuming column 0 corresponds to Id
			comparisonResult = strings.Compare(a.Id, b.Id)
		case 1: // Assuming column 1 corresponds to Name
			comparisonResult = strings.Compare(a.Name, b.Name)
		case 2: // Assuming column 2 corresponds to RiskLevel
			comparisonResult = strings.Compare(a.RiskLevel, b.RiskLevel)
		case 3: // Assuming column 3 corresponds to Url
			comparisonResult = strings.Compare(a.Url, b.Url)
		default: // Add more cases if there are more fields
			comparisonResult = 0
		}
		if order == sortAsc {
			return comparisonResult < 0
		}
		return comparisonResult > 0
	})
	t.Refresh()
}
