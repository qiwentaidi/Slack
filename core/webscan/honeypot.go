package webscan

import "strings"

var honeypotHeaders = []string{"Cacti", "grafana_session", "X-Jenkins", "Mime-Version", "Composed-By", "zbx_session", "akaunting_session", "DSSIGNIN", "X-Drupal", "drupal", "X-Influxdb", "X-Cmd-Response", "X-Root", "couchdb"}

func checkHoneypotWithHeaders(rawHeaders string) bool {
	var count int
	for _, header := range honeypotHeaders {
		if strings.Contains(rawHeaders, header) {
			count++
		}
		if count >= 3 {
			return true
		}
	}
	return false
}

func checkHoneypotWithFingerprintLength(length int) bool {
	return length > 15
}
