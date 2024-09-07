package crack

import (
	"fmt"
	"slack-wails/lib/clients"
	"testing"
)

func TestDdddOcr(t *testing.T) {
	_, body, _ := clients.NewSimpleGetRequest("https://faceself.myar.cc/captcha/default?rbRBMM1b", clients.DefaultClient())
	res, _ := DdddOcr(body)
	fmt.Printf("res4: %v\n", res)
}
