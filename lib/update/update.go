package update

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slack-wails/lib/util"
	"time"
)

const (
	LocalConfig   = "./config/"
	lastestPocUrl = "https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/"
)

// https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/v0.0.2/afrog-pocs.zip
func UpdatePoc(latestVersion string) error {
	temp := lastestPocUrl + latestVersion
	fmt.Printf("temp: %v\n", temp)
	workflow := temp + "/workflow.yaml"
	webfinger := temp + "/webfinger.yaml"
	pocs := temp + "/afrog-pocs.zip"
	if _, err := download(workflow, LocalConfig); err != nil {
		return errors.New("workflow.yaml更新失败")
	}
	if _, err := download(webfinger, LocalConfig); err != nil {
		return errors.New("webfinger.yaml更新失败")
	}
	fileName, err := download(pocs, LocalConfig)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(LocalConfig + "afrog-pocs"); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	uz := util.NewUnzip()
	if _, err := uz.Extract(LocalConfig+fileName, LocalConfig); err != nil {
		return fmt.Errorf("afrog-poc decompression failed. %v", err)
	}
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
