//go:build windows

package bridge

import (
	"syscall"
	"unsafe"
)

func ExecuteScriptCommand(command string) error {
	shell32 := syscall.MustLoadDLL("shell32")
	ShellExecuteW := shell32.MustFindProc("ShellExecuteW")

	args := []uintptr{
		0, // hwnd
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("open"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("cmd.exe"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("/k " + command))),
		0, // directory
		syscall.SW_SHOWDEFAULT,
	}

	_, _, err := ShellExecuteW.Call(args...)
	return err
}
