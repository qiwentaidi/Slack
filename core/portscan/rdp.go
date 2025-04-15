package portscan

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"sync"
	"time"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func RdpScan(ctx context.Context, host string, usernames, passwords []string) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	limiter := make(chan bool, 10) // 限制协程数量
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitFunc {
				close(limiter)
				wg.Wait()
				return
			}
			pass = strings.Replace(pass, "{user}", string(user), -1)
			wg.Add(1)
			limiter <- true
			go func(user, pass string) {
				wg.Done()
				flag, err := RdpConn(host, "", user, pass, 10)
				mutex.Lock()
				if flag && err == nil {
					runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
						ID:       "rdp weak password",
						Name:     "rdp weak password",
						URL:      host,
						Type:     "RDP",
						Severity: "CRITICAL",
						Extract:  user + "/" + pass,
					})
					return
				} else {
					gologger.Info(ctx, fmt.Sprintf("rdp://%s %s:%s is login failed", host, user, pass))
				}
				mutex.Unlock()
				<-limiter
			}(user, pass)
		}
	}
	wg.Wait()
}

func RdpConn(host, domain, user, password string, timeout int) (bool, error) {
	g := NewClient(host, glog.NONE)
	err := g.Login(domain, user, password, timeout)

	if err == nil {
		return true, nil
	}

	return false, err
}

type Client struct {
	Host string // ip:port
	tpkt *tpkt.TPKT
	x224 *x224.X224
	mcs  *t125.MCSClient
	sec  *sec.Client
	pdu  *pdu.Client
	//vnc  *rfb.RFB
}

func NewClient(host string, logLevel glog.LEVEL) *Client {
	glog.SetLevel(logLevel)
	logger := log.New(os.Stdout, "", 0)
	glog.SetLogger(logger)
	return &Client{
		Host: host,
	}
}

func (g *Client) Login(domain, user, pwd string, timeout int) error {
	conn, err := WrapperTcpWithTimeout("tcp", g.Host, time.Duration(timeout)*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}
	glog.Info(conn.LocalAddr().String())
	g.tpkt = tpkt.New(core.NewSocketLayer(conn), nla.NewNTLMv2(domain, user, pwd))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)
	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)
	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.pdu.SetFastPathSender(g.tpkt)
	err = g.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	glog.Info("wait connect ok")
	wg := &sync.WaitGroup{}
	breakFlag := false
	wg.Add(1)

	g.pdu.On("error", func(e error) {
		err = e
		glog.Error("error", e)
		g.pdu.Emit("done")
	})
	g.pdu.On("close", func() {
		err = errors.New("close")
		glog.Info("on close")
		g.pdu.Emit("done")
	})
	g.pdu.On("success", func() {
		err = nil
		glog.Info("on success")
		g.pdu.Emit("done")
	})
	g.pdu.On("ready", func() {
		glog.Info("on ready")
		g.pdu.Emit("done")
	})
	g.pdu.On("update", func(rectangles []pdu.BitmapData) {
		glog.Info("on update:", rectangles)
	})
	g.pdu.On("done", func() {
		if !breakFlag {
			breakFlag = true
			wg.Done()
		}
	})
	wg.Wait()
	return err
}
