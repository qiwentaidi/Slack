package info

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slack-wails/lib/clients"
	"slack-wails/lib/util"
	"time"

	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
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
		MatchType  string      `json:"matchType"`
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
			ID               int64       `json:"id"`
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
	company_name string
	company_id   string
	// tycTotal     = regexp.MustCompile(`beian-name">(\d+)`)
	// reg          = regexp.MustCompile(`(?s)ranking-keys">.*?<span class="ranking-ym" rel="nofollow">.*?</span>`) // 包含公司名称以及域名
	// regCompany   = regexp.MustCompile(`keys">(.*?)</a>`)                                                         // 公司名
	// regDomain    = regexp.MustCompile(`nofollow">(.*?)</span>`)                                                  // 域名
)

var (
	gethead  = http.Header{}
	posthead = http.Header{}
)

func InitHEAD(token string) {
	gethead.Set("Version", "TYC-Web")
	gethead.Set("X-Auth-Token", token)
	gethead.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	posthead.Set("Version", "TYC-Web")
	posthead.Set("X-Auth-Token", token)
	posthead.Set("Content-Type", "application/json")
	posthead.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
}

// 要根据ID值查子公司
func GetCompanyID(company string) (string, string) {
	var max, id int
	data := make(map[string]interface{})
	data["keyword"] = company
	bytesData, _ := json.Marshal(data)
	_, b, err := clients.NewRequest("POST", "https://capi.tianyancha.com/cloud-tempest/search/suggest/v3", posthead, bytes.NewReader(bytesData), 10, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	var qs TycSearchID
	json.Unmarshal([]byte(string(b)), &qs)
	time.Sleep(time.Second * 2)
	if len(qs.Data) > 0 { // 先走接口不会进行模糊匹配,如果匹配不到值那就走模糊查询
		return qs.Data[0].GraphID, qs.Data[0].ComName
	} else {
		_, b, err := clients.NewRequest("GET", "https://www.tianyancha.com/search?key="+url.QueryEscape(company), gethead, nil, 10, clients.DefaultClient())
		if err != nil {
			logger.NewDefaultLogger().Debug(err.Error())
		}
		fuzzy := regexp.MustCompile(`\d{10}" target="_blank">(.*?)</span></a>`)
		all := fuzzy.FindAllString(string(b), -1)
		for _, v := range all {
			s := strings.Split(v, `" target="_blank"><span>`)
			f := s[1][:len(s[1])-11] // 模糊匹配到的词绍兴市<em>公安</em>局<em>越城</em>区<em>分局</em>
			var temp string
			for _, keyword := range strings.Split(strings.ReplaceAll(f, "/", ""), "<em>") {
				if strings.Contains(company, keyword) {
					id++
				}
				temp += keyword
			}
			if max < id {
				max = id
				company_id = s[0]
				company_name = temp
			}
		}
		return company_id, company_name
	}
}

type CompanyInfo struct {
	CompanyName string
	Holding     string
	Investment  string
	Domains     []string
}

// 返回查询公司的名称和子公司的名称
func SearchSubsidiary(companyName, companyId string, ratio int) (holdasset []CompanyInfo, logs string) {
	data := make(map[string]interface{})
	data["gid"] = companyId
	data["pageSize"] = 100
	data["pageNum"] = 1
	data["province"] = "-100"
	data["percentLevel"] = "-100"
	data["category"] = "-100"
	bytesData, _ := json.Marshal(data)
	_, b, err := clients.NewRequest("POST", "https://capi.tianyancha.com/cloud-company-background/company/investListV2", posthead, bytes.NewReader(bytesData), 10, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
	}
	var qr TycResult
	json.Unmarshal(b, &qr)
	// 获取到本公司对应的域名
	domains := Beianx(companyName)
	holdasset = append(holdasset, CompanyInfo{companyName, "本公司", "", util.RemoveDuplicates(domains)})
	for _, result := range qr.Data.Result {
		gq, _ := strconv.Atoi(strings.TrimSuffix(result.Percent, "%"))
		if gq <= 100 && gq >= ratio { // 提取在控股范围内的子公司
			subsidiaryDomains := Beianx(result.Name)
			if len(subsidiaryDomains) > 0 {
				holdasset = append(holdasset, CompanyInfo{result.Name, result.Percent, result.Amount, util.RemoveDuplicates(subsidiaryDomains)})
			} else { // 没查到域名的子公司也要显示
				holdasset = append(holdasset, CompanyInfo{result.Name, result.Percent, result.Amount, []string{}})
			}
		}
	}
	return holdasset, logs
}

// 返回对应域名组
// func ICP2Domain(company string) (domains []string) {
// 	var pages int
// 	_, b, err := clients.NewRequest("GET", "https://beian.tianyancha.com/search/"+url.QueryEscape(company), gethead, nil, 10, clients.DefaultClient())
// 	if err != nil {
// 		logger.NewDefaultLogger().Debug(err.Error())
// 	}
// 	fmt.Printf("string(b): %v\n", string(b))
// 	t := tycTotal.FindStringSubmatch(string(b))
// 	if len(t) > 0 {
// 		if total, _ := strconv.Atoi(t[1]); total <= 20 {
// 			pages = 1
// 		} else {
// 			pages = total/20 + 1
// 		}
// 	} else {
// 		pages = 1
// 	}
// 	for _, v := range reg.FindAllString(string(b), -1) {
// 		companyName := regCompany.FindStringSubmatch(v)
// 		domain := regDomain.FindStringSubmatch(v)
// 		companyNames := strings.ReplaceAll(strings.ReplaceAll(companyName[1], "<em>", ""), "</em>", "")
// 		if companyNames == company {
// 			domains = append(domains, domain[1])
// 			d := Beianx(company)
// 			if len(d) > 0 {
// 				domains = append(domains, d...)
// 			}
// 		}
// 	}
// 	// 查询页码大于1时需要对其他页码也进行筛选
// 	if pages > 1 {
// 		for i := 2; i <= pages; i++ {
// 			_, b, err := clients.NewRequest("GET", fmt.Sprintf(`https://beian.tianyancha.com/search/%v/p%d`, url.QueryEscape(company), i), gethead, nil, 10, clients.DefaultClient())
// 			if err != nil {
// 				logger.NewDefaultLogger().Debug(err.Error())
// 			}
// 			for _, v := range reg.FindAllString(string(b), -1) {
// 				companyName := regCompany.FindStringSubmatch(v)
// 				domain := regDomain.FindStringSubmatch(v)
// 				companyNames := strings.ReplaceAll(strings.ReplaceAll(companyName[1], "<em>", ""), "</em>", "")
// 				if companyNames == company {
// 					domains = append(domains, domain[1])
// 					d := Beianx(company)
// 					if len(d) > 0 {
// 						domains = append(domains, d...)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return domains
// }

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
func WeChatOfficialAccounts(companyName, companyId string) (asset [][]string, logs string) {
	_, b, err := clients.NewRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize=1&pageNum=1", gethead, nil, 10, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return nil, logs
	}
	var oa OfficialAccounts
	json.Unmarshal(b, &oa)
	if oa.ErrorCode != 0 || oa.Data.Count == 0 {
		logs = "公众号查询出现错误或不存在公众号资产,公司名称: " + companyName
		return nil, logs
	}

	_, b, err = clients.NewRequest("GET", "https://capi.tianyancha.com/cloud-business-state/wechat/list?graphId="+companyId+"&pageSize="+fmt.Sprint(oa.Data.Count)+"&pageNum=1", gethead, nil, 10, clients.DefaultClient())
	if err != nil {
		logger.NewDefaultLogger().Debug(err.Error())
		return nil, logs
	}
	json.Unmarshal(b, &oa)
	for _, result := range oa.Data.ResultList {
		asset = append(asset, []string{result.Title, result.PublicNum, result.TitleImgURL, result.CodeImg, result.Recommend})
	}
	return asset, logs
}
