package space

import (
	"slack-wails/lib/structs"
	"testing"
)

func TestQuake(t *testing.T) {
	QuakeApiSearch(&structs.QuakeRequestOptions{
		Query:    ``,
		PageNum:  1,
		PageSize: 2,
		Token:    "",
	})
}
