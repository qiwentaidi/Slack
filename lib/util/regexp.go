package util

import "regexp"

var (
	RegIP             = regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	RegDomain         = regexp.MustCompile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?`)
	RegCompliance     = regexp.MustCompile(`(\w+)[!,=]{1,3}"([^"]+)"`)                                              // 匹配空间引擎输入内容是否合规
	RegIPCompleteMask = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`) // 192.168.0.0/255.255.0.0
	RegIPCIDR         = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})/(\d{1,2})`)                        // 192.168.0.0/24
)
