package module

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slack/common"
	"slack/common/logger"
	"slack/common/proxy"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/lib/util"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultDict = "./config/dirsearch/dicc.txt"
	springDict  = "./config/dirsearch/spring.txt"
	//backupDict  = "./config/dirsearch/backup.txt"

	dirsearchInfo = `1、状态码: 支持200,300或200-300,400-500两类形式过滤

2、拓展名: 会将字典中%EXT%字段替换，不指定默认替换成php, aspx, asp,jsp, html, js

3、字典路径: 指定字典后，优先级大于下方默认字典`
)

var (
	DirsearchProgress *widget.ProgressBar
	DirResult         = [][]string{{"#", "状态码", "长度", "目录", "跳转路径"}}
	lengthNums        = make(map[int]int)
	dirnum, id1       int64
)

func DirSearchUI() *fyne.Container {
	target := widget.NewEntry()
	target.PlaceHolder = "请输入URL地址"
	codeFilter := widget.NewEntry()
	ext := widget.NewEntry()
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
	sl := custom.NewSelectList("请选择字典,默认(扫描dicc.txt)", []string{defaultDict, springDict})
	info := widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
		custom.ShowCustomDialog(theme.InfoIcon(), "提示", "", widget.NewLabel(dirsearchInfo), nil, fyne.NewSize(400, 300))
	})
	rule := widget.NewForm(
		widget.NewFormItem("状态码:", container.NewBorder(nil, nil, nil, info, codeFilter)),
		widget.NewFormItem("拓展名:", ext),
		widget.NewFormItem("字典路径:", global.DirDictText),
	)
	t := custom.NewTableWithUpdateHeader1(&DirResult, []float32{50, 70, 100, 400, 500})
	scan.OnTapped = func() {
		go func() {
			s, err := common.ParseURLWithoutSlash(target.Text)
			if target.Text != "" && err == nil {
				if err := TestTarget(s); err != nil {
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
				} else {
					options = []string{"php", "aspx", "asp", "jsp", "html", "js"}
				}
				if global.DirDictText.Text == "" {
					if len(sl.Checked) > 0 {
						for _, d := range sl.Checked {
							count = append(count, common.ParseDirectoryDict(d, "%EXT%", options)...)
						}
					} else {
						count = common.ParseDirectoryDict(defaultDict, "%EXT%", options)
					}
				} else {
					count = common.ParseDirectoryDict(global.DirDictText.Text, "%EXT%", options)
				}
				DirsearchProgress.SetValue(0)
				DirsearchProgress.Max = float64(len(count))
				SimpleHttp(s, count, common.ParsePort(codeFilter.Text))
			} else {
				dialog.ShowError(fmt.Errorf("URL为空或存在错误 %v", err), global.Win)
			}
		}()
	}
	hbox := container.NewHSplit(rule, container.NewBorder(nil, nil, nil, scan, target))
	hbox.Offset = 0.3
	return container.NewBorder(container.NewBorder(nil, sl, nil, nil, hbox), DirsearchProgress, nil, nil, custom.Frame(t))
}

func SimpleHttp(url string, paths []string, filter []int) {
	var wg sync.WaitGroup
	limiter := make(chan bool, 10) // 限制协程数量
	client := proxy.NotFollowClient()
	if common.Profile.Proxy.Enable {
		client = proxy.SelectProxy(&common.Profile)
	}
	for _, path := range paths {
		wg.Add(1)
		limiter <- true
		go SimpleResp(url, path, client, filter, limiter, &wg)
	}
	wg.Wait()
}

func SimpleResp(url, path string, client *http.Client, filter []int, limiter chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	data := PathRequest(url+path, client)
	atomic.AddInt64(&id1, 1)
	DirsearchProgress.SetValue(float64(id1))
	DirsearchProgress.Refresh()
	if data.Status != 0 { // 响应正常
		mutex.Lock() // map不允许并发写入
		lengthNums[data.Length]++
		if lengthNums[data.Length] <= 5 && util.ArrayContains[int](data.Status, filter) {
			atomic.AddInt64(&dirnum, 1)
			DirResult = append(DirResult, []string{fmt.Sprintf("%v", dirnum), fmt.Sprintf("%v", data.Status), fmt.Sprintf("%v", data.Length), path, data.Location})
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

func PathRequest(url string, client *http.Client) *PathData {
	var pd PathData // 将响应头和响应头的数据存储到结构体中
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	if err != nil {
		logger.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Debug(err)
	}
	if resp != nil && resp.StatusCode != 302 { // 过滤重定向次数过多的
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Debug(err)
		}
		defer resp.Body.Close()
		pd.Status = resp.StatusCode
		pd.Length = len(body)
		pd.Location = resp.Header.Get("Location")
	}
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
