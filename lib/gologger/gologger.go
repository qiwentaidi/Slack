// 主要用于记录用户在客户端可查看的日志内容
package gologger

import (
	"context"
	"fmt"
	"slack-wails/lib/syslogger"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Level_INFO    = "[INF]"
	Level_WARN    = "[WRN]"
	Level_ERROR   = "[ERR]"
	Level_DEBUG   = "[DEB]"
	Level_Success = "[SUC]"
)

type MsgInfo struct {
	Level string
	Msg   string
}

func Info(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_INFO,
		Msg:   Msg(i),
	})
}

func Warning(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_WARN,
		Msg:   Msg(i),
	})
}

func Error(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_ERROR,
		Msg:   Msg(i),
	})
}

func Debug(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_DEBUG,
		Msg:   Msg(i),
	})
}

func Success(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gologger", &MsgInfo{
		Level: Level_Success,
		Msg:   Msg(i),
	})
}

func Msg(i interface{}) string {
	return fmt.Sprintf("[%s] %v", currentTime(), i)
}

func currentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 当出现大量错误时，每个5s打印一次错误日志，偷懒式写法
func IntervalError(ctx context.Context, i interface{}) {
	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()
	if unixTimestamp%5 == 0 {
		Error(ctx, i)
	}
}

// 同时使用gologger和syslogger
func DualLog(ctx context.Context, level string, i interface{}) {
	switch level {
	case Level_INFO:
		Info(ctx, i)
		syslogger.Info(ctx, i)
	case Level_WARN:
		Warning(ctx, i)
		syslogger.Warning(ctx, i)
	case Level_ERROR:
		Error(ctx, i)
		syslogger.Error(ctx, i)
	case Level_DEBUG:
		Debug(ctx, i)
		syslogger.Debug(ctx, i)
	case Level_Success:
		Success(ctx, i)
		syslogger.Success(ctx, i)
	}
}
