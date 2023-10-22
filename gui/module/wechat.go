package module

import (
	"encoding/json"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type wxAppidCheck struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

var describe = []wxAppidCheck{
	{-1, "系统繁忙，此时请开发者稍候再试"},
	{0, "请求成功"},
	{40001, "AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性"},
	{40002, "请确保grant_type字段值为client_credential"},
	{40013, "AppID错误"},
	{40164, "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置"},
	{89503, "此IP调用需要管理员确认,请联系管理员"},
	{89501, "此IP调用需要管理员确认,请联系管理员"},
	{89506, "24小时内该IP被管理员拒绝调用两次，24小时内不可再使用该IP调用"},
	{89507, "1小时内该IP被管理员拒绝调用一次，1小时内不可再使用该IP调用"},
}

func MakeWxUI() *fyne.Container {
	appid := widget.NewEntry()
	secert := widget.NewEntry()
	result := widget.NewMultiLineEntry()
	form := widget.NewForm(
		widget.NewFormItem("wx_appid", appid),
		widget.NewFormItem("wx_secert", secert),
	)
	button := &widget.Button{Text: "Check", Icon: theme.ConfirmIcon(), Importance: widget.WarningImportance, OnTapped: func() {
		resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid.Text + "&secret=" + secert.Text)
		if err != nil {
			result.SetText(err.Error())
		} else {
			b, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			var wx wxAppidCheck
			json.Unmarshal([]byte(string(b)), &wx)
			for _, d := range describe {
				if d.Errcode == wx.Errcode {
					result.SetText(string(b) + "，describe: " + d.Errmsg)
				}
			}
			if result.Text == "" {
				result.SetText(string(b))
			}
		}
	}}
	return container.NewBorder(container.NewBorder(nil, nil, nil, button, form), nil, nil, nil, result)
}
