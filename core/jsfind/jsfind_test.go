package jsfind

import (
	"context"
	"testing"
)

func TestFindInfo(t *testing.T) {
	AnalyzeAPI(context.Background(), "", "", []string{""}, nil)
}
