package module

import (
	"path"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	regGroup = map[string][]string{
		"FTP":       {"[+] ftp"},
		"SSH":       {"[+] SSH"},
		"Mssql":     {"[+] mssql"},
		"Oracle":    {"[+] oracle"},
		"Mysql":     {"[+] mysql"},
		"RDP":       {"[+] RDP"},
		"Redis":     {"[+] Redis"},
		"Postgres":  {"[+] Postgres"},
		"Mongodb":   {"[+] Mongodb"},
		"Memcached": {"[+] Memcached"},
		"MS17-010":  {"[+] MS17-010"},
		"POC":       {"poc"},
		"DC":        {"[+]DC"},
		"INFO":      {"[+] InfoScan"},
		"Vcenter":   {"ID_VC_Welcome"},
		"Camera":    {"len:2512", "len:600", "len:481", "len:480"},
	}
)

func Fscan2txt() *fyne.Container {
	global.FscanText = custom.NewFileEntry("请选择fscan结果文件")
	output := custom.NewMultiLineEntryPlaceHolder(`可提取内容如下:
FTP等协议暴破成功字段
MS17-010
POC字段
DC主机
INFO信息
Vcenter主机
海康摄像头主机	
`)
	do := &widget.Button{Text: "提取关键结果", Importance: widget.HighImportance, OnTapped: func() {
		if path.Ext(global.FscanText.Text) != ".txt" {
			return
		}
		output.SetText("")
		s := common.ParseFile(global.FscanText.Text)
		for name, reg := range regGroup {
			GetLine(name, reg, s, output)
		}
	}}
	return container.NewBorder(container.NewBorder(nil, nil, nil, do, global.FscanText), nil, nil, nil, output)
}

func GetLine(name string, contains, lines []string, output *widget.Entry) {
	var temp []string
	for _, v := range lines {
		for _, c := range contains {
			if strings.Contains(strings.ToLower(v), strings.ToLower(c)) {
				temp = append(temp, v)
			}
		}
	}
	if len(temp) > 0 {
		output.Text += "[" + name + "]\n"
		for _, v := range temp {
			output.Text += v + "\n"
		}
		output.Text += "\n\n"
		output.Refresh()
	}
}
