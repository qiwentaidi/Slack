package core

import (
	"fmt"
	"os"
	"testing"
)

func TestXxx(t *testing.T) {
	content, err := os.ReadFile("result.txt")
	if err != nil {
		t.Fatal(err)
	}
	tools := &Tools{}
	nodes := tools.FormatOutput(string(content))
	for _, node := range nodes {
		fmt.Printf("node: %v\n", node)
	}
}
