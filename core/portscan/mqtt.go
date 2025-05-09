package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MqttScan(ctx, ctrlCtx context.Context, taskId, host string, usernames, passwords []string) {
	flag, err := MqttUnauth(host)
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			TaskId:   taskId,
			ID:       "mqtt unauthorized",
			Name:     "mqtt unauthorized",
			URL:      host,
			Type:     "MQTT",
			Severity: "HIGH",
		})
		gologger.Success(ctx, fmt.Sprintf("mqtt://%s is unauthorized access", host))
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("mqtt://%s is no unauthorized access", host))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[mqtt] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MqttConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "mqtt weak password",
					Name:     "mqtt weak password",
					URL:      host,
					Type:     "mqtt",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
			} else {
				gologger.Info(ctx, fmt.Sprintf("mqtt://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func MqttUnauth(host string) (bool, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", host))
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	return true, nil
}

func MqttConn(host, user, pass string) (bool, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", host)).SetUsername(user).SetPassword(pass)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return false, token.Error()
	}
	return true, nil
}
