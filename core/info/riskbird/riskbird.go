package riskbird

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
)

type RiskbirdClient struct {
	ctx     context.Context
	Headers map[string]string
}

func NewClient(ctx context.Context, cookie string) *RiskbirdClient {
	headers := map[string]string{
		"Cookie":       cookie,
		"Content-Type": "application/json",
		"App-device":   "WEB",
		"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
	}
	return &RiskbirdClient{
		ctx:     ctx,
		Headers: headers,
	}
}

func (r *RiskbirdClient) WithXSHeader() map[string]string {
	headers := make(map[string]string)
	maps.Copy(headers, r.Headers)
	headers["xs-content-type"] = "application/json"
	return headers
}

type FuzzNameResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List        []FuzzNameCompany `json:"list"`
		Aggregation struct {
			Industry []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"industry"`
			EntStatus []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Count string `json:"count"`
				Sort  string `json:"sort"`
			} `json:"entStatus"`
			Region []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"region"`
		} `json:"aggregation"`
		IsList  bool   `json:"isList"`
		IsVip   bool   `json:"isVip"`
		Keyword string `json:"keyword"`
		Total   int    `json:"total"`
		FuzzyID int    `json:"fuzzyId"`
	} `json:"data"`
	Success bool `json:"success"`
}

type FuzzNameCompany struct {
	RegNo             string        `json:"regNo"`
	EntPic            string        `json:"entPic"`
	Dom               string        `json:"dom"`
	EntName           string        `json:"entName"`
	Entstatus         string        `json:"ENTSTATUS"`
	Entid             string        `json:"entid"`
	RiskTagNames      []string      `json:"risk_tag_names"`
	Faren             string        `json:"faren"`
	FarenID           string        `json:"farenId"`
	Uniscid           string        `json:"UNISCID"`
	Tels              []string      `json:"tels"`
	EntType           string        `json:"entType"`
	Emails            []string      `json:"emails"`
	RegCapCur         string        `json:"regCapCur"`
	HighlightIds      []interface{} `json:"highlightIds"`
	Regno             string        `json:"REGNO"`
	RegCap            string        `json:"regCap"`
	Personid          string        `json:"personid"`
	RegConcat         string        `json:"regConcat"`
	EmailCount        int           `json:"emailCount"`
	EntNameHistory    interface{}   `json:"entNameHistory"`
	RiskTagNames0     []string      `json:"riskTagNames"`
	Website           string        `json:"website"`
	Dom0              string        `json:"DOM"`
	EntnameSearch     string        `json:"ENTNAME_SEARCH"`
	DataType          string        `json:"dataType"`
	EntStatusColor    string        `json:"entStatusColor"`
	Tags              []string      `json:"tags"`
	Uniscid0          string        `json:"uniscid"`
	EsDate            string        `json:"esDate"`
	Esdate            string        `json:"ESDATE"`
	FarenType         int           `json:"farenType"`
	TelCount          int           `json:"telCount"`
	HighlightExtraKey string        `json:"highlightExtraKey"`
	Entname           string        `json:"ENTNAME"`
	EntStatus         string        `json:"entStatus"`
	HighlightNameType string        `json:"highlightNameType"`
	Website0          string        `json:"WEBSITE"`
}

const isVipURL = "https://riskbird.com/riskbird-api/user/checkIsVip"

func (r *RiskbirdClient) CheckLogin() bool {
	resp, err := clients.DoRequest("GET", isVipURL, r.Headers, nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		gologger.DualLog(r.ctx, gologger.Level_DEBUG, fmt.Sprintf("请求 %s 失败: %v", isVipURL, err))
		return false
	}
	return resp.StatusCode() == 200
}

// fuzz完整的企业名称，通过isNeedFuzz判断是否需要模糊查询
func (r *RiskbirdClient) FuzzCompanyName(company string) (*FuzzNameCompany, error) {
	// if !isNeedFuzz {
	// 	return company, "在营", nil
	// }
	data := map[string]string{
		"queryType":           "1",
		"searchKey":           company,
		"pageNo":              "1",
		"range":               "10",
		"selectConditionData": "{\"status\":\"\",\"sort_field\":\"\"}",
	}
	postData, _ := json.Marshal(data)
	resp, err := clients.DoRequest("POST", "https://riskbird.com/riskbird-api/newSearch", r.Headers, bytes.NewReader(postData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, err
	}
	var result FuzzNameResult
	json.Unmarshal(resp.Body(), &result)
	if result.Code != 20000 {
		return nil, fmt.Errorf("code: %d, msg: %s", result.Code, result.Msg)
	}
	if result.Data.Total == 0 {
		return nil, fmt.Errorf("没有找到相关企业")
	}
	return &result.Data.List[0], nil
}

var regOrderNo = regexp.MustCompile(`WEB20\d+`)

func (r *RiskbirdClient) FetchOrderNo(company string) (string, error) {
	searchURL := fmt.Sprintf("https://riskbird.com/ent/%s.html", url.QueryEscape(company))
	resp, err := clients.DoRequest("GET", searchURL, r.WithXSHeader(), nil, 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return "", err
	}
	body := string(resp.Body())
	orderNo := regOrderNo.FindString(body)
	if orderNo == "" {
		return "", fmt.Errorf("没有找到相关企业")
	}
	return orderNo, nil
}

func (r *RiskbirdClient) FetchBasicCompanyInfo(company string) (structs.CompanyInfo, string, error) {
	var companyInfo structs.CompanyInfo

	com, err := r.FuzzCompanyName(company)
	if err != nil {
		return companyInfo, "", fmt.Errorf("%s fuzz name error: %w", company, err)
	}
	orderNo, err := r.FetchOrderNo(com.EntName)
	if err != nil {
		return companyInfo, "", fmt.Errorf("%s fetch order no error: %w", com.EntName, err)
	}
	companyInfo.CompanyName = com.EntName
	companyInfo.RegStatus = com.EntStatus
	companyInfo.Amount = com.RegCap
	companyInfo.Investment = com.EsDate

	return companyInfo, orderNo, nil
}

var apiURL = "https://riskbird.com/riskbird-api/companyInfo/list"

type SubsidiaryInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		FilterData  interface{}         `json:"filterData"`
		TotalCount  int                 `json:"totalCount"`
		Aggregation interface{}         `json:"aggregation"`
		APIData     []SubsidiaryApiData `json:"apiData"`
	} `json:"data"`
	Success bool `json:"success"`
}

type SubsidiaryApiData struct {
	EntPic             string `json:"entPic"`
	EntName            string `json:"entName"`
	Entid              string `json:"entid"`
	FunderRatio        string `json:"funderRatio"`
	RegCapFormat       string `json:"regCapFormat"`
	PersonPic          string `json:"personPic"`
	EntStatusRiskLevel string `json:"entStatusRiskLevel"`
	PersonName         string `json:"personName"`
	SubConAm           string `json:"subConAm"`
	EsDate             string `json:"esDate"`
	RegCap             string `json:"regCap"`
	PersonID           string `json:"personId"`
	EntStatus          string `json:"entStatus"`
	PersonEntCount     int    `json:"personEntCount"`
}

func (r *RiskbirdClient) FetchSubsidiary(orderNo string) ([]SubsidiaryApiData, error) {
	data := map[string]string{
		"page":        "1",
		"size":        "100",
		"orderNo":     orderNo,
		"extractType": "companyInvest",
		"sortField":   "",
	}
	postData, _ := json.Marshal(data)
	resp, err := clients.DoRequest("POST", apiURL, r.WithXSHeader(), bytes.NewReader(postData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("%s 请求失败: %v", apiURL, err)
	}
	var subsidiaryInfo SubsidiaryInfo
	json.Unmarshal(resp.Body(), &subsidiaryInfo)
	if subsidiaryInfo.Code != 20000 {
		return nil, fmt.Errorf("code: %d, msg: %s", subsidiaryInfo.Code, subsidiaryInfo.Msg)
	}
	if subsidiaryInfo.Data.TotalCount == 0 {
		return nil, fmt.Errorf("没有找到相关子公司")
	}

	return subsidiaryInfo.Data.APIData, nil
}

type AppletInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		FilterData  interface{}      `json:"filterData"`
		TotalCount  int              `json:"totalCount"`
		Aggregation interface{}      `json:"aggregation"`
		APIData     []structs.Applet `json:"apiData"`
	} `json:"data"`
	Success bool `json:"success"`
}

// 小程序
func (r *RiskbirdClient) FetchApplet(orderNo string) ([]structs.Applet, error) {
	data := map[string]string{
		"page":        "1",
		"size":        "20",
		"orderNo":     orderNo,
		"extractType": "propertyMiniprogram",
		"sortField":   "",
	}
	postData, _ := json.Marshal(data)
	resp, err := clients.DoRequest("POST", apiURL, r.WithXSHeader(), bytes.NewReader(postData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("%s 请求失败: %v", apiURL, err)
	}
	var appletInfo AppletInfo
	json.Unmarshal(resp.Body(), &appletInfo)
	if appletInfo.Code != 20000 {
		return nil, fmt.Errorf("code: %d, msg: %s", appletInfo.Code, appletInfo.Msg)
	}
	if appletInfo.Data.TotalCount == 0 {
		return nil, fmt.Errorf("没有找到相关小程序")
	}
	return appletInfo.Data.APIData, nil
}

type AppInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		FilterData  interface{}   `json:"filterData"`
		TotalCount  int           `json:"totalCount"`
		Aggregation interface{}   `json:"aggregation"`
		APIData     []structs.App `json:"apiData"`
	} `json:"data"`
	Success bool `json:"success"`
}

// App
func (r *RiskbirdClient) FetchApp(orderNo string) ([]structs.App, error) {
	data := map[string]string{
		"page":        "1",
		"size":        "30",
		"orderNo":     orderNo,
		"extractType": "propertyApp",
		"sortField":   "",
	}
	postData, _ := json.Marshal(data)
	resp, err := clients.DoRequest("POST", apiURL, r.WithXSHeader(), bytes.NewReader(postData), 10, clients.NewRestyClient(nil, true))
	if err != nil {
		return nil, fmt.Errorf("%s 请求失败: %v", apiURL, err)
	}
	var appInfo AppInfo
	json.Unmarshal(resp.Body(), &appInfo)
	if appInfo.Code != 20000 || !appInfo.Success {
		return nil, fmt.Errorf("code: %d, msg: %s", appInfo.Code, appInfo.Msg)
	}
	if appInfo.Data.TotalCount == 0 {
		return nil, fmt.Errorf("没有找到相关App")
	}
	return appInfo.Data.APIData, nil
}
