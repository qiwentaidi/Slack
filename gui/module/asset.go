package module

import (
	"fmt"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/lib/result"
	"slack/lib/util"
	"slack/plugins/hunter"
	"slack/plugins/info"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func AssetItem() *fyne.Container {
	start := widget.NewButtonWithIcon("开始收集", theme.SearchIcon(), nil)
	sr1 := widget.NewCheck("反查域名&查询全资子公司", nil)
	sr2 := widget.NewCheck("查询HUNTER资产数量", nil)
	sr3 := widget.NewCheck("小程序&公众号(暂未开发)", nil)
	target := widget.NewMultiLineEntry()
	target.PlaceHolder = "ICP名称会进行模糊匹配\n目标仅支持换行分割"
	c := container.NewAppTabs(
		container.NewTabItem("控股企业", custom.NewTableWithUpdateHeader1(&common.HoldAsset, []float32{200, 100, 120, 380})),
		container.NewTabItem("HUNTER资产数量(资产/1积分)", custom.NewTableWithUpdateHeader1(&common.HunterAsset, []float32{400, 200})),
		container.NewTabItem("小程序&公众号", custom.NewTableWithUpdateHeader1(&common.WechatAsset, []float32{150, 150, 150, 100})),
	)
	c.SetTabLocation(2)
	start.OnTapped = func() {
		go func() {
			if !sr1.Checked && !sr2.Checked && !sr3.Checked {
				dialog.ShowInformation("", "请勾选需要查询的字段", global.Win)
				return
			}
			targets := common.ParseTarget(target.Text, common.Mode_Other)
			if len(targets) <= 0 {
				dialog.ShowInformation("", "需要查询的目标为空", global.Win)
				return
			}
			if sr1.Checked {
				common.HoldAsset = common.HoldAsset[:1]
				custom.Console.Append("[+] 正在查询全资子公司以及反查ICP域名备案信息\n")
				for _, t := range targets {
					fuzzname, subsidiaries := info.SearchSubsidiary(t)
					if t != fuzzname {
						targets = util.ReplaceElement(targets, t, fuzzname) // 将错误的名称替换成正确的ICP(含误报概率)
					}
					targets = append(targets, subsidiaries...)
				}
			}
			if sr2.Checked {
				custom.Console.Append(fmt.Sprintf("[+] 正在进行hunter资产数量查询,总数:%d\n", len(targets)+len(info.WaitSearchDomain)))
				common.HunterAsset = common.HunterAsset[:1]
				for _, company := range targets {
					hunter.SeachICP(company)
				}
				for _, domain := range info.WaitSearchDomain {
					hunter.SeachDomain(domain)
				}

				custom.Console.Append(hunter.Restquota + "\n")
			}
			custom.Console.Append("---资产收集，运行结束---\n")
		}()
	}
	download := widget.NewButtonWithIcon("保存结果", theme.DownloadIcon(), func() {
		go func() {
			combinedResult := [][]string{}
			combinedResult = append(append(append(combinedResult, common.HoldAsset...), common.WechatAsset...), common.HunterAsset...)
			result.ArrayOutput(combinedResult, "asset")
		}()
	})
	sbox := container.NewHSplit(container.NewBorder(container.NewBorder(nil, nil, nil, download, start), nil, nil, nil, target),
		container.NewBorder(container.NewHBox(widget.NewLabel("查询条件:"), sr1, sr2, sr3), nil, nil, nil, custom.Frame(c)))
	sbox.Offset = 0.3
	return container.NewBorder(nil, nil, nil, nil, sbox)
}
