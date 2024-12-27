package core

import "regexp"

var regexURL = regexp.MustCompile(`https?://[^\s"']+`)

func (t *Tools) ExtractURLs(input string) []string {
	return regexURL.FindAllString(input, -1)
}
