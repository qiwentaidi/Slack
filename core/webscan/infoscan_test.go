package webscan

import (
	"fmt"
	"slack-wails/lib/clients"
	"testing"
)

func TestInfoscan(t *testing.T) {
	// InitFingprintDB(util.HomeDir() + "/slack/config/webfinger.yaml")
	// title := ""
	// server := ""
	// ct := ""
	// banner := ""
	// resp, body, _ := clients.NewRequest("GET", "https://oa.shanghai-cjsw.com/", nil, nil, 10, clients.DefaultClient())
	// if match := util.RegTitle.FindSubmatch(body); len(match) > 1 {
	// 	title = util.Str2UTF8(string(match[1]))
	// }
	// for k, v := range resp.Header {
	// 	if k == "Server" {
	// 		server = strings.Join(v, ";")
	// 	}
	// 	if k == "Content-Type" {
	// 		ct = strings.Join(v, ";")
	// 	}
	// }
	// headers, _, _ := DumpResponseHeadersAndRaw(resp)

	// result := FingerScan(&TargetINFO{
	// 	HeadeString:   string(headers),
	// 	ContentType:   ct,
	// 	Cert:          GetTLSString("https", "oa.shanghai-cjsw.com:443"),
	// 	BodyString:    string(body),
	// 	Path:          "/",
	// 	Title:         title,
	// 	Server:        server,
	// 	ContentLength: len(body),
	// 	Port:          443,
	// 	IconHash:      FaviconHash("https://oa.shanghai-cjsw.com/", clients.DefaultClient()),
	// 	StatusCode:    resp.StatusCode,
	// 	Banner:        banner,
	// }, FingerprintDB)
	// fmt.Printf("result: %v\n", result)
	// InitWorkflow(util.HomeDir() + "/slack/config/workflow.yaml")
	// pe := TargetBindFingerPocs("https://sg.smt528.com/", []string{"layer.js", "Bootstrap", "Apache-Tomcat",
	// 	"Jeecg-Boot",
	// 	"JSP",
	// 	"Oracle-JAVA",
	// 	"jQuery"})
	// cmd := exec.Command("nuclei", "-duc", "-u", pe.URL, "-t", strings.Join(pe.PocFiles, ","))
	// out, _ := cmd.CombinedOutput()
	// fmt.Printf("string(out): %v\n", string(out))
	url := "https://oa.shanghai-cjsw.com/"
	hash := FaviconHash("https", url, clients.DefaultClient())
	fmt.Printf("hash: %v\n", hash)
}
