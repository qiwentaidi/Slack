package securitytrails

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	sh := FetchHosts(context.TODO(), "element-plus.org", "")
	fmt.Printf("sh.Subdomains: %v\n", sh.Subdomains)
}
