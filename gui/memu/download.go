package memu

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"slack/common"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/gui/global"
	"slack/gui/mytheme"
	"slack/lib/qqwry"
	"slack/lib/util"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	update "github.com/fynelabs/selfupdate"
)

const (
	remoteClientVersion = "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/version"
	updateClientContent = "https://gitee.com/the-temperature-is-too-low/Slack/raw/main/update" // 更新内容
	// https://gitee.com/the-temperature-is-too-low/Slack/releases/download/v1.4.3/slack.exe
	lastestClinetUrl = "https://gitee.com/the-temperature-is-too-low/Slack/releases/download/"

	localPocVersion  = "./config/afrog-pocs/version"
	localPocDir      = "./config/afrog-pocs/"
	remotePocVersion = "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version"
	updatePocContent = "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/update"
	lastestPocUrl    = "https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/"
)

func DowdloadQqwry() {
	if err := qqwry.Download(Qqwrypath); err == nil {
		dialog.ShowInformation("", "update success", global.Win)
	} else {
		dialog.ShowError(fmt.Errorf("update failed %v", err), global.Win)
	}
}

func CheckUpdate(removeTarget, localVersion string) (string, error) {
	r, err := http.Get(removeTarget)
	if err != nil {
		return "", err
	}
	b, err2 := io.ReadAll(r.Body)
	if err2 != nil {
		return "", err2
	}
	remoteVersion := string(b)
	if remoteVersion == localVersion {
		return "", errors.New("当前已是最新版本 " + string(b))
	}
	return remoteVersion, nil
}

// numOfTimes 为0时不显示dialog提示
func ConfrimUpdateClient(numOfTimes int) {
	if numOfTimes == 0 {
		return
	}
	if version, err := CheckUpdate(remoteClientVersion, common.Version); err != nil {
		dialog.ShowInformation("提示", "客户端"+err.Error(), global.Win)
	} else {
		r, err := http.Get(updateClientContent)
		if err != nil {
			return
		}
		b, err2 := io.ReadAll(r.Body)
		if err2 != nil {
			return
		}
		dp := widget.NewProgressBarInfinite()
		dp.Hide()
		l := &widget.Label{Text: string(b), Wrapping: fyne.TextWrapBreak}
		content := container.NewBorder(nil, dp, nil, nil, container.NewVScroll(l))
		custom.ShowCustomDialog(mytheme.UpdateIcon(), "更新提醒,最新版本:"+version, "立即更新", content, func() {
			dp.Show()
			content.Refresh()
			if err3 := UpdateClinet(version); err3 != nil {
				l.SetText(fmt.Sprintf("更新失败: %v", err3))
			}
		}, fyne.NewSize(400, 300))
	}
}

func UpdateClinet(latestVersion string) error {
	var binaryFileName string
	switch runtime.GOOS {
	case "windows":
		binaryFileName = "slack.exe"
	case "linux":
		binaryFileName = "slack_linux_amd64"
	case "darwin":
		binaryFileName = "slack_darwin_amd64"
	}
	if err := doUpdate(lastestClinetUrl + "v" + latestVersion + "/" + binaryFileName); err != nil {
		return err
	}
	custom.ShowCustomDialog(theme.InfoIcon(), "提示", "立即重启客户端", custom.NewCenterLable("更新成功!!"), func() {
		go func() {
			cmd := exec.Command(os.Args[0])
			err := cmd.Start()
			if err != nil {
				logger.Error(err)
			}
			// 退出当前的进程
			os.Exit(0)
		}()
	}, fyne.NewSize(100, 50))
	return nil
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Printf("Failed to rollback from bad update: %v", rerr)
		}
	}
	return err
}

func ConfrimUpdatePoc() {
	b, err := os.ReadFile(localPocVersion)
	if err != nil {
		dialog.ShowError(errors.New("检测本地漏洞版本失败,请检查afrog-pocs下的version文件是否存在"), global.Win)
		return
	}
	if version, err := CheckUpdate(remotePocVersion, string(b)); err != nil {
		dialog.ShowInformation("提示", "POC"+err.Error(), global.Win)
	} else {
		r, err := http.Get(updatePocContent)
		if err != nil {
			return
		}
		b, err2 := io.ReadAll(r.Body)
		if err2 != nil {
			return
		}
		dp := widget.NewProgressBarInfinite()
		dp.Hide()
		l := &widget.Label{Text: string(b), Wrapping: fyne.TextWrapBreak}
		content := container.NewBorder(nil, dp, nil, nil, container.NewVScroll(l))
		custom.ShowCustomDialog(mytheme.UpdateIcon(), "更新提醒,最新版本:"+version, "立即更新", content, func() {
			dp.Show()
			content.Refresh()
			if err := UpdatePoc(version); err != nil {
				l.SetText(fmt.Sprintf("更新失败: %v", err))
			}

		}, fyne.NewSize(400, 300))
	}
}

// https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/v0.0.2/afrog-pocs.zip
func UpdatePoc(latestVersion string) error {
	workflow := lastestPocUrl + "v" + latestVersion + "/workflow.yaml"
	webfinger := lastestPocUrl + "v" + latestVersion + "/webfiner.yaml"
	pocs := lastestPocUrl + "v" + latestVersion + "/afrog-pocs.zip"
	if _, err := download(workflow, "./config/"); err != nil {
		return err
	}
	if _, err := download(webfinger, "./config/"); err != nil {
		return err
	}
	fileName, err := download(pocs, "./config/")
	if err != nil {
		return err
	}
	if err := os.RemoveAll(localPocDir); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	uz := util.NewUnzip()
	if _, err := uz.Extract("./config/"+fileName, "./config/"); err != nil {
		return fmt.Errorf("afrog-poc decompression failed. %s", err.Error())
	}
	custom.ShowCustomDialog(theme.InfoIcon(), "提示", "重启客户端重新加载漏洞", custom.NewCenterLable("更新成功!!"), func() {
		go func() {
			cmd := exec.Command(os.Args[0])
			err := cmd.Start()
			if err != nil {
				logger.Error(err)
			}
			// 退出当前的进程
			os.Exit(0)
		}()
	}, fyne.NewSize(100, 50))
	return nil
}

func download(target, dest string) (string, error) {
	fileName := path.Base(target)
	res, err := http.Get(target)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(dest + fileName)
	if err != nil {
		return "", err
	}
	//获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	return fileName, nil
}
