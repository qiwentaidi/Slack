package jsfind

import (
	"context"
	"fmt"
	"testing"
)

func TestJSFInd(t *testing.T) {
	// result := detectContentType("http://api", nil)
	// fmt.Println(result)
	// str := `{"timestamp":"2025-05-20T02:40:05.456+0000","status":400,"error":"Bad Request","message":"Required String parameter 'userId' is not present","path":"/apphub-api/sprole/list/"}`
	// result2 := extractMissingParams(str)
	// fmt.Printf("result2: %v\n", result2)
	result := extractDynamicsJs(context.Background(), "http://195.195.32.97:59010/xs-addrmatch/#/login")
	fmt.Printf("result1: %v\n", result)
}
