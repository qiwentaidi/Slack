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

type TycSearchID struct {
	State      string `json:"state"`
	Message    string `json:"message"`
	Special    string `json:"special"`
	VipMessage string `json:"vipMessage"`
	IsLogin    int    `json:"isLogin"`
	ErrorCode  int    `json:"errorCode"`
	Data       []struct {
		ID         int         `json:"id"`
		GraphID    string      `json:"graphId"`
		Type       int         `json:"type"`
		MatchType  string      `json:"matchType"` // 商标信息匹配 > 股票简称匹配 > 公司名称匹配 > 公司品牌匹配 > 公司信息匹配
		ComName    string      `json:"comName"`
		Name       string      `json:"name"`
		Alias      string      `json:"alias"`
		Logo       string      `json:"logo"`
		ClaimLevel interface{} `json:"claimLevel"`
		RegStatus  int         `json:"regStatus"`
	} `json:"data"`
}

type TycResult struct {
	State      string `json:"state"`
	Message    string `json:"message"`
	Special    string `json:"special"`
	VipMessage string `json:"vipMessage"`
	IsLogin    int    `json:"isLogin"`
	ErrorCode  int    `json:"errorCode"`
	Data       struct {
		Result []struct {
			Name             string      `json:"name"` // 公司名称
			PersonType       int         `json:"personType"`
			ServiceType      interface{} `json:"serviceType"`
			RegStatus        string      `json:"regStatus"`
			Percent          string      `json:"percent"` // 股权比例
			LegalPersonTitle string      `json:"legalPersonTitle"`
			LegalPersonName  string      `json:"legalPersonName"`
			Logo             interface{} `json:"logo"`
			Alias            string      `json:"alias"`
			ID               int64       `json:"id"` // 子公司的companyId
			Amount           string      `json:"amount"`
			EstiblishTime    int64       `json:"estiblishTime"`
			LegalPersonID    int         `json:"legalPersonId"`
			ServiceCount     interface{} `json:"serviceCount"`
			LegalAlias       interface{} `json:"legalAlias"`
			LegalLogo        interface{} `json:"legalLogo"`
			JigouName        interface{} `json:"jigouName"`
			JigouLogo        interface{} `json:"jigouLogo"`
			JigouID          interface{} `json:"jigouId"`
			ProductName      interface{} `json:"productName"`
			ProductLogo      interface{} `json:"productLogo"`
			ProductID        interface{} `json:"productId"`
		} `json:"result"`
		SortField   interface{} `json:"sortField"`
		PercentList []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"percentList"`
		ProvinceList []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"provinceList"`
		CategoryList []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"categoryList"`
		Total int `json:"total"`
	} `json:"data"`
}

var (
	gethead   = map[string]string{}
	posthead  = map[string]string{}
	TycKeyMap = make(map[string]structs.TycCompanyInfo)
)

func InitHEAD(token string) {
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
	_, b, err := clients.NewRequest("POST", "https://capi.tianyancha.com/cloud-tempest/search/suggest/v3", posthead, bytes.NewReader(bytesData), 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, err)
	}
	var qs TycSearchID
	if err = json.Unmarshal(b, &qs); err != nil {
		gologger.Error(ctx, err)
	}
	if len(qs.Data) > 0 { // 接口会自动进行 商标信息匹配 > 股票简称匹配 > 公司名称匹配 > 公司品牌匹配 > 公司信息匹配 五种规则的匹配
		company_id = qs.Data[0].GraphID
		company_name = qs.Data[0].ComName
	}
	time.Sleep(time.Second * 2)
	return company_id, company_name
}

type CompanyInfo struct {
	CompanyName string
	Holding     string
	Investment  string // 投资比例
	RegStatus   string
	Domains     []string
	CompanyId   string
}

// 返回查询公司的名称和子公司的名称, isSecond 是否为二次查询
func SearchSubsidiary(ctx context.Context, companyName, companyId string, ratio int, isSecond bool) (Asset []CompanyInfo) {
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
		gologger.Error(ctx, err)
		return
	}
	var qr TycResult
	json.Unmarshal(b, &qr)
	// 获取到本公司对应的域名，若是二次查询跳过
	if !isSecond {
		var domains []string
		domains, _ = Beianx(companyName)
		Asset = append(Asset, CompanyInfo{companyName, "本公司", "", qr.State, util.RemoveDuplicates(domains), companyId})
	}
	for _, result := range qr.Data.Result {
		gq, _ := strconv.Atoi(strings.TrimSuffix(result.Percent, "%"))
		if gq <= 100 && gq >= ratio { // 提取在控股范围内的子公司
			gologger.Info(ctx, fmt.Sprintf("%v", result.Name))
			var subsidiaryDomains []string
			if result.RegStatus == "存续" || result.RegStatus == "ok" { // 注销的公司不用查备案
				subsidiaryDomains, _ = Beianx(result.Name)
			}
			Asset = append(Asset, CompanyInfo{result.Name, result.Percent, result.Amount, result.RegStatus, util.RemoveDuplicates(subsidiaryDomains), fmt.Sprint(result.ID)})
		}
	}
	return
}

type OfficialAccounts struct {
	State      string `json:"state"`
	Message    string `json:"message"`
	Special    string `json:"special"`
	VipMessage string `json:"vipMessage"`
	IsLogin    int    `json:"isLogin"`
	ErrorCode  int    `json:"errorCode"`
	Data       struct {
		Count      int `json:"count"`
		ResultList []struct {
			PublicNum   string `json:"publicNum"`   // 微信号
			CodeImg     string `json:"codeImg"`     // 二维码
			Recommend   string `json:"recommend"`   // 简介
			Title       string `json:"title"`       // 名称
			TitleImgURL string `json:"titleImgURL"` // 公众号LOGO
		} `json:"resultList"`
	} `json:"data"`
}

type WechatReulst struct {
	CompanyName  string
	WechatName   string
	WechatNums   string
	Logo         string
	Qrcode       string
	Introduction string
}

// 获取微信公众号信息
func WeChatOfficialAccounts(ctx context.Context, companyName, companyId string) (wr []WechatReulst) {
	_, b, err := clients.NewRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize=1&pageNum=1", gethead, nil, 10, true, clients.DefaultClient())
	if err != nil {
		gologger.Error(ctx, err)
		return
	}
	var oa OfficialAccounts
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
		wr = append(wr, WechatReulst{
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
