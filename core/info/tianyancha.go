package info

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"sync"
	"time"

	"strconv"
	"strings"
)

var (
	gethead   = map[string]string{}
	posthead  = map[string]string{}
	TycKeyMap = make(map[string]structs.TycCompanyInfo)
)

func initTycHeader(token string) {
	gethead = map[string]string{
		"Version":      "TYC-Web",
		"X-Auth-Token": token,
	}
	posthead = map[string]string{
		"Version":      "TYC-Web",
		"X-Auth-Token": token,
		"Content-Type": "application/json",
	}
}

// 要根据ID值查子公司，只有本部公司需要走这个接口查询
func GetCompanyID(ctx context.Context, company string) (string, string) {
	var company_id, company_name string
	data := make(map[string]interface{})
	data["keyword"] = company
	bytesData, _ := json.Marshal(data)
	_, body, err := clients.NewRequest("POST", "https://capi.tianyancha.com/cloud-tempest/search/suggest/v3", posthead, bytes.NewReader(bytesData), 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, err)
	}
	var qs structs.TycSearchID
	if err = json.Unmarshal(body, &qs); err != nil {
		gologger.Error(ctx, fmt.Sprintf("[tianyancha] company %s 请求过快导致触发人机校验", company))
	}
	if len(qs.Data) > 0 { // 接口会自动进行 商标信息匹配 > 股票简称匹配 > 公司名称匹配 > 公司品牌匹配 > 公司信息匹配 五种规则的匹配
		company_id = qs.Data[0].GraphID
		company_name = qs.Data[0].ComName
	}
	time.Sleep(time.Second * 2)
	return company_id, company_name
}

// 返回查询公司的名称和子公司的名称, isSecond 是否为二次查询
func SearchSubsidiary(ctx context.Context, companyName, companyId string, ratio int, isSecond bool, searchDomain bool, machine string) (Asset []structs.CompanyInfo) {
	data := make(map[string]interface{})
	data["gid"] = companyId
	data["pageSize"] = 100
	data["pageNum"] = 1
	data["province"] = "-100"
	data["percentLevel"] = "-100"
	data["category"] = "-100"
	bytesData, _ := json.Marshal(data)
	_, b, err := clients.NewRequest("POST", "https://capi.tianyancha.com/cloud-company-background/company/investListV2", posthead, bytes.NewReader(bytesData), 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, fmt.Sprintf("[tianyancha] company: %s request interface err: %v", companyName, err))
		return
	}
	var qr structs.TycResult
	json.Unmarshal(b, &qr)
	// 获取到本公司对应的域名，若是二次查询跳过
	if !isSecond {
		var domains []string
		if searchDomain {
			domains, err = Beianx(companyName, machine)
			if err != nil {
				gologger.Debug(ctx, fmt.Sprintf("[beianx] company: %s request interface err: %v", companyName, err))
			}
		}
		Asset = append(Asset, structs.CompanyInfo{
			CompanyName: companyName,
			Holding:     "本公司",
			Investment:  "",
			RegStatus:   qr.State,
			Domains:     util.RemoveDuplicates(domains),
			CompanyId:   companyId,
		})
	}
	for _, result := range qr.Data.Result {
		gq, _ := strconv.Atoi(strings.TrimSuffix(result.Percent, "%"))
		if gq <= 100 && gq >= ratio { // 提取在控股范围内的子公司
			gologger.Info(ctx, fmt.Sprintf("%v", result.Name))
			var subsidiaryDomains []string
			if (result.RegStatus == "存续" || result.RegStatus == "ok") && searchDomain { // 注销的公司不用查备案
				subsidiaryDomains, err = Beianx(result.Name, machine)
				if err != nil {
					gologger.Debug(ctx, fmt.Sprintf("[beianx] company: %s request interface err: %v", result.Name, err))
				}
			}
			Asset = append(Asset, structs.CompanyInfo{
				CompanyName: result.Name,
				Holding:     result.Percent,
				Investment:  result.Amount,
				RegStatus:   result.RegStatus,
				Domains:     util.RemoveDuplicates(subsidiaryDomains),
				CompanyId:   fmt.Sprint(result.ID),
			})
		}
	}
	return
}

// 获取微信公众号信息
func WeChatOfficialAccounts(ctx context.Context, companyName, companyId string) (wr []structs.WechatReulst) {
	_, b, err := clients.NewRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize=1&pageNum=1", gethead, nil, 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, err)
		return
	}
	var oa structs.OfficialAccounts
	json.Unmarshal(b, &oa)
	if oa.ErrorCode != 0 || oa.Data.Count == 0 {
		gologger.Info(ctx, "公众号查询出现错误或不存在公众号资产,公司名称: "+companyName)
		return
	}
	_, b, err = clients.NewRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize="+fmt.Sprint(oa.Data.Count)+"&pageNum=1", gethead, nil, 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, err)
		return
	}
	json.Unmarshal(b, &oa)
	for _, result := range oa.Data.ResultList {
		wr = append(wr, structs.WechatReulst{
			CompanyName:  companyName,
			WechatNums:   result.PublicNum,
			WechatName:   result.Title,
			Qrcode:       result.CodeImg,
			Introduction: result.Recommend,
			Logo:         result.TitleImgURL,
		})
	}
	return
}

func CheckLogin(token string) bool {
	initTycHeader(token)
	u := "https://capi.tianyancha.com/cloud-monitor-provider/v4/monitor/checkMonitorTip.json"
	_, body, err := clients.NewRequest("GET", u, gethead, nil, 10, true, clients.DefaultClient())
	if err != nil {
		return false
	}
	if strings.Contains(string(body), "mustlogin") {
		return false
	}
	return true
}

var me sync.RWMutex

func CheckKeyMap(ctx context.Context, query string) structs.TycCompanyInfo {
	if _, ok := TycKeyMap[query]; !ok {
		companyId, fuzzName := GetCompanyID(ctx, query) // 获得到一个模糊匹配后，关联度最高的名称
		if query != fuzzName {                          // 如果传进来的名称与模糊匹配的不相同
			var isFuzz = fmt.Sprintf("天眼查模糊匹配名称为%v ——> %v,已替换原有名称进行查.", query, fuzzName)
			gologger.Info(ctx, isFuzz)
		}
		me.Lock()
		TycKeyMap[query] = structs.TycCompanyInfo{
			CompanyName: fuzzName,
			CompanyId:   companyId,
		}
		me.Unlock()
	}
	return TycKeyMap[query]
}
