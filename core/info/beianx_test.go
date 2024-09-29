package info

import (
	"fmt"
	"testing"
)

func TestBeianx(t *testing.T) {
	domains, err := Beianx("苏州大学")
	fmt.Printf("err: %v\n", err)
	fmt.Printf("domains: %v\n", domains)
}
