package fofa

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slack/common"
	"slack/common/logger"
	"slack/gui/global"
	"slack/lib/util"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FofaResult struct {
	Error   bool       `json:"error"`
	Errmsg  string     `json:"errmsg"`
	Mode    string     `json:"mode"`
	Page    int64      `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int64      `json:"size"`
}

var (
	PageSize    int
	FofaTotal   int64 // 查询总量
	Fraud, Cert string
	temp        string
)

func FofaApiSearch(search, pageSize, pageNum string, data *[][]string) {
	address := "https://fofa.info/api/v1/search/all?email=" + common.Profile.Fofa.Email + "&key=" + common.Profile.Fofa.Api + "&qbase64=" +
		FOFABaseEncode(search) + "&cert.is_valid" + Cert + fmt.Sprintf("&is_fraud=%v&is_honeypot=%v", Fraud, Fraud) + "&page=" + pageNum + "&size=" + pageSize + "&fields=host,title,ip,domain,port,protocol,country_name,region,city,icp"
	r, err := http.Get(address)
	if err != nil {
		logger.Info(err)
	}
	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var fr FofaResult
	json.Unmarshal([]byte(string(b)), &fr)
	FofaTotal = fr.Size
	p, _ := strconv.Atoi(pageNum)
	t, _ := strconv.Atoi(pageSize)
	if fr.Error {
		dialog.ShowInformation("提示", fr.Errmsg, global.Win)
	} else {
		if fr.Size == 0 {
			dialog.ShowInformation("提示", "未查询到数据结果", global.Win)
		} else {
			for i := 0; i < len(fr.Results); i++ {
				if !strings.Contains(fr.Results[i][0], "://") && (fr.Results[i][5] == "http" || fr.Results[i][5] == "https") {
					temp = fr.Results[i][5] + "://" + fr.Results[i][0]
				} else {
					temp = fr.Results[i][0]
				}
				*data = append(*data, []string{
					strconv.Itoa(t*(p-1) + i + 1), temp, fr.Results[i][1], fr.Results[i][2],
					fr.Results[i][4], fr.Results[i][3], fr.Results[i][5], fr.Results[i][6] + " " + fr.Results[i][7] + " " + fr.Results[i][8], fr.Results[i][9],
				})
			}
		}
	}
	time.Sleep(time.Second)
}

// fofa base64加密接口
func FOFABaseEncode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func AssetExport(search string, total int) {
	os.Mkdir("reports", 0777)
	file, _ := os.OpenFile(fmt.Sprintf("./reports/fofa_asset_%v.csv", time.Now().Format("20060102_150405")), os.O_CREATE|os.O_RDWR, os.ModePerm) // 创建结果文件
	file.WriteString("\xEF\xBB\xBF")
	assetFile := csv.NewWriter(file)
	for page, size := range util.SplitInt(total, 1000) {
		addr := "https://fofa.info/api/v1/search/all?email=" + common.Profile.Fofa.Email + "&key=" + common.Profile.Fofa.Api + "&qbase64=" +
			FOFABaseEncode(search) + "&page=" + fmt.Sprint(page+1) + "&size=" + fmt.Sprint(size) + "&fields=host,title,ip,domain,port,protocol,country_name,region,city,icp"
		r, err := http.Get(addr)
		if err != nil {
			logger.Debug(err)
		}
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var fr FofaResult
		json.Unmarshal([]byte(string(b)), &fr)
		if fr.Error {
			dialog.ShowError(errors.New(fr.Errmsg), global.Win)
		} else {
			for i := range fr.Results {
				if !strings.Contains(fr.Results[i][0], "://") && (fr.Results[i][5] == "http" || fr.Results[i][5] == "https") {
					temp = fr.Results[i][5] + "://" + fr.Results[i][0]
				} else {
					temp = fr.Results[i][0]
				}
				assetFile.Write([]string{
					temp, fr.Results[i][1], fr.Results[i][2],
					fr.Results[i][4], fr.Results[i][3], fr.Results[i][5], fr.Results[i][6] + " " + fr.Results[i][7] + " " + fr.Results[i][8], fr.Results[i][9],
				})
				assetFile.Flush()
			}
		}
		time.Sleep(time.Millisecond * 2000)
	}
	dialog.ShowCustom("提示", "OK", widget.NewLabel("导出结束，文件名为"+file.Name()), global.Win)
}

// automation 自动化参数为true时需要额外获取ip与domain的对应关系
func Import(search, size string, automation bool) (urls []string, domainIP map[string][]string) {
	domainIP = make(map[string][]string)
	nums, _ := strconv.Atoi(size)
	var max int
	if nums%1000 == 0 {
		max = nums / 1000
	} else {
		max = nums/1000 + 1
	}
	if nums >= 1000 {
		size = "1000"
	}
	// max = 1-1000,循环一次 max 1001-2000,循环2次...
	for i := 1; i <= max; i++ {
		address := "https://fofa.info/api/v1/search/all?email=" + common.Profile.Fofa.Email + "&key=" + common.Profile.Fofa.Api + "&qbase64=" +
			FOFABaseEncode(search) + "&page=" + fmt.Sprintf("%v", i) + "&size=" + size + "&fields=host,port,protocol,ip"
		r, err := http.Get(address)
		if err != nil {
			logger.Debug(err)
			return urls, domainIP
		}
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var fr FofaResult
		json.Unmarshal([]byte(string(b)), &fr)
		FofaTotal = fr.Size
		if fr.Error {
			dialog.ShowInformation("提示", fr.Errmsg, global.Win)
			return urls, domainIP
		} else {
			if fr.Size == 0 {
				dialog.ShowInformation("提示", "未查询到数据结果", global.Win)
				return urls, domainIP
			} else {
				for i := 0; i < len(fr.Results); i++ {
					if !strings.Contains(fr.Results[i][0], "://") && (fr.Results[i][2] == "http" || fr.Results[i][2] == "https") {
						temp = fr.Results[i][2] + "://" + fr.Results[i][0]
					} else {
						temp = fr.Results[i][0]
					}
					urls = append(urls, temp)
					if automation {
						domain := fr.Results[i][0]
						if strings.Contains(domain, "://") {
							domain = strings.Split(domain, "://")[1]
						}
						if strings.Contains(domain, ":") {
							domain = strings.Split(domain, ":")[0]
						}
						domainIP[fr.Results[i][3]] = append(domainIP[fr.Results[i][3]], domain)
					}
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
	return urls, domainIP
}

func SearchTotal(search string) (total int64) {
	address := "https://fofa.info/api/v1/search/all?email=" + common.Profile.Fofa.Email + "&key=" + common.Profile.Fofa.Api + "&qbase64=" +
		FOFABaseEncode(search) + "&page=1&size=1&fields=host,port,protocol,ip"
	r, err := http.Get(address)
	if err != nil {
		logger.Debug(err)
		return 0
	}
	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var fr FofaResult
	json.Unmarshal([]byte(string(b)), &fr)
	if fr.Error {
		return 0
	} else {
		return fr.Size
	}
}
