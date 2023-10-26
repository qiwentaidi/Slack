package portscan

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slack/common"
	"slack/gui/custom"
	"slack/gui/global"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
)

type Brutelist struct {
	user string
	pass string
}

func RdpScan(host, domain string, associate bool, usertext, passtext *widget.Entry) {
	var wg sync.WaitGroup
	var signal bool
	var counter int64
	var mutex sync.Mutex
	var passwords []string
	custom.LogTime = time.Now().Unix()
	if associate {
		passwords = common.ParseTarget(global.ThinkDict.Text, common.Mode_Other)
	} else {
		passwords = common.ParseDict(passtext, common.Passwords)
	}
	usernames := common.ParseDict(usertext, common.Userdict["rdp"])
	all := len(passwords) * len(usernames)
	brlist := make(chan Brutelist, all)
	for _, user := range usernames {
		for _, pass := range passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			brlist <- Brutelist{user, pass}
		}
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(host, domain, &wg, brlist, &signal, counter, all, &mutex, common.Profile.PortScan.Timeout)
	}
	close(brlist)
	go func() {
		wg.Wait()
		signal = true
	}()
	for !signal {
		continue
	}
}

func worker(host, domain string, wg *sync.WaitGroup, brlist chan Brutelist, signal *bool, num int64, all int, mutex *sync.Mutex, timeout int) {
	defer wg.Done()
	for one := range brlist {
		if *signal {
			return
		}
		atomic.AddInt64(&num, 1)
		user, pass := one.user, one.pass
		flag, err := RdpConn(host, domain, user, pass, timeout)
		if flag && err == nil {
			if domain != "" {
				common.PortBurstResult = append(common.PortBurstResult, []string{"RDP", host + ":" + domain, user, pass, ""})
				custom.Console.Append(fmt.Sprintf("[+] rdp://%v:%v\\%v %v\n", host, domain, user, pass))
			} else {
				common.PortBurstResult = append(common.PortBurstResult, []string{"RDP", host, user, pass, ""})
				custom.Console.Append(fmt.Sprintf("[+] rdp://%v:%v %v\n", host, user, pass))
			}
			*signal = true
			return
		} else {
			custom.LogProgress(num, all, fmt.Sprintf("[-] (%v/%v) rdp %v %v %v %v", num, all, host, user, pass, err))
		}
	}
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
