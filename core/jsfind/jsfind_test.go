package jsfind

import (
	"context"
	"testing"
)

func TestFindInfo(t *testing.T) {
	authentication := []string{"token不能为空", "令牌不能为空", "令牌已过期", "Unauthorized", "Access Denied", "认证失败"}
	AnalyzeAPI(context.Background(), "", "", []string{""}, nil, authentication)
}
