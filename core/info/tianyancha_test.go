package info

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	token := ""
	InitHEAD(token)
	companyName := "杭州福斯达深冷装备股份有限公司"
	companyId, fuzzName := GetCompanyID(companyName) // 获得到一个模糊匹配后，关联度最高的名称
	fmt.Printf("fuzzName: %v\n", fuzzName)
	fmt.Printf("companyId: %v\n", companyId)
}
