package riskbird

import (
	"context"
	"fmt"
	"testing"
)

func TestFuzzName(t *testing.T) {
	cookies := `app-uuid=WEB-1C830717ADF94841BF21C9A77B86F8A2; app-device=WEB; web-version-code=1.0.29; first-authorization=1747929760509; token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJwYXNzd29yZCI6IjMzZWRkMTRjNjhlZTY5NTMwY2IxMTc4MGU3ZGE5ZGQ3IiwiZXhwIjoxNzQ4MDU0MDY5LCJ1c2VySWQiOjcwMDMzNSwidXVpZCI6IjBhZDFmNjNmLTZhNGUtNDYyYi1iMDI3LWYwODQxMWJlYWRiMSIsInVzZXJuYW1lIjoiMTM4NTg5MDQzODYifQ.7LTBsPFgxdi1wfLaZqkcfmJvG8hG4uLNsD8SpZvLWpI; userinfo=%7B%22userId%22%3A700335%2C%22inviteCode%22%3A%22E104308F0E1E937A%22%2C%22nickName%22%3A%2213858904386%22%2C%22unionid%22%3A%22oTZAV60XdhLj3Aqki5a3ezturBhs%22%2C%22isVip%22%3Atrue%2C%22vipStatus%22%3A%22vip%22%2C%22vipEndTime%22%3A%222030-05-20%22%2C%22mobile%22%3A%2213858904386%22%2C%22email%22%3Anull%2C%22timestamp%22%3A1748052233544%2C%22userNewType%22%3Atrue%2C%22vipTimeOut%22%3A1823%2C%22notGetLoginVip%22%3Afalse%2C%22vipExpireTime%22%3A1905523199000%2C%22isQueryRiskDoc%22%3Afalse%2C%22queryRiskDocSwitch%22%3A%221%22%2C%22status%22%3A%22vip%22%7D`
	client := NewClient(context.Background(), cookies)
	isLogin := client.CheckLogin()
	if !isLogin {
		fmt.Println("Login failed")
		return
	}
	// completeName, _ := client.FuzzCompanyName("建设银行", true)
	// fmt.Println("completeName: ", completeName)
	// orderNo, _ := client.FetchOrderNo(completeName)
	// fmt.Println("Order No:", orderNo)
	// client.FetchSubsidiary(orderNo)
	// time.Sleep(2 * time.Second)
	// client.FetchApplet(orderNo)
	// time.Sleep(2 * time.Second)
	// client.FetchApp(orderNo)
}
