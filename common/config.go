package common

import (
	"encoding/json"
	"os"
	"runtime"
	"runtime/debug"
	"slack/common/logger"
	"time"
)

const (
	Version                = "1.4.3"
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
