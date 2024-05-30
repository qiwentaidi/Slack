package webscan

import (
	"fmt"
	"slack-wails/lib/util"
	"strings"
	"testing"
)

func TestCaller(t *testing.T) {
	home := util.HomeDir()
	var workflowFile = home + "/slack/config/workflow.yaml"
	var webfingerFile = home + "/slack/config/webfinger.yaml"
	var activefingerFile = home + "/slack/config/dir.yaml"
	InitAll(webfingerFile, activefingerFile, workflowFile)
	keys := []string{}
	keywords := "sap"
	nc := NewNucleiCaller("", "", true, "")
	if keywords != "" {
		keys = strings.Split(keywords, ",")
	}
	v := nc.CallerAP("http://erpdev.wingtech.com:50000/", keys)
	fmt.Printf("v: %v\n", v)
}
