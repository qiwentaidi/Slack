package jsfind

import (
	"fmt"
	"testing"
)

func TestJSFInd(t *testing.T) {
	result := detectContentType("http://api", nil)
	fmt.Println(result)
}
