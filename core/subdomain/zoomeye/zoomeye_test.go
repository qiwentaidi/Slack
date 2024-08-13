package zoomeye

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	zh := FetchHosts(context.TODO(), "baidu.com", "")
	for _, item := range zh.List {
		fmt.Printf("item.Name: %v\n", item.Name)
	}
}
