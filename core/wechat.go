package core

import (
	"encoding/json"
	"io"
	"net/http"
)

type WxAppidMessage struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

var Describe = []WxAppidMessage{
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

func CheckSecert(appid, secert string) string {
	resp, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secert)
	if err != nil {
		return err.Error()
	} else {
		b, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		var wx WxAppidMessage
		json.Unmarshal([]byte(string(b)), &wx)
		for _, d := range Describe {
			if d.Errcode == wx.Errcode {
				return string(b) + "，describe: " + d.Errmsg
			}
		}
		return string(b)
	}
}
