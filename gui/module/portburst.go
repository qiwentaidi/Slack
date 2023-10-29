package module

import (
	"net/url"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/plugins/portscan"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const portburstInfo = `1、使用内置字典时，账号密码输入框无效，内置字典可以在菜单栏->配置中查看,
要使用自己的字典请取消勾选使用内置字典

2、账号密码框处支持单字段比如admin和txt字典路径，可以拖拽读取文件路径

3、联想模式意为根据关键字段如出生年月或公司名等固定规律组成字典

4、支持协议名: ftp、ssh、telnet、smb、oracle、mssql、mysql、rdp、
postgresql、vnc、redis、memcached、mongodb`

func PortBurstUI() *fyne.Container {
	global.PortBurstTarget = custom.NewMultiLineEntryPlaceHolder("目标请输入成如下格式\n协议://IP:端口\n\nredis://10.0.0.1:6379\nmysql://10.0.0.1:3306")
	global.UsernameText, global.PasswordText = custom.NewFileEntry("若未选中内置字典暴破，则输入要暴破的账号或者直接将txt字典拖入其中，密码框同理"), custom.NewFileEntry("")
	global.UsernameText.Disable()
	global.PasswordText.Disable()
	scan := &widget.Button{Text: "开始", Icon: theme.SearchIcon(), Importance: widget.HighImportance}
	builtin := widget.NewCheck("使用内置字典", func(b bool) {
		if b {
			global.UsernameText.Disable()
			global.PasswordText.Disable()
		} else {
			global.UsernameText.Enable()
			global.PasswordText.Enable()
		}
	})
	builtin.Checked = true
	l := custom.NewCenterLable("")
	associate := widget.NewCheck("联想模式", func(b bool) {
		if b {
			l.SetText("已选中联想模式，请前往 其他工具->联想字典生成器 生成字典，程序会采用文本框的内容进行密码暴破")
		} else {
			l.SetText("")
		}
	})
	info := widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		custom.ShowCustomDialog(theme.InfoIcon(), "提示", "", widget.NewLabel(portburstInfo), nil, fyne.NewSize(400, 300))
	})
	t := custom.NewTableWithUpdateHeader1(&common.PortBurstResult, []float32{80, 270, 200, 200, 0}, custom.SuperClick)
	scan.OnTapped = func() {
		hosts := common.ParseTarget(global.PortBurstTarget.Text, common.Mode_Other)
		go func() {
			for _, host := range hosts {
				GoGo(host, associate.Checked)
			}
		}()
	}
	hs := container.NewHSplit(global.PortBurstTarget, custom.Frame(t))
	hs.Offset = 0.3
	return container.NewBorder(container.NewGridWithRows(3,
		container.NewBorder(nil, nil, container.NewHBox(builtin, associate, info), scan, l),
		container.NewBorder(nil, nil, widget.NewLabel("账号: "), nil, global.UsernameText),
		container.NewBorder(nil, nil, widget.NewLabel("密码: "), nil, global.PasswordText)), nil, nil, nil, hs)
}

func GoGo(host string, associate bool) {
	u, err := url.Parse(host)
	if err != nil {
		custom.Console.Append(err.Error() + "\n")
		return
	}
	switch strings.Split(host, "://")[0] {
	case "ftp":
		portscan.FtpScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "ssh":
		portscan.SshScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "telnet":
		portscan.TelenetScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "smb":
		portscan.SmbScan(u.Host, "", associate, global.UsernameText, global.PasswordText)
	case "oracle":
		portscan.OracleScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "mssql":
		portscan.MssqlScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "mysql":
		portscan.MysqlScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "rdp":
		portscan.RdpScan(u.Host, "", associate, global.UsernameText, global.PasswordText)
	case "postgresql":
		portscan.PostgresScan(u.Host, associate, global.UsernameText, global.PasswordText)
	case "vnc":
		portscan.VncScan(u.Host, associate, global.PasswordText)
	case "redis":
		portscan.RedisScan(u.Host, associate, global.PasswordText)
	case "memcached":
		portscan.MemcachedScan(u.Host)
	case "mongodb":
		portscan.MongodbScan(u.Host)
	}
}
