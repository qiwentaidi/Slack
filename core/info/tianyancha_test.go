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
	if companyName != fuzzName {                                     // 如果传进来的名称与模糊匹配的不相同
		var isFuzz = fmt.Sprintf("天眼查模糊匹配名称为%v ——> %v,已替换原有名称进行查.", companyName, fuzzName)
		fmt.Printf("isFuzz: %v\n", isFuzz)
	}
	as := SearchSubsidiary(context.TODO(), fuzzName, companyId, 100)
	fmt.Printf("as: %v\n", as)
}
