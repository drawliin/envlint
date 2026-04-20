package parser

// File keeps the parsed env values plus a few details that are useful for reporting.
type File struct {
	Path       string            `json:"path"`
	Exists     bool              `json:"exists"`
	Values     map[string]string `json:"values"`
	Duplicates []string          `json:"duplicates"`
	Keys       []string          `json:"keys"`
}
