package portscan

import (
	"context"
	"fmt"
	"slack-wails/lib/gologger"
	"slack-wails/lib/structs"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func LdapScan(ctx context.Context, host string, usernames, passwords []string) {
	for _, user := range usernames {
		for _, pass := range passwords {
			if ExitFunc {
				return
			}
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(host, user, pass)
			if flag && err == nil {
				runtime.EventsEmit(ctx, "nucleiResult", structs.VulnerabilityInfo{
					ID:       "ldap weak password",
					Name:     "ldap weak password",
					URL:      host,
					Type:     "LDAP",
					Severity: "HIGH",
					Extract:  user + "/" + pass,
				})
				return
			} else {
				gologger.Info(ctx, fmt.Sprintf("ldap://%s %s:%s is login failed", host, user, pass))
			}
		}
	}
}

func Ldapconn(host, user, pass string) (bool, error) {
	conn, err := ldap.Dial("tcp", host)

	if err != nil {
		return false, err
	}

	err = conn.Bind(user, pass)
	if err != nil {
		return false, err
	}

	return true, nil
}
