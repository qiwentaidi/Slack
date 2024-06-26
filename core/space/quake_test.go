package space

import (
	"fmt"
	"testing"
)

func TestQuake(t *testing.T) {
	qk := QuakeApiSearch(&QuakeRequestOptions{
		`body="JHSoft.Web.AddMenu"|| app="Jinher-OA"`, []string{}, 1, 1, true, false, false, false, "",
	})
	fmt.Printf("qk: %v\n", qk)
	// qr := SearchQuakeTips("jeecg")
	// fmt.Printf("qr: %v\n", qr)
}
