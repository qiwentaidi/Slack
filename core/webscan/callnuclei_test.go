package webscan

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const text = `
				     __     _
   ____  __  _______/ /__  (_)
  / __ \/ / / / ___/ / _ \/ /
 / / / / /_/ / /__/ /  __/ /
/_/ /_/\__,_/\___/_/\___/_/   v3.2.4

			 projectdiscovery.io

[INF] Current nuclei version: v3.2.4 (outdated)
[INF] Current nuclei-templates version: v9.8.5 (latest)
[WRN] Scan results upload to cloud is disabled.
[INF] New templates added in latest release: 142
[INF] Templates loaded for current scan: 27
[INF] Executing 27 signed templates from projectdiscovery/nuclei-templates
[INF] Targets loaded for current scan: 1
[INF] Templates clustered: 3 (Reduced 2 Requests)
[sap-nw-webgui] [http] [info] http://erpdev.wingtech.com:50000/sap/bc/gui/sap/its/webgui
[CVE-2022-22536] [http] [critical] http://erpdev.wingtech.com:50000/sap/admin/public/default.html [sap_path="/sap/admin/public/default.html"]
[INF] Using Interactsh Server: oast.site
[CVE-2021-42063] [http] [medium] http://erpdev.wingtech.com:50000/SAPIrExtHelp/random/SAPIrExtHelp/random/%22%3e%3c%53%56%47%20%4f%4e%4c%4f%41%44%3d%26%23%39%37%26%23%31%30%38%26%23%31%30%31%26%23%31%31%34%26%23%31%31%36%28%26%23%78%36%34%26%23%78%36%66%26%23%78%36%33%26%23%78%37%35%26%23%78%36%64%26%23%78%36%35%26%23%78%36%65%26%23%78%37%34%26%23%78%32%65%26%23%78%36%34%26%23%78%36%66%26%23%78%36%64%26%23%78%36%31%26%23%78%36%39%26%23%78%36%65%29%3e.asp
[CVE-2020-6287] [http] [critical] http://erpdev.wingtech.com:50000/CTCWebService/CTCWebServiceBean/ConfigServlet
[sap-netweaver-detect] [http] [info] http://erpdev.wingtech.com:50000/startPage ["SAP NetWeaver Application Server 7.21 / AS Java 7.31"]
[sap-netweaver-detect] [http] [info] http://erpdev.wingtech.com:50000/index.jsp ["SAP NetWeaver Application Server 7.21 / AS Java 7.31"]
[sap-netweaver-detect] [http] [info] http://erpdev.wingtech.com:50000 ["SAP NetWeaver Application Server 7.21 / AS Java 7.31"]
`

func TestCaller(t *testing.T) {
	reResult := regexp.MustCompile(`\[(.*?)\]\s+\[(.*?)\]\s+\[(.*?)\]\s+(https?://[^\s]+)\s*\[(.*?)\]`)
	for _, line := range strings.Split(text, "\n") {
		vulinfo := reResult.FindString(line)
		matches := reResult.FindStringSubmatch(vulinfo)
		if len(matches) > 0 {
			fmt.Printf("漏洞名: %s\n", matches[1])
			fmt.Printf("漏洞类型: %s\n", matches[2])
			fmt.Printf("危害等级: %s\n", matches[3])
			fmt.Printf("漏洞地址: %s\n", matches[4])
			fmt.Printf("扩展信息: %s\n", matches[5])
		}
	}
}
