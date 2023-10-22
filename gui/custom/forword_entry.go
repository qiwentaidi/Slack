package custom

import (
	"slack/common"
	"slack/gui/global"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

type ForwordEntry struct {
	widget.Entry
}

func NewForwordEntry() *ForwordEntry {
	fe := &ForwordEntry{}
	fe.MultiLine = true
	fe.Wrapping = fyne.TextWrapBreak
	fe.ExtendBaseWidget(fe)
	return fe
}

// overwrite 文本框的右键事件
func (fe *ForwordEntry) TappedSecondary(ev *fyne.PointEvent) {
	m := &fyne.Menu{Items: []*fyne.MenuItem{
		{Label: "复制", Icon: theme.ContentCopyIcon(), Action: func() {
			clipboard.WriteAll(fe.Text)
		}},
		{Label: "剪切", Icon: theme.ContentCutIcon(), Action: func() {
			clipboard.WriteAll(fe.Text)
			fe.SetText("")
		}},
		{Label: "清空", Icon: theme.ContentClearIcon(), Action: func() {
			fe.SetText("")
		}},
		fyne.NewMenuItemSeparator(),
		{Label: "网站扫描", Icon: theme.MailSendIcon(), Action: func() {
			go func() {
				lines := common.ParseTarget(fe.Text, common.Mode_Other)
				if len(lines) > 1 {
					global.WebScanTarget.Text = ""
					global.WebScanTarget.Refresh()
					for _, line := range lines {
						if strings.HasPrefix(line, "http") {
							global.WebScanTarget.Text += line + "\n"
							global.WebScanTarget.Refresh()
						}
					}
				}
			}()
		}},
		{Label: "暴破与未授权检测", Icon: theme.MailSendIcon(), Action: func() {
			go func() {
				lines := common.ParseTarget(fe.Text, common.Mode_Other)
				if len(lines) <= 0 {
					return
				} else {
					global.PortBurstTarget.Text = ""
					global.PortBurstTarget.Refresh()
					for _, line := range lines {
						for protocol := range common.Userdict {
							if strings.HasPrefix(line, protocol) {
								global.PortBurstTarget.Text += line + "\n"
								global.PortBurstTarget.Refresh()
							}
						}
					}
				}
			}()
		}},
	}}
	rc := widget.NewPopUpMenu(m, global.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition)
}
