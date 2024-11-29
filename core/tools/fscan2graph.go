package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	// 调整后的端口服务映射
	portServiceMapAdjusted = map[string][]string{
		"SSH":               {":22 open"},
		"SMB":               {":445 open", ":139 open"},
		"RDP":               {":3389 open"},
		"Database":          {":1433 open", ":1521 open", ":3306 open", ":5432 open", ":27017 open", ":6379 open"},
		"K8S":               {":2375 open", ":2376 open", ":2379 open", "Kubernetes"},
		"Domain Controller": {"[+] DC"},
		"Code Repositories": {"Gitlab", "Nexus", "Harbor", "ATLASSIAN-Confluence", "jira"},
		"Vcenter && Exsi":   {"ID_VC_Welcome", "ID_EESX_Welcome"},
	}
)

// NetworkNode 表示网络图中的节点
type NetworkNode struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	IsDouble bool   `json:"isDouble,omitempty"` // 是否为多网卡
}

// NetworkLink 表示节点之间的连接
type NetworkLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Subnet string `json:"subnet"`
}

// NetworkCategory 表示一个类别
type NetworkCategory struct {
	Name string `json:"name"`
}

// FscanGraphData 表示整个图数据
type FscanGraphData struct {
	Nodes      []NetworkNode     `json:"nodes"`
	Links      []NetworkLink     `json:"links"`
	Categories []NetworkCategory `json:"categories"`
}

// 解析文件并生成关系数据
func (t *Tools) FscanToGraph(filePath string) (*FscanGraphData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	graph := &FscanGraphData{}
	nodeSet := make(map[string]bool)
	linkSet := make(map[string]bool)
	categorySet := make(map[string]bool)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析网卡信息
	nets := NetInfoReg.FindAllString(string(content), -1)
	netInfos := FormatNetInfo(nets)

	// 标记多网卡主机
	doubleNICHosts := make(map[string]NetInfo)
	for _, net := range netInfos {
		if len(net.IPs) > 1 {
			doubleNICHosts[net.Hostname] = net
		}
	}

	// 添加默认分类
	defaultCategories := []string{"Host", "Service", "Subnet"}
	for _, cat := range defaultCategories {
		graph.Categories = append(graph.Categories, NetworkCategory{Name: cat})
		categorySet[cat] = true
	}

	// 扫描文件逐行处理
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		for service, patterns := range portServiceMapAdjusted {
			for _, pattern := range patterns {
				if matched, _ := regexp.MatchString(pattern, line); matched {
					// 提取 IP 和服务信息
					re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
					ipMatches := re.FindStringSubmatch(line)
					if len(ipMatches) > 0 {
						ip := ipMatches[1]

						// 判断是否为多网卡
						isDouble := false
						if _, exists := doubleNICHosts[ip]; exists {
							isDouble = true
						}

						// 添加主机节点
						if !nodeSet[ip] {
							graph.Nodes = append(graph.Nodes, NetworkNode{
								Name:     ip,
								Category: "Host",
								IsDouble: isDouble,
							})
							nodeSet[ip] = true
						}

						// 添加服务节点
						if !nodeSet[service] {
							graph.Nodes = append(graph.Nodes, NetworkNode{
								Name:     service,
								Category: "Service",
							})
							nodeSet[service] = true
						}

						// 添加边
						linkKey := fmt.Sprintf("%s->%s", ip, service)
						if !linkSet[linkKey] {
							graph.Links = append(graph.Links, NetworkLink{
								Source: ip,
								Target: service,
								Subnet: "", // 如果需要子网信息，可通过额外逻辑填充
							})
							linkSet[linkKey] = true
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return graph, nil
}
