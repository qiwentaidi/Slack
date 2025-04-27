// 用于控制任务的启停模块
package control

import (
	"context"
	"sync"
)

type ControlType string

const (
	Webscan   ControlType = "webscan"
	Portscan  ControlType = "portscan"
	Dirseach  ControlType = "dirsearch"
	Subdomain ControlType = "subdomain"
	Crack     ControlType = "crack"
)

var (
	ctxMap     = make(map[ControlType]context.CancelFunc)
	ctxMapLock sync.RWMutex
)

func RegisterScanContext(scanType ControlType, cancel context.CancelFunc) {
	ctxMapLock.Lock()
	defer ctxMapLock.Unlock()
	ctxMap[scanType] = cancel
}

func CancelScanContext(scanType ControlType) {
	ctxMapLock.Lock()
	defer ctxMapLock.Unlock()
	if cancel, ok := ctxMap[scanType]; ok {
		cancel()
		delete(ctxMap, scanType)
	}
}

func GetScanContext(scanType ControlType) (context.Context, context.CancelFunc) {
	ctxMapLock.Lock()
	defer ctxMapLock.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	ctxMap[scanType] = cancel
	return ctx, cancel
}
