package gologger

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageLevel int

type Notification struct {
	Level MessageLevel
	Type  string // Notification || Message
	Msg   string
}

const (
	MessageInfo = iota
	MessageSuccess
	MessageWarning
	MessageError
)

func GoMessage(ctx context.Context, level MessageLevel, logtype string, i interface{}) {
	runtime.EventsEmit(ctx, "gomessage", &Notification{
		Level: level,
		Type:  logtype,
		Msg:   Msg(i),
	})
}
