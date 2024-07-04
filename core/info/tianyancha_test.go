package info

import (
	"context"
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	token := "eyJhbGciOiJIUzUxMiJ9..kDoQD7j_Kax9Fq0YTDus2veP7kS_7Z_7CT5p9rwmTvLyjm7xypHAPwSlEe6IDgh4ziLsFAeXxwNqgNQIGzDw2g"
	InitHEAD(token)
	companyName := "安恒信息"
	companyId, fuzzName := GetCompanyID(context.TODO(), companyName) // 获得到一个模糊匹配后，关联度最高的名称
	ss := SearchSubsidiary(context.TODO(), fuzzName, companyId, 100)
	fmt.Printf("ss: %v\n", ss)
}
