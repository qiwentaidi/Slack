package result

import (
	"fmt"
	"slack/common"
	"slack/lib/poc"
	"slack/lib/proto"
	"slack/lib/util"
	"strconv"
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

func (pr *PocResult) ReadFullResultRequestInfo() string {
	result := "\r\n" + pr.ResultRequest.Url.GetScheme() + "://" + pr.ResultRequest.Url.GetHost() + pr.ResultRequest.Url.GetPath()
	if len(pr.ResultRequest.Url.GetQuery()) > 0 {
		result += "?" + pr.ResultRequest.Url.GetQuery()
	} else if len(pr.ResultRequest.Url.Fragment) > 0 {
		result += "#" + pr.ResultRequest.Url.Fragment
	}
	result += "\r\n"

	for k, v := range pr.ResultRequest.Headers {
		result += k + ":" + v + "\r\n"
	}
	result += "\r\n\r\n" + string(pr.ResultRequest.GetBody())
	return result
}

func (pr *PocResult) ReadFullResultResponseInfo() string {
	return string(pr.ResultResponse.GetRaw())
}

func (r *Result) ReadPocInfo() string {
	result := "VulID: " + r.PocInfo.Id + "\r\n"
	result += "Name: " + r.PocInfo.Info.Name + "\r\n"
	result += "Author: " + r.PocInfo.Info.Author + "\r\n"
	result += "Severity: " + r.PocInfo.Info.Severity + "\r\n"
	if len(r.PocInfo.Info.Description) > 0 {
		result += "Description: " + r.PocInfo.Info.Description + "\r\n"
	}
	if len(r.PocInfo.Info.Reference) > 0 {
		result += "Reference: \r\n"
		for _, v := range r.PocInfo.Info.Reference {
			result += "    - " + v + "\r\n"
		}
	}
	if len(r.PocInfo.Info.Tags) > 0 {
		result += "Tags: " + r.PocInfo.Info.Tags + "\r\n"
	}
	if len(r.PocInfo.Info.Classification.CveId) > 0 {
		result += "Classification: \r\n"
		result += "    CveId: " + r.PocInfo.Info.Classification.CveId + "\r\n"
		result += "    CvssMetrics: " + r.PocInfo.Info.Classification.CvssMetrics + "\r\n"
		result += "    CweId: " + r.PocInfo.Info.Classification.CweId + "\r\n"
		result += "    CvssScore: " + strconv.FormatFloat(r.PocInfo.Info.Classification.CvssScore, 'f', 1, 64) + "\r\n"
	}
	return result
}

func (r *Result) WriteOutput() {
	util.BufferWriteAppend(r.Output, "["+util.GetNowDateTime()+"] ["+r.PocInfo.Id+"] ["+r.PocInfo.Info.Severity+"] "+r.Target) // output save to file
}

func (r *Result) PrintResultInfo() string {
	return "[" + util.GetNowDateTime() + "] [" + r.PocInfo.Id + "] [" + r.PocInfo.Info.Severity + "] " + r.Target
}

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

	fmt.Printf("\r%v \033[32m%v\033[37m %v %v\r\n", number+" "+util.GetNowDateTime(), r.PocInfo.Id+" "+strings.ToUpper(r.PocInfo.Info.Severity), r.FullTarget, extinfo)
	if strings.ToUpper(r.PocInfo.Info.Severity) == "INFO" { // 如果poc执行到等级为INFO视为指纹，需要根据指纹再去执行POC,指纹名截取-左侧的第一个字符串
		common.UrlFingerMap[r.Target] = append(common.UrlFingerMap[r.Target], strings.Split(r.PocInfo.Id, "-")[0])
	}
	common.Vulnerability = append(common.Vulnerability, common.VulnerabilityInfo{Id: fmt.Sprintf("%v", number), Name: r.PocInfo.Id, RiskLevel: strings.ToUpper(r.PocInfo.Info.Severity),
		Url: r.FullTarget, TransInfo: common.TransInfo{Request: string(r.AllPocResult[0].ResultRequest.Raw), Response: r.AllPocResult[0].ReadFullResultResponseInfo()}})
}

func (r *Result) Reset() {
	r.IsVul = false
	r.Target = ""
	*r.PocInfo = poc.Poc{}
	r.AllPocResult = nil
	r.Output = ""
}

func (pr *PocResult) Reset() {
	pr.IsVul = false
	pr.ResultRequest = &proto.Request{}
	pr.ResultResponse = &proto.Response{}
}
