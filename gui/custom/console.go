package custom

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const WaitTime = 5

var (
	LogTime int64
	Console *ConsoleLog // 全局日志
)

type ConsoleLog struct {
	widget.Label
}

func init() {
	Console = NewConsoleLog() // 初始化日志记录器
}

func NewConsoleLog() *ConsoleLog {
	cl := &ConsoleLog{}
	cl.Wrapping = fyne.TextWrapBreak
	cl.Text = "资产收集的运行日志以及暴破状况会在此显示"
	cl.ExtendBaseWidget(cl)
	return cl
}

func (cl *ConsoleLog) Append(text string) {
	cl.Text += text
	cl.Refresh()
}

// 轻量日志显示，日志数量过多的话用logger写入文件
func ConsoleWindow() {
	ShowCustomDialog(theme.FileTextIcon(), "日志中心", "清空日志", Frame(container.NewVScroll(Console)), func() {
		Console.SetText("")
	}, fyne.NewSize(500, 700))
}

func LogProgress(currentNum int64, countNum int, errinfo interface{}) {
	if (time.Now().Unix() - LogTime) > WaitTime {
		Console.Append(fmt.Sprintf("已完成 %v/%v %v\n", currentNum, countNum, errinfo))
		LogTime = time.Now().Unix()
	}
}
