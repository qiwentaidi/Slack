package custom

import (
	"net"
	"net/url"
	"os/exec"
	"runtime"
	"slack/gui/global"
	"slack/lib/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

type SuperLabel struct {
	widget.Label
}

func (sl *SuperLabel) TappedSecondary(ev *fyne.PointEvent) {
	rc := widget.NewPopUpMenu(fyne.NewMenu("",
		&fyne.MenuItem{Label: "复制", Icon: theme.ContentCopyIcon(), Action: func() {
			clipboard.WriteAll(sl.Text)
		}},
		&fyne.MenuItem{Label: "查看详情", Icon: theme.ZoomInIcon(), Action: func() {
			// if canvas.NewText(sl.Text, color.Black).MinSize().Width > sl.Size().Width-15 {
			// 	l := widget.NewLabel(sl.Text)
			// 	sl.p = widget.NewPopUpMenu(fyne.NewMenu("", fyne.NewMenuItem(l.Text, nil)), global.Win.Canvas())
			// 	sl.p.ShowAtPosition(ev.AbsolutePosition)
			// }
		}},
		&fyne.MenuItem{Label: "打开链接", Icon: theme.MailAttachmentIcon(), Action: func() {
			openlink(sl.Text)
		}},
		fyne.NewMenuItemSeparator(),
		&fyne.MenuItem{Label: "将URL发送到Web扫描", Icon: theme.MailSendIcon(), Action: func() {
			sendwebscan(sl.Text)
		}},
		&fyne.MenuItem{Label: "将IP发送到端口扫描", Icon: theme.MailSendIcon(), Action: func() {
			sendportscan(sl.Text)
		}},
		&fyne.MenuItem{Label: "将域名发送到子域名暴破", Icon: theme.MailSendIcon(), Action: func() {
			sendsubdomain(sl.Text)
		}},
	), global.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition) // 面板出现在鼠标点击位置
}

func (sl *SuperLabel) DoubleTapped(*fyne.PointEvent) {
	openlink(sl.Text)
}

func NewSuperLabel(text string) *SuperLabel {
	label := &SuperLabel{}
	label.Text = text
	label.Alignment = fyne.TextAlignCenter
	label.Truncation = fyne.TextTruncateEllipsis
	label.ExtendBaseWidget(label)
	return label
}

func openlink(text string) {
	if _, err := url.ParseRequestURI(text); err == nil {
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", text).Start()
		case "windows":
			err = exec.Command("cmd", "/c", "start", text).Start()
		case "darwin":
			err = exec.Command("open", text).Start()
		}
		if err != nil {
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
