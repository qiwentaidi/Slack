package space

import (
	"context"
	"fmt"
	"slack-wails/lib/structs"
	"slack-wails/lib/util"
	"strings"
)

type Result struct {
	URL        string
	IP         string
	Domain     string
	Port       string
	Protocol   string
	Title      string
	Components string
	Source     string
}

func Uncover(ctx context.Context, query, types string, size int, o structs.SpaceOption) []Result {
	var results []Result
	if o.FofaApi != "" && o.FofaEmail != "" && o.FofaKey != "" {
		config := NewFofaConfig(&FofaAuth{
			Address: o.FofaApi,
			Email:   o.FofaEmail,
			Key:     o.FofaKey,
		})
		fs := config.FofaApiSearch(ctx, FormatQuery("fofa", types, query), "10", "1", false, false)
		for _, r := range fs.Results {
			results = append(results, Result{
				URL:        r.URL,
				IP:         r.IP,
				Domain:     r.Domain,
				Port:       r.Port,
				Protocol:   r.Protocol,
				Title:      r.Title,
				Components: r.Product,
				Source:     "FOFA",
			})
		}
	}

	if o.HunterKey != "" {
		hs := HunterApiSearch(o.HunterKey, FormatQuery("hunter", types, query), fmt.Sprintf("%v", size), "1", "2", "", false)
		for _, r := range hs.Data.Arr {
			var component []string
			for _, c := range r.Component {
				component = append(component, util.MergeNonEmpty([]string{c.Name, c.Version}, "/"))
			}
			results = append(results, Result{
				URL:        r.URL,
				IP:         r.IP,
				Domain:     r.Domain,
				Port:       fmt.Sprintf("%v", r.Port),
				Protocol:   r.Protocol,
				Title:      r.WebTitle,
				Components: strings.Join(component, ","),
				Source:     "Hunter",
			})
		}
	}

	if o.QuakeKey != "" {
		for _, r := range QuakeApiSearch(&QuakeRequestOptions{
			Query:    FormatQuery("quake", types, query),
			PageNum:  1,
			PageSize: size,
			Token:    o.QuakeKey,
		}).Data {
			results = append(results, Result{
				IP:         r.IP,
				Domain:     r.Host,
				Port:       fmt.Sprintf("%v", r.Port),
				Protocol:   r.Protocol,
				Title:      r.Title,
				Components: strings.Join(r.Components, ","),
				Source:     "Quake",
			})
		}
	}
	uniqueResults := make([]Result, 0)
	seen := make(map[string]bool)

	for _, r := range results {
		key := fmt.Sprintf("%s:%s:%s", r.IP, r.Port, r.Title)
		if !seen[key] {
			seen[key] = true
			uniqueResults = append(uniqueResults, r)
		}
	}

	return uniqueResults
}

// IP 域名 Body 都可以自动识别
func formatQueryString(prefix string, queries []string, operator string) string {
	var temp []string
	for _, v := range queries {
		temp = append(temp, fmt.Sprintf("%s\"%s\"", prefix, v))
	}
	return strings.Join(temp, operator)
}

func FormatQuery(space, types, query string) string {
	qs := strings.Split(query, ",")
	operator := " || "
	var prefix string

	switch space {
	case "fofa":
		switch types {
		case "标题":
			prefix = "title="
		case "备案号":
			prefix = "icp="
		case "备案名称":
			return "不支持"
		case "域名":
			prefix = "domain="
		default:
			prefix = query
		}
	case "hunter":
		switch types {
		case "标题":
			prefix = "title="
		case "备案号":
			prefix = "icp.number="
		case "备案名称":
			prefix = "icp.name="
		default:
			prefix = query
		}
	case "quake":
		operator = " OR "
		switch types {
		case "标题":
			prefix = "title:"
		case "备案号":
			prefix = "icp:"
		case "备案名称":
			prefix = "icp_keywords:"
		case "域名":
			prefix = "domain:"
		default:
			prefix = query
		}
	default:
		return ""
	}

	return formatQueryString(prefix, qs, operator)
}
