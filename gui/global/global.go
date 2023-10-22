package global

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 存放输入、输出内容的控件，方便各模块间联动
var (
	WebScanTarget    *widget.Entry
	PortScanTarget   *widget.Entry
	PortBurstResult  *widget.Entry
	PortBurstTarget  *widget.Entry
	SubdomainTarget  *widget.Entry
	ProgerssWebscan  *widget.Label
	PortscanProgress *widget.Label
)

var (
	Win fyne.Window // 全局窗口

	ThinkDict *widget.Entry // 联想模式字典内容
	// 可拖拽文件的文本框
	FscanText         *widget.Entry
	SubdomainText     *widget.Entry
	VulnerabilityText *widget.Entry
	UsernameText      *widget.Entry
	PasswordText      *widget.Entry
	DirDictText       *widget.Entry
)

// 初始化拖拽方法
func InitDropped() {
	EntryDropped([]*widget.Entry{
		FscanText,     // fscan
		SubdomainText, // subdomain
		VulnerabilityText,
		UsernameText, // portburst_username
		PasswordText, // portburst_password
		DirDictText,  // dirsearch_path
	})
}

// 全局只调用一次该方法，不然会存在覆盖
func EntryDropped(ds []*widget.Entry) {
	Win.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		for _, d := range ds { // 判断鼠标拖动的范围是否在文本框内
			abp := fyne.CurrentApp().Driver().AbsolutePositionForObject(d)
			if !d.Disabled() && abp.X <= p.X && abp.X+d.Size().Width >= p.X && abp.Y <= p.Y && abp.Y+d.Size().Height >= p.Y {
				d.SetText(u[0].Path())
			}
		}
	})
}

// 切换标签时，需要重新设置每个标签的主题，不然会导致其他标签页某些控件不改变颜色
// https://github.com/fyne-io/fyne/issues/2996
func Refresh(tableapps ...*container.AppTabs) {
	for _, c := range tableapps {
		c.OnSelected = func(t *container.TabItem) {
			fyne.CurrentApp().Settings().SetTheme(fyne.CurrentApp().Settings().Theme())
		}
	}
}
