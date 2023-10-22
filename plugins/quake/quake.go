package quake

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slack/common"
	"slack/common/logger"
	"slack/common/proxy"
	"slack/gui/global"
	"time"

	"fyne.io/fyne/v2/dialog"
)

type QuakeResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Port    int `json:"port"`
		Service struct {
			Name string `json:"name"`
			HTTP struct {
				Server string `json:"server"`
				Host   string `json:"host"`
				Title  string `json:"title"`
			} `json:"http"`
		} `json:"service"`
		IP       string `json:"ip"`
		Location struct {
			CityCn string `json:"city_cn"`
		} `json:"location"`
	} `json:"data"`
	Meta struct {
		Pagination struct {
			Count     int `json:"count"`
			PageIndex int `json:"page_index"`
			PageSize  int `json:"page_size"`
			Total     int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

var (
	PageSize, StartIndex, FinalTotal int
	QuakeUsage, Latest               string
)

func QuakeApiSearch(query string, startIndex, pageSize int, result *[][]string) {
	if common.Profile.Quake.Api == "" {
		dialog.ShowError(errors.New("请配置quake API"), global.Win)
	}
	var temp string
	StartIndex = startIndex
	data := make(map[string]interface{})
	data["query"] = query
	data["start"] = startIndex
	data["size"] = pageSize
	data["ignore_cache"] = false
	data["include"] = []string{"ip", "port", "service.http.title", "service.name", "service.http.host", "service.http.server", "location.city_cn"}
	data["latest"] = Latest
	bytesData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", "https://quake.360.net/api/v3/search/quake_service", bytes.NewReader(bytesData))
	if err != nil {
		logger.Debug(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-QuakeToken", common.Profile.Quake.Api)
	resp, err := proxy.DefaultClient().Do(req)
	if err != nil {
		logger.Debug(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug(err)
	}
	var qk QuakeResult
	json.Unmarshal([]byte(string(b)), &qk)
	if qk.Code == 0 {
		for i, data := range qk.Data {
			if data.Service.Name == "http/ssl" {
				temp = "https://" + data.IP + ":" + fmt.Sprint(data.Port)
			} else if data.Service.Name == "http" {
				temp = "http://" + data.IP + ":" + fmt.Sprint(data.Port)
			} else {
				temp = ""
			}
			*result = append(*result, []string{fmt.Sprint(startIndex + i + 1), temp, data.IP, fmt.Sprint(data.Port), data.Service.Name,
				data.Service.HTTP.Title, data.Service.HTTP.Host, data.Service.HTTP.Server, data.Location.CityCn})
		}
		FinalTotal = qk.Meta.Pagination.Total
		QuakeUsage = fmt.Sprintf("查询状态:%v,共查询到资产数量:%v", qk.Message, qk.Meta.Pagination.Total)
	} else {
		dialog.ShowError(fmt.Errorf(qk.Message), global.Win)
	}
	time.Sleep(time.Second)
}
