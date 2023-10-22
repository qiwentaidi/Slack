package logger

import (
	"fmt"
	"os"

	"time"
)

const defaultLogDirectory = "./log/"

// INFO记录不影响结果的日志
func Info(message interface{}) {
	log("[INF]", fmt.Sprintf("%v", message))
}

// ERROR记录会退出的日志
func Error(message interface{}) {
	log("[ERR]", fmt.Sprintf("%v", message))
}

// DEBUG记录会影响结果的日志
func Debug(message interface{}) {
	log("[DEB]", fmt.Sprintf("%v", message))
}

// log 写入日志
func log(level string, message string) {
	if _, err := os.Stat(defaultLogDirectory); err != nil {
		os.Mkdir(defaultLogDirectory, os.ModePerm)
	}
	// 创建日志文件，名称为当前时间
	f, _ := os.OpenFile(defaultLogDirectory+time.Now().Format("2006-01-02")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	// 写入日志
	f.WriteString(fmt.Sprintf("%v %v: %v\n", level, time.Now().Format("2006-01-02 15:04:05"), message))
}
