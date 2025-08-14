package tianyancha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slack-wails/lib/gologger"
	"slack-wails/lib/gomessage"
	"slack-wails/lib/structs"
	"sync"
	"time"

	"github.com/qiwentaidi/clients"

	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TycClient struct {
	ctx     context.Context
	Headers map[string]string
	KeyMap  map[string]CompanySuggest
	mutex   sync.RWMutex
}

var HumanCheckChan = make(chan struct{}) // 放在全局或单例结构里
func NewClient(ctx context.Context, token, id string) *TycClient {
	headers := map[string]string{
		"Referer":      " https://www.tianyancha.com/",
		"Version":      "TYC-Web",
		"X-Auth-Token": token,
		"x-tycid":      id,
	}
	return &TycClient{
		ctx:     ctx,
		Headers: headers,
		KeyMap:  make(map[string]CompanySuggest),
	}
}

func (c *TycClient) WithContentTypeHeaders() map[string]string {
	headers := make(map[string]string)
	maps.Copy(headers, c.Headers)
	headers["Content-Type"] = "application/json"
	return headers
}

func (c *TycClient) CheckLogin() bool {
	loginURL := "https://capi.tianyancha.com/cloud-app-management/service/layout/vipStatus"
	resp, err := clients.DoRequest("GET", loginURL, c.Headers, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return false
	}
	return !strings.Contains(string(resp.Body()), "mustlogin")
}

func (c *TycClient) CheckKeyMap(companyName string) (CompanySuggest, error) {
	if _, ok := c.KeyMap[companyName]; !ok {
		suggest, err := c.fetchCompanySuggest(companyName) // 获得到一个模糊匹配后，关联度最高的名称
		if err != nil {
			return CompanySuggest{}, err
		}
		if companyName != suggest.ComName { // 如果传进来的名称与模糊匹配的不相同
			var isFuzz = fmt.Sprintf("天眼查模糊匹配名称为%v ——> %v,已替换原有名称进行查.", companyName, suggest.ComName)
			gologger.Info(c.ctx, isFuzz)
		}
		c.mutex.Lock()
		c.KeyMap[companyName] = *suggest
		c.mutex.Unlock()
	}
	return c.KeyMap[companyName], nil
}

type FuzzCompanyName struct {
	State      string `json:"state"`
	Message    string `json:"message"`
	Special    string `json:"special"`
	VipMessage string `json:"vipMessage"`
	IsLogin    int    `json:"isLogin"`
	ErrorCode  int    `json:"errorCode"`
	Data       struct {
		Keyword            string           `json:"keyword"`
		QuerySuggestList   []interface{}    `json:"querySuggestList"`
		CompanySuggestList []CompanySuggest `json:"companySuggestList"`
		OtherSuggestList   []struct {
			Keyword       interface{} `json:"keyword"`
			Alias         string      `json:"alias"`
			CompanyID     int64       `json:"companyId"`
			EntityID      string      `json:"entityId"`
			Type          int         `json:"type"`
			EntityName    string      `json:"entityName"`
			Logo          string      `json:"logo"`
			PromptContent string      `json:"promptContent"`
		} `json:"otherSuggestList"`
	} `json:"data"`
}

type CompanySuggest struct {
	ID            int64       `json:"id"`
	GraphID       string      `json:"graphId"`
	Type          int         `json:"type"`
	MatchType     string      `json:"matchType"`
	ComName       string      `json:"comName"`
	Name          string      `json:"name"`
	Alias         string      `json:"alias"`
	Logo          string      `json:"logo"`
	RegStatus     int         `json:"regStatus"`
	TaxCode       string      `json:"taxCode"`
	PromptContent interface{} `json:"promptContent"`
	Label         interface{} `json:"label"`
	SourceName    interface{} `json:"sourceName"`
	SourceURL     interface{} `json:"sourceUrl"`
}

// 比如输入公司简称，需要返回查询权重第一的公司基本信息
func (c *TycClient) fetchCompanySuggest(company string) (*CompanySuggest, error) {
	searchData := map[string]string{
		"keyword": company,
	}
	bytesData, _ := json.Marshal(searchData)
	resp, err := clients.DoRequest("POST", fmt.Sprintf("https://capi.tianyancha.com/cloud-tempest/search/suggest/company/main?_=%d", time.Now().UnixMicro()), c.WithContentTypeHeaders(), bytes.NewReader(bytesData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("[tianyancha] fuzz company: %s request interface err: %v", company, err)
	}

	var result FuzzCompanyName
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		runtime.EventsEmit(c.ctx, "tyc-human-check", "天眼查出现人机校验，请手动处理")
		gologger.DualLog(c.ctx, gologger.Level_DEBUG, "天眼查出现人机校验，请手动处理")
		<-HumanCheckChan // 挂起直到用户在前端点击“我已验证”
		gologger.DualLog(c.ctx, gologger.Level_DEBUG, "收到用户确认，继续查询")
		return c.fetchCompanySuggest(company) // 重新尝试
	}
	// 可能是没登录所以查询不到数据
	if len(result.Data.CompanySuggestList) == 0 {
		return nil, fmt.Errorf("[tianyancha] fuzz company %s not found", company)
	}
	return &result.Data.CompanySuggestList[0], nil
}

type SubsidiaryInfo struct {
	State      string `json:"state"`
	Message    string `json:"message"`
	Special    string `json:"special"`
	VipMessage string `json:"vipMessage"`
	IsLogin    int    `json:"isLogin"`
	ErrorCode  int    `json:"errorCode"`
	Data       struct {
		Result      []SubsidiaryApiData `json:"result"`
		SortField   interface{}         `json:"sortField"`
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

type SubsidiaryApiData struct {
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
}

// 返回子公司的信息
func (c *TycClient) FetchSubsidiary(companyName, companyId string, ratio int) ([]SubsidiaryApiData, error) {
	gologger.Info(c.ctx, fmt.Sprintf("[tianyancha] 正在查询%s的子公司信息 ", companyName))
	var assets []SubsidiaryApiData
	data := map[string]interface{}{
		"gid":          companyId,
		"pageSize":     100,
		"pageNum":      1,
		"province":     "-100",
		"percentLevel": "-100",
		"category":     "-100",
	}
	bytesData, _ := json.Marshal(data)

	resp, err := clients.DoRequest("POST", "https://capi.tianyancha.com/cloud-company-background/company/investListV2", c.WithContentTypeHeaders(), bytes.NewReader(bytesData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("[tianyancha] company: %s request subsidary error: %v", companyName, err)
	}

	var result SubsidiaryInfo
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("[tianyancha] company: %s json parse error: %v", companyName, err)
	}
	if result.State == "error" {
		gomessage.Warning(c.ctx, result.Message)
		return nil, fmt.Errorf("[tianyancha] company: %s request failed: %v", companyName, result.Message)
	}

	if len(result.Data.Result) == 0 {
		return nil, fmt.Errorf("[tianyancha] company: %s no subsidiary", companyName)
	}

	for _, item := range result.Data.Result {
		gq, _ := strconv.Atoi(strings.TrimSuffix(item.Percent, "%"))
		if gq < ratio || gq > 100 {
			continue
		}

		assets = append(assets, item)
	}
	return assets, nil
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

// 获取微信公众号信息
func (c *TycClient) FetchWeChatOfficialAccounts(companyName, companyId string) ([]structs.OfficialAccount, error) {
	gologger.Info(c.ctx, fmt.Sprintf("[tianyancha] 正在查询%s公众号信息 ", companyName))
	var result []structs.OfficialAccount
	resp, err := clients.DoRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize=50&pageNum=1", c.Headers, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("%s: 请求公众号接口出现错误 %v", companyName, err)
	}
	var oa OfficialAccounts
	json.Unmarshal(resp.Body(), &oa)
	if oa.ErrorCode != 0 {
		return nil, fmt.Errorf("%s: 获取公众号出现错误 %v", companyName, oa.Message)
	}
	if oa.Data.Count == 0 {
		return nil, fmt.Errorf(companyName + ": 不存在公众号资产不存在")
	}
	for _, item := range oa.Data.ResultList {
		result = append(result, structs.OfficialAccount{
			Numbers:      item.PublicNum,
			Name:         item.Title,
			Qrcode:       item.CodeImg,
			Introduction: item.Recommend,
			Logo:         item.TitleImgURL,
		})
	}
	return result, nil
}

func (c *TycClient) GetRegStatus(status int) string {
	switch status {
	case 0:
		return "存续"
	case 1:
		return "注销"
	case 2:
		return "吊销"
	default:
		return "其他"
	}
}
