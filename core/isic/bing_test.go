package isic

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	// u, _ := chooseBingEnvironment()
	// fmt.Printf("u: %v\n", u)
	result, _ := BingSearch("www.baidu.com")
	fmt.Printf("result1: %v\n", result)
}
