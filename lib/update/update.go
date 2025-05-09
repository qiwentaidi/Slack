package update

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"slack-wails/lib/clients"
	"slack-wails/lib/fileutil"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"

	"github.com/minio/selfupdate"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const LastestPocUrl = "https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/"

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

func UpdateClientWindows(ctx context.Context, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	totalSize := resp.ContentLength
	downloadedSize := int64(0)
	reader := bufio.NewReaderSize(resp.Body, 32*1024)
	// 创建一个缓冲区来模拟文件写入
	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			downloadedSize += int64(n)
			progress := float64(downloadedSize) / float64(totalSize) * 100
			roundedProgress := roundToTwoDecimals(progress)
			runtime.EventsEmit(ctx, "clientDownloadProgress", roundedProgress)
			// 将数据写入缓冲区，不然由于进度条消耗，resp.Body的数据会为空
			if _, err := buffer.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return selfupdate.Apply(&buffer, selfupdate.Options{})
}

// 通过execName执行文件名称是否为空判断，是否需要将远程下载的文件重命名为本地文件名称 -- 主要是支持Linux覆盖文件以及重新启动
func NewDownload(ctx context.Context, target, dest, events, execName string) (string, error) {
	fileName := path.Base(target)
	if execName != "" {
		fileName = execName
	}
	res, err := http.Get(target)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 468 {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Title:         "提示",
			Message:       "当前网络环境不支持请求更新文件, 请切换网络环境后重试",
			Type:          runtime.InfoDialog,
			DefaultButton: "Ok",
		})
		return "", errors.New("network error")
	}
	defer res.Body.Close()
	bufferSize := 128 * 1024
	totalSize := res.ContentLength
	downloadedSize := int64(0)

	reader := bufio.NewReaderSize(res.Body, bufferSize)
	file, err := os.Create(dest + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	buf := make([]byte, bufferSize)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			writer.Write(buf[:n])
			downloadedSize += int64(n)
			progress := float64(downloadedSize) / float64(totalSize) * 100
			roundedProgress := roundToTwoDecimals(progress)
			runtime.EventsEmit(ctx, events, roundedProgress)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	writer.Flush()
	return fileName, nil
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func InitConfig(ctx context.Context) bool {
	var defaultFile = util.HomeDir() + "/slack/"
	// os.MkdirAll(defaultFile, 0644)
	const latestConfigVersion = "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version"
	resp, err := clients.SimpleGet(latestConfigVersion, clients.NewRestyClient(nil, true))
	if err != nil {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Title:         "提示",
			Message:       "检查配置文件版本失败",
			Type:          runtime.InfoDialog,
			DefaultButton: "Ok",
		})
		return false
	}
	if resp.StatusCode() == 468 {
		runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Title:         "提示",
			Message:       "下载配置文件失败，请切换网络环境后重试",
			Type:          runtime.InfoDialog,
			DefaultButton: "Ok",
		})
		return false
	}
	configFileZip := fmt.Sprintf("%sv%s/config.zip", LastestPocUrl, string(resp.Body()))
	fileName, err := download(configFileZip, defaultFile)
	if err != nil {
		gologger.Error(ctx, err)
		return false
	}
	uz := fileutil.NewUnzip()
	if _, err := uz.Extract(defaultFile+fileName, defaultFile); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	os.Remove(util.HomeDir() + "/slack/config.zip")
	return true
}
