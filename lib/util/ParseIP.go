package util

import (
	"net"
	"regexp"

	"strconv"
	"strings"
)

var ipv4Regex = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

func ParseIPs(ipList []string) (ips []string) {
	var noip, temp []string
	r, _ := regexp.Compile(ipv4Regex)
	for _, line := range ipList {
		if strings.Contains(line, "!") {
			noip = append(noip, ParseIP(line[1:])...)
		} else {
			temp = append(ips, ParseIP(line)...)
		}
	}
	for _, np := range noip {
		temp = RemoveElement(ips, np)
	}
	for _, ip := range temp {
		if r.MatchString(ip) {
			ips = append(ips, ip)
		}
	}
	return ips
}

func ParseTarget(text string) []string {
	var temp, targets []string
	temp = strings.Split(text, "\n") // 通过查找换行去分割每个目标
	for _, t := range temp {
		if t != "" {
			targets = append(targets, t)
		}
	}
	return targets
}

func ParseIP(ipString string) []string {
	var result []string
	if strings.Contains(ipString, "-") {
		result = append(result, parseIP1(ipString)...)
	} else if strings.Contains(ipString, ",") {
		ipArray := strings.Split(ipString, ",")
		for _, ip := range ipArray {
			if strings.Contains(ip, "/") {
				parsedIPs, err := parseCIDR(ip)
				if err != nil {
					continue
				}
				result = append(result, parsedIPs...)
			} else {
				result = append(result, ip)
			}
		}
	} else if strings.Contains(ipString, "/") {
		parsedIPs, err := parseCIDR(ipString)
		if err == nil {
			result = append(result, parsedIPs...)
		}
	} else {
		result = append(result, ipString)
	}
	return result
}

func parseCIDR(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// 解析ip段: 192.168.111.1-255
// 192.168.111.1-192.168.112.255
func parseIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}
	return AllIP
}
