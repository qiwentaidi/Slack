package space

import (
	"fmt"
	"testing"
)

func TestQuake(t *testing.T) {
	qk, s := QuakeApiSearch(`ip:"183.131.92.179" AND port: "9088"`, 1, 1, "b13a9af7-f6b5-43cc-bff1-c807ea92dd72", true)
	fmt.Printf("qk: %v\n", qk)
	fmt.Printf("s: %v\n", s)
}
