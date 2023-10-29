package custom

import (
	"net"
	"net/url"
	"slack/gui/global"
	"slack/lib/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

type ClickMode int

const (
	SuperClick ClickMode = iota
	SimpleClick
)

type SuperLabel struct {
	widget.Label
	ClickMode
}

func (sl *SuperLabel) TappedSecondary(ev *fyne.PointEvent) {
	memu := fyne.NewMenu("",
		&fyne.MenuItem{Label: "复制", Icon: theme.ContentCopyIcon(), Action: func() {
			clipboard.WriteAll(sl.Text)
		}},
		&fyne.MenuItem{Label: "打开链接", Icon: theme.MailAttachmentIcon(), Action: func() {
			openlink(sl.Text)
		}},
	)
	if sl.ClickMode == SuperClick {
		memu.Items = append(memu.Items, fyne.NewMenuItemSeparator(),
			&fyne.MenuItem{Label: "将URL发送到Web扫描", Icon: theme.MailSendIcon(), Action: func() {
				sendwebscan(sl.Text)
			}},
			&fyne.MenuItem{Label: "将IP发送到端口扫描", Icon: theme.MailSendIcon(), Action: func() {
				sendportscan(sl.Text)
			}},
			&fyne.MenuItem{Label: "将域名发送到子域名暴破", Icon: theme.MailSendIcon(), Action: func() {
				sendsubdomain(sl.Text)
			}})
	}
	memu.Refresh()
	rc := widget.NewPopUpMenu(memu, global.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition) // 面板出现在鼠标点击位置
}

func (sl *SuperLabel) DoubleTapped(*fyne.PointEvent) {
	openlink(sl.Text)
}

// 给SuperLabel设置格式，比如发送模块只能给有用的地方使用
func NewSuperLabel(text string, mode ClickMode) *SuperLabel {
	label := &SuperLabel{ClickMode: mode}
	label.Text = text
	label.Alignment = fyne.TextAlignCenter
	label.Truncation = fyne.TextTruncateEllipsis
	label.ExtendBaseWidget(label)
	return label
}

func openlink(text string) {
	if u, err := url.ParseRequestURI(text); err == nil {
		if err := fyne.CurrentApp().OpenURL(u); err != nil {
			dialog.ShowInformation("提示", "浏览器打开失败！", global.Win)
		}
	}
}

func sendwebscan(text string) {
	if _, err := url.ParseRequestURI(text); err == nil {
		global.WebScanTarget.Text += text + "\n"
		global.WebScanTarget.Refresh()
	}
}

func sendportscan(text string) {
	if ip := net.ParseIP(text); ip != nil {
		global.PortScanTarget.Text += text + "\n"
		global.PortScanTarget.Refresh()
	}
}

func sendsubdomain(text string) {
	domains := util.RegDomain.FindAllString(text, -1)
	if len(domains) > 0 {
		for _, d := range domains {
			global.SubdomainTarget.Text += d + "\n"
			global.SubdomainTarget.Refresh()
		}
	}
}
