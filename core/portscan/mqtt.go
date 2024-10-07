package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func MqttScan(ctx context.Context, host string, usernames, passwords []string) {
	flag, err := MqttUnauth(host)
	if flag && err == nil {
		runtime.EventsEmit(ctx, "bruteResult", Burte{
			Status:   true,
			Host:     host,
			Protocol: "mqtt",
			Username: "unauthorized",
			Password: "",
		})
		gologger.Success(ctx, fmt.Sprintf("mqtt://%s is unauthorized access", host))
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("mqtt://%s is no unauthorized access", host))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitBruteFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MqttConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "bruteResult", Burte{
					Status:   true,
					Host:     host,
					Protocol: "mqtt",
					Username: user,
					Password: pass,
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
