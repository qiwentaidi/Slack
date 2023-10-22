package module

import (
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	ps "slack/plugins/portscan"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	PortScanResult *custom.ForwordEntry
)

func PortScanUI() *fyne.Container {
	global.PortScanTarget = custom.NewMultiLineEntryPlaceHolder(`目标支持换行分割,IP支持如下格式:
192.168.1.1
192.168.1.1/8
192.168.1.1/16
192.168.1.1/24
192.168.1.1,192.168.1.2
192.168.1.1-192.168.255.255
192.168.1.1-255

排除IP可以在可支持输入的IP格式前加!:
!192.168.1.6/28
...

如果端口遗漏多请在配置中调高端口超时时间,默认10s`)
	PortScanResult = custom.NewForwordEntry()
	icmp := widget.NewCheck("ICMP", nil)
	syn := widget.NewCheck("SYN(暂未完成)", nil)
	portlist := widget.NewEntry()
	portgroup := widget.NewRadioGroup([]string{"数据库", "企业端口", "高危端口", "全端口", "自定义"}, func(s string) {
		switch s {
		case "数据库":
			portlist.Text = common.Database
		case "企业端口":
			portlist.Text = common.Enterprise
		case "高危端口":
			portlist.Text = common.HighRisk
		case "全端口":
			portlist.Text = common.All
		default:
			portlist.Text = ""
		}
		portlist.Refresh()
	})
	portgroup.SetSelected("企业端口")
	portgroup.Horizontal = true
	global.PortscanProgress = custom.NewCenterLable("0/0")
	scan := &widget.Button{Text: "开始扫描", Icon: theme.SearchIcon(), Importance: widget.HighImportance, OnTapped: func() {
		go func() {
			var icmp_alive []string
			PortScanResult.Text = ""
			PortScanResult.Refresh()
			ports := common.ParsePort(portlist.Text)
			ips := common.ParseIPs(global.PortScanTarget.Text)
			if icmp.Checked {
				icmp_alive = ps.CheckLive(ips, false, PortScanResult)
				Run(icmp_alive, ports, global.PortscanProgress)
			} else {
				Run(ips, ports, global.PortscanProgress)
			}
		}()
	}}
	sbox := container.NewHSplit(
		container.NewBorder(nil, nil, nil, nil, global.PortScanTarget),
		container.NewBorder(container.NewVBox(container.NewBorder(nil, nil, widget.NewLabel("端口列表:"), scan, portlist), container.NewBorder(nil, nil, nil, nil, container.NewHBox(icmp, syn, widget.NewSeparator(), portgroup))), nil, nil, nil, PortScanResult))
	sbox.Offset = 0.3
	return container.NewBorder(nil, global.PortscanProgress, nil, nil, sbox)
}

func Run(ips []string, ports []int, progress *widget.Label) {
	if len(ports) > 0 && len(ips) > 0 {
		global.PortscanProgress.SetText("端口扫描任务正在初始化")
		ps.PortScanTCP(ips, ports, PortScanResult)
	}
}
