package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"testing"
)

func TestNucleiCaller(t *testing.T) {
	home := util.HomeDir()
	var templateDir = home + "/slack/config/pocs"
	var webfingerFile = home + "/slack/config/webfinger.yaml"
	var activefingerFile = home + "/slack/config/dir.yaml"
	NewConfig().InitAll(context.TODO(), webfingerFile, activefingerFile, templateDir)
	nc := NewNucleiCaller("", false, "", clients.Proxy{})
	v := nc.CallerFP(context.TODO(), FingerPoc{
		URL:  "http://erpdev.wingtech.com:50000/",
		Tags: []string{"SAP-Web-App-Server"},
	})
	fmt.Printf("v: %v\n", v)
}
