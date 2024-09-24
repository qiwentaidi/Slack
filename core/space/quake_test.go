package space

import (
	"context"
	"fmt"
	"testing"
)

func TestQuake(t *testing.T) {
	ports := GetShodanAllPort(context.Background(), "1.1.1.1")
	fmt.Printf("ports: %v\n", ports)
	// qr := SearchQuakeTips("jeecg")
	// fmt.Printf("qr: %v\n", qr)
}
