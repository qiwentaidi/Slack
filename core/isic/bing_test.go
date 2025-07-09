package isic

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	result, total, _ := GoogleHackerBingSearch("")
	fmt.Printf("total: %v\n", total)
	for _, item := range result {
		fmt.Printf("item: %v\n", item)
	}
}
