package info

import (
	"fmt"
	"testing"
)

func TestBeianx(t *testing.T) {
	d := Beianx("深圳市鲸瀚数字文化有限公司")
	fmt.Printf("d: %v\n", d)
}
