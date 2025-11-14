// AppLauncher.go | 2025-05-15 | 将原有的file.go中的启动器代码抽离成单独文件
package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	rt "runtime"
	"slack-wails/core/webscan"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/utils"
	"slack-wails/lib/utils/fileutil"
	"strings"

	"github.com/qiwentaidi/clients"

	"github.com/mat/besticon/v3/ico"
	"github.com/nfnt/resize"
	"github.com/orcastor/fico"
	lnk "github.com/parsiya/golnk"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	configFile    = utils.HomeDir() + "/slack/navogation.json"
	jdkConfigFile = utils.HomeDir() + "/slack/jdkConfig.json"
	navigation    []structs.Navigation
	jdkConfig     []structs.JdkConfig
)

func (f *File) GetLocalNaConfig() *[]structs.Navigation {
	if !f.CheckFileStat(configFile) {
		os.Create(configFile)
		gologger.Error(f.ctx, "[AppLauncher] Can't create navogation.json")
		return &navigation
	}
	b, err := os.ReadFile(configFile)
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

// 获取系统变量配置文件内容，主要用于配置JDK环境变量
func (f *File) GetJdkConfig() *[]structs.JdkConfig {
	if !f.CheckFileStat(jdkConfigFile) {
		os.Create(jdkConfigFile)
		gologger.Error(f.ctx, "[AppLauncher] Can't create jdkConfig.json")
	}
	b, err := os.ReadFile(jdkConfigFile)
	if err != nil {
		gologger.Error(f.ctx, fmt.Errorf("[AppLauncher] Read jdkConfig.json err: %w", err))
		return &jdkConfig
	}
	if err := json.Unmarshal(b, &jdkConfig); err != nil {
		gologger.Error(f.ctx, fmt.Errorf("[AppLauncher] Unmarshal jdkConfig.json err: %w", err))
		return &jdkConfig
	}
	return &jdkConfig
}

func (f *File) InsetJdkConfig(jdk structs.JdkConfig) bool {
	jdkConfig = append(jdkConfig, jdk)
	return fileutil.SaveJsonWithFormat(f.ctx, jdkConfigFile, jdkConfig)
}

func (f *File) SaveJdkConfig(jdks []structs.JdkConfig) bool {
	jdkConfig = jdks
	return fileutil.SaveJsonWithFormat(f.ctx, jdkConfigFile, jdkConfig)
}

func (f *File) InsetGroupNavigation(n structs.Navigation) bool {
	navigation = append(navigation, n)
	return fileutil.SaveJsonWithFormat(f.ctx, configFile, navigation)
}

func (f *File) InsetItemNavigation(groupName string, child structs.Children) bool {
	for i, n := range navigation {
		if n.Name == groupName {
			navigation[i].Children = append(n.Children, child)
		}
	}
	return fileutil.SaveJsonWithFormat(f.ctx, configFile, navigation)
}

func (f *File) SaveNavigation(n []structs.Navigation) bool {
	navigation = n
	return fileutil.SaveJsonWithFormat(f.ctx, configFile, navigation)
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

func (f *File) RunApp(item structs.Children) {
	var cmd *exec.Cmd
	name := filepath.Base(item.Path)
	switch item.Type {
	// 超链接类型
	case "Link":
		runtime.BrowserOpenURL(f.ctx, item.Path)
		return
	case "JAR":
		if item.Jdk != "" {
			cmd = exec.Command(item.Jdk, "-jar", name)
			break
		}
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
		switch filepath.Ext(item.Path) {
		case ".exe":
			if _, err := os.Stat(item.Path); err == nil {
				// 优先完整路径运行
				cmd = exec.Command(item.Path)
			} else {
				// fallback：用 cmd /c start 方式
				cmd = exec.Command("cmd", "/c", "start", name)
				cmd.Dir = filepath.Dir(item.Path)
			}
		default:
			cmd = exec.Command("cmd", "/c", "start", name)
		}
	default:
		if item.Target == "" {
			f.OpenTerminal(item.Path)
			return
		}
		item.Target = strings.ReplaceAll(item.Target, "%path%", name)
		args := strings.Split(item.Target, " ")
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Dir = filepath.Dir(item.Path)
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
