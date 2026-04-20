package scanner

import "regexp"

// pattern stores one regex and the capture group that contains the env key.
type pattern struct {
	expression *regexp.Regexp
	group      int
}

// Result is the raw output of the source-code scan before audit comparisons happen.
type Result struct {
	Referenced map[string][]string `json:"referenced"`
	Scanned    []string            `json:"scanned_files"`
}
