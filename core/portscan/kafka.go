// https://github.com/shadow1ng/fscan/blob/main/Plugins/Kafka.go
package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func KafkaScan(ctx, ctrlCtx context.Context, taskId, address string, usernames, passwords []string) {
	flag, err := KafkaConn(address, "", "")
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			TaskId:   taskId,
			ID:       "kafka unauthorized",
			Name:     "kafka unauthorized",
			URL:      address,
			Type:     "Kafka",
			Severity: "HIGH",
			Extract:  "",
		})
		gologger.Success(ctx, fmt.Sprintf("kafka://%s is unauthorized access", address))
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("kafka://%s is no unauthorized access", address))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ctrlCtx.Err() != nil {
				gologger.Warning(ctx, "[kafka] User exits crack scanning")
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := KafkaConn(address, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					TaskId:   taskId,
					ID:       "kafka weak password",
					Name:     "kafka weak password",
					URL:      address,
					Type:     "Kafka",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("kafka://%s %s:%s is login failed", address, user, pass))
			}
		}
	}
}

// KafkaConn 尝试 Kafka 连接
func KafkaConn(address, user, pass string) (bool, error) {
	timeout := time.Duration(10) * time.Second

	config := sarama.NewConfig()
	config.Net.DialTimeout = timeout
	config.Net.ReadTimeout = timeout
	config.Net.WriteTimeout = timeout
	config.Net.TLS.Enable = false
	config.Version = sarama.V2_0_0_0

	// 设置 SASL 配置
	if user != "" || pass != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = user
		config.Net.SASL.Password = pass
		config.Net.SASL.Handshake = true
	}

	brokers := []string{address}

	// 尝试作为消费者连接测试
	consumer, err := sarama.NewConsumer(brokers, config)
	if err == nil {
		defer consumer.Close()
		return true, nil
	}

	// 如果消费者连接失败，尝试作为客户端连接
	client, err := sarama.NewClient(brokers, config)
	if err == nil {
		defer client.Close()
		return true, nil
	}

	// 检查错误类型
	if strings.Contains(err.Error(), "SASL") ||
		strings.Contains(err.Error(), "authentication") ||
		strings.Contains(err.Error(), "credentials") {
		return false, fmt.Errorf("认证失败")
	}

	return false, err
}
