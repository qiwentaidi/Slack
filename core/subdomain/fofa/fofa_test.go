package fofa

import (
	"context"
	"fmt"
	"slack-wails/lib/structs"
	"testing"
)

func TestXxx(t *testing.T) {
	domains := FetchHosts(context.Background(), "baidu.com", structs.FofaAuth{
		Address: "https://fofa.info/",
		Key:     "",
		Email:   "",
	})
	fmt.Printf("domains: %v\n", domains)
}
