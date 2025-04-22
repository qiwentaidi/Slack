// 记录应用运行日志
package syslogger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"slack-wails/lib/util"
)

const (
	Level_INFO    = "[INF]"
	Level_WARN    = "[WRN]"
	Level_ERROR   = "[ERR]"
	Level_DEBUG   = "[DEB]"
	Level_Success = "[SUC]"
)

var sysLogPath = filepath.Join(util.HomeDir(), "syslog")
var logLock sync.Mutex

type MsgInfo struct {
	Level string
	Msg   string
}

func Info(ctx context.Context, i interface{}) {
	msg := Msg(i)
	writeToFile(Level_INFO, msg)
}

func Warning(ctx context.Context, i interface{}) {
	msg := Msg(i)
	writeToFile(Level_WARN, msg)
}

func Error(ctx context.Context, i interface{}) {
	msg := Msg(i)
	writeToFile(Level_ERROR, msg)
}

func Debug(ctx context.Context, i interface{}) {
	msg := Msg(i)
	writeToFile(Level_DEBUG, msg)
}

func Success(ctx context.Context, i interface{}) {
	msg := Msg(i)
	writeToFile(Level_Success, msg)
}

func Msg(i interface{}) string {
	return fmt.Sprintf("[%s] %v", currentTime(), i)
}

func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func currentLogFileName() string {
	return time.Now().Format("2006-01-02") + ".log"
}

func writeToFile(level string, msg string) {
	logLock.Lock()
	defer logLock.Unlock()

	// 确保日志目录存在
	if err := os.MkdirAll(sysLogPath, 0755); err != nil {
		return
	}

	filePath := filepath.Join(sysLogPath, currentLogFileName())
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	line := fmt.Sprintf("%s %s\n", level, msg)
	_, _ = f.WriteString(line)
}
