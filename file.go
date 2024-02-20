package main

import (
	"bufio"
	"encoding/base64"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slack-wails/lib/update"
	"slack-wails/lib/util"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// File struct 文件操作
type File struct {
	configPath string
}

func NewFile() *File {
	return &File{
		configPath: util.HomeDir() + "/slack/",
	}
}

// 开始就要检测
func (f *File) UserHomeDir() string {
	return util.HomeDir()
}

func (f *File) PathBase(p string) string {
	return filepath.Base(p)
}

func (f *File) OpenFolder(path string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
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

func (f *File) UpdatePocFile() string {
	if err := update.UpdatePoc(f.configPath); err != nil {
		return err.Error()
	}
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

func (f *File) InitConfig() bool {
	return update.InitConfig(f.configPath)
}

func (*File) InitMemo(filepath, content string) bool {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return false
	}
	_, err = f.WriteString(content)
	return err == nil
}

func (*File) ReadMemo(filepath string) map[string]string {
	file, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var key string
	var value strings.Builder
	keyValueMap := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// This is a key line
			if key != "" {
				// Save the previous key-value pair
				keyValueMap[key] = value.String()
				value.Reset()
			}
			key = line[1 : len(line)-1] // Remove brackets
		} else {
			// This is a value line
			value.WriteString(line + "\n")
		}
	}
	// Save the last key-value pair
	if key != "" {
		keyValueMap[key] = value.String()
	}
	return keyValueMap
}

func (*File) Mkdir(dir string) bool {
	return os.Mkdir(dir, 0644) == nil
}

func (*File) WriteFile(filetype, path, content string) bool {
	var buf []byte
	switch filetype {
	case "base64":
		buf, _ = base64.StdEncoding.DecodeString(content)
	// txt
	default:
		buf = []byte(content)
	}
	err := os.WriteFile(path, buf, 0644)
	return err == nil
}
