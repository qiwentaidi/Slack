package hunter

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

// 存储Hunter数据的结构体
type HunterResult struct {
	Code int64 `json:"code"`
	Data struct {
		AccountType string `json:"account_type"`
		Arr         []struct {
			AsOrg        string `json:"as_org"`
			Banner       string `json:"banner"`
			BaseProtocol string `json:"base_protocol"`
			City         string `json:"city"`
			Company      string `json:"company"`
			Component    []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"component"`
			Country        string `json:"country"`
			Domain         string `json:"domain"`
			IP             string `json:"ip"`
			IsRisk         string `json:"is_risk"`
			IsRiskProtocol string `json:"is_risk_protocol"`
			IsWeb          string `json:"is_web"`
			Isp            string `json:"isp"`
			Number         string `json:"number"`
			Os             string `json:"os"`
			Port           int64  `json:"port"`
			Protocol       string `json:"protocol"`
			Province       string `json:"province"`
			StatusCode     int64  `json:"status_code"`
			UpdatedAt      string `json:"updated_at"`
			URL            string `json:"url"`
			WebTitle       string `json:"web_title"`
		} `json:"arr"`
		ConsumeQuota string `json:"consume_quota"`
		RestQuota    string `json:"rest_quota"`
		SyntaxPrompt string `json:"syntax_prompt"`
		Time         int64  `json:"time"`
		Total        int64  `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}

var (
	StartTime, SelectAssets, DeDuplication string // hunter查询的全局参数
	PageSize                               int
	HunterSurplus                          string // 剩余积分
	HSearchDataSize                        string // 查询数量
	HunterTotal                            int64  // 查询到的数量
	HunterTime                             int64  // 所消耗时间
)

func HunterApiSearch(query, pageSize, page string, data *[][]string) {
	var assembly []string
	address := "https://hunter.qianxin.com/openApi/search?api-key=" + common.Profile.Hunter.Api + "&search=" + HunterBaseEncode(query) + "&page=" +
		page + "&page_size=" + pageSize + "&is_web=" + SelectAssets + "&port_filter=" + DeDuplication + "&start_time=" + StartTime + "&end_time=" + time.Now().Format("2006-01-02")
	r, err := http.Get(address)
	if err != nil {
		logger.Info(err)
	}
	b, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var hr HunterResult
	json.Unmarshal([]byte(string(b)), &hr)
	if hr.Code != 200 {
		if hr.Code == 40205 {
			dialog.ShowError(errors.New(hr.Message), global.Win)
		} else {
			dialog.ShowError(errors.New(hr.Message), global.Win)
			return
		}
	}
	p, _ := strconv.Atoi(page)
	t, _ := strconv.Atoi(pageSize)
	if len(hr.Data.Arr) == 0 {
		dialog.ShowInformation("提示", "查询数据结果为空", global.Win)
	} else {
		for i := 0; i < len(hr.Data.Arr); i++ {
			for _, v := range hr.Data.Arr[i].Component {
				assembly = append(assembly, v.Name+v.Version)
			}
			*data = append(*data, []string{
				strconv.Itoa(t*(p-1) + i + 1), hr.Data.Arr[i].URL, hr.Data.Arr[i].IP, strconv.FormatInt(hr.Data.Arr[i].Port, 10) + "/" + hr.Data.Arr[i].Protocol,
				hr.Data.Arr[i].Domain, strings.Join(assembly, " | "), hr.Data.Arr[i].WebTitle, strconv.FormatInt(hr.Data.Arr[i].StatusCode, 10), hr.Data.Arr[i].Company,
				hr.Data.Arr[i].Country + "" + hr.Data.Arr[i].Province + "" + hr.Data.Arr[i].City, hr.Data.Arr[i].UpdatedAt,
			})
			assembly = []string{}
		}
		HunterSurplus = hr.Data.RestQuota
		HunterTotal = hr.Data.Total
		HunterTime = hr.Data.Time
	}
}

// hunter base64加密接口
func HunterBaseEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func AssetExport(search string, total int) {
	var addr string
	if total <= 0 {
		dialog.ShowError(errors.New("导出数量不能为空"), global.Win)
		return
	}
	file, _ := os.OpenFile(fmt.Sprintf("./reports/hunter_asset_%v.csv", time.Now().Format("20060102_150405")), os.O_CREATE|os.O_RDWR, os.ModePerm) // 创建结果文件
	file.WriteString("\xEF\xBB\xBF")
	assetFile := csv.NewWriter(file)
	for page := range util.SplitInt(total, 100) {
		addr = "https://hunter.qianxin.com/openApi/search?api-key=" + common.Profile.Hunter.Api + "&search=" + HunterBaseEncode(search) + "&page=" + fmt.Sprint(page+1) + "&page_size=100&is_web=" + SelectAssets + "&port_filter=" + DeDuplication + "&start_time=" + StartTime + "&end_time=" + time.Now().Format("2006-01-02")
		r, err := http.Get(addr)
		if err != nil {
			logger.Debug(err)
		}
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var hr HunterResult
		json.Unmarshal([]byte(string(b)), &hr)
		if hr.Code != 200 {
			if hr.Code == 40205 {
				dialog.ShowError(errors.New(hr.Message), global.Win)
			} else {
				dialog.ShowError(errors.New(hr.Message), global.Win)
				return
			}
		} else {
			for _, arr := range hr.Data.Arr {
				assetFile.Write([]string{arr.URL, arr.IP, fmt.Sprintf("%v", arr.Port), arr.Protocol, arr.WebTitle, arr.Domain, arr.Company})
				assetFile.Flush()
			}
		}
		time.Sleep(time.Millisecond * 2000)
	}
	dialog.ShowCustom("提示", "OK", widget.NewLabel(fmt.Sprintf("[*] 导出结束,共计扣除积分%v(重复导出不会扣除积分)", total)), global.Win)
}

func Import(search, size string) (temps []string) {
	total, _ := strconv.Atoi(size)
	for page, size := range util.SplitInt(total, 100) {
		address := "https://hunter.qianxin.com/openApi/search?api-key=" + common.Profile.Hunter.Api + "&search=" + HunterBaseEncode(search) +
			"&page=" + fmt.Sprint(page+1) + "&page_size=" + fmt.Sprint(size) + "&is_web=1&port_filter=false&start_time=" + time.Now().AddDate(0, -1, 0).Format("2006-01-02") + "&end_time=" + time.Now().Format("2006-01-02")
		r, err := http.Get(address)
		if err != nil {
			logger.Info(err)
			return
		}
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		var hr HunterResult
		json.Unmarshal([]byte(string(b)), &hr)
		if hr.Code != 200 {
			if hr.Code == 40205 {
				dialog.ShowError(errors.New(hr.Message), global.Win)
			} else {
				dialog.ShowError(errors.New(hr.Message), global.Win)
				return
			}
		}
		if len(hr.Data.Arr) == 0 {
			dialog.ShowInformation("提示", "查询数据结果为空", global.Win)
			return
		} else {
			for i := 0; i < len(hr.Data.Arr); i++ {
				temps = append(temps, hr.Data.Arr[i].URL)
			}
		}
		time.Sleep(time.Millisecond * 2000)
	}
	return temps
}
