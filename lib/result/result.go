package result

import (
	"fmt"
	"slack/common"
	"slack/lib/poc"
	"slack/lib/proto"
	"strings"

	"gopkg.in/yaml.v2"
)

type Result struct {
	IsVul        bool
	Target       string
	FullTarget   string
	PocInfo      *poc.Poc
	AllPocResult []*PocResult
	Output       string
	FingerResult any
	Extractor    yaml.MapSlice
}

type PocResult struct {
	FullTarget     string
	ResultRequest  *proto.Request
	ResultResponse *proto.Response
	IsVul          bool
}

func (pr *PocResult) ReadFullResultResponseInfo() string {
	return string(pr.ResultResponse.GetRaw())
}

// 将漏洞信息输出到表格中
func (r *Result) PrintColorResultInfoConsole(number string) {
	extinfo := "" // 输出拓展信息
	if len(r.Extractor) > 0 {
		for _, v := range r.Extractor {
			switch value := v.Value.(type) {
			case map[string]string:
			case string:
				extinfo += "," + v.Key.(string) + "=\"" + fmt.Sprintf("%v", value) + "\""
			}
		}
		extinfo = "[" + strings.TrimLeft(extinfo, ",") + "]"
	}
	if strings.ToUpper(r.PocInfo.Info.Severity) == "INFO" { // 如果poc执行到等级为INFO视为指纹，需要根据指纹再去执行POC,指纹名截取-左侧的第一个字符串
		common.UrlFingerMap[r.Target] = append(common.UrlFingerMap[r.Target], strings.Split(r.PocInfo.Id, "-")[0])
	}
	common.Vulnerability = append(common.Vulnerability, common.VulnerabilityInfo{
		Id: number, Name: r.PocInfo.Id,
		RiskLevel: strings.ToUpper(r.PocInfo.Info.Severity),
		Url:       r.FullTarget,
		TransInfo: common.TransInfo{
			Request:  string(r.AllPocResult[0].ResultRequest.Raw),
			Response: r.AllPocResult[0].ReadFullResultResponseInfo(),
			ExtInfo:  extinfo,
		}})
}
