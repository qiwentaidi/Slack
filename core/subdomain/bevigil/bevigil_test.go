package bevigil

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	bh := FetchHosts(context.TODO(), "baidu.com", "")
	fmt.Printf("bh.Subdomains: %v\n", bh.Subdomains)
}
