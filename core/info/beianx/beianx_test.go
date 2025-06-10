package beianx

import (
	"fmt"
	"testing"
)

func TestBeianx(t *testing.T) {
	domains, err := FetchDomainByCompanyName("苏州大学", "d1794d1e-e24c-4d99-xxxx-177bd4d26b3c")
	fmt.Printf("err: %v\n", err)
	fmt.Printf("domains: %v\n", domains)
}
