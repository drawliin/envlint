package parser

type File struct {
	Path       string            `json:"path"`
	Exists     bool              `json:"exists"`
	Values     map[string]string `json:"values"`
	Duplicates []string          `json:"duplicates"`
	Keys       []string          `json:"keys"`
}
