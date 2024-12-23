package gomessage

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	Level_INFO    = "info"
	Level_WARN    = "warning"
	Level_ERROR   = "error"
	Level_Success = "success"
)

type MsgInfo struct {
	Level string
	Msg   string
}

func Info(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gomessage", &MsgInfo{
		Level: Level_INFO,
		Msg:   Msg(i),
	})
}

func Warning(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gomessage", &MsgInfo{
		Level: Level_WARN,
		Msg:   Msg(i),
	})
}

func Error(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gomessage", &MsgInfo{
		Level: Level_ERROR,
		Msg:   Msg(i),
	})
}

func Success(ctx context.Context, i interface{}) {
	runtime.EventsEmit(ctx, "gomessage", &MsgInfo{
		Level: Level_Success,
		Msg:   Msg(i),
	})
}

func Msg(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
