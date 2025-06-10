package fileutil

import (
	"context"
	"encoding/json"
	"os"
	"slack-wails/lib/gologger"
)

func SaveJsonWithFormat(ctx context.Context, filepath string, content interface{}) bool {
	b, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		gologger.Error(ctx, err)
		return false
	}
	if err := os.WriteFile(filepath, b, 0644); err != nil {
		gologger.Error(ctx, err)
		return false
	}
	return true
}
