package report

import (
	"slack-wails/core/webscan/poc"
	"slack-wails/core/webscan/proto"

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
