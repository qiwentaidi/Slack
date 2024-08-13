package chaos

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	domains := FetchHosts(context.TODO(), "element-plus.org", "")
	fmt.Printf("domains: %v\n", domains)
}
