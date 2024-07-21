package structs

// 返回后端执行状态
// Error: true/false
// Msg:   错误信息
type Status struct {
	Error bool
	Msg   string
}

type Response struct {
	Error  bool
	Proto  string
	Header []map[string]string
	Body   string
}

type Navigation struct {
	Name     string
	Children []Children
}

type Children struct {
	Name string
	Type string
	Path string
}

type TycCompanyInfo struct {
	CompanyID   string
	CompanyName string
}

type Position struct {
	Country   string
	Province  string
	City      string
	District  string
	Connector string // 连接符
}

type SpaceOption struct {
	FofaApi   string
	FofaEmail string
	FofaKey   string
	HunterKey string
	QuakeKey  string
}
