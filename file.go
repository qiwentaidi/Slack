package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	rt "runtime"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/update"
	"slack-wails/lib/util"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// File struct 文件操作
type File struct {
	ctx          context.Context
	configPath   string
	downloadPath string
}

func (f *File) startup(ctx context.Context) {
	f.ctx = ctx
}

func NewFile() *File {
	home := util.HomeDir()
	return &File{
		configPath:   home + "/slack/config",
		downloadPath: home + "/Downloads/",
	}
}

func (f *File) FileDialog(ext string) string {
	selection, err := runtime.OpenFileDialog(f.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     ext,
			},
		},
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// selection会返回保存的文件路径+文件名 例如/Users/xxx/Downloads/test.xlsx
func (f *File) SaveFile(filename string) string {
	selection, err := runtime.SaveFileDialog(f.ctx, runtime.SaveDialogOptions{
		Title:           "保存文件",
		DefaultFilename: filename,
	})
	if err != nil {
		return ""
	}
	return selection
}

// 开始就要检测
func (f *File) UserHomeDir() string {
	return util.HomeDir()
}

// 传入路径获取到的信息
type PathInfo struct {
	Name string
	Ext  string
	Dir  string
}

func (f *File) Path(p string) PathInfo {
	// 获取路径中的最后一个元素
	base := filepath.Base(p)
	// 如果有文件扩展名，则去除扩展名（例如 ".exe"）
	ext := filepath.Ext(base)
	if ext != "" {
		base = base[:len(base)-len(ext)]
	}
	return PathInfo{
		Name: base,
		Ext:  strings.ToUpper(strings.TrimLeft(ext, ".")),
		Dir:  filepath.Dir(p),
	}
}

func (f *File) List(path string) []string {
	var files []string
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	return files
}

func (f *File) CheckFileStat(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type FileInfo struct {
	Error   bool
	Message string
	Content string
}

func (f *File) ReadFile(filename string) *FileInfo {
	b, err := os.ReadFile(filename)
	if err != nil {
		return &FileInfo{
			Error:   true,
			Message: err.Error(),
			Content: "",
		}
	}
	if len(b) == 0 {
		return &FileInfo{
			Error:   true,
			Message: "Read file can't be empty",
			Content: "",
		}
	}
	return &FileInfo{
		Error:   false,
		Message: "",
		Content: string(b),
	}
}

func (f *File) UpdatePocFile() string {
	if err := update.UpdatePoc(f.configPath); err != nil {
		return err.Error()
	}
	return ""
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
	return os.Mkdir(dir, 0777) == nil
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

func (a *App) DownloadCyberChef(url string) error {
	cyber := "/slack/CyberChef.zip"
	fileName, err := update.NewDownload(a.ctx, url, a.defaultPath, "downloadProgress", "")
	if err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "downloadComplete", fileName)
	uz := util.NewUnzip()
	if _, err := uz.Extract(util.HomeDir()+cyber, a.defaultPath); err != nil {
		return err
	}
	return os.Remove(util.HomeDir() + cyber)
}

func (f *File) Restart() {
	if rt.GOOS == "darwin" {
		var filename string
		if rt.GOARCH == "arm64" {
			filename = "Slack-macos-arm64.dmg"
		} else {
			filename = "Slack-macos-amd64.dmg"
		}
		cmd := exec.Command("hdiutil", "attach", f.downloadPath+filename)
		if err := cmd.Run(); err == nil {
			cmd = exec.Command("Open", "/Volumes/Slack")
			cmd.Run()
		} else {
			gologger.Debug(f.ctx, err)
		}
	} else {
		cmd := exec.Command(os.Args[0])
		if err := cmd.Start(); err != nil {
			return
		}
		os.Exit(0)
	}
}

func (f *File) DownloadLastestClient() structs.Status {
	const (
		url           = "https://gitee.com/the-temperature-is-too-low/Slack/releases/download/v1/"
		darwin_amd64  = "Slack-macos-amd64.dmg"
		darwin_arm64  = "Slack-macos-arm64.dmg"
		windows_amd64 = "Slack-windows-amd64.exe"
		windows_arm64 = "Slack-windows-arm64.exe"
		linux_amd64   = "Slack-linux-amd64"
		linux_arm64   = "Slack-linux-arm64"
	)
	var filename string
	if rt.GOOS == "darwin" {
		if rt.GOARCH == "amd64" {
			filename = darwin_amd64
		} else {
			filename = darwin_arm64
		}
		_, err := update.NewDownload(f.ctx, url+filename, f.downloadPath, "clientDownloadProgress", "")
		if err != nil {
			return structs.Status{
				Error: true,
				Msg:   err.Error(),
			}
		}
		runtime.EventsEmit(f.ctx, "clientDownloadComplete", "mac-success")
		return structs.Status{
			Error: false,
			Msg:   "Update success!",
		}
	}
	if rt.GOOS == "windows" {
		if rt.GOARCH == "amd64" {
			filename = windows_amd64
		} else {
			filename = windows_arm64
		}
		if err := update.UpdateClientWindows(f.ctx, url+filename); err != nil {
			return structs.Status{
				Error: true,
				Msg:   err.Error(),
			}
		}
		runtime.EventsEmit(f.ctx, "clientDownloadComplete", "win-success")
		return structs.Status{
			Error: false,
			Msg:   "Update success!",
		}
	}
	if rt.GOOS == "linux" {
		if rt.GOARCH == "amd64" {
			filename = linux_amd64
		} else {
			filename = linux_arm64
		}
		dir, _ := os.Getwd()
		_, err := update.NewDownload(f.ctx, url+filename, dir+"/", "clientDownloadProgress", getExecName()+".new")
		if err != nil {
			return structs.Status{
				Error: true,
				Msg:   err.Error(),
			}
		}
		os.Rename(dir+"/"+getExecName()+".new", dir+"/"+getExecName()) // 下载完成就覆盖旧的文件
		os.Chmod(dir+"/"+getExecName(), 0755)                          // 赋予文件执行权限
		runtime.EventsEmit(f.ctx, "clientDownloadComplete", "linux-success")
		return structs.Status{
			Error: false,
			Msg:   "Update success!",
		}
	}
	return structs.Status{
		Error: true,
		Msg:   "Unsupported platform",
	}
}

func (f *File) RemoveOldConfig() error {
	err := os.RemoveAll(f.configPath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return err
}

// windows要移除.xxx.old文件
// mac需要推出挂载
func (f *File) RemoveOldClient() {
	if rt.GOOS == "windows" {
		filename := getExecName()
		if _, err := os.Stat(fmt.Sprintf(".%s.old", filename)); err == nil {
			os.Remove(fmt.Sprintf(".%s.old", filename))
		}
	} else if rt.GOOS == "darwin" {
		cmd := exec.Command("hdiutil", "detach", "/Volumes/Slack")
		if err := cmd.Run(); err != nil {
			gologger.Debug(f.ctx, err)
		}
	}
}

var (
	na         = util.HomeDir() + "/slack/navogation.json"
	navigation []structs.Navigation
)

func (f *File) GetLocalNaConfig() *[]structs.Navigation {
	if !f.CheckFileStat(na) {
		os.Create(na)
		gologger.Error(f.ctx, "Can't create navogation.json")
	}
	b, err := os.ReadFile(na)
	if err != nil {
		gologger.Error(f.ctx, err)
	}
	if err := json.Unmarshal(b, &navigation); err != nil {
		gologger.Error(f.ctx, err)
	}
	return &navigation
}

func (f *File) InsetGroupNavigation(n structs.Navigation) bool {
	navigation = append(navigation, n)
	return f.SaveJsonFile(navigation)
}

func (f *File) InsetItemNavigation(groupName string, child structs.Children) bool {
	for i, n := range navigation {
		if n.Name == groupName {
			navigation[i].Children = append(n.Children, child)
		}
	}
	return f.SaveJsonFile(navigation)
}

func (f *File) SaveNavigation(n []structs.Navigation) bool {
	navigation = n
	return f.SaveJsonFile(navigation)
}

func (f *File) SaveJsonFile(content interface{}) bool {
	b, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		gologger.Error(f.ctx, err)
		return false
	}
	if err := os.WriteFile(na, b, 0777); err != nil {
		gologger.Error(f.ctx, err)
		return false
	}
	return true
}

func (f *File) OpenFolder(filepath string) string {
	filepath, err := getDirectoryPath(filepath)
	if err != nil {
		return err.Error()
	}
	var cmd *exec.Cmd
	switch rt.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filepath)
	case "darwin":
		cmd = exec.Command("open", filepath)
	default:
		cmd = exec.Command("xdg-open", filepath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err.Error()
	}
	return ""
}

// JAR | EXE | LNK | Other
func (f *File) RunApp(jdk, types, filepath string) bool {
	var cmd *exec.Cmd
	// bridge.HideExecWindow(cmd)
	switch types {
	case "JAR":
		cmd = exec.Command(jdk, "-jar", filepath)
	case "APP":
		if rt.GOOS == "windows" {
			if path.Ext(filepath) == ".lnk" {
				cmd = exec.Command("cmd", "/c", "start", filepath)
			} else {
				cmd = exec.Command(filepath)
			}
		}
	default:
		filepath, _ = getDirectoryPath(filepath)
		if rt.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", "start", "cmd", "/K", "cd /d", filepath)
		} else if rt.GOOS == "darwin" {
			// Construct the osascript command to open a new iTerm2 window
			script := `tell application "iTerm"
                        activate
                        tell application "System Events"
                            keystroke "t" using {command down}
                        end tell
                        delay 0.2
                        tell current session of current window
                            write text "cd ` + filepath + `"
                        end tell
                    end tell`
			cmd = exec.Command("osascript", "-e", script)
		}
	}
	go func() {
		if err := cmd.Run(); err != nil {
			return
		}
	}()
	return true
}

func getDirectoryPath(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if fileInfo.IsDir() {
		// 如果是目录，直接返回路径
		return path, nil
	} else {
		// 如果是文件，返回其所在的目录
		return filepath.Dir(path), nil
	}
}

func getExecName() string {
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return ""
	}

	return filepath.Base(execPath)
}

// type Records struct {
// 	Fields []string
// }

// func (*File) HunterRemoveDuplicates(filename string) bool {
// 	// 打开 CSV 文件
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Printf("Cannot access file %s: %v\n", filename, err)
// 		return false
// 	}
// 	defer file.Close()

// 	// 创建 CSV 读取器
// 	reader := csv.NewReader(file)

// 	// 读取 CSV 文件头
// 	headers, err := reader.Read()
// 	if err != nil {
// 		fmt.Printf("Error reading headers: %v\n", err)
// 		return false
// 	}
// 	// 使用 map 去重
// 	urlRecords := make(map[string]Records)
// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			if err.Error() == "EOF" {
// 				break
// 			}
// 			fmt.Printf("Error reading record: %v\n", err)
// 		}
// 		// 获取URL字段进行初步去重
// 		url := record[4]
// 		urls = append(urls, url)
// 		if _, exists := urlRecords[url]; !exists {
// 			urlRecords[url] = Records{
// 				Fields: record,
// 			}
// 		}
// 	}

// 	// 第二轮去重，按 ip-port-title
// 	uniqueRecords := make(map[string]Records)
// 	for _, record := range urlRecords {
// 		fields := record.Fields
// 		ip, port, title := fields[0], fields[1], fields[5]
// 		key := fmt.Sprintf("%s-%s-%s", ip, port, title)
// 		ips = append(ips, ip)
// 		if _, exists := uniqueRecords[key]; !exists {
// 			uniqueRecords[key] = Records{
// 				Fields: fields,
// 			}
// 		}
// 	}

// 	// 创建输出文件
// 	outFile, err := os.Create(getOutputFilename(filename))
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return false
// 	}
// 	defer outFile.Close()

// 	// 创建 CSV 写入器
// 	writer := csv.NewWriter(outFile)
// 	defer writer.Flush()

// 	// 写入 CSV 文件头
// 	if err := writer.Write(headers); err != nil {
// 		fmt.Printf("Error writing headers: %v\n", err)
// 		return false
// 	}

// 	// 写入唯一记录
// 	for _, record := range uniqueRecords {
// 		if err := writer.Write(record.Fields); err != nil {
// 			fmt.Printf("Error writing record: %v\n", err)
// 			return false
// 		}
// 	}
// 	return true
// }

// func getOutputFilename(inputFile string) string {
// 	ext := filepath.Ext(inputFile)
// 	name := strings.TrimSuffix(inputFile, ext)
// 	return fmt.Sprintf("%s_removeDuplicates%s", name, ext)
// }
