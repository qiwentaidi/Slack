package main

import (
	"os"
	"os/exec"
	"runtime"
	"slack-wails/lib/update"
	"slack-wails/lib/util"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// File struct 文件操作
type File struct{}

func NewFile() *File {
	return &File{}
}

func (f *File) ExecutionPath() string {
	return util.ExecutionPath()
}

func (f *File) OpenFolder(path string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", f.ExecutionPath()+path)
	case "darwin":
		cmd = exec.Command("open", f.ExecutionPath()+path)
	default:
		cmd = exec.Command("xdg-open", f.ExecutionPath()+path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (f *File) CheckFileStat(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (f *File) GetFileContent(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "文件不存在"
	}
	return string(b)
}

func (f *File) UpdatePocFile(latestVersion string) string {
	if err := update.UpdatePoc(latestVersion); err != nil {
		return err.Error()
	}
	os.RemoveAll(f.ExecutionPath() + "/config/afrog-pocs.zip")
	return ""
}

func (f *File) UpdateClinetFile(latestVersion string) string {
	if err := update.UpdateClinet(latestVersion); err != nil {
		return err.Error()
	}
	return ""
}

func (f *File) Restart() {
	cmd := exec.Command(os.Args[0])
	err := cmd.Start()
	if err != nil {
		logger.NewDefaultLogger().Fatal(err.Error())
	}
	// 退出当前的进程
	os.Exit(0)
}
