package core

import "regexp"

var (
	regexAlibabaDruidWebSession = regexp.MustCompile(`"SESSIONID":"(?P<session>[^"]+)"`)
	regexAlibabaDruidWebURI     = regexp.MustCompile(`"URI":"(?P<uri>[^"]+)"`)
)

func (t *Tools) ExtractAlibabaDruidWebSession(input string) []string {
	matches := regexAlibabaDruidWebSession.FindAllStringSubmatch(input, -1)
	results := []string{}
	for _, match := range matches {
		results = append(results, match[regexAlibabaDruidWebSession.SubexpIndex("session")])
	}
	return results
}

func (t *Tools) ExtractAlibabaDruidWebURI(input string) []string {
	matches := regexAlibabaDruidWebURI.FindAllStringSubmatch(input, -1)
	results := []string{}
	for _, match := range matches {
		results = append(results, match[regexAlibabaDruidWebURI.SubexpIndex("uri")])
	}
	return results
}
