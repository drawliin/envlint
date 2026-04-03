package scanner

import "regexp"

type pattern struct {
	expression *regexp.Regexp
	group      int
}

type Result struct {
	Referenced map[string][]string `json:"referenced"`
	Scanned    []string            `json:"scanned_files"`
}
