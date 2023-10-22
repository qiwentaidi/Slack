package module

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

func TimestampUI() *fyne.Container {
	e := widget.NewEntry()
	e.ActionItem = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.WriteAll(e.Text)
	})
	go getTimestamp(e)
	t2n := widget.NewEntry()
	n2t := widget.NewEntry()
	r1 := widget.NewEntry()
	r2 := widget.NewEntry()
	timeTemplate1 := "2006-01-02 15:04:05" //常规类型
	card := widget.NewCard("时间戳转换", "", widget.NewForm(
		widget.NewFormItem("现在: ", e),
		widget.NewFormItem("时间戳>>北京时间: ", container.NewBorder(nil, nil, nil, widget.NewButton("转换", func() {
			i, err := strconv.ParseInt(t2n.Text, 10, 64)
			if err != nil {
				r2.SetText(err.Error())
			} else {
				r1.SetText(time.Unix(i, 0).Format(timeTemplate1))
			}
		}), container.NewGridWithColumns(2, t2n, r1))),
		widget.NewFormItem("北京时间>>时间戳: ", container.NewBorder(nil, nil, nil, widget.NewButton("转换", func() {
			stamp, _ := time.ParseInLocation(timeTemplate1, n2t.Text, time.Local)
			r2.SetText(fmt.Sprint(stamp.Unix()))
		}), container.NewGridWithColumns(2, n2t, r2))),
	))
	info := widget.NewRichTextFromMarkdown(`# 时间戳
Unix 时间戳是从1970年1月1日（UTC/GMT的午夜）开始所经过的秒数，不考虑闰秒。

# 北京时间
## 夏令时
1986年至1991年，中华人民共和国在全国范围实行了六年夏令时，每年从4月中旬的第一个星期日2时整(北京时间)到9月中旬第一个星期日的凌晨2时整(北京夏令时)。除1986年因是实行夏令时的第一年，从5月4日开始到9月14日结束外，其它年份均按规定的时段施行。夏令时实施期间，将时间向后调快一小时。1992年4月5日后不再实行。
	`)
	info.Wrapping = fyne.TextWrapBreak
	return container.NewBorder(card, nil, nil, nil, info)
}

func getTimestamp(e *widget.Entry) {
	for {
		e.SetText(fmt.Sprintf("%v", time.Now().Unix()))
		time.Sleep(time.Second * time.Duration(1))
	}
}
