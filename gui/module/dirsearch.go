package module

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slack/common"
	client2 "slack/common/client"
	"slack/common/logger"
	"slack/common/proxy"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/mytheme"
	"slack/lib/util"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultDict   = "./config/dirsearch/"
	dirsearchInfo = `1、状态码: 支持200,300或200-300,400-500两类形式过滤

2、拓展名: 会将字典中%EXT%字段替换，不指定则去除有关%EXT%字段

3、字典路径: 指定字典后，优先级大于下方内置字典

4、内置字典: 可以通过打开目录新增文件自定义内置字典

5、扫描结果会自动过滤重复长度出现5次以上的数据，防止数据显示过多`
)

var (
	DirsearchProgress *widget.ProgressBar
	DirResult         = [][]string{{"#", "状态码", "长度", "目录", "跳转路径", ""}}
	lengthNums        = make(map[int]int)
	dirnum, id1       int64
)

func DirSearchUI() *fyne.Container {
	target := &widget.Entry{PlaceHolder: "请输入URL地址"}
	codeFilter := widget.NewEntry()
	ext := &widget.Entry{Text: "php,aspx,asp,jsp,html,js"}
	scan := &widget.Button{Text: "开始任务", Importance: widget.HighImportance}
	DirsearchProgress = widget.NewProgressBar()
	DirsearchProgress.TextFormatter = func() string {
		if DirsearchProgress.Value < DirsearchProgress.Min {
			DirsearchProgress.Value = DirsearchProgress.Min
		}
		if DirsearchProgress.Value > DirsearchProgress.Max {
			DirsearchProgress.Value = DirsearchProgress.Max
		}
		delta := float32(DirsearchProgress.Max - DirsearchProgress.Min)
		ratio := float32(DirsearchProgress.Value-DirsearchProgress.Min) / delta
		return fmt.Sprintf("%v/%v | %v", DirsearchProgress.Value, DirsearchProgress.Max, strconv.Itoa(int(ratio*100))+"%")
	}
	global.DirDictText = custom.NewFileEntry("")
	sl := custom.NewCheckListBox(LoadLocalDict())
	info := widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		custom.ShowCustomDialog(theme.InfoIcon(), "提示", "", widget.NewLabel(dirsearchInfo), nil, fyne.NewSize(400, 300))
	})
	thread := custom.NewNumEntry("20")
	method := widget.NewSelect([]string{"GET", "POST", "HEAD", "OPTIONS"}, nil)
	method.SetSelected("GET")
	rule := container.NewBorder(nil, container.NewBorder(nil, nil, widget.NewLabel("字典路径:"), scan, global.DirDictText), nil, nil, container.NewGridWithColumns(3,
		custom.NewFormItem("选状态码:", container.NewBorder(nil, nil, nil, info, codeFilter)),
		custom.NewFormItem("拓展名:", ext),
		container.NewGridWithColumns(2,
			custom.NewFormItem("线程:", thread),
			custom.NewFormItem("模式:", method),
		),
	))

	t := custom.NewTableWithUpdateHeader1(&DirResult, []float32{50, 70, 100, 400, 500, 0})
	scan.OnTapped = func() {
		go func() {
			t, err := common.ParseURLWithoutSlash(target.Text)
			if target.Text != "" && err == nil {
				if err := TestTarget(t); err != nil {
					dialog.ShowError(err, global.Win)
					return
				}
				DirResult = DirResult[:1]
				var options = []string{}
				var count = []string{}
				dirnum = 0
				if ext.Text != "" && strings.Contains(ext.Text, ",") {
					options = strings.Split(ext.Text, ",")
				} else if ext.Text != "" {
					options = []string{ext.Text}
				}
				if global.DirDictText.Text == "" {
					if len(sl.Selected) > 0 {
						for _, d := range sl.Selected {
							count = append(count, common.ParseDirectoryDict(d, "%EXT%", options)...)
						}
					} else {
						count = common.ParseDirectoryDict(defaultDict+"dicc.txt", "%EXT%", options)
					}
				} else {
					count = common.ParseDirectoryDict(global.DirDictText.Text, "%EXT%", options)
				}
				DirsearchProgress.SetValue(0)
				DirsearchProgress.Max = float64(len(count))
				MultiThread(method.Selected, util.RemoveIllegalChar(t), count, common.ParsePort(codeFilter.Text), thread.Number)
			} else {
				dialog.ShowError(fmt.Errorf("URL为空或存在错误 %v", err), global.Win)
			}
		}()
	}
	open := widget.NewButtonWithIcon("打开字典路径", theme.FolderOpenIcon(), func() {
		dir, _ := os.Getwd()
		util.OpenFolder(dir + "\\config\\dirsearch\\")
	})
	update := widget.NewButtonWithIcon("刷新", mytheme.UpdateIcon(), func() {
		sl.Options = LoadLocalDict()
		sl.Refresh()
	})
	dict := container.NewBorder(container.NewHBox(&widget.Label{Text: "内置字典(默认扫描dicc.txt)", Truncation: fyne.TextTruncateOff}, layout.NewSpacer(), open, update), nil, nil, nil, sl)
	hbox := container.NewVSplit(container.NewBorder(rule, nil, nil, nil, container.NewBorder(nil, nil, nil, widget.NewCard("", "", dict), target)), custom.Frame(t))
	hbox.Offset = 0.3
	return container.NewBorder(nil, DirsearchProgress, nil, nil, hbox)
}

func MultiThread(method, url string, paths []string, filter []int, thread int) {
	var wg sync.WaitGroup
	limiter := make(chan bool, thread) // 限制协程数量
	c := client2.NotFollowClient()
	if common.Profile.Proxy.Enable {
		c = proxy.SelectProxy(&common.Profile)
	}
	for _, path := range paths {
		wg.Add(1)
		limiter <- true
		go SimpleResp(method, url, path, c, filter, limiter, &wg)
	}
	wg.Wait()
}

func SimpleResp(method, url, path string, client *http.Client, filter []int, limiter chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	data := PathRequest(method, url+path, client)
	atomic.AddInt64(&id1, 1)
	DirsearchProgress.SetValue(float64(id1))
	DirsearchProgress.Refresh()
	if data.Status != 0 { // 响应正常
		mutex.Lock() // map不允许并发写入
		lengthNums[data.Length]++
		if lengthNums[data.Length] <= 5 && util.ArrayContains[int](data.Status, filter) {
			atomic.AddInt64(&dirnum, 1)
			DirResult = append(DirResult, []string{fmt.Sprintf("%v", dirnum), fmt.Sprintf("%v", data.Status), fmt.Sprintf("%v", data.Length), path, data.Location, ""})
		}
		mutex.Unlock()
	}
	<-limiter
}

type PathData struct {
	Status   int    // 状态码
	Location string // server信息
	Length   int    // 主体内容
}

func PathRequest(method, url string, client *http.Client) *PathData {
	var pd PathData // 将响应头和响应头的数据存储到结构体中
	resp, body, err := client2.NewHttpWithDefaultHead("GET", url, client)
	if err != nil || resp == nil {
		logger.Error(err)
		return &pd
	}
	pd.Status = resp.StatusCode
	pd.Length = len(body)
	pd.Location = resp.Header.Get("Location")
	return &pd
}

// 测试URL是否可达
func TestTarget(target string) error {
	_, err := url.Parse(target)
	if err != nil {
		return err
	}
	if _, err2 := http.Get(target); err2 != nil {
		return err
	}
	return nil
}

func LoadLocalDict() []string {
	currentDicts := []string{}
	filepath.WalkDir(defaultDict, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".txt") {
			currentDicts = append(currentDicts, path)
		}
		return nil
	})
	return currentDicts
}
