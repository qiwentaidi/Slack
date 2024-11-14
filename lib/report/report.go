package report

import (
	"fmt"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
)

func GenerateReport(Fingerprints []structs.InfoResult, POCs []structs.VulnerabilityInfo) string {
	return defaultHeader() + reportBody(Fingerprints, POCs)
}

func reportBody(Fingerprints []structs.InfoResult, POCs []structs.VulnerabilityInfo) string {
	var allContent string
	fingerprintsSection := `<div onclick="$(this).next().toggle()" style="background:#1A2733; color:#FFF; padding:8px; cursor:pointer; font-weight:bold;">
		Fingerprints
	</div>
	<div style="background-color:#223B46; color:#DDE2DE; display:none;">`
	for _, fingerprint := range Fingerprints {
		fingerprintsSection += fmt.Sprintf(`
			<div style="padding:8px; border-bottom:1px solid #60786F;">
				<a href="%s" target="_blank" style="color:inherit; text-decoration:inherit;"><span style="color:inherit;">%s</span></a> &nbsp;
				<span style="color:#00FF00;">[%d]</span> &nbsp;
				<span>[%d]</span> &nbsp;
				<span>%s</span> &nbsp;
				<span style="color:#FF4C4C;">%s</span>
				%s
			</div>`,
			fingerprint.URL, fingerprint.URL, fingerprint.StatusCode, fingerprint.Length, fingerprint.Title, strings.Join(fingerprint.Fingerprints, ", "), showWafInfo(fingerprint.IsWAF, fingerprint.WAF))
	}
	fingerprintsSection += "</div>"

	allContent += fingerprintsSection
	for index, poc := range POCs {
		title := fmt.Sprintf(`<table>
		<thead onclick="$(this).next('tbody').toggle()" style="background:#DDE2DE">
			<td class="vuln">%d&nbsp;&nbsp;%s</td>
			<td class="security %s">%s</td>
			<td class="url">%s</td>
		</thead>`, index+1, poc.ID, strings.ToLower(poc.Risk), poc.Risk, util.GetBasicURL(poc.URL))
		info := fmt.Sprintf("<b>name:</b> %s&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<b>security:</b> %s",
			poc.Name, poc.Risk)
		if len(poc.Extract) > 0 {
			info += fmt.Sprintf("<br><b>extract:</b> %s", poc.Extract)
		}
		if len(poc.Description) > 0 {
			info += "<br/><b>description:</b> " + poc.Description
		}
		if len(poc.Reference) > 0 {
			info += "<br/><b>reference:</b> "
			if strings.Contains(poc.Reference, ",") {
				for _, rv := range strings.Split(poc.Reference, ",") {
					info += "<br/>&nbsp;&nbsp;- <a href='" + rv + "' target='_blank'>" + rv + "</a>"
				}
			} else {
				info += "<a href='" + poc.Reference + "' target='_blank'>" + poc.Reference + "</a>"
			}
		}

		header := "<tbody>"

		bodyinfo := fmt.Sprintf(`<tr>
			<td colspan="3">%s</td>
		</tr>`, info)

		body := ""

		// @edit 2023.5.15 10:36 because mysql-detect not html report
		// reqraw := []byte{}
		// respraw := []byte{}
		// if v.ResultResponse.Url != nil {
		reqraw := poc.Request
		respraw := poc.Response
		// }

		fullurl := xssfilter(poc.URL)

		body += fmt.Sprintf(`<tr>
		<td colspan="3"  style="border-top:1px solid #60786F"><a href="%s" target="_blank">%s</a></td>
	</tr><tr>
			<td colspan="3" style="background: #223B46; color: #DDE2DE;">
				<div class="clr">
				<div class="request w50">
				<div class="toggleR" onclick="$(this).parent().next('.response').toggle();if($(this).text()=='→'){$(this).text('←');$(this).css('background','red');$(this).parent().removeClass('w50').addClass('w100')}else{$(this).text('→');$(this).css('background','black');$(this).parent().removeClass('w100').addClass('w50')}">→</div>
				<div class="copy-button" onclick="copyXmpContent($(this));">Copy</div>
<xmp>%s</xmp>
				</div>
				<div class="response w50">
				<div class="toggleL" onclick="$(this).parent().prev('.request').toggle();if($(this).text()=='←'){$(this).text('→');$(this).css('background','red');$(this).parent().removeClass('w50').addClass('w100')}else{$(this).text('←');$(this).css('background','black');$(this).parent().removeClass('w100').addClass('w50')}">←</div>
				<div style="position: absolute;right: 0;"></div>
<xmp>%s</xmp>
				</div>
			</div>
			</td>
		</tr>
	`, fullurl, fullurl, reqraw, respraw)

		footer := "</tbody></table>"
		allContent += title + header + bodyinfo + body + footer
	}
	return allContent
}

func showWafInfo(isWaf bool, waf string) string {
	if isWaf {
		return fmt.Sprintf("<span style=\"#color: DCA550\">%s</span>", waf)
	} else {
		return ""
	}
}
