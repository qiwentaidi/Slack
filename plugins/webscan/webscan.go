package webscan

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slack/common"
	"slack/common/logger"
	"slack/gui/global"
	"slack/lib/poc"
	"slack/lib/result"
	"slack/lib/runner"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const defaultActiveDirectory = "./config/active-detect/"

func WebScan(target []string, keyword, severity string, active, finger2poc bool, vulpath *widget.Entry) {
	var proxys string
	if active || finger2poc { // 如果主动扫描或者指纹POC扫描开启则将severity的值初始化，为空时不会对POC进行筛选
		severity = ""
	}
	if common.Profile.Proxy.Enable {
		proxys = fmt.Sprintf("%v://%v:%v", common.Profile.Proxy.Mode, common.Profile.Proxy.Address, common.Profile.Proxy.Port)
	}
	options := common.NewOptions(target, keyword, severity, proxys)
	r, err := runner.NewRunner(options)
	if err != nil {
		logger.Error(err)
	}
	var (
		lock   = sync.Mutex{}
		number uint32
	)
	r.OnResult = func(result *result.Result) {
		defer func() {
			atomic.AddUint32(&options.CurrentCount, 1)
			global.ProgerssWebscan.SetText(fmt.Sprintf("%d/%d", options.CurrentCount, options.Count))
		}()
		if result.IsVul {
			lock.Lock()
			atomic.AddUint32(&number, 1)
			result.PrintColorResultInfoConsole(fmt.Sprint(number))
			r.Report.SetResult(result)
			r.Report.Append(fmt.Sprint(number))
			lock.Unlock()
		}
	}
	// 如果开启及主动探测，就只扫指纹和这些主动的POC即可 afrog的-P参数就是只扫指定目录的POC刚好符合主动探测只扫固定POC的逻辑
	if active {
		poc.LocalTestList = []string{}
		options.PocFile = defaultActiveDirectory
		poc.SelectFolderReadPocPath(options.PocFile)
		r.Execute(active, finger2poc, target)
	} else if finger2poc {
		for target, fingers := range common.UrlFingerMap { // url和指纹的键值对
			if len(fingers) > 0 {
				options.PocFile = "111" //此处一定得传入值，传啥无所谓
				poc.LocalTestList = []string{}
				FingerPocFile, _ := poc.FingerPocFilepath(fingers)
				poc.LocalTestList = append(poc.LocalTestList, FingerPocFile...)
				r.Execute(active, finger2poc, []string{target})
			}
		}
	} else if !vulpath.Disabled() && vulpath.Text != "" {
		file, err := os.Stat(vulpath.Text)
		if err != nil {
			dialog.ShowError(errors.New("所选文件不存在"), global.Win)
			return
		}
		if !file.IsDir() && (path.Ext(vulpath.Text) == ".yaml" || path.Ext(vulpath.Text) == ".yml") {
			poc.LocalTestList = []string{vulpath.Text} // 单poc扫描
		} else {
			poc.LocalTestList = []string{}
			poc.SelectFolderReadPocPath(vulpath.Text) // 文件夹poc扫描
		}
		r.Execute(true, false, target)
	} else {
		r.Execute(active, finger2poc, target)
	}
}
