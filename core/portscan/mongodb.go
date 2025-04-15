package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongodbScan(ctx context.Context, host string, usernames, passwords []string) {
	flag, err := MongodbConn(host, "", "")
	if flag && err == nil {
		runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
			ID:       "mongodb unauthorized",
			Name:     "mongodb unauthorized",
			URL:      host,
			Type:     "Mongodb",
			Severity: "HIGH",
		})
		return
	} else {
		gologger.Info(ctx, fmt.Sprintf("mongodb://%s is no unauthorized access", host))
	}
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := MongodbConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "mongodb weak password",
					Name:     "mongodb weak password",
					URL:      host,
					Type:     "mongodb",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("mongodb://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

// For higher versions of MongoDB, this function cannot be authenticated and will prompt to upgrade the driver. Abandoned on 1.6.6
//
// func MongodbUnauth(host string) (flag bool, err error) {
// 	flag = false
// 	senddata := []byte{72, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 212, 7, 0, 0, 0, 0, 0, 0, 97, 100, 109, 105, 110, 46, 36, 99, 109, 100, 0, 0, 0, 0, 0, 1, 0, 0, 0, 33, 0, 0, 0, 2, 103, 101, 116, 76, 111, 103, 0, 16, 0, 0, 0, 115, 116, 97, 114, 116, 117, 112, 87, 97, 114, 110, 105, 110, 103, 115, 0, 0}
// 	conn, err := WrapperTcpWithTimeout("tcp", host, 10*time.Second)
// 	defer func() {
// 		if conn != nil {
// 			conn.Close()
// 		}
// 	}()
// 	if err != nil {
// 		return flag, err
// 	}
// 	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
// 	if err != nil {
// 		return flag, err
// 	}
// 	_, err = conn.Write(senddata)
// 	if err != nil {
// 		return flag, err
// 	}
// 	buf := make([]byte, 1024)
// 	count, err := conn.Read(buf)
// 	if err != nil {
// 		return flag, err
// 	}
// 	text := string(buf[0:count])
// 	if strings.Contains(text, "totalLinesWritten") {
// 		flag = true
// 	}
// 	return flag, err
// }

func MongodbConn(host, user, pass string) (flag bool, err error) {
	flag = false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s", host)

	// Define client options with or without credentials
	clientOpts := options.Client().ApplyURI(mongoURI)
	if user != "" && pass != "" {
		credentials := options.Credential{
			Username: user,
			Password: pass,
		}
		clientOpts.SetAuth(credentials)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return flag, err
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, readpref.Primary())
	if err == nil {
		flag = true
	}
	return flag, err
}
