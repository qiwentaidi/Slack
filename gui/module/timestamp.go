package module

import (
	"fmt"
	"slack/gui/custom"
	"slack/gui/global"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TimestampUI() *fyne.Container {
	e := widget.NewEntry()
	timestamp2time, time2timestamp := widget.NewEntry(), widget.NewEntry()
	r1, r2 := widget.NewEntry(), widget.NewEntry()
	timeselect1 := widget.NewSelect([]string{"秒(s)", "毫秒(ms)"}, nil)
	timeselect1.SetSelectedIndex(0)
	timeselect2 := widget.NewSelect([]string{"秒(s)", "毫秒(ms)"}, nil)
	timeselect2.SetSelectedIndex(0)
	timeTemplate1 := "2006-01-02 15:04:05" //常规类型
	paused := false
	go func() {
		for {
			if !paused {
				if timeselect1.Selected == "秒(s)" {
					e.SetText(fmt.Sprintf("%v", time.Now().Unix()))
				} else {
					e.SetText(fmt.Sprintf("%v", time.Now().UnixMilli()))
				}
			}
			time.Sleep(time.Second)
		}
	}()
	ctrl := widget.NewButtonWithIcon("停止", theme.MediaStopIcon(), nil)
	ctrl.OnTapped = func() {
		if ctrl.Text == "停止" {
			paused = true
			ctrl.Icon = theme.MediaPlayIcon()
			ctrl.SetText("开始")
		} else {
			paused = false
			ctrl.Icon = theme.MediaStopIcon()
			ctrl.SetText("停止")
		}
	}
	card := widget.NewCard("时间戳转换", "", widget.NewForm(
		widget.NewFormItem("现在: ", container.NewBorder(nil, nil, nil, custom.NewFormItem("控制:", ctrl), e)),
		widget.NewFormItem("时间戳>>北京时间: ", container.NewBorder(nil, nil, nil, widget.NewButton("转换", func() {
			i, err := strconv.ParseInt(timestamp2time.Text, 10, 64)
			if err != nil {
				dialog.ShowError(err, global.Win)
			} else {
				if timeselect1.Selected == "秒(s)" {
					r1.SetText(time.Unix(i, 0).Format(timeTemplate1))
				} else {
					r1.SetText(time.UnixMilli(i).Format(timeTemplate1))
				}
			}
		}), container.NewGridWithColumns(3, timestamp2time, timeselect1, r1))),
		widget.NewFormItem("北京时间>>时间戳: ", container.NewBorder(nil, nil, nil, widget.NewButton("转换", func() {
			stamp, _ := time.ParseInLocation(timeTemplate1, time2timestamp.Text, time.Local)
			if timeselect2.Selected == "秒(s)" {
				r2.SetText(fmt.Sprint(stamp.Unix()))
			} else {
				r2.SetText(fmt.Sprint(stamp.UnixMilli()))
			}
		}), container.NewGridWithColumns(3, time2timestamp, timeselect2, r2))),
	))
	info := widget.NewRichTextFromMarkdown(`# 时间戳
Unix 时间戳是从1970年1月1日（UTC/GMT的午夜）开始所经过的秒数，不考虑闰秒。

---

# 北京时间
## 夏令时
1986年至1991年，中华人民共和国在全国范围实行了六年夏令时，每年从4月中旬的第一个星期日2时整(北京时间)到9月中旬第一个星期日的凌晨2时整(北京夏令时)。除1986年因是实行夏令时的第一年，从5月4日开始到9月14日结束外，其它年份均按规定的时段施行。夏令时实施期间，将时间向后调快一小时。1992年4月5日后不再实行。
	`)
	info.Wrapping = fyne.TextWrapBreak
	return container.NewBorder(card, nil, nil, nil, info)
}
