package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	rt "runtime"
	"slack-wails/core/webscan"
	"slack-wails/lib/clients"
	"slack-wails/lib/fileutil"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/update"
	"slack-wails/lib/util"
	"strings"
	"time"

	"github.com/mat/besticon/v3/ico"
	"github.com/nfnt/resize"
	"github.com/orcastor/fico"
	lnk "github.com/parsiya/golnk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var Userdict = map[string][]string{
	"ftp":        {"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"},
	"mysql":      {"root", "mysql"},
	"mssql":      {"sa", "sql"},
	"smb":        {"administrator", "admin", "guest"},
	"rdp":        {"administrator", "admin", "guest"},
	"postgresql": {"postgres", "admin"},
	"ssh":        {"root", "admin"},
	"mongodb":    {"root", "admin"},
	"oracle":     {"sys", "system", "admin", "test", "web", "orcl"},
	"ldap":       {"admin", "administrator", "root"},
	"socks5":     {"admin", "administrator"},
	"mqtt":       {"admin", "administrator"},
	"vnc":        {"admin", "administrator", "root"},
	"telnet":     {"root", "admin"},
}

var Passwords = []string{"123456", "admin", "admin123", "root", "", "pass123", "pass@123", "password", "123123", "654321", "111111", "123", "1", "admin@123", "Admin@123", "admin123!@#", "{user}", "{user}1", "{user}111", "{user}123", "{user}@123", "{user}_123", "{user}#123", "{user}@111", "{user}@2019", "{user}@123#4", "P@ssw0rd!", "P@ssw0rd", "Passw0rd", "qwe123", "12345678", "test", "test123", "123qwe", "123qwe!@#", "123456789", "123321", "666666", "a123456.", "123456~a", "123456!a", "000000", "1234567890", "8888888", "!QAZ2wsx", "1qaz2wsx", "abc123", "abc123456", "1qaz@WSX", "a11111", "a12345", "Aa1234", "Aa1234.", "Aa12345", "a123456", "a123123", "Aa123123", "Aa123456", "Aa12345.", "sysadmin", "system", "1qaz!QAZ", "2wsx@WSX", "qwe123!@#", "Aa123456!", "A123456s!", "sa123456", "1q2w3e", "Charge123", "Aa123456789"}

// File struct 文件操作
type File struct {
	ctx          context.Context
	configPath   string
	downloadPath string
}

func (f *File) Startup(ctx context.Context) {
	f.ctx = ctx
}

func NewFile() *File {
	home := util.HomeDir()
	return &File{
		configPath:   home + "/slack/config",
		downloadPath: home + "/Downloads/",
	}
}

// 创建爆破字典
func init() {
	var userPath = util.HomeDir() + "/slack/portburte/username"
	var passPath = util.HomeDir() + "/slack/portburte/password"
	os.MkdirAll(userPath, 0777)
	os.MkdirAll(passPath, 0777)
	for name, dict := range Userdict {
		file := fmt.Sprintf("%s/%s.txt", userPath, name)
		// 文件不存在则需要创建
		if _, err := os.Stat(file); err != nil {
			os.WriteFile(file, []byte(strings.Join(dict, "\n")), 0644)
		}
	}
	os.WriteFile(fmt.Sprintf("%s/password.txt", passPath), []byte(strings.Join(Passwords, "\n")), 0644)
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

func (f *File) DirectoryDialog() string {
	selection, err := runtime.OpenDirectoryDialog(f.ctx, runtime.OpenDialogOptions{
		Title: "选择文件夹",
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// selection会返回保存的文件路径+文件名 例如/Users/xxx/Downloads/test.xlsx
func (f *File) SaveFileDialog(filename string) string {
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

func (f *File) IsMacOS() bool {
	return rt.GOOS == "darwin"
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

type FileListInfo struct {
	Path     string // 完整路径
	Name     string // 带名称后缀
	BaseName string // 基础名称
	ModTime  string // 修改时间
	Size     int64  // 大小
}

func (f *File) List(folders []string) []FileListInfo {
	var files []FileListInfo
	for _, folder := range folders {
		if folder == "" {
			continue
		}
		fileinfo, err := os.Stat(folder)
		if os.IsNotExist(err) {
			gologger.Error(f.ctx, fmt.Sprintf("path %s not exist", folder))
			continue
		}

		if !fileinfo.IsDir() {
			filename := filepath.Base(folder)
			baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
			files = append(files, FileListInfo{
				Path:     folder,
				Name:     filename,
				BaseName: baseName,
				ModTime:  fileinfo.ModTime().Format("2006-01-02 15:04:05"),
				Size:     fileinfo.Size(),
			})
			continue
		}

		filepath.Walk(folder, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// 提取文件名
				filename := filepath.Base(p)
				// 去除文件后缀
				baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
				files = append(files, FileListInfo{
					Path:     p,
					Name:     filename,
					BaseName: baseName, // 存储去除后缀的文件名
					ModTime:  info.ModTime().Format("2006-01-02 15:04:05"),
					Size:     info.Size(),
				})
			}
			return nil
		})
	}
	return files
}

func (f *File) ListDir(folder string) []string {
	var dirs []string
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		gologger.Error(f.ctx, fmt.Sprintf("path %s not exist", folder))
		return nil
	}
	filepath.Walk(folder, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, p)
		}
		return nil
	})
	return dirs
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

func (f *File) UpdatePocFile(version string) bool {
	var defaultFile = util.HomeDir() + "/slack/"
	os.MkdirAll(defaultFile, 0777)
	configFileZip := fmt.Sprintf("%sv%s/config.zip", update.LastestPocUrl, version)
	_, err := update.NewDownload(f.ctx, configFileZip, defaultFile, "pocDownloadProgress", "")
	if err != nil {
		gologger.Error(f.ctx, err)
		return false
	}
	uz := fileutil.NewUnzip()
	if _, err := uz.Extract(defaultFile+"config.zip", defaultFile); err != nil {
		gologger.Error(f.ctx, err)
		return false
	}
	os.Remove(util.HomeDir() + "/slack/config.zip")
	return true
}

func (f *File) InitConfig() bool {
	return update.InitConfig(f.ctx)
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

func (*File) SaveToTempFile(content string) string {
	tempDir := os.TempDir()
	tempFileName := fmt.Sprintf("%stemp_%d.txt", tempDir, time.Now().UnixNano())
	if err := os.WriteFile(tempFileName, []byte(content), 0644); err != nil {
		return ""
	}
	return tempFileName
}

func (a *App) DownloadCyberChef(url string) error {
	cyber := util.HomeDir() + "/slack/CyberChef.zip"
	fileName, err := update.NewDownload(a.ctx, url, a.defaultPath, "downloadProgress", "")
	if err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "downloadComplete", fileName)
	uz := fileutil.NewUnzip()
	if _, err := uz.Extract(cyber, a.defaultPath); err != nil {
		return err
	}
	return os.Remove(cyber)
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
		exec.Command("xattr", "-c", f.downloadPath+filename).Run()
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
		cmd.Run()
	}
}

func (f *File) RemoveFile(file string) bool {
	return os.Remove(file) == nil
}

var (
	na         = util.HomeDir() + "/slack/navogation.json"
	navigation []structs.Navigation
)

func (f *File) GetLocalNaConfig() *[]structs.Navigation {
	if !f.CheckFileStat(na) {
		os.Create(na)
		gologger.Error(f.ctx, "[AppLauncher] Can't create navogation.json")
		return &navigation
	}
	b, err := os.ReadFile(na)
	if err != nil {
		gologger.Error(f.ctx, fmt.Errorf("[AppLauncher] Read navogation.json err: %w", err))
		return &navigation
	}
	if err := json.Unmarshal(b, &navigation); err != nil {
		gologger.Error(f.ctx, fmt.Errorf("[AppLauncher] Unmarshal navogation.json err: %w", err))
		return &navigation
	}
	return &navigation
}

func (f *File) InsetGroupNavigation(n structs.Navigation) bool {
	navigation = append(navigation, n)
	return fileutil.SaveJsonWithFormat(f.ctx, na, navigation)
}

func (f *File) InsetItemNavigation(groupName string, child structs.Children) bool {
	for i, n := range navigation {
		if n.Name == groupName {
			navigation[i].Children = append(n.Children, child)
		}
	}
	return fileutil.SaveJsonWithFormat(f.ctx, na, navigation)
}

func (f *File) SaveNavigation(n []structs.Navigation) bool {
	navigation = n
	return fileutil.SaveJsonWithFormat(f.ctx, na, navigation)
}

func (f *File) OpenFolder(filepath string) string {
	filepath, err := getDirectoryPath(filepath)
	if err != nil {
		return err.Error()
	}
	var cmd *exec.Cmd
	switch rt.GOOS {
	case "windows":
		cmd = exec.Command("explorer", "/select,"+filepath)
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

// 给定一个路径，在打开所在文件夹的命令行
func (f *File) OpenTerminal(filepath string) string {
	filepath, err := getDirectoryPath(filepath)
	if err != nil {
		return err.Error()
	}
	var cmd *exec.Cmd
	switch rt.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "cmd", "/K", "cd /d", filepath)
	case "darwin":
		cmd = exec.Command("open", "-a", "Terminal", filepath)
	default:
		cmd = exec.Command("gnome-terminal", "--working-directory", filepath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err.Error()
	}
	return ""
}

func (f *File) RunApp(types, file, target string) {
	var cmd *exec.Cmd
	name := filepath.Base(file)
	switch types {
	// 超链接类型
	case "Link":
		runtime.BrowserOpenURL(f.ctx, file)
		return
	case "JAR":
		cmd = exec.Command("java", "-jar", name)
	case "APP":
		if rt.GOOS != "windows" {
			runtime.MessageDialog(f.ctx, runtime.MessageDialogOptions{
				Title:         "提示",
				Message:       "仅提供Windows用户运行App类型",
				Type:          runtime.InfoDialog,
				DefaultButton: "Ok",
			})
			return
		}
		switch filepath.Ext(file) {
		case ".exe":
			cmd = exec.Command(name)
		default:
			cmd = exec.Command("cmd", "/c", "start", name)
		}
	default:
		if target == "" {
			f.OpenTerminal(file)
			return
		}
		target = strings.ReplaceAll(target, "%path%", name)
		args := strings.Split(target, " ")
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Dir = filepath.Dir(file)
	if _, err := os.Stat(cmd.Dir); os.IsPermission(err) {
		gologger.Debug(f.ctx, "Insufficient permissions for directory: "+cmd.Dir)
		return
	}
	go func(localCmd *exec.Cmd) {
		if err := localCmd.Run(); err != nil {
			gologger.Debug(f.ctx, err)
			return
		}
	}(cmd)
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

// 从网址获取到图标后保存在本地，并返回保存路径
// 为了防止图片过大，还需要压缩图片
func (f *File) GenerateFaviconBase64WithOnline(rawURL string) string {
	u, _ := url.Parse(rawURL)
	faviconURL, err := webscan.GetFaviconFullLink(u, clients.NewRestyClient(nil, true))
	if err != nil {
		return ""
	}
	resp, err := clients.SimpleGet(faviconURL, clients.NewRestyClient(nil, true))
	if err != nil {
		return ""
	}
	if resp.StatusCode() != 200 {
		return ""
	}
	return compressPictures(rawURL, resp.Body())
}

// 读取本读图片
func (f *File) GenerateFaviconBase64(filePath string) string {
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return compressPictures(filePath, data)
}

func (f *File) AutoGenerateFavicon(filePath string) string {
	if strings.HasSuffix(filePath, ".lnk") {
		Lnk, err := lnk.File(filePath)
		if err != nil {
			return ""
		}
		filePath = Lnk.LinkInfo.LocalBasePath
	}
	// 创建一个 Buffer 来存储图标数据
	var buf bytes.Buffer

	// 提取图标并写入到 Buffer
	err := fico.F2ICO(&buf, filePath, fico.Config{Format: "ico", Width: 128, Height: 128})
	if err != nil {
		gologger.Debug(f.ctx, "[AppLauncher] 读取应用图标失败 "+err.Error())
		return ""
	}
	// 将图标的字节数据进行 Base64 编码
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

// path可以是文件路径也可以是url,主要用来判断文件格式
func compressPictures(path string, data []byte) string {
	// 如果是svg图标，则直接返回base64，因为svg图标文件比较小
	if strings.HasSuffix(path, ".svg") {
		return fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString(data))
	}
	// 小于20kb可以直接返回
	if len(data) <= 1024*20 {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	}
	// 尝试解析不同格式的图像
	var img image.Image
	var err error
	if strings.HasSuffix(path, ".ico") {
		img, err = ico.Decode(bytes.NewReader(data)) // 处理ICO文件
	} else {
		img, _, err = image.Decode(bytes.NewReader(data)) // 默认解码
	}
	if err != nil {
		return ""
	}
	// 将图像缩放为64x64
	resizedImg := resize.Resize(64, 64, img, resize.Lanczos3)
	// 将缩放后的图像编码为Base64
	var buf bytes.Buffer
	err = png.Encode(&buf, resizedImg)
	if err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (f *File) SaveDataToFile(data interface{}) bool {
	content, _ := json.MarshalIndent(data, "", "  ")
	if err := os.WriteFile(f.UserHomeDir()+"/slack/config.json", content, 0777); err != nil {
		return false
	}
	return true
}

func (f *File) ReadLocalStore() map[string]interface{} {
	var data map[string]interface{}
	content, _ := os.ReadFile(f.UserHomeDir() + "/slack/config.json")
	if err := json.Unmarshal(content, &data); err != nil {
		return nil
	}
	return data
}

func (f *File) NetworkCardInfo() (networks []structs.NetwordCard) {
	ifaces, err := net.Interfaces()
	if err != nil {
		gologger.Error(f.ctx, err)
		return
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			gologger.Error(f.ctx, err)
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					networks = append(networks, structs.NetwordCard{
						Name: iface.Name,
						IP:   v.IP.String(),
					})
				}
			}
		}
	}
	return
}

type Tree struct {
	ID       string         `json:"id"`
	Label    string         `json:"label"`
	IsDir    bool           `json:"isDir"`
	Hits     map[string]int `json:"hits,omitempty"` // 记录命中次数
	Children []Tree         `json:"children,omitempty"`
}

func (f *File) BuildTree(root string, keywords, blackList []string) Tree {
	info, err := os.Stat(root)
	if err != nil {
		return Tree{}
	}

	rootNode := Tree{
		ID:    root,
		Label: info.Name(),
		IsDir: info.IsDir(),
	}

	// 不是文件夹就进行敏感词检测
	if !info.IsDir() {
		// 检测是否在黑名单中
		for _, black := range blackList {
			if strings.HasSuffix(root, black) {
				return rootNode
			}
		}
		rootNode.Hits = scanFileForKeywords(root, keywords)
		return rootNode
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return rootNode
	}

	for _, entry := range entries {
		childPath := filepath.Join(root, entry.Name())
		_, err := os.Stat(childPath)
		if err != nil {
			continue
		}

		childNode := f.BuildTree(childPath, keywords, blackList)
		rootNode.Children = append(rootNode.Children, childNode)
	}

	return rootNode
}

func scanFileForKeywords(filePath string, keywords []string) map[string]int {
	hitCounts := make(map[string]int)
	file, err := os.Open(filePath)
	if err != nil {
		return hitCounts // 读取失败直接返回
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				hitCounts[keyword]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("扫描文件错误:", err)
	}

	return hitCounts
}
