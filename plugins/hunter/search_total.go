package hunter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slack/common"
	"slack/common/logger"
	"slack/gui/custom"
	"slack/lib/util"
	"strings"
	"time"
)

var Restquota string

func SearchTotal(search string) {
	current_time := time.Now()
	before_time := current_time.AddDate(0, -1, 0)
	addr := "https://hunter.qianxin.com/openApi/search?api-key=" + common.Profile.Hunter.Api + "&search=" + HunterBaseEncode(search) +
		"&page=1&page_size=1&is_web=3&port_filter=false&start_time=" + before_time.Format("2006-01-02") + "&end_time=" + current_time.Format("2006-01-02")
	r, err := http.Get(addr)
	if err != nil {
		logger.Debug(err)
	}
	b, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		logger.Debug(err1)
	}
	defer r.Body.Close()
	var hr HunterResult
	json.Unmarshal([]byte(string(b)), &hr)
	if hr.Code != 200 {
		custom.Console.Append(fmt.Sprintf("[INF] %v\n", hr))
	} else {
		common.HunterAsset = append(common.HunterAsset, []string{strings.Split(search, "\"")[1], fmt.Sprintf("%v", hr.Data.Total)})
	}
	Restquota = hr.Data.RestQuota
}

func SeachICP(company string) {
	time.Sleep(time.Millisecond * 2000)
	str := fmt.Sprintf("icp.name=\"%v\"", company)
	SearchTotal(str)
}

func SeachDomain(domain string) {
	var str string
	time.Sleep(time.Millisecond * 2000)
	// 处理网站域名是IP的情况
	if util.RegIP.MatchString(domain) {
		str = fmt.Sprintf("ip=\"%v\"", domain)
	} else {
		str = fmt.Sprintf("domain.suffix=\"%v\"", domain)
	}
	SearchTotal(str)
}
