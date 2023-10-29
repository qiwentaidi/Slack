package common

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"slack/common/logger"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	defaultConfigFile      = "./config/config.json"
	defaultConfigDirectory = "./config"
	DefaultWebTimeout      = 10
)

type Profiles struct {
	WebScan struct {
		Thread int
	}
	Subdomain struct {
		DNS1 string
		DNS2 string
	}
	PortScan struct {
		Thread  int
		Timeout int
	}
	Proxy struct {
		Enable   bool
		Mode     string
		Address  string
		Port     int
		Username string
		Password string
	}
	Hunter struct {
		Api string
	}
	Fofa struct {
		Email string
		Api   string
	}
	Quake struct {
		Api string
	}
}

var Profile Profiles

func init() {
	go GC()
	if _, err := os.Stat(defaultConfigFile); err != nil {
		profile := Profiles{
			WebScan: struct {
				Thread int
			}{
				100,
			},
			Subdomain: struct {
				DNS1 string
				DNS2 string
			}{
				"223.5.5.5", "223.6.6.6",
			},
			PortScan: struct {
				Thread  int
				Timeout int
			}{
				2000, 10,
			},
			Proxy: struct {
				Enable   bool
				Mode     string
				Address  string
				Port     int
				Username string
				Password string
			}{
				false, "HTTP", "127.0.0.1", 8080, "", "",
			},
			Hunter: struct{ Api string }{
				Api: "",
			},
			Fofa: struct {
				Email string
				Api   string
			}{
				"", "",
			},
			Quake: struct{ Api string }{
				Api: "",
			},
		}
		b, err := json.MarshalIndent(profile, "", "    ")
		if err != nil {
			logger.Error(err)
		}
		if _, err := os.Stat(defaultConfigDirectory); err != nil {
			os.Mkdir(defaultConfigDirectory, 0777)
		}
		f, err := os.Create(defaultConfigFile)
		if err != nil {
			logger.Error(err)
			return
		}
		defer f.Close()
		_, err = f.WriteString(string(b))
		if err != nil {
			logger.Error(err)
			return
		}
	}
	file, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		logger.Error(err)
		return
	}
	// Unmarshal JSON content into struct
	err = json.Unmarshal(file, &Profile)
	if err != nil {
		logger.Error(err)
		return
	}
}

// 10s回收一次内存
func GC() {
	for {
		runtime.GC()
		debug.FreeOSMemory()
		time.Sleep(10 * time.Second)
	}
}

// 如果.old文件存在则删除
func init() {
	const oldPocZip = "./config/afrog-pocs.zip"
	currentMain := strings.Split(os.Args[0], "\\")
	dir, _ := os.Getwd()
	oldFile := fmt.Sprintf("%v\\.%v.old", dir, currentMain[len(currentMain)-1:][0])
	if _, err := os.Stat(oldFile); err == nil {
		if err2 := os.Remove(oldFile); err2 != nil {
			logger.Error(err)
		}
	}
	if _, err := os.Stat(oldPocZip); err == nil {
		if err2 := os.Remove(oldFile); err2 != nil {
			logger.Error(err)
		}
	}
}

func init() {
	// 初始化元数据
	app.SetMetadata(fyne.AppMetadata{
		Name:    "slack",
		Version: "1.4.5",
		Custom: map[string]string{
			"Issues": "https://github.com/qiwentaidi/Slack/issues/new",
		},
	})
}
