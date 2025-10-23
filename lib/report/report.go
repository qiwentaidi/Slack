package report

import (
	"fmt"
	"slack-wails/lib/structs"
	"strings"
)

func GenerateReport(Fingerprints []structs.InfoResult, POCs []structs.VulnerabilityInfo) string {
	return defaultHeader() + reportBody(Fingerprints, POCs) + defaultFooter()
}

func reportBody(Fingerprints []structs.InfoResult, POCs []structs.VulnerabilityInfo) string {
	var sb strings.Builder

	// Fingerprints
	for i, fp := range Fingerprints {
		statusClass := "ok"
		if fp.StatusCode >= 400 && fp.StatusCode < 500 {
			statusClass = "warn"
		} else if fp.StatusCode >= 500 {
			statusClass = "bad"
		}

		// WAF
		wafBadge := ""
		if fp.IsWAF && fp.WAF != "" {
			wafBadge = fmt.Sprintf(`<span class="badge warn">%s</span>`, fp.WAF)
		}

		// 指纹信息
		fpTags := ""
		if len(fp.Fingerprints) > 0 {
			fpTags += `<div class="fp-tags">`
			for _, f := range fp.Fingerprints {
				fpTags += fmt.Sprintf(`<span class="fp-tag">%s</span>`, f)
			}
			fpTags += `</div>`
		}

		sb.WriteString(fmt.Sprintf(`
<div class="fp-row">
  <div class="fp-id">%d</div>
  <a href="%s" target="_blank">%s</a>
  <span class="badge %s">%d</span>
  <span>%d</span>
  <span>%s</span>
  %s
  %s
</div>`,
			i+1, fp.URL, fp.URL, statusClass, fp.StatusCode, fp.Length, fp.Title, wafBadge, fpTags))
	}

	sb.WriteString(`</div><div id="vuln" class="tab-content">`)

	// Vulnerabilities
	// Vulnerabilities
	for i, poc := range POCs {
		severityClass := strings.ToLower(poc.Severity)
		fullurl := xssfilter(poc.URL) // 安全 URL

		sb.WriteString(fmt.Sprintf(`
<table>
  <thead class="vuln-head">
    <tr>
      <td class="vuln">%d %s</td>
      <td class="severity"><span class="security %s">%s</span></td>
      <td class="url"><a href="%s" target="_blank">%s</a></td>
    </tr>
  </thead>
  <tbody style="display:none;">
    <tr>
      <td colspan="3">
        <b>name:</b> %s &nbsp;&nbsp; <b>severity:</b> %s`,
			i+1, poc.ID, severityClass, poc.Severity, fullurl, fullurl, poc.Name, poc.Severity))

		if poc.Extract != "" {
			sb.WriteString("<br/><b>extract:</b> " + poc.Extract)
		}
		if poc.Description != "" {
			sb.WriteString("<br/><b>description:</b> " + poc.Description)
		}
		if poc.Reference != "" {
			refs := strings.Split(poc.Reference, ",")
			sb.WriteString("<br/><b>reference:</b>")
			for _, r := range refs {
				sb.WriteString(fmt.Sprintf(` - <a href="%s" target="_blank">%s</a><br/>`, r, r))
			}
		}

		// 请求/响应显示
		sb.WriteString(fmt.Sprintf(`
      </td>
    </tr>
    <tr>
      <td colspan="3" style="background:var(--panel-alt);">
        <div class="clr">
          <div class="request w50">
            <div class="reqresp-header">
              <div class="actions">
                <div class="toggle" onclick="$(this).closest('.request').toggleClass('w100').toggleClass('w50').siblings('.response').toggle();">Request</div>
              </div>
              <div class="time btn-mini" onclick="copyXmpContent($(this));">Copy</div>
            </div>
            <xmp>%s</xmp>
          </div>
          <div class="response w50">
            <div class="reqresp-header">
              <div class="actions">
                <div class="toggle" onclick="$(this).closest('.response').toggleClass('w100').toggleClass('w50').siblings('.request').toggle();">Response</div>
              </div>
              <div class="time">%ss</div>
            </div>
            <xmp>%s</xmp>
          </div>
        </div>
      </td>
    </tr>
  </tbody>
</table>`, poc.Request, poc.ResponseTime, poc.Response))
	}

	return sb.String()
}
