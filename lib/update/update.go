package update

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"time"

	update "github.com/fynelabs/selfupdate"
)

const (
	lastestPocUrl    = "https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/"
	lastestClinetUrl = "https://gitee.com/the-temperature-is-too-low/Slack/releases/download/"
)

// https://gitee.com/the-temperature-is-too-low/slack-poc/releases/download/v0.0.2/afrog-pocs.zip
func UpdatePoc() error {
	LocalConfig, _ := os.UserHomeDir()
	if err := os.RemoveAll(LocalConfig + "/afrog-pocs"); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	if !InitConfig() {
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

const (
	BinaryFile_Windows      = "slack-wails.exe"
	BinaryFile_Linux        = "slack-wails_linux_amd64"
	BinaryFile_Darwin_AMD64 = "slack-wails_darwin_amd64.app"
	BinaryFile_Darwin_ARM64 = "slack-wails_darwin_arm64.app"
)

func UpdateClinet(latestVersion string) error {
	var binaryFileName string
	switch runtime.GOOS {
	case "windows":
		binaryFileName = BinaryFile_Windows
	case "linux":
		binaryFileName = BinaryFile_Linux
	case "darwin":
		if runtime.GOARCH == "arm64" {
			binaryFileName = BinaryFile_Darwin_ARM64
		} else {
			binaryFileName = BinaryFile_Darwin_AMD64
		}
	}
	temp := lastestClinetUrl + latestVersion + "/" + binaryFileName
	if err := doUpdate(temp); err != nil {
		return err
	}
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

func InitConfig() bool {
	const latestConfigVersion = "https://gitee.com/the-temperature-is-too-low/slack-poc/raw/master/version"
	LocalConfig, _ := os.UserHomeDir()
	_, b, err := clients.NewRequest("GET", latestConfigVersion, nil, nil, 10, http.DefaultClient)
	if err != nil {
		return false
	}
	configFile := lastestPocUrl + "v" + string(b) + "/slack.zip"
	fileName, err := download(configFile, LocalConfig+"/")
	if err != nil {
		return false
	}
	uz := util.NewUnzip()
	if _, err := uz.Extract(LocalConfig+"/"+fileName, LocalConfig); err != nil {
		return false
	}
	return true
}
