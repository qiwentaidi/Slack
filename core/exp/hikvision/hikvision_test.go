package hikvision

import (
	"context"
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	result := CameraHandlessLogin(context.Background(), "http://xxxxxxxxx/", "admin", []string{"hik12345"})
	fmt.Printf("result: %v", result)
}
