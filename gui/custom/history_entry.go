package custom

import (
	"bufio"
	"fmt"
	"os"
	"slack/common"
	"slack/common/logger"
	"slack/gui/global"
	"slack/lib/util"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

const defaultSpaceDirectory = "./config/space/"

type HistoryEntry struct {
	widget.Entry
	moudle string
}

func NewHistoryEntry(m string) *HistoryEntry {
	he := &HistoryEntry{moudle: m}
	he.PlaceHolder = "Search..."
	he.ActionItem = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		he.SetText("")
	})
	he.ExtendBaseWidget(he)
	return he
}

// overwrite 文本框的右键事件
func (he *HistoryEntry) TappedSecondary(ev *fyne.PointEvent) {
	childmenu := fyne.NewMenu("")
	m := &fyne.Menu{Items: []*fyne.MenuItem{{Label: "历史记录", Icon: theme.HistoryIcon(), ChildMenu: childmenu},
		fyne.NewMenuItemSeparator(),
		{Label: "复制", Icon: theme.ContentCopyIcon(), Action: func() {
			clipboard.WriteAll(he.Text)
		}},
		{Label: "剪切", Icon: theme.ContentCutIcon(), Action: func() {
			clipboard.WriteAll(he.Text)
			he.SetText("")
		}},
		{Label: "清空", Icon: theme.ContentClearIcon(), Action: func() {
			he.SetText("")
		}},
	}}
	if _, err := os.Stat(defaultSpaceDirectory); err != nil {
		os.Mkdir(defaultSpaceDirectory, 0777)
	}
	if _, err := os.Stat(fmt.Sprintf("%v%v.txt", defaultSpaceDirectory, he.moudle)); err == nil {
		childmenu.Items = append(childmenu.Items, &fyne.MenuItem{Label: "清空查询记录", Icon: theme.DeleteIcon(), Action: func() {
			if err = os.Remove(fmt.Sprintf("%v%v.txt", defaultSpaceDirectory, he.moudle)); err != nil {
				dialog.ShowInformation("", "Failed!", global.Win)
			} else {
				dialog.ShowInformation("", "Success!", global.Win)
			}
		}})
		childmenu.Items = append(childmenu.Items, fyne.NewMenuItemSeparator())
		for _, line := range common.ParseFile(fmt.Sprintf("%v%v.txt", defaultSpaceDirectory, he.moudle)) {
			mi := fyne.NewMenuItem(line, nil)
			mi.Action = func() {
				he.Text = mi.Label
				he.Refresh()
			}
			childmenu.Items = append(childmenu.Items, mi)
		}
	} else {
		childmenu.Items = append(childmenu.Items, &fyne.MenuItem{Label: "未存在查询记录", Action: nil})
	}
	rc := widget.NewPopUpMenu(m, global.Win.Canvas())
	rc.ShowAtPosition(ev.AbsolutePosition)
}

func (he *HistoryEntry) WriteHistory(filename string) {
	file, err := os.OpenFile(defaultSpaceDirectory+filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logger.Info(err)
	}
	b, err := os.ReadFile(defaultSpaceDirectory + filename)
	if err != nil {
		logger.Info(err)
	}
	lines := strings.Split(string(b), "\n")
	write := bufio.NewWriter(file)
	if !util.ArrayContains[string](he.Text, lines) {
		if len(b) > 0 {
			write.WriteString("\n" + he.Text)
		} else {
			write.WriteString(he.Text)
		}
	}
	write.Flush()
	file.Close()
}
