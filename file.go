package main

import (
	// "database/sql"

	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slack-wails/lib/update"
	"slack-wails/lib/util"

	"github.com/wailsapp/wails/v2/pkg/logger"
	// _ "github.com/glebarez/go-sqlite"
)

// File struct 文件操作
type File struct {
	DB *sql.DB
}

func NewFile() *File {
	dp := util.HomeDir() + "/slack/config.db"
	db, _ := sql.Open("sqlite", dp)
	return &File{
		DB: db,
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
	if err := update.UpdatePoc(); err != nil {
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
	return update.InitConfig()
}

// func (f *File) CreateTable() bool {
// 	sql := `CREATE TABLE IF NOT EXISTS agent_pool ( hosts TEXT );`
// 	if _, err := f.DB.Exec(sql); err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return false
// 	}
// 	return true
// }
