package update

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"time"
)

const (
	lastestPocUrl    = "https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/"
	lastestClinetUrl = "https://gitee.com/the-temperature-is-too-low/Slack/releases/download/"
)

// https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/v0.0.2/afrog-pocs.zip
func UpdatePoc(configPath string) error {
	if err := os.RemoveAll(configPath + "config"); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	if !InitConfig(configPath) {
		return errors.New("update poc failed")
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

func InitConfig(configPath string) bool {
	os.Mkdir(util.HomeDir()+"/slack", 0777)
	const latestConfigVersion = "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version"
	_, b, err := clients.NewRequest("GET", latestConfigVersion, nil, nil, 10, http.DefaultClient)
	if err != nil {
		return false
	}
	configFileZip := lastestPocUrl + "v" + string(b) + "/config.zip"
	fileName, err := download(configFileZip, configPath)
	if err != nil {
		return false
	}
	uz := util.NewUnzip()
	if _, err := uz.Extract(configPath+fileName, configPath); err != nil {
		return false
	}
	os.Remove(configPath + "/config.zip")
	return true
}
