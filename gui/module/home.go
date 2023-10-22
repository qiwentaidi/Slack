package module

import (
	"slack/gui/global"
	"slack/gui/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func HomePage() *container.AppTabs {
	var c *container.AppTabs
	sunc1 := container.NewAppTabs(
		container.NewTabItem("网站扫描", WebScanUI()),
		container.NewTabItem("端口/指纹扫描", PortScanUI()),
		container.NewTabItem("暴破与未授权检测", PortBurstUI()),
		container.NewTabItem("漏洞利用", AttacktUI()),
		container.NewTabItem("目录扫描", DirSearchUI()),
		container.NewTabItem("漏洞详情", MakeReadPocUI()),
	)
	sunc2 := container.NewAppTabs(
		container.NewTabItem("公司名称查资产", AssetItem()),
		container.NewTabItem("子域名暴破", SubdomainUI()),
		container.NewTabItem("域名信息查询", IPAndCDN()),
	)
	sunc3 := container.NewAppTabs(
		container.NewTabItem("鹰图", HunterUI()),
		container.NewTabItem("FOFA", FofaUI()),
		container.NewTabItem("360夸克", QuakeUI()),
	)
	sunc4 := container.NewAppTabs(
		container.NewTabItem("编码转换", TranscodingUI()),
		container.NewTabItem("杀软识别&提权补丁", SystemUI()),
		container.NewTabItem("Fscan内容提取", Fscan2txt()),
		container.NewTabItem("反弹shell生成器", ReverseUI()),
		container.NewTabItem("联想字典生成器", DictUI()),
		container.NewTabItem("时间戳转换", TimestampUI()),
		container.NewTabItem("微信", MakeWxUI()),
	)
	card1 := widget.NewCard("渗透测试", "对网站的漏洞进行扫描，并且可以自定义添加YAML，可以直接使用AFROG POC", container.NewGridWrap(fyne.NewSize(190, 50),
		widget.NewButtonWithIcon("网站扫描", mytheme.WebIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(0)
		}),
		widget.NewButtonWithIcon("端口/指纹扫描", mytheme.FingerprintIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(1)
		}),
		widget.NewButtonWithIcon("暴破与未授权检测", mytheme.BombIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(2)
		}),
		widget.NewButtonWithIcon("漏洞利用", mytheme.AttackIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(3)
		}),
		widget.NewButtonWithIcon("目录扫描", mytheme.DirsearchIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(4)
		}),
		widget.NewButtonWithIcon("漏洞详情", mytheme.DetailIcon(), func() {
			c.SelectIndex(1)
			sunc1.SelectIndex(5)
		}),
	))
	card2 := widget.NewCard("资产收集", "可以从公司名称或者域名对资产进行收集", container.NewVBox(
		widget.NewButtonWithIcon("公司名称查资产", mytheme.CompnayIcon(), func() {
			c.SelectIndex(2)
			sunc2.SelectIndex(0)
		}),
		widget.NewButtonWithIcon("子域名暴破", mytheme.SubdomainIcon(), func() {
			c.SelectIndex(2)
			sunc2.SelectIndex(1)
		}),
		widget.NewButtonWithIcon("域名信息查询", mytheme.ChinazIcon(), func() {
			c.SelectIndex(2)
			sunc2.SelectIndex(2)
		}),
	))
	card3 := widget.NewCard("空间引擎", "可以使用鹰图和FOFA等空间引擎对资产进行测绘", container.NewVBox(
		widget.NewButtonWithIcon("鹰图", mytheme.HunterIcon(), func() {
			c.SelectIndex(3)
			sunc3.SelectIndex(0)
		}),
		widget.NewButtonWithIcon("FOFA", mytheme.FofaIcon(), func() {
			c.SelectIndex(3)
			sunc3.SelectIndex(1)
		}),
		widget.NewButtonWithIcon("360夸克", mytheme.QuakeIcon(), func() {
			c.SelectIndex(3)
			sunc3.SelectIndex(2)
		}),
	))
	card4 := widget.NewCard("小工具", "涵盖一些常用的小功能", container.NewGridWrap(fyne.NewSize(190, 50),
		widget.NewButtonWithIcon("编码转换", mytheme.CryptIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(0)
		}),
		widget.NewButtonWithIcon("杀软识别/提权补丁", mytheme.ProblemIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(1)
		}),
		widget.NewButtonWithIcon("Fscan内容提取", mytheme.ExtractionIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(2)
		}),
		widget.NewButtonWithIcon("反弹shell生成器", mytheme.ReserverIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(3)
		}),
		widget.NewButtonWithIcon("联想字典生成器", mytheme.DictIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(4)
		}),
		widget.NewButtonWithIcon("时间戳转换", mytheme.TimeIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(5)
		}),
		widget.NewButtonWithIcon("微信AppID校验", mytheme.WxIcon(), func() {
			c.SelectIndex(4)
			sunc4.SelectIndex(6)
		}),
		widget.NewButtonWithIcon("备忘录", mytheme.NoteIcon(), func() {
			c.SelectIndex(5)
		}),
	))
	c = container.NewAppTabs(
		container.NewTabItemWithIcon("首页", theme.HomeIcon(), container.NewAdaptiveGrid(2, card1, card2, card3, card4)),
		container.NewTabItemWithIcon("渗透测试", mytheme.WebscanIcon(), sunc1),
		container.NewTabItemWithIcon("资产收集", mytheme.AssetIcon(), sunc2),
		container.NewTabItemWithIcon("空间引擎", mytheme.MapIcon(), sunc3),
		container.NewTabItemWithIcon("其他工具", mytheme.ToolsIcon(), sunc4),
		container.NewTabItemWithIcon("备忘录", mytheme.NoteIcon(), MemoItem()),
	)
	global.Refresh(c, sunc1, sunc2, sunc3, sunc4)
	global.InitDropped()
	return c
}
